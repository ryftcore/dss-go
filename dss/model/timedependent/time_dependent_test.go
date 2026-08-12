// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/TimeDependent.java (DSS 6.5.RC1).
package timedependent

import "testing"

func TestTimeDependent_ImplementedByBaseTimeDependent(t *testing.T) {
	var td TimeDependent = NewBaseTimeDependent()
	if !td.StartDate().IsZero() || !td.EndDate().IsZero() {
		t.Fatalf("fresh BaseTimeDependent should report zero-value dates through the interface")
	}
}
