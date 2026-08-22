// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/detections/TLSignatureErrorDetection.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// TLSignatureErrorDetection detects if an error in a TL validation occurred.
type TLSignatureErrorDetection struct{}

var _ alert.AlertDetector[*tslmodel.TLInfo] = (*TLSignatureErrorDetection)(nil)

// NewTLSignatureErrorDetection is the default constructor.
func NewTLSignatureErrorDetection() *TLSignatureErrorDetection {
	return &TLSignatureErrorDetection{}
}

// Detect ports detect(TLInfo).
func (d *TLSignatureErrorDetection) Detect(info *tslmodel.TLInfo) bool {
	downloadCacheInfo := info.DownloadCacheInfo()
	if downloadCacheInfo != nil && downloadCacheInfo.IsDesynchronized() {
		validationCacheInfo := info.ValidationCacheInfo()
		if validationCacheInfo != nil && !validationCacheInfo.IsValid() {
			return true
		}
	}
	return false
}
