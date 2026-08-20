// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/detections/LOTLLocationChangeDetection.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// LOTLLocationChangeDetection detects the change of a LOTL location.
type LOTLLocationChangeDetection struct {
	// lotlSource is the LOTL source.
	lotlSource *LOTLSource
}

var _ alert.AlertDetector[tslmodel.LOTLInfo] = (*LOTLLocationChangeDetection)(nil)

// NewLOTLLocationChangeDetection is the default constructor.
func NewLOTLLocationChangeDetection(lotlSource *LOTLSource) *LOTLLocationChangeDetection {
	return &LOTLLocationChangeDetection{lotlSource: lotlSource}
}

// Detect ports detect(LOTLInfo).
func (d *LOTLLocationChangeDetection) Detect(info tslmodel.LOTLInfo) bool {
	if d.lotlSource.Url() == info.Url() && d.lotlSource.IsPivotSupport() {
		pivotInfos := info.PivotInfos()
		if len(pivotInfos) > 0 {
			lastPivotInfo := pivotInfos[len(pivotInfos)-1]
			if lastPivotInfo.LOTLLocation() != info.Url() {
				return true
			}
		}
	}
	return false
}
