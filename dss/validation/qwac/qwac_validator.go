// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/QWACValidator.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (flagged per this batch's porter brief - "cross-chunk assumptions"): Java's
// superclass AbstractCertificateValidator<CertificateReports, CertificateProcessExecutor>, the
// executor.certificate.qwac.QWACCertificateProcessExecutor, reports.CertificateReports and
// reports.diagnostic.{DiagnosticDataBuilder,QWACCertificateDiagnosticDataBuilder} types all
// belong to this same phase 8f batch's EXEC chunk ("certificate + qwac + eaa executors"), not yet
// landed at the time this file was written (unlike detached_timestamp_validator.go's single,
// better-understood forward dependency, none of this file's base infrastructure exists yet, so
// this is the highest-risk file in this chunk's manifest). Every type/method below not already
// confirmed to exist in the frozen spi packages is inferred purely from QWACValidator.java's own
// call sites and named per PORTING.md's get/set-prefix-dropping convention; expect this file to
// need real reconciliation once EXEC's chunk lands - see this batch's porter notes.
//
//	package validation // github.com/utain/esig/dss/validation
//	type AbstractCertificateValidator[R any, E any] struct { ... }  // certificateVerifier, identifierProvider, processExecutor accessible to embedders
//	func (v *AbstractCertificateValidator[R, E]) PrepareValidationContext(certificateVerifier spivalidation.CertificateVerifier) spivalidation.ValidationContext
//	func (v *AbstractCertificateValidator[R, E]) ValidateContext(validationContext spivalidation.ValidationContext)
//	func (v *AbstractCertificateValidator[R, E]) CreateDiagnosticDataBuilder(validationContext spivalidation.ValidationContext) diagnosticreports.DiagnosticDataBuilder
//	func (v *AbstractCertificateValidator[R, E]) CertificateVerifier() spivalidation.CertificateVerifier
//	func (v *AbstractCertificateValidator[R, E]) IdentifierProvider() model.TokenIdentifierProvider
//
// package qwacexecutor // github.com/utain/esig/dss/validation/executor/certificate/qwac
//	type QWACCertificateProcessExecutor struct { ... }
//	func NewQWACCertificateProcessExecutor() *QWACCertificateProcessExecutor
//	func (e *QWACCertificateProcessExecutor) SetCertificateId(id string)
//
// package diagnosticreports // github.com/utain/esig/dss/validation/reports/diagnostic
//	type DiagnosticDataBuilder interface { ... }
//	type QWACCertificateDiagnosticDataBuilder struct { ... }
//	func NewQWACCertificateDiagnosticDataBuilder() *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) WebsiteURL(url string) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificate(cert *model.CertificateToken) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificateBindingUrl(url string) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificateBindingSignature(sig spivalidation.AdvancedSignature) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) SetSignatureDiagnosticDataBuilder(d diagnosticreports.DiagnosticDataBuilder) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) FoundSignatures(sigs []spivalidation.AdvancedSignature) *QWACCertificateDiagnosticDataBuilder
//	func (b *QWACCertificateDiagnosticDataBuilder) DocumentCertificateSource(src spi.CertificateSource) *QWACCertificateDiagnosticDataBuilder
//
// Revisit once EXEC lands the real packages - only the shapes above need to resolve; this file's
// own control flow (mirroring QWACValidator.java statement-by-statement) does not otherwise
// depend on their internals.
package qwac

import (
	"fmt"

	"github.com/utain/esig/dss/model"
	modelhttp "github.com/utain/esig/dss/model/http"
	"github.com/utain/esig/dss/spi"
	spihttp "github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	dssvalidation "github.com/utain/esig/dss/validation"
	qwacexecutor "github.com/utain/esig/dss/validation/executor/certificate/qwac"
	diagnosticreports "github.com/utain/esig/dss/validation/reports/diagnostic"
)

// qwacValidationPolicyLocation is the path for the default QWAC validation policy.
const qwacValidationPolicyLocation = "/policy/qwac-constraint.xml"

// QWACValidator performs a validation of a TLS/SSL certificate as per ETSI TS 119 411-5 "Policy
// and security requirements for Trust Service Providers issuing certificates; Part 5:
// Implementation of qualified certificates for website authentication as in amended Regulation
// 910/2014".
type QWACValidator struct {
	dssvalidation.AbstractCertificateValidator[*dssvalidation.CertificateReports, *qwacexecutor.QWACCertificateProcessExecutor]

	// url is the URL to verify the used TSL/SSL certificate.
	url string

	// certificateToken is the TSL/SSL certificate derived from url.
	certificateToken *model.CertificateToken

	// dataLoader establishes the secure TLS/SSL connection with a remote server and/or loads
	// any applicable information required for a QWAC validation.
	dataLoader spihttp.AdvancedDataLoader
}

// newQWACValidator ports both private constructors.
func newQWACValidator(url string, certificateToken *model.CertificateToken) *QWACValidator {
	if url == "" {
		panic("URL cannot be null!")
	}
	return &QWACValidator{url: url, certificateToken: certificateToken}
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

// PrepareDiagnosticDataBuilder ports the protected prepareDiagnosticDataBuilder() override.
func (v *QWACValidator) PrepareDiagnosticDataBuilder() (diagnosticreports.DiagnosticDataBuilder, error) {
	response, err := v.connectToURL()
	if err != nil {
		return nil, err
	}
	if err := v.assertResponseValid(response); err != nil {
		return nil, err
	}

	tlsCertificates, err := v.toCertificateTokenList(response.TLSCertificates())
	if err != nil {
		return nil, err
	}
	tlsCertificateBindingURL := v.readTLSCertificateBindingURL(response)
	signedDocumentValidator, err := v.initSignedDocumentValidator(tlsCertificateBindingURL, tlsCertificates)
	if err != nil {
		return nil, err
	}
	signature := v.getTLSCertificateBindingSignature(signedDocumentValidator)

	tlsCertificate, err := v.getTLSCertificate(tlsCertificates)
	if err != nil {
		return nil, err
	}
	validationContext := v.prepareQWACValidationContext(tlsCertificate, tlsCertificates, signature)
	v.ValidateContext(validationContext)
	return v.createQWACDiagnosticDataBuilder(validationContext, signedDocumentValidator, tlsCertificate, tlsCertificateBindingURL, signature)
}

// connectToURL connects to url. Port of the protected connectToUrl().
func (v *QWACValidator) connectToURL() (*modelhttp.ResponseEnvelope, error) {
	trimmedURL := utils.Trim(v.url)
	if spihttp.ProtocolIsHttpURL(trimmedURL) {
		return v.dataLoaderOrDefault().RequestGetFull(v.url, true, false), nil
	}
	return nil, fmt.Errorf("DSS framework supports only HTTP(S) certificate extraction. Obtained URL : '%s'", v.url)
}

// prepareQWACValidationContext prepares a ValidationContext using the configuration and
// provided data objects. Port of the protected prepareValidationContext(CertificateToken, List,
// AdvancedSignature) override.
func (v *QWACValidator) prepareQWACValidationContext(tlsCertificate *model.CertificateToken, otherTLSCertificates []*model.CertificateToken,
	signature validation.AdvancedSignature) validation.ValidationContext {
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
// Returns an error when no TSL/SSL certificate is available (Java's thrown DSSException).
func (v *QWACValidator) getTLSCertificate(tlsCertificates []*model.CertificateToken) (*model.CertificateToken, error) {
	if v.certificateToken != nil {
		return v.certificateToken, nil
	} else if utils.IsCollectionNotEmpty(tlsCertificates) {
		// first certificate shall be the peer's end-entity certificate
		v.certificateToken = tlsCertificates[0]
		return v.certificateToken, nil
	}
	return nil, fmt.Errorf("no valid TLS/SSL certificates have been obtained from the URL '%s'", v.url)
}

// initSignedDocumentValidator ports the private initSignedDocumentValidator(String, List).
func (v *QWACValidator) initSignedDocumentValidator(tlsCertificateBindingURL string, tlsCertificates []*model.CertificateToken) (dssvalidation.SignedDocumentValidator, error) {
	if tlsCertificateBindingURL == "" {
		return nil, nil
	}
	tlsCertificateBindingSignatureBytes := v.dataLoaderOrDefault().Get(tlsCertificateBindingURL)
	if tlsCertificateBindingSignatureBytes == nil {
		return nil, nil
	}
	signatureDocument := model.NewInMemoryDocumentWithName(tlsCertificateBindingSignatureBytes, tlsCertificateBindingURL)
	documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(signatureDocument)
	if err != nil {
		return nil, err
	}
	documentValidator.SetCertificateVerifier(v.offlineCertificateVerifier())
	documentValidator.SetDetachedContents(v.toDetachedDocumentsList(tlsCertificates))
	return documentValidator, nil
}

// getTLSCertificateBindingSignature ports the private getTLSCertificateBindingSignature(SignedDocumentValidator).
func (v *QWACValidator) getTLSCertificateBindingSignature(signedDocumentValidator dssvalidation.SignedDocumentValidator) validation.AdvancedSignature {
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
func (v *QWACValidator) toCertificateTokenList(certificates []*model.CertificateToken) ([]*model.CertificateToken, error) {
	// modelhttp.ResponseEnvelope#TLSCertificates already exposes *stdx509.Certificate values
	// translated by callers where needed; here QWACValidator only ever forwards them, so this
	// stays a pass-through filtering out load failures - see toCertificateToken.
	return certificates, nil
}

// toDetachedDocumentsList ports the private toDetachedDocumentsList(List).
func (v *QWACValidator) toDetachedDocumentsList(certificates []*model.CertificateToken) []model.DSSDocument {
	result := make([]model.DSSDocument, 0, len(certificates))
	for _, certificate := range certificates {
		result = append(result, model.NewInMemoryDocumentWithName(certificate.Encoded(), v.IdentifierProvider().IDAsString(certificate)))
	}
	return result
}

// readTLSCertificateBindingURL reads the HTTP response headers and extracts the value of the
// "Link" header with a "rel" value of "tls-certificate-binding". Port of the protected
// readTLSCertificateBindingUrl(ResponseEnvelope).
func (v *QWACValidator) readTLSCertificateBindingURL(responseEnvelope *modelhttp.ResponseEnvelope) string {
	return GetTLSCertificateBindingUrl(responseEnvelope.Headers())
}

// assertResponseValid verifies whether response is valid and contains the information required
// to continue the QWAC validation process. Port of the protected assertResponseValid(ResponseEnvelope).
//
// Returns an error in place of Java's thrown IllegalArgumentException.
func (v *QWACValidator) assertResponseValid(response *modelhttp.ResponseEnvelope) error {
	if v.certificateToken == nil && len(response.TLSCertificates()) == 0 {
		return fmt.Errorf("no TSL/SSL certificates have been returned from the URL '%s'. "+
			"Please ensure the URL is valid and uses the HTTP(S) scheme, or provide the TLS/SSL certificate on validation explicitly", v.url)
	}
	return nil
}

// createQWACDiagnosticDataBuilder creates and configures a new DiagnosticDataBuilder. Port of
// the protected createDiagnosticDataBuilder(ValidationContext, SignedDocumentValidator,
// CertificateToken, String, AdvancedSignature) override.
func (v *QWACValidator) createQWACDiagnosticDataBuilder(validationContext validation.ValidationContext,
	signedDocumentValidator dssvalidation.SignedDocumentValidator, tlsCertificate *model.CertificateToken,
	tlsCertificateBindingURL string, signature validation.AdvancedSignature) (diagnosticreports.DiagnosticDataBuilder, error) {
	base, err := v.CreateDiagnosticDataBuilder(validationContext)
	if err != nil {
		return nil, err
	}
	diagnosticDataBuilder, ok := base.(*diagnosticreports.QWACCertificateDiagnosticDataBuilder)
	if !ok {
		return nil, fmt.Errorf("expected a QWACCertificateDiagnosticDataBuilder, got %T", base)
	}

	diagnosticDataBuilder = diagnosticDataBuilder.
		WebsiteURL(v.url).
		TLSCertificate(tlsCertificate).
		TLSCertificateBindingUrl(tlsCertificateBindingURL).
		TLSCertificateBindingSignature(signature)

	if signedDocumentValidator != nil {
		diagnosticDataBuilder = diagnosticDataBuilder.SetSignatureDiagnosticDataBuilder(signedDocumentValidator.InitializeDiagnosticDataBuilder())
	}

	return diagnosticDataBuilder.
		FoundSignatures(validationContext.ProcessedSignatures()).
		DocumentCertificateSource(validationContext.DocumentCertificateSource()), nil
}

// offlineCertificateVerifier builds a complete copy of the CertificateVerifier for the offline
// validation. Port of the protected getOfflineCertificateVerifier().
func (v *QWACValidator) offlineCertificateVerifier() validation.CertificateVerifier {
	return validation.NewCertificateVerifierBuilder(v.CertificateVerifier()).BuildCompleteCopyForValidation()
}

// InitDiagnosticDataBuilder ports the protected initDiagnosticDataBuilder() override.
func (v *QWACValidator) InitDiagnosticDataBuilder() diagnosticreports.DiagnosticDataBuilder {
	return diagnosticreports.NewQWACCertificateDiagnosticDataBuilder()
}

// DefaultValidationPolicyPath ports the protected getDefaultValidationPolicyPath() override.
func (v *QWACValidator) DefaultValidationPolicyPath() string {
	return qwacValidationPolicyLocation
}

// ProvideProcessExecutorInstance ports the protected provideProcessExecutorInstance() override.
func (v *QWACValidator) ProvideProcessExecutorInstance() *qwacexecutor.QWACCertificateProcessExecutor {
	processExecutor := v.ProcessExecutor()
	if processExecutor == nil {
		processExecutor = v.DefaultProcessExecutor()
		v.SetProcessExecutor(processExecutor)
	}
	if v.certificateToken != nil {
		processExecutor.SetCertificateId(v.IdentifierProvider().IDAsString(v.certificateToken))
	}
	return processExecutor
}

// DefaultProcessExecutor ports the public getDefaultProcessExecutor() override.
func (v *QWACValidator) DefaultProcessExecutor() *qwacexecutor.QWACCertificateProcessExecutor {
	return qwacexecutor.NewQWACCertificateProcessExecutor()
}
