// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/DSSReportException.java
// (DSS 6.5.RC1).

package reports

// DSSReportException is the error returned in case of a JAXB Report
// marshalling or unmarshalling error. Port of DSSReportException, whose four
// constructors (empty, message, cause, message+cause) collapse to the
// idiomatic Go wrapped-error constructors below.
type DSSReportException struct {
	msg   string
	cause error
}

// NewDSSReportException builds an empty DSSReportException. Port of the
// no-arg DSSReportException() constructor.
func NewDSSReportException() *DSSReportException {
	return &DSSReportException{}
}

// NewDSSReportExceptionMessage builds a DSSReportException with a message.
// Port of DSSReportException(String).
func NewDSSReportExceptionMessage(message string) *DSSReportException {
	return &DSSReportException{msg: message}
}

// NewDSSReportExceptionCause builds a DSSReportException wrapping a cause.
// Port of DSSReportException(Throwable).
func NewDSSReportExceptionCause(cause error) *DSSReportException {
	return &DSSReportException{cause: cause}
}

// NewDSSReportExceptionMessageCause builds a DSSReportException with a
// message and a wrapped cause. Port of DSSReportException(String,
// Throwable).
func NewDSSReportExceptionMessageCause(message string, cause error) *DSSReportException {
	return &DSSReportException{msg: message, cause: cause}
}

// Error implements the error interface, mirroring Throwable.getMessage()
// composed with the cause when present.
func (e *DSSReportException) Error() string {
	if e.msg != "" && e.cause != nil {
		return e.msg + ": " + e.cause.Error()
	}
	if e.msg != "" {
		return e.msg
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return "DSSReportException"
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As. Port of
// Throwable.getCause().
func (e *DSSReportException) Unwrap() error {
	return e.cause
}
