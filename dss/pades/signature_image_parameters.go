// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/SignatureImageParameters.java
// (DSS 6.5.RC1).
//
// java.awt.Color -> Go stdlib image/color.Color; see signature_image_text_parameters.go's file
// header. slf4j is dropped, per PORTING.md; the upstream log statement is kept as a comment
// where it fired. java.io.Serializable is dropped (no Go counterpart).
package pades

import (
	"fmt"
	"image/color"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureImageParametersNoScaling is the default zoom constraint. Port of NO_SCALING.
const SignatureImageParametersNoScaling = 100

// SignatureImageParameters holds parameters for a visible signature creation.
type SignatureImageParameters struct {
	// image contains the image to use (company logo,...).
	image model.DSSDocument

	// fieldParameters defines a SignatureFieldParameters like field positions and dimensions.
	fieldParameters *SignatureFieldParameters

	// zoom defines a percent to zoom the image (100% means no scaling). Note: this does not
	// touch zooming of the text representation.
	zoom int

	// backgroundColor defines the color of the image.
	backgroundColor color.Color

	// dpi defines the DPI of the image, nil when unset.
	dpi *int

	// legacyDPIHandling allows setting of the behavior on image scaling based on the DPI value.
	// This parameter is used to provide a smooth migration to the DSS version 6.5 and later.
	// When enabled (default in 6.5, to be changed later), the old behavior, hard-coded in DSS
	// 6.4 and before is applied, with the image is zoomed reversely, in the ratio of
	// PDF_DPI/Image_DPI. When disabled (recommended), the new behavior is used, with the image
	// being scaled in the ratio of Image_DPI/Target_DPI, effectively keeping the original image
	// size, unless a custom DPI value is defined.
	//
	// It is advisable to update the used implementation to the new behavior (with parameter
	// value set to true), as it will be enforced in future DSS versions.
	//
	// TODO : set to FALSE by default in DSS 6.6 (can change to boolean instead)
	legacyDPIHandling *bool

	// alignmentHorizontal is the horizontal alignment of the visual signature on the pdf page.
	alignmentHorizontal enumerations.VisualSignatureAlignmentHorizontal

	// alignmentVertical is the vertical alignment of the visual signature on the pdf page.
	alignmentVertical enumerations.VisualSignatureAlignmentVertical

	// imageScaling defines the image scaling behavior within a signature field with a fixed
	// size.
	//
	// DEFAULT : ImageScaling.STRETCH (stretches the image in both directions to fill the
	// signature field)
	imageScaling enumerations.ImageScaling

	// textParameters is used to define the text to generate on the image.
	textParameters *SignatureImageTextParameters
}

// NewSignatureImageParameters is the default constructor instantiating object with default
// parameters.
func NewSignatureImageParameters() *SignatureImageParameters {
	return &SignatureImageParameters{
		zoom:                SignatureImageParametersNoScaling,
		alignmentHorizontal: enumerations.VisualSignatureAlignmentHorizontal_NONE,
		alignmentVertical:   enumerations.VisualSignatureAlignmentVertical_NONE,
		imageScaling:        enumerations.ImageScaling_STRETCH,
	}
}

// Image returns a DSSDocument image defined for displaying on the signature field. Port of
// #getImage.
func (p *SignatureImageParameters) Image() model.DSSDocument {
	return p.image
}

// SetImage allows setting a custom image to display on a signature field. Panics with the Java
// message when image is a *model.DigestDocument. Port of #setImage.
func (p *SignatureImageParameters) SetImage(image model.DSSDocument) {
	if _, ok := image.(*model.DigestDocument); ok {
		panic("DigestDocument cannot be used as an image!")
	}
	p.image = image
}

// FieldParameters returns SignatureFieldParameters, lazily instantiating it. Port of
// #getFieldParameters.
func (p *SignatureImageParameters) FieldParameters() *SignatureFieldParameters {
	if p.fieldParameters == nil {
		p.fieldParameters = NewSignatureFieldParameters()
	}
	return p.fieldParameters
}

// SetFieldParameters sets SignatureFieldParameters, like signature field position and
// dimensions. Port of #setFieldParameters.
func (p *SignatureImageParameters) SetFieldParameters(fieldParameters *SignatureFieldParameters) {
	p.fieldParameters = fieldParameters
}

// Zoom returns the defined Zoom value in percentage. Port of #getZoom.
func (p *SignatureImageParameters) Zoom() int {
	return p.zoom
}

// SetZoom defines the signature field zoom in percentage (default value = 100). Port of
// #setZoom.
func (p *SignatureImageParameters) SetZoom(zoom int) {
	p.zoom = zoom
}

// BackgroundColor returns a specified background color for the signature field. Port of
// #getBackgroundColor.
func (p *SignatureImageParameters) BackgroundColor() color.Color {
	return p.backgroundColor
}

// SetBackgroundColor sets the background color for the signature field. Port of
// #setBackgroundColor.
func (p *SignatureImageParameters) SetBackgroundColor(backgroundColor color.Color) {
	p.backgroundColor = backgroundColor
}

// Dpi returns a defined DPI value. NOTE: can be nil. Port of #getDpi.
func (p *SignatureImageParameters) Dpi() *int {
	return p.dpi
}

// SetDpi sets an expected DPI value. If nil the default DPI of the provided image is applied and
// the image will take the exact space according to its dimensions. Otherwise, a ratio between
// the image DPI and the set value is to be applied.
//
// NOTE: images with a lower DPI will take more space on a PDF page
//
// Port of #setDpi.
func (p *SignatureImageParameters) SetDpi(dpi *int) {
	p.dpi = dpi
}

// IsLegacyDPIHandling gets whether the legacy DPI handling is to be applied. Port of
// #isLegacyDPIHandling.
//
// NOTE: upstream logs a warning when legacyDPIHandling is unset, pointing at
// SetLegacyDPIHandling; dropped as not load-bearing, per PORTING.md.
func (p *SignatureImageParameters) IsLegacyDPIHandling() bool {
	if p.legacyDPIHandling == nil {
		return true // TODO : Update in DSS 6.6
	}
	return *p.legacyDPIHandling
}

// SetLegacyDPIHandling sets the image scaling behavior based on the document DPI value. This
// parameter is used to provide a smooth migration to the DSS version 6.5 and later. When enabled
// (default in 6.5, to be changed later), the old behavior, hard-coded in DSS 6.4 and before is
// applied, with the image is zoomed reversely, in the ratio of PDF_DPI/Image_DPI. When disabled
// (recommended), the new behavior is used, with the image being scaled in the ratio of
// Image_DPI/Target_DPI, effectively keeping the original image size, unless a custom DPI value
// is defined.
//
// It is advisable to update the used implementation to the new behavior (with parameter value
// set to true), as it will be enforced in future DSS versions.
//
// Port of #setLegacyDPIHandling.
func (p *SignatureImageParameters) SetLegacyDPIHandling(legacyDPIHandling *bool) {
	p.legacyDPIHandling = legacyDPIHandling
}

// TextParameters returns text parameters, lazily instantiating them. Port of #getTextParameters.
func (p *SignatureImageParameters) TextParameters() *SignatureImageTextParameters {
	if p.textParameters == nil {
		p.textParameters = NewSignatureImageTextParameters()
	}
	return p.textParameters
}

// SetTextParameters sets text parameters. Port of #setTextParameters.
func (p *SignatureImageParameters) SetTextParameters(textParameters *SignatureImageTextParameters) {
	p.textParameters = textParameters
}

// VisualSignatureAlignmentHorizontal returns a horizontal alignment value of the signature
// field. Port of #getVisualSignatureAlignmentHorizontal.
func (p *SignatureImageParameters) VisualSignatureAlignmentHorizontal() enumerations.VisualSignatureAlignmentHorizontal {
	return p.alignmentHorizontal
}

// SetAlignmentHorizontal sets a horizontal alignment respectively to a page of the signature
// field. Port of #setAlignmentHorizontal.
func (p *SignatureImageParameters) SetAlignmentHorizontal(alignmentHorizontal enumerations.VisualSignatureAlignmentHorizontal) {
	p.alignmentHorizontal = alignmentHorizontal
}

// VisualSignatureAlignmentVertical returns a vertical alignment value of the signature field.
// Port of #getVisualSignatureAlignmentVertical.
func (p *SignatureImageParameters) VisualSignatureAlignmentVertical() enumerations.VisualSignatureAlignmentVertical {
	return p.alignmentVertical
}

// SetAlignmentVertical sets a vertical alignment respectively to a page of the signature field.
// Port of #setAlignmentVertical.
func (p *SignatureImageParameters) SetAlignmentVertical(alignmentVertical enumerations.VisualSignatureAlignmentVertical) {
	p.alignmentVertical = alignmentVertical
}

// ImageScaling gets the image scaling. Port of #getImageScaling.
func (p *SignatureImageParameters) ImageScaling() enumerations.ImageScaling {
	return p.imageScaling
}

// SetImageScaling sets the parameter used to define an image scaling behavior within a signature
// field with a fixed size.
//
// DEFAULT : ImageScaling.STRETCH (stretches the image in both directions in order to fill the
// signature field)
//
// Panics with the Java message when imageScaling is empty (Objects.requireNonNull). Port of
// #setImageScaling.
func (p *SignatureImageParameters) SetImageScaling(imageScaling enumerations.ImageScaling) {
	if imageScaling == "" {
		panic("ImageScaling parameter cannot be null!")
	}
	p.imageScaling = imageScaling
}

// IsEmpty checks if the SignatureImageParameters is empty (no image or text parameters are
// defined). Port of #isEmpty.
func (p *SignatureImageParameters) IsEmpty() bool {
	return p.image == nil && p.TextParameters().IsEmpty()
}

// String ports #toString.
func (p *SignatureImageParameters) String() string {
	return fmt.Sprintf("SignatureImageParameters [image=%v, fieldParameters=%v, zoom=%d, backgroundColor=%v, "+
		"dpi=%v, alignmentHorizontal=%v, alignmentVertical=%v, imageScaling=%v, textParameters=%v]",
		p.image, p.fieldParameters, p.zoom, p.backgroundColor, signatureImageParametersDpiString(p.dpi),
		p.alignmentHorizontal, p.alignmentVertical, p.imageScaling, p.textParameters)
}

// signatureImageParametersDpiString formats the optional dpi field for String.
func signatureImageParametersDpiString(dpi *int) string {
	if dpi == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%d", *dpi)
}

// Equals ports #equals. image is compared by identity (Go interface equality), matching the
// profileParametersDetachedContentsEqual precedent (document/profile_parameters.go) for
// interface-typed fields without a value Equals method.
func (p *SignatureImageParameters) Equals(other *SignatureImageParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.zoom != other.zoom {
		return false
	}
	if p.image != other.image {
		return false
	}
	if !signatureImageParametersFieldParametersEqual(p.fieldParameters, other.fieldParameters) {
		return false
	}
	if p.backgroundColor != other.backgroundColor {
		return false
	}
	if !signatureImageParametersDpiEqual(p.dpi, other.dpi) {
		return false
	}
	if p.alignmentHorizontal != other.alignmentHorizontal || p.alignmentVertical != other.alignmentVertical {
		return false
	}
	if p.imageScaling != other.imageScaling {
		return false
	}
	return signatureImageParametersTextParametersEqual(p.textParameters, other.textParameters)
}

func signatureImageParametersFieldParametersEqual(a, b *SignatureFieldParameters) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

func signatureImageParametersDpiEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func signatureImageParametersTextParametersEqual(a, b *SignatureImageTextParameters) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}
