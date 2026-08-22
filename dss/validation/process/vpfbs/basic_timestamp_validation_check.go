// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftsp/checks/BasicTimestampValidationCheck.java (DSS 6.5.RC1).
//
// PACKAGE-BOUNDARY DEVIATION (LTVB, phase 8e): the manifest places this file's
// goTargetDir at .../vpftsp, matching Java's
// eu.europa.esig.dss.validation.process.vpftsp.checks. It is filed under
// package vpfbs instead, alongside BasicTimestampValidationWithIdCheck, to
// break a Go import cycle the literal placement would create: Java's
// vpftsp.TimestampBasicValidationProcess extends
// vpfbs.AbstractBasicValidationProcess (vpftsp -> vpfbs), while
// vpfbs.AbstractBasicValidationProcess.timestampBasicValidation() and
// vpfltvd.ValidationProcessForSignaturesWithLongTermValidationData both
// construct a vpftsp.checks.BasicTimestampValidationWithIdCheck (vpfbs ->
// vpftsp and vpfltvd -> vpftsp). Go forbids package import cycles outright,
// unlike Java; only this class and its WithId subclass have no dependency of
// their own on the rest of vpftsp (they depend only on process.ChainItem), so
// relocating just these two into vpfbs - the package that actually needs to
// construct them - resolves the cycle while every caller (vpfbs, vpfltvd)
// still finds them, and vpftsp.TimestampBasicValidationProcess itself never
// referenced them to begin with. See abstract_basic_validation_process.go.
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicTimestampValidationCheck checks whether the validation result of EN
// 319 102-1 ch. "5.4 Time-stamp validation building block" process is valid.
type BasicTimestampValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// Timestamp is the timestamp to check. Exported because Java declares the
	// field protected (BasicTimestampValidationWithIdCheck reads it).
	Timestamp *diagnostic.TimestampWrapper

	// timestampValidationResult is the timestamp validation result.
	timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp
}

// NewBasicTimestampValidationCheck is the default constructor. Port of
// BasicTimestampValidationCheck(I18nProvider, T, TimestampWrapper, XmlValidationProcessBasicTimestamp, LevelRule).
func NewBasicTimestampValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp,
	constraint policy.LevelRule) *BasicTimestampValidationCheck[T] {
	return newBasicTimestampValidationCheck(i18nProvider, result, timestamp, timestampValidationResult, constraint, nil)
}

// NewBasicTimestampValidationCheckWithId is the constructor to instantiate the
// check with an Id provided. Port of
// BasicTimestampValidationCheck(I18nProvider, T, TimestampWrapper, XmlValidationProcessBasicTimestamp, LevelRule, String).
func NewBasicTimestampValidationCheckWithId[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp,
	constraint policy.LevelRule, bbbId string) *BasicTimestampValidationCheck[T] {
	return newBasicTimestampValidationCheck(i18nProvider, result, timestamp, timestampValidationResult, constraint, &bbbId)
}

func newBasicTimestampValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp,
	constraint policy.LevelRule, bbbId *string) *BasicTimestampValidationCheck[T] {
	var chainItemBase *process.ChainItemBase[T]
	if bbbId != nil {
		chainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *bbbId)
	} else {
		chainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &BasicTimestampValidationCheck[T]{
		ChainItemBase:             chainItemBase,
		Timestamp:                 timestamp,
		timestampValidationResult: timestampValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *BasicTimestampValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeTSTBBB
}

// Process performs the check. Port of process().
func (c *BasicTimestampValidationCheck[T]) Process() bool {
	return c.IsValid(&c.timestampValidationResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BasicTimestampValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTIBSVPTC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BasicTimestampValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTIBSVPTCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BasicTimestampValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.timestampValidationResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BasicTimestampValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.timestampValidationResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.timestampValidationResult.Conclusion.SubIndication.SubIndication()
}
