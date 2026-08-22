// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/PAdESUtils.java (DSS 6.5.RC1).
//
// slf4j is dropped, per PORTING.md; the sole upstream log statement (a WARN in
// #getRevocationInfoArchival) is kept as a comment where it fired. java.io.BufferedInputStream
// buffering is not reproduced: Go's bufio wrapping is unnecessary here since InMemoryDocument's
// OpenStream already returns an in-memory reader.
//
// FORWARD DEPENDENCIES (not in this chunk's manifest, flattened into this same package pades by
// sibling chunks, per internal/pdf/doc.go's Layering section - every function below was written
// against these shapes, already assumed by the landed call sites named per function):
//   - PdfCMSRevision (eu.europa.esig.dss.pdf.PdfCMSRevision) - an interface with
//     PreviousRevision() model.DSSDocument, satisfied by *PdfSignatureRevision (see
//     pades_signature.go's forward-dependency comment).
//   - PdfDict (eu.europa.esig.dss.pdf.PdfDict), already landed by the SIGN chunk
//     (pades/native_pdf_dict.go implements it) - NameValue(name string) string and
//     AsArray(name string) PdfArray.
//   - PdfArray (eu.europa.esig.dss.pdf.PdfArray), already landed - Size() int and
//     String(i int) string.
//   - SigFieldPermissions (eu.europa.esig.dss.pdf.SigFieldPermissions) - a struct with
//     SetAction(enumerations.PdfLockAction), SetFields([]string) and
//     SetCertificationPermission(enumerations.CertificationPermission), per
//     pdf_signature_field.go's forward-dependency comment.
//   - PdfDssDict / PdfVriDict (eu.europa.esig.dss.pdf) - per
//     pdf_composite_dss_dict_certificate_source.go's forward-dependency comment:
//     PdfDssDict.VRIs() []*PdfVriDict, PdfVriDict.Name() string.
//   - PAdESConstants<Name> (eu.europa.esig.dss.pdf.PAdESConstants), flattened into this package,
//     following the naming already observed in native_pdf_signature_service.go /
//     pades_baseline_requirements_checker.go / pdf_signature_field.go.
//   - PdfMemoryUsageSetting (eu.europa.esig.dss.pdf.PdfMemoryUsageSetting) - a value type;
//     PdfMemoryUsageSettingMemoryFull() PdfMemoryUsageSetting is its flattened static factory
//     eu.europa.esig.dss.pdf.PdfMemoryUsageSetting#memoryFull, following the PAdESUtils<Method>
//     naming precedent (pdf_composite_dss_dict_certificate_source.go) for a flattened static.
package pades

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/signature/resources"
	"github.com/utain/esig/dss/utils"
)

// PAdESUtilsDefaultResourcesHandlerBuilder is the default resources handler builder to be used
// across the code. Port of DEFAULT_RESOURCES_HANDLER_BUILDER.
//
// TODO : replace with cades.DEFAULT_RESOURCES_HANDLER_BUILDER ?
var PAdESUtilsDefaultResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder = document.NewInMemoryResourcesHandlerBuilder()

// PAdESUtilsDefaultPdfMemoryUsageSetting is the default PDF memory usage setting to be used on a
// PDF document reading. Port of DEFAULT_PDF_MEMORY_USAGE_SETTING.
var PAdESUtilsDefaultPdfMemoryUsageSetting PdfMemoryUsageSetting = PdfMemoryUsageSettingMemoryFull()

// pAdESUtilsPDFPreamble is the starting bytes of a PDF document. Port of PDF_PREAMBLE.
var pAdESUtilsPDFPreamble = []byte{'%', 'P', 'D', 'F', '-'}

// pAdESUtilsPDFEOFString is the string used to end a PDF revision. Port of PDF_EOF_STRING.
var pAdESUtilsPDFEOFString = []byte{'%', '%', 'E', 'O', 'F'}

// PAdESUtilsGetOriginalPDF returns the original signed content for the padesSignature. Port of
// #getOriginalPDF(PAdESSignature).
func PAdESUtilsGetOriginalPDF(padesSignature *PAdESSignature) model.DSSDocument {
	return PAdESUtilsGetOriginalPDFFromRevision(padesSignature.PdfRevision())
}

// PAdESUtilsGetOriginalPDFFromRevision returns the original signed content for the pdfRevision.
// Port of #getOriginalPDF(PdfCMSRevision).
func PAdESUtilsGetOriginalPDFFromRevision(pdfRevision PdfCMSRevision) model.DSSDocument {
	return pdfRevision.PreviousRevision()
}

// PAdESUtilsGetRevisionContent returns the complete revision content according to the provided
// byteRange ([0]-[3]). Panics with the Java messages when dssDocument or byteRange is nil
// (Objects.requireNonNull). Port of #getRevisionContent.
func PAdESUtilsGetRevisionContent(dssDocument model.DSSDocument, byteRange *ByteRange) model.DSSDocument {
	if dssDocument == nil {
		panic("DSSDocument cannot be null!")
	}
	if byteRange == nil {
		panic("ByteRange cannot be null!")
	}
	PAdESUtilsAssertPdfDocument(dssDocument)

	beginning := byteRange.FirstPartStart()
	endSigValueContent := byteRange.SecondPartStart()
	endValue := byteRange.SecondPartEnd()

	revisionByteRange := pAdESUtilsGetTwoIntegersByteRange(beginning, endSigValueContent+endValue-beginning)
	return NewPdfByteRangeDocument(dssDocument, revisionByteRange)
}

// PAdESUtilsGetPreviousRevision returns the best previous revision from revisions corresponding
// to the byteRange. Panics with the Java messages when byteRange or revisions is nil (Objects.
// requireNonNull; a nil slice is treated as Java's non-null empty/absent collection would not be,
// so callers pass a non-nil slice). Port of #getPreviousRevision.
func PAdESUtilsGetPreviousRevision(byteRange *ByteRange, revisions []*PdfByteRangeDocument) model.DSSDocument {
	if byteRange == nil {
		panic("ByteRange cannot be null!")
	}
	if revisions == nil {
		panic("Revisions cannot be null!")
	}

	var bestCandidate *PdfByteRangeDocument
	firstPartLength := byteRange.FirstPartStart() + byteRange.FirstPartEnd()
	for _, byteRangeDocument := range revisions {
		currentByteRange := byteRangeDocument.ByteRange()
		if firstPartLength > currentByteRange.Length() &&
			(bestCandidate == nil || currentByteRange.Length() > bestCandidate.ByteRange().Length()) {
			bestCandidate = byteRangeDocument
		}
	}
	if bestCandidate != nil {
		return bestCandidate
	}
	return model.CreateEmptyDocument()
}

// PAdESUtilsGetSignatureValue gets the SignatureValue from the dssDocument according to the
// byteRange.
//
// Example: extracts bytes from 841 to 959. [0, 840, 960, 1200]
//
// Panics with the Java messages when dssDocument or byteRange is nil, or on a hex-decoding
// failure (Java's unchecked exception path). Port of #getSignatureValue.
func PAdESUtilsGetSignatureValue(dssDocument model.DSSDocument, byteRange *ByteRange) []byte {
	if dssDocument == nil {
		panic("DSSDocument cannot be null!")
	}
	if byteRange == nil {
		panic("ByteRange cannot be null!")
	}
	PAdESUtilsAssertPdfDocument(dssDocument)

	startSigValueContent := byteRange.FirstPartStart() + byteRange.FirstPartEnd() + 1
	endSigValueContent := byteRange.SecondPartStart() - 1

	sigValueDocument := NewPdfByteRangeDocument(dssDocument, pAdESUtilsGetTwoIntegersByteRange(startSigValueContent, endSigValueContent))
	data, err := spi.DSSUtilsToByteArrayOfDocument(sigValueDocument)
	if err != nil {
		panic(err)
	}
	value, err := utils.FromHex(string(data))
	if err != nil {
		panic(err)
	}
	return value
}

// PAdESUtilsReplaceSignature replaces /Contents field value with the given cmsSignedData
// binaries.
//
// toBeSignedDocument represents a document to be signed with an empty signature value (Ex.:
// /Contents <00000 ... 000000>). resourcesHandlerBuilder is optional: if nil,
// PAdESUtilsDefaultResourcesHandlerBuilder is used.
//
// Returns the PDF document containing the inserted CMS signature. Port of #replaceSignature.
func PAdESUtilsReplaceSignature(toBeSignedDocument model.DSSDocument, cmsSignedData []byte,
	resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) (result model.DSSDocument, err error) {
	if toBeSignedDocument == nil {
		panic("DSSDocument cannot be null!")
	}
	if cmsSignedData == nil {
		panic("cmsSignedData cannot be null!")
	}
	PAdESUtilsAssertPdfDocument(toBeSignedDocument)

	if resourcesHandlerBuilder == nil {
		resourcesHandlerBuilder = PAdESUtilsDefaultResourcesHandlerBuilder
	}

	if utils.IsArrayEmpty(cmsSignedData) {
		panic("cmsSignedData cannot be empty!")
	}
	signature := []byte(utils.ToHex(cmsSignedData))

	resourcesHandler := resourcesHandlerBuilder.CreateResourcesHandler()
	defer utils.CloseQuietly(resourcesHandler)

	os, err := resourcesHandler.CreateOutputStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to replace /Contents value within a toBeSigned document. Reason : %s", err.Error()), err)
	}
	is, err := toBeSignedDocument.OpenStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to replace /Contents value within a toBeSigned document. Reason : %s", err.Error()), err)
	}
	defer utils.CloseQuietly(is)
	bis := bufio.NewReader(is)
	bos := bufio.NewWriter(os)

	const startSuspicion = '<'
	const continueSuspicion = '0'

	suspicion := false
	cmsPasted := false
	var temp bytes.Buffer

	for {
		b, readErr := bis.ReadByte()
		if readErr != nil {
			break
		}

		if suspicion {
			if continueSuspicion == b {
				temp.WriteByte(b)
				if len(signature) == temp.Len() {
					if cmsPasted {
						return nil, exception.NewIllegalInputException("PDF document contains more than one empty signature!")
					}
					if _, werr := bos.Write(signature); werr != nil {
						return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
							"Unable to replace /Contents value within a toBeSigned document. Reason : %s", werr.Error()), werr)
					}
					temp.Reset()
					suspicion = false
					cmsPasted = true
				}
				continue
			}
			if _, werr := bos.Write(temp.Bytes()); werr != nil {
				return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
					"Unable to replace /Contents value within a toBeSigned document. Reason : %s", werr.Error()), werr)
			}
			temp.Reset()
			suspicion = false
		}

		if werr := bos.WriteByte(b); werr != nil {
			return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
				"Unable to replace /Contents value within a toBeSigned document. Reason : %s", werr.Error()), werr)
		}

		if startSuspicion == b {
			temp.Reset()
			suspicion = true
		}
	}

	if !cmsPasted {
		return nil, exception.NewIllegalInputException("Reserved space to insert a signature was not found!")
	}

	if err := bos.Flush(); err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to replace /Contents value within a toBeSigned document. Reason : %s", err.Error()), err)
	}

	return resourcesHandler.WriteToDSSDocument()
}

// PAdESUtilsExtractRevisions parses document and extracts all revisions based on the %%EOF
// string. Port of #extractRevisions.
func PAdESUtilsExtractRevisions(document model.DSSDocument) []*PdfByteRangeDocument {
	PAdESUtilsAssertPdfDocument(document)

	var revisions []*PdfByteRangeDocument

	is, err := document.OpenStream()
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Unable to retrieve the last revision", err))
	}
	defer utils.CloseQuietly(is)
	bis := bufio.NewReader(is)

	position := 0
	var tempLine bytes.Buffer
	for {
		b, readErr := bis.ReadByte()
		if readErr != nil {
			break
		}
		position++
		tempLine.WriteByte(b)
		stringBytes := tempLine.Bytes()

		if bytes.Equal(pAdESUtilsPDFEOFString, stringBytes) {
			tempLine.Reset()

			eofPosition := position
			c, cErr := bis.ReadByte()
			if cErr == nil {
				position++
			}

			if cErr == nil && c == spi.DSSUtilsLineFeed {
				// if \n
				eofPosition++
			} else if cErr == nil && c == spi.DSSUtilsCarriageReturn {
				// if \r
				eofPosition++

				d, dErr := bis.ReadByte()
				if dErr == nil {
					position++
				}
				if dErr == nil && d == spi.DSSUtilsLineFeed {
					// if \r\n
					eofPosition++
				}
			}

			revisions = append(revisions, NewPdfByteRangeDocument(document, pAdESUtilsGetTwoIntegersByteRange(0, eofPosition)))

		} else if spi.DSSUtilsIsLineBreakByte(b) || len(stringBytes) > len(pAdESUtilsPDFEOFString) {
			tempLine.Reset()
		}
	}

	return revisions
}

// pAdESUtilsGetTwoIntegersByteRange ports the private #getTwoIntegersByteRange.
func pAdESUtilsGetTwoIntegersByteRange(offset, position int) *ByteRange {
	return NewByteRange([]int{offset, position - offset, position, 0})
}

// PAdESUtilsRevocationInfoArchival returns RevocationInfoArchival from the given encodable, nil
// if the parsing failed or encodable is nil. Port of #getRevocationInfoArchival.
func PAdESUtilsRevocationInfoArchival(encodable *asn1ber.Element) *RevocationInfoArchival {
	if encodable == nil {
		return nil
	}
	archival, err := RevocationInfoArchivalGetInstance(encodable)
	if err != nil {
		// Upstream logs "Unable to parse RevocationInfoArchival" (slf4j, dropped per
		// PORTING.md).
		return nil
	}
	return archival
}

// PAdESUtilsIsPDFDocument checks if the given DSSDocument represents a PDF document. Panics with
// the Java message on an I/O failure (Java's unchecked DSSException). Port of #isPDFDocument.
func PAdESUtilsIsPDFDocument(document model.DSSDocument) bool {
	is, err := document.OpenStream()
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot read a sequence of bytes from the InputStream.", err))
	}
	defer utils.CloseQuietly(is)

	ok, err := utils.StartsWithStream(is, pAdESUtilsPDFPreamble)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot read a sequence of bytes from the InputStream.", err))
	}
	return ok
}

// PAdESUtilsAssertPdfDocument verifies whether the provided document is a PDF. Panics with the
// Java message when document is nil (Objects.requireNonNull), is a *model.DigestDocument, or is
// not a PDF (Java's IllegalArgumentException / spi/exception.IllegalInputException). Port of
// #assertPdfDocument.
func PAdESUtilsAssertPdfDocument(document model.DSSDocument) {
	if document == nil {
		panic("DSSDocument cannot be null!")
	}
	if _, ok := document.(*model.DigestDocument); ok {
		panic("DigestDocument cannot be used! PDF document is expected!")
	}
	if !PAdESUtilsIsPDFDocument(document) {
		panic(exception.NewIllegalInputException(fmt.Sprintf(
			"The document with name '%s' is not a PDF. PDF document is expected!", document.Name())))
	}
}

// PAdESUtilsExtractPermissionsDictionary extracts SigFieldPermissions (for instance /Lock
// dictionary) from a wrapping dictionary. Panics on an unsupported /Action field value
// (PdfLockActionForName's error, mirroring Java's unchecked IllegalArgumentException). Port of
// #extractPermissionsDictionary.
func PAdESUtilsExtractPermissionsDictionary(wrapper PdfDict) *SigFieldPermissions {
	sigFieldPermissions := NewSigFieldPermissions()

	action := wrapper.NameValue(PAdESConstantsActionName)
	lockAction, err := enumerations.PdfLockActionForName(action)
	if err != nil {
		panic(err)
	}
	sigFieldPermissions.SetAction(lockAction)

	var fields []string
	fieldsArray := wrapper.AsArray(PAdESConstantsFieldsName)
	if fieldsArray != nil {
		for j := 0; j < fieldsArray.Size(); j++ {
			field := fieldsArray.String(j)
			if field != "" {
				fields = append(fields, field)
			}
		}
	}
	sigFieldPermissions.SetFields(fields)

	if PAdESConstantsSigFieldLockName == wrapper.NameValue(PAdESConstantsTypeName) {
		permissions := wrapper.NumberValue(PAdESConstantsPermissionsName)
		if permissions != nil {
			if code, ok := pdfNumberToInt(*permissions); ok {
				certificationPermission, err := enumerations.CertificationPermissionFromCode(code)
				if err == nil {
					sigFieldPermissions.SetCertificationPermission(certificationPermission)
				}
			}
		}
	}

	return sigFieldPermissions
}

// PAdESUtilsVRIsWithName returns a list of VRI dictionaries, corresponding to the given
// signature (VRI) SHA-1 name.
//
// NOTE: vriName can be the empty string. In this case all /VRI dictionaries are returned (Java's
// null vriName).
//
// Port of #getVRIsWithName.
func PAdESUtilsVRIsWithName(pdfDssDict PdfDssDict, vriName string) []*PdfVriDict {
	vris := pdfDssDict.VRIs()
	if utils.IsCollectionEmpty(vris) {
		return nil
	}
	if vriName == "" {
		return vris
	}
	for _, vriDict := range vris {
		if vriName == vriDict.Name() {
			return []*PdfVriDict{vriDict}
		}
	}
	return nil
}

// PAdESUtilsInitializeDSSResourcesHandler initializes a new DSSResourcesHandler object. Port of
// #initializeDSSResourcesHandler.
func PAdESUtilsInitializeDSSResourcesHandler() resources.DSSResourcesHandler {
	return PAdESUtilsDefaultResourcesHandlerBuilder.CreateResourcesHandler()
}
