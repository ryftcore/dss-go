// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against known org.apache.commons.lang3.ArrayUtils
// semantics.

package utils

import (
	"reflect"
	"testing"
)

func TestIsArrayEmpty(t *testing.T) {
	if !IsArrayEmpty[byte](nil) {
		t.Error("nil byte array should be empty")
	}
	if !IsArrayEmpty([]byte{}) {
		t.Error("empty byte array should be empty")
	}
	if IsArrayEmpty([]byte{1}) {
		t.Error("non-empty byte array should not be empty")
	}
	if !IsArrayEmpty[rune](nil) {
		t.Error("nil rune array should be empty")
	}
	if !IsArrayEmpty[string](nil) {
		t.Error("nil Object[] array should be empty")
	}
	if IsArrayNotEmpty[byte](nil) {
		t.Error("IsArrayNotEmpty(nil) should be false")
	}
	if !IsArrayNotEmpty([]byte{1}) {
		t.Error("IsArrayNotEmpty([1]) should be true")
	}
}

func TestArraySize(t *testing.T) {
	if got := ArraySize[byte](nil); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
	if got := ArraySize([]byte{1, 2, 3}); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestSubarray(t *testing.T) {
	src := []byte{0, 1, 2, 3, 4}
	cases := []struct {
		start, end int
		want       []byte
	}{
		{1, 3, []byte{1, 2}},
		{0, 5, []byte{0, 1, 2, 3, 4}},
		{-2, 3, []byte{0, 1, 2}},  // negative start clipped to 0
		{2, 100, []byte{2, 3, 4}}, // end beyond length clipped
		{3, 3, []byte{}},          // start == end
		{4, 1, []byte{}},          // start >= end after clipping
	}
	for _, c := range cases {
		got := Subarray(src, c.start, c.end)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Subarray(src, %d, %d) = %v, want %v", c.start, c.end, got, c.want)
		}
	}
	if got := Subarray(nil, 0, 1); got != nil {
		t.Errorf("Subarray(nil, ...) = %v, want nil", got)
	}
	// Returned slice must be an independent copy.
	got := Subarray(src, 0, 2)
	got[0] = 99
	if src[0] == 99 {
		t.Error("Subarray must return a copy, not a view")
	}
}

func TestConcat(t *testing.T) {
	got := Concat([]byte("Nowina"), []byte("123"))
	if string(got) != "Nowina123" {
		t.Errorf("got %q", got)
	}
	if got := Concat(); len(got) != 0 {
		t.Errorf("Concat() = %v, want empty", got)
	}
	// nil element treated as empty rather than panicking.
	got = Concat([]byte("a"), nil, []byte("b"))
	if string(got) != "ab" {
		t.Errorf("got %q, want \"ab\"", got)
	}
}

func TestStartsWith(t *testing.T) {
	cases := []struct {
		array, prefix []byte
		want          bool
	}{
		{[]byte("hello world"), []byte("hello"), true},
		{[]byte("hello"), []byte("hello world"), false},
		{[]byte("hello"), []byte(""), true},
		{nil, []byte("a"), false},
		{[]byte("a"), nil, false},
	}
	for _, c := range cases {
		if got := StartsWith(c.array, c.prefix); got != c.want {
			t.Errorf("StartsWith(%v, %v) = %v, want %v", c.array, c.prefix, got, c.want)
		}
	}
}
