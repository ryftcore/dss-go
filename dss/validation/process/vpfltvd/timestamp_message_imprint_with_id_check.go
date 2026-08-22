// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/TimestampMessageImprintWithIdCheck.java (DSS 6.5.RC1).
//
// This check class's immediate base,
// vpftspwatsp/checks.TimestampMessageImprintCheck (see that package's
// header), is a two-level dependency. This is the sole caller of both
// classes: bbb/sav's SignatureAcceptanceValidation.contentTimestampMessageImprint().
// Everything else in the real vpfltvd package tree
// (RevocationBasicValidationProcess and the rest of vpfltvd/checks) also
// lives in this same Go package; ValidationProcessForSignaturesWithLongTermValidationData
// is filed under the sibling package vpfltvdsig instead - see that
// package's header for why.
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
// TimestampMessageImprintWithIdCheck(Provider, T, TimestampWrapper, LevelRule).
func NewTimestampMessageImprintWithIdCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
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
