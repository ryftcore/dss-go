// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/RevocationFreshnessCheckerResultCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationFreshnessCheckerResultCheck checks if the revocation freshness
// checker's result is valid.
type RevocationFreshnessCheckerResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// rfcResult is the RFC result.
	rfcResult *jaxb.XmlRFC
}

// NewRevocationFreshnessCheckerResultCheck is the default constructor. Port
// of RevocationFreshnessCheckerResultCheck(Provider, T, XmlRFC, LevelRule).
func NewRevocationFreshnessCheckerResultCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	rfcResult *jaxb.XmlRFC, constraint policy.LevelRule) *RevocationFreshnessCheckerResultCheck[T] {
	var chainItemBase *process.ChainItemBase[T]
	if rfcResult.Id != nil {
		chainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *rfcResult.Id)
	} else {
		chainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &RevocationFreshnessCheckerResultCheck[T]{
		ChainItemBase: chainItemBase,
		rfcResult:     rfcResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *RevocationFreshnessCheckerResultCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeRFC
}

// Process performs the check. Port of process().
func (c *RevocationFreshnessCheckerResultCheck[T]) Process() bool {
	return c.IsValid(&c.rfcResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationFreshnessCheckerResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRFC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationFreshnessCheckerResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRFCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationFreshnessCheckerResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.rfcResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationFreshnessCheckerResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.rfcResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.rfcResult.Conclusion.SubIndication.SubIndication()
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of getPreviousErrors().
func (c *RevocationFreshnessCheckerResultCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.rfcResult.Conclusion.Errors
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationFreshnessCheckerResultCheck[T]) BuildAdditionalInfo() *string {
	if c.rfcResult.Id != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagTokenID, *c.rfcResult.Id)
		return &message
	}
	return nil
}
