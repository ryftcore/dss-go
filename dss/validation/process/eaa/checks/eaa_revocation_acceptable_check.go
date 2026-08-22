// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAARevocationAcceptableCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationAcceptableCheck checks whether the EAA revocation token is
// acceptable.
type EAARevocationAcceptableCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationWrapper

	// conclusion is the BBB validation result of the token.
	conclusion *jaxb.XmlConclusion
}

// NewEAARevocationAcceptableCheck is the default constructor.
func NewEAARevocationAcceptableCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationWrapper, conclusion *jaxb.XmlConclusion,
	constraint policy.LevelRule) *EAARevocationAcceptableCheck {
	c := &EAARevocationAcceptableCheck{
		ChainItemBase:  process.NewChainItemBaseWithId(i18nProvider, result, constraint, eaaStatusToken.Id()),
		eaaStatusToken: eaaStatusToken,
		conclusion:     conclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationAcceptableCheck) Process() bool {
	return c.IsValidConclusion(c.conclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationAcceptableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevACC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationAcceptableCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevACCANS
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAARevocationAcceptableCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagTokenID, c.eaaStatusToken.Id())
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationAcceptableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationAcceptableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
