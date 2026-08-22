// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESAttribute.java
// (DSS 6.5.RC1).
//
// hashDataInfoTransformPath ports the static field XAdES111Path.HASH_DATA_INFO_TRANSFORM_PATH
// (Java: "./xades111:HashDataInfo/xades111:Transforms/xades111:Transform"). That field is not an
// XAdESPath interface method - it is a class-qualified static constant of the concrete
// XAdES111Path type - and xades/definition (frozen, landed by a sibling chunk) exports only the
// XAdESPath interface surface, not per-format static XPathQuery fields such as this one. Building
// it locally with the same common.FromCurrentPosition/DSSElement machinery the frozen package
// itself is built from reproduces the same XPath string without editing that package; see this
// porter's report for the flagged gap.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/identifier"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// hashDataInfoTransformPath ports XAdES111Path.HASH_DATA_INFO_TRANSFORM_PATH.
var hashDataInfoTransformPath = common.FromCurrentPosition(
	definition.XAdES111ElementHashDataInfo,
	definition.XAdES111ElementTransforms,
	common.XMLDSigElementTransform,
)

// XAdESAttribute represents a XAdES attribute. Port of the class XAdESAttribute, implementing
// spi/validation.SignatureAttribute.
type Attribute struct {
	// element is the corresponding element.
	element *xmldom.Node

	// xadesPaths is the XPath list to use.
	xadesPaths definition.XAdESPath

	// localName is the tag name of the element, lazily computed.
	localName string

	// identifier identifies the instance, lazily computed.
	identifier *AttributeIdentifier
}

// newXAdESAttribute is the port of the package-private XAdESAttribute(Element, XAdESPath)
// constructor.
func newXAdESAttribute(element *xmldom.Node, xadesPaths definition.XAdESPath) *Attribute {
	return &Attribute{element: element, xadesPaths: xadesPaths}
}

// Name returns the local name of the element. Port of getName().
func (a *Attribute) Name() string {
	if a.localName == "" {
		a.localName = a.element.Name.Local
	}
	return a.localName
}

// Element returns the current element. Port of getElement().
func (a *Attribute) Element() *xmldom.Node {
	return a.element
}

// Namespace returns the namespace of the element. Port of getNamespace().
func (a *Attribute) Namespace() string {
	return a.element.Name.Space
}

// NodeList returns the node list found by the given XPath query. Port of getNodeList(XPathQuery).
func (a *Attribute) NodeList(xPathQuery common.XPathQuery) []*xmldom.Node {
	nodes, err := xmlutils.XPathUtilsGetNodeList(a.element, xPathQuery)
	if err != nil {
		return nil
	}
	return nodes
}

// TimestampCanonicalizationMethod returns the TimeStamp Canonicalization Method. Port of
// getTimestampCanonicalizationMethod().
//
// LOG.warn("Unable to retrieve the canonicalization algorithm") is dropped per PORTING.md.
func (a *Attribute) TimestampCanonicalizationMethod() string {
	canonicalizationMethod, err := xmlutils.XPathUtilsGetValue(a.element, common.XMLDSigPathCanonicalizationAlgorithmPath)
	if err != nil {
		canonicalizationMethod = ""
	}
	if utils.IsStringEmpty(canonicalizationMethod) {
		nodeList, err := xmlutils.XPathUtilsGetNodeList(a.element, hashDataInfoTransformPath)
		if err == nil && len(nodeList) == 1 {
			transform := nodeList[0]
			canonicalizationMethod = transform.AttrValue("", common.XMLDSigAttributeAlgorithm.AttributeName())
		}
	}
	return canonicalizationMethod
}

// TimestampIncludedReferences returns the list of TimestampIncludes in case of
// IndividualDataObjectsTimestamp, nil if it does not contain any includes. Port of
// getTimestampIncludedReferences().
func (a *Attribute) TimestampIncludedReferences() []*validation.TimestampInclude {
	currentIncludePath := a.xadesPaths.CurrentInclude()
	if currentIncludePath == nil {
		return nil
	}
	timestampIncludes, err := xmlutils.XPathUtilsGetNodeList(a.element, currentIncludePath)
	if err != nil || len(timestampIncludes) == 0 {
		return nil
	}
	includes := make([]*validation.TimestampInclude, 0, len(timestampIncludes))
	for _, include := range timestampIncludes {
		uri := xmlutils.DomUtilsGetId(include.AttrValue("", definition.XAdES132AttributeURI.AttributeName()))
		referencedData := include.AttrValue("", definition.XAdES132AttributeReferencedData.AttributeName())
		includes = append(includes, validation.NewTimestampIncludeWithURI(uri, referencedData == "true"))
	}
	return includes
}

// Identifier gets the attribute identifier. Port of getIdentifier(), implementing
// spi/validation.SignatureAttribute.
func (a *Attribute) Identifier() identifier.SignatureAttributeIdentifier {
	if a.identifier == nil {
		a.identifier = AttributeIdentifierBuild(a.element)
	}
	return a.identifier.SignatureAttributeIdentifier
}

// Equals ports equals(Object): two XAdESAttributes are equal when their identifiers are equal.
func (a *Attribute) Equals(other *Attribute) bool {
	if a == other {
		return true
	}
	if other == nil {
		return false
	}
	selfID := a.Identifier()
	otherID := other.Identifier()
	return selfID.Equals(&otherID)
}

// String ports toString().
func (a *Attribute) String() string {
	return a.Name()
}
