// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAARevocationStatusKnownCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationStatusKnownCheck verifies whether the obtained EAA revocation
// is known.
type EAARevocationStatusKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatus is the EAA revocation token to check.
	eaaStatus *diagnostic.EAARevocationWrapper
}

// NewEAARevocationStatusKnownCheck is the default constructor.
func NewEAARevocationStatusKnownCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatus *diagnostic.EAARevocationWrapper, constraint policy.LevelRule) *EAARevocationStatusKnownCheck {
	c := &EAARevocationStatusKnownCheck{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, eaaStatus.Id()),
		eaaStatus:     eaaStatus,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationStatusKnownCheck) Process() bool {
	return enumerations.EAAStatusValid == c.eaaStatus.Status() ||
		enumerations.EAAStatusInvalid == c.eaaStatus.Status() ||
		enumerations.EAAStatusSuspended == c.eaaStatus.Status()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationStatusKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_KNOWN
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationStatusKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_KNOWN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationStatusKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationStatusKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
