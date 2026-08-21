// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/RelatedRevocationWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// RelatedRevocationWrapper wraps an XmlRelatedRevocation object. Port of
// RelatedRevocationWrapper, which extends RevocationWrapper (DIAGWRAP_B).
type RelatedRevocationWrapper struct {
	RevocationWrapper

	// relatedRevocation is the wrapped XmlRelatedRevocation.
	relatedRevocation *jaxb.XmlRelatedRevocation
}

// NewRelatedRevocationWrapper is the default constructor.
func NewRelatedRevocationWrapper(relatedRevocation *jaxb.XmlRelatedRevocation) *RelatedRevocationWrapper {
	w := &RelatedRevocationWrapper{
		RevocationWrapper: *NewRevocationWrapper(relatedRevocation.Revocation),
		relatedRevocation: relatedRevocation,
	}
	w.InitTokenProxy(w)
	return w
}

// Origins returns a list of revocation token origins. Port of getOrigins().
func (w *RelatedRevocationWrapper) Origins() []enumerations.RevocationOrigin {
	values := w.relatedRevocation.Origin
	if values == nil {
		return nil
	}
	result := make([]enumerations.RevocationOrigin, len(values))
	for i, v := range values {
		result[i] = enumerations.RevocationOrigin(v)
	}
	return result
}

// References returns a list of revocation token references from the signature. Port of
// getReferences().
func (w *RelatedRevocationWrapper) References() []*RevocationRefWrapper {
	var references []*RevocationRefWrapper
	for _, revocationRef := range w.relatedRevocation.RevocationRef {
		references = append(references, NewRevocationRefWrapper(revocationRef, w.Id()))
	}
	return references
}
