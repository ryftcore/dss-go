// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCContentDocumentFilter.java (DSS 6.5.RC1).
package asic

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCContentDocumentFilter provides a configuration to filter the content of an ASiC
// container. Example: the class can be used to define type of documents to be protected by an
// evidence record (see ASiCEvidenceRecordDigestBuilder). See ASiCContentDocumentFilterFactory
// for creating pre-configured instances for basic usages.
type ASiCContentDocumentFilter struct {
	// mimetypeDocument defines whether the mimetype document within a root folder of the
	// container shall be returned.
	mimetypeDocument bool

	// signedDocuments defines whether the original signed documents from the container shall
	// be returned.
	signedDocuments bool

	// signatureDocuments defines whether the signature documents within a META-INF/ folder of
	// the container shall be returned.
	signatureDocuments bool

	// timestampDocuments defines whether the time-stamp documents within a META-INF/ folder
	// of the container shall be returned.
	timestampDocuments bool

	// evidenceRecordDocuments defines whether the evidence record documents within a
	// META-INF/ folder of the container shall be returned.
	evidenceRecordDocuments bool

	// manifestDocuments defines whether the manifest documents within a META-INF/ folder of
	// the container shall be returned.
	manifestDocuments bool

	// archiveManifestDocuments defines whether the archive manifest documents within a
	// META-INF/ folder of the container shall be returned.
	archiveManifestDocuments bool

	// evidenceRecordManifestDocuments defines whether the evidence record manifest documents
	// within a META-INF/ folder of the container shall be returned.
	evidenceRecordManifestDocuments bool

	// unsupportedDocuments defines whether other documents not directly supported by
	// EN 319 162-1 specification shall be returned.
	unsupportedDocuments bool

	// excludedFilenames contains a collection of filenames (including the directory path) to
	// be excluded from the returned result.
	excludedFilenames []string

	// includedFilenames contains a collection of filenames (including the directory path) to
	// be including despite the other settings.
	includedFilenames []string
}

// NewASiCContentDocumentFilter instantiates an ASiCContentDocumentFilter object with an empty
// configuration. Ports the default constructor.
func NewASiCContentDocumentFilter() *ASiCContentDocumentFilter {
	return &ASiCContentDocumentFilter{}
}

// SetMimetypeDocument sets whether the mimetype document present at the root level shall be
// returned. Ports setMimetypeDocument(boolean).
func (f *ASiCContentDocumentFilter) SetMimetypeDocument(mimetypeDocument bool) {
	f.mimetypeDocument = mimetypeDocument
}

// SetSignedDocuments sets whether the original signed documents shall be returned. Ports
// setSignedDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetSignedDocuments(signedDocuments bool) {
	f.signedDocuments = signedDocuments
}

// SetSignatureDocuments sets whether the signature documents present within a META-INF/
// folder shall be returned. Ports setSignatureDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetSignatureDocuments(signatureDocuments bool) {
	f.signatureDocuments = signatureDocuments
}

// SetTimestampDocuments sets whether the time-stamp documents present within a META-INF/
// folder shall be returned. Ports setTimestampDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetTimestampDocuments(timestampDocuments bool) {
	f.timestampDocuments = timestampDocuments
}

// SetEvidenceRecordDocuments sets whether the evidence record documents present within a
// META-INF/ folder shall be returned. Ports setEvidenceRecordDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetEvidenceRecordDocuments(evidenceRecordDocuments bool) {
	f.evidenceRecordDocuments = evidenceRecordDocuments
}

// SetManifestDocuments sets whether the ASiC manifest documents present within a META-INF/
// folder shall be returned. Ports setManifestDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetManifestDocuments(manifestDocuments bool) {
	f.manifestDocuments = manifestDocuments
}

// SetArchiveManifestDocuments sets whether the archive ASiC manifest documents present within
// a META-INF/ folder shall be returned. Ports setArchiveManifestDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetArchiveManifestDocuments(archiveManifestDocuments bool) {
	f.archiveManifestDocuments = archiveManifestDocuments
}

// SetEvidenceRecordManifestDocuments sets whether the evidence record manifest documents
// present within a META-INF/ folder shall be returned. Ports
// setEvidenceRecordManifestDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetEvidenceRecordManifestDocuments(evidenceRecordManifestDocuments bool) {
	f.evidenceRecordManifestDocuments = evidenceRecordManifestDocuments
}

// SetUnsupportedDocuments sets whether other documents not directly supported by
// EN 319 162-1 specification shall be returned. Ports setUnsupportedDocuments(boolean).
func (f *ASiCContentDocumentFilter) SetUnsupportedDocuments(unsupportedDocuments bool) {
	f.unsupportedDocuments = unsupportedDocuments
}

// SetExcludedFilenames sets a collection of document filenames to be excluded from the final
// return result. Ports setExcludedFilenames(Collection).
func (f *ASiCContentDocumentFilter) SetExcludedFilenames(excludedFilenames []string) {
	f.excludedFilenames = excludedFilenames
}

// SetIncludedFilenames sets a collection of document filenames to be included in the final
// return result despite other settings. NOTE: take precedence over all other constraints.
// Ports setIncludedFilenames(Collection).
func (f *ASiCContentDocumentFilter) SetIncludedFilenames(includedFilenames []string) {
	f.includedFilenames = includedFilenames
}

// Filter returns a list of filtered DSSDocuments according to the configuration. Ports
// filter(ASiCContent).
func (f *ASiCContentDocumentFilter) Filter(asicContent *ASiCContent) []model.DSSDocument {
	result := make([]model.DSSDocument, 0)
	if asicContent.MimeTypeDocument() != nil {
		result = append(result, f.filterDocuments([]model.DSSDocument{asicContent.MimeTypeDocument()}, f.mimetypeDocument)...)
	}
	result = append(result, f.filterDocuments(asicContent.SignedDocuments(), f.signedDocuments)...)
	result = append(result, f.filterDocuments(asicContent.SignatureDocuments(), f.signatureDocuments)...)
	result = append(result, f.filterDocuments(asicContent.TimestampDocuments(), f.timestampDocuments)...)
	result = append(result, f.filterDocuments(asicContent.EvidenceRecordDocuments(), f.evidenceRecordDocuments)...)
	result = append(result, f.filterDocuments(asicContent.ManifestDocuments(), f.manifestDocuments)...)
	result = append(result, f.filterDocuments(asicContent.ArchiveManifestDocuments(), f.archiveManifestDocuments)...)
	result = append(result, f.filterDocuments(asicContent.EvidenceRecordManifestDocuments(), f.evidenceRecordManifestDocuments)...)
	result = append(result, f.filterDocuments(asicContent.UnsupportedDocuments(), f.unsupportedDocuments)...)
	return result
}

// filterDocuments ports the private filterDocuments(Collection, boolean).
func (f *ASiCContentDocumentFilter) filterDocuments(documents []model.DSSDocument, formatSupported bool) []model.DSSDocument {
	result := make([]model.DSSDocument, 0)
	if utils.IsCollectionNotEmpty(f.includedFilenames) {
		for _, d := range documents {
			if containsString(f.includedFilenames, d.Name()) {
				result = append(result, d)
			}
		}
	}
	if !formatSupported {
		return result
	}
	if utils.IsCollectionEmpty(f.excludedFilenames) {
		return documents
	}
	filtered := make([]model.DSSDocument, 0)
	for _, d := range documents {
		if !containsDocument(result, d) && !containsString(f.excludedFilenames, d.Name()) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// containsString ports the Collection.contains(Object) used against a String collection.
func containsString(strs []string, s string) bool {
	for _, v := range strs {
		if v == s {
			return true
		}
	}
	return false
}

// containsDocument ports the List.contains(Object) used against a List<DSSDocument>, relying
// on Go interface equality (comparable pointer/value identity), matching Java's default
// Object.equals reference-identity semantics used here.
func containsDocument(documents []model.DSSDocument, d model.DSSDocument) bool {
	for _, v := range documents {
		if v == d {
			return true
		}
	}
	return false
}
