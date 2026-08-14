// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/BaseTimeDependent.java (DSS 6.5.RC1).
package timedependent

import (
	"testing"
	"time"
)

func TestBaseTimeDependent_RoundTrip(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	btd := NewBaseTimeDependentWithDates(start, end)
	if !btd.StartDate().Equal(start) {
		t.Fatalf("StartDate() = %v, want %v", btd.StartDate(), start)
	}
	if !btd.EndDate().Equal(end) {
		t.Fatalf("EndDate() = %v, want %v", btd.EndDate(), end)
	}

	empty := NewBaseTimeDependent()
	if !empty.StartDate().IsZero() || !empty.EndDate().IsZero() {
		t.Fatalf("empty constructor should leave both dates as the zero time.Time")
	}
}

func TestBaseTimeDependent_Equals(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	a := NewBaseTimeDependentWithDates(start, end)
	b := NewBaseTimeDependentWithDates(start, end)
	c := NewBaseTimeDependentWithDates(start, time.Time{})

	if !a.Equals(b) {
		t.Fatalf("Equals() = false for identical dates")
	}
	if a.Equals(c) {
		t.Fatalf("Equals() = true for differing end dates")
	}
	if a.Equals(nil) {
		t.Fatalf("Equals(nil) = true")
	}
}

func TestBaseTimeDependent_String(t *testing.T) {
	empty := NewBaseTimeDependent()
	if got, want := empty.String(), "[startDate=null, endDate=null]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
