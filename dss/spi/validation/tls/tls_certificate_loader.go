// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/tls/TLSCertificateLoader.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart), per PORTING.md.
package tls

import "github.com/ryftcore/dss-go/dss/model/tls"

// CertificateLoader is the data loader which includes server webpage certificates to the
// response context. Use GetTLSCertificates(url) to extract the data.
type CertificateLoader interface {
	// GetTLSCertificates extracts TLS/SSL-certificates from the given web page. Port of
	// getTLSCertificates(String).
	//
	// urlString represents a URL of a webpage with a secure connection (HTTPS); the return
	// value contains the chain of the TLS/SSL certificates and other supportive information.
	GetTLSCertificates(urlString string) *tls.Certificates
}
