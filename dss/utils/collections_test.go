// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against known
// org.apache.commons.collections4.CollectionUtils/MapUtils semantics.

package utils

import (
	"reflect"
	"testing"
)

func TestIsCollectionEmpty(t *testing.T) {
	if !IsCollectionEmpty[int](nil) {
		t.Error("nil collection should be empty")
	}
	if !IsCollectionEmpty([]int{}) {
		t.Error("empty collection should be empty")
	}
	if IsCollectionEmpty([]int{1}) {
		t.Error("non-empty collection should not be empty")
	}
	if IsCollectionNotEmpty[int](nil) {
		t.Error("IsCollectionNotEmpty(nil) should be false")
	}
	if CollectionSize[int](nil) != 0 {
		t.Error("CollectionSize(nil) should be 0")
	}
	if CollectionSize([]int{1, 2, 3}) != 3 {
		t.Error("CollectionSize should be 3")
	}
}

func TestIsMapEmpty(t *testing.T) {
	if !IsMapEmpty[string, int](nil) {
		t.Error("nil map should be empty")
	}
	if !IsMapEmpty(map[string]int{}) {
		t.Error("empty map should be empty")
	}
	if IsMapEmpty(map[string]int{"a": 1}) {
		t.Error("non-empty map should not be empty")
	}
	if IsMapNotEmpty[string, int](nil) {
		t.Error("IsMapNotEmpty(nil) should be false")
	}
	if MapSize[string, int](nil) != 0 {
		t.Error("MapSize(nil) should be 0")
	}
	if MapSize(map[string]int{"a": 1, "b": 2}) != 2 {
		t.Error("MapSize should be 2")
	}
}

func TestReverseList(t *testing.T) {
	in := []int{1, 2, 3}
	got := ReverseList(in)
	want := []int{3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReverseList(%v) = %v, want %v", in, got, want)
	}
	// Original must be untouched.
	if !reflect.DeepEqual(in, []int{1, 2, 3}) {
		t.Errorf("ReverseList mutated its input: %v", in)
	}
	if got := ReverseList[int](nil); len(got) != 0 {
		t.Errorf("ReverseList(nil) = %v, want empty", got)
	}
}

func TestContainsAny(t *testing.T) {
	if !ContainsAny([]string{"A", "B", "C"}, []string{"B", "C", "D"}) {
		t.Error("expected intersection")
	}
	if ContainsAny([]string{"A", "B"}, []string{"C", "D"}) {
		t.Error("expected no intersection")
	}
	if ContainsAny[string](nil, []string{"A"}) {
		t.Error("nil superCollection should yield false")
	}
	if ContainsAny([]string{"A"}, nil) {
		t.Error("nil subCollection should yield false")
	}
}
