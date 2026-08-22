// Ported from dss-policy-crypto-json/.../json/CryptographicSuiteJsonUtils.java (DSS 6.5.RC1).
//
// DEFERRED: Java's CryptographicSuiteJsonUtils (via its
// dss-json-common JSONSchemaAbstractUtils base) validates a document
// against the ETSI TS 119 322 JSON Schema (19322algocatjsonschema.json).
// No JSON Schema validator ships in the Go stdlib and none has been added
// without tech-lead sign-off (PORTING.md "Dependency
// policy"), the same way dss/diagnostic/diagnostic_data_xml_definer.go's
// deferred Schema() stub does. This file keeps the schema URI/location
// constants and embeds the schema document for documentation/testdata
// purposes, but ValidateAgainstSchema is a stub returning an error.
package cryptojson

import (
	_ "embed"
	"errors"
)

// CryptographicSuiteJsonSchemaURI is the cryptographic suite schema URI.
// Ports CryptographicSuiteJsonUtils#SCHEMA_URI.
const CryptographicSuiteJsonSchemaURI = "http://uri.etsi.org/19322/algoCatSchema"

// CryptographicSuiteJsonSchemaLocation is the cryptographic suite schema's
// location, relative to the embedded resource root. Ports
// CryptographicSuiteJsonUtils#SCHEMA_LOCATION.
const CryptographicSuiteJsonSchemaLocation = "/schema/19322algocatjsonschema.json"

// cryptographicSuiteJsonSchema is a byte-identical copy of upstream's
// src/main/resources/schema/19322algocatjsonschema.json, embedded for
// documentation/testdata purposes - see the file header for why
// ValidateAgainstSchema does not use it to actually validate.
//
//go:embed resources/19322algocatjsonschema.json
var cryptographicSuiteJsonSchema []byte

// ErrCryptographicSuiteJsonSchemaNotSupported is returned by
// ValidateAgainstSchema: no JSON Schema validator ships in the Go stdlib
// and none has been added (see the file header).
var ErrCryptographicSuiteJsonSchemaNotSupported = errors.New("cryptojson: 19322algocatjsonschema.json schema validation is not implemented in this port (deferred, see CryptographicSuiteJsonUtils)")

// ValidateAgainstSchema is a stub for JSON-Schema validation of a
// cryptographic suite document against 19322algocatjsonschema.json. Ports
// JSONSchemaAbstractUtils#validateAgainstSchema (as inherited by
// CryptographicSuiteJsonUtils); see the file header's DEFERRED note.
func ValidateAgainstSchema([]byte) ([]string, error) {
	return nil, ErrCryptographicSuiteJsonSchemaNotSupported
}
