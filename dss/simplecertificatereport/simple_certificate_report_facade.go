// Ported from dss-simple-certificate-report-jaxb/src/main/java/eu/europa/esig/dss/simplecertificatereport/SimpleCertificateReportFacade.java
// (DSS 6.5.RC1).
//
// Java's SimpleCertificateReportFacade extends dss-jaxb-common's AbstractJaxbFacade<T>; this
// port implements Marshal/Unmarshal directly with encoding/xml (the low-level byte-exact
// Marshal/Unmarshal that the KAT actually exercises lives in jaxb/xml.go, which this facade
// does not call). XSD-schema validation and HTML report generation are not implemented - see
// SimpleCertificateReportXmlDefiner's deferred Schema()/HtmlBootstrap4Templates() stubs.
package simplecertificatereport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
)

// Facade contains methods to generate a
// SimpleCertificateReport.
type Facade struct{}

// NewFacade instantiates a new
// SimpleCertificateReportFacade. Port of newFacade().
func NewFacade() *Facade {
	return &Facade{}
}

// Marshal returns the XML representation of simpleCertificateReport. Port
// of marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *Facade) Marshal(simpleCertificateReport *jaxb.XmlSimpleCertificateReport) (string, error) {
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
func (f *Facade) MarshalToWriter(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
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
func (f *Facade) Unmarshal(r io.Reader) (*jaxb.XmlSimpleCertificateReport, error) {
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
func (f *Facade) UnmarshalString(xmlObject string) (*jaxb.XmlSimpleCertificateReport, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}

// GenerateHtmlReport generates a Bootstrap 4 Simple certificate report.
// Port of generateHtmlReport(XmlSimpleCertificateReport); XSLT execution is
// deferred, see SimpleCertificateReportXmlDefiner.
func (f *Facade) GenerateHtmlReport(simpleCertificateReport *jaxb.XmlSimpleCertificateReport) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportToWriter generates a Bootstrap 4 Simple certificate
// report into w. Port of generateHtmlReport(XmlSimpleCertificateReport,
// Result); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *Facade) GenerateHtmlReportToWriter(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Simple
// certificate report from already-marshalled simple-certificate-report
// XML. Port of generateHtmlReport(String); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *Facade) GenerateHtmlReportFromMarshalled(marshalledSimpleCertificateReport string) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalledToWriter generates a Bootstrap 4 Simple
// certificate report from already-marshalled simple-certificate-report XML
// into w. Port of generateHtmlReport(String, Result); XSLT execution is
// deferred, see SimpleCertificateReportXmlDefiner.
func (f *Facade) GenerateHtmlReportFromMarshalledToWriter(marshalledSimpleCertificateReport string, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GeneratePdfReport generates a PDF Simple certificate report into w. Port
// of generatePdfReport(XmlSimpleCertificateReport, Result); XSLT execution
// is deferred, see SimpleCertificateReportXmlDefiner.
func (f *Facade) GeneratePdfReport(simpleCertificateReport *jaxb.XmlSimpleCertificateReport, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Simple certificate report
// from already-marshalled simple-certificate-report XML into w. Port of
// generatePdfReport(String, Result); XSLT execution is deferred, see
// SimpleCertificateReportXmlDefiner.
func (f *Facade) GeneratePdfReportFromMarshalled(marshalledSimpleCertificateReport string, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}
