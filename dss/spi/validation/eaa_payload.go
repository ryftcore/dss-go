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
	Identifier() claim.ClaimString
	// Issuer gets the EAA's issuer, when present. Port of getIssuer().
	Issuer() claim.ClaimString
	// Subject gets the EAA's subject, when present. Port of getSubject().
	Subject() claim.ClaimString
	// Audience gets the list of recipients the EAA is intended for, when present. Port of
	// getAudience().
	Audience() claim.ClaimArray
	// IssuedAtTime gets the time at which the EAA was issued, when present. Port of
	// getIssuedAtTime().
	IssuedAtTime() claim.ClaimDate
	// NotBeforeTime gets the time before which the EAA is not accepted for processing, when
	// present. Port of getNotBeforeTime().
	NotBeforeTime() claim.ClaimDate
	// ExpirationTime gets the expiration time of the EAA, after which the EAA is not accepted
	// for processing, when present. Port of getExpirationTime().
	ExpirationTime() claim.ClaimDate
	// UpdatedAtTime gets the time at which the information present within the EAA was the last
	// time updated, when present. Port of getUpdatedAtTime().
	UpdatedAtTime() claim.ClaimDate
	// DeviceKey gets the wallet holder's key. Port of getDeviceKey().
	DeviceKey() claim.ClaimDeviceKey
	// Category gets the EAA category URN, when present. Port of getCategory().
	Category() claim.ClaimString
	// VerifiableCredentialsType gets the EAA's Metadata type, when present. Port of
	// getVerifiableCredentialsType().
	VerifiableCredentialsType() claim.ClaimString
	// VerifiableCredentialsTypeIntegrity gets the EAA's Metadata integrity claim, when present.
	// Port of getVerifiableCredentialsTypeIntegrity().
	VerifiableCredentialsTypeIntegrity() claim.ClaimIntegrity
	// Status gets the EAA's Status value, when present. Port of getStatus().
	Status() claim.ClaimStatus
	// Nonce gets the EAA's nonce value, used to associate the Client's session Id with the EAA,
	// when present. Port of getNonce().
	Nonce() claim.ClaimString
	// FullName gets the user's full name information, when present. Port of getFullName().
	FullName() claim.ClaimString
	// GivenName gets the user's first or given name information, when present. Port of
	// getGivenName().
	GivenName() claim.ClaimString
	// FamilyName gets the user's last name or surname information, when present. Port of
	// getFamilyName().
	FamilyName() claim.ClaimString
	// MiddleName gets the user's middle name information, when present. Port of getMiddleName().
	MiddleName() claim.ClaimString
	// Nickname gets the user's casual name information, when present. Port of getNickname().
	Nickname() claim.ClaimString
	// ShortName gets the user's preferred name, usually a shorthand name, when present. Port of
	// getShortName().
	ShortName() claim.ClaimString
	// ProfileUrl gets the user's profile page URL, when present. Port of getProfileUrl().
	ProfileUrl() claim.ClaimString
	// PictureUrl gets the user's profile picture URL, when present. Port of getPictureUrl().
	PictureUrl() claim.ClaimString
	// WebsiteUrl gets the user's website or blog URL, when present. Port of getWebsiteUrl().
	WebsiteUrl() claim.ClaimString
	// Email gets the user's preferred email address, when present. Port of getEmail().
	Email() claim.ClaimString
	// EmailVerified gets whether the user's email address has been verified, when present. Port
	// of getEmailVerified().
	EmailVerified() claim.ClaimBoolean
	// Gender gets the user's gender, when present. Port of getGender().
	Gender() claim.Claim
	// Birthdate gets the user's birthdate, when present. Port of getBirthdate().
	Birthdate() claim.Claim
	// Timezone gets the user's TimeZone, when present. Port of getTimezone().
	Timezone() claim.ClaimString
	// Locale gets the user's locale, when present. Port of getLocale().
	Locale() claim.ClaimString
	// Address gets the user's full postal or physical address, when present. Port of
	// getAddress().
	Address() claim.ClaimAddress
	// PhoneNumber gets the user's preferred telephone number, when present. Port of
	// getPhoneNumber().
	PhoneNumber() claim.ClaimString
	// PhoneNumberVerified gets whether the user's preferred telephone number has been verified,
	// when present. Port of getPhoneNumberVerified().
	PhoneNumberVerified() claim.ClaimBoolean
	// PlaceOfBirth gets user's place of birth, when present. Port of getPlaceOfBirth().
	PlaceOfBirth() claim.Claim
	// Nationalities gets user's nationalities using ICAO 3-letter codes, when present. Port of
	// getNationalities().
	Nationalities() claim.Claim
	// BirthGivenName gets user's first or given name when they were born, when present. Port of
	// getBirthGivenName().
	BirthGivenName() claim.ClaimString
	// BirthFamilyName gets user's family or last name when they were born, when present. Port
	// of getBirthFamilyName().
	BirthFamilyName() claim.ClaimString
	// BirthMiddleName gets user's middle name when they were born, when present. Port of
	// getBirthMiddleName().
	BirthMiddleName() claim.ClaimString
	// Salutation gets user's salutation, e.g., "Mr", when present. Port of getSalutation().
	Salutation() claim.ClaimString
	// Title gets user's title, e.g., "Dr", when present. Port of getTitle().
	Title() claim.ClaimString
	// MobilePhoneNumber gets user's mobile phone number, when present. Port of
	// getMobilePhoneNumber().
	MobilePhoneNumber() claim.ClaimString
	// Pseudonym gets user's stage name, religious name or any other type of alias/pseudonym,
	// when present. Port of getPseudonym().
	Pseudonym() claim.ClaimString
	// CredentialSubjects returns a list of "4.8 Credential Subject" claims defined in W3C
	// Verifiable Credentials Data Model v2.0. Port of getCredentialSubjects().
	CredentialSubjects() []claim.ClaimCredentialSubject

	/* Mdoc specific payload headers as per ISO/IEC 18013-5 */

	// IssuingCountry gets alpha-2 country code, as defined in ISO 3166-1, of the issuing
	// authority's country or territory. Port of getIssuingCountry().
	IssuingCountry() claim.ClaimString
	// IssuingAuthority gets issuing authority name. The value shall only use latin1 characters
	// and shall have a maximum length of 150 characters. Port of getIssuingAuthority().
	IssuingAuthority() claim.ClaimString
	// DocumentNumber gets the number assigned or calculated by the issuing authority. The value
	// shall only use latin1 characters and shall have a maximum length of 150 characters. Port
	// of getDocumentNumber().
	DocumentNumber() claim.ClaimString
	// Portrait gets a reproduction of the mDL holder's portrait. Port of getPortrait().
	Portrait() claim.ClaimByteString
	// DrivingPrivileges gets driving privileges of the mDL holder. Port of
	// getDrivingPrivileges().
	DrivingPrivileges() claim.ClaimDrivingPrivileges
	// UNDistinguishingSign gets the distinguishing sign of the issuing country according to
	// ISO/IEC 18013-1:2018, Annex F. If no applicable distinguishing sign is available in
	// ISO/IEC 18013-1, an IA may use an empty identifier or another identifier by which it is
	// internationally recognized. In this case the IA should ensure there is no collision with
	// other IA's. Port of getUNDistinguishingSign().
	UNDistinguishingSign() claim.ClaimString
	// PersonalAdministrativeNumber is an audit control number assigned by the issuing
	// authority. The value shall only use latin1 characters and shall have a maximum length of
	// 150 characters. Port of getPersonalAdministrativeNumber().
	PersonalAdministrativeNumber() claim.ClaimString
	// Height gets the holder's height in centimetres. Port of getHeight().
	Height() claim.ClaimNumber
	// Weight gets the holder's weight in kilograms (Java's doc comment repeats "height in
	// centimetres" verbatim; kept as-is - this is a copy-paste doc bug in upstream, not
	// reproduced-then-fixed here per faithful porting). Port of getWeight().
	Weight() claim.ClaimNumber
	// EyeColour gets the mDL holder's eye colour. The value shall be one of the following:
	// "black", "blue", "brown", "dichromatic", "grey", "green", "hazel", "maroon", "pink",
	// "unknown". Port of getEyeColour().
	EyeColour() claim.ClaimString
	// HairColour gets the mDL holder's hair colour. The value shall be one of the following:
	// "bald", "black", "blond", "brown", "grey", "red", "auburn", "sandy", "white", "unknown".
	// Port of getHairColour().
	HairColour() claim.ClaimString
	// PostalAddress gets the place where the mDL holder resides and/or may be contacted
	// (street/house number, municipality etc.). The value shall only use latin1 characters and
	// shall have a maximum length of 150 characters. Port of getPostalAddress().
	PostalAddress() claim.ClaimString
	// PortraitCaptureDate gets the date when portrait was taken. Port of
	// getPortraitCaptureDate().
	PortraitCaptureDate() claim.ClaimDate
	// AgeInYears gets the date the age of the mDL holder. Port of getAgeInYears().
	AgeInYears() claim.ClaimNumber
	// AgeBirthYear gets the year when the mDL holder was born. Port of getAgeBirthYear().
	AgeBirthYear() claim.ClaimNumber
	// AgeEqualOrOver gets a map of elements attesting whether the User to whom the person
	// identification data relates is at least NN years old. N <> 18. Multiple instances of
	// this attribute may be present, provided the value of NN is different in each of them. If
	// present, the requirements in clause 7.2.5 of ISO/IEC 18013-5 are applicable for these
	// attributes. Port of getAgeEqualOrOver().
	AgeEqualOrOver() claim.ClaimAgeEqualOrOver
	// AgeOverNN gets a list of elements used to convey to an mDL verifier, in a data-minimized
	// fashion, if the mDL holder is as old or older than a specified age, or if the mDL holder
	// is younger than a specified age. To achieve this, the mDL contains age attestation
	// identifiers. An age attestation identifier has the format age_over_NN where NN is a
	// value from 00 to 99. The value of an age attestation identifier can be TRUE or FALSE.
	// Port of getAgeOverNN().
	AgeOverNN() []claim.ClaimAgeOverNN
	// IssuingJurisdiction gets a country subdivision code of the jurisdiction that issued the
	// mDL as defined in ISO 3166-2:2020, Clause 8. The first part of the code shall be the same
	// as the value for issuing_country. Port of getIssuingJurisdiction().
	IssuingJurisdiction() claim.ClaimString
	// ResidentAddressCity gets the city where the mDL holder lives. The value shall only use
	// latin1 characters and shall have a maximum length of 150 characters. Port of
	// getResidentAddressCity().
	ResidentAddressCity() claim.ClaimString
	// ResidentAddressState gets the state/province/district where the mDL holder lives. The
	// value shall only use latin1 characters and shall have a maximum length of 150 characters.
	// Port of getResidentAddressState().
	ResidentAddressState() claim.ClaimString
	// ResidentAddressPostalCode gets the postal code of the mDL holder. The value shall only
	// use latin1b characters and shall have a maximum length of 150 characters. Port of
	// getResidentAddressPostalCode().
	ResidentAddressPostalCode() claim.ClaimString
	// ResidentAddressCountry gets the country where the mDL holder lives as a two letter
	// country code (alpha-2 code) defined in ISO 3166-1. Port of getResidentAddressCountry().
	ResidentAddressCountry() claim.ClaimString
	// BiometricTemplate gets a list of elements containing optional facial, fingerprint, iris,
	// or other biometric information of the mDL holder. A biometric template identifier has the
	// format biometric_template_xx where xx shall be replaced with the corresponding "Abstract
	// value name" found in ISO/IEC 19785 3:2020, Table 7, according to the following
	// convention: capitalized characters are replaced with their lowercase equivalent and
	// spaces or non-alphanumeric characters are replaced by underscores (_). Port of
	// getBiometricTemplate().
	BiometricTemplate() []claim.ClaimBiometricTemplateXX
	// SignatureUsualMark gets an image of the signature or usual mark of the mDL holder, see
	// 7.2.7 ISO/IEC 18013-5. Port of getSignatureUsualMark().
	SignatureUsualMark() claim.ClaimByteString

	/* "9.1.2.4 Signing method and structure for MSO" headers as per ISO/IEC 18013-5 */

	// Version gets a version of the MobileSecurityObject. Port of getVersion().
	Version() claim.ClaimString
	// DocType gets a docType as used in Documents. NOTE: This a mandatory non-disclosable
	// property in comparison with DocumentType. Port of getDocType().
	DocType() claim.ClaimString
	// ValidityInfo gets the information related to the validity of the MSO and its signature.
	// Port of getValidityInfo().
	ValidityInfo() claim.ClaimValidityInfo

	/* Mdoc specific payload headers as per ISO/IEC 23220-2 */

	// Fingerprint gets a reproduction of the holder's fingerprint data (TBC). Port of
	// getFingerprint().
	Fingerprint() claim.ClaimByteString
	// BusinessName gets a business name of the holder. Port of getBusinessName().
	BusinessName() claim.ClaimString
	// OrganizationName gets a name of legal person. Port of getOrganizationName().
	OrganizationName() claim.ClaimString
	// BirthFullName gets the name(s) which holder was born. Port of getBirthFullName().
	BirthFullName() claim.ClaimString
	// Profession gets the profession of the holder. Port of getProfession().
	Profession() claim.ClaimString

	/* "6.3.2.3 Relationship attributes" headers as per ISO/IEC 23220-2 */

	// RelationshipFather gets the father of the holder. Port of getRelationshipFather().
	RelationshipFather() claim.ClaimString
	// RelationshipMother gets the mother of the holder. Port of getRelationshipMother().
	RelationshipMother() claim.ClaimString
	// RelationshipParent gets the parent of the holder. Port of getRelationshipParent().
	RelationshipParent() claim.ClaimString
	// RelationshipSon gets the son of the holder. Port of getRelationshipSon().
	RelationshipSon() claim.ClaimString
	// RelationshipDaughter gets the daughter of the holder. Port of getRelationshipDaughter().
	RelationshipDaughter() claim.ClaimString
	// RelationshipBrother gets the brother of the holder. Port of getRelationshipBrother().
	RelationshipBrother() claim.ClaimString
	// RelationshipSister gets the sister of the holder. Port of getRelationshipSister().
	RelationshipSister() claim.ClaimString
	// RelationshipSibling gets the sibling of the holder. Port of getRelationshipSibling().
	RelationshipSibling() claim.ClaimString
	// RelationshipSpouse gets the spouse of the holder. Port of getRelationshipSpouse().
	RelationshipSpouse() claim.ClaimString
	// RelationshipFatherInLaw gets the father-in-law of the holder. Port of
	// getRelationshipFatherInLaw().
	RelationshipFatherInLaw() claim.ClaimString
	// RelationshipMotherInLaw gets the mother-in-law of the holder. Port of
	// getRelationshipMotherInLaw().
	RelationshipMotherInLaw() claim.ClaimString
	// RelationshipParentInLaw gets the parent-in-law of the holder. Port of
	// getRelationshipParentInLaw().
	RelationshipParentInLaw() claim.ClaimString
	// RelationshipSonInLaw gets the son-in-law of the holder. Port of
	// getRelationshipSonInLaw().
	RelationshipSonInLaw() claim.ClaimString
	// RelationshipDaughterInLaw gets the daughter-in-law of the holder. Port of
	// getRelationshipDaughterInLaw().
	RelationshipDaughterInLaw() claim.ClaimString
	// RelationshipChildInLaw gets the child-in-law of the holder. Port of
	// getRelationshipChildInLaw().
	RelationshipChildInLaw() claim.ClaimString
	// RelationshipParentalAuthority gets the parental authority of the holder. Port of
	// getRelationshipParentalAuthority().
	RelationshipParentalAuthority() claim.ClaimString
	// RelationshipLegalRepresentative gets the legal representative of the holder. Port of
	// getRelationshipLegalRepresentative().
	RelationshipLegalRepresentative() claim.ClaimString
	// RelationshipAgent gets the voluntary agent of the holder. Port of
	// getRelationshipAgent().
	RelationshipAgent() claim.ClaimString

	/* "6.3.4 Data elements for document entity" headers as per ISO/IEC 23220-2 */

	// DocumentType gets the document type. NOTE: This a selectively disclosable property in
	// comparison with DocType. Port of getDocumentType().
	DocumentType() claim.ClaimString

	/* ARF PID Rulebook headers */

	// AdministrativeIssuanceDate gets the date when the data (e.g. a PID) was issued. Port of
	// getAdministrativeIssuanceDate().
	AdministrativeIssuanceDate() claim.ClaimDate
	// AdministrativeExpirationDate gets the date when the data (e.g. a PID) will expire. Port
	// of getAdministrativeExpirationDate().
	AdministrativeExpirationDate() claim.ClaimDate
	// TrustAnchor gets the URL at which a machine-readable version of the trust anchor to be
	// used for verifying the PID can be found or looked up. Port of getTrustAnchor().
	TrustAnchor() claim.ClaimString
	// ResidentAddressStreet gets the name of the street where the user to whom the person
	// identification data relates currently resides. Port of getResidentAddressStreet().
	ResidentAddressStreet() claim.ClaimString
	// ResidentAddressHouseNumber gets the house number where the user to whom the person
	// identification data relates currently resides, including any affix or suffix. Port of
	// getResidentAddressHouseNumber().
	ResidentAddressHouseNumber() claim.ClaimString

	/* ETSI TS 119 472-1 "5 Implementation of EAA based on SD-JWT VC" header parameters */

	// IssuingAuthorityRegistrationIdentifier gets the registration identifier of the legal
	// entity on whose behalf the EAA has been issued. Port of
	// getIssuingAuthorityRegistrationIdentifier().
	IssuingAuthorityRegistrationIdentifier() claim.ClaimString
	// OneTimeUse gets the signal indicating that the EAA shall be used only once, and that it
	// shall not be retained for future use. Port of getOneTimeUse().
	OneTimeUse() claim.Claim
	// ShortLived gets the EAA short-lived component indicating that the validity period of the
	// EAA is so short that it shall not be necessary to check its revocation status. Port of
	// getShortLived().
	ShortLived() claim.Claim
	// Evidence gets the array of evidence elements. Port of getEvidence().
	Evidence() claim.ClaimArray
	// AttestedAttributesSubject gets the claim for associating a set of attributes to one
	// entity different than the EAA subject. Port of getAttestedAttributesSubject().
	AttestedAttributesSubject() claim.ClaimAttestedAttributesSubject
}
