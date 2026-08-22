// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESCRLSource.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESCRLSource retrieves CRL values from an XAdES (-XL) signature. Port of the class
// CRLSource, extending spi.OfflineCRLSourceBase.
type CRLSource struct {
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
func NewCRLSource(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *CRLSource {
	if signatureElement == nil {
		panic("Signature element cannot be null")
	}
	if xadesPaths == nil {
		panic("XAdESPaths cannot be null")
	}

	s := &CRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		signatureElement:     signatureElement,
		xadesPaths:           xadesPaths,
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches OfflineCRLSourceBase.RevocationTokens; this class does
	// not override it.
	s.InitOfflineRevocationSource(s)

	// values
	s.collectValues(xadesPaths.RevocationValuesPath(), enumerations.RevocationOriginRevocationValues)
	s.collectValues(xadesPaths.AttributeRevocationValuesPath(), enumerations.RevocationOriginAttributeRevocationValues)
	s.collectValues(xadesPaths.TimeStampValidationDataRevocationValuesPath(), enumerations.RevocationOriginTimestampValidationData)
	s.collectValues(xadesPaths.AnyValidationDataRevocationValuesPath(), enumerations.RevocationOriginAnyValidationData)

	// references
	s.collectRefs(xadesPaths.CompleteRevocationRefsPath(), enumerations.RevocationRefOriginCompleteRevocationRefs)
	s.collectRefs(xadesPaths.AttributeRevocationRefsPath(), enumerations.RevocationRefOriginAttributeRevocationRefs)

	return s
}

// collectValues ports the private collectValues(XPathQuery, RevocationOrigin).
func (s *CRLSource) collectValues(revocationValuesPath common.XPathQuery, revocationOrigin enumerations.RevocationOrigin) {
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
func (s *CRLSource) collectRefs(revocationRefsPath common.XPathQuery, revocationRefOrigin enumerations.RevocationRefOrigin) {
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
			crlRef := RevocationRefExtractionUtilsCreateCRLRef(s.xadesPaths, crlRefNode)
			if crlRef != nil {
				s.AddRevocationReference(crlRef, revocationRefOrigin)
			}
		}
	}
}
