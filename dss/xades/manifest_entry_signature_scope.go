// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/ManifestEntrySignatureScope.java (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	spiscope "github.com/utain/esig/dss/spi/validation/scope"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// ManifestEntrySignatureScope represents a Manifest entry document. Port of the class
// ManifestEntrySignatureScope, extending spiscope.SignatureScopeWithTransformationsBase.
type ManifestEntrySignatureScope struct {
	spiscope.SignatureScopeWithTransformationsBase

	// manifestName is the name of the manifest containing the entry.
	manifestName string
}

// NewManifestEntrySignatureScope is the constructor with transformations (used in XAdES). Port
// of the public ManifestEntrySignatureScope(String, DSSDocument, String, List<String>)
// constructor.
func NewManifestEntrySignatureScope(entryName string, document model.DSSDocument, manifestName string,
	transformations []string) *ManifestEntrySignatureScope {
	return &ManifestEntrySignatureScope{
		SignatureScopeWithTransformationsBase: spiscope.NewSignatureScopeWithTransformationsBaseWithName(entryName, document, transformations),
		manifestName:                          manifestName,
	}
}

// Description returns the ManifestEntrySignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *ManifestEntrySignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	var description string
	if xmlutils.DomUtilsIsElementReference(s.DocumentName()) {
		description = fmt.Sprintf("The XML Manifest Entry with ID '%s' from a Manifest with name '%s'", s.DocumentName(), s.manifestName)
	} else {
		description = fmt.Sprintf("The File Manifest Entry with name '%s' from a Manifest with name '%s'", s.DocumentName(), s.manifestName)
	}
	return s.AddTransformationIfNeeded(description)
}

// Type returns the type of the signature scope. Port of getType().
func (s *ManifestEntrySignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// String ports toString().
func (s *ManifestEntrySignatureScope) String() string {
	return fmt.Sprintf("ManifestEntrySignatureScope{manifestName='%s'} %s", s.manifestName, s.SignatureScopeWithTransformationsBase.String())
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*ManifestEntrySignatureScope)(nil)
