// Ported from dss-model/.../ResponseEnvelope.java (DSS 6.5.RC1).
package http

import stdx509 "crypto/x509"

// ResponseEnvelope contains information retrieved from a Response (HTTP,
// HTTPS, etc.).
//
// java.security.cert.Certificate maps to the Go standard library's
// *crypto/x509.Certificate (aliased stdx509 to avoid confusion with a
// dss package of the same base name).
type ResponseEnvelope struct {
	// responseBody contains the response body binaries obtained from a
	// remote server.
	responseBody []byte

	// headers is the map of headers returned by the HTTP Response.
	headers map[string][]string

	// tlsCertificates are the TLS/SSL certificates used to establish a
	// secure connection, when applicable (for HTTPS calls).
	tlsCertificates []*stdx509.Certificate
}

// NewResponseEnvelope creates an empty ResponseEnvelope. Ports the empty
// constructor.
func NewResponseEnvelope() *ResponseEnvelope {
	return &ResponseEnvelope{headers: make(map[string][]string)}
}

// NewResponseEnvelopeWithBody creates a ResponseEnvelope with the given
// response message body. Ports ResponseEnvelope(byte[]).
func NewResponseEnvelopeWithBody(responseBody []byte) *ResponseEnvelope {
	return &ResponseEnvelope{responseBody: responseBody, headers: make(map[string][]string)}
}

// ResponseBody gets the response message body.
func (r *ResponseEnvelope) ResponseBody() []byte { return r.responseBody }

// SetResponseBody sets the response message body.
func (r *ResponseEnvelope) SetResponseBody(responseBody []byte) { r.responseBody = responseBody }

// Headers gets the response headers (e.g. HTTP(S) headers).
func (r *ResponseEnvelope) Headers() map[string][]string { return r.headers }

// SetHeaders sets the response headers (e.g. HTTP(S) headers). Ports
// ResponseEnvelope#setHeaders: merges headers into the existing map
// rather than replacing it (Java calls Map#putAll on the pre-initialized
// field).
func (r *ResponseEnvelope) SetHeaders(headers map[string][]string) {
	if r.headers == nil {
		r.headers = make(map[string][]string)
	}
	if headers != nil {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
}

// TLSCertificates gets the TLS/SSL certificates used by the remote server
// to establish a secure connection (e.g. for HTTPS).
func (r *ResponseEnvelope) TLSCertificates() []*stdx509.Certificate { return r.tlsCertificates }

// SetTLSCertificates sets the TLS/SSL certificates used by the remote
// server to establish a secure connection (e.g. for HTTPS).
func (r *ResponseEnvelope) SetTLSCertificates(tlsCertificates []*stdx509.Certificate) {
	r.tlsCertificates = tlsCertificates
}
