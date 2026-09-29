package jose

import "testing"

// TestJSONWritesTypedSlices pins the rendering of []string and []*Object, which writeJSONValue
// writes straight off the slice: they must stay what the []any conversion they replaced gave.
func TestJSONWritesTypedSlices(t *testing.T) {
	object := NewObject()
	object.Put("a", "b")

	cases := []struct {
		name  string
		value any
		want  string
	}{
		{"nil []string", []string(nil), `[]`},
		{"empty []string", []string{}, `[]`},
		{"[]string", []string{"a", `b"c`, "d/e"}, `["a","b\"c","d/e"]`},
		{"nil []*Object", []*Object(nil), `[]`},
		{"[]*Object", []*Object{object, nil, NewObject()}, `[{"a":"b"},null,{}]`},
		{"nested", []any{[]string{"x"}, []*Object{object}}, `[["x"],[{"a":"b"}]]`},
	}
	for _, testCase := range cases {
		if got := JSON(testCase.value); got != testCase.want {
			t.Errorf("%s: JSON = %s, want %s", testCase.name, got, testCase.want)
		}
	}
}
