// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/AnnotationBox.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Shape (NewAnnotationBox, NewAnnotationBoxFromFieldParameters,
// MinX/MinY/MaxX/MaxY/ToPdfPageCoordinates, used as a plain value type - see
// native_pdf_signature_service.go's `annotationBox := AnnotationBox{}`) confirmed against every
// already-landed call site.
//
// Java's isOverlap(PdfAnnotation) / isOverlap(AnnotationBox) overload pair has no Go overloading
// counterpart; this port only introduces the plain-AnnotationBox comparison (IsOverlap) plus an
// IsOverlapAnnotation(*PdfAnnotation) convenience that delegates to it, following the
// established "split by suffix" convention for de-overloaded Java methods (see
// pades_service.go's getAvailableSignatureFields{,WithPassword} precedent, cited in the SIGN
// chunk's handoff notes).
package pades

import (
	"fmt"
)

// AnnotationBox defines a PDF annotation dimension and position (note, shape, signature field,
// etc.). Port of the AnnotationBox class.
type AnnotationBox struct {
	// minX is the lower left X coordinate.
	minX float32
	// minY is the lower left Y coordinate.
	minY float32
	// maxX is the upper right X coordinate.
	maxX float32
	// maxY is the upper right Y coordinate.
	maxY float32
}

// NewAnnotationBox builds an AnnotationBox, normalizing the provided coordinates.
// Port of the AnnotationBox(float, float, float, float) constructor.
func NewAnnotationBox(minX, minY, maxX, maxY float32) AnnotationBox {
	box := AnnotationBox{minX: minX, minY: minY, maxX: maxX, maxY: maxY}
	if minX >= maxX {
		box.minX, box.maxX = maxX, minX
	}
	if minY >= maxY {
		box.minY, box.maxY = maxY, minY
	}
	return box
}

// NewAnnotationBoxFromFieldParameters instantiates an AnnotationBox from SignatureFieldParameters.
// Port of the AnnotationBox(SignatureFieldParameters) constructor.
func NewAnnotationBoxFromFieldParameters(fieldParameters *SignatureFieldParameters) AnnotationBox {
	return NewAnnotationBox(fieldParameters.OriginX(), fieldParameters.OriginY(),
		fieldParameters.OriginX()+fieldParameters.Width(), fieldParameters.OriginY()+fieldParameters.Height())
}

// MinX returns the lower left X coordinate. Port of #getMinX.
func (b AnnotationBox) MinX() float32 { return b.minX }

// MinY returns the lower left Y coordinate. Port of #getMinY.
func (b AnnotationBox) MinY() float32 { return b.minY }

// MaxX returns the upper right X coordinate. Port of #getMaxX.
func (b AnnotationBox) MaxX() float32 { return b.maxX }

// MaxY returns the upper right Y coordinate. Port of #getMaxY.
func (b AnnotationBox) MaxY() float32 { return b.maxY }

// Width returns the width of the box. Port of #getWidth.
func (b AnnotationBox) Width() float32 { return b.maxX - b.minX }

// Height returns the height of the box. Port of #getHeight.
func (b AnnotationBox) Height() float32 { return b.maxY - b.minY }

// ToPdfPageCoordinates creates a new AnnotationBox mirrored vertically relative to pageBox: in
// used PDF implementations the Y origin is bottom-based, while in DSS parameters it is
// top-based. This also accounts for non-zero upper-left corner coordinates, when applicable.
// Port of #toPdfPageCoordinates.
func (b AnnotationBox) ToPdfPageCoordinates(pageBox AnnotationBox) AnnotationBox {
	return NewAnnotationBox(pageBox.MinX()+b.minX, pageBox.MaxY()-b.maxY,
		pageBox.MinX()+b.maxX, pageBox.MaxY()-b.minY)
}

// IsOverlap checks if the current AnnotationBox overlaps with the given box. Port of
// #isOverlap(AnnotationBox).
func (b AnnotationBox) IsOverlap(box AnnotationBox) bool {
	if b.MinX() >= box.MaxX() || box.MinX() >= b.MaxX() {
		return false
	}
	if b.MinY() >= box.MaxY() || box.MinY() >= b.MaxY() {
		return false
	}
	return true
}

// IsOverlapAnnotation checks if the current AnnotationBox overlaps with the given PdfAnnotation.
// Port of #isOverlap(PdfAnnotation).
func (b AnnotationBox) IsOverlapAnnotation(pdfAnnotation *PdfAnnotation) bool {
	return b.IsOverlap(pdfAnnotation.AnnotationBox())
}

// String ports #toString.
func (b AnnotationBox) String() string {
	return fmt.Sprintf("AnnotationBox [minX=%v, minY=%v, maxX=%v, maxY=%v]", b.minX, b.minY, b.maxX, b.maxY)
}
