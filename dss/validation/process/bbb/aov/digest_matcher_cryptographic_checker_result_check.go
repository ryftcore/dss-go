// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/DigestMatcherCryptographicCheckerResultCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go. Java's
// class lives in aov.cc.checks and extends aov.checks.DigestAlgorithmCryptographicCheckerResultCheck
// across packages; flattening every aov subpackage into this one pkg aov
// removes that cross-package edge.
package aov

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// DigestMatcherCryptographicCheckerResultCheck verifies a DigestMatcher.
type DigestMatcherCryptographicCheckerResultCheck[T any] struct {
	*DigestAlgorithmCryptographicCheckerResultCheck[T]

	// referenceNames is the list of verifying reference names.
	referenceNames []string
}

// NewDigestMatcherCryptographicCheckerResultCheck is the default constructor.
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' BuildAdditionalInfo rather than the one
// inherited from DigestAlgorithmCryptographicCheckerResultCheck.
func NewDigestMatcherCryptographicCheckerResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validationDate time.Time, position i18n.MessageTag, referenceNames []string, ccResult *jaxb.XmlCC,
	constraint policy.LevelRule) *DigestMatcherCryptographicCheckerResultCheck[T] {
	c := &DigestMatcherCryptographicCheckerResultCheck[T]{
		DigestAlgorithmCryptographicCheckerResultCheck: NewDigestAlgorithmCryptographicCheckerResultCheck(
			i18nProvider, result, validationDate, position, ccResult, constraint),
		referenceNames: referenceNames,
	}
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *DigestMatcherCryptographicCheckerResultCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.validationDate)
	var message string
	if c.IsValid(&c.ccResult.XmlConstraintsConclusionContent) {
		switch utils.CollectionSize(c.referenceNames) {
		case 0:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM,
				c.ccResult.CryptographicValidation.Algorithm.Name, dateTime, c.position)
		case 1:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAME,
				c.ccResult.CryptographicValidation.Algorithm.Name, dateTime, c.position, c.referenceNames[0])
		default:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAMES,
				c.ccResult.CryptographicValidation.Algorithm.Name, dateTime, c.position, utils.JoinStrings(c.referenceNames, ", "))
		}
	} else {
		switch utils.CollectionSize(c.referenceNames) {
		case 0:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF, c.ErrorMessage(), dateTime)
		case 1:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAME,
				c.ErrorMessage(), dateTime, c.referenceNames[0])
		default:
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAMES,
				c.ErrorMessage(), dateTime, utils.JoinStrings(c.referenceNames, ", "))
		}
	}
	return &message
}
