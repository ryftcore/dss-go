// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// EvidenceRecordConstraints is the Go form of the generated JAXB class
// EvidenceRecordConstraints (complexType EvidenceRecordConstraints).
type EvidenceRecordConstraints struct {
	EvidenceRecordValid                       *LevelConstraint         `xml:"EvidenceRecordValid,omitempty"`
	DataObjectExistence                       *LevelConstraint         `xml:"DataObjectExistence,omitempty"`
	DataObjectIntact                          *LevelConstraint         `xml:"DataObjectIntact,omitempty"`
	DataObjectFound                           *LevelConstraint         `xml:"DataObjectFound,omitempty"`
	DataObjectGroup                           *LevelConstraint         `xml:"DataObjectGroup,omitempty"`
	SignedFilesCovered                        *LevelConstraint         `xml:"SignedFilesCovered,omitempty"`
	ContainerSignedAndTimestampedFilesCovered *LevelConstraint         `xml:"ContainerSignedAndTimestampedFilesCovered,omitempty"`
	HashTreeRenewal                           *LevelConstraint         `xml:"HashTreeRenewal,omitempty"`
	Cryptographic                             *CryptographicConstraint `xml:"Cryptographic,omitempty"`
	LevelConstraint
}
