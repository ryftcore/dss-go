// Ported from dss-enumerations/.../MessageType.java (DSS 6.5.RC1).
package enumerations

// MessageType defines possible levels for messages returned by the
// validation process. Implements UriBasedEnum.
type MessageType string

const (
	// MessageType_ERROR indicates a reason for validation process failure.
	MessageType_ERROR MessageType = "ERROR"
	// MessageType_WARN indicates a reason for an issue occurred during the
	// validation, not blocking the process.
	MessageType_WARN MessageType = "WARN"
	// MessageType_INFO is an additional informational message returned by
	// the validation process.
	MessageType_INFO MessageType = "INFO"
)

// messageTypeURIs holds the VR URI of the constraint for each constant.
var messageTypeURIs = map[MessageType]string{
	MessageType_ERROR: "urn:cef:dss:message:error",
	MessageType_WARN:  "urn:cef:dss:message:warning",
	MessageType_INFO:  "urn:cef:dss:message:information",
}

// MessageTypeValues returns all constants in declaration order.
func MessageTypeValues() []MessageType {
	return []MessageType{
		MessageType_ERROR,
		MessageType_WARN,
		MessageType_INFO,
	}
}

// URI returns the VR URI of the constraint. Implements UriBasedEnum.
func (m MessageType) URI() string {
	return messageTypeURIs[m]
}
