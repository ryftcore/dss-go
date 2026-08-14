// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/DefaultPdfObjectModificationsFinder.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). slf4j is dropped, per PORTING.md; the sole
// WARN-level messages that are not pure debug tracing (the maximum-deepness warning and the
// "unsupported comparison" ones) are kept as comments where they fired.
//
// Java's Set<ObjectModification>, populated with a LinkedHashSet to get a deterministic
// insertion-ordered, duplicate-free walk, becomes objectModificationSet: an insertion-ordered
// slice with an Add that skips an entry already present per ObjectModification.Equals (which,
// like Java's equals()/hashCode() override, compares only the object tree and the modification
// type - see object_modification.go).
//
// Java overloads find(PdfDocumentReader, PdfDocumentReader) [the PdfObjectModificationsFinder
// interface method] and find(PdfDict, PdfDict) [an additional public method called directly by
// PdfSigDictWrapper#checkConsistency, i.e. this port's PdfSignatureDictionary.CheckConsistency].
// Go has no overloading; Find keeps the interface's reader-based signature (the name every
// landed sibling that references PdfObjectModificationsFinder was already written against -
// native_pdf_signature_service.go, i_pdf_obj_factory.go, pdf_signature_service.go), and the
// dict-based overload is FindInDicts - pdf_signature_dictionary.go's CheckConsistency was
// written against a same-named Find(PdfDict, PdfDict) that cannot coexist with the interface
// method of the same name on the same type, so its call site was corrected to FindInDicts
// during integration.
package pades

import (
	"bytes"
	"io"

	"github.com/utain/esig/dss/utils"
)

// DefaultPdfObjectModificationsFinder is the default implementation used to find the
// differences between two provided PDF revisions. Port of the DefaultPdfObjectModificationsFinder
// class.
type DefaultPdfObjectModificationsFinder struct {
	// maximumObjectVerificationDeepness defines the maximum value of enveloped objects tree
	// deepness to be checked.
	maximumObjectVerificationDeepness int

	// laxNumericComparison defines whether an integer shall be promoted to a real for comparison
	// against a real number.
	laxNumericComparison bool

	// pdfObjectModificationsFilter is used to categorize found object modifications to
	// different groups.
	pdfObjectModificationsFilter *PdfObjectModificationsFilter
}

// NewDefaultPdfObjectModificationsFinder instantiates an object with the default configuration.
// Port of the default constructor.
func NewDefaultPdfObjectModificationsFinder() *DefaultPdfObjectModificationsFinder {
	return &DefaultPdfObjectModificationsFinder{
		maximumObjectVerificationDeepness: 500,
		laxNumericComparison:              true,
	}
}

// SetMaximumObjectVerificationDeepness sets the maximum objects verification deepness of
// enveloped objects to be compared. Default: 500. Port of #setMaximumObjectVerificationDeepness.
func (f *DefaultPdfObjectModificationsFinder) SetMaximumObjectVerificationDeepness(maximumObjectVerificationDeepness int) {
	f.maximumObjectVerificationDeepness = maximumObjectVerificationDeepness
}

// SetLaxNumericComparison sets whether an integer number shall be promoted to a real for
// comparison against a real number. Default: true. Port of #setLaxNumericComparison.
func (f *DefaultPdfObjectModificationsFinder) SetLaxNumericComparison(laxNumericComparison bool) {
	f.laxNumericComparison = laxNumericComparison
}

// PdfObjectModificationsFilter gets the PdfObjectModificationsFilter, creating a new instance
// when not set. Port of #getPdfObjectModificationsFilter.
func (f *DefaultPdfObjectModificationsFinder) PdfObjectModificationsFilter() *PdfObjectModificationsFilter {
	if f.pdfObjectModificationsFilter == nil {
		f.pdfObjectModificationsFilter = NewPdfObjectModificationsFilter()
	}
	return f.pdfObjectModificationsFilter
}

// SetPdfObjectModificationsFilter sets the PdfObjectModificationsFilter used to categorize
// found differences between PDF objects. Port of #setPdfObjectModificationsFilter.
func (f *DefaultPdfObjectModificationsFinder) SetPdfObjectModificationsFilter(pdfObjectModificationsFilter *PdfObjectModificationsFilter) {
	if pdfObjectModificationsFilter == nil {
		panic("PdfObjectModificationsFilter cannot be null!")
	}
	f.pdfObjectModificationsFilter = pdfObjectModificationsFilter
}

// Find returns found and categorized object modifications occurred between
// originalRevisionReader and finalRevisionReader. Port of #find(PdfDocumentReader, PdfDocumentReader).
func (f *DefaultPdfObjectModificationsFinder) Find(originalRevisionReader, finalRevisionReader PdfDocumentReader) PdfObjectModifications {
	objectModifications := f.findObjectModifications(originalRevisionReader, finalRevisionReader)
	return f.PdfObjectModificationsFilter().Filter(objectModifications.items)
}

// findObjectModifications ports the private #findObjectModifications.
func (f *DefaultPdfObjectModificationsFinder) findObjectModifications(originalRevisionReader,
	finalRevisionReader PdfDocumentReader) *objectModificationSet {
	modifications := &objectModificationSet{}
	signedCatalogDict := originalRevisionReader.CatalogDictionary()
	finalCatalogDict := finalRevisionReader.CatalogDictionary()
	f.compareObjectsRecursively(modifications, newPdfObjectTreeReferenceSet(), NewPdfObjectTree(PAdESConstantsCatalogName),
		PAdESConstantsCatalogName, pdfObjectOrNil(signedCatalogDict), pdfObjectOrNil(finalCatalogDict))
	return modifications
}

// FindInDicts returns found and categorized object differences between two provided PdfDict
// objects. Port of the additional public #find(PdfDict, PdfDict) overload.
func (f *DefaultPdfObjectModificationsFinder) FindInDicts(originalRevisionDict, finalRevisionDict PdfDict) PdfObjectModifications {
	objectModifications := &objectModificationSet{}
	f.compareDictsRecursively(objectModifications, newPdfObjectTreeReferenceSet(), NewPdfObjectTree(""),
		originalRevisionDict, finalRevisionDict)
	return f.PdfObjectModificationsFilter().Filter(objectModifications.items)
}

// pdfObjectOrNil returns a PdfDict as the PdfObject interface, keeping a Java-null catalog dict
// (never actually nil in practice, but matching the general PdfObject-typed call signature) nil.
func pdfObjectOrNil(dict PdfDict) PdfObject {
	if dict == nil {
		return nil
	}
	return dict
}

func (f *DefaultPdfObjectModificationsFinder) compareDictsRecursively(modifications *objectModificationSet,
	processedObjects *pdfObjectTreeReferenceSet, objectTree *PdfObjectTree, signedDict, finalDict PdfDict) {
	signedRevObjNames := signedDict.List()
	finalRevObjNames := finalDict.List()
	for _, objectName := range signedRevObjNames {
		currentObjectTree := objectTree.Copy()
		objectKey := signedDict.ObjectKey(objectName)
		if !f.isProcessedReference(processedObjects, currentObjectTree, objectName, objectKey) {
			currentObjectTree.AddKey(objectName)
			f.addProcessedReference(processedObjects, currentObjectTree, objectName, objectKey)
			f.compareObjectsRecursively(modifications, processedObjects, currentObjectTree, objectName,
				signedDict.Object(objectName), finalDict.Object(objectName))
		}
	}

	for _, objectName := range finalRevObjNames {
		currentObjectTree := objectTree.Copy()
		if !containsString(signedRevObjNames, objectName) {
			currentObjectTree.AddKey(objectName)
			finalObject := finalDict.Object(objectName)
			if isPdfDictOrArray(finalObject) {
				objectKey := finalDict.ObjectKey(objectName)
				f.addProcessedReference(processedObjects, currentObjectTree, objectName, objectKey)
				modifications.Add(NewObjectModificationCreate(currentObjectTree, finalDict.Object(objectName)))
				// Upstream logs "Added entry with key '{}'." at debug level.
			} else {
				modifications.Add(NewObjectModificationModify(currentObjectTree, nil, finalObject))
				// Upstream logs "Added parameter with key name '{}'." at debug level.
			}
		}
	}

	f.compareDictStreams(modifications, objectTree, signedDict, finalDict)
}

func (f *DefaultPdfObjectModificationsFinder) compareObjectsRecursively(modifications *objectModificationSet,
	processedObjects *pdfObjectTreeReferenceSet, objectTree *PdfObjectTree, name string, signedObject, finalObject PdfObject) {
	if f.maximumObjectVerificationDeepness < objectTree.ChainDeepness() {
		// Upstream logs, at WARN unless the limit is the intentional-skip value of 0 (then DEBUG),
		// "Maximum objects verification deepness has been reached : {}. Chain of objects is skipped."
		return
	}

	switch {
	case signedObject == nil && finalObject != nil:
		if isPdfDictOrArray(finalObject) {
			modifications.Add(NewObjectModificationCreate(objectTree, finalObject))
			// Upstream logs "Added entry with key '{}'." at debug level.
		} else {
			modifications.Add(NewObjectModificationModify(objectTree, nil, finalObject))
			// Upstream logs "Added parameter with key name '{}'." at debug level.
		}

	case signedObject != nil && finalObject == nil:
		if isPdfDictOrArray(signedObject) {
			modifications.Add(NewObjectModificationDelete(objectTree, signedObject))
			// Upstream logs "Deleted entry with key '{}'." at debug level.
		} else {
			modifications.Add(NewObjectModificationModify(objectTree, signedObject, nil))
			// Upstream logs "Deleted parameter with key name '{}'." at debug level.
		}

	case signedObject != nil && finalObject != nil:
		signedDict, signedIsDict := signedObject.(PdfDict)
		finalDict, finalIsDict := finalObject.(PdfDict)
		signedArray, signedIsArray := signedObject.(PdfArray)
		finalArray, finalIsArray := finalObject.(PdfArray)
		signedSimple, signedIsSimple := signedObject.(*pdfSimpleObject)
		finalSimple, finalIsSimple := finalObject.(*pdfSimpleObject)

		switch {
		case signedIsDict && finalIsDict:
			f.compareDictsRecursively(modifications, processedObjects, objectTree, signedDict, finalDict)

		case signedIsArray && finalIsArray:
			f.compareArraysRecursively(modifications, processedObjects, objectTree, name, signedArray, finalArray, true)
			f.compareArraysRecursively(modifications, processedObjects, objectTree, name, finalArray, signedArray, false)

		case signedIsSimple && finalIsSimple:
			f.compareSimpleObjects(modifications, objectTree, signedSimple, finalSimple)

		default:
			// Either of differing dynamic type, or an unsupported comparison; both branches
			// upstream record a MODIFICATION. Upstream additionally WARNs for the "same class,
			// unsupported comparison" branch, which cannot be reached here (every PdfObject
			// implementation is one of PdfDict/PdfArray/*pdfSimpleObject).
			modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
			// Upstream logs, when the dynamic types differ, "Object with key name '{}' of type
			// '{}' has been modified to type '{}'." at debug level.
		}
	}
}

func (f *DefaultPdfObjectModificationsFinder) compareSimpleObjects(modifications *objectModificationSet,
	objectTree *PdfObjectTree, signedObject, finalObject *pdfSimpleObject) {
	signedObjectValue := signedObject.Value()
	finalObjectValue := finalObject.Value()

	switch {
	case signedObjectValue == nil && finalObjectValue != nil:
		modifications.Add(NewObjectModificationModify(objectTree, nil, finalObject))
		// Upstream logs "Added object value with key '{}'." at debug level.

	case signedObjectValue != nil && finalObjectValue == nil:
		modifications.Add(NewObjectModificationModify(objectTree, signedObject, nil))
		// Upstream logs "Deleted object value with key '{}'." at debug level.

	case signedObjectValue != nil && finalObjectValue != nil:
		signedString, signedIsString := signedObjectValue.(string)
		finalString, finalIsString := finalObjectValue.(string)
		signedNumber, signedIsNumber := signedObjectValue.(*float64)
		finalNumber, finalIsNumber := finalObjectValue.(*float64)
		signedBool, signedIsBool := signedObjectValue.(bool)
		finalBool, finalIsBool := finalObjectValue.(bool)

		switch {
		case signedIsString && finalIsString:
			if signedString != finalString {
				modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
				// Upstream logs "Object changed with key '{}'." at debug level.
			}

		case signedIsNumber && finalIsNumber:
			f.compareNumericValues(modifications, objectTree, signedObject, finalObject, signedNumber, finalNumber)

		case signedIsBool && finalIsBool:
			if signedBool != finalBool {
				modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
				// Upstream logs "Object changed with key '{}'." at debug level.
			}

		default:
			// Either of differing dynamic type, or an unsupported comparison (upstream WARNs
			// for the latter, unreachable here since every PdfSimpleObject value produced by
			// this port is a string, *float64 or bool).
			modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
			// Upstream logs, when the dynamic types differ, "Object with key name '{}' of type
			// '{}' has been modified to type '{}'." at debug level.
		}
	}
}

func (f *DefaultPdfObjectModificationsFinder) compareNumericValues(modifications *objectModificationSet,
	objectTree *PdfObjectTree, signedObject, finalObject *pdfSimpleObject, signedNumber, finalNumber *float64) {
	if signedNumber == nil || finalNumber == nil || *signedNumber == *finalNumber {
		return
	}
	if !f.laxNumericComparison {
		modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
		// Upstream logs "Object changed with key '{}'." at debug level.
	} else if float32(*signedNumber) != float32(*finalNumber) {
		modifications.Add(NewObjectModificationModify(objectTree, signedObject, finalObject))
		// Upstream logs "Object changed with key '{}'." at debug level.
	}
	// Otherwise upstream logs "Number object with key changed type '{}'. Set
	// #setLaxNumericComparison(false) to return a warning." at debug level.
}

func (f *DefaultPdfObjectModificationsFinder) compareArraysRecursively(modifications *objectModificationSet,
	processedObjects *pdfObjectTreeReferenceSet, objectTree *PdfObjectTree, name string,
	firstArray, secondArray PdfArray, signedFirst bool) {
	for i := 0; i < firstArray.Size(); i++ {
		currentObjectTree := objectTree.Copy()

		signedRevObject := firstArray.Object(i)
		var finalRevObject PdfObject

		objectKey := firstArray.ObjectKey(i)
		if objectKey != nil {
			for j := 0; j < secondArray.Size(); j++ {
				finalObjectKey := secondArray.ObjectKey(j)
				if objectKey == finalObjectKey {
					finalRevObject = secondArray.Object(j)
				}
			}
		} else if i < secondArray.Size() {
			finalRevObject = secondArray.Object(i)
		}

		if !f.isProcessedReference(processedObjects, currentObjectTree, name, objectKey) {
			f.addProcessedReference(processedObjects, currentObjectTree, name, objectKey)
			if signedFirst {
				f.compareObjectsRecursively(modifications, processedObjects, currentObjectTree, name, signedRevObject, finalRevObject)
			} else {
				f.compareObjectsRecursively(modifications, processedObjects, currentObjectTree, name, finalRevObject, signedRevObject)
			}
		}
	}
}

func (f *DefaultPdfObjectModificationsFinder) isProcessedReference(processedObjects *pdfObjectTreeReferenceSet,
	objectTree *PdfObjectTree, name string, objectKey PdfObjectKey) bool {
	return processedObjects.Contains(name, objectKey) || objectTree.IsProcessedReference(objectKey)
}

func (f *DefaultPdfObjectModificationsFinder) addProcessedReference(processedObjects *pdfObjectTreeReferenceSet,
	objectTree *PdfObjectTree, name string, objectKey PdfObjectKey) {
	if objectKey != nil {
		processedObjects.Add(name, objectKey)
		objectTree.AddReference(objectKey)
	}
}

func (f *DefaultPdfObjectModificationsFinder) compareDictStreams(modifications *objectModificationSet,
	objectTree *PdfObjectTree, signedDict, finalDict PdfDict) {
	currentObjectTree := objectTree.Copy()
	currentObjectTree.SetStream()

	signedStreamSize := f.rawStreamSizeSecurely(signedDict)
	finalStreamSize := f.rawStreamSizeSecurely(finalDict)

	switch {
	case signedStreamSize == -1 && finalStreamSize > -1:
		modifications.Add(NewObjectModificationCreate(currentObjectTree, finalDict))
		// Upstream logs "A stream has been added '{}'." at debug level.

	case signedStreamSize > -1 && finalStreamSize == -1:
		modifications.Add(NewObjectModificationDelete(currentObjectTree, signedDict))
		// Upstream logs "A stream has been removed '{}'." at debug level.

	case signedStreamSize != finalStreamSize:
		modifications.Add(NewObjectModificationModify(currentObjectTree, signedDict, finalDict))
		// Upstream logs "A stream has been modified '{}'." at debug level.

	case signedStreamSize > -1 && finalStreamSize > -1:
		signedStream, err := f.rawInputStreamSecurely(signedDict)
		if err != nil {
			// Upstream logs "Unable to compare underlying stream binaries. Reason : {}".
			return
		}
		defer func() { _ = signedStream.Close() }()
		finalStream, err := f.rawInputStreamSecurely(finalDict)
		if err != nil {
			// Upstream logs "Unable to compare underlying stream binaries. Reason : {}".
			return
		}
		defer func() { _ = finalStream.Close() }()

		equal, err := utils.CompareInputStreams(signedStream, finalStream)
		if err != nil {
			// Upstream logs "Unable to compare underlying stream binaries. Reason : {}".
			return
		}
		if !equal {
			modifications.Add(NewObjectModificationModify(currentObjectTree, signedDict, finalDict))
			// Upstream logs "A stream has been modified '{}'." at debug level.
		}
	}
}

func (f *DefaultPdfObjectModificationsFinder) rawStreamSizeSecurely(pdfDict PdfDict) int64 {
	size, err := pdfDict.RawStreamSize()
	if err != nil {
		// Upstream logs "Unable to read the underlying stream binaries. Reason : {}".
		return -1
	}
	return size
}

func (f *DefaultPdfObjectModificationsFinder) rawInputStreamSecurely(pdfDict PdfDict) (io.ReadCloser, error) {
	stream, err := pdfDict.CreateRawInputStream()
	if err != nil {
		return nil, err
	}
	if stream != nil {
		return stream, nil
	}
	return io.NopCloser(bytes.NewReader(nil)), nil
}

// isPdfDictOrArray reports whether object is a PdfDict or a PdfArray, standing in for Java's
// "instanceof PdfDict || instanceof PdfArray".
func isPdfDictOrArray(object PdfObject) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case PdfDict, PdfArray:
		return true
	default:
		return false
	}
}

// objectModificationSet is an insertion-ordered, duplicate-free collection of ObjectModifications,
// standing in for Java's LinkedHashSet<ObjectModification>; see this file's header.
type objectModificationSet struct {
	items []ObjectModification
}

// Add appends objectModification unless an equal entry (per ObjectModification.Equals) is
// already present.
func (s *objectModificationSet) Add(objectModification ObjectModification) {
	for _, existing := range s.items {
		if existing.Equals(objectModification) {
			return
		}
	}
	s.items = append(s.items, objectModification)
}

// pdfObjectTreeReferenceSet is a set of (name, PdfObjectKey) pairs, standing in for Java's
// private PdfObjectTreeReference and the HashSet<PdfObjectTreeReference> holding them.
type pdfObjectTreeReferenceSet struct {
	entries map[string]map[PdfObjectKey]struct{}
}

func newPdfObjectTreeReferenceSet() *pdfObjectTreeReferenceSet {
	return &pdfObjectTreeReferenceSet{entries: make(map[string]map[PdfObjectKey]struct{})}
}

// Contains reports whether (name, objectKey) was already added.
func (s *pdfObjectTreeReferenceSet) Contains(name string, objectKey PdfObjectKey) bool {
	if objectKey == nil {
		return false
	}
	keys, ok := s.entries[name]
	if !ok {
		return false
	}
	_, ok = keys[objectKey]
	return ok
}

// Add records (name, objectKey).
func (s *pdfObjectTreeReferenceSet) Add(name string, objectKey PdfObjectKey) {
	keys, ok := s.entries[name]
	if !ok {
		keys = make(map[PdfObjectKey]struct{})
		s.entries[name] = keys
	}
	keys[objectKey] = struct{}{}
}

// containsString (List#contains(Object) over a plain string slice) is already declared in
// pdf_signature_dictionary.go, reused here rather than redeclared.
