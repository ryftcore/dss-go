// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/EAAWrapper.java (DSS 6.5.RC1).
//
// Depends on the `claim` Java subpackage (ClaimWrapper, AddressClaimWrapper, ...), which
// S8A_BRIEF.md flattens into this same Go package `diagnostic` ("wrappers + `claim` (mutual
// imports, collision-checked) -> ONE pkg diagnostic"); those types are assigned to a sibling
// chunk and are referenced here unqualified, matching that flattening. Also depends on
// EAAPayloadProxy (this file's sibling in the DIAGWRAP_A manifest, see eaa_payload_proxy.go) and
// on SignatureWrapper (DIAGWRAP_B).
package diagnostic

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// EAAWrapper provides a user-friendly interface for information extraction from a
// jaxb.XmlEAA JAXB object.
type EAAWrapper struct {
	AbstractTokenProxyBase

	// eaa is the wrapped EAA object.
	eaa *jaxb.XmlEAA
}

// NewEAAWrapper is the default constructor.
func NewEAAWrapper(eaa *jaxb.XmlEAA) *EAAWrapper {
	w := &EAAWrapper{eaa: eaa}
	w.InitTokenProxy(w)
	return w
}

// Id gets unique identifier. Port of getId().
func (w *EAAWrapper) Id() string {
	return w.eaa.Id
}

// Filename returns name of the EAA presentation's document, when applicable. Port of
// getFilename().
func (w *EAAWrapper) Filename() string {
	return w.eaa.DocumentName
}

// EAADocumentType gets claimed document type. NOTE: used in mdoc and the returned value
// corresponds to a string incorporated within a 'docType' element. Port of
// getEAADocumentType().
func (w *EAAWrapper) EAADocumentType() string {
	docType := payloadClaimTextValue(w.EAAPayload().EAADocType())
	if docType != "" {
		return docType
	}
	return w.eaa.DocumentType
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of
// getCurrentBasicSignature().
func (w *EAAWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	eaaSignature := w.getEAASignature()
	if eaaSignature != nil {
		return eaaSignature.CurrentBasicSignature()
	}
	return nil
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *EAAWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	eaaSignature := w.getEAASignature()
	if eaaSignature != nil {
		return eaaSignature.CurrentCertificateChain()
	}
	return nil
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *EAAWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	eaaSignature := w.getEAASignature()
	if eaaSignature != nil {
		return eaaSignature.CurrentSigningCertificate()
	}
	return nil
}

func (w *EAAWrapper) getEAASignature() *SignatureWrapper {
	eaaSignatures := w.EAASignatures()
	if len(eaaSignatures) == 1 {
		return eaaSignatures[0]
	}
	return nil
}

// SelectiveDisclosuresDigestAlgorithm gets digest algorithm used on hashes computation for
// selectively disclosable claims. Port of getSelectiveDisclosuresDigestAlgorithm().
func (w *EAAWrapper) SelectiveDisclosuresDigestAlgorithm() enumerations.DigestAlgorithm {
	return w.eaa.DigestMethod
}

// DigestMatchers gets a list of digest matchers representing the associated hashes and
// disclosures validation. Port of getDigestMatchers() (overridden).
func (w *EAAWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return w.eaa.DigestMatchers
}

// FoundCertificates is the AbstractTokenProxy default (not overridden in Java). Port of
// foundCertificates().
func (w *EAAWrapper) FoundCertificates() *FoundCertificatesProxy {
	return DefaultFoundCertificates()
}

// FoundRevocations is the AbstractTokenProxy default (not overridden in Java). Port of
// foundRevocations().
func (w *EAAWrapper) FoundRevocations() *FoundRevocationsProxy {
	return DefaultFoundRevocations()
}

// EAASignatures gets signatures used to create the EAA. NOTE: in most of the cases a single
// signature is expected, but it is possible for EAA to be signed by multiple signers. Port of
// getEAASignatures().
func (w *EAAWrapper) EAASignatures() []*SignatureWrapper {
	var result []*SignatureWrapper
	for _, xmlEAASignature := range w.eaa.EAASignature {
		result = append(result, NewSignatureWrapper(xmlEAASignature.Signature))
	}
	return result
}

// EAASignatureIds gets a list of identifiers of signatures used to create the EAA. Port of
// getEAASignatureIds().
func (w *EAAWrapper) EAASignatureIds() []string {
	eaaPresentationSignatures := w.EAASignatures()
	var result []string
	if len(eaaPresentationSignatures) != 0 {
		for _, s := range eaaPresentationSignatures {
			result = append(result, s.Id())
		}
		return result
	}
	return nil
}

// KeyBindingSignature gets a key binding signature, when present. Port of
// getKeyBindingSignature().
func (w *EAAWrapper) KeyBindingSignature() *SignatureWrapper {
	if w.eaa.KeyBindingSignature != nil {
		return NewSignatureWrapper(w.eaa.KeyBindingSignature.Signature)
	}
	return nil
}

// KeyBindingSignatureId gets unique identifier of the key binding signature, when present.
// Port of getKeyBindingSignatureId().
func (w *EAAWrapper) KeyBindingSignatureId() string {
	keyBindingSignature := w.KeyBindingSignature()
	if keyBindingSignature != nil {
		return keyBindingSignature.Id()
	}
	return ""
}

// EAAPayload gets access to the EAA payload, containing complete claims data. Port of
// getEAAPayload().
func (w *EAAWrapper) EAAPayload() *EAAPayloadProxy {
	return NewEAAPayloadProxy(w.eaa.EAAPayload)
}

// KeyBindingSignatureNonce gets the nonce provided in the key binding signature payload. Port
// of getKeyBindingSignatureNonce().
func (w *EAAWrapper) KeyBindingSignatureNonce() string {
	if w.eaa.KeyBindingPayload == nil {
		return ""
	}
	nonce := w.eaa.KeyBindingPayload.Nonce
	if nonce != nil {
		return nonce.Text
	}
	return ""
}

// KeyBindingSignatureAudience gets the audience provided in the key binding signature payload.
// Port of getKeyBindingSignatureAudience().
func (w *EAAWrapper) KeyBindingSignatureAudience() string {
	if w.eaa.KeyBindingPayload == nil {
		return ""
	}
	audience := w.eaa.KeyBindingPayload.Audience
	if audience != nil {
		return audience.Text
	}
	return ""
}

// KeyBindingSignatureIssuanceTime gets the issuance time provided in the key binding signature
// payload. Port of getKeyBindingSignatureIssuanceTime().
func (w *EAAWrapper) KeyBindingSignatureIssuanceTime() *time.Time {
	if w.eaa.KeyBindingPayload == nil {
		return nil
	}
	issuanceTime := w.eaa.KeyBindingPayload.IssuanceTime
	if issuanceTime != nil {
		return issuanceTime.DateTime
	}
	return nil
}

// OtherKeyBindingPayloadClaims gets a list of claims incorporated within the key binding
// payload, which are not (yet) directly supported by the implementation. Port of
// getOtherKeyBindingPayloadClaims().
func (w *EAAWrapper) OtherKeyBindingPayloadClaims() []*ClaimWrapper {
	if w.eaa.KeyBindingPayload != nil && w.eaa.KeyBindingPayload.OtherClaim != nil {
		var result []*ClaimWrapper
		for _, c := range w.eaa.KeyBindingPayload.OtherClaim {
			result = append(result, NewClaimWrapper(c))
		}
		return result
	}
	return nil
}

// EAAIdentifier gets EAA identifier provided in the EAA payload. Port of getEAAIdentifier().
func (w *EAAWrapper) EAAIdentifier() string {
	return payloadClaimTextValue(w.EAAPayload().EAAIdentifier())
}

// EAAIssuer gets EAA issuer as defined in the EAA payload. Port of getEAAIssuer().
func (w *EAAWrapper) EAAIssuer() string {
	return payloadClaimTextValue(w.EAAPayload().EAAIssuer())
}

// EAASubject gets EAA subject as defined in the EAA payload. Port of getEAASubject().
func (w *EAAWrapper) EAASubject() string {
	return payloadClaimTextValue(w.EAAPayload().EAASubject())
}

// EAAAudience gets EAA audience as defined in the EAA payload. Port of getEAAAudience().
func (w *EAAWrapper) EAAAudience() string {
	return payloadClaimTextValue(w.EAAPayload().EAAAudience())
}

// EAAIssuedAt gets EAA issuance time as defined in the EAA payload. Port of getEAAIssuedAt().
func (w *EAAWrapper) EAAIssuedAt() *time.Time {
	issuedAt := payloadClaimDateValue(w.EAAPayload().EAAIssuedAt())
	if issuedAt != nil {
		return issuedAt
	}
	eaaValidityInfo := w.EAAPayload().EAAValidityInfo()
	if eaaValidityInfo != nil {
		return payloadClaimDateValue(eaaValidityInfo.Signed())
	}
	return nil
}

// EAANotBefore gets EAA not before time as defined in the EAA payload. Port of
// getEAANotBefore().
func (w *EAAWrapper) EAANotBefore() *time.Time {
	notBefore := payloadClaimDateValue(w.EAAPayload().EAANotBefore())
	if notBefore != nil {
		return notBefore
	}
	eaaValidityInfo := w.EAAPayload().EAAValidityInfo()
	if eaaValidityInfo != nil {
		return payloadClaimDateValue(eaaValidityInfo.ValidFrom())
	}
	return nil
}

// EAAExpiration gets EAA expiration time as defined in the EAA payload. Port of
// getEAAExpiration().
func (w *EAAWrapper) EAAExpiration() *time.Time {
	expirationTime := payloadClaimDateValue(w.EAAPayload().EAAExpiration())
	if expirationTime != nil {
		return expirationTime
	}
	eaaValidityInfo := w.EAAPayload().EAAValidityInfo()
	if eaaValidityInfo != nil {
		return payloadClaimDateValue(eaaValidityInfo.ValidUntil())
	}
	return nil
}

// EAAUpdatedAt gets EAA update time as defined in the EAA payload. Port of getEAAUpdatedAt().
func (w *EAAWrapper) EAAUpdatedAt() *time.Time {
	return payloadClaimDateValue(w.EAAPayload().EAAUpdatedAt())
}

// EAANextUpdate gets EAA expected next update time. Port of getEAANextUpdate().
func (w *EAAWrapper) EAANextUpdate() *time.Time {
	eaaValidityInfo := w.EAAPayload().EAAValidityInfo()
	if eaaValidityInfo != nil {
		return payloadClaimDateValue(eaaValidityInfo.ExpectedUpdate())
	}
	return nil
}

// EAACategory gets category URN provided in the EAA payload. Port of getEAACategory().
func (w *EAAWrapper) EAACategory() string {
	return payloadClaimTextValue(w.EAAPayload().EAACategory())
}

// CategoryQualification gets the EAA Qualification based on the category URN provided in the
// EAA payload. Port of getCategoryQualification().
func (w *EAAWrapper) CategoryQualification() enumerations.EAAQualification {
	eaaCategory := w.EAACategory()
	switch {
	case string(enumerations.EAACategory_EU_QEAA.URN()) == eaaCategory:
		return enumerations.EAAQualification_QEAA
	case string(enumerations.EAACategory_EU_PUBEAA.URN()) == eaaCategory:
		return enumerations.EAAQualification_PUBEAA
	case eaaCategory == "":
		/*
		 * EAA-5.2.2.1-01: SD-JWT VC EAAs issued by EAAs issuers registered in the European Union,
		 * which are neither SD-JWT VC QEAAs nor SD-JWT VC PuB-EAAs, shall not include the category claim.
		 */
		return enumerations.EAAQualification_EAA
	default:
		return enumerations.EAAQualification_UNKNOWN
	}
}

// EAAVerifiableCredentialsTypeUri gets EAA metadata URI (e.g. 'vct' claim) as defined in the
// EAA payload. Port of getEAAVerifiableCredentialsTypeUri().
func (w *EAAWrapper) EAAVerifiableCredentialsTypeUri() string {
	return payloadClaimTextValue(w.EAAPayload().EAAVerifiableCredentialsType())
}

// EAAVerifiableCredentialsTypeIntegrityDigestAlgorithm gets Digest Algorithm used to compute
// the integrity material for the EAA metadata (when present). Port of
// getEAAVerifiableCredentialsTypeIntegrityDigestAlgorithm().
func (w *EAAWrapper) EAAVerifiableCredentialsTypeIntegrityDigestAlgorithm() enumerations.DigestAlgorithm {
	eaaVerifiableCredentialsTypeIntegrity := w.EAAPayload().EAAVerifiableCredentialsTypeIntegrity()
	if eaaVerifiableCredentialsTypeIntegrity != nil {
		return eaaVerifiableCredentialsTypeIntegrity.DigestAlgorithm()
	}
	return ""
}

// EAAVerifiableCredentialsTypeIntegrityBytes gets the integrity material for the EAA metadata
// (when present). Port of getEAAVerifiableCredentialsTypeIntegrityBytes().
func (w *EAAWrapper) EAAVerifiableCredentialsTypeIntegrityBytes() []byte {
	eaaVerifiableCredentialsTypeIntegrity := w.EAAPayload().EAAVerifiableCredentialsTypeIntegrity()
	if eaaVerifiableCredentialsTypeIntegrity != nil {
		return eaaVerifiableCredentialsTypeIntegrity.DigestValue()
	}
	return nil
}

// EAAStatusIndex gets EAA status index as defined in the EAA payload. Port of
// getEAAStatusIndex().
func (w *EAAWrapper) EAAStatusIndex() *int {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil {
		if eaaStatus.Index() != nil {
			return payloadClaimIntegerValue(eaaStatus.Index())
		} else if eaaStatus.StatusList() != nil {
			return payloadClaimIntegerValue(eaaStatus.StatusList().Index())
		}
	}
	return nil
}

// EAAStatusUri gets EAA status URI as defined in the EAA payload. Port of getEAAStatusUri().
func (w *EAAWrapper) EAAStatusUri() string {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil {
		if eaaStatus.Uri() != nil {
			return payloadClaimTextValue(eaaStatus.Uri())
		} else if eaaStatus.StatusList() != nil {
			return payloadClaimTextValue(eaaStatus.StatusList().Uri())
		}
	}
	return ""
}

// EAAStatusCertificate returns a certificate containing the public key that signed or sealed
// the top-level certificate in the x5chain element in the MSO revocation list structure. Port
// of getEAAStatusCertificate().
func (w *EAAWrapper) EAAStatusCertificate() []byte {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil && eaaStatus.StatusList() != nil {
		return payloadClaimByteValue(eaaStatus.StatusList().Certificate())
	}
	return nil
}

// EAAStatusType gets EAA status type as defined in the EAA payload. Port of getEAAStatusType().
func (w *EAAWrapper) EAAStatusType() string {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil {
		return payloadClaimTextValue(eaaStatus.Type())
	}
	return ""
}

// EAAStatusPurpose gets EAA status purpose as defined in the EAA payload. Port of
// getEAAStatusPurpose().
func (w *EAAWrapper) EAAStatusPurpose() string {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil {
		return payloadClaimTextValue(eaaStatus.Purpose())
	}
	return ""
}

// EAAIdentifierListId gets EAA identifier to be used for the EAA status verification using an
// Identifier List mechanism. Port of getEAAIdentifierListId().
func (w *EAAWrapper) EAAIdentifierListId() []byte {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil && eaaStatus.IdentifierList() != nil {
		return payloadClaimByteValue(eaaStatus.IdentifierList().Identifier())
	}
	return nil
}

// EAAIdentifierListUri gets EAA status URI as defined in the EAA payload. Port of
// getEAAIdentifierListUri().
func (w *EAAWrapper) EAAIdentifierListUri() string {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil && eaaStatus.IdentifierList() != nil {
		return payloadClaimTextValue(eaaStatus.IdentifierList().Uri())
	}
	return ""
}

// EAAIdentifierListCertificate returns a certificate containing the public key that signed or
// sealed the top-level certificate in the x5chain element in the MSO revocation list
// structure. Port of getEAAIdentifierListCertificate().
func (w *EAAWrapper) EAAIdentifierListCertificate() []byte {
	eaaStatus := w.EAAPayload().EAAStatus()
	if eaaStatus != nil && eaaStatus.IdentifierList() != nil {
		return payloadClaimByteValue(eaaStatus.IdentifierList().Certificate())
	}
	return nil
}

// EAANonce gets EAA nonce when defined in the EAA payload. Port of getEAANonce().
func (w *EAAWrapper) EAANonce() string {
	return payloadClaimTextValue(w.EAAPayload().EAANonce())
}

// EAADevicePublicKey gets EAA device public key when defined in the EAA payload. Port of
// getEAADevicePublicKey().
func (w *EAAWrapper) EAADevicePublicKey() []byte {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.PublicKey()
	}
	return nil
}

// EAADeviceCertificate gets EAA device certificate token when defined in the EAA payload. Port
// of getEAADeviceCertificate().
func (w *EAAWrapper) EAADeviceCertificate() *CertificateWrapper {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		certificates := eaaDeviceKey.Certificates()
		if len(certificates) != 0 {
			return certificates[0]
		}
	}
	return nil
}

// EAADeviceCertificateChain gets EAA device certificate chain when defined in the EAA payload.
// Port of getEAADeviceCertificateChain().
func (w *EAAWrapper) EAADeviceCertificateChain() []*CertificateWrapper {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.Certificates()
	}
	return nil
}

// EAADeviceCertificateChainDigests gets EAA device certificate chain digests when defined in
// the EAA payload. Port of getEAADeviceCertificateChainDigests().
func (w *EAAWrapper) EAADeviceCertificateChainDigests() []*jaxb.XmlDigestAlgoAndValue {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.CertificateDigests()
	}
	return nil
}

// EAADeviceCertificateKIDs gets EAA device certificate chain KIDs when defined in the EAA
// payload. Port of getEAADeviceCertificateKIDs().
func (w *EAAWrapper) EAADeviceCertificateKIDs() []string {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.KIDs()
	}
	return nil
}

// EAADeviceCertificateUrls gets EAA device certificate chain location URLs when defined in the
// EAA payload. Port of getEAADeviceCertificateUrls().
func (w *EAAWrapper) EAADeviceCertificateUrls() []string {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.X509URLs()
	}
	return nil
}

// EAADeviceKeyAuthorizedNamespaces gets a list of namespaces authorized for the device key to
// sign or MAC. Port of getEAADeviceKeyAuthorizedNamespaces().
func (w *EAAWrapper) EAADeviceKeyAuthorizedNamespaces() []string {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.AuthorizedNamespaces()
	}
	return nil
}

// EAADeviceKeyAuthorizedDataElements gets a map of namespaces and corresponding data elements
// the key is authorized to sign or MAC. Port of getEAADeviceKeyAuthorizedDataElements().
func (w *EAAWrapper) EAADeviceKeyAuthorizedDataElements() map[string][]string {
	eaaDeviceKey := w.EAAPayload().EAADeviceKey()
	if eaaDeviceKey != nil {
		return eaaDeviceKey.AuthorizedDataElements()
	}
	return nil
}

// EAARevocations returns a list of statuses for the EAA. Port of getEAARevocations().
func (w *EAAWrapper) EAARevocations() []*EAARevocationWrapper {
	var revocationWrappers []*EAARevocationWrapper
	for _, xmlEAARevocationStatus := range w.eaa.EAARevocations {
		revocationWrappers = append(revocationWrappers, NewEAARevocationWrapper(xmlEAARevocationStatus))
	}
	return revocationWrappers
}

// EAAVersion gets a version of the MobileSecurityObject. Port of getEAAVersion().
func (w *EAAWrapper) EAAVersion() string {
	return payloadClaimTextValue(w.EAAPayload().EAAVersion())
}

// HolderFullName gets user's full name when defined within EAA Payload claims. Port of
// getHolderFullName().
func (w *EAAWrapper) HolderFullName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderFullName())
}

// HolderGivenName gets user's first name when defined within EAA Payload claims. Port of
// getHolderGivenName().
func (w *EAAWrapper) HolderGivenName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderGivenName())
}

// HolderFamilyName gets user's last or family name when defined within EAA Payload claims.
// Port of getHolderFamilyName().
func (w *EAAWrapper) HolderFamilyName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderFamilyName())
}

// HolderMiddleName gets user's middle name when defined within EAA Payload claims. Port of
// getHolderMiddleName().
func (w *EAAWrapper) HolderMiddleName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderMiddleName())
}

// HolderNickname gets user's alternative name when defined within EAA Payload claims. Port of
// getHolderNickname().
func (w *EAAWrapper) HolderNickname() string {
	return payloadClaimTextValue(w.EAAPayload().HolderNickname())
}

// HolderShortName gets user's preferred or short name when defined within EAA Payload claims.
// Port of getHolderShortName().
func (w *EAAWrapper) HolderShortName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderShortName())
}

// HolderProfileUrl gets user's profile URL when defined within EAA Payload claims. Port of
// getHolderProfileUrl().
func (w *EAAWrapper) HolderProfileUrl() string {
	return payloadClaimTextValue(w.EAAPayload().HolderProfileUrl())
}

// HolderPictureUrl gets user's picture URL when defined within EAA Payload claims. Port of
// getHolderPictureUrl().
func (w *EAAWrapper) HolderPictureUrl() string {
	return payloadClaimTextValue(w.EAAPayload().HolderPictureUrl())
}

// HolderWebsiteUrl gets user's website when defined within EAA Payload claims. Port of
// getHolderWebsiteUrl().
func (w *EAAWrapper) HolderWebsiteUrl() string {
	return payloadClaimTextValue(w.EAAPayload().HolderWebsiteUrl())
}

// HolderEmail gets user's email when defined within EAA Payload claims. Port of
// getHolderEmail().
func (w *EAAWrapper) HolderEmail() string {
	return payloadClaimTextValue(w.EAAPayload().HolderEmail())
}

// HolderEmailVerified gets whether the user's website has been verified if defined within EAA
// Payload claims. Port of getHolderEmailVerified().
func (w *EAAWrapper) HolderEmailVerified() *bool {
	return payloadClaimBooleanValue(w.EAAPayload().HolderEmailVerified())
}

// HolderGender gets user's gender when defined within EAA Payload claims. Port of
// getHolderGender().
func (w *EAAWrapper) HolderGender() *int {
	return payloadClaimIntegerValue(w.EAAPayload().HolderGender())
}

// HolderBirthdate gets user's birthdate when defined within EAA Payload claims. Port of
// getHolderBirthdate().
func (w *EAAWrapper) HolderBirthdate() *time.Time {
	if w.EAAPayload().HolderBirthdate() != nil {
		return payloadClaimDateValue(w.EAAPayload().HolderBirthdate().Birthdate())
	}
	return nil
}

// HolderBirthdateApproximateMask gets an 8 digit flag to denote the location of the mask in
// YYYYMMDD format within the user's birthdate. Port of getHolderBirthdateApproximateMask().
func (w *EAAWrapper) HolderBirthdateApproximateMask() string {
	if w.EAAPayload().HolderBirthdate() != nil {
		return payloadClaimTextValue(w.EAAPayload().HolderBirthdate().ApproximateMask())
	}
	return ""
}

// HolderTimezone gets user's timezone when defined within EAA Payload claims. Port of
// getHolderTimezone().
func (w *EAAWrapper) HolderTimezone() string {
	return payloadClaimTextValue(w.EAAPayload().HolderTimezone())
}

// HolderLocale gets user's locale when defined within EAA Payload claims. Port of
// getHolderLocale().
func (w *EAAWrapper) HolderLocale() string {
	return payloadClaimTextValue(w.EAAPayload().HolderLocale())
}

// HolderPostalAddress gets user's full postal address, formatted, when defined within EAA
// Payload claims. Port of getHolderPostalAddress().
func (w *EAAWrapper) HolderPostalAddress() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.PostalAddress())
	}
	return ""
}

// HolderAddressCity gets user's city address when defined within EAA Payload claims. Port of
// getHolderAddressCity().
func (w *EAAWrapper) HolderAddressCity() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.City())
	}
	residentAddressCity := w.EAAPayload().HolderResidentAddressCity()
	if residentAddressCity != nil {
		return payloadClaimTextValue(residentAddressCity)
	}
	return ""
}

// HolderAddressStateOrProvince gets user's state or region address when defined within EAA
// Payload claims. Port of getHolderAddressStateOrProvince().
func (w *EAAWrapper) HolderAddressStateOrProvince() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.StateOrProvince())
	}
	residentAddressState := w.EAAPayload().HolderResidentAddressState()
	if residentAddressState != nil {
		return payloadClaimTextValue(residentAddressState)
	}
	return ""
}

// HolderAddressPostalCode gets user's postal code address when defined within EAA Payload
// claims. Port of getHolderAddressPostalCode().
func (w *EAAWrapper) HolderAddressPostalCode() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.PostalCode())
	}
	residentAddressPostalCode := w.EAAPayload().HolderResidentAddressPostalCode()
	if residentAddressPostalCode != nil {
		return payloadClaimTextValue(residentAddressPostalCode)
	}
	return ""
}

// HolderAddressCountry gets user's country address when defined within EAA Payload claims.
// NOTE: the returned value is usually represented by 2-letter ISO country code. Port of
// getHolderAddressCountry().
func (w *EAAWrapper) HolderAddressCountry() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.Country())
	}
	residentAddressCountry := w.EAAPayload().HolderResidentAddressCountry()
	if residentAddressCountry != nil {
		return payloadClaimTextValue(residentAddressCountry)
	}
	return ""
}

// HolderStreetAddress gets user's street address when defined within EAA Payload claims. Port
// of getHolderStreetAddress().
func (w *EAAWrapper) HolderStreetAddress() string {
	userAddress := w.EAAPayload().HolderAddress()
	if userAddress != nil {
		return payloadClaimTextValue(userAddress.StreetAddress())
	}
	postalAddress := w.EAAPayload().ResidentPostalAddress()
	if postalAddress != nil {
		return payloadClaimTextValue(postalAddress)
	}
	return ""
}

// HolderPhoneNumber gets user's phone number when defined within EAA Payload claims. Port of
// getHolderPhoneNumber().
func (w *EAAWrapper) HolderPhoneNumber() string {
	return payloadClaimTextValue(w.EAAPayload().HolderPhoneNumber())
}

// HolderPhoneNumberVerified gets whether the user's phone number has been verified if defined
// within EAA Payload claims. Port of getHolderPhoneNumberVerified().
func (w *EAAWrapper) HolderPhoneNumberVerified() *bool {
	return payloadClaimBooleanValue(w.EAAPayload().HolderPhoneNumberVerified())
}

// HolderPlaceOfBirth gets user's complete place of birth when defined within EAA Payload
// claims. Port of getHolderPlaceOfBirth().
func (w *EAAWrapper) HolderPlaceOfBirth() string {
	userPlaceOfBirth := w.EAAPayload().HolderPlaceOfBirth()
	if userPlaceOfBirth != nil && userPlaceOfBirth.IsText() {
		return payloadClaimTextValue(userPlaceOfBirth.AsClaim())
	}
	return ""
}

// HolderPlaceOfBirthCountry gets user's country of birth when defined within EAA Payload
// claims. Port of getHolderPlaceOfBirthCountry().
func (w *EAAWrapper) HolderPlaceOfBirthCountry() string {
	userPlaceOfBirth := w.EAAPayload().HolderPlaceOfBirth()
	if userPlaceOfBirth != nil {
		return payloadClaimTextValue(userPlaceOfBirth.Country())
	}
	return ""
}

// HolderPlaceOfBirthRegion gets user's state or region of birth when defined within EAA
// Payload claims. Port of getHolderPlaceOfBirthRegion().
func (w *EAAWrapper) HolderPlaceOfBirthRegion() string {
	userPlaceOfBirth := w.EAAPayload().HolderPlaceOfBirth()
	if userPlaceOfBirth != nil {
		return payloadClaimTextValue(userPlaceOfBirth.Region())
	}
	return ""
}

// HolderPlaceOfBirthCity gets user's city of birth when defined within EAA Payload claims.
// Port of getHolderPlaceOfBirthCity().
func (w *EAAWrapper) HolderPlaceOfBirthCity() string {
	userPlaceOfBirth := w.EAAPayload().HolderPlaceOfBirth()
	if userPlaceOfBirth != nil {
		return payloadClaimTextValue(userPlaceOfBirth.City())
	}
	return ""
}

// HolderNationalities gets user's nationalities list when defined within EAA Payload claims.
// NOTE: the values are usually represented by 3-letter nationality codes. Port of
// getHolderNationalities().
func (w *EAAWrapper) HolderNationalities() []string {
	return payloadClaimArrayAsStringsValue(w.EAAPayload().HolderNationalities())
}

// HolderBirthFamilyName gets user's last or family name at birth when defined within EAA
// Payload claims. Port of getHolderBirthFamilyName().
func (w *EAAWrapper) HolderBirthFamilyName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderBirthFamilyName())
}

// HolderBirthGivenName gets user's first name at birth when defined within EAA Payload claims.
// Port of getHolderBirthGivenName().
func (w *EAAWrapper) HolderBirthGivenName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderBirthGivenName())
}

// HolderBirthMiddleName gets user's middle name at birth when defined within EAA Payload
// claims. Port of getHolderBirthMiddleName().
func (w *EAAWrapper) HolderBirthMiddleName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderBirthMiddleName())
}

// HolderSalutation gets user's preferred salutation when defined within EAA Payload claims.
// Port of getHolderSalutation().
func (w *EAAWrapper) HolderSalutation() string {
	return payloadClaimTextValue(w.EAAPayload().HolderSalutation())
}

// HolderBirthFullName gets the name(s) which holder was born. Port of getHolderBirthFullName().
func (w *EAAWrapper) HolderBirthFullName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderBirthFullName())
}

// HolderTitle gets user's title when defined within EAA Payload claims. Port of
// getHolderTitle().
func (w *EAAWrapper) HolderTitle() string {
	return payloadClaimTextValue(w.EAAPayload().HolderTitle())
}

// HolderMobilePhoneNumber gets user's mobile phone number when defined within EAA Payload
// claims. Port of getHolderMobilePhoneNumber().
func (w *EAAWrapper) HolderMobilePhoneNumber() string {
	return payloadClaimTextValue(w.EAAPayload().HolderMobilePhoneNumber())
}

// HolderPseudonym gets user's scenic name or pseudonym, they are known as, when defined within
// EAA Payload claims. Port of getHolderPseudonym().
func (w *EAAWrapper) HolderPseudonym() string {
	return payloadClaimTextValue(w.EAAPayload().HolderPseudonym())
}

/* mdoc claims */

// DocumentIssuingAuthority gets issuing authority name. The value shall only use latin1
// characters and shall have a maximum length of 150 characters. Port of
// getDocumentIssuingAuthority().
func (w *EAAWrapper) DocumentIssuingAuthority() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentIssuingAuthority())
}

// DocumentIssuingAuthorityCountry gets alpha-2 country code, as defined in ISO 3166-1, of the
// issuing authority's country or territory. Port of getDocumentIssuingAuthorityCountry().
func (w *EAAWrapper) DocumentIssuingAuthorityCountry() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentIssuingAuthorityCountry())
}

// DocumentIssuingAuthorityJurisdiction gets a country subdivision code of the jurisdiction
// that issued the mDL as defined in ISO 3166-2:2020, Clause 8. The first part of the code
// shall be the same as the value for issuing_country. Port of
// getDocumentIssuingAuthorityJurisdiction().
func (w *EAAWrapper) DocumentIssuingAuthorityJurisdiction() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentIssuingAuthorityJurisdiction())
}

// PersonalAdministrativeNumber returns an audit control number assigned by the issuing
// authority. The value shall only use latin1 characters and shall have a maximum length of
// 150 characters. Port of getPersonalAdministrativeNumber().
func (w *EAAWrapper) PersonalAdministrativeNumber() string {
	return payloadClaimTextValue(w.EAAPayload().PersonalAdministrativeNumber())
}

// DocumentIssuingAuthorityCountryUNDistinguishingSign gets the distinguishing sign of the
// issuing country according to ISO/IEC 18013-1:2018, Annex F. If no applicable distinguishing
// sign is available in ISO/IEC 18013-1, an IA may use an empty identifier or another
// identifier by which it is internationally recognized. In this case the IA should ensure
// there is no collision with other IA's. Port of
// getDocumentIssuingAuthorityCountryUNDistinguishingSign().
func (w *EAAWrapper) DocumentIssuingAuthorityCountryUNDistinguishingSign() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentIssuingAuthorityUNDistinguishingSign())
}

// DocumentNumber gets the number assigned or calculated by the issuing authority. The value
// shall only use latin1 characters and shall have a maximum length of 150 characters. Port of
// getDocumentNumber().
func (w *EAAWrapper) DocumentNumber() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentNumber())
}

// DocumentType gets the document type. Port of getDocumentType().
func (w *EAAWrapper) DocumentType() string {
	return payloadClaimTextValue(w.EAAPayload().DocumentType())
}

// HolderPortrait gets a reproduction of the mDL holder's portrait. Port of
// getHolderPortrait().
func (w *EAAWrapper) HolderPortrait() []byte {
	return payloadClaimByteValue(w.EAAPayload().HolderPortrait())
}

// HolderDrivingPrivileges gets the categories of vehicles/restrictions/conditions contain
// information describing the driving privileges of the mDL holder. Port of
// getHolderDrivingPrivileges().
func (w *EAAWrapper) HolderDrivingPrivileges() *DrivingPrivilegesClaimWrapper {
	return w.EAAPayload().HolderDrivingPrivileges()
}

// HolderHeight gets the holder's height in centimetres. Port of getHolderHeight().
func (w *EAAWrapper) HolderHeight() *int {
	return payloadClaimIntegerValue(w.EAAPayload().HolderHeight())
}

// HolderWeight gets the holder's weight. Port of getHolderWeight().
func (w *EAAWrapper) HolderWeight() *int {
	return payloadClaimIntegerValue(w.EAAPayload().HolderWeight())
}

// HolderEyeColour gets the mDL holder's eye colour. The value shall be one of the following:
// "black", "blue", "brown", "dichromatic", "grey", "green", "hazel", "maroon", "pink",
// "unknown". Port of getHolderEyeColour().
func (w *EAAWrapper) HolderEyeColour() string {
	return payloadClaimTextValue(w.EAAPayload().HolderEyeColour())
}

// HolderHairColour gets the mDL holder's hair colour. The value shall be one of the following:
// "bald", "black", "blond", "brown", "grey", "red", "auburn", "sandy", "white", "unknown".
// Port of getHolderHairColour().
func (w *EAAWrapper) HolderHairColour() string {
	return payloadClaimTextValue(w.EAAPayload().HolderHairColour())
}

// HolderPortraitCaptureDate gets the date when portrait was taken. Port of
// getHolderPortraitCaptureDate().
func (w *EAAWrapper) HolderPortraitCaptureDate() *time.Time {
	return payloadClaimDateValue(w.EAAPayload().HolderPortraitCaptureDate())
}

// HolderAgeInYears gets the date the age of the mDL holder. Port of getHolderAgeInYears().
func (w *EAAWrapper) HolderAgeInYears() *int {
	return payloadClaimIntegerValue(w.EAAPayload().HolderAgeInYears())
}

// HolderAgeBirthYear gets the year when the mDL holder was born. Port of
// getHolderAgeBirthYear().
func (w *EAAWrapper) HolderAgeBirthYear() *int {
	return payloadClaimIntegerValue(w.EAAPayload().HolderAgeBirthYear())
}

// IsHolderAgeOver returns the claim value whether the age of an EAA's holder is over the
// given age. NOTE: if there is no claim provided for the requested age, nil is returned. Port
// of isHolderAgeOver(int).
func (w *EAAWrapper) IsHolderAgeOver(age int) *bool {
	holderAgeEqualOrOver := w.EAAPayload().HolderAgeEqualOrOver()
	if holderAgeEqualOrOver != nil {
		result := isHolderAgeOverList(holderAgeEqualOrOver.AgeEqualOrOverList(), age)
		if result != nil {
			return result
		}
	}
	return isHolderAgeOverList(w.EAAPayload().HolderAgeOverList(), age)
}

func isHolderAgeOverList(ageOverClaimsList []*AgeOverNNClaimWrapper, age int) *bool {
	if len(ageOverClaimsList) != 0 {
		for _, ageOverNNClaim := range ageOverClaimsList {
			if age == ageOverNNClaim.Age() {
				return payloadClaimBooleanValue(ageOverNNClaim.AsClaim())
			}
		}
	}
	return nil
}

// HolderBiometricTemplate returns the biometric template for the given value. The list of
// supported values is defined in ISO/IEC 18013-2:2020. NOTE: if there is no claim provided
// for the requested type, nil is returned. Port of getHolderBiometricTemplate(String).
func (w *EAAWrapper) HolderBiometricTemplate(templateType string) []byte {
	if templateType == "" {
		return nil
	}
	/*
	 * A biometric template identifier has the format biometric_template_xx
	 * where xx shall be replaced with the corresponding "Abstract value name" found in ISO/IEC 19785-3:2020,
	 * Table 7, according to the following convention: capitalized characters are replaced with their
	 * lowercase equivalent and spaces or non-alphanumeric characters are replaced by underscores (_).
	 */
	templateType = normalizeType(templateType)
	biometricTemplateList := w.EAAPayload().HolderBiometricTemplateList()
	for _, biometricTemplate := range biometricTemplateList {
		if templateType == normalizeType(biometricTemplate.Type()) {
			return payloadClaimByteValue(biometricTemplate.AsClaim())
		}
	}
	return nil
}

var nonAlphanumericTypeRE = regexp.MustCompile(`[^\p{L}\p{Nd}]+`)

func normalizeType(t string) string {
	if t == "" {
		return ""
	}
	t = strings.ToLower(t)
	return nonAlphanumericTypeRE.ReplaceAllString(t, "_")
}

// HolderSignatureUsualMark gets an image of the signature or usual mark of the mDL holder, see
// 7.2.7 ISO/IEC 18013-5. Port of getHolderSignatureUsualMark().
func (w *EAAWrapper) HolderSignatureUsualMark() []byte {
	return payloadClaimByteValue(w.EAAPayload().HolderSignatureUsualMark())
}

// HolderFingerprint gets a reproduction of the holder's fingerprint data (TBC). Port of
// getHolderFingerprint().
func (w *EAAWrapper) HolderFingerprint() []byte {
	return payloadClaimByteValue(w.EAAPayload().HolderFingerprint())
}

// HolderBusinessName gets a business name of the holder. Port of getHolderBusinessName().
func (w *EAAWrapper) HolderBusinessName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderBusinessName())
}

// HolderOrganizationName gets a name of legal person. Port of getHolderOrganizationName().
func (w *EAAWrapper) HolderOrganizationName() string {
	return payloadClaimTextValue(w.EAAPayload().HolderOrganizationName())
}

// HolderProfession gets the profession of the holder. Port of getHolderProfession().
func (w *EAAWrapper) HolderProfession() string {
	return payloadClaimTextValue(w.EAAPayload().HolderProfession())
}

// HolderRelationshipFather gets the father of the holder. Port of
// getHolderRelationshipFather().
func (w *EAAWrapper) HolderRelationshipFather() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipFather())
}

// HolderRelationshipMother gets the mother of the holder. Port of
// getHolderRelationshipMother().
func (w *EAAWrapper) HolderRelationshipMother() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipMother())
}

// HolderRelationshipParent gets the parent of the holder. Port of
// getHolderRelationshipParent().
func (w *EAAWrapper) HolderRelationshipParent() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipParent())
}

// HolderRelationshipSon gets the son of the holder. Port of getHolderRelationshipSon().
func (w *EAAWrapper) HolderRelationshipSon() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipSon())
}

// HolderRelationshipDaughter gets the daughter of the holder. Port of
// getHolderRelationshipDaughter().
func (w *EAAWrapper) HolderRelationshipDaughter() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipDaughter())
}

// HolderRelationshipBrother gets the brother of the holder. Port of
// getHolderRelationshipBrother().
func (w *EAAWrapper) HolderRelationshipBrother() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipBrother())
}

// HolderRelationshipSister gets the sister of the holder. Port of
// getHolderRelationshipSister().
func (w *EAAWrapper) HolderRelationshipSister() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipSister())
}

// HolderRelationshipSibling gets the sibling of the holder. Port of
// getHolderRelationshipSibling().
func (w *EAAWrapper) HolderRelationshipSibling() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipSibling())
}

// HolderRelationshipSpouse gets the spouse of the holder. Port of
// getHolderRelationshipSpouse().
func (w *EAAWrapper) HolderRelationshipSpouse() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipSpouse())
}

// HolderRelationshipFatherInLaw gets the father-in-law of the holder. Port of
// getHolderRelationshipFatherInLaw().
func (w *EAAWrapper) HolderRelationshipFatherInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipFatherInLaw())
}

// HolderRelationshipMotherInLaw gets the mother-in-law of the holder. Port of
// getHolderRelationshipMotherInLaw().
func (w *EAAWrapper) HolderRelationshipMotherInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipMotherInLaw())
}

// HolderRelationshipParentInLaw gets the parent-in-law of the holder. Port of
// getHolderRelationshipParentInLaw().
func (w *EAAWrapper) HolderRelationshipParentInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipParentInLaw())
}

// HolderRelationshipSonInLaw gets the son-in-law of the holder. Port of
// getHolderRelationshipSonInLaw().
func (w *EAAWrapper) HolderRelationshipSonInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipSonInLaw())
}

// HolderRelationshipDaughterInLaw gets the daughter-in-law of the holder. Port of
// getHolderRelationshipDaughterInLaw().
func (w *EAAWrapper) HolderRelationshipDaughterInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipDaughterInLaw())
}

// HolderRelationshipChildInLaw gets the child-in-law of the holder. Port of
// getHolderRelationshipChildInLaw().
func (w *EAAWrapper) HolderRelationshipChildInLaw() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipChildInLaw())
}

// HolderRelationshipParentalAuthority gets the parental authority of the holder. Port of
// getHolderRelationshipParentalAuthority().
func (w *EAAWrapper) HolderRelationshipParentalAuthority() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipParentalAuthority())
}

// HolderRelationshipLegalRepresentative gets the legal representative of the holder. Port of
// getHolderRelationshipLegalRepresentative().
func (w *EAAWrapper) HolderRelationshipLegalRepresentative() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipLegalRepresentative())
}

// HolderRelationshipAgent gets the voluntary agent of the holder. Port of
// getHolderRelationshipAgent().
func (w *EAAWrapper) HolderRelationshipAgent() string {
	return payloadClaimTextValue(w.EAAPayload().HolderRelationshipAgent())
}

// AdministrativeIssuanceDate gets the date when the data (e.g. a PID) was issued. Port of
// getAdministrativeIssuanceDate().
func (w *EAAWrapper) AdministrativeIssuanceDate() *time.Time {
	return payloadClaimDateValue(w.EAAPayload().AdministrativeIssuanceDate())
}

// AdministrativeExpirationDate gets the date when the data (e.g. a PID) will expire. Port of
// getAdministrativeExpirationDate().
func (w *EAAWrapper) AdministrativeExpirationDate() *time.Time {
	return payloadClaimDateValue(w.EAAPayload().AdministrativeExpirationDate())
}

// TrustAnchor gets the URL at which a machine-readable version of the trust anchor to be used
// for verifying the PID can be found or looked up. Port of getTrustAnchor().
func (w *EAAWrapper) TrustAnchor() string {
	return payloadClaimTextValue(w.EAAPayload().TrustAnchor())
}

// ResidentAddressStreet gets the name of the street where the user to whom the person
// identification data relates currently resides. Port of getResidentAddressStreet().
func (w *EAAWrapper) ResidentAddressStreet() string {
	return payloadClaimTextValue(w.EAAPayload().ResidentAddressStreet())
}

// ResidentAddressHouseNumber gets the house number where the user to whom the person
// identification data relates currently resides, including any affix or suffix. Port of
// getResidentAddressHouseNumber().
func (w *EAAWrapper) ResidentAddressHouseNumber() string {
	return payloadClaimTextValue(w.EAAPayload().ResidentAddressHouseNumber())
}

// ResidentAddressCity gets the name of the city where the user to whom the person
// identification data relates currently resides. Port of getResidentAddressCity().
func (w *EAAWrapper) ResidentAddressCity() string {
	return payloadClaimTextValue(w.EAAPayload().HolderResidentAddressCity())
}

// ResidentAddressState gets the name of the state where the user to whom the person
// identification data relates currently resides. Port of getResidentAddressState().
func (w *EAAWrapper) ResidentAddressState() string {
	return payloadClaimTextValue(w.EAAPayload().HolderResidentAddressState())
}

// ResidentAddressPostalCode gets the postal code of the address where the user to whom the
// person identification data relates currently resides. Port of
// getResidentAddressPostalCode().
func (w *EAAWrapper) ResidentAddressPostalCode() string {
	return payloadClaimTextValue(w.EAAPayload().HolderResidentAddressPostalCode())
}

// ResidentAddressCountry gets the name of the country where the user to whom the person
// identification data relates currently resides. Port of getResidentAddressCountry().
func (w *EAAWrapper) ResidentAddressCountry() string {
	return payloadClaimTextValue(w.EAAPayload().HolderResidentAddressCountry())
}

// ResidentPostalAddress gets the full address where the user to whom the person
// identification data relates currently resides. Port of getResidentPostalAddress().
func (w *EAAWrapper) ResidentPostalAddress() string {
	return payloadClaimTextValue(w.EAAPayload().ResidentPostalAddress())
}

/* ETSI TS 119 472-1 "5 Implementation of EAA based on SD-JWT VC" header parameters */

// IssuingRegistrationIdentifier gets the registration identifier of the legal entity on whose
// behalf the EAA has been issued. Port of getIssuingRegistrationIdentifier().
func (w *EAAWrapper) IssuingRegistrationIdentifier() string {
	return payloadClaimTextValue(w.EAAPayload().IssuingAuthorityRegistrationIdentifier())
}

// OneTimeUse gets the signal indicating that the EAA shall be used only once, and that it
// shall not be retained for future use. Port of getOneTimeUse().
func (w *EAAWrapper) OneTimeUse() *bool {
	return payloadClaimBooleanValue(w.EAAPayload().OneTimeUse())
}

// ShortLived gets the EAA short-lived component indicating that the validity period of the
// EAA is so short that it shall not be necessary to check its revocation status. Port of
// getShortLived().
func (w *EAAWrapper) ShortLived() *bool {
	return payloadClaimBooleanValue(w.EAAPayload().ShortLived())
}

// AttestedAttributesSubjectId gets the identifier of the attribute subject, which shall
// associate the attributes to this attribute subject. Port of
// getAttestedAttributesSubjectId().
func (w *EAAWrapper) AttestedAttributesSubjectId() string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil {
		return payloadClaimTextValue(attestedAttributesSubject.SubjectId())
	}
	return ""
}

// AttestedAttributesSubjectFamilyName gets the family name of the attribute subject, which
// shall associate the attributes to this attribute subject. Port of
// getAttestedAttributesSubjectFamilyName().
func (w *EAAWrapper) AttestedAttributesSubjectFamilyName() string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil && attestedAttributesSubject.SubjectId() != nil {
		return payloadClaimTextValue(attestedAttributesSubject.SubjectId().FamilyName())
	}
	return ""
}

// AttestedAttributesSubjectGivenName gets the given name of the attribute subject, which
// shall associate the attributes to this attribute subject. Port of
// getAttestedAttributesSubjectGivenName().
func (w *EAAWrapper) AttestedAttributesSubjectGivenName() string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil && attestedAttributesSubject.SubjectId() != nil {
		return payloadClaimTextValue(attestedAttributesSubject.SubjectId().GivenName())
	}
	return ""
}

// AttestedAttributesSubjectDocumentNumber gets the document number of the attribute subject,
// which shall associate the attributes to this attribute subject. Port of
// getAttestedAttributesSubjectDocumentNumber().
func (w *EAAWrapper) AttestedAttributesSubjectDocumentNumber() string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil && attestedAttributesSubject.SubjectId() != nil {
		return payloadClaimTextValue(attestedAttributesSubject.SubjectId().DocumentNumber())
	}
	return ""
}

// AttestedAttributesSubjectPseudonym gets the claim for associating a set of attributes to
// one entity different than the EAA subject. Port of
// getAttestedAttributesSubjectPseudonym().
func (w *EAAWrapper) AttestedAttributesSubjectPseudonym() string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil {
		return payloadClaimTextValue(attestedAttributesSubject.SubjectPseudonym())
	}
	return ""
}

// AttestedAttributes gets the attributes associated to the attribute subject whose identifier
// appears in the sub_id member or whose pseudonym appears in the sub_aka member. Port of
// getAttestedAttributes().
func (w *EAAWrapper) AttestedAttributes() []string {
	attestedAttributesSubject := w.EAAPayload().AttestedAttributesSubject()
	if attestedAttributesSubject != nil {
		return payloadClaimArrayAsStringsValue(attestedAttributesSubject.Attributes())
	}
	return nil
}

// OtherClaims gets a list of claims incorporated within the EAA Payload or provided as
// disclosures, which are not (yet) directly supported by the implementation. Port of
// getOtherClaims().
func (w *EAAWrapper) OtherClaims() []*ClaimWrapper {
	return w.EAAPayload().OtherClaims()
}

// ClaimByHeaderName returns a claim using the header name used within the EAA payload. Port
// of getClaimByHeaderName(String).
func (w *EAAWrapper) ClaimByHeaderName(headerName string) *ClaimWrapper {
	if headerName == "" {
		return nil
	}
	eaaPayloadClaims := w.AllEAAPayloadClaims()
	if w.eaa != nil && len(eaaPayloadClaims) != 0 {
		for _, c := range eaaPayloadClaims {
			if headerName == c.Name() {
				return c
			}
		}
	}
	return nil
}

// SelectivelyDisclosableClaims returns all claims that have been selectively disclosed and
// identified on the EAA (i.e. provided in the form of disclosures). Port of
// getSelectivelyDisclosableClaims().
func (w *EAAWrapper) SelectivelyDisclosableClaims() []*ClaimWrapper {
	var result []*ClaimWrapper
	eaaPayloadClaims := w.AllEAAPayloadClaims()
	if w.eaa != nil && len(eaaPayloadClaims) != 0 {
		for _, c := range eaaPayloadClaims {
			result = append(result, selectivelyDisclosableClaimsRecursively(c)...)
		}
	}
	return result
}

func selectivelyDisclosableClaimsRecursively(claimWrapper *ClaimWrapper) []*ClaimWrapper {
	var result []*ClaimWrapper
	if claimWrapper.IsSelectivelyDisclosable() {
		result = append(result, claimWrapper)
	}
	if claimWrapper.IsList() {
		for _, listItem := range claimWrapper.List() {
			result = append(result, selectivelyDisclosableClaimsRecursively(listItem)...)
		}
	} else if claimWrapper.IsMap() {
		for _, entryItem := range claimWrapper.Map() {
			result = append(result, selectivelyDisclosableClaimsRecursively(entryItem)...)
		}
	}
	return result
}

func payloadClaimTextValue(xmlDisclosableClaim *ClaimWrapper) string {
	if xmlDisclosableClaim == nil {
		return ""
	}
	return xmlDisclosableClaim.Text()
}

func payloadClaimNumberValue(xmlDisclosableClaim *ClaimWrapper) *big.Int {
	if xmlDisclosableClaim == nil {
		return nil
	}
	return xmlDisclosableClaim.Number()
}

func payloadClaimIntegerValue(xmlDisclosableClaim *ClaimWrapper) *int {
	bigIntegerValue := payloadClaimNumberValue(xmlDisclosableClaim)
	if bigIntegerValue != nil {
		v := int(bigIntegerValue.Int64())
		return &v
	}
	return nil
}

func payloadClaimDateValue(xmlDisclosableClaim *ClaimWrapper) *time.Time {
	if xmlDisclosableClaim == nil {
		return nil
	}
	return xmlDisclosableClaim.DateTime()
}

func payloadClaimBooleanValue(xmlDisclosableClaim *ClaimWrapper) *bool {
	if xmlDisclosableClaim == nil {
		return nil
	}
	if xmlDisclosableClaim.IsNull() {
		v := true // handle as a true flag
		return &v
	}
	v := xmlDisclosableClaim.IsBoolean() && xmlDisclosableClaim.Boolean() != nil && *xmlDisclosableClaim.Boolean()
	return &v
}

func payloadClaimArrayAsStringsValue(xmlDisclosableClaim *ClaimWrapper) []string {
	if xmlDisclosableClaim == nil {
		return nil
	}
	if xmlDisclosableClaim.IsText() {
		return []string{xmlDisclosableClaim.Text()}
	} else if xmlDisclosableClaim.IsList() {
		var result []string
		for _, item := range xmlDisclosableClaim.List() {
			if text := item.Text(); text != "" {
				result = append(result, text)
			}
		}
		return result
	}
	panic(fmt.Sprintf("Unsupported type '%T'!", xmlDisclosableClaim))
}

func payloadClaimByteValue(xmlDisclosableClaim *ClaimWrapper) []byte {
	if xmlDisclosableClaim == nil {
		return nil
	}
	return xmlDisclosableClaim.Binary()
}

// AllEAAPayloadClaims gets a list of all disclosable claims present within an EAA Payload.
// NOTE: the method retrieves claims from the root payload level only. Port of
// getAllEAAPayloadClaims().
func (w *EAAWrapper) AllEAAPayloadClaims() []*ClaimWrapper {
	return w.EAAPayload().AllEAAPayloadClaims()
}

// AllEAAPayloadClaimNames gets a list of names (keys) for all disclosable claims present
// within an EAA Payload. NOTE: the method retrieves names from the root payload level only.
// Port of getAllEAAPayloadClaimNames().
func (w *EAAWrapper) AllEAAPayloadClaimNames() []string {
	var result []string
	for _, c := range w.AllEAAPayloadClaims() {
		result = append(result, c.Name())
	}
	return result
}

// EAAType gets type of the EAA. Port of getEAAType().
func (w *EAAWrapper) EAAType() enumerations.EAAType {
	return w.eaa.EAAType
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries(). TODO: add support (per
// upstream Java comment).
func (w *EAAWrapper) Binaries() []byte {
	return nil
}
