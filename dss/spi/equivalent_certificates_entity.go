// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/EquivalentCertificatesEntity.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi. The Java
// class has default (package) visibility, so the Go type stays unexported: it is an internal
// detail of CommonCertificateSource (ported separately, chunk X509-B) and its subclasses.
package spi

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// equivalentCertificatesEntity re-groups equivalent certificates by a given property (e.g. a
// public key or entity key). All certificates for a given equivalentCertificatesEntity share
// (at least) the same public key.
type equivalentCertificatesEntity struct {
	// identifier is the unique Id for all certificates (SHA-256 of the common public key).
	identifier *model.EntityIdentifier

	// ski is the Subject Key Identifier (SHA-1 of the common public key).
	ski []byte

	// equivalentCertificates holds the equivalent certificates (certificates sharing the same
	// public key), keyed by DSSIDAsString(): Java relies on hashCode-based collections keyed
	// on tokens, the Go equivalent keys maps on the token's identifier string (see
	// CertificateToken.Equals in dss-model). Java's HashSet iteration order is arbitrary but
	// stable within a JVM run; a bare Go map is randomized on every run instead, so this is
	// kept insertion-ordered (slice + index map, PORTING.md's Collections rule) since
	// EquivalentCertificates() feeds CommonCertificateSource.Certificates()'s returned order.
	equivalentCertificates *utils.OrderedMap[string, *model.CertificateToken]
}

// newEquivalentCertificatesEntity builds the entity around its first member.
// Port of the package-private EquivalentCertificatesEntity(CertificateToken) constructor.
//
// DSSASN1UtilsComputeSkiFromCert can fail on a malformed public key encoding; Java's
// computeSkiFromCert wraps that failure as an unchecked DSSException, so it is exceptional
// here too - the error is data-dependent on the certificate's own encoding, so it is returned
// rather than panicking, letting the caller (CommonCertificateSource.AddCertificate) decide.
func newEquivalentCertificatesEntity(initialCert *model.CertificateToken) (*equivalentCertificatesEntity, error) {
	ski, err := DSSASN1UtilsComputeSkiFromCert(initialCert)
	if err != nil {
		return nil, err
	}
	equivalentCertificates := utils.NewOrderedMap[string, *model.CertificateToken]()
	equivalentCertificates.Set(initialCert.DSSIDAsString(), initialCert)
	return &equivalentCertificatesEntity{
		identifier:             initialCert.EntityKey(),
		ski:                    ski,
		equivalentCertificates: equivalentCertificates,
	}, nil
}

// addEquivalentCertificate adds a certificate token to the given list of equivalent
// certificates. Port of addEquivalentCertificate(CertificateToken).
//
// Java silently skips a token whose recomputed SKI does not match the entity's ("this should
// never happen"); the recomputation failure itself is returned to the caller rather than
// silently ignored, since - unlike the SKI mismatch - it does signal a genuinely malformed
// public key.
func (e *equivalentCertificatesEntity) addEquivalentCertificate(token *model.CertificateToken) error {
	if _, found := e.equivalentCertificates.Get(token.DSSIDAsString()); found {
		return nil
	}
	// we manually recompute the SKI (we had cases with wrongly encoded value in the certificate)
	newSKI, err := DSSASN1UtilsComputeSkiFromCert(token)
	if err != nil {
		return err
	}
	// This should never happen
	if !bytes.Equal(newSKI, e.ski) {
		return nil
	}
	e.equivalentCertificates.Set(token.DSSIDAsString(), token)
	return nil
}

// removeEquivalentCertificate removes a certificate token from the given list of equivalent
// certificates. Port of removeEquivalentCertificate(CertificateToken); Java refuses to empty
// the pool, since an empty pool is not a valid state for an entity to be in.
func (e *equivalentCertificatesEntity) removeEquivalentCertificate(token *model.CertificateToken) {
	key := token.DSSIDAsString()
	if _, found := e.equivalentCertificates.Get(key); !found {
		return
	}
	if e.equivalentCertificates.Len() == 1 {
		return
	}
	e.equivalentCertificates.Delete(key)
}

// Ski gets a Subject Key Identifier (SHA-1 of the common public key). Port of getSki().
func (e *equivalentCertificatesEntity) Ski() []byte {
	return e.ski
}

// EquivalentCertificates gets a set of equivalent certificate tokens present within the
// current instance, keyed by DSSIDAsString(). Port of getEquivalentCertificates(); the
// returned map is a defensive copy, standing in for Java's Collections.unmodifiableSet.
func (e *equivalentCertificatesEntity) EquivalentCertificates() map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken, e.equivalentCertificates.Len())
	for _, k := range e.equivalentCertificates.Keys() {
		v, _ := e.equivalentCertificates.Get(k)
		result[k] = v
	}
	return result
}

// orderedEquivalentCertificates returns the equivalent certificate tokens in insertion order,
// used by CommonCertificateSource.Certificates() so its returned slice is deterministic (the
// map-returning EquivalentCertificates() above has no order contract to preserve, unlike this).
func (e *equivalentCertificatesEntity) orderedEquivalentCertificates() []*model.CertificateToken {
	return e.equivalentCertificates.Values()
}

// Equals reports whether both entities were built from certificates sharing the same entity
// identifier. Port of equals(Object); Java's paired hashCode() has no Go counterpart.
func (e *equivalentCertificatesEntity) Equals(other *equivalentCertificatesEntity) bool {
	if e == other {
		return true
	}
	if other == nil {
		return false
	}
	if e.identifier == nil {
		return other.identifier == nil
	}
	return e.identifier.Equals(other.identifier)
}

// isCertificateSourceEntity implements the CertificateSourceEntity marker interface.
func (e *equivalentCertificatesEntity) isCertificateSourceEntity() {}

// compile-time assertion: an equivalentCertificatesEntity is a CertificateSourceEntity.
var _ CertificateSourceEntity = (*equivalentCertificatesEntity)(nil)
