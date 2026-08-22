// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/ContainerContentSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
)

// ContainerContentSignatureScope represents a content of a package.zip container file for an
// ASiC-S container.
type ContainerContentSignatureScope struct {
	ContainerSignatureScope
}

// NewContainerContentSignatureScope is the default constructor. Port of
// ContainerContentSignatureScope(DSSDocument).
func NewContainerContentSignatureScope(document model.DSSDocument) *ContainerContentSignatureScope {
	return &ContainerContentSignatureScope{ContainerSignatureScope: *NewContainerSignatureScope(document)}
}

// Description returns the ContainerContentSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *ContainerContentSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "ASiCS archive content"
}

// Type returns the type of the signature scope. Port of getType().
func (s *ContainerContentSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeTypeArchived
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*ContainerContentSignatureScope)(nil)
