// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/BaselineRequirementsChecker.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: CertificateVerifier, ValidationContext and SignatureValidationContext
// (Java spi.validation.*) are flattened into this same Go package by a sibling chunk of phase
// 2b and are therefore referenced unqualified, following the precedent set by
// timestamp_source.go's EvidenceRecord forward dependency. Their assumed shapes, inferred from
// the Java calls this file makes, are:
//
//	type CertificateVerifier interface { /* opaque; only passed through */ }
//
//	type ValidationContext interface {
//	    Initialize(certificateVerifier CertificateVerifier)
//	    AddDocumentCertificateSource(certificateSource spi.CertificateSource)
//	    AddDocumentCRLSourceFromList(crlSource *spi.ListRevocationSource[revocation.CRL])
//	    AddDocumentOCSPSourceFromList(ocspSource *spi.ListRevocationSource[revocation.OCSP])
//	    AddCertificateTokenForVerification(certificateToken *model.CertificateToken)
//	    AddTimestampTokenForVerification(timestampToken *TimestampToken)
//	    Validate()
//	    CheckAllRequiredRevocationDataPresent() bool
//	}
//
// Integration note: AdvancedSignature.CompleteCRLSource()/CompleteOCSPSource() return
// *spi.ListRevocationSource[R] (matching Java's getCompleteCRLSource(): ListRevocationSource<CRL>),
// which - like its Java counterpart - does not implement OfflineRevocationSource<R>; it only
// implements MultipleRevocationSource<R>. ValidationContext() below therefore calls the
// ...FromList overload (matching ValidationContext's addDocumentCRLSource(ListRevocationSource)
// Java overload), not AddDocumentCRLSource/AddDocumentOCSPSource.
//
//	func NewSignatureValidationContext() ValidationContext { ... }
//
// Java's two addDocumentCertificateSource/addDocumentCRLSource/addDocumentOCSPSource overloads
// (one taking the plural interface, one taking the List* aggregate) collapse to a single Go
// method each: *spi.ListCertificateSource and *spi.ListRevocationSource[R] already implement
// spi.CertificateSource and spi.OfflineRevocationSource[R] respectively (see their "compile-time
// assertion" comments), the same way ListCertificateSource itself does, so passing the plural
// value where the singular interface is expected type-checks without a second overload.
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// BaselineRequirementsCheckerSignature is what BaselineRequirementsChecker needs from its type
// parameter AS: the full AdvancedSignature contract plus the protected
// getCounterSignaturesCertificateSource() DefaultAdvancedSignature provides. Java expresses this
// as the bound "AS extends DefaultAdvancedSignature"; Go generics need the extra method spelled
// out explicitly since AdvancedSignature itself does not declare it.
type BaselineRequirementsCheckerSignature interface {
	AdvancedSignature

	// CounterSignaturesCertificateSource returns a merged certificate source for values
	// incorporated within counter signatures. Port of the protected
	// DefaultAdvancedSignature#getCounterSignaturesCertificateSource().
	CounterSignaturesCertificateSource() *spi.ListCertificateSource
}

// BaselineRequirementsCheckerContract is what DefaultAdvancedSignature needs from a
// BaselineRequirementsChecker[AS] instance, however AS is instantiated. Go generics cannot
// erase a type parameter into a shared, non-generic field the way Java's raw type
// BaselineRequirementsChecker<?> does; this non-generic interface is that abstraction.
// DefaultAdvancedSignature stores its cached checker with this type, and
// BaselineRequirementsChecker[AS] (via whatever concrete struct embeds it and supplies the
// abstract methods below) satisfies it automatically.
type BaselineRequirementsCheckerContract interface {
	// HasAdESProfile checks if the signature is conformant to a corresponding AdES profile.
	// Port of the abstract hasAdESProfile().
	HasAdESProfile() bool
	// HasBaselineBProfile checks if the signature has a corresponding BASELINE-B profile.
	// Port of the abstract hasBaselineBProfile().
	HasBaselineBProfile() bool
	// HasBaselineTProfile checks if the signature has a corresponding BASELINE-T profile.
	// Port of the abstract hasBaselineTProfile().
	HasBaselineTProfile() bool
	// HasBaselineLTProfile checks if the signature has a corresponding BASELINE-LT profile.
	// Port of the abstract hasBaselineLTProfile().
	HasBaselineLTProfile() bool
	// HasBaselineLTAProfile checks if the signature has a corresponding BASELINE-LTA profile.
	// Port of the abstract hasBaselineLTAProfile().
	HasBaselineLTAProfile() bool
	// HasExtendedBESProfile checks if the signature has a corresponding *AdES-BES profile.
	// Port of hasExtendedBESProfile() (default false, overridable).
	HasExtendedBESProfile() bool
	// HasExtendedEPESProfile checks if the signature has a corresponding *AdES-EPES profile.
	// Port of hasExtendedEPESProfile() (default false, overridable).
	HasExtendedEPESProfile() bool
	// HasExtendedTProfile checks if the signature has a corresponding *AdES-T profile.
	// Port of hasExtendedTProfile() (default false, overridable).
	HasExtendedTProfile() bool
	// HasExtendedCProfile checks if the signature has a corresponding *AdES-C profile.
	// Port of hasExtendedCProfile() (default false, overridable).
	HasExtendedCProfile() bool
	// HasExtendedXProfile checks if the signature has a corresponding *AdES-X profile.
	// Port of hasExtendedXProfile() (default false, overridable).
	HasExtendedXProfile() bool
	// HasExtendedXLProfile checks if the signature has a corresponding *AdES-XL profile.
	// Port of hasExtendedXLProfile() (default false, overridable).
	HasExtendedXLProfile() bool
	// HasExtendedAProfile checks if the signature has a corresponding *AdES-A profile.
	// Port of hasExtendedAProfile() (default false, overridable).
	HasExtendedAProfile() bool
	// HasExtendedERSProfile checks if the signature has a corresponding *AdES-E-ERS profile.
	// Port of hasExtendedERSProfile() (default false, overridable).
	HasExtendedERSProfile() bool
}

// BaselineRequirementsCheckerOverrides declares the single operation
// BaselineRequirementsChecker calls back into virtually: a concrete checker (a future phase's
// CAdESBaselineRequirementsChecker etc.) registers itself with InitBaselineRequirementsChecker
// so the base can dispatch it, the way model.TokenBase dispatches to model.TokenOverrides via
// InitToken. The five profile-detecting methods above (HasAdESProfile..HasBaselineLTAProfile)
// are genuinely abstract in Java - never called from within BaselineRequirementsChecker itself -
// so they need no dispatch entry here; a concrete checker simply defines them directly and Go's
// ordinary method promotion makes BaselineRequirementsCheckerContract see them.
type BaselineRequirementsCheckerOverrides interface {
	// ContainsLTLevelCertificates verifies whether the signature contains some of the LT-/XL
	// level attributes. Port of the protected containsLTLevelCertificates() (default false,
	// overridable), called back into by MinimalLTRequirement.
	ContainsLTLevelCertificates() bool

	// GetBaselineSignatureForm returns the signature form used to pick which CMS-attribute
	// cardinality rule applies where CAdES and PAdES disagree - the signing-time attribute is
	// the one case: EN 319 122-1 (CAdES-BASELINE-B) requires it present (cardinality == 1),
	// EN 319 142-1 (PAdES-BASELINE-B) requires it absent (cardinality == 0). Port of the
	// protected getBaselineSignatureForm(), which upstream's CAdESBaselineRequirementsChecker
	// (returning SignatureForm.CAdES) declares and both PAdESBaselineRequirementsChecker AND
	// CMSForPAdESBaselineRequirementsChecker override (returning SignatureForm.PAdES) purely
	// through ordinary Java virtual dispatch - not part of upstream's own BaselineRequirements
	// Checker base class at all. It is added to this port's cross-package override contract
	// instead, because cades.CAdESBaselineRequirementsChecker.cmsBaselineBRequirements() (the
	// shared CMS-attribute check both cades.CAdESSignature and, through
	// pades.CMSForPAdESBaselineRequirementsChecker, PAdES signatures run) needs to resolve it
	// virtually across the cades/pades package boundary the very same way MinimalLTRequirement
	// above resolves ContainsLTLevelCertificates - Go has no cross-package method-override
	// mechanism to fall back on. Defaults to "" (never consulted outside cades's own
	// cmsBaselineBRequirements(), so XAdES/JAdES concrete checkers need no override at all).
	GetBaselineSignatureForm() enumerations.SignatureForm
}

// BaselineRequirementsChecker checks conformance of a signature to the requested baseline
// format. AS is the Java type parameter "AS extends DefaultAdvancedSignature", bounded here by
// BaselineRequirementsCheckerSignature.
type BaselineRequirementsChecker[AS BaselineRequirementsCheckerSignature] struct {
	// overrides points back at the concrete checker; see InitBaselineRequirementsChecker.
	overrides BaselineRequirementsCheckerOverrides

	// signature is the signature object. Port of the protected final AS signature field; Go
	// callers reach it through Signature() since Go has no protected visibility.
	signature AS

	// offlineCertificateVerifier is the offline copy of a CertificateVerifier.
	// Port of the protected final CertificateVerifier offlineCertificateVerifier field.
	offlineCertificateVerifier CertificateVerifier

	// validationContext caches the ValidationContext so validation runs only once.
	// Port of the private ValidationContext validationContext field.
	validationContext ValidationContext
}

// NewBaselineRequirementsCheckerBase builds the base state a subclass embeds, without a
// CertificateVerifier (to be used for B-level validation only). Port of the protected
// BaselineRequirementsChecker(AS) constructor; the subclass constructor must follow it with
// InitBaselineRequirementsChecker.
func NewBaselineRequirementsCheckerBase[AS BaselineRequirementsCheckerSignature](signatureValue AS) BaselineRequirementsChecker[AS] {
	return BaselineRequirementsChecker[AS]{signature: signatureValue}
}

// NewBaselineRequirementsCheckerBaseWithVerifier builds the base state a subclass embeds. Port
// of the protected BaselineRequirementsChecker(AS, CertificateVerifier) constructor; the
// subclass constructor must follow it with InitBaselineRequirementsChecker.
func NewBaselineRequirementsCheckerBaseWithVerifier[AS BaselineRequirementsCheckerSignature](signatureValue AS, offlineCertificateVerifier CertificateVerifier) BaselineRequirementsChecker[AS] {
	return BaselineRequirementsChecker[AS]{signature: signatureValue, offlineCertificateVerifier: offlineCertificateVerifier}
}

// InitBaselineRequirementsChecker registers the concrete checker with its base so that the base
// can dispatch ContainsLTLevelCertificates. Every concrete subclass constructor must call this
// once.
func (b *BaselineRequirementsChecker[AS]) InitBaselineRequirementsChecker(overrides BaselineRequirementsCheckerOverrides) {
	b.overrides = overrides
}

// baselineRequirementsCheckerOverrides returns the registered overrides, panicking when the
// concrete checker forgot to call InitBaselineRequirementsChecker - Java's abstract class can
// never be instantiated bare, so there is no legitimate bare-instance fallback here (unlike
// e.g. TimestampIdentifierBuilder, which is a concrete, directly-instantiable Java class).
func (b *BaselineRequirementsChecker[AS]) baselineRequirementsCheckerOverrides() BaselineRequirementsCheckerOverrides {
	if b.overrides == nil {
		panic("BaselineRequirementsChecker was not initialised: the concrete checker must call InitBaselineRequirementsChecker in its constructor")
	}
	return b.overrides
}

// Signature returns the signature object being checked. The Go counterpart of reading Java's
// protected final signature field from a subclass.
func (b *BaselineRequirementsChecker[AS]) Signature() AS {
	return b.signature
}

// HasExtendedBESProfile checks if the signature has a corresponding *AdES-BES profile.
// Port of hasExtendedBESProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedBESProfile() bool { return false }

// HasExtendedEPESProfile checks if the signature has a corresponding *AdES-EPES profile.
// Port of hasExtendedEPESProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedEPESProfile() bool { return false }

// HasExtendedTProfile checks if the signature has a corresponding *AdES-T profile.
// Port of hasExtendedTProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedTProfile() bool { return false }

// HasExtendedCProfile checks if the signature has a corresponding *AdES-C profile.
// Port of hasExtendedCProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedCProfile() bool { return false }

// HasExtendedXProfile checks if the signature has a corresponding *AdES-X profile.
// Port of hasExtendedXProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedXProfile() bool { return false }

// HasExtendedXLProfile checks if the signature has a corresponding *AdES-XL profile.
// Port of hasExtendedXLProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedXLProfile() bool { return false }

// HasExtendedAProfile checks if the signature has a corresponding *AdES-A profile.
// Port of hasExtendedAProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedAProfile() bool { return false }

// HasExtendedERSProfile checks if the signature has a corresponding *AdES-E-ERS profile.
// Port of hasExtendedERSProfile(); not implemented by default.
func (b *BaselineRequirementsChecker[AS]) HasExtendedERSProfile() bool { return false }

// ContainsLTLevelCertificates verifies whether the signature contains some of the LT-/XL level
// attributes. Port of the protected containsLTLevelCertificates(); FALSE by default.
func (b *BaselineRequirementsChecker[AS]) ContainsLTLevelCertificates() bool { return false }

// GetBaselineSignatureForm is the default, unset ("") signature form; see the
// BaselineRequirementsCheckerOverrides doc comment on GetBaselineSignatureForm. Overridden by
// cades.CAdESBaselineRequirementsChecker, pades.PAdESBaselineRequirementsChecker, and
// pades.CMSForPAdESBaselineRequirementsChecker; never overridden (nor consulted) by XAdES/JAdES.
func (b *BaselineRequirementsChecker[AS]) GetBaselineSignatureForm() enumerations.SignatureForm {
	return ""
}

// BaselineSignatureForm resolves GetBaselineSignatureForm() through the registered overrides.
// Exported (unlike baselineRequirementsCheckerOverrides itself) so cross-package call sites -
// cades.CAdESBaselineRequirementsChecker.cmsBaselineBRequirements(), reached directly from the
// cades package and, through embedding, from pades.CMSForPAdESBaselineRequirementsChecker too -
// can consult it without reaching into this package's unexported overrides field/accessor.
func (b *BaselineRequirementsChecker[AS]) BaselineSignatureForm() enumerations.SignatureForm {
	return b.baselineRequirementsCheckerOverrides().GetBaselineSignatureForm()
}

// SignatureTimestampsCreatedBeforeSignCertExpiration checks whether signature timestamps have
// been created before expiration of the signing-certificate used to create the signature.
// Port of the protected signatureTimestampsCreatedBeforeSignCertExpiration().
func (b *BaselineRequirementsChecker[AS]) SignatureTimestampsCreatedBeforeSignCertExpiration() bool {
	signingCertificate := b.signature.SigningCertificateToken()
	if signingCertificate != nil {
		notAfter := signingCertificate.NotAfter()
		for _, timestampToken := range b.signature.SignatureTimestamps() {
			if notAfter.Before(timestampToken.GenerationTime()) {
				return false
			}
		}
	}
	return true
}

// MinimalTRequirement checks the minimal requirement to satisfy T-profile for AdES signatures.
// Port of the protected minimalTRequirement().
//
// slf4j logging (trace/warn) is dropped per PORTING.md; both branches keep their behaviour.
func (b *BaselineRequirementsChecker[AS]) MinimalTRequirement() bool {
	// SignatureTimeStamp (Cardinality >= 1)
	if utils.IsCollectionEmpty(b.signature.SignatureTimestamps()) {
		return false
	}
	signingCertificate := b.signature.SigningCertificateToken()
	if signingCertificate != nil {
		notAfter := signingCertificate.NotAfter()
		for _, timestampToken := range b.signature.SignatureTimestamps() {
			if !timestampToken.CreationDate().Before(notAfter) {
				return false
			}
		}
	}
	return true
}

// MinimalLTRequirement checks the minimal requirement to satisfy LT-profile for AdES signatures.
// Port of the public minimalLTRequirement().
//
// Panics with the Java message when offlineCertificateVerifier is missing (Objects.requireNonNull).
func (b *BaselineRequirementsChecker[AS]) MinimalLTRequirement() bool {
	if b.offlineCertificateVerifier == nil {
		panic("offlineCertificateVerifier cannot be null for LT-level verification!")
	}

	certificateSources := b.CertificateSourcesExceptLastArchiveTimestamp()
	certificateFound := certificateSources.NumberOfCertificates() > 0
	allSelfSigned := certificateFound && certificateSources.IsAllSelfSigned()

	emptyCRLs := len(b.signature.CompleteCRLSource().AllRevocationBinaries()) == 0
	emptyOCSPs := len(b.signature.CompleteOCSPSource().AllRevocationBinaries()) == 0
	emptyRevocation := emptyCRLs && emptyOCSPs

	minimalLTRequirement := !allSelfSigned && !emptyRevocation
	if minimalLTRequirement {
		// check presence of all revocation data
		return b.isAllRevocationDataPresent()
	}
	// Or one of the LT-/XL- profile properties are present (for self-signed cert chains only)
	if allSelfSigned {
		return b.baselineRequirementsCheckerOverrides().ContainsLTLevelCertificates()
	}
	return minimalLTRequirement
}

// CertificateSourcesExceptLastArchiveTimestamp returns a list of certificate sources with an
// exception to the last archive timestamp if applicable.
// Port of the protected getCertificateSourcesExceptLastArchiveTimestamp().
func (b *BaselineRequirementsChecker[AS]) CertificateSourcesExceptLastArchiveTimestamp() *spi.ListCertificateSource {
	certificateSource := spi.NewListCertificateSourceFromOne(b.signature.CertificateSource())
	certificateSource.AddAll(b.signature.TimestampSource().TimestampCertificateSourcesExceptLastArchiveTimestamp())
	certificateSource.AddAll(b.signature.CounterSignaturesCertificateSource())
	return certificateSource
}

// isAllRevocationDataPresent is the private isAllRevocationDataPresent().
func (b *BaselineRequirementsChecker[AS]) isAllRevocationDataPresent() bool {
	return b.ValidationContext().CheckAllRequiredRevocationDataPresent()
}

// ValidationContext returns a validated validation context. Port of the protected
// getValidationContext().
func (b *BaselineRequirementsChecker[AS]) ValidationContext() ValidationContext {
	if b.validationContext == nil {
		b.validationContext = NewSignatureValidationContext()
		b.validationContext.Initialize(b.offlineCertificateVerifier)

		b.validationContext.AddDocumentCertificateSource(b.signature.CompleteCertificateSource())
		b.validationContext.AddDocumentCRLSourceFromList(b.signature.CompleteCRLSource())
		b.validationContext.AddDocumentOCSPSourceFromList(b.signature.CompleteOCSPSource())

		baselineRequirementsCheckerAddSignatureForVerification(b.validationContext, b.signature)

		b.validationContext.Validate()
	}
	return b.validationContext
}

// baselineRequirementsCheckerAddSignatureForVerification is the private
// addSignatureForVerification(ValidationContext, AdvancedSignature). It is a free function
// (rather than a method) because Java declares the signature parameter as the AdvancedSignature
// interface, not AS - it is always called with b.signature, but its logic does not depend on AS.
func baselineRequirementsCheckerAddSignatureForVerification(validationContext ValidationContext, signatureValue AdvancedSignature) {
	signingCertificate := signatureValue.SigningCertificateToken()
	if signingCertificate != nil {
		validationContext.AddCertificateTokenForVerification(signingCertificate)
	} else {
		candidatesForSigningCertificate := signatureValue.CandidatesForSigningCertificate()
		certificateValidities := candidatesForSigningCertificate.CertificateValidityList()
		if utils.IsCollectionNotEmpty(certificateValidities) {
			for _, certificateValidity := range certificateValidities {
				if certificateValidity.IsValid() && certificateValidity.CertificateToken() != nil {
					validationContext.AddCertificateTokenForVerification(certificateValidity.CertificateToken())
				}
			}
		}
	}
	for _, timestampToken := range signatureValue.TimestampSource().AllTimestampsExceptLastArchiveTimestamp() {
		validationContext.AddTimestampTokenForVerification(timestampToken)
	}
}

// MinimalLTARequirement checks the minimal requirement to satisfy LTA-profile for AdES
// signatures. Port of the public minimalLTARequirement().
func (b *BaselineRequirementsChecker[AS]) MinimalLTARequirement() bool {
	// ArchiveTimeStamp (Cardinality >= 1)
	return utils.IsCollectionNotEmpty(b.signature.ArchiveTimestamps())
}

// ContainsSigningCertificate checks if the given collection of CertificateTokens contains the
// signing certificate for the signature.
// Port of the protected containsSigningCertificate(Collection).
func (b *BaselineRequirementsChecker[AS]) ContainsSigningCertificate(certificateTokens []*model.CertificateToken) bool {
	candidatesForSigningCertificate := b.signature.CandidatesForSigningCertificate()
	certificateValidity := candidatesForSigningCertificate.TheCertificateValidity()
	if certificateValidity != nil && certificateValidity.CertificateToken() != nil {
		signingCertificate := certificateValidity.CertificateToken()
		for _, certificate := range certificateTokens {
			if certificate.Equals(signingCertificate) {
				return true
			}
		}
	}
	return false
}

// IsSignaturePolicyIdentifierHashPresent checks if the signature contains a
// SignaturePolicyIdentifier containing a hash used to digest the signature policy.
// Port of the protected isSignaturePolicyIdentifierHashPresent().
func (b *BaselineRequirementsChecker[AS]) IsSignaturePolicyIdentifierHashPresent() bool {
	signaturePolicyIdentifier := b.signature.SignaturePolicy()
	if signaturePolicyIdentifier != nil {
		digest := signaturePolicyIdentifier.Digest()
		return digest.Algorithm() != ""
	}
	return false
}
