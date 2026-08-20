// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogOJUrlChangeAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// LogOJUrlChangeAlertHandler warns on the LOTL Official Journal URL change.
type LogOJUrlChangeAlertHandler struct{}

var _ alert.AlertHandler[tslmodel.LOTLInfo] = (*LogOJUrlChangeAlertHandler)(nil)

// NewLogOJUrlChangeAlertHandler is the default constructor.
func NewLogOJUrlChangeAlertHandler() *LogOJUrlChangeAlertHandler {
	return &LogOJUrlChangeAlertHandler{}
}

// Process ports process(LOTLInfo).
func (h *LogOJUrlChangeAlertHandler) Process(currentInfo tslmodel.LOTLInfo) error {
	if currentInfo.ParsingCacheInfo() != nil {
		slog.Warn("The Official Journal URL has changed - new location", "location", currentInfo.ParsingCacheInfo().SigningCertificateAnnouncementUrl())
	} else {
		slog.Warn("No parsing result found for a LOTL", "url", currentInfo.Url())
	}
	return nil
}
