// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/alerts/ProtectedDocumentExceptionOnStatusAlert.java
// (DSS 6.5.RC1).
//
// Java's `extends AbstractStatusAlert` with a lambda AlertHandler becomes embedding
// *alert.AbstractStatusAlert built from an alert.AlertHandler[alert.Status] closure, per the
// alert/exception_on_status_alert.go precedent.
package alerts

import (
	"github.com/utain/esig/dss/alert"
	"github.com/utain/esig/dss/pades/exception"
)

// ProtectedDocumentExceptionOnStatusAlert is used to throw a
// pades/exception.ProtectedDocumentException when the corresponding check fails.
type ProtectedDocumentExceptionOnStatusAlert struct {
	*alert.AbstractStatusAlert
}

// protectedDocumentExceptionOnStatusAlertHandler is the AlertHandler counterpart of the
// upstream lambda `object -> { throw new ProtectedDocumentException(object.getErrorString()); }`.
type protectedDocumentExceptionOnStatusAlertHandler struct{}

// Process ports the lambda's body.
func (protectedDocumentExceptionOnStatusAlertHandler) Process(status alert.Status) error {
	return exception.NewProtectedDocumentException(status.ErrorString())
}

// NewProtectedDocumentExceptionOnStatusAlert is the default constructor.
func NewProtectedDocumentExceptionOnStatusAlert() *ProtectedDocumentExceptionOnStatusAlert {
	return &ProtectedDocumentExceptionOnStatusAlert{
		AbstractStatusAlert: alert.NewAbstractStatusAlert(protectedDocumentExceptionOnStatusAlertHandler{}),
	}
}
