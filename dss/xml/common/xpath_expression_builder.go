// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/XPathExpressionBuilder.java (DSS 6.5.RC1).
package common

import "strings"

const (
	xPathExpressionBuilderAllPath                    = "//"
	xPathExpressionBuilderFromCurrentPositionPath    = "./"
	xPathExpressionBuilderAllFromCurrentPositionPath = ".//"
	xPathExpressionBuilderColonPath                  = ":"
	xPathExpressionBuilderSlashPath                  = "/"
	xPathExpressionBuilderAttributePath              = "@"
)

// XPathExpressionBuilder builds an XPath expression string directly (contrast
// XPathQueryBuilder, which builds a queryable XPathQuery chain).
//
// Deprecated: superseded by XPathQueryBuilder (xpath_query_builder.go), exactly as upstream
// marks this class @Deprecated. Ported in full nonetheless, per this phase's "one Go file
// per Java file" requirement; nothing in this codebase constructs one.
//
// Naming: Java overloads the no-argument fromCurrentPosition()/all() against the
// single-boolean-argument fromCurrentPosition(boolean)/all(boolean); Go cannot overload by
// arity, so the no-argument forms are FromCurrentPosition/All and the boolean forms are
// SetFromCurrentPosition/SetAll - the same convention used by XPathQueryBuilder.
type XPathExpressionBuilder struct {
	fromCurrentPosition bool
	all                 bool
	elements            []DSSElement
	attribute           DSSAttribute
	notParentOf         DSSElement
}

// NewXPathExpressionBuilder creates an XPathExpressionBuilder with empty configuration.
func NewXPathExpressionBuilder() *XPathExpressionBuilder {
	return &XPathExpressionBuilder{}
}

// FromCurrentPosition starts the XPath expression from the current position. Ports the
// no-argument fromCurrentPosition().
func (b *XPathExpressionBuilder) FromCurrentPosition() *XPathExpressionBuilder {
	return b.SetFromCurrentPosition(true)
}

// SetFromCurrentPosition defines whether to start the XPath expression from the current
// position. Ports fromCurrentPosition(boolean).
func (b *XPathExpressionBuilder) SetFromCurrentPosition(fromCurrentPosition bool) *XPathExpressionBuilder {
	b.fromCurrentPosition = fromCurrentPosition
	return b
}

// All defines that all element occurrences shall be searched. Ports the no-argument all().
func (b *XPathExpressionBuilder) All() *XPathExpressionBuilder {
	return b.SetAll(true)
}

// SetAll defines whether to search all element occurrences. Ports all(boolean).
func (b *XPathExpressionBuilder) SetAll(all bool) *XPathExpressionBuilder {
	b.all = all
	return b
}

// Element defines the element to search. Ports element(DSSElement).
func (b *XPathExpressionBuilder) Element(element DSSElement) *XPathExpressionBuilder {
	b.elements = []DSSElement{element}
	return b
}

// Elements defines the element path to search. Ports elements(DSSElement[]).
func (b *XPathExpressionBuilder) Elements(elements ...DSSElement) *XPathExpressionBuilder {
	b.elements = elements
	return b
}

// NotParentOf defines that the looked-up element shall not be a parent of notParentOf.
// Ports notParentOf(DSSElement).
func (b *XPathExpressionBuilder) NotParentOf(notParentOf DSSElement) *XPathExpressionBuilder {
	b.notParentOf = notParentOf
	return b
}

// Attribute defines the attribute to search. Ports attribute(DSSAttribute).
func (b *XPathExpressionBuilder) Attribute(attribute DSSAttribute) *XPathExpressionBuilder {
	b.attribute = attribute
	return b
}

// Build builds the XPath expression string. Panics if neither All nor
// SetFromCurrentPosition/FromCurrentPosition was set (Java throws
// UnsupportedOperationException("Unsupported operation")). Ports build().
func (b *XPathExpressionBuilder) Build() string {
	var sb strings.Builder

	switch {
	case b.all && b.fromCurrentPosition:
		sb.WriteString(xPathExpressionBuilderAllFromCurrentPositionPath)
	case b.fromCurrentPosition:
		sb.WriteString(xPathExpressionBuilderFromCurrentPositionPath)
	case b.all:
		sb.WriteString(xPathExpressionBuilderAllPath)
	default:
		panic("Unsupported operation")
	}

	nbElements := len(b.elements)
	for i, element := range b.elements {
		sb.WriteString(xPathExpressionBuilderElementString(element))
		if i < nbElements-1 {
			sb.WriteString(xPathExpressionBuilderSlashPath)
		}
	}

	if b.notParentOf != nil {
		sb.WriteString(xPathExpressionBuilderNotParent(b.notParentOf))
	}

	if b.attribute != nil {
		sb.WriteString(xPathExpressionBuilderSlashPath)
		sb.WriteString(xPathExpressionBuilderAttributeString(b.attribute))
	}

	return sb.String()
}

// xPathExpressionBuilderElementString ports the private getElement(DSSElement).
func xPathExpressionBuilderElementString(element DSSElement) string {
	var sb strings.Builder
	if namespace := element.Namespace(); namespace != nil {
		sb.WriteString(namespace.Prefix())
		sb.WriteString(xPathExpressionBuilderColonPath)
	}
	sb.WriteString(element.TagName())
	return sb.String()
}

// xPathExpressionBuilderNotParent ports the private getNotParent(DSSElement), whose Java
// doc-comment worked example is "//ds:Signature[not(parent::xades:CounterSignature)]".
func xPathExpressionBuilderNotParent(currentNotParentOf DSSElement) string {
	return "[not(parent::" + xPathExpressionBuilderElementString(currentNotParentOf) + ")]"
}

// xPathExpressionBuilderAttributeString ports the private getAttribute(DSSAttribute).
func xPathExpressionBuilderAttributeString(currentAttribute DSSAttribute) string {
	return xPathExpressionBuilderAttributePath + currentAttribute.AttributeName()
}
