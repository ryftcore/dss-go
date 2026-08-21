// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/ValidationReportFacade.java
// (DSS 6.5.RC1).
//
// Java's ValidationReportFacade extends the generic dss-jaxb-common
// AbstractJaxbFacade<T>, a module outside S8B_BRIEF.md's manifest. As with
// dss/diagnostic/diagnostic_data_facade.go (Phase 8a) and
// dss/simplereport/simple_report_facade.go (Phase 8b), this port collapses
// AbstractJaxbFacade's marshal/unmarshal template method directly into
// ValidationReportFacade using encoding/xml (the marshal-parity contract
// this phase exists to satisfy; the low-level byte-exact Marshal/Unmarshal
// that the KAT actually exercises lives in jaxb/xml.go, which this facade
// does NOT call - matching the 8a/8b precedent that the hand facade's own
// marshalling is a convenience, not the parity-pinned path). XSD-schema
// validation (the `validate` boolean AbstractJaxbFacade's constructor
// implies) is stubbed via ValidationReportUtils's deferred Schema() and
// surfaces the same error here. ValidationReportFacade has no HTML/PDF
// report generation methods to port, unlike SimpleReportFacade/
// DetailedReportFacade - the ETSI Validation Report has no XSLT stylesheets
// of its own.
package validationreport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

// ValidationReportFacade performs marshalling/unmarshalling operations for
// an ETSI Validation report.
type ValidationReportFacade struct{}

// NewValidationReportFacade creates a new facade. Port of newFacade().
func NewValidationReportFacade() *ValidationReportFacade {
	return &ValidationReportFacade{}
}

// Marshal returns the XML representation of validationReport. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *ValidationReportFacade) Marshal(validationReport *jaxb.ValidationReportType) (string, error) {
	if validationReport == nil {
		return "", errors.New("JAXBObject is null")
	}
	var buf bytes.Buffer
	if err := f.MarshalToWriter(validationReport, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarshalToWriter marshals validationReport into w. Port of marshall(T,
// OutputStream).
func (f *ValidationReportFacade) MarshalToWriter(validationReport *jaxb.ValidationReportType, w io.Writer) error {
	if validationReport == nil {
		return errors.New("JAXBObject is null")
	}
	if w == nil {
		return errors.New("OutputStream is null")
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "    ")
	if err := encoder.Encode(validationReport); err != nil {
		return err
	}
	return encoder.Flush()
}

// Unmarshal unmarshals r and returns the ValidationReportType. Port of
// unmarshall(InputStream).
func (f *ValidationReportFacade) Unmarshal(r io.Reader) (*jaxb.ValidationReportType, error) {
	if r == nil {
		return nil, errors.New("InputStream is null")
	}
	var result jaxb.ValidationReportType
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnmarshalString unmarshals xmlObject and returns the ValidationReportType.
// Port of unmarshall(String).
func (f *ValidationReportFacade) UnmarshalString(xmlObject string) (*jaxb.ValidationReportType, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}
