// Ported from dss-policy-crypto-xml/.../xml/CryptographicSuiteXmlUtils.java (DSS 6.5.RC1).
//
// DEFERRED: Java's CryptographicSuiteXmlUtils (via its dss-jaxb-common
// XSDAbstractUtils base) validates a document against rfc5698.xsd +
// 19322algocatxmlschema.xsd (plus xmldsig-core-schema.xsd, transitively).
// No XSD validator ships in the Go stdlib and none has been added to this
// port without tech-lead sign-off (PORTING.md "Dependency policy"), the
// same way dss/diagnostic/diagnostic_data_xml_definer.go's deferred
// Schema() stub and cryptojson.ValidateAgainstSchema do. This file keeps
// the schema locations and embeds both schema documents for
// documentation/testdata purposes, but ValidateAgainstSchema is a stub
// returning an error. getJAXBContext() has no port: it exists purely to
// drive JAXB's own (un)marshaller, which this package's
// CryptographicSuiteXmlFacade replaces directly with encoding/xml (see
// that file's header) - the same omission dss/diagnostic/
// diagnostic_data_xml_definer.go documents for its own sibling.
package cryptoxml

import (
	_ "embed"
	"errors"
)

// CryptographicSuiteXmlRFC5698SchemaLocation is the RFC 5698 base schema's
// location, relative to the embedded resource root. Ports
// CryptographicSuiteXmlUtils#CRYPTO_SUITES_CATALOGUES_SCHEMA_LOCATION.
const CryptographicSuiteXmlRFC5698SchemaLocation = "/xsd/rfc5698.xsd"

// CryptographicSuiteXmlAlgoCatSchemaLocation is the ETSI TS 119 322 algocat
// extension schema's location, relative to the embedded resource root.
// Ports CryptographicSuiteXmlUtils#CRYPTO_SUITES_ALGOCAT_SCHEMA_LOCATION.
const CryptographicSuiteXmlAlgoCatSchemaLocation = "/xsd/19322algocatxmlschema.xsd"

// cryptographicSuiteXmlRFC5698Schema/cryptographicSuiteXmlAlgoCatSchema are
// byte-identical copies of upstream's src/main/resources/xsd/rfc5698.xsd
// and .../19322algocatxmlschema.xsd, embedded for
// documentation/testdata purposes - see the file header for why
// ValidateAgainstSchema does not use them to actually validate.
//
//go:embed resources/rfc5698.xsd
var cryptographicSuiteXmlRFC5698Schema []byte

//go:embed resources/19322algocatxmlschema.xsd
var cryptographicSuiteXmlAlgoCatSchema []byte

// ErrCryptographicSuiteXmlSchemaNotSupported is returned by
// ValidateAgainstSchema: no XSD validator ships in the Go stdlib and none
// has been added to this port (see the file header).
var ErrCryptographicSuiteXmlSchemaNotSupported = errors.New("cryptoxml: rfc5698.xsd/19322algocatxmlschema.xsd schema validation is not implemented in this port (deferred, see CryptographicSuiteXmlUtils)")

// ValidateAgainstSchema is a stub for XSD validation of a cryptographic
// suite document against rfc5698.xsd/19322algocatxmlschema.xsd. Ports
// XSDAbstractUtils#validateAgainstSchema (as inherited by
// CryptographicSuiteXmlUtils, via its getXSDSources() override); see the
// file header's DEFERRED note.
func ValidateAgainstSchema([]byte) ([]string, error) {
	return nil, ErrCryptographicSuiteXmlSchemaNotSupported
}
