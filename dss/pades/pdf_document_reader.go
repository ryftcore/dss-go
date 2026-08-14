// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfDocumentReader.java (DSS 6.5.RC1).
//
// java.io.Closeable becomes a Close() error method (io.Closer). The two screenshot methods have
// no Go counterpart: rasterisation is an explicit non-goal of the native PDF engine
// (internal/pdf/DESIGN.md §0.2), so they are declared here to keep the SPI recognisable and the
// native implementation returns ErrRasterisationNotSupported for both. java.awt.image.BufferedImage
// therefore becomes model.DSSDocument (the PNG container the callers ultimately want).
//
// Map<PdfSignatureDictionary, List<PdfSignatureField>> becomes an ordered slice of pairs: Java's
// LinkedHashMap iteration order is part of the contract (AbstractPDFSignatureService sorts it
// with PdfSignatureDictionaryComparator), and a Go map has no order.
package pades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// PdfSignatureDictionaryFields pairs a signature dictionary with the signature fields that refer
// to it. It is one entry of Java's LinkedHashMap<PdfSignatureDictionary, List<PdfSignatureField>>.
type PdfSignatureDictionaryFields struct {
	// SignatureDictionary is the map key.
	SignatureDictionary *PdfSignatureDictionary

	// Fields are the signature fields whose /V points at SignatureDictionary.
	Fields []*PdfSignatureField
}

// PdfDocumentReader reads a PDF document.
type PdfDocumentReader interface {
	// Close releases the resources held by the reader. Port of Closeable#close.
	Close() error

	// DSSDictionary loads the last DSS dictionary from the document, nil when not present.
	// Port of #getDSSDictionary.
	DSSDictionary() PdfDssDict

	// ExtractSigDictionaries extracts the PdfSignatureDictionaries present in the document,
	// paired with the related signature fields, in document order.
	// Port of #extractSigDictionaries.
	ExtractSigDictionaries() ([]PdfSignatureDictionaryFields, error)

	// IsSignatureCoversWholeDocument checks whether the signature for the given signature
	// dictionary covers the whole document. Port of #isSignatureCoversWholeDocument.
	IsSignatureCoversWholeDocument(signatureDictionary *PdfSignatureDictionary) bool

	// NumberOfPages returns the number of pages found in the document.
	// Port of #getNumberOfPages.
	NumberOfPages() int

	// PageBox returns the dimensions of the given (1-based) page. Port of #getPageBox.
	PageBox(page int) AnnotationBox

	// PageRotation returns the rotation of the given (1-based) page, in degrees.
	// Port of #getPageRotation.
	PageRotation(page int) int

	// PdfAnnotations retrieves the annotations associated with the given (1-based) page.
	// Port of #getPdfAnnotations.
	PdfAnnotations(page int) ([]*PdfAnnotation, error)

	// GenerateImageScreenshot generates the image screenshot for the given page.
	// Port of #generateImageScreenshot; not supported by the native engine.
	GenerateImageScreenshot(page int) (model.DSSDocument, error)

	// GenerateImageScreenshotWithoutAnnotations generates the image screenshot for the given
	// page, hiding the given annotations.
	// Port of #generateImageScreenshotWithoutAnnotations; not supported by the native engine.
	GenerateImageScreenshotWithoutAnnotations(page int, addedAnnotations []*PdfAnnotation) (model.DSSDocument, error)

	// IsEncrypted checks whether the document is encrypted. Port of #isEncrypted.
	IsEncrypted() bool

	// IsOpenWithOwnerAccess verifies whether the document has been opened with full owner access
	// (all modifications are permitted). Port of #isOpenWithOwnerAccess.
	IsOpenWithOwnerAccess() bool

	// CanFillSignatureForm verifies whether filling in existing signature fields is allowed by
	// the PDF document permissions dictionary. Port of #canFillSignatureForm.
	CanFillSignatureForm() bool

	// CanCreateSignatureField verifies whether the creation of new signature fields is allowed
	// by the PDF permissions dictionary. Port of #canCreateSignatureField.
	CanCreateSignatureField() bool

	// CertificationPermission returns the value of the /DocMDP dictionary defining the permitted
	// modifications in the PDF, when present. Port of #getCertificationPermission.
	CertificationPermission() enumerations.CertificationPermission

	// IsUsageRightsSignaturePresent verifies whether the PDF contains a usage rights signature.
	// Port of #isUsageRightsSignaturePresent.
	IsUsageRightsSignaturePresent() bool

	// CatalogDictionary returns the document catalog as a dictionary.
	// Port of #getCatalogDictionary.
	CatalogDictionary() PdfDict

	// PdfHeaderVersion returns the version of the PDF document defined in the document's header.
	// Port of #getPdfHeaderVersion.
	PdfHeaderVersion() float32

	// Version returns the version of the PDF document: the catalog's /Version when present, the
	// header version otherwise. Port of #getVersion.
	Version() float32

	// SetVersion sets the PDF version number by upgrading the /Catalog dictionary's /Version
	// entry. Port of #setVersion.
	SetVersion(version float32)

	// CreatePdfDict creates an empty PdfDict. Port of #createPdfDict.
	CreatePdfDict() PdfDict

	// CreatePdfArray creates an empty PdfArray. Port of #createPdfArray.
	CreatePdfArray() PdfArray
}
