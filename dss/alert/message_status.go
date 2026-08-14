// Ported from dss-alert/src/main/java/eu/europa/esig/dss/alert/status/MessageStatus.java (DSS 6.5.RC1).
package alert

// MessageStatus carries a message describing an occurred event.
type MessageStatus struct {
	message string
}

// NewMessageStatus creates a MessageStatus with an empty message.
func NewMessageStatus() *MessageStatus {
	return &MessageStatus{}
}

// Message returns the message describing the occurred event.
func (m *MessageStatus) Message() string {
	return m.message
}

// SetMessage sets the message describing the occurred event.
func (m *MessageStatus) SetMessage(message string) {
	m.message = message
}

// RelatedObjectIds is unsupported for MessageStatus; it panics like Java's UnsupportedOperationException.
func (m *MessageStatus) RelatedObjectIds() []string {
	panic("RelatedObjectIds() is not supported for the current implementation!")
}

// IsEmpty reports whether the message is unset or empty.
func (m *MessageStatus) IsEmpty() bool {
	return m.message == ""
}

// ErrorString returns the message itself.
func (m *MessageStatus) ErrorString() string {
	return m.Message()
}

// String returns "Status : Valid" when empty, the error string otherwise.
func (m *MessageStatus) String() string {
	if m.IsEmpty() {
		return "Status : Valid"
	}
	return m.ErrorString()
}
