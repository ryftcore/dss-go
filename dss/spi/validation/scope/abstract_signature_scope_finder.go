// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/AbstractSignatureScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// AbstractSignatureScopeFinder is the base for SignatureScope finding.
type AbstractSignatureScopeFinder struct {
}

// NewAbstractSignatureScopeFinder is the default constructor instantiating the object with
// default values.
func NewAbstractSignatureScopeFinder() AbstractSignatureScopeFinder {
	return AbstractSignatureScopeFinder{}
}

// IsASiCSArchive checks if the given signature represents an ASiC-S container. Port of
// isASiCSArchive(AdvancedSignature).
func (f *AbstractSignatureScopeFinder) IsASiCSArchive(advancedSignature validation.AdvancedSignature) bool {
	return utils.IsCollectionNotEmpty(advancedSignature.ContainerContents())
}

// IsASiCEArchive checks if the given signature represents an ASiC-E container. Port of
// isASiCEArchive(AdvancedSignature).
func (f *AbstractSignatureScopeFinder) IsASiCEArchive(advancedSignature validation.AdvancedSignature) bool {
	return advancedSignature.ManifestFile() != nil
}

// CreateInMemoryDocument creates a DSSDocument from given binaries. Port of
// createInMemoryDocument(byte[]).
func (f *AbstractSignatureScopeFinder) CreateInMemoryDocument(binaries []byte) model.DSSDocument {
	return model.NewInMemoryDocument(binaries)
}

// CreateDigestDocument creates a DSSDocument from given digest. Port of
// createDigestDocument(Digest).
func (f *AbstractSignatureScopeFinder) CreateDigestDocument(digest model.Digest) model.DSSDocument {
	if digest.Algorithm() != "" && digest.Value() != nil {
		return model.NewDigestDocumentFromBase64(digest.Algorithm(), utils.ToBase64(digest.Value()))
	}
	return nil
}
