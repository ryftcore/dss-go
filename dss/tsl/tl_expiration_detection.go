// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/detections/TLExpirationDetection.java (DSS 6.5.RC1).
package tsl

import (
	"time"

	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// TLExpirationDetection detects an expiration of a TrustedList.
type TLExpirationDetection struct{}

var _ alert.AlertDetector[*tslmodel.TLInfo] = (*TLExpirationDetection)(nil)

// NewTLExpirationDetection is the default constructor.
func NewTLExpirationDetection() *TLExpirationDetection {
	return &TLExpirationDetection{}
}

// Detect ports detect(TLInfo).
func (d *TLExpirationDetection) Detect(info *tslmodel.TLInfo) bool {
	if parsingCacheInfo, ok := info.TLParsingCacheInfo(); ok {
		nextUpdateDate := parsingCacheInfo.NextUpdateDate()
		currentDate := time.Now()
		return !nextUpdateDate.IsZero() && nextUpdateDate.Before(currentDate)
	}
	return false
}
