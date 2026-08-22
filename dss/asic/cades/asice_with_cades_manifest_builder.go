// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/manifest/ASiCEWithCAdESManifestBuilder.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.cades.signature.manifest lands in
// this same Go package (dss/asic/cades).
//
// Java's `extends AbstractASiCManifestBuilder` becomes embedding plus the
// InitAbstractASiCManifestBuilderWithDigestAlgorithm(self, ...) registration, so the base
// dispatches getSigReferenceMimeType() / initDefaultAsicContentDocumentFilter() /
// getManifestFilename() / isRootfile() into the concrete builder instead of into itself. Since
// this class is abstract in Java, the registration happens in the leaf constructors (signature /
// timestamp manifest builders).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCEWithCAdESManifestBuilder generates the ASiCManifest.xml content (ASiC-E).
//
// Sample:
//
//	<asic:ASiCManifest xmlns:asic="http://uri.etsi.org/02918/v1.2.1#">
//		<asic:SigReference MimeType="application/pkcs7-signature" URI="META-INF/signature001.p7s">
//			<asic:DataObjectReference URI="document.txt">
//				<DigestMethod xmlns="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//				<DigestValue xmlns="http://www.w3.org/2000/09/xmldsig#">OuL0HMJE899y+uJtyNnTt5B/gFrrw8adNczI+9w9GDQ=</DigestValue>
//			</asic:DataObjectReference>
//		</asic:SigReference>
//	</asic:ASiCManifest>
type ASiCEWithCAdESManifestBuilder struct {
	asic.AbstractASiCManifestBuilder

	// asicFilenameFactory defines rules for filename creation for new manifest files.
	asicFilenameFactory ASiCWithCAdESFilenameFactory
}

// initASiCEWithCAdESManifestBuilder applies the protected
// ASiCEWithCAdESManifestBuilder(Content, String, DigestAlgorithm) constructor, which
// delegates to the filename-factory one with a DefaultASiCWithCAdESFilenameFactory. The leaf
// builder passes itself as overrides.
func (b *ASiCEWithCAdESManifestBuilder) initASiCEWithCAdESManifestBuilder(
	overrides asic.AbstractASiCManifestBuilderOverrides, asicContent *asic.Content,
	documentFilename string, digestAlgorithm enumerations.DigestAlgorithm) {
	b.initASiCEWithCAdESManifestBuilderWithFilenameFactory(overrides, asicContent, documentFilename,
		digestAlgorithm, NewDefaultASiCWithCAdESFilenameFactory())
}

// initASiCEWithCAdESManifestBuilderWithFilenameFactory applies the protected
// ASiCEWithCAdESManifestBuilder(Content, String, DigestAlgorithm, ASiCWithCAdESFilenameFactory)
// constructor.
func (b *ASiCEWithCAdESManifestBuilder) initASiCEWithCAdESManifestBuilderWithFilenameFactory(
	overrides asic.AbstractASiCManifestBuilderOverrides, asicContent *asic.Content,
	documentFilename string, digestAlgorithm enumerations.DigestAlgorithm,
	asicFilenameFactory ASiCWithCAdESFilenameFactory) {
	b.InitAbstractASiCManifestBuilderWithDigestAlgorithm(overrides, asicContent, documentFilename, digestAlgorithm)
	b.asicFilenameFactory = asicFilenameFactory
}

// InitDefaultAsicContentDocumentFilter ports the @Override protected
// initDefaultAsicContentDocumentFilter().
func (b *ASiCEWithCAdESManifestBuilder) InitDefaultAsicContentDocumentFilter() *asic.ContentDocumentFilter {
	return asic.SignedDocumentsOnlyFilter()
}

// ManifestFilename ports the @Override protected getManifestFilename().
func (b *ASiCEWithCAdESManifestBuilder) ManifestFilename() string {
	return b.asicFilenameFactory.ManifestFilename(b.AsicContent)
}
