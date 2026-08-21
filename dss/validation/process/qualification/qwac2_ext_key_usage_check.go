// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWAC2ExtKeyUsageCheck.java (DSS 6.5.RC1).
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

// QWAC2ExtKeyUsageCheck verifies conformance of the "extended-key-usage"
// certificate extension for the 2-QWAC profile.
type QWAC2ExtKeyUsageCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// certificate is the certificate to be validated.
	certificate *diagnostic.CertificateWrapper
}

// NewQWAC2ExtKeyUsageCheck is the default constructor. Port of
// QWAC2ExtKeyUsageCheck(I18nProvider, XmlValidationQWACProcess, CertificateWrapper, LevelRule).
func NewQWAC2ExtKeyUsageCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *QWAC2ExtKeyUsageCheck {
	c := &QWAC2ExtKeyUsageCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QWAC2ExtKeyUsageCheck) Process() bool {
	extendedKeyUsages := c.certificate.ExtendedKeyUsages()
	return utils.CollectionSize(extendedKeyUsages) == 1 &&
		enumerations.ExtendedKeyUsage_TSL_BINDING.OID() == extendedKeyUsages[0].Value
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QWAC2ExtKeyUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC2_EXT_KEY_USAGE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QWAC2ExtKeyUsageCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC2_EXT_KEY_USAGE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QWAC2ExtKeyUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QWAC2ExtKeyUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
