// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/AbstractStatusAlert.java (DSS 6.5.RC1).
package alert

// AbstractStatusAlert is the base for alerts that process a Status via a StatusDetector.
type AbstractStatusAlert struct {
	*AbstractAlert[Status]
}

// NewAbstractStatusAlert creates an AbstractStatusAlert running handler when a non-empty
// Status is detected.
func NewAbstractStatusAlert(handler Handler[Status]) *AbstractStatusAlert {
	return &AbstractStatusAlert{
		AbstractAlert: NewAbstractAlert[Status](NewStatusDetector(), handler),
	}
}
