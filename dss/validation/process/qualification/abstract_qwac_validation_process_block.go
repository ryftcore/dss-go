// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/AbstractQWACValidationProcessBlock.java (DSS 6.5.RC1).
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

// AbstractQWACValidationProcessBlock contains common methods for performing
// a QWAC certificate validation.
type AbstractQWACValidationProcessBlock struct {
	*process.ChainBase[*jaxb.XmlValidationQWACProcess]

	// validationTime is the validation time.
	validationTime time.Time

	// certificate is the certificate to determine qualification for.
	certificate *diagnostic.CertificateWrapper

	// buildingBlocksConclusion is the certificate's BasicBuildingBlock's
	// conclusion.
	buildingBlocksConclusion *jaxb.XmlConclusion

	// certificateQualification is the qualification validation process of
	// the certificate.
	certificateQualification *jaxb.XmlCertificateQualificationProcess

	// websiteUrl is the URL of the website to validate the QWAC certificate
	// against.
	websiteUrl string

	overrides AbstractQWACValidationProcessBlockOverrides
}

// AbstractQWACValidationProcessBlockOverrides captures the member Java's
// AbstractQWACValidationProcessBlock treats virtually: the abstract
// getQWACProfile(), self-called from buildChainTitle()/initChain().
// QWAC1ValidationProcessBlock/QWAC2ValidationProcessBlock implement it and
// register themselves via InitAbstractQWACValidationProcessBlock.
type AbstractQWACValidationProcessBlockOverrides interface {
	// QWACProfile gets the current QWAC profile. Port of the abstract public
	// QWACProfile getQWACProfile().
	QWACProfile() enumerations.QWACProfile
}

// InitAbstractQWACValidationProcessBlock wires the shared state; called by
// the concrete constructor before InitChainBase. Port of the common
// constructor
// AbstractQWACValidationProcessBlock(I18nProvider, Date, CertificateWrapper, XmlConclusion, XmlCertificateQualificationProcess, String).
func (c *AbstractQWACValidationProcessBlock) InitAbstractQWACValidationProcessBlock(i18nProvider *i18n.I18nProvider,
	validationTime time.Time, certificate *diagnostic.CertificateWrapper, buildingBlocksConclusion *jaxb.XmlConclusion,
	certificateQualification *jaxb.XmlCertificateQualificationProcess, websiteUrl string,
	overrides AbstractQWACValidationProcessBlockOverrides) {
	xmlResult := &jaxb.XmlValidationQWACProcess{}
	c.ChainBase = process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
		&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs))

	c.Result.Value.Id = certificate.Id()

	c.validationTime = validationTime
	c.certificate = certificate
	c.buildingBlocksConclusion = buildingBlocksConclusion
	c.certificateQualification = certificateQualification
	c.websiteUrl = websiteUrl
	c.overrides = overrides
}

// BuildChainTitle builds the chain title. Port of buildChainTitle().
func (c *AbstractQWACValidationProcessBlock) BuildChainTitle() string {
	message := i18n.MessageTag_QWAC_VALIDATION_PROFILE
	param, err := process.GetQWACValidationMessageTag(c.overrides.QWACProfile())
	if err != nil {
		panic(err)
	}
	return c.I18nProvider.GetMessage(message, param)
}

// InitChain initializes the chain. Port of initChain().
//
// ETSI TS 119 411-5 "Policy and security requirements for Trust Service
// Providers issuing certificates; Part 5: Implementation of qualified
// certificates for website authentication as in amended Regulation
// 910/2014"
//
// 6.1.2 Validation of QWACs
//
// For 1-QWACs and 2-QWACs, validation shall include:
//
//  1. that the QWAC includes QCStatements as specified in clause 4.2 of ETSI
//     EN 319 412-4 [4] and the appropriate Policy OID specified in ETSI EN
//     319 411-2 [3];
//
//  2. that the QWAC chains back through appropriate & valid digital
//     signatures to an issuer on the EU Trusted List which is authorized to
//     issue Qualified Certificates for Website Authentication as specified
//     in ETSI TS 119 615 [1];
func (c *AbstractQWACValidationProcessBlock) InitChain() {
	item := c.qwacCertificatePolicy()
	c.FirstItem = item

	item = item.SetNextItem(c.certificateQualificationConclusive())

	if c.certificateQualification != nil && utils.IsCollectionNotEmpty(c.certificateQualification.ValidationCertificateQualification) {
		for _, certQual := range c.certificateQualification.ValidationCertificateQualification {
			item = item.SetNextItem(c.certificateForWSAAtTime(certQual))
		}
	}

	// 3) that the QWAC's validity period covers the current date and time;
	item = item.SetNextItem(c.qwacValidityPeriod())

	// 4) that the website domain name in question appears in the QWAC's
	// subject alternative name(s); and
	item = item.SetNextItem(c.qwacDomainName())

	// 5) that the QWAC's certificate profile conforms with:
	// a) For a 1-QWAC, clause 4.1.2 of the present document; or
	// b) For a 2-QWAC, clause 4.2.2 of the present document.
	// TODO: include certificate policy requirements ?

	// 4.2.2 Certificate Profile Requirements
	// The 2-QWAC certificate shall be issued in accordance with ETSI EN 319
	// 412-4 [4] for the relevant certificate policy as identified in clause
	// 4.2.1 of the present document, except as described below:
	// - the extKeyUsage value shall only assert the extendedKeyUsage
	// purpose of id-kp-tls-binding as specified in Annex A.
	if enumerations.QWACProfileQWAC2 == c.overrides.QWACProfile() {
		item = item.SetNextItem(c.qwac2ExtKeyUsage())
	}

	// The web browser may also perform further checks on the security and
	// authenticity of the QWAC as appropriate (e.g. for checking revocation
	// status).
	item = item.SetNextItem(c.isAcceptableBuildingBlockConclusion()) //nolint:staticcheck // mirrors upstream AbstractQWACValidationProcessBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
}

func (c *AbstractQWACValidationProcessBlock) certificateQualificationConclusive() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewCertificateQualificationConclusiveCheck(c.I18nProvider, c.Result, c.certificateQualification, c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) certificateForWSAAtTime(certQual *jaxb.XmlValidationCertificateQualification) process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQualifiedCertificateForWSAAtTimeCheck(c.I18nProvider, c.Result, certQual, c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) qwacCertificatePolicy() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQWACCertificatePolicyCheck(c.I18nProvider, c.Result, c.certificate, c.overrides.QWACProfile(), c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) qwacValidityPeriod() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQWACValidityPeriodCheck(c.I18nProvider, c.Result, c.certificate, c.validationTime, c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) qwacDomainName() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQWACDomainNameCheck(c.I18nProvider, c.Result, c.certificate, c.websiteUrl, c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) qwac2ExtKeyUsage() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQWAC2ExtKeyUsageCheck(c.I18nProvider, c.Result, c.certificate, c.FailLevelRule())
}

func (c *AbstractQWACValidationProcessBlock) isAcceptableBuildingBlockConclusion() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewAcceptableBuildingBlockConclusionCheck(c.I18nProvider, c.Result, c.buildingBlocksConclusion, c.FailLevelRule())
}
