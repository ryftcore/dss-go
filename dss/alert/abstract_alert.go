// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/AbstractAlert.java (DSS 6.5.RC1).
package alert

// AbstractAlert contains the general logic for alert handling: detect, then process.
type AbstractAlert[T any] struct {
	detector Detector[T]
	handler  Handler[T]
}

// NewAbstractAlert creates an AbstractAlert with the given detector and handler.
func NewAbstractAlert[T any](detector Detector[T], handler Handler[T]) *AbstractAlert[T] {
	return &AbstractAlert[T]{detector: detector, handler: handler}
}

// Alert runs the handler on object when the detector fires.
func (a *AbstractAlert[T]) Alert(object T) error {
	if a.alertDetector().Detect(object) {
		return a.alertHandler().Process(object)
	}
	return nil
}

// alertDetector returns the configured detector, panicking (mirroring Java's
// Objects.requireNonNull NullPointerException) if it was never set.
func (a *AbstractAlert[T]) alertDetector() Detector[T] {
	if a.detector == nil {
		panic("AlertDetector shall be defined!")
	}
	return a.detector
}

// alertHandler returns the configured handler, panicking if it was never set.
func (a *AbstractAlert[T]) alertHandler() Handler[T] {
	if a.handler == nil {
		panic("AlertHandler shall be defined!")
	}
	return a.handler
}
