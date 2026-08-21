// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/CredentialSubjectProxy.java (DSS 6.5.RC1).
//
// Provides an NPE-safe initialization and always returns the first credential subject value,
// when applicable.
package diagnostic

import "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

// CredentialSubjectProxy provides an NPE-safe initialization and returns always the first
// credential subject value, when applicable.
type CredentialSubjectProxy struct {
	// xmlCredentialSubjectList is the wrapped list of credential subjects.
	xmlCredentialSubjectList []*jaxb.XmlCredentialSubjectClaim
}

// NewCredentialSubjectProxy is the default constructor. Port of
// CredentialSubjectProxy(List<XmlCredentialSubjectClaim>).
func NewCredentialSubjectProxy(xmlCredentialSubjectList []*jaxb.XmlCredentialSubjectClaim) *CredentialSubjectProxy {
	return &CredentialSubjectProxy{xmlCredentialSubjectList: xmlCredentialSubjectList}
}

// CredentialSubjects gets a list of credential subjects. Port of getCredentialSubjects().
func (p *CredentialSubjectProxy) CredentialSubjects() []*CredentialSubjectClaimWrapper {
	if len(p.xmlCredentialSubjectList) == 0 {
		return nil
	}
	result := make([]*CredentialSubjectClaimWrapper, 0, len(p.xmlCredentialSubjectList))
	for _, x := range p.xmlCredentialSubjectList {
		result = append(result, NewCredentialSubjectClaimWrapper(x))
	}
	return result
}

// FirstCredentialSubject gets the first credential subject, when defined. Returns nil
// otherwise. Port of getFirstCredentialSubject().
func (p *CredentialSubjectProxy) FirstCredentialSubject() *CredentialSubjectClaimWrapper {
	if len(p.xmlCredentialSubjectList) == 0 {
		return nil
	}
	return NewCredentialSubjectClaimWrapper(p.xmlCredentialSubjectList[0])
}

// FullName gets full name when defined within the first Credential Subject claim. Port of
// getFullName().
func (p *CredentialSubjectProxy) FullName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.FullName()
	}
	return nil
}

// GivenName gets first name when defined within the first Credential Subject claim. Port of
// getGivenName().
func (p *CredentialSubjectProxy) GivenName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.GivenName()
	}
	return nil
}

// FamilyName gets last or family name when defined within the first Credential Subject claim.
// Port of getFamilyName().
func (p *CredentialSubjectProxy) FamilyName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.FamilyName()
	}
	return nil
}

// MiddleName gets middle name when defined within the first Credential Subject claim. Port of
// getMiddleName().
func (p *CredentialSubjectProxy) MiddleName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.MiddleName()
	}
	return nil
}

// Nickname gets alternative name when defined within the first Credential Subject claim. Port
// of getNickname().
func (p *CredentialSubjectProxy) Nickname() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Nickname()
	}
	return nil
}

// ShortName gets preferred or short name when defined within the first Credential Subject
// claim. Port of getShortName().
func (p *CredentialSubjectProxy) ShortName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.ShortName()
	}
	return nil
}

// ProfileUrl gets profile URL when defined within the first Credential Subject claim. Port of
// getProfileUrl().
func (p *CredentialSubjectProxy) ProfileUrl() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.ProfileUrl()
	}
	return nil
}

// PictureUrl gets picture URL when defined within the first Credential Subject claim. Port of
// getPictureUrl().
func (p *CredentialSubjectProxy) PictureUrl() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.PictureUrl()
	}
	return nil
}

// WebsiteUrl gets website when defined within the first Credential Subject claim. Port of
// getWebsiteUrl().
func (p *CredentialSubjectProxy) WebsiteUrl() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.WebsiteUrl()
	}
	return nil
}

// Email gets email when defined within the first Credential Subject claim. Port of getEmail().
func (p *CredentialSubjectProxy) Email() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Email()
	}
	return nil
}

// EmailVerified gets whether the website has been verified if defined within the first
// Credential Subject claim. Port of getEmailVerified().
func (p *CredentialSubjectProxy) EmailVerified() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.EmailVerified()
	}
	return nil
}

// Gender gets gender when defined within the first Credential Subject claim. Port of
// getGender().
func (p *CredentialSubjectProxy) Gender() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Gender()
	}
	return nil
}

// Birthdate gets birthdate when defined within the first Credential Subject claim. Port of
// getBirthdate().
func (p *CredentialSubjectProxy) Birthdate() *BirthdateClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Birthdate()
	}
	return nil
}

// Timezone gets timezone when defined within the first Credential Subject claim. Port of
// getTimezone().
func (p *CredentialSubjectProxy) Timezone() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Timezone()
	}
	return nil
}

// Locale gets locale when defined within the first Credential Subject claim. Port of
// getLocale().
func (p *CredentialSubjectProxy) Locale() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Locale()
	}
	return nil
}

// Address gets full address, when defined within the first Credential Subject claim. Port of
// getAddress().
func (p *CredentialSubjectProxy) Address() *AddressClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Address()
	}
	return nil
}

// AddressCity gets city address when defined within the first Credential Subject claim. Port
// of getAddressCity().
func (p *CredentialSubjectProxy) AddressCity() *ClaimWrapper {
	if address := p.Address(); address != nil {
		return address.City()
	}
	return nil
}

// AddressStateOrProvince gets state or region address when defined within the first Credential
// Subject claim. Port of getAddressStateOrProvince().
func (p *CredentialSubjectProxy) AddressStateOrProvince() *ClaimWrapper {
	if address := p.Address(); address != nil {
		return address.StateOrProvince()
	}
	return nil
}

// AddressPostalCode gets postal code address when defined within the first Credential Subject
// claim. Port of getAddressPostalCode().
func (p *CredentialSubjectProxy) AddressPostalCode() *ClaimWrapper {
	if address := p.Address(); address != nil {
		return address.PostalCode()
	}
	return nil
}

// AddressCountry gets country address when defined within the first Credential Subject claim.
// NOTE: the returned value is usually represented by 2-letter ISO country code. Port of
// getAddressCountry().
func (p *CredentialSubjectProxy) AddressCountry() *ClaimWrapper {
	if address := p.Address(); address != nil {
		return address.Country()
	}
	return nil
}

// StreetAddress gets street address when defined within the first Credential Subject claim.
// Port of getStreetAddress().
func (p *CredentialSubjectProxy) StreetAddress() *ClaimWrapper {
	if address := p.Address(); address != nil {
		return address.StreetAddress()
	}
	return nil
}

// PhoneNumber gets phone number when defined within the first Credential Subject claim. Port
// of getPhoneNumber().
func (p *CredentialSubjectProxy) PhoneNumber() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.PhoneNumber()
	}
	return nil
}

// PhoneNumberVerified gets whether the phone number has been verified if defined within the
// first Credential Subject claim. Port of getPhoneNumberVerified().
func (p *CredentialSubjectProxy) PhoneNumberVerified() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.PhoneNumberVerified()
	}
	return nil
}

// PlaceOfBirth gets place of birth when defined within the first Credential Subject claim.
// Port of getPlaceOfBirth().
func (p *CredentialSubjectProxy) PlaceOfBirth() *PlaceOfBirthClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.PlaceOfBirth()
	}
	return nil
}

// Nationalities gets nationalities list when defined within the first Credential Subject
// claim. NOTE: the values are usually represented by 3-letter nationality codes. Port of
// getNationalities().
func (p *CredentialSubjectProxy) Nationalities() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Nationalities()
	}
	return nil
}

// BirthFamilyName gets last or family name at birth when defined within the first Credential
// Subject claim. Port of getBirthFamilyName().
func (p *CredentialSubjectProxy) BirthFamilyName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.BirthFamilyName()
	}
	return nil
}

// BirthGivenName gets first name at birth when defined within the first Credential Subject
// claim. Port of getBirthGivenName().
func (p *CredentialSubjectProxy) BirthGivenName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.BirthGivenName()
	}
	return nil
}

// BirthMiddleName gets middle name at birth when defined within the first Credential Subject
// claim. Port of getBirthMiddleName().
func (p *CredentialSubjectProxy) BirthMiddleName() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.BirthMiddleName()
	}
	return nil
}

// Salutation gets preferred salutation when defined within the first Credential Subject claim.
// Port of getSalutation().
func (p *CredentialSubjectProxy) Salutation() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Salutation()
	}
	return nil
}

// Title gets title when defined within the first Credential Subject claim. Port of getTitle().
func (p *CredentialSubjectProxy) Title() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Title()
	}
	return nil
}

// MobilePhoneNumber gets mobile phone number when defined within the first Credential Subject
// claim. Port of getMobilePhoneNumber().
func (p *CredentialSubjectProxy) MobilePhoneNumber() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.MobilePhoneNumber()
	}
	return nil
}

// Pseudonym gets scenic name or pseudonym, they are known as, when defined within the first
// Credential Subject claim. Port of getPseudonym().
func (p *CredentialSubjectProxy) Pseudonym() *ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.Pseudonym()
	}
	return nil
}

// OtherClaims gets a list of claims incorporated within the Credential Subject or provided as
// disclosures, which are not (yet) directly supported by the implementation. Port of
// getOtherClaims().
func (p *CredentialSubjectProxy) OtherClaims() []*ClaimWrapper {
	if credentialSubject := p.FirstCredentialSubject(); credentialSubject != nil {
		return credentialSubject.OtherClaims()
	}
	return nil
}
