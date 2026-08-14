// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/CertificateRefWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// CertificateRefWrapper represents a certificate reference wrapper.
type CertificateRefWrapper struct {
	// certificateRef is the wrapped XmlCertificateRef.
	certificateRef *jaxb.XmlCertificateRef

	// certificateId is the Id of the related certificate token.
	certificateId string
}

// NewCertificateRefWrapper is the default constructor.
func NewCertificateRefWrapper(certificateRef *jaxb.XmlCertificateRef, certificateId string) *CertificateRefWrapper {
	return &CertificateRefWrapper{certificateRef: certificateRef, certificateId: certificateId}
}

// Origin returns a certificate reference origin. Port of getOrigin().
func (w *CertificateRefWrapper) Origin() enumerations.CertificateRefOrigin {
	return w.certificateRef.Origin
}

// GetIssuerSerial returns IssuerSerial's binaries. Port of getIssuerSerial().
func (w *CertificateRefWrapper) GetIssuerSerial() []byte {
	if w.certificateRef.IssuerSerial != nil {
		return w.certificateRef.IssuerSerial.Value
	}
	return nil
}

// IsIssuerSerialPresent reports whether IssuerSerial is present within the certificate
// reference. Port of isIssuerSerialPresent().
func (w *CertificateRefWrapper) IsIssuerSerialPresent() bool {
	return w.certificateRef.IssuerSerial != nil
}

// IsIssuerSerialMatch reports whether IssuerSerial matches with a found certificate. Port of
// isIssuerSerialMatch().
func (w *CertificateRefWrapper) IsIssuerSerialMatch() bool {
	if w.certificateRef.IssuerSerial != nil && w.certificateRef.IssuerSerial.Match != nil {
		return *w.certificateRef.IssuerSerial.Match
	}
	return false
}

// GetIssuerName returns IssuerName. Port of getIssuerName().
func (w *CertificateRefWrapper) GetIssuerName() string {
	if w.certificateRef.SerialInfo != nil {
		return w.certificateRef.SerialInfo.IssuerName
	}
	return ""
}

// GetSki returns SKI of the certificate (SHA-1 of the certificate's public key). Port of
// getSki().
func (w *CertificateRefWrapper) GetSki() []byte {
	if w.certificateRef.SerialInfo != nil {
		return w.certificateRef.SerialInfo.Ski
	}
	return nil
}

// GetKid returns KID (key identifier) of the certificate reference. NOTE: if the original
// signature contains a byte string, the returned value of the method is base64-encoded
// result. Port of getKid().
func (w *CertificateRefWrapper) GetKid() string {
	return w.certificateRef.KID
}

// GetX509Url returns X.509 URL (key identifier) of the certificate reference. Port of
// getX509Url().
func (w *CertificateRefWrapper) GetX509Url() string {
	return w.certificateRef.X509Url
}

// GetDigestAlgoAndValue returns digest algo and value. Port of getDigestAlgoAndValue().
func (w *CertificateRefWrapper) GetDigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.certificateRef.DigestAlgoAndValue
}

// IsDigestValuePresent checks if the DigestAlgoAndValue of the reference is present in the
// certificate reference. Port of isDigestValuePresent().
func (w *CertificateRefWrapper) IsDigestValuePresent() bool {
	return w.GetDigestAlgoAndValue() != nil
}

// GetDigestMethod returns a used DigestAlgorithm for a certificate reference creation. Port of
// getDigestMethod().
func (w *CertificateRefWrapper) GetDigestMethod() enumerations.DigestAlgorithm {
	digestAlgoAndValue := w.GetDigestAlgoAndValue()
	if digestAlgoAndValue != nil && digestAlgoAndValue.DigestMethod != nil {
		return *digestAlgoAndValue.DigestMethod
	}
	return ""
}

// IsDigestValueMatch checks if the DigestAlgoAndValue of the reference matches one of the
// certificate. Port of isDigestValueMatch().
func (w *CertificateRefWrapper) IsDigestValueMatch() bool {
	digestAlgoAndValue := w.GetDigestAlgoAndValue()
	return digestAlgoAndValue != nil && digestAlgoAndValue.Match != nil && *digestAlgoAndValue.Match
}

// GetCertificateId returns Id of the referenced certificate token (when available) or the
// reference id. Port of getCertificateId().
func (w *CertificateRefWrapper) GetCertificateId() string {
	return w.certificateId
}

// String returns a string representation of the wrapper. Port of toString().
func (w *CertificateRefWrapper) String() string {
	if w.certificateRef != nil {
		return "CertificateRefWrapper Origin='" + string(w.certificateRef.Origin) + "'"
	}
	return "CertificateRefWrapper certificateRef=<nil>"
}
