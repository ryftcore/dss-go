// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/EAAPayloadProxy.java
// (DSS 6.5.RC1).
//
// Depends on the `claim` Java subpackage (ClaimWrapper, AddressClaimWrapper,
// BirthdateClaimWrapper, PlaceOfBirthClaimWrapper, DeviceKeyClaimWrapper, ValidityInfoClaimWrapper,
// IntegrityClaimWrapper, StatusClaimWrapper, DrivingPrivilegesClaimWrapper,
// AttestedAttributesSubjectClaimWrapper, AgeEqualOrOverClaimWrapper, AgeOverNNClaimWrapper,
// BiometricTemplateXXClaimWrapper, CredentialSubjectProxy, CredentialSubjectClaimWrapper), which
// S8A_BRIEF.md flattens into this same Go package `diagnostic` ("wrappers + `claim` (mutual
// imports, collision-checked) -> ONE pkg diagnostic"); those types are assigned to a sibling
// chunk and are referenced here unqualified, matching that flattening (see the same note in
// eaa_wrapper.go). CredentialSubjectProxy is assumed to expose accessors named after the Java
// getters (FullName(), GivenName(), Birthdate(), Address(), PlaceOfBirth(), ...), each returning
// the same wrapper type as the corresponding EAAPayloadProxy claim (nil when absent). Every Claim
// subtype wrapper is assumed to expose AsClaim() *ClaimWrapper to view itself as its Java base
// type ClaimWrapper (Go has no struct-inheritance upcast for heterogeneous []*ClaimWrapper
// slices/appends the way Java's covariant generics do); this mirrors the AsClaim() convention
// eaa_wrapper.go already relies on for AgeOverNNClaimWrapper/PlaceOfBirthClaimWrapper/
// BiometricTemplateXXClaimWrapper.
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// EAAPayloadProxy provides an interface for selectively disposable claims extraction.
type EAAPayloadProxy struct {
	// xmlEAAPayload is the wrapped EAA Payload to get access to.
	xmlEAAPayload *jaxb.XmlEAAPayload
}

// NewEAAPayloadProxy is the default constructor.
func NewEAAPayloadProxy(xmlEAAPayload *jaxb.XmlEAAPayload) *EAAPayloadProxy {
	return &EAAPayloadProxy{xmlEAAPayload: xmlEAAPayload}
}

// EAAIdentifier gets EAA identifier provided in the EAA payload. Port of getEAAIdentifier().
func (p *EAAPayloadProxy) EAAIdentifier() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Identifier)
	}
	return nil
}

// EAAIssuer gets EAA issuer as defined in the EAA payload. Port of getEAAIssuer().
func (p *EAAPayloadProxy) EAAIssuer() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Issuer)
	}
	return nil
}

// EAASubject gets EAA subject as defined in the EAA payload. Port of getEAASubject().
func (p *EAAPayloadProxy) EAASubject() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Subject)
	}
	return nil
}

// EAAAudience gets EAA audience as defined in the EAA payload. Port of getEAAAudience().
func (p *EAAPayloadProxy) EAAAudience() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Audience)
	}
	return nil
}

// EAAIssuedAt gets EAA issuance time as defined in the EAA payload. Port of getEAAIssuedAt().
func (p *EAAPayloadProxy) EAAIssuedAt() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.IssuedAt)
	}
	return nil
}

// EAANotBefore gets EAA not before time as defined in the EAA payload. Port of getEAANotBefore().
func (p *EAAPayloadProxy) EAANotBefore() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.NotBefore)
	}
	return nil
}

// EAAExpiration gets EAA expiration time as defined in the EAA payload. Port of getEAAExpiration().
func (p *EAAPayloadProxy) EAAExpiration() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Expiration)
	}
	return nil
}

// EAAUpdatedAt gets EAA update time as defined in the EAA payload. Port of getEAAUpdatedAt().
func (p *EAAPayloadProxy) EAAUpdatedAt() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.UpdatedAt)
	}
	return nil
}

// EAACategory gets category URN provided in the EAA payload. Port of getEAACategory().
func (p *EAAPayloadProxy) EAACategory() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Category)
	}
	return nil
}

// EAAVerifiableCredentialsType gets EAA metadata type (e.g. 'vct' claim) as defined in the EAA
// payload. Port of getEAAVerifiableCredentialsType().
func (p *EAAPayloadProxy) EAAVerifiableCredentialsType() *ClaimWrapper {
	if p.xmlEAAPayload != nil && p.xmlEAAPayload.VerifiableCredentialsType != nil {
		return p.getClaim(claimBase(p.xmlEAAPayload.VerifiableCredentialsType.XmlClaimContent, p.xmlEAAPayload.VerifiableCredentialsType.XmlClaimAttrs))
	}
	return nil
}

// EAAVerifiableCredentialsTypeIntegrity gets the integrity material for the EAA metadata (when
// present). Port of getEAAVerifiableCredentialsTypeIntegrity().
func (p *EAAPayloadProxy) EAAVerifiableCredentialsTypeIntegrity() *IntegrityClaimWrapper {
	if p.xmlEAAPayload != nil && p.xmlEAAPayload.VerifiableCredentialsType != nil {
		return p.getIntegrityClaim(p.xmlEAAPayload.VerifiableCredentialsType.Integrity)
	}
	return nil
}

// EAAStatus gets EAA status as defined in the EAA payload. Port of getEAAStatus().
func (p *EAAPayloadProxy) EAAStatus() *StatusClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getStatusClaim(p.xmlEAAPayload.Status)
	}
	return nil
}

// EAANonce gets EAA nonce when defined in the EAA payload. Port of getEAANonce().
func (p *EAAPayloadProxy) EAANonce() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Nonce)
	}
	return nil
}

// EAADeviceKey gets EAA device key when defined in the EAA payload. Port of getEAADeviceKey().
func (p *EAAPayloadProxy) EAADeviceKey() *DeviceKeyClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getDeviceKeyClaim(p.xmlEAAPayload.DeviceKey)
	}
	return nil
}

// EAAVersion gets a version of the MobileSecurityObject. Port of getEAAVersion().
func (p *EAAPayloadProxy) EAAVersion() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Version)
	}
	return nil
}

// EAADocType gets a docType as used in Documents. NOTE: This a mandatory non-disclosable
// property in comparison with DocumentType. Port of getEAADocType().
func (p *EAAPayloadProxy) EAADocType() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.DocType)
	}
	return nil
}

// EAAValidityInfo gets the information related to the validity of the MSO and its signature.
// Port of getEAAValidityInfo().
func (p *EAAPayloadProxy) EAAValidityInfo() *ValidityInfoClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getValidityInfoClaim(p.xmlEAAPayload.ValidityInfo)
	}
	return nil
}

// HolderFullName gets holder's full name when defined within EAA Payload claims. Port of
// getHolderFullName().
func (p *EAAPayloadProxy) HolderFullName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.FullName), p.getCredentialSubject().FullName())
	}
	return nil
}

// HolderGivenName gets holder's first name when defined within EAA Payload claims. Port of
// getHolderGivenName().
func (p *EAAPayloadProxy) HolderGivenName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.GivenName), p.getCredentialSubject().GivenName())
	}
	return nil
}

// HolderFamilyName gets holder's last or family name when defined within EAA Payload claims.
// Port of getHolderFamilyName().
func (p *EAAPayloadProxy) HolderFamilyName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.FamilyName), p.getCredentialSubject().FamilyName())
	}
	return nil
}

// HolderMiddleName gets holder's middle name when defined within EAA Payload claims. Port of
// getHolderMiddleName().
func (p *EAAPayloadProxy) HolderMiddleName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.MiddleName), p.getCredentialSubject().MiddleName())
	}
	return nil
}

// HolderNickname gets holder's alternative name when defined within EAA Payload claims. Port of
// getHolderNickname().
func (p *EAAPayloadProxy) HolderNickname() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Nickname), p.getCredentialSubject().Nickname())
	}
	return nil
}

// HolderShortName gets holder's preferred or short name when defined within EAA Payload claims.
// Port of getHolderShortName().
func (p *EAAPayloadProxy) HolderShortName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.ShortName), p.getCredentialSubject().ShortName())
	}
	return nil
}

// HolderProfileUrl gets holder's profile URL when defined within EAA Payload claims. Port of
// getHolderProfileUrl().
func (p *EAAPayloadProxy) HolderProfileUrl() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.ProfileUrl), p.getCredentialSubject().ProfileUrl())
	}
	return nil
}

// HolderPictureUrl gets holder's picture URL when defined within EAA Payload claims. Port of
// getHolderPictureUrl().
func (p *EAAPayloadProxy) HolderPictureUrl() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.PictureUrl), p.getCredentialSubject().PictureUrl())
	}
	return nil
}

// HolderWebsiteUrl gets holder's website when defined within EAA Payload claims. Port of
// getHolderWebsiteUrl().
func (p *EAAPayloadProxy) HolderWebsiteUrl() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.WebsiteUrl), p.getCredentialSubject().WebsiteUrl())
	}
	return nil
}

// HolderEmail gets holder's email when defined within EAA Payload claims. Port of
// getHolderEmail().
func (p *EAAPayloadProxy) HolderEmail() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Email), p.getCredentialSubject().Email())
	}
	return nil
}

// HolderEmailVerified gets whether the holder's website has been verified if defined within EAA
// Payload claims. Port of getHolderEmailVerified().
func (p *EAAPayloadProxy) HolderEmailVerified() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.EmailVerified), p.getCredentialSubject().EmailVerified())
	}
	return nil
}

// HolderGender gets holder's gender when defined within EAA Payload claims. Port of
// getHolderGender().
func (p *EAAPayloadProxy) HolderGender() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Gender), p.getCredentialSubject().Gender())
	}
	return nil
}

// HolderBirthdate gets holder's birthdate when defined within EAA Payload claims. Port of
// getHolderBirthdate().
func (p *EAAPayloadProxy) HolderBirthdate() *BirthdateClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getBirthdateClaim(p.xmlEAAPayload.Birthdate), p.getCredentialSubject().Birthdate())
	}
	return nil
}

// HolderTimezone gets holder's timezone when defined within EAA Payload claims. Port of
// getHolderTimezone().
func (p *EAAPayloadProxy) HolderTimezone() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Timezone), p.getCredentialSubject().Timezone())
	}
	return nil
}

// HolderLocale gets holder's locale when defined within EAA Payload claims. Port of
// getHolderLocale().
func (p *EAAPayloadProxy) HolderLocale() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Locale), p.getCredentialSubject().Locale())
	}
	return nil
}

// HolderAddress gets holder's full address, when defined within EAA Payload claims. Port of
// getHolderAddress().
func (p *EAAPayloadProxy) HolderAddress() *AddressClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getAddressClaim(p.xmlEAAPayload.Address), p.getCredentialSubject().Address())
	}
	return nil
}

// HolderPhoneNumber gets holder's phone number when defined within EAA Payload claims. Port of
// getHolderPhoneNumber().
func (p *EAAPayloadProxy) HolderPhoneNumber() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.PhoneNumber), p.getCredentialSubject().PhoneNumber())
	}
	return nil
}

// HolderPhoneNumberVerified gets whether the holder's phone number has been verified if defined
// within EAA Payload claims. Port of getHolderPhoneNumberVerified().
func (p *EAAPayloadProxy) HolderPhoneNumberVerified() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.PhoneNumberVerified), p.getCredentialSubject().PhoneNumberVerified())
	}
	return nil
}

// HolderPlaceOfBirth gets holder's place of birth when defined within EAA Payload claims. Port
// of getHolderPlaceOfBirth().
func (p *EAAPayloadProxy) HolderPlaceOfBirth() *PlaceOfBirthClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getPlaceOfBirthClaim(p.xmlEAAPayload.PlaceOfBirth), p.getCredentialSubject().PlaceOfBirth())
	}
	return nil
}

// HolderNationalities gets holder's nationalities list when defined within EAA Payload claims.
// NOTE: The values are usually represented by 3-letter nationality codes. Port of
// getHolderNationalities().
func (p *EAAPayloadProxy) HolderNationalities() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Nationalities), p.getCredentialSubject().Nationalities())
	}
	return nil
}

// HolderBirthFamilyName gets holder's last or family name at birth when defined within EAA
// Payload claims. Port of getHolderBirthFamilyName().
func (p *EAAPayloadProxy) HolderBirthFamilyName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.BirthFamilyName), p.getCredentialSubject().BirthFamilyName())
	}
	return nil
}

// HolderBirthGivenName gets holder's first name at birth when defined within EAA Payload
// claims. Port of getHolderBirthGivenName().
func (p *EAAPayloadProxy) HolderBirthGivenName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.BirthGivenName), p.getCredentialSubject().BirthGivenName())
	}
	return nil
}

// HolderBirthMiddleName gets holder's middle name at birth when defined within EAA Payload
// claims. Port of getHolderBirthMiddleName().
func (p *EAAPayloadProxy) HolderBirthMiddleName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.BirthMiddleName), p.getCredentialSubject().BirthMiddleName())
	}
	return nil
}

// HolderBirthFullName gets the name(s) which holder was born. Port of getHolderBirthFullName().
//
// NOTE: unlike its siblings, Java's getHolderBirthFullName() does not fall back to the
// CredentialSubject claim; ported faithfully.
func (p *EAAPayloadProxy) HolderBirthFullName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.BirthFullName)
	}
	return nil
}

// HolderSalutation gets holder's preferred salutation when defined within EAA Payload claims.
// Port of getHolderSalutation().
func (p *EAAPayloadProxy) HolderSalutation() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Salutation), p.getCredentialSubject().Salutation())
	}
	return nil
}

// HolderTitle gets holder's title when defined within EAA Payload claims. Port of
// getHolderTitle().
func (p *EAAPayloadProxy) HolderTitle() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Title), p.getCredentialSubject().Title())
	}
	return nil
}

// HolderMobilePhoneNumber gets holder's mobile phone number when defined within EAA Payload
// claims. Port of getHolderMobilePhoneNumber().
func (p *EAAPayloadProxy) HolderMobilePhoneNumber() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.MobilePhoneNumber), p.getCredentialSubject().MobilePhoneNumber())
	}
	return nil
}

// HolderPseudonym gets holder's scenic name or pseudonym, they are known as, when defined
// within EAA Payload claims. Port of getHolderPseudonym().
func (p *EAAPayloadProxy) HolderPseudonym() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return firstNonNil(p.getClaim(p.xmlEAAPayload.Pseudonym), p.getCredentialSubject().Pseudonym())
	}
	return nil
}

/* mdoc claims */

// DocumentIssuingAuthority gets issuing authority name. The value shall only use latin1
// characters and shall have a maximum length of 150 characters. Port of
// getDocumentIssuingAuthority().
func (p *EAAPayloadProxy) DocumentIssuingAuthority() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.IssuingAuthority)
	}
	return nil
}

// DocumentIssuingAuthorityCountry gets alpha-2 country code, as defined in ISO 3166-1, of the
// issuing authority's country or territory. Port of getDocumentIssuingAuthorityCountry().
func (p *EAAPayloadProxy) DocumentIssuingAuthorityCountry() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.IssuingCountry)
	}
	return nil
}

// DocumentIssuingAuthorityJurisdiction gets a country subdivision code of the jurisdiction that
// issued the mDL as defined in ISO 3166-2:2020, Clause 8. The first part of the code shall be
// the same as the value for issuing_country. Port of getDocumentIssuingAuthorityJurisdiction().
func (p *EAAPayloadProxy) DocumentIssuingAuthorityJurisdiction() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.IssuingJurisdiction)
	}
	return nil
}

// DocumentIssuingAuthorityUNDistinguishingSign gets the distinguishing sign of the issuing
// country according to ISO/IEC 18013-1:2018, Annex F. If no applicable distinguishing sign is
// available in ISO/IEC 18013-1, an IA may use an empty identifier or another identifier by which
// it is internationally recognized. In this case the IA should ensure there is no collision
// with other IA's. Port of getDocumentIssuingAuthorityUNDistinguishingSign().
func (p *EAAPayloadProxy) DocumentIssuingAuthorityUNDistinguishingSign() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.UNDistinguishingSign)
	}
	return nil
}

// PersonalAdministrativeNumber gets an audit control number assigned by the issuing authority.
// The value shall only use latin1 characters and shall have a maximum length of 150 characters.
// Port of getPersonalAdministrativeNumber().
func (p *EAAPayloadProxy) PersonalAdministrativeNumber() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.PersonalAdministrativeNumber)
	}
	return nil
}

// DocumentNumber gets the number assigned or calculated by the issuing authority. The value
// shall only use latin1 characters and shall have a maximum length of 150 characters. Port of
// getDocumentNumber().
func (p *EAAPayloadProxy) DocumentNumber() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.DocumentNumber)
	}
	return nil
}

// HolderPortrait gets a reproduction of the mDL holder's portrait. Port of getHolderPortrait().
func (p *EAAPayloadProxy) HolderPortrait() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Portrait)
	}
	return nil
}

// HolderDrivingPrivileges gets the categories of vehicles/restrictions/conditions contain
// information describing the driving privileges of the mDL holder. Port of
// getHolderDrivingPrivileges().
func (p *EAAPayloadProxy) HolderDrivingPrivileges() *DrivingPrivilegesClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getDrivingPrivilegesClaim(p.xmlEAAPayload.DrivingPrivileges)
	}
	return nil
}

// HolderHeight gets the holder's height in centimetres. Port of getHolderHeight().
func (p *EAAPayloadProxy) HolderHeight() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Height)
	}
	return nil
}

// HolderWeight gets the holder's weight in kilograms. Port of getHolderWeight().
func (p *EAAPayloadProxy) HolderWeight() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Weight)
	}
	return nil
}

// HolderEyeColour gets the mDL holder's eye colour. The value shall be one of the following:
// "black", "blue", "brown", "dichromatic", "grey", "green", "hazel", "maroon", "pink",
// "unknown". Port of getHolderEyeColour().
func (p *EAAPayloadProxy) HolderEyeColour() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.EyeColour)
	}
	return nil
}

// HolderHairColour gets the mDL holder's hair colour. The value shall be one of the following:
// "bald", "black", "blond", "brown", "grey", "red", "auburn", "sandy", "white", "unknown". Port
// of getHolderHairColour().
func (p *EAAPayloadProxy) HolderHairColour() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.HairColour)
	}
	return nil
}

// ResidentPostalAddress gets the place where the mDL holder resides and/or may be contacted
// (street/house number, municipality etc.). The value shall only use latin1 characters and
// shall have a maximum length of 150 characters. Port of getResidentPostalAddress().
func (p *EAAPayloadProxy) ResidentPostalAddress() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentPostalAddress)
	}
	return nil
}

// HolderPortraitCaptureDate gets the date when portrait was taken. Port of
// getHolderPortraitCaptureDate().
func (p *EAAPayloadProxy) HolderPortraitCaptureDate() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.PortraitCaptureDate)
	}
	return nil
}

// HolderAgeInYears gets the date the age of the mDL holder. Port of getHolderAgeInYears().
func (p *EAAPayloadProxy) HolderAgeInYears() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.AgeInYears)
	}
	return nil
}

// HolderAgeBirthYear gets the year when the mDL holder was born. Port of
// getHolderAgeBirthYear().
func (p *EAAPayloadProxy) HolderAgeBirthYear() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.AgeBirthYear)
	}
	return nil
}

// HolderAgeEqualOrOver gets the map of claims attesting whether the User to whom the person
// identification data relates is at least NN years old. N <> 18. Multiple instances of this
// attribute may be present, provided the value of NN is different in each of them. If present,
// the requirements in clause 7.2.5 of ISO/IEC 18013-5 are applicable for these attributes. Port
// of getHolderAgeEqualOrOver().
func (p *EAAPayloadProxy) HolderAgeEqualOrOver() *AgeEqualOrOverClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getAgeEqualOrOverClaim(p.xmlEAAPayload.AgeEqualOrOver)
	}
	return nil
}

// HolderAgeOverList gets a list of elements used to convey to an mDL verifier, in a
// data-minimized fashion, if the mDL holder is as old or older than a specified age, or if the
// mDL holder is younger than a specified age. To achieve this, the mDL contains age attestation
// identifiers. An age attestation identifier has the format age_over_NN where NN is a value
// from 00 to 99. The value of an age attestation identifier can be TRUE or FALSE. Port of
// getHolderAgeOverList().
func (p *EAAPayloadProxy) HolderAgeOverList() []*AgeOverNNClaimWrapper {
	if p.xmlEAAPayload != nil {
		ageOverNN := p.xmlEAAPayload.AgeOverNN
		if len(ageOverNN) != 0 {
			result := make([]*AgeOverNNClaimWrapper, 0, len(ageOverNN))
			for _, item := range ageOverNN {
				result = append(result, NewAgeOverNNClaimWrapper(item))
			}
			return result
		}
	}
	return nil
}

// HolderResidentAddressCity gets the city where the mDL holder lives. The value shall only use
// latin1 characters and shall have a maximum length of 150 characters. Port of
// getHolderResidentAddressCity().
func (p *EAAPayloadProxy) HolderResidentAddressCity() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressCity)
	}
	return nil
}

// HolderResidentAddressState gets the state/province/district where the mDL holder lives. The
// value shall only use latin1 characters and shall have a maximum length of 150 characters.
// Port of getHolderResidentAddressState().
func (p *EAAPayloadProxy) HolderResidentAddressState() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressState)
	}
	return nil
}

// HolderResidentAddressPostalCode gets the postal code of the mDL holder. The value shall only
// use latin1 characters and shall have a maximum length of 150 characters. Port of
// getHolderResidentAddressPostalCode().
func (p *EAAPayloadProxy) HolderResidentAddressPostalCode() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressPostalCode)
	}
	return nil
}

// HolderResidentAddressCountry gets the country where the mDL holder lives as a two letter
// country code (alpha-2 code) defined in ISO 3166-1. Port of getHolderResidentAddressCountry().
func (p *EAAPayloadProxy) HolderResidentAddressCountry() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressCountry)
	}
	return nil
}

// HolderBiometricTemplateList gets a list of elements containing optional facial, fingerprint,
// iris, or other biometric information of the mDL holder. A biometric template identifier has
// the format biometric_template_xx where xx shall be replaced with the corresponding "Abstract
// value name" found in ISO/IEC 19785-3:2020, Table 7, according to the following convention:
// capitalized characters are replaced with their lowercase equivalent and spaces or
// non-alphanumeric characters are replaced by underscores (_). Port of
// getHolderBiometricTemplateList().
func (p *EAAPayloadProxy) HolderBiometricTemplateList() []*BiometricTemplateXXClaimWrapper {
	if p.xmlEAAPayload != nil {
		biometricTemplateList := p.xmlEAAPayload.BiometricTemplate
		if len(biometricTemplateList) != 0 {
			result := make([]*BiometricTemplateXXClaimWrapper, 0, len(biometricTemplateList))
			for _, item := range biometricTemplateList {
				result = append(result, NewBiometricTemplateXXClaimWrapper(item))
			}
			return result
		}
	}
	return nil
}

// HolderSignatureUsualMark gets an image of the signature or usual mark of the mDL holder, see
// 7.2.7 ISO/IEC 18013-5. Port of getHolderSignatureUsualMark().
func (p *EAAPayloadProxy) HolderSignatureUsualMark() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.SignatureUsualMark)
	}
	return nil
}

// HolderFingerprint gets a reproduction of the holder's fingerprint data (TBC). Port of
// getHolderFingerprint().
func (p *EAAPayloadProxy) HolderFingerprint() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Fingerprint)
	}
	return nil
}

// HolderBusinessName gets a business name of the holder. Port of getHolderBusinessName().
func (p *EAAPayloadProxy) HolderBusinessName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.BusinessName)
	}
	return nil
}

// HolderOrganizationName gets a name of legal person. Port of getHolderOrganizationName().
func (p *EAAPayloadProxy) HolderOrganizationName() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.OrganizationName)
	}
	return nil
}

// HolderProfession gets the profession of the holder. Port of getHolderProfession().
func (p *EAAPayloadProxy) HolderProfession() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Profession)
	}
	return nil
}

// HolderRelationshipFather gets the father of the holder. Port of getHolderRelationshipFather().
func (p *EAAPayloadProxy) HolderRelationshipFather() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipFather)
	}
	return nil
}

// HolderRelationshipMother gets the mother of the holder. Port of getHolderRelationshipMother().
func (p *EAAPayloadProxy) HolderRelationshipMother() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipMother)
	}
	return nil
}

// HolderRelationshipParent gets the parent of the holder. Port of getHolderRelationshipParent().
func (p *EAAPayloadProxy) HolderRelationshipParent() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipParent)
	}
	return nil
}

// HolderRelationshipSon gets the son of the holder. Port of getHolderRelationshipSon().
func (p *EAAPayloadProxy) HolderRelationshipSon() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipSon)
	}
	return nil
}

// HolderRelationshipDaughter gets the daughter of the holder. Port of
// getHolderRelationshipDaughter().
func (p *EAAPayloadProxy) HolderRelationshipDaughter() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipDaughter)
	}
	return nil
}

// HolderRelationshipBrother gets the brother of the holder. Port of
// getHolderRelationshipBrother().
func (p *EAAPayloadProxy) HolderRelationshipBrother() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipBrother)
	}
	return nil
}

// HolderRelationshipSister gets the sister of the holder. Port of
// getHolderRelationshipSister().
func (p *EAAPayloadProxy) HolderRelationshipSister() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipSister)
	}
	return nil
}

// HolderRelationshipSibling gets the sibling of the holder. Port of
// getHolderRelationshipSibling().
func (p *EAAPayloadProxy) HolderRelationshipSibling() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipSibling)
	}
	return nil
}

// HolderRelationshipSpouse gets the spouse of the holder. Port of
// getHolderRelationshipSpouse().
func (p *EAAPayloadProxy) HolderRelationshipSpouse() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipSpouse)
	}
	return nil
}

// HolderRelationshipFatherInLaw gets the father-in-law of the holder. Port of
// getHolderRelationshipFatherInLaw().
func (p *EAAPayloadProxy) HolderRelationshipFatherInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipFatherInLaw)
	}
	return nil
}

// HolderRelationshipMotherInLaw gets the mother-in-law of the holder. Port of
// getHolderRelationshipMotherInLaw().
func (p *EAAPayloadProxy) HolderRelationshipMotherInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipMotherInLaw)
	}
	return nil
}

// HolderRelationshipParentInLaw gets the parent-in-law of the holder. Port of
// getHolderRelationshipParentInLaw().
func (p *EAAPayloadProxy) HolderRelationshipParentInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipParentInLaw)
	}
	return nil
}

// HolderRelationshipSonInLaw gets the son-in-law of the holder. Port of
// getHolderRelationshipSonInLaw().
func (p *EAAPayloadProxy) HolderRelationshipSonInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipSonInLaw)
	}
	return nil
}

// HolderRelationshipDaughterInLaw gets the daughter-in-law of the holder. Port of
// getHolderRelationshipDaughterInLaw().
func (p *EAAPayloadProxy) HolderRelationshipDaughterInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipDaughterInLaw)
	}
	return nil
}

// HolderRelationshipChildInLaw gets the child-in-law of the holder. Port of
// getHolderRelationshipChildInLaw().
func (p *EAAPayloadProxy) HolderRelationshipChildInLaw() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipChildInLaw)
	}
	return nil
}

// HolderRelationshipParentalAuthority gets the parental authority of the holder. Port of
// getHolderRelationshipParentalAuthority().
func (p *EAAPayloadProxy) HolderRelationshipParentalAuthority() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipParentalAuthority)
	}
	return nil
}

// HolderRelationshipLegalRepresentative gets the legal representative of the holder. Port of
// getHolderRelationshipLegalRepresentative().
func (p *EAAPayloadProxy) HolderRelationshipLegalRepresentative() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipLegalRepresentative)
	}
	return nil
}

// HolderRelationshipAgent gets the voluntary agent of the holder. Port of
// getHolderRelationshipAgent().
func (p *EAAPayloadProxy) HolderRelationshipAgent() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.RelationshipAgent)
	}
	return nil
}

// DocumentType gets the document type. Port of getDocumentType().
func (p *EAAPayloadProxy) DocumentType() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.DocumentType)
	}
	return nil
}

// AdministrativeIssuanceDate gets the date when the data (e.g. a PID) was issued. Port of
// getAdministrativeIssuanceDate().
func (p *EAAPayloadProxy) AdministrativeIssuanceDate() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.AdministrativeIssuanceDate)
	}
	return nil
}

// AdministrativeExpirationDate gets the date when the data (e.g. a PID) will expire. Port of
// getAdministrativeExpirationDate().
func (p *EAAPayloadProxy) AdministrativeExpirationDate() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.AdministrativeExpirationDate)
	}
	return nil
}

// TrustAnchor gets the URL at which a machine-readable version of the trust anchor to be used
// for verifying the PID can be found or looked up. Port of getTrustAnchor().
func (p *EAAPayloadProxy) TrustAnchor() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.TrustAnchor)
	}
	return nil
}

// ResidentAddressStreet gets the name of the street where the user to whom the person
// identification data relates currently resides. Port of getResidentAddressStreet().
func (p *EAAPayloadProxy) ResidentAddressStreet() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressStreet)
	}
	return nil
}

// ResidentAddressHouseNumber gets the house number where the user to whom the person
// identification data relates currently resides, including any affix or suffix. Port of
// getResidentAddressHouseNumber().
func (p *EAAPayloadProxy) ResidentAddressHouseNumber() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ResidentAddressHouseNumber)
	}
	return nil
}

/* ETSI TS 119 472-1 "5 Implementation of EAA based on SD-JWT VC" header parameters */

// IssuingAuthorityRegistrationIdentifier gets the registration identifier of the legal entity
// on whose behalf the EAA has been issued. Port of getIssuingAuthorityRegistrationIdentifier().
func (p *EAAPayloadProxy) IssuingAuthorityRegistrationIdentifier() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.IssuingAuthorityRegistrationIdentifier)
	}
	return nil
}

// OneTimeUse gets the signal indicating that the EAA shall be used only once, and that it shall
// not be retained for future use. Port of getOneTimeUse().
func (p *EAAPayloadProxy) OneTimeUse() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.OneTimeUse)
	}
	return nil
}

// ShortLived gets the EAA short-lived component indicating that the validity period of the EAA
// is so short that it shall not be necessary to check its revocation status. Port of
// getShortLived().
func (p *EAAPayloadProxy) ShortLived() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.ShortLived)
	}
	return nil
}

// Evidence gets the array of evidence elements. Port of getEvidence().
func (p *EAAPayloadProxy) Evidence() *ClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getClaim(p.xmlEAAPayload.Evidence)
	}
	return nil
}

// AttestedAttributesSubject gets the claim for associating a set of attributes to one entity
// different than the EAA subject. Port of getAttestedAttributesSubject().
func (p *EAAPayloadProxy) AttestedAttributesSubject() *AttestedAttributesSubjectClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getAttestedAttributesSubjectClaim(p.xmlEAAPayload.AttestedAttributesSubject)
	}
	return nil
}

// CredentialSubjectClaims gets a list of credential subject claims. Port of
// getCredentialSubjectClaims().
func (p *EAAPayloadProxy) CredentialSubjectClaims() []*CredentialSubjectClaimWrapper {
	if p.xmlEAAPayload != nil {
		return p.getCredentialSubject().CredentialSubjects()
	}
	return nil
}

// OtherClaims gets a list of claims incorporated within the EAA Payload or provided as
// disclosures, which are not (yet) directly supported by the implementation. Port of
// getOtherClaims().
func (p *EAAPayloadProxy) OtherClaims() []*ClaimWrapper {
	if p.xmlEAAPayload != nil && p.xmlEAAPayload.OtherClaim != nil {
		result := make([]*ClaimWrapper, 0, len(p.xmlEAAPayload.OtherClaim))
		for _, item := range p.xmlEAAPayload.OtherClaim {
			result = append(result, NewClaimWrapper(item))
		}
		return result
	}
	return nil
}

// AllEAAPayloadClaims gets a list of all claims present within an EAA Payload. Port of
// getAllEAAPayloadClaims().
func (p *EAAPayloadProxy) AllEAAPayloadClaims() []*ClaimWrapper {
	if p.xmlEAAPayload == nil {
		return nil
	}

	var claimList []*ClaimWrapper

	if p.xmlEAAPayload.Identifier != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Identifier))
	}
	if p.xmlEAAPayload.Issuer != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Issuer))
	}
	if p.xmlEAAPayload.Subject != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Subject))
	}
	if p.xmlEAAPayload.Audience != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Audience))
	}
	if p.xmlEAAPayload.IssuedAt != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.IssuedAt))
	}
	if p.xmlEAAPayload.NotBefore != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.NotBefore))
	}
	if p.xmlEAAPayload.Expiration != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Expiration))
	}
	if p.xmlEAAPayload.UpdatedAt != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.UpdatedAt))
	}
	if p.xmlEAAPayload.Category != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Category))
	}
	if p.xmlEAAPayload.VerifiableCredentialsType != nil {
		claimList = append(claimList, p.getClaim(claimBase(p.xmlEAAPayload.VerifiableCredentialsType.XmlClaimContent, p.xmlEAAPayload.VerifiableCredentialsType.XmlClaimAttrs)))
		if p.xmlEAAPayload.VerifiableCredentialsType.Integrity != nil {
			claimList = append(claimList, p.getIntegrityClaim(p.xmlEAAPayload.VerifiableCredentialsType.Integrity).AsClaim())
		}
	}
	if p.xmlEAAPayload.Status != nil {
		claimList = append(claimList, p.getStatusClaim(p.xmlEAAPayload.Status).AsClaim())
	}
	if p.xmlEAAPayload.Nonce != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Nonce))
	}
	if p.xmlEAAPayload.DeviceKey != nil {
		claimList = append(claimList, p.getDeviceKeyClaim(p.xmlEAAPayload.DeviceKey).AsClaim())
	}
	if p.xmlEAAPayload.Version != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Version))
	}
	if p.xmlEAAPayload.DocType != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.DocType))
	}
	if p.xmlEAAPayload.ValidityInfo != nil {
		claimList = append(claimList, p.getValidityInfoClaim(p.xmlEAAPayload.ValidityInfo).AsClaim())
	}
	if p.xmlEAAPayload.FullName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.FullName))
	}
	if p.xmlEAAPayload.GivenName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.GivenName))
	}
	if p.xmlEAAPayload.FamilyName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.FamilyName))
	}
	if p.xmlEAAPayload.MiddleName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.MiddleName))
	}
	if p.xmlEAAPayload.Nickname != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Nickname))
	}
	if p.xmlEAAPayload.ShortName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ShortName))
	}
	if p.xmlEAAPayload.ProfileUrl != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ProfileUrl))
	}
	if p.xmlEAAPayload.PictureUrl != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.PictureUrl))
	}
	if p.xmlEAAPayload.WebsiteUrl != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.WebsiteUrl))
	}
	if p.xmlEAAPayload.Email != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Email))
	}
	if p.xmlEAAPayload.EmailVerified != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.EmailVerified))
	}
	if p.xmlEAAPayload.Gender != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Gender))
	}
	if p.xmlEAAPayload.Birthdate != nil {
		claimList = append(claimList, p.getBirthdateClaim(p.xmlEAAPayload.Birthdate).AsClaim())
	}
	if p.xmlEAAPayload.Timezone != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Timezone))
	}
	if p.xmlEAAPayload.Locale != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Locale))
	}
	if p.xmlEAAPayload.Address != nil {
		claimList = append(claimList, p.getAddressClaim(p.xmlEAAPayload.Address).AsClaim())
	}
	if p.xmlEAAPayload.PhoneNumber != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.PhoneNumber))
	}
	if p.xmlEAAPayload.PhoneNumberVerified != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.PhoneNumberVerified))
	}
	if p.xmlEAAPayload.PlaceOfBirth != nil {
		claimList = append(claimList, p.getPlaceOfBirthClaim(p.xmlEAAPayload.PlaceOfBirth).AsClaim())
	}
	if p.xmlEAAPayload.Nationalities != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Nationalities))
	}
	if p.xmlEAAPayload.BirthFamilyName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.BirthFamilyName))
	}
	if p.xmlEAAPayload.BirthGivenName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.BirthGivenName))
	}
	if p.xmlEAAPayload.BirthMiddleName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.BirthMiddleName))
	}
	if p.xmlEAAPayload.Salutation != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Salutation))
	}
	if p.xmlEAAPayload.Title != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Title))
	}
	if p.xmlEAAPayload.MobilePhoneNumber != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.MobilePhoneNumber))
	}
	if p.xmlEAAPayload.Pseudonym != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Pseudonym))
	}
	if p.xmlEAAPayload.CredentialSubject != nil {
		for _, credentialSubjectClaim := range p.getCredentialSubject().CredentialSubjects() {
			claimList = append(claimList, credentialSubjectClaim.AsClaim())
		}
	}
	if p.xmlEAAPayload.IssuingCountry != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.IssuingCountry))
	}
	if p.xmlEAAPayload.IssuingAuthority != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.IssuingAuthority))
	}
	if p.xmlEAAPayload.DocumentNumber != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.DocumentNumber))
	}
	if p.xmlEAAPayload.Portrait != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Portrait))
	}
	if p.xmlEAAPayload.DrivingPrivileges != nil {
		claimList = append(claimList, p.getDrivingPrivilegesClaim(p.xmlEAAPayload.DrivingPrivileges).AsClaim())
	}
	if p.xmlEAAPayload.UNDistinguishingSign != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.UNDistinguishingSign))
	}
	if p.xmlEAAPayload.PersonalAdministrativeNumber != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.PersonalAdministrativeNumber))
	}
	if p.xmlEAAPayload.Height != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Height))
	}
	if p.xmlEAAPayload.Weight != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Weight))
	}
	if p.xmlEAAPayload.EyeColour != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.EyeColour))
	}
	if p.xmlEAAPayload.HairColour != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.HairColour))
	}
	if p.xmlEAAPayload.ResidentPostalAddress != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentPostalAddress))
	}
	if p.xmlEAAPayload.PortraitCaptureDate != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.PortraitCaptureDate))
	}
	if p.xmlEAAPayload.AgeInYears != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.AgeInYears))
	}
	if p.xmlEAAPayload.AgeBirthYear != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.AgeBirthYear))
	}
	if p.xmlEAAPayload.AgeEqualOrOver != nil {
		claimList = append(claimList, p.getClaim(claimBase(p.xmlEAAPayload.AgeEqualOrOver.XmlClaimContent, p.xmlEAAPayload.AgeEqualOrOver.XmlClaimAttrs)))
	}
	if p.xmlEAAPayload.AgeOverNN != nil {
		for _, item := range p.xmlEAAPayload.AgeOverNN {
			claimList = append(claimList, p.getClaim(claimBase(item.XmlClaimContent, item.XmlClaimAttrs)))
		}
	}
	if p.xmlEAAPayload.IssuingJurisdiction != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.IssuingJurisdiction))
	}
	if p.xmlEAAPayload.ResidentAddressCity != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressCity))
	}
	if p.xmlEAAPayload.ResidentAddressState != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressState))
	}
	if p.xmlEAAPayload.ResidentAddressPostalCode != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressPostalCode))
	}
	if p.xmlEAAPayload.ResidentAddressCountry != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressCountry))
	}
	if p.xmlEAAPayload.BiometricTemplate != nil {
		for _, item := range p.xmlEAAPayload.BiometricTemplate {
			claimList = append(claimList, p.getClaim(claimBase(item.XmlClaimContent, item.XmlClaimAttrs)))
		}
	}
	if p.xmlEAAPayload.SignatureUsualMark != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.SignatureUsualMark))
	}
	if p.xmlEAAPayload.Fingerprint != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Fingerprint))
	}
	if p.xmlEAAPayload.BusinessName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.BusinessName))
	}
	if p.xmlEAAPayload.OrganizationName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.OrganizationName))
	}
	if p.xmlEAAPayload.BirthFullName != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.BirthFullName))
	}
	if p.xmlEAAPayload.Profession != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Profession))
	}
	if p.xmlEAAPayload.RelationshipFather != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipFather))
	}
	if p.xmlEAAPayload.RelationshipMother != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipMother))
	}
	if p.xmlEAAPayload.RelationshipParent != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipParent))
	}
	if p.xmlEAAPayload.RelationshipSon != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipSon))
	}
	if p.xmlEAAPayload.RelationshipDaughter != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipDaughter))
	}
	if p.xmlEAAPayload.RelationshipBrother != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipBrother))
	}
	if p.xmlEAAPayload.RelationshipSister != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipSister))
	}
	if p.xmlEAAPayload.RelationshipSibling != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipSibling))
	}
	if p.xmlEAAPayload.RelationshipSpouse != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipSpouse))
	}
	if p.xmlEAAPayload.RelationshipFatherInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipFatherInLaw))
	}
	if p.xmlEAAPayload.RelationshipMotherInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipMotherInLaw))
	}
	if p.xmlEAAPayload.RelationshipParentInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipParentInLaw))
	}
	if p.xmlEAAPayload.RelationshipSonInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipSonInLaw))
	}
	if p.xmlEAAPayload.RelationshipDaughterInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipDaughterInLaw))
	}
	if p.xmlEAAPayload.RelationshipChildInLaw != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipChildInLaw))
	}
	if p.xmlEAAPayload.RelationshipParentalAuthority != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipParentalAuthority))
	}
	if p.xmlEAAPayload.RelationshipLegalRepresentative != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipLegalRepresentative))
	}
	if p.xmlEAAPayload.RelationshipAgent != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.RelationshipAgent))
	}
	if p.xmlEAAPayload.DocumentType != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.DocumentType))
	}
	if p.xmlEAAPayload.AdministrativeIssuanceDate != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.AdministrativeIssuanceDate))
	}
	if p.xmlEAAPayload.AdministrativeExpirationDate != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.AdministrativeExpirationDate))
	}
	if p.xmlEAAPayload.TrustAnchor != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.TrustAnchor))
	}
	if p.xmlEAAPayload.ResidentAddressStreet != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressStreet))
	}
	if p.xmlEAAPayload.ResidentAddressHouseNumber != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ResidentAddressHouseNumber))
	}
	if p.xmlEAAPayload.IssuingAuthorityRegistrationIdentifier != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.IssuingAuthorityRegistrationIdentifier))
	}
	if p.xmlEAAPayload.OneTimeUse != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.OneTimeUse))
	}
	if p.xmlEAAPayload.ShortLived != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.ShortLived))
	}
	if p.xmlEAAPayload.Evidence != nil {
		claimList = append(claimList, p.getClaim(p.xmlEAAPayload.Evidence))
	}
	if p.xmlEAAPayload.AttestedAttributesSubject != nil {
		claimList = append(claimList, p.getClaim(claimBase(p.xmlEAAPayload.AttestedAttributesSubject.XmlClaimContent, p.xmlEAAPayload.AttestedAttributesSubject.XmlClaimAttrs)))
	}
	if len(p.xmlEAAPayload.OtherClaim) != 0 {
		for _, item := range p.xmlEAAPayload.OtherClaim {
			claimList = append(claimList, NewClaimWrapper(item))
		}
	}

	return claimList
}

// claimBase reassembles a base jaxb.XmlClaim from the content/attrs embedded in one of its XSD
// extension types (XmlVerifiableCredentialsTypeClaim, XmlAgeEqualOrOverClaim,
// XmlAgeOverNNClaim, XmlBiometricTemplateXXClaim, XmlAttestedAttributesSubjectClaim, ...). The
// generated Go jaxb structs embed XmlClaimContent/XmlClaimAttrs directly rather than embedding
// XmlClaim itself (see jaxb_claim.go), so there is no single promoted XmlClaim field to take the
// address of the way Java's upcast to Claim/XmlClaim happens implicitly; this reconstructs the
// equivalent value from the shared embedded parts.
func claimBase(content jaxb.XmlClaimContent, attrs jaxb.XmlClaimAttrs) *jaxb.XmlClaim {
	return &jaxb.XmlClaim{XmlClaimContent: content, XmlClaimAttrs: attrs}
}

// getClaim is the private helper backing most claim accessors: nil in, nil out. Port of
// getClaim(XmlClaim).
func (p *EAAPayloadProxy) getClaim(xmlDisclosableClaim *jaxb.XmlClaim) *ClaimWrapper {
	if xmlDisclosableClaim == nil {
		return nil
	}
	return NewClaimWrapper(xmlDisclosableClaim)
}

// firstNonNil returns the first non-nil claim among its arguments, or nil. Port of the private
// generic helper get(T... claims).
func firstNonNil[T any](claims ...*T) *T {
	for _, claim := range claims {
		if claim != nil {
			return claim
		}
	}
	return nil
}

// getDeviceKeyClaim is the private helper for the device key claim. Port of
// getDeviceKeyClaim(XmlDeviceKeyClaim).
func (p *EAAPayloadProxy) getDeviceKeyClaim(xmlDeviceKeyClaim *jaxb.XmlDeviceKeyClaim) *DeviceKeyClaimWrapper {
	if xmlDeviceKeyClaim == nil {
		return nil
	}
	return NewDeviceKeyClaimWrapper(xmlDeviceKeyClaim)
}

// getValidityInfoClaim is the private helper for the validity info claim. Port of
// getValidityInfoClaim(XmlValidityInfoClaim).
func (p *EAAPayloadProxy) getValidityInfoClaim(xmlValidityInfoClaim *jaxb.XmlValidityInfoClaim) *ValidityInfoClaimWrapper {
	if xmlValidityInfoClaim == nil {
		return nil
	}
	return NewValidityInfoClaimWrapper(xmlValidityInfoClaim)
}

// getIntegrityClaim is the private helper for the integrity claim. Port of
// getIntegrityClaim(XmlIntegrityClaim).
func (p *EAAPayloadProxy) getIntegrityClaim(xmlIntegrityClaim *jaxb.XmlIntegrityClaim) *IntegrityClaimWrapper {
	if xmlIntegrityClaim == nil {
		return nil
	}
	return NewIntegrityClaimWrapper(xmlIntegrityClaim)
}

// getAddressClaim is the private helper for the address claim. Port of
// getAddressClaim(XmlAddressClaim).
func (p *EAAPayloadProxy) getAddressClaim(xmlAddressClaim *jaxb.XmlAddressClaim) *AddressClaimWrapper {
	if xmlAddressClaim == nil {
		return nil
	}
	return NewAddressClaimWrapper(xmlAddressClaim)
}

// getBirthdateClaim is the private helper for the birthdate claim. Port of
// getBirthdateClaim(XmlClaim); ported against the field's actual generated type XmlBirthdateClaim
// (Java's XmlBirthdateClaim extends the Claim complexType, and every call site passes the
// XmlBirthdateClaim-typed field; the Go jaxb structs are not generated with that inheritance, so
// the concrete type is used directly here instead of the declared-supertype signature).
func (p *EAAPayloadProxy) getBirthdateClaim(xmlBirthdateClaim *jaxb.XmlBirthdateClaim) *BirthdateClaimWrapper {
	if xmlBirthdateClaim == nil {
		return nil
	}
	return NewBirthdateClaimWrapper(xmlBirthdateClaim)
}

// getPlaceOfBirthClaim is the private helper for the place-of-birth claim. Port of
// getPlaceOfBirthClaim(XmlClaim); see the getBirthdateClaim note on the concrete-type deviation.
func (p *EAAPayloadProxy) getPlaceOfBirthClaim(xmlPlaceOfBirthClaim *jaxb.XmlPlaceOfBirthClaim) *PlaceOfBirthClaimWrapper {
	if xmlPlaceOfBirthClaim == nil {
		return nil
	}
	return NewPlaceOfBirthClaimWrapper(xmlPlaceOfBirthClaim)
}

// getStatusClaim is the private helper for the status claim. Port of
// getStatusClaim(XmlStatusClaim).
func (p *EAAPayloadProxy) getStatusClaim(xmlStatusClaim *jaxb.XmlStatusClaim) *StatusClaimWrapper {
	if xmlStatusClaim == nil {
		return nil
	}
	return NewStatusClaimWrapper(xmlStatusClaim)
}

// getCredentialSubject is the private helper wrapping the credential subject list. Port of
// getCredentialSubject().
func (p *EAAPayloadProxy) getCredentialSubject() *CredentialSubjectProxy {
	return NewCredentialSubjectProxy(p.xmlEAAPayload.CredentialSubject)
}

// getDrivingPrivilegesClaim is the private helper for the driving privileges claim. Port of
// getDrivingPrivilegesClaim(XmlDrivingPrivilegesClaim).
func (p *EAAPayloadProxy) getDrivingPrivilegesClaim(xmlDrivingPrivilegesClaim *jaxb.XmlDrivingPrivilegesClaim) *DrivingPrivilegesClaimWrapper {
	if xmlDrivingPrivilegesClaim == nil {
		return nil
	}
	return NewDrivingPrivilegesClaimWrapper(xmlDrivingPrivilegesClaim)
}

// getAttestedAttributesSubjectClaim is the private helper for the attested-attributes-subject
// claim. Port of getAttestedAttributesSubjectClaim(XmlAttestedAttributesSubjectClaim).
func (p *EAAPayloadProxy) getAttestedAttributesSubjectClaim(xmlAttestedAttributesSubjectClaim *jaxb.XmlAttestedAttributesSubjectClaim) *AttestedAttributesSubjectClaimWrapper {
	if xmlAttestedAttributesSubjectClaim == nil {
		return nil
	}
	return NewAttestedAttributesSubjectClaimWrapper(xmlAttestedAttributesSubjectClaim)
}

// getAgeEqualOrOverClaim is the private helper for the age-equal-or-over claim. Port of
// getAgeEqualOrOverClaim(XmlAgeEqualOrOverClaim).
func (p *EAAPayloadProxy) getAgeEqualOrOverClaim(xmlAgeEqualOrOverClaim *jaxb.XmlAgeEqualOrOverClaim) *AgeEqualOrOverClaimWrapper {
	if xmlAgeEqualOrOverClaim == nil {
		return nil
	}
	return NewAgeEqualOrOverClaimWrapper(xmlAgeEqualOrOverClaim)
}
