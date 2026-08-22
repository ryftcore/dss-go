// Ported from dss-enumerations/.../MessageType.java (DSS 6.5.RC1).
package enumerations

// MessageType defines possible levels for messages returned by the
// validation process. Implements UriBasedEnum.
type MessageType string

const (
	// MessageTypeError indicates a reason for validation process failure.
	MessageTypeError MessageType = "ERROR"
	// MessageTypeWarn indicates a reason for an issue occurred during the
	// validation, not blocking the process.
	MessageTypeWarn MessageType = "WARN"
	// MessageTypeInfo is an additional informational message returned by
	// the validation process.
	MessageTypeInfo MessageType = "INFO"
)

// messageTypeURIs holds the VR URI of the constraint for each constant.
var messageTypeURIs = map[MessageType]string{
	MessageTypeError: "urn:cef:dss:message:error",
	MessageTypeWarn:  "urn:cef:dss:message:warning",
	MessageTypeInfo:  "urn:cef:dss:message:information",
}

// MessageTypeValues returns all constants in declaration order.
func MessageTypeValues() []MessageType {
	return []MessageType{
		MessageTypeError,
		MessageTypeWarn,
		MessageTypeInfo,
	}
}

// URI returns the VR URI of the constraint. Implements UriBasedEnum.
func (m MessageType) URI() string {
	return messageTypeURIs[m]
}
