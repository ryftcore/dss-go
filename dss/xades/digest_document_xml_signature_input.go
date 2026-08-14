// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/DigestDocumentXMLSignatureInput.java
// (DSS 6.5.RC1).
//
// See dss_document_xml_signature_input.go's header for the internal/xmldsig.Data replacement
// context this class shares with its parent.
package xades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// DigestDocumentXMLSignatureInput is an XMLSignatureInput definition from a DigestDocument.
// Port of the class DigestDocumentXMLSignatureInput, extending DSSDocumentXMLSignatureInput.
type DigestDocumentXMLSignatureInput struct {
	DSSDocumentXMLSignatureInput
}

// NewDigestDocumentXMLSignatureInput is the constructor for an XMLSignatureInput from a
// DigestDocument. Port of the public DigestDocumentXMLSignatureInput(DigestDocument,
// DigestAlgorithm) constructor.
func NewDigestDocumentXMLSignatureInput(document *model.DigestDocument, digestAlgorithm enumerations.DigestAlgorithm) (*DigestDocumentXMLSignatureInput, error) {
	base, err := newDSSDocumentXMLSignatureInputWithDigest(document, digestAlgorithm)
	if err != nil {
		return nil, err
	}
	return &DigestDocumentXMLSignatureInput{DSSDocumentXMLSignatureInput: *base}, nil
}
