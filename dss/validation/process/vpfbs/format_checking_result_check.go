// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/FormatCheckingResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// FormatCheckingResultCheck verifies if the format checking process as per
// clause 5.2.2 succeeded.
type FormatCheckingResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlFC is the Format Checking process result.
	xmlFC *jaxb.XmlFC
}

// NewFormatCheckingResultCheck is the default constructor. Port of
// FormatCheckingResultCheck(I18nProvider, T, XmlFC, TokenProxy, LevelRule).
func NewFormatCheckingResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlFC *jaxb.XmlFC, token diagnostic.TokenProxy, constraint policy.LevelRule) *FormatCheckingResultCheck[T] {
	c := &FormatCheckingResultCheck[T]{
		// Format Checking building block suffix ("-FC"), a per-class private
		// constant in Java; inlined here since Go package-level constants share
		// one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-FC"),
		xmlFC:         xmlFC,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *FormatCheckingResultCheck[T]) Process() bool {
	return c.xmlFC != nil && c.IsValid(&c.xmlFC.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *FormatCheckingResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *FormatCheckingResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *FormatCheckingResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IFCRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *FormatCheckingResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IFCRC_ANS
}
