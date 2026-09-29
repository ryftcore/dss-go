package xmldom

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// Serialize is a byte-for-byte reimplementation of what DomUtils.serializeNode(Node)
// writes: an OpenJDK identity Transformer (METHOD=xml, no other output property except
// the ENCODING copied off the document) fed a DOMSource wrapping n.
//
// Byte parity is a REQUIREMENT, not a nicety. DSSXMLUtils.applyTransforms(Node, List) -
// dss-xades, the no-transform branch - returns DomUtils.getNodeBytes(node), and those
// bytes are what the ds:Reference DigestValue is computed over. A reference with no
// ds:Transforms therefore digests this serializer's output directly; anything this
// function emits differently from Java is a signature that verifies in one stack and not
// the other. (Serialize is still NOT canonicalization: it keeps prefixes and QNames
// verbatim, does not sort attributes into c14n order, and does not inherit an ancestor's
// namespace axis the way c14n does. Anything covered by a ds:Transform goes through
// internal/xmlc14n instead.)
//
// The implementation follows the two JDK classes the Transformer actually runs, and the
// method names below are theirs:
//
//   - com.sun.org.apache.xalan.internal.xsltc.trax.DOM2TO walks the DOM and turns it into
//     serializer events. It is where namespace fixup lives: declarations are replayed
//     from the attribute list first, then one is synthesized for the element's own prefix,
//     so a subtree serialized on its own regains the declarations it inherited.
//   - com.sun.org.apache.xml.internal.serializer.ToStream / ToXMLStream write the bytes.
//     They own the escaping rules, the empty-element form, the CDATA splitting and the
//     XML declaration.
//
// Consequences worth stating explicitly, because each was a divergence in the first port
// and each is pinned by xml/utils/testdata/serialize/goldens.txt:
//
//   - The XML declaration is written for EVERY node kind, not just for a Document:
//     DomUtils never sets OMIT_XML_DECLARATION. Serializing a text node emits
//     <?xml version="1.0" encoding="UTF-8"?> and then the text. getNodeBytes strips it
//     again, which is why the omission has gone unnoticed upstream.
//   - standalone="no" appears only when n is the Document node AND the parsed source did
//     not say standalone="yes"; a subtree never carries it (DOM2TO reads the document
//     info only for a DOCUMENT_NODE), and neither does a document that declared
//     standalone="yes" (setStandalone is called only for the false case).
//   - Attributes are written in the order Xerces' NamedNodeMap keeps them, which is
//     sorted by node name, not source order - then split into namespace declarations
//     first and everything else after.
//   - Empty elements are written <e/>.
//
// Parity means bug-compatibility where the JDK is buggy, because the bugs are inside the
// digest. Three are reproduced deliberately and each is documented where it lives:
// attributes are deduplicated by qualified name alone, so two that differ only in
// namespace collapse (addAttributeAlways); the UTF-8 writer mis-encodes code points from
// U+40000 up (writeXalanUTF8); and the US-ASCII writer truncates instead of substituting
// (writeTruncatingASCII).
//
// The line separator for a literal LF in text is Java's line.separator; the goldens were
// produced on Linux, so it is "\n" here.
func (n *Node) Serialize(w io.Writer, opts *SerializeOptions) error {
	b, err := n.Bytes(opts)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// SerializeOptions carries the two OutputKeys DomUtils.serializeNode can set. A nil
// *SerializeOptions is the DomUtils configuration itself.
type SerializeOptions struct {
	// XMLDeclaration false is OMIT_XML_DECLARATION=yes, which DomUtils never sets; it
	// exists for callers that want the markup alone.
	XMLDeclaration bool

	// Encoding overrides OutputKeys.ENCODING. Empty selects what DomUtils selects: the
	// document's own declared encoding when it has one, else UTF-8.
	Encoding string
}

var defaultSerializeOptions = SerializeOptions{XMLDeclaration: true}

// Bytes is Serialize into a buffer.
func (n *Node) Bytes(opts *SerializeOptions) ([]byte, error) {
	o := defaultSerializeOptions
	if opts != nil {
		o = *opts
	}

	// DomUtils: the ENCODING property is set from the OWNING document's xmlEncoding
	// (Utils.isStringNotBlank guards it), whether the node handed in is the document or
	// any node inside it. A detached subtree has no document and keeps the default.
	encName := o.Encoding
	if encName == "" {
		if declared := n.XMLEncoding(); strings.TrimSpace(declared) != "" {
			encName = declared
		}
	}
	if encName == "" {
		encName = "UTF-8"
	}
	enc, err := lookupOutEncoding(encName)
	if err != nil {
		return nil, err
	}

	s := &serializer{enc: enc, ns: map[string][]nsRec{}}
	// NamespaceMappings.initNamespaces: the default prefix maps to no namespace and
	// "xml" is predeclared, both at depth 0, under a sentinel that popNamespaces stops at.
	s.ns[""] = []nsRec{{uri: "", depth: 0}}
	s.ns["xml"] = []nsRec{{prefix: "xml", uri: XMLNamespace, depth: 0}}
	s.nsNode = []nsRec{{depth: -1}}

	s.out = append(s.out, enc.bom...)
	if o.XMLDeclaration {
		s.writeDeclaration(n)
	}
	s.walk(n)
	// endDocument -> flushPending closes a CDATA section left open by the last event.
	if s.cdataOpen {
		s.closeCDATA()
	}
	if s.err != nil {
		return nil, s.err
	}
	return s.out, nil
}

// writeDeclaration reproduces ToXMLStream.startDocumentInternal. m_doIndent is false, so
// no line separator ever follows it.
func (s *serializer) writeDeclaration(n *Node) {
	s.raw(`<?xml version="1.0" encoding="`)
	s.raw(s.enc.decl)
	s.raw(`"`)
	// DOM2TO.setDocumentInfo is reached only for a DOCUMENT_NODE, and calls
	// setStandalone only when getXmlStandalone() is false - the one call that sets
	// m_standaloneWasSpecified, which is what makes the pseudo-attribute appear.
	if n != nil && n.Kind == Document && !n.XMLStandalone() {
		s.raw(` standalone="no"`)
	}
	s.raw(`?>`)
}

// ------------------------------------------------------------------ serializer state

type nsRec struct {
	prefix string
	uri    string
	depth  int
}

// outAttr is one entry of ToStream's m_attributes, the list accumulated between
// startElement and the closing of the start tag.
type outAttr struct {
	uri   string
	local string
	qname string
	value string
}

type serializer struct {
	out []byte
	enc outEncoding
	err error

	// ToStream's m_elemContext, flattened: only the depth, the open-start-tag flag and
	// the pending attributes are ever read back.
	depth        int
	startTagOpen bool
	attrs        []outAttr
	attrIdx      map[string]int // qname -> index in attrs; built only for a wide start tag

	cdataOpen bool

	// NamespaceMappings: a stack per prefix plus the declaration-ordered node stack that
	// popNamespaces unwinds.
	ns     map[string][]nsRec
	nsNode []nsRec
}

func (s *serializer) fail(err error) {
	if s.err == nil {
		s.err = err
	}
}

// raw writes text that carries no escaping of its own - markup, names, comment and
// processing-instruction data. It still goes through the output encoding, and a character
// the encoding cannot represent becomes '?', which is what Java's CharsetEncoder does with
// its default REPLACE action.
func (s *serializer) raw(v string) {
	for _, r := range v {
		s.out = s.enc.write(s.out, r)
	}
}

func (s *serializer) char(r rune) { s.out = s.enc.write(s.out, r) }

// ------------------------------------------------------------------ DOM2TO traversal

func (s *serializer) walk(n *Node) {
	if n == nil || s.err != nil {
		return
	}
	switch n.Kind {
	// DOM2TO ignores ATTRIBUTE_NODE, DOCUMENT_TYPE_NODE, ENTITY_NODE,
	// ENTITY_REFERENCE_NODE and NOTATION_NODE: serializing an attribute yields the XML
	// declaration and nothing else.
	case Attribute:
		return

	case Document:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.walk(c)
		}

	case Element:
		s.element(n)

	case Text:
		s.characters(n.Value)

	case CDATA:
		s.cdataSection(n.Value)

	case Comment:
		s.comment(n.Value)

	case ProcInst:
		s.procInst(n.Name.Local, n.Value)
	}
}

func (s *serializer) element(n *Node) {
	qname := n.Name.QName()
	s.startElement(qname)

	attrs := domAttrOrder(n.Attrs)

	// Pass one: every attribute whose NODE NAME starts with "xmlns" is replayed as a
	// namespace declaration - including "xmlnsfoo" and "xmlns:", which is how a plain
	// setAttribute("xmlns:p", uri) still reaches the prefix map. The prefix is what
	// follows the LAST colon, so "xmlns" and "xmlnsfoo" both declare the default one.
	for _, a := range attrs {
		q := a.Name.QName()
		if !strings.HasPrefix(q, "xmlns") {
			continue
		}
		prefix := ""
		if colon := strings.LastIndex(q, ":"); colon > 0 {
			prefix = q[colon+1:]
		}
		s.startPrefixMapping(prefix, a.Value)
	}

	// Pass two: everything else. DOM2TO builds a THROWAWAY NamespaceMappings per element
	// and never pushes into it, so the generated prefixes restart at ns0 on every element
	// and repeat for two attributes that share a namespace.
	genCount := 0
	for _, a := range attrs {
		q := a.Name.QName()
		if strings.HasPrefix(q, "xmlns") {
			continue
		}
		uriAttr := a.Name.Space
		if uriAttr == "" {
			s.addAttribute(q, a.Value)
			continue
		}
		// Fix for bug 26319: an attribute carrying a namespace but no prefix has to be
		// given one. The throwaway mappings only ever match the predeclared "xml".
		newPrefix := ""
		if uriAttr == XMLNamespace {
			newPrefix = "xml"
		} else {
			newPrefix = "ns" + strconv.Itoa(genCount)
			genCount++
		}
		prefix := newPrefix
		if colon := strings.LastIndex(q, ":"); colon > 0 {
			prefix = q[:colon]
		}
		s.startPrefixMapping(prefix, uriAttr)
		// The concatenation is DOM2TO's, and it is patchName that turns the resulting
		// "n1:n1:Id" back into "n1:Id".
		s.addAttribute(prefix+":"+q, a.Value)
	}

	// The element's own namespace, last - which is why a fixup declaration for the
	// element's prefix trails the attributes when no attribute happened to declare it.
	if uri := n.Name.Space; uri != "" {
		prefix := ""
		if colon := strings.LastIndex(qname, ":"); colon > 0 {
			prefix = qname[:colon]
		}
		s.startPrefixMapping(prefix, uri)
	} else {
		// Fix for bug 26319: an element created with a null namespace must be able to
		// escape an inherited default declaration, so "" is mapped to "". It only
		// produces xmlns="" when something else had mapped the default namespace.
		s.startPrefixMapping("", "")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.walk(c)
	}
	s.endElement(qname)
}

// domAttrOrder returns the attributes in the order Xerces' NamedNodeMap holds them.
//
// AttributeMap keeps its entries sorted by node name: findNamePoint binary-searches the
// list and both setNamedItem and setNamedItemNS insert at the position it reports. So the
// serialized order has nothing to do with the source order - "b a c" comes back "a b c" -
// and the comparison is String.compareTo, i.e. UTF-16 code units.
func domAttrOrder(attrs []*Node) []*Node {
	if len(attrs) < 2 {
		return attrs
	}
	// One stable sort instead of a binary-search insertion per attribute, whose slice
	// shifts cost O(n^2) for an element carrying n attributes. Xerces inserts an
	// attribute BEFORE an equal-named one already present, so equal names would come out
	// in reverse source order; feeding the sort the reversed input reproduces exactly
	// that (the parser rejects duplicate QNames, so it only matters for a hand-built tree).
	type keyed struct {
		name string
		node *Node
	}
	items := make([]keyed, len(attrs))
	for i, a := range attrs {
		items[len(attrs)-1-i] = keyed{a.Name.QName(), a}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return compareUTF16(items[i].name, items[j].name) < 0
	})
	out := make([]*Node, len(items))
	for i, it := range items {
		out[i] = it.node
	}
	return out
}

// compareUTF16 orders two strings the way String.compareTo does, by UTF-16 code unit.
// This differs from Go's byte-wise ordering only when one string holds a supplementary
// character, whose surrogates sort below U+E000 in UTF-16 and above it in UTF-8.
func compareUTF16(a, b string) int {
	// Skip the common ASCII prefix: over ASCII the two orders agree, and almost every real
	// comparison is decided there without allocating.
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] && a[i] < utf8.RuneSelf {
		i++
	}
	if i == len(a) && i == len(b) {
		return 0
	}
	if i < len(a) && i < len(b) && a[i] < utf8.RuneSelf && b[i] < utf8.RuneSelf {
		return int(a[i]) - int(b[i])
	}
	ar, br := []rune(a[i:]), []rune(b[i:])
	au, bu := utf16.Encode(ar), utf16.Encode(br)
	for i := 0; i < len(au) && i < len(bu); i++ {
		if au[i] != bu[i] {
			return int(au[i]) - int(bu[i])
		}
	}
	return len(au) - len(bu)
}

// ------------------------------------------------------------- namespace bookkeeping

// lookupNamespace is NamespaceMappings.lookupNamespace; ok is false where Java returns null.
func (s *serializer) lookupNamespace(prefix string) (string, bool) {
	st := s.ns[prefix]
	if len(st) == 0 {
		return "", false
	}
	return st[len(st)-1].uri, true
}

// pushNamespace is NamespaceMappings.pushNamespace. It answers false - meaning "no
// declaration is needed" - for any prefix starting with "xml" (so "xmlfoo:a" is emitted
// with no declaration at all) and for a prefix already bound to this very URI, which is
// what suppresses a redundant redeclaration in the middle of a document.
func (s *serializer) pushNamespace(prefix, uri string, depth int) bool {
	if strings.HasPrefix(prefix, "xml") {
		return false
	}
	st := s.ns[prefix]
	if len(st) > 0 && st[len(st)-1].uri == uri {
		return false
	}
	rec := nsRec{prefix: prefix, uri: uri, depth: depth}
	s.ns[prefix] = append(st, rec)
	s.nsNode = append(s.nsNode, rec)
	return true
}

func (s *serializer) popNamespaces(depth int) {
	for len(s.nsNode) > 0 {
		rec := s.nsNode[len(s.nsNode)-1]
		if rec.depth < depth {
			return
		}
		s.nsNode = s.nsNode[:len(s.nsNode)-1]
		if strings.HasPrefix(rec.prefix, "xml") {
			continue
		}
		if st := s.ns[rec.prefix]; len(st) > 0 {
			s.ns[rec.prefix] = st[:len(st)-1]
		}
	}
}

// startPrefixMapping is ToStream.startPrefixMapping(prefix, uri, false), which is what
// namespaceAfterStartElement calls. The declaration it queues is an ordinary attribute of
// the open start tag, so it takes its place in the emitted order and is deduplicated
// against an existing one by (namespace, local name).
func (s *serializer) startPrefixMapping(prefix, uri string) {
	if !s.pushNamespace(prefix, uri, s.depth) {
		return
	}
	if prefix == "" {
		s.addAttributeAlways(XMLNSNamespace, "xmlns", "xmlns", uri)
		return
	}
	// A prefix mapped to the empty URI is remembered but not written: there is no legal
	// markup for undeclaring a prefix in XML 1.0.
	if uri != "" {
		s.addAttributeAlways(XMLNSNamespace, prefix, "xmlns:"+prefix, uri)
	}
}

// ------------------------------------------------------------------ attribute list

// addAttributeAlways is ToStream.doAddAttributeAlways, which overrides SerializerBase's
// and looks an attribute up by its QUALIFIED NAME alone - "Don't use 'localName' as it
// gives incorrect value, rely only on 'rawName'". Two attributes that differ only in
// namespace therefore collapse into one, the later value winning: an element carrying both
// setAttributeNS(A, "p:Id") and setAttributeNS(B, "p:Id") serializes with a single p:Id.
//
// The lookup is a linear scan for an ordinary start tag and switches to a map once the tag
// is wide, because n attributes cost O(n^2) comparisons otherwise - tens of seconds for an
// attacker-supplied element with 100 000 of them. The result is the same either way.
func (s *serializer) addAttributeAlways(uri, local, qname, value string) {
	if s.attrIdx != nil {
		if i, ok := s.attrIdx[qname]; ok {
			s.attrs[i].value = value
			return
		}
	} else {
		for i := range s.attrs {
			if s.attrs[i].qname == qname {
				s.attrs[i].value = value
				return
			}
		}
	}
	s.attrs = append(s.attrs, outAttr{uri: uri, local: local, qname: qname, value: value})
	switch {
	case s.attrIdx != nil:
		s.attrIdx[qname] = len(s.attrs) - 1
	case len(s.attrs) > attrIndexThreshold:
		s.attrIdx = make(map[string]int, 2*len(s.attrs))
		for i := range s.attrs {
			s.attrIdx[s.attrs[i].qname] = i
		}
	}
}

// attrIndexThreshold is the pending-attribute count above which addAttributeAlways indexes
// them by qname.
const attrIndexThreshold = 12

// addAttribute is SerializerBase.addAttribute(String, String).
func (s *serializer) addAttribute(name, value string) {
	if !s.startTagOpen {
		return
	}
	patched := s.patchName(name)
	uri, err := s.namespaceURIOf(patched, false)
	if err != nil {
		s.fail(err)
		return
	}
	s.addAttributeAlways(uri, localNameOf(patched), patched, value)
}

// patchName is SerializerBase.patchName: it collapses the doubled prefix DOM2TO builds
// ("n1:n1:Id" -> "n1:Id") and drops a prefix that is bound to the empty namespace.
func (s *serializer) patchName(qname string) string {
	lastColon := strings.LastIndex(qname, ":")
	if lastColon > 0 {
		firstColon := strings.Index(qname, ":")
		prefix := qname[:firstColon]
		local := qname[lastColon+1:]
		if uri, ok := s.lookupNamespace(prefix); ok && uri == "" {
			return local
		}
		if firstColon != lastColon {
			return prefix + ":" + local
		}
	}
	return qname
}

func localNameOf(qname string) string {
	if col := strings.LastIndex(qname, ":"); col > 0 {
		return qname[col+1:]
	}
	return qname
}

// namespaceURIOf is SerializerBase.getNamespaceURI. Its RuntimeException for an
// undeclared prefix is reachable: setAttribute("p:x", v) on an element whose own prefix
// is p still throws, because DOM2TO processes the attributes before it declares the
// element's own namespace.
func (s *serializer) namespaceURIOf(qname string, isElement bool) (string, error) {
	col := strings.LastIndex(qname, ":")
	prefix := ""
	if col > 0 {
		prefix = qname[:col]
	}
	if prefix == "" && !isElement {
		return "", nil
	}
	uri, ok := s.lookupNamespace(prefix)
	if !ok && prefix != "xmlns" {
		return "", fmt.Errorf("Namespace for prefix '%s' has not been declared.", prefix)
	}
	return uri, nil
}

// ------------------------------------------------------------------ ToStream events

func (s *serializer) startElement(qname string) {
	if s.cdataOpen {
		s.closeCDATA()
	}
	if s.startTagOpen {
		s.closeStartTag()
	}
	s.raw("<")
	s.raw(qname)
	s.depth++
	s.startTagOpen = true
}

func (s *serializer) closeStartTag() {
	s.processAttributes()
	s.raw(">")
	s.startTagOpen = false
}

func (s *serializer) processAttributes() {
	for _, a := range s.attrs {
		s.raw(" ")
		s.raw(a.qname)
		s.raw(`="`)
		s.writeAttrString(a.value)
		s.raw(`"`)
	}
	s.attrs = s.attrs[:0]
	s.attrIdx = nil
}

func (s *serializer) endElement(qname string) {
	// Namespaces declared at this depth stop applying here.
	s.popNamespaces(s.depth)
	if s.startTagOpen {
		s.processAttributes()
		s.raw("/>")
	} else {
		if s.cdataOpen {
			s.closeCDATA()
		}
		s.raw("</")
		s.raw(qname)
		s.raw(">")
	}
	s.startTagOpen = false
	s.depth--
}

// characters is ToStream.characters for a text node. An empty one produces nothing at
// all - not even the closing of an open start tag.
func (s *serializer) characters(text string) {
	if text == "" {
		return
	}
	if s.startTagOpen {
		s.closeStartTag()
	}
	if s.cdataOpen {
		s.closeCDATA()
	}
	s.outputCharacters(text)
}

// cdataSection is DOM2TO's startCDATA / characters / endCDATA triple.
func (s *serializer) cdataSection(text string) {
	if text != "" {
		if s.startTagOpen {
			s.closeStartTag()
		}
		s.cdata(text)
	}
	if s.cdataOpen {
		s.closeCDATA()
	}
}

func (s *serializer) closeCDATA() {
	s.raw("]]>")
	s.cdataOpen = false
}

// cdata is ToStream.cdata. The opening delimiter is written up front only when the first
// character is one the encoding can carry; writeNormalizedChars opens it later otherwise.
// A section whose last character is ']' is closed immediately so that it cannot merge
// with a following one into "]]>".
func (s *serializer) cdata(text string) {
	runes := []rune(text)
	writeBrackets := len(runes) >= 1 && s.escapingNotNeeded(runes[0])
	if writeBrackets && !s.cdataOpen {
		s.raw("<![CDATA[")
		s.cdataOpen = true
	}
	s.writeNormalizedChars(runes)
	if writeBrackets && runes[len(runes)-1] == ']' {
		s.closeCDATA()
	}
}

// writeNormalizedChars is ToStream.writeNormalizedChars with isCData true.
func (s *serializer) writeNormalizedChars(runes []rune) {
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case c == '\n':
			s.raw("\n")

		case !s.escapingNotNeeded(c):
			// A character the encoding cannot carry leaves the CDATA section as a
			// character reference - except a supplementary one, which is written as the
			// raw pair whenever the PAIR is encodable.
			if c > 0xFFFF && s.enc.pairInEncoding(c) {
				if !s.cdataOpen {
					s.raw("<![CDATA[")
					s.cdataOpen = true
				}
				s.char(c)
				continue
			}
			s.writeCharRef(c)

		case i < len(runes)-2 && c == ']' && runes[i+1] == ']' && runes[i+2] == '>':
			s.raw("]]]]><![CDATA[>")
			i += 2

		default:
			if !s.cdataOpen {
				s.raw("<![CDATA[")
				s.cdataOpen = true
			}
			s.char(c)
		}
	}
}

// comment is ToStream.comment. Comment data is written verbatim except that a run of two
// dashes is broken with a space and a trailing dash gains one, so that data no XML parser
// would have produced still comes back out as a well-formed comment.
func (s *serializer) comment(text string) {
	if s.startTagOpen {
		s.closeStartTag()
	}
	if s.cdataOpen {
		s.closeCDATA()
	}
	s.raw("<!--")
	runes := []rune(text)
	start := 0
	wasDash := false
	for i := 0; i < len(runes); i++ {
		if wasDash && runes[i] == '-' {
			s.rawRunes(runes[start:i])
			s.raw(" -")
			start = i + 1
		}
		wasDash = runes[i] == '-'
	}
	if len(runes) > 0 {
		if start < len(runes) {
			s.rawRunes(runes[start:])
		}
		if runes[len(runes)-1] == '-' {
			s.raw(" ")
		}
	}
	s.raw("-->")
}

func (s *serializer) rawRunes(runes []rune) {
	for _, r := range runes {
		s.char(r)
	}
}

// procInst is ToXMLStream.processingInstruction. The separating space is written only
// when the data does not already start with one - and "starts with one" is
// Character.isSpaceChar, which is the Unicode separator categories and therefore NOT tab
// or newline. A "?>" inside the data is broken as "? >".
func (s *serializer) procInst(target, data string) {
	if s.startTagOpen {
		s.closeStartTag()
	}
	s.raw("<?")
	s.raw(target)
	if data != "" && !isJavaSpaceChar([]rune(data)[0]) {
		s.raw(" ")
	}
	if idx := strings.Index(data, "?>"); idx >= 0 {
		if idx > 0 {
			s.raw(data[:idx])
		}
		s.raw("? >")
		if idx+2 < len(data) {
			s.raw(data[idx+2:])
		}
	} else {
		s.raw(data)
	}
	s.raw("?>")
}

// isJavaSpaceChar is Character.isSpaceChar: the three Unicode separator categories, which
// exclude the ASCII control whitespace.
func isJavaSpaceChar(r rune) bool {
	return unicode.In(r, unicode.Zs, unicode.Zl, unicode.Zp)
}

// ------------------------------------------------------------------------ escaping

// escapingNotNeeded is ToStream.escapingNotNeeded. Below 127 the answer is fixed; above
// it, it is the output encoding's - and for a supplementary character Java asks about the
// HIGH SURROGATE code unit on its own, which is why UTF-8 (where a lone surrogate encodes
// to '?') escapes every astral character while UTF-16 writes it literally.
func (s *serializer) escapingNotNeeded(r rune) bool {
	if r < 127 {
		return r >= 0x20 || r == 0x0A || r == 0x0D || r == 0x09
	}
	if r > 0xFFFF {
		h, _ := utf16.EncodeRune(r)
		return s.enc.inEncoding(uint16(h))
	}
	return s.enc.inEncoding(uint16(r))
}

// isSpecialTextChar and isSpecialAttrChar are CharInfo's two tables for METHOD=xml. The
// bit set holds the four XML entities plus LF and CR; the text table then re-clears the
// double quote and the attribute table additionally sets the tab.
func isSpecialTextChar(r rune) bool {
	switch r {
	case '&', '<', '>', 0x0A, 0x0D:
		return true
	}
	return r < 0x20 && r != 0x09
}

func isSpecialAttrChar(r rune) bool {
	switch r {
	case '"', '&', '<', '>', 0x09, 0x0A, 0x0D:
		return true
	}
	return false
}

func entityForChar(r rune) (string, bool) {
	switch r {
	case '"':
		return "&quot;", true
	case '&':
		return "&amp;", true
	case '<':
		return "&lt;", true
	case '>':
		return "&gt;", true
	}
	return "", false
}

// isCharacterInC0orC1Range is ToStream's, verbatim: tab, LF and CR are excluded, so what
// is left is C0 minus those, plus DEL and C1.
func isCharacterInC0orC1Range(r rune) bool {
	if r == 0x09 || r == 0x0A || r == 0x0D {
		return false
	}
	return (r >= 0x7F && r <= 0x9F) || (r >= 0x01 && r <= 0x1F)
}

// isInvalidXMLChar is XMLChar.isInvalid, the check that makes ToStream refuse a control
// character in XML 1.0 rather than escape it. Only a hand-built tree can reach it: Parse
// rejects the same characters.
func isInvalidXMLChar(r rune) bool {
	switch {
	case r == 0x09 || r == 0x0A || r == 0x0D:
		return false
	case r >= 0x20 && r <= 0xD7FF:
		return false
	case r >= 0xE000 && r <= 0xFFFD:
		return false
	case r >= 0x10000 && r <= 0x10FFFF:
		return false
	}
	return true
}

func (s *serializer) writeCharRef(r rune) {
	if s.cdataOpen {
		s.closeCDATA()
	}
	s.raw("&#")
	s.raw(strconv.Itoa(int(r)))
	s.raw(";")
}

// outputCharacters is ToStream.outputCharacters: the character data of a text node.
//
// The two loops in the original - one skipping leading whitespace, one for the rest -
// agree on every character they can both see (CharInfo's clean-text table marks LF and CR
// dirty and tab and space clean, which is exactly what the main loop's condition
// decides), so one pass reproduces both.
func (s *serializer) outputCharacters(text string) {
	for _, r := range text {
		clean := (!isCharacterInC0orC1Range(r) && s.escapingNotNeeded(r) && !isSpecialTextChar(r)) || r == '"'
		if clean {
			s.char(r)
			continue
		}
		if r == '\n' {
			s.raw("\n")
			continue
		}
		s.accumDefaultEscape(r, true, false)
	}
}

// writeAttrString is ToStream.writeAttrString: the value of an attribute.
func (s *serializer) writeAttrString(value string) {
	for _, r := range value {
		if s.escapingNotNeeded(r) && !isSpecialAttrChar(r) {
			s.char(r)
			continue
		}
		s.accumDefaultEscape(r, false, true)
	}
}

// accumDefaultEscape is ToStream.accumDefaultEscape over accumDefaultEntity.
//
// The named entity, when the character has one, wins first and does not consult the
// depth. Everything else - tab, LF, CR, a character outside the encoding - falls to the
// character-reference branch, which is guarded by m_currentElemDepth > 0. That guard is
// observable: serializing a bare text node puts its characters at depth 0, so a carriage
// return in it is written literally where the same text inside an element becomes &#13;.
func (s *serializer) accumDefaultEscape(r rune, fromTextNode, escLF bool) {
	if !escLF && r == '\n' {
		s.raw("\n")
		return
	}
	special := isSpecialAttrChar(r)
	if fromTextNode {
		special = isSpecialTextChar(r)
	}
	if special {
		if ent, ok := entityForChar(r); ok {
			s.raw(ent)
			return
		}
	}

	// A supplementary character is always a character reference here: Java has already
	// established that its high surrogate is not writable in this encoding.
	if r > 0xFFFF {
		s.writeCharRef(r)
		return
	}
	if isInvalidXMLChar(r) {
		s.fail(fmt.Errorf("invalid character 0x%x in the node text", r))
		return
	}
	if isCharacterInC0orC1Range(r) {
		s.writeCharRef(r)
		return
	}
	if (!s.escapingNotNeeded(r) || special) && s.depth > 0 {
		s.writeCharRef(r)
		return
	}
	s.char(r)
}
