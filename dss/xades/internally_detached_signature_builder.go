// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/InternallyDetachedSignatureBuilder.java (DSS 6.5.RC1).
//
// Java's class is package-private and overrides three of the hooks collected in
// XAdESSignatureBuilderOverrides - assertSignaturePossible(), buildRootDocumentDom() and
// incorporateFiles() - while inheriting getParentNodeOfSignature() and
// incorporateSignatureDom(Node) from XPathPlacementSignatureBuilder. Go has no method overriding
// across embedding, so this builder registers itself with InitXPathPlacementSignatureBuilder and
// the hooks it does not override are promoted from the embedded intermediate base, the
// TokenBase.InitToken(self) convention of PORTING.md.
//
// DEVIATION (flagged for the integrator, and consistent with detached_signature_builder.go): the
// Java class is package-private but is exported here; see enveloped_signature_builder.go for the
// reasoning.
//
// slf4j logging is dropped (PORTING.md).
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// internallyDetachedSignatureBuilderDefaultSignatureContainerName defines the name of the root
// signature container element, used when the root element is not provided.
// Port of the private DEFAULT_SIGNATURE_CONTAINER_NAME.
const internallyDetachedSignatureBuilderDefaultSignatureContainerName = "internally-detached"

// InternallyDetachedSignatureBuilder handles the specifics of the internally detached XML
// signature.
type InternallyDetachedSignatureBuilder struct {
	XPathPlacementSignatureBuilder
}

// NewInternallyDetachedSignatureBuilder is the constructor for a single-document signing. The
// internally detached signature uses by default the exclusive method of canonicalization.
// Port of InternallyDetachedSignatureBuilder(XAdESSignatureParameters, DSSDocument, CertificateVerifier).
func NewInternallyDetachedSignatureBuilder(params *XAdESSignatureParameters, document model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *InternallyDetachedSignatureBuilder {
	return NewInternallyDetachedSignatureBuilderForDocuments(params, []model.DSSDocument{document},
		certificateVerifier)
}

// NewInternallyDetachedSignatureBuilderForDocuments is the constructor for multiple documents
// signing. The internally detached signature uses by default the exclusive method of
// canonicalization.
// Port of InternallyDetachedSignatureBuilder(XAdESSignatureParameters, List<DSSDocument>, CertificateVerifier).
func NewInternallyDetachedSignatureBuilderForDocuments(params *XAdESSignatureParameters,
	documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *InternallyDetachedSignatureBuilder {
	builder := &InternallyDetachedSignatureBuilder{}
	builder.InitXPathPlacementSignatureBuilder(builder, params, documents, certificateVerifier)
	return builder
}

// AssertSignaturePossible verifies the documents allow an internally-detached signature; the
// parallel-signature check only applies when the signature is placed into the original document
// by XPath. Port of the overridden protected #assertSignaturePossible.
func (b *InternallyDetachedSignatureBuilder) AssertSignaturePossible() error {
	if err := b.XAdESSignatureBuilder.AssertSignaturePossible(); err != nil {
		return err
	}

	if b.Params.RootDocument() == nil && utils.IsStringNotEmpty(b.Params.XPathLocationString()) {
		return b.AssertOriginalXmlDocumentValid()
	}
	return nil
}

// BuildRootDocumentDom returns the root Document the signature is created in: the one handed over
// through the parameters, the parsed original document when an XPath location selects a place
// inside it, and otherwise a freshly created default container.
// Port of the overridden protected #buildRootDocumentDom.
func (b *InternallyDetachedSignatureBuilder) BuildRootDocumentDom() *xmldom.Node {
	if b.Params.RootDocument() != nil {
		return b.Params.RootDocument()
	} else if utils.IsStringNotEmpty(b.Params.XPathLocationString()) {
		// Java's DomUtils.buildDOM(DSSDocument) throws DSSException on unparsable input; this hook
		// has no error channel, and AssertSignaturePossible has already established that the
		// document is XML in exactly this branch, so a failure here is unreachable and panics.
		documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(b.Documents[0])
		if err != nil {
			panic(err)
		}
		return documentDom
	}
	return internallyDetachedSignatureBuilderCreateDefaultContainer()
}

// internallyDetachedSignatureBuilderCreateDefaultContainer ports the private
// createDefaultContainer. Java's Document#createElement builds a namespace-less element, which
// xmldom spells as a Name carrying only a local part.
func internallyDetachedSignatureBuilderCreateDefaultContainer() *xmldom.Node {
	rootDocument := xmlutils.DomUtilsBuildDOMEmpty()
	rootElement := xmldom.NewElement(
		xmldom.Name{Local: internallyDetachedSignatureBuilderDefaultSignatureContainerName})
	rootDocument.AppendChild(rootElement)
	return rootDocument
}

// IncorporateFiles adds the referenced contents into the signing document, skipping any element
// whose Id is already present. Port of the overridden protected #incorporateFiles.
func (b *InternallyDetachedSignatureBuilder) IncorporateFiles() error {
	references := b.Params.References()
	for _, ref := range references {
		elementId := xmlutils.DomUtilsGetId(ref.Uri())
		// the content shall be added only when it is not yet present in the document
		if xmlutils.XPathUtilsGetElementById(b.DocumentDom, elementId) == nil {
			doc, err := xmlutils.DomUtilsBuildDOMFromDocument(ref.Contents())
			if err != nil {
				return err
			}
			root := doc.DocumentElement()
			adopted := b.DocumentDom.Import(root, true)
			b.DocumentDom.DocumentElement().AppendChild(adopted)

		}
		// Upstream logs "The element with Id '{}' is already present in the signing document!
		// The addition is skipped." in the else branch.
	}
	return nil
}
