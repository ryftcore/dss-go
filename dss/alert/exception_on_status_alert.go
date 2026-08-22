// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/ExceptionOnStatusAlert.java (DSS 6.5.RC1).
package alert

// ExceptionOnStatusAlert returns an Error on a Status event.
type ExceptionOnStatusAlert struct {
	*AbstractStatusAlert
}

// NewExceptionOnStatusAlert creates an ExceptionOnStatusAlert.
func NewExceptionOnStatusAlert() *ExceptionOnStatusAlert {
	return &ExceptionOnStatusAlert{
		AbstractStatusAlert: NewAbstractStatusAlert(NewThrowAlertExceptionHandler[Status]()),
	}
}
