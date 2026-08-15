// Ported from dss-simple-report-jaxb/src/main/java/eu/europa/esig/dss/simplereport/SimpleReportFacade.java
// (DSS 6.5.RC1).
//
// Java's SimpleReportFacade extends the generic dss-jaxb-common
// AbstractJaxbFacade<T>, a module outside S8B_BRIEF.md's manifest. As with
// dss/diagnostic/diagnostic_data_facade.go (Phase 8a), this port collapses
// AbstractJaxbFacade's marshal/unmarshal template method directly into
// SimpleReportFacade using encoding/xml (the marshal-parity contract this
// phase exists to satisfy; the low-level byte-exact Marshal/Unmarshal that
// the KAT actually exercises lives in jaxb/xml.go, which this facade does
// NOT call - matching dss/diagnostic's precedent that the hand facade's own
// marshalling is a convenience, not the parity-pinned path). XSD-schema
// validation (the `validate` booleans on every Java overload) and the HTML
// Bootstrap4/PDF report generation methods are stubbed via
// SimpleReportXmlDefiner's deferred Schema()/HtmlBootstrap4Templates()/
// PdfTemplates() (see that file's header) and surface the same error here.
package simplereport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/utain/esig/dss/simplereport/jaxb"
)

// SimpleReportFacade contains methods to generate a SimpleReport.
type SimpleReportFacade struct{}

// NewSimpleReportFacade instantiates a new SimpleReportFacade. Port of
// newFacade().
func NewSimpleReportFacade() *SimpleReportFacade {
	return &SimpleReportFacade{}
}

// Marshal returns the XML representation of simpleReport. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *SimpleReportFacade) Marshal(simpleReport *jaxb.XmlSimpleReport) (string, error) {
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
func (f *SimpleReportFacade) MarshalToWriter(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
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
func (f *SimpleReportFacade) Unmarshal(r io.Reader) (*jaxb.XmlSimpleReport, error) {
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
func (f *SimpleReportFacade) UnmarshalString(xmlObject string) (*jaxb.XmlSimpleReport, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}

// GenerateHtmlReport generates a Bootstrap 4 Simple report. Port of
// generateHtmlReport(XmlSimpleReport); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GenerateHtmlReport(simpleReport *jaxb.XmlSimpleReport) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportToWriter generates a Bootstrap 4 Simple report into w.
// Port of generateHtmlReport(XmlSimpleReport, Result); XSLT execution is
// deferred, see SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GenerateHtmlReportToWriter(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Simple report
// from already-marshalled simple-report XML. Port of
// generateHtmlReport(String); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GenerateHtmlReportFromMarshalled(marshalledSimpleReport string) (string, error) {
	return "", ErrHtmlTemplatesNotSupported
}

// GenerateHtmlReportFromMarshalledToWriter generates a Bootstrap 4 Simple
// report from already-marshalled simple-report XML into w. Port of
// generateHtmlReport(String, Result); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GenerateHtmlReportFromMarshalledToWriter(marshalledSimpleReport string, w io.Writer) error {
	return ErrHtmlTemplatesNotSupported
}

// GeneratePdfReport generates a PDF Simple report. Port of
// generatePdfReport(XmlSimpleReport, Result); XSLT execution is deferred,
// see SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GeneratePdfReport(simpleReport *jaxb.XmlSimpleReport, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Simple report from
// already-marshalled simple-report XML into w. Port of
// generatePdfReport(String, Result); XSLT execution is deferred, see
// SimpleReportXmlDefiner.
func (f *SimpleReportFacade) GeneratePdfReportFromMarshalled(marshalledSimpleReport string, w io.Writer) error {
	return ErrPdfTemplatesNotSupported
}
