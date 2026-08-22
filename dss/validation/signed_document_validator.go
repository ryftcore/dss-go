// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/SignedDocumentValidator.java
// (DSS 6.5.RC1).
//
// Shape notes (all forced by the target language):
//
//   - Java's SignedDocumentValidator is an abstract class the format modules
//     extend. Go has no inheritance, so - following the
//     AbstractReportsBase/InitAbstractReports precedent used across this port -
//     the shared state and concrete methods live on SignedDocumentValidatorBase,
//     the one method the format validators actually override
//     (initializeDiagnosticDataBuilder()) is captured by
//     SignedDocumentValidatorOverrides, and a concrete validator registers
//     itself with InitSignedDocumentValidator in its constructor. Unlike the
//     Chain/Reports bases, registration is OPTIONAL here: Java's
//     initializeDiagnosticDataBuilder() is concrete ("default implementation"),
//     not abstract, so an unregistered base falls back to it.
//
//   - SignedDocumentValidator (the interface below) deliberately omits
//     getDocumentAnalyzer() and initializeDiagnosticDataBuilder(): both are
//     overridden in Java with COVARIANT return types
//     (CMSDocumentAnalyzer/CAdESDiagnosticDataBuilder, ...), which in Go is a
//     shadowing method with a different signature and would take the concrete
//     validator out of the interface's method set.
//
//   - Java's ServiceLoader lookup in fromDocument() becomes the explicit
//     registry in document_validator_factory.go.

package validation

import (
	"fmt"
	"io"
	"time"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	spi "github.com/ryftcore/dss-go/dss/spi"
	spiexception "github.com/ryftcore/dss-go/dss/spi/exception"
	spipolicy "github.com/ryftcore/dss-go/dss/spi/policy"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	spiexecutor "github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/validation/executor"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	"github.com/ryftcore/dss-go/dss/validation/reports"
	reportsdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// SignedDocumentValidator validates a signed document. The content of the
// document is determined automatically: it can be XML, CAdES(p7m), PDF or
// ASiC(zip). Port of the public surface of the abstract
// SignedDocumentValidator class; see the file header for the two members it
// omits.
type SignedDocumentValidator interface {
	DocumentValidator

	// IsSupported checks if the document is supported by the current
	// validator. Port of isSupported(DSSDocument).
	IsSupported(dssDocument model.DSSDocument) bool

	// SetLocale sets the Locale for report messages generation. Port of
	// setLocale(Locale).
	SetLocale(locale string)
}

// DiagnosticDataBuilderInitializer is the view of a SignedDocumentValidator
// that Java reaches through the class's own public
// initializeDiagnosticDataBuilder(): QWACValidator, for one, calls it on the
// validator fromDocument() handed it. It is a separate interface rather than a
// member of SignedDocumentValidator because Java's format-specific validators
// override the method with COVARIANT return types (CAdESDiagnosticDataBuilder,
// PAdESDiagnosticDataBuilder, ...) - the ported ones all narrow back to
// *SignedDocumentDiagnosticDataBuilder, but a future one need not, and would
// then silently drop out of SignedDocumentValidator's method set. Callers
// assert to this interface instead.
type DiagnosticDataBuilderInitializer interface {
	// InitializeDiagnosticDataBuilder creates a format-specific
	// implementation of the SignedDocumentDiagnosticDataBuilder. Port of
	// initializeDiagnosticDataBuilder().
	InitializeDiagnosticDataBuilder() *reportsdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// SignedDocumentValidatorOverrides captures the member Java's format-specific
// validators override and that SignedDocumentValidator self-calls from
// getDiagnosticData().
type SignedDocumentValidatorOverrides interface {
	// InitializeDiagnosticDataBuilder creates a format-specific
	// implementation of the SignedDocumentDiagnosticDataBuilder. Port of
	// initializeDiagnosticDataBuilder().
	InitializeDiagnosticDataBuilder() *reportsdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// SignedDocumentValidatorBase carries the state and the concrete methods of
// Java's abstract SignedDocumentValidator.
type SignedDocumentValidatorBase struct {
	// documentAnalyzer performs analysis of the document, tokens extraction
	// as well as cryptographic validation. Port of the protected final
	// documentAnalyzer field.
	documentAnalyzer analyzer.DocumentAnalyzer

	// ProcessExecutor can hold a specific DocumentProcessExecutor. Port of
	// the protected processExecutor field.
	ProcessExecutor executor.DocumentProcessExecutor

	defaultDigestAlgorithm     enumerations.DigestAlgorithm
	tokenExtractionStrategy    enumerations.TokenExtractionStrategy
	includeSemantics           bool
	validationLevel            enumerations.ValidationLevel
	locale                     string
	enableEtsiValidationReport bool

	overrides SignedDocumentValidatorOverrides
}

// NewSignedDocumentValidatorBase is the constructor taking the document
// analyzer. Port of the protected SignedDocumentValidator(DocumentAnalyzer)
// constructor, including the field initializers Java runs with it.
//
// Panics when documentAnalyzer is nil (Java's
// Objects.requireNonNull("DocumentAnalyzer cannot be null!")).
func NewSignedDocumentValidatorBase(documentAnalyzer analyzer.DocumentAnalyzer) SignedDocumentValidatorBase {
	if documentAnalyzer == nil {
		panic("DocumentAnalyzer cannot be null!")
	}
	return SignedDocumentValidatorBase{
		documentAnalyzer:        documentAnalyzer,
		defaultDigestAlgorithm:  enumerations.DigestAlgorithmSHA256,
		tokenExtractionStrategy: enumerations.TokenExtractionStrategyNone,
		includeSemantics:        false,
		validationLevel:         enumerations.ValidationLevelArchivalData,
		// Java's `private Locale locale = Locale.getDefault()`; the ported
		// i18n.NewI18nProviderForLocale documents "" as exactly that default.
		locale:                     "",
		enableEtsiValidationReport: true,
	}
}

// InitSignedDocumentValidator registers the concrete validator so that
// getDiagnosticData() dispatches initializeDiagnosticDataBuilder() onto it.
// Optional - see the file header.
func (v *SignedDocumentValidatorBase) InitSignedDocumentValidator(overrides SignedDocumentValidatorOverrides) {
	v.overrides = overrides
}

// SignedDocumentValidatorFromDocument guesses the document format and returns
// an appropriate document validator. Port of the static
// fromDocument(DSSDocument); Java's ServiceLoader iteration becomes the
// registry walk described in document_validator_factory.go, and the thrown
// UnsupportedOperationException becomes a returned error.
//
// Panics when dssDocument is nil (Java's Objects.requireNonNull("DSSDocument
// is null")).
func SignedDocumentValidatorFromDocument(dssDocument model.DSSDocument) (SignedDocumentValidator, error) {
	if dssDocument == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range documentValidatorFactoryRegistry {
		if factory.IsSupported(dssDocument) {
			return factory.Create(dssDocument), nil
		}
	}
	return nil, fmt.Errorf("Document format not recognized/handled")
}

// IsSupported checks if the document is supported by the current validator.
// Port of isSupported(DSSDocument).
func (v *SignedDocumentValidatorBase) IsSupported(dssDocument model.DSSDocument) bool {
	return v.documentAnalyzer.IsSupported(dssDocument)
}

// DocumentAnalyzer returns the current instance of DocumentAnalyzer. Port of
// getDocumentAnalyzer().
func (v *SignedDocumentValidatorBase) DocumentAnalyzer() analyzer.DocumentAnalyzer {
	return v.documentAnalyzer
}

// SetSigningCertificateSource is the port of
// setSigningCertificateSource(CertificateSource).
func (v *SignedDocumentValidatorBase) SetSigningCertificateSource(signingCertificateSource spi.CertificateSource) {
	v.documentAnalyzer.SetSigningCertificateSource(signingCertificateSource)
}

// SetCertificateVerifier is the port of
// setCertificateVerifier(CertificateVerifier).
func (v *SignedDocumentValidatorBase) SetCertificateVerifier(certificateVerifier spivalidation.CertificateVerifier) {
	v.documentAnalyzer.SetCertificateVerifier(certificateVerifier)
}

// SetValidationContextExecutor is the port of
// setValidationContextExecutor(ValidationContextExecutor).
func (v *SignedDocumentValidatorBase) SetValidationContextExecutor(validationContextExecutor spiexecutor.ValidationContextExecutor) {
	v.documentAnalyzer.SetValidationContextExecutor(validationContextExecutor)
}

// SetTokenIdentifierProvider is the port of
// setTokenIdentifierProvider(TokenIdentifierProvider).
func (v *SignedDocumentValidatorBase) SetTokenIdentifierProvider(tokenIdentifierProvider model.TokenIdentifierProvider) {
	v.documentAnalyzer.SetTokenIdentifierProvider(tokenIdentifierProvider)
}

// SetDetachedContents is the port of setDetachedContents(List).
func (v *SignedDocumentValidatorBase) SetDetachedContents(detachedContents []model.DSSDocument) {
	v.documentAnalyzer.SetDetachedContents(detachedContents)
}

// SetDetachedEvidenceRecordDocuments is the port of
// setDetachedEvidenceRecordDocuments(List).
func (v *SignedDocumentValidatorBase) SetDetachedEvidenceRecordDocuments(detachedEvidenceRecordDocuments []model.DSSDocument) {
	v.documentAnalyzer.SetDetachedEvidenceRecordDocuments(detachedEvidenceRecordDocuments)
}

// SetContainerContents is the port of setContainerContents(List).
func (v *SignedDocumentValidatorBase) SetContainerContents(containerContents []model.DSSDocument) {
	v.documentAnalyzer.SetContainerContents(containerContents)
}

// SetManifestFile is the port of setManifestFile(ManifestFile).
func (v *SignedDocumentValidatorBase) SetManifestFile(manifestFile *model.ManifestFile) {
	v.documentAnalyzer.SetManifestFile(manifestFile)
}

// SetValidationTime allows defining a custom validation time. Port of
// setValidationTime(Date).
func (v *SignedDocumentValidatorBase) SetValidationTime(validationTime time.Time) {
	v.documentAnalyzer.SetValidationTime(validationTime)
}

// SetSignaturePolicyProvider is the port of
// setSignaturePolicyProvider(SignaturePolicyProvider).
func (v *SignedDocumentValidatorBase) SetSignaturePolicyProvider(signaturePolicyProvider *spipolicy.SignaturePolicyProvider) {
	v.documentAnalyzer.SetSignaturePolicyProvider(signaturePolicyProvider)
}

// SetSignaturePolicyValidatorLoader is the port of
// setSignaturePolicyValidatorLoader(SignaturePolicyValidatorLoader).
func (v *SignedDocumentValidatorBase) SetSignaturePolicyValidatorLoader(signaturePolicyValidatorLoader spipolicy.SignaturePolicyValidatorLoader) {
	v.documentAnalyzer.SetSignaturePolicyValidatorLoader(signaturePolicyValidatorLoader)
}

// SetDefaultDigestAlgorithm is the port of
// setDefaultDigestAlgorithm(DigestAlgorithm).
//
// Panics when digestAlgorithm is empty (Java's
// Objects.requireNonNull("Default DigestAlgorithm cannot be nulL!")).
func (v *SignedDocumentValidatorBase) SetDefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	if digestAlgorithm == "" {
		panic("Default DigestAlgorithm cannot be nulL!")
	}
	v.defaultDigestAlgorithm = digestAlgorithm
}

// SetTokenExtractionStrategy is the port of
// setTokenExtractionStrategy(TokenExtractionStrategy).
//
// Panics when tokenExtractionStrategy is empty (Java's
// Objects.requireNonNull(tokenExtractionStrategy)).
func (v *SignedDocumentValidatorBase) SetTokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) {
	if tokenExtractionStrategy == "" {
		panic("tokenExtractionStrategy cannot be null")
	}
	v.tokenExtractionStrategy = tokenExtractionStrategy
}

// SetIncludeSemantics is the port of setIncludeSemantics(boolean).
func (v *SignedDocumentValidatorBase) SetIncludeSemantics(include bool) {
	v.includeSemantics = include
}

// SetValidationLevel is the port of setValidationLevel(ValidationLevel).
func (v *SignedDocumentValidatorBase) SetValidationLevel(validationLevel enumerations.ValidationLevel) {
	v.validationLevel = validationLevel
}

// SetEnableEtsiValidationReport is the port of
// setEnableEtsiValidationReport(boolean).
func (v *SignedDocumentValidatorBase) SetEnableEtsiValidationReport(enableEtsiValidationReport bool) {
	v.enableEtsiValidationReport = enableEtsiValidationReport
}

// SetProcessExecutor is the port of setProcessExecutor(DocumentProcessExecutor).
func (v *SignedDocumentValidatorBase) SetProcessExecutor(processExecutor executor.DocumentProcessExecutor) {
	v.ProcessExecutor = processExecutor
}

// ProvideProcessExecutorInstance returns the process executor, instantiating
// the default one on first use. Port of the protected
// provideProcessExecutorInstance().
func (v *SignedDocumentValidatorBase) ProvideProcessExecutorInstance() executor.DocumentProcessExecutor {
	if v.ProcessExecutor == nil {
		v.ProcessExecutor = v.DefaultProcessExecutor()
	}
	return v.ProcessExecutor
}

// DefaultProcessExecutor returns a default validator process executor. Port
// of getDefaultProcessExecutor().
func (v *SignedDocumentValidatorBase) DefaultProcessExecutor() executor.DocumentProcessExecutor {
	return executor.NewDefaultSignatureProcessExecutor()
}

// SetLocale sets the Locale for report messages generation. Port of
// setLocale(Locale).
func (v *SignedDocumentValidatorBase) SetLocale(locale string) {
	v.locale = locale
}

// ValidateDocument validates the document and all its signatures with the
// default validation policy. Port of the no-arg validateDocument().
func (v *SignedDocumentValidatorBase) ValidateDocument() (*reports.Reports, error) {
	return v.ValidateDocumentWithPolicyDocument(nil)
}

// ValidateDocumentWithPolicyPath is the port of the
// validateDocument(String) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyPath(policyResourcePath string) (*reports.Reports, error) {
	return v.ValidateDocumentWithPolicyAndCryptographicSuitePath(policyResourcePath, "")
}

// ValidateDocumentWithPolicyFile is the port of the validateDocument(File)
// overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyFile(policyFile string) (*reports.Reports, error) {
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteFile(policyFile, "")
}

// ValidateDocumentWithPolicyDocument is the port of the
// validateDocument(DSSDocument) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyDocument(policyDocument model.DSSDocument) (*reports.Reports, error) {
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteDocument(policyDocument, nil)
}

// ValidateDocumentWithPolicyReader is the port of the
// validateDocument(InputStream) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyReader(policyDataStream io.Reader) (*reports.Reports, error) {
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteReader(policyDataStream, nil)
}

// ValidateDocumentWithPolicyAndCryptographicSuitePath is the port of the
// validateDocument(String, String) overload; see document_validator.go on the
// classpath deviation. Java's Utils.isStringNotEmpty guard is kept: an empty
// path means "use the default".
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyAndCryptographicSuitePath(
	policyResourcePath, cryptographicSuitePath string) (*reports.Reports, error) {
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyResourcePath != "" {
		doc, err := model.NewFileDocument(policyResourcePath)
		if err != nil {
			return nil, spiexception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
				"Unable to load policy with URL '%s' and cryptographic suite '%s'. Reason : %s",
				policyResourcePath, cryptographicSuitePath, err.Error()), err)
		}
		policyDocument = doc
	}
	if cryptographicSuitePath != "" {
		doc, err := model.NewFileDocument(cryptographicSuitePath)
		if err != nil {
			return nil, spiexception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
				"Unable to load policy with URL '%s' and cryptographic suite '%s'. Reason : %s",
				policyResourcePath, cryptographicSuitePath, err.Error()), err)
		}
		cryptographicSuiteDocument = doc
	}
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// ValidateDocumentWithPolicyAndCryptographicSuiteFile is the port of the
// validateDocument(File, File) overload. Java tests `policyFile != null &&
// policyFile.exists()`: a missing file falls back to the default policy.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyAndCryptographicSuiteFile(
	policyFile, cryptographicSuiteFile string) (*reports.Reports, error) {
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyFile != "" {
		if doc, err := model.NewFileDocument(policyFile); err == nil {
			policyDocument = doc
		}
	}
	if cryptographicSuiteFile != "" {
		if doc, err := model.NewFileDocument(cryptographicSuiteFile); err == nil {
			cryptographicSuiteDocument = doc
		}
	}
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// ValidateDocumentWithPolicyAndCryptographicSuiteDocument is the port of the
// validateDocument(DSSDocument, DSSDocument) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyAndCryptographicSuiteDocument(
	policyDocument, cryptographicSuiteDocument model.DSSDocument) (*reports.Reports, error) {
	validationPolicy, err := v.LoadValidationPolicy(policyDocument, cryptographicSuiteDocument)
	if err != nil {
		return nil, err
	}
	return v.ValidateDocumentWithValidationPolicy(validationPolicy)
}

// ValidateDocumentWithPolicyAndCryptographicSuiteReader is the port of the
// validateDocument(InputStream, InputStream) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithPolicyAndCryptographicSuiteReader(
	policyDataStream, cryptographicSuiteStream io.Reader) (*reports.Reports, error) {
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyDataStream != nil {
		data, err := io.ReadAll(policyDataStream)
		if err != nil {
			return nil, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		policyDocument = model.NewInMemoryDocument(data)
	}
	if cryptographicSuiteStream != nil {
		data, err := io.ReadAll(cryptographicSuiteStream)
		if err != nil {
			return nil, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		cryptographicSuiteDocument = model.NewInMemoryDocument(data)
	}
	return v.ValidateDocumentWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// LoadValidationPolicy loads a validation policy from the policyDocument and
// the cryptographicSuiteDocument. When a document is not provided, a default
// policy or cryptographic suite is used, respectively. Port of the protected
// loadValidationPolicy(DSSDocument, DSSDocument).
//
// Java wraps every failure into IllegalInputException("Unable to load the
// policy"); the ported loaders signal the same failures by panicking, so the
// panic is recovered here and re-reported as that error.
func (v *SignedDocumentValidatorBase) LoadValidationPolicy(
	policyDocument, cryptographicSuiteDocument model.DSSDocument) (validationPolicy modelpolicy.ValidationPolicy, err error) {
	defer func() {
		if r := recover(); r != nil {
			validationPolicy = nil
			err = spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", fmt.Errorf("%v", r))
		}
	}()

	var validationPolicyLoader *validationpolicy.ValidationPolicyLoader
	if policyDocument == nil {
		// LOG.debug("No provided validation policy : use the default policy")
		validationPolicyLoader = v.FromDefaultValidationPolicyLoader()
	} else {
		validationPolicyLoader = validationpolicy.FromValidationPolicyDocument(policyDocument)
	}
	if cryptographicSuiteDocument != nil {
		return validationPolicyLoader.WithCryptographicSuiteDocument(cryptographicSuiteDocument).Create(), nil
	}

	return validationPolicyLoader.Create(), nil
}

// FromDefaultValidationPolicyLoader gets a default validation policy loader
// for a signature validation. Port of the protected
// fromDefaultValidationPolicyLoader().
func (v *SignedDocumentValidatorBase) FromDefaultValidationPolicyLoader() *validationpolicy.ValidationPolicyLoader {
	return validationpolicy.FromDefaultValidationPolicy()
}

// ValidateDocumentWithValidationPolicy validates the document and all its
// signatures against the given ValidationPolicy. Port of the
// validateDocument(ValidationPolicy) overload.
func (v *SignedDocumentValidatorBase) ValidateDocumentWithValidationPolicy(
	validationPolicy modelpolicy.ValidationPolicy) (*reports.Reports, error) {
	// LOG.info("Document validation...")
	if err := v.AssertConfigurationValid(); err != nil {
		return nil, err
	}

	diagnosticData := v.GetDiagnosticData()

	return v.ProcessValidationPolicy(diagnosticData, validationPolicy), nil
}

// AssertConfigurationValid verifies whether the configuration of the current
// instance of a document validator is valid. Port of the protected
// assertConfigurationValid(); Java's IllegalArgumentException becomes a
// returned error.
func (v *SignedDocumentValidatorBase) AssertConfigurationValid() error {
	if enumerations.ValidationLevelBasicSignatures == v.validationLevel &&
		(len(v.documentAnalyzer.DetachedTimestamps()) > 0 || len(v.documentAnalyzer.DetachedEvidenceRecords()) > 0) &&
		len(v.documentAnalyzer.Signatures()) == 0 {
		return fmt.Errorf("Basic Signatures validation cannot be used for timestamp or evidence record documents!")
	}
	return nil
}

// GetDiagnosticData retrieves the XmlDiagnosticData containing all
// information relevant for the validation process, including the certificate
// and revocation tokens obtained from online resources, e.g. AIA, CRL, OCSP
// (when applicable). Port of the final getDiagnosticData().
func (v *SignedDocumentValidatorBase) GetDiagnosticData() *diagnosticjaxb.XmlDiagnosticData {
	validationContext := v.documentAnalyzer.Validate()
	diagnosticDataBuilder := v.initializeDiagnosticDataBuilder()
	return v.InitDiagnosticDataFactory(diagnosticDataBuilder).
		SetDocument(v.documentAnalyzer.Document()).
		SetValidationTime(v.documentAnalyzer.ValidationTime()).
		SetTokenIdentifierProvider(v.documentAnalyzer.TokenIdentifierProvider()).
		SetValidationContext(validationContext).
		SetDefaultDigestAlgorithm(v.defaultDigestAlgorithm).
		SetTokenExtractionStrategy(v.tokenExtractionStrategy).
		Create()
}

// InitDiagnosticDataFactory creates a new instance of a factory used to
// create a Diagnostic Data. Port of the protected
// initDiagnosticDataFactory(SignedDocumentDiagnosticDataBuilder). No upstream
// subclass overrides it, so - unlike initializeDiagnosticDataBuilder() - it is
// not routed through SignedDocumentValidatorOverrides.
func (v *SignedDocumentValidatorBase) InitDiagnosticDataFactory(
	diagnosticDataBuilder *reportsdiagnostic.SignedDocumentDiagnosticDataBuilder) *reportsdiagnostic.XmlDiagnosticDataFactory {
	return reportsdiagnostic.NewXmlDiagnosticDataFactory(diagnosticDataBuilder)
}

// InitializeDiagnosticDataBuilder creates the default implementation of the
// SignedDocumentDiagnosticDataBuilder. Port of
// initializeDiagnosticDataBuilder().
func (v *SignedDocumentValidatorBase) InitializeDiagnosticDataBuilder() *reportsdiagnostic.SignedDocumentDiagnosticDataBuilder {
	// default implementation
	return reportsdiagnostic.NewSignedDocumentDiagnosticDataBuilder()
}

// initializeDiagnosticDataBuilder dispatches the virtual call onto the
// registered concrete validator, falling back to the base implementation when
// none registered (see the file header).
func (v *SignedDocumentValidatorBase) initializeDiagnosticDataBuilder() *reportsdiagnostic.SignedDocumentDiagnosticDataBuilder {
	if v.overrides != nil {
		return v.overrides.InitializeDiagnosticDataBuilder()
	}
	return v.InitializeDiagnosticDataBuilder()
}

// ProcessValidationPolicy executes the validation regarding the given
// validationPolicy. Port of the protected final
// processValidationPolicy(XmlDiagnosticData, ValidationPolicy).
func (v *SignedDocumentValidatorBase) ProcessValidationPolicy(
	diagnosticData *diagnosticjaxb.XmlDiagnosticData, validationPolicy modelpolicy.ValidationPolicy) *reports.Reports {
	exec := v.ProvideProcessExecutorInstance()
	exec.SetCurrentTime(v.documentAnalyzer.ValidationTime())
	exec.SetValidationPolicy(validationPolicy)
	exec.SetValidationLevel(v.validationLevel)
	exec.SetDiagnosticData(diagnosticData)
	exec.SetIncludeSemantics(v.includeSemantics)
	exec.SetEnableEtsiValidationReport(v.enableEtsiValidationReport)
	exec.SetLocale(v.locale)
	return exec.Execute()
}

// Signatures retrieves the signatures found in the document. Port of
// getSignatures().
func (v *SignedDocumentValidatorBase) Signatures() []spivalidation.AdvancedSignature {
	return v.documentAnalyzer.Signatures()
}

// signatureByIDAnalyzer stands in for Java's
// `documentAnalyzer instanceof DefaultDocumentAnalyzer` test. The ported
// DefaultDocumentAnalyzer leaves isSupported() abstract, so *DefaultDocumentAnalyzer
// does not itself satisfy analyzer.DocumentAnalyzer and cannot be the target of
// a type assertion; every concrete analyzer embeds it and therefore promotes
// SignatureByID, which is what the check is really after.
type signatureByIDAnalyzer interface {
	SignatureByID(signatureId string) spivalidation.AdvancedSignature
}

// SignatureByID returns the signature with the given id. Processes a custom
// TokenIdentifierProvider and counter signatures. Port of
// getSignatureById(String); Java's IllegalStateException becomes a panic, as
// it reports a wiring error rather than an input problem.
func (v *SignedDocumentValidatorBase) SignatureByID(signatureId string) spivalidation.AdvancedSignature {
	if defaultAnalyzer, ok := v.documentAnalyzer.(signatureByIDAnalyzer); ok {
		return defaultAnalyzer.SignatureByID(signatureId)
	}
	panic("The documentAnalyzer shall be an instance of DefaultDocumentAnalyzer to execute the method!")
}

// DetachedTimestamps is the port of getDetachedTimestamps().
func (v *SignedDocumentValidatorBase) DetachedTimestamps() []*spivalidation.TimestampToken {
	return v.documentAnalyzer.DetachedTimestamps()
}

// DetachedEvidenceRecords is the port of getDetachedEvidenceRecords().
func (v *SignedDocumentValidatorBase) DetachedEvidenceRecords() []spivalidation.EvidenceRecord {
	return v.documentAnalyzer.DetachedEvidenceRecords()
}

// OriginalDocuments is the port of the getOriginalDocuments(String) overload.
func (v *SignedDocumentValidatorBase) OriginalDocuments(signatureId string) []model.DSSDocument {
	return v.documentAnalyzer.OriginalDocuments(signatureId)
}

// OriginalDocumentsForSignature is the port of the
// getOriginalDocuments(AdvancedSignature) overload.
func (v *SignedDocumentValidatorBase) OriginalDocumentsForSignature(advancedSignature spivalidation.AdvancedSignature) []model.DSSDocument {
	return v.documentAnalyzer.OriginalDocumentsForSignature(advancedSignature)
}

// GetValidationData is the port of the getValidationData(Collection)
// overload.
func (v *SignedDocumentValidatorBase) GetValidationData(
	signatures []spivalidation.AdvancedSignature) (*spivalidation.ValidationDataContainer, error) {
	return v.documentAnalyzer.GetValidationData(signatures)
}

// GetValidationDataWithTimestamps is the port of the
// getValidationData(Collection, Collection) overload.
func (v *SignedDocumentValidatorBase) GetValidationDataWithTimestamps(signatures []spivalidation.AdvancedSignature,
	detachedTimestamps []*spivalidation.TimestampToken) (*spivalidation.ValidationDataContainer, error) {
	return v.documentAnalyzer.GetValidationDataWithTimestamps(signatures, detachedTimestamps)
}
