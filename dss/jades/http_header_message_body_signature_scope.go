// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/HTTPHeaderMessageBodySignatureScope.java (DSS 6.5.RC1).
//
// The signature scope used to define the HTTPHeader message body.
package jades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
)

// HTTPHeaderMessageBodySignatureScope extends HTTPHeaderSignatureScope (embedding, since Go has
// no class inheritance).
type HTTPHeaderMessageBodySignatureScope struct {
	HTTPHeaderSignatureScope
}

var _ scope.SignatureScope = (*HTTPHeaderMessageBodySignatureScope)(nil)

// NewHTTPHeaderMessageBodySignatureScope is the default constructor, taking the DSSDocument
// message body content.
func NewHTTPHeaderMessageBodySignatureScope(document model.DSSDocument) *HTTPHeaderMessageBodySignatureScope {
	return &HTTPHeaderMessageBodySignatureScope{HTTPHeaderSignatureScope: *NewHTTPHeaderSignatureScope(document)}
}

// Description returns the HTTPHeaderMessageBodySignatureScope description, overriding the
// embedded HTTPHeaderSignatureScope's. Port of getDescription(TokenIdentifierProvider).
func (s *HTTPHeaderMessageBodySignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "Message body value digest"
}
