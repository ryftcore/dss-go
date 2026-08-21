// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/DefaultAdvancedSignature.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: CertificateVerifier and CertificateVerifierBuilder (Java spi.validation.*)
// are flattened into this same Go package by a sibling chunk of phase 2b; see
// baseline_requirements_checker.go's header for the assumed CertificateVerifier/ValidationContext
// shapes. CertificateVerifierBuilder is additionally assumed to expose, matching
// dss-spi/.../spi/validation/CertificateVerifierBuilder.java:
//
//	func NewCertificateVerifierBuilder(certificateVerifier CertificateVerifier) *CertificateVerifierBuilder
//	func (b *CertificateVerifierBuilder) BuildOfflineAndSilentCopy() CertificateVerifier
package validation

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DefaultAdvancedSignatureOverrides declares every operation DefaultAdvancedSignature calls
// back into virtually: the AdvancedSignature interface methods it deliberately leaves
// unimplemented (deferred to the format-specific final signature type of a later phase, e.g.
// CAdESSignature/XAdESSignature) that it nonetheless calls on itself, plus its own
// protected-abstract and protected-overridable hooks. It stands in for the virtual dispatch a
// Java abstract class gets for free; a concrete signature registers itself with
// DefaultAdvancedSignature.InitDefaultAdvancedSignature so that the base can reach them, the way
// model.TokenBase dispatches to model.TokenOverrides via InitToken.
//
// AdvancedSignature methods DefaultAdvancedSignature never calls on itself (e.g. SigningTime,
// SignaturePolicyStore, ContentType, SignatureValue, ReferenceValidations, ...) have no entry
// here: Go's ordinary embedding means the eventual concrete signature type just defines them
// directly, and DefaultAdvancedSignature need not know about them at all - exactly as the Java
// abstract class leaves them abstract without a matching internal caller.
type DefaultAdvancedSignatureOverrides interface {
	// SignatureForm specifies the format of the signature. Port of the abstract
	// getSignatureForm(), called back into by String().
	SignatureForm() enumerations.SignatureForm

	// SignatureAlgorithm retrieves the signature algorithm used for generating the signature.
	// Port of the abstract getSignatureAlgorithm(), called back into by EncryptionAlgorithm and
	// DigestAlgorithm.
	SignatureAlgorithm() enumerations.SignatureAlgorithm

	// CertificateSource gets a certificate source which contains ALL certificates embedded in
	// the signature. Port of the abstract getCertificateSource().
	CertificateSource() *spi.SignatureCertificateSource

	// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature.
	// Port of the abstract getCRLSource().
	CRLSource() spi.OfflineRevocationSource[revocation.CRL]

	// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the
	// signature. Port of the abstract getOCSPSource().
	OCSPSource() spi.OfflineRevocationSource[revocation.OCSP]

	// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded
	// in the signature. Port of the abstract getTimestampSource().
	TimestampSource() TimestampSource

	// CounterSignatures returns a list of counter signatures applied to this signature.
	// Port of the abstract getCounterSignatures().
	CounterSignatures() []AdvancedSignature

	// CheckSignatureIntegrity verifies the signature integrity, storing the outcome with
	// SetSignatureCryptographicVerification. Port of the abstract checkSignatureIntegrity().
	CheckSignatureIntegrity()

	// ClaimedSignerRoles returns the claimed roles of the signer.
	// Port of the abstract getClaimedSignerRoles().
	ClaimedSignerRoles() []*signature.SignerRole

	// CertifiedSignerRoles returns the certified roles of the signer.
	// Port of the abstract getCertifiedSignerRoles().
	CertifiedSignerRoles() []*signature.SignerRole

	// SignedAssertions returns the list of embedded signed assertions.
	// Port of the abstract getSignedAssertions().
	SignedAssertions() []*signature.SignerRole

	// SignatureIdentifierBuilder returns a builder to define and build a signature Id.
	// Port of the protected abstract getSignatureIdentifierBuilder().
	SignatureIdentifierBuilder() SignatureIdentifierBuilder

	// ValidateStructure processes the structure validation of the signature. Port of the
	// protected validateStructure() (default: nil/empty, overridable).
	ValidateStructure() []string

	// FindSignatureScopes finds signature scopes. Port of the protected abstract
	// findSignatureScopes().
	FindSignatureScopes() []scope.SignatureScope

	// BuildSignaturePolicy extracts a signature policy from a signature and builds the object.
	// Port of the protected abstract buildSignaturePolicy().
	BuildSignaturePolicy() *signature.SignaturePolicy

	// BuildSignatureDigestReference builds a new SignatureDigestReference according to the
	// applicable signature format rules. Port of the protected abstract
	// buildSignatureDigestReference(DigestAlgorithm).
	BuildSignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.SignatureDigestReference

	// CreateBaselineRequirementsChecker instantiates a BaselineRequirementsChecker according to
	// the signature format. Port of the protected abstract
	// createBaselineRequirementsChecker(CertificateVerifier); Java's raw-typed return
	// (@SuppressWarnings("rawtypes")) becomes BaselineRequirementsCheckerContract, the same
	// erasure baseline_requirements_checker.go uses for DefaultAdvancedSignature's cached field.
	CreateBaselineRequirementsChecker(certificateVerifier CertificateVerifier) BaselineRequirementsCheckerContract
}

// DefaultAdvancedSignature is a common implementation of AdvancedSignature. Concrete signature
// types (CAdES/XAdES/JAdES/PAdES, ported in later phases) embed it and register themselves with
// InitDefaultAdvancedSignature; it does not, on its own, satisfy the AdvancedSignature interface
// (several of that interface's methods are left to the concrete signature type entirely, with no
// field or dispatch entry here at all - matching how the abstract Java class leaves them
// abstract without calling them internally).
//
// serialVersionUID and java.io.Serializable are dropped (no Go counterpart).
type DefaultAdvancedSignature struct {
	// overrides points back at the concrete signature; see InitDefaultAdvancedSignature.
	overrides DefaultAdvancedSignatureOverrides

	// detachedContents is, in case of a detached signature, the signed document.
	detachedContents []model.DSSDocument

	// containerContents is, in case of an ASiC-S signature, the archive or manifest content.
	containerContents []model.DSSDocument

	// manifestFile is, in case of an ASiC-E signature, the found related manifest file.
	manifestFile *model.ManifestFile

	// cachedReferenceValidations is a cache slot for a list of reference validations (reference
	// tag for XAdES or message-digest for CAdES). Port of the protected referenceValidations
	// field; DefaultAdvancedSignature itself never reads or writes it (Java leaves
	// getReferenceValidations() entirely abstract), so it is exposed only through
	// CachedReferenceValidations/SetCachedReferenceValidations for a concrete signature (in
	// another package) to use as it sees fit, the Go counterpart of inheriting a protected field.
	cachedReferenceValidations []*model.ReferenceValidation

	// signatureCryptographicVerification contains the result of the signature mathematical
	// validation. It is initialised when CheckSignatureIntegrity is called.
	signatureCryptographicVerification *signature.SignatureCryptographicVerification

	// structureValidationMessages is a list of error messages from a structure validation.
	structureValidationMessages []string

	// signingCertificateSource is the certificate source of a signing certificate.
	signingCertificateSource spi.CertificateSource

	// offlineCertificateSource is a cache slot for the offline signature certificate source.
	// Port of the protected offlineCertificateSource field, reset by ResetCertificateSource;
	// DefaultAdvancedSignature never populates it itself (a concrete CertificateSource()
	// override does), so it is exposed through OfflineCertificateSource/
	// SetOfflineCertificateSource for that override to use.
	offlineCertificateSource *spi.SignatureCertificateSource

	// signatureCRLSource is a cache slot for the offline signature CRL source.
	// Port of the protected signatureCRLSource field, reset by ResetRevocationSources; exposed
	// through SignatureCRLSource/SetSignatureCRLSource for the same reason as
	// offlineCertificateSource.
	signatureCRLSource spi.OfflineRevocationSource[revocation.CRL]

	// signatureOCSPSource is a cache slot for the offline signature OCSP source.
	// Port of the protected signatureOCSPSource field, reset by ResetRevocationSources; exposed
	// through SignatureOCSPSource/SetSignatureOCSPSource.
	signatureOCSPSource spi.OfflineRevocationSource[revocation.OCSP]

	// signatureTimestampSource is a cache slot for the offline signature timestamp source.
	// Port of the protected signatureTimestampSource field, reset by ResetTimestampSource;
	// exposed through SignatureTimestampSource/SetSignatureTimestampSource.
	signatureTimestampSource TimestampSource

	// cachedCounterSignatures is a cache slot for the list of embedded counter signatures.
	// Port of the protected counterSignatures field; DefaultAdvancedSignature never populates it
	// itself (Java leaves getCounterSignatures() entirely abstract), so it is exposed through
	// CachedCounterSignatures/SetCachedCounterSignatures.
	cachedCounterSignatures []AdvancedSignature

	// masterSignature is the master signature, in case the current signature is a counter
	// signature.
	masterSignature AdvancedSignature

	// eaa is the EAA in case of a key binding signature.
	eaa EAA

	// keyBindingSignature indicates whether the signature is a key binding signature.
	keyBindingSignature bool

	// signaturePolicy is the SignaturePolicy identifier.
	signaturePolicy *signature.SignaturePolicy

	// signatureScopes is a list of found SignatureScopes.
	signatureScopes []scope.SignatureScope

	// filename is the name of a signature file.
	filename string

	// signatureIdentifier is the cached unique signature identifier.
	signatureIdentifier *SignatureIdentifier

	// signatureDigestReferences caches computed SignatureDigestReference's as defined in ETSI
	// TS 119 102-2 ch. "4.1.1.5 Signature Reference".
	signatureDigestReferences map[enumerations.DigestAlgorithm]*signature.SignatureDigestReference

	// baselineRequirementsChecker performs a conformance check for the signature to a given
	// profile. "transient" (Java) has no Go counterpart since this port has no serialization.
	baselineRequirementsChecker BaselineRequirementsCheckerContract

	// signingCertificateToken caches the signing certificate token.
	signingCertificateToken *model.CertificateToken
}

// NewDefaultAdvancedSignatureBase instantiates the base state of a signature with null/zero
// values. Port of the protected default constructor. The concrete signature must still call
// InitDefaultAdvancedSignature.
func NewDefaultAdvancedSignatureBase() DefaultAdvancedSignature {
	return DefaultAdvancedSignature{}
}

// InitDefaultAdvancedSignature registers the concrete signature with its base so that the base
// can dispatch the operations listed on DefaultAdvancedSignatureOverrides. It must be called
// exactly once, by the concrete signature's constructor, before any other method.
func (s *DefaultAdvancedSignature) InitDefaultAdvancedSignature(overrides DefaultAdvancedSignatureOverrides) {
	s.overrides = overrides
}

// defaultAdvancedSignatureOverrides returns the registered overrides, panicking when the
// concrete signature forgot to call InitDefaultAdvancedSignature.
func (s *DefaultAdvancedSignature) defaultAdvancedSignatureOverrides() DefaultAdvancedSignatureOverrides {
	if s.overrides == nil {
		panic("DefaultAdvancedSignature was not initialised: the concrete signature must call InitDefaultAdvancedSignature in its constructor")
	}
	return s.overrides
}

// EncryptionAlgorithm retrieves the encryption algorithm used for generating the signature.
// Port of getEncryptionAlgorithm().
func (s *DefaultAdvancedSignature) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	signatureAlgorithm := s.defaultAdvancedSignatureOverrides().SignatureAlgorithm()
	if signatureAlgorithm == "" {
		return ""
	}
	return signatureAlgorithm.EncryptionAlgorithm()
}

// DigestAlgorithm retrieves the digest algorithm used for generating the signature.
// Port of getDigestAlgorithm().
func (s *DefaultAdvancedSignature) DigestAlgorithm() enumerations.DigestAlgorithm {
	signatureAlgorithm := s.defaultAdvancedSignatureOverrides().SignatureAlgorithm()
	if signatureAlgorithm == "" {
		return ""
	}
	return signatureAlgorithm.DigestAlgorithm()
}

// SigningCertificateSource gets a signing certificate source, when provided.
// Port of getSigningCertificateSource().
func (s *DefaultAdvancedSignature) SigningCertificateSource() spi.CertificateSource {
	return s.signingCertificateSource
}

// SetSigningCertificateSource sets a certificate source which allows finding the signing
// certificate by kid or certificate's digest. Port of setSigningCertificateSource(CertificateSource).
func (s *DefaultAdvancedSignature) SetSigningCertificateSource(signingCertificateSource spi.CertificateSource) {
	s.signingCertificateSource = signingCertificateSource
}

// Filename returns the signature filename. Port of getFilename().
func (s *DefaultAdvancedSignature) Filename() string {
	return s.filename
}

// SetFilename sets the signature filename. Port of setFilename(String).
func (s *DefaultAdvancedSignature) SetFilename(filename string) {
	s.filename = filename
}

// DetachedContents returns, in case of a detached signature, the signed contents.
// Port of getDetachedContents().
func (s *DefaultAdvancedSignature) DetachedContents() []model.DSSDocument {
	return s.detachedContents
}

// SetDetachedContents sets the signed contents in the case of the detached signature.
// Port of setDetachedContents(List).
func (s *DefaultAdvancedSignature) SetDetachedContents(detachedContents []model.DSSDocument) {
	s.detachedContents = detachedContents
}

// ContainerContents returns, in case of ASiC-S signature, the archive container documents.
// Port of getContainerContents().
func (s *DefaultAdvancedSignature) ContainerContents() []model.DSSDocument {
	return s.containerContents
}

// SetContainerContents sets the archive container contents in the case of ASiC-S signature.
// Port of setContainerContents(List).
func (s *DefaultAdvancedSignature) SetContainerContents(containerContents []model.DSSDocument) {
	s.containerContents = containerContents
}

// ManifestFile returns the related ManifestFile in the case of ASiC-E signature.
// Port of getManifestFile().
func (s *DefaultAdvancedSignature) ManifestFile() *model.ManifestFile {
	return s.manifestFile
}

// SetManifestFile sets a manifest file in the case of ASiC-E signature.
// Port of setManifestFile(ManifestFile).
func (s *DefaultAdvancedSignature) SetManifestFile(manifestFile *model.ManifestFile) {
	s.manifestFile = manifestFile
}

// DSSID returns the DSS unique identifier of the signature, building it via the registered
// SignatureIdentifierBuilder on first use. Port of getDSSId(); Java narrows the return type to
// SignatureIdentifier, which Go cannot express, so callers needing it assert on the result
// (matching the convention documented on model.IdentifierBasedObject).
func (s *DefaultAdvancedSignature) DSSID() model.Identifier {
	return s.dssID()
}

// dssID is the internal, narrowly typed counterpart of DSSID.
func (s *DefaultAdvancedSignature) dssID() *SignatureIdentifier {
	if s.signatureIdentifier == nil {
		s.signatureIdentifier = s.defaultAdvancedSignatureOverrides().SignatureIdentifierBuilder().BuildSignatureIdentifier()
	}
	return s.signatureIdentifier
}

// ID returns the DSS unique signature id. It allows unambiguously identifying each signature.
// Port of getId().
func (s *DefaultAdvancedSignature) ID() string {
	return s.dssID().AsXmlID()
}

// CompleteCertificateSource gets a ListCertificateSource representing a merged source from the
// signature certificate source and all included signature timestamp objects.
// Port of getCompleteCertificateSource().
func (s *DefaultAdvancedSignature) CompleteCertificateSource() *spi.ListCertificateSource {
	certificateSource := spi.NewListCertificateSourceFromOne(s.defaultAdvancedSignatureOverrides().CertificateSource())
	certificateSource.AddAll(s.defaultAdvancedSignatureOverrides().TimestampSource().TimestampCertificateSources())
	certificateSource.AddAll(s.CounterSignaturesCertificateSource())
	return certificateSource
}

// CompleteCRLSource gets a ListRevocationSource representing a merged source from the signature
// CRL source and all included signature timestamp objects. Port of getCompleteCRLSource().
func (s *DefaultAdvancedSignature) CompleteCRLSource() *spi.ListRevocationSource[revocation.CRL] {
	crlSource := spi.NewListRevocationSourceFrom(s.defaultAdvancedSignatureOverrides().CRLSource())
	crlSource.AddAllFrom(s.defaultAdvancedSignatureOverrides().TimestampSource().TimestampCRLSources())
	crlSource.AddAllFrom(s.CounterSignaturesCRLSource())
	return crlSource
}

// CompleteOCSPSource gets a ListRevocationSource representing a merged source from the signature
// OCSP source and all included signature timestamp objects. Port of getCompleteOCSPSource().
func (s *DefaultAdvancedSignature) CompleteOCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	ocspSource := spi.NewListRevocationSourceFrom(s.defaultAdvancedSignatureOverrides().OCSPSource())
	ocspSource.AddAllFrom(s.defaultAdvancedSignatureOverrides().TimestampSource().TimestampOCSPSources())
	ocspSource.AddAllFrom(s.CounterSignaturesOCSPSource())
	return ocspSource
}

// CounterSignaturesCertificateSource returns a merged certificate source for values incorporated
// within counter signatures. Port of the protected getCounterSignaturesCertificateSource().
func (s *DefaultAdvancedSignature) CounterSignaturesCertificateSource() *spi.ListCertificateSource {
	certificateSource := spi.NewListCertificateSource()
	for _, counterSignature := range s.defaultAdvancedSignatureOverrides().CounterSignatures() {
		certificateSource.AddAll(counterSignature.CompleteCertificateSource())
	}
	return certificateSource
}

// CounterSignaturesCRLSource returns a merged CRL source for values incorporated within counter
// signatures. Port of the protected getCounterSignaturesCRLSource().
func (s *DefaultAdvancedSignature) CounterSignaturesCRLSource() *spi.ListRevocationSource[revocation.CRL] {
	crlSource := spi.NewListRevocationSource[revocation.CRL]()
	for _, counterSignature := range s.defaultAdvancedSignatureOverrides().CounterSignatures() {
		crlSource.AddAllFrom(counterSignature.CompleteCRLSource())
	}
	return crlSource
}

// CounterSignaturesOCSPSource returns a merged OCSP source for values incorporated within
// counter signatures. Port of the protected getCounterSignaturesOCSPSource().
func (s *DefaultAdvancedSignature) CounterSignaturesOCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	ocspSource := spi.NewListRevocationSource[revocation.OCSP]()
	for _, counterSignature := range s.defaultAdvancedSignatureOverrides().CounterSignatures() {
		ocspSource.AddAllFrom(counterSignature.CompleteOCSPSource())
	}
	return ocspSource
}

// ResetCertificateSource resets the source of certificates. It must be called when any
// certificate is added to the KeyInfo or CertificateValues (XAdES), or 'xVals' (JAdES).
// NOTE: used in XAdES and JAdES. Port of resetCertificateSource().
func (s *DefaultAdvancedSignature) ResetCertificateSource() {
	s.offlineCertificateSource = nil
}

// ResetRevocationSources resets the sources of the revocation data. It must be called when the
// -LT level is created. NOTE: used in XAdES and JAdES. Port of resetRevocationSources().
func (s *DefaultAdvancedSignature) ResetRevocationSources() {
	s.signatureCRLSource = nil
	s.signatureOCSPSource = nil
}

// ResetTimestampSource resets the timestamp source. It must be called when the -LT level is
// created. NOTE: used in XAdES and JAdES. Port of resetTimestampSource().
func (s *DefaultAdvancedSignature) ResetTimestampSource() {
	s.signatureTimestampSource = nil
}

// OfflineCertificateSource returns the cached offline signature certificate source. The Go
// counterpart of reading Java's protected offlineCertificateSource field from a subclass.
func (s *DefaultAdvancedSignature) OfflineCertificateSource() *spi.SignatureCertificateSource {
	return s.offlineCertificateSource
}

// SetOfflineCertificateSource sets the cached offline signature certificate source. The Go
// counterpart of writing Java's protected offlineCertificateSource field from a subclass.
func (s *DefaultAdvancedSignature) SetOfflineCertificateSource(offlineCertificateSource *spi.SignatureCertificateSource) {
	s.offlineCertificateSource = offlineCertificateSource
}

// SignatureCRLSource returns the cached offline signature CRL source. The Go counterpart of
// reading Java's protected signatureCRLSource field from a subclass.
func (s *DefaultAdvancedSignature) SignatureCRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	return s.signatureCRLSource
}

// SetSignatureCRLSource sets the cached offline signature CRL source. The Go counterpart of
// writing Java's protected signatureCRLSource field from a subclass.
func (s *DefaultAdvancedSignature) SetSignatureCRLSource(signatureCRLSource spi.OfflineRevocationSource[revocation.CRL]) {
	s.signatureCRLSource = signatureCRLSource
}

// SignatureOCSPSource returns the cached offline signature OCSP source. The Go counterpart of
// reading Java's protected signatureOCSPSource field from a subclass.
func (s *DefaultAdvancedSignature) SignatureOCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	return s.signatureOCSPSource
}

// SetSignatureOCSPSource sets the cached offline signature OCSP source. The Go counterpart of
// writing Java's protected signatureOCSPSource field from a subclass.
func (s *DefaultAdvancedSignature) SetSignatureOCSPSource(signatureOCSPSource spi.OfflineRevocationSource[revocation.OCSP]) {
	s.signatureOCSPSource = signatureOCSPSource
}

// SignatureTimestampSource returns the cached offline signature timestamp source. The Go
// counterpart of reading Java's protected signatureTimestampSource field from a subclass.
func (s *DefaultAdvancedSignature) SignatureTimestampSource() TimestampSource {
	return s.signatureTimestampSource
}

// SetSignatureTimestampSource sets the cached offline signature timestamp source. The Go
// counterpart of writing Java's protected signatureTimestampSource field from a subclass.
func (s *DefaultAdvancedSignature) SetSignatureTimestampSource(signatureTimestampSource TimestampSource) {
	s.signatureTimestampSource = signatureTimestampSource
}

// CachedCounterSignatures returns the cached list of embedded counter signatures. The Go
// counterpart of reading Java's protected counterSignatures field from a subclass.
func (s *DefaultAdvancedSignature) CachedCounterSignatures() []AdvancedSignature {
	return s.cachedCounterSignatures
}

// SetCachedCounterSignatures sets the cached list of embedded counter signatures. The Go
// counterpart of writing Java's protected counterSignatures field from a subclass.
func (s *DefaultAdvancedSignature) SetCachedCounterSignatures(counterSignatures []AdvancedSignature) {
	s.cachedCounterSignatures = counterSignatures
}

// CachedReferenceValidations returns the cached list of reference validations. The Go
// counterpart of reading Java's protected referenceValidations field from a subclass.
func (s *DefaultAdvancedSignature) CachedReferenceValidations() []*model.ReferenceValidation {
	return s.cachedReferenceValidations
}

// SetCachedReferenceValidations sets the cached list of reference validations. The Go
// counterpart of writing Java's protected referenceValidations field from a subclass.
func (s *DefaultAdvancedSignature) SetCachedReferenceValidations(referenceValidations []*model.ReferenceValidation) {
	s.cachedReferenceValidations = referenceValidations
}

// ETSI TS 101 733 V2.2.1 (2013-04) 5.6.3 Signature Verification Process: the public key from the
// first certificate identified in the sequence of certificate identifiers from SigningCertificate
// shall be the key used to verify the digital signature.
//
// CandidatesForSigningCertificate gets an object containing the signing certificate or
// information indicating why it is impossible to extract it from the signature.
// Port of getCandidatesForSigningCertificate().
func (s *DefaultAdvancedSignature) CandidatesForSigningCertificate() *spi.CandidatesForSigningCertificate {
	return s.defaultAdvancedSignatureOverrides().CertificateSource().CandidatesForSigningCertificate(s.signingCertificateSource)
}

// InitBaselineRequirementsChecker creates an offline copy of certificateVerifier and
// instantiates a BaselineRequirementsChecker. Port of initBaselineRequirementsChecker(CertificateVerifier).
func (s *DefaultAdvancedSignature) InitBaselineRequirementsChecker(certificateVerifier CertificateVerifier) {
	offlineCertificateVerifier := NewCertificateVerifierBuilder(certificateVerifier).BuildOfflineAndSilentCopy()
	s.baselineRequirementsChecker = s.defaultAdvancedSignatureOverrides().CreateBaselineRequirementsChecker(offlineCertificateVerifier)
}

// Certificates returns an unmodifiable list of all certificate tokens encapsulated in the
// signature. Port of getCertificates(), see AdvancedSignature#getCertificates().
func (s *DefaultAdvancedSignature) Certificates() []*model.CertificateToken {
	return s.defaultAdvancedSignatureOverrides().CertificateSource().Certificates()
}

// SetMasterSignature indicates the master signature: the current signature is a counter
// signature. Port of setMasterSignature(AdvancedSignature).
func (s *DefaultAdvancedSignature) SetMasterSignature(masterSignature AdvancedSignature) {
	s.masterSignature = masterSignature
}

// MasterSignature gets the master signature. Port of getMasterSignature().
func (s *DefaultAdvancedSignature) MasterSignature() AdvancedSignature {
	return s.masterSignature
}

// IsCounterSignature checks if the current signature is a counter signature (i.e. has a master
// signature). Port of isCounterSignature().
func (s *DefaultAdvancedSignature) IsCounterSignature() bool {
	return s.masterSignature != nil
}

// EAA gets the EAA of an EAA issuing or key binding signature. Port of getEAA().
func (s *DefaultAdvancedSignature) EAA() EAA {
	return s.eaa
}

// SetEAA sets the EAA presentation of the EAA issuing or key binding signature. Port of setEAA(EAA).
func (s *DefaultAdvancedSignature) SetEAA(eaa EAA) {
	s.eaa = eaa
}

// IsKeyBindingSignature checks if the current signature is a key binding signature.
// NOTE: used for EAA tokens. Port of isKeyBindingSignature().
func (s *DefaultAdvancedSignature) IsKeyBindingSignature() bool {
	return s.keyBindingSignature
}

// SetKeyBindingSignature sets whether the current signature is a key binding signature.
// NOTE: used for EAA tokens. Port of setKeyBindingSignature(boolean).
func (s *DefaultAdvancedSignature) SetKeyBindingSignature(keyBindingSignature bool) {
	s.keyBindingSignature = keyBindingSignature
}

// SignatureCryptographicVerification gets the signature's cryptographic validation result,
// running CheckSignatureIntegrity on first use. Port of getSignatureCryptographicVerification().
func (s *DefaultAdvancedSignature) SignatureCryptographicVerification() *signature.SignatureCryptographicVerification {
	if s.signatureCryptographicVerification == nil {
		s.defaultAdvancedSignatureOverrides().CheckSignatureIntegrity()
	}
	return s.signatureCryptographicVerification
}

// SetSignatureCryptographicVerification sets the result of the signature mathematical
// validation. The Go counterpart of writing Java's protected signatureCryptographicVerification
// field from a subclass's CheckSignatureIntegrity implementation.
func (s *DefaultAdvancedSignature) SetSignatureCryptographicVerification(verification *signature.SignatureCryptographicVerification) {
	s.signatureCryptographicVerification = verification
}

// SignerRoles returns the list of roles of the signer. Port of getSignerRoles().
func (s *DefaultAdvancedSignature) SignerRoles() []*signature.SignerRole {
	var signerRoles []*signature.SignerRole
	claimedSignerRoles := s.defaultAdvancedSignatureOverrides().ClaimedSignerRoles()
	if utils.IsCollectionNotEmpty(claimedSignerRoles) {
		signerRoles = append(signerRoles, claimedSignerRoles...)
	}
	certifiedSignerRoles := s.defaultAdvancedSignatureOverrides().CertifiedSignerRoles()
	if utils.IsCollectionNotEmpty(certifiedSignerRoles) {
		signerRoles = append(signerRoles, certifiedSignerRoles...)
	}
	signedAssertionSignerRoles := s.defaultAdvancedSignatureOverrides().SignedAssertions()
	if utils.IsCollectionNotEmpty(signedAssertionSignerRoles) {
		signerRoles = append(signerRoles, signedAssertionSignerRoles...)
	}
	return signerRoles
}

// SigningCertificateToken returns the signing certificate token or nil if there is no valid
// signing certificate. Note that to determine the signing certificate the signature must be
// validated: CheckSignatureIntegrity must be called. Port of getSigningCertificateToken().
func (s *DefaultAdvancedSignature) SigningCertificateToken() *model.CertificateToken {
	if s.signingCertificateToken == nil {
		// This ensures that candidatesForSigningCertificate has been initialised.
		candidatesForSigningCertificate := s.CandidatesForSigningCertificate()
		// This ensures that signatureCryptographicVerification has been initialised.
		s.signatureCryptographicVerification = s.SignatureCryptographicVerification()
		theCertificateValidity := candidatesForSigningCertificate.TheCertificateValidity()
		if theCertificateValidity != nil && theCertificateValidity.IsValid() {
			return theCertificateValidity.CertificateToken()
		}
		theBestCandidate := candidatesForSigningCertificate.TheBestCandidate()
		if theBestCandidate != nil {
			s.signingCertificateToken = theBestCandidate.CertificateToken()
		}
	}
	return s.signingCertificateToken
}

// StructureValidationResult returns a message if the structure validation fails: a list of
// error messages if validation fails, empty if structural validation succeeds.
// Port of getStructureValidationResult().
func (s *DefaultAdvancedSignature) StructureValidationResult() []string {
	if utils.IsCollectionEmpty(s.structureValidationMessages) {
		s.structureValidationMessages = s.defaultAdvancedSignatureOverrides().ValidateStructure()
	}
	return s.structureValidationMessages
}

// ValidateStructure processes the structure validation of the signature. Port of the protected
// validateStructure(); not implemented by default.
func (s *DefaultAdvancedSignature) ValidateStructure() []string {
	return nil
}

// SignatureScopes returns a list of found SignatureScopes, computing them via FindSignatureScopes
// on first use. Port of getSignatureScopes().
func (s *DefaultAdvancedSignature) SignatureScopes() []scope.SignatureScope {
	if s.signatureScopes == nil {
		s.signatureScopes = s.defaultAdvancedSignatureOverrides().FindSignatureScopes()
	}
	return s.signatureScopes
}

// ContentTimestamps returns the content timestamps. Port of getContentTimestamps().
func (s *DefaultAdvancedSignature) ContentTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().ContentTimestamps()
}

// SignatureTimestamps returns the signature timestamps. Port of getSignatureTimestamps().
func (s *DefaultAdvancedSignature) SignatureTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().SignatureTimestamps()
}

// TimestampsX1 returns the timestamps covering the digital signature, the certification path
// references and the revocation status references. Port of getTimestampsX1().
func (s *DefaultAdvancedSignature) TimestampsX1() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().TimestampsX1()
}

// TimestampsX2 returns the timestamp computed over CompleteCertificateRefs and
// CompleteRevocationRefs elements (XAdES example). Port of getTimestampsX2().
func (s *DefaultAdvancedSignature) TimestampsX2() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().TimestampsX2()
}

// ArchiveTimestamps returns the archive timestamps. Port of getArchiveTimestamps().
func (s *DefaultAdvancedSignature) ArchiveTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().ArchiveTimestamps()
}

// DocumentTimestamps returns a list of timestamps defined with the 'DocTimeStamp' type.
// NOTE: applicable only for PAdES. Port of getDocumentTimestamps().
func (s *DefaultAdvancedSignature) DocumentTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().DocumentTimestamps()
}

// DetachedTimestamps returns a list of detached timestamps. NOTE: used for ASiC with CAdES only.
// Port of getDetachedTimestamps().
func (s *DefaultAdvancedSignature) DetachedTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().DetachedTimestamps()
}

// AllTimestamps returns a list of all timestamps found in the signature. Port of getAllTimestamps().
func (s *DefaultAdvancedSignature) AllTimestamps() []*TimestampToken {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().AllTimestamps()
}

// EmbeddedEvidenceRecords returns a list of embedded evidence records. Port of getEmbeddedEvidenceRecords().
func (s *DefaultAdvancedSignature) EmbeddedEvidenceRecords() []EvidenceRecord {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().EmbeddedEvidenceRecords()
}

// AddExternalEvidenceRecord adds an evidence record covering the signature file.
// Port of addExternalEvidenceRecord(EvidenceRecord).
func (s *DefaultAdvancedSignature) AddExternalEvidenceRecord(evidenceRecord EvidenceRecord) {
	s.defaultAdvancedSignatureOverrides().TimestampSource().AddExternalEvidenceRecord(evidenceRecord)
}

// DetachedEvidenceRecords returns a list of detached evidence records. Port of getDetachedEvidenceRecords().
func (s *DefaultAdvancedSignature) DetachedEvidenceRecords() []EvidenceRecord {
	return s.defaultAdvancedSignatureOverrides().TimestampSource().DetachedEvidenceRecords()
}

// AllEvidenceRecords returns a list of all evidence records. Port of getAllEvidenceRecords().
func (s *DefaultAdvancedSignature) AllEvidenceRecords() []EvidenceRecord {
	var evidenceRecords []EvidenceRecord
	evidenceRecords = append(evidenceRecords, s.EmbeddedEvidenceRecords()...)
	evidenceRecords = append(evidenceRecords, s.DetachedEvidenceRecords()...)
	return evidenceRecords
}

// SignaturePolicy returns the Signature Policy OID from the signature, building it via
// BuildSignaturePolicy on first use. Port of getSignaturePolicy().
func (s *DefaultAdvancedSignature) SignaturePolicy() *signature.SignaturePolicy {
	if s.signaturePolicy == nil {
		s.signaturePolicy = s.defaultAdvancedSignatureOverrides().BuildSignaturePolicy()
	}
	return s.signaturePolicy
}

// SignatureDigestReference returns a signature reference element as defined in TS 119 442 -
// V1.1.1, ch. 5.1.4.2.1.3 XML component, building and caching it per DigestAlgorithm on first
// use. Port of getSignatureDigestReference(DigestAlgorithm).
func (s *DefaultAdvancedSignature) SignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.SignatureDigestReference {
	if s.signatureDigestReferences == nil {
		s.signatureDigestReferences = make(map[enumerations.DigestAlgorithm]*signature.SignatureDigestReference)
	}
	if reference, ok := s.signatureDigestReferences[digestAlgorithm]; ok {
		return reference
	}
	// Java's Map#computeIfAbsent does not store a null mapping function result; mirrored here
	// by only inserting into the map when a reference was actually built.
	reference := s.defaultAdvancedSignatureOverrides().BuildSignatureDigestReference(digestAlgorithm)
	if reference != nil {
		s.signatureDigestReferences[digestAlgorithm] = reference
	}
	return reference
}

// BaselineRequirementsChecker returns the cached instance of the BaselineRequirementsChecker.
// Port of the protected getBaselineRequirementsChecker().
//
// Panics with the Java message (IllegalStateException) when InitBaselineRequirementsChecker was
// not called first.
func (s *DefaultAdvancedSignature) BaselineRequirementsChecker() BaselineRequirementsCheckerContract {
	if s.baselineRequirementsChecker == nil {
		panic("Please call #initBaselineRequirementsChecker method before!")
	}
	return s.baselineRequirementsChecker
}

// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile.
// Port of hasAdESProfile().
func (s *DefaultAdvancedSignature) HasAdESProfile() bool {
	return s.BaselineRequirementsChecker().HasAdESProfile()
}

// HasBProfile checks if the signature is conformant to AdES-BASELINE-B level. Port of hasBProfile().
func (s *DefaultAdvancedSignature) HasBProfile() bool {
	return s.BaselineRequirementsChecker().HasBaselineBProfile()
}

// HasTProfile checks if the T-level is present in the signature. Port of hasTProfile().
func (s *DefaultAdvancedSignature) HasTProfile() bool {
	return s.BaselineRequirementsChecker().HasBaselineTProfile()
}

// HasLTProfile checks if the LT-level is present in the signature. Port of hasLTProfile().
func (s *DefaultAdvancedSignature) HasLTProfile() bool {
	return s.BaselineRequirementsChecker().HasBaselineLTProfile()
}

// HasLTAProfile checks if the LTA-level is present in the signature. Port of hasLTAProfile().
func (s *DefaultAdvancedSignature) HasLTAProfile() bool {
	return s.BaselineRequirementsChecker().HasBaselineLTAProfile()
}

// HasBESProfile checks the presence of signing certificate covered by the signature, what is the
// proof of the -BES profile existence. Port of hasBESProfile().
func (s *DefaultAdvancedSignature) HasBESProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedBESProfile()
}

// HasEPESProfile checks the presence of SignaturePolicyIdentifier element in the signature, what
// is the proof of the -EPES profile existence. Port of hasEPESProfile().
func (s *DefaultAdvancedSignature) HasEPESProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedEPESProfile()
}

// HasExtendedTProfile checks the presence of SignatureTimeStamp element in the signature, what is
// the proof of the -T profile existence. Port of hasExtendedTProfile().
func (s *DefaultAdvancedSignature) HasExtendedTProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedTProfile()
}

// HasCProfile checks the presence of CompleteCertificateRefs and CompleteRevocationRefs segments
// in the signature, what is the proof of the -C profile existence. Port of hasCProfile().
func (s *DefaultAdvancedSignature) HasCProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedCProfile()
}

// HasXProfile checks the presence of SigAndRefsTimeStamp segment in the signature, what is the
// proof of the -X profile existence. Port of hasXProfile().
func (s *DefaultAdvancedSignature) HasXProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedXProfile()
}

// HasXLProfile checks the presence of CertificateValues/RevocationValues segment in the
// signature, what is the proof of the -XL profile existence. Port of hasXLProfile().
func (s *DefaultAdvancedSignature) HasXLProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedXLProfile()
}

// HasAProfile checks the presence of ArchiveTimeStamp element in the signature, what is the
// proof of the -A profile existence. Port of hasAProfile().
func (s *DefaultAdvancedSignature) HasAProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedAProfile()
}

// HasERSProfile checks the presence of SealingEvidenceRecord element in the signature, what is
// the proof of the -ERS profile existence. Port of hasERSProfile().
func (s *DefaultAdvancedSignature) HasERSProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedERSProfile()
}

// AreAllSelfSignedCertificates checks if all certificate chains present in the signature are
// self-signed. Port of areAllSelfSignedCertificates().
func (s *DefaultAdvancedSignature) AreAllSelfSignedCertificates() bool {
	certificateSources := s.CompleteCertificateSource()
	certificateFound := certificateSources.NumberOfCertificates() > 0
	return certificateFound && certificateSources.IsAllSelfSigned()
}

// IsDocHashOnlyValidation returns true if the validation of the signature has been performed
// only on Signer's Document Representation (SDR). Port of isDocHashOnlyValidation().
func (s *DefaultAdvancedSignature) IsDocHashOnlyValidation() bool {
	if utils.IsCollectionNotEmpty(s.detachedContents) {
		for _, dssDocument := range s.detachedContents {
			if _, ok := dssDocument.(*model.DigestDocument); !ok {
				return false
			}
		}
		return true
	}
	return false
}

// IsHashOnlyValidation returns true if the validation of the signature has been performed only
// on Data To Be Signed Representation (DTBSR). Port of isHashOnlyValidation().
//
// TODO: not implemented yet (matches upstream).
func (s *DefaultAdvancedSignature) IsHashOnlyValidation() bool {
	return false
}

// Equals reports whether the two signatures share the same DSS Id. Port of equals(Object); Go
// has no Object type, so the parameter is narrowed to AdvancedSignature, and the "this == obj"/
// instanceof dance collapses to a nil check plus the identifier comparison.
func (s *DefaultAdvancedSignature) Equals(other AdvancedSignature) bool {
	if other == nil {
		return false
	}
	return s.DSSID().Equals(other.DSSID())
}

// String returns a string representation of the signature. Port of toString().
//
// Java also overrides hashCode() (consistent with equals(), delegating to getDSSId().hashCode());
// Go's comparable maps/sets in this port key on the identifier's digest string instead of a
// numeric hash code (see model.Identifier.String / AsXmlID), so no HashCode method is ported.
func (s *DefaultAdvancedSignature) String() string {
	return fmt.Sprintf("%s Signature with Id : %s", s.defaultAdvancedSignatureOverrides().SignatureForm(), s.ID())
}
