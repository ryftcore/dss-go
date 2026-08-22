// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/RelatedToMraEnactedTrustServiceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RelatedToMraEnactedTrustServiceCheck verifies whether MRA enacted trusted
// services are present.
type RelatedToMraEnactedTrustServiceCheck[T any] struct {
	*process.ChainItemBase[T]

	// trustServicesAtTime is the list of TrustServiceWrappers at control time.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewRelatedToMraEnactedTrustServiceCheck is the default constructor. Port of
// RelatedToMraEnactedTrustServiceCheck(I18nProvider, T, List, LevelRule).
func NewRelatedToMraEnactedTrustServiceCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *RelatedToMraEnactedTrustServiceCheck[T] {
	c := &RelatedToMraEnactedTrustServiceCheck[T]{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RelatedToMraEnactedTrustServiceCheck[T]) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RelatedToMraEnactedTrustServiceCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_METS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RelatedToMraEnactedTrustServiceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_METS_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RelatedToMraEnactedTrustServiceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *RelatedToMraEnactedTrustServiceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
