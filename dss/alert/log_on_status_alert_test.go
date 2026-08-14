package alert

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestLogOnStatusAlert_LogsAtWarnByDefault(t *testing.T) {
	var buf bytes.Buffer
	restore := setDefaultLogger(&buf, slog.LevelWarn)
	defer restore()

	a := NewLogOnStatusAlert()
	status := NewMessageStatus()
	status.SetMessage("default level warning")

	if err := a.Alert(status); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	out := buf.String()
	if !strings.Contains(out, "default level warning") {
		t.Fatalf("log output = %q, want it to contain the status message", out)
	}
	if !strings.Contains(out, "level=WARN") {
		t.Fatalf("log output = %q, want WARN level", out)
	}
}

func TestLogOnStatusAlert_LogsAtGivenLevel(t *testing.T) {
	var buf bytes.Buffer
	restore := setDefaultLogger(&buf, slog.LevelDebug)
	defer restore()

	a := NewLogOnStatusAlertWithLevel(slog.LevelInfo)
	status := NewMessageStatus()
	status.SetMessage("info level message")

	if err := a.Alert(status); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=INFO") {
		t.Fatalf("log output = %q, want INFO level", out)
	}
}

func TestLogOnStatusAlert_NoLogWhenStatusEmpty(t *testing.T) {
	var buf bytes.Buffer
	restore := setDefaultLogger(&buf, slog.LevelWarn)
	defer restore()

	a := NewLogOnStatusAlert()
	if err := a.Alert(NewMessageStatus()); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("log output = %q, want empty (status was empty, detector should not fire)", buf.String())
	}
}

// setDefaultLogger swaps slog's default logger for a text handler writing into buf at the
// given minimum level, returning a func to restore the previous default.
func setDefaultLogger(buf *bytes.Buffer, level slog.Level) func() {
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: level})))
	return func() { slog.SetDefault(prev) }
}
