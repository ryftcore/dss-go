package enumerations

import "testing"

func TestTimestampedObjectTypeValues(t *testing.T) {
	want := []TimestampedObjectType{
		TimestampedObjectTypeSignedData,
		TimestampedObjectTypeSignature,
		TimestampedObjectTypeCertificate,
		TimestampedObjectTypeRevocation,
		TimestampedObjectTypeTimestamp,
		TimestampedObjectTypeEvidenceRecord,
		TimestampedObjectTypeOrphanCertificate,
		TimestampedObjectTypeOrphanRevocation,
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
