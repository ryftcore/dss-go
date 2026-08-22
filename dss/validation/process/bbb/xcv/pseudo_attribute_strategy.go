// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/pseudo/PseudoAttributeStrategy.java (DSS 6.5.RC1).
package xcv

import "github.com/ryftcore/dss-go/dss/diagnostic"

// pseudoAttributeStrategy extracts the pseudo string defined in X500
// Attributes for the certificate. Unexported: Java's package-private class
// (no public modifier) is only used within this file's package.
type pseudoAttributeStrategy struct{}

// newPseudoAttributeStrategy is the default constructor.
func newPseudoAttributeStrategy() *pseudoAttributeStrategy {
	return &pseudoAttributeStrategy{}
}

// GetPseudo gets pseudo for the certificate. Port of getPseudo(CertificateWrapper).
func (s *pseudoAttributeStrategy) GetPseudo(certificate *diagnostic.CertificateWrapper) string {
	return certificate.Pseudo()
}
