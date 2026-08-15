// Ported from the generated JAXB classes:
//   - AnyType.java
//   - VOReferenceType.java
//   - NsPrefixMappingType.java
//   - TypedDataType.java
//   - AdditionalValidationReportDataType.java
//
// (specs-validation-report, DSS 6.5.RC1).
//
// # Capturing an opaque element's own attributes
//
// AnyType and RawContent (below) round-trip an extensibility element's
// content via innerxml, but a real report can also put attributes - and the
// namespace declarations those attributes need - directly on the element
// itself: TypedDataType.Value commonly carries `xsi:type="xs:string"` with
// the `xmlns:xs`/`xmlns:xsi` declarations that make it resolvable. innerxml
// alone does not see those (they are outside the captured region), so both
// types also capture `,any,attr` and replay it on Marshal through
// flattenAttrs - see that function's doc comment for why a custom MarshalXML
// is needed to do so byte-exactly.
package jaxb

import "encoding/xml"

// AnyType is the Go form of the generated JAXB class AnyType (complexType
// AnyType, mixed="true" with a lax xs:any): used throughout this schema for
// every "OtherInformation" extensibility element. Content and Attrs are
// captured verbatim rather than interpreted - see doc.go's "xs:anyType /
// xs:any content" section and this file's header - so unmarshal->marshal
// reproduces whatever child elements/text/attributes a real report placed
// there byte-for-byte without this package needing to model them.
type AnyType struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element's own attributes (flattened back to literal
// prefixes, see flattenAttrs) followed by its raw captured content.
func (a *AnyType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, flattenAttrs(a.Attrs)...)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: a.Content}, start)
}

// RawContent is the Go stand-in for a JAXB property whose Java type is the
// bare `Object` (VOReferenceType.Any, TypedDataType.Value,
// IndividualValidationConstraintReportType.Indications,
// ValidationObjectRepresentationType's "direct" choice member): an
// xs:anyType/lax xs:any content model with no XSD complex type of its own.
// It captures whatever XML content and attributes were present without
// interpreting them, the same verbatim round-trip AnyType provides for its
// own, differently-named, use of the pattern.
type RawContent struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element's own attributes (flattened back to literal
// prefixes, see flattenAttrs) followed by its raw captured content.
func (r *RawContent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, flattenAttrs(r.Attrs)...)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: r.Content}, start)
}

// flattenAttrs rewrites the namespace-resolved attribute names the decoder
// produces (Name.Space holding a resolved URI) back into the literal
// "prefix:local" token JAXB wrote, with an empty Name.Space so
// encoding/xml's encoder writes it as-is instead of inventing its own
// prefix (encoding/xml always allocates a fresh prefix - or a fresh local
// `xmlns="..."` - for any attribute whose Name.Space is a real namespace
// URI; it never reuses one declared by a hand-written ancestor start tag,
// which is exactly the shape jaxbRootNamespaces produces at the document
// root - see xml.go). The prefix used for a real attribute is not itself
// captured by the decoder, only its resolved namespace URI is, so it is
// recovered from this same element's own xmlns:* pseudo-attributes
// (captured alongside it, unchanged) rather than from any document-wide
// table.
func flattenAttrs(attrs []xml.Attr) []xml.Attr {
	if len(attrs) == 0 {
		return nil
	}
	prefixOf := map[string]string{}
	for _, a := range attrs {
		if a.Name.Space == "xmlns" {
			prefixOf[a.Value] = a.Name.Local
		}
	}
	out := make([]xml.Attr, len(attrs))
	for i, a := range attrs {
		switch {
		case a.Name.Space == "xmlns":
			out[i] = xml.Attr{Name: xml.Name{Local: "xmlns:" + a.Name.Local}, Value: a.Value}
		case a.Name.Space == "" && a.Name.Local == "xmlns":
			out[i] = a
		case a.Name.Space != "":
			if p, ok := prefixOf[a.Name.Space]; ok {
				out[i] = xml.Attr{Name: xml.Name{Local: p + ":" + a.Name.Local}, Value: a.Value}
			} else {
				out[i] = a
			}
		default:
			out[i] = a
		}
	}
	return out
}

// VOReferenceType is the Go form of the generated JAXB class VOReferenceType
// (complexType VOReferenceType): a reference to one or more validation
// objects by xs:ID, optionally carrying an inline lax xs:any payload. No
// object-graph resolution of VOReference is performed by this package - see
// xml.go's header.
type VOReferenceType struct {
	Any         string `xml:",innerxml"`
	VOReference IDREFS `xml:"VOReference,attr"`
}

// NsPrefixMappingType is the Go form of the generated JAXB class
// NsPrefixMappingType (complexType NsPrefixMappingType).
type NsPrefixMappingType struct {
	NamespaceURI    string `xml:"NamespaceURI"`
	NamespacePrefix string `xml:"NamespacePrefix"`
}

// TypedDataType is the Go form of the generated JAXB class TypedDataType
// (complexType TypedDataType).
type TypedDataType struct {
	Type  string     `xml:"Type"`
	Value RawContent `xml:"Value"`
}

// AdditionalValidationReportDataType is the Go form of the generated JAXB
// class AdditionalValidationReportDataType (complexType
// AdditionalValidationReportDataType).
type AdditionalValidationReportDataType struct {
	ReportData []*TypedDataType `xml:"ReportData"`
}
