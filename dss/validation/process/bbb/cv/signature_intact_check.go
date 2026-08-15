// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/SignatureIntactCheck.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.bbb.cv and its .checks subpackage
// flatten into this single Go package cv (no name collisions), so the checks are
// referenced unqualified.
//
// Every check constructor takes *process.Result[T] where Java takes the result
// object itself: see chain.go for why the generated result object has to be
// bound to the JAXB base structs it embeds.
package cv

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureIntactCheck checks if the signature value is intact.
type SignatureIntactCheck[T any] struct {
	*process.ChainItemBase[T]

	// Token is the token to check. Exported because Java declares the field
	// protected and SignatureIntactWithIdCheck reads it.
	Token diagnostic.TokenProxy

	// context is the validation context.
	context enumerations.Context
}

// NewSignatureIntactCheck is the default constructor. Port of
// SignatureIntactCheck(I18nProvider, T, TokenProxy, Context, LevelRule).
func NewSignatureIntactCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	token diagnostic.TokenProxy, context enumerations.Context, constraint policy.LevelRule) *SignatureIntactCheck[T] {
	c := &SignatureIntactCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		Token:         token,
		context:       context,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignatureIntactCheck[T]) Process() bool {
	return c.Token.IsSignatureIntact()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureIntactCheck[T]) MessageTag() i18n.MessageTag {
	switch c.context {
	case enumerations.Context_CERTIFICATE:
		return i18n.MessageTag_BBB_CV_ISIC
	case enumerations.Context_REVOCATION:
		return i18n.MessageTag_BBB_CV_ISIR
	case enumerations.Context_TIMESTAMP:
		return i18n.MessageTag_BBB_CV_ISIT
	default:
		return i18n.MessageTag_BBB_CV_ISI
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignatureIntactCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ISI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignatureIntactCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignatureIntactCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CRYPTO_FAILURE
}
