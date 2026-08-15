// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// TimestampConstraints is the Go form of the generated JAXB class
// TimestampConstraints (complexType TimestampConstraints).
type TimestampConstraints struct {
	TimestampDelay                                            *TimeConstraint              `xml:"TimestampDelay,omitempty"`
	RevocationTimeAgainstBestSignatureTime                    *LevelConstraint             `xml:"RevocationTimeAgainstBestSignatureTime,omitempty"`
	BestSignatureTimeBeforeExpirationDateOfSigningCertificate *LevelConstraint             `xml:"BestSignatureTimeBeforeExpirationDateOfSigningCertificate,omitempty"`
	Coherence                                                 *LevelConstraint             `xml:"Coherence,omitempty"`
	TimestampValid                                            *LevelConstraint             `xml:"TimestampValid,omitempty"`
	BasicSignatureConstraints                                 *BasicSignatureConstraints   `xml:"BasicSignatureConstraints,omitempty"`
	SignedAttributes                                          *SignedAttributesConstraints `xml:"SignedAttributes,omitempty"`
	TSAGeneralNamePresent                                     *LevelConstraint             `xml:"TSAGeneralNamePresent,omitempty"`
	TSAGeneralNameContentMatch                                *LevelConstraint             `xml:"TSAGeneralNameContentMatch,omitempty"`
	TSAGeneralNameOrderMatch                                  *LevelConstraint             `xml:"TSAGeneralNameOrderMatch,omitempty"`
	AtsHashIndex                                              *LevelConstraint             `xml:"AtsHashIndex,omitempty"`
	ContainerSignedAndTimestampedFilesCovered                 *LevelConstraint             `xml:"ContainerSignedAndTimestampedFilesCovered,omitempty"`
}
