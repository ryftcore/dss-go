// Ported from
// dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// Performs checks according to EN 319 122-1 v1.1.1 "6.3 Requirements on components and
// services".
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug/LOG.trace call site is
// called out in the surrounding comment instead.
package cades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESBaselineRequirementsChecker checks conformance of a CAdES signature to the requested
// baseline format. Port of the class CAdESBaselineRequirementsChecker, extending
// validation.BaselineRequirementsChecker[CAdESSignature].
type CAdESBaselineRequirementsChecker struct {
	validation.BaselineRequirementsChecker[*CAdESSignature]
}

// newCAdESBaselineRequirementsChecker is used to verify conformance of a signature to
// Baseline-B level. Port of the protected CAdESBaselineRequirementsChecker(CAdESSignature)
// constructor.
func newCAdESBaselineRequirementsChecker(sig *CAdESSignature) *CAdESBaselineRequirementsChecker {
	return NewCAdESBaselineRequirementsChecker(sig, nil)
}

// NewCAdESBaselineRequirementsChecker is the default constructor.
// Port of the public CAdESBaselineRequirementsChecker(CAdESSignature, CertificateVerifier)
// constructor.
func NewCAdESBaselineRequirementsChecker(sig *CAdESSignature, offlineCertificateVerifier validation.CertificateVerifier) *CAdESBaselineRequirementsChecker {
	checker := &CAdESBaselineRequirementsChecker{
		BaselineRequirementsChecker: validation.NewBaselineRequirementsCheckerBaseWithVerifier[*CAdESSignature](sig, offlineCertificateVerifier),
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// GetBaselineSignatureForm returns the signature form corresponding to the signature: CAdES.
// Port of the protected getBaselineSignatureForm(), exported to satisfy this port's
// validation.BaselineRequirementsCheckerOverrides (see that interface's doc comment on
// GetBaselineSignatureForm for why: cmsBaselineBRequirements() below - unlike upstream, which
// resolves this through ordinary Java virtual dispatch - must reach
// pades.PAdESBaselineRequirementsChecker/pades.CMSForPAdESBaselineRequirementsChecker's own
// override across the cades/pades package boundary, and this is the only mechanism available).
func (b *CAdESBaselineRequirementsChecker) GetBaselineSignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormCAdES
}

// cmsBaselineBRequirements checks if BASELINE-B requirements satisfy for a CMS signature.
// Port of the protected cmsBaselineBRequirements().
func (b *CAdESBaselineRequirementsChecker) cmsBaselineBRequirements() bool {
	sig := b.Signature()
	cmsDoc := sig.CMS()
	signerInformation := sig.SignerInformation()

	// SignedData.certificates (Cardinality == 1)
	if utils.IsCollectionEmpty(cmsDoc.Certificates()) {
		// Upstream logs "SignedData.certificates shall be present for {}-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// content-type (Cardinality == 1)
	if !b.isContentTypeValid(signerInformation) {
		// Upstream logs "content-type attribute shall be present for {}-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// message-digest (Cardinality == 1)
	if !cadesBaselineIsMessageDigestPresent(signerInformation) {
		// Upstream logs "message-digest attribute shall be present for {}-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// signing-certificate/signing-certificate-v2 (Cardinality == 1)
	if !cadesBaselineIsOneSigningCertificatePresent(signerInformation) {
		// Upstream logs "signing-certificate(-v2) attribute shall be present for
		// {}-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// signing-time (Cardinality == 1 for CAdES, EN 319 122-1; Cardinality == 0 for PAdES,
	// EN 319 142-1 explicitly forbids it where CAdES requires it - getBaselineSignatureForm()/
	// GetBaselineSignatureForm is exactly how upstream's shared hasBaselineBProfile() (this
	// method) tells the two regimes apart, via ordinary Java virtual dispatch on whichever
	// concrete checker is running: CAdESBaselineRequirementsChecker itself (CAdES) or, when this
	// same method runs for a PDF's embedded CMS through
	// pades.CMSForPAdESBaselineRequirementsChecker (PAdES), that checker's own override. See
	// this port's spi/validation.BaselineRequirementsCheckerOverrides.GetBaselineSignatureForm
	// for the cross-package plumbing this needs in Go.
	signingTimeAttrs := CAdESUtilsSignedAttributesOfType(signerInformation, OIDPkcs9AtSigningTime)
	signingTimePresent := cadesBaselineAttributeValuesSize(signingTimeAttrs) == 1
	isCAdESForm := b.BaselineSignatureForm() == enumerations.SignatureFormCAdES
	if signingTimePresent != isCAdESForm {
		// Upstream logs "signing-time attribute shall be present for {}-BASELINE-B signature
		// (cardinality == 1})!" (CAdES) or "signing-time attribute shall not be present for
		// {}-BASELINE-B signature (cardinality == 0})!" (PAdES).
		return false
	}
	// signer-attributes (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSignerAttr))+
		utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaEtsSignerAttrV2)) > 1 {
		// Upstream logs "signer-attributes(-v2) attribute shall not be present multiple times
		// for {}-BASELINE-B signature (cardinality == 0 or 1)!".
		return false
	}
	// signature-policy-identifier (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSigPolicyId)) > 1 {
		// Upstream logs "signature-policy-identifier attribute shall not be present multiple
		// times for {}-BASELINE-B signature (cardinality == 0 or 1)!".
		return false
	}
	// Additional requirement (a)
	if !b.ContainsSigningCertificate(sig.CertificateSource().SignedDataCertificates()) {
		// Upstream logs "Signing certificate shall be present in SignedData.certificates for
		// {}-BASELINE-B signature (requirement (a))!".
		return false
	}
	// Additional requirement (f)
	if sig.ContentType() != "" && cmscore.OIDData.String() != sig.ContentType() {
		// Upstream logs "The content-type attribute shall have value id-data for
		// {}-BASELINE-B signature (requirement (f))!".
		return false
	}
	// Additional requirement (h) and (i)
	if !b.isSigningCertificateAttributeValid(signerInformation) {
		// Upstream logs "signing-certificate attribute shall be used for SHA1 hash algorithm
		// and signing-certificate-v2 for other hash algorithms for {}-BASELINE-B signature
		// (requirements (h) and (i) 319 122-1)!".
		return false
	}
	return true
}

// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile.
// Port of hasAdESProfile().
func (b *CAdESBaselineRequirementsChecker) HasAdESProfile() bool {
	return b.HasExtendedBESProfile() || b.HasBaselineBProfile()
}

// HasBaselineBProfile checks if the signature has a corresponding BASELINE-B profile.
// Port of hasBaselineBProfile().
func (b *CAdESBaselineRequirementsChecker) HasBaselineBProfile() bool {
	if !b.cmsBaselineBRequirements() {
		return false
	}
	// Additional requirement (k)
	signaturePolicyStore := b.Signature().SignaturePolicyStore()
	if signaturePolicyStore != nil && !b.IsSignaturePolicyIdentifierHashPresent() {
		// Upstream logs "signature-policy-store shall not be present for CAdES-BASELINE-B
		// signature with not defined signature-policy-identifier/sigPolicyHash (requirement
		// (k))!".
		return false
	}
	return true
}

// HasBaselineTProfile checks if the signature has a corresponding BASELINE-T profile.
// Port of hasBaselineTProfile().
func (b *CAdESBaselineRequirementsChecker) HasBaselineTProfile() bool {
	if !b.MinimalTRequirement() {
		return false
	}
	// Additional requirement (m)
	if !b.SignatureTimestampsCreatedBeforeSignCertExpiration() {
		// Upstream logs "signature-time-stamp shall be created before expiration of the
		// signing-certificate for CAdES-BASELINE-T signature (requirement (m))!".
		return false
	}
	return true
}

// HasBaselineLTProfile checks if the signature has a corresponding BASELINE-LT profile.
// Port of hasBaselineLTProfile().
func (b *CAdESBaselineRequirementsChecker) HasBaselineLTProfile() bool {
	if !b.MinimalLTRequirement() {
		return false
	}
	signerInformation := b.Signature().SignerInformation()
	// certificate-values (Cardinality == 0)
	if utils.IsArrayNotEmpty(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsCertValues)) {
		// Upstream logs "certificate-values attribute shall not be present for
		// CAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	// complete-certificate-references (Cardinality == 0)
	if utils.IsArrayNotEmpty(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsCertificateRefs)) {
		// Upstream logs "complete-certificate-references attribute shall not be present for
		// CAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	// revocation-values (Cardinality == 0)
	if utils.IsArrayNotEmpty(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsRevocationValues)) {
		// Upstream logs "revocation-values attribute shall not be present for
		// CAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	// complete-revocation-references (Cardinality == 0)
	if utils.IsArrayNotEmpty(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsRevocationRefs)) {
		// Upstream logs "complete-revocation-references attribute shall not be present for
		// CAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	// time-stamped-certs-crls-references (Cardinality == 0)
	if utils.IsArrayNotEmpty(CAdESUtilsUnsignedAttributesOfType(signerInformation, OIDIdAaEtsCertCRLTimestamp)) {
		// Upstream logs "time-stamped-certs-crls-references attribute shall not be present
		// for CAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	return true
}

// ContainsLTLevelCertificates verifies whether the signature contains some of the LT-/XL-
// level attributes. Port of the protected containsLTLevelCertificates() override.
func (b *CAdESBaselineRequirementsChecker) ContainsLTLevelCertificates() bool {
	sig := b.Signature()
	cmsDoc := sig.CMS()

	var signedDataCertificates []*model.CertificateToken
	for _, certificateDER := range cmsDoc.Certificates() {
		certificateToken, err := spi.DSSASN1UtilsCertificate(certificateDER)
		if err != nil {
			// A malformed SignedData.certificates member would make DSSASN1Utils.getCertificate
			// throw upstream, propagating uncaught through this method's boolean return; here
			// it is skipped instead (best-effort, matching every other malformed-input
			// degradation in this package rather than panicking for a condition this rare).
			continue
		}
		signedDataCertificates = append(signedDataCertificates, certificateToken)
	}

	// TimestampCertificateSourcesExceptLastArchiveTimestamp is part of the validation.TimestampSource
	// interface sig.TimestampSource() already returns, so no downcast is needed to reach it - and a
	// hard *CAdESTimestampSource assertion would be actively wrong here: Java's declared-CAdESSignature-
	// typed local still runs this method on the ACTUAL runtime object, so a PAdESSignature (whose
	// TimestampSource() is a *pades.PAdESTimestampSource, itself embedding CAdESTimestampSource
	// rather than being one) reaching this shared CAdES logic through
	// PAdESBaselineRequirementsChecker.ContainsLTLevelCertificates's delegation is exactly Java's
	// ordinary virtual dispatch, which a Go concrete-type assertion cannot reproduce. Confirmed by
	// pades/testdata/crossgen's own downstream cross-validation fixtures: a self-signed test
	// certificate chain drives MinimalLTRequirement into this exact
	// call, and the assertion below used to panic on every one of them.
	timestampListCertificateSource := sig.TimestampSource().TimestampCertificateSourcesExceptLastArchiveTimestamp()
	timestampCertificateSources := timestampListCertificateSource.Sources()
	if utils.IsCollectionEmpty(timestampCertificateSources) {
		return false
	}
	for _, timestampCertificateSource := range timestampCertificateSources {
		if !cadesBaselineContainsAnyCertificate(signedDataCertificates, timestampCertificateSource.Certificates()) {
			return false
		}
	}
	return true
}

// cadesBaselineContainsAnyCertificate reports whether any element of subset is Equals() to an
// element of superset, reproducing org.apache.commons.collections4.CollectionUtils#containsAny
// (which Java's Utils.containsAny delegates to) for a type with no comparable Go representation
// usable directly with utils.ContainsAny's `comparable` constraint (see spi/certificate_source.go's
// header for the same problem elsewhere in this port).
func cadesBaselineContainsAnyCertificate(superset, subset []*model.CertificateToken) bool {
	for _, candidate := range subset {
		for _, existing := range superset {
			if candidate.Equals(existing) {
				return true
			}
		}
	}
	return false
}

// HasBaselineLTAProfile checks if the signature has a corresponding BASELINE-LTA profile.
// Port of hasBaselineLTAProfile().
func (b *CAdESBaselineRequirementsChecker) HasBaselineLTAProfile() bool {
	sig := b.Signature()
	var timestampTokens []*validation.TimestampToken
	timestampTokens = append(timestampTokens, sig.ArchiveTimestamps()...)
	timestampTokens = append(timestampTokens, sig.DetachedTimestamps()...)
	if utils.IsCollectionEmpty(timestampTokens) {
		// Upstream logs "ArchiveTimeStamp shall be present for CAdES-BASELINE-LTA signature
		// (cardinality >= 1)!".
		return false
	}
	// archive-time-stamp-v3 / detached timestamps (Cardinality >= 1)
	validArcTstFound := false
	for _, timestampToken := range timestampTokens {
		if enumerations.ArchiveTimestampTypeCAdESV3 == timestampToken.ArchiveTimestampType() ||
			timestampToken.TimeStampType().IsContainerTimestamp() {
			validArcTstFound = true
			break
		}
	}
	if !validArcTstFound {
		// Upstream logs "archive-time-stamp-v3 attribute shall be present for
		// CAdES-BASELINE-LTA signature (cardinality == 1)!".
		return false
	}
	return true
}

// HasExtendedBESProfile checks if the signature has a corresponding *AdES-BES profile.
// Port of hasExtendedBESProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedBESProfile() bool {
	if !b.cmsExtendedBESRequirements() {
		return false
	}
	signerInformation := b.Signature().SignerInformation()

	// signing-time (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDPkcs9AtSigningTime)) > 1 {
		// Upstream logs "signing-time attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// commitment-time-indication (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsCommitmentType)) > 1 {
		// Upstream logs "commitment-time-indication attribute shall not be present multiple
		// times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// content-hints (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaContentHint)) > 1 {
		// Upstream logs "content-hints attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// mime-type (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaEtsMimeType)) > 1 {
		// Upstream logs "mime-type attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// signer-location (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSignerLocation)) > 1 {
		// Upstream logs "signer-location attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// signature-policy-identifier (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSigPolicyId)) > 1 {
		// Upstream logs "signature-policy-identifier attribute shall not be present multiple
		// times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// signature-policy-store (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsSigPolicyStore)) > 1 {
		// Upstream logs "signature-policy-store attribute shall not be present multiple times
		// for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// content-reference (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaContentReference)) > 1 {
		// Upstream logs "content-reference attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// content-identifier (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaContentIdentifier)) > 1 {
		// Upstream logs "content-identifier attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// complete-certificate-references (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsCertificateRefs)) > 1 {
		// Upstream logs "complete-certificate-references attribute shall not be present
		// multiple times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// complete-revocation-references (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsRevocationRefs)) > 1 {
		// Upstream logs "complete-revocation-references attribute shall not be present
		// multiple times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// attribute-certificate-references (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDAttributeCertificateRefsOid)) > 1 {
		// Upstream logs "attribute-certificate-references attribute shall not be present
		// multiple times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// attribute-revocation-references (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDAttributeRevocationRefsOid)) > 1 {
		// Upstream logs "attribute-revocation-references attribute shall not be present
		// multiple times for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// certificate-values (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsCertValues)) > 1 {
		// Upstream logs "certificate-values attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// revocation-values (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsRevocationValues)) > 1 {
		// Upstream logs "revocation-values attribute shall not be present multiple times for
		// {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	// Additional requirement (h) and (i)
	if !b.isSigningCertificateAttributeValid(signerInformation) {
		// Upstream logs "signing-certificate attribute shall be used for SHA1 hash algorithm
		// and signing-certificate-v2 for other hash algorithms for {}-BES signature
		// (requirements (a) and (b) 319 122-2)!".
		return false
	}
	return true
}

// cmsExtendedBESRequirements verifies whether CMS is conformant to the E-BES profile.
// Port of the protected cmsExtendedBESRequirements().
func (b *CAdESBaselineRequirementsChecker) cmsExtendedBESRequirements() bool {
	signerInformation := b.Signature().SignerInformation()

	// content-type (Cardinality == 1)
	if !b.isContentTypeValid(signerInformation) {
		// Upstream logs "content-type attribute shall be present for {}-BES signature
		// (cardinality == 1)!".
		return false
	}
	// message-digest (Cardinality == 1)
	if !cadesBaselineIsMessageDigestPresent(signerInformation) {
		// Upstream logs "message-digest attribute shall be present for {}-BES signature
		// (cardinality == 1)!".
		return false
	}
	// signing-certificate/signing-certificate-v2 (Cardinality == 1)
	if !cadesBaselineIsOneSigningCertificatePresent(signerInformation) {
		// Upstream logs "signing-certificate(-v2) attribute shall be present for {}-BES
		// signature (cardinality == 1)!".
		return false
	}
	// signer-attributes (Cardinality == 0 or 1)
	if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSignerAttr))+
		utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaEtsSignerAttrV2)) > 1 {
		// Upstream logs "signer-attributes(-v2) attribute shall not be present multiple times
		// for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	return true
}

// HasExtendedEPESProfile checks if the signature has a corresponding *AdES-EPES profile.
// Port of hasExtendedEPESProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedEPESProfile() bool {
	sig := b.Signature()
	signerInformation := sig.SignerInformation()

	// signature-policy-identifier (Cardinality == 1)
	sigPolicyIdAttrs := CAdESUtilsSignedAttributesOfType(signerInformation, OIDIdAaEtsSigPolicyId)
	if cadesBaselineAttributeValuesSize(sigPolicyIdAttrs) == 0 {
		// Upstream logs "signature-policy-identifier attribute shall be present for {}-EPES
		// signature (cardinality == 1)!".
		return false
	}
	signaturePolicyStore := sig.SignaturePolicyStore()
	if signaturePolicyStore != nil && !b.IsSignaturePolicyIdentifierHashPresent() {
		// Upstream logs "signature-policy-store may be present for {}-EPES signature only if
		// signature-policy-identifier is present and it contains sigPolicyHash element
		// (requirement (c))!".
		return false
	}
	return true
}

// HasExtendedTProfile checks if the signature has a corresponding *AdES-T profile.
// Port of hasExtendedTProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedTProfile() bool {
	if !b.MinimalTRequirement() {
		return false
	}
	// Additional requirement (f)
	if !b.SignatureTimestampsCreatedBeforeSignCertExpiration() {
		// Upstream logs "signature-time-stamp shall be created before expiration of the
		// signing-certificate for CAdES-T signature (requirement (f))!".
		return false
	}
	return true
}

// HasExtendedCProfile checks if the signature has a corresponding *AdES-C profile.
// Port of hasExtendedCProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedCProfile() bool {
	signerInformation := b.Signature().SignerInformation()

	// NOTE: at least complete-certificate-references shall be present for all self-signed
	// certificates
	// complete-certificate-references
	certificateRefAttrs := CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsCertificateRefs)
	completeCertificateRefsNumberOfOccurrences := cadesBaselineAttributeValuesSize(certificateRefAttrs)
	if completeCertificateRefsNumberOfOccurrences > 1 || completeCertificateRefsNumberOfOccurrences == 0 {
		// Upstream logs "complete-certificate-references attribute shall be present for
		// CAdES-C signature (cardinality == 1)!".
		return false
	}

	certificateSources := b.CertificateSourcesExceptLastArchiveTimestamp()
	certificateFound := certificateSources.NumberOfCertificates() > 0
	allSelfSigned := certificateFound && certificateSources.IsAllSelfSigned()

	// complete-revocation-references
	revocationRefAttrs := CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaEtsRevocationRefs)
	completeRevocationRefsNumberOfOccurrences := cadesBaselineAttributeValuesSize(revocationRefAttrs)
	if completeRevocationRefsNumberOfOccurrences > 1 || (!allSelfSigned && completeRevocationRefsNumberOfOccurrences == 0) {
		// Upstream logs "complete-revocation-references attribute shall be present for
		// CAdES-C signature (cardinality == 1)!".
		return false
	}
	return true
}

// HasExtendedXProfile checks if the signature has a corresponding *AdES-X profile.
// Port of hasExtendedXProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedXProfile() bool {
	signerInformation := b.Signature().SignerInformation()
	if utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, OIDIdAaEtsCertCRLTimestamp))+
		utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, OIDIdAaEtsEscTimeStamp)) != 1 {
		// Upstream logs "CAdES-C-timestamp or time-stamped-certs-crls-references attribute
		// shall be present for CAdES-X signature (cardinality == 1)!".
		return false
	}
	return true
}

// HasExtendedXLProfile checks if the signature has a corresponding *AdES-XL profile.
// Port of hasExtendedXLProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedXLProfile() bool {
	return b.MinimalLTRequirement()
}

// HasExtendedAProfile checks if the signature has a corresponding *AdES-A profile.
// Port of hasExtendedAProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedAProfile() bool {
	sig := b.Signature()
	var timestampTokens []*validation.TimestampToken
	timestampTokens = append(timestampTokens, sig.ArchiveTimestamps()...)
	timestampTokens = append(timestampTokens, sig.DetachedTimestamps()...)
	if utils.IsCollectionEmpty(timestampTokens) {
		// Upstream logs "ArchiveTimeStamp shall be present for CAdES-A signature
		// (cardinality >= 1)!".
		return false
	}
	return true
}

// HasExtendedERSProfile checks if the signature has a corresponding *AdES-E-ERS profile.
// Port of hasExtendedERSProfile().
func (b *CAdESBaselineRequirementsChecker) HasExtendedERSProfile() bool {
	// Validate for every signer, as in CMS an embedded ER covers all signatures
	sig := b.Signature()
	signerERSFound := false
	for _, signerInformation := range sig.CMS().SignerInfos() {
		// internal-evidence-record
		internalERNumber := utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaErInternal))
		// external-evidence-record
		externalERNumber := utils.ArraySize(CAdESUtilsUnsignedAttributesOfType(signerInformation, spi.OIDIdAaErExternal))
		if internalERNumber+externalERNumber == 0 {
			// Upstream logs "internal-evidence-records or external-evidence-records
			// attribute shall be present for CAdES-ERS signature (cardinality >= 1)!".
			continue
		}

		if sig.CMS().IsDetachedSignature() {
			if internalERNumber > 0 {
				// Upstream logs "In case a signature is detached, the
				// external-evidence-records attribute shall be used (requirement (q))!".
				continue
			}
		} else {
			if externalERNumber > 0 {
				// Upstream logs "In case a signature is attached, the
				// internal-evidence-records attribute shall be used (requirement (q))!".
				continue
			}
		}
		signerERSFound = true
	}
	return signerERSFound
}

// isContentTypeValid verifies whether the presence of content-type attribute is conformant to
// the given signature type. Port of the private isContentTypeValid(SignerInformation).
func (b *CAdESBaselineRequirementsChecker) isContentTypeValid(signerInformation *cmscore.SignerInfo) bool {
	contentTypeAttrs := CAdESUtilsSignedAttributesOfType(signerInformation, OIDPkcs9AtContentType)
	numberOfOccurrences := cadesBaselineAttributeValuesSize(contentTypeAttrs)
	if b.Signature().IsCounterSignature() && numberOfOccurrences == 0 {
		return true
	}
	return numberOfOccurrences == 1
}

// cadesBaselineIsMessageDigestPresent verifies the presence of a message-digest attribute.
// Port of the private isMessageDigestPresent(SignerInformation).
func cadesBaselineIsMessageDigestPresent(signerInformation *cmscore.SignerInfo) bool {
	messageDigestAttrs := CAdESUtilsSignedAttributesOfType(signerInformation, OIDPkcs9AtMessageDigest)
	return cadesBaselineAttributeValuesSize(messageDigestAttrs) == 1
}

// cadesBaselineIsOneSigningCertificatePresent ports the private
// isOneSigningCertificatePresent(SignerInformation).
func cadesBaselineIsOneSigningCertificatePresent(signerInformation *cmscore.SignerInfo) bool {
	signingCertAttrs := CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaSigningCertificate)
	signingCertV2Attrs := CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaSigningCertificateV2)
	return cadesBaselineAttributeValuesSize(signingCertAttrs)+cadesBaselineAttributeValuesSize(signingCertV2Attrs) == 1
}

// isSigningCertificateAttributeValid ports the private
// isSigningCertificateAttributeValid(SignerInformation).
func (b *CAdESBaselineRequirementsChecker) isSigningCertificateAttributeValid(signerInformation *cmscore.SignerInfo) bool {
	certificateRefs := b.Signature().CertificateSource().SigningCertificateRefs()
	if utils.IsCollectionNotEmpty(certificateRefs) {
		signingCertificateRef := certificateRefs[0] // only one shall be used
		certDigest := signingCertificateRef.CertDigest()
		if certDigest.Algorithm() != "" {
			digestAlgorithm := certDigest.Algorithm()
			if enumerations.DigestAlgorithmSHA1 == digestAlgorithm {
				if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaSigningCertificate)) == 0 {
					return false
				}
			} else {
				if utils.ArraySize(CAdESUtilsSignedAttributesOfType(signerInformation, spi.OIDIdAaSigningCertificateV2)) == 0 {
					return false
				}
			}
		}
	}
	return true
}

// cadesBaselineAttributeValuesSize ports the private getAttributeValuesSize(Attribute...).
func cadesBaselineAttributeValuesSize(attributes []*cmscore.Attribute) int {
	counter := 0
	for _, attribute := range attributes {
		attrValue := spi.DSSASN1UtilsAsn1Encodable(attribute)
		if attrValue != nil {
			counter++
		}
	}
	return counter
}
