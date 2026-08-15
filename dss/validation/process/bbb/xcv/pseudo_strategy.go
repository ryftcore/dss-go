// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/pseudo/PseudoStrategy.java (DSS 6.5.RC1).
package xcv

import "github.com/utain/esig/dss/diagnostic"

// PseudoStrategy is a strategy to extract a pseudo string from a given
// certificate.
type PseudoStrategy interface {
	// GetPseudo gets pseudo for the certificate. Port of getPseudo(CertificateWrapper).
	GetPseudo(certificate *diagnostic.CertificateWrapper) string
}
