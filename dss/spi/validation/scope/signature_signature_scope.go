// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/SignatureSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// SignatureSignatureScope defines a signature scope covering a signature.
type SignatureSignatureScope struct {
	scope.SignatureScopeBase

	// signature is the covered signature.
	signature validation.AdvancedSignature
}

// NewSignatureSignatureScope is the default constructor to create a signature scope. Port of
// SignatureSignatureScope(AdvancedSignature, DSSDocument). Panics if signature is nil,
// mirroring Java's Objects.requireNonNull(signature, "Signature shall be provided!").
//
// document represents the covered signature document (NOTE: not necessary to be a signature
// file).
func NewSignatureSignatureScope(signature validation.AdvancedSignature, document model.DSSDocument) *SignatureSignatureScope {
	if signature == nil {
		panic("Signature shall be provided!")
	}
	return &SignatureSignatureScope{
		SignatureScopeBase: scope.NewSignatureScopeBaseWithName(signature.ID(), document),
		signature:          signature,
	}
}

// DSSID returns the unique DSS Identifier. Port of getDSSId(), a no-op override in Java
// (the method body simply returns super.getDSSId()); kept here to mirror the source
// structure even though it adds no behaviour over the embedded SignatureScopeBase.DSSID().
func (s *SignatureSignatureScope) DSSID() model.Identifier {
	return s.SignatureScopeBase.DSSID()
}

// Name returns the covered signature's id. Port of getName(TokenIdentifierProvider).
func (s *SignatureSignatureScope) Name(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.signatureID(tokenIdentifierProvider)
}

// Description returns the SignatureSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *SignatureSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("Signature with Id : %s", s.signatureID(tokenIdentifierProvider))
}

func (s *SignatureSignatureScope) signatureID(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return tokenIdentifierProvider.IDAsString(s.signature)
}

// Type returns the type of the signature scope. Port of getType().
func (s *SignatureSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_SIGNATURE
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*SignatureSignatureScope)(nil)
