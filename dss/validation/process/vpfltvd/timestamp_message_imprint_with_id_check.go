// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/TimestampMessageImprintWithIdCheck.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): this is a deliberately minimal slice of
// Java's eu.europa.esig.dss.validation.process.vpfltvd package tree (assigned
// to phase 8e - "LTV+qualification" - per PORTING_PLAN.md). Only this one
// check class is ported here, plus its immediate base
// vpftspwatsp.TimestampMessageImprintCheck (see that package's header) - the
// two-level forward dependency the SAV porter flagged. It is the sole caller
// of both classes anywhere in phase 8c:
// bbb/sav's SignatureAcceptanceValidation.contentTimestampMessageImprint().
// Everything else in the real vpfltvd package (RevocationBasicValidationProcess,
// ValidationProcessForSignaturesWithLongTermValidationData, and the rest of
// vpfltvd/checks) remains unported and is left for 8e.
package vpfltvd

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/vpftspwatsp"
)

// TimestampMessageImprintWithIdCheck checks a timestamp's message-imprint and
// returns an Id of the provided token.
type TimestampMessageImprintWithIdCheck[T any] struct {
	*vpftspwatsp.TimestampMessageImprintCheck[T]
}

// NewTimestampMessageImprintWithIdCheck is the default constructor. Port of
// TimestampMessageImprintWithIdCheck(I18nProvider, T, TimestampWrapper, LevelRule).
func NewTimestampMessageImprintWithIdCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TimestampMessageImprintWithIdCheck[T] {
	tokenId := timestamp.Id()
	c := &TimestampMessageImprintWithIdCheck[T]{
		TimestampMessageImprintCheck: vpftspwatsp.NewTimestampMessageImprintCheckWithId(i18nProvider, result,
			timestamp, constraint, &tokenId),
	}
	// Re-register with the outer type so overridden methods (BuildAdditionalInfo
	// below) dispatch correctly - the same "standing bug class" pattern used by
	// bbb/cv's SignatureIntactWithIdCheck (see FRAME's notes).
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TimestampMessageImprintWithIdCheck[T]) BuildAdditionalInfo() *string {
	timestamp := c.Timestamp()
	date := process.GetFormattedDate(timestamp.ProductionTime())
	typeTag, err := process.GetTimestampTypeMessageTag(timestamp.Type())
	if err != nil {
		panic(err)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TIMESTAMP_VALIDATION, typeTag, timestamp.Id(), date)
	return &message
}
