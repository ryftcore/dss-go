// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESCRLSource.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESCRLSource retrieves CRL values from an XAdES (-XL) signature. Port of the class
// XAdESCRLSource, extending spi.OfflineCRLSourceBase.
type XAdESCRLSource struct {
	spi.OfflineCRLSourceBase

	// signatureElement is the current signature element.
	signatureElement *xmldom.Node

	// xadesPaths is the XAdES XPaths.
	xadesPaths definition.XAdESPath
}

// NewXAdESCRLSource is the port of the constructor XAdESCRLSource(Element, XAdESPath). All CRL
// values and references are extracted during instantiation.
//
// Panics with the Java messages when signatureElement or xadesPaths is missing
// (Objects.requireNonNull).
func NewXAdESCRLSource(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *XAdESCRLSource {
	if signatureElement == nil {
		panic("Signature element cannot be null")
	}
	if xadesPaths == nil {
		panic("XAdESPaths cannot be null")
	}

	s := &XAdESCRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		signatureElement:     signatureElement,
		xadesPaths:           xadesPaths,
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches OfflineCRLSourceBase.RevocationTokens; this class does
	// not override it.
	s.InitOfflineRevocationSource(s)

	// values
	s.collectValues(xadesPaths.RevocationValuesPath(), enumerations.RevocationOrigin_REVOCATION_VALUES)
	s.collectValues(xadesPaths.AttributeRevocationValuesPath(), enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES)
	s.collectValues(xadesPaths.TimeStampValidationDataRevocationValuesPath(), enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA)
	s.collectValues(xadesPaths.AnyValidationDataRevocationValuesPath(), enumerations.RevocationOrigin_ANY_VALIDATION_DATA)

	// references
	s.collectRefs(xadesPaths.CompleteRevocationRefsPath(), enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)
	s.collectRefs(xadesPaths.AttributeRevocationRefsPath(), enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)

	return s
}

// collectValues ports the private collectValues(XPathQuery, RevocationOrigin).
func (s *XAdESCRLSource) collectValues(revocationValuesPath common.XPathQuery, revocationOrigin enumerations.RevocationOrigin) {
	if revocationValuesPath == nil {
		return
	}

	revocationValuesNodeList, err := xmlutils.XPathUtilsGetNodeList(s.signatureElement, revocationValuesPath)
	if err != nil {
		return
	}
	for _, revocationValuesElement := range revocationValuesNodeList {
		crlValueNodes, err := xmlutils.XPathUtilsGetNodeList(revocationValuesElement, s.xadesPaths.CurrentCRLValuesChildren())
		if err != nil {
			continue
		}
		for _, crlValueEl := range crlValueNodes {
			base64EncodedCRL := crlValueEl.TextContent()
			crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(utils.FromBase64(base64EncodedCRL))
			if err != nil {
				// Upstream logs "Unable to build CRLBinary from an obtained element with origin '{}'.".
				continue
			}
			s.AddBinary(crlBinary, revocationOrigin)
		}
	}
}

// collectRefs ports the private collectRefs(XPathQuery, RevocationRefOrigin).
func (s *XAdESCRLSource) collectRefs(revocationRefsPath common.XPathQuery, revocationRefOrigin enumerations.RevocationRefOrigin) {
	if revocationRefsPath == nil {
		return
	}

	revocationRefsNodeList, err := xmlutils.XPathUtilsGetNodeList(s.signatureElement, revocationRefsPath)
	if err != nil {
		return
	}
	for _, revocationRefsElement := range revocationRefsNodeList {
		crlRefNodes, err := xmlutils.XPathUtilsGetNodeList(revocationRefsElement, s.xadesPaths.CurrentCRLRefsChildren())
		if err != nil {
			continue
		}
		for _, crlRefNode := range crlRefNodes {
			crlRef := XAdESRevocationRefExtractionUtilsCreateCRLRef(s.xadesPaths, crlRefNode)
			if crlRef != nil {
				s.AddRevocationReference(crlRef, revocationRefOrigin)
			}
		}
	}
}
