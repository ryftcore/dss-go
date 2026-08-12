package enumerations

import "testing"

func TestArchiveTimestampTypeValues(t *testing.T) {
	want := []ArchiveTimestampType{
		ArchiveTimestampType_XAdES,
		ArchiveTimestampType_XAdES_141,
		ArchiveTimestampType_CAdES,
		ArchiveTimestampType_CAdES_V2,
		ArchiveTimestampType_CAdES_V3,
		ArchiveTimestampType_CAdES_DETACHED,
		ArchiveTimestampType_JAdES,
		ArchiveTimestampType_CB_AdES,
		ArchiveTimestampType_PAdES,
		ArchiveTimestampType_XML_EVIDENCE_RECORD,
		ArchiveTimestampType_ASN1_EVIDENCE_RECORD,
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
		"XAdES": ArchiveTimestampType_XAdES, "XAdES_141": ArchiveTimestampType_XAdES_141,
		"CAdES": ArchiveTimestampType_CAdES, "CAdES_V2": ArchiveTimestampType_CAdES_V2,
		"CAdES_V3": ArchiveTimestampType_CAdES_V3, "CAdES_DETACHED": ArchiveTimestampType_CAdES_DETACHED,
		"JAdES": ArchiveTimestampType_JAdES, "CB_AdES": ArchiveTimestampType_CB_AdES,
		"PAdES": ArchiveTimestampType_PAdES, "XML_EVIDENCE_RECORD": ArchiveTimestampType_XML_EVIDENCE_RECORD,
		"ASN1_EVIDENCE_RECORD": ArchiveTimestampType_ASN1_EVIDENCE_RECORD,
	}
	for name, v := range names {
		if string(v) != name {
			t.Errorf("%v: string value = %q, want %q", v, string(v), name)
		}
	}
}
