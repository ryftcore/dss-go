// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/EncapsulatedRevocationTokenIdentifier.java (DSS 6.5.RC1).
package model

import (
	"github.com/utain/esig/dss/model/x509/revocation"
)

// EncapsulatedRevocationTokenIdentifier is a unique identifier for revocation data binaries.
// R is the revocation data implementation the identifier belongs to (revocation.CRL or
// revocation.OCSP), mirroring the Java type parameter; like in Java it is a phantom
// parameter that only constrains callers.
type EncapsulatedRevocationTokenIdentifier[R revocation.Revocation] struct {
	MultipleDigestIdentifier
}

// NewEncapsulatedRevocationTokenIdentifier builds the identifier over the revocation data
// binaries, rendered with the "R-" prefix.
func NewEncapsulatedRevocationTokenIdentifier[R revocation.Revocation](binaries []byte) *EncapsulatedRevocationTokenIdentifier[R] {
	return NewEncapsulatedRevocationTokenIdentifierWithClassName[R]("EncapsulatedRevocationTokenIdentifier", binaries)
}

// NewEncapsulatedRevocationTokenIdentifierWithClassName is the constructor subclasses in
// other packages (CRLBinary, OCSPResponseBinary, ...) use, so that the Java simple class
// name their toString() and equals() rely on stays correct.
func NewEncapsulatedRevocationTokenIdentifierWithClassName[R revocation.Revocation](className string, binaries []byte) *EncapsulatedRevocationTokenIdentifier[R] {
	return &EncapsulatedRevocationTokenIdentifier[R]{
		NewMultipleDigestIdentifier(className, "R-", binaries),
	}
}

// DSSID returns this identifier. Port of getDSSId().
func (e *EncapsulatedRevocationTokenIdentifier[R]) DSSID() Identifier {
	return e
}

// compile-time interface assertions.
var (
	_ Identifier            = (*EncapsulatedRevocationTokenIdentifier[revocation.CRL])(nil)
	_ IdentifierBasedObject = (*EncapsulatedRevocationTokenIdentifier[revocation.CRL])(nil)
)
