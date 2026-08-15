// Ported from DetailedReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.detailedreport.jaxb. See jaxb_common.go's header for
// the file-grouping rule and jaxb_process.go's header for the
// Content/Attrs extension-embedding pattern. This file holds the eIDAS
// qualification blocks: certificate/signature/timestamp/EAA/QWAC
// qualification, the TS 119 602 certificate approval status blocks, and the
// list-of-trusted-entities analysis blocks (TLAnalysis/LoTEAnalysis).

package jaxb

import (
	"encoding/xml"

	"github.com/utain/esig/dss/enumerations"
)

// XmlCertificateQualificationProcess is the Go form of the generated JAXB
// class XmlCertificateQualificationProcess
// (complexType CertificateQualificationProcess, extends ConstraintsConclusion).
type XmlCertificateQualificationProcess struct {
	XmlConstraintsConclusionContent
	ValidationCertificateQualification []*XmlValidationCertificateQualification `xml:"ValidationCertificateQualification,omitempty"`
	Id                                 *string                                  `xml:"Id,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlCertificateApprovalStatusProcess is the Go form of the generated JAXB
// class XmlCertificateApprovalStatusProcess
// (complexType CertificateApprovalStatusProcess, extends ConstraintsConclusion).
type XmlCertificateApprovalStatusProcess struct {
	XmlConstraintsConclusionContent
	ValidationCertificateApprovalStatus []*XmlValidationCertificateApprovalStatus `xml:"ValidationCertificateApprovalStatus,omitempty"`
	Id                                  *string                                   `xml:"Id,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationCertificateQualification is the Go form of the generated JAXB
// class XmlValidationCertificateQualification
// (complexType ValidationCertificateQualification, extends
// ConstraintsConclusion with no own elements).
type XmlValidationCertificateQualification struct {
	XmlConstraintsConclusionContent
	Id                       string                         `xml:"Id,attr"`
	DateTime                 *XSDateTime                    `xml:"DateTime,attr,omitempty"`
	ValidationTime           *ValidationTimeValue           `xml:"ValidationTime,attr,omitempty"`
	CertificateQualification *CertificateQualificationValue `xml:"CertificateQualification,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationCertificateApprovalStatus is the Go form of the generated JAXB
// class XmlValidationCertificateApprovalStatus
// (complexType ValidationCertificateApprovalStatus, extends
// ConstraintsConclusion).
type XmlValidationCertificateApprovalStatus struct {
	XmlConstraintsConclusionContent
	CertificateApprovalStatus *XmlCertificateApprovalStatus `xml:"CertificateApprovalStatus"`
	Id                        string                        `xml:"Id,attr"`
	DateTime                  *XSDateTime                   `xml:"DateTime,attr,omitempty"`
	ValidationTime            *ValidationTimeValue          `xml:"ValidationTime,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlCertificateApprovalStatus is the Go form of the generated JAXB class
// XmlCertificateApprovalStatus (complexType CertificateApprovalStatus).
//
// Its three elements are bound through Adapter12/13/14 to the LoTE-loader
// interfaces (enumerations.ListType, enumerations.LoTEServiceTypeIdentifier,
// enumerations.LoTEServiceStatus - see jaxb_adapters.go's header), which have
// no defined underlying type to hang MarshalText/UnmarshalText off. The type
// implements xml.Marshaler/xml.Unmarshaler directly instead.
type XmlCertificateApprovalStatus struct {
	ListType              enumerations.ListType
	ServiceTypeIdentifier enumerations.LoTEServiceTypeIdentifier
	ServiceStatus         enumerations.LoTEServiceStatus
	Label                 *string
}

// MarshalXML writes ListType/ServiceTypeIdentifier/ServiceStatus as their URI
// lexical form (ListTypeParser.print/etc. all print getUri()), and label as
// the plain attribute value.
func (c *XmlCertificateApprovalStatus) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if c == nil {
		return nil
	}
	if c.Label != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "label"}, Value: *c.Label})
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if c.ListType != nil {
		if err := e.EncodeElement(c.ListType.URI(), xml.StartElement{Name: xml.Name{Local: "ListType"}}); err != nil {
			return err
		}
	}
	if c.ServiceTypeIdentifier != nil {
		if err := e.EncodeElement(c.ServiceTypeIdentifier.URI(), xml.StartElement{Name: xml.Name{Local: "ServiceTypeIdentifier"}}); err != nil {
			return err
		}
	}
	if c.ServiceStatus != nil {
		if err := e.EncodeElement(c.ServiceStatus.URI(), xml.StartElement{Name: xml.Name{Local: "ServiceStatus"}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML resolves ListType/ServiceTypeIdentifier/ServiceStatus back to
// their loader-backed enumeration constant (ListTypeParser.parse/etc. all
// parse fromUri()).
func (c *XmlCertificateApprovalStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Local == "label" {
			v := a.Value
			c.Label = &v
		}
	}
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			var s string
			if err := d.DecodeElement(&s, &t); err != nil {
				return err
			}
			switch t.Name.Local {
			case "ListType":
				c.ListType = enumerations.ListTypeFromURI(s)
			case "ServiceTypeIdentifier":
				c.ServiceTypeIdentifier = enumerations.LoTEServiceTypeIdentifierFromURI(s)
			case "ServiceStatus":
				c.ServiceStatus = enumerations.LoTEServiceStatusFromURI(s)
			}
		case xml.EndElement:
			return nil
		}
	}
}

// XmlQWACProcess is the Go form of the generated JAXB class XmlQWACProcess
// (complexType QWACProcess, extends ConstraintsConclusion).
type XmlQWACProcess struct {
	XmlConstraintsConclusionContent
	ValidationQWACProcess []*XmlValidationQWACProcess `xml:"ValidationQWACProcess,omitempty"`
	Id                    string                      `xml:"Id,attr"`
	QWACType              *QWACProfileValue           `xml:"QWACType,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationQWACProcess is the Go form of the generated JAXB class
// XmlValidationQWACProcess (complexType ValidationQWACProcess, extends
// ConstraintsConclusion with no own elements).
type XmlValidationQWACProcess struct {
	XmlConstraintsConclusionContent
	Id       string            `xml:"Id,attr"`
	QWACType *QWACProfileValue `xml:"QWACType,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationSignatureQualification is the Go form of the generated JAXB
// class XmlValidationSignatureQualification
// (complexType ValidationSignatureQualification, extends ConstraintsConclusion).
type XmlValidationSignatureQualification struct {
	XmlConstraintsConclusionContent
	ValidationCertificateQualification []*XmlValidationCertificateQualification `xml:"ValidationCertificateQualification,omitempty"`
	QWACProcess                        *XmlQWACProcess                          `xml:"QWACProcess,omitempty"`
	SignatureQualification             SignatureQualificationValue              `xml:"SignatureQualification,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationTimestampQualificationAtTime is the Go form of the generated
// JAXB class XmlValidationTimestampQualificationAtTime
// (complexType ValidationTimestampQualificationAtTime, extends
// ConstraintsConclusion with no own elements).
type XmlValidationTimestampQualificationAtTime struct {
	XmlConstraintsConclusionContent
	DateTime               *XSDateTime                 `xml:"DateTime,attr,omitempty"`
	ValidationTime         *ValidationTimeValue        `xml:"ValidationTime,attr,omitempty"`
	TimestampQualification TimestampQualificationValue `xml:"TimestampQualification,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationTimestampQualification is the Go form of the generated JAXB
// class XmlValidationTimestampQualification
// (complexType ValidationTimestampQualification, extends ConstraintsConclusion).
type XmlValidationTimestampQualification struct {
	XmlConstraintsConclusionContent
	ValidationTimestampQualificationAtTime []*XmlValidationTimestampQualificationAtTime `xml:"ValidationTimestampQualificationAtTime,omitempty"`
	TimestampQualification                 TimestampQualificationValue                  `xml:"TimestampQualification,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationEAAQualificationProcess is the Go form of the generated JAXB
// class XmlValidationEAAQualificationProcess
// (complexType ValidationEAAQualificationProcess, extends ConstraintsConclusion
// with no own elements).
type XmlValidationEAAQualificationProcess struct {
	XmlConstraintsConclusionContent
	EAAQualification EAAQualificationValue `xml:"EAAQualification,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationPIDQualificationProcess is the Go form of the generated JAXB
// class XmlValidationPIDQualificationProcess
// (complexType ValidationPIDQualificationProcess, extends ConstraintsConclusion).
type XmlValidationPIDQualificationProcess struct {
	XmlConstraintsConclusionContent
	ValidationCertificateApprovalStatus []*XmlValidationCertificateApprovalStatus `xml:"ValidationCertificateApprovalStatus,omitempty"`
	EAAQualification                    EAAQualificationValue                     `xml:"EAAQualification,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationEAAQualification is the Go form of the generated JAXB class
// XmlValidationEAAQualification (complexType ValidationEAAQualification,
// extends ConstraintsConclusion with no own attributes).
type XmlValidationEAAQualification struct {
	XmlConstraintsConclusionContent
	ValidationEAAQualificationProcess *XmlValidationEAAQualificationProcess `xml:"ValidationEAAQualificationProcess,omitempty"`
	ValidationPIDQualificationProcess *XmlValidationPIDQualificationProcess `xml:"ValidationPIDQualificationProcess,omitempty"`
	EAAQualification                  []EAAQualificationValue               `xml:"EAAQualification"`
	XmlConstraintsConclusionAttrs
}

// XmlLoTEAnalysisAttrs carries the attributes LoTEAnalysis adds over
// ConstraintsConclusion. Own attrs first, base Title last (see
// jaxb_process.go's header for why). TLAnalysis extends LoTEAnalysis with no
// additions, so it reuses this struct rather than defining its own.
type XmlLoTEAnalysisAttrs struct {
	CountryCode *string `xml:"CountryCode,attr,omitempty"`
	URL         string  `xml:"URL,attr"`
	Id          string  `xml:"Id,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlLoTEAnalysis is the Go form of the generated JAXB class XmlLoTEAnalysis
// (complexType LoTEAnalysis, extends ConstraintsConclusion with no own
// elements).
type XmlLoTEAnalysis struct {
	XmlConstraintsConclusionContent
	XmlLoTEAnalysisAttrs
}

// XmlTLAnalysis is the Go form of the generated JAXB class XmlTLAnalysis
// (complexType TLAnalysis, extends LoTEAnalysis with no additions).
type XmlTLAnalysis struct {
	XmlConstraintsConclusionContent
	XmlLoTEAnalysisAttrs
}
