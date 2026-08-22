// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/DSSFont.java (DSS 6.5.RC1).
//
// java.awt.Font has no Go stdlib counterpart, and internal/pdf/DESIGN.md §0.2 declares the whole
// content-stream-interpretation stack (rasterisation, fonts, glyphs) a non-goal of the native PDF
// engine: there is no font-rendering/text-layout engine ported anywhere in this codebase (see
// native_signature_drawer.go's ErrVisualSignatureRenderingNotSupported for the same boundary
// applied to visible-signature painting). #getJavaFont therefore returns *NativeJavaFont, a small
// data-only placeholder carrying the same identity java.awt.Font exposed (name, style, size)
// without any glyph outline, metric or rendering capability - "pure data".
package pades

// NativeJavaFont is a minimal stand-in for java.awt.Font: name, style and size only, with no
// glyph/rendering capability (see the file header). It is not a port of any single upstream
// class; it exists because DSSFont#getJavaFont must return something.
type NativeJavaFont struct {
	// Name is the font family name.
	Name string

	// Style is the font style, one of the NativeJavaFontStyle* constants (mirroring
	// java.awt.Font's PLAIN/BOLD/ITALIC/BOLD+ITALIC bit flags).
	Style int

	// Size is the point size of the font.
	Size float32
}

// java.awt.Font style bit flags, reproduced verbatim since AbstractDSSFont/DSSJavaFont read them
// as plain ints (Java callers pass java.awt.Font.PLAIN etc. directly).
const (
	NativeJavaFontStylePlain  = 0
	NativeJavaFontStyleBold   = 1
	NativeJavaFontStyleItalic = 2
)

// DeriveFont returns a copy of f with the given size, mirroring java.awt.Font#deriveFont(float).
func (f *NativeJavaFont) DeriveFont(size float32) *NativeJavaFont {
	if f == nil {
		return nil
	}
	derived := *f
	derived.Size = size
	return &derived
}

// DSSFont defines a font used for a visual signature creation with text.
type DSSFont interface {
	// JavaFont gets the native font descriptor. Port of #getJavaFont.
	JavaFont() *NativeJavaFont

	// Size gets size of the font. Port of #getSize.
	Size() float32

	// SetSize sets size of the font. Port of #setSize.
	SetSize(size float32)
}
