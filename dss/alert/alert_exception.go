// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/exception/AlertException.java (DSS 6.5.RC1).
package alert

// AlertError is the error thrown/returned by an exception-raising alert (ports AlertException,
// a RuntimeException in Java). It implements the error interface and supports errors.Unwrap
// to reach a wrapped cause.
//
// Judgment call: Java's cause-only constructor sets the exception's message to
// cause.toString() (i.e. "java.lang.SomeException: cause message"). Go errors have no
// class-qualified toString(), so NewErrorWithCause uses cause.Error() as the message
// instead — the idiomatic Go equivalent, not a byte-exact reproduction of the Java string.
type Error struct {
	message string
	cause   error
}

// NewError creates an empty Error.
func NewError() *Error {
	return &Error{}
}

// NewErrorWithMessage creates an Error with the given message.
func NewErrorWithMessage(message string) *Error {
	return &Error{message: message}
}

// NewErrorWithCause creates a re-throwable Error wrapping cause.
func NewErrorWithCause(cause error) *Error {
	message := ""
	if cause != nil {
		message = cause.Error()
	}
	return &Error{message: message, cause: cause}
}

// NewErrorWithMessageAndCause creates a re-throwable Error with a custom message.
func NewErrorWithMessageAndCause(message string, cause error) *Error {
	return &Error{message: message, cause: cause}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.message
}

// Unwrap returns the wrapped cause, if any, enabling errors.Is/errors.As.
func (e *Error) Unwrap() error {
	return e.cause
}
