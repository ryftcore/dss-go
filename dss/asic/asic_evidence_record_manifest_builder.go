// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCEvidenceRecordManifestBuilder.java (DSS 6.5.RC1).
package asic

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
)

// ASiCEvidenceRecordManifestBuilder builds an ASiCManifest for an Evidence Record. Ports the
// Java class, which extends AbstractASiCManifestBuilder (ZIPCORE chunk, same asic package).
//
// AbstractASiCManifestBuilder is embedded by value (Go embedding, not inheritance) and
// initialised via InitAbstractASiCManifestBuilder(overrides, asicContent, sigReferenceUri); its
// Overrides interface methods are named SigReferenceMimeType/ManifestFilename (no "Get" prefix)
// and AsicContent is an exported field, not a method.
type ASiCEvidenceRecordManifestBuilder struct {
	AbstractASiCManifestBuilder

	// evidenceRecordFilenameFactory defines rules for filename creation for new manifest
	// files.
	evidenceRecordFilenameFactory ASiCEvidenceRecordFilenameFactory
}

// NewASiCEvidenceRecordManifestBuilderFromDocument builds a manifest from a DSSDocument
// representing the ASiC container. Ports
// ASiCEvidenceRecordManifestBuilder(DSSDocument, DigestAlgorithm, String).
func NewASiCEvidenceRecordManifestBuilderFromDocument(asicContainer model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm, evidenceRecordFilename string) *ASiCEvidenceRecordManifestBuilder {
	return NewASiCEvidenceRecordManifestBuilder(asicEvidenceRecordManifestToASiCContent(asicContainer), digestAlgorithm, evidenceRecordFilename)
}

// asicEvidenceRecordManifestToASiCContent ports the private static toASiCContent(DSSDocument).
func asicEvidenceRecordManifestToASiCContent(asicContainer model.DSSDocument) *ASiCContent {
	extractor, err := DefaultASiCContainerExtractorFromDocument(asicContainer)
	if err == nil {
		var content *ASiCContent
		content, err = extractor.Extract()
		if err == nil {
			return content
		}
	}
	panic(exception.NewIllegalInputExceptionWithCause(
		fmt.Sprintf("Unsupported ASiC or document type! Returned error : %s", err.Error()), err))
}

// NewASiCEvidenceRecordManifestBuilder builds a manifest from ASiCContent representing the
// ASiC container. Ports ASiCEvidenceRecordManifestBuilder(ASiCContent, DigestAlgorithm,
// String).
func NewASiCEvidenceRecordManifestBuilder(asicContent *ASiCContent, digestAlgorithm enumerations.DigestAlgorithm, evidenceRecordFilename string) *ASiCEvidenceRecordManifestBuilder {
	b := &ASiCEvidenceRecordManifestBuilder{}
	b.InitAbstractASiCManifestBuilderWithDigestAlgorithm(b, asicContent, evidenceRecordFilename, digestAlgorithm)
	return b
}

// SigReferenceMimeType ports the @Override protected getSigReferenceMimeType(): not required
// for an evidence record.
func (b *ASiCEvidenceRecordManifestBuilder) SigReferenceMimeType() enumerations.MimeType {
	return nil
}

// InitDefaultAsicContentDocumentFilter ports the @Override protected
// initDefaultAsicContentDocumentFilter().
func (b *ASiCEvidenceRecordManifestBuilder) InitDefaultAsicContentDocumentFilter() *ASiCContentDocumentFilter {
	return ArchiveDocumentsFilter()
}

// SetAsicContentDocumentFilter ports the @Override covariant-return
// setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *ASiCEvidenceRecordManifestBuilder) SetAsicContentDocumentFilter(asicContentDocumentFilter *ASiCContentDocumentFilter) *ASiCEvidenceRecordManifestBuilder {
	b.AbstractASiCManifestBuilder.SetAsicContentDocumentFilter(asicContentDocumentFilter)
	return b
}

// SetEvidenceRecordFilenameFactory sets an ASiC evidence record filename factory, used to
// provide a valid filename for the ASiC Evidence Record Manifest document to be created.
// Note: when not set, final DSSDocument will have name set to "" (Java: NULL). Ports
// setEvidenceRecordFilenameFactory(ASiCEvidenceRecordFilenameFactory).
func (b *ASiCEvidenceRecordManifestBuilder) SetEvidenceRecordFilenameFactory(evidenceRecordFilenameFactory ASiCEvidenceRecordFilenameFactory) *ASiCEvidenceRecordManifestBuilder {
	b.evidenceRecordFilenameFactory = evidenceRecordFilenameFactory
	return b
}

// ManifestFilename ports the @Override protected getManifestFilename().
func (b *ASiCEvidenceRecordManifestBuilder) ManifestFilename() string {
	if b.evidenceRecordFilenameFactory != nil {
		return b.evidenceRecordFilenameFactory.EvidenceRecordManifestFilename(b.AsicContent)
	}
	return ""
}
