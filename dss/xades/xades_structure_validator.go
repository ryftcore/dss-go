// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESStructureValidator.java
// (DSS 6.5.RC1).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S4D_BRIEF.md's "flag needs in notes" rule, for the
// integrator to arbitrate): upstream's getUtils(XAdESPath) resolves an XSDAbstractUtils
// implementation (eu.europa.esig.xades.XAdES111Utils / XAdES122Utils / XAdES319132Utils) drawn
// from the separate "dss-jaxb-xades"-family artifacts (XSD schemas bundled as JAXB resources)
// and hands it to DSSXMLUtils.validateAgainstXSD(XSDAbstractUtils, Source), which runs a
// javax.xml.validation.Validator against it. Nothing in PORTING_PLAN.md assigns those XSD
// schema resources or a javax.xml.validation-equivalent Go package a home yet - there is no
// bundled-schema or XML-Schema-validator package anywhere in this port. This file therefore
// assumes, without being able to point at a landed package, that a sibling phase supplies:
//
//	func XAdES111XSDUtils() XAdESStructureXSDUtils    // eu.europa.esig.xades.XAdES111Utils#getInstance()
//	func XAdES122XSDUtils() XAdESStructureXSDUtils    // eu.europa.esig.xades.XAdES122Utils#getInstance()
//	func XAdES319132XSDUtils() XAdESStructureXSDUtils // eu.europa.esig.xades.XAdES319132Utils#getInstance()
//	func DSSXMLUtilsValidateAgainstXSD(xsdUtils XAdESStructureXSDUtils, source *xmldom.Node) []string
//
// XAdESStructureXSDUtils stands in for eu.europa.esig.dss.jaxb.common.XSDAbstractUtils, an
// opaque marker type here since this file never calls a method on it directly (only getUtils
// selects one, and DSSXMLUtilsValidateAgainstXSD consumes it).
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
)

// XAdESStructureXSDUtils stands in for eu.europa.esig.dss.jaxb.common.XSDAbstractUtils; see the
// file header BLOCKED FORWARD DEPENDENCY note.
type XAdESStructureXSDUtils interface{}

// XAdESStructureValidator validates a structure of a XAdES signature against a corresponding
// XSD. Port of the class XAdESStructureValidator.
type XAdESStructureValidator struct {
	// signatureElement is the ds:Signature element to validate structure for.
	signatureElement *xmldom.Node

	// xadesPath is the corresponding XAdES path defining the format of the XAdES signature.
	xadesPath definition.XAdESPath

	// errors caches the list of validation errors.
	errors []string
	// errorsSet reports whether errors has been computed (Go has no null-slice/"uncomputed"
	// distinction as clean as Java's null field).
	errorsSet bool
}

// newXAdESStructureValidator is the port of the protected constructor
// XAdESStructureValidator(Element, XAdESPath).
func newXAdESStructureValidator(signatureElement *xmldom.Node, xadesPath definition.XAdESPath) *XAdESStructureValidator {
	return &XAdESStructureValidator{signatureElement: signatureElement, xadesPath: xadesPath}
}

// Validate validates the signature against the corresponding XSD and returns whether the
// signature has a valid XML structure. Port of validate().
func (v *XAdESStructureValidator) Validate() bool {
	v.errors = DSSXMLUtilsValidateAgainstXSD(v.getUtils(v.xadesPath), v.signatureElement)
	v.errorsSet = true
	return utils.IsCollectionEmpty(v.errors)
}

// ValidationErrors returns validation errors. WARN: Validate must be executed before. Port of
// getValidationErrors().
//
// Panics with the Java message when Validate has not been run yet (IllegalStateException).
func (v *XAdESStructureValidator) ValidationErrors() []string {
	if !v.errorsSet {
		panic("The method XAdESStructureValidator#validate shall be executed before accessing the validation messages!")
	}
	return v.errors
}

// getUtils gets a XAdES implementation of XAdESStructureXSDUtils corresponding to the given
// XAdESPath. Port of the protected getUtils(XAdESPath).
//
// Panics with the Java message when the namespace is not supported for structure validation
// (UnsupportedOperationException).
func (v *XAdESStructureValidator) getUtils(xadesPath definition.XAdESPath) XAdESStructureXSDUtils {
	namespace := xadesPath.Namespace()
	switch namespace {
	case definition.XAdESNamespace_XADES_111:
		return XAdES111XSDUtils()
	case definition.XAdESNamespace_XADES_122:
		return XAdES122XSDUtils()
	case definition.XAdESNamespace_XADES_132, definition.XAdESNamespace_XADES_141:
		return XAdES319132XSDUtils()
	default:
		uri := ""
		if namespace != nil {
			uri = namespace.Uri()
		}
		panic("The namespace '" + uri + "' is not supported for structure validation!")
	}
}
