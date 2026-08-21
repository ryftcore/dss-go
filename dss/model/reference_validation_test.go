package model

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestReferenceValidationRoundTrip(t *testing.T) {
	rv := NewReferenceValidation()
	rv.SetType(enumerations.DigestMatcherType_REFERENCE)
	rv.SetFound(true)
	rv.SetIntact(true)
	rv.SetId("r-id")
	rv.SetUri("#r-id")
	rv.SetTransformationNames([]string{"http://www.w3.org/TR/2001/REC-xml-c14n-20010315"})

	if rv.Type() != enumerations.DigestMatcherType_REFERENCE {
		t.Fatalf("Type() = %v", rv.Type())
	}
	if !rv.IsFound() || !rv.IsIntact() {
		t.Fatal("expected found/intact to be true")
	}
	if rv.Id() != "r-id" || rv.Uri() != "#r-id" {
		t.Fatalf("Id()/Uri() = %q/%q", rv.Id(), rv.Uri())
	}
	if !reflect.DeepEqual(rv.TransformationNames(), []string{"http://www.w3.org/TR/2001/REC-xml-c14n-20010315"}) {
		t.Fatalf("TransformationNames() = %v", rv.TransformationNames())
	}
}

func TestReferenceValidationLazyDependentValidations(t *testing.T) {
	rv := NewReferenceValidation()
	if rv.DependentValidations() == nil {
		t.Fatal("expected DependentValidations() to lazily initialize")
	}
	if rv.ErrorMessages() == nil {
		t.Fatal("expected ErrorMessages() to start as an empty (non-nil) slice, per the Java field initializer")
	}
}
