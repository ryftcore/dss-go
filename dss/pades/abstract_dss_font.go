// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/AbstractDSSFont.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart), as elsewhere in this port.
package pades

// AbstractDSSFontDefaultTextSize is the default font size (12f). Port of DEFAULT_TEXT_SIZE.
const AbstractDSSFontDefaultTextSize float32 = 12

// AbstractDSSFont is the abstract implementation of a DSSFont. Concrete font types (DSSJavaFont,
// DSSFileFont) embed it and shadow SetSize where they need to additionally re-derive their native
// font (Go has no method overriding across embedding; see those files).
type AbstractDSSFont struct {
	// size is the size of the font.
	size float32
}

// NewAbstractDSSFont instantiates the object with the default configuration. Port of the
// protected default constructor.
func NewAbstractDSSFont() AbstractDSSFont {
	return AbstractDSSFont{size: AbstractDSSFontDefaultTextSize}
}

// Size ports #getSize.
func (f *AbstractDSSFont) Size() float32 {
	return f.size
}

// SetSize ports #setSize.
func (f *AbstractDSSFont) SetSize(size float32) {
	f.size = size
}
