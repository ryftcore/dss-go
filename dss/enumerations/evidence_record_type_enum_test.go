package enumerations

import "testing"

func TestEvidenceRecordTypeEnumFromLabel(t *testing.T) {
	cases := []struct {
		v     EvidenceRecordTypeEnum
		label string
	}{
		{EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD, "XML Evidence Record"},
		{EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD, "ASN.1 Evidence Record"},
	}
	if len(EvidenceRecordTypeEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(EvidenceRecordTypeEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		got, err := EvidenceRecordTypeEnumFromLabel(c.label)
		if err != nil || got != c.v {
			t.Errorf("EvidenceRecordTypeEnumFromLabel(%q) = %v, %v; want %v, nil", c.label, got, err, c.v)
		}
	}
	if _, err := EvidenceRecordTypeEnumFromLabel("nope"); err == nil {
		t.Error("expected error for unknown label")
	}
}
