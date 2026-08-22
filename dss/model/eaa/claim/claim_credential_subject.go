// Ported from dss-model/.../claim/ClaimCredentialSubject.java (DSS 6.5.RC1).
package claim

// CredentialSubject represents a "4.8 Credential Subject" claim
// defined in W3C Verifiable Credentials Data Model v2.0.
type CredentialSubject interface {
	Claim

	// FullName gets the user's full name information, when present.
	// Ports ClaimCredentialSubject#getFullName.
	FullName() *String

	// GivenName gets the user's first or given name information, when
	// present. Ports ClaimCredentialSubject#getGivenName.
	GivenName() *String

	// FamilyName gets the user's last name or surname information, when
	// present. Ports ClaimCredentialSubject#getFamilyName.
	FamilyName() *String

	// MiddleName gets the user's middle name information, when present.
	// Ports ClaimCredentialSubject#getMiddleName.
	MiddleName() *String

	// Nickname gets the user's casual name information, when present.
	// Ports ClaimCredentialSubject#getNickname.
	Nickname() *String

	// ShortName gets the user's preferred name, usually a shorthand name,
	// when present. Ports ClaimCredentialSubject#getShortName.
	ShortName() *String

	// ProfileUrl gets the user's profile page URL, when present. Ports
	// ClaimCredentialSubject#getProfileUrl.
	ProfileUrl() *String

	// PictureUrl gets the user's profile picture URL, when present.
	// Ports ClaimCredentialSubject#getPictureUrl.
	PictureUrl() *String

	// WebsiteUrl gets the user's website or blog URL, when present.
	// Ports ClaimCredentialSubject#getWebsiteUrl.
	WebsiteUrl() *String

	// Email gets the user's preferred email address, when present. Ports
	// ClaimCredentialSubject#getEmail.
	Email() *String

	// EmailVerified gets whether the user's email address has been
	// verified, when present. Ports
	// ClaimCredentialSubject#getEmailVerified.
	EmailVerified() *Boolean

	// Gender gets the user's gender, when present. Ports
	// ClaimCredentialSubject#getGender.
	Gender() *String

	// Birthdate gets the user's birthdate, when present. Ports
	// ClaimCredentialSubject#getBirthdate.
	Birthdate() *Date

	// Timezone gets the user's TimeZone, when present. Ports
	// ClaimCredentialSubject#getTimezone.
	Timezone() *String

	// Locale gets the user's locale, when present. Ports
	// ClaimCredentialSubject#getLocale.
	Locale() *String

	// Address gets the user's full postal or physical address, when
	// present. Ports ClaimCredentialSubject#getAddress.
	Address() Address

	// PhoneNumber gets the user's preferred telephone number, when
	// present. Ports ClaimCredentialSubject#getPhoneNumber.
	PhoneNumber() *String

	// PhoneNumberVerified gets whether the user's preferred telephone
	// number has been verified, when present. Ports
	// ClaimCredentialSubject#getPhoneNumberVerified.
	PhoneNumberVerified() *Boolean

	// PlaceOfBirth gets user's place of birth, when present. Ports
	// ClaimCredentialSubject#getPlaceOfBirth.
	PlaceOfBirth() PlaceOfBirth

	// Nationalities gets user's nationalities using ICAO 3-letter codes,
	// when present. Ports ClaimCredentialSubject#getNationalities.
	Nationalities() *Array

	// BirthGivenName gets user's first or given name when they were born,
	// when present. Ports ClaimCredentialSubject#getBirthGivenName.
	BirthGivenName() *String

	// BirthFamilyName gets user's family or last name when they were
	// born, when present. Ports
	// ClaimCredentialSubject#getBirthFamilyName.
	BirthFamilyName() *String

	// BirthMiddleName gets user's middle name when they were born, when
	// present. Ports ClaimCredentialSubject#getBirthMiddleName.
	BirthMiddleName() *String

	// Salutation gets user's salutation, e.g., "Mr", when present. Ports
	// ClaimCredentialSubject#getSalutation.
	Salutation() *String

	// Title gets user's title, e.g., "Dr", when present. Ports
	// ClaimCredentialSubject#getTitle.
	Title() *String

	// MobilePhoneNumber gets user's mobile phone number, when present.
	// Ports ClaimCredentialSubject#getMobilePhoneNumber.
	MobilePhoneNumber() *String

	// Pseudonym gets user's stage name, religious name or any other type
	// of alias/pseudonym, when present. Ports
	// ClaimCredentialSubject#getPseudonym.
	Pseudonym() *String
}
