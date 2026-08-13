// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/AbstractXPathQuery.java (DSS 6.5.RC1).
package common

import "strings"

const (
	xPathQuerySlash        = "/"
	xPathQueryOpenBracket  = "["
	xPathQueryCloseBracket = "]"
)

// AbstractXPathQuery is the shared implementation behind every XPathQuery variant. What
// Java calls the abstract method getXPathPreamble() - overridden by AllXPathQuery,
// FromCurrentPositionXPathQuery and AllFromCurrentPositionXPathQuery to return their fixed
// "//", "./", ".//" strings - has no Go override point (see doc.go); the preamble and the
// isAll/isFromCurrentPosition flags are instead supplied once at construction by
// newAbstractXPathQuery, called from each of those three files.
type AbstractXPathQuery struct {
	firstItem   XPathQueryItem
	currentItem XPathQueryItem

	preamble            string
	all                 bool
	fromCurrentPosition bool
}

// newAbstractXPathQuery creates an AbstractXPathQuery with the given fixed preamble and
// isAll/isFromCurrentPosition flags.
func newAbstractXPathQuery(preamble string, all, fromCurrentPosition bool) *AbstractXPathQuery {
	return &AbstractXPathQuery{preamble: preamble, all: all, fromCurrentPosition: fromCurrentPosition}
}

// FirstXPathQueryItem implements XPathQuery.
func (q *AbstractXPathQuery) FirstXPathQueryItem() XPathQueryItem {
	return q.firstItem
}

// SetNextItem implements XPathQuery.
func (q *AbstractXPathQuery) SetNextItem(nextItem XPathQueryItem) XPathQuery {
	if q.firstItem == nil {
		q.firstItem = nextItem
		q.currentItem = nextItem
	} else {
		q.currentItem = q.currentItem.SetNextItem(nextItem)
	}
	return q
}

// IsEmpty implements XPathQuery.
func (q *AbstractXPathQuery) IsEmpty() bool {
	return q.firstItem == nil
}

// IsAll implements XPathQuery.
func (q *AbstractXPathQuery) IsAll() bool {
	return q.all
}

// IsFromCurrentPosition implements XPathQuery.
func (q *AbstractXPathQuery) IsFromCurrentPosition() bool {
	return q.fromCurrentPosition
}

// QueryString implements XPathQuery.
func (q *AbstractXPathQuery) QueryString() string {
	var sb strings.Builder

	if q.preamble != "" {
		sb.WriteString(q.preamble)
	}

	item := q.firstItem
	for item != nil {
		if item.IsEmpty() {
			item = item.NextItem()
			continue
		}

		if item != q.firstItem {
			sb.WriteString(xPathQuerySlash)
		}
		sb.WriteString(item.QueryString())

		for _, parameter := range item.Parameters() {
			sb.WriteString(xPathQueryOpenBracket)
			sb.WriteString(parameter.QueryString())
			sb.WriteString(xPathQueryCloseBracket)
		}
		item = item.NextItem()
	}

	return sb.String()
}
