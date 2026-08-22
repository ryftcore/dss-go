// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReportFacade.java
// (DSS 6.5.RC1).
//
// Java's DetailedReportFacade extends dss-jaxb-common's AbstractJaxbFacade<T>; this port
// implements Marshal/Unmarshal by delegating to jaxb.Marshal/jaxb.Unmarshal (jaxb/xml.go),
// which reproduce the JAXB reference implementation's bytes exactly. XSD-schema validation and
// HTML/PDF report generation are not implemented - see DetailedReportXmlDefiner's deferred
// Schema()/HtmlBootstrap4Templates()/PdfTemplates() stubs.
package detailedreport

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
)

// Facade contains methods for DetailedReport generation.
type Facade struct{}

// NewDetailedReportFacade creates a new DetailedReportFacade. Port of
// newFacade().
func NewDetailedReportFacade() *Facade {
	return &Facade{}
}

// Marshal returns the XML representation of detailedReportJaxb. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *Facade) Marshal(detailedReportJaxb *jaxb.XmlDetailedReport) (string, error) {
	if detailedReportJaxb == nil {
		return "", errors.New("JAXBObject is null")
	}
	out, err := jaxb.Marshal(detailedReportJaxb)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Unmarshal unmarshals xmlObject and returns the XmlDetailedReport. Port of
// unmarshall(String) (unmarshall(InputStream) is not ported separately: Go's
// []byte-based jaxb.Unmarshal already covers both Java overloads).
func (f *Facade) Unmarshal(xmlObject string) (*jaxb.XmlDetailedReport, error) {
	if xmlObject == "" {
		return nil, errors.New("InputStream is null")
	}
	return jaxb.Unmarshal([]byte(xmlObject))
}

// GenerateHtmlReport generates a Bootstrap 4 Detailed report. Port of
// generateHtmlReport(XmlDetailedReport); XSLT execution is deferred, see
// DetailedReportXmlDefiner.
func (f *Facade) GenerateHtmlReport(detailedReport *jaxb.XmlDetailedReport) (string, error) {
	return "", ErrHtmlBootstrap4TemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Detailed report
// from a string. Port of generateHtmlReport(String); XSLT execution is
// deferred, see DetailedReportXmlDefiner.
func (f *Facade) GenerateHtmlReportFromMarshalled(marshalledDetailedReport string) (string, error) {
	return "", ErrHtmlBootstrap4TemplatesNotSupported
}

// GeneratePdfReport generates a PDF Detailed report. Port of
// generatePdfReport(XmlDetailedReport, Result); XSLT execution is deferred,
// see DetailedReportXmlDefiner.
func (f *Facade) GeneratePdfReport(detailedReport *jaxb.XmlDetailedReport) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Detailed report from a
// string. Port of generatePdfReport(String, Result); XSLT execution is
// deferred, see DetailedReportXmlDefiner.
func (f *Facade) GeneratePdfReportFromMarshalled(marshalledDetailedReport string) error {
	return ErrPdfTemplatesNotSupported
}
