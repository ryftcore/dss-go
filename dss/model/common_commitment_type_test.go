package model

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestCommonCommitmentTypeRoundTrip(t *testing.T) {
	ct := NewCommonCommitmentType()
	ct.SetOid("1.2.840.113549.1.9.16.6.1")
	ct.SetSignedDataObjects("ref1", "ref2")

	q := NewCommitmentQualifier()
	q.SetOid("1.2.3")
	ct.SetCommitmentTypeQualifiers(q)

	if ct.OID() != "1.2.840.113549.1.9.16.6.1" {
		t.Fatalf("OID() = %q", ct.OID())
	}
	if !reflect.DeepEqual(ct.SignedDataObjects(), []string{"ref1", "ref2"}) {
		t.Fatalf("SignedDataObjects() = %v", ct.SignedDataObjects())
	}
	if len(ct.CommitmentTypeQualifiers()) != 1 || ct.CommitmentTypeQualifiers()[0] != q {
		t.Fatalf("CommitmentTypeQualifiers() did not round-trip")
	}

	var _ enumerations.CommitmentType = ct
}
