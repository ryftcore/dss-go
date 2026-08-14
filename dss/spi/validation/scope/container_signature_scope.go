// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/ContainerSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi"
)

// ContainerSignatureScope represents a "package.zip" container for ASiC-S signatures.
type ContainerSignatureScope struct {
	scope.SignatureScopeBase
}

// NewContainerSignatureScope is the default constructor, deriving the filename from the
// document's URI-decoded name. Port of ContainerSignatureScope(DSSDocument).
func NewContainerSignatureScope(document model.DSSDocument) *ContainerSignatureScope {
	return NewContainerSignatureScopeWithName(spi.DSSUtilsDecodeURI(document.Name()), document)
}

// NewContainerSignatureScopeWithName is the constructor with an explicit filename. Port of
// ContainerSignatureScope(String, DSSDocument).
func NewContainerSignatureScopeWithName(filename string, document model.DSSDocument) *ContainerSignatureScope {
	return &ContainerSignatureScope{SignatureScopeBase: scope.NewSignatureScopeBaseWithName(filename, document)}
}

// Description returns the ContainerSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *ContainerSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "ASiCS archive"
}

// Type returns the type of the signature scope. Port of getType().
func (s *ContainerSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*ContainerSignatureScope)(nil)
