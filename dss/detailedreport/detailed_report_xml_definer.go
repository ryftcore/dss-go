// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReportXmlDefiner.java
// (DSS 6.5.RC1).
//
// DEFERRED: Java's DetailedReportXmlDefiner centers on three thread-safe
// singletons built by javax.xml.validation/javax.xml.transform machinery that
// has no Go stdlib equivalent: an XSD Schema (for marshaller/unmarshaller
// validation) and two sets of XSLT Templates (Bootstrap 4 HTML and PDF report
// rendering). This port keeps the schema/XSLT resource-location constants for
// documentation/testdata purposes, but Schema()/HtmlBootstrap4Templates()/
// PdfTemplates() are stubs returning an error: no Go stdlib XSD validator or
// XSLT engine exists, and adding a third-party one is outside the
// stdlib-first dependency policy without tech-lead sign-off (PORTING.md
// "Dependency policy"). Facade.Marshal/Unmarshal (the half of
// this pair the marshal-parity KAT actually exercises) does not depend on any
// of the three.
package detailedreport

import "errors"

// DetailedReportSchemaLocation is the location of the DetailedReport XSD,
// relative to the embedded/testdata resource root. Port of
// DETAILED_REPORT_SCHEMA_LOCATION.
const DetailedReportSchemaLocation = "/xsd/DetailedReport.xsd"

// DetailedReportXsltHtmlBootstrap4Location is the location of the Bootstrap 4
// HTML XSLT template, relative to the embedded/testdata resource root. Port
// of DETAILED_REPORT_XSLT_HTML_BOOTSTRAP4_LOCATION.
const DetailedReportXsltHtmlBootstrap4Location = "/xslt/html/detailed-report-bootstrap4.xslt"

// DetailedReportXsltPdfLocation is the location of the PDF report XSLT
// template, relative to the embedded/testdata resource root. Port of
// DETAILED_REPORT_XSLT_PDF_LOCATION.
const DetailedReportXsltPdfLocation = "/xslt/pdf/detailed-report.xslt"

// ErrDetailedReportSchemaNotSupported is returned by Schema(): no XSD
// validator ships in the Go stdlib and none has been added to this port (see
// the file header).
var ErrDetailedReportSchemaNotSupported = errors.New("detailedreport: DetailedReport.xsd schema validation is not implemented in this port (deferred, see DetailedReportXmlDefiner)")

// ErrHtmlBootstrap4TemplatesNotSupported is returned by
// HtmlBootstrap4Templates(): no XSLT engine ships in the Go stdlib and none
// has been added to this port (see the file header).
var ErrHtmlBootstrap4TemplatesNotSupported = errors.New("detailedreport: detailed-report-bootstrap4.xslt HTML rendering is not implemented in this port (deferred, see DetailedReportXmlDefiner)")

// ErrPdfTemplatesNotSupported is returned by PdfTemplates(): no XSLT engine
// ships in the Go stdlib and none has been added to this port (see the file
// header).
var ErrPdfTemplatesNotSupported = errors.New("detailedreport: detailed-report.xslt PDF rendering is not implemented in this port (deferred, see DetailedReportXmlDefiner)")

// Schema is a stub for the XSD Schema for the DetailedReport. Port of
// getSchema(); see the DEFERRED note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrDetailedReportSchemaNotSupported
}

// HtmlBootstrap4Templates is a stub for the HTML Bootstrap 4 XSLT template.
// Port of getHtmlBootstrap4Templates(); see the DEFERRED note above.
func HtmlBootstrap4Templates() (struct{}, error) {
	return struct{}{}, ErrHtmlBootstrap4TemplatesNotSupported
}

// PdfTemplates is a stub for the PDF XSLT template. Port of getPdfTemplates();
// see the DEFERRED note above.
func PdfTemplates() (struct{}, error) {
	return struct{}{}, ErrPdfTemplatesNotSupported
}
