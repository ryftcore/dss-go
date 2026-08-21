package model

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestSpDocSpecificationRoundTripAndEquals(t *testing.T) {
	a := NewSpDocSpecification()
	a.SetId("2.2.25.1")
	a.SetDescription("desc")
	a.SetDocumentationReferences("http://doc")
	a.SetQualifier(enumerations.ObjectIdentifierQualifier_OID_AS_URN)

	b := NewSpDocSpecification()
	b.SetId("2.2.25.1")
	b.SetDescription("desc")
	b.SetDocumentationReferences("http://doc")
	b.SetQualifier(enumerations.ObjectIdentifierQualifier_OID_AS_URN)

	if a.Id() != "2.2.25.1" || a.Description() != "desc" {
		t.Fatalf("Id()/Description() = %q/%q", a.Id(), a.Description())
	}
	if !a.Equals(b) {
		t.Fatal("expected equal SpDocSpecification values to be Equals")
	}

	b.SetId("other")
	if a.Equals(b) {
		t.Fatal("expected differing SpDocSpecification values to not be Equals")
	}
}
