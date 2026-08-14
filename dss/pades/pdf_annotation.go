// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfAnnotation.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Shape (NewPdfAnnotation(AnnotationBox) *PdfAnnotation, plus
// SetName/SetSigned) confirmed against native_pdf_document_reader.go's PdfAnnotations.
package pades

// PdfAnnotation contains relative information about a PDF annotation. Port of the PdfAnnotation
// class.
type PdfAnnotation struct {
	// annotationBox defines the box of the annotation.
	annotationBox AnnotationBox

	// name is the name of the annotation.
	name string

	// signed defines whether the annotation is signed.
	signed bool
}

// NewPdfAnnotation builds a PdfAnnotation over the given box. Port of the
// PdfAnnotation(AnnotationBox) constructor.
func NewPdfAnnotation(annotationBox AnnotationBox) *PdfAnnotation {
	return &PdfAnnotation{annotationBox: annotationBox}
}

// AnnotationBox returns the AnnotationBox. Port of #getAnnotationBox.
func (a *PdfAnnotation) AnnotationBox() AnnotationBox { return a.annotationBox }

// Name returns the name of the annotation. Port of #getName.
func (a *PdfAnnotation) Name() string { return a.name }

// SetName sets the name of the annotation. Port of #setName.
func (a *PdfAnnotation) SetName(name string) { a.name = name }

// IsSigned gets whether the annotation field is signed. Port of #isSigned.
func (a *PdfAnnotation) IsSigned() bool { return a.signed }

// SetSigned sets whether the annotation field is signed. Port of #setSigned.
func (a *PdfAnnotation) SetSigned(signed bool) { a.signed = signed }

// Equals ports #equals: value equality over the box, name and signed flag. Named Equals (not
// Equal) to match this port's convention for a Java equals(Object) override with no direct Go
// operator counterpart (e.g. SignatureFieldParameters.Equals).
func (a *PdfAnnotation) Equals(other *PdfAnnotation) bool {
	if a == other {
		return true
	}
	if a == nil || other == nil {
		return false
	}
	return a.annotationBox == other.annotationBox && a.name == other.name && a.signed == other.signed
}
