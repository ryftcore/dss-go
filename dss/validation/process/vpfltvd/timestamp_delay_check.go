// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/TimestampDelayCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"math"
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampDelayCheck checks if the claimed signing time + timestamp's delay
// is after the best-signature-time.
type TimestampDelayCheck[T any] struct {
	*process.ChainItemBase[T]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper

	// bestSignatureTime is the best signature time.
	bestSignatureTime *time.Time

	// durationRule is the timestamp's delay constraint.
	durationRule policy.DurationRule
}

// NewTimestampDelayCheck is the default constructor. Port of
// TimestampDelayCheck(I18nProvider, T, SignatureWrapper, Date, DurationRule).
func NewTimestampDelayCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	signature *diagnostic.SignatureWrapper, bestSignatureTime *time.Time,
	durationRule policy.DurationRule) *TimestampDelayCheck[T] {
	c := &TimestampDelayCheck[T]{
		ChainItemBase:     process.NewChainItemBase(i18nProvider, result, durationRule),
		signature:         signature,
		bestSignatureTime: bestSignatureTime,
		durationRule:      durationRule,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TimestampDelayCheck[T]) Process() bool {
	signingTime := c.signature.ClaimedSigningTime()
	if signingTime == nil || c.bestSignatureTime == nil {
		return false
	}
	delayMilliseconds := c.durationRule.Duration()
	var limit time.Time
	if delayMilliseconds == math.MaxInt64 {
		limit = time.UnixMilli(math.MaxInt64)
	} else {
		limit = signingTime.Add(time.Duration(delayMilliseconds) * time.Millisecond)
	}
	return limit.After(*c.bestSignatureTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampDelayCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ISTPTDABST
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampDelayCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ISTPTDABST_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampDelayCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampDelayCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
