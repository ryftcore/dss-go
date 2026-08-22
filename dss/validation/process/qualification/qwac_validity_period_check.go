// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWACValidityPeriodCheck.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
)

// QWACValidityPeriodCheck verifies the validity period of the QWAC
// certificate against the current date and time.
type QWACValidityPeriodCheck struct {
	*xcv.CertificateValidityRangeCheck[*jaxb.XmlValidationQWACProcess]
}

// NewQWACValidityPeriodCheck is the default constructor. Port of
// QWACValidityPeriodCheck(I18nProvider, XmlValidationQWACProcess, CertificateWrapper, Date, LevelRule).
func NewQWACValidityPeriodCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, constraint policy.LevelRule) *QWACValidityPeriodCheck {
	c := &QWACValidityPeriodCheck{
		CertificateValidityRangeCheck: xcv.NewCertificateValidityRangeCheckMinimal(i18nProvider, result, certificate, currentTime, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden getMessageTag().
func (c *QWACValidityPeriodCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC_VAL_PERIOD
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *QWACValidityPeriodCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC_VAL_PERIOD_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the overridden getFailedIndicationForConclusion().
func (c *QWACValidityPeriodCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion(), whose
// default is null.
func (c *QWACValidityPeriodCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
