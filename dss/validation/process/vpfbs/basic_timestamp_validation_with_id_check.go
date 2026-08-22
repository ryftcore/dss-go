// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftsp/checks/BasicTimestampValidationWithIdCheck.java (DSS 6.5.RC1).
//
// See basic_timestamp_validation_check.go for the package-boundary deviation
// note (this class is filed under vpfbs, not vpftsp, to break an import
// cycle).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicTimestampValidationWithIdCheck verifies time-stamp's basic building
// block and returns its identifier within additional info.
type BasicTimestampValidationWithIdCheck[T any] struct {
	*BasicTimestampValidationCheck[T]
}

// NewBasicTimestampValidationWithIdCheck is the default constructor. Port of
// BasicTimestampValidationWithIdCheck(I18nProvider, T, TimestampWrapper, XmlValidationProcessBasicTimestamp, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' BuildAdditionalInfo rather than the one
// inherited from ChainItemBase.
func NewBasicTimestampValidationWithIdCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp,
	constraint policy.LevelRule) *BasicTimestampValidationWithIdCheck[T] {
	c := &BasicTimestampValidationWithIdCheck[T]{
		BasicTimestampValidationCheck: NewBasicTimestampValidationCheckWithId(i18nProvider, result, timestamp,
			timestampValidationResult, constraint, timestamp.Id()),
	}
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *BasicTimestampValidationWithIdCheck[T]) BuildAdditionalInfo() *string {
	date := process.GetFormattedDate(c.Timestamp.ProductionTime())
	typeTag, err := process.GetTimestampTypeMessageTag(c.Timestamp.Type())
	if err != nil {
		panic(err)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagTimestampValidation, typeTag, c.Timestamp.Id(), date)
	return &message
}
