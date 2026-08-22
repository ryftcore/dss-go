// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWACCertificatePolicyCheck.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QWACCertificatePolicyCheck verifies whether the certificate has been
// issued under the appropriate QWAC certificate policy.
type QWACCertificatePolicyCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// certificate is the certificate to be validated.
	certificate *diagnostic.CertificateWrapper

	// qwacProfile is the QWAC validation profile.
	qwacProfile enumerations.QWACProfile
}

// NewQWACCertificatePolicyCheck is the default constructor. Port of
// QWACCertificatePolicyCheck(Provider, XmlValidationQWACProcess, CertificateWrapper, QWACProfile, LevelRule).
func NewQWACCertificatePolicyCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificate *diagnostic.CertificateWrapper, qwacProfile enumerations.QWACProfile,
	constraint policy.LevelRule) *QWACCertificatePolicyCheck {
	c := &QWACCertificatePolicyCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		qwacProfile:   qwacProfile,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QWACCertificatePolicyCheck) Process() bool {
	switch c.qwacProfile {
	case enumerations.QWACProfileQWAC1:
		return c.certificatePolicyMatch(enumerations.CertificatePolicyQCPWeb, enumerations.CertificatePolicyQNCPWeb)
	case enumerations.QWACProfileQWAC2:
		return c.certificatePolicyMatch(enumerations.CertificatePolicyQNCPWebGen)
	default:
		panic(fmt.Sprintf("The QWAC profile '%s' is not supported!", c.qwacProfile))
	}
}

// certificatePolicyMatch ports the private certificatePolicyMatch(CertificatePolicy...).
func (c *QWACCertificatePolicyCheck) certificatePolicyMatch(policies ...enumerations.CertificatePolicy) bool {
	certificatePoliciesOids := c.certificate.CertificatePoliciesOids()
	if !utils.IsCollectionNotEmpty(certificatePoliciesOids) {
		return false
	}
	for _, k := range certificatePoliciesOids {
		for _, m := range policies {
			if m.OID() == k {
				return true
			}
		}
	}
	return false
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *QWACCertificatePolicyCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	tag, err := process.GetQWACValidationMessageTag(c.qwacProfile)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTagQWACCertPolicy, tag)
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *QWACCertificatePolicyCheck) BuildErrorMessage() *jaxb.XmlMessage {
	tag, err := process.GetQWACValidationMessageTag(c.qwacProfile)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTagQWACCertPolicyANS, tag)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QWACCertificatePolicyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QWACCertificatePolicyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
