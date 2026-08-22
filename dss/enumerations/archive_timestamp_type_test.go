package enumerations

import "testing"

func TestArchiveTimestampTypeValues(t *testing.T) {
	want := []ArchiveTimestampType{
		ArchiveTimestampTypeXAdES,
		ArchiveTimestampTypeXAdES141,
		ArchiveTimestampTypeCAdES,
		ArchiveTimestampTypeCAdESV2,
		ArchiveTimestampTypeCAdESV3,
		ArchiveTimestampTypeCAdESDetached,
		ArchiveTimestampTypeJAdES,
		ArchiveTimestampTypeCBAdES,
		ArchiveTimestampTypePAdES,
		ArchiveTimestampTypeXMLEvidenceRecord,
		ArchiveTimestampTypeASN1EvidenceRecord,
	}
	got := ArchiveTimestampTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	names := map[string]ArchiveTimestampType{
		"XAdES": ArchiveTimestampTypeXAdES, "XAdES_141": ArchiveTimestampTypeXAdES141,
		"CAdES": ArchiveTimestampTypeCAdES, "CAdES_V2": ArchiveTimestampTypeCAdESV2,
		"CAdES_V3": ArchiveTimestampTypeCAdESV3, "CAdES_DETACHED": ArchiveTimestampTypeCAdESDetached,
		"JAdES": ArchiveTimestampTypeJAdES, "CB_AdES": ArchiveTimestampTypeCBAdES,
		"PAdES": ArchiveTimestampTypePAdES, "XML_EVIDENCE_RECORD": ArchiveTimestampTypeXMLEvidenceRecord,
		"ASN1_EVIDENCE_RECORD": ArchiveTimestampTypeASN1EvidenceRecord,
	}
	for name, v := range names {
		if string(v) != name {
			t.Errorf("%v: string value = %q, want %q", v, string(v), name)
		}
	}
}
