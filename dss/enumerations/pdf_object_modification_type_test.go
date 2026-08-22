package enumerations

import "testing"

func TestPdfObjectModificationTypeValues(t *testing.T) {
	want := []PdfObjectModificationType{
		PdfObjectModificationTypeCreation,
		PdfObjectModificationTypeDeletion,
		PdfObjectModificationTypeModification,
	}
	got := PdfObjectModificationTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
