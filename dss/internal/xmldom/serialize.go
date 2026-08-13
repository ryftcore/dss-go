package xmldom

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// SerializeOptions tunes Serialize. A nil *SerializeOptions selects the defaults.
type SerializeOptions struct {
	XMLDeclaration bool   // default true via nil options
	Encoding       string // default "UTF-8"
}

var defaultSerializeOptions = SerializeOptions{XMLDeclaration: true, Encoding: "UTF-8"}

// Serialize writes n as XML. It is NOT canonicalization: attributes keep document
// order and no namespace inheritance is performed. Use internal/xmlc14n for anything
// that is hashed or signed.
//
// The output is always UTF-8. Encoding names the encoding in the XML declaration only;
// it is deliberately not used to transcode, because emitting non-UTF-8 bytes would
// make Parse(Serialize(d)) lossy and this serializer's whole job is round-tripping a
// tree after DOM surgery. Byte-parity with Java's Transformer is not a requirement:
// DomUtils.serializeNode output is never itself signed, only c14n output is.
func (n *Node) Serialize(w io.Writer, opts *SerializeOptions) error {
	o := defaultSerializeOptions
	if opts != nil {
		o = *opts
		if o.Encoding == "" {
			o.Encoding = defaultSerializeOptions.Encoding
		}
	}
	bw := bufio.NewWriter(w)
	s := &serializer{w: bw}
	if o.XMLDeclaration && n != nil && n.Kind == Document {
		s.str(`<?xml version="1.0" encoding="`)
		s.attrText(o.Encoding)
		s.str(`"?>`)
	}
	s.node(n)
	if s.err != nil {
		return s.err
	}
	return bw.Flush()
}

// Bytes is Serialize into a buffer.
func (n *Node) Bytes(opts *SerializeOptions) ([]byte, error) {
	var b bytes.Buffer
	if err := n.Serialize(&b, opts); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type serializer struct {
	w   *bufio.Writer
	err error
}

func (s *serializer) str(v string) {
	if s.err != nil {
		return
	}
	_, s.err = s.w.WriteString(v)
}

func (s *serializer) node(n *Node) {
	if n == nil || s.err != nil {
		return
	}
	switch n.Kind {
	case Document:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.node(c)
		}
	case Element:
		s.str("<")
		s.str(n.Name.QName())
		for _, a := range n.Attrs {
			s.str(" ")
			s.str(a.Name.QName())
			s.str(`="`)
			s.attrText(a.Value)
			s.str(`"`)
		}
		s.str(">")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.node(c)
		}
		// Empty elements are written as a start tag followed by an end tag, so that
		// <r/> and <r></r> serialize identically, exactly as c14n does.
		s.str("</")
		s.str(n.Name.QName())
		s.str(">")
	case Attribute:
		s.str(n.Name.QName())
		s.str(`="`)
		s.attrText(n.Value)
		s.str(`"`)
	case Text:
		s.text(n.Value)
	case CDATA:
		s.cdata(n.Value)
	case Comment:
		// Comment data is written verbatim. Parse rejects "--" inside a comment, so
		// the only way to reach an unserializable comment is to build one by hand.
		s.str("<!--")
		s.str(n.Value)
		s.str("-->")
	case ProcInst:
		s.str("<?")
		s.str(n.Name.Local)
		if n.Value != "" {
			s.str(" ")
			s.str(n.Value)
		}
		s.str("?>")
	}
}

// cdata writes a CDATA section, splitting on an embedded "]]>" into two sections.
func (s *serializer) cdata(v string) {
	s.str("<![CDATA[")
	s.str(strings.ReplaceAll(v, "]]>", "]]]]><![CDATA[>"))
	s.str("]]>")
}

// text escapes character data.
//
// The design note lists & < > for this context. CR is escaped as well, and must be:
// XML clause 2.11 turns a literal CR back into LF on the next parse, so a tree
// carrying a genuine CR - which can only have come from &#xD; - would silently change
// on a round trip. Java's Transformer escapes it here for the same reason.
func (s *serializer) text(v string) {
	if !strings.ContainsAny(v, "&<>\r") {
		s.str(v)
		return
	}
	for i := 0; i < len(v); i++ {
		switch c := v[i]; c {
		case '&':
			s.str("&amp;")
		case '<':
			s.str("&lt;")
		case '>':
			s.str("&gt;")
		case '\r':
			s.str("&#xD;")
		default:
			if s.err == nil {
				s.err = s.w.WriteByte(c)
			}
		}
	}
}

// attrText escapes an attribute value.
//
// The design note lists & < " for this context. TAB, LF and CR are escaped too, for
// the same round-trip reason as in text: unescaped, XML clause 3.3.3 attribute-value
// normalization would turn each of them into a space on the next parse and the tree
// would not survive Serialize followed by Parse.
func (s *serializer) attrText(v string) {
	if !strings.ContainsAny(v, "&<\"\t\n\r") {
		s.str(v)
		return
	}
	for i := 0; i < len(v); i++ {
		switch c := v[i]; c {
		case '&':
			s.str("&amp;")
		case '<':
			s.str("&lt;")
		case '"':
			s.str("&quot;")
		case '\t':
			s.str("&#x9;")
		case '\n':
			s.str("&#xA;")
		case '\r':
			s.str("&#xD;")
		default:
			if s.err == nil {
				s.err = s.w.WriteByte(c)
			}
		}
	}
}
