// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/DiagnosticDataXmlDefiner.java (DSS 6.5.RC1).
//
// DEFERRED: Java's DiagnosticDataXmlDefiner centers on three thread-safe singletons built by
// javax.xml.validation/javax.xml.transform machinery that has no Go stdlib equivalent: an XSD
// Schema (for marshaller/unmarshaller validation) and XSLT Templates (for SVG rendering). This
// port keeps the schema location constant and the resource paths for documentation/testdata
// purposes (the XSLT is not executed; nothing in scope consumes it), but
// Schema()/SvgTemplates() are stubs returning an error: no Go stdlib XSD validator or XSLT
// engine exists, and adding a third-party one is outside the stdlib-first dependency policy
// without tech-lead sign-off (PORTING.md "Dependency policy"). DiagnosticDataFacade.Marshal/
// Unmarshal (the half of this pair that the marshal-parity KAT actually exercises) does not
// depend on either stub.
package diagnostic

import "errors"

// DiagnosticDataSchemaLocation is the location of the DiagnosticData XSD, relative to the
// embedded/testdata resource root. Port of DIAGNOSTIC_DATA_SCHEMA_LOCATION.
const DiagnosticDataSchemaLocation = "/xsd/DiagnosticData.xsd"

// DiagnosticDataXsltSvgLocation is the location of the DiagnosticData SVG XSLT template,
// relative to the embedded/testdata resource root. Port of DIAGNOSTIC_DATA_XSLT_SVG_LOCATION.
const DiagnosticDataXsltSvgLocation = "/xslt/svg/diagnostic-data.xslt"

// ErrXSDSchemaNotSupported is returned by Schema(): no XSD validator ships in the Go stdlib and
// none has been added to this port (see the file header).
var ErrXSDSchemaNotSupported = errors.New("diagnostic: DiagnosticData.xsd schema validation is not implemented in this port (deferred, see DiagnosticDataXmlDefiner)")

// ErrSvgTemplatesNotSupported is returned by SvgTemplates(): no XSLT engine ships in the Go
// stdlib and none has been added to this port (see the file header).
var ErrSvgTemplatesNotSupported = errors.New("diagnostic: diagnostic-data.xslt SVG rendering is not implemented in this port (deferred, see DiagnosticDataXmlDefiner)")

// Schema is a stub for the XSD Schema for the DiagnosticData. Port of getSchema(); see the
// DEFERRED note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrXSDSchemaNotSupported
}

// SvgTemplates is a stub for the SVG XSLT template. Port of getSvgTemplates(); see the
// DEFERRED note above.
func SvgTemplates() (struct{}, error) {
	return struct{}{}, ErrSvgTemplatesNotSupported
}
