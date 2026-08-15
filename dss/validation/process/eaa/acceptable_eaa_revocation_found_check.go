// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/AcceptableEAARevocationFoundCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AcceptableEAARevocationFoundCheck checks whether an acceptable EAA
// revocation was found.
type AcceptableEAARevocationFoundCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationWrapper
}

// NewAcceptableEAARevocationFoundCheck is the default constructor.
func NewAcceptableEAARevocationFoundCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationWrapper, constraint policy.LevelRule) *AcceptableEAARevocationFoundCheck {
	c := &AcceptableEAARevocationFoundCheck{
		eaaStatusToken: eaaStatusToken,
	}
	if eaaStatusToken != nil {
		c.ChainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, eaaStatusToken.Id())
	} else {
		c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableEAARevocationFoundCheck) Process() bool {
	return c.eaaStatusToken != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableEAARevocationFoundCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_ACC_FND
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *AcceptableEAARevocationFoundCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_ACC_FND_ANS
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *AcceptableEAARevocationFoundCheck) BuildAdditionalInfo() *string {
	if c.eaaStatusToken != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_TOKEN_ID, c.eaaStatusToken.Id())
		return &message
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableEAARevocationFoundCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AcceptableEAARevocationFoundCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
