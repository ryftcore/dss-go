// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/status/Status.java (DSS 6.5.RC1).
package alert

import "fmt"

// Status contains information about an occurred event.
type Status interface {
	fmt.Stringer

	// Message returns the error message describing the occurred event.
	Message() string

	// RelatedObjectIds returns the identifiers of the objects associated with the event.
	RelatedObjectIds() []string

	// IsEmpty returns true when the Status is not filled (the underlying event did not occur).
	IsEmpty() bool

	// ErrorString returns the complete error message computed from the main message and,
	// where applicable, the sub-messages of the related objects.
	ErrorString() string
}
