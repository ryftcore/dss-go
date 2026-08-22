// Ported from dss-policy-crypto-xml/.../xml/CryptographicSuiteXmlFacade.java (DSS 6.5.RC1).
//
// Java's CryptographicSuiteXmlFacade extends dss-jaxb-common's
// AbstractJaxbFacade<SecuritySuitabilityPolicyType>; this port implements unmarshalling
// directly with encoding/xml against this package's own xml_types.go structs instead of a
// JAXBContext.
package cryptoxml

import (
	"encoding/xml"
	"errors"
	"io"

	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// CryptographicSuiteXmlFacade performs unmarshalling for the ETSI TS
// 119 312/322 XML schema. Ports CryptographicSuiteXmlFacade.
type CryptographicSuiteXmlFacade struct{}

// NewCryptographicSuiteXmlFacade initializes a new
// CryptographicSuiteXmlFacade. Ports CryptographicSuiteXmlFacade#newFacade.
func NewCryptographicSuiteXmlFacade() *CryptographicSuiteXmlFacade {
	return &CryptographicSuiteXmlFacade{}
}

// GetCryptographicSuite gets the cryptographic suite from r. Ports
// CryptographicSuiteXmlFacade#getCryptographicSuite(InputStream).
//
// Java's Objects.requireNonNull(is, "The provided cryptographic suite is
// null") becomes a returned error per PORTING.md; the inherited
// AbstractJaxbFacade#unmarshall(InputStream)'s schema-validated unmarshal
// is - as throughout this package, see ValidateAgainstSchema's doc comment
// - a schema-less encoding/xml.Decode instead.
func (f *CryptographicSuiteXmlFacade) GetCryptographicSuite(r io.Reader) (*modelpolicy.CryptographicSuiteCatalogue, error) {
	if r == nil {
		return nil, errors.New("the provided cryptographic suite is null")
	}
	var policy SecuritySuitabilityPolicyType
	if err := xml.NewDecoder(r).Decode(&policy); err != nil {
		return nil, err
	}
	return newCryptographicSuiteXmlCatalogue(&policy), nil
}
