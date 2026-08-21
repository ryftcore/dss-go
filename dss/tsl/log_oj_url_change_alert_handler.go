// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogOJUrlChangeAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// LogOJUrlChangeAlertHandler warns on the LOTL Official Journal URL change.
type LogOJUrlChangeAlertHandler struct{}

var _ alert.AlertHandler[*tslmodel.LOTLInfo] = (*LogOJUrlChangeAlertHandler)(nil)

// NewLogOJUrlChangeAlertHandler is the default constructor.
func NewLogOJUrlChangeAlertHandler() *LogOJUrlChangeAlertHandler {
	return &LogOJUrlChangeAlertHandler{}
}

// Process ports process(LOTLInfo).
func (h *LogOJUrlChangeAlertHandler) Process(currentInfo *tslmodel.LOTLInfo) error {
	if parsingCacheInfo, ok := currentInfo.TLParsingCacheInfo(); ok {
		slog.Warn("The Official Journal URL has changed - new location", "location", parsingCacheInfo.SigningCertificateAnnouncementUrl())
	} else {
		slog.Warn("No parsing result found for a LOTL", "url", currentInfo.Url())
	}
	return nil
}
