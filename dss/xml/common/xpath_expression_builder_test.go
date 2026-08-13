// KAT test for xpath_expression_builder.go, dumped from the same Java oracle approach
// described in xmldsig_path_test.go, run directly against XPathExpressionBuilder.
package common

import "testing"

func TestXPathExpressionBuilder_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{
			"all().element(SIGNATURE).build()",
			NewXPathExpressionBuilder().All().Element(XMLDSigElement_SIGNATURE).Build(),
			"//ds:Signature",
		},
		{
			"fromCurrentPosition().element(SIGNATURE).build()",
			NewXPathExpressionBuilder().FromCurrentPosition().Element(XMLDSigElement_SIGNATURE).Build(),
			"./ds:Signature",
		},
		{
			"all().fromCurrentPosition().elements(SIGNATURE,OBJECT).build()",
			NewXPathExpressionBuilder().All().FromCurrentPosition().Elements(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT).Build(),
			".//ds:Signature/ds:Object",
		},
		{
			"all().elements(SIGNATURE,OBJECT).notParentOf(MANIFEST).build()",
			NewXPathExpressionBuilder().All().Elements(XMLDSigElement_SIGNATURE, XMLDSigElement_OBJECT).NotParentOf(XMLDSigElement_MANIFEST).Build(),
			"//ds:Signature/ds:Object[not(parent::ds:Manifest)]",
		},
		{
			"all().element(SIGNATURE).notParentOf(OBJECT).build()",
			NewXPathExpressionBuilder().All().Element(XMLDSigElement_SIGNATURE).NotParentOf(XMLDSigElement_OBJECT).Build(),
			"//ds:Signature[not(parent::ds:Object)]",
		},
		{
			"fromCurrentPosition().element(SIGNATURE_VALUE).attribute(ID).build()",
			NewXPathExpressionBuilder().FromCurrentPosition().Element(XMLDSigElement_SIGNATURE_VALUE).Attribute(XMLDSigAttribute_ID).Build(),
			"./ds:SignatureValue/@Id",
		},
		{
			"all().fromCurrentPosition().element(SIGNATURE).build()",
			NewXPathExpressionBuilder().All().FromCurrentPosition().Element(XMLDSigElement_SIGNATURE).Build(),
			".//ds:Signature",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.got != c.want {
				t.Errorf("got %q, want %q", c.got, c.want)
			}
		})
	}
}

func TestXPathExpressionBuilder_BuildPanicsWithoutAllOrFromCurrentPosition(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic when neither All nor FromCurrentPosition was set")
		}
		if r != "Unsupported operation" {
			t.Fatalf("panic = %v, want %q", r, "Unsupported operation")
		}
	}()
	NewXPathExpressionBuilder().Element(XMLDSigElement_SIGNATURE).Build()
}
