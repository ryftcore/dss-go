// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/ASiCContentDocumentFilterFactory.java (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/utils"

// EmptyFilter creates an ASiCContentDocumentFilter with an empty configuration. Ports the
// static factory emptyFilter(). Java's private no-arg constructor has no Go analogue: this
// file's factory functions are free package functions, not methods on a value type.
func EmptyFilter() *ASiCContentDocumentFilter {
	return NewASiCContentDocumentFilter()
}

// SignedDocumentsOnlyFilter creates an ASiCContentDocumentFilter with a configuration to
// return only original signed documents. Note: this is a default configuration for an
// ASiCManifest (signed or time-stamped). Ports signedDocumentsOnlyFilter(String...).
func SignedDocumentsOnlyFilter(excludedFilenames ...string) *ASiCContentDocumentFilter {
	f := NewASiCContentDocumentFilter()
	f.SetSignedDocuments(true)
	if utils.IsArrayNotEmpty(excludedFilenames) {
		f.SetExcludedFilenames(excludedFilenames)
	}
	return f
}

// ArchiveDocumentsFilter creates an ASiCContentDocumentFilter with a configuration returning
// original signed documents, signature and time-stamp documents, as well as the corresponding
// manifest files. Note: this is a default configuration for an ASiCArchiveManifest. Ports
// archiveDocumentsFilter(String...).
func ArchiveDocumentsFilter(excludedFilenames ...string) *ASiCContentDocumentFilter {
	f := NewASiCContentDocumentFilter()
	f.SetSignedDocuments(true)
	f.SetSignatureDocuments(true)
	f.SetTimestampDocuments(true)
	f.SetManifestDocuments(true)
	f.SetArchiveManifestDocuments(true)
	if utils.IsArrayNotEmpty(excludedFilenames) {
		f.SetExcludedFilenames(excludedFilenames)
	}
	return f
}

// AllSupportedDocumentsFilter creates an ASiCContentDocumentFilter returning all recognized
// type documents available (without mimetype), within an ASiCContent, excluding the documents
// with filenames defined in excludedFilenames. Ports allSupportedDocumentsFilter(String...).
func AllSupportedDocumentsFilter(excludedFilenames ...string) *ASiCContentDocumentFilter {
	f := NewASiCContentDocumentFilter()
	f.SetSignedDocuments(true)
	f.SetSignatureDocuments(true)
	f.SetTimestampDocuments(true)
	f.SetEvidenceRecordDocuments(true)
	f.SetManifestDocuments(true)
	f.SetArchiveManifestDocuments(true)
	f.SetEvidenceRecordManifestDocuments(true)
	if utils.IsArrayNotEmpty(excludedFilenames) {
		f.SetExcludedFilenames(excludedFilenames)
	}
	return f
}

// AllDocumentsFilter creates an ASiCContentDocumentFilter returning all documents available,
// including unrecognized documents and mimetype document, within an ASiCContent, excluding the
// documents with filenames defined in excludedFilenames. Ports allDocumentsFilter(String...).
func AllDocumentsFilter(excludedFilenames ...string) *ASiCContentDocumentFilter {
	f := NewASiCContentDocumentFilter()
	f.SetMimetypeDocument(true)
	f.SetSignedDocuments(true)
	f.SetSignatureDocuments(true)
	f.SetTimestampDocuments(true)
	f.SetEvidenceRecordDocuments(true)
	f.SetManifestDocuments(true)
	f.SetArchiveManifestDocuments(true)
	f.SetEvidenceRecordManifestDocuments(true)
	f.SetUnsupportedDocuments(true)
	if utils.IsArrayNotEmpty(excludedFilenames) {
		f.SetExcludedFilenames(excludedFilenames)
	}
	return f
}

// AllowedFilenamesFilter creates an ASiCContentDocumentFilter returning all documents
// available, matching the array of allowedFilenames. Ignores all other documents. Ports
// allowedFilenamesFilter(String...).
func AllowedFilenamesFilter(allowedFilenames ...string) *ASiCContentDocumentFilter {
	f := NewASiCContentDocumentFilter()
	if utils.IsArrayNotEmpty(allowedFilenames) {
		f.SetIncludedFilenames(allowedFilenames)
	}
	return f
}
