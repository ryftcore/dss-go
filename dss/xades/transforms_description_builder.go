// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/TransformsDescriptionBuilder.java
// (DSS 6.5.RC1).
//
// Builds a user-friendly description for a provided ds:Transforms element. Already relied upon
// by the landed xades_reference_validation.go (NewTransformsDescriptionBuilder(...).Build()) and
// by the XAdESSignaturePolicy forward dependency (see xades_signature.go's header) for
// SignaturePolicy#getTransformsDescription().
//
// org.apache.xml.security.transforms.Transforms / org.apache.xml.security.c14n.Canonicalizer
// algorithm URI constants have no Go class to mirror one-to-one; their string values are the
// same internal/xmldsig already carries (TransformEnvelopedSignature, TransformBase64Decode,
// TransformXPath2Filter, TransformXPath, TransformXSLT and the various xmlc14n canonicalization
// algorithm URIs via internal/xmlc14n), reused here directly per PORTING.md ("never invent an
// algorithm URI - copy it verbatim").
package xades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// transformsDescriptionPresentableNames mirrors the static presentableTransformationNames map.
// The Canonicalizer.ALGO_ID_C14N* URIs Java names explicitly are the same strings
// internal/xmldsig already carries as its TransformC14N* constants (with/without the
// "#WithComments" suffix distinguishing the with/omit-comments variants).
var transformsDescriptionPresentableNames = map[string]string{
	xmldsig.TransformEnvelopedSignature: "Enveloped Signature Transform",
	xmldsig.TransformBase64Decode:       "Base64 Decoding",

	xmldsig.TransformXPath2Filter: "XPath Filter 2.0 Transform",
	xmldsig.TransformXPath:        "XPath filtering",
	xmldsig.TransformXSLT:         "XSLT Transform",

	xmldsig.TransformC14NWithComments:     "Canonical XML 1.0 with Comments",
	xmldsig.TransformC14N11WithComments:   "Canonical XML 1.1 with Comments",
	xmldsig.TransformC14NExclWithComments: "Exclusive XML Canonicalization 1.0 with Comments",

	xmldsig.TransformC14N:     "Canonical XML 1.0 (omits comments)",
	xmldsig.TransformC14N11:   "Canonical XML 1.1 (omits comments)",
	xmldsig.TransformC14NExcl: "Exclusive Canonical XML (omits comments)",
}

// TransformsDescriptionBuilder builds a user-friendly description for the provided
// 'ds:Transforms' element. Port of the class TransformsDescriptionBuilder.
type TransformsDescriptionBuilder struct {
	// transforms is the ds:Transforms element.
	transforms *xmldom.Node
}

// NewTransformsDescriptionBuilder is the default constructor. Port of
// TransformsDescriptionBuilder(Element).
func NewTransformsDescriptionBuilder(transforms *xmldom.Node) *TransformsDescriptionBuilder {
	return &TransformsDescriptionBuilder{transforms: transforms}
}

// Build builds a list of Strings describing the 'ds:Transforms' element. Returns an empty list
// if transforms are not found or cannot be extracted. Port of build().
func (b *TransformsDescriptionBuilder) Build() []string {
	transformsList := []string{}
	if b.transforms != nil {
		for child := b.transforms.FirstChild; child != nil; child = child.NextSibling {
			if child.Kind == xmldom.Element {
				transformsList = append(transformsList, b.buildTransformationName(child))
			}
		}
	}
	return transformsList
}

// buildTransformationName returns a complete description string for the given transformation
// node. Port of the private buildTransformationName(Element).
func (b *TransformsDescriptionBuilder) buildTransformationName(transformation *xmldom.Node) string {
	algorithmURI := transformation.AttrValue("", common.XMLDSigAttributeAlgorithm.AttributeName())
	algorithm := algorithmURI
	if presentable, ok := transformsDescriptionPresentableNames[algorithmURI]; ok {
		algorithm = presentable
	}

	var sb strings.Builder
	sb.WriteString(algorithm)
	if transformation.FirstChild != nil {
		sb.WriteString(" (")
		hasValues := false

		for parameterNode := transformation.FirstChild; parameterNode != nil; parameterNode = parameterNode.NextSibling {
			if parameterNode.Kind != xmldom.Element {
				continue
			}

			// attach attribute values
			for _, attribute := range parameterNode.Attrs {
				attrName := attribute.Name.Local
				attrValue := attribute.Value
				if algorithmURI == attrValue {
					continue // skip the case when the algorithm uri is defined in a child node
				}
				if hasValues {
					sb.WriteString("; ")
				}
				sb.WriteString(attrName)
				sb.WriteString(": ")
				sb.WriteString(attrValue)
				hasValues = true
			}

			// attach node value
			parameterValueNode := parameterNode.FirstChild
			if parameterValueNode != nil && parameterValueNode.Kind == xmldom.Text &&
				utils.IsStringNotBlank(parameterValueNode.TextContent()) {
				if hasValues {
					sb.WriteString("; ")
				}
				sb.WriteString(parameterNode.Name.Local)
				sb.WriteString(": ")
				sb.WriteString(parameterValueNode.TextContent())
				hasValues = true
			}
		}
		sb.WriteString(")")
	}
	return sb.String()
}
