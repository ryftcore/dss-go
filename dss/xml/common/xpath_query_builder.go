// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/XPathQueryBuilder.java (DSS 6.5.RC1).
package common

import "fmt"

// XPathQueryBuilder is a helper for building an XPathQuery.
//
// Naming: Java overloads the static factory all() (0-arg) against the instance method
// all(boolean), and fromCurrentPosition() (0-arg static) against fromCurrentPosition(boolean)
// (instance); Go cannot overload by arity, so the static factories are
// XPathQueryBuilderAll/XPathQueryBuilderFromCurrentPosition/
// XPathQueryBuilderAllFromCurrentPosition/XPathQueryBuilderFromXPathQuery (free functions)
// and the instance setters are SetAll/SetFromCurrentPosition.
//
// attributeValue and idValue are *string rather than string: Java's build() branches on
// "attributeValue != null" to decide between an @attribute predicate
// (XPathQueryAttributeParameter) and an independent @attribute node (XPathQueryAttributeItem)
// - see Attribute/AttributeWithValue below - and a bare Go string cannot represent "not set"
// the way Java null does.
type XPathQueryBuilder struct {
	fromCurrentPosition bool
	all                 bool
	elements            []DSSElement
	attribute           DSSAttribute
	attributeValue      *string
	notChildOf          DSSElement
	idValue             *string
}

// newXPathQueryBuilder creates an XPathQueryBuilder with empty configuration.
func newXPathQueryBuilder() *XPathQueryBuilder {
	return &XPathQueryBuilder{}
}

// XPathQueryBuilderFromCurrentPosition instantiates a builder starting from the current
// position. Ports the static factory fromCurrentPosition().
func XPathQueryBuilderFromCurrentPosition() *XPathQueryBuilder {
	return newXPathQueryBuilder().SetFromCurrentPosition(true)
}

// SetFromCurrentPosition defines whether to start the XPath expression from the current
// position. Ports the instance method fromCurrentPosition(boolean).
func (b *XPathQueryBuilder) SetFromCurrentPosition(fromCurrentPosition bool) *XPathQueryBuilder {
	b.fromCurrentPosition = fromCurrentPosition
	return b
}

// XPathQueryBuilderAll instantiates a builder searching all element occurrences. Ports the
// static factory all().
func XPathQueryBuilderAll() *XPathQueryBuilder {
	return newXPathQueryBuilder().SetAll(true)
}

// SetAll defines whether to search all element occurrences. Ports the instance method
// all(boolean).
func (b *XPathQueryBuilder) SetAll(all bool) *XPathQueryBuilder {
	b.all = all
	return b
}

// XPathQueryBuilderAllFromCurrentPosition instantiates a builder searching all element
// occurrences from the current position. Ports the static factory
// allFromCurrentPosition().
func XPathQueryBuilderAllFromCurrentPosition() *XPathQueryBuilder {
	return newXPathQueryBuilder().SetAll(true).SetFromCurrentPosition(true)
}

// XPathQueryBuilderFromXPathQuery instantiates a builder from the existing xPathQuery. This
// creates a new builder; any setter used on it overrides only that value. Panics if
// xPathQuery is nil (Java Objects.requireNonNull(xPathQuery, "XPathQuery cannot be null!")).
func XPathQueryBuilderFromXPathQuery(xPathQuery XPathQuery) *XPathQueryBuilder {
	if xPathQuery == nil {
		panic("XPathQuery cannot be null!")
	}

	builder := newXPathQueryBuilder().
		SetFromCurrentPosition(xPathQuery.IsFromCurrentPosition()).
		SetAll(xPathQuery.IsAll())

	if xPathQuery.IsEmpty() {
		return builder
	}

	var elementList []DSSElement

	queryItem := xPathQuery.FirstXPathQueryItem()
queryLoop:
	for queryItem != nil {
		switch item := queryItem.(type) {
		case *XPathQueryAnyItem:
			// continue
		case *XPathQueryAttributeItem:
			builder.Attribute(item.Attribute())
		case *XPathQueryElementItem:
			elementList = append(elementList, item.Element())
		case *XPathQueryEndItem:
			break queryLoop
		default:
			panic(fmt.Sprintf("The XPathQueryItem of type '%T' is not supported by the XPathQueryBuilder implementation!", queryItem))
		}

		nextItem := queryItem.NextItem()

		parameters := queryItem.Parameters()
		if len(parameters) > 0 {
			if nextItem == nil {
				for _, parameter := range parameters {
					switch p := parameter.(type) {
					case *XPathQueryNotChildOfParameter:
						builder.NotChildOf(p.ParentElement())
					case *XPathQueryAttributeParameter:
						builder.AttributeWithValue(p.Attribute(), p.AttributeValue())
					case *XPathQueryIdentifierParameter:
						builder.IdValue(p.Id())
					default:
						panic(fmt.Sprintf("The XPathQueryParameter of type '%T' is not supported by the XPathQueryBuilder implementation!", parameter))
					}
				}
			} else {
				panic("The XPathQueryBuilder does not support parameters handling for any item not in the last position. Please build the XPath query chain using other options.")
			}
		}

		queryItem = nextItem
	}

	if len(elementList) > 0 {
		builder.Elements(elementList...)
	}

	return builder
}

// Element defines the element to search. Ports element(DSSElement).
func (b *XPathQueryBuilder) Element(element DSSElement) *XPathQueryBuilder {
	b.elements = []DSSElement{element}
	return b
}

// Elements defines the element path to search. Ports elements(DSSElement...).
func (b *XPathQueryBuilder) Elements(elements ...DSSElement) *XPathQueryBuilder {
	b.elements = elements
	return b
}

// NotChildOf defines that the looked-up element shall not be a parent of notChildOf. Ports
// notChildOf(DSSElement).
func (b *XPathQueryBuilder) NotChildOf(notChildOf DSSElement) *XPathQueryBuilder {
	b.notChildOf = notChildOf
	return b
}

// Attribute defines the attribute to search, with no expected value: the built query
// carries it as an independent @attribute chain item. Ports attribute(DSSAttribute), which
// delegates to attribute(DSSAttribute, null).
func (b *XPathQueryBuilder) Attribute(attribute DSSAttribute) *XPathQueryBuilder {
	b.attribute = attribute
	b.attributeValue = nil
	return b
}

// AttributeWithValue defines the attribute and expected value to search: the built query
// carries it as a predicate on the last defined element. Ports
// attribute(DSSAttribute, String).
func (b *XPathQueryBuilder) AttributeWithValue(attribute DSSAttribute, attributeValue string) *XPathQueryBuilder {
	b.attribute = attribute
	b.attributeValue = &attributeValue
	return b
}

// IdValue defines the id attribute value to search. Ports idValue(String).
func (b *XPathQueryBuilder) IdValue(idValue string) *XPathQueryBuilder {
	b.idValue = &idValue
	return b
}

// Build builds the XPath expression query. Ports build().
func (b *XPathQueryBuilder) Build() XPathQuery {
	var xPathQuery XPathQuery
	switch {
	case b.all && b.fromCurrentPosition:
		xPathQuery = NewAllFromCurrentPositionXPathQuery()
	case b.fromCurrentPosition:
		xPathQuery = NewFromCurrentPositionXPathQuery()
	case b.all:
		xPathQuery = NewAllXPathQuery()
	default:
		panic("Unsupported operation")
	}

	var lastItem XPathQueryItem
	if len(b.elements) > 0 {
		for _, element := range b.elements {
			lastItem = NewXPathQueryElementItem(element)
			xPathQuery.SetNextItem(lastItem)
		}
	} else {
		lastItem = NewXPathQueryAnyItem()
		xPathQuery.SetNextItem(lastItem)
	}

	if b.notChildOf != nil {
		lastItem.AddParameter(NewXPathQueryNotChildOfParameter(b.notChildOf))
	}

	if b.attribute != nil {
		if b.attributeValue != nil {
			lastItem.AddParameter(NewXPathQueryAttributeParameter(b.attribute, *b.attributeValue))
		} else {
			// @attribute is considered as an independent node
			xPathQuery.SetNextItem(NewXPathQueryAttributeItem(b.attribute))
		}
	}

	if b.idValue != nil {
		lastItem.AddParameter(NewXPathQueryIdentifierParameter(*b.idValue))
	}

	xPathQuery.SetNextItem(NewXPathQueryEndItem())

	return xPathQuery
}
