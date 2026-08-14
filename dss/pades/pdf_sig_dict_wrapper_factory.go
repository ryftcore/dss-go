// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSigDictWrapperFactory.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). slf4j is dropped, per PORTING.md.
//
// Java's PdfSigDictWrapperFactory builds a PdfSigDictWrapper, the sole implementation of
// PdfSignatureDictionary; both already collapsed into this port's single *PdfSignatureDictionary
// type (see pdf_signature_dictionary.go's header). Create() therefore builds and returns a
// *PdfSignatureDictionary directly, matching native_pdf_document_reader.go's already-landed call
// site `NewPdfSigDictWrapperFactory(dictionary).Create()`. Java's DSSException, thrown unchecked
// from the private getContents/getByteRange getters Create() calls, becomes a returned error -
// the "protected builders keep (T, error)" convention the SIGN chunk's handoff notes document,
// and matching the (signature, err) pair the call site already destructures.
package pades

import (
	"fmt"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/enumerations"
)

// PdfSigDictWrapperFactory creates a PdfSignatureDictionary instance.
// Port of the PdfSigDictWrapperFactory class.
type PdfSigDictWrapperFactory struct {
	// sigFieldDictionary is the PDF dictionary representing the signature field.
	sigFieldDictionary PdfDict
}

// NewPdfSigDictWrapperFactory is the default constructor.
// Port of the PdfSigDictWrapperFactory(PdfDict) constructor.
func NewPdfSigDictWrapperFactory(sigFieldDictionary PdfDict) *PdfSigDictWrapperFactory {
	return &PdfSigDictWrapperFactory{sigFieldDictionary: sigFieldDictionary}
}

// Create creates a new PdfSignatureDictionary. Port of #create.
func (f *PdfSigDictWrapperFactory) Create() (*PdfSignatureDictionary, error) {
	contents, err := f.contents()
	if err != nil {
		return nil, err
	}
	cmsValue, err := cms.CMSUtilsParseToCMSBinaries(contents)
	if err != nil {
		return nil, err
	}
	byteRange, err := f.byteRange()
	if err != nil {
		return nil, err
	}

	pdfSignatureDictionary := NewPdfSignatureDictionary()
	pdfSignatureDictionary.SetDictionary(f.sigFieldDictionary)
	pdfSignatureDictionary.SetCMS(cmsValue)
	pdfSignatureDictionary.SetSignerName(f.sigFieldDictionary.StringValue(PAdESConstantsNameName))
	pdfSignatureDictionary.SetSigningDate(f.sigFieldDictionary.DateValue(PAdESConstantsSigningDateName))
	pdfSignatureDictionary.SetContactInfo(f.sigFieldDictionary.StringValue(PAdESConstantsContactInfoName))
	pdfSignatureDictionary.SetReason(f.sigFieldDictionary.StringValue(PAdESConstantsReasonName))
	pdfSignatureDictionary.SetLocation(f.sigFieldDictionary.StringValue(PAdESConstantsLocationName))
	pdfSignatureDictionary.SetType(f.sigFieldDictionary.NameValue(PAdESConstantsTypeName))
	pdfSignatureDictionary.SetFilter(f.sigFieldDictionary.NameValue(PAdESConstantsFilterName))
	pdfSignatureDictionary.SetSubFilter(f.sigFieldDictionary.NameValue(PAdESConstantsSubFilterName))
	pdfSignatureDictionary.SetContents(contents)
	pdfSignatureDictionary.SetByteRange(byteRange)
	pdfSignatureDictionary.SetDocMDP(f.docMDP())
	pdfSignatureDictionary.SetFieldMDP(f.fieldMDP())
	return pdfSignatureDictionary, nil
}

// contents ports the private getContents.
func (f *PdfSigDictWrapperFactory) contents() ([]byte, error) {
	contents, err := f.sigFieldDictionary.BinariesValue(PAdESConstantsContentsName)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the signature content: %w", err)
	}
	return contents, nil
}

// byteRange ports the private getByteRange.
func (f *PdfSigDictWrapperFactory) byteRange() (*ByteRange, error) {
	byteRangeArray := f.sigFieldDictionary.AsArray(PAdESConstantsByteRangeName)
	if byteRangeArray == nil {
		return nil, fmt.Errorf("unable to retrieve the '%s' field value", PAdESConstantsByteRangeName)
	}
	arraySize := byteRangeArray.Size()
	result := make([]int, arraySize)
	for i := 0; i < arraySize; i++ {
		number := byteRangeArray.Number(i)
		if number != nil {
			result[i] = int(*number)
		}
	}
	return NewByteRange(result), nil
}

// docMDP ports the private getDocMDP.
func (f *PdfSigDictWrapperFactory) docMDP() enumerations.CertificationPermission {
	referenceArray := f.sigFieldDictionary.AsArray(PAdESConstantsReferenceName)
	if referenceArray == nil {
		return ""
	}
	for i := 0; i < referenceArray.Size(); i++ {
		sigRef := referenceArray.AsDict(i)
		if sigRef == nil || sigRef.NameValue(PAdESConstantsTransformMethodName) != PAdESConstantsDocMdpName {
			continue
		}
		transformParams := sigRef.AsDict(PAdESConstantsTransformParamsName)
		if transformParams == nil {
			// Upstream logs "No '{}' dictionary found. Unable to perform a '{}' entry
			// validation!".
			continue
		}
		permissions := transformParams.NumberValue(PAdESConstantsPermissionsName)
		if permissions == nil {
			// Upstream logs "No '{}' parameter found. Unable to perform a '{}' entry
			// validation!".
			continue
		}
		certificationPermission, err := enumerations.CertificationPermissionFromCode(int(*permissions))
		if err != nil {
			continue
		}
		return certificationPermission
	}
	return ""
}

// fieldMDP ports the private getFieldMDP.
func (f *PdfSigDictWrapperFactory) fieldMDP() *SigFieldPermissions {
	referenceArray := f.sigFieldDictionary.AsArray(PAdESConstantsReferenceName)
	if referenceArray == nil {
		return nil
	}
	for i := 0; i < referenceArray.Size(); i++ {
		sigRef := referenceArray.AsDict(i)
		if sigRef == nil || sigRef.NameValue(PAdESConstantsTransformMethodName) != PAdESConstantsFieldMdpName {
			continue
		}
		dataDict := sigRef.AsDict(PAdESConstantsDataName)
		if dataDict == nil {
			// Upstream logs "No '{}' dictionary found. Unable to perform a '{}' entry
			// validation!".
			continue
		}
		dataDictType := dataDict.NameValue(PAdESConstantsTypeName)
		if dataDictType != PAdESConstantsCatalogName {
			// Upstream logs "Unsupported type of '{}' dictionary found : '{}'. The '{}'
			// validation skipped.".
			continue
		}
		transformParams := sigRef.AsDict(PAdESConstantsTransformParamsName)
		if transformParams == nil {
			// Upstream logs "No '{}' dictionary found. Unable to perform a '{}' entry
			// validation!".
			continue
		}
		return PAdESUtilsExtractPermissionsDictionary(transformParams)
	}
	return nil
}
