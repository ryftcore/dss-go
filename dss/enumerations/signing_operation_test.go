package enumerations

import "testing"

func TestSigningOperationValues(t *testing.T) {
	want := []SigningOperation{
		SigningOperationSign,
		SigningOperationCounterSign,
		SigningOperationTimestamp,
		SigningOperationExtend,
		SigningOperationAddSigPolicyStore,
		SigningOperationAddEvidenceRecord,
		SigningOperationEAAPresentation,
	}
	got := SigningOperationValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
