// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/HTTPHeader.java (DSS 6.5.RC1).
//
// The class represents an HTTP Header to be signed, per ETSI TS 119 182-1 "5.2.8.2 Mechanism
// HttpHeaders". It shall be used only for JAdES detached SigDMechanism_HTTP_HEADERS mechanism.
package jades

import (
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// HTTPHeader implements model.DSSDocument to model an HTTP Header entry (name/value pair) as a
// pseudo-document for signing.
type HTTPHeader struct {
	// name is the HTTP Header's name.
	name string
	// value is the HTTP Header's value.
	value string
}

var _ model.DSSDocument = (*HTTPHeader)(nil)

// NewHTTPHeader is the default constructor.
func NewHTTPHeader(name, value string) *HTTPHeader {
	return &HTTPHeader{name: name, value: value}
}

// Name returns a String name (key) of the HTTP Header. Port of #getName.
func (h *HTTPHeader) Name() string {
	return h.name
}

// SetName panics: not supported for HTTPHeader. Port of #setName.
func (h *HTTPHeader) SetName(name string) {
	panic("The setName(name) method is not supported for HTTPHeaderDocument.")
}

// Value returns a String value of the HTTP Header. Port of #getValue.
func (h *HTTPHeader) Value() string {
	return h.value
}

// SetValue sets a String value of HTTP Header. Port of #setValue.
func (h *HTTPHeader) SetValue(value string) {
	h.value = value
}

// OpenStream panics: not supported for HTTPHeader. Port of #openStream.
func (h *HTTPHeader) OpenStream() (io.ReadCloser, error) {
	panic("The openStream() method is not supported for HTTPHeaderDocument.")
}

// WriteTo panics: not supported for HTTPHeader. Port of #writeTo.
func (h *HTTPHeader) WriteTo(w io.Writer) (int64, error) {
	panic("The writeTo(stream) method is not supported for HTTPHeaderDocument.")
}

// MimeType returns nil: not applicable. Port of #getMimeType.
func (h *HTTPHeader) MimeType() enumerations.MimeType {
	return nil
}

// SetMimeType panics: not supported for HTTPHeader. Port of #setMimeType.
func (h *HTTPHeader) SetMimeType(mimeType enumerations.MimeType) {
	panic("The setMimeType(mimeType) method is not supported for HTTPHeaderDocument.")
}

// Save panics: not supported for HTTPHeader. Port of #save.
func (h *HTTPHeader) Save(filePath string) error {
	panic("The save(filePath) method is not supported for HTTPHeaderDocument.")
}

// Digest panics: not supported for HTTPHeader. Port of #getDigest.
func (h *HTTPHeader) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	panic("The getDigest(digestAlgorithm) method is not supported for HTTPHeaderDocument.")
}

// DigestValue panics: not supported for HTTPHeader. Port of #getDigestValue.
func (h *HTTPHeader) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	panic("The getDigestValue(digestAlgorithm) method is not supported for HTTPHeaderDocument.")
}
