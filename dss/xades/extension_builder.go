// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/ExtensionBuilder.java (DSS 6.5.RC1).
//
// Java's ExtensionBuilder is abstract: it is the base of every XAdES extension (-T, -C, -X, -XL,
// -A, -LT, -LTA, the SignaturePolicyStoreBuilder and the CounterSignatureBuilder) and overrides
// three XAdESBuilder members - alignNodes(), getXmldsigNamespace() and getXadesNamespace() - all
// three of which XAdESBuilder itself calls back into. Go has no method overriding across
// embedding, so this type registers the concrete extension with the base through
// InitExtensionBuilder / InitExtensionBuilderWithVerifier, which forward to
// InitXAdESBuilder / InitXAdESBuilderWithVerifier; the three overrides are promoted to the
// concrete extension and reached through XAdESBuilder's overrides field, the
// TokenBase.InitToken(self) convention of PORTING.md.
//
// The two namespace overrides delegate back to the base implementation - Java's
// params.getXmldsigNamespace() / params.getXadesNamespace() - when the signature being extended
// carries no namespace of its own, which is what XAdESBuilder.XmldsigNamespace/XadesNamespace do.
//
// Errors: IllegalInputException stays spi/exception.IllegalInputException and is returned
// (PORTING.md: throw -> (T, error)). slf4j logging is dropped.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ExtensionBuilder builds a XAdES signature extension.
type ExtensionBuilder struct {
	XAdESBuilder

	// XadesSignature allows accessing the DOM signature representation using XPath.
	// Port of the protected xadesSignature.
	XadesSignature *XAdESSignature

	// CurrentSignatureDom is the current signature being extended.
	// Port of the protected currentSignatureDom.
	CurrentSignatureDom *xmldom.Node

	// QualifyingPropertiesDom is the signature qualifying properties.
	// Port of the protected qualifyingPropertiesDom.
	QualifyingPropertiesDom *xmldom.Node

	// UnsignedPropertiesDom is the unsigned properties.
	// Port of the protected unsignedPropertiesDom.
	UnsignedPropertiesDom *xmldom.Node

	// UnsignedSignaturePropertiesDom contains the unsigned signature properties.
	// Port of the protected unsignedSignaturePropertiesDom.
	UnsignedSignaturePropertiesDom *xmldom.Node

	// DocumentAnalyzer is the used document analyzer.
	// Port of the protected documentAnalyzer.
	DocumentAnalyzer *XMLDocumentAnalyzer
}

// InitExtensionBuilder registers the concrete extension with its base.
// Port of the protected empty ExtensionBuilder() constructor.
func (b *ExtensionBuilder) InitExtensionBuilder(self XAdESBuilderOverrides) {
	b.InitXAdESBuilder(self)
}

// InitExtensionBuilderWithVerifier registers the concrete extension with its base and stores the
// CertificateVerifier. Port of the protected ExtensionBuilder(CertificateVerifier) constructor.
func (b *ExtensionBuilder) InitExtensionBuilderWithVerifier(self XAdESBuilderOverrides,
	certificateVerifier validation.CertificateVerifier) {
	b.InitXAdESBuilderWithVerifier(self, certificateVerifier)
}

// InitializeSignatureBuilder initializes all variables to be used for signature extension.
// Port of the protected #initializeSignatureBuilder.
func (b *ExtensionBuilder) InitializeSignatureBuilder(signature *XAdESSignature) (*XAdESSignature, error) {
	b.XadesSignature = signature
	b.CurrentSignatureDom = b.XadesSignature.SignatureElement()

	b.XadesPath = b.XadesSignature.XAdESPaths()

	// We ensure that all XML segments needed for the construction of the extension -T are
	// present. If a segment does not exist then it is created.
	if err := b.EnsureUnsignedProperties(); err != nil {
		return nil, err
	}
	if err := b.EnsureUnsignedSignatureProperties(); err != nil {
		return nil, err
	}
	if err := b.EnsureSignedDataObjectProperties(); err != nil {
		return nil, err
	}

	return b.XadesSignature, nil
}

// EnsureUnsignedProperties returns or creates (if it does not exist) the UnsignedPropertiesType
// DOM object. Port of the protected #ensureUnsignedProperties.
func (b *ExtensionBuilder) EnsureUnsignedProperties() error {
	qualifyingPropertiesNodeList, err := xmlutils.XPathUtilsGetNodeList(b.CurrentSignatureDom,
		b.XadesPath.QualifyingPropertiesPath())
	if err != nil {
		return err
	}
	if len(qualifyingPropertiesNodeList) != 1 {
		return exception.NewIllegalInputException("The signature does not contain QualifyingProperties " +
			"element (or contains more than one)! Extension is not possible.")
	}

	b.QualifyingPropertiesDom = qualifyingPropertiesNodeList[0]

	unsignedPropertiesNodeList, err := xmlutils.XPathUtilsGetNodeList(b.CurrentSignatureDom,
		b.XadesPath.UnsignedPropertiesPath())
	if err != nil {
		return err
	}
	length := len(unsignedPropertiesNodeList)
	if length == 1 {
		b.UnsignedPropertiesDom = unsignedPropertiesNodeList[0]
	} else if length == 0 {
		currentElements, err := b.CurrentXAdESElements()
		if err != nil {
			return err
		}
		b.UnsignedPropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.QualifyingPropertiesDom,
			b.overrides.XadesNamespace(), currentElements.ElementUnsignedProperties())
		if b.Params.IsPrettyPrint() {
			aligned, err := DSSXMLUtilsAlignChildrenIndents(b.QualifyingPropertiesDom)
			if err != nil {
				return err
			}
			b.QualifyingPropertiesDom = aligned
			b.UnsignedPropertiesDom, err = xmlutils.XPathUtilsGetNode(b.CurrentSignatureDom,
				b.XadesPath.UnsignedPropertiesPath())
			if err != nil {
				return err
			}
		}
	} else {
		return exception.NewIllegalInputException("The signature contains more then one " +
			"UnsignedProperties element! Extension is not possible.")
	}
	return nil
}

// EnsureUnsignedSignatureProperties returns or creates (if it does not exist) the
// UnsignedSignaturePropertiesType DOM object.
// Port of the protected #ensureUnsignedSignatureProperties.
func (b *ExtensionBuilder) EnsureUnsignedSignatureProperties() error {
	unsignedSignaturePropertiesNodeList, err := xmlutils.XPathUtilsGetNodeList(b.CurrentSignatureDom,
		b.XadesPath.UnsignedSignaturePropertiesPath())
	if err != nil {
		return err
	}
	length := len(unsignedSignaturePropertiesNodeList)
	if length == 1 {
		b.UnsignedSignaturePropertiesDom = unsignedSignaturePropertiesNodeList[0]
	} else if length == 0 {
		currentElements, err := b.CurrentXAdESElements()
		if err != nil {
			return err
		}
		b.UnsignedSignaturePropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom,
			b.UnsignedPropertiesDom, b.overrides.XadesNamespace(),
			currentElements.ElementUnsignedSignatureProperties())
		if b.Params.IsPrettyPrint() {
			indented, err := DSSXMLUtilsIndentAndReplace(b.DocumentDom, b.UnsignedPropertiesDom)
			if err != nil {
				return err
			}
			b.UnsignedPropertiesDom = indented
			b.UnsignedSignaturePropertiesDom, err = xmlutils.XPathUtilsGetNode(b.CurrentSignatureDom,
				b.XadesPath.UnsignedSignaturePropertiesPath())
			if err != nil {
				return err
			}
		}
	} else {
		return exception.NewIllegalInputException("The signature contains more than one " +
			"UnsignedSignatureProperties element! Extension is not possible.")
	}
	return nil
}

// EnsureSignedDataObjectProperties returns or creates (if it does not exist) the
// SignedDataObjectProperties DOM object.
// Port of the protected #ensureSignedDataObjectProperties.
func (b *ExtensionBuilder) EnsureSignedDataObjectProperties() error {
	signedDataObjectPropertiesNodeList, err := xmlutils.XPathUtilsGetNodeList(b.CurrentSignatureDom,
		b.XadesPath.SignedDataObjectPropertiesPath())
	if err != nil {
		return err
	}
	if len(signedDataObjectPropertiesNodeList) > 1 {
		return exception.NewIllegalInputException("The signature contains more than one " +
			"SignedDataObjectProperties element! Extension is not possible.")
	}
	return nil
}

// IndentIfPrettyPrint indents nodeToIndent if pretty-print is enabled.
// Port of the protected #indentIfPrettyPrint.
func (b *ExtensionBuilder) IndentIfPrettyPrint(nodeToIndent, oldNode *xmldom.Node) (*xmldom.Node, error) {
	if b.Params.IsPrettyPrint() {
		return DSSXMLUtilsIndentAndExtend(b.DocumentDom, nodeToIndent, oldNode)
	}
	return nodeToIndent, nil
}

// AlignNodes aligns children indents. Port of the overridden protected #alignNodes.
func (b *ExtensionBuilder) AlignNodes() {
	// Java calls alignChildrenIndents for effect and ignores the node it returns; the Go
	// signature adds an error the same way, and it is discarded here for the same reason.
	if b.UnsignedSignaturePropertiesDom != nil {
		DSSXMLUtilsAlignChildrenIndents(b.UnsignedSignaturePropertiesDom)
	}
	if b.QualifyingPropertiesDom != nil {
		DSSXMLUtilsAlignChildrenIndents(b.QualifyingPropertiesDom)
	}
}

// RemoveNode removes the given node from its parent and returns the text of the removed node's
// next TEXT sibling (empty when that sibling does not exist).
// Port of the protected #removeNode.
func (b *ExtensionBuilder) RemoveNode(nodeToRemove *xmldom.Node) string {
	text := ""
	if nodeToRemove != nil {
		nextSibling := nodeToRemove.NextSibling
		if nextSibling != nil && xmldom.Text == nextSibling.Kind {
			text = nextSibling.Value
			b.UnsignedSignaturePropertiesDom.RemoveChild(nextSibling)
		}
		b.UnsignedSignaturePropertiesDom.RemoveChild(nodeToRemove)
	}
	return text
}

// XmldsigNamespace returns the currently used XMLDSig namespace, determined from the signature,
// from the parameters, or falling back to the default.
// Port of the overridden protected #getXmldsigNamespace.
func (b *ExtensionBuilder) XmldsigNamespace() *common.DSSNamespace {
	xmldsigNamespace := b.XadesSignature.XmldSigNamespace()
	if xmldsigNamespace == nil {
		// Upstream warns "Current XMLDSig namespace not found in the signature".
		xmldsigNamespace = b.Params.XmldsigNamespace()
		if xmldsigNamespace == nil {
			// Upstream warns "Current XMLDSig namespace not found in the parameters
			// (use the default XMLDSig)".
			xmldsigNamespace = common.XMLDSigNS
		}
	}
	return xmldsigNamespace
}

// XadesNamespace returns the currently used XAdES namespace, determined from the signature, from
// the parameters, or falling back to the default XAdES 1.3.2.
// Port of the overridden protected #getXadesNamespace.
func (b *ExtensionBuilder) XadesNamespace() *common.DSSNamespace {
	xadesNamespace := b.XadesSignature.XadesNamespace()
	if xadesNamespace == nil {
		// Upstream warns "Current XAdES namespace not found in the signature".
		xadesNamespace = b.Params.XadesNamespace()
		if xadesNamespace == nil {
			// Upstream warns "Current XAdES namespace not found in the parameters
			// (use the default XAdES 1.3.2)".
			xadesNamespace = definition.XAdESNamespace_XADES_132
		}
	}
	return xadesNamespace
}

// AssertUnsignedPropertiesExtensionPossible verifies whether signature extension is possible, as
// it must not contain evidence records.
// Port of the protected #assertUnsignedPropertiesExtensionPossible.
func (b *ExtensionBuilder) AssertUnsignedPropertiesExtensionPossible(
	signature validation.AdvancedSignature) error {
	if utils.IsCollectionNotEmpty(signature.EmbeddedEvidenceRecords()) {
		return exception.NewIllegalInputException("Signature extension is not possible. " +
			"The signature already contains en embedded evidence record.")
	}
	return nil
}

// InitDocumentAnalyzer initializes the document analyzer and other signature properties.
// Port of the protected #initDocumentAnalyzer.
func (b *ExtensionBuilder) InitDocumentAnalyzer(document model.DSSDocument) (*XMLDocumentAnalyzer, error) {
	b.Params = NewXAdESSignatureParameters()

	documentAnalyzer, err := NewXMLDocumentAnalyzer(document)
	if err != nil {
		return nil, err
	}
	b.DocumentAnalyzer = documentAnalyzer
	b.DocumentDom = b.DocumentAnalyzer.RootElement()

	return b.DocumentAnalyzer, nil
}
