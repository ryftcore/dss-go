// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/QWACForTLSBindingCertificateValidationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QWACForTLSBindingCertificateValidationBlock runs a validation process for
// QWAC as per ETSI TS 119 411-5 for a TLS Binding certificate (e.g. 2-QWAC
// certificate of the TLS Certificate Binding JAdES signature).
type QWACForTLSBindingCertificateValidationBlock struct {
	*process.ChainBase[*jaxb.XmlQWACProcess]

	// bindingSignature is the TLS Certificate Binding signature.
	bindingSignature *diagnostic.SignatureWrapper

	// certificate is the certificate to determine qualification for.
	certificate *diagnostic.CertificateWrapper

	// validationTime is the validation time.
	validationTime time.Time

	// bbbs is the map of Basic Building Blocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// certificateQualification is the qualification status of the
	// certificate.
	certificateQualification *jaxb.XmlCertificateQualificationProcess

	// websiteUrl is the URL of the website to validate the QWAC certificate
	// against.
	websiteUrl string
}

// NewQWACForTLSBindingCertificateValidationBlock is the default constructor.
// Port of
// QWACForTLSBindingCertificateValidationBlock(I18nProvider, Date, SignatureWrapper, CertificateWrapper, Map, XmlCertificateQualificationProcess, String).
func NewQWACForTLSBindingCertificateValidationBlock(i18nProvider *i18n.I18nProvider, validationTime time.Time,
	bindingSignature *diagnostic.SignatureWrapper, certificate *diagnostic.CertificateWrapper,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, certificateQualification *jaxb.XmlCertificateQualificationProcess,
	websiteUrl string) *QWACForTLSBindingCertificateValidationBlock {
	xmlResult := &jaxb.XmlQWACProcess{}
	c := &QWACForTLSBindingCertificateValidationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		bindingSignature:         bindingSignature,
		certificate:              certificate,
		validationTime:           validationTime,
		bbbs:                     bbbs,
		certificateQualification: certificateQualification,
		websiteUrl:               websiteUrl,
	}
	c.Result.Value.Id = certificate.Id()
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *QWACForTLSBindingCertificateValidationBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_QWAC_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *QWACForTLSBindingCertificateValidationBlock) InitChain() {

	xmlConclusion := c.getSigningCertificateValidationProcessConclusion()
	qwac2Process := NewQWAC2ValidationProcessBlock(c.I18nProvider, c.validationTime, c.certificate, xmlConclusion,
		c.certificateQualification, c.websiteUrl)
	qwac2ValidationResult := qwac2Process.Execute()
	c.Result.Value.ValidationQWACProcess = append(c.Result.Value.ValidationQWACProcess, qwac2ValidationResult)

	item := c.qwacValidation(qwac2ValidationResult)
	c.FirstItem = item

	if c.IsValid(&qwac2ValidationResult.XmlConstraintsConclusionContent) {
		v := jaxb.QWACProfileValue(qwac2Process.QWACProfile())
		c.Result.Value.QWACType = &v
	} else {
		v := jaxb.QWACProfileValue(enumerations.QWACProfileNotQWAC)
		c.Result.Value.QWACType = &v
	}
}

func (c *QWACForTLSBindingCertificateValidationBlock) qwacValidation(qwacValidationProcesses ...*jaxb.XmlValidationQWACProcess) process.ChainItem[*jaxb.XmlQWACProcess] {
	return NewQWAC2ValidationResultCheck(c.I18nProvider, c.Result, qwacValidationProcesses, c.FailLevelRule())
}

// getSigningCertificateValidationProcessConclusion ports the private
// getSigningCertificateValidationProcessConclusion().
func (c *QWACForTLSBindingCertificateValidationBlock) getSigningCertificateValidationProcessConclusion() *jaxb.XmlConclusion {
	signatureBBB, ok := c.bbbs[c.bindingSignature.Id()]
	if !ok || signatureBBB == nil {
		panic(fmt.Sprintf("The signature basic validation process shall be performed! "+
			"No BasicBuildingBlock found for a signature with Id '%s'", c.certificate.Id()))
	}

	xcv := signatureBBB.XCV
	for _, subXCV := range xcv.SubXCV {
		if c.certificate.Id() == subXCV.Id {
			return subXCV.Conclusion
		}
	}
	return xcv.Conclusion
}
