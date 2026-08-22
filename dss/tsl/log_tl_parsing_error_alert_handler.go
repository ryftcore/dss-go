// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogTLParsingErrorAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// LogTLParsingErrorAlertHandler warns on TL parsing error.
type LogTLParsingErrorAlertHandler struct{}

var _ alert.AlertHandler[*tslmodel.TLInfo] = (*LogTLParsingErrorAlertHandler)(nil)

// NewLogTLParsingErrorAlertHandler is the default constructor.
func NewLogTLParsingErrorAlertHandler() *LogTLParsingErrorAlertHandler {
	return &LogTLParsingErrorAlertHandler{}
}

// Process ports process(TLInfo).
func (h *LogTLParsingErrorAlertHandler) Process(currentInfo *tslmodel.TLInfo) error {
	slog.Warn("There was an error while parsing a TL", "url", currentInfo.Url())
	return nil
}
