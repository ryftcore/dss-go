// Ported from dss-simple-report-jaxb/src/main/java/eu/europa/esig/dss/simplereport/SimpleReportXmlDefiner.java
// (DSS 6.5.RC1).
//
// DEFERRED: Java's SimpleReportXmlDefiner centers on three thread-safe
// singletons built by javax.xml.validation/javax.xml.transform machinery
// that has no Go stdlib equivalent: an XSD Schema (for
// marshaller/unmarshaller validation) and XSLT Templates (for Bootstrap 4
// HTML and PDF report rendering). This port keeps the schema/XSLT
// resource location constants for documentation/testdata purposes, but
// Schema()/HtmlBootstrap4Templates()/PdfTemplates() are stubs returning an
// error: no Go stdlib XSD validator or XSLT engine exists, and adding a
// third-party one is outside the stdlib-first dependency policy without
// tech-lead sign-off (PORTING.md "Dependency policy").
// Facade.Marshal/Unmarshal (the half of this pair that the
// marshal-parity KAT actually exercises, via jaxb.Marshal/Unmarshal) does
// not depend on any of the three.
package simplereport

import "errors"

// SimpleReportSchemaLocation is the location of the SimpleReport XSD,
// relative to the embedded/testdata resource root. Port of
// SIMPLE_REPORT_SCHEMA_LOCATION.
const SimpleReportSchemaLocation = "/xsd/SimpleReport.xsd"

// SimpleReportXsltHtmlBootstrap4Location is the location of the Bootstrap 4
// HTML XSLT template, relative to the embedded/testdata resource root. Port
// of SIMPLE_REPORT_XSLT_HTML_BOOTSTRAP4_LOCATION.
const SimpleReportXsltHtmlBootstrap4Location = "/xslt/html/simple-report-bootstrap4.xslt"

// SimpleReportXsltPdfLocation is the location of the PDF simple-report XSLT
// template, relative to the embedded/testdata resource root. Port of
// SIMPLE_REPORT_XSLT_PDF_LOCATION.
const SimpleReportXsltPdfLocation = "/xslt/pdf/simple-report.xslt"

// ErrXSDSchemaNotSupported is returned by Schema(): no XSD validator ships
// in the Go stdlib and none has been added to this port (see the file
// header).
var ErrXSDSchemaNotSupported = errors.New("simplereport: SimpleReport.xsd schema validation is not implemented in this port (deferred, see SimpleReportXmlDefiner)")

// ErrHtmlTemplatesNotSupported is returned by HtmlBootstrap4Templates()
// (and Facade's HTML report generators): no XSLT engine ships
// in the Go stdlib and none has been added to this port (see the file
// header).
var ErrHtmlTemplatesNotSupported = errors.New("simplereport: simple-report-bootstrap4.xslt HTML rendering is not implemented in this port (deferred, see SimpleReportXmlDefiner)")

// ErrPdfTemplatesNotSupported is returned by PdfTemplates() (and
// Facade's PDF report generators): no XSLT engine ships in the
// Go stdlib and none has been added to this port (see the file header).
var ErrPdfTemplatesNotSupported = errors.New("simplereport: simple-report.xslt PDF rendering is not implemented in this port (deferred, see SimpleReportXmlDefiner)")

// Schema is a stub for the XSD Schema for the SimpleReport. Port of
// getSchema(); see the DEFERRED note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrXSDSchemaNotSupported
}

// HtmlBootstrap4Templates is a stub for the Bootstrap 4 HTML XSLT template.
// Port of getHtmlBootstrap4Templates(); see the DEFERRED note above.
func HtmlBootstrap4Templates() (struct{}, error) {
	return struct{}{}, ErrHtmlTemplatesNotSupported
}

// PdfTemplates is a stub for the PDF XSLT template. Port of getPdfTemplates();
// see the DEFERRED note above.
func PdfTemplates() (struct{}, error) {
	return struct{}{}, ErrPdfTemplatesNotSupported
}
