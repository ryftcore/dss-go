// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/RevocationDataAcceptableCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationDataAcceptableCheck verifies the result of a basic revocation
// validation process.
type RevocationDataAcceptableCheck[T any] struct {
	*process.ChainItemBase[T]

	// revocationId is the Id of a revocation data to be checked.
	revocationId string

	// xmlConclusion is the revocation basic validation result.
	xmlConclusion *jaxb.XmlConclusion
}

// NewRevocationDataAcceptableCheck is the default constructor. Port of
// RevocationDataAcceptableCheck(I18nProvider, T, String, XmlConclusion, LevelRule).
func NewRevocationDataAcceptableCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	revocationId string, xmlConclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *RevocationDataAcceptableCheck[T] {
	c := &RevocationDataAcceptableCheck[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, revocationId),
		revocationId:  revocationId,
		xmlConclusion: xmlConclusion,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *RevocationDataAcceptableCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeRevBBB
}

// Process performs the check. Port of process().
func (c *RevocationDataAcceptableCheck[T]) Process() bool {
	return process.IsAllowedBasicRevocationDataValidation(c.xmlConclusion)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationDataAcceptableCheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagTokenID, c.revocationId)
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationDataAcceptableCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlConclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationDataAcceptableCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlConclusion.SubIndication == nil {
		return ""
	}
	return c.xmlConclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDataAcceptableCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTRORPIIC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationDataAcceptableCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTRORPIICANS
}
