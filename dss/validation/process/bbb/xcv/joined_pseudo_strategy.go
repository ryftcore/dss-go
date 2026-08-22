// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/pseudo/JoinedPseudoStrategy.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JoinedPseudoStrategy represents a merged strategy to extract pseudo
// information, accepting the certificate's pseudo attribute and custom
// German pseudo processing algorithm.
type JoinedPseudoStrategy struct {
	// strategies is the list of strategies to be used to extract pseudo.
	strategies []PseudoStrategy
}

// NewJoinedPseudoStrategy is the default constructor.
func NewJoinedPseudoStrategy() *JoinedPseudoStrategy {
	return &JoinedPseudoStrategy{
		strategies: []PseudoStrategy{
			newPseudoAttributeStrategy(),
			NewPseudoGermanyStrategy(),
		},
	}
}

// GetPseudo gets pseudo for the certificate. Port of getPseudo(CertificateWrapper).
func (s *JoinedPseudoStrategy) GetPseudo(certificate *diagnostic.CertificateWrapper) string {
	for _, strategy := range s.strategies {
		pseudo := strategy.GetPseudo(certificate)
		if utils.IsStringNotEmpty(pseudo) {
			return pseudo
		}
	}
	return ""
}
