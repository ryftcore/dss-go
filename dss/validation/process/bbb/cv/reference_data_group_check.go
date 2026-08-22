// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ReferenceDataGroupCheck.java (DSS 6.5.RC1).
package cv

import (
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ReferenceDataGroupCheck checks if only hashes of only provided archive data
// objects are present at the first level of the reduced hash tree.
type ReferenceDataGroupCheck[T any] struct {
	*process.ChainItemBase[T]

	// digestMatchers is the collection of DigestMatchers.
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher
}

// NewReferenceDataGroupCheck is the default constructor. Port of
// ReferenceDataGroupCheck(I18nProvider, T, List, LevelRule).
func NewReferenceDataGroupCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ReferenceDataGroupCheck[T] {
	c := &ReferenceDataGroupCheck[T]{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatchers: digestMatchers,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ReferenceDataGroupCheck[T]) Process() bool {
	for _, d := range c.digestMatchers {
		if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference == digestMatcherType(d) {
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ReferenceDataGroupCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_DFHVLCDOG
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataGroupCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_DFHVLCDOG_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ReferenceDataGroupCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ReferenceDataGroupCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignedDataNotFound
}
