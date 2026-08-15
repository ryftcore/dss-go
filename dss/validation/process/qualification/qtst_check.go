// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/timestamp/checks/QTSTCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// QTSTCheck checks whether the certificate used to issue a timestamp is
// QTST.
type QTSTCheck[T any] struct {
	*process.ChainItemBase[T]

	// trustServicesAtTime is the list of TrustServices declaring QTST status
	// for the certificate.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewQTSTCheck is the default constructor. Port of
// QTSTCheck(I18nProvider, T, List, LevelRule).
func NewQTSTCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *QTSTCheck[T] {
	c := &QTSTCheck[T]{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QTSTCheck[T]) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QTSTCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_QTST
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QTSTCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_QTST_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QTSTCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QTSTCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
