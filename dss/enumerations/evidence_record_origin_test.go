// Ported from dss-enumerations/.../EvidenceRecordOrigin.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestEvidenceRecordOriginValueOf(t *testing.T) {
	for _, v := range EvidenceRecordOriginValues() {
		got, err := EvidenceRecordOriginValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("EvidenceRecordOriginValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if len(EvidenceRecordOriginValues()) != 3 {
		t.Errorf("expected 3 values, got %d", len(EvidenceRecordOriginValues()))
	}
	if _, err := EvidenceRecordOriginValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
