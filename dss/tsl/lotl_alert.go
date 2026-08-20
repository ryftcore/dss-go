// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/LOTLAlert.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see xml_download_result.go's header): extends
// job.DocumentAlert[*tslmodel.LOTLInfo, *tslmodel.LOTLInfo], itself embedding alert.AbstractAlert[T]
// (dss/alert, already ported).
package tsl

import (
	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/validation/job"
)

// LOTLAlert processes events on LOTL.
type LOTLAlert struct {
	job.DocumentAlert[*tslmodel.LOTLInfo, *tslmodel.LOTLInfo]
}

var _ alert.Alert[tslmodel.LOTLInfo] = (*LOTLAlert)(nil)

// NewLOTLAlert is the default constructor.
func NewLOTLAlert(detection alert.AlertDetector[tslmodel.LOTLInfo], handler alert.AlertHandler[tslmodel.LOTLInfo]) *LOTLAlert {
	return &LOTLAlert{DocumentAlert: job.NewDocumentAlert(detection, handler)}
}
