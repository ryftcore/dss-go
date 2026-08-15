// Ported from DetailedReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.detailedreport.jaxb. See jaxb_common.go's header for
// the file-grouping rule. This file holds the ConstraintsConclusion extension
// family: every complexType of the schema that extends ConstraintsConclusion,
// ConstraintsConclusionWithControlTime or
// ConstraintsConclusionWithProofOfExistence, following the Content/Attrs
// embedding pattern dss/diagnostic/jaxb/jaxb_common.go establishes for JAXB
// class extension:
//
//   - the JAXB RI marshals base-class content before extension content, so an
//     extension type embeds the base "Content" struct first and its own
//     elements after;
//   - the JAXB RI marshals extension attributes before base-class attributes,
//     so an extension type declares its own attributes first and embeds the
//     base "Attrs" struct last.

package jaxb

// XmlConstraintsConclusionContent carries the element content of
// ConstraintsConclusion. Extension types embed it first.
type XmlConstraintsConclusionContent struct {
	Constraint []*XmlConstraint `xml:"Constraint,omitempty"`
	Conclusion *XmlConclusion   `xml:"Conclusion"`
}

// XmlConstraintsConclusionAttrs carries the attributes of ConstraintsConclusion.
// Extension types embed it last.
type XmlConstraintsConclusionAttrs struct {
	Title string `xml:"Title,attr"`
}

// XmlConstraintsConclusion is the Go form of the generated JAXB class
// XmlConstraintsConclusion (complexType ConstraintsConclusion).
type XmlConstraintsConclusion struct {
	XmlConstraintsConclusionContent
	XmlConstraintsConclusionAttrs
}

// XmlProofOfExistence is the Go form of the generated JAXB class
// XmlProofOfExistence (complexType ProofOfExistence).
type XmlProofOfExistence struct {
	Time        XSDateTime `xml:"Time"`
	TimestampId *string    `xml:"TimestampId,omitempty"`
}

// XmlConstraintsConclusionWithControlTimeContent carries the element content of
// ConstraintsConclusionWithControlTime. Extension types embed it first.
type XmlConstraintsConclusionWithControlTimeContent struct {
	XmlConstraintsConclusionContent
	ControlTime *XSDateTime `xml:"ControlTime,omitempty"`
}

// XmlConstraintsConclusionWithControlTime is the Go form of the generated JAXB
// class XmlConstraintsConclusionWithControlTime
// (complexType ConstraintsConclusionWithControlTime).
type XmlConstraintsConclusionWithControlTime struct {
	XmlConstraintsConclusionWithControlTimeContent
	XmlConstraintsConclusionAttrs
}

// XmlConstraintsConclusionWithProofOfExistenceContent carries the element
// content of ConstraintsConclusionWithProofOfExistence. Extension types embed
// it first.
type XmlConstraintsConclusionWithProofOfExistenceContent struct {
	XmlConstraintsConclusionContent
	ProofOfExistence *XmlProofOfExistence `xml:"ProofOfExistence,omitempty"`
}

// XmlConstraintsConclusionWithProofOfExistence is the Go form of the generated
// JAXB class XmlConstraintsConclusionWithProofOfExistence
// (complexType ConstraintsConclusionWithProofOfExistence).
type XmlConstraintsConclusionWithProofOfExistence struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessBasicSignature is the Go form of the generated JAXB
// class XmlValidationProcessBasicSignature
// (complexType ValidationProcessBasicSignature, extends
// ConstraintsConclusionWithProofOfExistence with no additions).
type XmlValidationProcessBasicSignature struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessBasicTimestamp is the Go form of the generated JAXB
// class XmlValidationProcessBasicTimestamp
// (complexType ValidationProcessBasicTimestamp, extends ConstraintsConclusion).
type XmlValidationProcessBasicTimestamp struct {
	XmlConstraintsConclusionContent
	Type           string     `xml:"Type,attr"`
	ProductionTime XSDateTime `xml:"ProductionTime,attr"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessArchivalDataTimestamp is the Go form of the generated
// JAXB class XmlValidationProcessArchivalDataTimestamp
// (complexType ValidationProcessArchivalDataTimestamp, extends
// ConstraintsConclusionWithProofOfExistence with no additions).
type XmlValidationProcessArchivalDataTimestamp struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessEvidenceRecord is the Go form of the generated JAXB
// class XmlValidationProcessEvidenceRecord
// (complexType ValidationProcessEvidenceRecord, extends
// ConstraintsConclusionWithProofOfExistence).
type XmlValidationProcessEvidenceRecord struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	AOV *XmlAOV `xml:"AOV"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessEAA is the Go form of the generated JAXB class
// XmlValidationProcessEAA (complexType ValidationProcessEAA, extends
// ConstraintsConclusionWithProofOfExistence with no additions).
type XmlValidationProcessEAA struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessLongTermData is the Go form of the generated JAXB class
// XmlValidationProcessLongTermData (complexType ValidationProcessLongTermData,
// extends ConstraintsConclusionWithProofOfExistence).
type XmlValidationProcessLongTermData struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	CRS []*XmlCRS `xml:"CRS,omitempty"`
	RFC []*XmlRFC `xml:"RFC,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlValidationProcessArchivalData is the Go form of the generated JAXB class
// XmlValidationProcessArchivalData (complexType ValidationProcessArchivalData,
// extends ConstraintsConclusionWithProofOfExistence with no additions).
type XmlValidationProcessArchivalData struct {
	XmlConstraintsConclusionWithProofOfExistenceContent
	XmlConstraintsConclusionAttrs
}

// XmlRevocationBasicValidation is the Go form of the generated JAXB class
// XmlRevocationBasicValidation (complexType RevocationBasicValidation, extends
// ConstraintsConclusion).
type XmlRevocationBasicValidation struct {
	XmlConstraintsConclusionContent
	Id *string `xml:"Id,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}
