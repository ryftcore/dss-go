package lote

import "testing"

func TestNewLoTEValidationJobSummaryErrorsWhenEmpty(t *testing.T) {
	if _, err := NewLoTEValidationJobSummary(nil, nil); err == nil {
		t.Fatal("expected an error when both loloteInfos and otherLoTEInfos are empty")
	}
}

func TestLoTEValidationJobSummaryRoundTrip(t *testing.T) {
	lolote := NewLoLoTEInfo(nil, nil, nil, "https://example.org/lolote.xml")
	child := NewLoTEInfoWithParent(nil, nil, nil, "https://example.org/child.xml", lolote)
	lolote.SetChildrenInfos([]*LoTEInfo{child})

	other := NewLoTEInfo(nil, nil, nil, "https://example.org/other.xml")

	summary, err := NewLoTEValidationJobSummary([]*LoLoTEInfo{lolote}, []*LoTEInfo{other})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.NumberOfProcessedLoLoTEs() != 1 {
		t.Fatalf("NumberOfProcessedLoLoTEs() = %d, want 1", summary.NumberOfProcessedLoLoTEs())
	}
	// One "other" LoTE plus one child of the LoLoTE.
	if summary.NumberOfProcessedLoTEs() != 2 {
		t.Fatalf("NumberOfProcessedLoTEs() = %d, want 2", summary.NumberOfProcessedLoTEs())
	}
	if summary.LoTEInfoByID(other.DSSID()) != other {
		t.Fatal("expected LoTEInfoByID() to find the other LoTEInfo")
	}
	if summary.LoTEInfoByID(child.DSSID()) != child {
		t.Fatal("expected LoTEInfoByID() to find the LoLoTE's child LoTEInfo")
	}
	if summary.LoLoTEInfoByID(lolote.DSSID()) != lolote {
		t.Fatal("expected LoLoTEInfoByID() to find the LoLoTEInfo")
	}
	if len(summary.DocumentListInfos()) != 1 || summary.DocumentListInfos()[0] != lolote {
		t.Fatalf("unexpected DocumentListInfos(): %v", summary.DocumentListInfos())
	}
	if len(summary.OtherDocumentInfos()) != 1 || summary.OtherDocumentInfos()[0] != other {
		t.Fatalf("unexpected OtherDocumentInfos(): %v", summary.OtherDocumentInfos())
	}
}
