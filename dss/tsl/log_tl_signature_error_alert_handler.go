// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/handlers/log/LogTLSignatureErrorAlertHandler.java (DSS 6.5.RC1).
package tsl

import (
	"log/slog"

	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// LogTLSignatureErrorAlertHandler warns on TL validation error.
type LogTLSignatureErrorAlertHandler struct{}

var _ alert.AlertHandler[*tslmodel.TLInfo] = (*LogTLSignatureErrorAlertHandler)(nil)

// NewLogTLSignatureErrorAlertHandler is the default constructor.
func NewLogTLSignatureErrorAlertHandler() *LogTLSignatureErrorAlertHandler {
	return &LogTLSignatureErrorAlertHandler{}
}

// Process ports process(TLInfo).
func (h *LogTLSignatureErrorAlertHandler) Process(currentInfo *tslmodel.TLInfo) error {
	slog.Warn("There is a problem in the TL signature", "url", currentInfo.Url())
	return nil
}
