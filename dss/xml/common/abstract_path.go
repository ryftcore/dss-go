// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/AbstractPath.java (DSS 6.5.RC1).
package common

// This file contains a set of common XML xpath builders. Upstream is an abstract class
// meant to be extended by definition-specific Path classes (e.g. XMLDSigPath, see
// xmldsig_path.go); since these are pure static helper methods with no per-subclass state,
// Go ports them as free package functions rather than as a base type other files embed.
//
// Naming: Java overloads all/fromCurrentPosition by both a single DSSElement and a
// DSSElement... varargs; a single-argument call and XPathQueryBuilder.element(x) versus
// .elements(new DSSElement[]{x}) build identical chains, so Go's variadic parameter merges
// the two overloads below without any behavioural difference. The (DSSElement, DSSAttribute)
// overload of fromCurrentPosition, which Go cannot also name FromCurrentPosition without a
// second overload Go doesn't support, is FromCurrentPositionAttribute.

// All builds the XPath expression to return entries of the given element(s), searching the
// whole document. Ports both all(DSSElement) and all(DSSElement...).
func All(elements ...DSSElement) XPathQuery {
	return XPathQueryBuilderAll().Elements(elements...).Build()
}

// FromCurrentPosition builds the XPath expression to return entries of the given element(s)
// from the current position. Ports both fromCurrentPosition(DSSElement) and
// fromCurrentPosition(DSSElement...).
func FromCurrentPosition(elements ...DSSElement) XPathQuery {
	return XPathQueryBuilderFromCurrentPosition().Elements(elements...).Build()
}

// AllFromCurrentPosition builds the XPath expression to return entries of element from the
// current position, searching all descendants. Ports allFromCurrentPosition(DSSElement).
func AllFromCurrentPosition(element DSSElement) XPathQuery {
	return XPathQueryBuilderAllFromCurrentPosition().Element(element).Build()
}

// AllNotParent builds the XPath expression to return entries of element which are not a
// child of notChildOf. Ports allNotParent(DSSElement, DSSElement).
func AllNotParent(element, notChildOf DSSElement) XPathQuery {
	return XPathQueryBuilderAll().Element(element).NotChildOf(notChildOf).Build()
}

// FromCurrentPositionAttribute builds the XPath expression to return entries starting from
// the current position with the given attribute. Ports the (DSSElement, DSSAttribute)
// overload of fromCurrentPosition (renamed - see the file doc comment).
func FromCurrentPositionAttribute(element DSSElement, attribute DSSAttribute) XPathQuery {
	return XPathQueryBuilderFromCurrentPosition().Element(element).Attribute(attribute).Build()
}
