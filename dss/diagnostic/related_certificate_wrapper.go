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

// Origins returns a list of certificate token origins. Port of getOrigins().
func (w *RelatedCertificateWrapper) Origins() []enumerations.CertificateOrigin {
	values := w.relatedCertificate.Origin
	if values == nil {
		return nil
	}
	result := make([]enumerations.CertificateOrigin, len(values))
	for i, v := range values {
		result[i] = enumerations.CertificateOrigin(v)
	}
	return result
}

// References returns a list of certificate token references from the signature. Port of
// getReferences().
func (w *RelatedCertificateWrapper) References() []*CertificateRefWrapper {
	var references []*CertificateRefWrapper
	for _, certificateRef := range w.relatedCertificate.CertificateRef {
		references = append(references, NewCertificateRefWrapper(certificateRef, w.Id()))
	}
	return references
}
