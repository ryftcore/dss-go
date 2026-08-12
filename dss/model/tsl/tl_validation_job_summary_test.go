package tsl

import "testing"

func TestNewTLValidationJobSummaryRequiresLOTLOrTL(t *testing.T) {
	if _, err := NewTLValidationJobSummary(nil, nil); err == nil {
		t.Fatalf("expected an error when neither LOTL nor TL info is provided")
	}
}

func TestTLValidationJobSummaryCounts(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	lotlInfo.SetTlInfos([]*TLInfo{
		NewTLInfo(nil, nil, nil, "https://example.org/tl1.xml"),
		NewTLInfo(nil, nil, nil, "https://example.org/tl2.xml"),
	})
	otherTL := NewTLInfo(nil, nil, nil, "https://example.org/other.xml")

	summary, err := NewTLValidationJobSummary([]*LOTLInfo{&lotlInfo}, []*TLInfo{otherTL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.NumberOfProcessedLOTLs() != 1 {
		t.Fatalf("NumberOfProcessedLOTLs() = %d, want 1", summary.NumberOfProcessedLOTLs())
	}
	if summary.NumberOfProcessedTLs() != 3 {
		t.Fatalf("NumberOfProcessedTLs() = %d, want 3 (2 from LOTL + 1 other)", summary.NumberOfProcessedTLs())
	}
	if len(summary.DocumentListInfos()) != 1 {
		t.Fatalf("unexpected DocumentListInfos: %v", summary.DocumentListInfos())
	}
	if len(summary.OtherDocumentInfos()) != 1 {
		t.Fatalf("unexpected OtherDocumentInfos: %v", summary.OtherDocumentInfos())
	}
}

func TestTLValidationJobSummaryLookupByID(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	tlInfo := NewTLInfo(nil, nil, nil, "https://example.org/tl1.xml")
	lotlInfo.SetTlInfos([]*TLInfo{tlInfo})
	otherTL := NewTLInfo(nil, nil, nil, "https://example.org/other.xml")

	summary, err := NewTLValidationJobSummary([]*LOTLInfo{&lotlInfo}, []*TLInfo{otherTL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.LOTLInfoByID(lotlInfo.DSSID()) != &lotlInfo {
		t.Fatalf("expected LOTLInfoByID to find the LOTLInfo by its identifier")
	}
	if summary.TLInfoByID(tlInfo.DSSID()) != tlInfo {
		t.Fatalf("expected TLInfoByID to find a nested LOTL TLInfo by its identifier")
	}
	if summary.TLInfoByID(otherTL.DSSID()) != otherTL {
		t.Fatalf("expected TLInfoByID to find an other TLInfo by its identifier")
	}

	unknown := NewTLInfo(nil, nil, nil, "https://example.org/unknown.xml")
	if summary.TLInfoByID(unknown.DSSID()) != nil {
		t.Fatalf("expected TLInfoByID to return nil for an unknown identifier")
	}
}
