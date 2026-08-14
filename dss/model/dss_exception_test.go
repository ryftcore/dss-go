package model

import (
	"errors"
	"testing"
)

func TestDSSErrorMessageOnly(t *testing.T) {
	err := NewDSSError("boom")
	if err.Error() != "boom" {
		t.Fatalf("got %q, want %q", err.Error(), "boom")
	}
}

func TestDSSErrorCauseOnly(t *testing.T) {
	cause := errors.New("root cause")
	err := NewDSSErrorWithCause(cause)
	if err.Error() != "root cause" {
		t.Fatalf("got %q, want %q", err.Error(), "root cause")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected errors.Is to unwrap to cause")
	}
}

func TestDSSErrorMessageAndCause(t *testing.T) {
	cause := errors.New("root cause")
	err := NewDSSErrorMessageCause("wrapper", cause)
	want := "wrapper: root cause"
	if err.Error() != want {
		t.Fatalf("got %q, want %q", err.Error(), want)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected errors.Is to unwrap to cause")
	}
}
