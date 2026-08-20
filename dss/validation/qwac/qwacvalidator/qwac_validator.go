// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/QWACValidator.java (DSS 6.5.RC1).
//
// See doc.go for why this file lives in its own package instead of dss/validation/qwac.
//
// FORWARD DEPENDENCY (narrow, flagged for the integrator): executor.QWACCertificateProcessExecutor
// and its constructor are part of this same phase 8f batch's EXEC chunk manifest
// (dss-validation/src/main/java/.../executor/certificate/qwac/QWACCertificateProcessExecutor.java,
// flattened into dss/validation/executor per that chunk's own manifest note) and had not landed
// at the time this file was written; DefaultProcessExecutor below is the only call site that
// needs it. Every other type this file touches (AbstractCertificateValidator,
// SignedDocumentValidator, CertificateProcessExecutor, CertificateReports,
// QWACCertificateDiagnosticDataBuilder, ...) is confirmed landed and used against its real
// signature.
//
// CROSS-CHUNK GAP (flagged for the integrator): createQWACDiagnosticDataBuilder needs to call
// initializeDiagnosticDataBuilder() on an arbitrary SignedDocumentValidator obtained from
// dssvalidation.SignedDocumentValidatorFromDocument (the TLS Certificate Binding signature's
// own validator - CAdES, XAdES, JAdES, whichever format the binding signature happens to use).
// The exported dssvalidation.SignedDocumentValidator interface deliberately omits this method
// (see signed_document_validator.go's file header: format validators override it with a
// covariant return type, which Go cannot express in the interface's method set). This file
// reaches it through the exported dssvalidation.SignedDocumentValidatorOverrides interface
// instead (the same method Java calls, at its base, non-covariant, return type - which is all
// this call site ever needs, since the result is only fed to
// QWACCertificateDiagnosticDataBuilder.SetSignatureDiagnosticDataBuilder, itself base-typed).
// This assumes every concrete per-format validator (still gated behind UNGATE_A/UNGATE_B at the
// time of writing) registers an adapter satisfying SignedDocumentValidatorOverrides verbatim, as
// SignedDocumentValidatorBase.InitSignedDocumentValidator's own doc comment prescribes; if a
// concrete validator's own covariant-return method instead shadows the name without also
// exposing a base-typed adapter, the type assertion below fails and this port falls back to
// Java's other reachable behavior for a null diagnostic data builder (an empty signature
// section), which is a very close approximation - never a nil pointer panic downstream. Revisit
// once UNGATE lands the concrete validators, and drop the type assertion here for a direct
// interface call if SignedDocumentValidatorOverrides is folded into the public interface later.
package qwacvalidator

import (
	stdx509 "crypto/x509"
	"fmt"

	"github.com/utain/esig/dss/model"
	modelhttp "github.com/utain/esig/dss/model/http"
	"github.com/utain/esig/dss/spi"
	spihttp "github.com/utain/esig/dss/spi/client/http"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	dssvalidation "github.com/utain/esig/dss/validation"
	"github.com/utain/esig/dss/validation/executor"
	"github.com/utain/esig/dss/validation/qwac"
	"github.com/utain/esig/dss/validation/reports"
	reportsdiagnostic "github.com/utain/esig/dss/validation/reports/diagnostic"
)

// qwacValidationPolicyLocation is the path for the default QWAC validation policy.
const qwacValidationPolicyLocation = "/policy/qwac-constraint.xml"

// QWACValidator performs a validation of a TLS/SSL certificate as per ETSI TS 119 411-5 "Policy
// and security requirements for Trust Service Providers issuing certificates; Part 5:
// Implementation of qualified certificates for website authentication as in amended Regulation
// 910/2014".
type QWACValidator struct {
	dssvalidation.AbstractCertificateValidator[*reports.CertificateReports, executor.CertificateProcessExecutor]

	// url is the URL to verify the used TSL/SSL certificate.
	url string

	// certificateToken is the TSL/SSL certificate derived from url.
	certificateToken *model.CertificateToken

	// dataLoader establishes the secure TLS/SSL connection with a remote server and/or loads
	// any applicable information required for a QWAC validation.
	dataLoader spihttp.AdvancedDataLoader
}

// compile-time interface assertion: QWACValidator supplies every member
// AbstractCertificateValidator dispatches virtually, either by overriding it below or, for the
// four it does not (AssertConfigurationValid, PrepareValidationContext, CreateValidationContext,
// CreateDiagnosticDataBuilder), by inheriting the promoted base implementation unchanged -
// exactly as Java's QWACValidator inherits those same four from AbstractCertificateValidator.
var _ dssvalidation.AbstractCertificateValidatorOverrides[executor.CertificateProcessExecutor] = (*QWACValidator)(nil)

// newQWACValidator ports both private constructors.
//
// Panics when url is empty (Java's Objects.requireNonNull("URL cannot be null!")).
func newQWACValidator(url string, certificateToken *model.CertificateToken) *QWACValidator {
	if url == "" {
		panic("URL cannot be null!")
	}
	v := &QWACValidator{
		AbstractCertificateValidator: dssvalidation.NewAbstractCertificateValidator[*reports.CertificateReports, executor.CertificateProcessExecutor](),
		url:                          url,
		certificateToken:             certificateToken,
	}
	v.InitAbstractCertificateValidator(v)
	return v
}

// QWACValidatorFromURL instantiates a new QWACValidator to verify the TSL/SSL certificate from
// the specified url. When loaded with this function, QWACValidator performs a request to the
// remote url to retrieve the actual TSL/SSL certificate and perform its validation. Port of
// fromUrl(String).
func QWACValidatorFromURL(url string) *QWACValidator {
	return newQWACValidator(url, nil)
}

// QWACValidatorFromURLAndCertificate instantiates a new QWACValidator to verify the provided
// TSL/SSL certificateToken against the specified url. When loaded with this function,
// QWACValidator validates the provided certificateToken whether it can be used as a QWAC TSL/SSL
// certificate for the url. Port of fromUrlAndCertificate(String, CertificateToken).
func QWACValidatorFromURLAndCertificate(url string, certificateToken *model.CertificateToken) *QWACValidator {
	return newQWACValidator(url, certificateToken)
}

// dataLoaderOrDefault gets the data loader to be used for accessing the information from remote
// sources. If not defined with SetDataLoader, instantiates a NativeHTTPDataLoader by default.
// Port of the protected getDataLoader().
func (v *QWACValidator) dataLoaderOrDefault() spihttp.AdvancedDataLoader {
	if v.dataLoader == nil {
		v.dataLoader = spihttp.NewNativeHTTPDataLoader()
	}
	return v.dataLoader
}

// SetDataLoader sets a data loader used to establish a TLS/SSL connection and retrieve any
// related information over. If not set, a default instance of NativeHTTPDataLoader is used for
// remote calls, if any.
func (v *QWACValidator) SetDataLoader(dataLoader spihttp.AdvancedDataLoader) {
	v.dataLoader = dataLoader
}

// PrepareDiagnosticDataBuilder ports the protected prepareDiagnosticDataBuilder() override.
//
// Every failure below is an unchecked exception in Java (DSSException,
// UnsupportedOperationException, IllegalArgumentException) propagating straight out of
// prepareDiagnosticDataBuilder(), which itself declares no checked exception; each is
// reproduced as a panic, matching this interface method's own error-free signature.
func (v *QWACValidator) PrepareDiagnosticDataBuilder() dssvalidation.DiagnosticDataBuilderRef {
	response := v.connectToURL()
	v.assertResponseValid(response)

	tlsCertificates := v.toCertificateTokenList(response.TLSCertificates())
	tlsCertificateBindingURL := v.readTLSCertificateBindingURL(response)
	signedDocumentValidator := v.initSignedDocumentValidator(tlsCertificateBindingURL, tlsCertificates)
	signature := v.getTLSCertificateBindingSignature(signedDocumentValidator)

	tlsCertificate := v.getTLSCertificate(tlsCertificates)
	validationContext := v.prepareQWACValidationContext(tlsCertificate, tlsCertificates, signature)
	v.ValidateContext(validationContext)
	return v.createQWACDiagnosticDataBuilder(validationContext, signedDocumentValidator, tlsCertificate, tlsCertificateBindingURL, signature)
}

// connectToURL connects to url. Port of the protected connectToUrl().
//
// Panics for a non-HTTP(S) URL (Java's thrown UnsupportedOperationException).
func (v *QWACValidator) connectToURL() *modelhttp.ResponseEnvelope {
	trimmedURL := utils.Trim(v.url)
	if !spihttp.ProtocolIsHttpURL(trimmedURL) {
		panic(fmt.Sprintf("DSS framework supports only HTTP(S) certificate extraction. Obtained URL : '%s'", v.url))
	}
	return v.dataLoaderOrDefault().RequestGetFull(v.url, true, false)
}

// prepareQWACValidationContext prepares a ValidationContext using the configuration and
// provided data objects. Port of the protected prepareValidationContext(CertificateToken, List,
// AdvancedSignature) - an overload of, not an override of, the base 1-arg
// prepareValidationContext(CertificateVerifier) that PrepareValidationContext (promoted,
// unshadowed, from AbstractCertificateValidator below) already satisfies.
func (v *QWACValidator) prepareQWACValidationContext(tlsCertificate *model.CertificateToken, otherTLSCertificates []*model.CertificateToken,
	signature spivalidation.AdvancedSignature) spivalidation.ValidationContext {
	certificateVerifierForValidation := v.offlineCertificateVerifier()

	validationContext := v.PrepareValidationContext(certificateVerifierForValidation)
	validationContext.AddCertificateTokenForVerification(tlsCertificate)

	adjunctCertificateSource := spi.NewCommonCertificateSource()
	for _, certificate := range otherTLSCertificates {
		adjunctCertificateSource.AddCertificate(certificate)
	}
	certificateVerifierForValidation.AddAdjunctCertSources(&adjunctCertificateSource)

	if signature != nil {
		validationContext.AddSignatureForVerification(signature)
	}

	return validationContext
}

// getTLSCertificate ports the private getTLSCertificate(List).
//
// Panics when no TSL/SSL certificate is available (Java's thrown
// `new DSSException("No valid TLS/SSL certificates have been obtained from the URL '{}'.")` -
// note the unsubstituted "'{}'." literal is upstream's own bug: DSSException's message
// constructor does not do SLF4J-style {}-placeholder substitution, so the url argument is
// simply never supplied. Reproduced verbatim.)
func (v *QWACValidator) getTLSCertificate(tlsCertificates []*model.CertificateToken) *model.CertificateToken {
	if v.certificateToken != nil {
		return v.certificateToken
	} else if utils.IsCollectionNotEmpty(tlsCertificates) {
		// first certificate shall be the peer's end-entity certificate
		v.certificateToken = tlsCertificates[0]
		return v.certificateToken
	}
	panic("No valid TLS/SSL certificates have been obtained from the URL '{}'.")
}

// initSignedDocumentValidator ports the private initSignedDocumentValidator(String, List).
func (v *QWACValidator) initSignedDocumentValidator(tlsCertificateBindingURL string, tlsCertificates []*model.CertificateToken) dssvalidation.SignedDocumentValidator {
	if tlsCertificateBindingURL == "" {
		return nil
	}
	tlsCertificateBindingSignatureBytes := v.dataLoaderOrDefault().Get(tlsCertificateBindingURL)
	if tlsCertificateBindingSignatureBytes == nil {
		return nil
	}
	signatureDocument := model.NewInMemoryDocumentWithName(tlsCertificateBindingSignatureBytes, tlsCertificateBindingURL)
	documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(signatureDocument)
	if err != nil {
		panic(err.Error())
	}
	documentValidator.SetCertificateVerifier(v.offlineCertificateVerifier())
	documentValidator.SetDetachedContents(v.toDetachedDocumentsList(tlsCertificates))
	return documentValidator
}

// getTLSCertificateBindingSignature ports the private getTLSCertificateBindingSignature(SignedDocumentValidator).
func (v *QWACValidator) getTLSCertificateBindingSignature(signedDocumentValidator dssvalidation.SignedDocumentValidator) spivalidation.AdvancedSignature {
	if signedDocumentValidator != nil {
		signatures := signedDocumentValidator.Signatures()
		if len(signatures) == 1 {
			return signatures[0]
		}
		// Java: LOG.warn when more than one signature is found; dropped per PORTING.md.
	}
	return nil
}

// toCertificateTokenList ports the private toCertificateTokenList(Certificate[]).
func (v *QWACValidator) toCertificateTokenList(certificates []*stdx509.Certificate) []*model.CertificateToken {
	if len(certificates) == 0 {
		return nil
	}
	result := make([]*model.CertificateToken, 0, len(certificates))
	for _, certificate := range certificates {
		token, err := model.NewCertificateToken(certificate)
		if err != nil {
			// Java: LOG.warn("Unable to load certificate : ...") and skip; dropped per
			// PORTING.md, the entry is still skipped.
			continue
		}
		result = append(result, token)
	}
	return result
}

// toDetachedDocumentsList ports the private toDetachedDocumentsList(List).
func (v *QWACValidator) toDetachedDocumentsList(certificates []*model.CertificateToken) []model.DSSDocument {
	result := make([]model.DSSDocument, 0, len(certificates))
	for _, certificate := range certificates {
		result = append(result, model.NewInMemoryDocumentWithName(certificate.Encoded(), v.IdentifierProvider.IDAsString(certificate)))
	}
	return result
}

// readTLSCertificateBindingURL reads the HTTP response headers and extracts the value of the
// "Link" header with a "rel" value of "tls-certificate-binding". Port of the protected
// readTLSCertificateBindingUrl(ResponseEnvelope).
func (v *QWACValidator) readTLSCertificateBindingURL(responseEnvelope *modelhttp.ResponseEnvelope) string {
	return qwac.GetTLSCertificateBindingUrl(responseEnvelope.Headers())
}

// assertResponseValid verifies whether response is valid and contains the information required
// to continue the QWAC validation process. Port of the protected assertResponseValid(ResponseEnvelope).
//
// Panics when invalid (Java's thrown IllegalArgumentException).
func (v *QWACValidator) assertResponseValid(response *modelhttp.ResponseEnvelope) {
	if v.certificateToken == nil && len(response.TLSCertificates()) == 0 {
		panic(fmt.Sprintf("No TSL/SSL certificates have been returned from the URL '%s'. "+
			"Please ensure the URL is valid and uses the HTTP(S) scheme, or provide the TLS/SSL certificate on validation explicitly.", v.url))
	}
}

// createQWACDiagnosticDataBuilder creates and configures a new DiagnosticDataBuilder. Port of
// the protected createDiagnosticDataBuilder(ValidationContext, SignedDocumentValidator,
// CertificateToken, String, AdvancedSignature) - an overload of, not an override of, the base
// 1-arg createDiagnosticDataBuilder(ValidationContext) that CreateDiagnosticDataBuilder
// (promoted, unshadowed, from AbstractCertificateValidator below) already satisfies; this
// method calls that promoted 1-arg method directly, exactly as Java's
// `super.createDiagnosticDataBuilder(validationContext)` does.
func (v *QWACValidator) createQWACDiagnosticDataBuilder(validationContext spivalidation.ValidationContext,
	signedDocumentValidator dssvalidation.SignedDocumentValidator, tlsCertificate *model.CertificateToken,
	tlsCertificateBindingURL string, signature spivalidation.AdvancedSignature) dssvalidation.DiagnosticDataBuilderRef {
	base := v.CreateDiagnosticDataBuilder(validationContext)
	diagnosticDataBuilder, ok := base.(*reportsdiagnostic.QWACCertificateDiagnosticDataBuilder)
	if !ok {
		panic(fmt.Sprintf("qwacvalidator: expected a *QWACCertificateDiagnosticDataBuilder from InitDiagnosticDataBuilder, got %T", base))
	}

	diagnosticDataBuilder = diagnosticDataBuilder.
		WebsiteUrl(v.url).
		TLSCertificate(tlsCertificate).
		TLSCertificateBindingUrl(tlsCertificateBindingURL).
		TLSCertificateBindingSignature(signature)

	if signedDocumentValidator != nil {
		if overrides, ok := signedDocumentValidator.(dssvalidation.SignedDocumentValidatorOverrides); ok {
			// See this file's header on why this reaches through
			// SignedDocumentValidatorOverrides rather than the public SignedDocumentValidator
			// interface.
			diagnosticDataBuilder = diagnosticDataBuilder.SetSignatureDiagnosticDataBuilder(overrides.InitializeDiagnosticDataBuilder())
		}
	}

	return diagnosticDataBuilder.
		FoundSignatures(validationContext.GetProcessedSignatures()).
		DocumentCertificateSource(validationContext.GetDocumentCertificateSource())
}

// offlineCertificateVerifier builds a complete copy of the CertificateVerifier for the offline
// validation. Port of the protected getOfflineCertificateVerifier().
func (v *QWACValidator) offlineCertificateVerifier() spivalidation.CertificateVerifier {
	return spivalidation.NewCertificateVerifierBuilder(v.CertificateVerifier()).BuildCompleteCopyForValidation()
}

// InitDiagnosticDataBuilder ports the protected initDiagnosticDataBuilder() override.
func (v *QWACValidator) InitDiagnosticDataBuilder() (dssvalidation.DiagnosticDataBuilderRef, *reportsdiagnostic.DiagnosticDataBuilder) {
	builder := reportsdiagnostic.NewQWACCertificateDiagnosticDataBuilder()
	return builder, &builder.DiagnosticDataBuilder
}

// DefaultValidationPolicyPath ports the protected getDefaultValidationPolicyPath() override.
func (v *QWACValidator) DefaultValidationPolicyPath() string {
	return qwacValidationPolicyLocation
}

// ProvideProcessExecutorInstance ports the protected provideProcessExecutorInstance() override.
func (v *QWACValidator) ProvideProcessExecutorInstance() executor.CertificateProcessExecutor {
	processExecutor := v.ProcessExecutor
	if processExecutor == nil {
		processExecutor = v.DefaultProcessExecutor()
		v.SetProcessExecutor(processExecutor)
	}
	if v.certificateToken != nil {
		processExecutor.SetCertificateId(v.IdentifierProvider.IDAsString(v.certificateToken))
	}
	return processExecutor
}

// DefaultProcessExecutor ports the public getDefaultProcessExecutor() override.
func (v *QWACValidator) DefaultProcessExecutor() executor.CertificateProcessExecutor {
	return executor.NewQWACCertificateProcessExecutor()
}
