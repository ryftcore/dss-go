// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/StatusAlert.java (DSS 6.5.RC1).
package alert

// StatusAlert is an Alert typed with a Status object.
type StatusAlert interface {
	Alert[Status]
}
