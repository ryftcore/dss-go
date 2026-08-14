package model

import (
	"reflect"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestCommonObjectIdentifierRoundTrip(t *testing.T) {
	o := NewCommonObjectIdentifier()
	o.SetUri("http://example.org/oid")
	o.SetOid("1.2.3")
	o.SetQualifier(enumerations.ObjectIdentifierQualifier_OID_AS_URI)
	o.SetDescription("desc")
	o.SetDocumentationReferences("http://doc1", "http://doc2")

	if o.URI() != "http://example.org/oid" {
		t.Fatalf("URI() = %q", o.URI())
	}
	if o.OID() != "1.2.3" {
		t.Fatalf("OID() = %q", o.OID())
	}
	if o.Qualifier() != enumerations.ObjectIdentifierQualifier_OID_AS_URI {
		t.Fatalf("Qualifier() = %v", o.Qualifier())
	}
	if o.Description() != "desc" {
		t.Fatalf("Description() = %q", o.Description())
	}
	want := []string{"http://doc1", "http://doc2"}
	if !reflect.DeepEqual(o.DocumentationReferences(), want) {
		t.Fatalf("DocumentationReferences() = %v, want %v", o.DocumentationReferences(), want)
	}

	var _ enumerations.ObjectIdentifier = o
}
