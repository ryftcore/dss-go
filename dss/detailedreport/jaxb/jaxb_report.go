// Ported from DetailedReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.detailedreport.jaxb. See jaxb_common.go's header for
// the file-grouping rule and jaxb_process.go's header for the
// Content/Attrs extension-embedding pattern. This file holds the token-level
// report blocks (Signature, Timestamp, EvidenceRecord, EAA, Certificate) and
// the DetailedReport document root that lists them.

package jaxb

import "encoding/xml"

// xsiNamespace is the XML Schema instance namespace JAXB declares when it
// needs to write an xsi:type override (see XmlLoTEAnalysisEntry).
const xsiNamespace = "http://www.w3.org/2001/XMLSchema-instance"

// XmlSignature is the Go form of the generated JAXB class XmlSignature
// (complexType Signature).
type XmlSignature struct {
	ValidationProcessBasicSignature  *XmlValidationProcessBasicSignature  `xml:"ValidationProcessBasicSignature,omitempty"`
	Timestamp                        []*XmlTimestamp                      `xml:"Timestamp,omitempty"`
	EvidenceRecord                   []*XmlEvidenceRecord                 `xml:"EvidenceRecord,omitempty"`
	ValidationProcessLongTermData    *XmlValidationProcessLongTermData    `xml:"ValidationProcessLongTermData,omitempty"`
	ValidationProcessArchivalData    *XmlValidationProcessArchivalData    `xml:"ValidationProcessArchivalData,omitempty"`
	ValidationSignatureQualification *XmlValidationSignatureQualification `xml:"ValidationSignatureQualification,omitempty"`
	Conclusion                       *XmlConclusion                       `xml:"Conclusion"`
	Id                               *string                              `xml:"Id,attr,omitempty"`
	CounterSignature                 *bool                                `xml:"CounterSignature,attr,omitempty"`
}

// reportItemName implements XmlReportItem.
func (*XmlSignature) reportItemName() string { return "Signature" }

// XmlTimestamp is the Go form of the generated JAXB class XmlTimestamp
// (complexType Timestamp).
type XmlTimestamp struct {
	ValidationProcessBasicTimestamp        *XmlValidationProcessBasicTimestamp        `xml:"ValidationProcessBasicTimestamp"`
	EvidenceRecord                         []*XmlEvidenceRecord                       `xml:"EvidenceRecord,omitempty"`
	ValidationProcessArchivalDataTimestamp *XmlValidationProcessArchivalDataTimestamp `xml:"ValidationProcessArchivalDataTimestamp,omitempty"`
	ValidationTimestampQualification       *XmlValidationTimestampQualification       `xml:"ValidationTimestampQualification,omitempty"`
	Conclusion                             *XmlConclusion                             `xml:"Conclusion"`
	Id                                     *string                                    `xml:"Id,attr,omitempty"`
}

// reportItemName implements XmlReportItem.
func (*XmlTimestamp) reportItemName() string { return "Timestamp" }

// XmlEvidenceRecord is the Go form of the generated JAXB class
// XmlEvidenceRecord (complexType EvidenceRecord).
type XmlEvidenceRecord struct {
	ValidationProcessEvidenceRecord *XmlValidationProcessEvidenceRecord `xml:"ValidationProcessEvidenceRecord"`
	Timestamp                       []*XmlTimestamp                     `xml:"Timestamp,omitempty"`
	Conclusion                      *XmlConclusion                      `xml:"Conclusion"`
	Id                              *string                             `xml:"Id,attr,omitempty"`
}

// reportItemName implements XmlReportItem.
func (*XmlEvidenceRecord) reportItemName() string { return "EvidenceRecord" }

// XmlEAA is the Go form of the generated JAXB class XmlEAA (complexType EAA).
type XmlEAA struct {
	ValidationProcessEAA       *XmlValidationProcessEAA       `xml:"ValidationProcessEAA"`
	Signature                  []*XmlSignature                `xml:"Signature,omitempty"`
	KeyBindingSignature        *XmlSignature                  `xml:"KeyBindingSignature,omitempty"`
	ValidationEAAQualification *XmlValidationEAAQualification `xml:"ValidationEAAQualification,omitempty"`
	Conclusion                 *XmlConclusion                 `xml:"Conclusion"`
	Id                         *string                        `xml:"Id,attr,omitempty"`
}

// reportItemName implements XmlReportItem.
func (*XmlEAA) reportItemName() string { return "EAA" }

// XmlCertificate is the Go form of the generated JAXB class XmlCertificate
// (complexType Certificate).
type XmlCertificate struct {
	CertificateQualificationProcess  *XmlCertificateQualificationProcess  `xml:"CertificateQualificationProcess"`
	CertificateApprovalStatusProcess *XmlCertificateApprovalStatusProcess `xml:"CertificateApprovalStatusProcess,omitempty"`
	QWACProcess                      *XmlQWACProcess                      `xml:"QWACProcess,omitempty"`
	Id                               *string                              `xml:"Id,attr,omitempty"`
}

// reportItemName implements XmlReportItem.
func (*XmlCertificate) reportItemName() string { return "Certificate" }

// XmlReportItem is satisfied by every type that may appear in DetailedReport's
// top-level choice group (Signature|Timestamp|EvidenceRecord|EAA|Certificate),
// generated as `List<Serializable> signatureOrTimestampOrEvidenceRecord` with
// an @XmlElements choice. Go has no such choice-of-named-elements construct, so
// XmlDetailedReport implements xml.Marshaler/xml.Unmarshaler by hand, the same
// approach dss/diagnostic/jaxb/xml.go's CertificateExtensionsWrapper takes for
// its own named-element choice - except here the choice is the top-level
// document's own field, not a nested wrapper element, so the whole struct
// carries the custom (Un)MarshalXML rather than a dedicated wrapper type.
type XmlReportItem interface {
	// reportItemName returns the element name this type is bound to in the
	// choice group.
	reportItemName() string
}

// reportItemElements is the @XmlElements table of
// DetailedReport.signatureOrTimestampOrEvidenceRecord, in schema order.
var reportItemElements = []struct {
	name string
	new  func() XmlReportItem
}{
	{"Signature", func() XmlReportItem { return &XmlSignature{} }},
	{"Timestamp", func() XmlReportItem { return &XmlTimestamp{} }},
	{"EvidenceRecord", func() XmlReportItem { return &XmlEvidenceRecord{} }},
	{"EAA", func() XmlReportItem { return &XmlEAA{} }},
	{"Certificate", func() XmlReportItem { return &XmlCertificate{} }},
}

func newReportItem(name string) XmlReportItem {
	for _, e := range reportItemElements {
		if e.name == name {
			return e.new()
		}
	}
	return nil
}

// XmlDetailedReport is the Go form of the generated JAXB class
// XmlDetailedReport (complexType DetailedReport, the document root).
//
// MarshalXML/UnmarshalXML below implement the actual (un)marshalling by hand
// (see XmlReportItem's doc comment), so encoding/xml never consults these xml
// struct tags; they are kept accurate regardless so that
// jaxb_content_model.go's reflection-based sweep - which reads the tags
// directly - still checks this type like every other one.
type XmlDetailedReport struct {
	// SignatureOrTimestampOrEvidenceRecord holds the interleaved choice-group
	// items in document order (Signature, Timestamp, EvidenceRecord, EAA and
	// Certificate may occur in any order and repeat).
	SignatureOrTimestampOrEvidenceRecord []XmlReportItem           `xml:"SignatureOrTimestampOrEvidenceRecord"`
	BasicBuildingBlocks                  []*XmlBasicBuildingBlocks `xml:"BasicBuildingBlocks,omitempty"`
	TLAnalysis                           []*XmlTLAnalysis          `xml:"TLAnalysis,omitempty"`
	LoTEAnalysis                         []*XmlLoTEAnalysisEntry   `xml:"LoTEAnalysis,omitempty"`
	Semantic                             []*XmlSemantic            `xml:"Semantic,omitempty"`
	ValidationTime                       *XSDateTime               `xml:"ValidationTime,attr,omitempty"`
}

// XmlLoTEAnalysisEntry is one entry of DetailedReport.LoTEAnalysis. The Java
// field is typed List<XmlLoTEAnalysis>, but callers may place an XmlTLAnalysis
// there (TLAnalysis IS-A LoTEAnalysis in Java); JAXB then marshals it under
// the LoTEAnalysis element name but adds an "xsi:type" override naming the
// runtime type, since it differs from the field's declared type. Go has no
// upcast between the two distinct XmlLoTEAnalysis/XmlTLAnalysis struct types,
// so the choice is recorded explicitly instead
// (see TestMarshalParity/dr-eaa-pid.xml, whose LoTEAnalysis list holds a
// TLAnalysis this way).
type XmlLoTEAnalysisEntry struct {
	XmlLoTEAnalysis *XmlLoTEAnalysis
	// AsTLAnalysis marshals this entry under "xsi:type=TLAnalysis" and
	// records that its element carried that override on unmarshal.
	AsTLAnalysis bool
}

// MarshalXML writes the document element, its namespace, the ValidationTime
// attribute, and every child in propOrder
// (signatureOrTimestampOrEvidenceRecord, basicBuildingBlocks, tlAnalysis,
// loTEAnalysis, semantic).
func (r *XmlDetailedReport) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Space: Namespace, Local: "DetailedReport"}
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
	for _, item := range r.SignatureOrTimestampOrEvidenceRecord {
		if err := e.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: item.reportItemName()}}); err != nil {
			return err
		}
	}
	for _, bbb := range r.BasicBuildingBlocks {
		if err := e.EncodeElement(bbb, xml.StartElement{Name: xml.Name{Local: "BasicBuildingBlocks"}}); err != nil {
			return err
		}
	}
	for _, tl := range r.TLAnalysis {
		if err := e.EncodeElement(tl, xml.StartElement{Name: xml.Name{Local: "TLAnalysis"}}); err != nil {
			return err
		}
	}
	for _, lote := range r.LoTEAnalysis {
		loteStart := xml.StartElement{Name: xml.Name{Local: "LoTEAnalysis"}}
		if lote.AsTLAnalysis {
			loteStart.Attr = append(loteStart.Attr,
				xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: xsiNamespace},
				xml.Attr{Name: xml.Name{Local: "xsi:type"}, Value: "TLAnalysis"},
			)
		}
		if err := e.EncodeElement(lote.XmlLoTEAnalysis, loteStart); err != nil {
			return err
		}
	}
	for _, sem := range r.Semantic {
		if err := e.EncodeElement(sem, xml.StartElement{Name: xml.Name{Local: "Semantic"}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML reads the ValidationTime attribute and rebuilds every child,
// routing the choice-group elements to the type their name is bound to.
func (r *XmlDetailedReport) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Local == "ValidationTime" {
			var v XSDateTime
			if err := v.UnmarshalText([]byte(a.Value)); err != nil {
				return err
			}
			r.ValidationTime = &v
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
			case "BasicBuildingBlocks":
				var bbb XmlBasicBuildingBlocks
				if err := d.DecodeElement(&bbb, &t); err != nil {
					return err
				}
				r.BasicBuildingBlocks = append(r.BasicBuildingBlocks, &bbb)
			case "TLAnalysis":
				var tl XmlTLAnalysis
				if err := d.DecodeElement(&tl, &t); err != nil {
					return err
				}
				r.TLAnalysis = append(r.TLAnalysis, &tl)
			case "LoTEAnalysis":
				var lote XmlLoTEAnalysis
				if err := d.DecodeElement(&lote, &t); err != nil {
					return err
				}
				entry := &XmlLoTEAnalysisEntry{XmlLoTEAnalysis: &lote}
				for _, a := range t.Attr {
					if a.Name.Local == "type" && a.Value == "TLAnalysis" {
						entry.AsTLAnalysis = true
					}
				}
				r.LoTEAnalysis = append(r.LoTEAnalysis, entry)
			case "Semantic":
				var sem XmlSemantic
				if err := d.DecodeElement(&sem, &t); err != nil {
					return err
				}
				r.Semantic = append(r.Semantic, &sem)
			default:
				item := newReportItem(t.Name.Local)
				if item == nil {
					if err := d.Skip(); err != nil {
						return err
					}
					continue
				}
				if err := d.DecodeElement(item, &t); err != nil {
					return err
				}
				r.SignatureOrTimestampOrEvidenceRecord = append(r.SignatureOrTimestampOrEvidenceRecord, item)
			}
		case xml.EndElement:
			return nil
		}
	}
}
