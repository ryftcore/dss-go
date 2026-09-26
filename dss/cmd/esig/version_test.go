package main

import (
	"strings"
	"testing"
)

// TestVersionStamped pins that a release build's -X main.version stamp
// takes precedence over the module version from the build information.
func TestVersionStamped(t *testing.T) {
	old := version
	version = "v9.9.9-test"
	t.Cleanup(func() { version = old })

	code, stdout, _ := runOut(t, "version")
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (exitOK)", code, exitOK)
	}
	if !strings.HasPrefix(stdout, "esig v9.9.9-test\n") {
		t.Errorf("stdout = %q, want it to start with %q", stdout, "esig v9.9.9-test\n")
	}
}
