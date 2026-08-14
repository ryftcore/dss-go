// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/CertificateRevocationWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// CertificateRevocationWrapper is a complete revocation wrapper, containing detailed
// certificate revocation and common information. Port of CertificateRevocationWrapper, which
// extends RevocationWrapper (DIAGWRAP_B).
type CertificateRevocationWrapper struct {
	RevocationWrapper

	// certificateRevocation is the wrapped XmlCertificateRevocation.
	certificateRevocation *jaxb.XmlCertificateRevocation
}

// NewCertificateRevocationWrapper is the default constructor.
func NewCertificateRevocationWrapper(certificateRevocation *jaxb.XmlCertificateRevocation) *CertificateRevocationWrapper {
	w := &CertificateRevocationWrapper{
		RevocationWrapper:     *NewRevocationWrapper(certificateRevocation.Revocation),
		certificateRevocation: certificateRevocation,
	}
	w.InitTokenProxy(w)
	return w
}

// Status returns the revocation status of the concerned certificate. Port of getStatus().
func (w *CertificateRevocationWrapper) Status() enumerations.CertificateStatus {
	if w.certificateRevocation.Status != nil {
		return enumerations.CertificateStatus(*w.certificateRevocation.Status)
	}
	return ""
}

// Reason returns the revocation reason for the concerned certificate. Port of getReason().
func (w *CertificateRevocationWrapper) Reason() enumerations.RevocationReason {
	if w.certificateRevocation.Reason != nil {
		return enumerations.RevocationReason(*w.certificateRevocation.Reason)
	}
	return ""
}

// RevocationDate returns the revocation time for the concerned certificate. Port of
// getRevocationDate().
func (w *CertificateRevocationWrapper) RevocationDate() *time.Time {
	if w.certificateRevocation.RevocationDate == nil {
		return nil
	}
	t := w.certificateRevocation.RevocationDate.Time()
	return &t
}

// IsRevoked reports whether the concerned certificate has been revoked. Port of isRevoked().
func (w *CertificateRevocationWrapper) IsRevoked() bool {
	return w.Status().IsRevoked()
}

// IsKnown reports whether the revocation status is known. Port of isKnown().
func (w *CertificateRevocationWrapper) IsKnown() bool {
	return w.Status().IsKnown()
}
