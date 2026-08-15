// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/KeyBindingSignaturePresentCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// KeyBindingSignaturePresentCheck verifies whether the EAA contains an
// attached key binding signature.
type KeyBindingSignaturePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlFC]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewKeyBindingSignaturePresentCheck is the default constructor.
func NewKeyBindingSignaturePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlFC],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *KeyBindingSignaturePresentCheck {
	c := &KeyBindingSignaturePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *KeyBindingSignaturePresentCheck) Process() bool {
	return c.eaa.KeyBindingSignature() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *KeyBindingSignaturePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_KBSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *KeyBindingSignaturePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_KBSP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *KeyBindingSignaturePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *KeyBindingSignaturePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
