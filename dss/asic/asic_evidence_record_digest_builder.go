// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
package asic

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/x509/evidencerecord/digest"
)

// ASiCEvidenceRecordDigestBuilder is used to build hashes for data objects within an ASiC
// container for potential evidence-record incorporation. Ports the Java class, which extends
// ZipContentEvidenceRecordDigestBuilder; Go embeds it. This type fully reimplements
// BuildDigestGroup and AssertConfigurationValid (rather than relying on the embedded base's
// dispatch), so the virtual-dispatch base-calls-overridden-subclass-method hazard flagged in
// S7_BRIEF.md does not arise here: BuildDigestGroup below calls b.AssertConfigurationValid()
// on the receiver directly (Go static dispatch resolves to this type's own method, matching
// what Java's dynamic dispatch would resolve to for a call site typed as
// ASiCEvidenceRecordDigestBuilder).
type ASiCEvidenceRecordDigestBuilder struct {
	*ZipContentEvidenceRecordDigestBuilder

	// asicContent is the content of an ASiC container.
	asicContent *ASiCContent

	// asicContentDocumentFilter is used to filter the documents to compute hashes for.
	asicContentDocumentFilter *ASiCContentDocumentFilter
}

// NewASiCEvidenceRecordDigestBuilderFromDocument creates a ASiCEvidenceRecordDigestBuilder to
// build hashes from a DSSDocument, represented by an ASiC container, using a default SHA-256
// digest algorithm. Ports ASiCEvidenceRecordDigestBuilder(DSSDocument), which may panic with
// an *exception.IllegalInputException-wrapping error surfaced via panic(error) if the document
// is not a supported ASiC or document type - see toASiCContent.
func NewASiCEvidenceRecordDigestBuilderFromDocument(asicContainer model.DSSDocument) *ASiCEvidenceRecordDigestBuilder {
	return NewASiCEvidenceRecordDigestBuilderFromDocumentWithAlgorithm(asicContainer, enumerations.DigestAlgorithm_SHA256)
}

// NewASiCEvidenceRecordDigestBuilderFromDocumentWithAlgorithm creates a
// ASiCEvidenceRecordDigestBuilder to build hashes with the provided DigestAlgorithm from a
// DSSDocument, represented by an ASiC container. Ports
// ASiCEvidenceRecordDigestBuilder(DSSDocument, DigestAlgorithm).
func NewASiCEvidenceRecordDigestBuilderFromDocumentWithAlgorithm(asicContainer model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *ASiCEvidenceRecordDigestBuilder {
	return NewASiCEvidenceRecordDigestBuilderWithAlgorithm(asicEvidenceRecordToASiCContent(asicContainer), digestAlgorithm)
}

// asicEvidenceRecordToASiCContent ports the private static toASiCContent(DSSDocument).
//
// Cross-chunk assumption (ZIPCORE): DefaultASiCContainerExtractor is expected to expose a
// package-level constructor DefaultASiCContainerExtractorFromDocument(model.DSSDocument)
// (*DefaultASiCContainerExtractor, error), mirroring
// DefaultASiCContainerExtractor.fromDocument(DSSDocument), and an Extract() (*ASiCContent,
// error) method on the ASiCContainerExtractor interface.
func asicEvidenceRecordToASiCContent(asicContainer model.DSSDocument) *ASiCContent {
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

// NewASiCEvidenceRecordDigestBuilder creates a ASiCEvidenceRecordDigestBuilder to build hashes
// from ASiCContent, using a default SHA-256 digest algorithm. Ports
// ASiCEvidenceRecordDigestBuilder(ASiCContent).
func NewASiCEvidenceRecordDigestBuilder(asicContent *ASiCContent) *ASiCEvidenceRecordDigestBuilder {
	return NewASiCEvidenceRecordDigestBuilderWithAlgorithm(asicContent, enumerations.DigestAlgorithm_SHA256)
}

// NewASiCEvidenceRecordDigestBuilderWithAlgorithm creates a ASiCEvidenceRecordDigestBuilder to
// build hashes with the provided DigestAlgorithm from ASiCContent. Ports
// ASiCEvidenceRecordDigestBuilder(ASiCContent, DigestAlgorithm).
func NewASiCEvidenceRecordDigestBuilderWithAlgorithm(asicContent *ASiCContent, digestAlgorithm enumerations.DigestAlgorithm) *ASiCEvidenceRecordDigestBuilder {
	return &ASiCEvidenceRecordDigestBuilder{
		ZipContentEvidenceRecordDigestBuilder: newZipContentEvidenceRecordDigestBuilderBase(digestAlgorithm),
		asicContent:                           asicContent,
	}
}

// SetDataObjectDigestBuilderFactory sets a factory to instantiate a new
// DataObjectDigestBuilder for hashes computation of the given evidence record type (e.g.
// XMLERS or ASN.1 ERS). Ports the @Override setDataObjectDigestBuilderFactory(...), which
// covariantly returns ASiCEvidenceRecordDigestBuilder.
func (b *ASiCEvidenceRecordDigestBuilder) SetDataObjectDigestBuilderFactory(dataObjectDigestBuilderFactory digest.DataObjectDigestBuilderFactory) *ASiCEvidenceRecordDigestBuilder {
	b.ZipContentEvidenceRecordDigestBuilder.SetDataObjectDigestBuilderFactory(dataObjectDigestBuilderFactory)
	return b
}

// SetAsicContentDocumentFilter sets an ASiCContentDocumentFilter used to filter the documents
// to compute hashes for. Ports setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *ASiCEvidenceRecordDigestBuilder) SetAsicContentDocumentFilter(asicContentDocumentFilter *ASiCContentDocumentFilter) *ASiCEvidenceRecordDigestBuilder {
	b.asicContentDocumentFilter = asicContentDocumentFilter
	return b
}

// BuildDigestGroup ports the @Override buildDigestGroup(). Shadows (does not call) the
// embedded ZipContentEvidenceRecordDigestBuilder.BuildDigestGroup.
func (b *ASiCEvidenceRecordDigestBuilder) BuildDigestGroup() []model.Digest {
	b.AssertConfigurationValid()

	documents := b.GetDocumentListToComputeDigest()
	return b.ComputeDigestForDocuments(documents)
}

// AssertConfigurationValid ports the @Override protected assertConfigurationValid(), which
// calls super.assertConfigurationValid() then adds its own check.
func (b *ASiCEvidenceRecordDigestBuilder) AssertConfigurationValid() {
	b.ZipContentEvidenceRecordDigestBuilder.AssertConfigurationValid()
	if b.asicContentDocumentFilter == nil {
		panic("ASiCContentDocumentFilter shall be set to continue! " +
			"Use ASiCContentDocumentFilterFactory to facilitate configuration.")
	}
}

// GetDocumentListToComputeDigest executes an ASiCContentDocumentFilter and returns a list of
// documents to compute hashes for. Ports the protected getDocumentListToComputeDigest().
func (b *ASiCEvidenceRecordDigestBuilder) GetDocumentListToComputeDigest() []model.DSSDocument {
	return b.asicContentDocumentFilter.Filter(b.asicContent)
}
