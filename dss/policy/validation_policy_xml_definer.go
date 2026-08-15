// Ported from dss-policy-jaxb/.../policy/ValidationPolicyXmlDefiner.java (DSS
// 6.5.RC1).
//
// DEFERRED: Java's ValidationPolicyXmlDefiner centers on two thread-safe
// singletons built by javax.xml.validation/jakarta.xml.bind machinery that
// has no Go stdlib equivalent: a JAXBContext (needed only to drive JAXB's own
// (un)marshaller, which policy/jaxb's Marshal/Unmarshal replace directly with
// encoding/xml - see that package's xml.go) and an XSD Schema used to
// validate a validation-policy document against policy.xsd. Per
// S8A_BRIEF.md's marshal-parity contract ("policy.xsd ... copied under
// testdata/ + embedded where upstream loads them at runtime"), this port
// keeps the schema resource path and embeds its bytes, but Schema() is a stub
// returning an error, following the precedent set by
// dss/diagnostic/diagnostic_data_xml_definer.go: no XSD validator ships in
// the Go stdlib and none has been added to this port without tech-lead
// sign-off (PORTING.md "Dependency policy"). ValidationPolicyFacade (the half
// of this pair the marshal-parity KAT actually exercises) does not depend on
// Schema().
package policy

import (
	_ "embed"
	"errors"
)

// ValidationPolicySchemaLocation is the location of the Validation Policy XSD
// schema, relative to the embedded resource root. Ports
// ValidationPolicyXmlDefiner#VALIDATION_POLICY_SCHEMA_LOCATION.
const ValidationPolicySchemaLocation = "/xsd/policy.xsd"

// validationPolicyXSD is a byte-identical copy of upstream's
// src/main/resources/xsd/policy.xsd, embedded for documentation/testdata
// purposes per S8A_BRIEF.md - see the file header for why Schema() does not
// use it to actually validate.
//
//go:embed jaxb/testdata/xsd/policy.xsd
var validationPolicyXSD string

// ErrValidationPolicySchemaNotSupported is returned by Schema(): no XSD
// validator ships in the Go stdlib and none has been added to this port (see
// the file header).
var ErrValidationPolicySchemaNotSupported = errors.New("policy: policy.xsd schema validation is not implemented in this port (deferred, see ValidationPolicyXmlDefiner)")

// Schema is a stub for the XSD Schema used to validate a validation-policy
// document. Ports ValidationPolicyXmlDefiner#getSchema(); see the DEFERRED
// note above.
func Schema() (struct{}, error) {
	return struct{}{}, ErrValidationPolicySchemaNotSupported
}
