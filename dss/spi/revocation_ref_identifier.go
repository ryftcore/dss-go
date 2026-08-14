// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationRefIdentifier.java (DSS 6.5.RC1).
//
// OCSPRefIdentifier (chunk CRLOCSP, a sibling of this phase 2a chunk) already embeds this type
// by value and calls NewRevocationRefIdentifierFromDigest("OCSPRefIdentifier", ...); this file
// defines it to match that shape.
package spi

import (
	"github.com/utain/esig/dss/model"
)

// RevocationRefIdentifier is a unique id for a revocation reference.
type RevocationRefIdentifier struct {
	model.IdentifierBase
}

// NewRevocationRefIdentifier builds the identifier over the given digest, with the Java simple
// class name "RevocationRefIdentifier" and the "R-" prefix.
//
// Java's public constructor RevocationRefIdentifier(RevocationRef<?>) only ever forwards to the
// protected RevocationRefIdentifier(Digest) overload (`this(revocationRef.getDigest())`); since
// that is all the Go port needs (see RevocationRefBase's default CreateIdentifier), a single
// digest-based constructor replaces both.
func NewRevocationRefIdentifier(digest model.Digest) *RevocationRefIdentifier {
	identifier := NewRevocationRefIdentifierFromDigest("RevocationRefIdentifier", digest)
	return &identifier
}

// NewRevocationRefIdentifierFromDigest builds the identifier over the given digest with an
// explicit Java simple class name, rendered with the "R-" prefix.
//
// A subclass overriding createIdentifier() to build its own identifier type (OCSPRefIdentifier
// does, embedding this type by value) reports its own class name here: Java's getClass().
// getSimpleName() finds it polymorphically at runtime, but IdentifierBase stores the class name
// as an explicit field (see model/identifier.go), so nothing in this Go port can infer it from
// the embedding alone - it must be supplied by the caller, matching the caller's own type name.
func NewRevocationRefIdentifierFromDigest(className string, digest model.Digest) RevocationRefIdentifier {
	return RevocationRefIdentifier{model.NewIdentifierBaseFromDigest(className, "R-", digest)}
}

// compile-time assertion: a RevocationRefIdentifier is an Identifier.
var _ model.Identifier = (*RevocationRefIdentifier)(nil)
