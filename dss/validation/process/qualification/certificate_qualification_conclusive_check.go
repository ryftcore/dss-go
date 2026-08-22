// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/CertificateQualificationConclusiveCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateQualificationConclusiveCheck checks whether the certificate
// qualification process has been performed successfully.
type CertificateQualificationConclusiveCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// certificateQualification is the qualification validation process of
	// the certificate.
	certificateQualification *jaxb.XmlCertificateQualificationProcess
}

// NewCertificateQualificationConclusiveCheck is the default constructor.
// Port of
// CertificateQualificationConclusiveCheck(Provider, XmlValidationQWACProcess, XmlCertificateQualificationProcess, LevelRule).
func NewCertificateQualificationConclusiveCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificateQualification *jaxb.XmlCertificateQualificationProcess, constraint policy.LevelRule) *CertificateQualificationConclusiveCheck {
	c := &CertificateQualificationConclusiveCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		certificateQualification: certificateQualification,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQualificationConclusiveCheck) Process() bool {
	return c.IsValidConclusion(c.certificateQualification.Conclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQualificationConclusiveCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQWACCertQualConclusive
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQualificationConclusiveCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQWACCertQualConclusiveANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQualificationConclusiveCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CertificateQualificationConclusiveCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
