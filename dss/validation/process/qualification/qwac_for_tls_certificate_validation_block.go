// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/QWACForTLSCertificateValidationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/validation/process"
)

// QWACForTLSCertificateValidationBlock runs a validation process for QWAC
// as per ETSI TS 119 411-5 for a direct TLS/SSL certificate.
type QWACForTLSCertificateValidationBlock struct {
	*process.ChainBase[*jaxb.XmlQWACProcess]

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// certificate is the certificate to determine qualification for.
	certificate *diagnostic.CertificateWrapper

	// bbbs is the map of Basic Building Blocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// certificateQualification is the qualification status of the
	// certificate.
	certificateQualification *jaxb.XmlCertificateQualificationProcess

	// bindingCertificateProfile is the QWAC profile of the binding
	// certificate.
	bindingCertificateProfile enumerations.QWACProfile

	// websiteUrl is the URL of the website to validate the QWAC certificate
	// against.
	websiteUrl string
}

// NewQWACForTLSCertificateValidationBlock is the default constructor. Port
// of
// QWACForTLSCertificateValidationBlock(I18nProvider, DiagnosticData, CertificateWrapper, Map, XmlCertificateQualificationProcess, QWACProfile, String).
func NewQWACForTLSCertificateValidationBlock(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	certificate *diagnostic.CertificateWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	certificateQualification *jaxb.XmlCertificateQualificationProcess, bindingCertificateProfile enumerations.QWACProfile,
	websiteUrl string) *QWACForTLSCertificateValidationBlock {
	xmlResult := &jaxb.XmlQWACProcess{}
	c := &QWACForTLSCertificateValidationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		certificate:               certificate,
		diagnosticData:            diagnosticData,
		bbbs:                      bbbs,
		certificateQualification:  certificateQualification,
		bindingCertificateProfile: bindingCertificateProfile,
		websiteUrl:                websiteUrl,
	}
	c.Result.Value.Id = certificate.Id()
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *QWACForTLSCertificateValidationBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_QWAC_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *QWACForTLSCertificateValidationBlock) InitChain() {

	certBBB, ok := c.bbbs[c.certificate.Id()]
	if !ok || certBBB == nil {
		panic(fmt.Sprintf("The certificate basic validation process shall be performed! "+
			"No BasicBuildingBlock found for a certificate with Id '%s'", c.certificate.Id()))
	}

	validationDate := c.diagnosticData.ValidationDate()
	qwac1Process := NewQWAC1ValidationProcessBlock(c.I18nProvider, timeOrZero(validationDate), c.certificate,
		certBBB.Conclusion, c.certificateQualification, c.websiteUrl)
	qwac1ValidationResult := qwac1Process.Execute()
	c.Result.Value.ValidationQWACProcess = append(c.Result.Value.ValidationQWACProcess, qwac1ValidationResult)

	// Java's getTokenValidationConclusion(TokenProxy) is called polymorphically
	// on a CertificateWrapper and (possibly null) SignatureWrapper; each is
	// resolved to its concrete type before the shared TokenProxy-typed helper
	// runs, so a nil *SignatureWrapper never gets boxed into a non-nil
	// TokenProxy interface value (the classic Go "typed nil in interface" trap).
	var bindingSignatureConclusion *jaxb.XmlConclusion
	if bindingSignature := c.diagnosticData.TLSCertificateBindingSignature(); bindingSignature != nil {
		bindingSignatureConclusion = c.getTokenValidationConclusion(bindingSignature)
	}
	tlsCertificateProcess := NewTLSCertificateSupportedByQWAC2ValidationProcessBlock(c.I18nProvider, c.diagnosticData, c.certificate,
		c.getTokenValidationConclusion(c.certificate), bindingSignatureConclusion,
		c.bindingCertificateProfile, c.websiteUrl)
	tlsCertificateValidationResult := tlsCertificateProcess.Execute()
	c.Result.Value.ValidationQWACProcess = append(c.Result.Value.ValidationQWACProcess, tlsCertificateValidationResult)

	item := c.qwacValidation(qwac1ValidationResult, tlsCertificateValidationResult)
	c.FirstItem = item

	if c.IsValid(&qwac1ValidationResult.XmlConstraintsConclusionContent) {
		v := jaxb.QWACProfileValue(qwac1Process.QWACProfile())
		c.Result.Value.QWACType = &v
	} else if c.IsValid(&tlsCertificateValidationResult.XmlConstraintsConclusionContent) {
		v := jaxb.QWACProfileValue(tlsCertificateProcess.QWACProfile())
		c.Result.Value.QWACType = &v
	} else {
		v := jaxb.QWACProfileValue(enumerations.QWACProfile_NOT_QWAC)
		c.Result.Value.QWACType = &v
	}
}

func (c *QWACForTLSCertificateValidationBlock) qwacValidation(qwacValidationProcesses ...*jaxb.XmlValidationQWACProcess) process.ChainItem[*jaxb.XmlQWACProcess] {
	return NewQWACValidationResultCheck(c.I18nProvider, c.Result, qwacValidationProcesses, c.FailLevelRule())
}

// getTokenValidationConclusion ports the private getTokenValidationConclusion(TokenProxy).
// Callers only ever pass a genuinely non-nil TokenProxy (see InitChain); a
// nil interface value here would indicate a porting bug, not Java's null
// short-circuit, so it is not defended against.
func (c *QWACForTLSCertificateValidationBlock) getTokenValidationConclusion(token diagnostic.TokenProxy) *jaxb.XmlConclusion {
	bbb, ok := c.bbbs[token.Id()]
	if !ok || bbb == nil {
		panic("The Basic validation shall be performed!")
	}
	return bbb.Conclusion
}

// timeOrZero dereferences a possibly-nil *time.Time to its zero value; a nil
// DiagnosticData#getValidationDate() should not occur for a schema-valid
// dump (ValidationDate is effectively required), so the zero-time fallback
// only avoids a nil-pointer panic and is not expected to be exercised.
func timeOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
