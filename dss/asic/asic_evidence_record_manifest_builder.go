// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCEvidenceRecordManifestBuilder.java (DSS 6.5.RC1).
package asic

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// ASiCEvidenceRecordManifestBuilder builds an ASiCManifest for an Evidence Record. Ports the
// Java class, which extends AbstractASiCManifestBuilder (ZIPCORE chunk, same asic package).
//
// AbstractASiCManifestBuilder is embedded by value (Go embedding, not inheritance) and
// initialised via InitAbstractASiCManifestBuilder(overrides, asicContent, sigReferenceUri); its
// Overrides interface methods are named SigReferenceMimeType/ManifestFilename (no "Get" prefix)
// and AsicContent is an exported field, not a method.
type EvidenceRecordManifestBuilder struct {
	AbstractASiCManifestBuilder

	// evidenceRecordFilenameFactory defines rules for filename creation for new manifest
	// files.
	evidenceRecordFilenameFactory EvidenceRecordFilenameFactory
}

// NewEvidenceRecordManifestBuilderFromDocument builds a manifest from a DSSDocument
// representing the ASiC container. Ports
// EvidenceRecordManifestBuilder(DSSDocument, DigestAlgorithm, String).
func NewEvidenceRecordManifestBuilderFromDocument(asicContainer model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm, evidenceRecordFilename string) *EvidenceRecordManifestBuilder {
	return NewEvidenceRecordManifestBuilder(asicEvidenceRecordManifestToASiCContent(asicContainer), digestAlgorithm, evidenceRecordFilename)
}

// asicEvidenceRecordManifestToASiCContent ports the private static toASiCContent(DSSDocument).
func asicEvidenceRecordManifestToASiCContent(asicContainer model.DSSDocument) *Content {
	extractor, err := DefaultContainerExtractorFromDocument(asicContainer)
	if err == nil {
		var content *Content
		content, err = extractor.Extract()
		if err == nil {
			return content
		}
	}
	panic(exception.NewIllegalInputExceptionWithCause(
		fmt.Sprintf("Unsupported ASiC or document type! Returned error : %s", err.Error()), err))
}

// NewEvidenceRecordManifestBuilder builds a manifest from Content representing the
// ASiC container. Ports ASiCEvidenceRecordManifestBuilder(ASiCContent, DigestAlgorithm,
// String).
func NewEvidenceRecordManifestBuilder(asicContent *Content, digestAlgorithm enumerations.DigestAlgorithm, evidenceRecordFilename string) *EvidenceRecordManifestBuilder {
	b := &EvidenceRecordManifestBuilder{}
	b.InitAbstractASiCManifestBuilderWithDigestAlgorithm(b, asicContent, evidenceRecordFilename, digestAlgorithm)
	return b
}

// SigReferenceMimeType ports the @Override protected getSigReferenceMimeType(): not required
// for an evidence record.
func (b *EvidenceRecordManifestBuilder) SigReferenceMimeType() enumerations.MimeType {
	return nil
}

// InitDefaultAsicContentDocumentFilter ports the @Override protected
// initDefaultAsicContentDocumentFilter().
func (b *EvidenceRecordManifestBuilder) InitDefaultAsicContentDocumentFilter() *ContentDocumentFilter {
	return ArchiveDocumentsFilter()
}

// SetAsicContentDocumentFilter ports the @Override covariant-return
// setAsicContentDocumentFilter(ContentDocumentFilter).
func (b *EvidenceRecordManifestBuilder) SetAsicContentDocumentFilter(asicContentDocumentFilter *ContentDocumentFilter) *EvidenceRecordManifestBuilder {
	b.AbstractASiCManifestBuilder.SetAsicContentDocumentFilter(asicContentDocumentFilter)
	return b
}

// SetEvidenceRecordFilenameFactory sets an ASiC evidence record filename factory, used to
// provide a valid filename for the ASiC Evidence Record Manifest document to be created.
// Note: when not set, final DSSDocument will have name set to "" (Java: NULL). Ports
// setEvidenceRecordFilenameFactory(EvidenceRecordFilenameFactory).
func (b *EvidenceRecordManifestBuilder) SetEvidenceRecordFilenameFactory(evidenceRecordFilenameFactory EvidenceRecordFilenameFactory) *EvidenceRecordManifestBuilder {
	b.evidenceRecordFilenameFactory = evidenceRecordFilenameFactory
	return b
}

// ManifestFilename ports the @Override protected getManifestFilename().
func (b *EvidenceRecordManifestBuilder) ManifestFilename() string {
	if b.evidenceRecordFilenameFactory != nil {
		return b.evidenceRecordFilenameFactory.EvidenceRecordManifestFilename(b.AsicContent)
	}
	return ""
}
