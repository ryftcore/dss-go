// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryIdentifierParameter.java (DSS 6.5.RC1).
package common

// XPathQueryIdentifierParameter builds an XPath expression part for an element retrieval by
// any identifier type (suitable types are: 'Id', 'ID', 'id').
type XPathQueryIdentifierParameter struct {
	*XPathQueryAttributeParameter
}

// NewXPathQueryIdentifierParameter creates a case-insensitive Id-attribute parameter for the
// given ID value.
func NewXPathQueryIdentifierParameter(idValue string) *XPathQueryIdentifierParameter {
	return &XPathQueryIdentifierParameter{
		XPathQueryAttributeParameter: NewXPathQueryAttributeParameterIgnoreCase(XMLDSigAttribute_ID, idValue, true),
	}
}

// Id returns the ID string value. Ports getId().
func (p *XPathQueryIdentifierParameter) Id() string {
	return p.AttributeValue()
}

var _ XPathQueryParameter = (*XPathQueryIdentifierParameter)(nil)
