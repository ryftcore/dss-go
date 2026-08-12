// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/SilentOnStatusAlert.java (DSS 6.5.RC1).
package alert

// SilentOnStatusAlert processes a Status alert silently.
type SilentOnStatusAlert struct {
	*AbstractStatusAlert
}

// NewSilentOnStatusAlert creates a SilentOnStatusAlert.
func NewSilentOnStatusAlert() *SilentOnStatusAlert {
	return &SilentOnStatusAlert{
		AbstractStatusAlert: NewAbstractStatusAlert(NewSilentHandler[Status]()),
	}
}
