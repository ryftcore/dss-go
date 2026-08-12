// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/status/ObjectStatus.java (DSS 6.5.RC1).
package alert

import (
	"sort"
	"strings"
)

// ObjectStatus is a Status implementation holding a list of concerned objects'
// identifiers and their corresponding error messages.
//
// Judgment call: upstream keys the related objects in a java.util.HashMap, whose
// iteration order is unspecified. objectMapToString here sorts keys so String()/
// ErrorString() are deterministic; this is an intentional, documented deviation
// from (undefined) Java map ordering, not a behavior change in observable content.
type ObjectStatus struct {
	*MessageStatus
	relatedObjectMap map[string]string
}

// NewObjectStatus creates an ObjectStatus with an empty objects map.
func NewObjectStatus() *ObjectStatus {
	return &ObjectStatus{
		MessageStatus:    NewMessageStatus(),
		relatedObjectMap: make(map[string]string),
	}
}

// AddRelatedObjectIdentifierAndErrorMessage adds a concerned object identifier and
// information about the occurred event.
func (o *ObjectStatus) AddRelatedObjectIdentifierAndErrorMessage(objectId, errorMessage string) {
	o.relatedObjectMap[objectId] = errorMessage
}

// MessageForObjectWithId returns the error message for the object with the given id.
func (o *ObjectStatus) MessageForObjectWithId(objectId string) string {
	return o.relatedObjectMap[objectId]
}

// RelatedObjectIds returns the identifiers of the objects held in the related objects map.
func (o *ObjectStatus) RelatedObjectIds() []string {
	ids := make([]string, 0, len(o.relatedObjectMap))
	for id := range o.relatedObjectMap {
		ids = append(ids, id)
	}
	return ids
}

// IsEmpty reports whether both the message and the related objects map are empty.
func (o *ObjectStatus) IsEmpty() bool {
	return o.MessageStatus.IsEmpty() && len(o.relatedObjectMap) == 0
}

// ErrorString returns the message followed by the formatted related-objects list.
//
// Known deviation: when no message has been set, upstream concatenates a Java null and
// produces "null [id: msg]", whereas this returns " [id: msg]". Unlike the same quirk in
// SignatureAlgorithm.URIBasedOnOID — whose output is embedded in emitted signatures and so
// is reproduced verbatim — this string is only ever surfaced in logs and alert exception
// text, where the literal "null" is a display bug with no consumer parsing it. Left
// deviating deliberately; revisit if any caller ever compares this string.
func (o *ObjectStatus) ErrorString() string {
	return o.Message() + " " + o.objectMapToString()
}

// objectMapToString returns a string listing the occurred errors for each concerned object,
// formatted as "[id1: msg1; id2: msg2]".
func (o *ObjectStatus) objectMapToString() string {
	var sb strings.Builder
	if len(o.relatedObjectMap) > 0 {
		sb.WriteString("[")
		keys := make([]string, 0, len(o.relatedObjectMap))
		for k := range o.relatedObjectMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(o.relatedObjectMap[k])
			if i < len(keys)-1 {
				sb.WriteString("; ")
			}
		}
		sb.WriteString("]")
	}
	return sb.String()
}

// String returns "Status : Valid" when empty, the error string otherwise.
func (o *ObjectStatus) String() string {
	if o.IsEmpty() {
		return "Status : Valid"
	}
	return o.ErrorString()
}
