// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/AcceptableBuildingBlockConclusionCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableBuildingBlockConclusionCheck verifies whether the
// BasicBuildingBlock's validation succeeded.
type AcceptableBuildingBlockConclusionCheck[T any] struct {
	*process.ChainItemBase[T]

	// buildingBlockConclusion is the BasicBuildingBlock's validation
	// conclusion.
	buildingBlockConclusion *jaxb.XmlConclusion
}

// NewAcceptableBuildingBlockConclusionCheck is the default constructor. Port
// of AcceptableBuildingBlockConclusionCheck(I18nProvider, T, XmlConclusion, LevelRule).
func NewAcceptableBuildingBlockConclusionCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	buildingBlockConclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *AcceptableBuildingBlockConclusionCheck[T] {
	c := &AcceptableBuildingBlockConclusionCheck[T]{
		ChainItemBase:           process.NewChainItemBase(i18nProvider, result, constraint),
		buildingBlockConclusion: buildingBlockConclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableBuildingBlockConclusionCheck[T]) Process() bool {
	return c.IsValidConclusion(c.buildingBlockConclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableBuildingBlockConclusionCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ACCEPT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AcceptableBuildingBlockConclusionCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ACCEPT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableBuildingBlockConclusionCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.buildingBlockConclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AcceptableBuildingBlockConclusionCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.buildingBlockConclusion.SubIndication == nil {
		return ""
	}
	return c.buildingBlockConclusion.SubIndication.SubIndication()
}
