// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XmlElementSignatureScope.java (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
)

// XmlElementSignatureScope is the XML element signature scope. Port of the class
// XmlElementSignatureScope, extending spiscope.SignatureScopeWithTransformationsBase.
type XmlElementSignatureScope struct {
	spiscope.SignatureScopeWithTransformationsBase
}

// newXmlElementSignatureScope is the port of the protected XmlElementSignatureScope(String,
// DSSDocument, List<String>) constructor.
func newXmlElementSignatureScope(xmlId string, document model.DSSDocument, transformations []string) *XmlElementSignatureScope {
	return &XmlElementSignatureScope{
		SignatureScopeWithTransformationsBase: spiscope.NewSignatureScopeWithTransformationsBaseWithName(xmlId, document, transformations),
	}
}

// Description returns the XmlElementSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *XmlElementSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	description := fmt.Sprintf("The XML element with ID '%s'", s.DocumentName())
	return s.AddTransformationIfNeeded(description)
}

// Type returns the type of the signature scope. Port of getType().
func (s *XmlElementSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_PARTIAL
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*XmlElementSignatureScope)(nil)
