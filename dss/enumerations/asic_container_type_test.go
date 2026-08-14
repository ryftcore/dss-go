package enumerations

import "testing"

func TestASiCContainerTypeValueOf(t *testing.T) {
	cases := []struct {
		name string
		want ASiCContainerType
	}{
		{"ASiC_S", ASiCContainerType_ASiC_S},
		{"ASiC_E", ASiCContainerType_ASiC_E},
	}
	for _, c := range cases {
		got, err := ASiCContainerTypeValueOf(c.name)
		if err != nil {
			t.Fatalf("ASiCContainerTypeValueOf(%q) returned error: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("ASiCContainerTypeValueOf(%q) = %q, want %q", c.name, got, c.want)
		}
	}

	if _, err := ASiCContainerTypeValueOf("bogus"); err == nil {
		t.Error("ASiCContainerTypeValueOf(\"bogus\") expected error, got nil")
	}
}

func TestASiCContainerTypeValueByName(t *testing.T) {
	cases := []struct {
		name string
		want ASiCContainerType
	}{
		{"ASiC-S", ASiCContainerType_ASiC_S},
		{"ASiC-E", ASiCContainerType_ASiC_E},
		{"ASiC_S", ASiCContainerType_ASiC_S},
	}
	for _, c := range cases {
		got, err := ASiCContainerTypeValueByName(c.name)
		if err != nil {
			t.Fatalf("ASiCContainerTypeValueByName(%q) returned error: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("ASiCContainerTypeValueByName(%q) = %q, want %q", c.name, got, c.want)
		}
	}

	if _, err := ASiCContainerTypeValueByName("bogus"); err == nil {
		t.Error("ASiCContainerTypeValueByName(\"bogus\") expected error, got nil")
	}
}

func TestASiCContainerTypeString(t *testing.T) {
	if got := ASiCContainerType_ASiC_S.String(); got != "ASiC-S" {
		t.Errorf("ASiC_S.String() = %q, want %q", got, "ASiC-S")
	}
	if got := ASiCContainerType_ASiC_E.String(); got != "ASiC-E" {
		t.Errorf("ASiC_E.String() = %q, want %q", got, "ASiC-E")
	}
}

func TestASiCContainerTypeValues(t *testing.T) {
	want := []ASiCContainerType{ASiCContainerType_ASiC_S, ASiCContainerType_ASiC_E}
	got := ASiCContainerTypeValues()
	if len(got) != len(want) {
		t.Fatalf("ASiCContainerTypeValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ASiCContainerTypeValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
