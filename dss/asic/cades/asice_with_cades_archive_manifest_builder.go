// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/manifest/ASiCEWithCAdESArchiveManifestBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// ASiCEWithCAdESArchiveManifestBuilder generates the ASiCArchiveManifest.xml content (ASiC-E).
//
// Sample:
//
//	<asic:ASiCManifest xmlns:asic="http://uri.etsi.org/02918/v1.2.1#">
//		<asic:SigReference URI="META-INF/archive_timestamp.tst" MimeType="application/vnd.etsi.timestamp-token"/>
//		<asic:DataObjectReference URI="META-INF/signature.p7s" MimeType="application/x-pkcs7-signature">
//			<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//			<DigestValue>3Qeos8...</DigestValue>
//		</asic:DataObjectReference>
//		<asic:DataObjectReference URI="toBeSigned.txt" MimeType="text/plain">
//			<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//			<DigestValue>JJZt...</DigestValue>
//		</asic:DataObjectReference>
//		<asic:DataObjectReference URI="META-INF/ASiCManifest_1.xml" MimeType="text/xml">
//			<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//			<DigestValue>g5dY...</DigestValue>
//		</asic:DataObjectReference>
//	</asic:ASiCManifest>
type ASiCEWithCAdESArchiveManifestBuilder struct {
	asic.AbstractASiCManifestBuilder

	// lastArchiveManifest is the "ASiCArchiveManifest.xml" document (root manifest).
	lastArchiveManifest model.DSSDocument
}

var _ asic.AbstractASiCManifestBuilderOverrides = (*ASiCEWithCAdESArchiveManifestBuilder)(nil)

// NewASiCEWithCAdESArchiveManifestBuilder is the default constructor. Ports
// ASiCEWithCAdESArchiveManifestBuilder(ASiCContent, DSSDocument, DigestAlgorithm, String).
func NewASiCEWithCAdESArchiveManifestBuilder(asicContent *asic.ASiCContent, lastArchiveManifest model.DSSDocument,
	digestAlgorithm enumerations.DigestAlgorithm, timestampFilename string) *ASiCEWithCAdESArchiveManifestBuilder {
	builder := &ASiCEWithCAdESArchiveManifestBuilder{lastArchiveManifest: lastArchiveManifest}
	builder.InitAbstractASiCManifestBuilderWithDigestAlgorithm(builder, asicContent, timestampFilename, digestAlgorithm)
	return builder
}

// IsRootfile ports the @Override protected isRootfile(DSSDocument). Java compares by reference
// (`lastArchiveManifest == document`); Go's interface comparison is the same reference test for
// the pointer-backed DSSDocument implementations this port uses.
func (b *ASiCEWithCAdESArchiveManifestBuilder) IsRootfile(document model.DSSDocument) bool {
	return b.lastArchiveManifest == document
}

// SigReferenceMimeType ports the @Override protected getSigReferenceMimeType().
func (b *ASiCEWithCAdESArchiveManifestBuilder) SigReferenceMimeType() enumerations.MimeType {
	return enumerations.MimeTypeEnum_TST
}

// InitDefaultAsicContentDocumentFilter ports the @Override protected
// initDefaultAsicContentDocumentFilter().
func (b *ASiCEWithCAdESArchiveManifestBuilder) InitDefaultAsicContentDocumentFilter() *asic.ASiCContentDocumentFilter {
	return asic.ArchiveDocumentsFilter()
}

// SetAsicContentDocumentFilter ports the @Override covariant-return
// setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *ASiCEWithCAdESArchiveManifestBuilder) SetAsicContentDocumentFilter(
	asicContentDocumentFilter *asic.ASiCContentDocumentFilter) *ASiCEWithCAdESArchiveManifestBuilder {
	b.AbstractASiCManifestBuilder.SetAsicContentDocumentFilter(asicContentDocumentFilter)
	return b
}

// ManifestFilename ports the @Override protected getManifestFilename().
func (b *ASiCEWithCAdESArchiveManifestBuilder) ManifestFilename() string {
	return ASiCWithCAdESUtilsDefaultArchiveManifestFilename
}
