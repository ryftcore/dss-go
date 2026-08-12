// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/Alert.java (DSS 6.5.RC1).
package alert

// Alert detects an event on an object and, if warranted, executes the corresponding handler.
//
// Judgment call: Alert(T) returns error (Java's alert(T) is void). See AlertHandler for the
// rationale — the returned error carries what Java would throw as an AlertException.
type Alert[T any] interface {
	// Alert detects and, if needed, executes the alert on the provided object.
	Alert(object T) error
}
