// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/DSSJavaFont.java (DSS 6.5.RC1).
//
// java.awt.Font -> *NativeJavaFont (see dss_font.go's file header for why: no font engine is
// ported in this codebase). #getJavaFont/#setSize therefore work on NativeJavaFont's plain
// name/style/size fields instead of deriving a real, glyph-capable font.
package pades

// DSSJavaFontDefaultFontStyle is the default font style (java.awt.Font.PLAIN). Port of
// DEFAULT_FONT_STYLE.
const DSSJavaFontDefaultFontStyle = NativeJavaFontStylePlain

// DSSJavaFont represents the JAVA implementation of the DSSFont.
type DSSJavaFont struct {
	AbstractDSSFont

	// javaFont is the native font descriptor.
	javaFont *NativeJavaFont
}

// NewDSSJavaFont is the default constructor, from a native font descriptor. Port of
// DSSJavaFont(Font).
func NewDSSJavaFont(javaFont *NativeJavaFont) *DSSJavaFont {
	font := &DSSJavaFont{AbstractDSSFont: NewAbstractDSSFont(), javaFont: javaFont}
	if javaFont != nil {
		font.size = javaFont.Size
	}
	return font
}

// NewDSSJavaFontFromName is the constructor from a font's name. Port of DSSJavaFont(String).
func NewDSSJavaFontFromName(fontName string) *DSSJavaFont {
	font := &DSSJavaFont{AbstractDSSFont: NewAbstractDSSFont()}
	font.javaFont = &NativeJavaFont{Name: fontName, Style: DSSJavaFontDefaultFontStyle, Size: AbstractDSSFontDefaultTextSize}
	return font
}

// NewDSSJavaFontFromNameSize is the constructor from a font's name with size. Port of
// DSSJavaFont(String, int).
func NewDSSJavaFontFromNameSize(fontName string, size int) *DSSJavaFont {
	font := &DSSJavaFont{AbstractDSSFont: NewAbstractDSSFont()}
	font.size = float32(size)
	font.javaFont = &NativeJavaFont{Name: fontName, Style: DSSJavaFontDefaultFontStyle, Size: float32(size)}
	return font
}

// NewDSSJavaFontFromNameStyleSize is the constructor from a font's name with a style and size.
// Port of DSSJavaFont(String, int, int).
func NewDSSJavaFontFromNameStyleSize(fontName string, style, size int) *DSSJavaFont {
	font := &DSSJavaFont{AbstractDSSFont: NewAbstractDSSFont()}
	font.size = float32(size)
	font.javaFont = &NativeJavaFont{Name: fontName, Style: style, Size: float32(size)}
	return font
}

// JavaFont ports the overridden #getJavaFont.
func (f *DSSJavaFont) JavaFont() *NativeJavaFont {
	return f.javaFont
}

// Name gets the name of the font. Port of #getName.
func (f *DSSJavaFont) Name() string {
	return f.javaFont.Name
}

// SetSize ports the overridden #setSize, additionally re-deriving javaFont for the new size.
func (f *DSSJavaFont) SetSize(size float32) {
	f.AbstractDSSFont.SetSize(size)
	f.javaFont = f.javaFont.DeriveFont(size)
}

// Equals ports #equals.
func (f *DSSJavaFont) Equals(other *DSSJavaFont) bool {
	if f == other {
		return true
	}
	if other == nil {
		return false
	}
	if f.javaFont == nil || other.javaFont == nil {
		return f.javaFont == other.javaFont
	}
	return *f.javaFont == *other.javaFont
}

var _ DSSFont = (*DSSJavaFont)(nil)
