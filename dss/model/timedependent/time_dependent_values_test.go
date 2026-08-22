// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/TimeDependentValues.java (DSS 6.5.RC1).
package timedependent

import (
	"testing"
	"time"
)

func TestTimeDependentValues_LatestAndCurrent(t *testing.T) {
	d2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	d2021 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	d2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	// Latest first, matching the Java contract.
	latest := NewBaseTimeDependentWithDates(d2021, time.Time{})
	oldest := NewBaseTimeDependentWithDates(d2020, d2021)

	values := NewValuesFrom[*BaseTimeDependent]([]*BaseTimeDependent{latest, oldest})

	if got := values.Latest(); got != latest {
		t.Fatalf("Latest() = %v, want %v", got, latest)
	}

	if got := values.Current(d2020); got != oldest {
		t.Fatalf("Current(2020) = %v, want the oldest entry", got)
	}
	if got := values.Current(d2022); got != latest {
		t.Fatalf("Current(2022) = %v, want the open-ended latest entry", got)
	}

	empty := NewValues[*BaseTimeDependent]()
	if got := empty.Latest(); got != nil {
		t.Fatalf("Latest() on an empty list = %v, want nil", got)
	}
	if got := empty.Current(d2020); got != nil {
		t.Fatalf("Current() on an empty list = %v, want nil", got)
	}
}

func TestTimeDependentValues_After(t *testing.T) {
	d2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	d2021 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	d2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	openEnded := NewBaseTimeDependentWithDates(d2021, time.Time{})
	closedBefore := NewBaseTimeDependentWithDates(d2020, d2020.AddDate(0, 6, 0))
	closedAfter := NewBaseTimeDependentWithDates(d2020, d2022)

	values := NewValuesFrom[*BaseTimeDependent]([]*BaseTimeDependent{openEnded, closedBefore, closedAfter})

	after := values.After(d2021)
	if len(after) != 2 {
		t.Fatalf("After(2021) returned %d entries, want 2: %v", len(after), after)
	}
	for _, x := range after {
		if x == closedBefore {
			t.Fatalf("After(2021) should not include an entry that ended before notBefore")
		}
	}
}

func TestTimeDependentValues_Iterator(t *testing.T) {
	a := NewBaseTimeDependentWithDates(time.Now(), time.Time{})
	values := NewValuesFrom[*BaseTimeDependent]([]*BaseTimeDependent{a})

	var seen []*BaseTimeDependent
	for x := range values.Iterator() {
		seen = append(seen, x)
	}
	if len(seen) != 1 || seen[0] != a {
		t.Fatalf("Iterator() yielded %v, want [%v]", seen, a)
	}
}

func TestTimeDependentValues_String(t *testing.T) {
	empty := NewValues[*BaseTimeDependent]()
	if got, want := empty.String(), "[]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	one := NewValuesFrom[*BaseTimeDependent]([]*BaseTimeDependent{NewBaseTimeDependent()})
	if got, want := one.String(), "[[startDate=null, endDate=null]]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
