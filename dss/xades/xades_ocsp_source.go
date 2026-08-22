// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESOCSPSource.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	"github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESOCSPSource retrieves OCSP values from an XAdES (XL/LT) signature. Port of the class
// OCSPSource, extending spi.OfflineOCSPSourceBase.
type OCSPSource struct {
	spi.OfflineOCSPSourceBase

	// signatureElement is the current signature element.
	signatureElement *xmldom.Node

	// xadesPaths are the XAdES XPaths.
	xadesPaths definition.XAdESPath
}

// NewXAdESOCSPSource is the default constructor for XAdESOCSPSource. Port of the constructor
// OCSPSource(Element, XAdESPath).
//
// Panics with the Java messages when signatureElement or xadesPaths is missing
// (Objects.requireNonNull).
func NewOCSPSource(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *OCSPSource {
	if signatureElement == nil {
		panic("Signature element cannot be null")
	}
	if xadesPaths == nil {
		panic("XAdESPaths cannot be null")
	}

	s := &OCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
		signatureElement:      signatureElement,
		xadesPaths:            xadesPaths,
	}
	s.InitOfflineRevocationSource(s)
	s.appendContainedOCSPResponses()
	return s
}

// appendContainedOCSPResponses ports the private appendContainedOCSPResponses().
func (s *OCSPSource) appendContainedOCSPResponses() {
	// values
	s.collectValues(s.xadesPaths.RevocationValuesPath(), enumerations.RevocationOriginRevocationValues)
	s.collectValues(s.xadesPaths.AttributeRevocationValuesPath(), enumerations.RevocationOriginAttributeRevocationValues)
	s.collectValues(s.xadesPaths.TimeStampValidationDataRevocationValuesPath(), enumerations.RevocationOriginTimestampValidationData)
	s.collectValues(s.xadesPaths.AnyValidationDataRevocationValuesPath(), enumerations.RevocationOriginAnyValidationData)

	// references
	s.collectRefs(s.xadesPaths.CompleteRevocationRefsPath(), enumerations.RevocationRefOriginCompleteRevocationRefs)
	s.collectRefs(s.xadesPaths.AttributeRevocationRefsPath(), enumerations.RevocationRefOriginAttributeRevocationRefs)
}

// collectValues ports the private collectValues(XPathQuery, RevocationOrigin).
func (s *OCSPSource) collectValues(revocationValuesPath common.XPathQuery, origin enumerations.RevocationOrigin) {
	if revocationValuesPath == nil {
		return
	}

	revocationValuesNodeList, err := utils.XPathUtilsGetNodeList(s.signatureElement, revocationValuesPath)
	if err != nil {
		return
	}
	for _, revocationValuesElement := range revocationValuesNodeList {
		ocspValueNodes, err := utils.XPathUtilsGetNodeList(revocationValuesElement, s.xadesPaths.CurrentOCSPValuesChildren())
		if err != nil {
			continue
		}
		for _, ocspValueEl := range ocspValueNodes {
			if ocspValueEl != nil {
				s.convertAndAppend(ocspValueEl.TextContent(), origin)
			}
		}
	}
}

// collectRefs ports the private collectRefs(XPathQuery, RevocationRefOrigin).
func (s *OCSPSource) collectRefs(revocationRefsPath common.XPathQuery, revocationRefOrigin enumerations.RevocationRefOrigin) {
	if revocationRefsPath == nil {
		return
	}

	revocationRefsNodeList, err := utils.XPathUtilsGetNodeList(s.signatureElement, revocationRefsPath)
	if err != nil {
		return
	}
	for _, revocationRefsElement := range revocationRefsNodeList {
		ocspRefNodes, err := utils.XPathUtilsGetNodeList(revocationRefsElement, s.xadesPaths.CurrentOCSPRefsChildren())
		if err != nil {
			continue
		}
		for _, ocspRefElement := range ocspRefNodes {
			if ocspRefElement != nil {
				ocspRef := RevocationRefExtractionUtilsCreateOCSPRef(s.xadesPaths, ocspRefElement)
				if ocspRef != nil {
					s.AddRevocationReference(ocspRef, revocationRefOrigin)
				}
			}
		}
	}
}

// convertAndAppend ports the private convertAndAppend(String, RevocationOrigin).
//
// Upstream logs "Cannot retrieve OCSP response from '{}' : {}" and swallows the error; slf4j
// logging is dropped per PORTING.md, so a load failure is simply skipped here too.
func (s *OCSPSource) convertAndAppend(ocspValue string, origin enumerations.RevocationOrigin) {
	basicOCSPResp, err := spi.DSSRevocationUtilsLoadOCSPBase64Encoded(ocspValue)
	if err != nil {
		return
	}
	binary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
	if err != nil {
		return
	}
	s.AddBinary(binary, origin)
}
