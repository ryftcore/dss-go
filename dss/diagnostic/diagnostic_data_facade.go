// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/DiagnosticDataFacade.java (DSS 6.5.RC1).
//
// Java's DiagnosticDataFacade extends the generic dss-jaxb-common AbstractJaxbFacade<T>, a
// module outside S8A_BRIEF.md's manifest (dss-i18n, dss-policy-jaxb(+crypto-json/xml),
// dss-diagnostic-jaxb only). Rather than leaving DiagnosticDataFacade unbuildable on an unported
// base class, this port collapses AbstractJaxbFacade's marshal/unmarshal template method
// directly into DiagnosticDataFacade using encoding/xml (the marshal-parity contract this phase
// exists to satisfy), and flags the collapse here per S8A_BRIEF.md's "Report: cross-chunk
// assumptions" instruction. XSD-schema validation (the `validate` booleans on every Java
// overload) and the SVG XSLT generation methods are stubbed via DiagnosticDataXmlDefiner's
// deferred Schema()/SvgTemplates() (see that file's header) and surface the same error here.
package diagnostic

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// DiagnosticDataFacade is used to marshal/unmarshal a DiagnosticData report.
type DiagnosticDataFacade struct{}

// NewDiagnosticDataFacade creates a new instance of DiagnosticDataFacade. Port of newFacade().
func NewDiagnosticDataFacade() *DiagnosticDataFacade {
	return &DiagnosticDataFacade{}
}

// Marshal returns the XML representation of diagnosticDataJaxb. Port of marshall(T) (schema
// validation always requested in Java; see the file header for why this port does not perform
// it).
func (f *DiagnosticDataFacade) Marshal(diagnosticDataJaxb *jaxb.XmlDiagnosticData) (string, error) {
	if diagnosticDataJaxb == nil {
		return "", errors.New("JAXBObject is null")
	}
	var buf bytes.Buffer
	if err := f.MarshalToWriter(diagnosticDataJaxb, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarshalToWriter marshals diagnosticDataJaxb into w. Port of marshall(T, OutputStream).
func (f *DiagnosticDataFacade) MarshalToWriter(diagnosticDataJaxb *jaxb.XmlDiagnosticData, w io.Writer) error {
	if diagnosticDataJaxb == nil {
		return errors.New("JAXBObject is null")
	}
	if w == nil {
		return errors.New("OutputStream is null")
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "    ")
	if err := encoder.Encode(diagnosticDataJaxb); err != nil {
		return err
	}
	return encoder.Flush()
}

// Unmarshal unmarshals r and returns the XmlDiagnosticData. Port of unmarshall(InputStream).
func (f *DiagnosticDataFacade) Unmarshal(r io.Reader) (*jaxb.XmlDiagnosticData, error) {
	if r == nil {
		return nil, errors.New("InputStream is null")
	}
	var result jaxb.XmlDiagnosticData
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	// The JAXB unmarshaller resolves the document's IDREF attributes natively;
	// encoding/xml does not, so the ported model exposes jaxb.Link for it (and
	// jaxb.Unmarshal, which every in-tree caller uses, already calls it). Doing
	// it here too keeps this facade - the public entry point Java callers use -
	// behaviourally identical: without it every ChainItem/SigningCertificate
	// reference stays an unresolved stub carrying only the raw id, and the whole
	// validation engine reads empty certificate data off it.
	jaxb.Link(&result)
	return &result, nil
}

// UnmarshalString unmarshals xmlObject and returns the XmlDiagnosticData. Port of
// unmarshall(String).
func (f *DiagnosticDataFacade) UnmarshalString(xmlObject string) (*jaxb.XmlDiagnosticData, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}

// GenerateSVG generates a SVG representation of the diagnostic data. Port of
// generateSVG(XmlDiagnosticData); XSLT execution is deferred, see DiagnosticDataXmlDefiner.
func (f *DiagnosticDataFacade) GenerateSVG(diagnosticDataJaxb *jaxb.XmlDiagnosticData) (string, error) {
	return "", ErrSvgTemplatesNotSupported
}

// GenerateSVGToWriter generates a SVG representation of the diagnostic data into w. Port of
// generateSVG(XmlDiagnosticData, Result); XSLT execution is deferred, see
// DiagnosticDataXmlDefiner.
func (f *DiagnosticDataFacade) GenerateSVGToWriter(diagnosticDataJaxb *jaxb.XmlDiagnosticData, w io.Writer) error {
	return ErrSvgTemplatesNotSupported
}

// GenerateSVGFromMarshalled generates a SVG representation from already-marshalled diagnostic
// data. Port of generateSVG(String); XSLT execution is deferred, see DiagnosticDataXmlDefiner.
func (f *DiagnosticDataFacade) GenerateSVGFromMarshalled(marshalledDiagnosticData string) (string, error) {
	return "", ErrSvgTemplatesNotSupported
}

// GenerateSVGFromMarshalledToWriter generates a SVG representation from already-marshalled
// diagnostic data into w. Port of generateSVG(String, Result); XSLT execution is deferred, see
// DiagnosticDataXmlDefiner.
func (f *DiagnosticDataFacade) GenerateSVGFromMarshalledToWriter(marshalledDiagnosticData string, w io.Writer) error {
	return ErrSvgTemplatesNotSupported
}
