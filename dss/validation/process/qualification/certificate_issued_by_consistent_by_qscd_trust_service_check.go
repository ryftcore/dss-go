// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/CertificateIssuedByConsistentByQSCDTrustServiceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateIssuedByConsistentByQSCDTrustServiceCheck checks if there are
// consistent by QSCD TrustServices issued the certificate in question at
// control time.
type CertificateIssuedByConsistentByQSCDTrustServiceCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServicesAtTime is the list of consistent Trusted Services issued
	// the certificate at control time.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewCertificateIssuedByConsistentByQSCDTrustServiceCheck is the default
// constructor. Port of
// CertificateIssuedByConsistentByQSCDTrustServiceCheck(I18nProvider, XmlValidationCertificateQualification, List, LevelRule).
func NewCertificateIssuedByConsistentByQSCDTrustServiceCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification], trustServicesAtTime []*diagnostic.TrustServiceWrapper,
	constraint policy.LevelRule) *CertificateIssuedByConsistentByQSCDTrustServiceCheck {
	c := &CertificateIssuedByConsistentByQSCDTrustServiceCheck{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateIssuedByConsistentByQSCDTrustServiceCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateIssuedByConsistentByQSCDTrustServiceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasConsistentByQSCD
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateIssuedByConsistentByQSCDTrustServiceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasConsistentByQSCDANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateIssuedByConsistentByQSCDTrustServiceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CertificateIssuedByConsistentByQSCDTrustServiceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
