// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogTLExpirationAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// LogTLExpirationAlertHandler warns on the TL expiration.
type LogTLExpirationAlertHandler struct{}

var _ alert.AlertHandler[tslmodel.TLInfo] = (*LogTLExpirationAlertHandler)(nil)

// NewLogTLExpirationAlertHandler is the default constructor.
func NewLogTLExpirationAlertHandler() *LogTLExpirationAlertHandler {
	return &LogTLExpirationAlertHandler{}
}

// Process ports process(TLInfo).
func (h *LogTLExpirationAlertHandler) Process(currentInfo tslmodel.TLInfo) error {
	if currentInfo.ParsingCacheInfo() != nil {
		slog.Warn("The TL has expired", "url", currentInfo.Url(), "lastUpdate", currentInfo.ParsingCacheInfo().NextUpdateDate())
	} else {
		slog.Warn("No parsing result found for a LOTL", "url", currentInfo.Url())
	}
	return nil
}
