// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationRef.java (DSS 6.5.RC1).
//
// RevocationRef<R> is the abstract superclass CRLRef and OCSPRef (chunk CRLOCSP, a sibling of
// this phase 2a chunk) already embed via RevocationRefBase[R] and register with
// InitRevocationRef, mirroring the InitToken self-registration pattern.
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// RevocationRef represents an extracted revocation reference from a signature.
type RevocationRef[R revocation.Revocation] interface {
	model.IdentifierBasedObject

	// Digest returns digests of the reference. Port of getDigest().
	Digest() model.Digest
	// DSSIDAsString returns the revocation reference string id. Port of getDSSIdAsString().
	DSSIDAsString() string
	// String renders the reference. Port of toString(): `Utils.toBase64(digest.getValue())`,
	// unless the concrete reference overrides it (CRLRef and OCSPRef both do).
	String() string
}

// RevocationRefOverrides captures what RevocationRefBase needs to reach through virtual
// dispatch: building the reference's unique identifier. Port of the protected
// createIdentifier() override point; RevocationRefBase itself supplies the default body
// (`return new RevocationRefIdentifier(this);`), promoted to a concrete reference unless it
// overrides CreateIdentifier itself (OCSPRef does, to build an OCSPRefIdentifier instead).
type RevocationRefOverrides interface {
	CreateIdentifier() model.Identifier
}

// RevocationRefBase carries the state and the concrete behaviour of the Java abstract class
// RevocationRef<R>. Concrete references (CRLRef, OCSPRef) embed it and register themselves
// with InitRevocationRef.
type RevocationRefBase[R revocation.Revocation] struct {
	// overrides points back at the concrete reference; see InitRevocationRef.
	overrides RevocationRefOverrides

	// digest is the digest within the reference.
	digest model.Digest
	// identifier caches the value DSSID returns, computed on first use.
	identifier model.Identifier
}

// NewRevocationRefBase instantiates the base state of a revocation reference with the Java
// default (null) values. Port of the protected default constructor; the concrete reference
// must still call InitRevocationRef.
func NewRevocationRefBase[R revocation.Revocation]() RevocationRefBase[R] {
	return RevocationRefBase[R]{}
}

// InitRevocationRef registers the concrete reference with its base so that the base can
// dispatch CreateIdentifier the way Java reaches an overridden createIdentifier() through
// virtual dispatch. It must be called exactly once, by the concrete reference's constructor,
// before any other method.
func (r *RevocationRefBase[R]) InitRevocationRef(overrides RevocationRefOverrides) {
	r.overrides = overrides
}

// revocationRefBaseOverrides returns the registered overrides, panicking when the concrete
// reference forgot to call InitRevocationRef.
func (r *RevocationRefBase[R]) revocationRefBaseOverrides() RevocationRefOverrides {
	if r.overrides == nil {
		panic("RevocationRef was not initialised: the concrete reference must call InitRevocationRef in its constructor")
	}
	return r.overrides
}

// Digest returns digests of the reference. Port of getDigest().
func (r *RevocationRefBase[R]) Digest() model.Digest {
	return r.digest
}

// SetDigest sets the digest within the reference; the Go counterpart of writing Java's
// protected digest field directly from a subclass.
func (r *RevocationRefBase[R]) SetDigest(digest model.Digest) {
	r.digest = digest
}

// CreateIdentifier creates the reference's unique identifier: `new RevocationRefIdentifier(this)`.
// Port of the protected createIdentifier() default body; promoted to a concrete reference that
// does not override it (CRLRef).
func (r *RevocationRefBase[R]) CreateIdentifier() model.Identifier {
	return NewRevocationRefIdentifier(r.digest)
}

// DSSID returns the revocation ref DSS Identifier, building it (through the registered
// overrides, so a reference overriding CreateIdentifier is honoured) and caching it on first
// use. Port of getDSSId().
func (r *RevocationRefBase[R]) DSSID() model.Identifier {
	if r.identifier == nil {
		r.identifier = r.revocationRefBaseOverrides().CreateIdentifier()
	}
	return r.identifier
}

// DSSIDAsString returns revocation reference string id. Port of getDSSIdAsString().
func (r *RevocationRefBase[R]) DSSIDAsString() string {
	return r.DSSID().AsXmlID()
}

// String renders the reference. Port of toString(): `Utils.toBase64(digest.getValue())`.
// Promoted to a concrete reference that does not override it; CRLRef and OCSPRef both do,
// appending this rendering after their own fields (`... + r.RevocationRefBase.String()`).
func (r *RevocationRefBase[R]) String() string {
	return utils.ToBase64(r.digest.Value())
}

// compile-time assertion: a RevocationRefBase satisfies RevocationRefOverrides, i.e. it is a
// valid default CreateIdentifier provider for a concrete reference that does not override it.
var _ RevocationRefOverrides = (*RevocationRefBase[revocation.CRL])(nil)
