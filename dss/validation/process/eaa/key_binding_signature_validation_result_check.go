// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/KeyBindingSignatureValidationResultCheck.java (DSS 6.5.RC1).
//
// NECESSARY ADDITION (LTVB, phase 8e): not itself a line in this porter's
// manifest (only EAAValidationBlock.java and EAAValidationProcess.java are),
// but doc.go already flagged it as the one eaa.checks class the phase 8d
// integration pass deliberately left unported, "wires
// eu.europa.esig.dss.validation.process.qualification.signature.checks.SignatureValidationResultCheck
// and eu.europa.esig.dss.detailedreport.jaxb.XmlValidationProcessEAA... left
// unported, consistent with 'port minimally from upstream'" - i.e. left for
// whichever future pass ports EAAValidationProcess, which is this one:
// EAAValidationProcess.keyBindingSignatureValidationConclusive() constructs it
// directly. Porting it here keeps the eaa package buildable.
//
// FORWARD DEPENDENCY: qualification.SignatureValidationResultCheck is ported
// by the shared QCERT/QTRUST/QSIG porters into package
// github.com/ryftcore/dss-go/dss/validation/process/qualification, assumed to have
// the shape
//
//	type SignatureValidationResultCheck[T any] struct { *process.ChainItemBase[T]; ... }
//	func NewSignatureValidationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
//	    basicValidationConclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *SignatureValidationResultCheck[T]
//
// mirroring Java's
// SignatureValidationResultCheck(I18nProvider, T, XmlConclusion, LevelRule).
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
