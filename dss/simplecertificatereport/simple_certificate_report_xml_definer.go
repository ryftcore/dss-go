// Ported from dss-simple-certificate-report-jaxb/src/main/java/eu/europa/esig/dss/simplecertificatereport/SimpleCertificateReportXmlDefiner.java
// (DSS 6.5.RC1).
//
// DEFERRED: Java's SimpleCertificateReportXmlDefiner centers on three
// thread-safe singletons built by javax.xml.validation/javax.xml.transform
// machinery that has no Go stdlib equivalent: an XSD Schema (for
// marshaller/unmarshaller validation) and XSLT Templates (for Bootstrap 4
// HTML and PDF report rendering). This port keeps the schema/XSLT
// resource location constants for documentation/testdata purposes, but
// Schema()/HtmlBootstrap4Templates()/PdfTemplates() are stubs returning an
// error. Facade.Marshal/Unmarshal (the half of this pair
// that the marshal-parity KAT actually exercises, via jaxb.Marshal/
// Unmarshal) does not depend on any of the three.
package simplecertificatereport

import "errors"

// SimpleCertificateReportSchemaLocation is the location of the
// SimpleCertificateReport XSD, relative to the embedded/testdata resource
// root. Port of SIMPLE_CERTIFICATE_REPORT_SCHEMA_LOCATION.
const SimpleCertificateReportSchemaLocation = "/xsd/SimpleCertificateReport.xsd"

// SimpleCertificateReportXsltHtmlBootstrap4Location is the location of the
// Bootstrap 4 HTML XSLT template, relative to the embedded/testdata
// resource root. Port of
// SIMPLE_CERTIFICATE_REPORT_XSLT_HTML_BOOTSTRAP4_LOCATION.
const SimpleCertificateReportXsltHtmlBootstrap4Location = "/xslt/html/simple-certificate-report-bootstrap4.xslt"

// SimpleCertificateReportXsltPdfLocation is the location of the PDF
// simple-certificate-report XSLT template, relative to the
// embedded/testdata resource root. Port of
// SIMPLE_CERTIFICATE_REPORT_XSLT_PDF_LOCATION.
const SimpleCertificateReportXsltPdfLocation = "/xslt/pdf/simple-certificate-report.xslt"

// ErrXSDSchemaNotSupported is returned by Schema(): no XSD validator ships
// in the Go stdlib and none has been added to this port (see the file
// header).
var ErrXSDSchemaNotSupported = errors.New("simplecertificatereport: SimpleCertificateReport.xsd schema validation is not implemented in this port (deferred, see SimpleCertificateReportXmlDefiner)")

// ErrHtmlTemplatesNotSupported is returned by HtmlBootstrap4Templates()
// (and Facade's HTML report generators): no XSLT
// engine ships in the Go stdlib and none has been added to this port (see
// the file header).
var ErrHtmlTemplatesNotSupported = errors.New("simplecertificatereport: simple-certificate-report-bootstrap4.xslt HTML rendering is not implemented in this port (deferred, see SimpleCertificateReportXmlDefiner)")

// ErrPdfTemplatesNotSupported is returned by PdfTemplates() (and
// Facade's PDF report generators): no XSLT engine
// ships in the Go stdlib and none has been added to this port (see the
// file header).
var ErrPdfTemplatesNotSupported = errors.New("simplecertificatereport: simple-certificate-report.xslt PDF rendering is not implemented in this port (deferred, see SimpleCertificateReportXmlDefiner)")

// Schema is a stub for the XSD Schema for the SimpleCertificateReport. Port
// of getSchema(); see the DEFERRED note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrXSDSchemaNotSupported
}

// HtmlBootstrap4Templates is a stub for the Bootstrap 4 HTML XSLT template.
// Port of getHtmlBootstrap4Templates(); see the DEFERRED note above.
func HtmlBootstrap4Templates() (struct{}, error) {
	return struct{}{}, ErrHtmlTemplatesNotSupported
}

// PdfTemplates is a stub for the PDF XSLT template. Port of
// getPdfTemplates(); see the DEFERRED note above.
func PdfTemplates() (struct{}, error) {
	return struct{}{}, ErrPdfTemplatesNotSupported
}
