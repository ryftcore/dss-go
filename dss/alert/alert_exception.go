// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/exception/AlertException.java (DSS 6.5.RC1).
package alert

// AlertError is the error thrown/returned by an exception-raising alert (ports AlertException,
// a RuntimeException in Java). It implements the error interface and supports errors.Unwrap
// to reach a wrapped cause.
//
// Judgment call: Java's cause-only constructor sets the exception's message to
// cause.toString() (i.e. "java.lang.SomeException: cause message"). Go errors have no
// class-qualified toString(), so NewAlertErrorWithCause uses cause.Error() as the message
// instead — the idiomatic Go equivalent, not a byte-exact reproduction of the Java string.
type AlertError struct {
	message string
	cause   error
}

// NewAlertError creates an empty AlertError.
func NewAlertError() *AlertError {
	return &AlertError{}
}

// NewAlertErrorWithMessage creates an AlertError with the given message.
func NewAlertErrorWithMessage(message string) *AlertError {
	return &AlertError{message: message}
}

// NewAlertErrorWithCause creates a re-throwable AlertError wrapping cause.
func NewAlertErrorWithCause(cause error) *AlertError {
	message := ""
	if cause != nil {
		message = cause.Error()
	}
	return &AlertError{message: message, cause: cause}
}

// NewAlertErrorWithMessageAndCause creates a re-throwable AlertError with a custom message.
func NewAlertErrorWithMessageAndCause(message string, cause error) *AlertError {
	return &AlertError{message: message, cause: cause}
}

// Error returns the error message.
func (e *AlertError) Error() string {
	return e.message
}

// Unwrap returns the wrapped cause, if any, enabling errors.Is/errors.As.
func (e *AlertError) Unwrap() error {
	return e.cause
}
