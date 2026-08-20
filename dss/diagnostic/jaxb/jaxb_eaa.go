// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

import ()

// XmlAuthorizedDataElements is the Go form of the generated JAXB class XmlAuthorizedDataElements
// (complexType AuthorizedDataElements).
type XmlAuthorizedDataElements struct {
	DataElement []string `xml:"DataElement"`
	Namespace   *string  `xml:"namespace,attr,omitempty"`
}

// XmlEAADocument is the Go form of the generated JAXB class XmlEAADocument
// (complexType EAADocument).
type XmlEAADocument struct {
	DocumentType *string      `xml:"DocumentType,omitempty"`
	Errors       []*XmlErrors `xml:"Errors"`
	EAA          *XmlEAA      `xml:"EAA,attr,omitempty"`
}

// XmlEAAPayload is the Go form of the generated JAXB class XmlEAAPayload
// (complexType EAAPayload).
type XmlEAAPayload struct {
	Identifier                             *XmlClaim                          `xml:"Identifier,omitempty"`
	Issuer                                 *XmlClaim                          `xml:"Issuer,omitempty"`
	Subject                                *XmlClaim                          `xml:"Subject,omitempty"`
	Audience                               *XmlClaim                          `xml:"Audience,omitempty"`
	IssuedAt                               *XmlClaim                          `xml:"IssuedAt,omitempty"`
	NotBefore                              *XmlClaim                          `xml:"NotBefore,omitempty"`
	Expiration                             *XmlClaim                          `xml:"Expiration,omitempty"`
	UpdatedAt                              *XmlClaim                          `xml:"UpdatedAt,omitempty"`
	Category                               *XmlClaim                          `xml:"Category,omitempty"`
	VerifiableCredentialsType              *XmlVerifiableCredentialsTypeClaim `xml:"VerifiableCredentialsType,omitempty"`
	Status                                 *XmlStatusClaim                    `xml:"Status,omitempty"`
	Nonce                                  *XmlClaim                          `xml:"Nonce,omitempty"`
	DeviceKey                              *XmlDeviceKeyClaim                 `xml:"DeviceKey,omitempty"`
	Version                                *XmlClaim                          `xml:"Version,omitempty"`
	DocType                                *XmlClaim                          `xml:"DocType,omitempty"`
	ValidityInfo                           *XmlValidityInfoClaim              `xml:"ValidityInfo,omitempty"`
	AdministrativeIssuanceDate             *XmlClaim                          `xml:"AdministrativeIssuanceDate,omitempty"`
	AdministrativeExpirationDate           *XmlClaim                          `xml:"AdministrativeExpirationDate,omitempty"`
	OneTimeUse                             *XmlClaim                          `xml:"OneTimeUse,omitempty"`
	ShortLived                             *XmlClaim                          `xml:"ShortLived,omitempty"`
	Evidence                               *XmlClaim                          `xml:"Evidence,omitempty"`
	AttestedAttributesSubject              *XmlAttestedAttributesSubjectClaim `xml:"AttestedAttributesSubject,omitempty"`
	FullName                               *XmlClaim                          `xml:"FullName,omitempty"`
	GivenName                              *XmlClaim                          `xml:"GivenName,omitempty"`
	FamilyName                             *XmlClaim                          `xml:"FamilyName,omitempty"`
	MiddleName                             *XmlClaim                          `xml:"MiddleName,omitempty"`
	Nickname                               *XmlClaim                          `xml:"Nickname,omitempty"`
	ShortName                              *XmlClaim                          `xml:"ShortName,omitempty"`
	ProfileUrl                             *XmlClaim                          `xml:"ProfileUrl,omitempty"`
	PictureUrl                             *XmlClaim                          `xml:"PictureUrl,omitempty"`
	WebsiteUrl                             *XmlClaim                          `xml:"WebsiteUrl,omitempty"`
	Email                                  *XmlClaim                          `xml:"Email,omitempty"`
	EmailVerified                          *XmlClaim                          `xml:"EmailVerified,omitempty"`
	Gender                                 *XmlClaim                          `xml:"Gender,omitempty"`
	Birthdate                              *XmlBirthdateClaim                 `xml:"Birthdate,omitempty"`
	Timezone                               *XmlClaim                          `xml:"Timezone,omitempty"`
	Locale                                 *XmlClaim                          `xml:"Locale,omitempty"`
	Address                                *XmlAddressClaim                   `xml:"Address,omitempty"`
	PhoneNumber                            *XmlClaim                          `xml:"PhoneNumber,omitempty"`
	PhoneNumberVerified                    *XmlClaim                          `xml:"PhoneNumberVerified,omitempty"`
	PlaceOfBirth                           *XmlPlaceOfBirthClaim              `xml:"PlaceOfBirth,omitempty"`
	Nationalities                          *XmlClaim                          `xml:"Nationalities,omitempty"`
	BirthFamilyName                        *XmlClaim                          `xml:"BirthFamilyName,omitempty"`
	BirthGivenName                         *XmlClaim                          `xml:"BirthGivenName,omitempty"`
	BirthMiddleName                        *XmlClaim                          `xml:"BirthMiddleName,omitempty"`
	Salutation                             *XmlClaim                          `xml:"Salutation,omitempty"`
	Title                                  *XmlClaim                          `xml:"Title,omitempty"`
	MobilePhoneNumber                      *XmlClaim                          `xml:"MobilePhoneNumber,omitempty"`
	Pseudonym                              *XmlClaim                          `xml:"Pseudonym,omitempty"`
	CredentialSubject                      []*XmlCredentialSubjectClaim       `xml:"CredentialSubject"`
	IssuingCountry                         *XmlClaim                          `xml:"IssuingCountry,omitempty"`
	IssuingAuthority                       *XmlClaim                          `xml:"IssuingAuthority,omitempty"`
	DocumentNumber                         *XmlClaim                          `xml:"DocumentNumber,omitempty"`
	Portrait                               *XmlClaim                          `xml:"Portrait,omitempty"`
	DrivingPrivileges                      *XmlDrivingPrivilegesClaim         `xml:"DrivingPrivileges,omitempty"`
	UNDistinguishingSign                   *XmlClaim                          `xml:"UNDistinguishingSign,omitempty"`
	PersonalAdministrativeNumber           *XmlClaim                          `xml:"PersonalAdministrativeNumber,omitempty"`
	Height                                 *XmlClaim                          `xml:"Height,omitempty"`
	Weight                                 *XmlClaim                          `xml:"Weight,omitempty"`
	EyeColour                              *XmlClaim                          `xml:"EyeColour,omitempty"`
	HairColour                             *XmlClaim                          `xml:"HairColour,omitempty"`
	ResidentPostalAddress                  *XmlClaim                          `xml:"ResidentPostalAddress,omitempty"`
	PortraitCaptureDate                    *XmlClaim                          `xml:"PortraitCaptureDate,omitempty"`
	AgeInYears                             *XmlClaim                          `xml:"AgeInYears,omitempty"`
	AgeBirthYear                           *XmlClaim                          `xml:"AgeBirthYear,omitempty"`
	AgeEqualOrOver                         *XmlAgeEqualOrOverClaim            `xml:"AgeEqualOrOver,omitempty"`
	AgeOverNN                              []*XmlAgeOverNNClaim               `xml:"AgeOverNN"`
	IssuingJurisdiction                    *XmlClaim                          `xml:"IssuingJurisdiction,omitempty"`
	ResidentAddressCity                    *XmlClaim                          `xml:"ResidentAddressCity,omitempty"`
	ResidentAddressState                   *XmlClaim                          `xml:"ResidentAddressState,omitempty"`
	ResidentAddressPostalCode              *XmlClaim                          `xml:"ResidentAddressPostalCode,omitempty"`
	ResidentAddressCountry                 *XmlClaim                          `xml:"ResidentAddressCountry,omitempty"`
	BiometricTemplate                      []*XmlBiometricTemplateXXClaim     `xml:"BiometricTemplate"`
	SignatureUsualMark                     *XmlClaim                          `xml:"SignatureUsualMark,omitempty"`
	Fingerprint                            *XmlClaim                          `xml:"Fingerprint,omitempty"`
	BusinessName                           *XmlClaim                          `xml:"BusinessName,omitempty"`
	OrganizationName                       *XmlClaim                          `xml:"OrganizationName,omitempty"`
	BirthFullName                          *XmlClaim                          `xml:"BirthFullName,omitempty"`
	Profession                             *XmlClaim                          `xml:"Profession,omitempty"`
	RelationshipFather                     *XmlClaim                          `xml:"RelationshipFather,omitempty"`
	RelationshipMother                     *XmlClaim                          `xml:"RelationshipMother,omitempty"`
	RelationshipParent                     *XmlClaim                          `xml:"RelationshipParent,omitempty"`
	RelationshipSon                        *XmlClaim                          `xml:"RelationshipSon,omitempty"`
	RelationshipDaughter                   *XmlClaim                          `xml:"RelationshipDaughter,omitempty"`
	RelationshipBrother                    *XmlClaim                          `xml:"RelationshipBrother,omitempty"`
	RelationshipSister                     *XmlClaim                          `xml:"RelationshipSister,omitempty"`
	RelationshipSibling                    *XmlClaim                          `xml:"RelationshipSibling,omitempty"`
	RelationshipSpouse                     *XmlClaim                          `xml:"RelationshipSpouse,omitempty"`
	RelationshipFatherInLaw                *XmlClaim                          `xml:"RelationshipFatherInLaw,omitempty"`
	RelationshipMotherInLaw                *XmlClaim                          `xml:"RelationshipMotherInLaw,omitempty"`
	RelationshipParentInLaw                *XmlClaim                          `xml:"RelationshipParentInLaw,omitempty"`
	RelationshipSonInLaw                   *XmlClaim                          `xml:"RelationshipSonInLaw,omitempty"`
	RelationshipDaughterInLaw              *XmlClaim                          `xml:"RelationshipDaughterInLaw,omitempty"`
	RelationshipChildInLaw                 *XmlClaim                          `xml:"RelationshipChildInLaw,omitempty"`
	RelationshipParentalAuthority          *XmlClaim                          `xml:"RelationshipParentalAuthority,omitempty"`
	RelationshipLegalRepresentative        *XmlClaim                          `xml:"RelationshipLegalRepresentative,omitempty"`
	RelationshipAgent                      *XmlClaim                          `xml:"RelationshipAgent,omitempty"`
	DocumentType                           *XmlClaim                          `xml:"DocumentType,omitempty"`
	IssuingAuthorityRegistrationIdentifier *XmlClaim                          `xml:"IssuingAuthorityRegistrationIdentifier,omitempty"`
	TrustAnchor                            *XmlClaim                          `xml:"TrustAnchor,omitempty"`
	ResidentAddressStreet                  *XmlClaim                          `xml:"ResidentAddressStreet,omitempty"`
	ResidentAddressHouseNumber             *XmlClaim                          `xml:"ResidentAddressHouseNumber,omitempty"`
	OtherClaim                             []*XmlClaim                        `xml:"OtherClaim"`
}

// XmlEAAPresentationInfo is the Go form of the generated JAXB class XmlEAAPresentationInfo
// (complexType EAAPresentationInfo).
type XmlEAAPresentationInfo struct {
	EAAPresentationType *EAAPresentationTypeValue `xml:"EAAPresentationType,omitempty"`
	Version             *string                   `xml:"Version,omitempty"`
	Documents           *DocumentsWrapper         `xml:"Documents"`
	Errors              *XmlErrors                `xml:"Errors,omitempty"`
	Status              *BigInteger               `xml:"Status,omitempty"`
}

// XmlEAARevocationStatus is the Go form of the generated JAXB class XmlEAARevocationStatus
// (complexType EAARevocationStatus).
type XmlEAARevocationStatus struct {
	Status             *EAAStatusValue        `xml:"Status,omitempty"`
	EAARevocationToken *XmlEAARevocationToken `xml:"EAARevocationToken,attr,omitempty"`
}

// XmlEAASignature is the Go form of the generated JAXB class XmlEAASignature
// (complexType anonymous).
type XmlEAASignature struct {
	Signature *XmlSignature `xml:"Signature,attr,omitempty"`
}

// XmlEAASubject is the Go form of the generated JAXB class XmlEAASubject
// (complexType EAASubject).
type XmlEAASubject struct {
	Value string `xml:",chardata"`
	Match *bool  `xml:"match,attr,omitempty"`
}

// XmlError is the Go form of the generated JAXB class XmlError
// (complexType Error).
type XmlError struct {
	Label *string     `xml:"label,attr,omitempty"`
	Code  *BigInteger `xml:"code,attr,omitempty"`
}

// XmlErrors is the Go form of the generated JAXB class XmlErrors
// (complexType Errors).
type XmlErrors struct {
	Error     []*XmlError `xml:"Error"`
	Namespace *string     `xml:"namespace,attr,omitempty"`
}

// XmlKeyAuthorizations is the Go form of the generated JAXB class XmlKeyAuthorizations
// (complexType KeyAuthorizations).
type XmlKeyAuthorizations struct {
	AuthorizedNamespace    []string                     `xml:"AuthorizedNamespace"`
	AuthorizedDataElements []*XmlAuthorizedDataElements `xml:"AuthorizedDataElements"`
}

// XmlKeyBindingPayload is the Go form of the generated JAXB class XmlKeyBindingPayload
// (complexType KeyBindingPayload).
type XmlKeyBindingPayload struct {
	IssuanceTime *XmlClaim   `xml:"IssuanceTime,omitempty"`
	Nonce        *XmlClaim   `xml:"Nonce,omitempty"`
	Audience     *XmlClaim   `xml:"Audience,omitempty"`
	OtherClaim   []*XmlClaim `xml:"OtherClaim"`
}

// XmlKeyBindingSignature is the Go form of the generated JAXB class XmlKeyBindingSignature
// (complexType anonymous).
type XmlKeyBindingSignature struct {
	Signature *XmlSignature `xml:"Signature,attr,omitempty"`
}

// XmlEAA is the Go form of the generated JAXB class XmlEAA
// (complexType EAA).
type XmlEAA struct {
	DocumentName         *string                  `xml:"DocumentName,omitempty"`
	DocumentType         *string                  `xml:"DocumentType,omitempty"`
	EAAType              *EAATypeValue            `xml:"EAAType,omitempty"`
	StructuralValidation *XmlStructuralValidation `xml:"StructuralValidation,omitempty"`
	DigestMethod         *DigestAlgorithmValue    `xml:"DigestMethod,omitempty"`
	DigestMatchers       *DigestMatchersWrapper   `xml:"DigestMatchers"`
	EAASignature         []*XmlEAASignature       `xml:"EAASignature"`
	EAAPayload           *XmlEAAPayload           `xml:"EAAPayload,omitempty"`
	KeyBindingSignature  *XmlKeyBindingSignature  `xml:"KeyBindingSignature,omitempty"`
	KeyBindingPayload    *XmlKeyBindingPayload    `xml:"KeyBindingPayload,omitempty"`
	EAARevocations       *EAARevocationsWrapper   `xml:"EAARevocations"`
	XmlAbstractTokenAttrs
}

// XmlEAARevocationToken is the Go form of the generated JAXB class XmlEAARevocationToken
// (complexType EAARevocationToken).
type XmlEAARevocationToken struct {
	Origin             *EAARevocationOriginValue `xml:"Origin,omitempty"`
	Type               *string                   `xml:"Type,omitempty"`
	SourceAddress      *string                   `xml:"SourceAddress,omitempty"`
	Subject            *XmlEAASubject            `xml:"Subject,omitempty"`
	IssuedAt           *XSDateTime               `xml:"IssuedAt,omitempty"`
	ExpirationTime     *XSDateTime               `xml:"ExpirationTime,omitempty"`
	TimeToLive         *BigInteger               `xml:"TimeToLive,omitempty"`
	BasicSignature     *XmlBasicSignature        `xml:"BasicSignature,omitempty"`
	SigningCertificate *XmlSigningCertificate    `xml:"SigningCertificate,omitempty"`
	CertificateChain   *CertificateChainWrapper  `xml:"CertificateChain"`
	FoundCertificates  *XmlFoundCertificates     `xml:"FoundCertificates,omitempty"`
	Base64Encoded      *Base64Binary             `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue    `xml:"DigestAlgoAndValue,omitempty"`
	XmlAbstractTokenAttrs
}
