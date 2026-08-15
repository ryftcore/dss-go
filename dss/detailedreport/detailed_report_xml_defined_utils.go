// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReportXmlDefinedUtils.java
// (DSS 6.5.RC1).
//
// Java's DetailedReportXmlDefinedUtils is a thin holder for a
// dss-xml-common SchemaFactoryBuilder/TransformerFactoryBuilder pair, neither
// of which has a Go stdlib equivalent (see DetailedReportXmlDefiner's header
// for the same DEFERRED rationale dss/diagnostic's DiagnosticDataXmlDefiner
// documents). Kept as a type for API shape parity; both accessors return the
// same stub error DetailedReportXmlDefiner.Schema()/loadTemplates would.
package detailedreport

import "errors"

// ErrSecureTransformerFactoryNotSupported is returned by
// SecureTransformerFactory(): no XSLT engine ships in the Go stdlib and none
// has been added to this port.
var ErrSecureTransformerFactoryNotSupported = errors.New("detailedreport: secure TransformerFactory is not implemented in this port (deferred, see DetailedReportXmlDefinedUtils)")

// ErrSecureSchemaFactoryNotSupported is returned by SecureSchemaFactory(): no
// XSD validator ships in the Go stdlib and none has been added to this port.
var ErrSecureSchemaFactoryNotSupported = errors.New("detailedreport: secure SchemaFactory is not implemented in this port (deferred, see DetailedReportXmlDefinedUtils)")

// DetailedReportXmlDefinedUtils provides access to XML Securities
// configuration required for processing and building of the DSS XML Detailed
// Report.
type DetailedReportXmlDefinedUtils struct{}

// detailedReportXmlDefinedUtilsSingleton is the package-level singleton. Port
// of the private static `singleton` field.
var detailedReportXmlDefinedUtilsSingleton = &DetailedReportXmlDefinedUtils{}

// DetailedReportXmlDefinedUtilsInstance returns the DetailedReportXmlDefinedUtils
// singleton. Port of getInstance().
func DetailedReportXmlDefinedUtilsInstance() *DetailedReportXmlDefinedUtils {
	return detailedReportXmlDefinedUtilsSingleton
}

// SecureTransformerFactory is a stub for a TransformerFactory with enabled
// security features (disabled external DTD/XSD + secure processing). Port of
// getSecureTransformerFactory(); see the DEFERRED note above.
func (u *DetailedReportXmlDefinedUtils) SecureTransformerFactory() (struct{}, error) {
	return struct{}{}, ErrSecureTransformerFactoryNotSupported
}

// SecureSchemaFactory is a stub for a SchemaFactory with enabled security
// features (disabled external DTD/XSD + secure processing). Port of
// getSecureSchemaFactory(); see the DEFERRED note above.
func (u *DetailedReportXmlDefinedUtils) SecureSchemaFactory() (struct{}, error) {
	return struct{}{}, ErrSecureSchemaFactoryNotSupported
}
