// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESOCSPSource.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: XAdESRevocationRefExtractionUtils (Java eu.europa.esig.dss.xades.
// validation.XAdESRevocationRefExtractionUtils, same "validation" SCC, hence the same Go
// package) is owned by a sibling chunk not in this manifest. The single entry point this file
// calls, inferred from its Java source and every call site here:
//
//	func XAdESRevocationRefExtractionUtilsCreateOCSPRef(xadesPaths definition.XAdESPath, ocspRefElement *xmldom.Node) *spi.OCSPRef
//	  - port of the static createOCSPRef(XAdESPath, Element); nil on a malformed/incomplete
//	    OCSPRef element (Java returns null after logging a warning).
package xades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	"github.com/utain/esig/dss/xml/utils"
)

// XAdESOCSPSource retrieves OCSP values from an XAdES (XL/LT) signature. Port of the class
// XAdESOCSPSource, extending spi.OfflineOCSPSourceBase.
type XAdESOCSPSource struct {
	spi.OfflineOCSPSourceBase

	// signatureElement is the current signature element.
	signatureElement *xmldom.Node

	// xadesPaths are the XAdES XPaths.
	xadesPaths definition.XAdESPath
}

// NewXAdESOCSPSource is the default constructor for XAdESOCSPSource. Port of the constructor
// XAdESOCSPSource(Element, XAdESPath).
//
// Panics with the Java messages when signatureElement or xadesPaths is missing
// (Objects.requireNonNull).
func NewXAdESOCSPSource(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *XAdESOCSPSource {
	if signatureElement == nil {
		panic("Signature element cannot be null")
	}
	if xadesPaths == nil {
		panic("XAdESPaths cannot be null")
	}

	s := &XAdESOCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
		signatureElement:      signatureElement,
		xadesPaths:            xadesPaths,
	}
	s.InitOfflineRevocationSource(s)
	s.appendContainedOCSPResponses()
	return s
}

// appendContainedOCSPResponses ports the private appendContainedOCSPResponses().
func (s *XAdESOCSPSource) appendContainedOCSPResponses() {
	// values
	s.collectValues(s.xadesPaths.RevocationValuesPath(), enumerations.RevocationOrigin_REVOCATION_VALUES)
	s.collectValues(s.xadesPaths.AttributeRevocationValuesPath(), enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES)
	s.collectValues(s.xadesPaths.TimeStampValidationDataRevocationValuesPath(), enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA)
	s.collectValues(s.xadesPaths.AnyValidationDataRevocationValuesPath(), enumerations.RevocationOrigin_ANY_VALIDATION_DATA)

	// references
	s.collectRefs(s.xadesPaths.CompleteRevocationRefsPath(), enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)
	s.collectRefs(s.xadesPaths.AttributeRevocationRefsPath(), enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)
}

// collectValues ports the private collectValues(XPathQuery, RevocationOrigin).
func (s *XAdESOCSPSource) collectValues(revocationValuesPath common.XPathQuery, origin enumerations.RevocationOrigin) {
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
func (s *XAdESOCSPSource) collectRefs(revocationRefsPath common.XPathQuery, revocationRefOrigin enumerations.RevocationRefOrigin) {
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
				ocspRef := XAdESRevocationRefExtractionUtilsCreateOCSPRef(s.xadesPaths, ocspRefElement)
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
func (s *XAdESOCSPSource) convertAndAppend(ocspValue string, origin enumerations.RevocationOrigin) {
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
