// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteFactory.java (DSS 6.5.RC1).
package policy

import (
	"io"

	"github.com/utain/esig/dss/model"
)

// CryptographicSuiteFactory contains methods to load a CryptographicSuite
// object.
//
// Java overloads loadCryptographicSuite for DSSDocument and InputStream;
// Go cannot overload by parameter type, so the InputStream variant is
// named LoadCryptographicSuiteFromReader.
type CryptographicSuiteFactory interface {
	// IsSupported evaluates whether the cryptographic suite DSSDocument
	// is supported by the current implementation.
	IsSupported(cryptographicSuiteDocument model.DSSDocument) bool

	// LoadDefaultCryptographicSuite loads a default cryptographic suite
	// provided by the application.
	LoadDefaultCryptographicSuite() *CryptographicSuiteCatalogue

	// LoadCryptographicSuite loads a cryptographic suite from a
	// DSSDocument provided to the method.
	LoadCryptographicSuite(cryptographicSuiteDocument model.DSSDocument) *CryptographicSuiteCatalogue

	// LoadCryptographicSuiteFromReader loads a cryptographic suite from
	// an io.Reader provided to the method. Ports the
	// loadCryptographicSuite(InputStream) overload.
	LoadCryptographicSuiteFromReader(cryptographicSuiteInputStream io.Reader) *CryptographicSuiteCatalogue
}
