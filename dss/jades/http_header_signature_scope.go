// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/HTTPHeaderSignatureScope.java (DSS 6.5.RC1).
//
// The signature scope used to define the signed payload with the HTTPHeader SigD Mechanism.
package jades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
)

// HTTPHeaderSignatureScope embeds scope.SignatureScopeBase, the Go counterpart of Java's
// abstract SignatureScope class.
type HTTPHeaderSignatureScope struct {
	scope.SignatureScopeBase
}

var _ scope.SignatureScope = (*HTTPHeaderSignatureScope)(nil)

// NewHTTPHeaderSignatureScope is the default constructor, taking the DSSDocument representing
// the computed JWS Payload.
func NewHTTPHeaderSignatureScope(document model.DSSDocument) *HTTPHeaderSignatureScope {
	name := "HttpHeaders payload"
	if document.Name() != "" {
		name = document.Name()
	}
	return &HTTPHeaderSignatureScope{SignatureScopeBase: scope.NewSignatureScopeBaseWithName(name, document)}
}

// Description returns the HTTPHeaderSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *HTTPHeaderSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "Payload value digest"
}

// Type returns the type of the signature scope. Port of getType().
func (s *HTTPHeaderSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}
