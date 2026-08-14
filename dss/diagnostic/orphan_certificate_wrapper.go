// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanCertificateWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
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

// GetOrigins returns a list of orphan certificate origins. Port of getOrigins().
func (w *OrphanCertificateWrapper) GetOrigins() []enumerations.CertificateOrigin {
	return w.orphanCertificate.Origins
}

// GetReferences returns a list of orphan certificate references. Port of getReferences().
func (w *OrphanCertificateWrapper) GetReferences() []*CertificateRefWrapper {
	var certificateRefWrappers []*CertificateRefWrapper
	for _, certificateRef := range w.orphanCertificate.CertificateRefs {
		certificateRefWrappers = append(certificateRefWrappers, NewCertificateRefWrapper(certificateRef, w.Id()))
	}
	return certificateRefWrappers
}
