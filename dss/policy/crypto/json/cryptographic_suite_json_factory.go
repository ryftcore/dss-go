// Ported from dss-policy-crypto-json/.../json/CryptographicSuiteJsonFactory.java (DSS 6.5.RC1).
package cryptojson

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// defaultCryptographicSuite is a byte-identical copy of upstream's
// src/main/resources/suite/dss-crypto-suite.json, embedded because
// CryptographicSuiteJsonFactory#loadDefaultCryptographicSuite loads it from
// the classpath at runtime (DEFAULT_CRYPTOGRAPHIC_SUITES_LOCATION =
// "/suite/dss-crypto-suite.json") - per S8A_BRIEF.md's "Schema/suite
// resources embedded" instruction.
//
//go:embed resources/dss-crypto-suite.json
var defaultCryptographicSuite []byte

// CryptographicSuiteJsonFactory is an implementation of a cryptographic
// suite using the JSON schema defined in ETSI TS 119 322. Ports
// CryptographicSuiteJsonFactory.
type CryptographicSuiteJsonFactory struct{}

var _ modelpolicy.CryptographicSuiteFactory = (*CryptographicSuiteJsonFactory)(nil)

// NewCryptographicSuiteJsonFactory is the default constructor.
func NewCryptographicSuiteJsonFactory() *CryptographicSuiteJsonFactory {
	return &CryptographicSuiteJsonFactory{}
}

// IsSupported ports CryptographicSuiteJsonFactory#isSupported.
//
// Java delegates to CryptographicSuiteJsonUtils#validateAgainstSchema,
// full ETSI TS 119 322 JSON Schema validation; this port instead performs
// the structural check every successful load already requires - the
// document parses as JSON and its root object carries a
// "SecuritySuitabilityPolicy" property - since no JSON Schema validator is
// available (see ValidateAgainstSchema's doc comment). This is a strictly
// weaker check than upstream's (accepts some malformed-but-structurally-
// shaped documents upstream would reject), documented here per
// S8A_BRIEF.md's "JAXB quirks / deviations" reporting instruction.
func (f *CryptographicSuiteJsonFactory) IsSupported(cryptographicSuiteDocument model.DSSDocument) bool {
	rc, err := cryptographicSuiteDocument.OpenStream()
	if err != nil {
		return false
	}
	defer rc.Close()

	obj, err := parseJSONObject(rc)
	if err != nil {
		return false
	}
	return obj.getAsObject(jsonConstraintSecuritySuitabilityPolicy) != nil
}

// LoadDefaultCryptographicSuite ports
// CryptographicSuiteJsonFactory#loadDefaultCryptographicSuite.
func (f *CryptographicSuiteJsonFactory) LoadDefaultCryptographicSuite() *modelpolicy.CryptographicSuiteCatalogue {
	return f.LoadCryptographicSuiteFromReader(bytes.NewReader(defaultCryptographicSuite))
}

// LoadCryptographicSuite ports
// CryptographicSuiteJsonFactory#loadCryptographicSuite(DSSDocument).
func (f *CryptographicSuiteJsonFactory) LoadCryptographicSuite(cryptographicSuiteDocument model.DSSDocument) *modelpolicy.CryptographicSuiteCatalogue {
	rc, err := cryptographicSuiteDocument.OpenStream()
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	return f.LoadCryptographicSuiteFromReader(rc)
}

// LoadCryptographicSuiteFromReader ports
// CryptographicSuiteJsonFactory#loadCryptographicSuite(InputStream). Named
// per PORTING_PLAN's Java-overload convention (see
// model/policy.CryptographicSuiteFactory's doc comment).
//
// Java wraps any failure (including the try-with-resources close) in an
// unchecked UnsupportedOperationException; ported as a panic, matching the
// model/policy.CryptographicSuiteFactory interface, whose loader methods
// return no error.
func (f *CryptographicSuiteJsonFactory) LoadCryptographicSuiteFromReader(cryptographicSuiteInputStream io.Reader) *modelpolicy.CryptographicSuiteCatalogue {
	if closer, ok := cryptographicSuiteInputStream.(io.Closer); ok {
		defer closer.Close()
	}
	jsonObj, err := parseJSONObject(cryptographicSuiteInputStream)
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	securitySuitabilityPolicyType := jsonObj.getAsObject(jsonConstraintSecuritySuitabilityPolicy)
	if securitySuitabilityPolicyType == nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : The root element of JSON shall be a JSON object of '%s' type!", jsonConstraintSecuritySuitabilityPolicy))
	}
	return newCryptographicSuiteJsonCatalogue(securitySuitabilityPolicyType)
}
