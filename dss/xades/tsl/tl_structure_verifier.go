// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/tsl/TLStructureVerifier.java (DSS 6.5.RC1).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S4D_BRIEF.md's "flag needs in notes" rule, for the
// integrator to arbitrate): upstream's XSD validation (TrustedListUtils.getInstance() /
// TrustedList211Utils.getInstance(), both eu.europa.esig.dss.jaxb.common.XSDAbstractUtils
// implementations bundling the TL V5/V6 XSD schemas as JAXB resources) has no landed Go home,
// mirroring xades_structure_validator.go's identical BLOCKED FORWARD DEPENDENCY note for the
// XAdES schema family. That file's XAdESStructureXSDUtils marker type and
// DSSXMLUtilsValidateAgainstXSD(xsdUtils, source) helper are reused verbatim here (both already
// resolve to an untyped interface{} + a single wrapping call, so no new machinery is needed);
// this file only adds the two TL-specific singleton getters:
//
//	func TrustedListXSDUtils() xades.XAdESStructureXSDUtils    // eu.europa.esig.trustedlist.TrustedListUtils#getInstance() (TL V6)
//	func TrustedList211XSDUtils() xades.XAdESStructureXSDUtils // eu.europa.esig.trustedlist211.TrustedList211Utils#getInstance() (TL V5)
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// trustedListParentElement is the expected local name of a Trusted List document's root
// element.
const trustedListParentElement = "TrustServiceStatusList"

func init() {
	xmlutils.XPathUtilsRegisterNamespace(definition.TrustedListNamespaceNS)
}

// TLStructureVerifier verifies conformity of a TL to the defined TLVersion.
// NOTE: The class currently handles validation of only V5 and V6 Trusted List versions.
type TLStructureVerifier struct {
	// acceptedTLVersions is the list of acceptable TL versions.
	acceptedTLVersions []int

	// signingMode defines whether the current validation of the XML Trusted List is performed
	// for signing.
	signingMode bool
}

// NewTLStructureVerifier is the default constructor. Accepts TL V5 and TL V6.
func NewTLStructureVerifier() *TLStructureVerifier {
	return &TLStructureVerifier{}
}

// SetAcceptedTLVersions sets a list of acceptable TL version. When defined, an error message
// will be produced for Trusted Lists with a different version. Port of
// #setAcceptedTLVersions, chainable.
func (v *TLStructureVerifier) SetAcceptedTLVersions(acceptedTLVersions []int) *TLStructureVerifier {
	v.acceptedTLVersions = acceptedTLVersions
	return v
}

// SetSigningMode sets whether the current operation is the XML Trusted List signing. If
// enabled, verifies that no ds:Signature element is present within the XML Trusted List.
// Otherwise, verifies presence and validity of the ds:Signature element.
// Default: FALSE (verifies that the signature is not present). Port of #setSigningMode,
// chainable.
func (v *TLStructureVerifier) SetSigningMode(signingMode bool) *TLStructureVerifier {
	v.signingMode = signingMode
	return v
}

// Validate validates the Trusted List's conformity to the specified TLVersion. Port of the
// #validate(DSSDocument, Integer) overload.
//
// Panics when dssDocument is nil (Java's Objects.requireNonNull).
func (v *TLStructureVerifier) Validate(dssDocument model.DSSDocument, tlVersion *int) ([]string, error) {
	if dssDocument == nil {
		panic("Document to be validated cannot be null!")
	}
	documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(dssDocument)
	if err != nil {
		return nil, err
	}
	return v.ValidateDocument(documentDom, tlVersion), nil
}

// ValidateDocument validates the Trusted List's conformity to the specified TLVersion. Port of
// the #validate(Document, Integer) overload.
//
// Panics when document is nil (Java's Objects.requireNonNull).
func (v *TLStructureVerifier) ValidateDocument(documentDom *xmldom.Node, tlVersion *int) []string {
	if documentDom == nil {
		panic("Document to be validated cannot be null!")
	}
	if tlVersion == nil {
		return []string{"No TLVersion has been found!"}
	}

	if utils.IsCollectionEmpty(v.acceptedTLVersions) {
		// Upstream logs "No acceptable TL Versions have been defined. The structural
		// validation is skipped.".
		return nil
	}

	var errors []string
	errors = append(errors, v.validateNamespace(documentDom)...)

	if !xadesTrustedListUtilsContainsInt(v.acceptedTLVersions, *tlVersion) {
		errors = append(errors, fmt.Sprintf("The TL Version '%d' is not acceptable!", *tlVersion))
	} else if TLV5Identifier == *tlVersion {
		errors = append(errors, v.validateTrustedListV5(documentDom)...)
	} else if TLV6Identifier == *tlVersion {
		errors = append(errors, v.validateTrustedListV6(documentDom)...)
	}

	return errors
}

// validateTrustedListV5 validates the Trusted List XML document against the TL V5 definition.
// Port of the protected #validateTrustedListV5(Document).
func (v *TLStructureVerifier) validateTrustedListV5(documentDom *xmldom.Node) []string {
	var errors []string
	xsdValidationErrors := v.validateAgainstXSD(documentDom, TrustedList211XSDUtils())
	if utils.IsCollectionNotEmpty(xsdValidationErrors) {
		errors = append(errors, xsdValidationErrors...)
	}
	signatureElement := v.getSignatureElement(documentDom)
	errors = append(errors, v.verifySignatureElementPresence(signatureElement)...)
	return errors
}

// validateTrustedListV6 validates the Trusted List XML document against the TL V6 definition.
// Port of the protected #validateTrustedListV6(Document).
func (v *TLStructureVerifier) validateTrustedListV6(documentDom *xmldom.Node) []string {
	var errors []string
	xsdValidationErrors := v.validateAgainstXSD(documentDom, TrustedListXSDUtils())
	if utils.IsCollectionNotEmpty(xsdValidationErrors) {
		errors = append(errors, xsdValidationErrors...)
	}

	signatureElement := v.getSignatureElement(documentDom)
	errors = append(errors, v.verifySignatureElementPresence(signatureElement)...)

	v2ConformityErrors := v.validateSignatureElement(signatureElement, true)
	if utils.IsCollectionNotEmpty(v2ConformityErrors) {
		errors = append(errors, v2ConformityErrors...)
	}
	return errors
}

// validateAgainstXSD ports the private validateAgainstXSD(Document, XSDAbstractUtils).
func (v *TLStructureVerifier) validateAgainstXSD(documentDom *xmldom.Node, xsdUtils xades.XAdESStructureXSDUtils) []string {
	return xades.DSSXMLUtilsValidateAgainstXSD(xsdUtils, documentDom)
}

// trustedListXSDUtils and trustedList211XSDUtils are opaque placeholders standing in for
// eu.europa.esig.trustedlist.TrustedListUtils#getInstance() (TL V6) and
// eu.europa.esig.trustedlist211.TrustedList211Utils#getInstance() (TL V5). See the file header's
// BLOCKED FORWARD DEPENDENCY note: no bundled TL XSD schema resources exist anywhere in this
// port yet, so xades.DSSXMLUtilsValidateAgainstXSD (which never inspects its xsdUtils argument)
// always reports no errors; these placeholders are never dereferenced.
var (
	trustedListXSDUtils    xades.XAdESStructureXSDUtils = struct{}{}
	trustedList211XSDUtils xades.XAdESStructureXSDUtils = struct{}{}
)

// TrustedListXSDUtils ports eu.europa.esig.trustedlist.TrustedListUtils#getInstance() (TL V6).
// See the note above.
func TrustedListXSDUtils() xades.XAdESStructureXSDUtils {
	return trustedListXSDUtils
}

// TrustedList211XSDUtils ports eu.europa.esig.trustedlist211.TrustedList211Utils#getInstance()
// (TL V5). See the note above.
func TrustedList211XSDUtils() xades.XAdESStructureXSDUtils {
	return trustedList211XSDUtils
}

// validateNamespace ports the private validateNamespace(Document).
func (v *TLStructureVerifier) validateNamespace(documentDom *xmldom.Node) []string {
	documentElement := documentDom.DocumentElement()
	if trustedListParentElement != documentElement.Name.Local ||
		definition.TrustedListNamespaceNS.Uri() != documentElement.Name.Space {
		return []string{fmt.Sprintf("The root of XML Trusted List shall be %s:%s element!",
			definition.TrustedListNamespaceNS.Prefix(), trustedListParentElement)}
	}
	return nil
}

// getSignatureElement ports the private getSignatureElement(Document).
func (v *TLStructureVerifier) getSignatureElement(documentDom *xmldom.Node) *xmldom.Node {
	documentElement := documentDom.DocumentElement()
	return v.getChildElement(documentElement, common.XMLDSigElementSignature)
}

// verifySignatureElementPresence ports the private verifySignatureElementPresence(Element).
func (v *TLStructureVerifier) verifySignatureElementPresence(dsSignature *xmldom.Node) []string {
	if v.signingMode {
		if dsSignature != nil {
			return []string{"The ds:Signature element shall not be present for XML Trusted List signing!"}
		}
		// no ds:Signature is expected on signing
	} else if dsSignature == nil {
		return []string{"No ds:Signature element is present!"}
	}
	return nil
}

// validateSignatureElement ports the private validateSignatureElement(Element, boolean).
func (v *TLStructureVerifier) validateSignatureElement(dsSignature *xmldom.Node, v2Expected bool) []string {
	if dsSignature == nil || v.signingMode {
		return nil
	}

	// NOTE: manual parsing is used for performance reasons
	objects := v.getChildElements(dsSignature, common.XMLDSigElementObject)
	if utils.IsCollectionEmpty(objects) {
		return []string{"No ds:Object elements are present!"}
	}
	qualifyingProperties := v.getQualifyingPropertiesElement(objects)
	if qualifyingProperties == nil {
		return []string{"No xades:QualifyingProperties element has been found!"}
	}
	signedProperties := v.getChildElement(qualifyingProperties, definition.XAdES132ElementSignedProperties)
	if signedProperties == nil {
		return []string{"No xades:SignedProperties element has been found!"}
	}
	signedSignatureProperties := v.getChildElement(signedProperties, definition.XAdES132ElementSignedSignatureProperties)
	if signedSignatureProperties == nil {
		return []string{"No xades:SignedSignatureProperties element has been found!"}
	}

	var errorMessages []string

	signingCertificateElements := v.getMultipleElements(signedSignatureProperties,
		definition.XAdES132ElementSigningCertificate, definition.XAdES132ElementSigningCertificateV2)
	if utils.IsCollectionNotEmpty(signingCertificateElements) {
		for _, signingCertificate := range signingCertificateElements {
			if v2Expected != v.doesMatch(signingCertificate, definition.XAdES132ElementSigningCertificateV2) {
				errorMessages = append(errorMessages, fmt.Sprintf("%s element shall not be present!", signingCertificate.Name.Local))
			}
		}
	} else {
		v2Suffix := ""
		if v2Expected {
			v2Suffix = "V2"
		}
		errorMessages = append(errorMessages, fmt.Sprintf("No xades:SigningCertificate%s element has been found!", v2Suffix))
	}

	signatureProductionPlaceElements := v.getMultipleElements(signedSignatureProperties,
		definition.XAdES132ElementSignatureProductionPlace, definition.XAdES132ElementSignatureProductionPlaceV2)
	if utils.IsCollectionNotEmpty(signatureProductionPlaceElements) {
		for _, signatureProductionPlace := range signatureProductionPlaceElements {
			if v2Expected != v.doesMatch(signatureProductionPlace, definition.XAdES132ElementSignatureProductionPlaceV2) {
				errorMessages = append(errorMessages, fmt.Sprintf("%s element shall not be present!", signatureProductionPlace.Name.Local))
			}
		}
	}

	signerRoleElements := v.getMultipleElements(signedSignatureProperties,
		definition.XAdES132ElementSignerRole, definition.XAdES132ElementSignerRoleV2)
	if utils.IsCollectionNotEmpty(signerRoleElements) {
		for _, signerRole := range signerRoleElements {
			if v2Expected != v.doesMatch(signerRole, definition.XAdES132ElementSignerRoleV2) {
				errorMessages = append(errorMessages, fmt.Sprintf("%s element shall not be present!", signerRole.Name.Local))
			}
		}
	}

	return errorMessages
}

// getQualifyingPropertiesElement ports the private getQualifyingPropertiesElement(List<Element>).
func (v *TLStructureVerifier) getQualifyingPropertiesElement(objects []*xmldom.Node) *xmldom.Node {
	for _, object := range objects {
		qualifyingProperties := v.getChildElement(object, definition.XAdES132ElementQualifyingProperties)
		if qualifyingProperties != nil {
			return qualifyingProperties
		}
	}
	return nil
}

// getMultipleElements ports the private getMultipleElements(Element, XAdESElement...).
func (v *TLStructureVerifier) getMultipleElements(signedSignatureProperties *xmldom.Node, targetElements ...common.DSSElement) []*xmldom.Node {
	var result []*xmldom.Node
	for _, targetElement := range targetElements {
		signingCertificates := v.getChildElements(signedSignatureProperties, targetElement)
		if utils.IsCollectionNotEmpty(signingCertificates) {
			result = append(result, signingCertificates...)
		}
	}
	return result
}

// getChildElement ports the private getChildElement(Element, DSSElement).
func (v *TLStructureVerifier) getChildElement(parentElement *xmldom.Node, targetElement common.DSSElement) *xmldom.Node {
	childrenList := v.getChildElements(parentElement, targetElement)
	if len(childrenList) == 1 {
		return childrenList[0]
	}
	return nil
}

// getChildElements ports the private getChildElements(Element, DSSElement).
func (v *TLStructureVerifier) getChildElements(parentElement *xmldom.Node, targetElement common.DSSElement) []*xmldom.Node {
	var children []*xmldom.Node
	for childNode := parentElement.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if xmldom.Element != childNode.Kind {
			continue
		}
		if v.doesMatch(childNode, targetElement) {
			children = append(children, childNode)
		}
	}
	return children
}

// doesMatch ports the private doesMatch(Element, DSSElement).
func (v *TLStructureVerifier) doesMatch(element *xmldom.Node, dssElement common.DSSElement) bool {
	return dssElement.IsSameTagName(element.Name.Local) && dssElement.URI() == element.Name.Space
}

// xadesTrustedListUtilsContainsInt reports whether values contains target, mirroring Java's
// List<Integer>#contains(Integer).
func xadesTrustedListUtilsContainsInt(values []int, target int) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
