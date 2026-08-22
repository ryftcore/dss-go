// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

// XmlDigestAlgoAndValueContent carries the element content of the XmlDigestAlgoAndValue base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlDigestAlgoAndValueContent struct {
	DigestMethod *DigestAlgorithmValue `xml:"DigestMethod,omitempty"`
	DigestValue  *Base64Binary         `xml:"DigestValue,omitempty"`
}

// XmlDigestAlgoAndValueAttrs carries the attributes of the XmlDigestAlgoAndValue base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlDigestAlgoAndValueAttrs struct {
	Match *bool `xml:"match,attr,omitempty"`
}

// XmlDigestAlgoAndValue is the Go form of the generated JAXB class XmlDigestAlgoAndValue
// (complexType DigestAlgoAndValue).
type XmlDigestAlgoAndValue struct {
	XmlDigestAlgoAndValueContent
	XmlDigestAlgoAndValueAttrs
}

// XmlDigestAlgoAndValues is the Go form of the generated JAXB class XmlDigestAlgoAndValues
// (complexType DigestAlgoAndValues).
type XmlDigestAlgoAndValues struct {
	DigestAlgoAndValue []*XmlDigestAlgoAndValue `xml:"DigestAlgoAndValue"`
}

// XmlDigestMatcher is the Go form of the generated JAXB class XmlDigestMatcher
// (complexType DigestMatcher).
type XmlDigestMatcher struct {
	XmlDigestAlgoAndValueContent
	DataFound            bool                         `xml:"DataFound"`
	DataIntact           bool                         `xml:"DataIntact"`
	DataObjectReferences *DataObjectReferencesWrapper `xml:"DataObjectReferences"`
	DisclosableClaim     *XmlDisclosableClaim         `xml:"DisclosableClaim,omitempty"`
	Type                 *DigestMatcherTypeValue      `xml:"type,attr,omitempty"`
	Id                   *string                      `xml:"id,attr,omitempty"`
	Uri                  *string                      `xml:"uri,attr,omitempty"`
	DocumentName         *string                      `xml:"documentName,attr,omitempty"`
	Duplicated           *bool                        `xml:"duplicated,attr,omitempty"`
	XmlDigestAlgoAndValueAttrs
}
