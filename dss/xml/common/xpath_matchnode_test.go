// Tests XPathQueryItem/XPathQueryParameter.MatchNode against real *xmldom.Node trees, i.e.
// the "process" half of the framework that dss-xml-utils's NativeDOMXPathQueryExecutor (a
// later phase) will drive - see doc.go and xpath_query_item.go.
package common

import (
	"testing"

	"github.com/utain/esig/dss/internal/xmldom"
)

func parseMatchNodeTestDoc(t *testing.T) *xmldom.Node {
	t.Helper()
	src := []byte(`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="sig1">
  <ds:SignedInfo>
    <ds:Reference Type="http://www.w3.org/2000/09/xmldsig#Object" URI="#o1"/>
  </ds:SignedInfo>
  <ds:Object ID="o1"/>
</ds:Signature>`)
	doc, err := xmldom.Parse(src, nil)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	return doc
}

func TestXPathQueryElementItem_MatchNode(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()

	sigItem := NewXPathQueryElementItem(XMLDSigElement_SIGNATURE)
	if !sigItem.MatchNode(root) {
		t.Errorf("SIGNATURE item should match the ds:Signature root element")
	}

	objItem := NewXPathQueryElementItem(XMLDSigElement_OBJECT)
	if objItem.MatchNode(root) {
		t.Errorf("OBJECT item should not match the ds:Signature root element")
	}

	object := root.Elements()[1]
	if !objItem.MatchNode(object) {
		t.Errorf("OBJECT item should match the ds:Object element")
	}

	// An attribute node is never Element-related.
	idAttr := root.Attr("", "Id")
	if idAttr == nil {
		t.Fatalf("Id attribute not found on root")
	}
	if sigItem.MatchNode(idAttr) {
		t.Errorf("element item should not match an attribute node")
	}
}

func TestXPathQueryAnyItem_MatchNode(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()
	any := NewXPathQueryAnyItem()
	if !any.MatchNode(root) {
		t.Errorf("AnyItem should match any element node")
	}
	idAttr := root.Attr("", "Id")
	if any.MatchNode(idAttr) {
		t.Errorf("AnyItem should not match an attribute node")
	}
}

func TestXPathQueryAttributeItem_MatchNode(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()
	idItem := NewXPathQueryAttributeItem(XMLDSigAttribute_ID)

	idAttr := root.Attr("", "Id")
	if idAttr == nil {
		t.Fatalf("Id attribute not found")
	}
	if !idItem.MatchNode(idAttr) {
		t.Errorf("attribute item for Id should match the Id attribute node")
	}
	if idItem.MatchNode(root) {
		t.Errorf("attribute item should not match an element node")
	}

	uriItem := NewXPathQueryAttributeItem(XMLDSigAttribute_URI)
	if uriItem.MatchNode(idAttr) {
		t.Errorf("URI attribute item should not match the Id attribute node")
	}
}

func TestXPathQueryAttributeParameter_MatchNode(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()
	signedInfo := root.Elements()[0]
	reference := signedInfo.Elements()[0]

	typeParam := NewXPathQueryAttributeParameter(XMLDSigAttribute_TYPE, XMLDSigPath_OBJECT_TYPE)
	if !typeParam.MatchNode(reference) {
		t.Errorf("Type attribute parameter should match ds:Reference with the expected Type value")
	}

	wrongValueParam := NewXPathQueryAttributeParameter(XMLDSigAttribute_TYPE, XMLDSigPath_MANIFEST_TYPE)
	if wrongValueParam.MatchNode(reference) {
		t.Errorf("Type attribute parameter with a mismatched value should not match")
	}
}

func TestXPathQueryIdentifierParameter_MatchNode_CaseInsensitive(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()
	object := root.Elements()[1] // carries ID="o1" (uppercase attribute name)

	idParam := NewXPathQueryIdentifierParameter("o1")
	if !idParam.MatchNode(object) {
		t.Errorf("identifier parameter should match ID=o1 case-insensitively")
	}

	wrongParam := NewXPathQueryIdentifierParameter("no-such-id")
	if wrongParam.MatchNode(object) {
		t.Errorf("identifier parameter with the wrong id value should not match")
	}
}

func TestXPathQueryNotChildOfParameter_MatchNode(t *testing.T) {
	doc := parseMatchNodeTestDoc(t)
	root := doc.DocumentElement()
	signedInfo := root.Elements()[0]
	reference := signedInfo.Elements()[0]
	object := root.Elements()[1]

	notChildOfSignedInfo := NewXPathQueryNotChildOfParameter(XMLDSigElement_SIGNED_INFO)
	if notChildOfSignedInfo.MatchNode(reference) {
		t.Errorf("ds:Reference is a child of ds:SignedInfo, so notChildOf(SIGNED_INFO) should be false")
	}
	if !notChildOfSignedInfo.MatchNode(object) {
		t.Errorf("ds:Object is not a child of ds:SignedInfo, so notChildOf(SIGNED_INFO) should be true")
	}

	notChildOfManifest := NewXPathQueryNotChildOfParameter(XMLDSigElement_MANIFEST)
	if !notChildOfManifest.MatchNode(object) {
		t.Errorf("ds:Object's parent is ds:Signature, not ds:Manifest, so notChildOf(MANIFEST) should be true")
	}
}
