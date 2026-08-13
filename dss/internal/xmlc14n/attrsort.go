// Ported from org.apache.xml.security.c14n.helper.AttrCompare (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import (
	"bufio"
	"sort"
	"unicode/utf16"
	"unicode/utf8"
)

// xmlnsPrefix is the key under which the default namespace declaration lives, both in
// NameSpaceSymbTable and as the local name of an xmlns="..." attribute node.
const xmlnsPrefix = "xmlns"

// xmlnsNamespace and xmlNamespace repeat xmldom's constants so that the files which do not
// depend on xmldom stay independent of it. They are asserted equal in c14n.go.
const (
	xmlnsNamespace = "http://www.w3.org/2000/xmlns/"
	xmlNamespace   = "http://www.w3.org/XML/1998/namespace"
	xmlPrefix      = "xml"
)

// outAttr is an attribute about to be emitted: everything the comparator and the writer need,
// and nothing else. Namespace declarations reach here either from an attribute node or
// synthesized from a (prefix, uri) pair in the symbol table; the two forms are identical by
// construction, because a declaration attribute's QName is always "xmlns" or "xmlns:"+prefix
// and its value is always the URI.
type outAttr struct {
	space string // namespace URI; "" for an attribute in no namespace
	local string // local name; "xmlns" for the default declaration
	qname string // the QName as it must be emitted
	value string
}

// nsAttr builds the declaration attribute for a (prefix, uri) mapping.
func nsAttr(prefix, uri string) outAttr {
	a := outAttr{space: xmlnsNamespace, local: prefix, value: uri, qname: xmlnsPrefix}
	if prefix != xmlnsPrefix {
		a.qname = xmlnsPrefix + ":" + prefix
	}
	return a
}

// attrLess reports whether a sorts before b, and attrCompare is the three-way form.
//
//   - namespace declarations sort before every other attribute;
//   - among declarations, by local name, with "xmlns" (the default declaration) mapped to ""
//     so that it sorts first;
//   - among the rest, no-namespace attributes first, compared by QName; then by namespace URI
//     and, on a tie, by local name.
func attrLess(a, b outAttr) bool { return attrCompare(a, b) < 0 }

func attrCompare(a, b outAttr) int {
	isNSa := a.space == xmlnsNamespace
	isNSb := b.space == xmlnsNamespace
	if isNSa {
		if !isNSb {
			return -1
		}
		la, lb := a.local, b.local
		if la == xmlnsPrefix {
			la = ""
		}
		if lb == xmlnsPrefix {
			lb = ""
		}
		return compareUTF16(la, lb)
	}
	if isNSb {
		return 1
	}
	// Java reports namespaceURI == null for an attribute in no namespace and sorts null before
	// any non-null URI; "" is lexicographically least among all strings, so the two agree and
	// no "has namespace" flag is needed.
	if a.space == "" {
		if b.space == "" {
			return compareUTF16(a.qname, b.qname)
		}
		return -1
	}
	if b.space == "" {
		return 1
	}
	if c := compareUTF16(a.space, b.space); c != 0 {
		return c
	}
	return compareUTF16(a.local, b.local)
}

// compareUTF16 orders two strings the way java.lang.String.compareTo does, by UTF-16 code
// unit. This is NOT Go's byte-wise string order: a code point above U+FFFF is a surrogate pair
// beginning at U+D800 in UTF-16 and so sorts BELOW U+E000..U+FFFD, while in UTF-8 it sorts
// above them. The corpus document sort-astral.xml straddles the boundary and Santuario emits
//
//	<r xmlns:a="urn:xU+FDF0" xmlns:b="urn:xU+10000" Id="r" b:z="2" a:z="1"></r>
//
// - b:z before a:z - which byte-wise ordering gets backwards. Namespace URIs and NCNames in
// that range do not occur in practice, but the KAT does, so the comparator is Java's.
func compareUTF16(a, b string) int {
	// Skip the common ASCII prefix: over ASCII the two orders agree, and almost every real
	// comparison is decided here. i stays on a rune boundary in both strings.
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
	ua := utf16.Encode([]rune(a[i:]))
	ub := utf16.Encode([]rune(b[i:]))
	for k := 0; k < len(ua) && k < len(ub); k++ {
		if ua[k] != ub[k] {
			return int(ua[k]) - int(ub[k])
		}
	}
	return len(ua) - len(ub)
}

// attrSet is the sorted attribute set the emitters build, a stand-in for Santuario's
// TreeSet<Attr>(COMPARE). TreeSet semantics matter twice: it sorts by the comparator alone,
// and add() on an element that compares equal to one already present keeps the first and
// discards the second.
type attrSet struct {
	items []outAttr
}

func (s *attrSet) add(a outAttr) {
	i := sort.Search(len(s.items), func(i int) bool { return attrCompare(s.items[i], a) >= 0 })
	if i < len(s.items) && attrCompare(s.items[i], a) == 0 {
		return // TreeSet.add: an equal element is already present.
	}
	s.items = append(s.items, outAttr{})
	copy(s.items[i+1:], s.items[i:])
	s.items[i] = a
}

func (s *attrSet) len() int { return len(s.items) }

// writeTo emits the set, which is CanonicalizerBase.outputAttrToWriter over the sorted result.
func (s *attrSet) writeTo(w *bufio.Writer) {
	for _, a := range s.items {
		w.WriteByte(' ')
		w.WriteString(a.qname)
		w.WriteString(`="`)
		writeAttrValueEscaped(w, a.value)
		w.WriteByte('"')
	}
}
