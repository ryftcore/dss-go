// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanRevocationWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// OrphanRevocationWrapper wraps document-embedded revocation data.
type OrphanRevocationWrapper struct {
	OrphanRevocationTokenWrapper

	// orphanRevocation is the wrapped XmlOrphanRevocation.
	orphanRevocation *jaxb.XmlOrphanRevocation
}

// NewOrphanRevocationWrapper is the default constructor.
func NewOrphanRevocationWrapper(orphanRevocation *jaxb.XmlOrphanRevocation) *OrphanRevocationWrapper {
	w := &OrphanRevocationWrapper{
		OrphanRevocationTokenWrapper: *NewOrphanRevocationTokenWrapper(orphanRevocation.Token),
		orphanRevocation:             orphanRevocation,
	}
	w.InitOrphanTokenWrapper(w)
	return w
}

// Origins returns a list of orphan revocation origins. Port of getOrigins().
func (w *OrphanRevocationWrapper) Origins() []enumerations.RevocationOrigin {
	values := w.orphanRevocation.Origin
	if values == nil {
		return nil
	}
	result := make([]enumerations.RevocationOrigin, len(values))
	for i, v := range values {
		result[i] = enumerations.RevocationOrigin(v)
	}
	return result
}

// References returns a list of orphan revocation references. Port of getReferences().
func (w *OrphanRevocationWrapper) References() []*RevocationRefWrapper {
	var revocationRefWrappers []*RevocationRefWrapper
	for _, revocationRef := range w.orphanRevocation.RevocationRef {
		revocationRefWrappers = append(revocationRefWrappers, NewRevocationRefWrapper(revocationRef, w.Id()))
	}
	return revocationRefWrappers
}
