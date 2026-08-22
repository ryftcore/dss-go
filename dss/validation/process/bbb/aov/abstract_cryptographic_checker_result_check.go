// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/AbstractCryptographicCheckerResultCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go. Java's
// AbstractCryptographicCheckerResultCheck lives in aov.cc.checks and is
// extended cross-package by aov.checks.SignatureAlgorithmCryptographicCheckerResultCheck
// / aov.checks.DigestAlgorithmCryptographicCheckerResultCheck; flattening
// every aov subpackage into this one pkg aov removes that cross-package edge.
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractCryptographicCheckerResultCheck performs cryptographic validation.
// A concrete check embeds it instead of process.UninterruptedChainItemBase
// and registers itself with InitChainItem, the same pattern as
// bbb.AbstractMultiValuesCheckItem.
type AbstractCryptographicCheckerResultCheck[T any] struct {
	*process.UninterruptedChainItemBase[T]

	// position is the cryptographic constraint position to be validated.
	position i18n.MessageTag

	// ccResult is the Cryptographic Check result.
	ccResult *jaxb.XmlCC

	// checkerResultMessage is the checker result message.
	checkerResultMessage *jaxb.XmlMessage
}

// NewAbstractCryptographicCheckerResultCheck is the default constructor. Port
// of AbstractCryptographicCheckerResultCheck(I18nProvider, T, MessageTag,
// XmlCC, LevelRule): a Java caller passing a null tokenId maps to this
// constructor, matching process.NewChainItemBase's own convention.
func NewAbstractCryptographicCheckerResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	position i18n.MessageTag, ccResult *jaxb.XmlCC, constraint policy.LevelRule) *AbstractCryptographicCheckerResultCheck[T] {
	return &AbstractCryptographicCheckerResultCheck[T]{
		UninterruptedChainItemBase: process.NewUninterruptedChainItemBase(i18nProvider, result, constraint),
		position:                   position,
		ccResult:                   ccResult,
		checkerResultMessage:       extractCCXmlMessage(ccResult, constraint),
	}
}

// NewAbstractCryptographicCheckerResultCheckWithId is the constructor
// carrying a token identifier. Port of
// AbstractCryptographicCheckerResultCheck(I18nProvider, T, MessageTag, XmlCC, LevelRule, String).
func NewAbstractCryptographicCheckerResultCheckWithId[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	position i18n.MessageTag, ccResult *jaxb.XmlCC, constraint policy.LevelRule, tokenId string) *AbstractCryptographicCheckerResultCheck[T] {
	return &AbstractCryptographicCheckerResultCheck[T]{
		UninterruptedChainItemBase: process.NewUninterruptedChainItemBaseWithId(i18nProvider, result, constraint, tokenId),
		position:                   position,
		ccResult:                   ccResult,
		checkerResultMessage:       extractCCXmlMessage(ccResult, constraint),
	}
}

// extractCCXmlMessage ports the private static extractXmlMessage(XmlCC, LevelRule).
func extractCCXmlMessage(ccResult *jaxb.XmlCC, constraint policy.LevelRule) *jaxb.XmlMessage {
	conclusion := ccResult.Conclusion
	if conclusion != nil && constraint != nil && constraint.Level() != "" {
		// Collects messages from all levels (required for generic crypto check)
		var messages []*jaxb.XmlMessage
		messages = append(messages, conclusion.Errors...)
		messages = append(messages, conclusion.Warnings...)
		messages = append(messages, conclusion.Infos...)
		if utils.IsCollectionNotEmpty(messages) {
			return messages[0] // take the first one
		}
	}
	return nil
}

// Process performs the check. Port of the overridden process().
func (c *AbstractCryptographicCheckerResultCheck[T]) Process() bool {
	return c.IsValid(&c.ccResult.XmlConstraintsConclusionContent) && allCCConstraintsValid(c.ccResult)
}

// allCCConstraintsValid ports the private allConstraintsValid(XmlConstraintsConclusion).
func allCCConstraintsValid(result *jaxb.XmlCC) bool {
	for _, constraint := range result.Constraint {
		if jaxb.XmlStatus_OK != constraint.Status && jaxb.XmlStatus_IGNORED != constraint.Status {
			return false
		}
	}
	return true
}

// Level returns an execution Level of the chain item. Port of the overridden
// getLevel().
func (c *AbstractCryptographicCheckerResultCheck[T]) Level() enumerations.Level {
	conclusion := c.ccResult.Conclusion
	if conclusion != nil {
		if utils.IsCollectionNotEmpty(conclusion.Errors) {
			return enumerations.Level_FAIL
		} else if utils.IsCollectionNotEmpty(conclusion.Warnings) {
			return enumerations.Level_WARN
		} else if utils.IsCollectionNotEmpty(conclusion.Infos) {
			return enumerations.Level_INFORM
		}
	}
	return c.UninterruptedChainItemBase.Level()
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *AbstractCryptographicCheckerResultCheck[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ACCM, c.position)
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *AbstractCryptographicCheckerResultCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.checkerResultMessage
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port
// of the overridden getFailedIndicationForConclusion().
func (c *AbstractCryptographicCheckerResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.ccResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion().
func (c *AbstractCryptographicCheckerResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.ccResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.ccResult.Conclusion.SubIndication.SubIndication()
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of the overridden getPreviousErrors().
func (c *AbstractCryptographicCheckerResultCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.ccResult.Conclusion.Errors
}

// ErrorMessage gets error message. Port of the protected getErrorMessage():
// buildErrorMessage() is fixed (never further overridden by a subclass in
// this port), so it is read directly rather than dispatched through the
// overrides interface.
func (c *AbstractCryptographicCheckerResultCheck[T]) ErrorMessage() string {
	if c.checkerResultMessage != nil {
		return c.checkerResultMessage.Value
	}
	return utils.EmptyString
}
