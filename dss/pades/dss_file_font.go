// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/DSSFileFont.java (DSS 6.5.RC1).
//
// DEVIATION: #initializeDefault loads dss-pades/src/main/resources/fonts/PTSerifRegular.ttf as a
// classpath resource. This port has no font-rendering engine to consume real glyph data (see
// dss_font.go's file header) and no established go:embed convention elsewhere in the codebase, so
// the default font document carries the correct name (for anything that inspects it, e.g.
// #getName) but empty content instead of the real TrueType bytes.
//
// java.awt.Font -> *NativeJavaFont; #deriveJavaFont therefore cannot parse the font file's bytes
// (no TrueType parser is ported either) and instead derives a NativeJavaFont from the font
// document's name and the current size only.
package pades

import (
	"io"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DSSFileFontDefaultFontName is the default font name. Port of DEFAULT_FONT_NAME.
const DSSFileFontDefaultFontName = "PTSerifRegular.ttf"

// DSSFileFontDefaultFontExtension is the font file extension. Port of DEFAULT_FONT_EXTENSION.
const DSSFileFontDefaultFontExtension = ".ttf"

// dssFileFontDefaultFont is the default font resource. Port of DEFAULT_FONT; see the file header
// for why its content is empty rather than the real TrueType bytes.
//
// NewInMemoryDocumentWithName panics on a nil byte slice (Java's Objects.requireNonNull, ported
// literally) - an empty, non-nil []byte{} is what "empty content" actually means here, not nil,
// which would fail every package-level init depending on this var (i.e. everything importing
// pades) before a single test could even run.
var dssFileFontDefaultFont model.DSSDocument = model.NewInMemoryDocumentWithName([]byte{}, DSSFileFontDefaultFontName)

// DSSFileFont is the Font created from a file.
type DSSFileFont struct {
	AbstractDSSFont

	// fileFont is the font document.
	fileFont model.DSSDocument

	// javaFont is the Java implementation of the font (lazily derived).
	javaFont *NativeJavaFont

	// embedFontSubset defines whether only a subset of used glyphs should be embedded to a PDF,
	// when a font file is used with a text information defined within a signature field.
	//
	// DEFAULT : FALSE (all glyphs from a font file are embedded to a PDF)
	//
	// NOTE : this parameter will not take effect for DefaultPdfBoxVisibleSignatureDrawer
	embedFontSubset bool
}

// InitializeDefaultDSSFileFont initializes the default DSSFileFont. Port of #initializeDefault.
func InitializeDefaultDSSFileFont() *DSSFileFont {
	return NewDSSFileFontFromDocument(dssFileFontDefaultFont)
}

// NewDSSFileFontFromDocument is the constructor to load the font from a DSSDocument. Port of
// DSSFileFont(DSSDocument).
func NewDSSFileFontFromDocument(dssDocument model.DSSDocument) *DSSFileFont {
	return NewDSSFileFontFromDocumentWithSize(dssDocument, AbstractDSSFontDefaultTextSize)
}

// NewDSSFileFontFromDocumentWithSize is the constructor to load the font from a DSSDocument with
// a size. Port of DSSFileFont(DSSDocument, float).
func NewDSSFileFontFromDocumentWithSize(dssDocument model.DSSDocument, size float32) *DSSFileFont {
	if dssDocument == nil {
		panic("Font document cannot be null!")
	}
	if _, ok := dssDocument.(*model.DigestDocument); ok {
		panic("DigestDocument cannot be used as a font document!")
	}

	font := &DSSFileFont{AbstractDSSFont: NewAbstractDSSFont(), fileFont: dssDocument}
	font.size = size
	font.initFontName(dssDocument)
	return font
}

// initFontName ports the private #initFontName.
func (f *DSSFileFont) initFontName(fileFont model.DSSDocument) {
	if utils.IsStringBlank(fileFont.Name()) {
		data, err := spi.DSSUtilsToByteArrayOfDocument(fileFont)
		if err != nil {
			panic(err)
		}
		digest, err := spi.DSSUtilsMD5Digest(data)
		if err != nil {
			panic(err)
		}
		fileFont.SetName(digest + DSSFileFontDefaultFontExtension)
	}
}

// JavaFont ports the overridden #getJavaFont, lazily deriving the native font descriptor.
func (f *DSSFileFont) JavaFont() *NativeJavaFont {
	if f.javaFont == nil {
		f.javaFont = f.deriveJavaFont()
	}
	return f.javaFont
}

// deriveJavaFont ports the private #deriveJavaFont. See the file header for why this derives a
// NativeJavaFont from the document's name rather than parsing the font file's bytes.
func (f *DSSFileFont) deriveJavaFont() *NativeJavaFont {
	return &NativeJavaFont{Name: f.fileFont.Name(), Style: NativeJavaFontStylePlain, Size: f.size}
}

// InputStream gets the font's content InputStream. Port of #getInputStream.
func (f *DSSFileFont) InputStream() (io.ReadCloser, error) {
	return f.fileFont.OpenStream()
}

// Name gets the name of the font document. Port of #getName.
func (f *DSSFileFont) Name() string {
	return f.fileFont.Name()
}

// SetEmbedFontSubset sets whether only a subset of used glyphs should be embedded to a PDF, when
// a DSSFileFont is used.
//
// When set to TRUE, only the used glyphs will be embedded to a font.
// When set to FALSE, all glyphs from a font will be embedded to a PDF.
//
// DEFAULT : FALSE (the whole font file is embedded to a PDF)
//
// NOTE : this parameter will not take effect for DefaultPdfBoxVisibleSignatureDrawer
//
// Port of #setEmbedFontSubset.
func (f *DSSFileFont) SetEmbedFontSubset(embedFontSubset bool) {
	f.embedFontSubset = embedFontSubset
}

// IsEmbedFontSubset returns whether a font subset should be included into a PDF. Port of
// #isEmbedFontSubset.
func (f *DSSFileFont) IsEmbedFontSubset() bool {
	return f.embedFontSubset
}

// Equals ports #equals.
func (f *DSSFileFont) Equals(other *DSSFileFont) bool {
	if f == other {
		return true
	}
	if other == nil {
		return false
	}
	if f.embedFontSubset != other.embedFontSubset {
		return false
	}
	if f.fileFont != other.fileFont {
		return false
	}
	if f.javaFont == nil || other.javaFont == nil {
		return f.javaFont == other.javaFont
	}
	return *f.javaFont == *other.javaFont
}

var _ DSSFont = (*DSSFileFont)(nil)
