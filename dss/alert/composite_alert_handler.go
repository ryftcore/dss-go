// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/handler/CompositeAlertHandler.java (DSS 6.5.RC1).
package alert

// CompositeHandler runs multiple AlertHandlers in sequence.
//
// Judgment call: in Java, an exception thrown by one handler propagates immediately and
// stops the loop (no try/catch). Process mirrors that by returning on the first handler
// error rather than running all handlers unconditionally.
type CompositeHandler[T any] struct {
	handlers []Handler[T]
}

// NewCompositeAlertHandler creates a CompositeHandler running the given handlers in order.
func NewCompositeAlertHandler[T any](handlers []Handler[T]) *CompositeHandler[T] {
	return &CompositeHandler[T]{handlers: handlers}
}

// Process runs each handler in order, stopping and returning the first error encountered.
func (h *CompositeHandler[T]) Process(object T) error {
	for _, handler := range h.handlers {
		if err := handler.Process(object); err != nil {
			return err
		}
	}
	return nil
}
