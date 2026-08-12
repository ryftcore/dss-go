package enumerations

import "testing"

func TestSigningOperationValues(t *testing.T) {
	want := []SigningOperation{
		SigningOperation_SIGN,
		SigningOperation_COUNTER_SIGN,
		SigningOperation_TIMESTAMP,
		SigningOperation_EXTEND,
		SigningOperation_ADD_SIG_POLICY_STORE,
		SigningOperation_ADD_EVIDENCE_RECORD,
		SigningOperation_EAA_PRESENTATION,
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
