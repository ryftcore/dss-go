// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/detector/StatusDetector.java (DSS 6.5.RC1).
package alert

// StatusDetector detects a custom event associated with token(s) processing.
type StatusDetector struct{}

// NewStatusDetector creates a StatusDetector.
func NewStatusDetector() *StatusDetector {
	return &StatusDetector{}
}

// Detect returns true when the status is not empty.
func (d *StatusDetector) Detect(object Status) bool {
	return !object.IsEmpty()
}
