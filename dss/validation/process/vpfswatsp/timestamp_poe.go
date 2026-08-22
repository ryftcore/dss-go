// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/TimestampPOE.java (DSS 6.5.RC1).
//
// See poe.go for the POE hierarchy note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TimestampPOE is a POE provided by a time-stamp token.
type TimestampPOE struct {
	*POEBase

	// timestampWrapper is the timestamp.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewTimestampPOE is the constructor to instantiate POE by a timestamp. Port of
// TimestampPOE(TimestampWrapper).
func NewTimestampPOE(timestampWrapper *diagnostic.TimestampWrapper) *TimestampPOE {
	return &TimestampPOE{
		POEBase:          NewPOE(timestampPOETime(timestampWrapper)),
		timestampWrapper: timestampWrapper,
	}
}

// timestampPOETime ports the private static getPOETime(TimestampWrapper).
//
// Java raises Objects.requireNonNull on a null wrapper; the null production time
// it may return then trips POE(Date)'s own requireNonNull. Both are raised here
// with the upstream messages (see NewPOE).
func timestampPOETime(timestampWrapper *diagnostic.TimestampWrapper) time.Time {
	if timestampWrapper == nil {
		panic("The timestampWrapper must be defined!")
	}
	productionTime := timestampWrapper.ProductionTime()
	if productionTime == nil {
		panic("The controlTime must be defined!")
	}
	return *productionTime
}

// POEProviderId returns the time-stamp's Id. Port of the overridden
// getPOEProviderId().
func (p *TimestampPOE) POEProviderId() *string {
	id := p.timestampWrapper.Id()
	return &id
}

// TimestampType returns timestamp type if the POE defined by a timestamp. Port
// of getTimestampType(): Java's documented NULL for a control-time POE cannot
// occur on this type, and the wrapper's own getType() maps a missing type to
// the empty TimestampType.
func (p *TimestampPOE) TimestampType() enumerations.TimestampType {
	return p.timestampWrapper.Type()
}

// POEObjects returns the objects covered by the time-stamp. Port of the
// overridden getPOEObjects().
func (p *TimestampPOE) POEObjects() []*diagnosticjaxb.XmlTimestampedObject {
	return p.timestampWrapper.TimestampedObjects()
}

// IsTokenProvided returns whether the POE is provided by a token. Port of the
// overridden isTokenProvided().
func (p *TimestampPOE) IsTokenProvided() bool {
	return true
}
