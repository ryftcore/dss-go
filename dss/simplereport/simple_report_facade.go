// Ported from dss-simple-report-jaxb/src/main/java/eu/europa/esig/dss/simplereport/SimpleReportFacade.java
// (DSS 6.5.RC1).
//
// Java's SimpleReportFacade extends dss-jaxb-common's AbstractJaxbFacade<T>; this port
// implements Marshal/Unmarshal directly with encoding/xml (the low-level byte-exact
// Marshal/Unmarshal that the KAT actually exercises lives in jaxb/xml.go, which this facade
// does not call). XSD-schema validation and HTML Bootstrap4/PDF report generation are not
// implemented - see SimpleReportXmlDefiner's deferred
// Schema()/HtmlBootstrap4Templates()/PdfTemplates() stubs.
package simplereport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/simplereport/jaxb"
)

// Facade contains methods to generate a SimpleReport.
type Facade struct{}

// NewSimpleReportFacade instantiates a new SimpleReportFacade. Port of
// newFacade().
func NewFacade() *Facade {
	return &Facade{}
}

// Marshal returns the XML representation of simpleReport. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *Facade) Marshal(simpleReport *jaxb.XmlSimpleReport) (string, error) {
	if simpleReport == nil {
		return "", errors.New("JAXBObject is null")
	}
	var buf bytes.Buffer
	if err := f.MarshalToWriter(simpleReport, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarshalToWriter marshals simpleReport into w. Port of marshall(T,
// OutputStream).
func (f *Facade) MarshalToWriter(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
	if simpleReport == nil {
		return errors.New("JAXBObject is null")
	}
	if w == nil {
		return errors.New("OutputStream is null")
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "    ")
	if err := encoder.Encode(simpleReport); err != nil {
		return err
	}
	return encoder.Flush()
}

// Unmarshal unmarshals r and returns the XmlSimpleReport. Port of
// unmarshall(InputStream).
func (f *Facade) Unmarshal(r io.Reader) (*jaxb.XmlSimpleReport, error) {
	if r == nil {
		return nil, errors.New("InputStream is null")
	}
	var result jaxb.XmlSimpleReport
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnmarshalString unmarshals xmlObject and returns the XmlSimpleReport.
// Port of unmarshall(String).
func (f *Facade) UnmarshalString(xmlObject string) (*jaxb.XmlSimpleReport, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}

// GenerateHtmlReport generates a Bootstrap 4 Simple report. Port of
// generateHtmlReport(XmlSimpleReport); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *Facade) GenerateHtmlReport(simpleReport *jaxb.XmlSimpleReport) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportToWriter generates a Bootstrap 4 Simple report into w.
// Port of generateHtmlReport(XmlSimpleReport, Result); XSLT execution is
// deferred, see SimpleReportXmlDefiner.
func (f *Facade) GenerateHtmlReportToWriter(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Simple report
// from already-marshalled simple-report XML. Port of
// generateHtmlReport(String); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *Facade) GenerateHtmlReportFromMarshalled(marshalledSimpleReport string) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalledToWriter generates a Bootstrap 4 Simple
// report from already-marshalled simple-report XML into w. Port of
// generateHtmlReport(String, Result); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *Facade) GenerateHtmlReportFromMarshalledToWriter(marshalledSimpleReport string, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GeneratePdfReport generates a PDF Simple report. Port of
// generatePdfReport(XmlSimpleReport, Result); XSLT execution is deferred,
// see SimpleReportXmlDefiner.
func (f *Facade) GeneratePdfReport(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Simple report from
// already-marshalled simple-report XML into w. Port of
// generatePdfReport(String, Result); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *Facade) GeneratePdfReportFromMarshalled(marshalledSimpleReport string, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}
