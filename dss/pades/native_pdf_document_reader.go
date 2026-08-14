// Native implementation of eu.europa.esig.dss.pdf.PdfDocumentReader on internal/pdf.
//
// Behavioural reference: eu.europa.esig.dss.pdf.pdfbox.PdfBoxDocumentReader. Every method below
// names the PdfBox method it reproduces. The object model underneath is internal/pdf's
// (*pdf.Document, *pdf.Dict, pdf.ObjectKey) instead of PDFBox's COS classes; the leniency of the
// reader is pinned to pdfbox 3.0.7 by internal/pdf itself (DESIGN.md §2.7), so this layer only
// translates, it never re-implements parsing.
//
// Deviations, all of them consequences of internal/pdf's declared non-goals (DESIGN.md §0.2):
//   - generateImageScreenshot / generateImageScreenshotWithoutAnnotations return
//     ErrRasterisationNotSupported: there is no rasteriser.
//   - PdfMemoryUsageSetting is accepted and ignored: the engine always works from a byte slice /
//     io.ReaderAt.
package pades

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/pdf"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/pades/exception"
	"github.com/utain/esig/dss/spi"
)

// ErrRasterisationNotSupported is returned by the page-screenshot methods. Upstream renders with
// PDFBox's PDFRenderer; the native engine deliberately has no rasteriser, no font engine and no
// content-stream interpreter (internal/pdf/DESIGN.md §0.2), so visual-difference detection and
// the signature-field previews are unavailable.
var ErrRasterisationNotSupported = errors.New("pades: PDF rasterisation is not supported by the native PDF engine")

// NativePdfDocumentReader reads a PDF document with the native internal/pdf engine.
type NativePdfDocumentReader struct {
	// document is the internal/pdf view of the PDF.
	document *pdf.Document

	// dssDocument is the PDF document being read.
	dssDocument model.DSSDocument

	// binaries is the document's content, kept because internal/pdf works over a byte slice and
	// isSignatureCoversWholeDocument needs the total length.
	binaries []byte

	// signatureDictionaries caches the result of ExtractSigDictionaries, as PdfBoxDocumentReader
	// caches signatureDictionaryMap.
	signatureDictionaries []PdfSignatureDictionaryFields

	// pendingVersion holds the value passed to SetVersion until the increment is laid out; see
	// SetVersion.
	pendingVersion *float32
}

// NewNativePdfDocumentReader opens the given document, optionally protected by
// passwordProtection. Port of the PdfBoxDocumentReader(DSSDocument, String, PdfMemoryUsageSetting)
// constructor; the memory-usage setting has no effect here (see the file header).
func NewNativePdfDocumentReader(dssDocument model.DSSDocument,
	passwordProtection []byte) (*NativePdfDocumentReader, error) {
	if dssDocument == nil {
		panic("The document must be defined!")
	}
	binaries, err := spi.DSSUtilsToByteArrayOfDocument(dssDocument)
	if err != nil {
		return nil, err
	}
	document, err := pdf.OpenBytes(binaries, &pdf.Options{Password: passwordProtection})
	if err != nil {
		if errors.Is(err, pdf.ErrInvalidPassword) {
			// Port of the catch of pdfbox's InvalidPasswordException.
			return nil, exception.NewInvalidPasswordException(fmt.Sprintf("Encrypted document : %s", err.Error()))
		}
		return nil, err
	}
	return &NativePdfDocumentReader{document: document, dssDocument: dssDocument, binaries: binaries}, nil
}

// Document returns the underlying internal/pdf document. It is the counterpart of
// PdfBoxDocumentReader#getPDDocument, i.e. the escape hatch the signature service uses.
func (r *NativePdfDocumentReader) Document() *pdf.Document {
	return r.document
}

// Binaries returns the document's content. There is no upstream counterpart: PDFBox keeps the
// source stream open, the native engine parses a byte slice, and the writer needs it back.
func (r *NativePdfDocumentReader) Binaries() []byte {
	return r.binaries
}

// Close releases the reader. Port of #close.
func (r *NativePdfDocumentReader) Close() error {
	return r.document.Close()
}

// DSSDictionary loads the last DSS dictionary from the document, nil when absent.
// Port of #getDSSDictionary.
func (r *NativePdfDocumentReader) DSSDictionary() PdfDssDict {
	return SingleDssDictExtract(r.CatalogDictionary())
}

// ExtractSigDictionaries extracts the signature dictionaries and the fields referring to them.
// Port of #extractSigDictionaries: one entry per distinct signature dictionary object, in field
// order, with every field that points at the same object grouped under it.
func (r *NativePdfDocumentReader) ExtractSigDictionaries() ([]PdfSignatureDictionaryFields, error) {
	if r.signatureDictionaries != nil {
		return r.signatureDictionaries, nil
	}
	signatureDictionaries := make([]PdfSignatureDictionaryFields, 0)
	byObjectNumber := make(map[int64]int)

	fields, err := r.document.SignatureFields()
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		sigFieldDict := newNativePdfDict(r.document, field.Dict, field.Key)
		pdfSignatureField := NewPdfSignatureField(sigFieldDict)

		if field.Value == nil || field.ValueKey.IsZero() {
			// Upstream logs "Signature field with name '{}' does not contain a signature".
			continue
		}

		if index, known := byObjectNumber[field.ValueKey.Num]; known {
			signatureDictionaries[index].Fields = append(signatureDictionaries[index].Fields, pdfSignatureField)
			// Upstream logs "More than one field refers to the same signature dictionary: {}!".
			continue
		}

		dictionary := newNativePdfDict(r.document, field.Value, field.ValueKey)
		signature, err := NewPdfSigDictWrapperFactory(dictionary).Create()
		if err != nil {
			// Upstream logs "Unable to create a PdfSignatureDictionary for field with name '{}'".
			continue
		}
		byObjectNumber[field.ValueKey.Num] = len(signatureDictionaries)
		signatureDictionaries = append(signatureDictionaries, PdfSignatureDictionaryFields{
			SignatureDictionary: signature,
			Fields:              []*PdfSignatureField{pdfSignatureField},
		})
	}
	r.signatureDictionaries = signatureDictionaries
	return r.signatureDictionaries, nil
}

// IsSignatureCoversWholeDocument checks whether the given signature covers the whole document.
// Port of #isSignatureCoversWholeDocument, arithmetic included:
// /ByteRange [0 575649 632483 10206] covers the file when
// (end1-start1) + (start2-end1-start1) + end2 equals the file length.
func (r *NativePdfDocumentReader) IsSignatureCoversWholeDocument(signatureDictionary *PdfSignatureDictionary) bool {
	byteRange := signatureDictionary.ByteRange()
	originalBytesLength := int64(len(r.binaries))
	beforeSignatureLength := int64(byteRange.FirstPartEnd()) - int64(byteRange.FirstPartStart())
	expectedCMSLength := int64(byteRange.SecondPartStart()) - int64(byteRange.FirstPartEnd()) -
		int64(byteRange.FirstPartStart())
	afterSignatureLength := int64(byteRange.SecondPartEnd())
	totalCoveredByByteRange := beforeSignatureLength + expectedCMSLength + afterSignatureLength

	return originalBytesLength == totalCoveredByByteRange
}

// NumberOfPages returns the number of pages of the document. Port of #getNumberOfPages.
func (r *NativePdfDocumentReader) NumberOfPages() int {
	return r.document.NumberOfPages()
}

// PageBox returns the /MediaBox of the given (1-based) page. Port of #getPageBox.
func (r *NativePdfDocumentReader) PageBox(page int) AnnotationBox {
	rect, err := r.document.PageBox(page)
	if err != nil {
		panic(err)
	}
	return NewAnnotationBox(float32(rect.MinX), float32(rect.MinY), float32(rect.MaxX), float32(rect.MaxY))
}

// PageRotation returns the rotation of the given (1-based) page. Port of #getPageRotation.
func (r *NativePdfDocumentReader) PageRotation(page int) int {
	return r.document.PageRotation(page)
}

// PdfAnnotations retrieves the annotations of the given (1-based) page.
// Port of #getPdfAnnotations together with the private toPdfAnnotation.
func (r *NativePdfDocumentReader) PdfAnnotations(page int) ([]*PdfAnnotation, error) {
	annotations, err := r.document.Annotations(page)
	if err != nil {
		return nil, err
	}
	pageRotation := r.PageRotation(page)
	result := make([]*PdfAnnotation, 0, len(annotations))
	for _, annotation := range annotations {
		annotationBox := NewAnnotationBox(float32(annotation.Rect.MinX), float32(annotation.Rect.MinY),
			float32(annotation.Rect.MaxX), float32(annotation.Rect.MaxY))
		if annotation.NoRotate {
			annotationBox = ImageRotationUtilsEnsureNoRotate(annotationBox, pageRotation)
		}
		pdfAnnotation := NewPdfAnnotation(annotationBox)
		pdfAnnotation.SetName(annotation.Name)
		pdfAnnotation.SetSigned(annotation.Signed)
		result = append(result, pdfAnnotation)
	}
	return result, nil
}

// GenerateImageScreenshot is not supported; see ErrRasterisationNotSupported.
// Port of #generateImageScreenshot.
func (r *NativePdfDocumentReader) GenerateImageScreenshot(page int) (model.DSSDocument, error) {
	return nil, ErrRasterisationNotSupported
}

// GenerateImageScreenshotWithoutAnnotations is not supported; see ErrRasterisationNotSupported.
// Port of #generateImageScreenshotWithoutAnnotations.
func (r *NativePdfDocumentReader) GenerateImageScreenshotWithoutAnnotations(page int,
	addedAnnotations []*PdfAnnotation) (model.DSSDocument, error) {
	return nil, ErrRasterisationNotSupported
}

// IsEncrypted checks whether the document is encrypted. Port of #isEncrypted.
func (r *NativePdfDocumentReader) IsEncrypted() bool {
	return r.document.IsEncrypted()
}

// IsOpenWithOwnerAccess verifies whether the document was opened with full owner access.
// Port of #isOpenWithOwnerAccess.
func (r *NativePdfDocumentReader) IsOpenWithOwnerAccess() bool {
	return r.document.Permissions().OwnerAccess
}

// CanFillSignatureForm verifies whether filling in existing signature fields is permitted.
// Port of #canFillSignatureForm.
func (r *NativePdfDocumentReader) CanFillSignatureForm() bool {
	permissions := r.document.Permissions()
	return permissions.CanModifyAnnots || permissions.CanFillInForm
}

// CanCreateSignatureField verifies whether the creation of new signature fields is permitted.
// Port of #canCreateSignatureField.
func (r *NativePdfDocumentReader) CanCreateSignatureField() bool {
	permissions := r.document.Permissions()
	return permissions.CanModify && permissions.CanModifyAnnots
}

// CertificationPermission returns the /DocMDP permission defined in /Root /Perms, when present.
// Port of #getCertificationPermission together with the private getPermissionFromReference.
func (r *NativePdfDocumentReader) CertificationPermission() enumerations.CertificationPermission {
	catalog, err := r.document.Catalog()
	if err != nil || catalog == nil {
		return ""
	}
	permsDict, ok := r.document.GetDict(catalog, "Perms")
	if !ok {
		return ""
	}
	signatureDict, ok := r.document.GetDict(permsDict, "DocMDP")
	if !ok {
		return ""
	}
	refArray, ok := r.document.GetArray(signatureDict, "Reference")
	if !ok {
		return ""
	}
	for i := range refArray {
		sigRefDict, ok := r.document.Resolve(refArray[i]).(*pdf.Dict)
		if !ok {
			continue
		}
		if transformMethod, ok := r.document.GetName(sigRefDict, "TransformMethod"); !ok ||
			transformMethod != "DocMDP" {
			continue
		}
		transformDict, ok := r.document.GetDict(sigRefDict, "TransformParams")
		if !ok {
			continue
		}
		accessPermissions, ok := r.document.GetInt(transformDict, "P")
		if !ok {
			accessPermissions = 2 // transformDict.getInt(COSName.P, 2)
		}
		if accessPermissions < 1 || accessPermissions > 3 {
			accessPermissions = 2
		}
		certificationPermission, err := enumerations.CertificationPermissionFromCode(int(accessPermissions))
		if err == nil {
			return certificationPermission
		}
	}
	return ""
}

// IsUsageRightsSignaturePresent verifies whether the PDF contains a usage rights signature.
// Port of #isUsageRightsSignaturePresent.
func (r *NativePdfDocumentReader) IsUsageRightsSignaturePresent() bool {
	catalog, err := r.document.Catalog()
	if err != nil || catalog == nil {
		return false
	}
	permsDict, ok := r.document.GetDict(catalog, "Perms")
	if !ok {
		return false
	}
	return permsDict.Has(pdf.Name(PAdESConstantsURName)) || permsDict.Has(pdf.Name(PAdESConstantsUR3Name))
}

// CatalogDictionary returns the document catalog. Port of #getCatalogDictionary.
func (r *NativePdfDocumentReader) CatalogDictionary() PdfDict {
	catalog, err := r.document.Catalog()
	if err != nil {
		panic(err)
	}
	return newNativePdfDict(r.document, catalog, pdf.ObjectKey{})
}

// PdfHeaderVersion returns the version from the file header. Port of #getPdfHeaderVersion.
func (r *NativePdfDocumentReader) PdfHeaderVersion() float32 {
	return r.document.HeaderVersion()
}

// Version returns the catalog's /Version when present, the header version otherwise.
// Port of #getVersion.
func (r *NativePdfDocumentReader) Version() float32 {
	return r.document.Version()
}

// SetVersion upgrades the /Catalog dictionary's /Version entry. Port of #setVersion.
//
// DEVIATION: internal/pdf is append-only and the catalog is rewritten by the Updater, not by the
// reader, so the value is recorded on the reader and applied by the signature service when it
// lays out the increment.
func (r *NativePdfDocumentReader) SetVersion(version float32) {
	r.pendingVersion = &version
}

// CreatePdfDict creates an empty dictionary. Port of #createPdfDict.
func (r *NativePdfDocumentReader) CreatePdfDict() PdfDict {
	return newNativePdfDict(r.document, pdf.NewDict(), pdf.ObjectKey{})
}

// CreatePdfArray creates an empty array. Port of #createPdfArray.
func (r *NativePdfDocumentReader) CreatePdfArray() PdfArray {
	return newNativePdfArray(r.document, pdf.Array{})
}

// GenerateDocumentID computes a document identifier in a deterministic way from the given
// parameters and the document, per ISO 32000-1 "14.4 File identifiers".
// Port of #generateDocumentId.
//
// DEVIATION: upstream feeds the computed long into PDFBox's COSWriter, which derives the 16-byte
// /ID entries from it with MD5. internal/pdf takes the second /ID element verbatim, so the same
// deterministic seed is hashed here instead. Nothing in PAdES reads /ID - it only has to be
// stable between the message-digest pass and the signing pass, which it is.
func (r *NativePdfDocumentReader) GenerateDocumentID(parameters PAdESCommonParameters) []byte {
	deterministicID := parameters.DeterministicId()
	if r.dssDocument != nil && r.dssDocument.Name() != "" {
		deterministicID = deterministicID + "-" + r.dssDocument.Name()
	}

	documentID := int64(len(r.binaries))
	// attach the String to the long value in a deterministic way. Java's "<<" on a long uses
	// only the low six bits of the shift distance, whereas Go shifts a 64-bit value by 64 or
	// more to zero, hence the explicit &63.
	for i, c := range []rune(deterministicID) {
		documentID += (int64(c) & 0xFF) << (uint(8*i) & 63)
	}
	seed := fmt.Sprintf("%d", documentID)
	digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_MD5, []byte(seed))
	if err != nil {
		panic(err)
	}
	return digest
}

// Compile-time assertion standing in for Java's "implements PdfDocumentReader".
var _ PdfDocumentReader = (*NativePdfDocumentReader)(nil)
