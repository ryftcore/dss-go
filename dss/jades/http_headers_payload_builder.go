// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/HttpHeadersPayloadBuilder.java (DSS 6.5.RC1).
//
// The payload this builder produces is signed (or time-stamped) verbatim, so its bytes are a hard
// contract: header names lowercased, ": " between name and value, values of repeated headers
// joined with ", " in encounter order, '\n' between fields and none after the last, and - for a
// timestamp - the 'Digest' field replaced by the raw message-body octets.
//
// # HTTPHeaderDigest, and why the working list is not []*HTTPHeader
//
// Java collects the documents into a List<HTTPHeader> and later asks `instanceof
// HTTPHeaderDigest` of its elements: the subclass survives the widening. Go's HTTPHeaderDigest
// embeds HTTPHeader, and narrowing a *HTTPHeaderDigest to its embedded *HTTPHeader would DROP the
// message-body document the timestamp path needs. The working list therefore carries the local
// httpHeaderDocument interface, which both types satisfy through the promoted Value/SetValue, so
// a later type assertion recovers the digest header exactly where Java's instanceof does.
//
// Java's ByteArrayOutputStream becomes a bytes.Buffer, so the IOException branch (and the
// DSSException wrapping it) is unreachable in Go; the error return is kept because reading the
// message-body document can fail. Every IllegalArgumentException becomes a returned error.
package jades

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// httpHeaderDocument is the Go stand-in for Java's `HTTPHeader` static type in this file: the
// common surface of HTTPHeader and its HTTPHeaderDigest subclass. See the file header.
type httpHeaderDocument interface {
	model.DSSDocument

	// Value returns the HTTP header's value. Port of HTTPHeader#getValue.
	Value() string

	// SetValue sets the HTTP header's value. Port of HTTPHeader#setValue.
	SetValue(value string)
}

// HttpHeadersPayloadBuilder builds payload binaries from HTTPHeaderDocuments for the 'sigD'
// HttpHeaders mechanism.
type HttpHeadersPayloadBuilder struct {
	// detachedContents holds the provided detached documents.
	detachedContents []model.DSSDocument

	// isTimestamp tells whether the payload shall be computed for a timestamp (which defines a
	// different processing).
	isTimestamp bool
}

// NewHttpHeadersPayloadBuilder is the default constructor.
// Port of HttpHeadersPayloadBuilder(List<DSSDocument>, boolean).
func NewHttpHeadersPayloadBuilder(detachedContents []model.DSSDocument,
	isTimestamp bool) *HttpHeadersPayloadBuilder {
	return &HttpHeadersPayloadBuilder{detachedContents: detachedContents, isTimestamp: isTimestamp}
}

// Build builds the payload from the HTTPHeaderDocuments. Port of #build.
func (b *HttpHeadersPayloadBuilder) Build() ([]byte, error) {
	if err := b.assertHttpHeadersConfigurationIsValid(); err != nil {
		return nil, err
	}

	httpHeaderDocuments, err := httpHeadersPayloadBuilderToHTTPHeaders(b.detachedContents)
	if err != nil {
		return nil, err
	}

	/*
	 * Signing HTTP Messages draft-cavage-http-signatures-10
	 *
	 * To include the HTTP request target in the signature calculation, use the
	 * special `(request-target)` header field name.
	 *
	 * 1. If the header field name is `(request-target)` then generate the header
	 * field value by concatenating the lowercased :method, an ASCII space, and the
	 * :path pseudo-headers (as specified in HTTP/2, Section 8.1.2.3 [7]). Note: For
	 * the avoidance of doubt, lowercasing only applies to the :method pseudo-header
	 * and not to the :path pseudo-header.
	 *
	 * 2. Create the header field string by concatenating the lowercased header
	 * field name followed with an ASCII colon `:`, an ASCII space ` `, and the
	 * header field value. Leading and trailing optional whitespace (OWS) in the
	 * header field value MUST be omitted (as specified in RFC7230 [RFC7230],
	 * Section 3.2.4 [8]). If there are multiple instances of the same header field,
	 * all header field values associated with the header field MUST be
	 * concatenated, separated by a ASCII comma and an ASCII space `, `, and used in
	 * the order in which they will appear in the transmitted HTTP message. Any
	 * other modification to the header field value MUST NOT be made.
	 *
	 * 3. If value is not the last value then append an ASCII newline `\n`.
	 */

	concatenatedHttpFields := make([]httpHeaderDocument, 0)

	for _, httpHeader := range httpHeaderDocuments {
		headerName := utils.Trim(httpHeader.Name())
		headerValue := utils.Trim(httpHeader.Value())

		concatenatedHttpHeader := httpHeadersPayloadBuilderHTTPHeaderWithName(concatenatedHttpFields, headerName)

		if DSSJsonUtilsHTTPHeaderDigest == headerName && b.isTimestamp {
			if concatenatedHttpHeader != nil {
				return nil, fmt.Errorf("Only one HTTPHeader with the name '%s' is allowed!",
					DSSJsonUtilsHTTPHeaderDigest)
			}
			if _, ok := httpHeader.(*HTTPHeaderDigest); !ok {
				return nil, errors.New("Unable to compute message-imprint for an Archive Timestamp! " +
					"'Digest' header must be an instance of HTTPHeaderDigest class.")
			}
			concatenatedHttpFields = append(concatenatedHttpFields, httpHeader)

		} else if concatenatedHttpHeader != nil {
			var stringBuilder bytes.Buffer
			stringBuilder.WriteString(concatenatedHttpHeader.Value())
			stringBuilder.WriteString(", ")
			stringBuilder.WriteString(headerValue)
			headerValue = stringBuilder.String()

			concatenatedHttpHeader.SetValue(headerValue)

		} else {
			concatenatedHttpFields = append(concatenatedHttpFields, NewHTTPHeader(headerName, headerValue))
		}
	}

	var baos bytes.Buffer
	for index, header := range concatenatedHttpFields {
		if DSSJsonUtilsHTTPHeaderDigest == header.Name() && b.isTimestamp {
			httpHeaderDigest, ok := header.(*HTTPHeaderDigest)
			if !ok {
				// Java's (HTTPHeaderDigest) cast; the loop above already guarantees it.
				return nil, errors.New("Unable to compute message-imprint for a Timestamp! " +
					"'Digest' header must be an instance of HTTPHeaderDigest class.")
			}
			messageBodyDocument := httpHeaderDigest.MessageBodyDocument()
			binaries, err := spi.DSSUtilsToByteArrayOfDocument(messageBodyDocument)
			if err != nil {
				return nil, err
			}
			baos.Write(binaries)
		} else {
			var stringBuilder bytes.Buffer
			stringBuilder.WriteString(utils.LowerCase(header.Name()))
			stringBuilder.WriteString(":")
			stringBuilder.WriteString(" ")
			stringBuilder.WriteString(header.Value())
			baos.WriteString(stringBuilder.String())
		}
		if index != len(concatenatedHttpFields)-1 {
			baos.WriteString("\n")
		}
	}
	return baos.Bytes(), nil
}

// httpHeadersPayloadBuilderHTTPHeaderWithName ports the private getHTTPHeaderWithName.
func httpHeadersPayloadBuilderHTTPHeaderWithName(httpHeaders []httpHeaderDocument,
	name string) httpHeaderDocument {
	for _, httpHeader := range httpHeaders {
		if name == httpHeader.Name() {
			return httpHeader
		}
	}
	return nil
}

// httpHeadersPayloadBuilderToHTTPHeaders casts a list of DSSDocuments to a list of HTTPHeaders,
// returning an error if a document of another class is found.
// Port of the private toHTTPHeaders.
func httpHeadersPayloadBuilderToHTTPHeaders(dssDocuments []model.DSSDocument) ([]httpHeaderDocument, error) {
	httpHeaderDocuments := make([]httpHeaderDocument, 0, len(dssDocuments))
	for _, document := range dssDocuments {
		httpHeaderDocument, ok := httpHeadersPayloadBuilderAsHTTPHeader(document)
		if !ok {
			return nil, fmt.Errorf("The document with name '%s' is not of type HTTPHeader!", document.Name())
		}
		httpHeaderDocuments = append(httpHeaderDocuments, httpHeaderDocument)
	}
	return httpHeaderDocuments, nil
}

// httpHeadersPayloadBuilderAsHTTPHeader is Go's `instanceof HTTPHeader`: it matches the class and
// its HTTPHeaderDigest subclass, and nothing else.
func httpHeadersPayloadBuilderAsHTTPHeader(document model.DSSDocument) (httpHeaderDocument, bool) {
	switch header := document.(type) {
	case *HTTPHeaderDigest:
		return header, true
	case *HTTPHeader:
		return header, true
	default:
		return nil, false
	}
}

// assertHttpHeadersConfigurationIsValid checks that a valid detached content is provided for the
// "HTTPHeaders" "sigD" mechanism. Port of the private assertHttpHeadersConfigurationIsValid.
func (b *HttpHeadersPayloadBuilder) assertHttpHeadersConfigurationIsValid() error {
	if !utils.IsCollectionNotEmpty(b.detachedContents) {
		return errors.New("Unable to compute HTTPHeaders payload! The list of detached documents is empty.")
	}
	digestDocumentFound := false
	for _, document := range b.detachedContents {
		digestHTTPHeaderDocument, err := b.checkIfDigestHTTPHeaderDocument(document)
		if err != nil {
			return err
		}
		if digestHTTPHeaderDocument {
			if digestDocumentFound {
				return errors.New("Only one 'Digest' header or HTTPHeaderDigest object is allowed!")
			}
			digestDocumentFound = true
		}
	}
	return nil
}

// checkIfDigestHTTPHeaderDocument ports the private checkIfDigestHTTPHeaderDocument.
func (b *HttpHeadersPayloadBuilder) checkIfDigestHTTPHeaderDocument(
	document model.DSSDocument) (bool, error) {
	if _, ok := httpHeadersPayloadBuilderAsHTTPHeader(document); !ok {
		return false, errors.New("The documents to sign must have " +
			"a type of HTTPHeader for 'sigD' HttpHeaders mechanism!")
	}
	if DSSJsonUtilsHTTPHeaderDigest == document.Name() {
		if _, isDigest := document.(*HTTPHeaderDigest); !isDigest && b.isTimestamp {
			return false, errors.New("Unable to compute message-imprint for a Timestamp! " +
				"'Digest' header must be an instance of HTTPHeaderDigest class.")
		}
		return true, nil
	}
	return false, nil
}
