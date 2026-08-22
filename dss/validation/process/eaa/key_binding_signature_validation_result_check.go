// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/KeyBindingSignatureValidationResultCheck.java (DSS 6.5.RC1).
//
// KeyBindingSignatureValidationResultCheck wraps
// qualification.SignatureValidationResultCheck for the key-binding
// signature's own validation result.
package eaa

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/qualification"
)

// KeyBindingSignatureValidationResultCheck performs verification of the
// validation result of the key binding signature.
type KeyBindingSignatureValidationResultCheck struct {
	*qualification.SignatureValidationResultCheck[*jaxb.XmlValidationProcessEAA]
}

// NewKeyBindingSignatureValidationResultCheck is the default constructor.
// Port of
// KeyBindingSignatureValidationResultCheck(I18nProvider, XmlValidationProcessEAA, XmlConclusion, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' MessageTag/ErrorMessageTag rather than
// SignatureValidationResultCheck's.
func NewKeyBindingSignatureValidationResultCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessEAA], basicValidationConclusion *jaxb.XmlConclusion,
	constraint policy.LevelRule) *KeyBindingSignatureValidationResultCheck {
	c := &KeyBindingSignatureValidationResultCheck{
		SignatureValidationResultCheck: qualification.NewSignatureValidationResultCheck(i18nProvider, result, basicValidationConclusion, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden
// getMessageTag().
func (c *KeyBindingSignatureValidationResultCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAKBRC
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *KeyBindingSignatureValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAKBRCANS
}
