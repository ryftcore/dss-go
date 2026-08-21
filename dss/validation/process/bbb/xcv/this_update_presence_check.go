// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/ThisUpdatePresenceCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ThisUpdatePresenceCheck verifies whether the 'thisUpdate' field is defined
// within the revocation information.
type ThisUpdatePresenceCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewThisUpdatePresenceCheck is the default constructor. Port of
// ThisUpdatePresenceCheck(I18nProvider, XmlRAC, RevocationWrapper, LevelRule).
func NewThisUpdatePresenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *ThisUpdatePresenceCheck {
	c := &ThisUpdatePresenceCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ThisUpdatePresenceCheck) Process() bool {
	return c.revocationData.ThisUpdate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ThisUpdatePresenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ThisUpdatePresenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ThisUpdatePresenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ThisUpdatePresenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
