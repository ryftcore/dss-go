// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/AtLeastOneReferenceDataObjectFoundCheck.java (DSS 6.5.RC1).
package cv

import (
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AtLeastOneReferenceDataObjectFoundCheck checks if at least one covered data
// object has been found.
type AtLeastOneReferenceDataObjectFoundCheck[T any] struct {
	*process.ChainItemBase[T]

	// digestMatchers is the collection of DigestMatchers.
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher
}

// NewAtLeastOneReferenceDataObjectFoundCheck is the default constructor. Port of
// AtLeastOneReferenceDataObjectFoundCheck(I18nProvider, T, List, LevelRule).
func NewAtLeastOneReferenceDataObjectFoundCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher,
	constraint policy.LevelRule) *AtLeastOneReferenceDataObjectFoundCheck[T] {
	c := &AtLeastOneReferenceDataObjectFoundCheck[T]{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatchers: digestMatchers,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AtLeastOneReferenceDataObjectFoundCheck[T]) Process() bool {
	for _, d := range c.digestMatchers {
		if d.DataFound {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AtLeastOneReferenceDataObjectFoundCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_IODOF
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *AtLeastOneReferenceDataObjectFoundCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_IODOF_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AtLeastOneReferenceDataObjectFoundCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *AtLeastOneReferenceDataObjectFoundCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignedDataNotFound
}
