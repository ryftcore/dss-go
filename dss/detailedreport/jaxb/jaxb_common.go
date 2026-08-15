// Ported from DetailedReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.detailedreport.jaxb. Per the phase-8a
// generated-JAXB rule the generated classes are grouped into schema-area files
// rather than one file per class; every Java class keeps its name and its
// exact field order so that encoding/xml reproduces the JAXB element sequence
// byte for byte. This file holds the leaf types shared across the schema:
// messages, conclusions, constraints, the two simple-type enums, and the
// cryptographic/revocation value types.

package jaxb

import "fmt"

// XmlStatus is the Go form of the generated JAXB enum XmlStatus (simpleType Status).
type XmlStatus string

// The Status values, with the exact lexical form of the schema.
const (
	XmlStatus_OK          XmlStatus = "OK"
	XmlStatus_NOT_OK      XmlStatus = "NOT OK"
	XmlStatus_IGNORED     XmlStatus = "IGNORED"
	XmlStatus_INFORMATION XmlStatus = "INFORMATION"
	XmlStatus_WARNING     XmlStatus = "WARNING"
)

// MarshalText writes the enum value, as the generated value() method does.
func (v XmlStatus) MarshalText() ([]byte, error) { return []byte(v), nil }

// UnmarshalText resolves the lexical form, mirroring fromValue().
func (v *XmlStatus) UnmarshalText(text []byte) error {
	switch s := XmlStatus(text); s {
	case XmlStatus_OK, XmlStatus_NOT_OK, XmlStatus_IGNORED, XmlStatus_INFORMATION, XmlStatus_WARNING:
		*v = s
		return nil
	default:
		return fmt.Errorf("no enum constant XmlStatus.%s", text)
	}
}

// XmlBlockType is the Go form of the generated JAXB enum XmlBlockType (simpleType BlockType).
type XmlBlockType string

// The BlockType values, with the exact lexical form of the schema. The schema
// enumerates "ER" twice; the duplicate collapses to one Go constant.
const (
	XmlBlockType_SIG_BBB     XmlBlockType = "SIG_BBB"
	XmlBlockType_REV_BBB     XmlBlockType = "REV_BBB"
	XmlBlockType_TST_BBB     XmlBlockType = "TST_BBB"
	XmlBlockType_CNT_TST_BBB XmlBlockType = "CNT_TST_BBB"
	XmlBlockType_AOV         XmlBlockType = "AOV"
	XmlBlockType_AOV_XCV     XmlBlockType = "AOV_XCV"
	XmlBlockType_CRS         XmlBlockType = "CRS"
	XmlBlockType_PSV_CRS     XmlBlockType = "PSV_CRS"
	XmlBlockType_RAC         XmlBlockType = "RAC"
	XmlBlockType_RAC_SUB_XCV XmlBlockType = "RAC_SUB_XCV"
	XmlBlockType_RFC         XmlBlockType = "RFC"
	XmlBlockType_SUB_XCV     XmlBlockType = "SUB_XCV"
	XmlBlockType_SUB_XCV_TA  XmlBlockType = "SUB_XCV_TA"
	XmlBlockType_REV_CC      XmlBlockType = "REV_CC"
	XmlBlockType_ER          XmlBlockType = "ER"
	XmlBlockType_PSV         XmlBlockType = "PSV"
	XmlBlockType_PCV         XmlBlockType = "PCV"
	XmlBlockType_VTS         XmlBlockType = "VTS"
	XmlBlockType_TST_PSV     XmlBlockType = "TST_PSV"
	XmlBlockType_TST         XmlBlockType = "TST"
	XmlBlockType_LTV         XmlBlockType = "LTV"
	XmlBlockType_LTV_SUB_XCV XmlBlockType = "LTV_SUB_XCV"
	XmlBlockType_LTA         XmlBlockType = "LTA"
)

// blockTypeValues lists every constant once, for UnmarshalText validation.
var blockTypeValues = []XmlBlockType{
	XmlBlockType_SIG_BBB, XmlBlockType_REV_BBB, XmlBlockType_TST_BBB, XmlBlockType_CNT_TST_BBB,
	XmlBlockType_AOV, XmlBlockType_AOV_XCV, XmlBlockType_CRS, XmlBlockType_PSV_CRS,
	XmlBlockType_RAC, XmlBlockType_RAC_SUB_XCV, XmlBlockType_RFC, XmlBlockType_SUB_XCV,
	XmlBlockType_SUB_XCV_TA, XmlBlockType_REV_CC, XmlBlockType_ER, XmlBlockType_PSV,
	XmlBlockType_PCV, XmlBlockType_VTS, XmlBlockType_TST_PSV, XmlBlockType_TST,
	XmlBlockType_LTV, XmlBlockType_LTV_SUB_XCV, XmlBlockType_LTA,
}

// MarshalText writes the enum value, as the generated value() method does.
func (v XmlBlockType) MarshalText() ([]byte, error) { return []byte(v), nil }

// UnmarshalText resolves the lexical form, mirroring fromValue().
func (v *XmlBlockType) UnmarshalText(text []byte) error {
	s := XmlBlockType(text)
	for _, c := range blockTypeValues {
		if c == s {
			*v = s
			return nil
		}
	}
	return fmt.Errorf("no enum constant XmlBlockType.%s", text)
}

// XmlMessage is the Go form of the generated JAXB class XmlMessage
// (complexType Message, simpleContent extending xs:string).
type XmlMessage struct {
	Value string  `xml:",chardata"`
	Key   *string `xml:"Key,attr,omitempty"`
}

// XmlConclusion is the Go form of the generated JAXB class XmlConclusion
// (complexType Conclusion).
type XmlConclusion struct {
	Indication    IndicationValue     `xml:"Indication"`
	SubIndication *SubIndicationValue `xml:"SubIndication,omitempty"`
	Errors        []*XmlMessage       `xml:"Errors,omitempty"`
	Warnings      []*XmlMessage       `xml:"Warnings,omitempty"`
	Infos         []*XmlMessage       `xml:"Infos,omitempty"`
}

// XmlConstraint is the Go form of the generated JAXB class XmlConstraint
// (complexType Constraint).
type XmlConstraint struct {
	Name           *XmlMessage   `xml:"Name"`
	Status         XmlStatus     `xml:"Status"`
	Error          *XmlMessage   `xml:"Error,omitempty"`
	Warning        *XmlMessage   `xml:"Warning,omitempty"`
	Info           *XmlMessage   `xml:"Info,omitempty"`
	AdditionalInfo *string       `xml:"AdditionalInfo,omitempty"`
	Id             *string       `xml:"Id,attr,omitempty"`
	BlockType      *XmlBlockType `xml:"BlockType,attr,omitempty"`
}

// XmlRevocationInformation is the Go form of the generated JAXB class
// XmlRevocationInformation (complexType RevocationInformation).
type XmlRevocationInformation struct {
	CertificateId  string                 `xml:"CertificateId"`
	RevocationId   string                 `xml:"RevocationId"`
	Reason         *RevocationReasonValue `xml:"Reason,omitempty"`
	RevocationDate XSDateTime             `xml:"RevocationDate"`
}

// XmlCryptographicAlgorithm is the Go form of the generated JAXB class
// XmlCryptographicAlgorithm (complexType CryptographicAlgorithm).
type XmlCryptographicAlgorithm struct {
	Name      string  `xml:"Name"`
	Uri       string  `xml:"Uri"`
	KeyLength *string `xml:"KeyLength,omitempty"`
}

// XmlCryptographicValidation is the Go form of the generated JAXB class
// XmlCryptographicValidation (complexType CryptographicValidation).
type XmlCryptographicValidation struct {
	Algorithm                    *XmlCryptographicAlgorithm `xml:"Algorithm"`
	NotAfter                     *XSDateTime                `xml:"NotAfter,omitempty"`
	ConcernedMaterialDescription *string                    `xml:"ConcernedMaterialDescription,omitempty"`
	Conclusion                   *XmlConclusion             `xml:"Conclusion"`
	TokenId                      *string                    `xml:"TokenId,attr,omitempty"`
}

// XmlSemantic is the Go form of the generated JAXB class XmlSemantic
// (complexType Semantic, simpleContent extending xs:string).
type XmlSemantic struct {
	Value string `xml:",chardata"`
	Key   string `xml:"Key,attr"`
}
