// Ported from dss-enumerations/.../EvidenceRecordIncorporationType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestEvidenceRecordIncorporationTypeValueOf(t *testing.T) {
	for _, v := range EvidenceRecordIncorporationTypeValues() {
		got, err := EvidenceRecordIncorporationTypeValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("EvidenceRecordIncorporationTypeValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if len(EvidenceRecordIncorporationTypeValues()) != 2 {
		t.Errorf("expected 2 values, got %d", len(EvidenceRecordIncorporationTypeValues()))
	}
	if _, err := EvidenceRecordIncorporationTypeValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
