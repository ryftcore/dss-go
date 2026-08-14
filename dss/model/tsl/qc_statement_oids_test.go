package tsl

import "testing"

func TestQCStatementOidsRoundTrip(t *testing.T) {
	q := NewQCStatementOids()
	q.SetQcStatementIds([]string{"stmt1"})
	q.SetQcTypeIds([]string{"type1"})
	q.SetQcCClegislations([]string{"FR"})
	q.SetQcStatementIdsToRemove([]string{"stmt2"})
	q.SetQcTypeIdsToRemove([]string{"type2"})
	q.SetQcCClegislationsToRemove([]string{"DE"})

	if got := q.QcStatementIds(); len(got) != 1 || got[0] != "stmt1" {
		t.Fatalf("unexpected QcStatementIds: %v", got)
	}
	if got := q.QcTypeIds(); len(got) != 1 || got[0] != "type1" {
		t.Fatalf("unexpected QcTypeIds: %v", got)
	}
	if got := q.QcCClegislations(); len(got) != 1 || got[0] != "FR" {
		t.Fatalf("unexpected QcCClegislations: %v", got)
	}
	if got := q.QcStatementIdsToRemove(); len(got) != 1 || got[0] != "stmt2" {
		t.Fatalf("unexpected QcStatementIdsToRemove: %v", got)
	}
	if got := q.QcTypeIdsToRemove(); len(got) != 1 || got[0] != "type2" {
		t.Fatalf("unexpected QcTypeIdsToRemove: %v", got)
	}
	if got := q.QcCClegislationsToRemove(); len(got) != 1 || got[0] != "DE" {
		t.Fatalf("unexpected QcCClegislationsToRemove: %v", got)
	}
}

func TestQCStatementOidsEquals(t *testing.T) {
	a := NewQCStatementOids()
	a.SetQcStatementIds([]string{"stmt1"})
	b := NewQCStatementOids()
	b.SetQcStatementIds([]string{"stmt1"})

	if !a.Equals(b) {
		t.Fatal("expected equal QCStatementOids to be Equals()")
	}

	b.SetQcStatementIds([]string{"other"})
	if a.Equals(b) {
		t.Fatal("expected different QcStatementIds to not be Equals()")
	}
	if a.Equals(nil) {
		t.Fatal("expected Equals(nil) to be false")
	}
}

func TestQCStatementOidsStringNilVsEmpty(t *testing.T) {
	q := NewQCStatementOids()
	if got := q.String(); got == "" {
		t.Fatal("expected non-empty String()")
	}
	q.SetQcStatementIds([]string{})
	if got := q.QcStatementIds(); got == nil {
		t.Fatal("expected empty (non-nil) slice to stay non-nil")
	}
}
