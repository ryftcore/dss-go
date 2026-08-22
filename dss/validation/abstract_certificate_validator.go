// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/AbstractCertificateValidator.java
// (DSS 6.5.RC1).
//
// Shape notes:
//
//   - Java's AbstractCertificateValidator<R extends AbstractReports,
//     PE extends ProcessExecutor<R>> becomes a generic Go struct with the same
//     two parameters. The abstract members and every member an in-tree
//     subclass overrides are captured by AbstractCertificateValidatorOverrides;
//     a concrete validator registers itself with InitAbstractCertificateValidator
//     in its constructor. Because a concrete validator embeds this struct, Go
//     method promotion supplies the base implementation of every member it does
//     not itself redeclare, so the overrides interface is satisfied for free -
//     the same arrangement as process/chain_item.go's InitChainItem.
//
//   - Java uses DiagnosticDataBuilder polymorphically (initDiagnosticDataBuilder()
//     returns the base type, build() dispatches to the subclass, and
//     QWACValidator downcasts the result back). The ported builders shadow the
//     base fluent setters with covariant return types, which Go cannot express
//     in one interface, so InitDiagnosticDataBuilder hands back a pair: the
//     concrete builder (for the virtual build() and for the downcast) plus a
//     pointer to its embedded DataBuilder (carrying the base fluent
//     setters Java's statically-typed chain calls). Both refer to the same
//     object.
//
//   - Java's classpath lookup
//     CertificateValidator.class.getResourceAsStream(getDefaultValidationPolicyPath())
//     has no Go counterpart. The dss-policy-jaxb policy resources are embedded
//     in resources/ (byte-identical copies) and the "/policy/<name>.xml"
//     classpath paths upstream returns are resolved against them.

package validation

import (
	"embed"
	"fmt"
	"io"
	"path"
	"time"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	spiexception "github.com/ryftcore/dss-go/dss/spi/exception"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	spiexecutor "github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/validation/executor"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	reportsdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// defaultValidationPolicyResources holds byte-identical copies of the
// dss-policy-jaxb src/main/resources/policy/*.xml files that upstream reaches
// through the classpath. See the file header.
//
//go:embed resources/certificate-constraint.xml resources/qwac-constraint.xml resources/eaa-constraint.xml
var defaultValidationPolicyResources embed.FS

// DiagnosticDataBuilderRef is the Go stand-in for the polymorphic use Java
// makes of the DataBuilder base type: the only member the validator
// calls through it is build().
type DiagnosticDataBuilderRef interface {
	// Build builds the XmlDiagnosticData. Port of build().
	Build() *diagnosticjaxb.XmlDiagnosticData
}

// AbstractCertificateValidatorOverrides captures the members Java declares
// abstract on AbstractCertificateValidator, plus those an in-tree subclass
// (CertificateValidator, QWACValidator) overrides and the base self-calls.
type AbstractCertificateValidatorOverrides[PE any] interface {
	// DefaultValidationPolicyPath returns the default validation policy path.
	// Port of the abstract getDefaultValidationPolicyPath().
	DefaultValidationPolicyPath() string

	// ProvideProcessExecutorInstance gets the ProcessExecutor. Port of the
	// abstract provideProcessExecutorInstance().
	ProvideProcessExecutorInstance() PE

	// DefaultProcessExecutor returns a default validator process executor.
	// Port of getDefaultProcessExecutor(), declared by
	// ProcessExecutorProvider and left abstract here.
	DefaultProcessExecutor() PE

	// AssertConfigurationValid checks if the Validator configuration is
	// valid. Port of the protected assertConfigurationValid().
	AssertConfigurationValid() error

	// PrepareValidationContext initializes and fills the Context
	// for a certificate token validation. Port of the protected
	// prepareValidationContext(CertificateVerifier).
	PrepareValidationContext(certificateVerifier spivalidation.CertificateVerifier) spivalidation.Context

	// CreateValidationContext creates a new instance of Context.
	// Port of the protected createValidationContext().
	CreateValidationContext() spivalidation.Context

	// PrepareDiagnosticDataBuilder creates a DiagnosticDataBuilder. Port of
	// the protected prepareDiagnosticDataBuilder().
	PrepareDiagnosticDataBuilder() DiagnosticDataBuilderRef

	// CreateDiagnosticDataBuilder creates and fills the DataBuilder
	// with the relevant data. Port of the protected
	// createDiagnosticDataBuilder(Context).
	CreateDiagnosticDataBuilder(validationContext spivalidation.Context) DiagnosticDataBuilderRef

	// InitDiagnosticDataBuilder instantiates a new DataBuilder.
	// Port of the protected initDiagnosticDataBuilder(); see the file header
	// on the returned pair.
	InitDiagnosticDataBuilder() (DiagnosticDataBuilderRef, *reportsdiagnostic.DataBuilder)
}

// AbstractCertificateValidator contains common configuration and methods for
// a CertificateToken validation. Port of the abstract
// AbstractCertificateValidator<R, PE> (implements ProcessExecutorProvider<PE>).
type AbstractCertificateValidator[R any, PE executor.ProcessExecutor[R]] struct {
	// ValidationTimeValue is the validation time. Port of the protected
	// validationTime field.
	ValidationTimeValue time.Time

	// CertificateVerifierValue is the CertificateVerifier to use. Port of the
	// protected certificateVerifier field.
	CertificateVerifierValue spivalidation.CertificateVerifier

	// TokenExtractionStrategy is the used TokenExtractionStrategy. Port of
	// the protected tokenExtractionStrategy field.
	TokenExtractionStrategy enumerations.TokenExtractionStrategy

	// IdentifierProvider is the token identifier provider to use. Port of the
	// protected identifierProvider field.
	IdentifierProvider model.TokenIdentifierProvider

	// ValidationContextExecutor performs validation of the Context.
	// Port of the protected validationContextExecutor field.
	ValidationContextExecutor spiexecutor.ValidationContextExecutor

	// Locale is the locale to use for reports generation. Port of the
	// protected locale field; "" is the ported equivalent of Java's
	// Locale.getDefault() (see i18n.NewProviderForLocale).
	Locale string

	// ProcessExecutor is the certificate process executor. Port of the
	// protected processExecutor field.
	ProcessExecutor PE

	// DefaultDigestAlgorithm is the DigestAlgorithm used for the calculation
	// of digests for validation tokens and signed data. Port of the protected
	// defaultDigestAlgorithm field.
	DefaultDigestAlgorithm enumerations.DigestAlgorithm

	overrides AbstractCertificateValidatorOverrides[PE]
}

// NewAbstractCertificateValidator is the default constructor. Port of the
// protected AbstractCertificateValidator() constructor plus the field
// initializers Java runs with it.
func NewAbstractCertificateValidator[R any, PE executor.ProcessExecutor[R]]() AbstractCertificateValidator[R, PE] {
	return AbstractCertificateValidator[R, PE]{
		TokenExtractionStrategy:   enumerations.TokenExtractionStrategyNone,
		IdentifierProvider:        model.NewOriginalIdentifierProvider(),
		ValidationContextExecutor: spiexecutor.DefaultValidationContextExecutorInstance,
		Locale:                    "",
		DefaultDigestAlgorithm:    enumerations.DigestAlgorithmSHA256,
	}
}

// InitAbstractCertificateValidator registers the concrete validator so the
// base methods can dispatch the abstract and overridden members onto it.
// Every concrete constructor must call this before use.
func (v *AbstractCertificateValidator[R, PE]) InitAbstractCertificateValidator(
	overrides AbstractCertificateValidatorOverrides[PE]) {
	v.overrides = overrides
}

// certificateValidatorOverrides returns the registered overrides, panicking
// when the concrete validator failed to call InitAbstractCertificateValidator.
func (v *AbstractCertificateValidator[R, PE]) certificateValidatorOverrides() AbstractCertificateValidatorOverrides[PE] {
	if v.overrides == nil {
		panic("validation: AbstractCertificateValidator used without InitAbstractCertificateValidator")
	}
	return v.overrides
}

// SetCertificateVerifier sets the CertificateVerifier. Port of
// setCertificateVerifier(CertificateVerifier).
func (v *AbstractCertificateValidator[R, PE]) SetCertificateVerifier(certificateVerifier spivalidation.CertificateVerifier) {
	v.CertificateVerifierValue = certificateVerifier
}

// CertificateVerifier returns the configured CertificateVerifier. It has no
// Java counterpart (the field is protected and read directly by the
// subclasses); Go's separate packages need an accessor.
func (v *AbstractCertificateValidator[R, PE]) CertificateVerifier() spivalidation.CertificateVerifier {
	return v.CertificateVerifierValue
}

// SetTokenExtractionStrategy sets the TokenExtractionStrategy. Port of
// setTokenExtractionStrategy(TokenExtractionStrategy).
//
// Panics when tokenExtractionStrategy is empty (Java's
// Objects.requireNonNull(tokenExtractionStrategy)).
func (v *AbstractCertificateValidator[R, PE]) SetTokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) {
	if tokenExtractionStrategy == "" {
		panic("tokenExtractionStrategy cannot be null")
	}
	v.TokenExtractionStrategy = tokenExtractionStrategy
}

// SetTokenIdentifierProvider sets the TokenIdentifierProvider. Port of
// setTokenIdentifierProvider(TokenIdentifierProvider).
//
// Panics when identifierProvider is nil (Java's
// Objects.requireNonNull(identifierProvider)).
func (v *AbstractCertificateValidator[R, PE]) SetTokenIdentifierProvider(identifierProvider model.TokenIdentifierProvider) {
	if identifierProvider == nil {
		panic("identifierProvider cannot be null")
	}
	v.IdentifierProvider = identifierProvider
}

// SetValidationTime sets the validationTime. Port of setValidationTime(Date).
func (v *AbstractCertificateValidator[R, PE]) SetValidationTime(validationTime time.Time) {
	v.ValidationTimeValue = validationTime
}

// SetLocale sets the Locale to use for messages in reports. Port of
// setLocale(Locale).
func (v *AbstractCertificateValidator[R, PE]) SetLocale(locale string) {
	v.Locale = locale
}

// ValidationTime gets the validation time, instantiating it to the current
// time when not provided explicitly. Port of the protected
// getValidationTime().
func (v *AbstractCertificateValidator[R, PE]) ValidationTime() time.Time {
	if v.ValidationTimeValue.IsZero() {
		v.ValidationTimeValue = time.Now()
	}
	return v.ValidationTimeValue
}

// SetValidationContextExecutor sets the ValidationContextExecutor used for
// validation of the prepared ValidationContext. Port of
// setValidationContextExecutor(ValidationContextExecutor).
//
// Panics when validationContextExecutor is nil (Java's
// Objects.requireNonNull(validationContextExecutor)).
func (v *AbstractCertificateValidator[R, PE]) SetValidationContextExecutor(validationContextExecutor spiexecutor.ValidationContextExecutor) {
	if validationContextExecutor == nil {
		panic("validationContextExecutor cannot be null")
	}
	v.ValidationContextExecutor = validationContextExecutor
}

// SetDefaultDigestAlgorithm changes the DigestAlgorithm used for tokens'
// digest calculation. Port of setDefaultDigestAlgorithm(DigestAlgorithm).
//
// Panics when digestAlgorithm is empty (Java's
// Objects.requireNonNull("Default DigestAlgorithm cannot be nulL!")).
func (v *AbstractCertificateValidator[R, PE]) SetDefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	if digestAlgorithm == "" {
		panic("Default DigestAlgorithm cannot be nulL!")
	}
	v.DefaultDigestAlgorithm = digestAlgorithm
}

// SetProcessExecutor sets a specific process executor. Port of
// setProcessExecutor(PE).
func (v *AbstractCertificateValidator[R, PE]) SetProcessExecutor(processExecutor PE) {
	v.ProcessExecutor = processExecutor
}

// Validate validates the certificate with a default validation policy. Port
// of the no-arg validate().
func (v *AbstractCertificateValidator[R, PE]) Validate() (R, error) {
	return v.ValidateWithPolicyDocument(nil)
}

// ValidateWithPolicyPath validates the certificate with the policy obtained
// from policyResourcePath. Port of the validate(String) overload; see
// document_validator.go on the classpath deviation.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyPath(policyResourcePath string) (R, error) {
	return v.ValidateWithPolicyAndCryptographicSuitePath(policyResourcePath, "")
}

// ValidateWithPolicyFile validates the certificate with the policy obtained
// from policyFile. Port of the validate(File) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyFile(policyFile string) (R, error) {
	return v.ValidateWithPolicyAndCryptographicSuiteFile(policyFile, "")
}

// ValidateWithPolicyDocument validates the certificate with the policy
// carried by policyDocument. Port of the validate(DSSDocument) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyDocument(policyDocument model.DSSDocument) (R, error) {
	return v.ValidateWithPolicyAndCryptographicSuiteDocument(policyDocument, nil)
}

// ValidateWithPolicyReader validates the certificate with the policy read
// from policyDataStream. Port of the validate(InputStream) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyReader(policyDataStream io.Reader) (R, error) {
	return v.ValidateWithPolicyAndCryptographicSuiteReader(policyDataStream, nil)
}

// ValidateWithPolicyAndCryptographicSuitePath is the port of the
// validate(String, String) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyAndCryptographicSuitePath(
	policyResourcePath, cryptographicSuitePath string) (R, error) {
	var zero R
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyResourcePath != "" {
		doc, err := model.NewFileDocument(policyResourcePath)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
				"Unable to load policy with URL '%s' and cryptographic suite '%s'. Reason : %s",
				policyResourcePath, cryptographicSuitePath, err.Error()), err)
		}
		policyDocument = doc
	}
	if cryptographicSuitePath != "" {
		doc, err := model.NewFileDocument(cryptographicSuitePath)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
				"Unable to load policy with URL '%s' and cryptographic suite '%s'. Reason : %s",
				policyResourcePath, cryptographicSuitePath, err.Error()), err)
		}
		cryptographicSuiteDocument = doc
	}
	return v.ValidateWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// ValidateWithPolicyAndCryptographicSuiteFile is the port of the
// validate(File, File) overload. Unlike SignedDocumentValidator, Java does not
// test File.exists() here: a missing file surfaces as a load failure.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyAndCryptographicSuiteFile(
	policyFile, cryptographicSuiteFile string) (R, error) {
	var zero R
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyFile != "" {
		doc, err := model.NewFileDocument(policyFile)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		policyDocument = doc
	}
	if cryptographicSuiteFile != "" {
		doc, err := model.NewFileDocument(cryptographicSuiteFile)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		cryptographicSuiteDocument = doc
	}
	return v.ValidateWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// ValidateWithPolicyAndCryptographicSuiteDocument is the port of the
// validate(DSSDocument, DSSDocument) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyAndCryptographicSuiteDocument(
	policyDocument, cryptographicSuiteDocument model.DSSDocument) (R, error) {
	validationPolicy, err := v.LoadValidationPolicy(policyDocument, cryptographicSuiteDocument)
	if err != nil {
		var zero R
		return zero, err
	}
	return v.ValidateWithValidationPolicy(validationPolicy)
}

// ValidateWithPolicyAndCryptographicSuiteReader is the port of the
// validate(InputStream, InputStream) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithPolicyAndCryptographicSuiteReader(
	policyDataStream, cryptographicSuiteStream io.Reader) (R, error) {
	var zero R
	var policyDocument, cryptographicSuiteDocument model.DSSDocument
	if policyDataStream != nil {
		data, err := io.ReadAll(policyDataStream)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		policyDocument = model.NewInMemoryDocument(data)
	}
	if cryptographicSuiteStream != nil {
		data, err := io.ReadAll(cryptographicSuiteStream)
		if err != nil {
			return zero, spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", err)
		}
		cryptographicSuiteDocument = model.NewInMemoryDocument(data)
	}
	return v.ValidateWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument)
}

// LoadValidationPolicy loads a validation policy from the policyDocument and
// the cryptographicSuiteDocument. Port of the protected
// loadValidationPolicy(DSSDocument, DSSDocument); the ported loaders signal
// failures by panicking, so the panic is recovered and re-reported as Java's
// IllegalInputException("Unable to load the policy").
func (v *AbstractCertificateValidator[R, PE]) LoadValidationPolicy(
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
		validationPolicyLoader = v.FromDefaultCertificateValidationPolicyLoader()
	} else {
		validationPolicyLoader = validationpolicy.FromValidationPolicyDocument(policyDocument)
	}
	if cryptographicSuiteDocument != nil {
		return validationPolicyLoader.WithCryptographicSuiteDocument(cryptographicSuiteDocument).Create(), nil
	}

	return validationPolicyLoader.Create(), nil
}

// FromDefaultCertificateValidationPolicyLoader gets a default validation
// policy loader for a certificate validation. Port of the protected
// fromDefaultCertificateValidationPolicyLoader(); see the file header on the
// classpath deviation.
func (v *AbstractCertificateValidator[R, PE]) FromDefaultCertificateValidationPolicyLoader() *validationpolicy.ValidationPolicyLoader {
	policyPath := v.certificateValidatorOverrides().DefaultValidationPolicyPath()
	data, err := defaultValidationPolicyResources.ReadFile("resources/" + path.Base(policyPath))
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default validation policy '%s'. Reason : %s", policyPath, err.Error()))
	}
	return validationpolicy.FromValidationPolicyDocument(model.NewInMemoryDocument(data))
}

// ValidateWithValidationPolicy validates the certificate with a custom
// validation policy. Port of the validate(ValidationPolicy) overload.
func (v *AbstractCertificateValidator[R, PE]) ValidateWithValidationPolicy(
	validationPolicy modelpolicy.ValidationPolicy) (R, error) {
	overrides := v.certificateValidatorOverrides()
	if err := overrides.AssertConfigurationValid(); err != nil {
		var zero R
		return zero, err
	}

	diagnosticData := v.GetDiagnosticData()

	exec := overrides.ProvideProcessExecutorInstance()
	exec.SetValidationPolicy(validationPolicy)
	exec.SetDiagnosticData(diagnosticData)
	exec.SetLocale(v.Locale)
	exec.SetCurrentTime(v.ValidationTime())
	return exec.Execute(), nil
}

// AssertConfigurationValid checks if the Validator configuration is valid.
// Port of the protected assertConfigurationValid(); Java's
// Objects.requireNonNull becomes a returned error so the subclass override can
// chain onto it.
func (v *AbstractCertificateValidator[R, PE]) AssertConfigurationValid() error {
	if v.CertificateVerifierValue == nil {
		return fmt.Errorf("CertificateVerifier is not defined")
	}
	return nil
}

// GetDiagnosticData retrieves the XmlDiagnosticData containing all
// information relevant for the validation process, including the certificate
// and revocation tokens obtained from online resources, e.g. AIA, CRL, OCSP
// (when applicable). Port of the final getDiagnosticData().
func (v *AbstractCertificateValidator[R, PE]) GetDiagnosticData() *diagnosticjaxb.XmlDiagnosticData {
	return v.certificateValidatorOverrides().PrepareDiagnosticDataBuilder().Build()
}

// PrepareValidationContext initializes and fills the Context for a
// certificate token validation. Port of the protected
// prepareValidationContext(CertificateVerifier).
func (v *AbstractCertificateValidator[R, PE]) PrepareValidationContext(
	certificateVerifier spivalidation.CertificateVerifier) spivalidation.Context {
	svc := v.certificateValidatorOverrides().CreateValidationContext()
	svc.Initialize(certificateVerifier)
	return svc
}

// CreateValidationContext creates a new instance of Context
// performing preparation of validation data, certificate chain building,
// revocation request, as well as custom validation checks execution. Port of
// the protected createValidationContext().
func (v *AbstractCertificateValidator[R, PE]) CreateValidationContext() spivalidation.Context {
	return spivalidation.NewSignatureValidationContextAtTime(v.ValidationTime())
}

// PrepareDiagnosticDataBuilder creates a DiagnosticDataBuilder. Port of the
// protected prepareDiagnosticDataBuilder().
func (v *AbstractCertificateValidator[R, PE]) PrepareDiagnosticDataBuilder() DiagnosticDataBuilderRef {
	overrides := v.certificateValidatorOverrides()
	certificateVerifierForValidation := spivalidation.NewCertificateVerifierBuilder(
		v.CertificateVerifierValue).BuildCompleteCopyForValidation()
	validationContext := overrides.PrepareValidationContext(certificateVerifierForValidation)
	v.ValidateContext(validationContext)
	return overrides.CreateDiagnosticDataBuilder(validationContext)
}

// ValidateContext processes the validation. Port of the protected
// validateContext(Context).
func (v *AbstractCertificateValidator[R, PE]) ValidateContext(validationContext spivalidation.Context) {
	v.ValidationContextExecutor.Validate(validationContext)
}

// CreateDiagnosticDataBuilder creates and fills the DataBuilder
// with the relevant data. Port of the protected
// createDiagnosticDataBuilder(Context); the fluent chain runs on the
// embedded base builder, exactly as Java's statically-typed chain does.
func (v *AbstractCertificateValidator[R, PE]) CreateDiagnosticDataBuilder(
	validationContext spivalidation.Context) DiagnosticDataBuilderRef {
	ref, base := v.certificateValidatorOverrides().InitDiagnosticDataBuilder()
	base.UsedCertificates(validationContext.GetProcessedCertificates()).
		UsedRevocations(validationContext.GetProcessedRevocations()).
		AllCertificateSources(validationContext.GetAllCertificateSources()).
		DefaultDigestAlgorithm(v.DefaultDigestAlgorithm).
		TokenExtractionStrategy(v.TokenExtractionStrategy).
		TokenIdentifierProvider(v.IdentifierProvider).
		ValidationDate(v.ValidationTime())
	return ref
}

// InitDiagnosticDataBuilder instantiates a new DiagnosticDataBuilder. Port of
// the protected initDiagnosticDataBuilder(); see the file header on the
// returned pair.
func (v *AbstractCertificateValidator[R, PE]) InitDiagnosticDataBuilder() (DiagnosticDataBuilderRef, *reportsdiagnostic.DataBuilder) {
	builder := reportsdiagnostic.NewCertificateDiagnosticDataBuilder()
	return builder, &builder.DataBuilder
}
