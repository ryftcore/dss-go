// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanCertificateWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// OrphanCertificateWrapper wraps an orphan certificate token extracted from a document
// (signature/timestamp).
type OrphanCertificateWrapper struct {
	OrphanCertificateTokenWrapper

	// orphanCertificate is the wrapped XmlOrphanCertificate.
	orphanCertificate *jaxb.XmlOrphanCertificate
}

// NewOrphanCertificateWrapper is the default constructor.
func NewOrphanCertificateWrapper(orphanCertificate *jaxb.XmlOrphanCertificate) *OrphanCertificateWrapper {
	w := &OrphanCertificateWrapper{
		OrphanCertificateTokenWrapper: *NewOrphanCertificateTokenWrapper(orphanCertificate.Token),
		orphanCertificate:             orphanCertificate,
	}
	w.InitOrphanTokenWrapper(w)
	return w
}

// Origins returns a list of orphan certificate origins. Port of getOrigins().
func (w *OrphanCertificateWrapper) Origins() []enumerations.CertificateOrigin {
	values := w.orphanCertificate.Origin
	if values == nil {
		return nil
	}
	result := make([]enumerations.CertificateOrigin, len(values))
	for i, v := range values {
		result[i] = enumerations.CertificateOrigin(v)
	}
	return result
}

// References returns a list of orphan certificate references. Port of getReferences().
func (w *OrphanCertificateWrapper) References() []*CertificateRefWrapper {
	var certificateRefWrappers []*CertificateRefWrapper
	for _, certificateRef := range w.orphanCertificate.CertificateRef {
		certificateRefWrappers = append(certificateRefWrappers, NewCertificateRefWrapper(certificateRef, w.Id()))
	}
	return certificateRefWrappers
}
