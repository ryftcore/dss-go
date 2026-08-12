package enumerations

import "testing"

func TestJWSSerializationTypeValues(t *testing.T) {
	want := []JWSSerializationType{
		JWSSerializationType_COMPACT_SERIALIZATION,
		JWSSerializationType_JSON_SERIALIZATION,
		JWSSerializationType_FLATTENED_JSON_SERIALIZATION,
	}
	got := JWSSerializationTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	names := map[JWSSerializationType]string{
		JWSSerializationType_COMPACT_SERIALIZATION:        "COMPACT_SERIALIZATION",
		JWSSerializationType_JSON_SERIALIZATION:           "JSON_SERIALIZATION",
		JWSSerializationType_FLATTENED_JSON_SERIALIZATION: "FLATTENED_JSON_SERIALIZATION",
	}
	for v, name := range names {
		if string(v) != name {
			t.Errorf("%v: string value = %q, want %q", v, string(v), name)
		}
	}
}
