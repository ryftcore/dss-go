// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// EAAConstraints is the Go form of the generated JAXB class EAAConstraints
// (complexType EAAConstraints): group of constraints for the Electronic
// Attestation Of Attributes validation process.
type EAAConstraints struct {
	EAASignatureUnicity                       *LevelConstraint         `xml:"EAASignatureUnicity,omitempty"`
	EAASignatureValid                         *LevelConstraint         `xml:"EAASignatureValid,omitempty"`
	DisclosurePresent                         *LevelConstraint         `xml:"DisclosurePresent,omitempty"`
	DisclosureFound                           *LevelConstraint         `xml:"DisclosureFound,omitempty"`
	DisclosureIntact                          *LevelConstraint         `xml:"DisclosureIntact,omitempty"`
	DisclosureListExhaustive                  *LevelConstraint         `xml:"DisclosureListExhaustive,omitempty"`
	KeyBindingSignaturePresent                *LevelConstraint         `xml:"KeyBindingSignaturePresent,omitempty"`
	KeyBindingSignatureValid                  *LevelConstraint         `xml:"KeyBindingSignatureValid,omitempty"`
	ETSI194721Conformance                     *LevelConstraint         `xml:"ETSI194721Conformance,omitempty"`
	EAAType                                   *MultiValuesConstraint   `xml:"EAAType,omitempty"`
	EAATypeIntegrityPresent                   *LevelConstraint         `xml:"EAATypeIntegrityPresent,omitempty"`
	EAAIdentifierPresent                      *LevelConstraint         `xml:"EAAIdentifierPresent,omitempty"`
	EAAIssuanceDatePresent                    *LevelConstraint         `xml:"EAAIssuanceDatePresent,omitempty"`
	EAANotBeforePresent                       *LevelConstraint         `xml:"EAANotBeforePresent,omitempty"`
	EAAExpirationPresent                      *LevelConstraint         `xml:"EAAExpirationPresent,omitempty"`
	EAANotExpired                             *LevelConstraint         `xml:"EAANotExpired,omitempty"`
	EAAAdministrativeIssuanceDatePresent      *LevelConstraint         `xml:"EAAAdministrativeIssuanceDatePresent,omitempty"`
	EAAAdministrativeExpirationDatePresent    *LevelConstraint         `xml:"EAAAdministrativeExpirationDatePresent,omitempty"`
	EAAAdministrativePeriodNotExpired         *LevelConstraint         `xml:"EAAAdministrativePeriodNotExpired,omitempty"`
	EAACategory                               *MultiValuesConstraint   `xml:"EAACategory,omitempty"`
	EAASubject                                *MultiValuesConstraint   `xml:"EAASubject,omitempty"`
	EAASubjectPseudonym                       *MultiValuesConstraint   `xml:"EAASubjectPseudonym,omitempty"`
	EAAIssuingCountry                         *MultiValuesConstraint   `xml:"EAAIssuingCountry,omitempty"`
	EAAIssuingAuthority                       *MultiValuesConstraint   `xml:"EAAIssuingAuthority,omitempty"`
	EAAIssuingAuthorityRegistrationIdentifier *MultiValuesConstraint   `xml:"EAAIssuingAuthorityRegistrationIdentifier,omitempty"`
	EAAOneTimeUse                             *LevelConstraint         `xml:"EAAOneTimeUse,omitempty"`
	EAAShortLived                             *LevelConstraint         `xml:"EAAShortLived,omitempty"`
	EAARevocationPresent                      *LevelConstraint         `xml:"EAARevocationPresent,omitempty"`
	EAARevocationAvailable                    *LevelConstraint         `xml:"EAARevocationAvailable,omitempty"`
	AcceptableEAARevocationFound              *LevelConstraint         `xml:"AcceptableEAARevocationFound,omitempty"`
	NotRevoked                                *LevelConstraint         `xml:"NotRevoked,omitempty"`
	NotOnHold                                 *LevelConstraint         `xml:"NotOnHold,omitempty"`
	EAAUsePseudonym                           *LevelConstraint         `xml:"EAAUsePseudonym,omitempty"`
	EAAClaims                                 *MultiValuesConstraint   `xml:"EAAClaims,omitempty"`
	EAASupportedClaims                        *MultiValuesConstraint   `xml:"EAASupportedClaims,omitempty"`
	Cryptographic                             *CryptographicConstraint `xml:"Cryptographic,omitempty"`
	LevelConstraint
}

// EAARevocationConstraints is the Go form of the generated JAXB class
// EAARevocationConstraints (complexType EAARevocationConstraints).
type EAARevocationConstraints struct {
	Type                      *MultiValuesConstraint     `xml:"Type,omitempty"`
	UnknownStatus             *LevelConstraint           `xml:"UnknownStatus,omitempty"`
	IssuanceTime              *LevelConstraint           `xml:"IssuanceTime,omitempty"`
	ExpirationTime            *LevelConstraint           `xml:"ExpirationTime,omitempty"`
	NotExpired                *LevelConstraint           `xml:"NotExpired,omitempty"`
	Subject                   *MultiValuesConstraint     `xml:"Subject,omitempty"`
	SubjectMatch              *LevelConstraint           `xml:"SubjectMatch,omitempty"`
	IssuerValidAtIssuanceTime *LevelConstraint           `xml:"IssuerValidAtIssuanceTime,omitempty"`
	BasicSignatureConstraints *BasicSignatureConstraints `xml:"BasicSignatureConstraints,omitempty"`
	LevelConstraint
}
