// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/LOTLAlert.java (DSS 6.5.RC1).
//
// Extends job.DocumentAlert[*tslmodel.LOTLInfo, *tslmodel.LOTLInfo] (see
// xml_download_result.go's header for the wider job.* convention), itself embedding
// alert.AbstractAlert[T] (dss/alert).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// LOTLAlert processes events on LOTL.
type LOTLAlert struct {
	*job.DocumentAlert[*tslmodel.LOTLInfo, *tslmodel.LOTLInfo]
}

var _ alert.Alert[*tslmodel.LOTLInfo] = (*LOTLAlert)(nil)

// NewLOTLAlert is the default constructor.
func NewLOTLAlert(detection alert.Detector[*tslmodel.LOTLInfo], handler alert.Handler[*tslmodel.LOTLInfo]) *LOTLAlert {
	return &LOTLAlert{DocumentAlert: job.NewDocumentAlert(detection, handler)}
}
