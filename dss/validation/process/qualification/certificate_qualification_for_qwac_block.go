// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/CertificateQualificationForQWACBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
)

// CertificateQualificationForQWACBlock performs qualification determination
// for a QWAC certificate according to the TS 119 615 process.
type CertificateQualificationForQWACBlock struct {
	*CertificateQualificationBlock
}

// NewCertificateQualificationForQWACBlock is the default constructor. Port
// of CertificateQualificationForQWACBlock(Provider, XmlConclusion, Date, CertificateWrapper, List).
func NewCertificateQualificationForQWACBlock(i18nProvider *i18n.Provider, buildingBlocksConclusion *jaxb.XmlConclusion,
	validationTime time.Time, signingCertificate *diagnostic.CertificateWrapper,
	tlAnalysis []*jaxb.XmlTLAnalysis) *CertificateQualificationForQWACBlock {
	c := &CertificateQualificationForQWACBlock{
		CertificateQualificationBlock: NewCertificateQualificationBlock(i18nProvider, buildingBlocksConclusion, validationTime,
			signingCertificate, tlAnalysis),
	}
	c.InitCertificateQualificationBlock(c)
	return c
}

// CertQualificationAtIssuanceTimeBlock builds a QWAC-specific
// certificate-issuance-time qualification block for the given acceptable
// services. Port of the overridden getCertQualificationAtIssuanceTimeBlock(List).
func (c *CertificateQualificationForQWACBlock) CertQualificationAtIssuanceTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTimeCertificateIssuanceTime,
		c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}

// CertQualificationAtValidationTimeBlock builds a QWAC-specific
// validation-time qualification block for the given acceptable services.
// Port of the overridden getCertQualificationAtValidationTimeBlock(List).
func (c *CertificateQualificationForQWACBlock) CertQualificationAtValidationTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlock(c.I18nProvider, enumerations.ValidationTimeValidationTime, &c.ValidationTime,
		c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}
