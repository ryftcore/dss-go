// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanCertificateTokenWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
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
	if w.orphanToken.Id != nil {
		return string(*w.orphanToken.Id)
	}
	return ""
}

// CertificateDN returns the certificate's Distinguished Name (by RFC 2253). Port of
// getCertificateDN().
func (w *OrphanCertificateTokenWrapper) CertificateDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.orphanToken.SubjectDistinguishedName)
	return distinguishedNameListWrapper.Value("RFC2253")
}

// CertificateIssuerDN returns the certificate issuer's Distinguished Name (by RFC 2253).
// Port of getCertificateIssuerDN().
func (w *OrphanCertificateTokenWrapper) CertificateIssuerDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.orphanToken.IssuerDistinguishedName)
	return distinguishedNameListWrapper.Value("RFC2253")
}

// SerialNumber returns the serial number of the certificate. Port of getSerialNumber().
func (w *OrphanCertificateTokenWrapper) SerialNumber() string {
	serialNumber := w.orphanToken.SerialNumber
	if serialNumber == nil {
		return ""
	}
	return serialNumber.String()
}

// NotBefore returns the certificate's notBefore date. Port of getNotBefore().
func (w *OrphanCertificateTokenWrapper) NotBefore() *time.Time {
	if w.orphanToken.NotBefore == nil {
		return nil
	}
	t := w.orphanToken.NotBefore.Time()
	return &t
}

// NotAfter returns the certificate's notAfter date. Port of getNotAfter().
func (w *OrphanCertificateTokenWrapper) NotAfter() *time.Time {
	if w.orphanToken.NotAfter == nil {
		return nil
	}
	t := w.orphanToken.NotAfter.Time()
	return &t
}

// EntityKey returns a string identifier of the certificate's public key. Port of
// getEntityKey().
func (w *OrphanCertificateTokenWrapper) EntityKey() string {
	if w.orphanToken.EntityKey != nil {
		return *w.orphanToken.EntityKey
	}
	return ""
}

// IsTrusted reports whether the certificate is trusted. Port of isTrusted().
func (w *OrphanCertificateTokenWrapper) IsTrusted() bool {
	return w.orphanToken.Trusted != nil && *w.orphanToken.Trusted
}

// IsSelfSigned reports whether the certificate is self-signed. Port of isSelfSigned().
func (w *OrphanCertificateTokenWrapper) IsSelfSigned() bool {
	return w.orphanToken.SelfSigned != nil && *w.orphanToken.SelfSigned
}

// Binaries returns base64-encoded byte array of the token. Port of getBinaries().
func (w *OrphanCertificateTokenWrapper) Binaries() []byte {
	if w.orphanToken.Base64Encoded == nil {
		return nil
	}
	return []byte(*w.orphanToken.Base64Encoded)
}

// DigestAlgoAndValue returns digest of the token. Port of getDigestAlgoAndValue().
func (w *OrphanCertificateTokenWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.orphanToken.DigestAlgoAndValue
}
