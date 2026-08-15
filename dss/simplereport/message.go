// Ported from dss-jaxb-parsers/src/main/java/eu/europa/esig/dss/jaxb/object/Message.java
// (DSS 6.5.RC1).
//
// dss-jaxb-parsers is outside S8B_BRIEF.md's manifest (it belongs to the
// generic JAXB machinery dss-diagnostic-jaxb's port already collapsed past,
// see dss/diagnostic/diagnostic_data_facade.go's header), but SimpleReport's
// public method surface returns []Message, so this file collapses the one
// class it actually needs directly into this package rather than leaving
// SimpleReport unbuildable - the same cross-chunk assumption
// diagnostic_data_facade.go documents for AbstractJaxbFacade, flagged here
// per S8B_BRIEF.md's "Report: cross-chunk assumptions" instruction.
package simplereport

// Message represents the Message returned in the validation process. Port
// of eu.europa.esig.dss.jaxb.object.Message.
type Message struct {
	// Key is the message key.
	Key string
	// Value is the message text value.
	Value string
}

// NewMessage builds a Message. Port of the Message(String, String)
// constructor.
func NewMessage(key, value string) Message {
	return Message{Key: key, Value: value}
}
