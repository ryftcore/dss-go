// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// CertificateConstraints is the Go form of the generated JAXB class
// CertificateConstraints (complexType CertificateConstraints).
type CertificateConstraints struct {
	Recognition                    *LevelConstraint             `xml:"Recognition,omitempty"`
	Signature                      *LevelConstraint             `xml:"Signature,omitempty"`
	NotExpired                     *LevelConstraint             `xml:"NotExpired,omitempty"`
	SunsetDate                     *LevelConstraint             `xml:"SunsetDate,omitempty"`
	AuthorityInfoAccessPresent     *LevelConstraint             `xml:"AuthorityInfoAccessPresent,omitempty"`
	RevocationDataSkip             *CertificateValuesConstraint `xml:"RevocationDataSkip,omitempty"`
	RevocationInfoAccessPresent    *LevelConstraint             `xml:"RevocationInfoAccessPresent,omitempty"`
	RevocationDataAvailable        *LevelConstraint             `xml:"RevocationDataAvailable,omitempty"`
	AcceptableRevocationDataFound  *LevelConstraint             `xml:"AcceptableRevocationDataFound,omitempty"`
	CRLNextUpdatePresent           *LevelConstraint             `xml:"CRLNextUpdatePresent,omitempty"`
	OCSPNextUpdatePresent          *LevelConstraint             `xml:"OCSPNextUpdatePresent,omitempty"`
	RevocationFreshness            *TimeConstraint              `xml:"RevocationFreshness,omitempty"`
	RevocationFreshnessNextUpdate  *LevelConstraint             `xml:"RevocationFreshnessNextUpdate,omitempty"`
	CA                             *LevelConstraint             `xml:"CA,omitempty"`
	MaxPathLength                  *LevelConstraint             `xml:"MaxPathLength,omitempty"`
	KeyUsage                       *MultiValuesConstraint       `xml:"KeyUsage,omitempty"`
	ExtendedKeyUsage               *MultiValuesConstraint       `xml:"ExtendedKeyUsage,omitempty"`
	PolicyTree                     *LevelConstraint             `xml:"PolicyTree,omitempty"`
	NameConstraints                *LevelConstraint             `xml:"NameConstraints,omitempty"`
	AuthorityKeyIdentifierPresent  *LevelConstraint             `xml:"AuthorityKeyIdentifierPresent,omitempty"`
	SubjectKeyIdentifierPresent    *LevelConstraint             `xml:"SubjectKeyIdentifierPresent,omitempty"`
	NoRevAvail                     *LevelConstraint             `xml:"NoRevAvail,omitempty"`
	SupportedCriticalExtensions    *MultiValuesConstraint       `xml:"SupportedCriticalExtensions,omitempty"`
	ForbiddenExtensions            *MultiValuesConstraint       `xml:"ForbiddenExtensions,omitempty"`
	IssuerName                     *LevelConstraint             `xml:"IssuerName,omitempty"`
	Surname                        *MultiValuesConstraint       `xml:"Surname,omitempty"`
	GivenName                      *MultiValuesConstraint       `xml:"GivenName,omitempty"`
	CommonName                     *MultiValuesConstraint       `xml:"CommonName,omitempty"`
	Pseudonym                      *MultiValuesConstraint       `xml:"Pseudonym,omitempty"`
	Title                          *MultiValuesConstraint       `xml:"Title,omitempty"`
	Email                          *MultiValuesConstraint       `xml:"Email,omitempty"`
	OrganizationIdentifier         *MultiValuesConstraint       `xml:"OrganizationIdentifier,omitempty"`
	OrganizationUnit               *MultiValuesConstraint       `xml:"OrganizationUnit,omitempty"`
	OrganizationName               *MultiValuesConstraint       `xml:"OrganizationName,omitempty"`
	Country                        *MultiValuesConstraint       `xml:"Country,omitempty"`
	Locality                       *MultiValuesConstraint       `xml:"Locality,omitempty"`
	State                          *MultiValuesConstraint       `xml:"State,omitempty"`
	SerialNumberPresent            *LevelConstraint             `xml:"SerialNumberPresent,omitempty"`
	NotRevoked                     *LevelConstraint             `xml:"NotRevoked,omitempty"`
	NotOnHold                      *LevelConstraint             `xml:"NotOnHold,omitempty"`
	RevocationIssuerNotExpired     *LevelConstraint             `xml:"RevocationIssuerNotExpired,omitempty"`
	SelfSigned                     *LevelConstraint             `xml:"SelfSigned,omitempty"`
	NotSelfSigned                  *LevelConstraint             `xml:"NotSelfSigned,omitempty"`
	PolicyIds                      *MultiValuesConstraint       `xml:"PolicyIds,omitempty"`
	PolicyQualificationIds         *LevelConstraint             `xml:"PolicyQualificationIds,omitempty"`
	PolicySupportedByQSCDIds       *LevelConstraint             `xml:"PolicySupportedByQSCDIds,omitempty"`
	QcCompliance                   *LevelConstraint             `xml:"QcCompliance,omitempty"`
	QcEuLimitValueCurrency         *ValueConstraint             `xml:"QcEuLimitValueCurrency,omitempty"`
	MinQcEuLimitValue              *IntValueConstraint          `xml:"MinQcEuLimitValue,omitempty"`
	MinQcEuRetentionPeriod         *IntValueConstraint          `xml:"MinQcEuRetentionPeriod,omitempty"`
	QcSSCD                         *LevelConstraint             `xml:"QcSSCD,omitempty"`
	QcEuPDSLocation                *MultiValuesConstraint       `xml:"QcEuPDSLocation,omitempty"`
	QcType                         *MultiValuesConstraint       `xml:"QcType,omitempty"`
	QcLegislationCountryCodes      *MultiValuesConstraint       `xml:"QcLegislationCountryCodes,omitempty"`
	IssuedToNaturalPerson          *LevelConstraint             `xml:"IssuedToNaturalPerson,omitempty"`
	IssuedToLegalPerson            *LevelConstraint             `xml:"IssuedToLegalPerson,omitempty"`
	SemanticsIdentifier            *MultiValuesConstraint       `xml:"SemanticsIdentifier,omitempty"`
	PSD2QcTypeRolesOfPSP           *MultiValuesConstraint       `xml:"PSD2QcTypeRolesOfPSP,omitempty"`
	PSD2QcCompetentAuthorityName   *MultiValuesConstraint       `xml:"PSD2QcCompetentAuthorityName,omitempty"`
	PSD2QcCompetentAuthorityId     *MultiValuesConstraint       `xml:"PSD2QcCompetentAuthorityId,omitempty"`
	QcQSCDLegislation              *MultiValuesConstraint       `xml:"QcQSCDLegislation,omitempty"`
	QcIdentificationMethod         *MultiValuesConstraint       `xml:"QcIdentificationMethod,omitempty"`
	QcPSBCountryOfLegislation      *MultiValuesConstraint       `xml:"QcPSBCountryOfLegislation,omitempty"`
	QcPSBAuthSourceIdentification  *MultiValuesConstraint       `xml:"QcPSBAuthSourceIdentification,omitempty"`
	QcPSBLegislationIdentification *MultiValuesConstraint       `xml:"QcPSBLegislationIdentification,omitempty"`
	UsePseudonym                   *LevelConstraint             `xml:"UsePseudonym,omitempty"`
	Cryptographic                  *CryptographicConstraint     `xml:"Cryptographic,omitempty"`
}
