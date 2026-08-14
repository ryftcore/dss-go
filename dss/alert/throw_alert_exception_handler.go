// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/handler/ThrowAlertExceptionHandler.java (DSS 6.5.RC1).
package alert

import "fmt"

// ThrowAlertExceptionHandler is an AlertHandler which reports an AlertError built from the
// object's string representation. Java throws the AlertException; Go returns it.
type ThrowAlertExceptionHandler[T any] struct{}

// NewThrowAlertExceptionHandler creates a ThrowAlertExceptionHandler.
func NewThrowAlertExceptionHandler[T any]() *ThrowAlertExceptionHandler[T] {
	return &ThrowAlertExceptionHandler[T]{}
}

// Process returns an AlertError whose message is fmt.Sprint(object) (mirroring Java's
// object.toString()).
func (h *ThrowAlertExceptionHandler[T]) Process(object T) error {
	return NewAlertErrorWithMessage(fmt.Sprint(object))
}
