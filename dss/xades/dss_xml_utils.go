// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/DSSXMLUtils.java (DSS 6.5.RC1).
//
// # Santuario replacement
//
// org.apache.xml.security.signature.{Reference,Manifest,XMLSignatureInput} are xmldsig.Reference/
// Manifest/Data per internal/xmldsig's doc.go mapping table; org.apache.xml.security.keys.KeyInfo
// and org.apache.xml.security.transforms.{Transform,Transforms} have no internal/xmldsig
// counterpart (that package covers only reading an existing signature's references, not KeyInfo
// public-key extraction) and are handled directly against *xmldom.Node/encoding-asn1/crypto/x509
// below. SantuarioInitializer.init() and the JCEMapper/ResourceResolver static registration have
// no Go counterpart (see xades_signature.go's file header) and this file's static block is
// therefore reduced to the transform-URI registries and DSSXMLUtilsRegisterXAdESNamespaces,
// invoked once from xades_signature.go's init().
//
// # FLAGGED FOR INTEGRATOR
//
//   - DSSXMLUtilsGetDocument(reference) always returns nil: dss_document_xml_signature_input.go's
//     file header already documents why the Java type-assertion
//     (XMLSignatureInput -> DSSDocumentXMLSignatureInput -> getDocument()) has no Go equivalent
//     against the frozen internal/xmldsig.Data type, which carries no DSSDocument backreference.
//   - DSSXMLUtilsValidateAgainstXSD always returns nil (no errors): xades_structure_validator.go's
//     file header already documents the same blocked forward dependency (no bundled-schema or
//     javax.xml.validation-equivalent Go package exists yet in this port).
//   - DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey supports ds:X509Data/ds:X509Certificate and
//     ds:KeyValue/ds:RSAKeyValue only; DSAKeyValue and ECKeyValue extraction (upstream's
//     org.apache.xml.security.keys.KeyInfo#getPublicKey covers all three) are out of scope for
//     this pass - a signature whose KeyInfo carries only a DSA/EC KeyValue (no certificate) will
//     not resolve its signing-certificate candidate through this path.
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug call site noted in a comment
// becomes a Go catch-and-continue (best effort), matching every other file of this port.
package xades

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"math/big"
	"strings"
	"sync"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// DSSXMLUtilsTransformerIndentNumber is the value used to pretty print a xades signature. Ports
// TRANSFORMER_INDENT_NUMBER.
const DSSXMLUtilsTransformerIndentNumber = 4

// dssXMLUtilsTransformationExcludeSignature is the Enveloped-signature transformation. Ports the
// private TRANSFORMATION_EXCLUDE_SIGNATURE.
const dssXMLUtilsTransformationExcludeSignature = "not(ancestor-or-self::ds:Signature)"

// dssXMLUtilsTransformationXPathNodeName is the XPath transform name. Ports the private
// TRANSFORMATION_XPATH_NODE_NAME.
const dssXMLUtilsTransformationXPathNodeName = "XPath"

// DSSXMLUtilsSPDocDigestAsInSpecificationAlgorithmURI is the SPDocDigestAsInSpecification
// transform algorithm URI for a custom SignaturePolicy processing. Ports
// SP_DOC_DIGEST_AS_IN_SPECIFICATION_ALGORITHM_URI.
const DSSXMLUtilsSPDocDigestAsInSpecificationAlgorithmURI = "http://uri.etsi.org/01903/v1.3.2/SignaturePolicy/SPDocDigestAsInSpecification"

// DSSXMLUtilsSAMLNamespace is the SAML namespace definition. Ports SAML_NAMESPACE.
var DSSXMLUtilsSAMLNamespace = common.NewDSSNamespace("urn:oasis:names:tc:SAML:2.0:assertion", "saml2")

var (
	dssXMLUtilsMu sync.Mutex

	// dssXMLUtilsTransforms is the set of supported transforms. Ports the private static
	// transforms, seeded by registerDefaultTransforms().
	dssXMLUtilsTransforms = map[string]struct{}{
		string(xmldsig.TransformBase64Decode):       {},
		string(xmldsig.TransformEnvelopedSignature): {},
		string(xmldsig.TransformXPath):              {},
		string(xmldsig.TransformXPath2Filter):       {},
		string(xmldsig.TransformXPointer):           {},
		string(xmldsig.TransformXSLT):               {},
	}

	// dssXMLUtilsTransformsWithNodeSetOutput is the set of transforms resulting in a NodeSet
	// output. Ports the private static transformsWithNodeSetOutput, seeded by
	// registerTransformsWithNodeSetOutput().
	dssXMLUtilsTransformsWithNodeSetOutput = map[string]struct{}{
		string(xmldsig.TransformEnvelopedSignature): {},
		string(xmldsig.TransformXPath):              {},
		string(xmldsig.TransformXPath2Filter):       {},
	}
)

// DSSXMLUtilsRegisterTransform registers a transformation, reporting whether the set did not
// already contain it. Ports registerTransform(String).
func DSSXMLUtilsRegisterTransform(transformURI string) bool {
	dssXMLUtilsMu.Lock()
	defer dssXMLUtilsMu.Unlock()
	if _, ok := dssXMLUtilsTransforms[transformURI]; ok {
		return false
	}
	dssXMLUtilsTransforms[transformURI] = struct{}{}
	return true
}

// DSSXMLUtilsRegisterTransformWithNodeSetOutput registers a transformation resulting in a
// node-set output, reporting whether the set did not already contain it. Ports
// registerTransformWithNodeSetOutput(String).
func DSSXMLUtilsRegisterTransformWithNodeSetOutput(transformURI string) bool {
	dssXMLUtilsMu.Lock()
	defer dssXMLUtilsMu.Unlock()
	if _, ok := dssXMLUtilsTransformsWithNodeSetOutput[transformURI]; ok {
		return false
	}
	dssXMLUtilsTransformsWithNodeSetOutput[transformURI] = struct{}{}
	return true
}

// DSSXMLUtilsRegisterXAdESNamespaces registers the XAdES namespaces for XPath evaluation. Ports
// registerXAdESNamespaces().
func DSSXMLUtilsRegisterXAdESNamespaces() {
	xmlutils.XPathUtilsRegisterNamespace(common.XMLDSigNS)

	xmlutils.XPathUtilsRegisterNamespace(definition.XAdESNamespace_XADES_111)
	xmlutils.XPathUtilsRegisterNamespace(definition.XAdESNamespace_XADES_122)
	xmlutils.XPathUtilsRegisterNamespace(definition.XAdESNamespace_XADES_132)
	xmlutils.XPathUtilsRegisterNamespace(definition.XAdESNamespace_XADES_141)
	xmlutils.XPathUtilsRegisterNamespace(definition.XAdESNamespace_XADES_EVIDENCERECORD_NAMESPACE)
	// DO NOT register "xades"

	xmlutils.XPathUtilsRegisterNamespace(DSSXMLUtilsSAMLNamespace)
}

// DSSXMLUtilsIndentAndReplace indents the given node and replaces it with a new one on the
// document. Ports indentAndReplace(Document, Node).
func DSSXMLUtilsIndentAndReplace(documentDom, node *xmldom.Node) (*xmldom.Node, error) {
	indentedNode, err := DSSXMLUtilsGetIndentedNode(documentDom, node)
	if err != nil {
		return nil, err
	}
	importedNode := documentDom.Import(indentedNode, true)
	node.Parent.ReplaceChild(importedNode, node)
	return importedNode, nil
}

// DSSXMLUtilsIndentAndExtend extends the given oldNode by appending new indented children from
// the given newNode. Ports indentAndExtend(Document, Node, Node).
func DSSXMLUtilsIndentAndExtend(documentDom, newNode, oldNode *xmldom.Node) (*xmldom.Node, error) {
	indentedNode, err := DSSXMLUtilsGetIndentedNode(documentDom, newNode)
	if err != nil {
		return nil, err
	}
	indentedNode, err = DSSXMLUtilsAlignChildrenIndents(indentedNode)
	if err != nil {
		return nil, err
	}
	importedNode := documentDom.Import(indentedNode, true)
	children := dssXMLUtilsChildNodes(importedNode)
	startPosition := dssXMLUtilsGetPositionToStartExtension(oldNode, importedNode)
	for i := startPosition; i < len(children); i++ {
		nodeToAppend := children[i].Clone(true)
		if nodeToAppend.Kind != xmldom.Element || !dssXMLUtilsCheckIfExists(oldNode, nodeToAppend) {
			oldNode.AppendChild(nodeToAppend)
		}
	}
	newNode.Parent.ReplaceChild(oldNode, newNode)
	return oldNode, nil
}

// dssXMLUtilsGetPositionToStartExtension ports the private static
// getPositionToStartExtension(Node, Node).
func dssXMLUtilsGetPositionToStartExtension(oldNode, indentedNode *xmldom.Node) int {
	startPosition := len(dssXMLUtilsChildNodes(oldNode))
	var child *xmldom.Node
	for oldNode.FirstChild != nil {
		child = oldNode.LastChild
		if child.Kind == xmldom.Text {
			oldNode.RemoveChild(child)
		} else {
			break
		}
	}
	if position, ok := dssXMLUtilsGetPosition(indentedNode, child); ok {
		return position
	}
	return startPosition
}

// dssXMLUtilsCheckIfExists ports the private static checkIfExists(Node, Node).
func dssXMLUtilsCheckIfExists(parentNode, childToCheck *xmldom.Node) bool {
	_, ok := dssXMLUtilsGetPosition(parentNode, childToCheck)
	return ok
}

// dssXMLUtilsGetPosition ports the private static getPosition(Node, Node). Java's
// childToCheck.getLocalName() can be null for a non-element childToCheck, which the "Integer"
// return value never exercises against a null comparison target in practice (the only callers
// pass an element or the deliberately-final non-text `child` from
// getPositionToStartExtension) - reproduced here with Name.Local's zero-value "" standing in for
// Java's null, which compares safely rather than risking a NullPointerException.
func dssXMLUtilsGetPosition(parentNode, childToCheck *xmldom.Node) (int, bool) {
	if parentNode == nil || childToCheck == nil {
		return 0, false
	}
	nodeName := childToCheck.Name.Local
	i := 0
	for newChildNode := parentNode.FirstChild; newChildNode != nil; newChildNode = newChildNode.NextSibling {
		if nodeName == newChildNode.Name.Local {
			idIdentifier := DSSXMLUtilsGetIDIdentifier(childToCheck)
			if idIdentifier == "" || idIdentifier == DSSXMLUtilsGetIDIdentifier(newChildNode) {
				return i + 1, true
			}
		}
		i++
	}
	return 0, false
}

// dssXMLUtilsChildNodes snapshots a node's children into a slice, standing in for
// Node#getChildNodes()'s indexable NodeList.
func dssXMLUtilsChildNodes(node *xmldom.Node) []*xmldom.Node {
	var children []*xmldom.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

// DSSXMLUtilsGetDocWithIndentedSignature pretty prints a signature in the given document. Ports
// getDocWithIndentedSignature(Document, String, List<String>).
//
// DEVIATION preserved verbatim: upstream re-queries the ORIGINAL (by-then-detached) `signature`
// node for its new UnsignedSignatureProperties location, rather than `importedSignature`; the
// second lookup can only ever find the untouched, pre-indentation node inside the detached
// subtree, making the whole `if (unsignedSignatureProperties != null)` branch a no-op against
// the live document. Reproduced exactly rather than fixed, per PORTING.md.
func DSSXMLUtilsGetDocWithIndentedSignature(documentDom *xmldom.Node, signatureId string, noIndentObjectIds []string) (*xmldom.Node, error) {
	signatures, err := xmlutils.XPathUtilsGetNodeList(documentDom, common.XMLDSigPath_ALL_SIGNATURES_PATH)
	if err != nil {
		return nil, err
	}
	for _, signature := range signatures {
		signatureAttrIdValue := DSSXMLUtilsGetIDIdentifier(signature)
		if utils.IsStringNotEmpty(signatureAttrIdValue) && strings.Contains(signatureAttrIdValue, signatureId) {
			unsignedSignatureProperties, err := xmlutils.XPathUtilsGetNode(signature,
				common.AllFromCurrentPosition(definition.XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES))
			if err != nil {
				return nil, err
			}
			indentedSignature := dssXMLUtilsGetIndentedSignature(signature, noIndentObjectIds)
			importedSignature := documentDom.Import(indentedSignature, true)
			signature.Parent.ReplaceChild(importedSignature, signature)
			if unsignedSignatureProperties != nil {
				newUnsignedSignatureProperties, err := xmlutils.XPathUtilsGetNode(signature,
					common.AllFromCurrentPosition(definition.XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES))
				if err != nil {
					return nil, err
				}
				if newUnsignedSignatureProperties != nil && newUnsignedSignatureProperties.Parent != nil {
					newUnsignedSignatureProperties.Parent.ReplaceChild(unsignedSignatureProperties, newUnsignedSignatureProperties)
				}
			}
		}
	}
	return documentDom, nil
}

// dssXMLUtilsGetIndentedSignature ports the private static getIndentedSignature(Node,
// List<String>).
func dssXMLUtilsGetIndentedSignature(signature *xmldom.Node, noIndentObjectIds []string) *xmldom.Node {
	indentedSignature := dssXMLUtilsPrettyPrintNode(signature)
	for sigChild := signature.FirstChild; sigChild != nil; sigChild = sigChild.NextSibling {
		if sigChild.Kind != xmldom.Element {
			continue
		}
		idAttribute := DSSXMLUtilsGetIDIdentifier(sigChild)
		if dssXMLUtilsContainsString(noIndentObjectIds, idAttribute) {
			nodeToReplace := xmlutils.XPathUtilsGetElementById(indentedSignature, idAttribute)
			importedNode := indentedSignature.OwnerDocument().Import(sigChild, true)
			indentedSignature.ReplaceChild(importedNode, nodeToReplace)
		}
	}
	return indentedSignature
}

// dssXMLUtilsContainsString reports whether value is present in values, standing in for
// List#contains.
func dssXMLUtilsContainsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// DSSXMLUtilsGetIndentedNode returns an indented xmlNode. Ports getIndentedNode(Node, Node).
func DSSXMLUtilsGetIndentedNode(documentDom, xmlNode *xmldom.Node) (*xmldom.Node, error) {
	signatures, err := xmlutils.XPathUtilsGetNodeList(documentDom, common.XMLDSigPath_ALL_SIGNATURES_PATH)
	if err != nil {
		return nil, err
	}

	// TODO handle by namespace
	element := definition.XAdES132ElementFromTagName(xmlNode.Name.Local)
	if element == nil {
		element = common.DSSElementFromDefinition(xmlNode.Name.Local, common.NewDSSNamespace(xmlNode.Name.Space, xmlNode.Name.Prefix))
	}
	pathAllFromCurrentPosition := common.XPathQueryBuilderAllFromCurrentPosition().Element(element).Build()

	for _, signature := range signatures {
		idAttribute := DSSXMLUtilsGetIDIdentifier(xmlNode)
		var candidateList []*xmldom.Node
		if idAttribute != "" {
			candidateList, err = xmlutils.XPathUtilsGetNodeList(signature,
				common.XPathQueryBuilderAllFromCurrentPosition().IdValue(idAttribute).Build())
		} else {
			candidateList, err = xmlutils.XPathUtilsGetNodeList(signature, pathAllFromCurrentPosition)
		}
		if err != nil {
			continue
		}
		// TODO : review necessity of node comparison
		if dssXMLUtilsIsNodeListContains(candidateList, xmlNode) {
			indentedSignature := dssXMLUtilsPrettyPrintNode(signature)
			var indentedXmlNode *xmldom.Node
			if idAttribute != "" {
				indentedXmlNode = xmlutils.XPathUtilsGetElementById(indentedSignature, idAttribute)
			} else {
				indentedXmlNodes, err := xmlutils.XPathUtilsGetNodeList(indentedSignature, pathAllFromCurrentPosition)
				if err != nil {
					return nil, err
				}
				if len(indentedXmlNodes) == 0 {
					return nil, fmt.Errorf("no elements found matching the '%s' XPath expression!", pathAllFromCurrentPosition.QueryString())
				}
				// return the last item
				indentedXmlNode = indentedXmlNodes[len(indentedXmlNodes)-1]
			}
			if indentedXmlNode != nil {
				return indentedXmlNode, nil
			}
		}
	}
	return xmlNode, nil
}

// dssXMLUtilsPrettyPrintNode ports the private static getIndentedNode(Node).
func dssXMLUtilsPrettyPrintNode(xmlNode *xmldom.Node) *xmldom.Node {
	return NewPrettyPrintTransformer().Transform(xmlNode)
}

// dssXMLUtilsIsNodeListContains ports the private static isNodeListContains(NodeList, Node).
func dssXMLUtilsIsNodeListContains(nodeList []*xmldom.Node, node *xmldom.Node) bool {
	for _, n := range nodeList {
		if n == node {
			return true
		}
	}
	return false
}

// DSSXMLUtilsAlignChildrenIndents aligns indents for all children of the given node. Ports
// alignChildrenIndents(Node). The error return is always nil here: Java's method cannot fail
// either, but every other pretty-print helper in this file needs to propagate an XPath/DOM
// failure, and the caller-side pattern (extension_builder.go) expects a uniform (Node, error)
// shape across all of them.
func DSSXMLUtilsAlignChildrenIndents(parentNode *xmldom.Node) (*xmldom.Node, error) {
	if parentNode.FirstChild != nil {
		nodeChildren := dssXMLUtilsChildNodes(parentNode)
		targetIndent, ok := dssXMLUtilsGetTargetIndent(nodeChildren)
		if ok {
			for i := 0; i < len(nodeChildren)-1; i++ {
				node := nodeChildren[i]
				if node.Kind == xmldom.Text {
					dssXMLUtilsReplaceTextValue(node, targetIndent)
				}
			}
			lastChild := parentNode.LastChild
			targetIndent = targetIndent[:len(targetIndent)-DSSXMLUtilsTransformerIndentNumber]
			switch lastChild.Kind {
			case xmldom.Element:
				xmlutils.DomUtilsSetTextNode(parentNode.OwnerDocument(), parentNode, targetIndent)
			case xmldom.Text:
				dssXMLUtilsReplaceTextValue(lastChild, targetIndent)
			}
		}
	}
	return parentNode, nil
}

// dssXMLUtilsReplaceTextValue replaces a text node's value in place by swapping in a freshly
// constructed text node, following pretty_print_transformer.go's established idiom of always
// constructing new *xmldom.Node text nodes rather than mutating an existing node's Value field
// directly.
func dssXMLUtilsReplaceTextValue(node *xmldom.Node, value string) {
	node.Parent.ReplaceChild(xmldom.NewText(value), node)
}

// dssXMLUtilsGetTargetIndent ports the private static getTargetIndent(NodeList).
func dssXMLUtilsGetTargetIndent(nodeChildren []*xmldom.Node) (string, bool) {
	for i := 0; i < len(nodeChildren)-1; i++ {
		node := nodeChildren[i]
		if node.Kind == xmldom.Text {
			return node.Value, true
		}
	}
	return "", false
}

// DSSXMLUtilsGetIDIdentifier finds the value of an attribute named ID (case-insensitive), or ""
// if none. If there is more than one ID attribute then the first one is returned. Ports
// getIDIdentifier(Node).
func DSSXMLUtilsGetIDIdentifier(node *xmldom.Node) string {
	return DSSXMLUtilsGetAttribute(node, common.XMLDSigAttribute_ID.AttributeName())
}

// DSSXMLUtilsGetAttribute returns the attribute value for the given attribute name if it
// exists, otherwise "". Ports getAttribute(Node, String).
func DSSXMLUtilsGetAttribute(node *xmldom.Node, attributeName string) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attrs {
		localName := attr.Name.Local
		if utils.AreStringsEqualIgnoreCase(attributeName, localName) {
			return attr.TextContent()
		}
	}
	return ""
}

// DSSXMLUtilsValidateAgainstXSD validates an XML against the XAdES XSD schema. Ports
// validateAgainstXSD(XSDAbstractUtils, Source). See the file header's "FLAGGED FOR INTEGRATOR"
// note: no bundled-schema/javax.xml.validation-equivalent Go package exists yet in this port, so
// this always reports no errors (i.e. the structure is always considered valid).
func DSSXMLUtilsValidateAgainstXSD(xsdUtils XAdESStructureXSDUtils, source *xmldom.Node) []string {
	return nil
}

// XAdES111XSDUtils, XAdES122XSDUtils and XAdES319132XSDUtils stand in for
// eu.europa.esig.xades.XAdES111Utils#getInstance(), XAdES122Utils#getInstance() and
// XAdES319132Utils#getInstance(). Per xades_structure_validator.go's file header BLOCKED
// FORWARD DEPENDENCY note, no bundled XAdES XSD schema resources or javax.xml.validation
// equivalent exists anywhere in this port yet, so these return an opaque placeholder value;
// DSSXMLUtilsValidateAgainstXSD above never inspects it (it always reports no errors), so the
// placeholder is never dereferenced. XAdESStructureValidator.getUtils's namespace dispatch -
// including its UnsupportedOperationException panic for unrecognized namespaces - is still
// exercised faithfully; only the actual schema validation is a no-op.
var (
	xades111XSDUtils    XAdESStructureXSDUtils = struct{}{}
	xades122XSDUtils    XAdESStructureXSDUtils = struct{}{}
	xades319132XSDUtils XAdESStructureXSDUtils = struct{}{}
)

// XAdES111XSDUtils ports eu.europa.esig.xades.XAdES111Utils#getInstance(). See the note above.
func XAdES111XSDUtils() XAdESStructureXSDUtils {
	return xades111XSDUtils
}

// XAdES122XSDUtils ports eu.europa.esig.xades.XAdES122Utils#getInstance(). See the note above.
func XAdES122XSDUtils() XAdESStructureXSDUtils {
	return xades122XSDUtils
}

// XAdES319132XSDUtils ports eu.europa.esig.xades.XAdES319132Utils#getInstance(). See the note
// above.
func XAdES319132XSDUtils() XAdESStructureXSDUtils {
	return xades319132XSDUtils
}

// DSSXMLUtilsIsDuplicateIdsDetected detects duplicate id values. Ports
// isDuplicateIdsDetected(DSSDocument).
func DSSXMLUtilsIsDuplicateIdsDetected(doc model.DSSDocument) bool {
	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(doc)
	if err != nil {
		return false
	}
	dom.RegisterIDs()
	duplicates := dom.DuplicateIDs()
	if len(duplicates) > 0 {
		// Upstream logs "Duplicated identifier '{}' detected".
		return true
	}
	return false
}

// DSSXMLUtilsGetReferenceOriginalContentBytes returns bytes of the original referenced data.
// Ports getReferenceOriginalContentBytes(Reference).
func DSSXMLUtilsGetReferenceOriginalContentBytes(reference *xmldsig.Reference) []byte {
	// returns bytes after transformation in case of enveloped signature
	transformsElement := reference.TransformsElement()
	if transformsElement != nil {
		for transformation := transformsElement.FirstChild; transformation != nil; transformation = transformation.NextSibling {
			if dssXMLUtilsIsEnvelopedTransform(transformation) {
				referencedBytes, err := reference.ReferencedBytes()
				if err == nil {
					return referencedBytes
				}
				// if an exception occurs during the transformations, fall through to the
				// pre-transformation bytes below, mirroring upstream's catch(XMLSecurityException).
				break
			}
			// if enveloped transformations are not applied to the signature go further and
			// return bytes before transformation
		}
	}
	// otherwise bytes before transformation
	return dssXMLUtilsGetBytesBeforeTransformation(reference)
}

// dssXMLUtilsIsEnvelopedTransform ports the private static isEnvelopedTransform(Node).
func dssXMLUtilsIsEnvelopedTransform(transformation *xmldom.Node) bool {
	if transformation.Kind != xmldom.Element {
		return false
	}
	algorithm := transformation.AttrValue("", common.XMLDSigAttribute_ALGORITHM.AttributeName())
	if string(xmldsig.TransformEnvelopedSignature) == algorithm {
		return true
	}
	if string(xmldsig.TransformXPath) == algorithm || string(xmldsig.TransformXPath2Filter) == algorithm {
		for item := transformation.FirstChild; item != nil; item = item.NextSibling {
			if item.Kind == xmldom.Element && dssXMLUtilsTransformationXPathNodeName == item.Name.Local &&
				dssXMLUtilsTransformationExcludeSignature == item.TextContent() {
				return true
			}
		}
	}
	return false
}

// dssXMLUtilsGetBytesBeforeTransformation ports the private static
// getBytesBeforeTransformation(Reference).
func dssXMLUtilsGetBytesBeforeTransformation(reference *xmldsig.Reference) []byte {
	data, err := reference.ContentsBeforeTransformation()
	if err != nil || data == nil {
		// Upstream logs "Original data is not provided..." / "Unable to retrieve the content...".
		return nil
	}
	b, err := data.Bytes()
	if err != nil {
		return nil
	}
	return b
}

// DSSXMLUtilsGetDigestAndValue extracts the Digest algorithm and value from an element of type
// DigestAlgAndValueType. Ports getDigestAndValue(Element). Returns a nil-Digest (IsEmpty()) on a
// malformed/absent element.
func DSSXMLUtilsGetDigestAndValue(element *xmldom.Node) model.Digest {
	if element == nil {
		return model.Digest{}
	}

	var digestAlgorithmUri, digestValueBase64 string
	if definition.XAdESNamespace_XADES_111.IsSameUri(element.Name.Space) {
		digestAlgorithmUri, _ = xmlutils.XPathUtilsGetValue(element,
			common.FromCurrentPositionAttribute(definition.XAdES111Element_DIGEST_METHOD, common.XMLDSigAttribute_ALGORITHM))
		digestValueBase64, _ = xmlutils.XPathUtilsGetValue(element, common.FromCurrentPosition(definition.XAdES111Element_DIGEST_VALUE))
	} else {
		digestAlgorithmUri, _ = xmlutils.XPathUtilsGetValue(element, common.XMLDSigPath_DIGEST_METHOD_ALGORITHM_PATH)
		digestValueBase64, _ = xmlutils.XPathUtilsGetValue(element, common.XMLDSigPath_DIGEST_VALUE_PATH)
	}

	digestAlgorithm := dssXMLUtilsGetDigestAlgorithm(digestAlgorithmUri)
	digestValue := dssXMLUtilsGetDigestValue(digestValueBase64)

	if digestAlgorithm == "" || utils.IsArrayEmpty(digestValue) {
		// Upstream logs "Unable to read object DigestAlgAndValueType. An error occurred during
		// processing.".
		return model.Digest{}
	}
	return model.NewDigest(digestAlgorithm, digestValue)
}

// dssXMLUtilsGetDigestValue ports the private static getDigestValue(String).
func dssXMLUtilsGetDigestValue(digestValueBase64 string) []byte {
	if utils.IsStringEmpty(digestValueBase64) {
		// Upstream logs "An empty DigestValue obtained!".
		return nil
	}
	if !utils.IsBase64Encoded(digestValueBase64) {
		// Upstream logs "The DigestValue is not base64 encoded! Obtained string : {}".
		return nil
	}
	return utils.FromBase64(digestValueBase64)
}

// dssXMLUtilsGetDigestAlgorithm ports the private static getDigestAlgorithm(String).
func dssXMLUtilsGetDigestAlgorithm(digestAlgorithmUri string) enumerations.DigestAlgorithm {
	if utils.IsStringNotEmpty(digestAlgorithmUri) {
		result, err := enumerations.DigestAlgorithmForXML(digestAlgorithmUri)
		if err == nil {
			return result
		}
		// Upstream logs "Unable to retrieve the used digest algorithm".
	}
	return ""
}

// DSSXMLUtilsContainsTransforms checks if the reference element contains any transformation.
// Ports containsTransforms(Element).
func DSSXMLUtilsContainsTransforms(referenceElement *xmldom.Node) bool {
	transforms, err := xmlutils.XPathUtilsGetElement(referenceElement, common.XMLDSigPath_TRANSFORMS_PATH)
	if err != nil {
		// Upstream logs "Unable to detect Transforms".
		return false
	}
	return transforms != nil
}

// DSSXMLUtilsIsSignedProperties determines if the given reference refers to the SignedProperties
// element. Ports isSignedProperties(Reference, XAdESPath).
func DSSXMLUtilsIsSignedProperties(reference *xmldsig.Reference, xadesPath definition.XAdESPath) bool {
	return xadesPath.SignedPropertiesUri() == reference.Type()
}

// DSSXMLUtilsIsCounterSignatureReference determines if the given reference refers to a
// CounterSignature element within the signature. Ports isCounterSignatureReference(Reference,
// XAdESSignature).
func DSSXMLUtilsIsCounterSignatureReference(reference *xmldsig.Reference, sig *XAdESSignature) bool {
	masterSignature, _ := sig.MasterSignature().(*XAdESSignature)
	if masterSignature != nil {
		return DSSXMLUtilsIsCounterSignatureReferenceType(reference.Type()) ||
			dssXMLUtilsIsSignatureValueReferenced(masterSignature, reference)
	} else if DSSXMLUtilsIsCounterSignatureReferenceType(reference.Type()) { //nolint:staticcheck // mirrors upstream DSSXMLUtils#isCounterSignatureReference: the branch is kept because Java's body is `LOG.warn("Master signature is not found! ...")` only.
		// Upstream logs "Master signature is not found! Unable to verify counter signed
		// SignatureValue for detached signatures.".
	}
	return false
}

// dssXMLUtilsIsSignatureValueReferenced ports the private static
// isSignatureValueReferenced(XAdESSignature, Reference).
func dssXMLUtilsIsSignatureValueReferenced(masterSignature *XAdESSignature, reference *xmldsig.Reference) bool {
	sigValueID := masterSignature.SignatureValueId()
	return sigValueID != "" && sigValueID == xmlutils.DomUtilsGetId(reference.URI())
}

// DSSXMLUtilsIsKeyInfoReference checks if the given reference is linked to a KeyInfo element.
// Ports isKeyInfoReference(Reference, Element).
func DSSXMLUtilsIsKeyInfoReference(reference *xmldsig.Reference, signatureElement *xmldom.Node) bool {
	id := xmlutils.DomUtilsGetId(reference.URI())
	keyInfoElement := xmlutils.XPathUtilsGetElementByIdWithQuery(signatureElement, common.XMLDSigPath_KEY_INFO_PATH, id)
	return keyInfoElement != nil
}

// DSSXMLUtilsIsSignaturePropertiesReference checks if the given reference is linked to a
// SignatureProperties element or one of its SignatureProperty children. Ports
// isSignaturePropertiesReference(Reference, Element).
func DSSXMLUtilsIsSignaturePropertiesReference(reference *xmldsig.Reference, signatureElement *xmldom.Node) bool {
	id := xmlutils.DomUtilsGetId(reference.URI())
	signaturePropertiesElement := xmlutils.XPathUtilsGetElementByIdWithQuery(signatureElement, common.XMLDSigPath_SIGNATURE_PROPERTIES_PATH, id)
	signaturePropertyElement := xmlutils.XPathUtilsGetElementByIdWithQuery(signatureElement, common.XMLDSigPath_SIGNATURE_PROPERTY_PATH, id)
	return signaturePropertiesElement != nil || signaturePropertyElement != nil
}

// DSSXMLUtilsIsObjectReferenceType checks if the given referenceType is an xmldsig Object type.
// Ports isObjectReferenceType(String).
func DSSXMLUtilsIsObjectReferenceType(referenceType string) bool {
	return common.XMLDSigPath_OBJECT_TYPE == referenceType
}

// DSSXMLUtilsIsManifestReferenceType checks if the given referenceType is an xmldsig Manifest
// type. Ports isManifestReferenceType(String).
func DSSXMLUtilsIsManifestReferenceType(referenceType string) bool {
	return common.XMLDSigPath_MANIFEST_TYPE == referenceType
}

// DSSXMLUtilsIsCounterSignatureReferenceType checks if the given referenceType is an etsi
// Countersignature type. Ports isCounterSignatureReferenceType(String).
func DSSXMLUtilsIsCounterSignatureReferenceType(referenceType string) bool {
	return common.XMLDSigPath_COUNTER_SIGNATURE_TYPE == referenceType
}

// DSSXMLUtilsIsSameDocumentReference determines whether referenceUri points to a
// 'same-document' reference: a URI-Reference consisting of a hash sign ('#') followed by a
// fragment, or an empty URI. Ports isSameDocumentReference(String).
func DSSXMLUtilsIsSameDocumentReference(referenceUri string) bool {
	return referenceUri == "" || xmlutils.DomUtilsStartsFromHash(referenceUri)
}

// DSSXMLUtilsGetObjectById gets ds:Object by its Id from the ds:Signature element. Ports
// getObjectById(Element, String).
func DSSXMLUtilsGetObjectById(signatureElement *xmldom.Node, id string) *xmldom.Node {
	if !utils.IsStringNotBlank(id) {
		return nil
	}
	// Upstream logs "An error occurred on attempt to extract Object element with Id '{}' : {}" on
	// failure; XPathUtilsGetElementByIdWithQuery already reproduces that catch-and-return-nil.
	return xmlutils.XPathUtilsGetElementByIdWithQuery(signatureElement, common.XMLDSigPath_OBJECT_PATH, id)
}

// DSSXMLUtilsGetManifestById gets ds:Manifest by its Id from the ds:Signature element. Ports
// getManifestById(Element, String).
func DSSXMLUtilsGetManifestById(signatureElement *xmldom.Node, id string) *xmldom.Node {
	if !utils.IsStringNotBlank(id) {
		return nil
	}
	return xmlutils.XPathUtilsGetElementByIdWithQuery(signatureElement, common.XMLDSigPath_MANIFEST_PATH, id)
}

// DSSXMLUtilsInitManifest initializes a Manifest object from the provided ds:Manifest element.
// Ports initManifest(Element).
func DSSXMLUtilsInitManifest(manifestElement *xmldom.Node) (*xmldsig.Manifest, error) {
	return xmldsig.NewManifest(manifestElement, nil)
}

// DSSXMLUtilsInitManifestWithDetachedContent initializes a Manifest object from the provided
// ds:Manifest element with the provided detachedContents. Ports
// initManifestWithDetachedContent(Element, List<DSSDocument>).
func DSSXMLUtilsInitManifestWithDetachedContent(manifestElement *xmldom.Node, detachedContents []model.DSSDocument) (*xmldsig.Manifest, error) {
	manifest, err := DSSXMLUtilsInitManifest(manifestElement)
	if err != nil {
		return nil, err
	}
	DSSXMLUtilsInitManifestDetachedContent(manifest, detachedContents)
	return manifest, nil
}

// DSSXMLUtilsInitManifestDetachedContent initializes detached content within the given manifest,
// registering one DetachedSignatureResolver per distinct digest algorithm found among its
// ds:Reference elements. Ports initManifestDetachedContent(Manifest, List<DSSDocument>).
func DSSXMLUtilsInitManifestDetachedContent(manifest *xmldsig.Manifest, detachedContents []model.DSSDocument) {
	if utils.IsCollectionNotEmpty(detachedContents) {
		for _, digestAlgorithm := range DSSXMLUtilsGetReferenceDigestAlgos(manifest.Element()) {
			manifest.AddResourceResolver(&xmldsig.DetachedSignatureResolver{
				Documents:       detachedContents,
				DigestAlgorithm: digestAlgorithm,
			})
		}
	}
}

// DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey extracts the signing certificate's public key
// from the KeyInfo element of a given signature if present. NOTE: can return nil (the value is
// optional). See the file header's "FLAGGED FOR INTEGRATOR" note on scope. Ports
// getKeyInfoSigningCertificatePublicKey(Element).
func DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey(signatureElement *xmldom.Node) *model.PublicKey {
	keyInfoElement, err := xmlutils.XPathUtilsGetElement(signatureElement, common.XMLDSigPath_KEY_INFO_PATH)
	if err != nil || keyInfoElement == nil {
		// Upstream logs "Unable to extract the public key. Reason : KeyInfo element is null".
		return nil
	}
	if publicKey := dssXMLUtilsKeyInfoPublicKeyFromCertificate(keyInfoElement); publicKey != nil {
		return publicKey
	}
	return dssXMLUtilsKeyInfoPublicKeyFromRSAKeyValue(keyInfoElement)
}

// dssXMLUtilsKeyInfoPublicKeyFromCertificate extracts the public key from a
// ds:KeyInfo/ds:X509Data/ds:X509Certificate child, when present.
func dssXMLUtilsKeyInfoPublicKeyFromCertificate(keyInfoElement *xmldom.Node) *model.PublicKey {
	certElement, err := xmlutils.XPathUtilsGetElement(keyInfoElement,
		common.FromCurrentPosition(common.XMLDSigElement_X509_DATA, common.XMLDSigElement_X509_CERTIFICATE))
	if err != nil || certElement == nil {
		return nil
	}
	der := utils.FromBase64(utils.Trim(certElement.TextContent()))
	cert, err := spi.DSSUtilsLoadCertificateFromBinary(der)
	if err != nil {
		// Upstream logs "Unable to extract signing certificate's public key. Reason : {}".
		return nil
	}
	return cert.PublicKey()
}

// dssXMLUtilsKeyInfoPublicKeyFromRSAKeyValue extracts the public key from a
// ds:KeyInfo/ds:KeyValue/ds:RSAKeyValue child, when present.
func dssXMLUtilsKeyInfoPublicKeyFromRSAKeyValue(keyInfoElement *xmldom.Node) *model.PublicKey {
	modulusElement, err := xmlutils.XPathUtilsGetElement(keyInfoElement,
		common.FromCurrentPosition(common.XMLDSigElement_KEY_VALUE, common.XMLDSigElement_RSA_KEY_VALUE, common.XMLDSigElement_MODULUS))
	if err != nil || modulusElement == nil {
		return nil
	}
	exponentElement, err := xmlutils.XPathUtilsGetElement(keyInfoElement,
		common.FromCurrentPosition(common.XMLDSigElement_KEY_VALUE, common.XMLDSigElement_RSA_KEY_VALUE, common.XMLDSigElement_EXPONENT))
	if err != nil || exponentElement == nil {
		return nil
	}
	modulusBytes := utils.FromBase64(utils.Trim(modulusElement.TextContent()))
	exponentBytes := utils.FromBase64(utils.Trim(exponentElement.TextContent()))
	if utils.IsArrayEmpty(modulusBytes) || utils.IsArrayEmpty(exponentBytes) {
		return nil
	}
	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: int(new(big.Int).SetBytes(exponentBytes).Int64()),
	}
	spki, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}
	return model.NewPublicKeyFromEncoded(spki, pub)
}

// DSSXMLUtilsCreateCounterSignature creates and returns a counter signature found in
// counterSignatureElement. Returns nil on failure. Ports createCounterSignature(Element,
// XAdESSignature).
func DSSXMLUtilsCreateCounterSignature(counterSignatureElement *xmldom.Node, masterSignature *XAdESSignature) (result validation.AdvancedSignature) {
	defer func() {
		if recover() != nil {
			// Upstream logs "An error occurred during counter signature extraction. The element
			// entry is skipped. Reason : {}".
			result = nil
		}
	}()
	/*
	 * 5.2.7.2 Enveloped countersignatures: the CounterSignature qualifying property
	 *
	 * The CounterSignature qualifying property shall contain one countersignature
	 * of the XAdES signature where CounterSignature is incorporated.
	 */
	counterSignatureNode, err := xmlutils.XPathUtilsGetNode(counterSignatureElement, common.XMLDSigPath_SIGNATURE_PATH)
	if err != nil || counterSignatureNode == nil {
		return nil
	}
	counterSigDOMElement := NewXAdESDOMElement(counterSignatureNode, masterSignature.OwnerDocument())

	// Verify that the element is a proper signature by trying to build a XAdESSignature out of it
	xadesCounterSignature := NewXAdESSignatureFromDOMElement(counterSigDOMElement)
	xadesCounterSignature.SetFilename(masterSignature.Filename())
	xadesCounterSignature.SetDetachedContents(masterSignature.DetachedContents())
	xadesCounterSignature.SetMasterSignature(masterSignature)
	return xadesCounterSignature
}

// DSSXMLUtilsGetAllSignaturesExceptCounterSignatures returns a NodeList of all "ds:Signature"
// elements found in documentNode. Ports getAllSignaturesExceptCounterSignatures(Node).
func DSSXMLUtilsGetAllSignaturesExceptCounterSignatures(documentNode *xmldom.Node) ([]*xmldom.Node, error) {
	return xmlutils.XPathUtilsGetNodeList(documentNode,
		common.AllNotParent(common.XMLDSigElement_SIGNATURE, definition.XAdES132Element_COUNTER_SIGNATURE))
}

// DSSXMLUtilsGetAllEncapsulatedTimestamps returns a NodeList of all
// "xades:EncapsulatedTimeStamp" elements found in documentNode. Ports
// getAllEncapsulatedTimestamps(Node).
func DSSXMLUtilsGetAllEncapsulatedTimestamps(documentNode *xmldom.Node) ([]*xmldom.Node, error) {
	return xmlutils.XPathUtilsGetNodeList(documentNode, common.All(definition.XAdES132Element_ENCAPSULATED_TIMESTAMP))
}

// DSSXMLUtilsGetReferenceNodeList returns a NodeList of "ds:Reference" elements. Ports
// getReferenceNodeList(Node).
func DSSXMLUtilsGetReferenceNodeList(signatureElement *xmldom.Node) ([]*xmldom.Node, error) {
	return xmlutils.XPathUtilsGetNodeList(signatureElement, common.XMLDSigPath_SIGNED_INFO_REFERENCE_PATH)
}

// DSSXMLUtilsGetReferenceOutputType returns the expected dereferencing output for the provided
// reference. Generic over *DSSReference and *xmldsig.Reference, unifying Java's two overloads
// (getReferenceOutputType(DSSReference) and getReferenceOutputType(Reference) throws
// XMLSecurityException) into one Go name, since both landed call sites
// (all_data_objects_time_stamp_builder.go, xades_timestamp_message_digest_builder.go) already
// invoke DSSXMLUtilsGetReferenceOutputType with either concrete type. Internal/xmldsig's
// Reference.TransformsElement() cannot fail (per its own doc comment), so this port drops the
// checked exception entirely rather than returning an error nobody can produce.
func DSSXMLUtilsGetReferenceOutputType[T *DSSReference | *xmldsig.Reference](reference T) ReferenceOutputType {
	switch r := any(reference).(type) {
	case *DSSReference:
		outputType := dssXMLUtilsGetDereferenceOutputType(r.Uri())
		if utils.IsCollectionNotEmpty(r.Transforms()) {
			for _, transform := range r.Transforms() {
				outputType = dssXMLUtilsGetTransformOutputType(transform.Algorithm())
			}
		}
		return outputType
	case *xmldsig.Reference:
		outputType := dssXMLUtilsGetDereferenceOutputType(r.URI())
		transformsElement := r.TransformsElement()
		if transformsElement != nil {
			for transform := transformsElement.FirstChild; transform != nil; transform = transform.NextSibling {
				if transform.Kind != xmldom.Element {
					continue
				}
				outputType = dssXMLUtilsGetTransformOutputType(transform.AttrValue("", common.XMLDSigAttribute_ALGORITHM.AttributeName()))
			}
		}
		return outputType
	}
	return ReferenceOutputType_OCTET_STREAM
}

// dssXMLUtilsGetDereferenceOutputType ports the private static getDereferenceOutputType(String).
func dssXMLUtilsGetDereferenceOutputType(referenceUri string) ReferenceOutputType {
	if DSSXMLUtilsIsSameDocumentReference(referenceUri) {
		return ReferenceOutputType_NODE_SET
	}
	return ReferenceOutputType_OCTET_STREAM
}

// dssXMLUtilsGetTransformOutputType ports the private static getTransformOutputType(String).
func dssXMLUtilsGetTransformOutputType(algorithmUri string) ReferenceOutputType {
	dssXMLUtilsMu.Lock()
	_, ok := dssXMLUtilsTransformsWithNodeSetOutput[algorithmUri]
	dssXMLUtilsMu.Unlock()
	if ok {
		return ReferenceOutputType_NODE_SET
	}
	return ReferenceOutputType_OCTET_STREAM
}

// DSSXMLUtilsApplyTransforms applies transforms on the node and returns the byte array to be
// used for a digest computation. NOTE: returns the original node binaries, if the list of
// transforms is empty. Ports applyTransforms(Node, List<DSSTransform>).
func DSSXMLUtilsApplyTransforms(node *xmldom.Node, transforms []DSSTransform) ([]byte, error) {
	if utils.IsCollectionNotEmpty(transforms) {
		output := NewDSSTransformOutputFromNode(node)
		var bytesValue []byte
		for i, transform := range transforms {
			var err error
			output, err = transform.PerformTransform(output)
			if err != nil {
				return nil, err
			}
			bytesValue, err = output.Bytes()
			if err != nil {
				return nil, err
			}
			if i < len(transforms)-1 && utils.IsArrayEmpty(bytesValue) {
				return nil, exception.NewIllegalInputException(fmt.Sprintf(
					"Unable to perform the next transform. The %v produced an empty output!", transform))
			}
		}
		// Upstream logs "Reference bytes after transforms: " + new String(bytes) when debug
		// logging is enabled, and warns "The output of reference transforms processing is an
		// empty byte array!" when it is empty.
		return bytesValue, nil
	}
	return xmlutils.DomUtilsGetNodeBytes(node)
}

// DSSXMLUtilsApplyTransformsToDocument applies transforms on document content and returns the
// byte array to be used for a digest computation. NOTE: returns the original document binaries,
// if the list of transforms is empty. document shall represent an XML content. Ports
// applyTransforms(DSSDocument, List<DSSTransform>).
func DSSXMLUtilsApplyTransformsToDocument(document model.DSSDocument, transforms []DSSTransform) ([]byte, error) {
	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(document)
	if err != nil {
		return nil, err
	}
	return DSSXMLUtilsApplyTransforms(dom, transforms)
}

// DSSXMLUtilsGetReferenceDigestAlgos returns the list of DigestAlgorithms for all references
// contained inside referenceContainer, deduplicated and in first-seen order (the determinism
// sweep's slice-of-first-seen substitute for Java's iteration-order-unstable HashSet). Ports
// getReferenceDigestAlgos(Element).
func DSSXMLUtilsGetReferenceDigestAlgos(referenceContainer *xmldom.Node) []enumerations.DigestAlgorithm {
	var digestAlgorithms []enumerations.DigestAlgorithm
	seen := make(map[enumerations.DigestAlgorithm]struct{})
	referenceNodeList, err := xmlutils.XPathUtilsGetNodeList(referenceContainer, common.XMLDSigPath_REFERENCE_PATH)
	if err != nil {
		return digestAlgorithms
	}
	for _, referenceElement := range referenceNodeList {
		digest := DSSXMLUtilsGetDigestAndValue(referenceElement)
		if !digest.IsEmpty() {
			if _, ok := seen[digest.Algorithm()]; !ok {
				seen[digest.Algorithm()] = struct{}{}
				digestAlgorithms = append(digestAlgorithms, digest.Algorithm())
			}
		}
	}
	return digestAlgorithms
}

// DSSXMLUtilsGetReferenceTypes returns a list of reference types. Ports getReferenceTypes(Element).
func DSSXMLUtilsGetReferenceTypes(referenceContainer *xmldom.Node) []string {
	var referenceTypes []string
	referenceNodeList, err := xmlutils.XPathUtilsGetNodeList(referenceContainer, common.XMLDSigPath_REFERENCE_PATH)
	if err != nil {
		return referenceTypes
	}
	for _, referenceElement := range referenceNodeList {
		typ := referenceElement.AttrValue("", common.XMLDSigAttribute_TYPE.AttributeName())
		if utils.IsStringNotEmpty(typ) {
			referenceTypes = append(referenceTypes, typ)
		}
	}
	return referenceTypes
}

// DSSXMLUtilsExtractReferences extracts a list of References from the given Manifest object.
// NOTE: can be used also for a SignedInfo element. Ports extractReferences(Manifest).
func DSSXMLUtilsExtractReferences(manifest *xmldsig.Manifest) []*xmldsig.Reference {
	references, err := manifest.References()
	if err != nil {
		// Upstream logs "Unable to retrieve reference #{} : {}" per failing index; this port's
		// Manifest.References() reports the batch as a whole rather than per-index, so the
		// failure is reported for the whole manifest instead.
		return []*xmldsig.Reference{}
	}
	return references
}

// DSSXMLUtilsGetReferenceDigest returns the Digest extracted from the provided reference. Ports
// getReferenceDigest(Reference).
func DSSXMLUtilsGetReferenceDigest(reference *xmldsig.Reference) model.Digest {
	digestValue, err := reference.DigestValue()
	if err != nil {
		// Upstream logs "Unable to extract Digest from a reference with Id [{}] : {}".
		return model.Digest{}
	}
	digestAlgorithm, err := reference.DigestAlgorithm()
	if err != nil {
		return model.Digest{}
	}
	return model.NewDigest(digestAlgorithm, digestValue)
}

// DSSXMLUtilsGetReferenceId retrieves the Id attribute value of the given reference, when
// applicable. NOTE: used because Apache Santuario Signature returns an empty string instead of a
// null result. Ports getReferenceId(Reference).
func DSSXMLUtilsGetReferenceId(reference *xmldsig.Reference) string {
	if reference == nil {
		return ""
	}
	element := reference.Element()
	if element == nil {
		return ""
	}
	return DSSXMLUtilsGetAttribute(element, common.XMLDSigAttribute_ID.AttributeName())
}

// DSSXMLUtilsGetReferenceURI retrieves the URI attribute value of the given reference, when
// applicable. NOTE: used because Apache Santuario Signature returns an empty string instead of a
// null result. Ports getReferenceURI(Reference).
func DSSXMLUtilsGetReferenceURI(reference *xmldsig.Reference) string {
	if reference == nil {
		return ""
	}
	element := reference.Element()
	if element == nil {
		return ""
	}
	referenceUri := DSSXMLUtilsGetAttribute(element, common.XMLDSigAttribute_URI.AttributeName())
	if referenceUri == "" {
		return ""
	}
	return spi.DSSUtilsDecodeURI(referenceUri)
}

// DSSXMLUtilsIsAbleToDeReferenceContent checks if the original reference document content can be
// obtained (de-referenced). Ports isAbleToDeReferenceContent(Reference).
//
// DEVIATION: upstream's private getClosedContentsBeforeTransformation additionally closes the
// XMLSignatureInput's octet stream (a SANTUARIO-622 workaround for a leaked InputStream); this
// port's xmldsig.Reference.ContentsBeforeTransformation buffers eagerly and has no such stream to
// close.
func DSSXMLUtilsIsAbleToDeReferenceContent(reference *xmldsig.Reference) bool {
	data, err := reference.ContentsBeforeTransformation()
	return err == nil && data != nil
}

// DSSXMLUtilsIsReferencedContentAmbiguous checks if the reference with the given uri occurs
// multiple times in document (a wrapping attack indicator). Ports
// isReferencedContentAmbiguous(Document, String).
func DSSXMLUtilsIsReferencedContentAmbiguous(document *xmldom.Node, uri string) bool {
	if utils.IsStringEmpty(uri) {
		// empty URI means enveloped signature (unambiguous)
		return false
	}
	docNode := document
	if docNode != nil && docNode.Kind != xmldom.Document {
		docNode = docNode.OwnerDocument()
	}
	if docNode == nil {
		return false
	}
	id := xmlutils.DomUtilsGetId(uri)
	docNode.RegisterIDs()
	for _, dup := range docNode.DuplicateIDs() {
		if dup == id {
			return true
		}
	}
	return false
}

// DSSXMLUtilsIncorporateTransforms incorporates a ds:Transforms element into the given parent
// element. Ports incorporateTransforms(Element, List<DSSTransform>, DSSNamespace).
func DSSXMLUtilsIncorporateTransforms(parentElement *xmldom.Node, transforms []DSSTransform, namespace *common.DSSNamespace) {
	if utils.IsCollectionNotEmpty(transforms) {
		documentDom := parentElement.OwnerDocument()
		transformsDom := xmlutils.DomUtilsCreateElementNS(documentDom, namespace, common.XMLDSigElement_TRANSFORMS)
		parentElement.AppendChild(transformsDom)
		for _, dssTransform := range transforms {
			dssTransform.CreateTransform(documentDom, transformsDom)
		}
	}
}

// DSSXMLUtilsIncorporateDigestMethod creates the ds:DigestMethod DOM object:
//
//	<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//
// Ports incorporateDigestMethod(Element, DigestAlgorithm, DSSNamespace).
func DSSXMLUtilsIncorporateDigestMethod(parentElement *xmldom.Node, digestAlgorithm enumerations.DigestAlgorithm, namespace *common.DSSNamespace) {
	documentDom := parentElement.OwnerDocument()
	digestMethodDom := xmlutils.DomUtilsAddElement(documentDom, parentElement, namespace, common.XMLDSigElement_DIGEST_METHOD)
	digestMethodDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ALGORITHM.AttributeName()}, digestAlgorithm.URI())
}

// DSSXMLUtilsIncorporateDigestValue creates the ds:DigestValue DOM object:
//
//	<ds:DigestValue>fj8SJujSXU4fi342bdtiKVbglA0=</ds:DigestValue>
//
// Ports incorporateDigestValue(Element, String, DSSNamespace).
func DSSXMLUtilsIncorporateDigestValue(parentDom *xmldom.Node, base64EncodedDigestBytes string, namespace *common.DSSNamespace) {
	documentDom := parentDom.OwnerDocument()
	digestValueDom := xmlutils.DomUtilsCreateElementNS(documentDom, namespace, common.XMLDSigElement_DIGEST_VALUE)
	digestValueDom.AppendChild(xmldom.NewText(base64EncodedDigestBytes))
	parentDom.AppendChild(digestValueDom)
}

// DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault returns params.ReferenceDigestAlgorithm() if
// set, params.DigestAlgorithm() otherwise. Panics with the Java message when neither algorithm is
// usable (IllegalArgumentException upstream). Ports getReferenceDigestAlgorithmOrDefault(
// XAdESSignatureParameters).
func DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault(params *XAdESSignatureParameters) enumerations.DigestAlgorithm {
	digestAlgorithm := params.ReferenceDigestAlgorithm()
	if digestAlgorithm == "" {
		digestAlgorithm = params.DigestAlgorithm()
	}
	if digestAlgorithm == "" || digestAlgorithm.URI() == "" {
		panic(fmt.Sprintf("The Reference DigestAlgorithm '%v' is not supported for XAdES creation! "+
			"Define another algorithm within #setReferenceDigestAlgorithm method.", digestAlgorithm))
	}
	return digestAlgorithm
}

// DSSXMLUtilsEnsureNamespacesDefined produces a copy of the document and returns an element by
// the defined xpathString. This method can be used as a workaround for canonicalization, as
// namespaces are not added to the canonicalizer for newly created elements
// (SANTUARIO-139). Deprecated since DSS 6.5; use DSSXMLUtilsEnsureNamespacesDefinedWithQuery
// instead. Ports the deprecated ensureNamespacesDefined(Document, String, String).
func DSSXMLUtilsEnsureNamespacesDefined(document *xmldom.Node, elementId, xpathString string) (*xmldom.Node, error) {
	element, err := dssXMLUtilsRecreateAndLocate(document, elementId)
	if err != nil {
		return nil, err
	}
	return xmlutils.DomUtilsGetElement(element, xpathString)
}

// DSSXMLUtilsEnsureNamespacesDefinedWithQuery is the XPathQuery-based overload of
// DSSXMLUtilsEnsureNamespacesDefined. Ports ensureNamespacesDefined(Document, String, XPathQuery).
func DSSXMLUtilsEnsureNamespacesDefinedWithQuery(document *xmldom.Node, elementId string, xpathQuery common.XPathQuery) (*xmldom.Node, error) {
	element, err := dssXMLUtilsRecreateAndLocate(document, elementId)
	if err != nil {
		return nil, err
	}
	return xmlutils.XPathUtilsGetElement(element, xpathQuery)
}

// dssXMLUtilsRecreateAndLocate serializes and rebuilds document, then locates its document
// element or, when elementId is set, the element carrying that Id. Shared by both
// DSSXMLUtilsEnsureNamespacesDefined overloads.
func dssXMLUtilsRecreateAndLocate(document *xmldom.Node, elementId string) (*xmldom.Node, error) {
	// TODO : consider switching to DomUtilsCreateDeepCopy
	serializedDoc, err := xmlutils.DomUtilsSerializeNode(document)
	if err != nil {
		return nil, err
	}
	recreatedDocument, err := xmlutils.DomUtilsBuildDOMFromBytes(serializedDoc)
	if err != nil {
		return nil, err
	}
	element := recreatedDocument.DocumentElement()
	if utils.IsStringNotEmpty(elementId) {
		element = xmlutils.XPathUtilsGetElementById(recreatedDocument, elementId)
	}
	return element, nil
}

// DSSXMLUtilsGetDocument returns the linked document to the reference (when applicable). Always
// returns nil in this port: see the file header's "FLAGGED FOR INTEGRATOR" note. Ports
// getDocument(Reference).
func DSSXMLUtilsGetDocument(reference *xmldsig.Reference) model.DSSDocument {
	return nil
}

// DSSXMLUtilsGetDigestOnCanonicalizedBytes computes a digest on a canonicalized value of
// binaries using digestAlgorithm and canonicalizationAlgorithm. The digest is computed "on the
// fly" using stream functionality. Panics with the Java messages when an argument is missing
// (Objects.requireNonNull upstream). Ports getDigestOnCanonicalizedBytes(byte[],
// DigestAlgorithm, String).
func DSSXMLUtilsGetDigestOnCanonicalizedBytes(binaries []byte, digestAlgorithm enumerations.DigestAlgorithm, canonicalizationAlgorithm string) (model.DSSMessageDigest, error) {
	if binaries == nil {
		panic("Binaries cannot be null!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	if canonicalizationAlgorithm == "" {
		panic("Canonicalization algorithm cannot be null!")
	}
	calculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	writer := calculator.Writer()
	cErr := canonicalizer.CanonicalizeBytesTo(binaries, writer)
	_ = writer.Close()
	if cErr != nil {
		return model.DSSMessageDigest{}, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to canonicalize a node : %s", cErr.Error()), cErr)
	}
	return calculator.MessageDigest(digestAlgorithm), nil
}

// DSSXMLUtilsGetDigestOnCanonicalizedNode computes a digest on a canonicalized value of node
// using digestAlgorithm and canonicalizationAlgorithm. The digest is computed "on the fly" using
// stream functionality. Panics with the Java messages when an argument is missing
// (Objects.requireNonNull upstream). Ports getDigestOnCanonicalizedNode(Node, DigestAlgorithm,
// String).
func DSSXMLUtilsGetDigestOnCanonicalizedNode(node *xmldom.Node, digestAlgorithm enumerations.DigestAlgorithm, canonicalizationAlgorithm string) (model.DSSMessageDigest, error) {
	if node == nil {
		panic("Node cannot be null!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	if canonicalizationAlgorithm == "" {
		panic("Canonicalization algorithm cannot be null!")
	}
	calculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	writer := calculator.Writer()
	cErr := canonicalizer.CanonicalizeNodeTo(node, writer)
	_ = writer.Close()
	if cErr != nil {
		return model.DSSMessageDigest{}, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to canonicalize a node : %s", cErr.Error()), cErr)
	}
	return calculator.MessageDigest(digestAlgorithm), nil
}

// DSSXMLUtilsGetDigestOnCanonicalizedInputStream computes a digest on a canonicalized value of
// inputStream using digestAlgorithm and canonicalizationAlgorithm. The digest is computed "on the
// fly" using stream functionality. This method closes inputStream afterwards. Panics with the
// Java messages when an argument is missing (Objects.requireNonNull upstream). Ports
// getDigestOnCanonicalizedInputStream(InputStream, DigestAlgorithm, String).
func DSSXMLUtilsGetDigestOnCanonicalizedInputStream(inputStream io.Reader, digestAlgorithm enumerations.DigestAlgorithm, canonicalizationAlgorithm string) (model.DSSMessageDigest, error) {
	if inputStream == nil {
		panic("InputStream cannot be null!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	if canonicalizationAlgorithm == "" {
		panic("Canonicalization algorithm cannot be null!")
	}
	if closer, ok := inputStream.(io.Closer); ok {
		defer func() { _ = closer.Close() }()
	}
	calculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	writer := calculator.Writer()
	cErr := canonicalizer.CanonicalizeStreamTo(inputStream, writer)
	_ = writer.Close()
	if cErr != nil {
		return model.DSSMessageDigest{}, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to canonicalize a node : %s", cErr.Error()), cErr)
	}
	return calculator.MessageDigest(digestAlgorithm), nil
}

// DSSXMLUtilsGetReferenceDigestAlgos deliberately does not sort its output: it returns
// DigestAlgorithms in first-seen order (the determinism sweep's slice-of-first-seen substitute
// for Java's iteration-order-unstable HashSet), which is deterministic and matches upstream's
// only real use of the result (registering one DetachedSignatureResolver per algorithm, where
// order does not affect the outcome).
