// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/detector/AlertDetector.java (DSS 6.5.RC1).
package alert

// AlertDetector detects, on an object, whether an alert must be executed.
type AlertDetector[T any] interface {
	// Detect returns true if the alert must be executed, false otherwise.
	Detect(object T) bool
}
