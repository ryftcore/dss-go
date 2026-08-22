// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ComplexTransform.java (DSS 6.5.RC1).
//
// # Santuario replacement
//
// Java's ComplexTransform is the one transform family that DSS does not execute itself: it
// hands the work to org.apache.xml.security.transforms.Transform. Per internal/xmldsig's
// doc.go mapping table that becomes, member for member:
//
//	org.apache.xml.security.transforms.Transform (the algorithm registry)   xmldsig.Registry / DefaultRegistry
//	Transform#performTransform(XMLSignatureInput, boolean)                  xmldsig.Registry.Perform
//	org.apache.xml.security.signature.XMLSignatureInput                     *xmldsig.Data
//
// Nothing here canonicalizes, digests or serializes on its own; every octet a ComplexTransform
// produces comes out of internal/xmldsig, which is the Santuario byte-exact core.
//
// # buildTransformObject
//
// Santuario's Transform(Document, String algorithmURI, NodeList contextNodes) constructor does
// three things this port reproduces literally, because they are what the transform SPIs then
// read:
//
//  1. it creates a FRESH ds:Transform element - not the one DSS just built into transformsDom -
//     in the XMLDSIG namespace under the prefix ElementProxy.registerDefaultPrefixes binds to
//     it, which SantuarioInitializer.init() leaves at the stock "ds", and declares that prefix
//     on the element itself (ElementProxy#createElementForFamilyLocal);
//  2. it copies the Algorithm attribute and DEEP CLONES the context nodes onto it - so the
//     ds:XPath / xpf:XPath parameter children the concrete subclass wrote are what the SPI
//     sees;
//  3. DSS then calls setXPathNamespaceContext for every binding in XPathUtils's namespace
//     context map, which is ElementProxy#setXPathNamespaceContext: an xmlns:<prefix>
//     declaration set on that same element.
//
// Step 3 is load-bearing rather than cosmetic. xmldsig's ds:XPath and XPath Filter 2.0
// transforms resolve the prefixes of their expression against the declarations in scope on the
// XPath element (xpath10.NamespaceContextOf), exactly as Santuario's DOMNamespaceContext does,
// and the XPath element is a child of this ds:Transform - so "not(ancestor-or-self::ds:Signature)"
// resolves precisely because the registered "ds" binding was written here.
//
// The prefixes are applied in sorted order. Java iterates a HashMap, which fixes only the
// attribute order on a throwaway element that is never serialized or canonicalized, so sorting
// changes nothing observable and keeps the port deterministic (per the repository's determinism
// convention).
//
// # Errors
//
// Java throws DSSException from both methods; those become returned errors carrying the same
// messages. Santuario raises InvalidTransformException from
// the constructor when the algorithm URI has no registered TransformSpi - which is how
// CanonicalizationTransform(physical c14n) fails, since XMLCanonicalizer can canonicalize with
// the Santuario "physical" method but Santuario never registers it as a TRANSFORM - so the
// registry lookup is done here, at buildTransformObject time, where Java detects it.
package xades

import (
	"fmt"
	"sort"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// complexTransformRegistry is the transform-algorithm registry Santuario keeps statically on
// org.apache.xml.security.transforms.Transform and Transform#registerDefaultAlgorithms fills
// (SantuarioInitializer.init() -> Transform.registerDefaultAlgorithms()). One instance for the
// package, as in upstream.
var complexTransformRegistry = xmldsig.DefaultRegistry()

// complexTransformDefaultPrefix is the prefix ElementProxy.registerDefaultPrefixes binds to
// the XMLDSIG namespace and therefore the prefix Santuario gives the ds:Transform element it
// builds internally, whatever prefix the signature itself uses.
const complexTransformDefaultPrefix = "ds"

// ComplexTransform is a transform processed by the Santuario replacement in internal/xmldsig.
//
// Java declares it abstract; Go has no abstract structs, so the type is instantiable but never
// instantiated directly - CanonicalizationTransform, XPathTransform and XsltTransform embed it,
// and each registers itself through InitAbstractTransform (see abstract_transform.go).
type ComplexTransform struct {
	AbstractTransform

	// transformElement is the internal object used to build the transformation. Port of the
	// private Transform transformObject field, which xmldsig models as the ds:Transform
	// element the registry is driven with.
	transformElement *xmldom.Node
}

// newComplexTransform ports the protected ComplexTransform(DSSNamespace, String).
func newComplexTransform(xmlDSigNamespace *common.DSSNamespace, algorithm string) ComplexTransform {
	return ComplexTransform{AbstractTransform: newAbstractTransformWithNamespace(xmlDSigNamespace, algorithm)}
}

// transformObject ports the private getTransformObject(): build once, then reuse. Java caches
// the Transform in the same way, and - as upstream - a SetNamespace call made after the first
// use does not invalidate the cache.
func (t *ComplexTransform) transformObject() (*xmldom.Node, error) {
	if t.transformElement == nil {
		element, err := t.BuildTransformObject()
		if err != nil {
			return nil, err
		}
		t.transformElement = element
	}
	return t.transformElement, nil
}

// BuildTransformObject builds the ds:Transform element that drives the xmldsig registry.
// Ports the protected buildTransformObject(); see the file header for the three Santuario
// constructor steps it reproduces.
func (t *ComplexTransform) BuildTransformObject() (*xmldom.Node, error) {
	document := xmlutils.DomUtilsBuildDOMEmpty()
	transformsDom := xmlutils.DomUtilsCreateElementNS(document, t.namespace, common.XMLDSigElementTransforms)
	document.AppendChild(transformsDom)
	t.self.CreateTransform(document, transformsDom)

	// Santuario's Transform constructor resolves the TransformSpi eagerly and raises
	// InvalidTransformException when the URI is not registered; DSS wraps that in a
	// DSSException with this message.
	if _, ok := complexTransformRegistry.Lookup(t.algorithm); !ok {
		return nil, model.NewDSSError(fmt.Sprintf("Cannot initialize a transform [%s]", t.algorithm))
	}

	transform := xmldom.NewElement(xmldom.Name{
		Space:  common.XMLDSigNS.Uri(),
		Local:  common.XMLDSigElementTransform.TagName(),
		Prefix: complexTransformDefaultPrefix,
	})
	transform.SetAttr(xmldom.Name{
		Space:  xmldom.XMLNSNamespace,
		Local:  complexTransformDefaultPrefix,
		Prefix: "xmlns",
	}, common.XMLDSigNS.Uri())
	transform.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()}, t.algorithm)

	// transformsDom.getFirstChild().getChildNodes(): the parameter children of the
	// ds:Transform the subclass just created, deep-cloned onto the new element.
	for child := transformsDom.FirstChild.FirstChild; child != nil; child = child.NextSibling {
		transform.AppendChild(child.Clone(true))
	}

	prefixMap := xmlutils.XPathUtilsGetNamespaceContextMap().PrefixMap()
	prefixes := make([]string, 0, len(prefixMap))
	for prefix := range prefixMap {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	for _, prefix := range prefixes {
		if err := complexTransformSetXPathNamespaceContext(transform, prefix, prefixMap[prefix]); err != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Cannot initialize a transform [%s]", t.algorithm), err)
		}
	}
	return transform, nil
}

// PerformTransform executes the transform. Ports performTransform(DSSTransformOutput), which
// calls Transform#performTransform(input, true), i.e. with Santuario's secureValidation flag
// on. That flag reaches the TransformSpi and only hardens it (TransformXSLT closes off external
// DTDs and stylesheets with it); the outright refusal of the XSLT algorithm lives in
// Transforms#checkSecureValidation, which is the VALIDATION-side chain runner and is not on
// this path - so upstream really does run a stylesheet here. See xslt_transform.go for the one
// place where this port therefore diverges.
func (t *ComplexTransform) PerformTransform(transformOutput *DSSTransformOutput) (*DSSTransformOutput, error) {
	transform, err := t.transformObject()
	if err != nil {
		return nil, err
	}
	out, err := complexTransformRegistry.Perform(transformOutput.xmlSignatureInput(), transform, "", true)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Cannot process transformation [%s] on the given DOM object. Reason : [%s]",
			t.algorithm, err.Error()), err)
	}
	return NewDSSTransformOutput(out), nil
}

// complexTransformSetXPathNamespaceContext ports
// org.apache.xml.security.utils.ElementProxy#setXPathNamespaceContext.
//
// The conflict check upstream performs before writing - getAttributeNodeNS(xmlns-URI,
// "xmlns:<prefix>") - can never fire, because a namespace declaration's DOM local name is the
// bare prefix, not the "xmlns:"-qualified name it is looked up under. The lookup is therefore
// omitted rather than reproduced as dead code, and the write is unconditional, which is what
// upstream ends up doing for every binding.
func complexTransformSetXPathNamespaceContext(element *xmldom.Node, prefix, uri string) error {
	if prefix == "" || prefix == "xmlns" {
		return fmt.Errorf("defaultNamespaceCannotBeSetHere")
	}
	local := prefix
	if len(prefix) > len("xmlns:") && prefix[:len("xmlns:")] == "xmlns:" {
		local = prefix[len("xmlns:"):]
	}
	element.SetAttr(xmldom.Name{Space: xmldom.XMLNSNamespace, Local: local, Prefix: "xmlns"}, uri)
	return nil
}
