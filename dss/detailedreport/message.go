// Ported from dss-jaxb-parsers/src/main/java/eu/europa/esig/dss/jaxb/object/Message.java
// (DSS 6.5.RC1).
//
// dss-jaxb-parsers has not been ported as its own package, but
// DetailedReportMessageCollector returns this DTO from every one of its
// public methods, so it has no working port without it. Message is small,
// self-contained and carries no dependency of its own, so it is implemented
// directly here since it has no other consumer in this package; a dedicated
// jaxbobject package could hold it once dss-jaxb-parsers itself is ported.
package detailedreport

// Message represents the message returned in the validation process. Value
// semantics (not a pointer) so that equality mirrors Java's equals()/hashCode()
// override, which addMessages's dedup relies on the way Java's List.contains
// does.
type Message struct {
	// Key represents the message key.
	Key string
	// Value represents the message text value.
	Value string
}

// NewMessage is the default constructor.
func NewMessage(key, value string) Message {
	return Message{Key: key, Value: value}
}
