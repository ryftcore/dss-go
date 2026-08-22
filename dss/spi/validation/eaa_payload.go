// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAPayload.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.eaa.EAAPayload lands in this same Go package ("spi.eaa" is one of
// the packages flattened into dss/spi/validation).
package validation

import "github.com/ryftcore/dss-go/dss/model/eaa/claim"

// EAAPayload provides an interface for accessing the content of the EAA payload.
type EAAPayload interface {
	claim.Claim

	// Identifier gets the EAA's unique identifier, when present. Port of getIdentifier().
	Identifier() claim.String
	// Issuer gets the EAA's issuer, when present. Port of getIssuer().
	Issuer() claim.String
	// Subject gets the EAA's subject, when present. Port of getSubject().
	Subject() claim.String
	// Audience gets the list of recipients the EAA is intended for, when present. Port of
	// getAudience().
	Audience() claim.Array
	// IssuedAtTime gets the time at which the EAA was issued, when present. Port of
	// getIssuedAtTime().
	IssuedAtTime() claim.Date
	// NotBeforeTime gets the time before which the EAA is not accepted for processing, when
	// present. Port of getNotBeforeTime().
	NotBeforeTime() claim.Date
	// ExpirationTime gets the expiration time of the EAA, after which the EAA is not accepted
	// for processing, when present. Port of getExpirationTime().
	ExpirationTime() claim.Date
	// UpdatedAtTime gets the time at which the information present within the EAA was the last
	// time updated, when present. Port of getUpdatedAtTime().
	UpdatedAtTime() claim.Date
	// DeviceKey gets the wallet holder's key. Port of getDeviceKey().
	DeviceKey() claim.DeviceKey
	// Category gets the EAA category URN, when present. Port of getCategory().
	Category() claim.String
	// VerifiableCredentialsType gets the EAA's Metadata type, when present. Port of
	// getVerifiableCredentialsType().
	VerifiableCredentialsType() claim.String
	// VerifiableCredentialsTypeIntegrity gets the EAA's Metadata integrity claim, when present.
	// Port of getVerifiableCredentialsTypeIntegrity().
	VerifiableCredentialsTypeIntegrity() claim.Integrity
	// Status gets the EAA's Status value, when present. Port of getStatus().
	Status() claim.Status
	// Nonce gets the EAA's nonce value, used to associate the Client's session Id with the EAA,
	// when present. Port of getNonce().
	Nonce() claim.String
	// FullName gets the user's full name information, when present. Port of getFullName().
	FullName() claim.String
	// GivenName gets the user's first or given name information, when present. Port of
	// getGivenName().
	GivenName() claim.String
	// FamilyName gets the user's last name or surname information, when present. Port of
	// getFamilyName().
	FamilyName() claim.String
	// MiddleName gets the user's middle name information, when present. Port of getMiddleName().
	MiddleName() claim.String
	// Nickname gets the user's casual name information, when present. Port of getNickname().
	Nickname() claim.String
	// ShortName gets the user's preferred name, usually a shorthand name, when present. Port of
	// getShortName().
	ShortName() claim.String
	// ProfileUrl gets the user's profile page URL, when present. Port of getProfileUrl().
	ProfileUrl() claim.String
	// PictureUrl gets the user's profile picture URL, when present. Port of getPictureUrl().
	PictureUrl() claim.String
	// WebsiteUrl gets the user's website or blog URL, when present. Port of getWebsiteUrl().
	WebsiteUrl() claim.String
	// Email gets the user's preferred email address, when present. Port of getEmail().
	Email() claim.String
	// EmailVerified gets whether the user's email address has been verified, when present. Port
	// of getEmailVerified().
	EmailVerified() claim.Boolean
	// Gender gets the user's gender, when present. Port of getGender().
	Gender() claim.Claim
	// Birthdate gets the user's birthdate, when present. Port of getBirthdate().
	Birthdate() claim.Claim
	// Timezone gets the user's TimeZone, when present. Port of getTimezone().
	Timezone() claim.String
	// Locale gets the user's locale, when present. Port of getLocale().
	Locale() claim.String
	// Address gets the user's full postal or physical address, when present. Port of
	// getAddress().
	Address() claim.Address
	// PhoneNumber gets the user's preferred telephone number, when present. Port of
	// getPhoneNumber().
	PhoneNumber() claim.String
	// PhoneNumberVerified gets whether the user's preferred telephone number has been verified,
	// when present. Port of getPhoneNumberVerified().
	PhoneNumberVerified() claim.Boolean
	// PlaceOfBirth gets user's place of birth, when present. Port of getPlaceOfBirth().
	PlaceOfBirth() claim.Claim
	// Nationalities gets user's nationalities using ICAO 3-letter codes, when present. Port of
	// getNationalities().
	Nationalities() claim.Claim
	// BirthGivenName gets user's first or given name when they were born, when present. Port of
	// getBirthGivenName().
	BirthGivenName() claim.String
	// BirthFamilyName gets user's family or last name when they were born, when present. Port
	// of getBirthFamilyName().
	BirthFamilyName() claim.String
	// BirthMiddleName gets user's middle name when they were born, when present. Port of
	// getBirthMiddleName().
	BirthMiddleName() claim.String
	// Salutation gets user's salutation, e.g., "Mr", when present. Port of getSalutation().
	Salutation() claim.String
	// Title gets user's title, e.g., "Dr", when present. Port of getTitle().
	Title() claim.String
	// MobilePhoneNumber gets user's mobile phone number, when present. Port of
	// getMobilePhoneNumber().
	MobilePhoneNumber() claim.String
	// Pseudonym gets user's stage name, religious name or any other type of alias/pseudonym,
	// when present. Port of getPseudonym().
	Pseudonym() claim.String
	// CredentialSubjects returns a list of "4.8 Credential Subject" claims defined in W3C
	// Verifiable Credentials Data Model v2.0. Port of getCredentialSubjects().
	CredentialSubjects() []claim.CredentialSubject

	/* Mdoc specific payload headers as per ISO/IEC 18013-5 */

	// IssuingCountry gets alpha-2 country code, as defined in ISO 3166-1, of the issuing
	// authority's country or territory. Port of getIssuingCountry().
	IssuingCountry() claim.String
	// IssuingAuthority gets issuing authority name. The value shall only use latin1 characters
	// and shall have a maximum length of 150 characters. Port of getIssuingAuthority().
	IssuingAuthority() claim.String
	// DocumentNumber gets the number assigned or calculated by the issuing authority. The value
	// shall only use latin1 characters and shall have a maximum length of 150 characters. Port
	// of getDocumentNumber().
	DocumentNumber() claim.String
	// Portrait gets a reproduction of the mDL holder's portrait. Port of getPortrait().
	Portrait() claim.ByteString
	// DrivingPrivileges gets driving privileges of the mDL holder. Port of
	// getDrivingPrivileges().
	DrivingPrivileges() claim.DrivingPrivileges
	// UNDistinguishingSign gets the distinguishing sign of the issuing country according to
	// ISO/IEC 18013-1:2018, Annex F. If no applicable distinguishing sign is available in
	// ISO/IEC 18013-1, an IA may use an empty identifier or another identifier by which it is
	// internationally recognized. In this case the IA should ensure there is no collision with
	// other IA's. Port of getUNDistinguishingSign().
	UNDistinguishingSign() claim.String
	// PersonalAdministrativeNumber is an audit control number assigned by the issuing
	// authority. The value shall only use latin1 characters and shall have a maximum length of
	// 150 characters. Port of getPersonalAdministrativeNumber().
	PersonalAdministrativeNumber() claim.String
	// Height gets the holder's height in centimetres. Port of getHeight().
	Height() claim.Number
	// Weight gets the holder's weight in kilograms (Java's doc comment repeats "height in
	// centimetres" verbatim; kept as-is - this is a copy-paste doc bug in upstream, not
	// reproduced-then-fixed here per faithful porting). Port of getWeight().
	Weight() claim.Number
	// EyeColour gets the mDL holder's eye colour. The value shall be one of the following:
	// "black", "blue", "brown", "dichromatic", "grey", "green", "hazel", "maroon", "pink",
	// "unknown". Port of getEyeColour().
	EyeColour() claim.String
	// HairColour gets the mDL holder's hair colour. The value shall be one of the following:
	// "bald", "black", "blond", "brown", "grey", "red", "auburn", "sandy", "white", "unknown".
	// Port of getHairColour().
	HairColour() claim.String
	// PostalAddress gets the place where the mDL holder resides and/or may be contacted
	// (street/house number, municipality etc.). The value shall only use latin1 characters and
	// shall have a maximum length of 150 characters. Port of getPostalAddress().
	PostalAddress() claim.String
	// PortraitCaptureDate gets the date when portrait was taken. Port of
	// getPortraitCaptureDate().
	PortraitCaptureDate() claim.Date
	// AgeInYears gets the date the age of the mDL holder. Port of getAgeInYears().
	AgeInYears() claim.Number
	// AgeBirthYear gets the year when the mDL holder was born. Port of getAgeBirthYear().
	AgeBirthYear() claim.Number
	// AgeEqualOrOver gets a map of elements attesting whether the User to whom the person
	// identification data relates is at least NN years old. N <> 18. Multiple instances of
	// this attribute may be present, provided the value of NN is different in each of them. If
	// present, the requirements in clause 7.2.5 of ISO/IEC 18013-5 are applicable for these
	// attributes. Port of getAgeEqualOrOver().
	AgeEqualOrOver() claim.AgeEqualOrOver
	// AgeOverNN gets a list of elements used to convey to an mDL verifier, in a data-minimized
	// fashion, if the mDL holder is as old or older than a specified age, or if the mDL holder
	// is younger than a specified age. To achieve this, the mDL contains age attestation
	// identifiers. An age attestation identifier has the format age_over_NN where NN is a
	// value from 00 to 99. The value of an age attestation identifier can be TRUE or FALSE.
	// Port of getAgeOverNN().
	AgeOverNN() []claim.AgeOverNN
	// IssuingJurisdiction gets a country subdivision code of the jurisdiction that issued the
	// mDL as defined in ISO 3166-2:2020, Clause 8. The first part of the code shall be the same
	// as the value for issuing_country. Port of getIssuingJurisdiction().
	IssuingJurisdiction() claim.String
	// ResidentAddressCity gets the city where the mDL holder lives. The value shall only use
	// latin1 characters and shall have a maximum length of 150 characters. Port of
	// getResidentAddressCity().
	ResidentAddressCity() claim.String
	// ResidentAddressState gets the state/province/district where the mDL holder lives. The
	// value shall only use latin1 characters and shall have a maximum length of 150 characters.
	// Port of getResidentAddressState().
	ResidentAddressState() claim.String
	// ResidentAddressPostalCode gets the postal code of the mDL holder. The value shall only
	// use latin1b characters and shall have a maximum length of 150 characters. Port of
	// getResidentAddressPostalCode().
	ResidentAddressPostalCode() claim.String
	// ResidentAddressCountry gets the country where the mDL holder lives as a two letter
	// country code (alpha-2 code) defined in ISO 3166-1. Port of getResidentAddressCountry().
	ResidentAddressCountry() claim.String
	// BiometricTemplate gets a list of elements containing optional facial, fingerprint, iris,
	// or other biometric information of the mDL holder. A biometric template identifier has the
	// format biometric_template_xx where xx shall be replaced with the corresponding "Abstract
	// value name" found in ISO/IEC 19785 3:2020, Table 7, according to the following
	// convention: capitalized characters are replaced with their lowercase equivalent and
	// spaces or non-alphanumeric characters are replaced by underscores (_). Port of
	// getBiometricTemplate().
	BiometricTemplate() []claim.BiometricTemplateXX
	// SignatureUsualMark gets an image of the signature or usual mark of the mDL holder, see
	// 7.2.7 ISO/IEC 18013-5. Port of getSignatureUsualMark().
	SignatureUsualMark() claim.ByteString

	/* "9.1.2.4 Signing method and structure for MSO" headers as per ISO/IEC 18013-5 */

	// Version gets a version of the MobileSecurityObject. Port of getVersion().
	Version() claim.String
	// DocType gets a docType as used in Documents. NOTE: This a mandatory non-disclosable
	// property in comparison with DocumentType. Port of getDocType().
	DocType() claim.String
	// ValidityInfo gets the information related to the validity of the MSO and its signature.
	// Port of getValidityInfo().
	ValidityInfo() claim.ValidityInfo

	/* Mdoc specific payload headers as per ISO/IEC 23220-2 */

	// Fingerprint gets a reproduction of the holder's fingerprint data (TBC). Port of
	// getFingerprint().
	Fingerprint() claim.ByteString
	// BusinessName gets a business name of the holder. Port of getBusinessName().
	BusinessName() claim.String
	// OrganizationName gets a name of legal person. Port of getOrganizationName().
	OrganizationName() claim.String
	// BirthFullName gets the name(s) which holder was born. Port of getBirthFullName().
	BirthFullName() claim.String
	// Profession gets the profession of the holder. Port of getProfession().
	Profession() claim.String

	/* "6.3.2.3 Relationship attributes" headers as per ISO/IEC 23220-2 */

	// RelationshipFather gets the father of the holder. Port of getRelationshipFather().
	RelationshipFather() claim.String
	// RelationshipMother gets the mother of the holder. Port of getRelationshipMother().
	RelationshipMother() claim.String
	// RelationshipParent gets the parent of the holder. Port of getRelationshipParent().
	RelationshipParent() claim.String
	// RelationshipSon gets the son of the holder. Port of getRelationshipSon().
	RelationshipSon() claim.String
	// RelationshipDaughter gets the daughter of the holder. Port of getRelationshipDaughter().
	RelationshipDaughter() claim.String
	// RelationshipBrother gets the brother of the holder. Port of getRelationshipBrother().
	RelationshipBrother() claim.String
	// RelationshipSister gets the sister of the holder. Port of getRelationshipSister().
	RelationshipSister() claim.String
	// RelationshipSibling gets the sibling of the holder. Port of getRelationshipSibling().
	RelationshipSibling() claim.String
	// RelationshipSpouse gets the spouse of the holder. Port of getRelationshipSpouse().
	RelationshipSpouse() claim.String
	// RelationshipFatherInLaw gets the father-in-law of the holder. Port of
	// getRelationshipFatherInLaw().
	RelationshipFatherInLaw() claim.String
	// RelationshipMotherInLaw gets the mother-in-law of the holder. Port of
	// getRelationshipMotherInLaw().
	RelationshipMotherInLaw() claim.String
	// RelationshipParentInLaw gets the parent-in-law of the holder. Port of
	// getRelationshipParentInLaw().
	RelationshipParentInLaw() claim.String
	// RelationshipSonInLaw gets the son-in-law of the holder. Port of
	// getRelationshipSonInLaw().
	RelationshipSonInLaw() claim.String
	// RelationshipDaughterInLaw gets the daughter-in-law of the holder. Port of
	// getRelationshipDaughterInLaw().
	RelationshipDaughterInLaw() claim.String
	// RelationshipChildInLaw gets the child-in-law of the holder. Port of
	// getRelationshipChildInLaw().
	RelationshipChildInLaw() claim.String
	// RelationshipParentalAuthority gets the parental authority of the holder. Port of
	// getRelationshipParentalAuthority().
	RelationshipParentalAuthority() claim.String
	// RelationshipLegalRepresentative gets the legal representative of the holder. Port of
	// getRelationshipLegalRepresentative().
	RelationshipLegalRepresentative() claim.String
	// RelationshipAgent gets the voluntary agent of the holder. Port of
	// getRelationshipAgent().
	RelationshipAgent() claim.String

	/* "6.3.4 Data elements for document entity" headers as per ISO/IEC 23220-2 */

	// DocumentType gets the document type. NOTE: This a selectively disclosable property in
	// comparison with DocType. Port of getDocumentType().
	DocumentType() claim.String

	/* ARF PID Rulebook headers */

	// AdministrativeIssuanceDate gets the date when the data (e.g. a PID) was issued. Port of
	// getAdministrativeIssuanceDate().
	AdministrativeIssuanceDate() claim.Date
	// AdministrativeExpirationDate gets the date when the data (e.g. a PID) will expire. Port
	// of getAdministrativeExpirationDate().
	AdministrativeExpirationDate() claim.Date
	// TrustAnchor gets the URL at which a machine-readable version of the trust anchor to be
	// used for verifying the PID can be found or looked up. Port of getTrustAnchor().
	TrustAnchor() claim.String
	// ResidentAddressStreet gets the name of the street where the user to whom the person
	// identification data relates currently resides. Port of getResidentAddressStreet().
	ResidentAddressStreet() claim.String
	// ResidentAddressHouseNumber gets the house number where the user to whom the person
	// identification data relates currently resides, including any affix or suffix. Port of
	// getResidentAddressHouseNumber().
	ResidentAddressHouseNumber() claim.String

	/* ETSI TS 119 472-1 "5 Implementation of EAA based on SD-JWT VC" header parameters */

	// IssuingAuthorityRegistrationIdentifier gets the registration identifier of the legal
	// entity on whose behalf the EAA has been issued. Port of
	// getIssuingAuthorityRegistrationIdentifier().
	IssuingAuthorityRegistrationIdentifier() claim.String
	// OneTimeUse gets the signal indicating that the EAA shall be used only once, and that it
	// shall not be retained for future use. Port of getOneTimeUse().
	OneTimeUse() claim.Claim
	// ShortLived gets the EAA short-lived component indicating that the validity period of the
	// EAA is so short that it shall not be necessary to check its revocation status. Port of
	// getShortLived().
	ShortLived() claim.Claim
	// Evidence gets the array of evidence elements. Port of getEvidence().
	Evidence() claim.Array
	// AttestedAttributesSubject gets the claim for associating a set of attributes to one
	// entity different than the EAA subject. Port of getAttestedAttributesSubject().
	AttestedAttributesSubject() claim.AttestedAttributesSubject
}
