// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/CounterSignatureResolver.java
// (DSS 6.5.RC1).
//
// Resolver for a counter signature only, used for a counter signature extension. Per
// internal/xmldsig's doc.go ("What is out of scope" section), this class is explicitly NOT part
// of that package - it needs xml/utils.DomUtilsSerializeNode and XPathUtilsGetElementByIdWithQuery,
// both of which live in DSS packages internal/ may not import - and is therefore the XAdES
// port's own responsibility, implementing internal/xmldsig.URIResolver directly.
package xades

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// CounterSignatureResolver resolves the counter-signed ds:SignatureValue document only. Port of
// the class CounterSignatureResolver, extending org.apache.xml.security's ResourceResolverSpi;
// satisfies internal/xmldsig.URIResolver.
type CounterSignatureResolver struct {
	// document is the counter signed SignatureValue document.
	document model.DSSDocument
}

// NewCounterSignatureResolver is the port of the constructor CounterSignatureResolver(DSSDocument).
func NewCounterSignatureResolver(document model.DSSDocument) *CounterSignatureResolver {
	return &CounterSignatureResolver{document: document}
}

// Resolve ports engineResolveURI(ResourceResolverContext).
func (r *CounterSignatureResolver) Resolve(ctx *xmldsig.ResolverContext) (*xmldsig.Data, error) {
	uriValue := r.uriValue(ctx)
	node := r.resolveNode(uriValue)

	if node != nil {
		return r.createFromNode(node), nil
	}

	return nil, fmt.Errorf("%w: unable to find a signed content by URI : '%s'", errCounterSignatureResolverNotFound, uriValue)
}

// errCounterSignatureResolverNotFound is the ResourceResolverException("generic.EmptyMessage")
// CounterSignatureResolver raises when the URI cannot be resolved.
var errCounterSignatureResolverNotFound = errors.New("xades: unable to find a signed content by URI")

// createFromNode ports the private createFromNode(Node).
func (r *CounterSignatureResolver) createFromNode(node *xmldom.Node) *xmldsig.Data {
	serialized, err := xmlutils.DomUtilsSerializeNode(node)
	if err != nil {
		// DomUtils.serializeNode declares no checked exception in Java; a Go serialization
		// failure here has no legitimate cause other than a malformed in-memory DOM, treated the
		// same "cannot happen in practice" way xades_attribute_identifier.go documents.
		panic(fmt.Sprintf("Unable to serialize node : %s", err.Error()))
	}
	result := xmldsig.NewOctetData(serialized)
	result.SetMIMEType(enumerations.MimeTypeEnum_XML.MimeTypeString())
	return result
}

// isXPointerSlash ports the private isXPointerSlash(String).
func (r *CounterSignatureResolver) isXPointerSlash(uri string) bool {
	return uri == "#xpointer(/)"
}

// CanResolve ports engineCanResolveURI(ResourceResolverContext).
func (r *CounterSignatureResolver) CanResolve(ctx *xmldsig.ResolverContext) bool {
	uriValue := r.uriValue(ctx)
	return (xmlutils.DomUtilsIsXPointerQuery(uriValue) || xmlutils.DomUtilsIsElementReference(uriValue)) && r.resolveNode(uriValue) != nil
}

// uriValue ports the private getURIValue(ResourceResolverContext).
func (r *CounterSignatureResolver) uriValue(ctx *xmldsig.ResolverContext) string {
	if ctx.Attr == nil {
		return ""
	}
	return spi.DSSUtilsDecodeURI(ctx.Attr.Value)
}

// resolveNode ports the private resolveNode(String).
func (r *CounterSignatureResolver) resolveNode(uriValue string) *xmldom.Node {
	if uriValue == "" {
		return nil
	}

	documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(r.document)
	if err != nil {
		return nil
	}
	node := xmlutils.XPathUtilsGetElementByIdWithQuery(documentDom, common.XMLDSigPath_ALL_SIGNATURE_VALUES_PATH, uriValue)

	if node == nil && r.isXPointerSlash(uriValue) && common.XMLDSigElement_SIGNATURE_VALUE.IsSameTagName(documentDom.Name.Local) {
		node = documentDom
	}

	return node
}
