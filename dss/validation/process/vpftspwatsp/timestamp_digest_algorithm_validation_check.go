// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftspwatsp/checks/TimestampDigestAlgorithmValidationCheck.java (DSS 6.5.RC1).
package vpftspwatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampDigestAlgorithmValidationCheck verifies whether the result of
// MessageImprintDigestAlgorithmValidation is valid.
type TimestampDigestAlgorithmValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// timestamp is the timestamp to check.
	timestamp *diagnostic.TimestampWrapper

	// cvResult is the message-imprint Digest Algorithm validation result.
	cvResult *jaxb.XmlCryptographicValidation

	// currentTime defines the validation time.
	currentTime time.Time
}

// NewTimestampDigestAlgorithmValidationCheck is the default constructor. Port
// of TimestampDigestAlgorithmValidationCheck(I18nProvider, T, TimestampWrapper, XmlCryptographicValidation, Date, LevelRule).
func NewTimestampDigestAlgorithmValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, cvResult *jaxb.XmlCryptographicValidation, currentTime time.Time,
	constraint policy.LevelRule) *TimestampDigestAlgorithmValidationCheck[T] {
	c := &TimestampDigestAlgorithmValidationCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		timestamp:     timestamp,
		cvResult:      cvResult,
		currentTime:   currentTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TimestampDigestAlgorithmValidationCheck[T]) Process() bool {
	return c.cvResult != nil && c.IsValidConclusion(c.cvResult.Conclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampDigestAlgorithmValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_ICHFCRLPOET
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampDigestAlgorithmValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_ICHFCRLPOET_ANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(): Java dereferences getMessageImprint() unguarded, so a
// time-stamp without one raises a NullPointerException; the nil pointer panics
// here in its place. The generated DigestMethod member is a
// *DigestAlgorithmValue, whose nil java.text.MessageFormat renders as "null".
func (c *TimestampDigestAlgorithmValidationCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.currentTime)
	digestMethod := "null"
	if c.timestamp.MessageImprint().DigestMethod != nil {
		digestMethod = string(c.timestamp.MessageImprint().DigestMethod.DigestAlgorithm())
	}
	var message string
	if c.Process() {
		message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_ID,
			digestMethod, dateTime, i18n.MessageTag_ACCM_POS_MESS_IMP, c.timestamp.Id())
	} else {
		message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_ID,
			digestMethod, dateTime, i18n.MessageTag_ACCM_POS_MESS_IMP, c.timestamp.Id())
	}
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampDigestAlgorithmValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.cvResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(): the generated SubIndication
// member is a pointer, whose nil is Java's null.
func (c *TimestampDigestAlgorithmValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.cvResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.cvResult.Conclusion.SubIndication.SubIndication()
}
