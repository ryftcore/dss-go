// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESStructureValidator.java
// (DSS 6.5.RC1).
//
// XSD structure validation is a permanent no-op in this port. Upstream's
// getUtils(XAdESPath) resolves an XSDAbstractUtils implementation
// (eu.europa.esig.xades.XAdES111Utils / XAdES122Utils / XAdES319132Utils) that
// bundles the XAdES XSD schemas as JAXB resources and hands it to
// DSSXMLUtils.validateAgainstXSD(XSDAbstractUtils, Source), which runs a
// javax.xml.validation.Validator. There is no bundled-schema or
// XML-Schema-validator package in this port, so dss_xml_utils.go's
// XAdES111XSDUtils/XAdES122XSDUtils/XAdES319132XSDUtils return opaque
// placeholders and DSSXMLUtilsValidateAgainstXSD always reports no errors -
// i.e. every structure is considered valid. getUtils's namespace dispatch
// below, including its panic on an unrecognized namespace, is still exercised
// faithfully; only the schema check itself is skipped.
//
// XAdESStructureXSDUtils stands in for
// eu.europa.esig.dss.jaxb.common.XSDAbstractUtils: an opaque marker type, since
// nothing here calls a method on it.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
)

// XAdESStructureXSDUtils stands in for eu.europa.esig.dss.jaxb.common.XSDAbstractUtils; see the
// file header's note on XSD validation.
type XAdESStructureXSDUtils any

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
	case definition.XAdESNamespaceXAdES111:
		return XAdES111XSDUtils()
	case definition.XAdESNamespaceXAdES122:
		return XAdES122XSDUtils()
	case definition.XAdESNamespaceXAdES132, definition.XAdESNamespaceXAdES141:
		return XAdES319132XSDUtils()
	default:
		uri := ""
		if namespace != nil {
			uri = namespace.Uri()
		}
		panic("The namespace '" + uri + "' is not supported for structure validation!")
	}
}
