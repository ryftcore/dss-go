// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfObjectModificationsFilter.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). slf4j is dropped, per PORTING.md; the sole
// WARN (in the private isStreamFill) is kept as a comment where it fired.
package pades

import "github.com/ryftcore/dss-go/dss/enumerations"

// PdfObjectModificationsFilter is used to categorize ObjectModifications into four different
// categories. Port of the PdfObjectModificationsFilter class.
type PdfObjectModificationsFilter struct{}

// NewPdfObjectModificationsFilter is the default constructor.
func NewPdfObjectModificationsFilter() *PdfObjectModificationsFilter {
	return &PdfObjectModificationsFilter{}
}

// Filter categorizes the given collection of ObjectModifications into various categories and
// returns the result of filtering. Port of #filter.
func (f *PdfObjectModificationsFilter) Filter(objectModifications []ObjectModification) PdfObjectModifications {
	var pdfObjectModifications PdfObjectModifications
	for _, objectModification := range objectModifications {
		if f.skipChange(objectModification) {
			continue
		}
		if f.isExtensionChange(objectModification) {
			pdfObjectModifications.AddSecureChange(objectModification)
		} else if f.isSignatureOrFormFillChange(objectModification) {
			pdfObjectModifications.AddFormFillInAndSignatureCreationChange(objectModification)
		} else if f.isAnnotationChange(objectModification) {
			pdfObjectModifications.AddAnnotCreationChange(objectModification)
		} else {
			pdfObjectModifications.AddUndefinedChange(objectModification)
		}
	}
	return pdfObjectModifications
}

// skipChange allows skipping some modifications occurring in PdfBox and OpenPDF.
// Port of #skipChange.
func (f *PdfObjectModificationsFilter) skipChange(objectModification ObjectModification) bool {
	lastKey := objectModification.ObjectTree().LastKey()
	actionType := objectModification.ActionType()
	if actionType == enumerations.PdfObjectModificationTypeDeletion && lastKey == PAdESConstantsAppearanceDictionaryName {
		return true
	} else if actionType == enumerations.PdfObjectModificationTypeModification && lastKey == PAdESConstantsAnnotFlag {
		return true
	} else if actionType == enumerations.PdfObjectModificationTypeModification && lastKey == PAdESConstantsTypeName {
		return true
	} else if actionType == enumerations.PdfObjectModificationTypeModification && lastKey == PAdESConstantsItextName {
		return true
	}
	return false
}

// isExtensionChange returns whether the modification corresponds to a signature augmentation
// (such as DocTimeStamp or DSS dictionary creation). Port of #isExtensionChange.
func (f *PdfObjectModificationsFilter) isExtensionChange(objectModification ObjectModification) bool {
	return f.isDSSDictionaryChange(objectModification) ||
		f.isDocTimeStampAdded(objectModification) ||
		f.isDocTimeStampEmptyFieldFill(objectModification) ||
		f.isDocTimeStampEmptyFieldFontCreation(objectModification) ||
		f.isDocumentExtension(objectModification) ||
		f.isVersionChange(objectModification) ||
		f.isExtensionsChange(objectModification) ||
		f.isMetaDataChange(objectModification) ||
		f.isStructTreeRootChange(objectModification)
}

func (f *PdfObjectModificationsFilter) isDSSDictionaryChange(objectModification ObjectModification) bool {
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if key == PAdESConstantsDssDictionaryName {
			return true
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isDocTimeStampAdded(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	if f.isAnnotsKey(key) {
		if addedDict, ok := objectModification.FinalObject().(PdfDict); ok {
			if valueDict, ok := addedDict.Object(PAdESConstantsValueName).(PdfDict); ok && f.isDocTimeStamp(valueDict) {
				return true
			}
		}
	} else if f.isValueKey(key) {
		if addedDict, ok := objectModification.FinalObject().(PdfDict); ok && f.isDocTimeStamp(addedDict) {
			return true
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isSignature(pdfDict PdfDict) bool {
	return f.isObjectOfType(pdfDict, PAdESConstantsSignatureType)
}

func (f *PdfObjectModificationsFilter) isDocTimeStamp(pdfDict PdfDict) bool {
	return f.isObjectOfType(pdfDict, PAdESConstantsTimestampType)
}

func (f *PdfObjectModificationsFilter) isObjectOfType(pdfDict PdfDict, typeValue string) bool {
	if pdfDict == nil {
		return false
	}
	typeObject := pdfDict.Object(PAdESConstantsTypeName)
	if typeObject == nil {
		return false
	}
	value, ok := typeObject.Value().(string)
	return ok && typeValue == value
}

func (f *PdfObjectModificationsFilter) isDocumentExtension(objectModification ObjectModification) bool {
	// can be relevant for /DSS or/and /DocTimeStamp incorporation
	key := objectModification.ObjectTree().LastKey()
	parentKey := f.getParentKey(objectModification)
	return key == PAdESConstantsExtensionsName && parentKey == PAdESConstantsCatalogName
}

// isSignatureOrFormFillChange returns whether the modification corresponds to a signature
// addition or a form fill. Port of #isSignatureOrFormFillChange.
func (f *PdfObjectModificationsFilter) isSignatureOrFormFillChange(objectModification ObjectModification) bool {
	return f.isFieldFilled(objectModification) ||
		f.isAnnotsArrayCreation(objectModification) ||
		f.isEmptyAnnotFill(objectModification) ||
		f.isAnnotsFill(objectModification) ||
		f.isFieldAppearanceCreationChange(objectModification) ||
		f.isFieldValueAssignmentChange(objectModification) ||
		f.isSignatureEmptyFieldFill(objectModification) ||
		f.isCatalogPieceInfoChange(objectModification) ||
		f.isCatalogPermsCreationChange(objectModification) ||
		f.isCatalogNamesChange(objectModification) ||
		f.isCatalogOutputIntentsChange(objectModification) ||
		f.isAcroFormDictionaryChange(objectModification) ||
		f.isSignatureEmptyFieldFontCreation(objectModification)
}

func (f *PdfObjectModificationsFilter) isFieldFilled(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	if key == PAdESConstantsValueName && f.isAnnotChange(objectModification) {
		return true
	} else if f.isAnnotChange(objectModification) {
		addedDict, ok := objectModification.FinalObject().(PdfDict)
		return ok && f.isValueChange(addedDict)
	}
	return false
}

func (f *PdfObjectModificationsFilter) isValueChange(pdfDict PdfDict) bool {
	return pdfDict.AsDict(PAdESConstantsValueName) != nil
}

func (f *PdfObjectModificationsFilter) isValueKey(key string) bool {
	return f.isOneOf(key, PAdESConstantsValueName)
}

func (f *PdfObjectModificationsFilter) isAnnotChange(objectModification ObjectModification) bool {
	lastKey := objectModification.ObjectTree().LastKey()
	parentKey := f.getParentKey(objectModification)
	if f.isAnnotsKey(lastKey) || f.isAnnotsKey(parentKey) {
		return true
	}

	isAnnotProcessing := false
	resetChain := false
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if resetChain {
			isAnnotProcessing = false
			resetChain = false
		}
		if f.isAnnotsKey(key) {
			isAnnotProcessing = true
		} else if !f.isParentOrKid(key) {
			resetChain = true
		}
	}
	return isAnnotProcessing && (f.isParentOrKid(lastKey) || f.isParentOrKid(parentKey))
}

func (f *PdfObjectModificationsFilter) isAnnotsKey(key string) bool {
	return f.isOneOf(key, PAdESConstantsAnnotsName, PAdESConstantsFieldsName)
}

func (f *PdfObjectModificationsFilter) isParentOrKid(key string) bool {
	return f.isOneOf(key, PAdESConstantsParentName, PAdESConstantsKidsName)
}

func (f *PdfObjectModificationsFilter) isAnnotsFill(objectModification ObjectModification) bool {
	if objectModification.ActionType() != enumerations.PdfObjectModificationTypeDeletion {
		lastKey := objectModification.ObjectTree().LastKey()
		return f.isAnnotChange(objectModification) && !f.isAnnotsKey(lastKey)
	}
	return false
}

func (f *PdfObjectModificationsFilter) isAnnotsArrayCreation(objectModification ObjectModification) bool {
	lastKey := objectModification.ObjectTree().LastKey()
	if objectModification.ActionType() != enumerations.PdfObjectModificationTypeCreation || !f.isAnnotsKey(lastKey) {
		return false
	}
	_, ok := objectModification.FinalObject().(PdfArray)
	return ok
}

func (f *PdfObjectModificationsFilter) isFieldAppearanceCreationChange(objectModification ObjectModification) bool {
	appearanceDictChangeFound := false
	annotChangeFound := false
	if objectModification.ActionType() == enumerations.PdfObjectModificationTypeCreation {
		for _, chainKey := range objectModification.ObjectTree().KeyChain() {
			if f.isAnnotsKey(chainKey) {
				annotChangeFound = true
			} else if chainKey == PAdESConstantsAppearanceDictionaryName {
				appearanceDictChangeFound = true
				break
			}
		}
	}
	return appearanceDictChangeFound && annotChangeFound
}

func (f *PdfObjectModificationsFilter) isFieldValueAssignmentChange(objectModification ObjectModification) bool {
	appearanceDictChangeFound := false
	annotChangeFound := false
	if objectModification.ActionType() == enumerations.PdfObjectModificationTypeCreation {
		for _, chainKey := range objectModification.ObjectTree().KeyChain() {
			if f.isAnnotsKey(chainKey) {
				annotChangeFound = true
			} else if chainKey == PAdESConstantsValueName {
				appearanceDictChangeFound = true
				break
			}
		}
	}
	return appearanceDictChangeFound && annotChangeFound
}

func (f *PdfObjectModificationsFilter) isSignatureEmptyFieldFill(objectModification ObjectModification) bool {
	return f.isEmptyFieldFill(objectModification, PAdESConstantsSignatureType)
}

func (f *PdfObjectModificationsFilter) isDocTimeStampEmptyFieldFill(objectModification ObjectModification) bool {
	return f.isEmptyFieldFill(objectModification, PAdESConstantsTimestampType)
}

func (f *PdfObjectModificationsFilter) isEmptyFieldFill(objectModification ObjectModification, signatureType string) bool {
	return f.isEmptyAnnotFill(objectModification) && f.checkRecursivelyForNewSignatureCreation(
		objectModification.OriginalObject(), objectModification.FinalObject(), signatureType)
}

func (f *PdfObjectModificationsFilter) isEmptyAnnotFill(objectModification ObjectModification) bool {
	appearanceDictChangeFound := false
	normalAppearanceFound := false
	if objectModification.ActionType() == enumerations.PdfObjectModificationTypeModification {
		for _, chainKey := range objectModification.ObjectTree().KeyChain() {
			if chainKey == PAdESConstantsAppearanceDictionaryName {
				appearanceDictChangeFound = true
			} else if appearanceDictChangeFound && chainKey == PAdESConstantsNormalAppearanceName {
				normalAppearanceFound = true
			}

			if normalAppearanceFound {
				if chainKey == PAdESConstantsLengthName {
					return true
				} else if f.isStreamFill(objectModification) {
					return true
				}
			}
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) checkRecursivelyForNewSignatureCreation(originalObject, finalObject PdfObject, targetType string) bool {
	var originalSigValue, finalSigValue PdfObject
	if originalDict, ok := originalObject.(PdfDict); ok {
		originalSigValue = originalDict.Object(PAdESConstantsValueName)
	}
	if finalDict, ok := finalObject.(PdfDict); ok {
		finalSigValue = finalDict.Object(PAdESConstantsValueName)
	}
	if originalSigValue == nil {
		if finalSigValueDict, ok := finalSigValue.(PdfDict); ok && f.isObjectOfType(finalSigValueDict, targetType) {
			return true
		}
	}

	var originalParent, finalParent PdfObject
	if originalObject != nil {
		originalParent = originalObject.Parent()
	}
	if finalObject != nil {
		finalParent = finalObject.Parent()
	}
	if originalParent == nil && finalParent == nil {
		return false
	}
	return f.checkRecursivelyForNewSignatureCreation(originalParent, finalParent, targetType)
}

func (f *PdfObjectModificationsFilter) isStreamFill(objectModification ObjectModification) bool {
	finalDict, ok := objectModification.FinalObject().(PdfDict)
	if _, okOriginal := objectModification.OriginalObject().(PdfDict); !okOriginal || !ok {
		return false
	}
	finalBytes, err := finalDict.StreamBytes()
	if err != nil {
		// Upstream logs "Unable to evaluate stream modification from path '{}'. Reason : {}".
		return false
	}
	return len(finalBytes) != 0
}

func (f *PdfObjectModificationsFilter) isVersionChange(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	parentKey := f.getParentKey(objectModification)
	return objectModification.ActionType() == enumerations.PdfObjectModificationTypeModification &&
		key == PAdESConstantsVersionName &&
		f.isOneOf(parentKey, PAdESConstantsCatalogName, PAdESConstantsDataName, PAdESConstantsRootName)
}

func (f *PdfObjectModificationsFilter) isExtensionsChange(objectModification ObjectModification) bool {
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if key == PAdESConstantsExtensionsName {
			return true
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isCatalogPieceInfoChange(objectModification ObjectModification) bool {
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if key == PAdESConstantsPieceInfoName {
			return true
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isCatalogPermsCreationChange(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	return objectModification.ActionType() == enumerations.PdfObjectModificationTypeCreation && key == PAdESConstantsPermsName
}

func (f *PdfObjectModificationsFilter) isCatalogNamesChange(objectModification ObjectModification) bool {
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if key == PAdESConstantsNamesName {
			return true
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isCatalogOutputIntentsChange(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	return objectModification.ActionType() == enumerations.PdfObjectModificationTypeCreation && key == PAdESConstantsOutputIntentsName
}

func (f *PdfObjectModificationsFilter) isMetaDataChange(objectModification ObjectModification) bool {
	key := objectModification.ObjectTree().LastKey()
	parentKey := f.getParentKey(objectModification)
	return key == PAdESConstantsMetadataName || parentKey == PAdESConstantsMetadataName
}

// isStructTreeRootChange's validation of the StructTreeRoot modifications is aligned with
// "allowed to be modified" changes of the iText library's DocumentRevisionsValidator
// (github.com/itext/itext-java, sign/.../DocumentRevisionsValidator.java, per upstream's comment).
func (f *PdfObjectModificationsFilter) isStructTreeRootChange(objectModification ObjectModification) bool {
	isStructTreeRoot := false
	keyChain := objectModification.ObjectTree().KeyChain()
	lastKey := objectModification.ObjectTree().LastKey()
	actionType := objectModification.ActionType()
	for _, key := range keyChain {
		if key == PAdESConstantsStructTreeRootName {
			isStructTreeRoot = true
		} else if isStructTreeRoot && f.isOneOf(key,
			PAdESConstantsStructTreeRootIDTreeName,
			PAdESConstantsStructTreeRootParentTreeName,
			PAdESConstantsStructTreeRootParentTreeNextKeyName) {
			return true
		} else if isStructTreeRoot && key == PAdESConstantsStructTreeRootKName && lastKey == PAdESConstantsStructTreeRootKName &&
			(actionType == enumerations.PdfObjectModificationTypeCreation || actionType == enumerations.PdfObjectModificationTypeDeletion) {
			return true
		} else {
			isStructTreeRoot = false
		}
	}
	return false
}

func (f *PdfObjectModificationsFilter) isAcroFormDictionaryChange(objectModification ObjectModification) bool {
	containsAcroForm := false
	containsResourceDict := false
	for _, key := range objectModification.ObjectTree().KeyChain() {
		if key == PAdESConstantsAcroFormName {
			containsAcroForm = true
		} else if f.isOneOf(key, PAdESConstantsDocumentAppearanceName,
			PAdESConstantsDocumentResourcesName, PAdESConstantsSigFlagsName) {
			containsResourceDict = true
		}
	}
	return containsAcroForm && containsResourceDict
}

func (f *PdfObjectModificationsFilter) isSignatureEmptyFieldFontCreation(objectModification ObjectModification) bool {
	return f.isFontCreationChange(objectModification, PAdESConstantsSignatureType)
}

func (f *PdfObjectModificationsFilter) isDocTimeStampEmptyFieldFontCreation(objectModification ObjectModification) bool {
	return f.isFontCreationChange(objectModification, PAdESConstantsTimestampType)
}

func (f *PdfObjectModificationsFilter) isFontCreationChange(objectModification ObjectModification, signatureType string) bool {
	key := objectModification.ObjectTree().LastKey()
	parentKey := f.getParentKey(objectModification)
	return objectModification.ActionType() == enumerations.PdfObjectModificationTypeCreation &&
		(key == PAdESConstantsFontName || parentKey == PAdESConstantsFontName) &&
		f.checkRecursivelyForNewSignatureCreation(objectModification.OriginalObject(), objectModification.FinalObject(), signatureType)
}

func (f *PdfObjectModificationsFilter) getParentKey(objectModification ObjectModification) string {
	keyChain := objectModification.ObjectTree().KeyChain()
	if len(keyChain) > 1 {
		return keyChain[len(keyChain)-2]
	}
	return ""
}

// isAnnotationChange returns whether the modification corresponds to an annotation change.
// Port of #isAnnotationChange.
func (f *PdfObjectModificationsFilter) isAnnotationChange(objectModification ObjectModification) bool {
	return f.isOtherAnnotChange(objectModification)
}

func (f *PdfObjectModificationsFilter) isOtherAnnotChange(objectModification ObjectModification) bool {
	if f.isAnnotChange(objectModification) {
		if objectModification.ActionType() == enumerations.PdfObjectModificationTypeDeletion {
			if pdfDict, ok := objectModification.OriginalObject().(PdfDict); ok {
				if valueDict, ok := pdfDict.Object(PAdESConstantsValueName).(PdfDict); ok {
					return !f.isSignature(valueDict) && !f.isDocTimeStamp(valueDict)
				}
			}
		}
		return true
	}
	return false
}

func (f *PdfObjectModificationsFilter) isOneOf(key string, toCompare ...string) bool {
	for _, candidate := range toCompare {
		if key == candidate {
			return true
		}
	}
	return false
}
