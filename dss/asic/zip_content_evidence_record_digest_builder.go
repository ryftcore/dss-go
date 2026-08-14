// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ZipContentEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
package asic

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/x509/evidencerecord/digest"
	"github.com/utain/esig/dss/utils"
)

// ZipContentEvidenceRecordDigestBuilder builds hashes for all documents present within a ZIP
// archive. Note: for covering an ASiC container with an evidence record, please use
// ASiCEvidenceRecordDigestBuilder.
type ZipContentEvidenceRecordDigestBuilder struct {
	// documents is the list of documents to compute hashes for.
	documents []model.DSSDocument

	// DigestAlgorithm is the digest algorithm to be used on hash computation. Default:
	// DigestAlgorithm_SHA256.
	DigestAlgorithm enumerations.DigestAlgorithm

	// DataObjectDigestBuilderFactory is the factory to be used to instantiate a new
	// DataObjectDigestBuilder for hashes computation.
	DataObjectDigestBuilderFactory digest.DataObjectDigestBuilderFactory
}

// newZipContentEvidenceRecordDigestBuilderBase ports the two protected constructors
// (no-arg, delegating to SHA256, and the single-DigestAlgorithm form): both leave documents
// nil, for use by subclasses (e.g. ASiCEvidenceRecordDigestBuilder) that supply their own
// document list via an overridden BuildDigestGroup.
func newZipContentEvidenceRecordDigestBuilderBase(digestAlgorithm enumerations.DigestAlgorithm) *ZipContentEvidenceRecordDigestBuilder {
	return &ZipContentEvidenceRecordDigestBuilder{
		documents:       nil,
		DigestAlgorithm: digestAlgorithm,
	}
}

// NewZipContentEvidenceRecordDigestBuilder creates a ZipContentEvidenceRecordDigestBuilder to
// build hashes from a DSSDocument, represented by a ZIP container, using a default SHA-256
// digest algorithm. Ports ZipContentEvidenceRecordDigestBuilder(DSSDocument).
func NewZipContentEvidenceRecordDigestBuilder(zipContainer model.DSSDocument) *ZipContentEvidenceRecordDigestBuilder {
	return NewZipContentEvidenceRecordDigestBuilderAlgo(zipContainer, enumerations.DigestAlgorithm_SHA256)
}

// NewZipContentEvidenceRecordDigestBuilderAlgo creates a ZipContentEvidenceRecordDigestBuilder
// to build hashes with the provided DigestAlgorithm from a DSSDocument, represented by a ZIP
// container. Ports ZipContentEvidenceRecordDigestBuilder(DSSDocument, DigestAlgorithm).
func NewZipContentEvidenceRecordDigestBuilderAlgo(zipContainer model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *ZipContentEvidenceRecordDigestBuilder {
	return &ZipContentEvidenceRecordDigestBuilder{
		documents:       extractZipContentDocuments(zipContainer),
		DigestAlgorithm: digestAlgorithm,
	}
}

// extractZipContentDocuments ports the private static extractDocuments(DSSDocument). Panics on
// an extraction error, matching the panic-on-Build-error convention used throughout this port.
func extractZipContentDocuments(zipContainer model.DSSDocument) []model.DSSDocument {
	documents, err := ZipUtilsInstance().ExtractContainerContent(zipContainer)
	if err != nil {
		panic(err)
	}
	return documents
}

// SetDataObjectDigestBuilderFactory sets a factory to instantiate a new
// DataObjectDigestBuilder for hashes computation of the given evidence record type (e.g.
// XMLERS or ASN.1 ERS). Ports setDataObjectDigestBuilderFactory(DataObjectDigestBuilderFactory).
func (b *ZipContentEvidenceRecordDigestBuilder) SetDataObjectDigestBuilderFactory(dataObjectDigestBuilderFactory digest.DataObjectDigestBuilderFactory) *ZipContentEvidenceRecordDigestBuilder {
	b.DataObjectDigestBuilderFactory = dataObjectDigestBuilderFactory
	return b
}

// BuildDigestGroup builds a list of hashes for the content files of the provided ZIP
// container. Ports buildDigestGroup(). Overridden by ASiCEvidenceRecordDigestBuilder; per the
// virtual-dispatch warning in S7_BRIEF.md, this base method is only invoked directly on a
// *ZipContentEvidenceRecordDigestBuilder value (never through an overrides interface), since
// no code in this package holds a ZipContentEvidenceRecordDigestBuilder-typed reference to an
// embedded ASiCEvidenceRecordDigestBuilder and expects override dispatch.
func (b *ZipContentEvidenceRecordDigestBuilder) BuildDigestGroup() []model.Digest {
	b.AssertConfigurationValid()
	return b.ComputeDigestForDocuments(b.documents)
}

// AssertConfigurationValid verifies whether the configuration of the current builder is
// valid. Ports the protected assertConfigurationValid().
func (b *ZipContentEvidenceRecordDigestBuilder) AssertConfigurationValid() {
	if b.DataObjectDigestBuilderFactory == nil {
		panic("DataObjectDigestBuilderFactory shall be set to continue! " +
			"Please choose the corresponding implementation for your evidence record type (e.g. XMLERS or ASN.1).")
	}
}

// ComputeDigestForDocuments computes a list of digests for the given list of DSSDocuments.
// Ports the protected computeDigestForDocuments(List).
func (b *ZipContentEvidenceRecordDigestBuilder) ComputeDigestForDocuments(documents []model.DSSDocument) []model.Digest {
	if utils.IsCollectionEmpty(documents) {
		return []model.Digest{}
	}

	result := make([]model.Digest, 0, len(documents))
	for _, document := range documents {
		dataObjectDigestBuilder := b.DataObjectDigestBuilderFactory.CreateWithAlgorithm(document, b.DigestAlgorithm)
		d := dataObjectDigestBuilder.Build()
		result = append(result, *d)
	}
	return result
}
