package common

import "testing"

func TestDSSNamespace(t *testing.T) {
	ns := NewDSSNamespace("urn:test", "t")
	if ns.Uri() != "urn:test" {
		t.Errorf("Uri() = %q, want %q", ns.Uri(), "urn:test")
	}
	if ns.Prefix() != "t" {
		t.Errorf("Prefix() = %q, want %q", ns.Prefix(), "t")
	}
	if !ns.IsSameUri("urn:test") {
		t.Errorf("IsSameUri(urn:test) = false, want true")
	}
	if ns.IsSameUri("urn:other") {
		t.Errorf("IsSameUri(urn:other) = true, want false")
	}
	// Verbatim upstream toString() quirk: no closing quote after the uri value.
	want := "DSSNamespace [uri='urn:test, prefix='t]"
	if got := ns.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestDSSElementFromDefinition(t *testing.T) {
	ns := NewDSSNamespace("urn:test", "t")
	e := DSSElementFromDefinition("Foo", ns)
	if e.TagName() != "Foo" {
		t.Errorf("TagName() = %q, want %q", e.TagName(), "Foo")
	}
	if e.Namespace() != ns {
		t.Errorf("Namespace() did not return the same DSSNamespace instance")
	}
	if e.URI() != "urn:test" {
		t.Errorf("URI() = %q, want %q", e.URI(), "urn:test")
	}
	if !e.IsSameTagName("Foo") || e.IsSameTagName("Bar") {
		t.Errorf("IsSameTagName behaved incorrectly")
	}

	noNS := DSSElementFromDefinition("Foo", nil)
	if noNS.Namespace() != nil {
		t.Errorf("Namespace() = %v, want nil", noNS.Namespace())
	}
	if noNS.URI() != "" {
		t.Errorf("URI() = %q, want \"\" when Namespace is nil", noNS.URI())
	}
}

func TestDSSAttributeFromDefinition(t *testing.T) {
	a := DSSAttributeFromDefinition("Foo")
	if a.AttributeName() != "Foo" {
		t.Errorf("AttributeName() = %q, want %q", a.AttributeName(), "Foo")
	}
}

func TestXMLDSigElement_AttributeNames(t *testing.T) {
	if XMLDSigElement_SIGNATURE.TagName() != "Signature" {
		t.Errorf("SIGNATURE.TagName() = %q, want %q", XMLDSigElement_SIGNATURE.TagName(), "Signature")
	}
	if XMLDSigElement_SIGNATURE.URI() != "http://www.w3.org/2000/09/xmldsig#" {
		t.Errorf("SIGNATURE.URI() = %q", XMLDSigElement_SIGNATURE.URI())
	}
	if !XMLDSigElement_SIGNATURE.IsSameTagName("Signature") {
		t.Errorf("IsSameTagName(Signature) = false, want true")
	}
	if XMLDSigAttribute_MIME_TYPE.AttributeName() != "MimeType" {
		t.Errorf("MIME_TYPE.AttributeName() = %q, want %q", XMLDSigAttribute_MIME_TYPE.AttributeName(), "MimeType")
	}
}
