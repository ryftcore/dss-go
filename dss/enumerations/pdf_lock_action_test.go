package enumerations

import "testing"

func TestPdfLockActionForName(t *testing.T) {
	tests := []struct {
		v    PdfLockAction
		name string
	}{
		{PdfLockAction_ALL, "All"},
		{PdfLockAction_INCLUDE, "Include"},
		{PdfLockAction_EXCLUDE, "Exclude"},
	}
	for _, tt := range tests {
		if got := tt.v.Name(); got != tt.name {
			t.Errorf("%v.Name() = %q, want %q", tt.v, got, tt.name)
		}
		got, err := PdfLockActionForName(tt.name)
		if err != nil {
			t.Errorf("PdfLockActionForName(%q) unexpected error: %v", tt.name, err)
		}
		if got != tt.v {
			t.Errorf("PdfLockActionForName(%q) = %q, want %q", tt.name, got, tt.v)
		}
	}
}

func TestPdfLockActionForNameUnknown(t *testing.T) {
	if _, err := PdfLockActionForName("bogus"); err == nil {
		t.Error("expected error for unknown name, got nil")
	}
}

func TestPdfLockActionValueOf(t *testing.T) {
	for _, v := range PdfLockActionValues() {
		got, err := PdfLockActionValueOf(string(v))
		if err != nil {
			t.Errorf("PdfLockActionValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("PdfLockActionValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}
