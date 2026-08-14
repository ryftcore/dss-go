// Ported from dss-model/.../DSSException.java (DSS 6.5.RC1).
package model

// DSSError is the error raised for data-dependent failures during DSS
// framework processing. Ports DSSException (a RuntimeException in Java;
// callers in other Go packages match it with errors.As).
type DSSError struct {
	// Message is the error message.
	Message string

	// Cause is the wrapped underlying error, if any.
	Cause error
}

// Error implements the error interface.
func (e *DSSError) Error() string {
	if e.Message != "" && e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "DSSError"
}

// Unwrap allows errors.Is/errors.As to reach the wrapped Cause.
func (e *DSSError) Unwrap() error { return e.Cause }

// NewDSSError creates a DSSError with a message. Ports
// DSSException(String message).
func NewDSSError(message string) *DSSError {
	return &DSSError{Message: message}
}

// NewDSSErrorWithCause creates a re-throwable DSSError wrapping cause.
// Ports DSSException(Throwable cause).
func NewDSSErrorWithCause(cause error) *DSSError {
	return &DSSError{Cause: cause}
}

// NewDSSErrorMessageCause creates a DSSError with a custom message wrapping
// cause. Ports DSSException(String message, Throwable cause).
func NewDSSErrorMessageCause(message string, cause error) *DSSError {
	return &DSSError{Message: message, Cause: cause}
}
