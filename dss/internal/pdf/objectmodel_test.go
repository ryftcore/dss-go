// Tests for the joint object model (object.go, §4.1). Named objectmodel_test.go
// rather than object_test.go so the two implementers cannot collide on the file.

package pdf

import (
	"math"
	"strings"
	"testing"
)

// TestDictPreservesInsertionOrder is review-checklist item 1: Set on an existing
// key updates in place and keeps its position, the same as LinkedHashMap.put.
// PdfDict.list() exposes this order to SingleDssDict.extractVRIs, so a map alone
// would be a determinism bug.
func TestDictPreservesInsertionOrder(t *testing.T) {
	d := NewDict()
	for _, k := range []Name{"Z", "Y", "X", "W"} {
		d.Set(k, Integer(1))
	}
	d.Set("Y", Integer(99)) // update, not append
	want := []Name{"Z", "Y", "X", "W"}
	got := d.Keys()
	if len(got) != len(want) {
		t.Fatalf("keys = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %v, want %v", got, want)
		}
	}
	if v, _ := d.GetRaw("Y").(Integer); v != 99 {
		t.Errorf("/Y = %v, want the updated value", v)
	}
	if d.Len() != 4 {
		t.Errorf("Len = %d", d.Len())
	}
	// Keys returns a fresh slice.
	got[0] = "mutated"
	if d.Keys()[0] != "Z" {
		t.Error("Keys() handed out the internal slice")
	}
}

func TestDictDeleteKeepsOrder(t *testing.T) {
	d := DictOf(Name("A"), Integer(1), Name("B"), Integer(2), Name("C"), Integer(3))
	d.Delete("B")
	got := d.Keys()
	if len(got) != 2 || got[0] != "A" || got[1] != "C" {
		t.Fatalf("keys after delete = %v", got)
	}
	if d.Has("B") || d.GetRaw("B") != nil {
		t.Error("/B is still present")
	}
	// The index map must have been rebuilt: /C still resolves.
	if v, _ := d.GetRaw("C").(Integer); v != 3 {
		t.Errorf("/C = %v after deleting /B", v)
	}
	d.Delete("missing") // no-op
	d.Set("D", Integer(4))
	if got := d.Keys(); got[2] != "D" {
		t.Errorf("keys after re-add = %v", got)
	}
}

func TestDictCloneIsShallow(t *testing.T) {
	inner := DictOf(Name("X"), Integer(1))
	d := DictOf(Name("A"), Integer(1), Name("Inner"), inner)
	c := d.Clone()
	c.Set("A", Integer(2))
	if v, _ := d.GetRaw("A").(Integer); v != 1 {
		t.Error("Clone is not independent at the top level")
	}
	if c.GetRaw("Inner") != Object(inner) {
		t.Error("Clone must share value objects")
	}
	var nilDict *Dict
	if nilDict.Clone() != nil {
		t.Error("nil.Clone() must be nil")
	}
}

func TestNilDictAccessorsAreSafe(t *testing.T) {
	var d *Dict
	if d.Len() != 0 || d.Has("A") || d.GetRaw("A") != nil || d.Keys() != nil {
		t.Error("nil dictionary accessors must be safe")
	}
	d.Delete("A")
	if !strings.Contains(d.String(), "nil") {
		t.Errorf("nil dictionary String = %q", d.String())
	}
}

func TestDictOfPanicsOnMalformedPairs(t *testing.T) {
	for _, kv := range [][]any{
		{Name("A")},
		{"A", Integer(1)},
		{Name("A"), 1},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("DictOf(%v) did not panic", kv)
				}
			}()
			DictOf(kv...)
		}()
	}
}

func TestObjectKindsImplementObject(t *testing.T) {
	// The compiler enforces this; the test documents the closed set of §4.1.
	var objs = []Object{
		Null{}, Bool(true), Integer(1), Real{Val: 1},
		String{Bytes: []byte("x")}, Name("N"), Array{}, NewDict(),
		NewStream(nil, nil), Ref{Num: 1},
	}
	if len(objs) != 10 {
		t.Fatalf("the object model has %d kinds", len(objs))
	}
}

func TestRealKeepsRawAndCoercesSpecials(t *testing.T) {
	if got, raw, ok := parseLenientReal("595.276"); !ok || raw != "595.276" || math.Abs(got-595.276) > 1e-3 {
		t.Errorf("parseLenientReal = %v %q %v", got, raw, ok)
	}
	// A literal that does not round-trip as a float32 keeps no raw form, so the
	// writer formats it instead of echoing something misleading.
	if _, raw, _ := parseLenientReal("1e400"); raw != "" {
		t.Errorf("an out-of-range literal kept its raw form: %q", raw)
	}
}

func TestRect(t *testing.T) {
	r, ok := RectFromArray(Array{Integer(10), Integer(20), Integer(110), Integer(220)})
	if !ok {
		t.Fatal("RectFromArray failed")
	}
	if r.Width() != 100 || r.Height() != 200 {
		t.Errorf("rect = %+v", r)
	}
	// corners are normalised
	flipped, ok := RectFromArray(Array{Integer(110), Integer(220), Integer(10), Integer(20)})
	if !ok || flipped != r {
		t.Errorf("flipped rect = %+v, want %+v", flipped, r)
	}
	if _, ok := RectFromArray(Array{Integer(1), Integer(2)}); ok {
		t.Error("a short array must not make a rectangle")
	}
	if _, ok := RectFromArray(Array{Name("x"), Integer(2), Integer(3), Integer(4)}); ok {
		t.Error("a non-numeric element must not make a rectangle")
	}
	arr := Rect{1, 2, 3, 4}.Array()
	if len(arr) != 4 {
		t.Fatalf("Rect.Array = %v", arr)
	}
	if _, isReal := arr[0].(Real); !isReal {
		t.Errorf("Rect.Array must emit reals, got %T", arr[0])
	}
}

func TestStreamConstructor(t *testing.T) {
	s := NewStream(nil, []byte("x"))
	if s.Dict == nil {
		t.Error("NewStream(nil, …) must still have a dictionary")
	}
	if s.Offset != 0 || s.Length != 0 {
		t.Error("an in-memory stream has no source range")
	}
}

func TestObjectStringRendering(t *testing.T) {
	d := DictOf(
		Name("N"), Name("Foo"),
		Name("A"), Array{Integer(1), Real{Val: 2.5, Raw: "2.5"}, Bool(true), Null{}},
		Name("S"), String{Bytes: []byte("hi")},
		Name("R"), Ref{Num: 3, Gen: 1},
		Name("Stream"), NewStream(NewDict(), []byte("abc")),
	)
	got := d.String()
	for _, want := range []string{"/N /Foo", "1 2.5 true null", "(hi)", "3 1 R", "stream(3)"} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %s, missing %q", got, want)
		}
	}
}
