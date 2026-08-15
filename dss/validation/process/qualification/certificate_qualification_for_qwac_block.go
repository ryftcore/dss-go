// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/CertificateQualificationForQWACBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
)

// CertificateQualificationForQWACBlock performs qualification determination
// for a QWAC certificate according to the TS 119 615 process.
type CertificateQualificationForQWACBlock struct {
	*CertificateQualificationBlock
}

// NewCertificateQualificationForQWACBlock is the default constructor. Port
// of CertificateQualificationForQWACBlock(I18nProvider, XmlConclusion, Date, CertificateWrapper, List).
func NewCertificateQualificationForQWACBlock(i18nProvider *i18n.I18nProvider, buildingBlocksConclusion *jaxb.XmlConclusion,
	validationTime time.Time, signingCertificate *diagnostic.CertificateWrapper,
	tlAnalysis []*jaxb.XmlTLAnalysis) *CertificateQualificationForQWACBlock {
	c := &CertificateQualificationForQWACBlock{
		CertificateQualificationBlock: NewCertificateQualificationBlock(i18nProvider, buildingBlocksConclusion, validationTime,
			signingCertificate, tlAnalysis),
	}
	c.InitCertificateQualificationBlock(c)
	return c
}

// CertQualificationAtIssuanceTimeBlock is the port of the overridden
// getCertQualificationAtIssuanceTimeBlock(List).
func (c *CertificateQualificationForQWACBlock) CertQualificationAtIssuanceTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTime_CERTIFICATE_ISSUANCE_TIME,
		c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}

// CertQualificationAtValidationTimeBlock is the port of the overridden
// getCertQualificationAtValidationTimeBlock(List).
func (c *CertificateQualificationForQWACBlock) CertQualificationAtValidationTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeForQWACBlock(c.I18nProvider, enumerations.ValidationTime_VALIDATION_TIME, &c.ValidationTime,
		c.SigningCertificate, acceptableServices).CertQualificationAtTimeBlock
}
