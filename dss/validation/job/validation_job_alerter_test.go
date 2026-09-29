package job

import (
	"errors"
	"testing"
)

// testAlert is an alert.Alert[string] whose Alert either returns err, panics, or records info.
type testAlert struct {
	err       error
	panicWith any
	seen      *[]string
}

func (a testAlert) Alert(info string) error {
	if a.seen != nil {
		*a.seen = append(*a.seen, info)
	}
	if a.panicWith != nil {
		panic(a.panicWith)
	}
	return a.err
}

// TestExecuteAlert_SwallowsErrorsAndPanics ports ValidationJobAlerter.execute's
// `catch (Exception e)`: a handler failure, whether it is the error Alert returns (Java's
// checked / handler-raised exceptions) or a panic (an unchecked exception such as a
// NullPointerException in a detector), must not escape and abort the alerting pass.
func TestExecuteAlert_SwallowsErrorsAndPanics(t *testing.T) {
	var seen []string
	alerts := []testAlert{
		{err: errors.New("handler failed"), seen: &seen},
		{panicWith: "boom", seen: &seen},
		{panicWith: errors.New("nil dereference"), seen: &seen},
		{seen: &seen},
	}
	for _, a := range alerts {
		executeDocumentAlert[string](a, "doc")
		executeDocumentListAlert[string](a, "list")
	}
	if len(seen) != 2*len(alerts) {
		t.Fatalf("every alert must have run (and none aborted the pass): saw %v", seen)
	}
}
