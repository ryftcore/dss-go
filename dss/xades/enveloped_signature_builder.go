// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/EnvelopedSignatureBuilder.java (DSS 6.5.RC1).
//
// Java's class is package-private and overrides two of the hooks collected in
// XAdESSignatureBuilderOverrides - assertSignaturePossible() and buildRootDocumentDom() - while
// inheriting getParentNodeOfSignature() and incorporateSignatureDom(Node) from
// XPathPlacementSignatureBuilder. Go has no method overriding across embedding, so this builder
// registers itself with InitXPathPlacementSignatureBuilder and the two it does not override are
// promoted from the embedded intermediate base, the TokenBase.InitToken(self) convention of
// PORTING.md.
//
// DEVIATION (flagged for the integrator, and consistent with detached_signature_builder.go): the
// Java class is package-private. The whole XAdES SCC is one Go package, so unexported would
// preserve that exactly, but the type is named by XAdESSignatureBuilderGetSignatureBuilder's
// factory in a sibling file and reads better under the "exported identifiers keep the Java name"
// rule of PORTING.md, so it is exported here. Nothing outside dss-xades constructs it upstream.
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// EnvelopedSignatureBuilder handles the specifics of the enveloped XML signature.
type EnvelopedSignatureBuilder struct {
	XPathPlacementSignatureBuilder
}

// NewEnvelopedSignatureBuilder is the constructor for a single-document signing. The enveloped
// signature uses by default the exclusive method of canonicalization.
// Port of EnvelopedSignatureBuilder(XAdESSignatureParameters, DSSDocument, CertificateVerifier).
func NewEnvelopedSignatureBuilder(params *XAdESSignatureParameters, document model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *EnvelopedSignatureBuilder {
	return NewEnvelopedSignatureBuilderForDocuments(params, []model.DSSDocument{document}, certificateVerifier)
}

// NewEnvelopedSignatureBuilderForDocuments is the constructor for multiple documents signing. The
// enveloped signature uses by default the exclusive method of canonicalization.
// Port of EnvelopedSignatureBuilder(XAdESSignatureParameters, List<DSSDocument>, CertificateVerifier).
func NewEnvelopedSignatureBuilderForDocuments(params *XAdESSignatureParameters,
	documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *EnvelopedSignatureBuilder {
	builder := &EnvelopedSignatureBuilder{}
	builder.InitXPathPlacementSignatureBuilder(builder, params, documents, certificateVerifier)
	return builder
}

// AssertSignaturePossible verifies that exactly one XML document was provided and that adding a
// parallel signature to it is possible. Port of the overridden protected #assertSignaturePossible.
func (b *EnvelopedSignatureBuilder) AssertSignaturePossible() error {
	if err := b.XAdESSignatureBuilder.AssertSignaturePossible(); err != nil {
		return err
	}
	return b.AssertOriginalXmlDocumentValid()
}

// BuildRootDocumentDom returns the parsed original document: in the case of an enveloped
// signature the document has to be the original file, which matters for inclusive
// canonicalization and namespaces.
// Port of the overridden protected #buildRootDocumentDom.
func (b *EnvelopedSignatureBuilder) BuildRootDocumentDom() *xmldom.Node {
	// Java's DomUtils.buildDOM(DSSDocument) throws DSSException on unparsable input; this hook
	// has no error channel, and AssertSignaturePossible has already established that the document
	// is XML (DomUtils.isDOM), so a failure here is unreachable and panics.
	documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(b.Documents[0])
	if err != nil {
		panic(err)
	}
	return documentDom
}
