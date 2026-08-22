// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// Performs checks according to EN 319 142-1 v1.1.1 "6.3 PAdES baseline signatures".
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug call site is called out in
// the surrounding comment instead.
//
// STRUCTURE DEVIATION (documented for the integrator): Java's PAdESBaselineRequirementsChecker
// extends CAdESBaselineRequirementsChecker, itself hard-typed as
// BaselineRequirementsChecker<CAdESSignature> (not re-parametrised for PAdESSignature - Java
// relies on PAdESSignature being a CAdESSignature subtype and on virtual dispatch to route every
// signature.getXxx() call the shared/CAdES code makes back to PAdESSignature's own overrides).
// Go generics have no covariant substitution for a type parameter and no virtual dispatch
// through an embedded base, so this port cannot embed cades.CAdESBaselineRequirementsChecker
// (pinned to *cades.CAdESSignature) and still have MinimalLTRequirement/MinimalTRequirement/
// ContainsSigningCertificate/ValidationContext/CertificateSourcesExceptLastArchiveTimestamp
// resolve against PAdESSignature's own (DSS-dictionary-aware) certificate/CRL/OCSP sources - the
// behaviour every method PAdESBaselineRequirementsChecker.java itself defines actually needs.
// Instead this type embeds its own validation.BaselineRequirementsChecker[*PAdESSignature]
// (a sibling of cades.CAdESBaselineRequirementsChecker, not a specialisation of it), giving
// correct polymorphism for those methods. The six profile methods
// (HasExtendedTProfile/C/X/XL/A/ERS) that PAdESBaselineRequirementsChecker.java does not
// override at all (pure Java inheritance of CAdESBaselineRequirementsChecker's CMS/
// SignerInformation-only logic) are reproduced by delegating to a CAdES-level checker built over
// this signature's embedded CAdES CMS handling (see cadesChecker() below); this is faithful for
// every CMS-attribute-only accessor those six methods use, with one narrow, documented exception
// (HasExtendedXLProfile's internal MinimalLTRequirement call - see cadesChecker()'s comment).
//
// FORWARD DEPENDENCIES (not in this chunk's manifest, owned by sibling chunks):
//
//   - PAdESSignature (eu.europa.esig.dss.pades.validation.PAdESSignature). Per this phase's
//     layout note ("PAdESSignature embeds CAdES-like CMS handling"), it is assumed to embed
//     *cades.CAdESSignature (a promoted field literally named CAdESSignature), and to expose:
//
//     func (s *PAdESSignature) PdfSignatureDictionary() *PdfSignatureDictionary
//     func (s *PAdESSignature) DssDictionary() PdfDssDict
//     func (s *PAdESSignature) MessageDigestValue() []byte
//     func (s *PAdESSignature) CMS() *cms.CMS                     // promoted from *cades.CAdESSignature
//     func (s *PAdESSignature) SignerInformation() *cmscore.SignerInfo // promoted
//     func (s *PAdESSignature) ContentType() string                // promoted (AdvancedSignature)
//
//     with CertificateSource()/CompleteCertificateSource()/CompleteCRLSource()/CompleteOCSPSource()
//     etc. shadowed on *PAdESSignature itself (not merely promoted) to satisfy AdvancedSignature
//     with PAdES-aware (DSS-dictionary-including) sources - see the STRUCTURE DEVIATION note.
//
//   - PdfSignatureDictionary (eu.europa.esig.dss.pdf.PdfSignatureDictionary), already used as a
//     forward dependency by the landed SIGN chunk (native_pdf_signature_service.go), confirming
//     Type()/SubFilter()/ByteRange()/Contents(); this file additionally uses
//     SigningDate() time.Time (the zero value stands for Java's null), Filter() string and
//     Reason() string.
//
//   - PdfTimestampToken (eu.europa.esig.dss.pades.validation.timestamp.PdfTimestampToken) and
//     PdfDocTimestampRevision (eu.europa.esig.dss.pdf.PdfDocTimestampRevision). Java
//     distinguishes a PDF document-timestamp's TimestampToken from any other with
//     "instanceof PdfTimestampToken"; validation.TimestampToken is a concrete struct (not an
//     interface), so a *validation.TimestampToken value carries no Go-visible tag for which
//     wrapper (if any) built it. PdfTimestampTokenOf below is assumed to be the sibling chunk's
//     answer to that gap - a lookup that succeeds exactly when timestampToken was produced by
//     the PAdES timestamp source:
//
//     func PdfTimestampTokenOf(timestampToken *validation.TimestampToken) (*PdfTimestampToken, bool)
//     func (t *PdfTimestampToken) PdfRevision() *PdfDocTimestampRevision
//     func (r *PdfDocTimestampRevision) PdfSigDictInfo() *PdfSignatureDictionary
//
//   - CMSForPAdESBaselineRequirementsChecker
//     (eu.europa.esig.dss.pades.signature.CMSForPAdESBaselineRequirementsChecker, not in this
//     chunk's manifest; landed as cms_for_pades_baseline_requirements_checker.go, confirming the
//     shape below - its constructor takes *cades.CAdESSignature, not *PAdESSignature, hence
//     b.Signature().CAdESSignature at the one call site below):
//
//     func NewCMSForPAdESBaselineRequirementsChecker(signature *cades.CAdESSignature) *CMSForPAdESBaselineRequirementsChecker
//     func (c *CMSForPAdESBaselineRequirementsChecker) IsValidForPAdESBaselineBProfile() bool
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PAdESBaselineRequirementsChecker checks conformance of a PAdES signature to the requested
// baseline format.
type PAdESBaselineRequirementsChecker struct {
	validation.BaselineRequirementsChecker[*PAdESSignature]

	// offlineCertificateVerifier mirrors the base's own copy, kept here as well so cadesChecker
	// can build a fresh cades-level checker with the same verifier.
	offlineCertificateVerifier validation.CertificateVerifier
}

// NewPAdESBaselineRequirementsChecker is the default constructor.
// Port of the constructor PAdESBaselineRequirementsChecker(PAdESSignature, CertificateVerifier).
func NewPAdESBaselineRequirementsChecker(signature *PAdESSignature, offlineCertificateVerifier validation.CertificateVerifier) *PAdESBaselineRequirementsChecker {
	checker := &PAdESBaselineRequirementsChecker{
		BaselineRequirementsChecker: validation.NewBaselineRequirementsCheckerBaseWithVerifier[*PAdESSignature](signature, offlineCertificateVerifier),
		offlineCertificateVerifier:  offlineCertificateVerifier,
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// cadesChecker builds a CAdES-level requirements checker over this signature's embedded CAdES
// CMS handling; see this file's header ("STRUCTURE DEVIATION") for why and its limits.
func (b *PAdESBaselineRequirementsChecker) cadesChecker() *cades.CAdESBaselineRequirementsChecker {
	return cades.NewCAdESBaselineRequirementsChecker(b.Signature().CAdESSignature, b.offlineCertificateVerifier)
}

// GetBaselineSignatureForm returns the signature form corresponding to the signature: PAdES.
// Port of the protected getBaselineSignatureForm(), exported to satisfy this port's
// validation.BaselineRequirementsCheckerOverrides (see that interface's doc comment on
// GetBaselineSignatureForm). b's own HasBaselineBProfile() override never actually reaches
// cades.CAdESBaselineRequirementsChecker.cmsBaselineBRequirements() through this method -
// cmsBaselineBRequirements() below builds a fresh CMSForPAdESBaselineRequirementsChecker, whose
// own GetBaselineSignatureForm() override is what that call chain actually resolves - but this
// override still has to exist so b satisfies BaselineRequirementsCheckerOverrides for
// InitBaselineRequirementsChecker(b) in NewPAdESBaselineRequirementsChecker, matching Java's own
// PAdESBaselineRequirementsChecker override of the identical method for the identical reason.
func (b *PAdESBaselineRequirementsChecker) GetBaselineSignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureForm_PAdES
}

// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile.
// Port of hasAdESProfile().
func (b *PAdESBaselineRequirementsChecker) HasAdESProfile() bool {
	return b.HasExtendedBESProfile() || b.HasBaselineBProfile()
}

// HasBaselineBProfile checks if the signature has a corresponding BASELINE-B profile.
// Port of hasBaselineBProfile().
func (b *PAdESBaselineRequirementsChecker) HasBaselineBProfile() bool {
	if !b.cmsBaselineBRequirements() {
		return false
	}
	padesSignature := b.Signature()
	pdfSignatureDictionary := padesSignature.PdfSignatureDictionary()
	// SPO: entry with the key M in the Signature Dictionary (Cardinality == 1)
	if pdfSignatureDictionary.SigningDate().IsZero() {
		// Upstream logs "Entry with the key M in the Signature Dictionary shall be present for
		// PAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Contents in the Signature Dictionary (Cardinality == 1)
	if utils.IsArrayEmpty(pdfSignatureDictionary.Contents()) {
		// Upstream logs "Entry with the key Contents in the Signature Dictionary shall be
		// present for PAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Filter in the Signature Dictionary (Cardinality == 1)
	if utils.IsStringEmpty(pdfSignatureDictionary.Filter()) {
		// Upstream logs "Entry with the key Filter in the Signature Dictionary shall be present
		// for PAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key ByteRange in the Signature Dictionary (Cardinality == 1)
	if pdfSignatureDictionary.ByteRange() == nil {
		// Upstream logs "Entry with the key ByteRange in the Signature Dictionary shall be
		// present for PAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key SubFilter in the Signature Dictionary (Cardinality == 1)
	if utils.IsStringEmpty(pdfSignatureDictionary.SubFilter()) {
		// Upstream logs "Entry with the key SubFilter in the Signature Dictionary shall be
		// present for PAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Cert in the Signature Dictionary (Cardinality == 0) // TODO : (not supported)
	// Additional requirement (c)
	if cmscore.OIDData.String() != padesSignature.ContentType() {
		// Upstream logs "content-type attribute shall have value 'id-data' for PAdES-BASELINE-B
		// signature! (requirement (c))".
		return false
	}
	// Additional requirement (d)
	if utils.IsStringNotEmpty(pdfSignatureDictionary.Reason()) &&
		utils.IsCollectionNotEmpty(padesSignature.CommitmentTypeIndications()) {
		// Upstream logs "commitment-type-indication attribute shall not be incorporated in the
		// CMS signature when entry with a key Reason is used for PAdES-BASELINE-B signature!
		// (requirement (d))".
		return false
	}
	// Additional requirement (l)
	if PAdESConstantsSignatureDefaultSubFilter != pdfSignatureDictionary.SubFilter() {
		// Upstream logs "Entry with a key SubFilter shall contain a value 'ETSI.CAdES.detached'
		// for PAdES-BASELINE-B signature! (requirement (l))".
		return false
	}
	// Additional requirement (m)
	if (utils.IsCollectionNotEmpty(padesSignature.CommitmentTypeIndications()) ||
		padesSignature.SignaturePolicy() != nil) && utils.IsStringNotEmpty(pdfSignatureDictionary.Reason()) {
		// Upstream logs "Entry with a key Reason shall not be used when
		// commitment-type-attribute or signature-policy-identifier is present in the CMS
		// signature for PAdES-BASELINE-B signature! (requirement (m))".
		return false
	}
	return true
}

// cmsBaselineBRequirements checks if BASELINE-B requirements are satisfied for a CMS signature.
// Port of the protected cmsBaselineBRequirements() override.
func (b *PAdESBaselineRequirementsChecker) cmsBaselineBRequirements() bool {
	cmsRequirementsChecker := NewCMSForPAdESBaselineRequirementsChecker(b.Signature().CAdESSignature)
	return cmsRequirementsChecker.IsValidForPAdESBaselineBProfile()
}

// HasBaselineTProfile checks if the signature has a corresponding BASELINE-T profile.
// Port of hasBaselineTProfile().
func (b *PAdESBaselineRequirementsChecker) HasBaselineTProfile() bool {
	// signature-time-stamp or document-time-stamp (Cardinality >= 1)
	if utils.IsCollectionEmpty(b.Signature().SignatureTimestamps()) &&
		utils.IsCollectionEmpty(b.Signature().DocumentTimestamps()) {
		// Upstream logs "SignatureTimeStamp shall be present for BASELINE-T signature
		// (cardinality >= 1)!".
		return false
	}
	return true
}

// HasBaselineLTProfile checks if the signature has a corresponding BASELINE-LT profile.
// Port of hasBaselineLTProfile().
func (b *PAdESBaselineRequirementsChecker) HasBaselineLTProfile() bool {
	return b.hasLTProfile()
}

// hasLTProfile verifies a presence of LT-profile for a PDF signature.
// Port of the protected hasLTProfile().
func (b *PAdESBaselineRequirementsChecker) hasLTProfile() bool {
	if !b.MinimalLTRequirement() {
		return false
	}
	padesSignature := b.Signature()
	allSelfSigned := b.CertificateSourcesExceptLastArchiveTimestamp().IsAllSelfSigned()
	// SPO: DSS
	if !allSelfSigned && padesSignature.DssDictionary() == nil {
		// Upstream logs "DSS dictionary shall be present for PAdES-BASELINE-LT signature!
		// (cardinality >= 1)".
		return false
	}
	return true
}

// HasBaselineLTAProfile checks if the signature has a corresponding BASELINE-LTA profile.
// Port of hasBaselineLTAProfile().
func (b *PAdESBaselineRequirementsChecker) HasBaselineLTAProfile() bool {
	// Additional requirement (y)
	if !b.isBaselineLTATimestampPresent() {
		// Upstream logs "document-time-stamp covering LT-level and containing a key SubFilter
		// with value ETSI.RFC3161 shall be present for PAdES-BASELINE-LTA signature
		// (cardinality >= 1, requirement (y))!".
		return false
	}
	return true
}

// isBaselineLTATimestampPresent ports the private isBaselineLTATimestampPresent().
func (b *PAdESBaselineRequirementsChecker) isBaselineLTATimestampPresent() bool {
	for _, timestampToken := range b.Signature().DocumentTimestamps() {
		if b.isBaselineLTATimestamp(timestampToken) {
			return true
		}
	}
	return false
}

// isBaselineLTATimestamp ports the private isBaselineLTATimestamp(TimestampToken).
func (b *PAdESBaselineRequirementsChecker) isBaselineLTATimestamp(timestampToken *validation.TimestampToken) bool {
	return b.containsRFC3161SubFilter(timestampToken) && b.coversLTLevelData(timestampToken)
}

// HasExtendedBESProfile checks if the signature has a corresponding *AdES-BES profile.
// Port of hasExtendedBESProfile().
func (b *PAdESBaselineRequirementsChecker) HasExtendedBESProfile() bool {
	if !b.cmsExtendedBESRequirements() {
		return false
	}
	padesSignature := b.Signature()
	cmsDoc := padesSignature.CMS()
	pdfSignatureDictionary := padesSignature.PdfSignatureDictionary()
	// PAdES Part 1 : a DER-encoded SignedData object as specified in ETSI EN 319 122-1 [2]
	// shall be included
	if cmsDoc == nil {
		// Upstream logs "DER-encoded SignedData object shall be included as the PDF signature
		// in the entry with the key Contents of the Signature Dictionary for PAdES-BES
		// signature (PAdES Part 1, requirement (a))!".
		return false
	}
	// PAdES Part 1 : There shall only be a single signer in any PDF Signature
	if len(cmsDoc.SignerInfos()) != 1 {
		// Upstream logs "There shall be a single signer for any PAdES-BES signature (PAdES
		// Part 1, requirement (a))!".
		return false
	}
	// PAdES Part 1 : "No data shall be encapsulated in the PKCS#7 SignedData field
	if !padesSignature.CMS().IsDetachedSignature() {
		// Upstream logs "No data shall be encapsulated in the PKCS#7 SignedData field for
		// PAdES-BES signature (PAdES Part 1, requirement (b))!".
		return false
	}
	// SPO: entry with the key Filter in the Signature Dictionary (Cardinality == 1)
	if utils.IsStringEmpty(pdfSignatureDictionary.Filter()) {
		// Upstream logs "Entry with the key Filter in the Signature Dictionary shall be present
		// for PAdES-BES signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key ByteRange in the Signature Dictionary (Cardinality == 1)
	if pdfSignatureDictionary.ByteRange() == nil {
		// Upstream logs "Entry with the key ByteRange in the Signature Dictionary shall be
		// present for PAdES-BES signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key SubFilter in the Signature Dictionary (Cardinality == 1)
	if utils.IsStringEmpty(pdfSignatureDictionary.SubFilter()) {
		// Upstream logs "Entry with the key SubFilter in the Signature Dictionary shall be
		// present for PAdES-BES signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Contents in the Signature Dictionary (Cardinality == 1)
	if utils.IsArrayEmpty(pdfSignatureDictionary.Contents()) {
		// Upstream logs "Entry with the key Contents in the Signature Dictionary shall be
		// present for PAdES-BES signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Cert in the Signature Dictionary (Cardinality == 0) // TODO : (not supported)
	// Additional requirement (a)
	if cmscore.OIDData.String() != padesSignature.ContentType() {
		// Upstream logs "content-type attribute shall have value 'id-data' for PAdES-BES
		// signature! (requirement (a))".
		return false
	}
	// Additional requirement (e)
	if utils.IsStringNotEmpty(pdfSignatureDictionary.Reason()) &&
		utils.IsCollectionNotEmpty(padesSignature.CommitmentTypeIndications()) {
		// Upstream logs "commitment-type-indication attribute shall not be incorporated in the
		// CMS signature when entry with a key Reason is used for PAdES-BES signature!
		// (requirement (e))".
		return false
	}
	// Additional requirement (j)
	if PAdESConstantsSignatureDefaultSubFilter != pdfSignatureDictionary.SubFilter() {
		// Upstream logs "Entry with a key SubFilter shall contain a value 'ETSI.CAdES.detached'
		// for PAdES-BES signature! (requirement (j))".
		return false
	}
	return true
}

// cmsExtendedBESRequirements verifies whether CMS is conformant to the E-BES profile.
//
// Port of the protected cmsExtendedBESRequirements(); not overridden by upstream
// PAdESBaselineRequirementsChecker, so inherited unchanged from
// CAdESBaselineRequirementsChecker. Unlike the six HasExtendedXxxProfile methods above, that
// inherited method (and the private helpers it calls: isContentTypeValid,
// isMessageDigestPresent, isOneSigningCertificatePresent, getAttributeValuesSize) are
// unexported in package cades, so delegating to cadesChecker() is not possible here; this
// reproduces the same three cardinality checks directly instead, against
// cades.CAdESUtilsSignedAttributesOfType/spi.DSSASN1UtilsAsn1Encodable, the same building
// blocks cades_baseline_requirements_checker.go itself is built from.
func (b *PAdESBaselineRequirementsChecker) cmsExtendedBESRequirements() bool {
	signerInformation := b.Signature().SignerInformation()

	// content-type (Cardinality == 1)
	if !b.isContentTypeValid(signerInformation) {
		// Upstream logs "content-type attribute shall be present for {}-BES signature
		// (cardinality == 1)!".
		return false
	}
	// message-digest (Cardinality == 1)
	if padesBaselineAttributeValuesSize(cades.CAdESUtilsSignedAttributesOfType(signerInformation, cades.OID_pkcs_9_at_messageDigest)) != 1 {
		// Upstream logs "message-digest attribute shall be present for {}-BES signature
		// (cardinality == 1)!".
		return false
	}
	// signing-certificate/signing-certificate-v2 (Cardinality == 1)
	signingCertAttrs := cades.CAdESUtilsSignedAttributesOfType(signerInformation, spi.OID_id_aa_signingCertificate)
	signingCertV2Attrs := cades.CAdESUtilsSignedAttributesOfType(signerInformation, spi.OID_id_aa_signingCertificateV2)
	if padesBaselineAttributeValuesSize(signingCertAttrs)+padesBaselineAttributeValuesSize(signingCertV2Attrs) != 1 {
		// Upstream logs "signing-certificate(-v2) attribute shall be present for {}-BES
		// signature (cardinality == 1)!".
		return false
	}
	// signer-attributes (Cardinality == 0 or 1)
	if len(cades.CAdESUtilsSignedAttributesOfType(signerInformation, cades.OID_id_aa_ets_signerAttr))+
		len(cades.CAdESUtilsSignedAttributesOfType(signerInformation, spi.OID_id_aa_ets_signerAttrV2)) > 1 {
		// Upstream logs "signer-attributes(-v2) attribute shall not be present multiple times
		// for {}-BES signature (cardinality == 0 or 1)!".
		return false
	}
	return true
}

// isContentTypeValid verifies whether the presence of content-type attribute is conformant to
// the given signature type. Port of the private isContentTypeValid(SignerInformation); see
// cmsExtendedBESRequirements's doc comment for why it is reproduced here rather than reused.
func (b *PAdESBaselineRequirementsChecker) isContentTypeValid(signerInformation *cmscore.SignerInfo) bool {
	contentTypeAttrs := cades.CAdESUtilsSignedAttributesOfType(signerInformation, cades.OID_pkcs_9_at_contentType)
	numberOfOccurrences := padesBaselineAttributeValuesSize(contentTypeAttrs)
	if b.Signature().IsCounterSignature() && numberOfOccurrences == 0 {
		return true
	}
	return numberOfOccurrences == 1
}

// padesBaselineAttributeValuesSize ports the private getAttributeValuesSize(Attribute...); see
// cmsExtendedBESRequirements's doc comment for why it is reproduced here rather than reused.
func padesBaselineAttributeValuesSize(attributes []*cmscore.Attribute) int {
	counter := 0
	for _, attribute := range attributes {
		if spi.DSSASN1UtilsAsn1Encodable(attribute) != nil {
			counter++
		}
	}
	return counter
}

// HasExtendedEPESProfile checks if the signature has a corresponding *AdES-EPES profile.
// Port of hasExtendedEPESProfile().
func (b *PAdESBaselineRequirementsChecker) HasExtendedEPESProfile() bool {
	padesSignature := b.Signature()
	signerInformation := padesSignature.SignerInformation()
	pdfSignatureDictionary := padesSignature.PdfSignatureDictionary()
	// signature-policy-identifier (Cardinality == 1)
	if len(cades.CAdESUtilsSignedAttributesOfType(signerInformation, cades.OID_id_aa_ets_sigPolicyId)) == 0 {
		// Upstream logs "signature-policy-identifier attribute shall be present for PAdES-EPES
		// signature (cardinality == 1)!".
		return false
	}
	// SPO: entry with the key Reason in the Signature Dictionary (Cardinality == 0)
	if utils.IsStringNotEmpty(pdfSignatureDictionary.Reason()) {
		// Upstream logs "Entry with the key Reason in the Signature Dictionary shall not be
		// present for PAdES-EPES signature (cardinality == 0)!".
		return false
	}
	return true
}

// HasExtendedTProfile checks if the signature has a corresponding *AdES-T profile. Port of
// hasExtendedTProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) HasExtendedTProfile() bool {
	return b.cadesChecker().HasExtendedTProfile()
}

// HasExtendedCProfile checks if the signature has a corresponding *AdES-C profile. Port of
// hasExtendedCProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) HasExtendedCProfile() bool {
	return b.cadesChecker().HasExtendedCProfile()
}

// HasExtendedXProfile checks if the signature has a corresponding *AdES-X profile. Port of
// hasExtendedXProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) HasExtendedXProfile() bool {
	return b.cadesChecker().HasExtendedXProfile()
}

// HasExtendedXLProfile checks if the signature has a corresponding *AdES-XL profile. Port of
// hasExtendedXLProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header for the
// narrow deviation this delegation carries (MinimalLTRequirement resolves against the embedded
// CAdESSignature's own sources rather than PAdESSignature's DSS-dictionary-aware ones here).
func (b *PAdESBaselineRequirementsChecker) HasExtendedXLProfile() bool {
	return b.cadesChecker().HasExtendedXLProfile()
}

// HasExtendedAProfile checks if the signature has a corresponding *AdES-A profile. Port of
// hasExtendedAProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) HasExtendedAProfile() bool {
	return b.cadesChecker().HasExtendedAProfile()
}

// HasExtendedERSProfile checks if the signature has a corresponding *AdES-E-ERS profile. Port of
// hasExtendedERSProfile(); not overridden by upstream PAdESBaselineRequirementsChecker, so
// inherited unchanged from CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) HasExtendedERSProfile() bool {
	return b.cadesChecker().HasExtendedERSProfile()
}

// ContainsLTLevelCertificates verifies whether the signature contains some of the LT-/XL-level
// attributes. Port of the protected containsLTLevelCertificates() override; not defined by
// upstream PAdESBaselineRequirementsChecker, so inherited unchanged from
// CAdESBaselineRequirementsChecker - see this file's header.
func (b *PAdESBaselineRequirementsChecker) ContainsLTLevelCertificates() bool {
	return b.cadesChecker().ContainsLTLevelCertificates()
}

// HasExtendedLTVProfile checks if the signature has a corresponding PAdES-E-LTV profile.
// Port of hasExtendedLTVProfile().
func (b *PAdESBaselineRequirementsChecker) HasExtendedLTVProfile() bool {
	// a) Validation data check
	if !b.MinimalLTRequirement() {
		return false
	}
	padesSignature := b.Signature()
	allSelfSigned := b.CertificateSourcesExceptLastArchiveTimestamp().IsAllSelfSigned()
	// Validation data shall be carried by values within the DSS.
	if !allSelfSigned && padesSignature.DssDictionary() == nil {
		// Upstream logs "DSS dictionary shall be present for PAdES-LTV signature!".
		return false
	}
	if !b.isLTVTimestampPresent() {
		// Upstream logs "document-time-stamp covering validation data shall be present for
		// PAdES-LTV signature!".
		return false
	}
	return true
}

// isLTVTimestampPresent ports the private isLTVTimestampPresent().
func (b *PAdESBaselineRequirementsChecker) isLTVTimestampPresent() bool {
	for _, timestampToken := range b.Signature().DocumentTimestamps() {
		if b.coversLTLevelData(timestampToken) {
			return true
		}
	}
	return false
}

// coversLTLevelData ports the private coversLTLevelData(TimestampToken).
func (b *PAdESBaselineRequirementsChecker) coversLTLevelData(timestampToken *validation.TimestampToken) bool {
	if enumerations.ArchiveTimestampType_PAdES == timestampToken.ArchiveTimestampType() {
		signatureValidationData := b.ValidationContext().GetValidationData(b.Signature())
		certificateTokens := signatureValidationData.CertificateTokens()
		crlTokens := signatureValidationData.CrlTokens()
		ocspTokens := signatureValidationData.OcspTokens()
		allTimestamps := b.Signature().AllTimestamps()

		if utils.IsCollectionEmpty(crlTokens) && utils.IsCollectionEmpty(ocspTokens) {
			return b.coversDSSCertificateTokens(timestampToken, certificateTokens)
		}
		return b.coversRevocationTokens(timestampToken, crlTokens, ocspTokens) &&
			(b.coversTimestampTokens(timestampToken, allTimestamps) || b.coversOwnRevocationData(timestampToken))
	}
	return false
}

// coversDSSCertificateTokens ports the private coversDSSCertificateTokens(TimestampToken,
// Collection<CertificateToken>).
func (b *PAdESBaselineRequirementsChecker) coversDSSCertificateTokens(timestampToken *validation.TimestampToken,
	certificateTokens []*model.CertificateToken) bool {
	dssCertificates := make([]*model.CertificateToken, 0)
	dssCertificates = append(dssCertificates, b.Signature().CertificateSource().DSSDictionaryCertValues()...)
	dssCertificates = append(dssCertificates, b.Signature().CertificateSource().VRIDictionaryCertValues()...)
	if utils.IsCollectionNotEmpty(dssCertificates) {
		for _, certificateToken := range certificateTokens {
			if padesBaselineContainsCertificate(dssCertificates, certificateToken) &&
				b.coversToken(timestampToken, certificateToken) {
				return true
			}
		}
	}
	return false
}

// padesBaselineContainsCertificate reports whether certificates contains candidate, reproducing
// java.util.List#contains (equals()-based) semantics.
func padesBaselineContainsCertificate(certificates []*model.CertificateToken, candidate *model.CertificateToken) bool {
	for _, c := range certificates {
		if c.Equals(candidate) {
			return true
		}
	}
	return false
}

// coversRevocationTokens verifies whether all the revocation data is covered by the given
// timestamp to fulfil the minimum requirement for LT-level.
//
// NOTE: This method checks coverage of the available revocation data, and the actual LT-level
// shall be determined in prior using HasBaselineLTProfile()!
// Port of the private coversRevocationTokens(TimestampToken, Collection<CRLToken>,
// Collection<OCSPToken>).
func (b *PAdESBaselineRequirementsChecker) coversRevocationTokens(timestampToken *validation.TimestampToken,
	crlTokens []*spi.CRLToken, ocspTokens []*spi.OCSPToken) bool {
	revocationsByCertificate := padesBaselineRevocationsByCertificate(crlTokens, ocspTokens)
	for _, revocationTokens := range revocationsByCertificate {
		revocationForCertificateIsCovered := false
		for _, revocationToken := range revocationTokens {
			if b.coversToken(timestampToken, revocationToken) {
				revocationForCertificateIsCovered = true
				break
			}
		}
		if !revocationForCertificateIsCovered {
			return false
		}
	}
	return true
}

// coversTimestampTokens ports the private coversTimestampTokens(TimestampToken,
// Collection<TimestampToken>).
func (b *PAdESBaselineRequirementsChecker) coversTimestampTokens(timestampToken *validation.TimestampToken,
	signatureTimestampTokens []*validation.TimestampToken) bool {
	timestampedReferences := timestampToken.TimestampedReferences()
	if utils.IsCollectionNotEmpty(timestampedReferences) {
		for _, sigTst := range signatureTimestampTokens {
			for _, r := range timestampedReferences {
				if sigTst.DSSIDAsString() == r.ObjectId() {
					return true
				}
			}
		}
	}
	return false
}

// coversOwnRevocationData ports the private coversOwnRevocationData(TimestampToken).
func (b *PAdESBaselineRequirementsChecker) coversOwnRevocationData(timestampToken *validation.TimestampToken) bool {
	sig := b.Signature()
	validationContext := validation.NewSignatureValidationContext()
	validationContext.Initialize(b.offlineCertificateVerifier)

	validationContext.AddDocumentCertificateSource(sig.CompleteCertificateSource())
	validationContext.AddDocumentCRLSourceFromList(sig.CompleteCRLSource())
	validationContext.AddDocumentOCSPSourceFromList(sig.CompleteOCSPSource())

	validationContext.AddTimestampTokenForVerification(timestampToken)
	validationContext.Validate()

	validationData := validationContext.GetValidationDataForTimestamp(timestampToken)
	revocationTokenIdentifiers := make(map[string]struct{})
	for _, crlToken := range validationData.CrlTokens() {
		revocationTokenIdentifiers[crlToken.DSSIDAsString()] = struct{}{}
	}
	for _, ocspToken := range validationData.OcspTokens() {
		revocationTokenIdentifiers[ocspToken.DSSIDAsString()] = struct{}{}
	}

	if len(revocationTokenIdentifiers) == 0 {
		return validationContext.CheckAllRequiredRevocationDataPresent()
	}

	timestampedReferences := timestampToken.TimestampedReferences()
	for _, r := range timestampedReferences {
		if _, ok := revocationTokenIdentifiers[r.ObjectId()]; ok {
			return true
		}
	}
	return false
}

// padesBaselineRevocationsByCertificate ports the private getRevocationsByCertificate(
// Collection<CRLToken>, Collection<OCSPToken>) together with the private
// enrichRevocationDataMap(Map, Collection) helper it calls twice.
func padesBaselineRevocationsByCertificate(crlTokens []*spi.CRLToken, ocspTokens []*spi.OCSPToken) map[string][]model.Token {
	result := make(map[string][]model.Token)
	for _, token := range crlTokens {
		padesBaselineEnrichRevocationDataMap(result, token.RelatedCertificateID(), token)
	}
	for _, token := range ocspTokens {
		padesBaselineEnrichRevocationDataMap(result, token.RelatedCertificateID(), token)
	}
	return result
}

// padesBaselineEnrichRevocationDataMap adds token under relatedCertificateID, deduplicating by
// pointer identity (the Set<RevocationToken<?>> semantics of the Java value).
func padesBaselineEnrichRevocationDataMap(revocationDataMap map[string][]model.Token, relatedCertificateID string, token model.Token) {
	for _, existing := range revocationDataMap[relatedCertificateID] {
		if existing == token {
			return
		}
	}
	revocationDataMap[relatedCertificateID] = append(revocationDataMap[relatedCertificateID], token)
}

// coversToken ports the private coversToken(TimestampToken, Token).
func (b *PAdESBaselineRequirementsChecker) coversToken(timestampToken *validation.TimestampToken, token model.Token) bool {
	for _, timestampedReference := range timestampToken.TimestampedReferences() {
		if token.DSSIDAsString() == timestampedReference.ObjectId() {
			return true
		}
	}
	return false
}

// containsRFC3161SubFilter ports the private containsRFC3161SubFilter(TimestampToken).
func (b *PAdESBaselineRequirementsChecker) containsRFC3161SubFilter(timestampToken *validation.TimestampToken) bool {
	pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken)
	if !ok {
		return false
	}
	pdfRevision := pdfTimestampToken.PdfRevision()
	if pdfRevision != nil {
		pdfSigDictInfo := pdfRevision.PdfSigDictInfo()
		return pdfSigDictInfo != nil && PAdESConstantsTimestampDefaultSubFilter == pdfSigDictInfo.SubFilter()
	}
	return false
}

// HasPKCS7Profile checks if the signature has PKCS#7 profile (according to ISO 32000-1).
// Port of hasPKCS7Profile().
func (b *PAdESBaselineRequirementsChecker) HasPKCS7Profile() bool {
	padesSignature := b.Signature()
	pdfSignatureDictionary := padesSignature.PdfSignatureDictionary()
	// SubFilter shall take one of the following values: (adbe.pkcs7.detached, adbe.pkcs7.sha1)
	if PAdESConstantsSignaturePKCS7SubFilter != pdfSignatureDictionary.SubFilter() &&
		PAdESConstantsSignaturePKCS7SHA1SubFilter != pdfSignatureDictionary.SubFilter() {
		// Upstream logs "Entry with a key SubFilter shall have a value adbe.pkcs7.detached or
		// adbe.pkcs7.sha1 for PKCS#7 signature!".
		return false
	}
	// At minimum the CMS object shall include the signer's X.509 signing certificate.
	if !b.ContainsSigningCertificate(padesSignature.CertificateSource().Certificates()) {
		// Upstream logs "PKCS#7 signature shall include signing certificate!".
		return false
	}
	// SubFilter adbe.pkcs7.detached
	if PAdESConstantsSignaturePKCS7SubFilter == pdfSignatureDictionary.SubFilter() {
		// The original signed message digest over the document's byte range shall be
		// incorporated as the normal CMS SignedData field.
		if utils.IsArrayEmpty(padesSignature.MessageDigestValue()) {
			// Upstream logs "PKCS#7 signature shall include message digest!".
			return false
		}
		// No data shall be encapsulated in the CMS SignedData field.
		if !padesSignature.CMS().IsDetachedSignature() {
			// Upstream logs "No data shall be encapsulated in the CMS SignedData field for
			// PKCS#7 signature!".
			return false
		}
	}
	// SubFilter adbe.pkcs7.sha1
	if PAdESConstantsSignaturePKCS7SHA1SubFilter == pdfSignatureDictionary.SubFilter() {
		signedContent := padesSignature.CMS().SignedContent()
		if signedContent == nil {
			// Upstream logs "ContentInfo of type Data shall be encapsulated in the CMS
			// SignedData field for PKCS#7 signature with SHA-1 SubFilter!".
			return false
		}
		signedContentBytes, err := spi.DSSUtilsToByteArrayOfDocument(signedContent)
		if err != nil || !spi.DSSUtilsIsSHA1Digest(utils.ToHex(signedContentBytes)) {
			// Upstream logs "The SHA-1 digest of the document's byte range shall be
			// encapsulated in the CMS SignedData field with ContentInfo of type Data for
			// PKCS#7 signature with SHA-1 SubFilter!".
			return false
		}
	}
	return true
}

// HasPKCS7TProfile checks if the signature has a PKCS#7-T profile. Port of hasPKCS7TProfile().
func (b *PAdESBaselineRequirementsChecker) HasPKCS7TProfile() bool {
	return b.HasBaselineTProfile()
}

// HasPKCS7LTProfile checks if the signature has a PKCS#7-LT profile. Port of hasPKCS7LTProfile().
func (b *PAdESBaselineRequirementsChecker) HasPKCS7LTProfile() bool {
	return b.hasLTProfile()
}

// HasPKCS7LTAProfile checks if the signature has a PKCS#7-LTA profile.
// Port of hasPKCS7LTAProfile().
func (b *PAdESBaselineRequirementsChecker) HasPKCS7LTAProfile() bool {
	ltaTimestampFound := false
	for _, timestampToken := range b.Signature().DocumentTimestamps() {
		if b.coversLTLevelData(timestampToken) {
			ltaTimestampFound = true
			break
		}
	}
	if !ltaTimestampFound {
		// Upstream logs "document-time-stamp covering LT-level shall be present for
		// PKCS#7-LTA signature!".
		return false
	}
	return true
}
