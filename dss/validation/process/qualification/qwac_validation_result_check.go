// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWACValidationResultCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QWACValidationResultCheck verifies whether the QWAC validation has
// succeeded at least for one QWAC profile.
type QWACValidationResultCheck struct {
	*process.ChainItemBase[*jaxb.XmlQWACProcess]

	// qwacValidationProcesses is the array of XmlValidationQWACProcesses.
	qwacValidationProcesses []*jaxb.XmlValidationQWACProcess
}

// NewQWACValidationResultCheck is the default constructor. Port of
// QWACValidationResultCheck(Provider, XmlQWACProcess, XmlValidationQWACProcess[], LevelRule).
func NewQWACValidationResultCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlQWACProcess],
	qwacValidationProcesses []*jaxb.XmlValidationQWACProcess, constraint policy.LevelRule) *QWACValidationResultCheck {
	c := &QWACValidationResultCheck{
		ChainItemBase:           process.NewChainItemBase(i18nProvider, result, constraint),
		qwacValidationProcesses: qwacValidationProcesses,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QWACValidationResultCheck) Process() bool {
	if utils.IsCollectionNotEmpty(c.qwacValidationProcesses) {
		for _, qwacProcess := range c.qwacValidationProcesses {
			if c.IsValid(&qwacProcess.XmlConstraintsConclusionContent) {
				return true
			}
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QWACValidationResultCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQWACValid
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QWACValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQWACValidANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QWACValidationResultCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QWACValidationResultCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
