// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAANotOnHoldCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAANotOnHoldCheck checks whether the corresponding status declares that the
// EAA is not on hold.
type EAANotOnHoldCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationWrapper
}

// NewEAANotOnHoldCheck is the default constructor.
func NewEAANotOnHoldCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationWrapper, constraint policy.LevelRule) *EAANotOnHoldCheck {
	c := &EAANotOnHoldCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAANotOnHoldCheck) Process() bool {
	return c.eaaStatusToken == nil || enumerations.EAAStatusSuspended != c.eaaStatusToken.Status()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAANotOnHoldCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_NOT_ON_HOLD
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAANotOnHoldCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_NOT_ON_HOLD_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAANotOnHoldCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAANotOnHoldCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTryLater
}
