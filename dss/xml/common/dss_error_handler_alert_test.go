package common

import (
	"errors"
	"testing"
)

func TestDSSErrorHandlerAlert_ValidHandlerDoesNotAlert(t *testing.T) {
	a := NewDSSErrorHandlerAlert()
	h := NewDSSErrorHandler()
	if err := a.Alert(h); err != nil {
		t.Fatalf("Alert() error = %v, want nil for a valid (empty) handler", err)
	}
}

func TestDSSErrorHandlerAlert_ErrorsAndFatalErrorsAreReported(t *testing.T) {
	a := NewDSSErrorHandlerAlert()
	h := NewDSSErrorHandler()
	h.RecordFatalError(&XMLParseError{Message: "fatal one"})
	h.RecordError(&XMLParseError{Message: "error one"})
	h.RecordWarning(&XMLParseError{Message: "warning one"})

	err := a.Alert(h)
	if err == nil {
		t.Fatalf("Alert() error = nil, want an XSDValidationException")
	}
	var xsdErr *XSDValidationException
	if !errors.As(err, &xsdErr) {
		t.Fatalf("Alert() error = %T, want *XSDValidationException", err)
	}
	// Order: fatal errors, then errors, then (if enabled) warnings.
	want := []string{"fatal one", "error one"}
	if len(xsdErr.Messages) != len(want) {
		t.Fatalf("Messages = %v, want %v", xsdErr.Messages, want)
	}
	for i := range want {
		if xsdErr.Messages[i] != want[i] {
			t.Errorf("Messages[%d] = %q, want %q", i, xsdErr.Messages[i], want[i])
		}
	}
}

func TestDSSErrorHandlerAlert_EnableWarnings(t *testing.T) {
	a := NewDSSErrorHandlerAlert()
	a.SetEnableWarnings(true)
	h := NewDSSErrorHandler()
	h.RecordWarning(&XMLParseError{Message: "warning one"})

	err := a.Alert(h)
	var xsdErr *XSDValidationException
	if !errors.As(err, &xsdErr) {
		t.Fatalf("Alert() error = %T, want *XSDValidationException", err)
	}
	if len(xsdErr.Messages) != 1 || xsdErr.Messages[0] != "warning one" {
		t.Errorf("Messages = %v, want [\"warning one\"]", xsdErr.Messages)
	}
}

func TestDSSErrorHandlerAlert_EnablePosition(t *testing.T) {
	a := NewDSSErrorHandlerAlert()
	a.SetEnablePosition(true)
	h := NewDSSErrorHandler()
	h.RecordError(&XMLParseError{Message: "bad element", LineNumber: 3, ColumnNumber: 7})

	err := a.Alert(h)
	var xsdErr *XSDValidationException
	if !errors.As(err, &xsdErr) {
		t.Fatalf("Alert() error = %T, want *XSDValidationException", err)
	}
	want := "bad element (Line: 3, Column: 7)"
	if len(xsdErr.Messages) != 1 || xsdErr.Messages[0] != want {
		t.Errorf("Messages = %v, want [%q]", xsdErr.Messages, want)
	}
}

func TestXSDValidationException_ErrorJoinsMessages(t *testing.T) {
	e := NewXSDValidationException([]string{"a", "b"})
	if got, want := e.Error(), "a; b"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	empty := NewXSDValidationException(nil)
	if got := empty.Error(); got != "" {
		t.Errorf("Error() = %q, want \"\" for an empty message list", got)
	}
}

func TestSecurityConfigurationException_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	e := NewSecurityConfigurationException(cause)
	if !errors.Is(e, cause) {
		t.Errorf("errors.Is(e, cause) = false, want true")
	}
	if e.Error() != cause.Error() {
		t.Errorf("Error() = %q, want %q", e.Error(), cause.Error())
	}
}
