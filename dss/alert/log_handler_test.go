package alert

import (
	"log/slog"
	"testing"
)

func TestLogHandler_RejectsUnsupportedLevel(t *testing.T) {
	h := NewLogHandlerWithLevel[Status](42) // not one of Trace/Debug/Info/Warn/Error

	status := NewMessageStatus()
	status.SetMessage("x")

	if err := h.Process(status); err == nil {
		t.Fatalf("Process() error = nil, want error for unsupported level")
	}
}

func TestLogHandler_AcceptsAllStandardLevels(t *testing.T) {
	status := NewMessageStatus()
	status.SetMessage("x")

	for _, lvl := range []slog.Level{LevelTrace, slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		h := NewLogHandlerWithLevel[Status](lvl)
		if err := h.Process(status); err != nil {
			t.Fatalf("Process() with level %v error = %v, want nil", lvl, err)
		}
	}
}
