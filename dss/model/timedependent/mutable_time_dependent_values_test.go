// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/MutableTimeDependentValues.java (DSS 6.5.RC1).
package timedependent

import (
	"testing"
	"time"
)

func TestMutableTimeDependentValues_AddOldestAndClear(t *testing.T) {
	d2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	d2021 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	values := NewMutableTimeDependentValues[*BaseTimeDependent]()
	values.AddOldest(NewBaseTimeDependentWithDates(d2021, time.Time{}))
	values.AddOldest(NewBaseTimeDependentWithDates(d2020, d2021))

	if got := len(values.List()); got != 2 {
		t.Fatalf("List() has %d entries, want 2", got)
	}
	if got := values.Latest().StartDate(); !got.Equal(d2021) {
		t.Fatalf("Latest().StartDate() = %v, want %v (insertion order is preserved)", got, d2021)
	}

	values.Clear()
	if got := len(values.List()); got != 0 {
		t.Fatalf("List() after Clear() has %d entries, want 0", got)
	}
}

func TestMutableTimeDependentValues_AddOldestPanicsOnOverlap(t *testing.T) {
	d2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	d2021 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	d2019 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	values := NewMutableTimeDependentValuesFrom([]*BaseTimeDependent{NewBaseTimeDependentWithDates(d2020, d2021)})

	defer func() {
		if r := recover(); r != "Cannot add overlapping item" {
			t.Fatalf("recover() = %v, want the Java IllegalArgumentException message", r)
		}
	}()
	// Ends after the existing entry's start date: overlaps.
	values.AddOldest(NewBaseTimeDependentWithDates(d2019, d2021))
}

func TestMutableTimeDependentValues_AddOldestPanicsOnNil(t *testing.T) {
	values := NewMutableTimeDependentValues[*BaseTimeDependent]()

	defer func() {
		if r := recover(); r != "Cannot add null" {
			t.Fatalf("recover() = %v, want the Java requireNonNull message", r)
		}
	}()
	values.AddOldest(nil)
}
