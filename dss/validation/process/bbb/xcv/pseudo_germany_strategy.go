// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/pseudo/PseudoGermanyStrategy.java (DSS 6.5.RC1).
package xcv

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/utils"
)

// pseudoGermanyCountryCode is the Germany country code.
const pseudoGermanyCountryCode = "DE"

// pseudoGermanySuffix is the suffix used for certificate's pseudo definition
// within X500 Attributes.
const pseudoGermanySuffix = ":PN"

// PseudoGermanyStrategy extracts pseudo information for German certificates.
type PseudoGermanyStrategy struct{}

// NewPseudoGermanyStrategy is the default constructor.
func NewPseudoGermanyStrategy() *PseudoGermanyStrategy {
	return &PseudoGermanyStrategy{}
}

// GetPseudo gets pseudo for the certificate. Port of getPseudo(CertificateWrapper).
func (s *PseudoGermanyStrategy) GetPseudo(certificate *diagnostic.CertificateWrapper) string {
	if pseudoGermanyCountryCode == certificate.CountryName() {
		cn := certificate.CommonName()
		if utils.EndsWithIgnoreCase(cn, pseudoGermanySuffix) {
			return strings.ReplaceAll(cn, pseudoGermanySuffix, "")
		}
	}
	return ""
}
