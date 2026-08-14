// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/LogOnStatusAlert.java (DSS 6.5.RC1).
package alert

import "log/slog"

// LogOnStatusAlert logs a message on a Status alert.
type LogOnStatusAlert struct {
	*AbstractStatusAlert
}

// NewLogOnStatusAlert creates a LogOnStatusAlert logging at Java's default level, WARN.
func NewLogOnStatusAlert() *LogOnStatusAlert {
	return NewLogOnStatusAlertWithLevel(slog.LevelWarn)
}

// NewLogOnStatusAlertWithLevel creates a LogOnStatusAlert logging at the given level.
func NewLogOnStatusAlertWithLevel(level slog.Level) *LogOnStatusAlert {
	return &LogOnStatusAlert{
		AbstractStatusAlert: NewAbstractStatusAlert(NewLogHandlerWithLevel[Status](level)),
	}
}
