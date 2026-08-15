// Ported from SimpleReport.xsd (DSS 6.5.RC1) via the JAXB class generated
// into eu.europa.esig.dss.simplereport.jaxb.XmlSimpleReport. Per the
// phase-8a generated-JAXB rule the generated classes are grouped into
// schema-area files rather than one file per class.
//
// The SimpleReport complexType's Signature/Timestamp/EvidenceRecord/EAA
// choice sits directly in its element sequence (no wrapping element), so
// encoding/xml's declarative struct-tag matching - which dispatches on a
// single tag name per field - cannot express it. XmlSimpleReport therefore
// implements MarshalXML/UnmarshalXML by hand for the whole element, the way
// dss/diagnostic/jaxb's CertificateExtensionsWrapper does for its
// (dedicated-wrapper) choice group, generalised to a struct with ordinary
// sibling elements before and after the choice.

package jaxb

import (
	"encoding/xml"
	"fmt"
)

// XmlSimpleReport is the Go form of the generated JAXB class
// XmlSimpleReport (complexType SimpleReport, the document element).
type XmlSimpleReport struct {
	ValidationPolicy                     *XmlValidationPolicy
	DocumentName                         *string
	ValidSignaturesCount                 int
	SignaturesCount                      int
	ContainerType                        *ASiCContainerTypeValue
	PDFAInfo                             *XmlPDFAInfo
	SignatureOrTimestampOrEvidenceRecord []XmlTokenItem
	Semantic                             []*XmlSemantic
	ValidationTime                       *XSDateTime
}

// MarshalXML writes the SimpleReport document element, streaming the choice
// group as interleaved Signature/Timestamp/EvidenceRecord/EAA elements
// between PDFAInfo and Semantic, per the schema sequence.
func (r *XmlSimpleReport) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Space: Namespace, Local: "SimpleReport"}
	start.Attr = nil
	if r.ValidationTime != nil {
		text, err := r.ValidationTime.MarshalText()
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "ValidationTime"}, Value: string(text)})
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeElement(r.ValidationPolicy, xml.StartElement{Name: xml.Name{Local: "ValidationPolicy"}}); err != nil {
		return err
	}
	if r.DocumentName != nil {
		if err := e.EncodeElement(*r.DocumentName, xml.StartElement{Name: xml.Name{Local: "DocumentName"}}); err != nil {
			return err
		}
	}
	if err := e.EncodeElement(r.ValidSignaturesCount, xml.StartElement{Name: xml.Name{Local: "ValidSignaturesCount"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(r.SignaturesCount, xml.StartElement{Name: xml.Name{Local: "SignaturesCount"}}); err != nil {
		return err
	}
	if r.ContainerType != nil {
		if err := e.EncodeElement(r.ContainerType, xml.StartElement{Name: xml.Name{Local: "ContainerType"}}); err != nil {
			return err
		}
	}
	if r.PDFAInfo != nil {
		if err := e.EncodeElement(r.PDFAInfo, xml.StartElement{Name: xml.Name{Local: "PDFAInfo"}}); err != nil {
			return err
		}
	}
	for _, item := range r.SignatureOrTimestampOrEvidenceRecord {
		if item == nil {
			continue
		}
		name := item.xmlTokenElementName()
		if err := e.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: name}}); err != nil {
			return err
		}
	}
	for _, s := range r.Semantic {
		if err := e.EncodeElement(s, xml.StartElement{Name: xml.Name{Local: "Semantic"}}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML reads a SimpleReport document element, routing
// Signature/Timestamp/EvidenceRecord/EAA children into the choice-group
// slice regardless of their interleaving.
func (r *XmlSimpleReport) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "ValidationTime" {
			var t XSDateTime
			if err := t.UnmarshalText([]byte(attr.Value)); err != nil {
				return err
			}
			r.ValidationTime = &t
		}
	}
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ValidationPolicy":
				var v XmlValidationPolicy
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				r.ValidationPolicy = &v
			case "DocumentName":
				var v string
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				r.DocumentName = &v
			case "ValidSignaturesCount":
				if err := d.DecodeElement(&r.ValidSignaturesCount, &t); err != nil {
					return err
				}
			case "SignaturesCount":
				if err := d.DecodeElement(&r.SignaturesCount, &t); err != nil {
					return err
				}
			case "ContainerType":
				var v ASiCContainerTypeValue
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				r.ContainerType = &v
			case "PDFAInfo":
				var v XmlPDFAInfo
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				r.PDFAInfo = &v
			case "Semantic":
				var v XmlSemantic
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				r.Semantic = append(r.Semantic, &v)
			case "Signature", "Timestamp", "EvidenceRecord", "EAA":
				item := tokenItemFor(t.Name.Local)
				if item == nil {
					return fmt.Errorf("jaxb: unreachable token element %q", t.Name.Local)
				}
				if err := d.DecodeElement(item, &t); err != nil {
					return err
				}
				r.SignatureOrTimestampOrEvidenceRecord = append(r.SignatureOrTimestampOrEvidenceRecord, item)
			default:
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}
