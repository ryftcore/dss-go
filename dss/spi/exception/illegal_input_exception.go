// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/exception/IllegalInputException.java (DSS 6.5.RC1).
package exception

// IllegalInputException indicates that a provided user input or file is not valid for a
// particular operation. Ports IllegalInputException (extends RuntimeException directly, not
// DSSException); callers match it with errors.As per PORTING.md.
type IllegalInputException struct {
	// Message describes the exception.
	Message string

	// Cause is the wrapped underlying error, if any.
	Cause error
}

// Error implements the error interface.
func (e *IllegalInputException) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Unwrap allows errors.Is/errors.As to reach the wrapped Cause.
func (e *IllegalInputException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewIllegalInputException creates an IllegalInputException with a message. Port of
// IllegalInputException(String).
func NewIllegalInputException(message string) *IllegalInputException {
	return &IllegalInputException{Message: message}
}

// NewIllegalInputExceptionWithCause creates an IllegalInputException with a message and a
// wrapped cause. Port of IllegalInputException(String, Throwable).
func NewIllegalInputExceptionWithCause(message string, cause error) *IllegalInputException {
	return &IllegalInputException{Message: message, Cause: cause}
}
