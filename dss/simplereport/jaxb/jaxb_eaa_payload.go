// Ported from SimpleReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.simplereport.jaxb. Per the phase-8a generated-JAXB
// rule the generated classes are grouped into schema-area files rather than
// one file per class; every Java class keeps its name and its exact field
// order so that encoding/xml reproduces the JAXB element sequence byte for
// byte.
//
// EAAPayload carries every ETSI TS 119 472-1 / OpenID4VC / mdoc claim the
// simple report can present, one optional element per claim, all bound to
// the same DisclosableClaim (or, for the two parametrised claims,
// ParametrizedDisclosableClaim) simpleContent type; xjc generates one field
// and one accessor pair per claim name rather than a single map, and this
// port keeps that shape so the field list stays a direct, greppable mirror
// of the schema.

package jaxb

// XmlEAAPayload is the Go form of the generated JAXB class XmlEAAPayload
// (complexType EAAPayload).
type XmlEAAPayload struct {
	Identifier                              *XmlDisclosableClaim               `xml:"Identifier,omitempty"`
	Issuer                                  *XmlDisclosableClaim               `xml:"Issuer,omitempty"`
	Subject                                 *XmlDisclosableClaim               `xml:"Subject,omitempty"`
	Audience                                *XmlDisclosableClaim               `xml:"Audience,omitempty"`
	IssuedAt                                *XmlDisclosableClaim               `xml:"IssuedAt,omitempty"`
	NotBefore                               *XmlDisclosableClaim               `xml:"NotBefore,omitempty"`
	Expiration                              *XmlDisclosableClaim               `xml:"Expiration,omitempty"`
	UpdatedAt                               *XmlDisclosableClaim               `xml:"UpdatedAt,omitempty"`
	NextUpdate                              *XmlDisclosableClaim               `xml:"NextUpdate,omitempty"`
	Category                                *XmlDisclosableClaim               `xml:"Category,omitempty"`
	VerifiableCredentialsType               *XmlDisclosableClaim               `xml:"VerifiableCredentialsType,omitempty"`
	StatusIndex                             *XmlDisclosableClaim               `xml:"StatusIndex,omitempty"`
	StatusUri                               *XmlDisclosableClaim               `xml:"StatusUri,omitempty"`
	StatusType                              *XmlDisclosableClaim               `xml:"StatusType,omitempty"`
	StatusPurpose                           *XmlDisclosableClaim               `xml:"StatusPurpose,omitempty"`
	Nonce                                   *XmlDisclosableClaim               `xml:"Nonce,omitempty"`
	DeviceKey                               *XmlDisclosableClaim               `xml:"DeviceKey,omitempty"`
	Version                                 *XmlDisclosableClaim               `xml:"Version,omitempty"`
	DocType                                 *XmlDisclosableClaim               `xml:"DocType,omitempty"`
	AdministrativeIssuanceDate              *XmlDisclosableClaim               `xml:"AdministrativeIssuanceDate,omitempty"`
	AdministrativeExpirationDate            *XmlDisclosableClaim               `xml:"AdministrativeExpirationDate,omitempty"`
	OneTimeUse                              *XmlDisclosableClaim               `xml:"OneTimeUse,omitempty"`
	ShortLived                              *XmlDisclosableClaim               `xml:"ShortLived,omitempty"`
	Evidence                                *XmlDisclosableClaim               `xml:"Evidence,omitempty"`
	AttestedAttributesSubjectId             *XmlDisclosableClaim               `xml:"AttestedAttributesSubjectId,omitempty"`
	AttestedAttributesSubjectFamilyName     *XmlDisclosableClaim               `xml:"AttestedAttributesSubjectFamilyName,omitempty"`
	AttestedAttributesSubjectGivenName      *XmlDisclosableClaim               `xml:"AttestedAttributesSubjectGivenName,omitempty"`
	AttestedAttributesSubjectDocumentNumber *XmlDisclosableClaim               `xml:"AttestedAttributesSubjectDocumentNumber,omitempty"`
	AttestedAttributesSubjectPseudonym      *XmlDisclosableClaim               `xml:"AttestedAttributesSubjectPseudonym,omitempty"`
	AttestedAttributes                      *XmlDisclosableClaim               `xml:"AttestedAttributes,omitempty"`
	FullName                                *XmlDisclosableClaim               `xml:"FullName,omitempty"`
	GivenName                               *XmlDisclosableClaim               `xml:"GivenName,omitempty"`
	FamilyName                              *XmlDisclosableClaim               `xml:"FamilyName,omitempty"`
	MiddleName                              *XmlDisclosableClaim               `xml:"MiddleName,omitempty"`
	Nickname                                *XmlDisclosableClaim               `xml:"Nickname,omitempty"`
	ShortName                               *XmlDisclosableClaim               `xml:"ShortName,omitempty"`
	ProfileUrl                              *XmlDisclosableClaim               `xml:"ProfileUrl,omitempty"`
	PictureUrl                              *XmlDisclosableClaim               `xml:"PictureUrl,omitempty"`
	WebsiteUrl                              *XmlDisclosableClaim               `xml:"WebsiteUrl,omitempty"`
	Email                                   *XmlDisclosableClaim               `xml:"Email,omitempty"`
	EmailVerified                           *XmlDisclosableClaim               `xml:"EmailVerified,omitempty"`
	Gender                                  *XmlDisclosableClaim               `xml:"Gender,omitempty"`
	Birthdate                               *XmlDisclosableClaim               `xml:"Birthdate,omitempty"`
	BirthdateApproximateMask                *XmlDisclosableClaim               `xml:"BirthdateApproximateMask,omitempty"`
	Timezone                                *XmlDisclosableClaim               `xml:"Timezone,omitempty"`
	Locale                                  *XmlDisclosableClaim               `xml:"Locale,omitempty"`
	AddressPostalAddress                    *XmlDisclosableClaim               `xml:"AddressPostalAddress,omitempty"`
	AddressCity                             *XmlDisclosableClaim               `xml:"AddressCity,omitempty"`
	AddressStateOrProvince                  *XmlDisclosableClaim               `xml:"AddressStateOrProvince,omitempty"`
	AddressPostalCode                       *XmlDisclosableClaim               `xml:"AddressPostalCode,omitempty"`
	AddressCountryName                      *XmlDisclosableClaim               `xml:"AddressCountryName,omitempty"`
	AddressStreetAddress                    *XmlDisclosableClaim               `xml:"AddressStreetAddress,omitempty"`
	PhoneNumber                             *XmlDisclosableClaim               `xml:"PhoneNumber,omitempty"`
	PhoneNumberVerified                     *XmlDisclosableClaim               `xml:"PhoneNumberVerified,omitempty"`
	PlaceOfBirth                            *XmlDisclosableClaim               `xml:"PlaceOfBirth,omitempty"`
	PlaceOfBirthCountry                     *XmlDisclosableClaim               `xml:"PlaceOfBirthCountry,omitempty"`
	PlaceOfBirthRegion                      *XmlDisclosableClaim               `xml:"PlaceOfBirthRegion,omitempty"`
	PlaceOfBirthCity                        *XmlDisclosableClaim               `xml:"PlaceOfBirthCity,omitempty"`
	Nationalities                           *XmlDisclosableClaim               `xml:"Nationalities,omitempty"`
	BirthFamilyName                         *XmlDisclosableClaim               `xml:"BirthFamilyName,omitempty"`
	BirthGivenName                          *XmlDisclosableClaim               `xml:"BirthGivenName,omitempty"`
	BirthMiddleName                         *XmlDisclosableClaim               `xml:"BirthMiddleName,omitempty"`
	Salutation                              *XmlDisclosableClaim               `xml:"Salutation,omitempty"`
	Title                                   *XmlDisclosableClaim               `xml:"Title,omitempty"`
	MobilePhoneNumber                       *XmlDisclosableClaim               `xml:"MobilePhoneNumber,omitempty"`
	Pseudonym                               *XmlDisclosableClaim               `xml:"Pseudonym,omitempty"`
	IssuingCountry                          *XmlDisclosableClaim               `xml:"IssuingCountry,omitempty"`
	IssuingAuthority                        *XmlDisclosableClaim               `xml:"IssuingAuthority,omitempty"`
	DocumentNumber                          *XmlDisclosableClaim               `xml:"DocumentNumber,omitempty"`
	Portrait                                *XmlDisclosableClaim               `xml:"Portrait,omitempty"`
	DrivingPrivileges                       *XmlDisclosableClaim               `xml:"DrivingPrivileges,omitempty"`
	UNDistinguishingSign                    *XmlDisclosableClaim               `xml:"UNDistinguishingSign,omitempty"`
	PersonalAdministrativeNumber            *XmlDisclosableClaim               `xml:"PersonalAdministrativeNumber,omitempty"`
	Height                                  *XmlDisclosableClaim               `xml:"Height,omitempty"`
	Weight                                  *XmlDisclosableClaim               `xml:"Weight,omitempty"`
	EyeColour                               *XmlDisclosableClaim               `xml:"EyeColour,omitempty"`
	HairColour                              *XmlDisclosableClaim               `xml:"HairColour,omitempty"`
	ResidentPostalAddress                   *XmlDisclosableClaim               `xml:"ResidentPostalAddress,omitempty"`
	PortraitCaptureDate                     *XmlDisclosableClaim               `xml:"PortraitCaptureDate,omitempty"`
	AgeInYears                              *XmlDisclosableClaim               `xml:"AgeInYears,omitempty"`
	AgeBirthYear                            *XmlDisclosableClaim               `xml:"AgeBirthYear,omitempty"`
	AgeOverNN                               []*XmlParametrizedDisclosableClaim `xml:"AgeOverNN,omitempty"`
	IssuingJurisdiction                     *XmlDisclosableClaim               `xml:"IssuingJurisdiction,omitempty"`
	ResidentAddressCity                     *XmlDisclosableClaim               `xml:"ResidentAddressCity,omitempty"`
	ResidentAddressState                    *XmlDisclosableClaim               `xml:"ResidentAddressState,omitempty"`
	ResidentAddressPostalCode               *XmlDisclosableClaim               `xml:"ResidentAddressPostalCode,omitempty"`
	ResidentAddressCountry                  *XmlDisclosableClaim               `xml:"ResidentAddressCountry,omitempty"`
	BiometricTemplate                       []*XmlParametrizedDisclosableClaim `xml:"BiometricTemplate,omitempty"`
	SignatureUsualMark                      *XmlDisclosableClaim               `xml:"SignatureUsualMark,omitempty"`
	Fingerprint                             *XmlDisclosableClaim               `xml:"Fingerprint,omitempty"`
	BusinessName                            *XmlDisclosableClaim               `xml:"BusinessName,omitempty"`
	OrganizationName                        *XmlDisclosableClaim               `xml:"OrganizationName,omitempty"`
	BirthFullName                           *XmlDisclosableClaim               `xml:"BirthFullName,omitempty"`
	Profession                              *XmlDisclosableClaim               `xml:"Profession,omitempty"`
	RelationshipFather                      *XmlDisclosableClaim               `xml:"RelationshipFather,omitempty"`
	RelationshipMother                      *XmlDisclosableClaim               `xml:"RelationshipMother,omitempty"`
	RelationshipParent                      *XmlDisclosableClaim               `xml:"RelationshipParent,omitempty"`
	RelationshipSon                         *XmlDisclosableClaim               `xml:"RelationshipSon,omitempty"`
	RelationshipDaughter                    *XmlDisclosableClaim               `xml:"RelationshipDaughter,omitempty"`
	RelationshipBrother                     *XmlDisclosableClaim               `xml:"RelationshipBrother,omitempty"`
	RelationshipSister                      *XmlDisclosableClaim               `xml:"RelationshipSister,omitempty"`
	RelationshipSibling                     *XmlDisclosableClaim               `xml:"RelationshipSibling,omitempty"`
	RelationshipSpouse                      *XmlDisclosableClaim               `xml:"RelationshipSpouse,omitempty"`
	RelationshipFatherInLaw                 *XmlDisclosableClaim               `xml:"RelationshipFatherInLaw,omitempty"`
	RelationshipMotherInLaw                 *XmlDisclosableClaim               `xml:"RelationshipMotherInLaw,omitempty"`
	RelationshipParentInLaw                 *XmlDisclosableClaim               `xml:"RelationshipParentInLaw,omitempty"`
	RelationshipSonInLaw                    *XmlDisclosableClaim               `xml:"RelationshipSonInLaw,omitempty"`
	RelationshipDaughterInLaw               *XmlDisclosableClaim               `xml:"RelationshipDaughterInLaw,omitempty"`
	RelationshipChildInLaw                  *XmlDisclosableClaim               `xml:"RelationshipChildInLaw,omitempty"`
	RelationshipParentalAuthority           *XmlDisclosableClaim               `xml:"RelationshipParentalAuthority,omitempty"`
	RelationshipLegalRepresentative         *XmlDisclosableClaim               `xml:"RelationshipLegalRepresentative,omitempty"`
	RelationshipAgent                       *XmlDisclosableClaim               `xml:"RelationshipAgent,omitempty"`
	DocumentType                            *XmlDisclosableClaim               `xml:"DocumentType,omitempty"`
	IssuingAuthorityRegistrationIdentifier  *XmlDisclosableClaim               `xml:"IssuingAuthorityRegistrationIdentifier,omitempty"`
	TrustAnchor                             *XmlDisclosableClaim               `xml:"TrustAnchor,omitempty"`
	ResidentAddressStreet                   *XmlDisclosableClaim               `xml:"ResidentAddressStreet,omitempty"`
	ResidentAddressHouseNumber              *XmlDisclosableClaim               `xml:"ResidentAddressHouseNumber,omitempty"`
	OtherClaim                              []*XmlDisclosableClaim             `xml:"OtherClaim,omitempty"`
}

// eaaPayloadClaimElements lists every element name XmlEAAPayload binds, in
// schema order. jaxb_content_model.go consults it to mark all of them as
// simpleContent (DisclosableClaim / ParametrizedDisclosableClaim are both
// simpleContent extensions of xs:string).
var eaaPayloadClaimElements = []string{
	"Identifier",
	"Issuer",
	"Subject",
	"Audience",
	"IssuedAt",
	"NotBefore",
	"Expiration",
	"UpdatedAt",
	"NextUpdate",
	"Category",
	"VerifiableCredentialsType",
	"StatusIndex",
	"StatusUri",
	"StatusType",
	"StatusPurpose",
	"Nonce",
	"DeviceKey",
	"Version",
	"DocType",
	"AdministrativeIssuanceDate",
	"AdministrativeExpirationDate",
	"OneTimeUse",
	"ShortLived",
	"Evidence",
	"AttestedAttributesSubjectId",
	"AttestedAttributesSubjectFamilyName",
	"AttestedAttributesSubjectGivenName",
	"AttestedAttributesSubjectDocumentNumber",
	"AttestedAttributesSubjectPseudonym",
	"AttestedAttributes",
	"FullName",
	"GivenName",
	"FamilyName",
	"MiddleName",
	"Nickname",
	"ShortName",
	"ProfileUrl",
	"PictureUrl",
	"WebsiteUrl",
	"Email",
	"EmailVerified",
	"Gender",
	"Birthdate",
	"BirthdateApproximateMask",
	"Timezone",
	"Locale",
	"AddressPostalAddress",
	"AddressCity",
	"AddressStateOrProvince",
	"AddressPostalCode",
	"AddressCountryName",
	"AddressStreetAddress",
	"PhoneNumber",
	"PhoneNumberVerified",
	"PlaceOfBirth",
	"PlaceOfBirthCountry",
	"PlaceOfBirthRegion",
	"PlaceOfBirthCity",
	"Nationalities",
	"BirthFamilyName",
	"BirthGivenName",
	"BirthMiddleName",
	"Salutation",
	"Title",
	"MobilePhoneNumber",
	"Pseudonym",
	"IssuingCountry",
	"IssuingAuthority",
	"DocumentNumber",
	"Portrait",
	"DrivingPrivileges",
	"UNDistinguishingSign",
	"PersonalAdministrativeNumber",
	"Height",
	"Weight",
	"EyeColour",
	"HairColour",
	"ResidentPostalAddress",
	"PortraitCaptureDate",
	"AgeInYears",
	"AgeBirthYear",
	"AgeOverNN",
	"IssuingJurisdiction",
	"ResidentAddressCity",
	"ResidentAddressState",
	"ResidentAddressPostalCode",
	"ResidentAddressCountry",
	"BiometricTemplate",
	"SignatureUsualMark",
	"Fingerprint",
	"BusinessName",
	"OrganizationName",
	"BirthFullName",
	"Profession",
	"RelationshipFather",
	"RelationshipMother",
	"RelationshipParent",
	"RelationshipSon",
	"RelationshipDaughter",
	"RelationshipBrother",
	"RelationshipSister",
	"RelationshipSibling",
	"RelationshipSpouse",
	"RelationshipFatherInLaw",
	"RelationshipMotherInLaw",
	"RelationshipParentInLaw",
	"RelationshipSonInLaw",
	"RelationshipDaughterInLaw",
	"RelationshipChildInLaw",
	"RelationshipParentalAuthority",
	"RelationshipLegalRepresentative",
	"RelationshipAgent",
	"DocumentType",
	"IssuingAuthorityRegistrationIdentifier",
	"TrustAnchor",
	"ResidentAddressStreet",
	"ResidentAddressHouseNumber",
	"OtherClaim",
}

// XmlDisclosableClaim is the Go form of the generated JAXB class
// XmlDisclosableClaim (complexType DisclosableClaim).
type XmlDisclosableClaim struct {
	Value      string  `xml:",chardata"`
	Name       *string `xml:"name,attr,omitempty"`
	Disclosure *bool   `xml:"disclosure,attr,omitempty"`
}

// XmlParametrizedDisclosableClaim is the Go form of the generated JAXB class
// XmlParametrizedDisclosableClaim (complexType ParametrizedDisclosableClaim,
// extension of DisclosableClaim). JAXB emits base content before extension
// content and extension attributes before base attributes; DisclosableClaim
// has no element content (its "content" is the chardata Value) so only the
// attribute ordering applies here: Parameter first, then the embedded
// DisclosableClaim's Name/Disclosure.
type XmlParametrizedDisclosableClaim struct {
	Value      string  `xml:",chardata"`
	Parameter  *string `xml:"parameter,attr,omitempty"`
	Name       *string `xml:"name,attr,omitempty"`
	Disclosure *bool   `xml:"disclosure,attr,omitempty"`
}
