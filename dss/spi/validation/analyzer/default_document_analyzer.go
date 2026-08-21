// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/DefaultDocumentAnalyzer.java (DSS 6.5.RC1).
//
// This class contains a common code for processing of signed documents.
//
// # Virtual dispatch
//
// Java's DefaultDocumentAnalyzer is abstract, and format-specific subclasses (CAdES/XAdES/
// JAdES/PAdES/ASiC/EAA/evidence-record analyzers, ported in later phases) override several of
// its protected methods; the base class calls those methods on `this`, so Java's dynamic
// dispatch reaches the override. Go has no such call-site polymorphism for methods invoked from
// within an embedded base, so - following the AbstractSignatureIdentifierBuilder /
// AbstractSignatureEvidenceRecordDigestBuilder precedent already established in this package's
// sibling files - every protected method that (a) upstream leaves without a body (truly
// abstract) or (b) at least one upstream subclass (surveyed: PDFDocumentAnalyzer,
// AbstractJWSDocumentAnalyzer, DefaultEvidenceRecordAnalyzer, AbstractASiCContainerAnalyzer,
// DefaultEAAPresentationAnalyzer, XMLDocumentAnalyzer, COSEDocumentAnalyzer) actually overrides,
// is collected into DefaultDocumentAnalyzerOverrides; DefaultDocumentAnalyzer dispatches to it
// through the overrides field set by InitDefaultDocumentAnalyzer, which every concrete
// subclass's constructor must call once. Protected methods no upstream subclass overrides (e.g.
// appendCounterSignatures, processSignaturesValidation, validateSignaturePolicy) stay ordinary
// methods on the base; a Go embedder can still shadow a *top-level* DocumentAnalyzer interface
// method (e.g. OriginalDocuments) simply by defining its own method of that name, which is how
// PDFDocumentAnalyzer.getOriginalDocuments(String)-style full overrides are reproduced without
// needing an overrides-interface entry.
//
// getTimestampReaders() is `@Deprecated ... to be removed`, and no surveyed subclass overrides
// it; it is ported as a concrete (non-virtual) method returning an empty slice, matching the
// upstream default body, and is not part of DefaultDocumentAnalyzerOverrides.
//
// # Data-dependent throws
//
// DSSException("At least one signature or a timestamp shall be provided to extract the
// validation data!") in getValidationData is data-dependent on the caller-supplied signatures/
// detachedTimestamps, so it is returned as an error (per PORTING.md's "data-dependent throw ->
// (T, error)" rule) rather than raised as a panic.
//
// # Dropped / adapted members
//
// The static initializer block (`static { DSSSecurityProvider.initSystemProviders(); }`) calls
// through to a documented no-op (see spi.DSSSecurityProviderInitSystemProviders's own doc
// comment); it is still invoked, from this file's init(), for fidelity, though it currently has
// no observable effect.
//
// slf4j logging (the LOG.info/LOG.warn calls) is dropped per the phase 2a handoff fact ("slf4j
// dropped unless load-bearing").
//
// FORWARD DEPENDENCY (flagged per S2B_BRIEF.md): EvidenceRecordScopeFinder (Java
// spi.validation.scope.EvidenceRecordScopeFinder) is owned by a sibling chunk of phase 2b
// (package scope, dss/spi/validation/scope) not yet landed. It is referenced here by name only,
// with the shape inferred from every call this file makes to it (the Java source itself, out of
// this manifest's scope, was read for accuracy):
//
//	func scope.NewEvidenceRecordScopeFinder(evidenceRecord validation.EvidenceRecord) *scope.EvidenceRecordScopeFinder
//	func (*scope.EvidenceRecordScopeFinder) FindEvidenceRecordScope() []modelscope.SignatureScope
package analyzer

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/spi"
	spihttp "github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/spi/policy"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer/timestamp"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/spi/validation/scope"
	timestampsrc "github.com/ryftcore/dss-go/dss/spi/validation/timestamp"
	"github.com/ryftcore/dss-go/dss/utils"
)

func init() {
	// Port of the static { DSSSecurityProvider.initSystemProviders(); } block; see the file
	// header for why this currently has no observable effect.
	spi.DSSSecurityProviderInitSystemProviders()
}

// DefaultDocumentAnalyzerOverrides declares the operations DefaultDocumentAnalyzer calls back
// into virtually; every concrete analyzer (a later phase's CAdES/XAdES/JAdES/PAdES/ASiC/EAA/
// evidence-record analyzer) registers itself with InitDefaultDocumentAnalyzer so the base can
// dispatch them, the way model.TokenBase dispatches to model.TokenOverrides via InitToken. See
// the file header's "Virtual dispatch" section for which protected methods qualify.
type DefaultDocumentAnalyzerOverrides interface {
	// IsSupported checks if the document is supported by the current validator. Left abstract
	// by DocumentAnalyzer's Java interface with no body in DefaultDocumentAnalyzer.Port of
	// isSupported(DSSDocument).
	IsSupported(dssDocument model.DSSDocument) bool

	// OriginalDocumentsForSignature returns the signed document(s) without their signature(s).
	// Left abstract: DefaultDocumentAnalyzer.java declares no body for the
	// getOriginalDocuments(AdvancedSignature) overload, only for getOriginalDocuments(String)
	// (which calls this one). Port of the getOriginalDocuments(AdvancedSignature) overload.
	OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument

	// BuildSignatures builds a list of signatures to be extracted from a document. Default: not
	// implemented, returns an empty slice. Port of buildSignatures().
	BuildSignatures() []validation.AdvancedSignature

	// BuildDetachedTimestamps builds a list of detached TimestampTokens extracted from the
	// document. Default: empty slice. Port of buildDetachedTimestamps().
	BuildDetachedTimestamps() []*validation.TimestampToken

	// BuildDetachedEvidenceRecords builds a list of detached EvidenceRecords extracted from the
	// document. Default: from the detachedEvidenceRecordDocuments field (see
	// DefaultDocumentAnalyzer's own buildDetachedEvidenceRecords()). Port of
	// buildDetachedEvidenceRecords().
	BuildDetachedEvidenceRecords() []validation.EvidenceRecord

	// GetDefaultSignaturePolicyValidator returns a signature format dependent
	// SignaturePolicyValidator, based on the specification. Default: nil. Port of
	// getDefaultSignaturePolicyValidator().
	GetDefaultSignaturePolicyValidator() policy.SignaturePolicyValidator

	// CoversSignature verifies whether evidenceRecord covers signature. Default: true. Port of
	// coversSignature(AdvancedSignature, EvidenceRecord).
	CoversSignature(signature validation.AdvancedSignature, evidenceRecord validation.EvidenceRecord) bool

	// AddReference checks if the signature scope shall be added as a timestamped reference.
	// NOTE: used to avoid duplicates in ASiC with CAdES validator, due to covered signature/
	// timestamp files. Default: true. Port of addReference(SignatureScope).
	AddReference(signatureScope modelscope.SignatureScope) bool

	// IsTimestampCoveredByEvidenceRecord checks whether timestampToken is covered by
	// evidenceRecord. Default: true. Port of
	// isTimestampCoveredByEvidenceRecord(TimestampToken, EvidenceRecord).
	IsTimestampCoveredByEvidenceRecord(timestampToken *validation.TimestampToken, evidenceRecord validation.EvidenceRecord) bool

	// GetAllSignatures returns a list of all signatures from the validating document. Default:
	// Signatures() plus their counter signatures, with external evidence records appended. Port
	// of getAllSignatures().
	GetAllSignatures() []validation.AdvancedSignature

	// PrepareValidationContext initializes and fills a ValidationContext with the necessary
	// data sources. Port of prepareValidationContext(Collection, Collection, Collection,
	// CertificateVerifier).
	PrepareValidationContext(signatures []validation.AdvancedSignature, detachedTimestamps []*validation.TimestampToken,
		detachedEvidenceRecords []validation.EvidenceRecord, certificateVerifier validation.CertificateVerifier) validation.ValidationContext

	// CreateValidationContext creates a new instance of ValidationContext performing
	// preparation of validation data, certificate chain building, revocation request, as well as
	// custom validation checks execution. Default: a SignatureValidationContext at
	// ValidationTime(). Port of createValidationContext().
	CreateValidationContext() validation.ValidationContext

	// InstantiateValidationDataContainer creates a new instance of ValidationDataContainer.
	// Port of instantiateValidationDataContainer().
	InstantiateValidationDataContainer() *validation.ValidationDataContainer

	// AppendExternalEvidenceRecords appends the detached evidence record provided to the
	// validator to the corresponding signatures covered by the evidence record document.
	// Default: see DefaultDocumentAnalyzer.AppendExternalEvidenceRecords. Port of
	// appendExternalEvidenceRecords(List). ADDITIVE (see that method's doc comment): promoted
	// into this interface in phase 3 (S3_BRIEF.md) once dss-cades.CMSDocumentAnalyzer surfaced
	// a real override this port had not yet accounted for.
	AppendExternalEvidenceRecords(allSignatureList []validation.AdvancedSignature) []validation.AdvancedSignature
}

// DefaultDocumentAnalyzer contains a common code for processing of signed documents. It is
// embedded by every concrete format-specific analyzer, which must call
// InitDefaultDocumentAnalyzer once (typically from its own constructor).
type DefaultDocumentAnalyzer struct {
	// overrides points back at the concrete analyzer; see InitDefaultDocumentAnalyzer.
	overrides DefaultDocumentAnalyzerOverrides

	// document is the document to be validated (with the signature(s) or timestamp(s)). Java
	// declares this protected; a Go subclass in another package reaches it through Document()/
	// SetDocument().
	document model.DSSDocument

	// detachedContents contains the signed documents, in case of a detached signature.
	detachedContents []model.DSSDocument

	// detachedEvidenceRecordDocuments contains a list of evidence record documents detached
	// from the signature.
	detachedEvidenceRecordDocuments []model.DSSDocument

	// containerContents is the list of container documents, in case of an ASiC signature.
	containerContents []model.DSSDocument

	// manifestFile is a related ManifestFile to the provided document.
	manifestFile *model.ManifestFile

	// signingCertificateSource finds the signing certificate.
	signingCertificateSource spi.CertificateSource

	// validationTime is a time to validate the document against.
	validationTime *time.Time

	// certificateVerifier is the reference to the certificate verifier. The current DSS
	// implementation proposes validation.CommonCertificateVerifier. This verifier encapsulates
	// the references to different sources used in the signature validation process.
	certificateVerifier validation.CertificateVerifier

	// validationContextExecutor performs validation of ValidationContext.
	// Default: executor.DefaultValidationContextExecutorInstance.
	validationContextExecutor executor.ValidationContextExecutor

	// tokenIdentifierProvider is the implementation to be used for identifiers generation.
	tokenIdentifierProvider model.TokenIdentifierProvider

	// signaturePolicyProvider provides methods to extract a policy content by its identifier.
	signaturePolicyProvider *policy.SignaturePolicyProvider

	// signaturePolicyValidatorLoader loads a SignaturePolicyValidator to perform a signature
	// policy validation.
	signaturePolicyValidatorLoader policy.SignaturePolicyValidatorLoader

	// signatures caches the list of signatures extracted from the document.
	signatures []validation.AdvancedSignature
	// signaturesSet reports whether signatures has been computed (Go has no null slice/
	// "uncomputed" distinction as clean as Java's null field, since an empty slice is a valid
	// cached result).
	signaturesSet bool

	// detachedTimestamps caches the list of detached timestamps extracted from the document.
	detachedTimestamps    []*validation.TimestampToken
	detachedTimestampsSet bool

	// evidenceRecords caches the list of detached evidence records extracted from the document.
	evidenceRecords    []validation.EvidenceRecord
	evidenceRecordsSet bool
}

// NewDefaultDocumentAnalyzerBase builds the base state a subclass embeds. Port of the protected
// default DefaultDocumentAnalyzer() constructor; the subclass constructor must follow it with
// InitDefaultDocumentAnalyzer.
func NewDefaultDocumentAnalyzerBase() DefaultDocumentAnalyzer {
	return DefaultDocumentAnalyzer{
		detachedContents:                []model.DSSDocument{},
		detachedEvidenceRecordDocuments: []model.DSSDocument{},
		validationContextExecutor:       executor.DefaultValidationContextExecutorInstance,
		tokenIdentifierProvider:         model.NewOriginalIdentifierProvider(),
	}
}

// InitDefaultDocumentAnalyzer registers the concrete analyzer with its base so that the base can
// dispatch DefaultDocumentAnalyzerOverrides. Every concrete subclass constructor must call this
// once.
func (a *DefaultDocumentAnalyzer) InitDefaultDocumentAnalyzer(overrides DefaultDocumentAnalyzerOverrides) {
	a.overrides = overrides
}

// defaultDocumentAnalyzerOverrides returns the registered overrides, panicking when the
// concrete analyzer forgot to call InitDefaultDocumentAnalyzer - Java's abstract class can never
// be instantiated bare.
func (a *DefaultDocumentAnalyzer) defaultDocumentAnalyzerOverrides() DefaultDocumentAnalyzerOverrides {
	if a.overrides == nil {
		panic("DefaultDocumentAnalyzer was not initialised: the concrete analyzer must call InitDefaultDocumentAnalyzer in its constructor")
	}
	return a.overrides
}

// DocumentAnalyzerFromDocument guesses the document format and returns an appropriate document
// reader, by finding a registered DocumentAnalyzerFactory that supports dssDocument. Port of the
// static fromDocument(DSSDocument).
//
// Panics when dssDocument is nil (Objects.requireNonNull("DSSDocument is null")). Java's
// UnsupportedOperationException("Document format not recognized/handled"), thrown when no
// registered factory supports the document, is returned as an error: whether a document format
// is recognized is data-dependent on which analyzer factories happen to be registered (loaded
// modules).
func DocumentAnalyzerFromDocument(dssDocument model.DSSDocument) (DocumentAnalyzer, error) {
	if dssDocument == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range documentAnalyzerFactoryRegistry {
		if factory.IsSupported(dssDocument) {
			return factory.Create(dssDocument), nil
		}
	}
	return nil, errDocumentFormatNotRecognized
}

// Document gets the document to be validated. Port of getDocument().
//
// Panics when no document has been provided (Java's IllegalStateException("Document is not
// provided! Please use a different constructor to extract the document.")).
func (a *DefaultDocumentAnalyzer) Document() model.DSSDocument {
	if a.document == nil {
		panic("Document is not provided! Please use a different constructor to extract the document.")
	}
	return a.document
}

// HasDocument reports whether a document has been provided, without Document's panic. Go
// counterpart of a subclass testing its protected `document` field directly against null - e.g.
// CMSDocumentAnalyzer.buildSignatures()'s own "if (document != null)" in
// dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentAnalyzer.java, which
// reads the field rather than calling the accessor precisely to avoid getDocument()'s throw.
func (a *DefaultDocumentAnalyzer) HasDocument() bool {
	return a.document != nil
}

// SetDocument sets the document to be validated. Go counterpart of assigning the protected
// `document` field directly, as concrete subclass constructors do in Java.
func (a *DefaultDocumentAnalyzer) SetDocument(document model.DSSDocument) {
	a.document = document
}

// SetSigningCertificateSource sets a certificate source which allows finding the signing
// certificate by kid or certificate's digest. Port of
// setSigningCertificateSource(CertificateSource).
func (a *DefaultDocumentAnalyzer) SetSigningCertificateSource(signingCertificateSource spi.CertificateSource) {
	a.signingCertificateSource = signingCertificateSource
}

// SetCertificateVerifier passes the CertificateVerifier used to carry out the validation
// process of the signature(s): some external sources of certificates and of revocation data can
// be needed. Note that once this setter is called any change in the content of the
// CommonTrustedCertificateSource or in adjunct certificate source is not taken into account.
// Port of setCertificateVerifier(CertificateVerifier).
//
// Panics when certificateVerifier is nil (Objects.requireNonNull).
func (a *DefaultDocumentAnalyzer) SetCertificateVerifier(certificateVerifier validation.CertificateVerifier) {
	if certificateVerifier == nil {
		panic("certificateVerifier cannot be null!")
	}
	a.certificateVerifier = certificateVerifier
}

// SetValidationContextExecutor sets the ValidationContextExecutor. Port of
// setValidationContextExecutor(ValidationContextExecutor).
func (a *DefaultDocumentAnalyzer) SetValidationContextExecutor(validationContextExecutor executor.ValidationContextExecutor) {
	a.validationContextExecutor = validationContextExecutor
}

// TokenIdentifierProvider gets the TokenIdentifierProvider. Port of getTokenIdentifierProvider().
func (a *DefaultDocumentAnalyzer) TokenIdentifierProvider() model.TokenIdentifierProvider {
	return a.tokenIdentifierProvider
}

// SetTokenIdentifierProvider sets the TokenIdentifierProvider. Port of
// setTokenIdentifierProvider(TokenIdentifierProvider).
//
// Panics when tokenIdentifierProvider is nil (Objects.requireNonNull).
func (a *DefaultDocumentAnalyzer) SetTokenIdentifierProvider(tokenIdentifierProvider model.TokenIdentifierProvider) {
	if tokenIdentifierProvider == nil {
		panic("tokenIdentifierProvider cannot be null!")
	}
	a.tokenIdentifierProvider = tokenIdentifierProvider
}

// SetDetachedContents sets the list of DSSDocument containing the original contents to sign,
// for detached signature scenarios. Port of setDetachedContents(List).
func (a *DefaultDocumentAnalyzer) SetDetachedContents(detachedContents []model.DSSDocument) {
	a.detachedContents = detachedContents
}

// DetachedContents returns the signed documents, in case of a detached signature. Exported
// accessor for the protected `detachedContents` field: Java lets a subclass in another package
// read a protected field through inheritance, which a Go subclass in another package reaches
// only through a getter (ADDITIVE, S3_BRIEF.md porter, phase 3 - see
// AppendExternalEvidenceRecords's doc comment on why dss-cades surfaces gaps like this one).
func (a *DefaultDocumentAnalyzer) DetachedContents() []model.DSSDocument {
	return a.detachedContents
}

// ContainerContents returns the list of container documents, in case of an ASiC signature.
// ADDITIVE accessor; see DetachedContents.
func (a *DefaultDocumentAnalyzer) ContainerContents() []model.DSSDocument {
	return a.containerContents
}

// ManifestFile returns the related ManifestFile to the provided document. ADDITIVE accessor;
// see DetachedContents.
func (a *DefaultDocumentAnalyzer) ManifestFile() *model.ManifestFile {
	return a.manifestFile
}

// SigningCertificateSource returns the certificate source that finds the signing certificate.
// ADDITIVE accessor; see DetachedContents.
func (a *DefaultDocumentAnalyzer) SigningCertificateSource() spi.CertificateSource {
	return a.signingCertificateSource
}

// CertificateVerifier returns the reference to the certificate verifier. ADDITIVE accessor; see
// DetachedContents.
func (a *DefaultDocumentAnalyzer) CertificateVerifier() validation.CertificateVerifier {
	return a.certificateVerifier
}

// ValidateSignaturePolicy performs validation of the signature policy's identifier, when
// present. Exported so a subclass in another package can call it as CMSDocumentAnalyzer's
// buildSignatures() does (Java: protected inherited method). ADDITIVE, see DetachedContents;
// forwards to the private validateSignaturePolicy this file already defines.
func (a *DefaultDocumentAnalyzer) ValidateSignaturePolicy(sig validation.AdvancedSignature) {
	a.validateSignaturePolicy(sig)
}

// SetDetachedEvidenceRecordDocuments sets a list of DSSDocument containing the evidence record
// documents covering the signature document. Port of
// setDetachedEvidenceRecordDocuments(List).
func (a *DefaultDocumentAnalyzer) SetDetachedEvidenceRecordDocuments(detachedEvidenceRecordDocuments []model.DSSDocument) {
	a.detachedEvidenceRecordDocuments = detachedEvidenceRecordDocuments
}

// SetContainerContents sets the list of DSSDocument containing the original container content
// for ASiC-S signatures. Port of setContainerContents(List).
func (a *DefaultDocumentAnalyzer) SetContainerContents(containerContents []model.DSSDocument) {
	a.containerContents = containerContents
}

// SetManifestFile sets a related ManifestFile to the document to be validated. Port of
// setManifestFile(ManifestFile).
func (a *DefaultDocumentAnalyzer) SetManifestFile(manifestFile *model.ManifestFile) {
	a.manifestFile = manifestFile
}

// ValidationTime returns validation time. In case the validation time is not provided,
// initializes the current time value from the system. Port of getValidationTime().
func (a *DefaultDocumentAnalyzer) ValidationTime() time.Time {
	if a.validationTime == nil {
		now := time.Now()
		a.validationTime = &now
	}
	return *a.validationTime
}

// SetValidationTime allows defining a custom validation time. Port of setValidationTime(Date).
func (a *DefaultDocumentAnalyzer) SetValidationTime(validationTime time.Time) {
	a.validationTime = &validationTime
}

// SetDetachedEvidenceRecords sets a list of detached evidence records. Port of
// setDetachedEvidenceRecords(List).
func (a *DefaultDocumentAnalyzer) SetDetachedEvidenceRecords(evidenceRecords []validation.EvidenceRecord) {
	a.evidenceRecords = evidenceRecords
	a.evidenceRecordsSet = true
}

// signaturePolicyProviderOrDefault returns signaturePolicyProvider; if not defined, returns a
// default provider. Port of the protected getSignaturePolicyProvider().
func (a *DefaultDocumentAnalyzer) signaturePolicyProviderOrDefault() *policy.SignaturePolicyProvider {
	if a.signaturePolicyProvider == nil {
		// LOG.info("Default SignaturePolicyProvider instantiated with NativeHTTPDataLoader.") dropped: not load-bearing.
		a.signaturePolicyProvider = policy.NewSignaturePolicyProvider()
		a.signaturePolicyProvider.SetDataLoader(spihttp.NewNativeHTTPDataLoader())
	}
	return a.signaturePolicyProvider
}

// SetSignaturePolicyProvider allows setting a provider for Signature policies. Port of
// setSignaturePolicyProvider(SignaturePolicyProvider).
func (a *DefaultDocumentAnalyzer) SetSignaturePolicyProvider(signaturePolicyProvider *policy.SignaturePolicyProvider) {
	a.signaturePolicyProvider = signaturePolicyProvider
}

// SetSignaturePolicyValidatorLoader sets a loader for a SignaturePolicyValidator. Port of
// setSignaturePolicyValidatorLoader(SignaturePolicyValidatorLoader).
func (a *DefaultDocumentAnalyzer) SetSignaturePolicyValidatorLoader(signaturePolicyValidatorLoader policy.SignaturePolicyValidatorLoader) {
	a.signaturePolicyValidatorLoader = signaturePolicyValidatorLoader
}

// Validate performs validation of the document. Port of validate().
//
// Panics when certificateVerifier or document is missing (Objects.requireNonNull).
func (a *DefaultDocumentAnalyzer) Validate() validation.ValidationContext {
	if a.certificateVerifier == nil {
		panic("CertificateVerifier is not defined")
	}
	if a.document == nil {
		panic("Document is not provided to the validator")
	}

	overrides := a.defaultDocumentAnalyzerOverrides()

	allSignatures := overrides.GetAllSignatures()
	allDetachedTimestamps := a.DetachedTimestamps()
	allDetachedEvidenceRecords := a.DetachedEvidenceRecords()

	certificateVerifierForValidation := validation.NewCertificateVerifierBuilder(a.certificateVerifier).BuildCompleteCopyForValidation()
	validationContext := overrides.PrepareValidationContext(
		allSignatures, allDetachedTimestamps, allDetachedEvidenceRecords, certificateVerifierForValidation)
	a.validateContext(validationContext)
	return validationContext
}

// prepareValidationContext initializes and fills a ValidationContext with the necessary data
// sources. This is the default body dispatched to by
// DefaultDocumentAnalyzerOverrides.PrepareValidationContext for concrete analyzers that do not
// override it. Port of prepareValidationContext(Collection, Collection, Collection,
// CertificateVerifier).
func (a *DefaultDocumentAnalyzer) prepareValidationContext(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken, detachedEvidenceRecords []validation.EvidenceRecord,
	certificateVerifier validation.CertificateVerifier) validation.ValidationContext {
	overrides := a.defaultDocumentAnalyzerOverrides()
	validationContext := overrides.CreateValidationContext()
	validationContext.Initialize(certificateVerifier)
	a.prepareSignatureValidationContext(validationContext, signatures)
	a.prepareDetachedTimestampValidationContext(validationContext, detachedTimestamps)
	a.prepareDetachedEvidenceRecordValidationContext(validationContext, detachedEvidenceRecords)
	return validationContext
}

// PrepareValidationContext is DefaultDocumentAnalyzerOverrides' default body: see
// prepareValidationContext.
func (a *DefaultDocumentAnalyzer) PrepareValidationContext(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken, detachedEvidenceRecords []validation.EvidenceRecord,
	certificateVerifier validation.CertificateVerifier) validation.ValidationContext {
	return a.prepareValidationContext(signatures, detachedTimestamps, detachedEvidenceRecords, certificateVerifier)
}

// CreateValidationContext is DefaultDocumentAnalyzerOverrides' default body. Port of
// createValidationContext().
func (a *DefaultDocumentAnalyzer) CreateValidationContext() validation.ValidationContext {
	return validation.NewSignatureValidationContextAtTime(a.ValidationTime())
}

// GetValidationData extracts a validation data for the provided collection of signatures. Port
// of the getValidationData(Collection) overload.
func (a *DefaultDocumentAnalyzer) GetValidationData(signatures []validation.AdvancedSignature) (*validation.ValidationDataContainer, error) {
	return a.GetValidationDataWithTimestamps(signatures, nil)
}

// GetValidationDataWithTimestamps extracts a validation data for the provided collection of
// signatures and/or timestamps. Port of the getValidationData(Collection, Collection) overload.
//
// Java's DSSException("At least one signature or a timestamp shall be provided to extract the
// validation data!") is data-dependent on the caller-supplied signatures/detachedTimestamps, so
// it is returned as an error.
func (a *DefaultDocumentAnalyzer) GetValidationDataWithTimestamps(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken) (*validation.ValidationDataContainer, error) {
	if utils.IsCollectionEmpty(signatures) && utils.IsCollectionEmpty(detachedTimestamps) {
		return nil, model.NewDSSError("At least one signature or a timestamp shall be provided to extract the validation data!")
	}

	overrides := a.defaultDocumentAnalyzerOverrides()

	// TODO : review use of evidence records
	validationContext := overrides.PrepareValidationContext(signatures, detachedTimestamps, nil, a.certificateVerifier)
	a.validateContext(validationContext)

	validationDataContainer := overrides.InstantiateValidationDataContainer()
	for _, sig := range signatures {
		signatureValidationData := validationContext.GetValidationData(sig)
		validationDataContainer.AddValidationDataForSignature(sig, signatureValidationData)
		for _, timestampToken := range sig.AllTimestamps() {
			timestampValidationData := validationContext.GetValidationDataForTimestamp(timestampToken)
			validationDataContainer.AddValidationDataForTimestamp(timestampToken, timestampValidationData)
		}
		for _, counterSignature := range sig.CounterSignatures() {
			counterSignatureValidationData := validationContext.GetValidationData(counterSignature)
			validationDataContainer.AddValidationDataForSignature(counterSignature, counterSignatureValidationData)
			for _, timestampToken := range counterSignature.AllTimestamps() {
				timestampValidationData := validationContext.GetValidationDataForTimestamp(timestampToken)
				validationDataContainer.AddValidationDataForTimestamp(timestampToken, timestampValidationData)
			}
		}
	}
	for _, detachedTimestamp := range detachedTimestamps {
		timestampValidationData := validationContext.GetValidationDataForTimestamp(detachedTimestamp)
		validationDataContainer.AddValidationDataForTimestamp(detachedTimestamp, timestampValidationData)
	}

	return validationDataContainer, nil
}

// InstantiateValidationDataContainer is DefaultDocumentAnalyzerOverrides' default body. Port of
// instantiateValidationDataContainer().
func (a *DefaultDocumentAnalyzer) InstantiateValidationDataContainer() *validation.ValidationDataContainer {
	return validation.NewValidationDataContainer()
}

// getAllEvidenceRecords returns a list of all found evidence records (embedded and detached).
// Port of getAllEvidenceRecords(List, List). Currently unused by any code path in this file
// (upstream keeps it as a protected helper for subclasses; it is preserved 1:1 for API parity).
func (a *DefaultDocumentAnalyzer) getAllEvidenceRecords(signatures []validation.AdvancedSignature,
	detachedEvidenceRecords []validation.EvidenceRecord) []validation.EvidenceRecord {
	result := make([]validation.EvidenceRecord, 0, len(detachedEvidenceRecords))
	for _, sig := range signatures {
		result = append(result, sig.EmbeddedEvidenceRecords()...)
	}
	result = append(result, detachedEvidenceRecords...)
	return result
}

// prepareSignatureValidationContext prepares validationContext for the signature validation
// process. Port of prepareSignatureValidationContext(ValidationContext, Collection).
func (a *DefaultDocumentAnalyzer) prepareSignatureValidationContext(validationContext validation.ValidationContext,
	allSignatures []validation.AdvancedSignature) {
	a.prepareSignatureForVerification(validationContext, allSignatures)
	a.processSignaturesValidation(allSignatures)
}

// prepareSignatureForVerification prepares a SignatureValidationContext for signatures
// validation. Port of prepareSignatureForVerification(ValidationContext, Collection).
func (a *DefaultDocumentAnalyzer) prepareSignatureForVerification(validationContext validation.ValidationContext,
	allSignatureList []validation.AdvancedSignature) {
	for _, sig := range allSignatureList {
		validationContext.AddSignatureForVerification(sig)
	}
}

// prepareDetachedTimestampValidationContext prepares validationContext for a timestamp
// validation process. Port of prepareDetachedTimestampValidationContext(ValidationContext,
// Collection).
func (a *DefaultDocumentAnalyzer) prepareDetachedTimestampValidationContext(validationContext validation.ValidationContext,
	timestamps []*validation.TimestampToken) {
	for _, timestampToken := range timestamps {
		validationContext.AddTimestampTokenForVerification(timestampToken)
	}
}

// prepareDetachedEvidenceRecordValidationContext prepares validationContext for the evidence
// record validation process. Port of prepareDetachedEvidenceRecordValidationContext(
// ValidationContext, Collection).
func (a *DefaultDocumentAnalyzer) prepareDetachedEvidenceRecordValidationContext(validationContext validation.ValidationContext,
	evidenceRecords []validation.EvidenceRecord) {
	for _, evidenceRecord := range evidenceRecords {
		validationContext.AddEvidenceRecordForVerification(evidenceRecord)
	}
}

// validateContext processes the validation. Port of validateContext(ValidationContext).
func (a *DefaultDocumentAnalyzer) validateContext(validationContext validation.ValidationContext) {
	a.validationContextExecutor.Validate(validationContext)
}

// SignaturePolicyValidatorLoader returns an instance of a corresponding SignaturePolicyValidatorLoader.
// Port of getSignaturePolicyValidatorLoader().
func (a *DefaultDocumentAnalyzer) SignaturePolicyValidatorLoader() policy.SignaturePolicyValidatorLoader {
	if a.signaturePolicyValidatorLoader == nil {
		overrides := a.defaultDocumentAnalyzerOverrides()
		a.signaturePolicyValidatorLoader = policy.DefaultSignaturePolicyValidatorLoaderDefaultUnlessSpecified(
			overrides.GetDefaultSignaturePolicyValidator())
	}
	return a.signaturePolicyValidatorLoader
}

// GetDefaultSignaturePolicyValidator is DefaultDocumentAnalyzerOverrides' default body. Port of
// getDefaultSignaturePolicyValidator().
func (a *DefaultDocumentAnalyzer) GetDefaultSignaturePolicyValidator() policy.SignaturePolicyValidator {
	return nil
}

// GetAllSignatures is DefaultDocumentAnalyzerOverrides' default body. Port of
// getAllSignatures().
func (a *DefaultDocumentAnalyzer) GetAllSignatures() []validation.AdvancedSignature {
	allSignatureList := make([]validation.AdvancedSignature, 0)
	for _, sig := range a.Signatures() {
		allSignatureList = append(allSignatureList, sig)
		allSignatureList = a.appendCounterSignatures(allSignatureList, sig)
	}
	allSignatureList = a.defaultDocumentAnalyzerOverrides().AppendExternalEvidenceRecords(allSignatureList)
	return allSignatureList
}

// appendCounterSignatures links counter signatures with the related master signatures. Port of
// appendCounterSignatures(List, AdvancedSignature).
func (a *DefaultDocumentAnalyzer) appendCounterSignatures(allSignatureList []validation.AdvancedSignature,
	sig validation.AdvancedSignature) []validation.AdvancedSignature {
	for _, counterSignature := range sig.CounterSignatures() {
		counterSignature.InitBaselineRequirementsChecker(a.certificateVerifier)
		a.validateSignaturePolicy(counterSignature)
		allSignatureList = append(allSignatureList, counterSignature)

		allSignatureList = a.appendCounterSignatures(allSignatureList, counterSignature)
	}
	return allSignatureList
}

// AppendExternalEvidenceRecords is DefaultDocumentAnalyzerOverrides' default body: appends the
// detached evidence record provided to the validator to the corresponding signatures covered by
// the evidence record document. Port of appendExternalEvidenceRecords(List).
//
// ADDITIVE FIX (S3_BRIEF.md porter, phase 3): this method was originally ported as a concrete,
// non-virtual method (unexported appendExternalEvidenceRecords), on the survey in this file's
// header of upstream subclasses known to override protected DefaultDocumentAnalyzer methods at
// the time of that port - a survey that could not yet include dss-cades's
// CMSDocumentAnalyzer.appendExternalEvidenceRecords(List), a real override phase 3 needs to
// reproduce. Promoted into DefaultDocumentAnalyzerOverrides/GetAllSignatures's dispatch so a
// concrete analyzer registered via InitDefaultDocumentAnalyzer can override it; every other
// analyzer's behaviour is unchanged since this is exactly its former body.
func (a *DefaultDocumentAnalyzer) AppendExternalEvidenceRecords(allSignatureList []validation.AdvancedSignature) []validation.AdvancedSignature {
	overrides := a.defaultDocumentAnalyzerOverrides()
	detachedEvidenceRecords := a.DetachedEvidenceRecords()
	if utils.IsCollectionNotEmpty(detachedEvidenceRecords) && utils.IsCollectionNotEmpty(allSignatureList) {
		for _, sig := range allSignatureList {
			for _, evidenceRecord := range detachedEvidenceRecords {
				if overrides.CoversSignature(sig, evidenceRecord) {
					sig.AddExternalEvidenceRecord(evidenceRecord)
				}
			}
		}
	}
	return allSignatureList
}

// appendExternalEvidenceRecordsToTimestamp appends the detached evidence records covering the
// time-stamp. Port of appendExternalEvidenceRecords(TimestampToken).
func (a *DefaultDocumentAnalyzer) appendExternalEvidenceRecordsToTimestamp(timestampToken *validation.TimestampToken) {
	overrides := a.defaultDocumentAnalyzerOverrides()
	detachedTimestampSource := timestampsrc.NewDetachedTimestampSourceWithTimestamp(timestampToken)
	for _, evidenceRecord := range a.DetachedEvidenceRecords() {
		if overrides.IsTimestampCoveredByEvidenceRecord(timestampToken, evidenceRecord) {
			timestampToken.AddDetachedEvidenceRecord(evidenceRecord)
			_ = detachedTimestampSource.AddExternalEvidenceRecord(evidenceRecord)
		}
	}
}

// IsTimestampCoveredByEvidenceRecord is DefaultDocumentAnalyzerOverrides' default body: true.
// Port of isTimestampCoveredByEvidenceRecord(TimestampToken, EvidenceRecord).
func (a *DefaultDocumentAnalyzer) IsTimestampCoveredByEvidenceRecord(timestampToken *validation.TimestampToken,
	evidenceRecord validation.EvidenceRecord) bool {
	return true
}

// CoversSignature is DefaultDocumentAnalyzerOverrides' default body: true. Port of
// coversSignature(AdvancedSignature, EvidenceRecord).
func (a *DefaultDocumentAnalyzer) CoversSignature(signature validation.AdvancedSignature, evidenceRecord validation.EvidenceRecord) bool {
	return true
}

// Signatures returns a list of all signatures from the validating document, building and
// caching them on first access. Port of getSignatures().
func (a *DefaultDocumentAnalyzer) Signatures() []validation.AdvancedSignature {
	if !a.signaturesSet {
		a.signatures = a.defaultDocumentAnalyzerOverrides().BuildSignatures()
		a.signaturesSet = true
	}
	return a.signatures
}

// BuildSignatures is DefaultDocumentAnalyzerOverrides' default body: not implemented, returns
// an empty slice. Port of buildSignatures().
func (a *DefaultDocumentAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	return nil
}

// DetachedTimestamps returns a list of the detached timestamps found in the document, building
// and caching them on first access. Port of getDetachedTimestamps().
func (a *DefaultDocumentAnalyzer) DetachedTimestamps() []*validation.TimestampToken {
	if !a.detachedTimestampsSet {
		a.detachedTimestamps = a.defaultDocumentAnalyzerOverrides().BuildDetachedTimestamps()
		a.detachedTimestampsSet = true
	}
	return a.detachedTimestamps
}

// BuildDetachedTimestamps is DefaultDocumentAnalyzerOverrides' default body: empty. Port of
// buildDetachedTimestamps().
func (a *DefaultDocumentAnalyzer) BuildDetachedTimestamps() []*validation.TimestampToken {
	return nil
}

// TimestampReaders returns a list of timestamp validators for timestamps embedded into the
// container.
//
// Deprecated: since DSS 6.5. To be removed. Port of the deprecated protected
// getTimestampReaders(); unlike the methods above, no upstream subclass overrides this one, so
// it is a concrete (non-virtual) method rather than part of DefaultDocumentAnalyzerOverrides.
func (a *DefaultDocumentAnalyzer) TimestampReaders() []timestamp.TimestampAnalyzer {
	return nil
}

// DetachedEvidenceRecords returns a list of the detached evidence records found in the
// document, building and caching them on first access. Port of getDetachedEvidenceRecords().
func (a *DefaultDocumentAnalyzer) DetachedEvidenceRecords() []validation.EvidenceRecord {
	if !a.evidenceRecordsSet {
		a.evidenceRecords = a.defaultDocumentAnalyzerOverrides().BuildDetachedEvidenceRecords()
		a.evidenceRecordsSet = true
	}
	return a.evidenceRecords
}

// BuildDetachedEvidenceRecords is DefaultDocumentAnalyzerOverrides' default body: builds
// evidence records from detachedEvidenceRecordDocuments. Port of
// buildDetachedEvidenceRecords().
func (a *DefaultDocumentAnalyzer) BuildDetachedEvidenceRecords() []validation.EvidenceRecord {
	if utils.IsCollectionNotEmpty(a.detachedEvidenceRecordDocuments) {
		result := make([]validation.EvidenceRecord, 0, len(a.detachedEvidenceRecordDocuments))
		for _, evidenceRecordDocument := range a.detachedEvidenceRecordDocuments {
			if evidenceRecord := a.buildEvidenceRecord(evidenceRecordDocument); evidenceRecord != nil {
				result = append(result, evidenceRecord)
			}
		}
		return result
	}
	return nil
}

// buildEvidenceRecord builds an evidence record from the given DSSDocument. Port of
// buildEvidenceRecord(DSSDocument).
//
// Java catches (and logs) UnsupportedOperationException and any other Exception, returning null
// in either case; the Go port reproduces that by simply returning the two error results as nil,
// since the logging is dropped per PORTING.md and both error paths already produce the "not
// found" outcome upstream converges on.
func (a *DefaultDocumentAnalyzer) buildEvidenceRecord(evidenceRecordDocument model.DSSDocument) validation.EvidenceRecord {
	evidenceRecordAnalyzer, err := EvidenceRecordAnalyzerFromDocument(evidenceRecordDocument)
	if err != nil {
		return nil
	}
	evidenceRecordAnalyzer.SetDetachedContents(a.signatureEvidenceRecordDetachedContents())
	evidenceRecordAnalyzer.SetCertificateVerifier(a.certificateVerifier)
	return a.getEvidenceRecord(evidenceRecordAnalyzer)
}

// signatureEvidenceRecordDetachedContents ports the private getSignatureEvidenceRecordDetachedContents().
func (a *DefaultDocumentAnalyzer) signatureEvidenceRecordDetachedContents() []model.DSSDocument {
	erDetachedContents := []model.DSSDocument{a.document}
	if utils.IsCollectionNotEmpty(a.detachedContents) {
		erDetachedContents = append(erDetachedContents, a.detachedContents...)
	}
	return erDetachedContents
}

// getEvidenceRecord gets an evidence record from evidenceRecordAnalyzer. Port of
// getEvidenceRecord(EvidenceRecordAnalyzer).
func (a *DefaultDocumentAnalyzer) getEvidenceRecord(evidenceRecordAnalyzer EvidenceRecordAnalyzer) validation.EvidenceRecord {
	evidenceRecord := evidenceRecordAnalyzer.EvidenceRecord()
	if evidenceRecord != nil {
		evidenceRecordScopes := a.getEvidenceRecordScopes(evidenceRecord)
		evidenceRecord.SetEvidenceRecordScopes(evidenceRecordScopes)
		evidenceRecord.SetTimestampedReferences(a.getTimestampedReferences(evidenceRecordScopes))
		return evidenceRecord
	}
	return nil
}

// getEvidenceRecordScopes finds evidence record scopes. Port of
// getEvidenceRecordScopes(EvidenceRecord). See the file header's FORWARD DEPENDENCY note on
// scope.EvidenceRecordScopeFinder.
func (a *DefaultDocumentAnalyzer) getEvidenceRecordScopes(evidenceRecord validation.EvidenceRecord) []modelscope.SignatureScope {
	return scope.NewEvidenceRecordScopeFinder(evidenceRecord).FindEvidenceRecordScope()
}

// processSignaturesValidation performs cryptographic validation of the signatures. Port of
// processSignaturesValidation(Collection).
func (a *DefaultDocumentAnalyzer) processSignaturesValidation(allSignatureList []validation.AdvancedSignature) {
	for _, sig := range allSignatureList {
		a.processSignatureValidation(sig)
	}
}

// processSignatureValidation performs cryptographic validation of the signature. Port of
// processSignatureValidation(AdvancedSignature).
func (a *DefaultDocumentAnalyzer) processSignatureValidation(sig validation.AdvancedSignature) {
	if sig != nil {
		sig.CheckSignatureIntegrity()
	}
}

// getTimestampedReferences returns a list of timestamped references from the given list of
// SignatureScopes. Port of getTimestampedReferences(List).
func (a *DefaultDocumentAnalyzer) getTimestampedReferences(signatureScopes []modelscope.SignatureScope) []*validation.TimestampedReference {
	overrides := a.defaultDocumentAnalyzerOverrides()
	timestampedReferences := make([]*validation.TimestampedReference, 0, len(signatureScopes))
	if utils.IsCollectionNotEmpty(signatureScopes) {
		for _, signatureScope := range signatureScopes {
			if overrides.AddReference(signatureScope) {
				timestampedReferences = append(timestampedReferences,
					validation.NewTimestampedReference(signatureScope.DSSIDAsString(), enumerations.TimestampedObjectType_SIGNED_DATA))
			}
		}
	}
	return timestampedReferences
}

// AddReference is DefaultDocumentAnalyzerOverrides' default body: true (accept all). Port of
// addReference(SignatureScope).
func (a *DefaultDocumentAnalyzer) AddReference(signatureScope modelscope.SignatureScope) bool {
	return true
}

// OriginalDocuments returns the signed document(s) without their signature(s), given a
// signature's DSS ID. Port of the getOriginalDocuments(String) overload.
func (a *DefaultDocumentAnalyzer) OriginalDocuments(signatureId string) []model.DSSDocument {
	advancedSignature := a.SignatureByID(signatureId)
	if advancedSignature != nil {
		return a.defaultDocumentAnalyzerOverrides().OriginalDocumentsForSignature(advancedSignature)
	}
	return nil
}

// SignatureByID returns the signature with the given id. Processes a custom
// TokenIdentifierProvider and counter signatures. Port of getSignatureById(String).
//
// Panics when signatureId is empty (Objects.requireNonNull("Signature Id cannot be null!")).
func (a *DefaultDocumentAnalyzer) SignatureByID(signatureId string) validation.AdvancedSignature {
	if signatureId == "" {
		panic("Signature Id cannot be null!")
	}
	for _, advancedSignature := range a.Signatures() {
		if sig := a.findSignatureRecursively(advancedSignature, signatureId); sig != nil {
			return sig
		}
	}
	return nil
}

// findSignatureRecursively ports the private findSignatureRecursively(AdvancedSignature,
// String).
func (a *DefaultDocumentAnalyzer) findSignatureRecursively(sig validation.AdvancedSignature, signatureId string) validation.AdvancedSignature {
	if a.doesIdMatch(sig, signatureId) {
		return sig
	}
	for _, counterSignature := range sig.CounterSignatures() {
		if advancedSignature := a.findSignatureRecursively(counterSignature, signatureId); advancedSignature != nil {
			return advancedSignature
		}
	}
	return nil
}

// doesIdMatch ports the private doesIdMatch(AdvancedSignature, String).
func (a *DefaultDocumentAnalyzer) doesIdMatch(sig validation.AdvancedSignature, signatureId string) bool {
	return signatureId == sig.ID() || signatureId == sig.DAIdentifier() ||
		signatureId == a.tokenIdentifierProvider.IDAsString(sig)
}

// validateSignaturePolicy performs validation of the signature policy's identifier, when
// present. Port of validateSignaturePolicy(AdvancedSignature).
func (a *DefaultDocumentAnalyzer) validateSignaturePolicy(sig validation.AdvancedSignature) {
	signaturePolicy := sig.SignaturePolicy()
	if signaturePolicy != nil {
		signaturePolicyStore := sig.SignaturePolicyStore()
		policyContent := a.extractSignaturePolicyContent(signaturePolicy, signaturePolicyStore)
		signaturePolicy.SetPolicyContent(policyContent)

		signaturePolicyValidator := a.SignaturePolicyValidatorLoader().LoadValidator(signaturePolicy)

		validationResult := signaturePolicyValidator.Validate(signaturePolicy)
		signaturePolicy.SetValidationResult(validationResult)
	}
}

// extractSignaturePolicyContent ports the private extractSignaturePolicyContent(SignaturePolicy,
// SignaturePolicyStore).
func (a *DefaultDocumentAnalyzer) extractSignaturePolicyContent(signaturePolicy *signature.SignaturePolicy,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	if signaturePolicyStore != nil {
		if signaturePolicyStore.SignaturePolicyContent() != nil {
			return signaturePolicyStore.SignaturePolicyContent()
		} else if signaturePolicyStore.SigPolDocLocalURI() != "" && a.signaturePolicyProvider != nil {
			return a.signaturePolicyProvider.GetSignaturePolicyByURL(signaturePolicyStore.SigPolDocLocalURI())
		}
	}
	if a.signaturePolicyProvider != nil {
		return a.signaturePolicyProvider.GetSignaturePolicy(signaturePolicy.Identifier(), signaturePolicy.URI())
	}
	return nil
}
