// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/EAAQualificationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAQualificationBlock is used to verify qualification status of a
// signature used to create the EAA.
type EAAQualificationBlock struct {
	*process.ChainBase[*jaxb.XmlValidationEAAQualification]

	// eaa is the EAA to be validated.
	eaa *diagnostic.EAAWrapper

	// eaaConclusion is the conclusion of EAA validation.
	eaaConclusion *jaxb.XmlConclusion

	// signatureMap is a map of signature validation processes.
	signatureMap map[string]*jaxb.XmlSignature

	// tlAnalysis is the list of all TL analyses.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// loteAnalysis is the list of List of Trusted Entities validations.
	loteAnalysis []*jaxb.XmlLoTEAnalysis

	// currentTime is the validation time.
	currentTime time.Time
}

// NewEAAQualificationBlock is the default constructor. Port of
// EAAQualificationBlock(I18nProvider, EAAWrapper, XmlConclusion, Map, List, List, Date).
func NewEAAQualificationBlock(i18nProvider *i18n.I18nProvider, eaa *diagnostic.EAAWrapper, eaaConclusion *jaxb.XmlConclusion,
	signatureMap map[string]*jaxb.XmlSignature, tlAnalysis []*jaxb.XmlTLAnalysis, loteAnalysis []*jaxb.XmlLoTEAnalysis,
	currentTime time.Time) *EAAQualificationBlock {
	xmlResult := &jaxb.XmlValidationEAAQualification{}
	c := &EAAQualificationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		eaa:           eaa,
		eaaConclusion: eaaConclusion,
		signatureMap:  signatureMap,
		tlAnalysis:    tlAnalysis,
		loteAnalysis:  loteAnalysis,
		currentTime:   currentTime,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *EAAQualificationBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_EAA_QUALIFICATION
}

// InitChain initializes the chain. Port of initChain().
func (c *EAAQualificationBlock) InitChain() {

	var eaaQualificationProcess *jaxb.XmlValidationEAAQualificationProcess
	var pidQualificationProcess *jaxb.XmlValidationPIDQualificationProcess

	if utils.CollectionSize(c.eaa.EAASignatures()) == 1 {

		signingCertificate := c.getSigningCertificate()

		item := c.isTrustAnchorListReachedForCertificateChain(signingCertificate)
		c.FirstItem = item

		eaaQualificationProcessBlock := NewEAAQualificationProcessBlock(c.I18nProvider, c.eaa, c.eaaConclusion,
			c.signatureMap, c.tlAnalysis, c.currentTime)
		eaaQualificationProcess = eaaQualificationProcessBlock.Execute()
		c.Result.Value.ValidationEAAQualificationProcess = eaaQualificationProcess

		pidQualificationProcessBlock := NewPIDQualificationProcessBlock(c.I18nProvider, c.eaa, c.eaaConclusion,
			c.loteAnalysis, c.currentTime)
		pidQualificationProcess = pidQualificationProcessBlock.Execute()
		c.Result.Value.ValidationPIDQualificationProcess = pidQualificationProcess

		item = item.SetNextItem(c.eaaQualificationProcessConclusiveCheck(eaaQualificationProcess, pidQualificationProcess)) //nolint:staticcheck // mirrors upstream EAAQualificationBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}

	c.determineFinalQualification(eaaQualificationProcess, pidQualificationProcess)
}

func (c *EAAQualificationBlock) isTrustAnchorListReachedForCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualification] {
	return NewTrustAnchorListReachedForCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, c.FailLevelRule())
}

func (c *EAAQualificationBlock) eaaQualificationProcessConclusiveCheck(
	eaaQualificationProcess *jaxb.XmlValidationEAAQualificationProcess,
	pidQualificationProcess *jaxb.XmlValidationPIDQualificationProcess) process.ChainItem[*jaxb.XmlValidationEAAQualification] {
	conclusions := []*jaxb.XmlConstraintsConclusionContent{
		&eaaQualificationProcess.XmlConstraintsConclusionContent,
		&pidQualificationProcess.XmlConstraintsConclusionContent,
	}
	return NewEAAQualificationProcessConclusiveCheck(c.I18nProvider, c.Result, conclusions, c.FailLevelRule())
}

// getSigningCertificate ports the private getSigningCertificate().
func (c *EAAQualificationBlock) getSigningCertificate() *diagnostic.CertificateWrapper {
	eaaSignature := c.eaa.EAASignatures()[0]
	return eaaSignature.SigningCertificate()
}

// determineFinalQualification ports the private
// determineFinalQualification(XmlValidationEAAQualificationProcess, XmlValidationPIDQualificationProcess).
func (c *EAAQualificationBlock) determineFinalQualification(eaaQualificationProcess *jaxb.XmlValidationEAAQualificationProcess,
	pidQualificationProcess *jaxb.XmlValidationPIDQualificationProcess) {
	eaaQualification := enumerations.EAAQualificationNA
	if eaaQualificationProcess != nil {
		eaaQualification = eaaQualificationProcess.EAAQualification.EAAQualification()
	}
	if enumerations.EAAQualificationNA != eaaQualification {
		c.Result.Value.EAAQualification = append(c.Result.Value.EAAQualification, jaxb.EAAQualificationValue(eaaQualification))
	}
	pidQualification := enumerations.EAAQualificationNA
	if pidQualificationProcess != nil {
		pidQualification = pidQualificationProcess.EAAQualification.EAAQualification()
	}
	if (enumerations.EAAQualificationPID == pidQualification || enumerations.EAAQualificationIndeterminatePID == pidQualification) &&
		pidQualification != eaaQualification {
		c.Result.Value.EAAQualification = append(c.Result.Value.EAAQualification, jaxb.EAAQualificationValue(pidQualification))
	} else if (enumerations.EAAQualificationUnknown == pidQualification || enumerations.EAAQualificationIndeterminateUnknown == pidQualification) &&
		enumerations.EAAQualificationNA == eaaQualification {
		c.Result.Value.EAAQualification = append(c.Result.Value.EAAQualification, jaxb.EAAQualificationValue(pidQualification))
	}
	if utils.IsCollectionEmpty(c.Result.Value.EAAQualification) {
		c.Result.Value.EAAQualification = append(c.Result.Value.EAAQualification, jaxb.EAAQualificationValue(enumerations.EAAQualificationNA))
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden protected void
// collectAdditionalMessages(XmlConclusion).
func (c *EAAQualificationBlock) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	signingCertificate := c.getSigningCertificate()
	if signingCertificate != nil && (signingCertificate.IsTrustedListReached() || signingCertificate.IsListOfTrustedEntitiesReached()) {
		if signingCertificate.IsTrustedListReached() {
			eaaQualificationProcess := c.Result.Value.ValidationEAAQualificationProcess
			c.CollectAllMessages(conclusion, eaaQualificationProcess.Conclusion)
		}
		if signingCertificate.IsListOfTrustedEntitiesReached() {
			pidQualificationProcess := c.Result.Value.ValidationPIDQualificationProcess
			c.CollectAllMessages(conclusion, pidQualificationProcess.Conclusion)
		}
	}
}
