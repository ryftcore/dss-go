// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/CertQualificationAtTimeForQWACBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
)

// CertQualificationAtTimeForQWACBlock determines certificates qualification
// status as per ETSI TS 119 615 at the given time for a QWAC.
type CertQualificationAtTimeForQWACBlock struct {
	*CertQualificationAtTimeBlock
}

// NewCertQualificationAtTimeForQWACBlockAtIssuanceTime is the constructor to
// instantiate the validation at the certificate's issuance time. Port of
// CertQualificationAtTimeForQWACBlock(I18nProvider, ValidationTime, CertificateWrapper, List).
func NewCertQualificationAtTimeForQWACBlockAtIssuanceTime(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime,
	signingCertificate *diagnostic.CertificateWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeForQWACBlock {
	c := &CertQualificationAtTimeForQWACBlock{
		CertQualificationAtTimeBlock: NewCertQualificationAtTimeBlockAtIssuanceTime(i18nProvider, validationTime, signingCertificate, acceptableServices),
	}
	c.InitCertQualificationAtTimeBlock(c)
	return c
}

// NewCertQualificationAtTimeForQWACBlock is the constructor to instantiate
// the validation at the validation time. Port of
// CertQualificationAtTimeForQWACBlock(I18nProvider, ValidationTime, Date, CertificateWrapper, List).
func NewCertQualificationAtTimeForQWACBlock(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime, date *time.Time,
	signingCertificate *diagnostic.CertificateWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeForQWACBlock {
	c := &CertQualificationAtTimeForQWACBlock{
		CertQualificationAtTimeBlock: NewCertQualificationAtTimeBlock(i18nProvider, validationTime, date, signingCertificate, acceptableServices),
	}
	c.InitCertQualificationAtTimeBlock(c)
	return c
}

// ExecuteQSCDCheck defines whether a QSCD check should be processed for
// certificate qualification determination. Port of the overridden
// executeQSCDCheck(): not required for a QWAC.
func (c *CertQualificationAtTimeForQWACBlock) ExecuteQSCDCheck() bool {
	return false
}
