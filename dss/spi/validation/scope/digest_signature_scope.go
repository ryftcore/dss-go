// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/DigestSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
)

// DigestSignatureScope defines a Digest document scope.
type DigestSignatureScope struct {
	scope.SignatureScopeBase
}

// NewDigestSignatureScope is the default constructor with a filename. Port of
// DigestSignatureScope(String, DSSDocument).
func NewDigestSignatureScope(filename string, digestDocument model.DSSDocument) *DigestSignatureScope {
	return &DigestSignatureScope{SignatureScopeBase: scope.NewSignatureScopeBaseWithName(filename, digestDocument)}
}

// Description returns the DigestSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *DigestSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "Digest of the document content"
}

// Type returns the type of the signature scope. Port of getType().
func (s *DigestSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_DIGEST
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*DigestSignatureScope)(nil)
