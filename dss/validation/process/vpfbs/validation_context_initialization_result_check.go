// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/ValidationContextInitializationResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationContextInitializationResultCheck verifies if the Validation
// Context Initialization as per clause 5.2.4 succeeded.
type ValidationContextInitializationResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlVCI is the Validation Context Initialization result.
	xmlVCI *jaxb.XmlVCI
}

// NewValidationContextInitializationResultCheck is the default constructor.
// Port of
// ValidationContextInitializationResultCheck(I18nProvider, T, XmlVCI, TokenProxy, LevelRule).
func NewValidationContextInitializationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlVCI *jaxb.XmlVCI, token diagnostic.TokenProxy, constraint policy.LevelRule) *ValidationContextInitializationResultCheck[T] {
	c := &ValidationContextInitializationResultCheck[T]{
		// Validation Context Initialization building block suffix ("-VCI"), a
		// per-class private constant in Java; inlined here since Go
		// package-level constants share one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-VCI"),
		xmlVCI:        xmlVCI,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ValidationContextInitializationResultCheck[T]) Process() bool {
	return c.xmlVCI != nil && c.IsValid(&c.xmlVCI.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ValidationContextInitializationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlVCI.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ValidationContextInitializationResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlVCI.Conclusion.SubIndication == nil {
		return ""
	}
	return c.xmlVCI.Conclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidationContextInitializationResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIVCIRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ValidationContextInitializationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIVCIRCANS
}
