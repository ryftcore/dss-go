// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/handler/AlertHandler.java (DSS 6.5.RC1).
package alert

// Handler executes a process on an object after an alert-worthy event has been detected.
//
// Judgment call: Java's process(T) is void and reports failure by throwing an unchecked
// AlertException. Per PORTING.md ("dss-alert handlers keep their semantics: alerts receive
// a status and decide to log/throw; in Go they receive the status and may return an error"),
// Process here returns error instead of panicking, so ThrowAlertExceptionHandler's
// Error propagates through normal Go error handling.
type Handler[T any] interface {
	// Process alerts on the object after some change or problem has been detected.
	Process(object T) error
}
