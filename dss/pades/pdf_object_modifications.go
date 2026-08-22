// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfObjectModifications.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf package implemented
// here (see pdf_object.go's header). java.io.Serializable is dropped, as elsewhere. Used by
// value (not pointer) throughout, matching pdf_signature_dictionary.go's usage
// (`func (m PdfObjectModifications) UndefinedChanges() []ObjectModification`).
package pades

// PdfObjectModifications contains a collection of ObjectModifications categorized by different
// groups. Port of the PdfObjectModifications class.
type PdfObjectModifications struct {
	// secureChanges are modifications not considered as changes in an incremental update (DSS
	// dictionary, DocTimeStamp creation).
	secureChanges []ObjectModification

	// formFillInAndSignatureCreationChanges are changes acceptable for /DocMDP P=2 (filling in
	// forms, instantiating page templates, and signing).
	formFillInAndSignatureCreationChanges []ObjectModification

	// annotCreationChanges are changes acceptable for /DocMDP P=3 (annotation creation, deletion
	// and modification).
	annotCreationChanges []ObjectModification

	// undefinedChanges are other changes that may invalidate the signature.
	undefinedChanges []ObjectModification
}

// AddSecureChange adds a secure change concerning signature augmentation (DSS dictionary,
// DocTimeStamp). Port of #addSecureChange.
func (m *PdfObjectModifications) AddSecureChange(objectModification ObjectModification) {
	m.secureChanges = append(m.secureChanges, objectModification)
}

// AddFormFillInAndSignatureCreationChange adds a modification concerning form filling or a
// signature creation. Port of #addFormFillInAndSignatureCreationChange.
func (m *PdfObjectModifications) AddFormFillInAndSignatureCreationChange(objectModification ObjectModification) {
	m.formFillInAndSignatureCreationChanges = append(m.formFillInAndSignatureCreationChanges, objectModification)
}

// AddAnnotCreationChange adds a modification concerning annotation creation, modification or
// deletion. Port of #addAnnotCreationChange.
func (m *PdfObjectModifications) AddAnnotCreationChange(objectModification ObjectModification) {
	m.annotCreationChanges = append(m.annotCreationChanges, objectModification)
}

// AddUndefinedChange adds an undefined modification. Port of #addUndefinedChange.
func (m *PdfObjectModifications) AddUndefinedChange(objectModification ObjectModification) {
	m.undefinedChanges = append(m.undefinedChanges, objectModification)
}

// SecureChanges returns the list of secure changes. Port of #getSecureChanges.
func (m PdfObjectModifications) SecureChanges() []ObjectModification { return m.secureChanges }

// FormFillInAndSignatureCreationChanges returns the list of form filling and signature creation
// related changes. Port of #getFormFillInAndSignatureCreationChanges.
func (m PdfObjectModifications) FormFillInAndSignatureCreationChanges() []ObjectModification {
	return m.formFillInAndSignatureCreationChanges
}

// AnnotCreationChanges returns the list of annotation creation/modification/deletion changes.
// Port of #getAnnotCreationChanges.
func (m PdfObjectModifications) AnnotCreationChanges() []ObjectModification {
	return m.annotCreationChanges
}

// UndefinedChanges returns the list of undefined changes. Port of #getUndefinedChanges.
func (m PdfObjectModifications) UndefinedChanges() []ObjectModification { return m.undefinedChanges }

// IsEmpty checks whether the object is empty. Port of #isEmpty.
func (m PdfObjectModifications) IsEmpty() bool {
	return len(m.secureChanges) == 0 && len(m.formFillInAndSignatureCreationChanges) == 0 &&
		len(m.annotCreationChanges) == 0 && len(m.undefinedChanges) == 0
}
