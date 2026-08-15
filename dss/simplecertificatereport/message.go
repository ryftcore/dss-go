// Ported from dss-jaxb-parsers/src/main/java/eu/europa/esig/dss/jaxb/object/Message.java
// (DSS 6.5.RC1).
//
// dss-jaxb-parsers is outside S8B_BRIEF.md's manifest; see
// dss/simplereport/message.go, which collapses the same class for the same
// reason (SimpleCertificateReport's public method surface returns
// []Message too).
package simplecertificatereport

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
