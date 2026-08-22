// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWAC2ValidationResultCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QWAC2ValidationResultCheck verifies whether the QWAC validation has
// succeeded for the 2-QWAC profile.
type QWAC2ValidationResultCheck struct {
	*QWACValidationResultCheck
}

// NewQWAC2ValidationResultCheck is the default constructor. Port of
// QWAC2ValidationResultCheck(Provider, XmlQWACProcess, XmlValidationQWACProcess[], LevelRule).
func NewQWAC2ValidationResultCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlQWACProcess],
	qwacValidationProcesses []*jaxb.XmlValidationQWACProcess, constraint policy.LevelRule) *QWAC2ValidationResultCheck {
	c := &QWAC2ValidationResultCheck{
		QWACValidationResultCheck: NewQWACValidationResultCheck(i18nProvider, result, qwacValidationProcesses, constraint),
	}
	c.InitChainItem(c)
	return c
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *QWAC2ValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQWACValidANS2
}
