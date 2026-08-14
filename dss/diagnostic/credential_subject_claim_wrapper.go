// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/CredentialSubjectClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// CredentialSubjectClaimWrapper wraps a jaxb.XmlCredentialSubjectClaim. Unlike most claim
// subtype wrappers, Java declares only the single-argument constructor here (no parent-taking
// overload).
type CredentialSubjectClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlCredentialSubjectClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlCredentialSubjectClaim
}

// NewCredentialSubjectClaimWrapper is the default constructor. Port of
// CredentialSubjectClaimWrapper(XmlCredentialSubjectClaim).
func NewCredentialSubjectClaimWrapper(wrapped *jaxb.XmlCredentialSubjectClaim) *CredentialSubjectClaimWrapper {
	return &CredentialSubjectClaimWrapper{
		ClaimWrapper: *NewClaimWrapper(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs)),
		wrapped:      wrapped,
	}
}

func (w *CredentialSubjectClaimWrapper) claim(xmlClaim *jaxb.XmlClaim) *ClaimWrapper {
	if xmlClaim != nil {
		return NewClaimWrapperWithParent(xmlClaim, &w.ClaimWrapper)
	}
	return nil
}

// FullName gets user's full name when defined within Credential Subject claim. Port of getFullName().
func (w *CredentialSubjectClaimWrapper) FullName() *ClaimWrapper { return w.claim(w.wrapped.FullName) }

// GivenName gets user's first name when defined within Credential Subject claim. Port of getGivenName().
func (w *CredentialSubjectClaimWrapper) GivenName() *ClaimWrapper {
	return w.claim(w.wrapped.GivenName)
}

// FamilyName gets user's last or family name when defined within Credential Subject claim. Port
// of getFamilyName().
func (w *CredentialSubjectClaimWrapper) FamilyName() *ClaimWrapper {
	return w.claim(w.wrapped.FamilyName)
}

// MiddleName gets user's middle name when defined within Credential Subject claim. Port of
// getMiddleName().
func (w *CredentialSubjectClaimWrapper) MiddleName() *ClaimWrapper {
	return w.claim(w.wrapped.MiddleName)
}

// Nickname gets user's alternative name when defined within Credential Subject claim. Port of
// getNickname().
func (w *CredentialSubjectClaimWrapper) Nickname() *ClaimWrapper { return w.claim(w.wrapped.Nickname) }

// ShortName gets user's preferred or short name when defined within Credential Subject claim.
// Port of getShortName().
func (w *CredentialSubjectClaimWrapper) ShortName() *ClaimWrapper {
	return w.claim(w.wrapped.ShortName)
}

// ProfileUrl gets user's profile URL when defined within Credential Subject claim. Port of
// getProfileUrl().
func (w *CredentialSubjectClaimWrapper) ProfileUrl() *ClaimWrapper {
	return w.claim(w.wrapped.ProfileUrl)
}

// PictureUrl gets user's picture URL when defined within Credential Subject claim. Port of
// getPictureUrl().
func (w *CredentialSubjectClaimWrapper) PictureUrl() *ClaimWrapper {
	return w.claim(w.wrapped.PictureUrl)
}

// WebsiteUrl gets user's website when defined within Credential Subject claim. Port of
// getWebsiteUrl().
func (w *CredentialSubjectClaimWrapper) WebsiteUrl() *ClaimWrapper {
	return w.claim(w.wrapped.WebsiteUrl)
}

// Email gets user's email when defined within Credential Subject claim. Port of getEmail().
func (w *CredentialSubjectClaimWrapper) Email() *ClaimWrapper { return w.claim(w.wrapped.Email) }

// EmailVerified gets whether the user's website has been verified if defined within Credential
// Subject claim. Port of getEmailVerified().
func (w *CredentialSubjectClaimWrapper) EmailVerified() *ClaimWrapper {
	return w.claim(w.wrapped.EmailVerified)
}

// Gender gets user's gender when defined within Credential Subject claim. Port of getGender().
func (w *CredentialSubjectClaimWrapper) Gender() *ClaimWrapper { return w.claim(w.wrapped.Gender) }

// Birthdate gets user's birthdate when defined within Credential Subject claim. Port of
// getBirthdate().
func (w *CredentialSubjectClaimWrapper) Birthdate() *BirthdateClaimWrapper {
	if w.wrapped.Birthdate != nil {
		return NewBirthdateClaimWrapperWithParent(w.wrapped.Birthdate, &w.ClaimWrapper)
	}
	return nil
}

// Timezone gets user's timezone when defined within Credential Subject claim. Port of
// getTimezone().
func (w *CredentialSubjectClaimWrapper) Timezone() *ClaimWrapper { return w.claim(w.wrapped.Timezone) }

// Locale gets user's locale when defined within Credential Subject claim. Port of getLocale().
func (w *CredentialSubjectClaimWrapper) Locale() *ClaimWrapper { return w.claim(w.wrapped.Locale) }

// Address gets user's full address, when defined within Credential Subject claim. Port of
// getAddress().
func (w *CredentialSubjectClaimWrapper) Address() *AddressClaimWrapper {
	if w.wrapped.Address != nil {
		return NewAddressClaimWrapperWithParent(w.wrapped.Address, &w.ClaimWrapper)
	}
	return nil
}

// PhoneNumber gets user's phone number when defined within Credential Subject claim. Port of
// getPhoneNumber().
func (w *CredentialSubjectClaimWrapper) PhoneNumber() *ClaimWrapper {
	return w.claim(w.wrapped.PhoneNumber)
}

// PhoneNumberVerified gets whether the user's phone number has been verified if defined within
// Credential Subject claim. Port of getPhoneNumberVerified().
func (w *CredentialSubjectClaimWrapper) PhoneNumberVerified() *ClaimWrapper {
	return w.claim(w.wrapped.PhoneNumberVerified)
}

// PlaceOfBirth gets user's country of birth when defined within Credential Subject claim. Port
// of getPlaceOfBirth().
func (w *CredentialSubjectClaimWrapper) PlaceOfBirth() *PlaceOfBirthClaimWrapper {
	if w.wrapped.PlaceOfBirth != nil {
		return newPlaceOfBirthClaimWrapperFromClaim(w.wrapped.PlaceOfBirth, &w.ClaimWrapper)
	}
	return nil
}

// Nationalities gets user's nationalities list when defined within Credential Subject claim.
// NOTE: the values are usually represented by 3-letter nationality codes. Port of
// getNationalities().
func (w *CredentialSubjectClaimWrapper) Nationalities() *ClaimWrapper {
	return w.claim(w.wrapped.Nationalities)
}

// BirthFamilyName gets user's last or family name at birth when defined within Credential
// Subject claim. Port of getBirthFamilyName().
func (w *CredentialSubjectClaimWrapper) BirthFamilyName() *ClaimWrapper {
	return w.claim(w.wrapped.BirthFamilyName)
}

// BirthGivenName gets user's first name at birth when defined within Credential Subject claim.
// Port of getBirthGivenName().
func (w *CredentialSubjectClaimWrapper) BirthGivenName() *ClaimWrapper {
	return w.claim(w.wrapped.BirthGivenName)
}

// BirthMiddleName gets user's middle name at birth when defined within Credential Subject
// claim. Port of getBirthMiddleName().
func (w *CredentialSubjectClaimWrapper) BirthMiddleName() *ClaimWrapper {
	return w.claim(w.wrapped.BirthMiddleName)
}

// Salutation gets user's preferred salutation when defined within Credential Subject claim.
// Port of getSalutation().
func (w *CredentialSubjectClaimWrapper) Salutation() *ClaimWrapper {
	return w.claim(w.wrapped.Salutation)
}

// Title gets user's title when defined within Credential Subject claim. Port of getTitle().
func (w *CredentialSubjectClaimWrapper) Title() *ClaimWrapper { return w.claim(w.wrapped.Title) }

// MobilePhoneNumber gets user's mobile phone number when defined within Credential Subject
// claim. Port of getMobilePhoneNumber().
func (w *CredentialSubjectClaimWrapper) MobilePhoneNumber() *ClaimWrapper {
	return w.claim(w.wrapped.MobilePhoneNumber)
}

// Pseudonym gets user's scenic name or pseudonym, they are known as, when defined within
// Credential Subject claim. Port of getPseudonym().
func (w *CredentialSubjectClaimWrapper) Pseudonym() *ClaimWrapper {
	return w.claim(w.wrapped.Pseudonym)
}

// AllCredentialSubjectClaims gets a list of all claims present with the Credential Subject
// claim. Port of getAllCredentialSubjectClaims().
func (w *CredentialSubjectClaimWrapper) AllCredentialSubjectClaims() []*ClaimWrapper {
	var claimList []*ClaimWrapper
	if fullName := w.FullName(); fullName != nil {
		claimList = append(claimList, fullName)
	}
	if givenName := w.GivenName(); givenName != nil {
		claimList = append(claimList, givenName)
	}
	if familyName := w.FamilyName(); familyName != nil {
		claimList = append(claimList, familyName)
	}
	if middleName := w.MiddleName(); middleName != nil {
		claimList = append(claimList, middleName)
	}
	if nickname := w.Nickname(); nickname != nil {
		claimList = append(claimList, nickname)
	}
	if shortName := w.ShortName(); shortName != nil {
		claimList = append(claimList, shortName)
	}
	if profileUrl := w.ProfileUrl(); profileUrl != nil {
		claimList = append(claimList, profileUrl)
	}
	if pictureUrl := w.PictureUrl(); pictureUrl != nil {
		claimList = append(claimList, pictureUrl)
	}
	if websiteUrl := w.WebsiteUrl(); websiteUrl != nil {
		claimList = append(claimList, websiteUrl)
	}
	if email := w.Email(); email != nil {
		claimList = append(claimList, email)
	}
	if emailVerified := w.EmailVerified(); emailVerified != nil {
		claimList = append(claimList, emailVerified)
	}
	if gender := w.Gender(); gender != nil {
		claimList = append(claimList, gender)
	}
	if birthdate := w.Birthdate(); birthdate != nil {
		claimList = append(claimList, birthdate.AsClaim())
	}
	if timezone := w.Timezone(); timezone != nil {
		claimList = append(claimList, timezone)
	}
	if locale := w.Locale(); locale != nil {
		claimList = append(claimList, locale)
	}
	if address := w.Address(); address != nil {
		claimList = append(claimList, address.AsClaim())
	}
	if phoneNumber := w.PhoneNumber(); phoneNumber != nil {
		claimList = append(claimList, phoneNumber)
	}
	if phoneNumberVerified := w.PhoneNumberVerified(); phoneNumberVerified != nil {
		claimList = append(claimList, phoneNumberVerified)
	}
	if placeOfBirth := w.PlaceOfBirth(); placeOfBirth != nil {
		claimList = append(claimList, placeOfBirth.AsClaim())
	}
	if nationalities := w.Nationalities(); nationalities != nil {
		claimList = append(claimList, nationalities)
	}
	if birthFamilyName := w.BirthFamilyName(); birthFamilyName != nil {
		claimList = append(claimList, birthFamilyName)
	}
	if birthGivenName := w.BirthGivenName(); birthGivenName != nil {
		claimList = append(claimList, birthGivenName)
	}
	if birthMiddleName := w.BirthMiddleName(); birthMiddleName != nil {
		claimList = append(claimList, birthMiddleName)
	}
	if salutation := w.Salutation(); salutation != nil {
		claimList = append(claimList, salutation)
	}
	if title := w.Title(); title != nil {
		claimList = append(claimList, title)
	}
	if mobilePhoneNumber := w.MobilePhoneNumber(); mobilePhoneNumber != nil {
		claimList = append(claimList, mobilePhoneNumber)
	}
	if pseudonym := w.Pseudonym(); pseudonym != nil {
		claimList = append(claimList, pseudonym)
	}
	claimList = append(claimList, w.OtherClaims()...)
	return claimList
}

// OtherClaims gets a list of claims incorporated within the Credential Subject or provided as
// disclosures, which are not (yet) directly supported by the implementation. Port of
// getOtherClaims().
func (w *CredentialSubjectClaimWrapper) OtherClaims() []*ClaimWrapper {
	if w.wrapped.OtherClaim == nil {
		return nil
	}
	result := make([]*ClaimWrapper, 0, len(w.wrapped.OtherClaim))
	for _, c := range w.wrapped.OtherClaim {
		result = append(result, NewClaimWrapperWithParent(c, &w.ClaimWrapper))
	}
	return result
}

// IsMap is the override: a CredentialSubjectClaimWrapper is unconditionally a map claim. Port
// of the overridden isMap().
func (w *CredentialSubjectClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the named credential-subject fields rather than
// the generic Entry list. Port of the overridden getMap().
func (w *CredentialSubjectClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if fullName := w.FullName(); fullName != nil {
		result[fullName.Name()] = fullName
	}
	if givenName := w.GivenName(); givenName != nil {
		result[givenName.Name()] = givenName
	}
	if familyName := w.FamilyName(); familyName != nil {
		result[familyName.Name()] = familyName
	}
	if middleName := w.MiddleName(); middleName != nil {
		result[middleName.Name()] = middleName
	}
	if nickname := w.Nickname(); nickname != nil {
		result[nickname.Name()] = nickname
	}
	if shortName := w.ShortName(); shortName != nil {
		result[shortName.Name()] = shortName
	}
	if profileUrl := w.ProfileUrl(); profileUrl != nil {
		result[profileUrl.Name()] = profileUrl
	}
	if pictureUrl := w.PictureUrl(); pictureUrl != nil {
		result[pictureUrl.Name()] = pictureUrl
	}
	if websiteUrl := w.WebsiteUrl(); websiteUrl != nil {
		result[websiteUrl.Name()] = websiteUrl
	}
	if email := w.Email(); email != nil {
		result[email.Name()] = email
	}
	if emailVerified := w.EmailVerified(); emailVerified != nil {
		result[emailVerified.Name()] = emailVerified
	}
	if gender := w.Gender(); gender != nil {
		result[gender.Name()] = gender
	}
	if birthdate := w.Birthdate(); birthdate != nil {
		result[birthdate.Name()] = birthdate.AsClaim()
	}
	if timezone := w.Timezone(); timezone != nil {
		result[timezone.Name()] = timezone
	}
	if locale := w.Locale(); locale != nil {
		result[locale.Name()] = locale
	}
	if address := w.Address(); address != nil {
		result[address.Name()] = address.AsClaim()
	}
	if phoneNumber := w.PhoneNumber(); phoneNumber != nil {
		result[phoneNumber.Name()] = phoneNumber
	}
	if phoneNumberVerified := w.PhoneNumberVerified(); phoneNumberVerified != nil {
		result[phoneNumberVerified.Name()] = phoneNumberVerified
	}
	if placeOfBirth := w.PlaceOfBirth(); placeOfBirth != nil {
		result[placeOfBirth.Name()] = placeOfBirth.AsClaim()
	}
	if nationalities := w.Nationalities(); nationalities != nil {
		result[nationalities.Name()] = nationalities
	}
	if birthFamilyName := w.BirthFamilyName(); birthFamilyName != nil {
		result[birthFamilyName.Name()] = birthFamilyName
	}
	if birthGivenName := w.BirthGivenName(); birthGivenName != nil {
		result[birthGivenName.Name()] = birthGivenName
	}
	if birthMiddleName := w.BirthMiddleName(); birthMiddleName != nil {
		result[birthMiddleName.Name()] = birthMiddleName
	}
	if salutation := w.Salutation(); salutation != nil {
		result[salutation.Name()] = salutation
	}
	if title := w.Title(); title != nil {
		result[title.Name()] = title
	}
	if mobilePhoneNumber := w.MobilePhoneNumber(); mobilePhoneNumber != nil {
		result[mobilePhoneNumber.Name()] = mobilePhoneNumber
	}
	if pseudonym := w.Pseudonym(); pseudonym != nil {
		result[pseudonym.Name()] = pseudonym
	}
	for _, otherClaim := range w.OtherClaims() {
		if otherClaim != nil {
			result[otherClaim.Name()] = otherClaim
		}
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *CredentialSubjectClaimWrapper) Wrapped() *jaxb.XmlCredentialSubjectClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *CredentialSubjectClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
