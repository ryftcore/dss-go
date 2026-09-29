package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DefaultMaxDepth is the element nesting limit applied when ParseOptions.MaxDepth is
// zero. It exists because encoding/xml recurses no further than we ask it to, but our
// own tree walks do, and an unbounded document is an easy denial of service.
const DefaultMaxDepth = 500

// ParseOptions tunes Parse. The zero value is the DSS-equivalent secure posture.
type ParseOptions struct {
	AllowDoctype  bool  // default false: any DOCTYPE declaration is a SyntaxError
	MaxDepth      int   // default 500 when zero
	MaxBytes      int64 // default 0 = unlimited
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)
}

func (o *ParseOptions) normalized() ParseOptions {
	var out ParseOptions
	if o != nil {
		out = *o
	}
	if out.MaxDepth <= 0 {
		out.MaxDepth = DefaultMaxDepth
	}
	return out
}

// Parse builds a document from src. opts may be nil, which selects the DSS-equivalent
// secure defaults (no DOCTYPE, no external anything, UTF-8/US-ASCII/ISO-8859-1/UTF-16).
// The returned node has Kind Document.
func Parse(src []byte, opts *ParseOptions) (*Node, error) {
	o := opts.normalized()
	if o.MaxBytes > 0 && int64(len(src)) > o.MaxBytes {
		return nil, &SyntaxError{Line: 1, Column: 1, Offset: int64(len(src)),
			Msg: fmt.Sprintf("document is %d bytes, exceeding MaxBytes %d", len(src), o.MaxBytes)}
	}
	buf, decl, err := decodeSource(src, o.CharsetReader)
	if err != nil {
		return nil, err
	}
	p := &parser{buf: buf, opts: o, decl: decl}
	return p.run()
}

// ParseReader reads r to completion and calls Parse.
//
// It buffers the whole stream, and with nil options or MaxBytes == 0 it does so WITHOUT a size
// bound (as DocumentBuilder does upstream): a caller reading an untrusted stream must set
// ParseOptions.MaxBytes, which caps what is read (at most MaxBytes+1 bytes) as well as what
// is parsed. No caller inside this module uses ParseReader; the DSS entry points read their
// input into a []byte and hand it to Parse.
func ParseReader(r io.Reader, opts *ParseOptions) (*Node, error) {
	o := opts.normalized()
	var (
		src []byte
		err error
	)
	if o.MaxBytes > 0 {
		src, err = io.ReadAll(io.LimitReader(r, o.MaxBytes+1))
	} else {
		src, err = io.ReadAll(r)
	}
	if err != nil {
		return nil, err
	}
	return Parse(src, opts)
}

// nsBinding is one prefix-to-URI mapping. The default namespace uses prefix "".
type nsBinding struct{ prefix, uri string }

type parser struct {
	buf  []byte
	opts ParseOptions
	decl xmlDecl
	dec  *xml.Decoder

	doc   *Node
	stack []*Node // stack[0] is the document node

	ns     []nsBinding
	nsMark []int

	prev  int64 // end offset of the previous token = start offset of the next
	roots int
}

func (p *parser) run() (*Node, error) {
	p.dec = xml.NewDecoder(bytes.NewReader(p.buf))
	// The buffer is already UTF-8; a declaration naming anything else still makes
	// encoding/xml ask for a converter, and the identity one is the right answer
	// because decodeSource has done the conversion. Probed: InputOffset stays
	// consistent across the switch, which is what the raw-span recovery depends on.
	p.dec.CharsetReader = func(_ string, in io.Reader) (io.Reader, error) { return in, nil }
	p.doc = NewDocument()
	// Xerces keeps the declaration's encoding and standalone on the Document; the
	// serializer reads both back (see XMLEncoding/XMLStandalone).
	p.doc.doc.xmlEncoding = p.decl.encoding
	p.doc.doc.xmlStandalone = p.decl.standalone == "yes"
	p.stack = []*Node{p.doc}

	for {
		start := p.prev
		tok, err := p.dec.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, p.wrap(err)
		}
		end := p.dec.InputOffset()
		raw := p.buf[start:end]
		p.prev = end

		if err := p.token(tok, raw, start); err != nil {
			return nil, err
		}
	}

	if len(p.stack) > 1 {
		el := p.stack[len(p.stack)-1]
		return nil, p.errAt(p.prev, fmt.Sprintf("unexpected end of document: <%s> is not closed", el.Name.QName()))
	}
	if p.roots != 1 {
		return nil, p.errAt(p.prev, "document must have exactly one element child")
	}
	return p.doc, nil
}

func (p *parser) token(tok xml.Token, raw []byte, off int64) error {
	switch t := tok.(type) {
	case xml.StartElement:
		return p.startElement(t, raw, off)
	case xml.EndElement:
		return p.endElement(t, raw, off)
	case xml.CharData:
		return p.charData(t, raw, off)
	case xml.Comment:
		p.parent().link(NewComment(string(t)), nil)
		return nil
	case xml.ProcInst:
		return p.procInst(t, off)
	case xml.Directive:
		return p.directive(t, off)
	}
	return nil
}

func (p *parser) parent() *Node { return p.stack[len(p.stack)-1] }

// ---------------------------------------------------------------- elements

func (p *parser) startElement(se xml.StartElement, raw []byte, off int64) error {
	if len(p.stack) > p.opts.MaxDepth {
		return p.errAt(off, fmt.Sprintf("element nesting deeper than MaxDepth %d", p.opts.MaxDepth))
	}
	if p.parent().Kind == Document {
		p.roots++
		if p.roots > 1 {
			return p.errAt(off, "document must have exactly one element child")
		}
	}

	tagName, rattrs, _, err := scanStartTag(raw)
	if err != nil {
		if ave, ok := err.(*attValueError); ok {
			return p.errAt(off+int64(ave.off), ave.Error())
		}
		return p.errAt(off, err.Error())
	}
	elPrefix, elLocal, err := splitQName(tagName)
	if err != nil {
		return p.errAt(off+1, err.Error())
	}
	// Cross-check the rescan against encoding/xml's own tokenization. This is the
	// only place an xml.Attr is looked at, and only its Name: reading Attr.Value
	// would defeat the entire point of the rescan.
	if elPrefix != se.Name.Space || elLocal != se.Name.Local || len(rattrs) != len(se.Attr) {
		return p.errAt(off, errInternalScan.Error())
	}
	for i, a := range rattrs {
		if a.prefix != se.Attr[i].Name.Space || a.local != se.Attr[i].Name.Local {
			return p.errAt(off, errInternalScan.Error())
		}
	}

	p.pushScope()
	if err := p.declareNamespaces(rattrs, off); err != nil {
		p.popScope()
		return err
	}

	el := &Node{Kind: Element}
	if elPrefix == "xmlns" {
		p.popScope()
		return p.errAt(off+1, `"xmlns" is not a usable element name prefix`)
	}
	uri, ok := p.resolve(elPrefix)
	if !ok {
		p.popScope()
		return p.errAt(off+1, fmt.Sprintf("undeclared namespace prefix %q on element <%s>", elPrefix, tagName))
	}
	el.Name = Name{Space: uri, Local: elLocal, Prefix: elPrefix}

	if err := p.buildAttrs(el, rattrs, tagName, off); err != nil {
		p.popScope()
		return err
	}

	p.parent().link(el, nil)
	p.stack = append(p.stack, el)
	return nil
}

// buildAttrs resolves every attribute name and installs the attribute nodes in
// document order, rejecting duplicates by literal QName and by expanded name.
func (p *parser) buildAttrs(el *Node, rattrs []rawAttr, tagName string, off int64) error {
	if len(rattrs) == 0 {
		return nil
	}
	el.Attrs = make([]*Node, 0, len(rattrs))
	byQName := make(map[string]struct{}, len(rattrs))
	byExpanded := make(map[Name]struct{}, len(rattrs))

	for _, a := range rattrs {
		var name Name
		switch {
		case a.prefix == "" && a.local == "xmlns":
			name = Name{Space: XMLNSNamespace, Local: "xmlns"}
		case a.prefix == "xmlns":
			name = Name{Space: XMLNSNamespace, Local: a.local, Prefix: "xmlns"}
		case a.prefix == "":
			// An unprefixed attribute is in no namespace - never in the default one.
			name = Name{Local: a.local}
		default:
			uri, ok := p.resolve(a.prefix)
			if !ok {
				return p.errAt(off+int64(a.off), fmt.Sprintf("undeclared namespace prefix %q on attribute %s of <%s>", a.prefix, a.qname(), tagName))
			}
			name = Name{Space: uri, Local: a.local, Prefix: a.prefix}
		}

		q := a.qname()
		if _, dup := byQName[q]; dup {
			return p.errAt(off+int64(a.off), fmt.Sprintf("duplicate attribute %s on <%s>", q, tagName))
		}
		byQName[q] = struct{}{}
		expanded := Name{Space: name.Space, Local: name.Local}
		if _, dup := byExpanded[expanded]; dup {
			return p.errAt(off+int64(a.off), fmt.Sprintf("duplicate attribute %s on <%s>: {%s}%s is already present under another prefix", q, tagName, name.Space, name.Local))
		}
		byExpanded[expanded] = struct{}{}

		el.Attrs = append(el.Attrs, &Node{Kind: Attribute, Name: name, Value: a.value, Parent: el})
	}
	return nil
}

func (p *parser) endElement(ee xml.EndElement, raw []byte, off int64) error {
	if len(p.stack) <= 1 {
		return p.errAt(off, fmt.Sprintf("unexpected end tag </%s>", ee.Name.Local))
	}
	el := p.parent()
	// A self-closing tag yields a synthetic EndElement whose source span is empty;
	// there is nothing to check because encoding/xml manufactured it from the start
	// tag. A real end tag is checked, because RawToken accepts <r></s> silently.
	if len(raw) > 0 && (ee.Name.Space != el.Name.Prefix || ee.Name.Local != el.Name.Local) {
		closing := ee.Name.Local
		if ee.Name.Space != "" {
			closing = ee.Name.Space + ":" + ee.Name.Local
		}
		return p.errAt(off, fmt.Sprintf("end tag </%s> does not match start tag <%s>", closing, el.Name.QName()))
	}
	p.stack = p.stack[:len(p.stack)-1]
	p.popScope()
	return nil
}

// ---------------------------------------------------------------- leaves

func (p *parser) charData(cd xml.CharData, raw []byte, off int64) error {
	// encoding/xml reuses the buffer behind a CharData across calls, so the value
	// must be copied out; the []byte-to-string conversion does that.
	isCDATA := bytes.HasPrefix(raw, []byte("<![CDATA["))
	data := string(cd)

	if p.parent().Kind == Document {
		if isCDATA {
			return p.errAt(off, "CDATA section outside the document element")
		}
		if strings.TrimLeft(data, " \t\n\r") != "" {
			return p.errAt(off, "character data outside the document element")
		}
		return nil // whitespace between prolog/epilog nodes produces no node
	}
	// A U+FFFD in the decoded run is the only signal that encoding/xml may have
	// swallowed a surrogate character reference; see checkTextRefs.
	if !isCDATA && strings.ContainsRune(data, '�') {
		if err := checkTextRefs(raw); err != nil {
			return p.errAt(off, err.Error())
		}
	}
	kind := Text
	if isCDATA {
		kind = CDATA
	}
	p.parent().link(&Node{Kind: kind, Value: data}, nil)
	return nil
}

func (p *parser) procInst(pi xml.ProcInst, off int64) error {
	if strings.EqualFold(pi.Target, "xml") {
		// The XML declaration is not a node: DOM consumes it for its version,
		// encoding and standalone pseudo-attributes and discards it, and so has
		// decodeSource. Anywhere else the target is reserved and Xerces refuses it.
		if off == 0 && pi.Target == "xml" {
			return nil
		}
		return p.errAt(off, fmt.Sprintf("processing-instruction target %q is reserved", pi.Target))
	}
	p.parent().link(NewProcInst(pi.Target, string(pi.Inst)), nil)
	return nil
}

func (p *parser) directive(d xml.Directive, off int64) error {
	body := strings.TrimLeft(string(d), " \t\n\r")
	if !strings.HasPrefix(body, "DOCTYPE") {
		return p.errAt(off, "markup declarations are not allowed outside a DTD")
	}
	if !p.opts.AllowDoctype {
		// DocumentBuilderFactoryBuilder sets disallow-doctype-decl=true and
		// Santuario's XMLParserImpl passes disallowDocTypeDeclarations=true, so
		// upstream refuses the document outright. Nothing in DSS opts out.
		return p.errAt(off, "DOCTYPE declarations are not allowed")
	}
	// Tolerated but never processed: no external entities, no parameter entities, no
	// internal subset, no attribute defaulting, no DTD-declared ID attributes. An
	// entity reference other than the five predefined ones stays a hard error.
	return nil
}

// ---------------------------------------------------------------- namespaces

func (p *parser) pushScope() { p.nsMark = append(p.nsMark, len(p.ns)) }

func (p *parser) popScope() {
	mark := p.nsMark[len(p.nsMark)-1]
	p.ns = p.ns[:mark]
	p.nsMark = p.nsMark[:len(p.nsMark)-1]
}

// resolve maps a literal prefix to a namespace URI. The empty prefix asks for the
// default namespace, whose absence is not an error - it simply means no namespace.
func (p *parser) resolve(prefix string) (string, bool) {
	switch prefix {
	case "xml":
		return XMLNamespace, true
	case "xmlns":
		return XMLNSNamespace, true
	}
	for i := len(p.ns) - 1; i >= 0; i-- {
		if p.ns[i].prefix == prefix {
			return p.ns[i].uri, true
		}
	}
	return "", prefix == ""
}

// declareNamespaces installs this element's declarations and enforces the constraints
// of Namespaces in XML 1.0 third edition that encoding/xml does not know about.
func (p *parser) declareNamespaces(rattrs []rawAttr, off int64) error {
	for _, a := range rattrs {
		at := off + int64(a.off)
		switch {
		case a.prefix == "" && a.local == "xmlns":
			switch a.value {
			case XMLNSNamespace:
				return p.errAt(at, "the xmlns namespace name must not be declared as the default namespace")
			case XMLNamespace:
				return p.errAt(at, "the xml namespace name must not be declared as the default namespace")
			}
			p.ns = append(p.ns, nsBinding{prefix: "", uri: a.value})

		case a.prefix == "xmlns":
			switch {
			case a.local == "xmlns":
				return p.errAt(at, `the "xmlns" prefix must not be declared`)
			case a.value == "":
				// Xerces: "Prefixed namespace bindings may not be empty."
				// encoding/xml accepts it silently.
				return p.errAt(at, fmt.Sprintf("prefixed namespace binding xmlns:%s may not be empty", a.local))
			case a.local == "xml" && a.value != XMLNamespace:
				return p.errAt(at, fmt.Sprintf("the xml prefix must be bound to %s, not %q", XMLNamespace, a.value))
			case a.local != "xml" && a.value == XMLNamespace:
				return p.errAt(at, fmt.Sprintf("only the xml prefix may be bound to %s", XMLNamespace))
			case a.value == XMLNSNamespace:
				return p.errAt(at, fmt.Sprintf("no prefix may be bound to %s", XMLNSNamespace))
			}
			p.ns = append(p.ns, nsBinding{prefix: a.local, uri: a.value})
		}
	}
	return nil
}

// ---------------------------------------------------------------- errors

func (p *parser) errAt(off int64, msg string) *SyntaxError {
	line, col := lineCol(p.buf, off)
	return &SyntaxError{Line: line, Column: col, Offset: off, Msg: msg}
}

func (p *parser) wrap(err error) error {
	msg := err.Error()
	if se, ok := err.(*xml.SyntaxError); ok {
		msg = se.Msg
	}
	return p.errAt(p.dec.InputOffset(), msg)
}
