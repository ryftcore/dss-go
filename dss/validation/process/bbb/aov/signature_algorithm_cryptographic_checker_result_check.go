// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/SignatureAlgorithmCryptographicCheckerResultCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: Java's eu.europa.esig.dss.validation.process.bbb.aov.checks
// package is flattened into this pkg aov, per the phase 8d porter brief.
package aov

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureAlgorithmCryptographicCheckerResultCheck validates the result of a
// Signature Algorithm cryptographic checker.
type SignatureAlgorithmCryptographicCheckerResultCheck[T any] struct {
	*AbstractCryptographicCheckerResultCheck[T]

	// validationDate is the validation time.
	validationDate time.Time

	// context is the validation context.
	context enumerations.Context
}

// NewSignatureAlgorithmCryptographicCheckerResultCheck is the default
// constructor. Port of SignatureAlgorithmCryptographicCheckerResultCheck(
// I18nProvider, T, Date, MessageTag, XmlCC, LevelRule): a Java caller passing
// a null context maps to this constructor.
func NewSignatureAlgorithmCryptographicCheckerResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validationDate time.Time, position i18n.MessageTag, ccResult *jaxb.XmlCC,
	constraint policy.LevelRule) *SignatureAlgorithmCryptographicCheckerResultCheck[T] {
	return newSignatureAlgorithmCryptographicCheckerResultCheck(i18nProvider, result, validationDate, "", position,
		ccResult, constraint, "")
}

// NewSignatureAlgorithmCryptographicCheckerResultCheckWithContext is the full
// constructor. Port of SignatureAlgorithmCryptographicCheckerResultCheck(
// I18nProvider, T, Date, Context, MessageTag, XmlCC, LevelRule, String).
func NewSignatureAlgorithmCryptographicCheckerResultCheckWithContext[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validationDate time.Time, context enumerations.Context, position i18n.MessageTag, ccResult *jaxb.XmlCC,
	constraint policy.LevelRule, tokenId string) *SignatureAlgorithmCryptographicCheckerResultCheck[T] {
	return newSignatureAlgorithmCryptographicCheckerResultCheck(i18nProvider, result, validationDate, context, position,
		ccResult, constraint, tokenId)
}

func newSignatureAlgorithmCryptographicCheckerResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validationDate time.Time, context enumerations.Context, position i18n.MessageTag, ccResult *jaxb.XmlCC,
	constraint policy.LevelRule, tokenId string) *SignatureAlgorithmCryptographicCheckerResultCheck[T] {
	c := &SignatureAlgorithmCryptographicCheckerResultCheck[T]{
		validationDate: validationDate,
		context:        context,
	}
	if tokenId == "" {
		c.AbstractCryptographicCheckerResultCheck = NewAbstractCryptographicCheckerResultCheck(i18nProvider, result, position, ccResult, constraint)
	} else {
		c.AbstractCryptographicCheckerResultCheck = NewAbstractCryptographicCheckerResultCheckWithId(i18nProvider, result, position, ccResult, constraint, tokenId)
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of the overridden
// getBlockType().
func (c *SignatureAlgorithmCryptographicCheckerResultCheck[T]) BlockType() jaxb.XmlBlockType {
	if enumerations.Context_CERTIFICATE == c.context {
		return jaxb.XmlBlockType_AOV_XCV
	}
	return ""
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *SignatureAlgorithmCryptographicCheckerResultCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.validationDate)
	if c.IsValid(&c.ccResult.XmlConstraintsConclusionContent) {
		cryptographicValidation := c.ccResult.CryptographicValidation
		algorithm := cryptographicValidation.Algorithm
		var message string
		if algorithm.KeyLength != nil && utils.IsStringNotEmpty(*algorithm.KeyLength) {
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE,
				algorithm.Name, *algorithm.KeyLength, dateTime)
		} else {
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS, algorithm.Name, dateTime)
		}
		return &message
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE, c.ErrorMessage(), dateTime)
	return &message
}
