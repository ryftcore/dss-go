// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/SignatureAcceptanceValidationResultCheck.java (DSS 6.5.RC1).
//
// The class is designed for extension - SignatureAcceptanceValidationNoCryptoResultCheck
// overrides process() - so the self-call is routed through
// SignatureAcceptanceValidationResultCheckOverrides, the way Chain routes its
// own through ChainOverrides. A subclass registers itself once, with
// InitChainItem, following the embed-and-re-register technique used
// throughout this port (e.g. bbb/xcv's trustAnchorCheckSubXCVResult).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureAcceptanceValidationResultCheck verifies if the format Signature
// Acceptance Validation process as per clause 5.2.8 succeeded.
type SignatureAcceptanceValidationResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// XmlSAV is the Signature Acceptance Validation result. Exported because
	// Java declares the field protected (SignatureAcceptanceValidationNoCryptoResultCheck
	// reads it).
	XmlSAV *jaxb.XmlSAV
}

// NewSignatureAcceptanceValidationResultCheck is the default constructor.
// Port of
// SignatureAcceptanceValidationResultCheck(I18nProvider, T, XmlSAV, TokenProxy, LevelRule).
func NewSignatureAcceptanceValidationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlSAV *jaxb.XmlSAV, token diagnostic.TokenProxy, constraint policy.LevelRule) *SignatureAcceptanceValidationResultCheck[T] {
	c := &SignatureAcceptanceValidationResultCheck[T]{
		// Signature Acceptance Validation building block suffix ("-SAV"), a
		// per-class private constant in Java; inlined here since Go
		// package-level constants share one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-SAV"),
		XmlSAV:        xmlSAV,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignatureAcceptanceValidationResultCheck[T]) Process() bool {
	return c.XmlSAV != nil && c.IsValid(&c.XmlSAV.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignatureAcceptanceValidationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.XmlSAV.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SignatureAcceptanceValidationResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.XmlSAV.Conclusion.SubIndication == nil {
		return ""
	}
	return c.XmlSAV.Conclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureAcceptanceValidationResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ISAVRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignatureAcceptanceValidationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ISAVRC_ANS
}
