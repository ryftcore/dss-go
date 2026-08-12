package enumerations

import "testing"

func TestTimestampContainerFormReadable(t *testing.T) {
	tests := []struct {
		v    TimestampContainerForm
		want string
	}{
		{TimestampContainerForm_PDF, "PDF"},
		{TimestampContainerForm_ASiC_E, "ASiC-E"},
		{TimestampContainerForm_ASiC_S, "ASiC-S"},
	}
	for _, tt := range tests {
		if got := tt.v.Readable(); got != tt.want {
			t.Errorf("%v.Readable() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestTimestampContainerFormValueOf(t *testing.T) {
	for _, v := range TimestampContainerFormValues() {
		got, err := TimestampContainerFormValueOf(string(v))
		if err != nil {
			t.Errorf("TimestampContainerFormValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("TimestampContainerFormValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestTimestampContainerFormValueOfUnknown(t *testing.T) {
	if _, err := TimestampContainerFormValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
