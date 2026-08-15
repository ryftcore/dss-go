// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/StructuralValidationCheck.java (DSS 6.5.RC1).
package sav

import (
	"strings"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// StructuralValidationCheck checks if the structural validation of the signature succeeds.
type StructuralValidationCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewStructuralValidationCheck is the default constructor. Port of
// StructuralValidationCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewStructuralValidationCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *StructuralValidationCheck {
	c := &StructuralValidationCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *StructuralValidationCheck) Process() bool {
	return c.signature.IsStructuralValidationValid()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *StructuralValidationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISSV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *StructuralValidationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISSV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *StructuralValidationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *StructuralValidationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *StructuralValidationCheck) BuildAdditionalInfo() *string {
	errorMessages := c.signature.StructuralValidationMessages()
	if utils.IsCollectionNotEmpty(errorMessages) {
		// Java passes errorMessages.toString() (java.util.List#toString: "[a,
		// b, c]"); replicated explicitly here since messageFormatArgString
		// would otherwise render a []string with Go's "[a b c]" format.
		message := c.I18nProvider.GetMessage(i18n.MessageTag_STRUCTURAL_VALIDATION_FAILURE,
			"["+strings.Join(errorMessages, ", ")+"]")
		return &message
	}
	return nil
}
