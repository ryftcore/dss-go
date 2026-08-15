// Ported from dss-simple-certificate-report-jaxb/src/main/java/eu/europa/esig/dss/simplecertificatereport/SimpleCertificateReportFacade.java
// (DSS 6.5.RC1).
//
// Java's SimpleCertificateReportFacade extends the generic dss-jaxb-common
// AbstractJaxbFacade<T>, a module outside S8B_BRIEF.md's manifest. As with
// dss/diagnostic/diagnostic_data_facade.go (Phase 8a) and
// dss/simplereport/simple_report_facade.go (Phase 8b), this port collapses
// AbstractJaxbFacade's marshal/unmarshal template method directly into
// SimpleCertificateReportFacade using encoding/xml (the marshal-parity
// contract this phase exists to satisfy; the low-level byte-exact
// Marshal/Unmarshal that the KAT actually exercises lives in jaxb/xml.go,
// which this facade does NOT call). XSD-schema validation and the HTML
// report generation methods are stubbed via
// SimpleCertificateReportXmlDefiner's deferred
// Schema()/HtmlBootstrap4Templates() (see that file's header) and surface
// the same error here.
package simplecertificatereport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/utain/esig/dss/simplecertificatereport/jaxb"
)

// SimpleCertificateReportFacade contains methods to generate a
// SimpleCertificateReport.
type SimpleCertificateReportFacade struct{}

// NewSimpleCertificateReportFacade instantiates a new
// SimpleCertificateReportFacade. Port of newFacade().
func NewSimpleCertificateReportFacade() *SimpleCertificateReportFacade {
	return &SimpleCertificateReportFacade{}
}

// Marshal returns the XML representation of simpleCertificateReport. Port
// of marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *SimpleCertificateReportFacade) Marshal(simpleCertificateReport *jaxb.XmlSimpleCertificateReport) (string, error) {
	if simpleCertificateReport == nil {
		return "", errors.New("JAXBObject is null")
	}
	var buf bytes.Buffer
	if err := f.MarshalToWriter(simpleCertificateReport, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarshalToWriter marshals simpleCertificateReport into w. Port of
// marshall(T, OutputStream).
func (f *SimpleCertificateReportFacade) MarshalToWriter(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
	if simpleCertificateReport == nil {
		return errors.New("JAXBObject is null")
	}
	if w == nil {
		return errors.New("OutputStream is null")
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "    ")
	if err := encoder.Encode(simpleCertificateReport); err != nil {
		return err
	}
	return encoder.Flush()
}

// Unmarshal unmarshals r and returns the XmlSimpleCertificateReport. Port
// of unmarshall(InputStream).
func (f *SimpleCertificateReportFacade) Unmarshal(r io.Reader) (*jaxb.XmlSimpleCertificateReport, error) {
	if r == nil {
		return nil, errors.New("InputStream is null")
	}
	var result jaxb.XmlSimpleCertificateReport
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnmarshalString unmarshals xmlObject and returns the
// XmlSimpleCertificateReport. Port of unmarshall(String).
func (f *SimpleCertificateReportFacade) UnmarshalString(xmlObject string) (*jaxb.XmlSimpleCertificateReport, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}

// GenerateHtmlReport generates a Bootstrap 4 Simple certificate report.
// Port of generateHtmlReport(XmlSimpleCertificateReport); XSLT execution is
// deferred, see SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GenerateHtmlReport(simpleCertificateReport *jaxb.XmlSimpleCertificateReport) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportToWriter generates a Bootstrap 4 Simple certificate
// report into w. Port of generateHtmlReport(XmlSimpleCertificateReport,
// Result); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GenerateHtmlReportToWriter(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Simple
// certificate report from already-marshalled simple-certificate-report
// XML. Port of generateHtmlReport(String); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GenerateHtmlReportFromMarshalled(marshalledSimpleCertificateReport string) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalledToWriter generates a Bootstrap 4 Simple
// certificate report from already-marshalled simple-certificate-report XML
// into w. Port of generateHtmlReport(String, Result); XSLT execution is
// deferred, see SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GenerateHtmlReportFromMarshalledToWriter(marshalledSimpleCertificateReport string, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GeneratePdfReport generates a PDF Simple certificate report into w. Port
// of generatePdfReport(XmlSimpleCertificateReport, Result); XSLT execution
// is deferred, see SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GeneratePdfReport(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Simple certificate report
// from already-marshalled simple-certificate-report XML into w. Port of
// generatePdfReport(String, Result); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *SimpleCertificateReportFacade) GeneratePdfReportFromMarshalled(marshalledSimpleCertificateReport string, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}
