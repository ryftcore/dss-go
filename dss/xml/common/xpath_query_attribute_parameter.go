// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryAttributeParameter.java (DSS 6.5.RC1).
package common

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// XPathQueryAttributeParameter allows extraction of an element by a given attribute and
// expected value.
type XPathQueryAttributeParameter struct {
	*AbstractXPathQueryParameter
	attribute      DSSAttribute
	attributeValue string
	ignoreCase     bool
}

// NewXPathQueryAttributeParameter creates a parameter with an expected attribute value,
// enforcing case-sensitive comparison. Ports the two-argument constructor.
func NewXPathQueryAttributeParameter(attribute DSSAttribute, attributeValue string) *XPathQueryAttributeParameter {
	return NewXPathQueryAttributeParameterIgnoreCase(attribute, attributeValue, false)
}

// NewXPathQueryAttributeParameterIgnoreCase creates a parameter with an expected attribute
// value and an explicit case-sensitivity choice. Panics if attribute is nil (Java
// Objects.requireNonNull(attribute, "DSSAttribute cannot be null!")); the matching
// requireNonNull on attributeValue has no Go counterpart, since a Go string can never be
// nil.
func NewXPathQueryAttributeParameterIgnoreCase(attribute DSSAttribute, attributeValue string, ignoreCase bool) *XPathQueryAttributeParameter {
	if attribute == nil {
		panic("DSSAttribute cannot be null!")
	}
	return &XPathQueryAttributeParameter{
		AbstractXPathQueryParameter: newAbstractXPathQueryParameter(),
		attribute:                   attribute,
		attributeValue:              attributeValue,
		ignoreCase:                  ignoreCase,
	}
}

// Attribute returns the corresponding DSSAttribute. Ports getAttribute().
func (p *XPathQueryAttributeParameter) Attribute() DSSAttribute {
	return p.attribute
}

// AttributeValue returns the expected attribute value. Ports getAttributeValue().
func (p *XPathQueryAttributeParameter) AttributeValue() string {
	return p.attributeValue
}

// process reports whether node is an Element node carrying an attribute matching both name
// (respecting ignoreCase) and value. Ports process(Node).
func (p *XPathQueryAttributeParameter) process(node *xmldom.Node) bool {
	if node.Kind != xmldom.Element {
		return false
	}
	for _, attributeNode := range node.Attrs {
		var nameMatch bool
		if p.ignoreCase {
			nameMatch = strings.EqualFold(p.attribute.AttributeName(), localName(attributeNode))
		} else {
			nameMatch = p.attribute.AttributeName() == localName(attributeNode)
		}
		if nameMatch && p.attributeValue == attributeNode.Value {
			return true
		}
	}
	return false
}

// MatchNode implements XPathQueryItem. Ports the parameter's matchNode(Node), which skips
// the (always-empty, for a parameter) parameter loop and calls process directly.
func (p *XPathQueryAttributeParameter) MatchNode(node *xmldom.Node) bool {
	return p.process(node)
}

// IsElementRelated implements XPathQueryItem.
func (p *XPathQueryAttributeParameter) IsElementRelated() bool {
	return true
}

// IsAttributeRelated implements XPathQueryItem.
func (p *XPathQueryAttributeParameter) IsAttributeRelated() bool {
	return false
}

// QueryString implements XPathQueryItem.
func (p *XPathQueryAttributeParameter) QueryString() string {
	var sb strings.Builder
	xPathQueryAppendAttributeCondition(&sb, p.attribute.AttributeName(), p.attributeValue)
	if p.ignoreCase {
		sb.WriteString(" or ")
		xPathQueryAppendAttributeCondition(&sb, strings.ToLower(p.attribute.AttributeName()), p.attributeValue)
		sb.WriteString(" or ")
		xPathQueryAppendAttributeCondition(&sb, strings.ToUpper(p.attribute.AttributeName()), p.attributeValue)
	}
	return sb.String()
}

// xPathQueryAppendAttributeCondition appends "@*[local-name()='attrName']='value'". Ports
// the private addQueryForAttributeWithName(StringBuilder, String); its
// "if (attributeValue != null)" guard is unreachable here since every
// XPathQueryAttributeParameter is constructed with a non-nil (in Go: always-present) value,
// so the value clause is emitted unconditionally.
func xPathQueryAppendAttributeCondition(sb *strings.Builder, attrName, attributeValue string) {
	sb.WriteString("@*[local-name()='")
	sb.WriteString(attrName)
	sb.WriteString("']")
	sb.WriteString("='")
	sb.WriteString(attributeValue)
	sb.WriteString("'")
}

var _ XPathQueryParameter = (*XPathQueryAttributeParameter)(nil)
