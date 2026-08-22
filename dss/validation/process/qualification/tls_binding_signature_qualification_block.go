// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/qwac/TLSBindingSignatureQualificationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TLSBindingSignatureQualificationBlock executes qualification determination
// process for a TLS/SSL Certificate Binding signature.
type TLSBindingSignatureQualificationBlock struct {
	*SignatureQualificationBlock

	// bbbs is the map of Basic Building Blocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// bindingSignature is the TLS Certificate Binding signature.
	bindingSignature *diagnostic.SignatureWrapper

	// websiteUrl is the URL of the website to validate the QWAC certificate
	// against.
	websiteUrl string
}

// NewTLSBindingSignatureQualificationBlock is the default constructor. Port
// of
// TLSBindingSignatureQualificationBlock(I18nProvider, Map, XmlConstraintsConclusionWithProofOfExistence, SignatureWrapper, List, String).
func NewTLSBindingSignatureQualificationBlock(i18nProvider *i18n.I18nProvider, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	etsi319102validation *jaxb.XmlConstraintsConclusionWithProofOfExistence, bindingSignature *diagnostic.SignatureWrapper,
	tlAnalysis []*jaxb.XmlTLAnalysis, websiteUrl string) *TLSBindingSignatureQualificationBlock {
	c := &TLSBindingSignatureQualificationBlock{
		SignatureQualificationBlock: NewSignatureQualificationBlock(i18nProvider, etsi319102validation,
			bindingSignature.SigningCertificate(), tlAnalysis),
		bbbs:             bbbs,
		bindingSignature: bindingSignature,
		websiteUrl:       websiteUrl,
	}
	c.InitSignatureQualificationBlock(c)
	return c
}

// InitChain initializes the chain. Port of the overridden protected void
// initChain().
func (c *TLSBindingSignatureQualificationBlock) InitChain() {
	c.SignatureQualificationBlock.InitChain()

	qwacForTLSBindingValidationBlock := NewQWACForTLSBindingCertificateValidationBlock(c.I18nProvider, c.BestSignatureTime,
		c.bindingSignature, c.SigningCertificate, c.bbbs, c.getCertificateQualification(), c.websiteUrl)
	qwacProcess := qwacForTLSBindingValidationBlock.Execute()
	c.Result.Value.QWACProcess = qwacProcess
}

// getCertificateQualification ports the private getCertificateQualification().
func (c *TLSBindingSignatureQualificationBlock) getCertificateQualification() *jaxb.XmlCertificateQualificationProcess {
	xmlCertificateQualificationProcess := &jaxb.XmlCertificateQualificationProcess{}
	xmlCertificateQualificationProcess.ValidationCertificateQualification = append(
		xmlCertificateQualificationProcess.ValidationCertificateQualification, c.Result.Value.ValidationCertificateQualification...)
	if utils.IsCollectionEmpty(c.Result.Value.ValidationCertificateQualification) {
		xmlCertificateQualificationProcess.Conclusion = getConclusionByIndication(enumerations.IndicationFailed)
		return xmlCertificateQualificationProcess
	}

	for _, certificateQualificationAtTime := range c.Result.Value.ValidationCertificateQualification {
		if !c.IsValid(&certificateQualificationAtTime.XmlConstraintsConclusionContent) {
			c.Result.SetConclusion(getConclusionFrom(certificateQualificationAtTime.Conclusion))
			return xmlCertificateQualificationProcess
		}
	}
	xmlCertificateQualificationProcess.Conclusion = getConclusionByIndication(enumerations.IndicationPassed)
	return xmlCertificateQualificationProcess
}

// getConclusionByIndication ports the private getConclusion(Indication).
func getConclusionByIndication(indication enumerations.Indication) *jaxb.XmlConclusion {
	xmlConclusion := &jaxb.XmlConclusion{}
	xmlConclusion.Indication = jaxb.IndicationValue(indication)
	return xmlConclusion
}

// getConclusionFrom ports the private getConclusion(XmlConclusion).
func getConclusionFrom(conclusion *jaxb.XmlConclusion) *jaxb.XmlConclusion {
	xmlConclusion := &jaxb.XmlConclusion{}
	xmlConclusion.Indication = conclusion.Indication
	xmlConclusion.SubIndication = conclusion.SubIndication
	xmlConclusion.Errors = append(xmlConclusion.Errors, conclusion.Errors...)
	xmlConclusion.Warnings = append(xmlConclusion.Warnings, conclusion.Warnings...)
	xmlConclusion.Infos = append(xmlConclusion.Infos, conclusion.Infos...)
	return xmlConclusion
}

// CertQualificationAtIssuanceTimeBlock gets a certificate qualification
// determination process for validation at the certificate issuance time.
// Port of the overridden getCertQualificationAtIssuanceTimeBlock(List).
func (c *TLSBindingSignatureQualificationBlock) CertQualificationAtIssuanceTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlockAtIssuanceTime(c.I18nProvider,
		enumerations.ValidationTimeCertificateIssuanceTime, c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}

// CertQualificationAtSigningTimeBlock gets a certificate qualification
// determination process for validation at the certificate signing time.
// Port of the overridden getCertQualificationAtSigningTimeBlock(List, Date).
func (c *TLSBindingSignatureQualificationBlock) CertQualificationAtSigningTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper, signingTime time.Time) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlock(c.I18nProvider, enumerations.ValidationTimeValidationTime, &signingTime,
		c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}
