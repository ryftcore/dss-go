package alert

import (
	"errors"
	"testing"
)

func TestExceptionOnStatusAlert_NoErrorWhenStatusEmpty(t *testing.T) {
	a := NewExceptionOnStatusAlert()
	if err := a.Alert(NewMessageStatus()); err != nil {
		t.Fatalf("Alert() error = %v, want nil for an empty status", err)
	}
}

func TestExceptionOnStatusAlert_ReturnsAlertErrorWhenStatusNotEmpty(t *testing.T) {
	a := NewExceptionOnStatusAlert()

	status := NewMessageStatus()
	status.SetMessage("token is not valid")

	err := a.Alert(status)
	if err == nil {
		t.Fatalf("Alert() error = nil, want AlertError")
	}

	var alertErr *Error
	if !errors.As(err, &alertErr) {
		t.Fatalf("Alert() error type = %T, want *AlertError", err)
	}
	if got, want := alertErr.Error(), "token is not valid"; got != want {
		t.Fatalf("AlertError.Error() = %q, want %q", got, want)
	}
}

func TestExceptionOnStatusAlert_ObjectStatusMessage(t *testing.T) {
	a := NewExceptionOnStatusAlert()

	status := NewObjectStatus()
	status.SetMessage("Signature is not valid")
	status.AddRelatedObjectIdentifierAndErrorMessage("sig-1", "expired")

	err := a.Alert(status)
	want := "Signature is not valid [sig-1: expired]"
	if err == nil || err.Error() != want {
		t.Fatalf("Alert() error = %v, want %q", err, want)
	}
}

func TestExceptionOnStatusAlert_ImplementsStatusAlert(t *testing.T) {
	var _ StatusAlert = NewExceptionOnStatusAlert()
}
