// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// EIDAS is the Go form of the generated JAXB class EIDAS (complexType
// eIDAS): constraints related to the European context (trusted list
// validity,...).
type EIDAS struct {
	TLFreshness    *TimeConstraint        `xml:"TLFreshness,omitempty"`
	TLNotExpired   *LevelConstraint       `xml:"TLNotExpired,omitempty"`
	TLWellSigned   *LevelConstraint       `xml:"TLWellSigned,omitempty"`
	TLVersion      *MultiValuesConstraint `xml:"TLVersion,omitempty"`
	TLStructure    *LevelConstraint       `xml:"TLStructure,omitempty"`
	LoTEFreshness  *TimeConstraint        `xml:"LoTEFreshness,omitempty"`
	LoTENotExpired *LevelConstraint       `xml:"LoTENotExpired,omitempty"`
	LoTEWellSigned *LevelConstraint       `xml:"LoTEWellSigned,omitempty"`
	LoTEVersion    *MultiValuesConstraint `xml:"LoTEVersion,omitempty"`
	LoTEStructure  *LevelConstraint       `xml:"LoTEStructure,omitempty"`
}

// ModelConstraint is the Go form of the generated JAXB class ModelConstraint
// (complexType ModelConstraint): a boolean check that follows the specified
// validation model.
type ModelConstraint struct {
	Value ValidationModelValue `xml:"Value,attr,omitempty"`
}
