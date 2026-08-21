// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogLOTLLocationChangeAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// LogLOTLLocationChangeAlertHandler warns on the LOTL location change.
type LogLOTLLocationChangeAlertHandler struct{}

var _ alert.AlertHandler[*tslmodel.LOTLInfo] = (*LogLOTLLocationChangeAlertHandler)(nil)

// NewLogLOTLLocationChangeAlertHandler is the default constructor.
func NewLogLOTLLocationChangeAlertHandler() *LogLOTLLocationChangeAlertHandler {
	return &LogLOTLLocationChangeAlertHandler{}
}

// Process ports process(LOTLInfo).
func (h *LogLOTLLocationChangeAlertHandler) Process(currentInfo *tslmodel.LOTLInfo) error {
	pivotInfos := currentInfo.PivotInfos()
	if len(pivotInfos) > 0 {
		lastPivotInfo := pivotInfos[len(pivotInfos)-1]
		if lastPivotInfo.LOTLLocation() != currentInfo.Url() {
			slog.Warn("The LOTL Location has changed - new location", "location", lastPivotInfo.LOTLLocation())
		}
	}
	return nil
}
