// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/KeyBindingSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// KeyBindingSignatureScope is used to refer a key binding signature's signed content.
type KeyBindingSignatureScope struct {
	scope.SignatureScopeBase

	// EAA is the counter-signed parent signature.
	EAA validation.EAA
}

// NewKeyBindingSignatureScope is the default constructor. Port of
// KeyBindingSignatureScope(EAA, DSSDocument). Panics if eaa is nil, mirroring Java's
// Objects.requireNonNull(eaa, "EAA cannot be null!").
func NewKeyBindingSignatureScope(eaa validation.EAA, originalDocument model.DSSDocument) *KeyBindingSignatureScope {
	if eaa == nil {
		panic("EAA cannot be null!")
	}
	return &KeyBindingSignatureScope{
		SignatureScopeBase: scope.NewSignatureScopeBase(originalDocument),
		EAA:                eaa,
	}
}

// Name returns the EAA presentation id. Port of getName(TokenIdentifierProvider).
func (s *KeyBindingSignatureScope) Name(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.eaaPresentationID(tokenIdentifierProvider)
}

// Description returns the KeyBindingSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *KeyBindingSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("Key binding for EAA with Id : %s", s.eaaPresentationID(tokenIdentifierProvider))
}

func (s *KeyBindingSignatureScope) eaaPresentationID(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return tokenIdentifierProvider.IDAsString(s.EAA)
}

// Type returns the type of the signature scope. Port of getType().
func (s *KeyBindingSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeTypeKeyBindingSignature
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*KeyBindingSignatureScope)(nil)
