// Ported from ts_119612v020401_xsd.xsd (DSS 6.5.RC1) via the generated JAXB
// class TrustStatusListType.java (eu.europa.esig.trustedlist.jaxb.tsl), plus
// the marshalling behaviour of eu.europa.esig.trustedlist.TrustedListFacade
// and eu.europa.esig.trustedlist.mra.MRAFacade (dss-jaxb-common's
// AbstractJaxbFacade with JAXB_FORMATTED_OUTPUT=true).
//
// # Root namespace declarations
//
// The JAXB reference implementation pre-declares, once at the document
// element, every namespace known to the JAXBContext the facade builds -
// see eu.europa.esig.trustedlist.TrustedListUtils.getJAXBContext()/
// eu.europa.esig.trustedlist.mra.MRAUtils.getJAXBContext() - with a "nsN"
// prefix assigned by the RI's own internal ordering, REGARDLESS of whether
// that namespace's elements actually appear in a given document (the same
// behaviour dss/validationreport/jaxb/jaxb_root.go documents for
// ValidationReportType). This was verified empirically against this
// package's own oracle (a small Java program driving TrustedListFacade/MRAFacade
// directly, run over the real fixtures this package's KATs use - see
// xml_kat_test.go): TrustedListFacade.marshall always writes
//
//	ns2 -> XMLDSig, ns3 -> XAdES 1.3.2, ns4 -> tslx, ns5 -> ecc, ns6 -> XAdES 1.4.1
//
// and MRAFacade.marshall always writes
//
//	ns2 -> XMLDSig, ns3 -> XAdES 1.3.2, ns4 -> tslx, ns5 -> ecc, ns6 -> mra, ns7 -> XAdES 1.4.1
//
// with the tsl namespace itself as the unprefixed default (the document's
// own elementFormDefault="qualified" namespace). Neither set depends on the
// input document's own xmlns:* spelling or order (real eu-lotl/country TL
// files commonly declare ns3/ns4 the other way around, or use "no" instead
// of "yes" for the XML declaration's standalone attribute - none of that
// survives a JAXB round trip), so plainRootNamespaces/mraRootNamespaces are
// simply fixed constants rather than a captured-and-replayed set: unlike
// ValidationReportType, which round-trips reports it did not itself
// produce, this package always reproduces the FACADE's canonical bytes for
// marshal-parity, never the original download's.
//
// Marshal (TrustedListFacade's set) and MarshalMRA (MRAFacade's) are
// otherwise identical; wildcardElements (jaxb_common.go) dispatches
// mra:MutualRecognitionAgreementInformation regardless of which one is
// used, since the two facades' JAXBContext otherwise differ only in whether
// mra.ObjectFactory is registered - content is what determines whether it
// legitimately appears, not the facade.
//
// # Known deviation: JAXB RI indentation desync around list/xs:any content
//
// eclipse-ee4j jaxb-ri 3.0.2's indenting output writer tracks nesting depth
// with a counter that this package's own oracle (see xml_kat_test.go) shows
// getting thrown off, for a stretch of following elements whose length
// varies per document, by array-valued properties (List<T>, i.e. every
// repeatable element) and even more so by @XmlAnyElement/@XmlMixed content
// (AnyType) - both go through a different internal serialization path
// (ArrayERProperty/ArrayElementNodeProperty) than a plain single-valued
// property. The result is not a clean, position-independent rule (the SAME
// class, eu.europa.esig.trustedlist.jaxb.tsl.InternationalNamesType, comes
// out correctly indented at one position in a real fixture and
// flush-left - depth reset towards zero - at another, later one in the
// SAME document, once something earlier has "poisoned" the RI's counter),
// so it is a genuine implementation quirk of that JAXB RI version, not a
// documented or portable behaviour. Reproducing it bit-for-bit would mean
// reverse-engineering jaxb-ri's ArrayERProperty/IndentingUTF8XmlOutput
// depth bookkeeping rather than porting DSS - out of proportion to what it
// buys, since it is PURE, semantically-insignificant whitespace: the
// information content (element/attribute names, order, text, structure) is
// unaffected. Marshal/MarshalMRA therefore always produce clean,
// consistently depth-indented output instead of reproducing the RI's own
// desync, and xml_kat_test.go's marshal-parity KAT compares the two
// canonically (ignoring insignificant whitespace) rather than byte-for-byte
// for this reason.
//
// # xs:dateTime properties
//
// NextUpdateType.DateTime, TSLSchemeInformationType.ListIssueDateTime,
// (TSP)ServiceInformationType/ServiceHistoryInstanceType.StatusStartingTime
// and mra's TrustServiceEquivalenceStatusStartingTime are all plain
// XMLGregorianCalendar/@XmlSchemaType(name="dateTime") properties - the
// JAXB built-in xs:dateTime binding, not routed through a custom
// XmlAdapter the way dss/diagnostic/jaxb's DiagnosticData.xsd is (see that
// package's XSDateTime). javax.xml.datatype.XMLGregorianCalendar's
// lexical round trip through DatatypeFactory is loss-and-reformat-free: it
// preserves the exact digits (fractional seconds, timezone offset or its
// absence) it was parsed from. Reproducing that losslessly needs nothing
// more than a plain Go string, so every one of these properties is bound
// `*string`: the raw xs:dateTime lexical form, unparsed and
// unreformatted. Downstream time semantics (parsing into time.Time,
// comparing against now for cache/expiry decisions) belong to dss/tsl's
// parsing tasks, not this structural model.
package jaxb

import (
	"bytes"
	"encoding/xml"
	"strings"
)

// xmlDeclaration is the declaration the RI's marshaller emits - see
// dss/diagnostic/jaxb/xml.go's identical constant.
const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// plainRootNamespaces is TrustedListFacade's fixed root xmlns:* set - see
// this file's header.
var plainRootNamespaces = []xml.Attr{
	{Name: xml.Name{Local: "ns2"}, Value: NamespaceDSig},
	{Name: xml.Name{Local: "ns3"}, Value: NamespaceXAdES132},
	{Name: xml.Name{Local: "ns4"}, Value: NamespaceTSLX},
	{Name: xml.Name{Local: "ns5"}, Value: NamespaceECC},
	{Name: xml.Name{Local: "ns6"}, Value: NamespaceXAdES141},
}

// mraRootNamespaces is MRAFacade's fixed root xmlns:* set - see this file's
// header.
var mraRootNamespaces = []xml.Attr{
	{Name: xml.Name{Local: "ns2"}, Value: NamespaceDSig},
	{Name: xml.Name{Local: "ns3"}, Value: NamespaceXAdES132},
	{Name: xml.Name{Local: "ns4"}, Value: NamespaceTSLX},
	{Name: xml.Name{Local: "ns5"}, Value: NamespaceECC},
	{Name: xml.Name{Local: "ns6"}, Value: NamespaceMRA},
	{Name: xml.Name{Local: "ns7"}, Value: NamespaceXAdES141},
}

// TrustStatusListType is the Go form of the generated JAXB class
// TrustStatusListType (complexType TrustStatusListType, the document
// element "TrustServiceStatusList"). Signature is
// eu.europa.esig.xmldsig.jaxb.SignatureType (ds:Signature), outside this
// manifest - see jaxb_common.go's header - captured verbatim; nothing in
// this package's scope validates it (TLValidatorTask verifies the ORIGINAL
// document bytes through the frozen xades/xmldsig validator, never this
// unmarshalled tree - see jaxb_common.go's header).
type TrustStatusListType struct {
	SchemeInformation        *TSLSchemeInformationType     `xml:"SchemeInformation"`
	TrustServiceProviderList *TrustServiceProviderListType `xml:"TrustServiceProviderList,omitempty"`
	Signature                *dsigSignature                `xml:"http://www.w3.org/2000/09/xmldsig# Signature,omitempty"`
	TSLTag                   string                        `xml:"TSLTag,attr"`
	Id                       *string                       `xml:"Id,attr,omitempty"`
}

// Unmarshal parses a trusted-list document, the way TrustedListFacade's
// unmarshall (used by both TrustedListFacade and MRAFacade - see
// dss/trustedlist) does.
func Unmarshal(data []byte) (*TrustStatusListType, error) {
	tsl := &TrustStatusListType{}
	if err := xml.Unmarshal(data, tsl); err != nil {
		return nil, err
	}
	return tsl, nil
}

// Marshal writes a trusted-list document byte-for-byte the way
// TrustedListFacade.marshall does: the XML declaration, four-space indented
// output, a trailing newline, plainRootNamespaces at the root, and the JAXB
// spellings jaxbCanonical restores.
func Marshal(tsl *TrustStatusListType) ([]byte, error) {
	return marshal(tsl, plainRootNamespaces)
}

// MarshalMRA writes a trusted-list document the way MRAFacade.marshall
// does: identical to Marshal but with mraRootNamespaces (adding mra's own
// prefix) at the root - see this file's header.
func MarshalMRA(tsl *TrustStatusListType) ([]byte, error) {
	return marshal(tsl, mraRootNamespaces)
}

func marshal(tsl *TrustStatusListType, extraNamespaces []xml.Attr) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "    ")
	start := xml.StartElement{
		Name: xml.Name{Local: "TrustServiceStatusList"},
		Attr: make([]xml.Attr, 0, 1+len(extraNamespaces)),
	}
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: NamespaceTSL})
	for _, a := range extraNamespaces {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:" + a.Name.Local}, Value: a.Value})
	}
	if err := enc.EncodeElement(tsl, start); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	body := jaxbCanonical(buf.Bytes())
	out := make([]byte, 0, len(xmlDeclaration)+len(body)+1)
	out = append(out, xmlDeclaration...)
	out = append(out, body...)
	out = append(out, '\n')
	return out, nil
}

// ------------------------------------------------------------- canonicalising

// jaxbCanonical rewrites encoding/xml output into the spelling the JAXB RI
// produces - the same two syntax-only normalisations as
// dss/diagnostic/jaxb/xml.go's jaxbCanonical (self-closing empty tags per
// jaxb_model.go's content model; " and ' left alone in character data,
// numeric references lower-cased/relocated per JAXB's escaper), reproduced
// here rather than shared because each jaxb package is self-contained (see
// PORTING.md's layout rules).
func jaxbCanonical(in []byte) []byte {
	var out bytes.Buffer
	out.Grow(len(in))
	var stack []string
	for i := 0; i < len(in); {
		if in[i] != '<' {
			j := bytes.IndexByte(in[i:], '<')
			if j < 0 {
				j = len(in) - i
			}
			writeCharData(&out, in[i:i+j])
			i += j
			continue
		}
		end := tagEnd(in, i)
		tag := in[i:end]
		if len(tag) > 1 && tag[1] == '/' {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			out.Write(tag)
			i = end
			continue
		}
		name := tagName(tag)
		closing := append(append([]byte("</"), name...), '>')
		if bytes.HasPrefix(in[end:], closing) && !carriesCharData(name, stack) {
			writeStartTag(&out, tag[:len(tag)-1])
			out.WriteString("/>")
			i = end + len(closing)
			continue
		}
		writeStartTag(&out, tag)
		stack = append(stack, name)
		i = end
	}
	return out.Bytes()
}

// tagEnd returns the index just past the '>' closing the tag starting at i.
// encoding/xml escapes '<' and '>' inside attribute values, so no
// quoting-aware scan is needed.
func tagEnd(in []byte, i int) int {
	j := bytes.IndexByte(in[i:], '>')
	if j < 0 {
		return len(in)
	}
	return i + j + 1
}

func tagName(tag []byte) string {
	s := tag[1:]
	for k := 0; k < len(s); k++ {
		switch s[k] {
		case ' ', '\t', '\n', '\r', '>', '/':
			return string(s[:k])
		}
	}
	return string(s)
}

// writeStartTag copies a start tag, rewriting the escapes inside attribute
// values to the RI's spelling.
func writeStartTag(out *bytes.Buffer, tag []byte) {
	inValue := false
	start := 0
	for k := 0; k < len(tag); k++ {
		if tag[k] != '"' {
			continue
		}
		if inValue {
			writeAttrValue(out, tag[start:k])
		} else {
			out.Write(tag[start:k])
		}
		out.WriteByte('"')
		inValue = !inValue
		start = k + 1
	}
	if start < len(tag) {
		out.Write(tag[start:])
	}
}

// charDataEscapes maps the numeric references encoding/xml emits in
// character data to the spelling the RI uses there.
var charDataEscapes = []struct{ from, to string }{
	{"&#34;", `"`},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "\n"},
	{"&#xD;", "&#13;"},
}

// attrValueEscapes maps the same references to the spelling the RI uses
// inside attribute values, where a double quote must stay escaped.
var attrValueEscapes = []struct{ from, to string }{
	{"&#34;", "&quot;"},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "&#10;"},
	{"&#xD;", "&#13;"},
}

func writeCharData(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, charDataEscapes)
}

func writeAttrValue(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, attrValueEscapes)
}

func writeRewritten(out *bytes.Buffer, b []byte, table []struct{ from, to string }) {
	if bytes.IndexByte(b, '&') < 0 {
		out.Write(b)
		return
	}
	s := string(b)
	for _, e := range table {
		if e.from != e.to {
			s = strings.ReplaceAll(s, e.from, e.to)
		}
	}
	out.WriteString(s)
}
