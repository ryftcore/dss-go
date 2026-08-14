// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/visible/ImageRotationUtils.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.visible is part of the "not ported" pdf/visible SPI stack per
// internal/pdf/DESIGN.md §0.2 (native_signature_drawer.go's header), but ImageRotationUtils is a
// pure geometry helper with no dependency on rasterisation, and BuildSignatureFieldBox (the
// geometry half of visible-signature placement, which IS implemented - see
// native_pdf_signature_service.go) needs it. Only the two-argument getRotation(VisualSignatureRotation,
// int) overload and the int overload of isSwapOfDimensionsRequired are called anywhere in this
// codebase, so only those are ported (the VisualSignatureRotation overload of
// isSwapOfDimensionsRequired and the one-argument getRotation are Java conveniences with no
// call site here); function names follow the ImageRotationUtils<Method> convention the SIGN
// chunk's handoff notes already assumed (Rotation/IsSwapOfDimensionsRequired/SwapDimensions/
// RotateRelativelyWrappingBox/EnsureNoRotate).
package pades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// Rotation angle constants. Port of ImageRotationUtils.ANGLE_0/90/180/270/360.
const (
	ImageRotationUtilsAngle0   = 0
	ImageRotationUtilsAngle90  = 90
	ImageRotationUtilsAngle180 = 180
	ImageRotationUtilsAngle270 = 270
	ImageRotationUtilsAngle360 = 360
)

// imageRotationUtilsSupportedAnglesErrorMessage is the message used when an unsupported rotation
// degree is encountered. Port of ImageRotationUtils.SUPPORTED_ANGLES_ERROR_MESSAGE.
const imageRotationUtilsSupportedAnglesErrorMessage = "rotation angle must be 90, 180, 270 or 360 (0)"

// imageRotationUtilsNeedRotation reports whether a rotation is requested at all.
// Port of the private needRotation(VisualSignatureRotation).
func imageRotationUtilsNeedRotation(visualSignatureRotation enumerations.VisualSignatureRotation) bool {
	return visualSignatureRotation != "" && visualSignatureRotation != enumerations.VisualSignatureRotation_NONE
}

// ImageRotationUtilsRotation returns the rotation based on the page's default rotation
// parameter. Port of the two-argument #getRotation(VisualSignatureRotation, int).
func ImageRotationUtilsRotation(visualSignatureRotation enumerations.VisualSignatureRotation, pageRotation int) int {
	rotate := ImageRotationUtilsAngle360
	if imageRotationUtilsNeedRotation(visualSignatureRotation) {
		switch visualSignatureRotation {
		case enumerations.VisualSignatureRotation_AUTOMATIC:
			rotate = ImageRotationUtilsAngle360 - pageRotation
		case enumerations.VisualSignatureRotation_ROTATE_90:
			rotate = ImageRotationUtilsAngle90
		case enumerations.VisualSignatureRotation_ROTATE_180:
			rotate = ImageRotationUtilsAngle180
		case enumerations.VisualSignatureRotation_ROTATE_270:
			rotate = ImageRotationUtilsAngle270
		default:
			panic(imageRotationUtilsSupportedAnglesErrorMessage)
		}
	}
	return rotate
}

// ImageRotationUtilsIsSwapOfDimensionsRequired verifies if a swap of dimensions is required with
// the current rotation. Port of the int overload of #isSwapOfDimensionsRequired.
func ImageRotationUtilsIsSwapOfDimensionsRequired(rotation int) bool {
	return ImageRotationUtilsAngle90 == rotation || ImageRotationUtilsAngle270 == rotation
}

// ImageRotationUtilsSwapDimensions swaps the dimensions of the given AnnotationBox.
// Port of #swapDimensions.
func ImageRotationUtilsSwapDimensions(annotationBox AnnotationBox) AnnotationBox {
	return NewAnnotationBox(annotationBox.MinY(), annotationBox.MinX(), annotationBox.MaxY(), annotationBox.MaxX())
}

// ImageRotationUtilsRotateRelativelyWrappingBox rotates the given annotationBox relative to
// wrappingBox according to the given rotation. Port of #rotateRelativelyWrappingBox.
func ImageRotationUtilsRotateRelativelyWrappingBox(annotationBox, wrappingBox AnnotationBox, rotation int) AnnotationBox {
	switch rotation {
	case ImageRotationUtilsAngle90:
		return NewAnnotationBox(annotationBox.MinY(),
			wrappingBox.Width()-annotationBox.MaxX(),
			annotationBox.MaxY(),
			wrappingBox.Width()-annotationBox.MinX())
	case ImageRotationUtilsAngle180:
		return NewAnnotationBox(wrappingBox.Width()-annotationBox.MaxX(),
			wrappingBox.Height()-annotationBox.MaxY(),
			wrappingBox.Width()-annotationBox.MinX(),
			wrappingBox.Height()-annotationBox.MinY())
	case ImageRotationUtilsAngle270:
		return NewAnnotationBox(wrappingBox.Height()-annotationBox.MaxY(),
			annotationBox.MinX(),
			wrappingBox.Height()-annotationBox.MinY(),
			annotationBox.MaxX())
	case ImageRotationUtilsAngle0, ImageRotationUtilsAngle360:
		return annotationBox
	default:
		panic(imageRotationUtilsSupportedAnglesErrorMessage)
	}
}

// ImageRotationUtilsEnsureNoRotate ensures the annotation wrapping box defines correct
// coordinates relative to the "noRotate" flag. Port of #ensureNoRotate.
func ImageRotationUtilsEnsureNoRotate(annotationBox AnnotationBox, pageRotation int) AnnotationBox {
	switch pageRotation {
	case ImageRotationUtilsAngle90:
		return NewAnnotationBox(
			annotationBox.MinX(),
			annotationBox.MaxY(),
			annotationBox.MinX()+annotationBox.Height(),
			annotationBox.MaxY()+annotationBox.Width())
	case ImageRotationUtilsAngle180:
		return NewAnnotationBox(
			annotationBox.MinX()-annotationBox.Width(),
			annotationBox.MaxY(),
			annotationBox.MinX(),
			annotationBox.MaxY()+annotationBox.Height())
	case ImageRotationUtilsAngle270:
		return NewAnnotationBox(
			annotationBox.MinX()-annotationBox.Height(),
			annotationBox.MaxY()-annotationBox.Width(),
			annotationBox.MinX(),
			annotationBox.MaxY())
	case ImageRotationUtilsAngle0, ImageRotationUtilsAngle360:
		return annotationBox
	default:
		panic(fmt.Sprintf("The rotation degree '%d' is not supported!", pageRotation))
	}
}
