// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/alert/DSSErrorHandlerAlert.java (DSS 6.5.RC1).
package common

import (
	"fmt"

	"github.com/utain/esig/dss/alert"
)

// dssErrorHandlerDetector fires whenever the handler recorded any error, fatal error or
// warning. Ports the lambda errorHandler -> !errorHandler.isValid().
type dssErrorHandlerDetector struct{}

func (dssErrorHandlerDetector) Detect(h *DSSErrorHandler) bool {
	return !h.IsValid()
}

// DSSErrorHandlerAlert is the default DSSErrorHandler alert: it builds an
// XSDValidationException from the collected messages. Ports
// "class DSSErrorHandlerAlert extends AbstractAlert<DSSErrorHandler>" - a plain AbstractAlert
// subclass, not an AbstractStatusAlert one, since DSSErrorHandler does not implement the
// Status interface.
type DSSErrorHandlerAlert struct {
	*alert.AbstractAlert[*DSSErrorHandler]

	// enableWarnings indicates whether warning messages are reported within the error list.
	// Default: false (warnings are dropped, not logged - see the doc comment on
	// process/getWarnings below for why the Java "else log" branch has no port).
	enableWarnings bool

	// enablePosition indicates whether the position (line and column number) of the failed
	// condition is reported within the error list. Default: false.
	enablePosition bool
}

// NewDSSErrorHandlerAlert creates a DSSErrorHandlerAlert.
func NewDSSErrorHandlerAlert() *DSSErrorHandlerAlert {
	a := &DSSErrorHandlerAlert{}
	handler := dssErrorHandlerAlertHandlerFunc(a.process)
	a.AbstractAlert = alert.NewAbstractAlert[*DSSErrorHandler](dssErrorHandlerDetector{}, handler)
	return a
}

// dssErrorHandlerAlertHandlerFunc adapts a plain function to alert.AlertHandler[*DSSErrorHandler].
type dssErrorHandlerAlertHandlerFunc func(*DSSErrorHandler) error

func (f dssErrorHandlerAlertHandlerFunc) Process(h *DSSErrorHandler) error {
	return f(h)
}

// SetEnableWarnings sets whether validation warnings shall be returned within the error
// list. Ports setEnableWarnings(boolean).
func (a *DSSErrorHandlerAlert) SetEnableWarnings(enableWarnings bool) {
	a.enableWarnings = enableWarnings
}

// SetEnablePosition sets whether the position (line and column number) of the error shall
// be extracted into returned validation messages. Ports setEnablePosition(boolean).
func (a *DSSErrorHandlerAlert) SetEnablePosition(enablePosition bool) {
	a.enablePosition = enablePosition
}

// process builds the XSDValidationException from the handler's fatal errors, errors and
// (when enabled) warnings, in that order. Ports the anonymous AlertHandler<DSSErrorHandler>;
// Java throws the built exception, Go returns it as the error per PORTING.md.
//
// Java's "else if LOG.isDebugEnabled()" branch, which logs dropped warnings, has no port
// (slf4j is dropped per PORTING.md); dropped warnings are simply discarded.
func (a *DSSErrorHandlerAlert) process(h *DSSErrorHandler) error {
	var exceptions []*XMLParseError
	exceptions = append(exceptions, h.FatalErrors()...)
	exceptions = append(exceptions, h.Errors()...)
	if a.enableWarnings {
		exceptions = append(exceptions, h.Warnings()...)
	}

	messages := make([]string, 0, len(exceptions))
	for _, e := range exceptions {
		messages = append(messages, a.validationMessage(e))
	}

	return NewXSDValidationException(messages)
}

// validationMessage builds a validation message from e. Ports getValidationMessage(SAXParseException).
func (a *DSSErrorHandlerAlert) validationMessage(e *XMLParseError) string {
	if a.enablePosition {
		return fmt.Sprintf("%s (Line: %d, Column: %d)", e.Message, e.LineNumber, e.ColumnNumber)
	}
	return e.Message
}
