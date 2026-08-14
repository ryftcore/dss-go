// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XPointerSignatureScope.java (DSS 6.5.RC1).
package xades

import (
	"strings"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	spiscope "github.com/utain/esig/dss/spi/validation/scope"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XPointerSignatureScope is the XPointer signature scope. Port of the class
// XPointerSignatureScope, extending spiscope.SignatureScopeWithTransformationsBase.
type XPointerSignatureScope struct {
	spiscope.SignatureScopeWithTransformationsBase

	// uri is the XPointer query.
	uri string
}

// newXPointerSignatureScope is the port of the protected XPointerSignatureScope(String,
// DSSDocument, List<String>) constructor.
func newXPointerSignatureScope(uri string, document model.DSSDocument, transformations []string) *XPointerSignatureScope {
	return &XPointerSignatureScope{
		SignatureScopeWithTransformationsBase: spiscope.NewSignatureScopeWithTransformationsBaseWithName(
			xPointerSignatureScopeDocumentNameFromXPointer(uri), document, transformations),
		uri: uri,
	}
}

// xPointerSignatureScopeDocumentNameFromXPointer ports the private static
// getDocumentNameFromXPointer(String).
func xPointerSignatureScopeDocumentNameFromXPointer(uri string) string {
	if xmlutils.DomUtilsIsRootXPointer(uri) {
		return "Full XML file"
	}
	xpointerID, _ := xmlutils.DomUtilsGetXPointerId(uri)
	return xpointerID
}

// Description returns the XPointerSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *XPointerSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	var sb strings.Builder
	sb.WriteString("XPointer query to ")
	if xmlutils.DomUtilsIsRootXPointer(s.uri) {
		sb.WriteString("root XML element")
	} else {
		sb.WriteString("element with Id '")
		sb.WriteString(s.DocumentName())
		sb.WriteString("'")
	}
	return s.AddTransformationIfNeeded(sb.String())
}

// Type returns the type of the signature scope. Port of getType().
func (s *XPointerSignatureScope) Type() enumerations.SignatureScopeType {
	if xmlutils.DomUtilsIsRootXPointer(s.uri) {
		return enumerations.SignatureScopeType_FULL
	}
	return enumerations.SignatureScopeType_PARTIAL
}

// String ports toString().
func (s *XPointerSignatureScope) String() string {
	return "XPointerSignatureScope{uri='" + s.uri + "'} " + s.SignatureScopeWithTransformationsBase.String()
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*XPointerSignatureScope)(nil)
