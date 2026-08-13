// KAT test for xpath_query_builder.go and abstract_path.go's general XPath-building
// behaviour (not just XMLDSigPath's fixed constants): dumped from the same Java oracle run
// described in xmldsig_path_test.go, exercising XPathQueryBuilder directly.
package common

import "testing"

func TestAbstractPath_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  XPathQuery
		want string
	}{
		{"all(SIGNATURE)", All(XMLDSigElement_SIGNATURE), "//ds:Signature"},
		{"all(SIGNATURE,OBJECT)", All(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT), "//ds:Signature/ds:Object"},
		{"fromCurrentPosition(SIGNATURE)", FromCurrentPosition(XMLDSigElement_SIGNATURE), "./ds:Signature"},
		{"fromCurrentPosition(SIGNATURE,OBJECT)", FromCurrentPosition(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT), "./ds:Signature/ds:Object"},
		{"allFromCurrentPosition(SIGNATURE)", AllFromCurrentPosition(XMLDSigElement_SIGNATURE), ".//ds:Signature"},
		{"allNotParent(SIGNATURE,OBJECT)", AllNotParent(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT), "//ds:Signature[not(parent::ds:Object)]"},
		{"fromCurrentPosition(SIGNATURE,ID)", FromCurrentPositionAttribute(XMLDSigElement_SIGNATURE, XMLDSigAttribute_ID), "./ds:Signature/@Id"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.got.QueryString(); got != c.want {
				t.Errorf("%s.QueryString() = %q, want %q", c.name, got, c.want)
			}
		})
	}
}

func TestXPathQueryBuilder_KAT(t *testing.T) {
	t.Run("attribute with value", func(t *testing.T) {
		q := XPathQueryBuilderFromCurrentPosition().
			Element(XMLDSigElement_REFERENCE).
			AttributeWithValue(XMLDSigAttribute_TYPE, XMLDSigPath_OBJECT_TYPE).
			Build()
		want := "./ds:Reference[@*[local-name()='Type']='http://www.w3.org/2000/09/xmldsig#Object']"
		if got := q.QueryString(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("idValue", func(t *testing.T) {
		q := XPathQueryBuilderAll().
			Element(XMLDSigElement_SIGNATURE).
			IdValue("abc123").
			Build()
		want := "//ds:Signature[@*[local-name()='Id']='abc123' or @*[local-name()='id']='abc123' or @*[local-name()='ID']='abc123']"
		if got := q.QueryString(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("fromCurrentPosition no elements uses AnyItem", func(t *testing.T) {
		q := XPathQueryBuilderFromCurrentPosition().Build()
		if got := q.QueryString(); got != "./*" {
			t.Errorf("got %q, want %q", got, "./*")
		}
	})

	t.Run("all no elements uses AnyItem", func(t *testing.T) {
		q := XPathQueryBuilderAll().Build()
		if got := q.QueryString(); got != "//*" {
			t.Errorf("got %q, want %q", got, "//*")
		}
	})

	t.Run("attribute without value is an independent node", func(t *testing.T) {
		q := XPathQueryBuilderFromCurrentPosition().
			Element(XMLDSigElement_SIGNATURE).
			Attribute(XMLDSigAttribute_ID).
			Build()
		if got := q.QueryString(); got != "./ds:Signature/@Id" {
			t.Errorf("got %q, want %q", got, "./ds:Signature/@Id")
		}
	})

	t.Run("elements + notChildOf + attribute with value", func(t *testing.T) {
		q := XPathQueryBuilderAll().
			Elements(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT).
			NotChildOf(XMLDSigElement_MANIFEST).
			AttributeWithValue(XMLDSigAttribute_ID, "xyz").
			Build()
		want := "//ds:Signature/ds:Object[not(parent::ds:Manifest)][@*[local-name()='Id']='xyz']"
		if got := q.QueryString(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("allFromCurrentPosition via SetAll no elements", func(t *testing.T) {
		q := XPathQueryBuilderFromCurrentPosition().SetAll(true).Build()
		if got := q.QueryString(); got != ".//*" {
			t.Errorf("got %q, want %q", got, ".//*")
		}
	})

	t.Run("flags", func(t *testing.T) {
		q := XPathQueryBuilderFromCurrentPosition().
			Element(XMLDSigElement_REFERENCE).
			AttributeWithValue(XMLDSigAttribute_TYPE, XMLDSigPath_OBJECT_TYPE).
			Build()
		if q.IsAll() {
			t.Errorf("IsAll() = true, want false")
		}
		if !q.IsFromCurrentPosition() {
			t.Errorf("IsFromCurrentPosition() = false, want true")
		}
		if q.IsEmpty() {
			t.Errorf("IsEmpty() = true, want false")
		}
	})
}
