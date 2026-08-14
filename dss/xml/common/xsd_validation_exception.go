// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/exception/XSDValidationException.java (DSS 6.5.RC1).
package common

import "strings"

// XSDValidationException is raised for XSD validation error(s). Ports the unchecked
// RuntimeException; callers match it with errors.As per PORTING.md.
//
// Java's addSuppressed(exception) chain (one suppressed exception per collected
// SAXParseException) has no port: Go has no suppressed-exception mechanism, and nothing in
// this codebase currently constructs the SAXParseException-equivalent XMLParseError that
// would have fed it (see dss_error_handler.go) - a documented, accepted gap.
type XSDValidationException struct {
	// Messages are the individual XSD validation error messages.
	Messages []string
}

// NewXSDValidationException creates an XSDValidationException. Ports
// XSDValidationException(List<String>).
func NewXSDValidationException(messages []string) *XSDValidationException {
	return &XSDValidationException{Messages: messages}
}

// AllMessages returns the XSD validation error messages. Ports getAllMessages(), whose
// Collections.emptyList() fallback for a nil list is simply Messages itself (nil is
// equivalent to an empty slice for every use of AllMessages).
func (e *XSDValidationException) AllMessages() []string {
	return e.Messages
}

// Error implements the error interface, joining every message with "; ". Ports
// getMessage(), whose Java falls back to returning null when Messages is empty; the closest
// Go equivalent for an error's message is "".
func (e *XSDValidationException) Error() string {
	if len(e.Messages) == 0 {
		return ""
	}
	return strings.Join(e.Messages, "; ")
}
