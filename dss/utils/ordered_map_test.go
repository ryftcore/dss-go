package utils

import (
	"reflect"
	"testing"
)

func TestOrderedMapInsertionOrder(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)
	// re-Set an existing key: value updates, position does not move.
	m.Set("a", 10)

	if got, want := m.Keys(), []string{"c", "a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
	if got, want := m.Values(), []int{3, 10, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Values() = %v, want %v", got, want)
	}
	if v, found := m.Get("a"); !found || v != 10 {
		t.Fatalf("Get(a) = %v, %v, want 10, true", v, found)
	}
	if m.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", m.Len())
	}
}

func TestOrderedMapDelete(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Delete("b")

	if got, want := m.Keys(), []string{"a", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
	if _, found := m.Get("b"); found {
		t.Fatalf("Get(b) found after Delete")
	}
	// Deleting a missing key is a no-op.
	m.Delete("missing")
	if got, want := m.Keys(), []string{"a", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() after no-op delete = %v, want %v", got, want)
	}
}

func TestOrderedMapReset(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("a", 1)
	m.Reset()
	if m.Len() != 0 {
		t.Fatalf("Len() after Reset = %d, want 0", m.Len())
	}
	m.Set("z", 26)
	if got, want := m.Keys(), []string{"z"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() after Reset+Set = %v, want %v", got, want)
	}
}

func TestOrderedMapRange(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("x", 1)
	m.Set("y", 2)
	m.Set("z", 3)

	var seen []string
	m.Range(func(k string, v int) bool {
		seen = append(seen, k)
		return v != 2
	})
	if got, want := seen, []string{"x", "y"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Range early-stop = %v, want %v", got, want)
	}
}

func TestOrderedMapNilReceiverReadsAsEmpty(t *testing.T) {
	var m *OrderedMap[string, int]
	if got := m.Len(); got != 0 {
		t.Fatalf("nil.Len() = %d, want 0", got)
	}
	if got := m.Keys(); got != nil {
		t.Fatalf("nil.Keys() = %v, want nil", got)
	}
	if got := m.Values(); got != nil {
		t.Fatalf("nil.Values() = %v, want nil", got)
	}
	if _, found := m.Get("a"); found {
		t.Fatalf("nil.Get() found = true, want false")
	}
	// Delete and Range on a nil receiver must not panic.
	m.Delete("a")
	m.Range(func(string, int) bool { t.Fatal("Range on nil map should not call f"); return true })
}

// TestOrderedMapDeterministic25Runs is the shared determinism check pattern used at every
// map-iteration-ordering site: populate, snapshot output, compare across many iterations.
func TestOrderedMapDeterministic25Runs(t *testing.T) {
	build := func() []string {
		m := NewOrderedMap[string, int]()
		for i, k := range []string{"delta", "alpha", "charlie", "echo", "bravo"} {
			m.Set(k, i)
		}
		return m.Keys()
	}
	want := build()
	for i := 0; i < 25; i++ {
		if got := build(); !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: Keys() = %v, want %v", i, got, want)
		}
	}
}
