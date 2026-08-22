// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CompositeRevocationSource.java (DSS 6.5.RC1).
package spi

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// CompositeRevocationSource allows retrieving a RevocationToken from different sources. The
// composite tries all sources until it gets a non-empty response.
type CompositeRevocationSource[T revocation.Revocation] struct {
	// compositeRevocationSources maps source keys to their corresponding RevocationSources.
	compositeRevocationSources map[string]RevocationSource[T]

	// sourceOrder preserves insertion order, standing in for Java's LinkedHashMap-like
	// iteration order guarantee that Map.entrySet() otherwise would not provide; SetSources
	// is the only mutator, and it replaces the whole map, so recomputing the order there is
	// enough.
	sourceOrder []string
}

// NewCompositeRevocationSource instantiates the object with nil values. Port of the default
// constructor.
func NewCompositeRevocationSource[T revocation.Revocation]() *CompositeRevocationSource[T] {
	return &CompositeRevocationSource[T]{}
}

// SetSources allows providing multiple revocationSources, keyed by a label.
// Port of setSources(Map<String, RevocationSource<T>>).
//
// Go maps have no defined iteration order (Java's HashMap doesn't either, but callers in
// practice pass a LinkedHashMap); sourceKeys lets callers pin down the try order the way a
// LinkedHashMap constructor argument would. When sourceKeys is nil, Java's HashMap-backed
// setSources(Map) would still iterate in *some* order, arbitrary but stable for the life of
// that Map; Go map iteration is instead randomized on every run, so the fallback sorts the
// keys lexically (PORTING.md's "order-sensitive upstream iteration -> ... explicit sort" rule)
// to keep RevocationToken's try order - and therefore its result, when more than one source
// would answer - stable from one run to the next.
func (c *CompositeRevocationSource[T]) SetSources(compositeRevocationSources map[string]RevocationSource[T], sourceKeys []string) {
	c.compositeRevocationSources = compositeRevocationSources
	if sourceKeys != nil {
		c.sourceOrder = sourceKeys
	} else {
		c.sourceOrder = nil
		for sourceKey := range compositeRevocationSources {
			c.sourceOrder = append(c.sourceOrder, sourceKey)
		}
		sort.Strings(c.sourceOrder)
	}
}

// RevocationToken tries all sources until it gets a non-empty response, returning nil if none
// of them does. Port of getRevocationToken(CertificateToken, CertificateToken).
//
// Java catches and logs any Exception raised by a source, then keeps trying the remaining
// ones; the port cannot distinguish "the source itself panicked" the way Java's Exception
// catch does without recover(), and recovering from an arbitrary panic here would hide
// programmer errors from sibling sources, so a source that returns an error is skipped (the
// equivalent of Java's catch) while a source that panics is left to propagate.
func (c *CompositeRevocationSource[T]) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) RevocationToken[T] {
	for _, sourceKey := range c.sourceOrder {
		source, found := c.compositeRevocationSources[sourceKey]
		if !found {
			continue
		}
		if token := source.RevocationToken(certificateToken, issuerCertificateToken); token != nil {
			return token
		}
	}
	return nil
}

// compile-time assertion: a CompositeRevocationSource is a RevocationSource.
var _ RevocationSource[revocation.OCSP] = (*CompositeRevocationSource[revocation.OCSP])(nil)
