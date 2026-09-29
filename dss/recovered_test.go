package dss

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// TestRecoveredConvertsPanicNil pins the assumption recovered's "p != nil" guard rests on:
// with this module's Go version (go.mod: go >= 1.21), panic(nil) is turned into a
// *runtime.PanicNilError, so recover() returns a non-nil value and the facade reports an error
// rather than success. Were go.mod's language version ever lowered below 1.21, panic(nil) would
// recover as nil and this test would fail.
func TestRecoveredConvertsPanicNil(t *testing.T) {
	err := recovered("sign", func() error {
		panic(nil)
	})
	if err == nil {
		t.Fatal("a panic(nil) inside the facade must surface as an error, not as success")
	}
	var panicNil *runtime.PanicNilError
	if !errors.As(err, &panicNil) {
		t.Fatalf("panic(nil) error = %v (%T), want it to wrap *runtime.PanicNilError", err, err)
	}
}

func TestRecoveredWrapsPanicValues(t *testing.T) {
	cause := model.NewDSSError("boom")
	err := recovered("validate", func() error { panic(cause) })
	var dssErr *model.DSSError
	if !errors.As(err, &dssErr) {
		t.Fatalf("panic(*model.DSSError) must stay reachable through errors.As, got %v", err)
	}

	err = recovered("extend", func() error { panic("plain message") })
	if err == nil || !strings.Contains(err.Error(), "dss: extend: plain message") {
		t.Fatalf("panic(string) error = %v, want it to carry the operation and the message", err)
	}

	if err := recovered("sign", func() error { return nil }); err != nil {
		t.Fatalf("no panic and no error must stay nil, got %v", err)
	}
}
