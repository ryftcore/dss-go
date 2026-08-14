// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/RelatedCertificateWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// RelatedCertificateWrapper wraps an XmlRelatedCertificate object.
type RelatedCertificateWrapper struct {
	CertificateWrapper

	// relatedCertificate is the wrapped XmlRelatedCertificate.
	relatedCertificate *jaxb.XmlRelatedCertificate
}

// NewRelatedCertificateWrapper is the default constructor.
func NewRelatedCertificateWrapper(relatedCertificate *jaxb.XmlRelatedCertificate) *RelatedCertificateWrapper {
	w := &RelatedCertificateWrapper{
		CertificateWrapper: *NewCertificateWrapper(relatedCertificate.Certificate),
		relatedCertificate: relatedCertificate,
	}
	w.InitTokenProxy(w)
	return w
}

// GetOrigins returns a list of certificate token origins. Port of getOrigins().
func (w *RelatedCertificateWrapper) GetOrigins() []enumerations.CertificateOrigin {
	return w.relatedCertificate.Origins
}

// GetReferences returns a list of certificate token references from the signature. Port of
// getReferences().
func (w *RelatedCertificateWrapper) GetReferences() []*CertificateRefWrapper {
	var references []*CertificateRefWrapper
	for _, certificateRef := range w.relatedCertificate.CertificateRefs {
		references = append(references, NewCertificateRefWrapper(certificateRef, w.Id()))
	}
	return references
}
