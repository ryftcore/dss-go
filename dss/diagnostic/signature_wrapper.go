// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/SignatureWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// SignatureWrapper contains user-friendly methods to extract information from an
// jaxb.XmlSignature.
type SignatureWrapper struct {
	AbstractSignatureWrapperBase

	// signature is the wrapped XmlSignature.
	signature *jaxb.XmlSignature
}

// NewSignatureWrapper is the default constructor. Port of SignatureWrapper(XmlSignature);
// panics per Objects.requireNonNull(signature, "XmlSignature cannot be null!").
func NewSignatureWrapper(signature *jaxb.XmlSignature) *SignatureWrapper {
	if signature == nil {
		panic("XmlSignature cannot be null!")
	}
	w := &SignatureWrapper{signature: signature}
	w.InitSignatureWrapper(w)
	return w
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *SignatureWrapper) Id() string {
	if w.signature.Id != nil {
		return string(*w.signature.Id)
	}
	return ""
}

// DAIdentifier returns the signature document identifier of the signature. Port of
// getDAIdentifier().
func (w *SignatureWrapper) DAIdentifier() string {
	if w.signature.DAIdentifier != nil {
		return *w.signature.DAIdentifier
	}
	return ""
}

// DigestMatchers is the AbstractTokenProxy override. Port of getDigestMatchers().
func (w *SignatureWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return w.signature.DigestMatchers.All()
}

// MessageDigest returns the message-digest for a CMS signature. Port of getMessageDigest().
func (w *SignatureWrapper) MessageDigest() *jaxb.XmlDigestMatcher {
	for _, digestMatcher := range w.DigestMatchers() {
		if digestMatcher.Type != nil && enumerations.DigestMatcherTypeMessageDigest == enumerations.DigestMatcherType(*digestMatcher.Type) {
			return digestMatcher
		}
	}
	return nil
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of getCurrentBasicSignature().
func (w *SignatureWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	return w.signature.BasicSignature
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *SignatureWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	return w.signature.CertificateChain.All()
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *SignatureWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	return w.signature.SigningCertificate
}

// FoundCertificates is the AbstractTokenProxy override. Port of foundCertificates().
func (w *SignatureWrapper) FoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(w.signature.FoundCertificates)
}

// FoundRevocations is the AbstractTokenProxy override. Port of foundRevocations().
func (w *SignatureWrapper) FoundRevocations() *FoundRevocationsProxy {
	return NewFoundRevocationsProxy(w.signature.FoundRevocations)
}

// Filename is the AbstractSignatureWrapper override. Port of getFilename().
func (w *SignatureWrapper) Filename() string {
	if w.signature.SignatureFilename != nil {
		return *w.signature.SignatureFilename
	}
	return ""
}

// IsStructuralValidationValid gets if a structural validation of the signature is valid. Port
// of isStructuralValidationValid().
func (w *SignatureWrapper) IsStructuralValidationValid() bool {
	return w.signature.StructuralValidation != nil && w.signature.StructuralValidation.Valid
}

// StructuralValidationMessages returns structural validation error messages, when applicable.
// Port of getStructuralValidationMessages().
func (w *SignatureWrapper) StructuralValidationMessages() []string {
	if w.signature.StructuralValidation != nil {
		return w.signature.StructuralValidation.Message
	}
	return nil
}

// ClaimedSigningTime returns the claimed signing time extracted from the signature. Port of
// getClaimedSigningTime().
func (w *SignatureWrapper) ClaimedSigningTime() *time.Time {
	if w.signature.ClaimedSigningTime != nil {
		t := w.signature.ClaimedSigningTime.Time()
		return &t
	}
	return nil
}

// ExpirationTime gets the expiration time of the signature, after which it should not be
// accepted for processing. NOTE: the value is currently used only for the ETSI TS 119 411-5 TLS
// Certificate Binding signature. The maximum effective expiry time is whichever is soonest of
// this field, the longest-lived TLS certificate identified in the sigD member payload (below),
// or the notAfter time of the signing certificate. Port of getExpirationTime().
func (w *SignatureWrapper) ExpirationTime() *time.Time {
	if w.signature.ExpirationTime != nil {
		t := w.signature.ExpirationTime.Time()
		return &t
	}
	return nil
}

// ContentType returns the content type. Port of getContentType().
func (w *SignatureWrapper) ContentType() string {
	if w.signature.ContentType != nil {
		return *w.signature.ContentType
	}
	return ""
}

// MimeType returns the MimeType. Port of getMimeType().
func (w *SignatureWrapper) MimeType() string {
	if w.signature.MimeType != nil {
		return *w.signature.MimeType
	}
	return ""
}

// ContentHints returns the content hints string. Port of getContentHints().
func (w *SignatureWrapper) ContentHints() string {
	if w.signature.ContentHints != nil {
		return *w.signature.ContentHints
	}
	return ""
}

// ContentIdentifier returns the content identifier. Port of getContentIdentifier().
func (w *SignatureWrapper) ContentIdentifier() string {
	if w.signature.ContentIdentifier != nil {
		return *w.signature.ContentIdentifier
	}
	return ""
}

// IsCounterSignature gets if the current signature counter-signs another signature within the
// document. Port of isCounterSignature().
func (w *SignatureWrapper) IsCounterSignature() bool {
	return w.signature.CounterSignature != nil && *w.signature.CounterSignature
}

// IsKeyBindingSignature gets if the current signature is a key binding signature used to
// verify the authenticity of the token's holder. Port of isKeyBindingSignature().
func (w *SignatureWrapper) IsKeyBindingSignature() bool {
	return w.signature.KeyBindingSignature != nil && *w.signature.KeyBindingSignature
}

// IsSignatureDuplicated checks if the signature's Id is duplicated within the validating
// document. Port of isSignatureDuplicated().
func (w *SignatureWrapper) IsSignatureDuplicated() bool {
	return w.signature.Duplicated != nil && *w.signature.Duplicated
}

// SignatureDigestReference returns Signature Digest Reference. Port of
// getSignatureDigestReference().
func (w *SignatureWrapper) SignatureDigestReference() *jaxb.XmlSignatureDigestReference {
	return w.signature.SignatureDigestReference
}

// DataToBeSignedRepresentation returns a DataToBeSigned digest. Port of
// getDataToBeSignedRepresentation().
func (w *SignatureWrapper) DataToBeSignedRepresentation() *jaxb.XmlDigestAlgoAndValue {
	return w.signature.DataToBeSignedRepresentation
}

// TimestampList returns a list of associated timestamps. Port of getTimestampList().
func (w *SignatureWrapper) TimestampList() []*TimestampWrapper {
	var tsps []*TimestampWrapper
	for _, xmlFoundTimestamp := range w.signature.FoundTimestamps.All() {
		tsps = append(tsps, NewTimestampWrapper(xmlFoundTimestamp.Timestamp))
	}
	return tsps
}

// TimestampListByType returns a list of associated timestamps by type. Port of
// getTimestampListByType(TimestampType).
func (w *SignatureWrapper) TimestampListByType(timestampType enumerations.TimestampType) []*TimestampWrapper {
	var result []*TimestampWrapper
	for _, tsp := range w.TimestampList() {
		if timestampType == tsp.Type() {
			result = append(result, tsp)
		}
	}
	return result
}

// EvidenceRecords returns a list of EvidenceRecordWrapper associated with the signature. Port
// of getEvidenceRecords().
func (w *SignatureWrapper) EvidenceRecords() []*EvidenceRecordWrapper {
	var evidenceRecords []*EvidenceRecordWrapper
	for _, xmlFoundEvidenceRecord := range w.signature.FoundEvidenceRecords.All() {
		evidenceRecords = append(evidenceRecords, NewEvidenceRecordWrapper(xmlFoundEvidenceRecord.EvidenceRecord))
	}
	return evidenceRecords
}

// EvidenceRecordIdsList returns a list of associated evidence record identifiers. Port of
// getEvidenceRecordIdsList().
func (w *SignatureWrapper) EvidenceRecordIdsList() []string {
	var result []string
	for _, evidenceRecordWrapper := range w.EvidenceRecords() {
		result = append(result, evidenceRecordWrapper.Id())
	}
	return result
}

// EvidenceRecordTimestampIds returns identifiers of all embedded evidence record time-stamps.
// Port of getEvidenceRecordTimestampIds().
func (w *SignatureWrapper) EvidenceRecordTimestampIds() []string {
	var result []string
	for _, evidenceRecordWrapper := range w.EvidenceRecords() {
		result = append(result, evidenceRecordWrapper.TimestampIdsList()...)
	}
	return result
}

// EmbeddedEvidenceRecords returns a list of EvidenceRecordWrapper embedded within the
// signature. Port of getEmbeddedEvidenceRecords().
func (w *SignatureWrapper) EmbeddedEvidenceRecords() []*EvidenceRecordWrapper {
	var embeddedEvidenceRecords []*EvidenceRecordWrapper
	for _, evidenceRecord := range w.EvidenceRecords() {
		if evidenceRecord.IsEmbedded() {
			embeddedEvidenceRecords = append(embeddedEvidenceRecords, evidenceRecord)
		}
	}
	return embeddedEvidenceRecords
}

// IsSignatureProductionPlacePresent gets if the signature production place is claimed within
// the signature. Port of isSignatureProductionPlacePresent().
func (w *SignatureWrapper) IsSignatureProductionPlacePresent() bool {
	return w.signature.SignatureProductionPlace != nil
}

// StreetAddress returns the signature production place's street address, when present. Port of
// getStreetAddress().
func (w *SignatureWrapper) StreetAddress() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.StreetAddress != nil {
		return *w.signature.SignatureProductionPlace.StreetAddress
	}
	return ""
}

// City returns the signature production place's city, when present. Port of getCity().
func (w *SignatureWrapper) City() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.City != nil {
		return *w.signature.SignatureProductionPlace.City
	}
	return ""
}

// CountryName returns the signature production place's country name, when present. Port of
// getCountryName().
func (w *SignatureWrapper) CountryName() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.CountryName != nil {
		return *w.signature.SignatureProductionPlace.CountryName
	}
	return ""
}

// PostOfficeBoxNumber returns the signature production place's post office box number, when
// present. Port of getPostOfficeBoxNumber().
func (w *SignatureWrapper) PostOfficeBoxNumber() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.PostOfficeBoxNumber != nil {
		return *w.signature.SignatureProductionPlace.PostOfficeBoxNumber
	}
	return ""
}

// PostalCode returns the signature production place's postal code, when present. Port of
// getPostalCode().
func (w *SignatureWrapper) PostalCode() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.PostalCode != nil {
		return *w.signature.SignatureProductionPlace.PostalCode
	}
	return ""
}

// StateOrProvince returns the signature production place's state or province, when present.
// Port of getStateOrProvince().
func (w *SignatureWrapper) StateOrProvince() string {
	if w.IsSignatureProductionPlacePresent() && w.signature.SignatureProductionPlace.StateOrProvince != nil {
		return *w.signature.SignatureProductionPlace.StateOrProvince
	}
	return ""
}

// PostalAddress returns the signature production place's postal address, when present. Port of
// getPostalAddress().
func (w *SignatureWrapper) PostalAddress() []string {
	if w.IsSignatureProductionPlacePresent() {
		return w.signature.SignatureProductionPlace.PostalAddress
	}
	return nil
}

// SignatureFormat returns the signature level (format). Port of getSignatureFormat().
func (w *SignatureWrapper) SignatureFormat() enumerations.SignatureLevel {
	if w.signature.SignatureFormat != nil {
		return enumerations.SignatureLevel(*w.signature.SignatureFormat)
	}
	return ""
}

// SignatureType returns the signature media type. NOTE: currently used only in JAdES. Port of
// getSignatureType().
func (w *SignatureWrapper) SignatureType() string {
	if w.signature.SignatureType != nil {
		return *w.signature.SignatureType
	}
	return ""
}

// JWSSerializationType gets the JWS Serialization type. NOTE: JAdES only. Port of
// getJWSSerializationType().
func (w *SignatureWrapper) JWSSerializationType() enumerations.JWSSerializationType {
	if w.signature.JWSSerializationType != nil {
		return enumerations.JWSSerializationType(*w.signature.JWSSerializationType)
	}
	return ""
}

// COSESignatureType gets COSE signature structure's type, when applicable (CB-AdES only). Port
// of getCOSESignatureType().
func (w *SignatureWrapper) COSESignatureType() enumerations.COSESignatureType {
	if w.signature.COSESignatureType != nil {
		return enumerations.COSESignatureType(w.signature.COSESignatureType.Value)
	}
	return ""
}

// IsCOSETagged gets whether the COSE signature structure is tagged. Port of isCOSETagged().
func (w *SignatureWrapper) IsCOSETagged() bool {
	if w.signature.COSESignatureType != nil {
		return w.signature.COSESignatureType.Tagged != nil && *w.signature.COSESignatureType.Tagged
	}
	return false
}

// ErrorMessage returns an error message. Port of getErrorMessage().
func (w *SignatureWrapper) ErrorMessage() string {
	if w.signature.ErrorMessage != nil {
		return *w.signature.ErrorMessage
	}
	return ""
}

// IsSigningCertificateIdentified gets if a signing certificate has been unambiguously
// identified. Port of isSigningCertificateIdentified().
func (w *SignatureWrapper) IsSigningCertificateIdentified() bool {
	signingCertificate := w.SigningCertificate()
	signingCertificateReference := w.SigningCertificateReference()
	if signingCertificate != nil && signingCertificateReference != nil {
		return signingCertificateReference.IsDigestValueMatch() &&
			(!signingCertificateReference.IsIssuerSerialPresent() || signingCertificateReference.IsIssuerSerialMatch())
	}
	return false
}

// PolicyId returns the signature policy Id, when present. Port of getPolicyId().
func (w *SignatureWrapper) PolicyId() string {
	if policy := w.signature.Policy; policy != nil {
		if policy.Id != nil {
			return *policy.Id
		}
	}
	return ""
}

// IsPolicyZeroHash returns if the signature policy's hash should not be compared (zero hash is
// used). Port of isPolicyZeroHash().
func (w *SignatureWrapper) IsPolicyZeroHash() bool {
	policy := w.signature.Policy
	if policy != nil && policy.DigestAlgoAndValue != nil {
		return policy.DigestAlgoAndValue.ZeroHash != nil && *policy.DigestAlgoAndValue.ZeroHash
	}
	return false
}

// PolicyDigestAlgoAndValue returns the signature policy digest. Port of
// getPolicyDigestAlgoAndValue().
func (w *SignatureWrapper) PolicyDigestAlgoAndValue() *jaxb.XmlPolicyDigestAlgoAndValue {
	if policy := w.signature.Policy; policy != nil {
		return policy.DigestAlgoAndValue
	}
	return nil
}

// IsPolicyStorePresent checks if a SignaturePolicyStore unsigned property is present. Port of
// isPolicyStorePresent().
func (w *SignatureWrapper) IsPolicyStorePresent() bool {
	return w.signature.SignaturePolicyStore != nil
}

// PolicyStoreId gets the signature policy store id. Port of getPolicyStoreId().
func (w *SignatureWrapper) PolicyStoreId() string {
	if policyStore := w.signature.SignaturePolicyStore; policyStore != nil && policyStore.Id != nil {
		return *policyStore.Id
	}
	return ""
}

// PolicyStoreDescription gets the signature policy store description. Port of
// getPolicyStoreDescription().
func (w *SignatureWrapper) PolicyStoreDescription() string {
	if policyStore := w.signature.SignaturePolicyStore; policyStore != nil && policyStore.Description != nil {
		return *policyStore.Description
	}
	return ""
}

// PolicyStoreDigestAlgoAndValue gets the digest of a signature policy containing within the
// signature policy store. Port of getPolicyStoreDigestAlgoAndValue().
func (w *SignatureWrapper) PolicyStoreDigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	if policyStore := w.signature.SignaturePolicyStore; policyStore != nil {
		return policyStore.DigestAlgoAndValue
	}
	return nil
}

// PolicyStoreDocumentationReferences returns a signature policy store documentation
// references. Port of getPolicyStoreDocumentationReferences().
func (w *SignatureWrapper) PolicyStoreDocumentationReferences() []string {
	if policyStore := w.signature.SignaturePolicyStore; policyStore != nil {
		return policyStore.DocumentationReferences.All()
	}
	return nil
}

// PolicyStoreLocalURI returns a signature policy store local URI. Port of
// getPolicyStoreLocalURI().
func (w *SignatureWrapper) PolicyStoreLocalURI() string {
	if policyStore := w.signature.SignaturePolicyStore; policyStore != nil && policyStore.SigPolDocLocalURI != nil {
		return *policyStore.SigPolDocLocalURI
	}
	return ""
}

// IsBLevelTechnicallyValid gets if the B-level of the signature is valid. Port of
// isBLevelTechnicallyValid().
func (w *SignatureWrapper) IsBLevelTechnicallyValid() bool {
	return w.IsSignatureValid()
}

// IsThereXLevel returns if there is the X-Level within the signature. Port of isThereXLevel().
func (w *SignatureWrapper) IsThereXLevel() bool {
	return len(w.TimestampLevelX()) != 0
}

// IsXLevelTechnicallyValid gets if the X-level of the signature is valid. Port of
// isXLevelTechnicallyValid().
func (w *SignatureWrapper) IsXLevelTechnicallyValid() bool {
	return w.isAtLeastOneTimestampValid(w.TimestampLevelX())
}

// TimestampLevelX returns a list of validation-data-refs-only- and validation-data-
// time-stamps for the signature. Port of getTimestampLevelX().
func (w *SignatureWrapper) TimestampLevelX() []*TimestampWrapper {
	timestamps := w.TimestampListByType(enumerations.TimestampTypeValidationDataRefsOnlyTimestamp)
	timestamps = append(timestamps, w.TimestampListByType(enumerations.TimestampTypeValidationDataTimestamp)...)
	return timestamps
}

// IsThereALevel returns if there is the A-Level within the signature. Port of isThereALevel().
func (w *SignatureWrapper) IsThereALevel() bool {
	return len(w.ALevelTimestamps()) != 0
}

// IsALevelTechnicallyValid gets if the A-level of the signature is valid. Port of
// isALevelTechnicallyValid().
func (w *SignatureWrapper) IsALevelTechnicallyValid() bool {
	return w.isAtLeastOneTimestampValid(w.ALevelTimestamps())
}

// ALevelTimestamps returns a list of archive timestamps for the signature. Port of
// getALevelTimestamps().
func (w *SignatureWrapper) ALevelTimestamps() []*TimestampWrapper {
	timestamps := append([]*TimestampWrapper{}, w.ArchiveTimestamps()...)
	timestamps = append(timestamps, w.documentTimestamps(true)...)
	timestamps = append(timestamps, w.ContainerTimestamps()...)
	return timestamps
}

// ArchiveTimestamps returns a list of archive timestamps for the signature. Port of
// getArchiveTimestamps().
func (w *SignatureWrapper) ArchiveTimestamps() []*TimestampWrapper {
	return w.TimestampListByType(enumerations.TimestampTypeArchiveTimestamp)
}

// IsThereTLevel returns if there is the T-Level within the signature. Port of isThereTLevel().
func (w *SignatureWrapper) IsThereTLevel() bool {
	return len(w.TLevelTimestamps()) != 0
}

// IsTLevelTechnicallyValid gets if the T-level of the signature is valid. Port of
// isTLevelTechnicallyValid().
func (w *SignatureWrapper) IsTLevelTechnicallyValid() bool {
	return w.isAtLeastOneTimestampValid(w.TLevelTimestamps())
}

// TLevelTimestamps returns a list of signature timestamps for the signature. Port of
// getTLevelTimestamps().
func (w *SignatureWrapper) TLevelTimestamps() []*TimestampWrapper {
	timestamps := append([]*TimestampWrapper{}, w.SignatureTimestamps()...)
	timestamps = append(timestamps, w.DocumentTimestamps()...)
	return timestamps
}

// ContentTimestamps returns a list of content timestamps of the signature. Port of
// getContentTimestamps().
func (w *SignatureWrapper) ContentTimestamps() []*TimestampWrapper {
	timestamps := w.TimestampListByType(enumerations.TimestampTypeContentTimestamp)
	timestamps = append(timestamps, w.TimestampListByType(enumerations.TimestampTypeIndividualDataObjectsTimestamp)...)
	timestamps = append(timestamps, w.TimestampListByType(enumerations.TimestampTypeAllDataObjectsTimestamp)...)
	return timestamps
}

// AllTimestampsProducedAfterSignatureCreation returns all non-content timestamps. Port of
// getAllTimestampsProducedAfterSignatureCreation().
func (w *SignatureWrapper) AllTimestampsProducedAfterSignatureCreation() []*TimestampWrapper {
	var timestamps []*TimestampWrapper
	for _, timestampType := range enumerations.TimestampTypeValues() {
		if !timestampType.IsContentTimestamp() {
			timestamps = append(timestamps, w.TimestampListByType(timestampType)...)
		}
	}
	return timestamps
}

// SignatureTimestamps returns all signature timestamps. Port of getSignatureTimestamps().
func (w *SignatureWrapper) SignatureTimestamps() []*TimestampWrapper {
	return w.TimestampListByType(enumerations.TimestampTypeSignatureTimestamp)
}

// DocumentTimestamps returns all PDF document timestamps. Port of getDocumentTimestamps().
func (w *SignatureWrapper) DocumentTimestamps() []*TimestampWrapper {
	return w.TimestampListByType(enumerations.TimestampTypeDocumentTimestamp)
}

// ContainerTimestamps returns all container detached timestamps (used for ASiC containers).
// Port of getContainerTimestamps().
func (w *SignatureWrapper) ContainerTimestamps() []*TimestampWrapper {
	return w.TimestampListByType(enumerations.TimestampTypeContainerTimestamp)
}

// VRITimestamps returns all corresponding VRI timestamps (PAdES only). Port of
// getVRITimestamps().
func (w *SignatureWrapper) VRITimestamps() []*TimestampWrapper {
	return w.TimestampListByType(enumerations.TimestampTypeVRITimestamp)
}

// documentTimestamps is the private helper backing ALevelTimestamps/TLevelTimestamps. Port of
// the private getDocumentTimestamps(boolean).
func (w *SignatureWrapper) documentTimestamps(coversLTLevel bool) []*TimestampWrapper {
	var timestampWrappers []*TimestampWrapper
	for _, timestampWrapper := range w.DocumentTimestamps() {
		if coversLTLevel == w.coversLTLevel(timestampWrapper) {
			timestampWrappers = append(timestampWrappers, timestampWrapper)
		}
	}
	return timestampWrappers
}

// coversLTLevel is the private helper backing documentTimestamps. Port of the private
// coversLTLevel(TimestampWrapper).
func (w *SignatureWrapper) coversLTLevel(timestampWrapper *TimestampWrapper) bool {
	if enumerations.ArchiveTimestampTypePAdES == timestampWrapper.ArchiveTimestampType() {
		signatureCertificateChain := w.CertificateChain()
		relatedRevocationData := w.FoundRevocations().RelatedRevocationData()
		if len(relatedRevocationData) == 0 {
			return w.coversDSSCertificateDataForCertificateChain(timestampWrapper, signatureCertificateChain)
		}
		return w.coversRevocationDataForCertificateChain(timestampWrapper, signatureCertificateChain) &&
			(w.coversTimestampTokens(timestampWrapper, w.TimestampList()) || w.coversOwnRevocationData(timestampWrapper))
	}
	return false
}

// coversDSSCertificateDataForCertificateChain is the private helper backing coversLTLevel. Port
// of the private coversDSSCertificateDataForCertificateChain(TimestampWrapper,
// List<CertificateWrapper>).
func (w *SignatureWrapper) coversDSSCertificateDataForCertificateChain(timestampWrapper *TimestampWrapper, certificateChain []*CertificateWrapper) bool {
	var dssCertificates []*CertificateWrapper
	for _, c := range w.FoundCertificates().RelatedCertificatesByOrigin(enumerations.CertificateOriginDSSDictionary) {
		dssCertificates = append(dssCertificates, &c.CertificateWrapper)
	}
	for _, c := range w.FoundCertificates().RelatedCertificatesByOrigin(enumerations.CertificateOriginVRIDictionary) {
		dssCertificates = append(dssCertificates, &c.CertificateWrapper)
	}
	if len(dssCertificates) != 0 {
		timestampedCertificates := timestampWrapper.TimestampedCertificates()
		for _, certificateWrapper := range certificateChain {
			if certificateWrapperContains(dssCertificates, certificateWrapper) && certificateWrapperContains(timestampedCertificates, certificateWrapper) {
				return true
			}
		}
	}
	return false
}

// coversRevocationDataForCertificateChain is the private helper backing coversLTLevel and
// coversOwnRevocationData. Port of the private coversRevocationDataForCertificateChain(
// TimestampWrapper, List<CertificateWrapper>): Java's
// `certificateRevocationData.stream().anyMatch(timestampedRevocations::contains)` checks each
// CertificateRevocationWrapper (which extends RevocationWrapper) against the
// List<RevocationWrapper> using RevocationWrapper.equals() (Id-only, no dynamic-type check); the
// embedded RevocationWrapper field is compared here the same way via revocationWrapperListContains.
func (w *SignatureWrapper) coversRevocationDataForCertificateChain(timestampWrapper *TimestampWrapper, certificateChain []*CertificateWrapper) bool {
	timestampedRevocations := timestampWrapper.TimestampedRevocations()
	for _, certificateWrapper := range certificateChain {
		certificateRevocationData := certificateWrapper.CertificateRevocationData()
		if len(certificateRevocationData) != 0 {
			for _, revocation := range certificateRevocationData {
				if revocationWrapperListContains(timestampedRevocations, &revocation.RevocationWrapper) {
					return true
				}
			}
			return false
		}
	}
	return false
}

// coversTimestampTokens is the private helper backing coversLTLevel. Port of the private
// coversTimestampTokens(TimestampWrapper, List<TimestampWrapper>).
func (w *SignatureWrapper) coversTimestampTokens(timestamp *TimestampWrapper, timestampWrappers []*TimestampWrapper) bool {
	timestampedTimestamps := timestamp.TimestampedTimestamps()
	if timestampedTimestamps == nil {
		return false
	}
	for _, tsp := range timestampWrappers {
		if timestampWrapperContains(timestampedTimestamps, tsp) {
			return true
		}
	}
	return false
}

// coversOwnRevocationData is the private helper backing coversLTLevel. Port of the private
// coversOwnRevocationData(TimestampWrapper).
func (w *SignatureWrapper) coversOwnRevocationData(timestampWrapper *TimestampWrapper) bool {
	signingCertificate := timestampWrapper.SigningCertificate()
	if signingCertificate.IsSelfSigned() || signingCertificate.IsTrusted() {
		return true // no revocation data required
	}
	return w.coversRevocationDataForCertificateChain(timestampWrapper, timestampWrapper.CertificateChain())
}

// isAtLeastOneTimestampValid is the private helper backing the level-validity checks. Port of
// the private isAtLeastOneTimestampValid(List<TimestampWrapper>).
func (w *SignatureWrapper) isAtLeastOneTimestampValid(timestampList []*TimestampWrapper) bool {
	for _, timestamp := range timestampList {
		signatureValid := timestamp.IsSignatureValid()
		messageImprint := timestamp.MessageImprint()
		messageImprintIntact := messageImprint != nil && messageImprint.DataFound && messageImprint.DataIntact
		if signatureValid && messageImprintIntact {
			return true
		}
	}
	return false
}

// certificateWrapperContains reports whether certificate appears in list, using CertificateWrapper's
// Java equals() semantics (compare by Id, since AbstractTokenProxy.equals additionally requires
// the same dynamic type, which for two CertificateWrapper values is always satisfied here).
func certificateWrapperContains(list []*CertificateWrapper, certificate *CertificateWrapper) bool {
	for _, c := range list {
		if c.Equals(certificate) {
			return true
		}
	}
	return false
}

// revocationWrapperListContains reports whether item appears in list, using
// RevocationWrapper.Equals (Id-only) semantics.
func revocationWrapperListContains(list []*RevocationWrapper, item *RevocationWrapper) bool {
	for _, r := range list {
		if r.Equals(item) {
			return true
		}
	}
	return false
}

// timestampWrapperContains reports whether timestamp appears in list, using
// AbstractTokenProxy.equals (Id plus dynamic type) semantics.
func timestampWrapperContains(list []*TimestampWrapper, timestamp *TimestampWrapper) bool {
	for _, t := range list {
		if t.Equals(timestamp) {
			return true
		}
	}
	return false
}

// TimestampIdsList returns a list of timestamp IDs. Port of getTimestampIdsList().
func (w *SignatureWrapper) TimestampIdsList() []string {
	var result []string
	for _, tsp := range w.TimestampList() {
		result = append(result, tsp.Id())
	}
	return result
}

// IsThereERSLevel returns if there is the ERS-Level within the signature. Port of
// isThereERSLevel().
func (w *SignatureWrapper) IsThereERSLevel() bool {
	return len(w.EmbeddedEvidenceRecords()) != 0
}

// KeyIdentifierReference returns a reference extracted from a 'kid' (key identifier) header
// (used in JAdES, CB-AdES). Port of getKeyIdentifierReference().
func (w *SignatureWrapper) KeyIdentifierReference() *CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	certificateRefs = append(certificateRefs, w.FoundCertificates().RelatedCertificateRefsByRefOrigin(enumerations.CertificateRefOriginKeyIdentifier)...)
	certificateRefs = append(certificateRefs, w.FoundCertificates().OrphanCertificateRefsByRefOrigin(enumerations.CertificateRefOriginKeyIdentifier)...)
	if len(certificateRefs) != 0 {
		// only one shall be present
		return certificateRefs[0]
	}
	return nil
}

// X509UrlReferences returns a list of references extracted from a 'x5u' (X.509 URL) header
// (used in JAdES, CB-AdES). Port of getX509UrlReferences().
func (w *SignatureWrapper) X509UrlReferences() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	certificateRefs = append(certificateRefs, w.FoundCertificates().RelatedCertificateRefsByRefOrigin(enumerations.CertificateRefOriginX509URL)...)
	certificateRefs = append(certificateRefs, w.FoundCertificates().OrphanCertificateRefsByRefOrigin(enumerations.CertificateRefOriginX509URL)...)
	return certificateRefs
}

// Parent returns a master-signature in case of a counter-signature. Port of getParent().
func (w *SignatureWrapper) Parent() *SignatureWrapper {
	if w.signature.Parent != nil {
		return NewSignatureWrapper(w.signature.Parent)
	}
	return nil
}

// SignatureScopes returns Signature Scopes. Port of getSignatureScopes().
func (w *SignatureWrapper) SignatureScopes() []*jaxb.XmlSignatureScope {
	return w.signature.SignatureScopes.All()
}

// SignerRoles returns list of all found SignerRoles. Port of getSignerRoles().
func (w *SignatureWrapper) SignerRoles() []*jaxb.XmlSignerRole {
	return w.signature.SignerRole
}

// ClaimedRoles returns list of found ClaimedRoles. Port of getClaimedRoles().
func (w *SignatureWrapper) ClaimedRoles() []*jaxb.XmlSignerRole {
	return w.signerRolesByCategory(enumerations.EndorsementTypeClaimed)
}

// CertifiedRoles returns list of found CertifiedRoles. Port of getCertifiedRoles().
func (w *SignatureWrapper) CertifiedRoles() []*jaxb.XmlSignerRole {
	return w.signerRolesByCategory(enumerations.EndorsementTypeCertified)
}

// SignedAssertions returns list of all found SignedAssertions. Port of getSignedAssertions().
func (w *SignatureWrapper) SignedAssertions() []*jaxb.XmlSignerRole {
	return w.signerRolesByCategory(enumerations.EndorsementTypeSigned)
}

// SignerRoleDetails returns a list of strings describing the role for the given
// listOfSignerRoles. Port of getSignerRoleDetails(List<XmlSignerRole>).
func (w *SignatureWrapper) SignerRoleDetails(listOfSignerRoles []*jaxb.XmlSignerRole) []string {
	var roles []string
	for _, xmlSignerRole := range listOfSignerRoles {
		if xmlSignerRole.Role != nil {
			roles = append(roles, *xmlSignerRole.Role)
		} else {
			roles = append(roles, "")
		}
	}
	return roles
}

// signerRolesByCategory is the private helper backing ClaimedRoles/CertifiedRoles/
// SignedAssertions. Port of the private getSignerRolesByCategory(EndorsementType).
func (w *SignatureWrapper) signerRolesByCategory(category enumerations.EndorsementType) []*jaxb.XmlSignerRole {
	var roles []*jaxb.XmlSignerRole
	for _, xmlSignerRole := range w.SignerRoles() {
		if xmlSignerRole.Category != nil && category == enumerations.EndorsementType(*xmlSignerRole.Category) {
			roles = append(roles, xmlSignerRole)
		}
	}
	return roles
}

// CommitmentTypeIndications returns a list of commitment type indications. Port of
// getCommitmentTypeIndications().
func (w *SignatureWrapper) CommitmentTypeIndications() []*jaxb.XmlCommitmentTypeIndication {
	return w.signature.CommitmentTypeIndications.All()
}

// IsPolicyPresent checks if a SignaturePolicyIdentifier is present. Port of isPolicyPresent().
func (w *SignatureWrapper) IsPolicyPresent() bool {
	return w.signature.Policy != nil
}

// PolicyProcessingError returns an error string occurred during a SignaturePolicy proceeding,
// when applicable. Port of getPolicyProcessingError().
func (w *SignatureWrapper) PolicyProcessingError() string {
	if policy := w.signature.Policy; policy != nil && policy.ProcessingError != nil {
		return *policy.ProcessingError
	}
	return ""
}

// PolicyDescription returns XMLPolicy description if it is not empty. Port of
// getPolicyDescription().
func (w *SignatureWrapper) PolicyDescription() string {
	if policy := w.signature.Policy; policy != nil && policy.Description != nil {
		return *policy.Description
	}
	return ""
}

// PolicyDocumentationReferences returns DocumentationReferences defined for the signature
// policy. Port of getPolicyDocumentationReferences().
func (w *SignatureWrapper) PolicyDocumentationReferences() []string {
	if policy := w.signature.Policy; policy != nil {
		return policy.DocumentationReferences.All()
	}
	return nil
}

// PolicyTransforms returns a list of Policy transformations. NOTE: used only for XAdES
// signatures. Port of getPolicyTransforms().
func (w *SignatureWrapper) PolicyTransforms() []string {
	if policy := w.signature.Policy; policy != nil {
		return policy.Transformations.All()
	}
	return nil
}

// PolicyUrl returns the signature policy url. Port of getPolicyUrl().
func (w *SignatureWrapper) PolicyUrl() string {
	if policy := w.signature.Policy; policy != nil && policy.Url != nil {
		return *policy.Url
	}
	return ""
}

// PolicyUserNotice returns the policy UserNotice. Port of getPolicyUserNotice().
func (w *SignatureWrapper) PolicyUserNotice() *jaxb.XmlUserNotice {
	if policy := w.signature.Policy; policy != nil {
		return policy.UserNotice
	}
	return nil
}

// PolicyDocSpecification returns the signature policy document specification. Port of
// getPolicyDocSpecification().
func (w *SignatureWrapper) PolicyDocSpecification() *jaxb.XmlSPDocSpecification {
	if policy := w.signature.Policy; policy != nil {
		return policy.DocSpecification
	}
	return nil
}

// IsPolicyAsn1Processable gets if the signature policy is ASN.1 processable. Port of
// isPolicyAsn1Processable().
func (w *SignatureWrapper) IsPolicyAsn1Processable() bool {
	if policy := w.signature.Policy; policy != nil {
		return policy.Asn1Processable != nil && *policy.Asn1Processable
	}
	return false
}

// IsPolicyIdentified gets if the signature policy has been found. Port of
// isPolicyIdentified().
func (w *SignatureWrapper) IsPolicyIdentified() bool {
	if policy := w.signature.Policy; policy != nil {
		return policy.Identified != nil && *policy.Identified
	}
	return false
}

// IsPolicyDigestValid gets if the signature policy digest validation succeeds. Port of
// isPolicyDigestValid().
func (w *SignatureWrapper) IsPolicyDigestValid() bool {
	policy := w.signature.Policy
	if policy != nil && policy.DigestAlgoAndValue != nil {
		return policy.DigestAlgoAndValue.Match != nil && *policy.DigestAlgoAndValue.Match
	}
	return false
}

// IsPolicyDigestAlgorithmsEqual gets if the validated signature policy algorithm match. Port of
// isPolicyDigestAlgorithmsEqual().
func (w *SignatureWrapper) IsPolicyDigestAlgorithmsEqual() bool {
	policy := w.signature.Policy
	if policy != nil && policy.DigestAlgoAndValue != nil {
		return policy.DigestAlgoAndValue.DigestAlgorithmsEqual != nil && *policy.DigestAlgoAndValue.DigestAlgorithmsEqual
	}
	return false
}

// PDFRevision returns a PAdES-specific PDF Revision info. NOTE: applicable only for PAdES. Port
// of getPDFRevision().
func (w *SignatureWrapper) PDFRevision() *PDFRevisionWrapper {
	if w.signature.PDFRevision != nil {
		return NewPDFRevisionWrapper(w.signature.PDFRevision)
	}
	return nil
}

// SignatureInformationStore returns a list if Signer Infos (Signer Information Store) from
// CAdES CMS Signed Data. Port of getSignatureInformationStore().
func (w *SignatureWrapper) SignatureInformationStore() []*jaxb.XmlSignerInfo {
	return w.signature.SignerInformationStore.All()
}

// VRIDictionaryCreationTime returns time of /VRI dictionary creation, when 'TU' attribute is
// present (PAdES only). Port of getVRIDictionaryCreationTime().
func (w *SignatureWrapper) VRIDictionaryCreationTime() *time.Time {
	if w.signature.VRIDictionaryCreationTime != nil {
		t := w.signature.VRIDictionaryCreationTime.Time()
		return &t
	}
	return nil
}

// SignatureValue gets the SignatureValue. Port of getSignatureValue().
func (w *SignatureWrapper) SignatureValue() []byte {
	if w.signature.SignatureValue != nil {
		return []byte(*w.signature.SignatureValue)
	}
	return nil
}

// IsDocHashOnly gets if the signature is a document hash only. Port of isDocHashOnly().
func (w *SignatureWrapper) IsDocHashOnly() bool {
	if signerDocumentRepresentation := w.signature.SignerDocumentRepresentations; signerDocumentRepresentation != nil {
		return signerDocumentRepresentation.DocHashOnly
	}
	return false
}

// IsHashOnly gets if the signature is a hash only. Port of isHashOnly().
func (w *SignatureWrapper) IsHashOnly() bool {
	if signerDocumentRepresentation := w.signature.SignerDocumentRepresentations; signerDocumentRepresentation != nil {
		return signerDocumentRepresentation.HashOnly
	}
	return false
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *SignatureWrapper) Binaries() []byte {
	return w.SignatureValue()
}
