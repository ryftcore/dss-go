// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/ObjectModification.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header).
package pades

import "github.com/ryftcore/dss-go/dss/enumerations"

// ObjectModification represents a modification that occurred in a PDF document.
// Port of the ObjectModification class.
type ObjectModification struct {
	// objectTree is a collection of string keys representing a chain from a root to the actual object.
	objectTree *PdfObjectTree

	// originalObject is the signed revision modified object.
	originalObject PdfObject

	// finalObject is the final revision modified object.
	finalObject PdfObject

	// objectModificationType is the type of modification.
	objectModificationType enumerations.PdfObjectModificationType
}

// NewObjectModificationCreate creates an ObjectModification for a new object creation change.
// Port of the static #create.
func NewObjectModificationCreate(objectTree *PdfObjectTree, finalObject PdfObject) ObjectModification {
	return ObjectModification{objectTree: objectTree, finalObject: finalObject,
		objectModificationType: enumerations.PdfObjectModificationType_CREATION}
}

// NewObjectModificationDelete creates an ObjectModification for an object removal change.
// Port of the static #delete.
func NewObjectModificationDelete(objectTree *PdfObjectTree, originalObject PdfObject) ObjectModification {
	return ObjectModification{objectTree: objectTree, originalObject: originalObject,
		objectModificationType: enumerations.PdfObjectModificationType_DELETION}
}

// NewObjectModificationModify creates an ObjectModification for an object modification change.
// Port of the static #modify.
func NewObjectModificationModify(objectTree *PdfObjectTree, originalObject, finalObject PdfObject) ObjectModification {
	return ObjectModification{objectTree: objectTree, originalObject: originalObject, finalObject: finalObject,
		objectModificationType: enumerations.PdfObjectModificationType_MODIFICATION}
}

// ObjectTree returns the object tree. Port of #getObjectTree.
func (o ObjectModification) ObjectTree() *PdfObjectTree { return o.objectTree }

// OriginalObject gets the signed revision object. Port of #getOriginalObject.
func (o ObjectModification) OriginalObject() PdfObject { return o.originalObject }

// FinalObject gets the final document revision object. Port of #getFinalObject.
func (o ObjectModification) FinalObject() PdfObject { return o.finalObject }

// ActionType returns the corresponding object modification type. Port of #getActionType.
func (o ObjectModification) ActionType() enumerations.PdfObjectModificationType {
	return o.objectModificationType
}

// FieldName returns the name of the changed field object, when applicable. The object shall be
// of type field; returns "" for other objects (Java's null). Port of #getFieldName.
func (o ObjectModification) FieldName() string {
	if dict, ok := o.originalObject.(PdfDict); ok {
		return dict.StringValue(PAdESConstantsFieldNameName)
	} else if dict, ok := o.finalObject.(PdfDict); ok {
		return dict.StringValue(PAdESConstantsFieldNameName)
	}
	return ""
}

// Type returns the type of the concerned object, when applicable. Port of #getType.
func (o ObjectModification) Type() string {
	if dict, ok := o.originalObject.(PdfDict); ok {
		return objectModificationDictType(dict)
	} else if dict, ok := o.finalObject.(PdfDict); ok {
		return objectModificationDictType(dict)
	}
	return ""
}

// objectModificationDictType ports the private getType(PdfDict).
func objectModificationDictType(pdfDict PdfDict) string {
	valueDict := pdfDict.AsDict(PAdESConstantsValueName)
	if valueDict != nil {
		return valueDict.NameValue(PAdESConstantsTypeName)
	}
	return pdfDict.NameValue(PAdESConstantsTypeName)
}

// Equals ports #equals: two ObjectModifications are equal when their object tree and
// modification type match (the original/final objects are deliberately excluded, as upstream).
func (o ObjectModification) Equals(other ObjectModification) bool {
	if o.objectModificationType != other.objectModificationType {
		return false
	}
	return o.objectTree.Equals(other.objectTree)
}
