// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/TimestampCoherenceOrderCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampCoherenceOrderCheck checks if the timestamp's order is coherent.
type TimestampCoherenceOrderCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessLongTermData]

	// timestamps is the list of timestamps to check.
	timestamps []*diagnostic.TimestampWrapper
}

// NewTimestampCoherenceOrderCheck is the default constructor. Port of
// TimestampCoherenceOrderCheck(I18nProvider, XmlValidationProcessLongTermData, List, LevelRule).
func NewTimestampCoherenceOrderCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData], timestamps []*diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *TimestampCoherenceOrderCheck {
	c := &TimestampCoherenceOrderCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		timestamps:    timestamps,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TimestampCoherenceOrderCheck) Process() bool {
	return utils.CollectionSize(c.timestamps) <= 1 || c.checkTimestampCoherenceOrderByType()
}

// checkTimestampCoherenceOrderByType ports the private
// checkTimestampCoherenceOrderByType().
func (c *TimestampCoherenceOrderCheck) checkTimestampCoherenceOrderByType() bool {
	toBeCheckedTimestamps := make([]*diagnostic.TimestampWrapper, len(c.timestamps))
	copy(toBeCheckedTimestamps, c.timestamps)
	for len(toBeCheckedTimestamps) > 0 {
		timestamp := toBeCheckedTimestamps[0]
		toBeCheckedTimestamps = toBeCheckedTimestamps[1:] // in order do not re-validate the same pairs
		if !c.isValidAgainstList(timestamp, toBeCheckedTimestamps) {
			return false
		}
	}
	return true
}

// isValidAgainstList ports the private
// isValidAgainstList(TimestampWrapper, List).
func (c *TimestampCoherenceOrderCheck) isValidAgainstList(timestamp *diagnostic.TimestampWrapper,
	timestampList []*diagnostic.TimestampWrapper) bool {
	for _, timestampToCompare := range timestampList {
		typeResult := timestamp.Type().Compare(timestampToCompare.Type())

		var productionTimeResult int
		pt, ptc := timestamp.ProductionTime(), timestampToCompare.ProductionTime()
		switch {
		case pt == nil && ptc == nil:
			productionTimeResult = 0
		case pt == nil:
			productionTimeResult = -1
		case ptc == nil:
			productionTimeResult = 1
		case pt.Before(*ptc):
			productionTimeResult = -1
		case pt.After(*ptc):
			productionTimeResult = 1
		default:
			productionTimeResult = 0
		}

		// if time is different and types are in a wrong order
		if productionTimeResult != 0 && typeResult != productionTimeResult {
			// if the type is the same, but time is different, check the references
			if typeResult == 0 {
				// if the first timestamp is created earlier and it covers the next timestamp
				if productionTimeResult < 0 && c.coversTheTimestamp(timestamp, timestampToCompare) {
					return false
					// if the first timestamp is created after and its covered by the previous timestamp
				} else if productionTimeResult > 0 && c.coversTheTimestamp(timestampToCompare, timestamp) {
					return false
				}

			} else {
				return false
			}
		}
	}
	return true
}

// coversTheTimestamp ports the private
// coversTheTimestamp(TimestampWrapper, TimestampWrapper).
func (c *TimestampCoherenceOrderCheck) coversTheTimestamp(timestamp *diagnostic.TimestampWrapper,
	timestampToCompare *diagnostic.TimestampWrapper) bool {
	for _, timestamped := range timestamp.TimestampedTimestamps() {
		if timestamped.Id() == timestampToCompare.Id() {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampCoherenceOrderCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTSVASTPTCT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampCoherenceOrderCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTSVASTPTCTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampCoherenceOrderCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampCoherenceOrderCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTimestampOrderFailure
}
