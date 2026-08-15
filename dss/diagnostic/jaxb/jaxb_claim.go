// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

import (
	"math/big"
)

// XmlClaimContent carries the element content of the XmlClaim base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlClaimContent struct {
	Text     *string       `xml:"Text,omitempty"`
	Number   *big.Int      `xml:"Number,omitempty"`
	Boolean  *bool         `xml:"Boolean,omitempty"`
	Binary   *Base64Binary `xml:"Binary,omitempty"`
	DateTime *XSDateTime   `xml:"DateTime,omitempty"`
	Item     []*XmlClaim   `xml:"Item"`
	Entry    []*XmlClaim   `xml:"Entry"`
}

// XmlClaimAttrs carries the attributes of the XmlClaim base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlClaimAttrs struct {
	Name       *string `xml:"name,attr,omitempty"`
	Namespace  *string `xml:"namespace,attr,omitempty"`
	Disclosure *bool   `xml:"disclosure,attr,omitempty"`
}

// XmlClaim is the Go form of the generated JAXB class XmlClaim
// (complexType Claim).
type XmlClaim struct {
	XmlClaimContent
	XmlClaimAttrs
}

// XmlDisclosableClaim is the Go form of the generated JAXB class XmlDisclosableClaim
// (complexType DisclosableClaim).
type XmlDisclosableClaim struct {
	Value     string   `xml:",chardata"`
	Id        *big.Int `xml:"id,attr,omitempty"`
	Name      *string  `xml:"name,attr,omitempty"`
	Namespace *string  `xml:"namespace,attr,omitempty"`
}

// XmlAddressClaim is the Go form of the generated JAXB class XmlAddressClaim
// (complexType AddressClaim).
type XmlAddressClaim struct {
	XmlClaimContent
	PostalAddress   *XmlClaim `xml:"PostalAddress,omitempty"`
	City            *XmlClaim `xml:"City,omitempty"`
	StateOrProvince *XmlClaim `xml:"StateOrProvince,omitempty"`
	PostalCode      *XmlClaim `xml:"PostalCode,omitempty"`
	CountryName     *XmlClaim `xml:"CountryName,omitempty"`
	StreetAddress   *XmlClaim `xml:"StreetAddress,omitempty"`
	XmlClaimAttrs
}

// XmlAgeEqualOrOverClaim is the Go form of the generated JAXB class XmlAgeEqualOrOverClaim
// (complexType AgeEqualOrOverClaim).
type XmlAgeEqualOrOverClaim struct {
	XmlClaimContent
	AgeOverNNClaim []*XmlAgeOverNNClaim `xml:"AgeOverNNClaim"`
	XmlClaimAttrs
}

// XmlAgeOverNNClaim is the Go form of the generated JAXB class XmlAgeOverNNClaim
// (complexType AgeOverNNClaim).
type XmlAgeOverNNClaim struct {
	XmlClaimContent
	Age *int `xml:"age,attr,omitempty"`
	XmlClaimAttrs
}

// XmlAttestedAttributesSubjectClaim is the Go form of the generated JAXB class XmlAttestedAttributesSubjectClaim
// (complexType AttestedAttributesSubjectClaim).
type XmlAttestedAttributesSubjectClaim struct {
	XmlClaimContent
	SubjectId        *XmlAttestedAttributesSubjectIdClaim `xml:"SubjectId,omitempty"`
	SubjectPseudonym *XmlClaim                            `xml:"SubjectPseudonym,omitempty"`
	Attributes       *XmlClaim                            `xml:"Attributes,omitempty"`
	XmlClaimAttrs
}

// XmlAttestedAttributesSubjectIdClaim is the Go form of the generated JAXB class XmlAttestedAttributesSubjectIdClaim
// (complexType AttestedAttributesSubjectIdClaim).
type XmlAttestedAttributesSubjectIdClaim struct {
	XmlClaimContent
	FamilyName     *XmlClaim `xml:"FamilyName,omitempty"`
	GivenName      *XmlClaim `xml:"GivenName,omitempty"`
	DocumentNumber *XmlClaim `xml:"DocumentNumber,omitempty"`
	XmlClaimAttrs
}

// XmlBiometricTemplateXXClaim is the Go form of the generated JAXB class XmlBiometricTemplateXXClaim
// (complexType BiometricTemplateXXClaim).
type XmlBiometricTemplateXXClaim struct {
	XmlClaimContent
	Type *string `xml:"type,attr,omitempty"`
	XmlClaimAttrs
}

// XmlBirthdateClaim is the Go form of the generated JAXB class XmlBirthdateClaim
// (complexType BirthdateClaim).
type XmlBirthdateClaim struct {
	XmlClaimContent
	Birthdate       *XmlClaim `xml:"Birthdate,omitempty"`
	ApproximateMask *XmlClaim `xml:"ApproximateMask,omitempty"`
	XmlClaimAttrs
}

// XmlCredentialSubjectClaim is the Go form of the generated JAXB class XmlCredentialSubjectClaim
// (complexType CredentialSubjectClaim).
type XmlCredentialSubjectClaim struct {
	XmlClaimContent
	FullName            *XmlClaim          `xml:"FullName,omitempty"`
	GivenName           *XmlClaim          `xml:"GivenName,omitempty"`
	FamilyName          *XmlClaim          `xml:"FamilyName,omitempty"`
	MiddleName          *XmlClaim          `xml:"MiddleName,omitempty"`
	Nickname            *XmlClaim          `xml:"Nickname,omitempty"`
	ShortName           *XmlClaim          `xml:"ShortName,omitempty"`
	ProfileUrl          *XmlClaim          `xml:"ProfileUrl,omitempty"`
	PictureUrl          *XmlClaim          `xml:"PictureUrl,omitempty"`
	WebsiteUrl          *XmlClaim          `xml:"WebsiteUrl,omitempty"`
	Email               *XmlClaim          `xml:"Email,omitempty"`
	EmailVerified       *XmlClaim          `xml:"EmailVerified,omitempty"`
	Gender              *XmlClaim          `xml:"Gender,omitempty"`
	Birthdate           *XmlBirthdateClaim `xml:"Birthdate,omitempty"`
	Timezone            *XmlClaim          `xml:"Timezone,omitempty"`
	Locale              *XmlClaim          `xml:"Locale,omitempty"`
	Address             *XmlAddressClaim   `xml:"Address,omitempty"`
	PhoneNumber         *XmlClaim          `xml:"PhoneNumber,omitempty"`
	PhoneNumberVerified *XmlClaim          `xml:"PhoneNumberVerified,omitempty"`
	PlaceOfBirth        *XmlClaim          `xml:"PlaceOfBirth,omitempty"`
	Nationalities       *XmlClaim          `xml:"Nationalities,omitempty"`
	BirthFamilyName     *XmlClaim          `xml:"BirthFamilyName,omitempty"`
	BirthGivenName      *XmlClaim          `xml:"BirthGivenName,omitempty"`
	BirthMiddleName     *XmlClaim          `xml:"BirthMiddleName,omitempty"`
	Salutation          *XmlClaim          `xml:"Salutation,omitempty"`
	Title               *XmlClaim          `xml:"Title,omitempty"`
	MobilePhoneNumber   *XmlClaim          `xml:"MobilePhoneNumber,omitempty"`
	Pseudonym           *XmlClaim          `xml:"Pseudonym,omitempty"`
	OtherClaim          []*XmlClaim        `xml:"OtherClaim"`
	XmlClaimAttrs
}

// XmlDeviceKeyClaim is the Go form of the generated JAXB class XmlDeviceKeyClaim
// (complexType DeviceKeyClaim).
type XmlDeviceKeyClaim struct {
	XmlClaimContent
	PublicKey          *Base64Binary            `xml:"PublicKey,omitempty"`
	X509Certificate    []*XmlX509Certificate    `xml:"X509Certificate"`
	DigestAlgoAndValue []*XmlDigestAlgoAndValue `xml:"DigestAlgoAndValue"`
	KID                []string                 `xml:"KID"`
	X509Url            []string                 `xml:"X509Url"`
	KeyAuthorizations  *XmlKeyAuthorizations    `xml:"KeyAuthorizations,omitempty"`
	XmlClaimAttrs
}

// XmlDrivingPrivilegeClaim is the Go form of the generated JAXB class XmlDrivingPrivilegeClaim
// (complexType DrivingPrivilegeClaim).
type XmlDrivingPrivilegeClaim struct {
	XmlClaimContent
	VehicleCategoryCode *XmlClaim                      `xml:"VehicleCategoryCode,omitempty"`
	IssueDate           *XmlClaim                      `xml:"IssueDate,omitempty"`
	ExpiryDate          *XmlClaim                      `xml:"ExpiryDate,omitempty"`
	Codes               *XmlDrivingPrivilegeCodesClaim `xml:"Codes,omitempty"`
	XmlClaimAttrs
}

// XmlDrivingPrivilegeCodeClaim is the Go form of the generated JAXB class XmlDrivingPrivilegeCodeClaim
// (complexType DrivingPrivilegeCodeClaim).
type XmlDrivingPrivilegeCodeClaim struct {
	XmlClaimContent
	Code  *XmlClaim `xml:"Code,omitempty"`
	Sign  *XmlClaim `xml:"Sign,omitempty"`
	Value *XmlClaim `xml:"Value,omitempty"`
	XmlClaimAttrs
}

// XmlDrivingPrivilegeCodesClaim is the Go form of the generated JAXB class XmlDrivingPrivilegeCodesClaim
// (complexType DrivingPrivilegeCodesClaim).
type XmlDrivingPrivilegeCodesClaim struct {
	XmlClaimContent
	Code []*XmlDrivingPrivilegeCodeClaim `xml:"Code"`
	XmlClaimAttrs
}

// XmlDrivingPrivilegesClaim is the Go form of the generated JAXB class XmlDrivingPrivilegesClaim
// (complexType DrivingPrivilegesClaim).
type XmlDrivingPrivilegesClaim struct {
	XmlClaimContent
	DrivingPrivilege []*XmlDrivingPrivilegeClaim `xml:"DrivingPrivilege"`
	XmlClaimAttrs
}

// XmlIdentifierListClaim is the Go form of the generated JAXB class XmlIdentifierListClaim
// (complexType IdentifierListClaim).
type XmlIdentifierListClaim struct {
	XmlClaimContent
	Identifier  *XmlClaim `xml:"Identifier,omitempty"`
	Uri         *XmlClaim `xml:"Uri,omitempty"`
	Certificate *XmlClaim `xml:"Certificate,omitempty"`
	XmlClaimAttrs
}

// XmlIntegrityClaim is the Go form of the generated JAXB class XmlIntegrityClaim
// (complexType IntegrityClaim).
type XmlIntegrityClaim struct {
	XmlClaimContent
	DigestMethod *DigestAlgorithmValue `xml:"DigestMethod,omitempty"`
	DigestValue  *Base64Binary         `xml:"DigestValue,omitempty"`
	XmlClaimAttrs
}

// XmlPlaceOfBirthClaim is the Go form of the generated JAXB class XmlPlaceOfBirthClaim
// (complexType PlaceOfBirthClaim).
type XmlPlaceOfBirthClaim struct {
	XmlClaimContent
	Country *XmlClaim `xml:"Country,omitempty"`
	Region  *XmlClaim `xml:"Region,omitempty"`
	City    *XmlClaim `xml:"City,omitempty"`
	XmlClaimAttrs
}

// XmlStatusClaim is the Go form of the generated JAXB class XmlStatusClaim
// (complexType StatusClaim).
type XmlStatusClaim struct {
	XmlClaimContent
	StatusList     *XmlStatusListClaim     `xml:"StatusList,omitempty"`
	IdentifierList *XmlIdentifierListClaim `xml:"IdentifierList,omitempty"`
	Index          *XmlClaim               `xml:"Index,omitempty"`
	Uri            *XmlClaim               `xml:"Uri,omitempty"`
	Type           *XmlClaim               `xml:"Type,omitempty"`
	Purpose        *XmlClaim               `xml:"Purpose,omitempty"`
	XmlClaimAttrs
}

// XmlStatusListClaim is the Go form of the generated JAXB class XmlStatusListClaim
// (complexType StatusListClaim).
type XmlStatusListClaim struct {
	XmlClaimContent
	Index       *XmlClaim `xml:"Index,omitempty"`
	Uri         *XmlClaim `xml:"Uri,omitempty"`
	Certificate *XmlClaim `xml:"Certificate,omitempty"`
	XmlClaimAttrs
}

// XmlValidityInfoClaim is the Go form of the generated JAXB class XmlValidityInfoClaim
// (complexType ValidityInfoClaim).
type XmlValidityInfoClaim struct {
	XmlClaimContent
	Signed         *XmlClaim `xml:"Signed,omitempty"`
	ValidFrom      *XmlClaim `xml:"ValidFrom,omitempty"`
	ValidUntil     *XmlClaim `xml:"ValidUntil,omitempty"`
	ExpectedUpdate *XmlClaim `xml:"ExpectedUpdate,omitempty"`
	XmlClaimAttrs
}

// XmlVerifiableCredentialsTypeClaim is the Go form of the generated JAXB class XmlVerifiableCredentialsTypeClaim
// (complexType VerifiableCredentialsTypeClaim).
type XmlVerifiableCredentialsTypeClaim struct {
	XmlClaimContent
	Integrity *XmlIntegrityClaim `xml:"Integrity,omitempty"`
	XmlClaimAttrs
}
