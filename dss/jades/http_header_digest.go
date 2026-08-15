// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/HTTPHeaderDigest.java (DSS 6.5.RC1).
//
// An HTTP message body, which 'Digest' representation is being signed with the 'sigD'
// HTTP_HEADERS mechanism.
package jades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// HTTPHeaderDigest extends HTTPHeader (embeds it, since Go has no class inheritance) to carry the
// original signing message-body document alongside the computed 'Digest' header value.
type HTTPHeaderDigest struct {
	HTTPHeader

	// messageBodyDocument is the message body content document.
	messageBodyDocument model.DSSDocument
}

var _ model.DSSDocument = (*HTTPHeaderDigest)(nil)

// NewHTTPHeaderDigest is the default constructor, taking the signing message body document
// content and the DigestAlgorithm to use to compute the document's Digest element. Panics with
// the Java message when document or digestAlgorithm is nil/empty (Objects.requireNonNull
// upstream), or when the DigestAlgorithm has no RFC 5843 'sigD' HTTP_HEADERS mapping
// (IllegalArgumentException upstream).
func NewHTTPHeaderDigest(messageBodyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *HTTPHeaderDigest {
	value := httpHeaderDigestBuildInstanceDigestValue(messageBodyDocument, digestAlgorithm)
	return &HTTPHeaderDigest{
		HTTPHeader:          *NewHTTPHeader(DSSJsonUtilsHTTPHeaderDigest, value),
		messageBodyDocument: messageBodyDocument,
	}
}

// httpHeaderDigestBuildInstanceDigestValue ports the private static #buildInstanceDigestValue.
func httpHeaderDigestBuildInstanceDigestValue(document model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) string {
	if document == nil {
		panic("DSSDocument shall be provided!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm shall be provided!")
	}

	jwsHTTPHeaderAlgo := digestAlgorithm.HttpHeaderAlgo()
	if jwsHTTPHeaderAlgo == "" {
		panic(fmt.Sprintf("The DigestAlgorithm '%s' is not supported for 'sigD' HTTP_HEADERS mechanism. "+
			"See RFC 5843 for more information.", digestAlgorithm))
	}

	// RFC 3230 "Instance Digests in HTTP"
	//
	// 4.2 Instance digests
	//
	// An instance digest is the representation of the output of a digest algorithm, together
	// with an indication of the algorithm used (and any parameters).
	//
	// instance-digest = digest-algorithm "=" <encoded digest output>

	digest, err := document.DigestValue(digestAlgorithm)
	if err != nil {
		panic(err)
	}
	base64EncodedDigest := utils.ToBase64(digest)

	return jwsHTTPHeaderAlgo + "=" + base64EncodedDigest
}

// MessageBodyDocument returns the original HTTP Message Body Document. Port of
// #getMessageBodyDocument.
func (h *HTTPHeaderDigest) MessageBodyDocument() model.DSSDocument {
	return h.messageBodyDocument
}
