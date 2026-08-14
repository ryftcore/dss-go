// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanCertificateTokenWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
)

// OrphanCertificateTokenWrapper wraps an XML orphan certificate. Port of
// OrphanTokenWrapper<XmlOrphanCertificateToken>'s certificate-token specialization.
type OrphanCertificateTokenWrapper struct {
	OrphanTokenWrapperBase

	// orphanToken is the wrapped XmlOrphanCertificateToken.
	orphanToken *jaxb.XmlOrphanCertificateToken
}

// NewOrphanCertificateTokenWrapper is the default constructor. Port of
// OrphanCertificateTokenWrapper(XmlOrphanCertificateToken).
func NewOrphanCertificateTokenWrapper(orphanToken *jaxb.XmlOrphanCertificateToken) *OrphanCertificateTokenWrapper {
	w := &OrphanCertificateTokenWrapper{orphanToken: orphanToken}
	w.InitOrphanTokenWrapper(w)
	return w
}

// Id returns identifier of the orphan token. Port of getId() (inherited from
// OrphanTokenWrapper).
func (w *OrphanCertificateTokenWrapper) Id() string {
	return w.orphanToken.Id
}

// GetCertificateDN returns the certificate's Distinguished Name (by RFC 2253). Port of
// getCertificateDN().
func (w *OrphanCertificateTokenWrapper) GetCertificateDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.orphanToken.SubjectDistinguishedName)
	return distinguishedNameListWrapper.GetValue("RFC2253")
}

// GetCertificateIssuerDN returns the certificate issuer's Distinguished Name (by RFC 2253).
// Port of getCertificateIssuerDN().
func (w *OrphanCertificateTokenWrapper) GetCertificateIssuerDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.orphanToken.IssuerDistinguishedName)
	return distinguishedNameListWrapper.GetValue("RFC2253")
}

// GetSerialNumber returns the serial number of the certificate. Port of getSerialNumber().
func (w *OrphanCertificateTokenWrapper) GetSerialNumber() string {
	serialNumber := w.orphanToken.SerialNumber
	if serialNumber == nil {
		return ""
	}
	return serialNumber.String()
}

// GetNotBefore returns the certificate's notBefore date. Port of getNotBefore().
func (w *OrphanCertificateTokenWrapper) GetNotBefore() *time.Time {
	return w.orphanToken.NotBefore
}

// GetNotAfter returns the certificate's notAfter date. Port of getNotAfter().
func (w *OrphanCertificateTokenWrapper) GetNotAfter() *time.Time {
	return w.orphanToken.NotAfter
}

// GetEntityKey returns a string identifier of the certificate's public key. Port of
// getEntityKey().
func (w *OrphanCertificateTokenWrapper) GetEntityKey() string {
	return w.orphanToken.EntityKey
}

// IsTrusted reports whether the certificate is trusted. Port of isTrusted().
func (w *OrphanCertificateTokenWrapper) IsTrusted() bool {
	return w.orphanToken.Trusted
}

// IsSelfSigned reports whether the certificate is self-signed. Port of isSelfSigned().
func (w *OrphanCertificateTokenWrapper) IsSelfSigned() bool {
	return w.orphanToken.SelfSigned
}

// Binaries returns base64-encoded byte array of the token. Port of getBinaries().
func (w *OrphanCertificateTokenWrapper) Binaries() []byte {
	return w.orphanToken.Base64Encoded
}

// DigestAlgoAndValue returns digest of the token. Port of getDigestAlgoAndValue().
func (w *OrphanCertificateTokenWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.orphanToken.DigestAlgoAndValue
}
