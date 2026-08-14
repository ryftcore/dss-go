// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule and its LevelConstraint attribute-order
// finding, which this file's RevocationConstraints (extends LevelConstraint,
// no own attributes) follows by embedding LevelConstraint last.
package jaxb

// RevocationConstraints is the Go form of the generated JAXB class
// RevocationConstraints (complexType RevocationConstraints).
type RevocationConstraints struct {
	UnknownStatus                            *LevelConstraint           `xml:"UnknownStatus,omitempty"`
	ThisUpdatePresent                        *LevelConstraint           `xml:"ThisUpdatePresent,omitempty"`
	RevocationIssuerKnown                    *LevelConstraint           `xml:"RevocationIssuerKnown,omitempty"`
	RevocationIssuerValidAtProductionTime    *LevelConstraint           `xml:"RevocationIssuerValidAtProductionTime,omitempty"`
	RevocationAfterCertificateIssuance       *LevelConstraint           `xml:"RevocationAfterCertificateIssuance,omitempty"`
	RevocationHasInformationAboutCertificate *LevelConstraint           `xml:"RevocationHasInformationAboutCertificate,omitempty"`
	OCSPResponderIdMatch                     *LevelConstraint           `xml:"OCSPResponderIdMatch,omitempty"`
	OCSPCertHashPresent                      *LevelConstraint           `xml:"OCSPCertHashPresent,omitempty"`
	OCSPCertHashMatch                        *LevelConstraint           `xml:"OCSPCertHashMatch,omitempty"`
	SelfIssuedOCSP                           *LevelConstraint           `xml:"SelfIssuedOCSP,omitempty"`
	BasicSignatureConstraints                *BasicSignatureConstraints `xml:"BasicSignatureConstraints,omitempty"`
	LevelConstraint
}
