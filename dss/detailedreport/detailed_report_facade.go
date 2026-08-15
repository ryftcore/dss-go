// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReportFacade.java
// (DSS 6.5.RC1).
//
// Java's DetailedReportFacade extends the generic dss-jaxb-common
// AbstractJaxbFacade<T>, a module outside this phase's manifest (only
// dss-detailed-report-jaxb's own files are). dss/diagnostic/jaxb's
// DiagnosticDataFacade precedent collapses AbstractJaxbFacade's marshal/
// unmarshal template method into a direct encoding/xml call instead - this
// port does one better, since jaxb.Marshal/jaxb.Unmarshal (jaxb/xml.go)
// already exist and reproduce the JAXB reference implementation's bytes
// exactly (the marshal-parity KAT's whole purpose): Marshal/Unmarshal below
// call them directly, so this facade is actually byte-exact where the
// diagnostic-module precedent it otherwise mirrors is not. Flagged here per
// the porting brief's "judgment calls, cross-chunk assumptions" instruction.
//
// XSD-schema validation (the `validate` booleans on every Java overload) and
// the HTML/PDF XSLT generation methods are stubbed via
// DetailedReportXmlDefiner's deferred Schema()/HtmlBootstrap4Templates()/
// PdfTemplates() (see that file's header) and surface the same error here.
package detailedreport

import (
	"errors"

	"github.com/utain/esig/dss/detailedreport/jaxb"
)

// DetailedReportFacade contains methods for DetailedReport generation.
type DetailedReportFacade struct{}

// NewDetailedReportFacade creates a new DetailedReportFacade. Port of
// newFacade().
func NewDetailedReportFacade() *DetailedReportFacade {
	return &DetailedReportFacade{}
}

// Marshal returns the XML representation of detailedReportJaxb. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *DetailedReportFacade) Marshal(detailedReportJaxb *jaxb.XmlDetailedReport) (string, error) {
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
func (f *DetailedReportFacade) Unmarshal(xmlObject string) (*jaxb.XmlDetailedReport, error) {
	if xmlObject == "" {
		return nil, errors.New("InputStream is null")
	}
	return jaxb.Unmarshal([]byte(xmlObject))
}

// GenerateHtmlReport generates a Bootstrap 4 Detailed report. Port of
// generateHtmlReport(XmlDetailedReport); XSLT execution is deferred, see
// DetailedReportXmlDefiner.
func (f *DetailedReportFacade) GenerateHtmlReport(detailedReport *jaxb.XmlDetailedReport) (string, error) {
	return "", ErrHtmlBootstrap4TemplatesNotSupported
}

// GenerateHtmlReportFromMarshalled generates a Bootstrap 4 Detailed report
// from a string. Port of generateHtmlReport(String); XSLT execution is
// deferred, see DetailedReportXmlDefiner.
func (f *DetailedReportFacade) GenerateHtmlReportFromMarshalled(marshalledDetailedReport string) (string, error) {
	return "", ErrHtmlBootstrap4TemplatesNotSupported
}

// GeneratePdfReport generates a PDF Detailed report. Port of
// generatePdfReport(XmlDetailedReport, Result); XSLT execution is deferred,
// see DetailedReportXmlDefiner.
func (f *DetailedReportFacade) GeneratePdfReport(detailedReport *jaxb.XmlDetailedReport) error {
	return ErrPdfTemplatesNotSupported
}

// GeneratePdfReportFromMarshalled generates a PDF Detailed report from a
// string. Port of generatePdfReport(String, Result); XSLT execution is
// deferred, see DetailedReportXmlDefiner.
func (f *DetailedReportFacade) GeneratePdfReportFromMarshalled(marshalledDetailedReport string) error {
	return ErrPdfTemplatesNotSupported
}
