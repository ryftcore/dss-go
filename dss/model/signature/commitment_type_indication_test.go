// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/CommitmentTypeIndication.java (DSS 6.5.RC1).
package signature

import (
	"reflect"
	"testing"
)

func TestCommitmentTypeIndication_RoundTrip(t *testing.T) {
	cti := NewCommitmentTypeIndication("1.2.3.4")
	if got, want := cti.Identifier(), "1.2.3.4"; got != want {
		t.Fatalf("Identifier() = %q, want %q", got, want)
	}

	cti.SetDescription("proof of origin")
	cti.SetDocumentReferences([]string{"doc1", "doc2"})
	cti.SetObjectReferences([]string{"ref1"})
	cti.SetAllDataSignedObjects(true)

	if got, want := cti.Description(), "proof of origin"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
	if got, want := cti.DocumentReferences(), []string{"doc1", "doc2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("DocumentReferences() = %v, want %v", got, want)
	}
	if got, want := cti.ObjectReferences(), []string{"ref1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ObjectReferences() = %v, want %v", got, want)
	}
	if !cti.IsAllDataSignedObjects() {
		t.Fatalf("IsAllDataSignedObjects() = false, want true")
	}
}
