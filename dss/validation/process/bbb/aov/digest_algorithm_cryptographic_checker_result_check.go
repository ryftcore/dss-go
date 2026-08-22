// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/DigestAlgorithmCryptographicCheckerResultCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DigestAlgorithmCryptographicCheckerResultCheck validates a Digest
// cryptographic constraint.
type DigestAlgorithmCryptographicCheckerResultCheck[T any] struct {
	*AbstractCryptographicCheckerResultCheck[T]

	// validationDate is the validation time.
	validationDate time.Time
}

// NewDigestAlgorithmCryptographicCheckerResultCheck is the default
// constructor.
func NewDigestAlgorithmCryptographicCheckerResultCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	validationDate time.Time, position i18n.MessageTag, ccResult *jaxb.XmlCC,
	constraint policy.LevelRule) *DigestAlgorithmCryptographicCheckerResultCheck[T] {
	c := &DigestAlgorithmCryptographicCheckerResultCheck[T]{validationDate: validationDate}
	c.AbstractCryptographicCheckerResultCheck = NewAbstractCryptographicCheckerResultCheck(i18nProvider, result, position, ccResult, constraint)
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *DigestAlgorithmCryptographicCheckerResultCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.validationDate)
	if c.IsValid(&c.ccResult.XmlConstraintsConclusionContent) {
		message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckSuccessDM,
			c.ccResult.CryptographicValidation.Algorithm.Name, dateTime, c.position)
		return &message
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckFailureWithRef, c.ErrorMessage(), dateTime)
	return &message
}
