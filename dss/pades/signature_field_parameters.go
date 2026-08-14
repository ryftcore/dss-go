// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/SignatureFieldParameters.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart).
//
// eu.europa.esig.dss.pdf.visible.ImageUtils.DEFAULT_FIRST_PAGE (= 1) is a non-goal package per
// internal/pdf/DESIGN.md §0.2 (no rasteriser/renderer is ported); the literal 1 is inlined below
// rather than referencing an unported constant.
package pades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// SignatureFieldParameters holds parameters which allow creating a new signature field in a PDF
// document.
type SignatureFieldParameters struct {
	// fieldId is the signature field id/name (optional).
	fieldId string

	// page is the page number where the signature field is added. Port of ImageUtils's
	// DEFAULT_FIRST_PAGE default (see file header).
	page int

	// originX is the coordinate X where to add the signature field (origin is top/left corner).
	originX float32

	// originY is the coordinate Y where to add the signature field (origin is top/left corner).
	originY float32

	// width is the signature field width.
	width float32

	// height is the signature field height.
	height float32

	// rotation is the rotation to use on the PDF page, where the signature field will be
	// created.
	rotation enumerations.VisualSignatureRotation
}

// NewSignatureFieldParameters is the default constructor, instantiating the object with null
// values (page defaults to the first page, per ImageUtils.DEFAULT_FIRST_PAGE).
func NewSignatureFieldParameters() *SignatureFieldParameters {
	return &SignatureFieldParameters{page: 1}
}

// FieldId gets signature field id. Port of #getFieldId.
func (p *SignatureFieldParameters) FieldId() string {
	return p.fieldId
}

// SetFieldId sets a signature field id/name to place a signature into. Port of #setFieldId.
func (p *SignatureFieldParameters) SetFieldId(fieldId string) {
	p.fieldId = fieldId
}

// Page gets a page where the signature should be placed. Port of #getPage.
func (p *SignatureFieldParameters) Page() int {
	return p.page
}

// SetPage sets a page number where the signature field should be placed.
//
// NOTE: the counting starts from 1 (one) for the first page of the document
//
// Port of #setPage.
func (p *SignatureFieldParameters) SetPage(page int) {
	p.page = page
}

// OriginX gets an upper left X coordinate. Port of #getOriginX.
func (p *SignatureFieldParameters) OriginX() float32 {
	return p.originX
}

// SetOriginX sets a upper left X coordinate of the signature field. Port of #setOriginX.
func (p *SignatureFieldParameters) SetOriginX(originX float32) {
	p.originX = originX
}

// OriginY gets a upper left Y coordinate. Port of #getOriginY.
func (p *SignatureFieldParameters) OriginY() float32 {
	return p.originY
}

// SetOriginY sets a upper left Y coordinate of the signature field. Port of #setOriginY.
func (p *SignatureFieldParameters) SetOriginY(originY float32) {
	p.originY = originY
}

// Width gets a width of the signature field. Port of #getWidth.
func (p *SignatureFieldParameters) Width() float32 {
	return p.width
}

// SetWidth sets a width of the signature field. Port of #setWidth.
func (p *SignatureFieldParameters) SetWidth(width float32) {
	p.width = width
}

// Height gets a height of the signature field. Port of #getHeight.
func (p *SignatureFieldParameters) Height() float32 {
	return p.height
}

// SetHeight sets a height of the signature field. Port of #setHeight.
func (p *SignatureFieldParameters) SetHeight(height float32) {
	p.height = height
}

// Rotation returns rotation value for a signature field relatively the PDF page. Port of
// #getRotation.
func (p *SignatureFieldParameters) Rotation() enumerations.VisualSignatureRotation {
	return p.rotation
}

// SetRotation sets a rotation value for the signature field relatively the PDF page.
//
// rotation can be one of the following values:
//
//	NONE (DEFAULT value. No rotation is applied. The origin of coordinates begins from the top
//	left corner of a page);
//	AUTOMATIC (Rotates a signature field respectively to the page's rotation. Rotates the
//	signature field on the same value as a defined in a PDF page);
//	ROTATE_90 (Rotates a signature field for a 90° clockwise. Coordinates' origin begins from
//	top right page corner);
//	ROTATE_180 (Rotates a signature field for a 180° clockwise. Coordinates' origin begins from
//	the bottom right page corner);
//	ROTATE_270 (Rotates a signature field for a 270° clockwise. Coordinates' origin begins from
//	the bottom left page corner).
//
// Port of #setRotation.
func (p *SignatureFieldParameters) SetRotation(rotation enumerations.VisualSignatureRotation) {
	p.rotation = rotation
}

// String ports #toString.
func (p *SignatureFieldParameters) String() string {
	return fmt.Sprintf("SignatureFieldParameters [fieldId='%s', page=%d, originX=%v, originY=%v, width=%v, height=%v, rotation=%v]",
		p.fieldId, p.page, p.originX, p.originY, p.width, p.height, p.rotation)
}

// Equals ports #equals.
func (p *SignatureFieldParameters) Equals(other *SignatureFieldParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	return p.page == other.page &&
		p.originX == other.originX &&
		p.originY == other.originY &&
		p.width == other.width &&
		p.height == other.height &&
		p.fieldId == other.fieldId &&
		p.rotation == other.rotation
}
