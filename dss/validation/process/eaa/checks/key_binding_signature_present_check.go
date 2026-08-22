// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/KeyBindingSignaturePresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// KeyBindingSignaturePresentCheck verifies whether the EAA contains an
// attached key binding signature.
type KeyBindingSignaturePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlFC]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewKeyBindingSignaturePresentCheck is the default constructor.
func NewKeyBindingSignaturePresentCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlFC],
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
	return i18n.MessageTagEAAKBSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *KeyBindingSignaturePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAKBSPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *KeyBindingSignaturePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *KeyBindingSignaturePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
