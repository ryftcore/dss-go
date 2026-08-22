package enumerations

import "testing"

func TestCOSEStructureTypeValueOf(t *testing.T) {
	for _, v := range COSEStructureTypeValues() {
		got, err := COSEStructureTypeValueOf(string(v))
		if err != nil {
			t.Fatalf("COSEStructureTypeValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("COSEStructureTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := COSEStructureTypeValueOf("bogus"); err == nil {
		t.Error("COSEStructureTypeValueOf(\"bogus\") expected error, got nil")
	}
}

func TestCOSEStructureTypeValues(t *testing.T) {
	want := []COSEStructureType{COSEStructureTypeCoseSign, COSEStructureTypeCoseSign1}
	got := COSEStructureTypeValues()
	if len(got) != len(want) {
		t.Fatalf("COSEStructureTypeValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("COSEStructureTypeValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
