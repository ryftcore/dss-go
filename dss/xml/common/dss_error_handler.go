// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/DSSErrorHandler.java (DSS 6.5.RC1).
package common

// XMLParseError is a documented stand-in for org.xml.sax.SAXParseException: the (message,
// line, column, public/system ID) tuple XSD schema validation reports a condition with. No
// XSD validation is implemented (SchemaFactory/Validator are placeholder stubs - see
// schema_factory_builder.go's doc comment), so nothing in this package constructs one today;
// the type exists so DSSErrorHandler's public contract is ready for the phase that adds
// real schema validation.
type XMLParseError struct {
	Message      string
	LineNumber   int
	ColumnNumber int
	PublicID     string
	SystemID     string
}

// Error implements the error interface.
func (e *XMLParseError) Error() string {
	return e.Message
}

// DSSErrorHandler is the default error handler used to collect errors occurring during
// validation. Ports the org.xml.sax.ErrorHandler implementation.
//
// Naming: the Java error(SAXParseException)/fatalError(SAXParseException)/
// warning(SAXParseException) methods implement org.xml.sax.ErrorHandler, a real Java
// interface a Validator invokes; Go has no such interface to satisfy (no schema validation
// is implemented, see the type doc above), and a method literally named Error(*XMLParseError)
// would read as - without being - an implementation of Go's error interface. The mutators
// are renamed RecordError/RecordFatalError/RecordWarning to avoid that collision; the
// getters keep their Java names.
type DSSErrorHandler struct {
	errors      []*XMLParseError
	fatalErrors []*XMLParseError
	warnings    []*XMLParseError
}

// NewDSSErrorHandler creates a DSSErrorHandler with empty message lists.
func NewDSSErrorHandler() *DSSErrorHandler {
	return &DSSErrorHandler{}
}

// RecordError records an error occurred during validation. Ports error(SAXParseException).
func (h *DSSErrorHandler) RecordError(e *XMLParseError) {
	h.errors = append(h.errors, e)
}

// RecordFatalError records a fatal error occurred during validation. Ports
// fatalError(SAXParseException).
func (h *DSSErrorHandler) RecordFatalError(e *XMLParseError) {
	h.fatalErrors = append(h.fatalErrors, e)
}

// RecordWarning records a warning occurred during validation. Ports
// warning(SAXParseException).
func (h *DSSErrorHandler) RecordWarning(e *XMLParseError) {
	h.warnings = append(h.warnings, e)
}

// Errors returns the errors occurred during validation. An empty slice means validation
// succeeded. Ports getErrors().
func (h *DSSErrorHandler) Errors() []*XMLParseError {
	return h.errors
}

// FatalErrors returns the fatal errors occurred during validation. Ports getFatalErrors().
func (h *DSSErrorHandler) FatalErrors() []*XMLParseError {
	return h.fatalErrors
}

// Warnings returns the warnings occurred during validation. Ports getWarnings().
func (h *DSSErrorHandler) Warnings() []*XMLParseError {
	return h.warnings
}

// IsValid reports whether validation succeeded (no errors, fatal errors or warnings
// recorded). Ports isValid().
func (h *DSSErrorHandler) IsValid() bool {
	return len(h.errors) == 0 && len(h.fatalErrors) == 0 && len(h.warnings) == 0
}
