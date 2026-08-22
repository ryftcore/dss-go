// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftspwatsp/checks/TimestampMessageImprintCheck.java (DSS 6.5.RC1).
//
// This is a deliberately minimal slice of Java's
// eu.europa.esig.dss.validation.process.vpftspwatsp package tree. Only
// this one check class is ported here: it is a plain leaf ChainItem with no
// dependency on the rest of vpftspwatsp (the "5.6.2.4 Validation process for
// timestamps" orchestration classes), and it is the base class
// vpfltvd.TimestampMessageImprintWithIdCheck (used unconditionally by
// bbb/sav's SignatureAcceptanceValidation.contentTimestampMessageImprint())
// needs to extend.
//
// Package placement deviation: Java's vpftspwatsp.checks is a package of its
// own, distinct from vpftspwatsp; every other ...:checks subpackage in this
// port flattens into its parent, but doing that here would close an import
// cycle. bbb/sav imports vpfltvd for TimestampMessageImprintWithIdCheck,
// vpfltvd imports this class, and the rest of vpftspwatsp
// (ValidationProcessForTimestampsWithArchivalData and the three checks around
// it) imports vpfswatsp, which imports vpfltvd and bbb/sav back. This one
// class - a leaf ChainItem needing nothing but process and diagnostic -
// therefore keeps Java's own vpftspwatsp/checks package boundary, while
// everything else of vpftspwatsp is flattened into vpftspwatsp. The same
// relocation technique vpfbs and vpfltvdsig use for their own cycles.
package checks

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampMessageImprintCheck checks message-imprint validity for a
// timestamp token.
type TimestampMessageImprintCheck[T any] struct {
	*process.ChainItemBase[T]

	// timestamp is the timestamp to check.
	timestamp *diagnostic.TimestampWrapper
}

// NewTimestampMessageImprintCheck is the default constructor. Port of
// TimestampMessageImprintCheck(Provider, T, TimestampWrapper, LevelRule).
func NewTimestampMessageImprintCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TimestampMessageImprintCheck[T] {
	return NewTimestampMessageImprintCheckWithId(i18nProvider, result, timestamp, constraint, nil)
}

// NewTimestampMessageImprintCheckWithId is the constructor accepting an
// explicit bbbId, for subclasses. Port of the protected
// TimestampMessageImprintCheck(Provider, T, TimestampWrapper, LevelRule, String).
// A nil bbbId matches the Java null passed by the public constructor.
func NewTimestampMessageImprintCheckWithId[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, constraint policy.LevelRule, bbbId *string) *TimestampMessageImprintCheck[T] {
	var base *process.ChainItemBase[T]
	if bbbId != nil {
		base = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *bbbId)
	} else {
		base = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &TimestampMessageImprintCheck[T]{
		ChainItemBase: base,
		timestamp:     timestamp,
	}
	c.InitChainItem(c)
	return c
}

// Timestamp returns the timestamp under check, for subclasses in other
// packages (Go has no protected-field access across packages).
func (c *TimestampMessageImprintCheck[T]) Timestamp() *diagnostic.TimestampWrapper {
	return c.timestamp
}

// Process performs the check. Port of process().
func (c *TimestampMessageImprintCheck[T]) Process() bool {
	return c.timestamp.IsMessageImprintDataFound() && c.timestamp.IsMessageImprintDataIntact()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampMessageImprintCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVDMICTSTMCMI
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampMessageImprintCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVDMICTSTMCMIANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port
// of getFailedIndicationForConclusion().
func (c *TimestampMessageImprintCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampMessageImprintCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
