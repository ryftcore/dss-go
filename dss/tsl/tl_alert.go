// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/TLAlert.java (DSS 6.5.RC1).
//
// Extends job.DocumentAlert[*tslmodel.TLInfo, *tslmodel.LOTLInfo] (see lotl_alert.go's header
// for the wider job.* convention).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// TLAlert processes events on TL.
type TLAlert struct {
	*job.DocumentAlert[*tslmodel.TLInfo, *tslmodel.LOTLInfo]
}

var _ alert.Alert[*tslmodel.TLInfo] = (*TLAlert)(nil)

// NewTLAlert is the default constructor.
func NewTLAlert(detection alert.AlertDetector[*tslmodel.TLInfo], handler alert.AlertHandler[*tslmodel.TLInfo]) *TLAlert {
	return &TLAlert{DocumentAlert: job.NewDocumentAlert(detection, handler)}
}
