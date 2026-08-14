// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/tsl/XAdESTrustedListUtils.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"strconv"

	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// TLV5Identifier is the identifier used for a TL version 5. Port of TL_V5_IDENTIFIER.
const TLV5Identifier = 5

// TLV6Identifier is the identifier used for a TL version 6. Port of TL_V6_IDENTIFIER.
const TLV6Identifier = 6

func init() {
	xmlutils.XPathUtilsRegisterNamespace(definition.TrustedListNamespace_NS)
}

// XAdESTrustedListUtilsGetTSLVersionIdentifier parses the XML Trusted List and returns the TL
// version, when present. Port of the static getTSLVersionIdentifier(Document).
func XAdESTrustedListUtilsGetTSLVersionIdentifier(documentDom *xmldom.Node) (*int, error) {
	tlVersionIdentifierElement, err := xmlutils.XPathUtilsGetElement(documentDom.DocumentElement(), definition.TrustedListPath_TSL_VERSION_IDENTIFIER_PATH)
	if err != nil {
		return nil, err
	}
	if tlVersionIdentifierElement == nil {
		return nil, nil
	}
	tlVersionIdentifier := tlVersionIdentifierElement.TextContent()
	if !utils.IsStringDigits(tlVersionIdentifier) {
		return nil, nil
	}
	value, err := strconv.Atoi(tlVersionIdentifier)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// XAdESTrustedListUtilsValidateUnsignedTrustedList helps to determine whether the chosen
// signature parameters builder is applicable to the given DSSDocument. It verifies whether the
// provided document representing the XML Trusted List is conformant to the definition and the
// target version. Returns a list of errors if problems have been found during the validation.
// Returns an empty list in case of a valid XML Trusted List conformant to the specified
// version. Port of the static validateUnsignedTrustedList(DSSDocument, int).
//
// Panics when tlDocument is nil (Java's Objects.requireNonNull).
func XAdESTrustedListUtilsValidateUnsignedTrustedList(tlDocument model.DSSDocument, targetTLVersion int) ([]string, error) {
	if tlDocument == nil {
		panic("XML Trusted List cannot be null!")
	}

	if !xmlutils.DomUtilsIsDOM(tlDocument) {
		return []string{"The document is not XML!"}, nil
	}
	documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(tlDocument)
	if err != nil {
		return nil, err
	}
	return XAdESTrustedListUtilsValidateUnsignedTrustedListDocument(documentDom, targetTLVersion)
}

// XAdESTrustedListUtilsValidateUnsignedTrustedListDocument helps to determine whether the
// chosen signature parameters builder is applicable to the given Document. It verifies whether
// the provided document representing the XML Trusted List is conformant to the definition and
// the target version. Returns a list of errors if problems have been found during the
// validation. Returns an empty list in case of a valid XML Trusted List conformant to the
// specified version. Port of the static validateUnsignedTrustedList(Document, int).
//
// Panics when documentDom is nil (Java's Objects.requireNonNull).
func XAdESTrustedListUtilsValidateUnsignedTrustedListDocument(documentDom *xmldom.Node, targetTLVersion int) ([]string, error) {
	if documentDom == nil {
		panic("XML Trusted List cannot be null!")
	}
	var errorMessages []string

	tlVersionIdentifier, err := XAdESTrustedListUtilsGetTSLVersionIdentifier(documentDom)
	if err != nil {
		return nil, err
	}
	if tlVersionIdentifier == nil {
		errorMessages = append(errorMessages, "No TSLVersionIdentifier has been found!")

	} else if *tlVersionIdentifier != targetTLVersion {
		errorMessages = append(errorMessages, fmt.Sprintf("TSL Version '%d' found in the XML Trusted List "+
			"does not correspond to the target version defined by the builder '%d'! "+
			"Please modify the document or change to the appropriate builder.", *tlVersionIdentifier, targetTLVersion))

	} else {
		xsdErrors := NewTLStructureVerifier().
			SetSigningMode(true).
			SetAcceptedTLVersions([]int{targetTLVersion}).
			ValidateDocument(documentDom, &targetTLVersion)
		if utils.IsCollectionNotEmpty(xsdErrors) {
			errorMessages = append(errorMessages, xsdErrors...)
		}
	}

	return errorMessages, nil
}
