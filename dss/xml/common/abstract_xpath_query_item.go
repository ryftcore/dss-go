// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/AbstractXPathQueryItem.java (DSS 6.5.RC1).
package common

import "github.com/utain/esig/dss/internal/xmldom"

// AbstractXPathQueryItem is the shared implementation behind every non-parameter
// XPathQueryItem: chain linkage and parameter storage. What Java calls process(Node) - the
// abstract method matchNode(Node) delegates to after checking parameters - has no Go
// override point (see doc.go); each concrete item instead implements MatchNode itself,
// calling matchParameters after its own test passes.
type AbstractXPathQueryItem struct {
	nextItem   XPathQueryItem
	parameters []XPathQueryParameter
}

// NextItem returns the next XPath chain item, if any. Ports nextItem().
func (i *AbstractXPathQueryItem) NextItem() XPathQueryItem {
	return i.nextItem
}

// SetNextItem sets the next chain item and returns it. Ports setNextItem(XPathQueryItem).
func (i *AbstractXPathQueryItem) SetNextItem(nextItem XPathQueryItem) XPathQueryItem {
	i.nextItem = nextItem
	return nextItem
}

// AddParameter adds a parameter to this item. Ports addParameter(XPathQueryParameter).
func (i *AbstractXPathQueryItem) AddParameter(parameter XPathQueryParameter) {
	i.parameters = append(i.parameters, parameter)
}

// Parameters returns the parameters attached to this item. Ports getParameters().
func (i *AbstractXPathQueryItem) Parameters() []XPathQueryParameter {
	return i.parameters
}

// IsEmpty always reports false for a concrete chain item. Ports isEmpty().
func (i *AbstractXPathQueryItem) IsEmpty() bool {
	return false
}

// matchParameters checks every configured parameter against node. Concrete items call this
// after their own process(node) test passes, together implementing what was one polymorphic
// Java method (matchNode(Node), see the type doc comment).
func (i *AbstractXPathQueryItem) matchParameters(node *xmldom.Node) bool {
	for _, parameter := range i.parameters {
		if !parameter.MatchNode(node) {
			return false
		}
	}
	return true
}

// localName returns node's local name. Ports the protected getLocalName(Node) helper; its
// Java fallback to getNodeName() when getLocalName() is null has no Go counterpart, since
// xmldom.Node.Name.Local is always populated (never null).
func localName(node *xmldom.Node) string {
	return node.Name.Local
}
