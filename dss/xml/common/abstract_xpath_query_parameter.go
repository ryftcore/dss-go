// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/AbstractXPathQueryParameter.java (DSS 6.5.RC1).
package common

// AbstractXPathQueryParameter is the shared implementation behind every XPathQueryParameter:
// unlike a plain XPathQueryItem, a parameter cannot itself carry parameters.
type AbstractXPathQueryParameter struct {
	*AbstractXPathQueryItem
}

// newAbstractXPathQueryParameter creates an AbstractXPathQueryParameter.
func newAbstractXPathQueryParameter() *AbstractXPathQueryParameter {
	return &AbstractXPathQueryParameter{AbstractXPathQueryItem: &AbstractXPathQueryItem{}}
}

// AddParameter always panics: a parameter cannot carry a parameter of its own. Ports
// addParameter(XPathQueryParameter), which Java throws UnsupportedOperationException from.
func (p *AbstractXPathQueryParameter) AddParameter(parameter XPathQueryParameter) {
	panic("Unable to set a parameter for XPathQueryParameter!")
}

// Parameters always returns nil: parameters are not supported. Ports getParameters().
func (p *AbstractXPathQueryParameter) Parameters() []XPathQueryParameter {
	return nil
}
