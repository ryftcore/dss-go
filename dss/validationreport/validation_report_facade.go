// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/ValidationReportFacade.java
// (DSS 6.5.RC1).
//
// Java's ValidationReportFacade extends dss-jaxb-common's AbstractJaxbFacade<T>; this port
// implements Marshal/Unmarshal directly with encoding/xml (the low-level byte-exact
// Marshal/Unmarshal that the KAT actually exercises lives in jaxb/xml.go, which this facade
// does not call). XSD-schema validation is not implemented - see ValidationReportUtils's
// deferred Schema() stub. Facade has no HTML/PDF report generation methods to
// port, unlike Facade/Facade - the ETSI Validation Report has no
// XSLT stylesheets of its own.
package validationreport

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

// Facade performs marshalling/unmarshalling operations for
// an ETSI Validation report.
type Facade struct{}

// NewValidationReportFacade creates a new facade. Port of newFacade().
func NewValidationReportFacade() *Facade {
	return &Facade{}
}

// Marshal returns the XML representation of validationReport. Port of
// marshall(T) (schema validation always requested in Java; see the file
// header for why this port does not perform it).
func (f *Facade) Marshal(validationReport *jaxb.ValidationReportType) (string, error) {
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
func (f *Facade) MarshalToWriter(validationReport *jaxb.ValidationReportType, w io.Writer) error {
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
func (f *Facade) Unmarshal(r io.Reader) (*jaxb.ValidationReportType, error) {
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
func (f *Facade) UnmarshalString(xmlObject string) (*jaxb.ValidationReportType, error) {
	return f.Unmarshal(bytes.NewReader([]byte(xmlObject)))
}
