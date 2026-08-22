// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/ContentTimestampsCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ContentTimestampsCheck checks if a collection of content timestamps is not
// empty.
type ContentTimestampsCheck[T any] struct {
	*process.ChainItemBase[T]

	// contentTimestamps is the content timestamps collection.
	contentTimestamps []*diagnostic.TimestampWrapper
}

// NewContentTimestampsCheck is the default constructor. Port of
// ContentTimestampsCheck(I18nProvider, T, List, LevelRule).
func NewContentTimestampsCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	contentTimestamps []*diagnostic.TimestampWrapper, constraint policy.LevelRule) *ContentTimestampsCheck[T] {
	c := &ContentTimestampsCheck[T]{
		ChainItemBase:     process.NewChainItemBase(i18nProvider, result, constraint),
		contentTimestamps: contentTimestamps,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ContentTimestampsCheck[T]) Process() bool {
	return utils.IsCollectionNotEmpty(c.contentTimestamps)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion(), whose default is null.
func (c *ContentTimestampsCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return ""
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ContentTimestampsCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ContentTimestampsCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVISCCTC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ContentTimestampsCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVISCCTCANS
}
