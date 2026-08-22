// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/ManifestSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
)

// ManifestSignatureScope represents a Manifest document.
//
// Java declares this class final; there is no Go language-level equivalent, so it is
// documented here instead.
type ManifestSignatureScope struct {
	SignatureScopeWithTransformationsBase
}

// NewManifestSignatureScope is the default constructor. Port of
// ManifestSignatureScope(ManifestFile).
func NewManifestSignatureScope(manifestFile *model.ManifestFile) *ManifestSignatureScope {
	return NewManifestSignatureScopeWithTransformations(manifestFile, nil)
}

// NewManifestSignatureScopeWithTransformations is the constructor with a list of XML
// transformations (to be used for XAdES only). Port of ManifestSignatureScope(ManifestFile,
// List<String>).
func NewManifestSignatureScopeWithTransformations(manifestFile *model.ManifestFile, transformations []string) *ManifestSignatureScope {
	return NewManifestSignatureScopeWithNameAndTransformations(manifestFile.Filename(), manifestFile.Document(), transformations)
}

// NewManifestSignatureScopeWithNameAndTransformations is the constructor with an explicit
// filename, document, and list of XML transformations (to be used for XAdES only). Port of
// ManifestSignatureScope(String, DSSDocument, List<String>).
func NewManifestSignatureScopeWithNameAndTransformations(filename string, document model.DSSDocument, transformations []string) *ManifestSignatureScope {
	return &ManifestSignatureScope{
		SignatureScopeWithTransformationsBase: NewSignatureScopeWithTransformationsBaseWithName(filename, document, transformations),
	}
}

// Description returns the ManifestSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *ManifestSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return "Manifest document"
}

// Type returns the type of the signature scope. Port of getType().
func (s *ManifestSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*ManifestSignatureScope)(nil)
