// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/TokenIssuerSelector.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// TokenIssuerSelector is used to select an issuer of the provided Token.
type TokenIssuerSelector struct {
	// token is the token to get an issuer for.
	token model.Token

	// certificateTokens is the collection of issuer candidates.
	certificateTokens []*model.CertificateToken
}

// NewTokenIssuerSelector builds the selector. Port of the default constructor.
func NewTokenIssuerSelector(token model.Token, certificateTokens []*model.CertificateToken) *TokenIssuerSelector {
	return &TokenIssuerSelector{token: token, certificateTokens: certificateTokens}
}

// Issuer filters certificateTokens and returns the best issuer candidate for token, or nil if
// none was found. Port of getIssuer().
func (s *TokenIssuerSelector) Issuer() *model.CertificateToken {
	if len(s.certificateTokens) == 0 {
		return nil
	}
	candidates := tokenIssuerSelectorFilterByPublicKey(s.token, s.certificateTokens)
	candidates = tokenIssuerSelectorFilterByIssuerSubjectName(s.token, candidates)

	for _, candidate := range candidates {
		if candidate.IsValidOn(s.token.CreationDate()) {
			return candidate
		}
	}
	if len(candidates) > 0 {
		// No matching issuer found for the token creation date; the process continues with an
		// issuer which has the same public key.
		return candidates[0]
	}
	return nil
}

// tokenIssuerSelectorFilterByPublicKey filters a collection of candidates having a public key
// matching the one used to sign token. Port of the private static filterIssuersByPublicKey.
func tokenIssuerSelectorFilterByPublicKey(token model.Token, candidates []*model.CertificateToken) []*model.CertificateToken {
	var issuers []*model.CertificateToken
	for _, candidate := range candidates {
		if token.IsSignedByToken(candidate) {
			issuers = append(issuers, candidate)
		}
	}
	return issuers
}

// tokenIssuerSelectorFilterByIssuerSubjectName filters a collection of candidates with a
// Subject name matching the Issuer Subject name of token. When none of the candidates match
// the Issuer Subject name, the original collection of candidates is returned.
// Port of the private static filterIssuersByIssuerSubjectName.
func tokenIssuerSelectorFilterByIssuerSubjectName(token model.Token, candidates []*model.CertificateToken) []*model.CertificateToken {
	var issuers []*model.CertificateToken
	issuerX500Principal := token.IssuerX500Principal()
	if issuerX500Principal != nil {
		issuerX500PrincipalHelper := model.NewX500PrincipalHelper(issuerX500Principal)
		for _, candidate := range candidates {
			if issuerX500PrincipalHelper.Equals(candidate.Subject()) {
				issuers = append(issuers, candidate)
			}
			// else: the issuer subject name and subject name do not match (upstream logs a
			// warning here; dropped, see PORTING.md).
		}
	}
	if len(issuers) > 0 {
		return issuers
	}
	return candidates
}
