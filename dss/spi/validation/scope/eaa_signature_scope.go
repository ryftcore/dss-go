// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/EAASignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EAASignatureScope is used to refer a signature used to issue EAA.
type EAASignatureScope struct {
	scope.SignatureScopeBase

	// EAA is the counter-signed parent signature.
	EAA validation.EAA
}

// NewEAASignatureScope is the default constructor. Port of EAASignatureScope(EAA,
// DSSDocument). Panics if eaa is nil, mirroring Java's Objects.requireNonNull(eaa, "EAA
// cannot be null!").
func NewEAASignatureScope(eaa validation.EAA, originalDocument model.DSSDocument) *EAASignatureScope {
	if eaa == nil {
		panic("EAA cannot be null!")
	}
	return &EAASignatureScope{
		SignatureScopeBase: scope.NewSignatureScopeBase(originalDocument),
		EAA:                eaa,
	}
}

// Name returns the EAA presentation id. Port of getName(TokenIdentifierProvider).
func (s *EAASignatureScope) Name(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.eaaPresentationID(tokenIdentifierProvider)
}

// Description returns the EAASignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *EAASignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("EAA with Id : %s", s.eaaPresentationID(tokenIdentifierProvider))
}

func (s *EAASignatureScope) eaaPresentationID(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return tokenIdentifierProvider.IDAsString(s.EAA)
}

// Type returns the type of the signature scope. Port of getType().
func (s *EAASignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_EAA_SIGNATURE
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*EAASignatureScope)(nil)
