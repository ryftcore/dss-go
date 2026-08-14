// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/OpenDocumentSupportUtils.java (DSS 6.5.RC1).
package xades

import (
	"strings"

	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
)

// openDocumentSupportUtilsExternalData defines the external data directory name.
const openDocumentSupportUtilsExternalData = "external-data/"

// OpenDocumentSupportUtilsGetOpenDocumentCoverage returns the list of covered documents.
//
// # ODF 1.2 ch 3.16
//
// An OpenDocument document that is stored in a package may have one or more digital signatures
// applied to the package.
//
// Document signatures shall be stored in a file called META-INF/documentsignatures.xml in the
// package as described in section 3.5 of the OpenDocument specification part 3. Document
// signatures shall contain a <ds:Reference> element for each file within the package, with the
// exception that <ds:Reference> elements for the META-INF/documentsignatures.xml file containing
// the signature, and any files contained in the package whose relative path starts with
// "external-data/" should be omitted.
//
// Ports the static getOpenDocumentCoverage(ASiCContent).
func OpenDocumentSupportUtilsGetOpenDocumentCoverage(asicContent *asic.ASiCContent) []model.DSSDocument {
	docs := make([]model.DSSDocument, 0)
	docs = append(docs, asicContent.SignedDocuments()...)
	docs = append(docs, asicContent.ManifestDocuments()...)
	docs = append(docs, asicContent.ArchiveManifestDocuments()...)
	docs = append(docs, asicContent.TimestampDocuments()...)
	docs = append(docs, asicContent.UnsupportedDocuments()...)
	docs = append(docs, asicContent.MimeTypeDocument())

	result := make([]model.DSSDocument, 0, len(docs))
	for _, doc := range docs {
		if !OpenDocumentSupportUtilsIsExternalDataDocument(doc) {
			result = append(result, doc)
		}
	}

	return result
}

// OpenDocumentSupportUtilsIsExternalDataDocument verifies whether the document is located within
// "external-data/" folder within the archive. Ports the static isExternalDataDocument(DSSDocument).
//
// Matches Java's null-check scope exactly: only document.getName() is null-checked, so a nil
// document panics here (nil interface method call) the same way Java's NPE would.
func OpenDocumentSupportUtilsIsExternalDataDocument(document model.DSSDocument) bool {
	return document.Name() != "" && strings.HasPrefix(document.Name(), openDocumentSupportUtilsExternalData)
}
