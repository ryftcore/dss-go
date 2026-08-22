// Ported from dss-policy-crypto-xml/.../xml/CryptographicSuiteXmlFactory.java (DSS 6.5.RC1).
package cryptoxml

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// defaultCryptographicSuite is a byte-identical copy of upstream's
// src/main/resources/suite/dss-crypto-suite.xml, embedded because
// CryptographicSuiteXmlFactory#loadDefaultCryptographicSuite loads it from
// the classpath at runtime (DEFAULT_CRYPTOGRAPHIC_SUITES_LOCATION =
// "/suite/dss-crypto-suite.xml").
//
//go:embed resources/dss-crypto-suite.xml
var defaultCryptographicSuite []byte

// CryptographicSuiteXmlFactory is an implementation of a cryptographic
// suite using the XML schema defined in ETSI TS 119 322. Ports
// CryptographicSuiteXmlFactory.
type CryptographicSuiteXmlFactory struct{}

var _ modelpolicy.CryptographicSuiteFactory = (*CryptographicSuiteXmlFactory)(nil)

// NewCryptographicSuiteXmlFactory is the default constructor.
func NewCryptographicSuiteXmlFactory() *CryptographicSuiteXmlFactory {
	return &CryptographicSuiteXmlFactory{}
}

// IsSupported ports CryptographicSuiteXmlFactory#isSupported.
//
// Java delegates to CryptographicSuiteXmlFacade#unmarshall(is, false) - an
// unvalidated (schema-less) unmarshal, matching this port's plain
// encoding/xml.Decode (see ValidateAgainstSchema's doc comment for why
// full XSD validation is out of scope for this port, same as
// cryptojson.ValidateAgainstSchema).
func (f *CryptographicSuiteXmlFactory) IsSupported(cryptographicSuiteDocument model.DSSDocument) bool {
	rc, err := cryptographicSuiteDocument.OpenStream()
	if err != nil {
		return false
	}
	defer rc.Close()

	var policy SecuritySuitabilityPolicyType
	return xml.NewDecoder(rc).Decode(&policy) == nil
}

// LoadDefaultCryptographicSuite ports
// CryptographicSuiteXmlFactory#loadDefaultCryptographicSuite.
func (f *CryptographicSuiteXmlFactory) LoadDefaultCryptographicSuite() *modelpolicy.CryptographicSuiteCatalogue {
	return f.LoadCryptographicSuiteFromReader(bytes.NewReader(defaultCryptographicSuite))
}

// LoadCryptographicSuite ports
// CryptographicSuiteXmlFactory#loadCryptographicSuite(DSSDocument).
func (f *CryptographicSuiteXmlFactory) LoadCryptographicSuite(cryptographicSuiteDocument model.DSSDocument) *modelpolicy.CryptographicSuiteCatalogue {
	rc, err := cryptographicSuiteDocument.OpenStream()
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	return f.LoadCryptographicSuiteFromReader(rc)
}

// LoadCryptographicSuiteFromReader ports
// CryptographicSuiteXmlFactory#loadCryptographicSuite(InputStream). Named
// per PORTING.md's Java-overload convention (see
// model/policy.CryptographicSuiteFactory's doc comment).
//
// Java wraps any failure (including the try-with-resources close) in an
// unchecked UnsupportedOperationException; ported as a panic, matching the
// model/policy.CryptographicSuiteFactory interface, whose loader methods
// return no error.
func (f *CryptographicSuiteXmlFactory) LoadCryptographicSuiteFromReader(cryptographicSuiteInputStream io.Reader) *modelpolicy.CryptographicSuiteCatalogue {
	if closer, ok := cryptographicSuiteInputStream.(io.Closer); ok {
		defer closer.Close()
	}
	catalogue, err := NewCryptographicSuiteXmlFacade().GetCryptographicSuite(cryptographicSuiteInputStream)
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	return catalogue
}
