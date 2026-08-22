// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/ValidationReportUtils.java
// (DSS 6.5.RC1).
//
// DEFERRED: Java's ValidationReportUtils (an XSDAbstractUtils singleton)
// builds two thread-safe resources this port has no stdlib equivalent for:
// a JAXBContext (needed only to drive JAXB's own marshaller/unmarshaller,
// which Facade.go replaces with encoding/xml) and an XSD
// Schema assembled from 1910202xmlSchema.xsd plus the trusted-list, XAdES
// and XMLDSig schemas it imports (via TrustedList211Utils.getXSDSources()).
// The schema location constant is kept for documentation/testdata
// purposes, but Schema() is a stub returning an error: no Go stdlib XSD
// validator exists, and adding a third-party one is outside the
// stdlib-first dependency policy without tech-lead sign-off (PORTING.md
// "Dependency policy") - the same deferral dss/simplereport's
// SimpleReportXmlDefiner documents for SimpleReport.xsd.
// Facade.Marshal/Unmarshal (the half of this pair the
// marshal-parity KAT actually exercises, via jaxb.Marshal/Unmarshal) does
// not depend on it. ObjectFactory.java (ValidationReportUtils.OBJECT_FACTORY)
// has no Go counterpart either, for the reason jaxb/doc.go gives.
package validationreport

import "errors"

// ValidationReportSchemaLocation is the location of the ETSI Validation
// Report XSD schema, relative to the embedded/testdata resource root. Port
// of ValidationReportUtils.VALIDATION_REPORT_SCHEMA_LOCATION.
const ValidationReportSchemaLocation = "/xsd/1910202xmlSchema.xsd"

// ErrXSDSchemaNotSupported is returned by Schema(): no XSD validator ships
// in the Go stdlib and none has been added to this port (see the file
// header).
var ErrXSDSchemaNotSupported = errors.New("validationreport: 1910202xmlSchema.xsd schema validation is not implemented in this port (deferred, see ValidationReportUtils)")

// Schema is a stub for the XSD Schema of the ETSI Validation Report. Port of
// getSchema() (XSDAbstractUtils.getSchema()); see the DEFERRED note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrXSDSchemaNotSupported
}
