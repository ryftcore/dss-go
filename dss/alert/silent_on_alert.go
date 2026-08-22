// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/SilentOnAlert.java (DSS 6.5.RC1).
package alert

// SilentOnAlert processes an alert silently (does nothing) once detected.
type SilentOnAlert[T any] struct {
	*AbstractAlert[T]
}

// NewSilentOnAlert creates a SilentOnAlert using detector to decide when to fire.
func NewSilentOnAlert[T any](detector Detector[T]) *SilentOnAlert[T] {
	return &SilentOnAlert[T]{
		AbstractAlert: NewAbstractAlert[T](detector, NewSilentHandler[T]()),
	}
}
