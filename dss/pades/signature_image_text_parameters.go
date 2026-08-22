// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/SignatureImageTextParameters.java
// (DSS 6.5.RC1).
//
// java.awt.Color -> Go stdlib image/color.Color (the closest stdlib counterpart: a 32-bit ARGB
// pixel), consistent with there being no rendering engine ported (dss_font.go's file header).
// java.io.Serializable is dropped (no Go counterpart).
package pades

import (
	"fmt"
	"image/color"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SignatureImageTextParametersDefaultPadding is the default padding (5 pixels). Port of
// DEFAULT_PADDING.
const SignatureImageTextParametersDefaultPadding float32 = 5

// SignatureImageTextParametersDefaultBackgroundColor is the default background color (white).
// Port of DEFAULT_BACKGROUND_COLOR.
var SignatureImageTextParametersDefaultBackgroundColor color.Color = color.White

// SignatureImageTextParametersDefaultTextColor is the default text color (black). Port of
// DEFAULT_TEXT_COLOR.
var SignatureImageTextParametersDefaultTextColor color.Color = color.Black

// SignatureImageTextParameters allows custom text generation in the PAdES visible signature.
type SignatureImageTextParameters struct {
	// signerTextPosition allows adding signer name on the image (by default, LEFT).
	signerTextPosition enumerations.SignerTextPosition

	// signerTextVerticalAlignment defines the image from text vertical alignment in connection
	// with the image.
	//
	// It has effect when the SignerTextPosition is LEFT or RIGHT.
	signerTextVerticalAlignment enumerations.SignerTextVerticalAlignment

	// signerTextHorizontalAlignment sets the more line text horizontal alignment.
	signerTextHorizontalAlignment enumerations.SignerTextHorizontalAlignment

	// text defines the text to sign.
	text string

	// dssFont defines the font to use (default is PTSerifRegular).
	dssFont DSSFont

	// textWrapping defines how the given text should be wrapped within the signature field's
	// box.
	//
	// Default : TextWrapping.FONT_BASED - the text is computed based on the dssFont
	// configuration
	textWrapping enumerations.TextWrapping

	// padding defines a padding in pixels to bound text around (default is 5).
	padding float32

	// textColor defines the text color to use (default is BLACK)
	// (PAdES visual appearance: allow nil as text color, preventing graphic operators)
	textColor color.Color

	// backgroundColor defines the background of a text bounding box.
	backgroundColor color.Color
}

// NewSignatureImageTextParameters is the default constructor instantiating object with default
// configuration.
func NewSignatureImageTextParameters() *SignatureImageTextParameters {
	return &SignatureImageTextParameters{
		signerTextPosition:            enumerations.SignerTextPositionLeft,
		signerTextVerticalAlignment:   enumerations.SignerTextVerticalAlignmentMiddle,
		signerTextHorizontalAlignment: enumerations.SignerTextHorizontalAlignmentLeft,
		textWrapping:                  enumerations.TextWrappingFontBased,
		padding:                       SignatureImageTextParametersDefaultPadding,
		textColor:                     SignatureImageTextParametersDefaultTextColor,
		backgroundColor:               SignatureImageTextParametersDefaultBackgroundColor,
	}
}

// SignerTextPosition returns a signer text position respectively to an image. Port of
// #getSignerTextPosition.
func (p *SignatureImageTextParameters) SignerTextPosition() enumerations.SignerTextPosition {
	return p.signerTextPosition
}

// SetSignerTextPosition specifies a text position respectively to an image inside the signature
// field area (TOP, BOTTOM, RIGHT, LEFT). Port of #setSignerTextPosition.
func (p *SignatureImageTextParameters) SetSignerTextPosition(signerTextPosition enumerations.SignerTextPosition) {
	p.signerTextPosition = signerTextPosition
}

// SignerTextVerticalAlignment returns a signer text vertical alignment value. Port of
// #getSignerTextVerticalAlignment.
func (p *SignatureImageTextParameters) SignerTextVerticalAlignment() enumerations.SignerTextVerticalAlignment {
	return p.signerTextVerticalAlignment
}

// SetSignerTextVerticalAlignment defines a vertical alignment (positioning) of signer text
// inside the signature field (TOP, MIDDLE, BOTTOM). Port of #setSignerTextVerticalAlignment.
func (p *SignatureImageTextParameters) SetSignerTextVerticalAlignment(signerTextVerticalAlignment enumerations.SignerTextVerticalAlignment) {
	p.signerTextVerticalAlignment = signerTextVerticalAlignment
}

// SignerTextHorizontalAlignment returns a signer text horizontal alignment value. Port of
// #getSignerTextHorizontalAlignment.
func (p *SignatureImageTextParameters) SignerTextHorizontalAlignment() enumerations.SignerTextHorizontalAlignment {
	return p.signerTextHorizontalAlignment
}

// SetSignerTextHorizontalAlignment allows a horizontal alignment of a text with respect to its
// area (LEFT, CENTER, RIGHT). Port of #setSignerTextHorizontalAlignment.
func (p *SignatureImageTextParameters) SetSignerTextHorizontalAlignment(signerTextHorizontalAlignment enumerations.SignerTextHorizontalAlignment) {
	p.signerTextHorizontalAlignment = signerTextHorizontalAlignment
}

// Font returns the specified text font. If not defined, returns a default font instance
// (PTSerifRegular). Port of #getFont.
func (p *SignatureImageTextParameters) Font() DSSFont {
	if p.dssFont == nil {
		p.dssFont = InitializeDefaultDSSFileFont()
	}
	return p.dssFont
}

// SetFont sets a text font. Port of #setFont.
func (p *SignatureImageTextParameters) SetFont(dssFont DSSFont) {
	p.dssFont = dssFont
}

// TextWrapping gets TextWrapping. Port of #getTextWrapping.
func (p *SignatureImageTextParameters) TextWrapping() enumerations.TextWrapping {
	return p.textWrapping
}

// SetTextWrapping sets the TextWrapping parameter, defining a way the text will be generated.
//
// Default : TextWrapping.FONT_BASED - text is generated based in the provided dssFont
// configuration
//
// Panics with the Java message when textWrapping is empty (Objects.requireNonNull). Port of
// #setTextWrapping.
func (p *SignatureImageTextParameters) SetTextWrapping(textWrapping enumerations.TextWrapping) {
	if textWrapping == "" {
		panic("TextWrapping cannot be null!")
	}
	p.textWrapping = textWrapping
}

// Padding returns padding between text and its area. Port of #getPadding.
func (p *SignatureImageTextParameters) Padding() float32 {
	return p.padding
}

// SetPadding sets a padding between text and its area. Port of #setPadding.
func (p *SignatureImageTextParameters) SetPadding(padding float32) {
	p.padding = padding
}

// TextColor returns text color parameter. Port of #getTextColor.
func (p *SignatureImageTextParameters) TextColor() color.Color {
	return p.textColor
}

// SetTextColor sets color for the text.
//
// NOTE: use nil for the default text color (if supported by a selected implementation)
// DEFAULT: black
// (PAdES visual appearance: allow nil as text color, preventing graphic operators)
//
// Port of #setTextColor.
func (p *SignatureImageTextParameters) SetTextColor(textColor color.Color) {
	p.textColor = textColor
}

// BackgroundColor returns background color for the text's area. Port of #getBackgroundColor.
func (p *SignatureImageTextParameters) BackgroundColor() color.Color {
	return p.backgroundColor
}

// SetBackgroundColor sets the provided background color for a text's area.
//
// NOTE: use nil for a transparent background (if supported by a selected implementation)
// DEFAULT: white
// (PAdES visual appearance: allow nil as text color, preventing graphic operators)
//
// Port of #setBackgroundColor.
func (p *SignatureImageTextParameters) SetBackgroundColor(backgroundColor color.Color) {
	p.backgroundColor = backgroundColor
}

// Text returns defined text content. Port of #getText.
func (p *SignatureImageTextParameters) Text() string {
	return p.text
}

// SetText sets a text content parameter. Port of #setText.
func (p *SignatureImageTextParameters) SetText(text string) {
	p.text = text
}

// IsEmpty checks if the text property is set for the parameters. Port of #isEmpty.
func (p *SignatureImageTextParameters) IsEmpty() bool {
	return utils.IsStringEmpty(p.text)
}

// String ports #toString.
func (p *SignatureImageTextParameters) String() string {
	return fmt.Sprintf("SignatureImageTextParameters [signerTextPosition=%v, signerTextVerticalAlignment=%v, "+
		"signerTextHorizontalAlignment=%v, text='%s', dssFont=%v, textWrapping=%v, padding=%v, textColor=%v, "+
		"backgroundColor=%v]",
		p.signerTextPosition, p.signerTextVerticalAlignment, p.signerTextHorizontalAlignment, p.text, p.dssFont,
		p.textWrapping, p.padding, p.textColor, p.backgroundColor)
}

// Equals ports #equals. dssFont is compared by identity (Go interface equality), not by the
// structural DSSJavaFont/DSSFileFont#equals Java would dispatch to, per the
// profileParametersDetachedContentsEqual precedent (document/profile_parameters.go) for
// interface-typed fields without a value Equals method.
func (p *SignatureImageTextParameters) Equals(other *SignatureImageTextParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	return p.padding == other.padding &&
		p.signerTextPosition == other.signerTextPosition &&
		p.signerTextVerticalAlignment == other.signerTextVerticalAlignment &&
		p.signerTextHorizontalAlignment == other.signerTextHorizontalAlignment &&
		p.text == other.text &&
		p.dssFont == other.dssFont &&
		p.textWrapping == other.textWrapping &&
		p.textColor == other.textColor &&
		p.backgroundColor == other.backgroundColor
}
