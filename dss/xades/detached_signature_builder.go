// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/DetachedSignatureBuilder.java (DSS 6.5.RC1).
//
// Java's class is package-private and abstract-free: it only overrides the two hooks
// AbstractSignatureBuilder leaves for its packaging-specific subclasses, buildRootDocumentDom()
// and getParentNodeOfSignature(). Go has no method overriding across embedding, so the two are
// members of SignatureBuilderOverrides and this builder registers itself with
// InitXAdESSignatureBuilder, the TokenBase.InitToken(self) convention of PORTING.md.
//
// DEVIATION (flagged for the integrator): the Java class is package-private. Go's unexported
// form would preserve that exactly - the whole XAdES SCC is one Go package - but the type is
// named by AbstractSignatureBuilder.getSignatureBuilder's factory in a sibling file and reads far
// better under the "exported identifiers keep the Java name" rule of PORTING.md, so it is
// exported here. Nothing outside dss-xades constructs it upstream.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// DetachedSignatureBuilder handles the specifics of the detached XML signature.
type DetachedSignatureBuilder struct {
	AbstractSignatureBuilder
}

// NewDetachedSignatureBuilder is the constructor for a single-document signing. The detached
// signature uses by default the exclusive method of canonicalization.
// Port of DetachedSignatureBuilder(XAdESSignatureParameters, DSSDocument, CertificateVerifier).
func NewDetachedSignatureBuilder(params *SignatureParameters, document model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *DetachedSignatureBuilder {
	return NewDetachedSignatureBuilderForDocuments(params, []model.DSSDocument{document}, certificateVerifier)
}

// NewDetachedSignatureBuilderForDocuments is the constructor for multiple documents signing.
// The detached signature uses by default the exclusive method of canonicalization.
// Port of DetachedSignatureBuilder(XAdESSignatureParameters, List<DSSDocument>, CertificateVerifier).
func NewDetachedSignatureBuilderForDocuments(params *SignatureParameters,
	documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *DetachedSignatureBuilder {
	builder := &DetachedSignatureBuilder{}
	builder.InitXAdESSignatureBuilder(builder, params, documents, certificateVerifier)
	return builder
}

// BuildRootDocumentDom builds the root Document DOM the signature is created in: the root
// document handed over through the parameters when there is one, an empty document otherwise.
// Port of the overridden protected #buildRootDocumentDom.
func (b *DetachedSignatureBuilder) BuildRootDocumentDom() *xmldom.Node {
	if b.Params.RootDocument() != nil {
		return b.Params.RootDocument()
	}
	return xmlutils.DomUtilsBuildDOMEmpty()
}

// ParentNodeOfSignature returns the node the ds:Signature element is appended to: the document
// element when signing into an existing root document, the document node itself otherwise.
// Port of the overridden protected #getParentNodeOfSignature.
func (b *DetachedSignatureBuilder) ParentNodeOfSignature() *xmldom.Node {
	if b.Params.RootDocument() != nil {
		return b.DocumentDom.DocumentElement()
	}
	return b.DocumentDom
}
