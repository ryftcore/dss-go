package enumerations

import "testing"

func TestTimestampedObjectTypeValues(t *testing.T) {
	want := []TimestampedObjectType{
		TimestampedObjectType_SIGNED_DATA,
		TimestampedObjectType_SIGNATURE,
		TimestampedObjectType_CERTIFICATE,
		TimestampedObjectType_REVOCATION,
		TimestampedObjectType_TIMESTAMP,
		TimestampedObjectType_EVIDENCE_RECORD,
		TimestampedObjectType_ORPHAN_CERTIFICATE,
		TimestampedObjectType_ORPHAN_REVOCATION,
	}
	got := TimestampedObjectTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
