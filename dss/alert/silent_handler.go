// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/handler/SilentHandler.java (DSS 6.5.RC1).
package alert

// SilentHandler is an Handler which does nothing.
type SilentHandler[T any] struct{}

// NewSilentHandler creates a SilentHandler.
func NewSilentHandler[T any]() *SilentHandler[T] {
	return &SilentHandler[T]{}
}

// Process does nothing and always returns nil.
func (h *SilentHandler[T]) Process(object T) error {
	return nil
}
