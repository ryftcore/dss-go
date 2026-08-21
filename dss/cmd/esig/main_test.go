// Argument-parsing coverage: every usage error each subcommand's flag
// parsing rejects, and the smallest paths that need no fixture at all
// (version, tl's own usage). End-to-end coverage of the paths that need a
// signed document lives in e2e_test.go.
package main

import (
	"bytes"
	"strings"
	"testing"
)

// runOut runs args and returns its exit code and combined stdout/stderr, for
// tests that only care whether the right thing landed on the right stream at
// a glance; tests that need to tell the streams apart call run directly.
func runOut(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	code = run(args, &outBuf, &errBuf)
	return code, outBuf.String(), errBuf.String()
}

func TestNoArguments(t *testing.T) {
	code, _, stderr := runOut(t)
	if code != exitUsage {
		t.Errorf("exit code = %d, want %d (exitUsage)", code, exitUsage)
	}
	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("stderr = %q, want it to contain the usage banner", stderr)
	}
}

func TestUnknownCommand(t *testing.T) {
	code, _, stderr := runOut(t, "frobnicate")
	if code != exitUsage {
		t.Errorf("exit code = %d, want %d (exitUsage)", code, exitUsage)
	}
	if !strings.Contains(stderr, `"frobnicate"`) {
		t.Errorf("stderr = %q, want it to name the unknown command", stderr)
	}
}

func TestTopLevelHelp(t *testing.T) {
	code, stdout, _ := runOut(t, "-h")
	if code != exitOK {
		t.Errorf("exit code = %d, want %d (exitOK)", code, exitOK)
	}
	if !strings.Contains(stdout, "validate   validate") {
		t.Errorf("stdout = %q, want the command list", stdout)
	}
}

func TestVersion(t *testing.T) {
	code, stdout, _ := runOut(t, "version")
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (exitOK)", code, exitOK)
	}
	if !strings.HasPrefix(stdout, "esig ") {
		t.Errorf("stdout = %q, want it to start with %q", stdout, "esig ")
	}
	if !strings.Contains(stdout, "go:") {
		t.Errorf("stdout = %q, want a go: line", stdout)
	}
}

func TestVersionTakesNoArguments(t *testing.T) {
	code, _, _ := runOut(t, "version", "extra")
	if code != exitUsage {
		t.Errorf("exit code = %d, want %d (exitUsage)", code, exitUsage)
	}
}

// TestUsageErrors covers every subcommand's argument-parsing rejections: a
// missing required flag, a bad enum value, and the wrong number of
// positional arguments. None of these reach the facade, so none need a
// fixture.
func TestUsageErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"validate: no file", []string{"validate"}},
		{"validate: two files", []string{"validate", "a", "b"}},
		{"validate: bad -at", []string{"validate", "a", "-at", "not-a-time"}},
		{"validate: bad -format", []string{"validate", "a", "-format", "bogus"}},
		{"sign: no file", []string{"sign", "-format", "cades", "-level", "B", "-p12", "k", "-p12-pass", "env:X"}},
		{"sign: missing required flags", []string{"sign", "a"}},
		{"sign: bad -format", []string{"sign", "a", "-format", "bogus", "-level", "B", "-p12", "k", "-p12-pass", "env:X"}},
		{"sign: bad -level", []string{"sign", "a", "-format", "cades", "-level", "bogus", "-p12", "k", "-p12-pass", "env:X"}},
		{"sign: bad -digest", []string{"sign", "a", "-format", "cades", "-level", "B", "-p12", "k", "-p12-pass", "env:X", "-digest", "md5"}},
		{"sign: detached and packaging", []string{"sign", "a", "-format", "cades", "-level", "B", "-p12", "k", "-p12-pass", "env:X", "-detached", "-packaging", "enveloped"}},
		{"sign: level T needs -tsa", []string{"sign", "a", "-format", "cades", "-level", "T", "-p12", "k", "-p12-pass", "env:X"}},
		{"extend: no file", []string{"extend", "-format", "cades", "-level", "T", "-tsa", "http://x"}},
		{"extend: missing required flags", []string{"extend", "a"}},
		{"extend: level B rejected", []string{"extend", "a", "-format", "cades", "-level", "B", "-tsa", "http://x"}},
		{"extend: missing -tsa", []string{"extend", "a", "-format", "cades", "-level", "T"}},
		{"report: no file", []string{"report", "-render"}},
		{"report: missing -render", []string{"report", "a.xml"}},
		{"inspect: no file", []string{"inspect"}},
		{"tl: no subcommand", []string{"tl"}},
		{"tl: unknown subcommand", []string{"tl", "bogus"}},
		{"tl refresh: missing -cache", []string{"tl", "refresh"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, _, stderr := runOut(t, tc.args...)
			if code != exitUsage {
				t.Errorf("exit code = %d, want %d (exitUsage); stderr:\n%s", code, exitUsage, stderr)
			}
		})
	}
}

func TestTLHelp(t *testing.T) {
	code, stdout, _ := runOut(t, "tl", "-h")
	if code != exitOK {
		t.Errorf("exit code = %d, want %d (exitOK)", code, exitOK)
	}
	if !strings.Contains(stdout, "esig tl refresh") {
		t.Errorf("stdout = %q, want it to mention \"esig tl refresh\"", stdout)
	}
}
