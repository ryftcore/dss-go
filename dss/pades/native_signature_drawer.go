// The visible-signature seam of the native PDF backend.
//
// internal/pdf/DESIGN.md §0.2 declares the whole eu.europa.esig.dss.pdf.visible drawer stack a
// non-goal: there is no rasteriser, no font engine and no content-stream interpreter, so an
// image or a text appearance cannot be rendered. The four one-method SPI interfaces are
// nevertheless declared here, because AbstractPDFSignatureService is written against them
// (loadSignatureDrawer, buildSignatureFieldBox, getVisibleSignatureFieldBoxPosition) and the
// signature service would otherwise lose the shape of the upstream flow:
//
//	eu.europa.esig.dss.pdf.visible.SignatureDrawerFactory        -> SignatureDrawerFactory
//	eu.europa.esig.dss.pdf.visible.SignatureDrawer               -> SignatureDrawer
//	eu.europa.esig.dss.pdf.visible.SignatureFieldBoxBuilder      -> SignatureFieldBoxBuilder
//	eu.europa.esig.dss.pdf.visible.VisualSignatureFieldAppearance -> VisualSignatureFieldAppearance
//
// The native drawer implements the geometry half fully - it reports the field box built from
// SignatureFieldParameters, which is what the field placement, the /Rect and the empty
// appearance XObject need - and refuses the rendering half with
// ErrVisualSignatureRenderingNotSupported.
package pades

import "errors"

// ErrVisualSignatureRenderingNotSupported is returned by the native drawer when the signature
// image parameters ask for content to be painted (an image or a text block). The signature field,
// its /Rect and an empty appearance stream are still produced; only the painting is unavailable.
var ErrVisualSignatureRenderingNotSupported = errors.New(
	"pades: rendering of image/text visual signatures is not supported by the native PDF engine")

// SignatureDrawer draws a visual signature into the document being signed.
// Port of eu.europa.esig.dss.pdf.visible.SignatureDrawer.
type SignatureDrawer interface {
	// Draw paints the visual signature. Port of #draw.
	Draw() error
}

// SignatureDrawerFactory loads the relevant SignatureDrawer for the given image parameters.
// Port of eu.europa.esig.dss.pdf.visible.SignatureDrawerFactory.
type SignatureDrawerFactory interface {
	// SignatureDrawer returns the drawer to be used for the given image parameters.
	// Port of #getSignatureDrawer.
	SignatureDrawer(imageParameters *SignatureImageParameters) SignatureDrawer
}

// VisualSignatureFieldAppearance carries the box a visual signature field occupies.
// Port of eu.europa.esig.dss.pdf.visible.VisualSignatureFieldAppearance.
type VisualSignatureFieldAppearance interface {
	// AnnotationBox returns the field's box. Port of #getAnnotationBox.
	AnnotationBox() AnnotationBox
}

// SignatureFieldBoxBuilder builds the box of the signature field a drawer would paint into.
// Port of eu.europa.esig.dss.pdf.visible.SignatureFieldBoxBuilder.
type SignatureFieldBoxBuilder interface {
	// BuildSignatureFieldBox returns the field appearance, nil when the drawer cannot compute
	// one. Port of #buildSignatureFieldBox.
	BuildSignatureFieldBox() (VisualSignatureFieldAppearance, error)
}

// NativeSignatureDrawerFactory is the native counterpart of
// eu.europa.esig.dss.pdf.pdfbox.visible.defaultdrawer.PdfBoxDefaultSignatureDrawerFactory: it
// always returns the single native drawer, since there is no image/text renderer to choose
// between.
type NativeSignatureDrawerFactory struct{}

// NewNativeSignatureDrawerFactory creates the factory.
func NewNativeSignatureDrawerFactory() *NativeSignatureDrawerFactory {
	return &NativeSignatureDrawerFactory{}
}

// SignatureDrawer returns the native drawer for the given image parameters.
// Port of PdfBoxDefaultSignatureDrawerFactory#getSignatureDrawer.
func (f *NativeSignatureDrawerFactory) SignatureDrawer(imageParameters *SignatureImageParameters) SignatureDrawer {
	return &NativeSignatureDrawer{imageParameters: imageParameters}
}

// NativeSignatureDrawer computes the geometry of a visible signature field and refuses to render
// its content. It is the native stand-in for
// eu.europa.esig.dss.pdf.pdfbox.visible.defaultdrawer.DefaultPdfBoxVisibleSignatureDrawer.
type NativeSignatureDrawer struct {
	imageParameters *SignatureImageParameters
}

// ImageParameters returns the parameters the drawer was created for.
func (d *NativeSignatureDrawer) ImageParameters() *SignatureImageParameters {
	return d.imageParameters
}

// BuildSignatureFieldBox returns the field box taken from the signature field parameters, i.e.
// the origin/width/height the caller configured. Upstream derives it from the rendered image or
// text when those are set; without a renderer, only an explicitly sized field can be honoured,
// and a zero-sized one yields nil - the same "no box available" answer
// AbstractPDFSignatureService#buildSignatureFieldBox gives for a drawer that is not a
// SignatureFieldBoxBuilder. Port of #buildSignatureFieldBox.
func (d *NativeSignatureDrawer) BuildSignatureFieldBox() (VisualSignatureFieldAppearance, error) {
	if d.imageParameters == nil {
		return nil, nil
	}
	fieldParameters := d.imageParameters.FieldParameters()
	if fieldParameters == nil || fieldParameters.Width() == 0 || fieldParameters.Height() == 0 {
		return nil, nil
	}
	return &nativeVisualSignatureFieldAppearance{
		box: NewAnnotationBoxFromFieldParameters(fieldParameters),
	}, nil
}

// Draw refuses to paint the visual signature; see ErrVisualSignatureRenderingNotSupported.
// Port of #draw.
func (d *NativeSignatureDrawer) Draw() error {
	return ErrVisualSignatureRenderingNotSupported
}

// nativeVisualSignatureFieldAppearance is the trivial VisualSignatureFieldAppearance the native
// drawer returns.
type nativeVisualSignatureFieldAppearance struct {
	box AnnotationBox
}

// AnnotationBox returns the field's box. Port of #getAnnotationBox.
func (a *nativeVisualSignatureFieldAppearance) AnnotationBox() AnnotationBox {
	return a.box
}

var (
	_ SignatureDrawerFactory   = (*NativeSignatureDrawerFactory)(nil)
	_ SignatureDrawer          = (*NativeSignatureDrawer)(nil)
	_ SignatureFieldBoxBuilder = (*NativeSignatureDrawer)(nil)
)
