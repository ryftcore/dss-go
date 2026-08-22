// Ported from dss-model/.../TLSCertificates.java (DSS 6.5.RC1).
package tls

import "github.com/ryftcore/dss-go/dss/model"

// Certificates represents information obtained from a remote server as
// the result of the TLS/SSL handshake.
//
// model.CertificateToken (ported from model.x509, flattened into the
// model package) is outside this manifest; assumed to already exist
// there.
type Certificates struct {
	// certificates is the collection of certificate tokens returned
	// after the TLS/SSL handshake.
	certificates []*model.CertificateToken

	// tlsCertificateBindingUrl is the value of the "Link" response header
	// with a rel value of tls-certificate-binding.
	tlsCertificateBindingUrl string
}

// NewTLSCertificates creates an empty TLSCertificates. Ports the empty
// constructor.
func NewCertificates() *Certificates {
	return &Certificates{}
}

// Certificates gets a list of certificates returned by a remote server
// during the TLS/SSL handshake.
func (t *Certificates) Certificates() []*model.CertificateToken { return t.certificates }

// SetCertificates sets a list of certificates returned by a remote server
// during the TLS/SSL handshake.
func (t *Certificates) SetCertificates(certificates []*model.CertificateToken) {
	t.certificates = certificates
}

// TLSCertificateBindingUrl gets the value of the "Link" response header
// with a rel value of tls-certificate-binding. This URL is used to
// extract a TLS/SSL binding signature.
func (t *Certificates) TLSCertificateBindingUrl() string { return t.tlsCertificateBindingUrl }

// SetTLSCertificateBindingUrl sets the value of the "Link" response
// header with a rel value of tls-certificate-binding. This URL is used to
// extract a TLS/SSL binding signature.
func (t *Certificates) SetTLSCertificateBindingUrl(tlsCertificateBindingUrl string) {
	t.tlsCertificateBindingUrl = tlsCertificateBindingUrl
}
