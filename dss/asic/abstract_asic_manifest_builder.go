// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/AbstractASiCManifestBuilder.java
// (DSS 6.5.RC1).
//
// The manifest bytes this builder produces are covered by the CAdES signature over
// ASiCManifest*.xml, so the DOM is assembled through the frozen xml/utils DomUtils serializer, in
// upstream's exact element and attribute order - do not reorder anything below.
//
// Java's protected/abstract methods are exported: the concrete builders live in the
// dss/asic/{cades,xades} packages. Every self-call runs through AbstractASiCManifestBuilderOverrides
// (the Overrides + Init pattern), since the whole point of this class is that subclasses replace
// createRootElement/getSigReferenceMimeType/isRootfile/getManifestFilename/... .
package asic

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// AbstractASiCManifestBuilderOverrides captures the abstract and protected methods of
// AbstractASiCManifestBuilder that the class calls on itself.
type AbstractASiCManifestBuilderOverrides interface {
	// BuildDom builds the initial XML document. Port of the protected buildDom().
	BuildDom() *xmldom.Node

	// CreateRootElement creates the root element
	// {@code <asic:ASiCManifest xmlns:asic="http://uri.etsi.org/02918/v1.2.1#">}. Port of the
	// protected createRootElement(Document).
	CreateRootElement(documentDom *xmldom.Node) *xmldom.Node

	// AddSigReference adds a {@code <SigReference>} element. Port of the protected
	// addSigReference(Document, Element).
	AddSigReference(documentDom, asicManifestDom *xmldom.Node)

	// AddDataObjectReferences adds references to data objects, corresponding to the
	// ASiCContentDocumentFilter configuration. Port of the protected
	// addDataObjectReferences(Document, Element).
	AddDataObjectReferences(documentDom, asicManifestDom *xmldom.Node) error

	// SigReferenceMimeType optionally returns the MimeType to be used for a signature
	// reference (signature or timestamp). Port of the protected abstract
	// getSigReferenceMimeType().
	SigReferenceMimeType() enumerations.MimeType

	// AsicContentDocumentFilter gets the ASiCContentDocumentFilter used to filter the documents
	// to be referenced within the ASiC Manifest. Port of the protected
	// getAsicContentDocumentFilter().
	AsicContentDocumentFilter() *ASiCContentDocumentFilter

	// InitDefaultAsicContentDocumentFilter builds the default ASiCContentDocumentFilter for the
	// given manifest type. Port of the protected abstract
	// initDefaultAsicContentDocumentFilter().
	InitDefaultAsicContentDocumentFilter() *ASiCContentDocumentFilter

	// AddDataObjectReference adds a {@code <DataObjectReference>} element. Port of the
	// protected addDataObjectReference(Document, Element, DSSDocument, DigestAlgorithm).
	AddDataObjectReference(documentDom, asicManifestDom *xmldom.Node, doc model.DSSDocument,
		digestAlgorithm enumerations.DigestAlgorithm) (*xmldom.Node, error)

	// IsRootfile specifies whether the document is a Rootfile document. Port of the protected
	// isRootfile(DSSDocument).
	IsRootfile(doc model.DSSDocument) bool

	// ToDSSDocument transforms the DOM Document to a DSSDocument. Port of the protected
	// toDSSDocument(Document).
	ToDSSDocument(documentDom *xmldom.Node) (model.DSSDocument, error)

	// ManifestFilename returns a final filename of the manifest. Port of the protected abstract
	// getManifestFilename().
	ManifestFilename() string
}

// AbstractASiCManifestBuilder is the abstract class to build a Manifest for ASiC.
type AbstractASiCManifestBuilder struct {
	// overrides points back at the concrete builder; see InitAbstractASiCManifestBuilder.
	overrides AbstractASiCManifestBuilderOverrides

	// AsicContent is the container representation. Port of the protected final asicContent.
	AsicContent *ASiCContent

	// SigReferenceUri is the URI of a document signing the manifest. Port of the protected
	// final sigReferenceUri.
	SigReferenceUri string

	// DigestAlgorithm is the DigestAlgorithm to use for reference digests computation. Port of
	// the protected final digestAlgorithm.
	DigestAlgorithm enumerations.DigestAlgorithm

	// asicContentDocumentFilter is used to filter the documents to compute hashes for.
	asicContentDocumentFilter *ASiCContentDocumentFilter
}

// InitAbstractASiCManifestBuilder instantiates the builder with a default SHA-256 digest algorithm.
// Port of the protected AbstractASiCManifestBuilder(ASiCContent, String) constructor.
func (b *AbstractASiCManifestBuilder) InitAbstractASiCManifestBuilder(overrides AbstractASiCManifestBuilderOverrides,
	asicContent *ASiCContent, sigReferenceUri string) {
	b.InitAbstractASiCManifestBuilderWithDigestAlgorithm(overrides, asicContent, sigReferenceUri, enumerations.DigestAlgorithmSHA256)
}

// InitAbstractASiCManifestBuilderWithDigestAlgorithm instantiates the builder with a provided digest
// algorithm. Port of the protected AbstractASiCManifestBuilder(ASiCContent, String,
// DigestAlgorithm) constructor.
func (b *AbstractASiCManifestBuilder) InitAbstractASiCManifestBuilderWithDigestAlgorithm(overrides AbstractASiCManifestBuilderOverrides,
	asicContent *ASiCContent, sigReferenceUri string, digestAlgorithm enumerations.DigestAlgorithm) {
	b.overrides = overrides
	b.AsicContent = asicContent
	b.SigReferenceUri = sigReferenceUri
	b.DigestAlgorithm = digestAlgorithm
}

// Build builds the ArchiveManifest and returns the Document Node. Port of build().
func (b *AbstractASiCManifestBuilder) Build() (model.DSSDocument, error) {
	documentDom := b.overrides.BuildDom()
	asicManifestDom := b.overrides.CreateRootElement(documentDom)

	b.overrides.AddSigReference(documentDom, asicManifestDom)
	if err := b.overrides.AddDataObjectReferences(documentDom, asicManifestDom); err != nil {
		return nil, err
	}

	return b.overrides.ToDSSDocument(documentDom)
}

// BuildDom builds the initial XML document. Port of the protected buildDom().
func (b *AbstractASiCManifestBuilder) BuildDom() *xmldom.Node {
	return xmlutils.DomUtilsBuildDOMEmpty()
}

// CreateRootElement creates a root element
// {@code <asic:ASiCManifest xmlns:asic="http://uri.etsi.org/02918/v1.2.1#">}. Port of the protected
// createRootElement(Document).
func (b *AbstractASiCManifestBuilder) CreateRootElement(documentDom *xmldom.Node) *xmldom.Node {
	asicManifestDom := xmlutils.DomUtilsCreateElementNS(documentDom, ASiCManifestNS, ASiCManifestElementASiCManifest)
	documentDom.AppendChild(asicManifestDom)
	return asicManifestDom
}

// AddSigReference adds a {@code <SigReference>} element. Port of the protected
// addSigReference(Document, Element).
func (b *AbstractASiCManifestBuilder) AddSigReference(documentDom, asicManifestDom *xmldom.Node) {
	sigReferenceDom := xmlutils.DomUtilsAddElement(documentDom, asicManifestDom, ASiCManifestNS, ASiCManifestElementSigReference)
	sigReferenceDom.SetAttr(xmldom.Name{Local: ASiCManifestAttributeURI.AttributeName()}, spi.DSSUtilsEncodeURI(b.SigReferenceUri))
	sigReferenceMimeType := b.overrides.SigReferenceMimeType()
	if sigReferenceMimeType != nil {
		sigReferenceDom.SetAttr(xmldom.Name{Local: ASiCManifestAttributeMIMEType.AttributeName()}, sigReferenceMimeType.MimeTypeString())
	}
}

// AddDataObjectReferences adds references to data objects, corresponding to the
// ASiCContentDocumentFilter configuration.
//
// Panics with the Java message when the filter is nil (Objects.requireNonNull); returns an error
// when a filtered document has no name (IllegalArgumentException).
//
// Port of the protected addDataObjectReferences(Document, Element).
func (b *AbstractASiCManifestBuilder) AddDataObjectReferences(documentDom, asicManifestDom *xmldom.Node) error {
	documentFilter := b.overrides.AsicContentDocumentFilter()
	if documentFilter == nil {
		panic("ASiCContentDocumentFilter cannot be null!")
	}

	for _, doc := range documentFilter.Filter(b.AsicContent) {
		if err := b.assertDocumentNameDefined(doc); err != nil {
			return err
		}
		if _, err := b.overrides.AddDataObjectReference(documentDom, asicManifestDom, doc, b.DigestAlgorithm); err != nil {
			return err
		}
	}
	return nil
}

// assertDocumentNameDefined is the port of the private assertDocumentNameDefined(DSSDocument).
func (b *AbstractASiCManifestBuilder) assertDocumentNameDefined(doc model.DSSDocument) error {
	if doc.Name() == "" {
		return model.NewDSSError("Document name shall be defined for an ASiC Manifest building!")
	}
	return nil
}

// AsicContentDocumentFilter gets the ASiCContentDocumentFilter used to filter the documents to be
// referenced within the ASiC Manifest. Port of the protected getAsicContentDocumentFilter().
func (b *AbstractASiCManifestBuilder) AsicContentDocumentFilter() *ASiCContentDocumentFilter {
	if b.asicContentDocumentFilter == nil {
		b.asicContentDocumentFilter = b.overrides.InitDefaultAsicContentDocumentFilter()
	}
	return b.asicContentDocumentFilter
}

// SetAsicContentDocumentFilter sets the ASiCContentDocumentFilter used to filter the documents to
// compute hashes for. When not set, a default ASiCContentDocumentFilter is used for the given
// manifest type.
//
// Java returns `this` for chaining; Go returns the embedded struct pointer, which a subclass value
// cannot stand in for, so callers chain on their own receiver instead.
//
// Port of setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *AbstractASiCManifestBuilder) SetAsicContentDocumentFilter(asicContentDocumentFilter *ASiCContentDocumentFilter) *AbstractASiCManifestBuilder {
	b.asicContentDocumentFilter = asicContentDocumentFilter
	return b
}

// AddDataObjectReference adds a {@code <DataObjectReference>} element. Port of the protected
// addDataObjectReference(Document, Element, DSSDocument, DigestAlgorithm).
func (b *AbstractASiCManifestBuilder) AddDataObjectReference(documentDom, asicManifestDom *xmldom.Node,
	doc model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) (*xmldom.Node, error) {
	dataObjectReferenceDom := xmlutils.DomUtilsAddElement(documentDom, asicManifestDom, ASiCManifestNS, ASiCManifestElementDataObjectReference)

	dataObjectReferenceDom.SetAttr(xmldom.Name{Local: ASiCManifestAttributeURI.AttributeName()}, spi.DSSUtilsEncodeURI(doc.Name()))

	mimeType := doc.MimeType()
	if mimeType != nil {
		dataObjectReferenceDom.SetAttr(xmldom.Name{Local: ASiCManifestAttributeMIMEType.AttributeName()}, mimeType.MimeTypeString())
	}

	if b.overrides.IsRootfile(doc) {
		dataObjectReferenceDom.SetAttr(xmldom.Name{Local: ASiCManifestAttributeRootFile.AttributeName()}, "true")
	}

	digestMethodDom := xmlutils.DomUtilsAddElement(documentDom, dataObjectReferenceDom, common.XMLDSigNS, common.XMLDSigElementDigestMethod)
	digestMethodDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()}, digestAlgorithm.URI())

	digestValueDom := xmlutils.DomUtilsAddElement(documentDom, dataObjectReferenceDom, common.XMLDSigNS, common.XMLDSigElementDigestValue)
	digestValue, err := doc.DigestValue(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	textNode := xmldom.NewText(utils.ToBase64(digestValue))
	digestValueDom.AppendChild(textNode)

	return dataObjectReferenceDom, nil
}

// IsRootfile specifies whether the document is a Rootfile document. Port of the protected
// isRootfile(DSSDocument).
func (b *AbstractASiCManifestBuilder) IsRootfile(doc model.DSSDocument) bool {
	// FALSE by default
	return false
}

// ToDSSDocument transforms the DOM Document to a DSSDocument. Port of the protected
// toDSSDocument(Document).
func (b *AbstractASiCManifestBuilder) ToDSSDocument(documentDom *xmldom.Node) (model.DSSDocument, error) {
	newManifestName := b.overrides.ManifestFilename()
	return xmlutils.DomUtilsCreateDssDocumentFromDomDocument(documentDom, newManifestName)
}
