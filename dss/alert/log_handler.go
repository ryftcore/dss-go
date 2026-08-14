// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/handler/LogHandler.java (DSS 6.5.RC1).
package alert

import (
	"context"
	"fmt"
	"log/slog"
)

// LevelTrace is a slog level below slog.LevelDebug, mirroring slf4j's Level.TRACE (which
// log/slog does not define out of the box). It sits one "step" (4, following slog's own
// spacing between Debug/Info/Warn/Error) below LevelDebug.
const LevelTrace slog.Level = slog.LevelDebug - 4

// LogHandler is an AlertHandler which logs the object at the configured slog.Level.
type LogHandler[T any] struct {
	level  slog.Level
	logger *slog.Logger
}

// NewLogHandler creates a LogHandler logging at Java's default level, WARN.
func NewLogHandler[T any]() *LogHandler[T] {
	return NewLogHandlerWithLevel[T](slog.LevelWarn)
}

// NewLogHandlerWithLevel creates a LogHandler logging at the given level.
func NewLogHandlerWithLevel[T any](level slog.Level) *LogHandler[T] {
	return &LogHandler[T]{level: level, logger: slog.Default()}
}

// Process logs object's string representation at the configured level. Levels outside
// {TRACE, DEBUG, INFO, WARN, ERROR} report an error, mirroring Java's IllegalArgumentException.
func (h *LogHandler[T]) Process(object T) error {
	switch h.level {
	case LevelTrace, slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError:
		logger := h.logger
		if logger == nil {
			logger = slog.Default()
		}
		logger.Log(context.Background(), h.level, fmt.Sprint(object))
		return nil
	default:
		return fmt.Errorf("the LogLevel [%v] is not allowed configuration", h.level)
	}
}
