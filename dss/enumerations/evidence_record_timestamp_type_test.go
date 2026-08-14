// Ported from dss-enumerations/.../EvidenceRecordTimestampType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestEvidenceRecordTimestampTypeValueOf(t *testing.T) {
	for _, v := range EvidenceRecordTimestampTypeValues() {
		got, err := EvidenceRecordTimestampTypeValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("EvidenceRecordTimestampTypeValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if len(EvidenceRecordTimestampTypeValues()) != 3 {
		t.Errorf("expected 3 values, got %d", len(EvidenceRecordTimestampTypeValues()))
	}
	if _, err := EvidenceRecordTimestampTypeValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
