// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCContentDocumentFilter.java (DSS 6.5.RC1).
package asic

import (
	"slices"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ContentDocumentFilter provides a configuration to filter the content of an ASiC
// container. Example: the class can be used to define type of documents to be protected by an
// evidence record (see EvidenceRecordDigestBuilder). See ASiCContentDocumentFilterFactory
// for creating pre-configured instances for basic usages.
type ContentDocumentFilter struct {
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

// NewASiCContentDocumentFilter instantiates an ContentDocumentFilter object with an empty
// configuration. Ports the default constructor.
func NewASiCContentDocumentFilter() *ContentDocumentFilter {
	return &ContentDocumentFilter{}
}

// SetMimetypeDocument sets whether the mimetype document present at the root level shall be
// returned. Ports setMimetypeDocument(boolean).
func (f *ContentDocumentFilter) SetMimetypeDocument(mimetypeDocument bool) {
	f.mimetypeDocument = mimetypeDocument
}

// SetSignedDocuments sets whether the original signed documents shall be returned. Ports
// setSignedDocuments(boolean).
func (f *ContentDocumentFilter) SetSignedDocuments(signedDocuments bool) {
	f.signedDocuments = signedDocuments
}

// SetSignatureDocuments sets whether the signature documents present within a META-INF/
// folder shall be returned. Ports setSignatureDocuments(boolean).
func (f *ContentDocumentFilter) SetSignatureDocuments(signatureDocuments bool) {
	f.signatureDocuments = signatureDocuments
}

// SetTimestampDocuments sets whether the time-stamp documents present within a META-INF/
// folder shall be returned. Ports setTimestampDocuments(boolean).
func (f *ContentDocumentFilter) SetTimestampDocuments(timestampDocuments bool) {
	f.timestampDocuments = timestampDocuments
}

// SetEvidenceRecordDocuments sets whether the evidence record documents present within a
// META-INF/ folder shall be returned. Ports setEvidenceRecordDocuments(boolean).
func (f *ContentDocumentFilter) SetEvidenceRecordDocuments(evidenceRecordDocuments bool) {
	f.evidenceRecordDocuments = evidenceRecordDocuments
}

// SetManifestDocuments sets whether the ASiC manifest documents present within a META-INF/
// folder shall be returned. Ports setManifestDocuments(boolean).
func (f *ContentDocumentFilter) SetManifestDocuments(manifestDocuments bool) {
	f.manifestDocuments = manifestDocuments
}

// SetArchiveManifestDocuments sets whether the archive ASiC manifest documents present within
// a META-INF/ folder shall be returned. Ports setArchiveManifestDocuments(boolean).
func (f *ContentDocumentFilter) SetArchiveManifestDocuments(archiveManifestDocuments bool) {
	f.archiveManifestDocuments = archiveManifestDocuments
}

// SetEvidenceRecordManifestDocuments sets whether the evidence record manifest documents
// present within a META-INF/ folder shall be returned. Ports
// setEvidenceRecordManifestDocuments(boolean).
func (f *ContentDocumentFilter) SetEvidenceRecordManifestDocuments(evidenceRecordManifestDocuments bool) {
	f.evidenceRecordManifestDocuments = evidenceRecordManifestDocuments
}

// SetUnsupportedDocuments sets whether other documents not directly supported by
// EN 319 162-1 specification shall be returned. Ports setUnsupportedDocuments(boolean).
func (f *ContentDocumentFilter) SetUnsupportedDocuments(unsupportedDocuments bool) {
	f.unsupportedDocuments = unsupportedDocuments
}

// SetExcludedFilenames sets a collection of document filenames to be excluded from the final
// return result. Ports setExcludedFilenames(Collection).
func (f *ContentDocumentFilter) SetExcludedFilenames(excludedFilenames []string) {
	f.excludedFilenames = excludedFilenames
}

// SetIncludedFilenames sets a collection of document filenames to be included in the final
// return result despite other settings. NOTE: take precedence over all other constraints.
// Ports setIncludedFilenames(Collection).
func (f *ContentDocumentFilter) SetIncludedFilenames(includedFilenames []string) {
	f.includedFilenames = includedFilenames
}

// Filter returns a list of filtered DSSDocuments according to the configuration. Ports
// filter(Content).
func (f *ContentDocumentFilter) Filter(asicContent *Content) []model.DSSDocument {
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
func (f *ContentDocumentFilter) filterDocuments(documents []model.DSSDocument, formatSupported bool) []model.DSSDocument {
	result := make([]model.DSSDocument, 0)
	if utils.IsCollectionNotEmpty(f.includedFilenames) {
		for _, d := range documents {
			if slices.Contains(f.includedFilenames, d.Name()) {
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
		if !slices.Contains(result, d) && !slices.Contains(f.excludedFilenames, d.Name()) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}
