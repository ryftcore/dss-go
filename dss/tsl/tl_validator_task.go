// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/validation/TLValidatorTask.java (DSS 6.5.RC1).
package tsl

import (
	"bytes"
	"embed"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	spipolicy "github.com/ryftcore/dss-go/dss/spi/policy"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	spiexecutor "github.com/ryftcore/dss-go/dss/spi/validation/executor"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	"github.com/ryftcore/dss-go/dss/validation/reports"
	"github.com/ryftcore/dss-go/dss/xades"
	xadesdefinition "github.com/ryftcore/dss-go/dss/xades/definition"
)

// trustedListValidationPolicyLocation is the path for a LOTL/TL validation policy. Port of the
// private TRUSTED_LIST_VALIDATION_POLICY_LOCATION constant; Java reaches it with
// TLValidatorTask.class.getResourceAsStream("/policy/tsl-constraint.xml"), Go with the embedded
// copy below.
const trustedListValidationPolicyLocation = "resources/policy/tsl-constraint.xml"

// trustedListValidationPolicyResource embeds dss-tsl-validation's
// src/main/resources/policy/tsl-constraint.xml, which upstream reaches through the classpath.
// The file is a byte-for-byte copy.
//
//go:embed resources/policy/tsl-constraint.xml
var trustedListValidationPolicyResource embed.FS

// TLValidatorTask validates a TL or LOTL.
//
// It implements eu.europa.esig.dss.validation.job.validation.ValidationTask (Go package
// dss/validation/job), whose single member is Supplier#get; Go satisfies interfaces
// structurally, so that package is deliberately NOT imported here.
type TLValidatorTask struct {
	// trustedList is the Trusted List document to validate.
	trustedList model.DSSDocument

	// certificateSource is the certificate source to use.
	certificateSource spi.CertificateSource
}

// NewTLValidatorTask is the constructor used to instantiate a validator for a trusted list,
// given the DSSDocument with a trusted list and a certificate source with the allowed
// certificates to sign this TL. Port of TLValidatorTask(DSSDocument, CertificateSource).
//
// Panics with the Java messages when either argument is nil (Objects.requireNonNull).
func NewTLValidatorTask(trustedList model.DSSDocument, certificateSource spi.CertificateSource) *TLValidatorTask {
	if trustedList == nil {
		panic("The document is null")
	}
	if certificateSource == nil {
		panic("The certificate source is null")
	}
	return &TLValidatorTask{trustedList: trustedList, certificateSource: certificateSource}
}

// Get performs the validation and returns its result. Port of the get() override
// (Supplier#get); Java's DSSException becomes a returned error, per PORTING.md.
func (t *TLValidatorTask) Get() (*TLValidationResult, error) {
	validationReports, err := t.validateTL()
	if err != nil {
		return nil, err
	}
	return t.fillResult(validationReports)
}

// validateTL wires the frozen XAdES validator exactly as upstream's private validateTL() does.
func (t *TLValidatorTask) validateTL() (*reports.Reports, error) {
	certificateVerifier := spivalidation.NewCommonCertificateVerifierSimple(true)
	trustedCertSource, err := t.buildTrustedCertificateSource(t.certificateSource)
	if err != nil {
		return nil, err
	}
	certificateVerifier.SetTrustedCertSources(trustedCertSource)

	// To increase the security: the default XAdESPaths is used.
	xmlDocumentValidator, err := xades.NewXMLDocumentValidatorWithPathHolders(t.trustedList,
		[]xadesdefinition.XAdESPath{xadesdefinition.NewXAdES132Path()})
	if err != nil {
		return nil, err
	}

	xmlDocumentValidator.SetCertificateVerifier(certificateVerifier)
	xmlDocumentValidator.SetTokenExtractionStrategy(enumerations.TokenExtractionStrategyExtractCertificatesOnly)
	xmlDocumentValidator.SetEnableEtsiValidationReport(false)                            // Ignore ETSI VR
	xmlDocumentValidator.SetValidationLevel(enumerations.ValidationLevelBasicSignatures) // Timestamps,... are ignored
	// Only need to validate against the trusted certificate source
	xmlDocumentValidator.SetValidationContextExecutor(spiexecutor.SkipValidationContextExecutorInstance)
	xmlDocumentValidator.SetSignaturePolicyProvider(spipolicy.NewSignaturePolicyProvider()) // ignore signature policy loading

	validationPolicy, err := t.trustedListValidationPolicy()
	if err != nil {
		return nil, err
	}
	return xmlDocumentValidator.ValidateDocumentWithValidationPolicy(validationPolicy)
}

// fillResult maps the reports onto a TLValidationResult. Port of the private
// fillResult(Reports).
func (t *TLValidatorTask) fillResult(validationReports *reports.Reports) (*TLValidationResult, error) {
	simpleReport := validationReports.GetSimpleReport()
	if simpleReport.GetSignaturesCount() != 1 {
		return nil, model.NewDSSError(fmt.Sprintf(
			"Number of signatures must be equal to 1 (currently : %d)", simpleReport.GetSignaturesCount()))
	}

	indication := simpleReport.GetIndication(simpleReport.GetFirstSignatureId())
	subIndication := simpleReport.GetSubIndication(simpleReport.GetFirstSignatureId())

	diagnosticData := validationReports.GetDiagnosticData()
	signatureWrapper := diagnosticData.SignatureById(diagnosticData.FirstSignatureId())
	var signingTime time.Time
	if claimedSigningTime := signatureWrapper.ClaimedSigningTime(); claimedSigningTime != nil {
		signingTime = *claimedSigningTime
	}
	signingCertificateWrapper := signatureWrapper.SigningCertificate()
	var signingCertificate *model.CertificateToken
	if signingCertificateWrapper != nil {
		var err error
		signingCertificate, err = spi.DSSUtilsLoadCertificateFromBinary(signingCertificateWrapper.Binaries())
		if err != nil {
			return nil, err
		}
	}

	return NewTLValidationResult(indication, subIndication, signingTime, signingCertificate, t.certificateSource), nil
}

// buildTrustedCertificateSource imports the announced certificates into a trusted source. Port
// of the private buildTrustedCertificateSource(CertificateSource).
func (t *TLValidatorTask) buildTrustedCertificateSource(certificateSource spi.CertificateSource) (spi.TrustedCertificateSource, error) {
	commonTrustedCertificateSource := spi.NewCommonTrustedCertificateSource()
	commonTrustedCertificateSource.ImportAsTrusted(certificateSource)
	return commonTrustedCertificateSource, nil
}

// trustedListValidationPolicy loads the trusted-list validation policy. Port of the private
// getTrustedListValidationPolicy(); Java's DSSException("Unable to load the validation policy
// for trusted list") becomes a returned error, per PORTING.md.
func (t *TLValidatorTask) trustedListValidationPolicy() (policy modelpolicy.ValidationPolicy, err error) {
	// The Go loader signals a malformed policy with a panic rather than a returned error, so
	// the Java try/catch is reproduced with a recover.
	defer func() {
		if recovered := recover(); recovered != nil {
			policy = nil
			err = model.NewDSSErrorMessageCause("Unable to load the validation policy for trusted list",
				fmt.Errorf("%v", recovered))
		}
	}()
	data, readErr := trustedListValidationPolicyResource.ReadFile(trustedListValidationPolicyLocation)
	if readErr != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to load the validation policy for trusted list", readErr)
	}
	return validationpolicy.FromValidationPolicyReader(bytes.NewReader(data)).Create(), nil
}
