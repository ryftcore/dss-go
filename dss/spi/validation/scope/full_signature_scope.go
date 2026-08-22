// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/FullSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
)

// FullSignatureScope is a marker type meaning that the signature covers the whole document.
//
// Java declares this class final; there is no Go language-level equivalent, so it is
// documented here instead - no other file in this package embeds FullSignatureScope.
type FullSignatureScope struct {
	scope.SignatureScopeBase
}

// NewFullSignatureScope is the default constructor. Port of FullSignatureScope(String,
// DSSDocument).
func NewFullSignatureScope(filename string, document model.DSSDocument) *FullSignatureScope {
	return &FullSignatureScope{SignatureScopeBase: scope.NewSignatureScopeBaseWithName(filename, document)}
}

// Description returns the FullSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *FullSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "Full document"
}

// Type returns the type of the signature scope. Port of getType().
func (s *FullSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*FullSignatureScope)(nil)
