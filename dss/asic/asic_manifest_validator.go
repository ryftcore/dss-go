// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/ASiCManifestValidator.java (DSS 6.5.RC1).
package asic

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ManifestValidator performs validation of an ASiC Manifest entries.
type ManifestValidator struct {
	// manifest is the manifest to validate.
	manifest *model.ManifestFile

	// signedDocuments is a list of documents covered by the manifest.
	signedDocuments []model.DSSDocument
}

// NewASiCManifestValidator is the default constructor. Ports
// ManifestValidator(ManifestFile, List).
//
// Panics with the Java message when manifest is nil (Objects.requireNonNull).
func NewASiCManifestValidator(manifest *model.ManifestFile, signedDocuments []model.DSSDocument) *ManifestValidator {
	if manifest == nil {
		panic("ManifestFile must be defined!")
	}
	return &ManifestValidator{
		manifest:        manifest,
		signedDocuments: signedDocuments,
	}
}

// ValidateEntries validates the manifest entries and returns the list of validated
// ManifestEntrys. Ports validateEntries(). Logging (LOG.warn) is dropped per PORTING.md.
func (v *ManifestValidator) ValidateEntries() []*model.ManifestEntry {
	manifestEntries := v.manifest.Entries()
	if utils.IsCollectionEmpty(v.signedDocuments) {
		// no signed data to validate on
		return manifestEntries
	}
	for _, entry := range manifestEntries {
		digest := entry.Digest()
		if !digest.IsEmpty() {
			// Use strict by name handling, as document names are predefined within an ASiC container
			signedDocument := spi.DSSUtilsDocumentWithName(v.signedDocuments, entry.Uri())
			if signedDocument != nil {
				entry.SetFound(true)
				entry.SetDocument(signedDocument)
				computedDigest, err := signedDocument.DigestValue(digest.Algorithm())
				if err == nil && bytes.Equal(digest.Value(), computedDigest) {
					entry.SetIntact(true)
				}
			}
		}
	}

	return manifestEntries
}
