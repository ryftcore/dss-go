// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/exception/DSSExternalResourceException.java (DSS 6.5.RC1).
package exception

import (
	"github.com/utain/esig/dss/model"
)

// DSSExternalResourceException is thrown in case of an external error arisen during a data
// loader request. Ports DSSExternalResourceException (extends DSSException); per PORTING.md's
// "DSSException -> model.DSSError; spi.exception classes -> error types in pkg exception", it
// embeds *model.DSSError for the Message/Cause/Error()/Unwrap() behaviour rather than
// duplicating it.
type DSSExternalResourceException struct {
	*model.DSSError
}

// newDSSExternalResourceExceptionEmpty creates an empty DSSExternalResourceException. Port of
// the package-private no-arg constructor DSSExternalResourceException(), used only by
// DSSDataLoaderMultipleException's Go counterpart to embed this type without a message.
func newDSSExternalResourceExceptionEmpty() *DSSExternalResourceException {
	return &DSSExternalResourceException{DSSError: &model.DSSError{}}
}

// NewDSSExternalResourceException creates a DSSExternalResourceException with a message. Port
// of DSSExternalResourceException(String).
func NewDSSExternalResourceException(message string) *DSSExternalResourceException {
	return &DSSExternalResourceException{DSSError: model.NewDSSError(message)}
}

// NewDSSExternalResourceExceptionWithCause creates a re-throwable DSSExternalResourceException
// wrapping cause. Port of DSSExternalResourceException(Throwable).
func NewDSSExternalResourceExceptionWithCause(cause error) *DSSExternalResourceException {
	return &DSSExternalResourceException{DSSError: model.NewDSSErrorWithCause(cause)}
}

// NewDSSExternalResourceExceptionMessageCause creates a DSSExternalResourceException with a
// custom message wrapping cause. Port of DSSExternalResourceException(String, Throwable).
func NewDSSExternalResourceExceptionMessageCause(message string, cause error) *DSSExternalResourceException {
	return &DSSExternalResourceException{DSSError: model.NewDSSErrorMessageCause(message, cause)}
}

// causeMessage returns the wrapped cause's message, or the empty string when there is none.
// Package-private port of DSSExternalResourceException#getCauseMessage(), used by
// DSSDataLoaderMultipleException to flatten a multi-URL failure report.
func (e *DSSExternalResourceException) causeMessage() string {
	if e == nil || e.DSSError == nil || e.Cause == nil {
		return ""
	}
	return e.Cause.Error()
}
