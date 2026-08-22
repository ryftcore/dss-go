// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/OCSPFirstRevocationDataLoadingStrategy.java (DSS 6.5.RC1).
//
// This is the concrete strategy already forward-declared and consumed by
// ocsp_first_revocation_data_loading_strategy_factory.go (landed earlier in this same chunk),
// which expects OCSPFirstRevocationDataLoadingStrategy to embed RevocationDataLoadingStrategy
// exactly like CRLFirstRevocationDataLoadingStrategy does (see
// crl_first_revocation_data_loading_strategy.go), so that Create() can return the embedded
// *RevocationDataLoadingStrategy uniformly. RevocationDataLoadingStrategy itself (the base
// abstract class, with its checkCRL/checkOCSP/isAcceptableToken/InitRevocationDataLoadingStrategy
// members) is owned by sibling chunk VAL-D and is not yet landed in this package; it is used
// here exactly as already assumed by crl_first_revocation_data_loading_strategy.go's header
// comment.
//
// slf4j LOG.debug call in the Java source (fired only when neither OCSP nor CRL responded) is
// dropped per phase 2a handoff fact "slf4j dropped unless load-bearing" - it carries no
// behavior, matching the same drop already made for CRLFirstRevocationDataLoadingStrategy's
// analogous branch.
package validation

import "github.com/ryftcore/dss-go/dss/model"

// OCSPFirstRevocationDataLoadingStrategy fetches firstly the OCSP token response, and if not
// available, tries CRL, returning the first succeeded result.
//
// NOTE: This implementation is used by default for revocation retrieving.
type OCSPFirstRevocationDataLoadingStrategy struct {
	RevocationDataLoadingStrategy
}

// NewOCSPFirstRevocationDataLoadingStrategy is the default constructor.
func NewOCSPFirstRevocationDataLoadingStrategy() *OCSPFirstRevocationDataLoadingStrategy {
	s := &OCSPFirstRevocationDataLoadingStrategy{}
	s.InitRevocationDataLoadingStrategy(s)
	return s
}

// RevocationToken retrieves the revocation token for the given certificate, trying OCSP first
// and falling back to CRL. Port of getRevocationToken(CertificateToken, CertificateToken).
func (s *OCSPFirstRevocationDataLoadingStrategy) RevocationToken(certificateToken, issuerToken *model.CertificateToken) AnyRevocationToken {
	ocspToken := s.checkOCSP(certificateToken, issuerToken)
	if ocspToken != nil && s.isAcceptableToken(ocspToken) {
		return ocspToken
	}
	crlToken := s.checkCRL(certificateToken, issuerToken)
	if crlToken != nil && s.isAcceptableToken(crlToken) {
		return crlToken
	}
	if s.fallbackEnabled {
		// return first successful result
		if ocspToken != nil {
			return ocspToken
		}
		if crlToken != nil {
			return crlToken
		}
	}
	return nil
}
