// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESRevocationRefExtractionUtils.java
// (DSS 6.5.RC1).
package xades

import (
	"math/big"
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// RevocationRefExtractionUtilsCreateOCSPRef extracts an OCSPRef from an ocspRefElement.
// Port of the public static createOCSPRef(XAdESPath, Element).
func RevocationRefExtractionUtilsCreateOCSPRef(xadesPaths definition.XAdESPath, ocspRefElement *xmldom.Node) *spi.OCSPRef {
	digestElement, _ := xmlutils.XPathUtilsGetElement(ocspRefElement, xadesPaths.CurrentDigestAlgAndValue())
	digest := DSSXMLUtilsGetDigestAndValue(digestElement)

	responderId := xadesRevocationRefExtractionUtilsOCSPResponderID(xadesPaths, ocspRefElement)
	if responderId == nil {
		// Upstream logs "Skipped OCSPRef (missing OCSPIdentifier / ResponderID)".
		return nil
	}

	producedAtDate, ok := xadesRevocationRefExtractionUtilsOCSPProducedAtDate(xadesPaths, ocspRefElement)
	if !ok {
		// Upstream logs "Skipped OCSPRef (missing OCSPIdentifier / ProducedAt)".
		return nil
	}

	return spi.NewOCSPRef(digest, producedAtDate, responderId)
}

// xadesRevocationRefExtractionUtilsOCSPProducedAtDate ports the private static
// getOCSPProducedAtDate(XAdESPath, Element). The second return reports whether a valid date was
// found (Java's null return).
func xadesRevocationRefExtractionUtilsOCSPProducedAtDate(xadesPaths definition.XAdESPath, ocspRefElement *xmldom.Node) (time.Time, bool) {
	producedAtEl, _ := xmlutils.XPathUtilsGetElement(ocspRefElement, xadesPaths.CurrentOCSPRefProducedAt())
	if producedAtEl == nil {
		return time.Time{}, false
	}
	producedAtDate := xmlutils.DomUtilsGetDate(producedAtEl.TextContent())
	if producedAtDate.IsZero() {
		return time.Time{}, false
	}
	return producedAtDate, true
}

// xadesRevocationRefExtractionUtilsOCSPResponderID ports the private static
// getOCSPResponderId(XAdESPath, Element).
func xadesRevocationRefExtractionUtilsOCSPResponderID(xadesPaths definition.XAdESPath, ocspRefElement *xmldom.Node) *spi.ResponderId {
	var responderName *model.X500Principal
	var ski []byte

	currentOCSPRefResponderIDByName := xadesPaths.CurrentOCSPRefResponderIDByName()
	currentOCSPRefResponderIDByKey := xadesPaths.CurrentOCSPRefResponderIDByKey()
	if currentOCSPRefResponderIDByName != nil && currentOCSPRefResponderIDByKey != nil {
		responderIdByName, _ := xmlutils.XPathUtilsGetElement(ocspRefElement, currentOCSPRefResponderIDByName)
		if responderIdByName != nil {
			responderName = xadesASN1UtilsGetX500PrincipalOrNull(responderIdByName.TextContent())
		}

		responderIdByKey, _ := xmlutils.XPathUtilsGetElement(ocspRefElement, currentOCSPRefResponderIDByKey)
		if responderIdByKey != nil {
			base64EncodedResponderId := responderIdByKey.TextContent()
			if utils.IsBase64Encoded(base64EncodedResponderId) {
				ski = utils.FromBase64(base64EncodedResponderId)
			}
			// Upstream logs "OCSP ResponderId value is not represented by a base64-encoded string!".
		}
	} else {
		responderIdElement, _ := xmlutils.XPathUtilsGetElement(ocspRefElement, xadesPaths.CurrentOCSPRefResponderID())
		if responderIdElement != nil {
			responderName = xadesASN1UtilsGetX500PrincipalOrNull(responderIdElement.TextContent())
		}
	}

	if responderName != nil || utils.IsArrayNotEmpty(ski) {
		return spi.NewResponderId(responderName, ski)
	}
	return nil
}

// XAdESRevocationRefExtractionUtilsCreateCRLRef extracts a CRLRef from a crlRefElement. Port of
// the public static createCRLRef(XAdESPath, Element).
func RevocationRefExtractionUtilsCreateCRLRef(xadesPaths definition.XAdESPath, crlRefElement *xmldom.Node) *spi.CRLRef {
	digestElement, _ := xmlutils.XPathUtilsGetElement(crlRefElement, xadesPaths.CurrentDigestAlgAndValue())
	digest := DSSXMLUtilsGetDigestAndValue(digestElement)
	if digest.IsEmpty() {
		// Upstream logs "Skipped CRLRef (missing DigestAlgAndValue)".
		return nil
	}

	issuer := xadesRevocationRefExtractionUtilsCRLIssuer(xadesPaths, crlRefElement)
	if issuer == nil {
		// Upstream logs "Skipped CRLRef (missing CRLIdentifier / Issuer)".
		return nil
	}

	issueTime, ok := xadesRevocationRefExtractionUtilsCRLIssueTime(xadesPaths, crlRefElement)
	if !ok {
		// Upstream logs "Skipped CRLRef (missing CRLIdentifier / IssueTime)".
		return nil
	}

	crlNumber := xadesRevocationRefExtractionUtilsCRLNumber(xadesPaths, crlRefElement) // can be nil
	return spi.NewCRLRefWithNumber(digest, issuer, issueTime, crlNumber)
}

// xadesRevocationRefExtractionUtilsCRLIssuer ports the private static getCRLIssuer(XAdESPath,
// Element).
func xadesRevocationRefExtractionUtilsCRLIssuer(xadesPaths definition.XAdESPath, crlRefElement *xmldom.Node) *model.X500Principal {
	issuerEl, _ := xmlutils.XPathUtilsGetElement(crlRefElement, xadesPaths.CurrentCRLRefCRLIdentifierIssuer())
	if issuerEl == nil {
		return nil
	}
	return xadesASN1UtilsGetX500PrincipalOrNull(issuerEl.TextContent())
}

// xadesRevocationRefExtractionUtilsCRLIssueTime ports the private static
// getCRLIssueTime(XAdESPath, Element).
func xadesRevocationRefExtractionUtilsCRLIssueTime(xadesPaths definition.XAdESPath, crlRefElement *xmldom.Node) (time.Time, bool) {
	issueTimeEl, _ := xmlutils.XPathUtilsGetElement(crlRefElement, xadesPaths.CurrentCRLRefCRLIdentifierIssueTime())
	if issueTimeEl == nil {
		return time.Time{}, false
	}
	issueTime := xmlutils.DomUtilsGetDate(issueTimeEl.TextContent())
	if issueTime.IsZero() {
		return time.Time{}, false
	}
	return issueTime, true
}

// xadesRevocationRefExtractionUtilsCRLNumber ports the private static getCRLNumber(XAdESPath,
// Element).
func xadesRevocationRefExtractionUtilsCRLNumber(xadesPaths definition.XAdESPath, crlRefElement *xmldom.Node) *big.Int {
	numberEl, _ := xmlutils.XPathUtilsGetElement(crlRefElement, xadesPaths.CurrentCRLRefCRLIdentifierNumber())
	if numberEl == nil {
		return nil
	}
	number, ok := new(big.Int).SetString(strings.TrimSpace(numberEl.TextContent()), 10)
	if !ok {
		// Upstream logs "Unable to read CRL number".
		return nil
	}
	return number
}
