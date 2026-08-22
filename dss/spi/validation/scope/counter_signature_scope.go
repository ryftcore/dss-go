// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/CounterSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CounterSignatureScope is used to refer a counter-signed SignatureValue.
type CounterSignatureScope struct {
	scope.SignatureScopeBase

	// MasterSignature is the counter-signed parent signature.
	MasterSignature validation.AdvancedSignature
}

// NewCounterSignatureScope is the default constructor. Port of
// CounterSignatureScope(AdvancedSignature, DSSDocument). Panics if masterSignature is nil,
// mirroring Java's Objects.requireNonNull(masterSignature, "Master signature cannot be
// null!").
func NewCounterSignatureScope(masterSignature validation.AdvancedSignature, originalDocument model.DSSDocument) *CounterSignatureScope {
	if masterSignature == nil {
		panic("Master signature cannot be null!")
	}
	return &CounterSignatureScope{
		SignatureScopeBase: scope.NewSignatureScopeBase(originalDocument),
		MasterSignature:    masterSignature,
	}
}

// Name returns the master signature id. Port of getName(TokenIdentifierProvider).
func (s *CounterSignatureScope) Name(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.masterSignatureID(tokenIdentifierProvider)
}

// Description returns the CounterSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *CounterSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("Master signature with Id : %s", s.masterSignatureID(tokenIdentifierProvider))
}

func (s *CounterSignatureScope) masterSignatureID(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return tokenIdentifierProvider.IDAsString(s.MasterSignature)
}

// Type returns the type of the signature scope. Port of getType().
func (s *CounterSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeTypeCounterSignature
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*CounterSignatureScope)(nil)
