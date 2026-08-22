// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/TimestampMessageImprintWithIdCheck.java (DSS 6.5.RC1).
//
// This is a deliberately minimal slice of Java's
// eu.europa.esig.dss.validation.process.vpfltvd package tree. Only this one
// check class is ported here, plus its immediate base
// vpftspwatsp/checks.TimestampMessageImprintCheck (see that package's header) - the
// two-level forward dependency the SAV porter flagged. It is the sole caller
// of both classes anywhere:
// bbb/sav's SignatureAcceptanceValidation.contentTimestampMessageImprint().
// Everything else in the real vpfltvd package (RevocationBasicValidationProcess,
// ValidationProcessForSignaturesWithLongTermValidationData, and the rest of
// vpfltvd/checks) remains unported and is left for 8e.
package vpfltvd

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	vpftspwatspchecks "github.com/ryftcore/dss-go/dss/validation/process/vpftspwatsp/checks"
)

// TimestampMessageImprintWithIdCheck checks a timestamp's message-imprint and
// returns an Id of the provided token.
type TimestampMessageImprintWithIdCheck[T any] struct {
	*vpftspwatspchecks.TimestampMessageImprintCheck[T]
}

// NewTimestampMessageImprintWithIdCheck is the default constructor. Port of
// TimestampMessageImprintWithIdCheck(I18nProvider, T, TimestampWrapper, LevelRule).
func NewTimestampMessageImprintWithIdCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TimestampMessageImprintWithIdCheck[T] {
	tokenId := timestamp.Id()
	c := &TimestampMessageImprintWithIdCheck[T]{
		TimestampMessageImprintCheck: vpftspwatspchecks.NewTimestampMessageImprintCheckWithId(i18nProvider, result,
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
	message := c.I18nProvider.GetMessage(i18n.MessageTagTimestampValidation, typeTag, timestamp.Id(), date)
	return &message
}
