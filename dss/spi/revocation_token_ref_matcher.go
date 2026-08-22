// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationTokenRefMatcher.java (DSS 6.5.RC1).
//
// CRLTokenRefMatcher and OCSPTokenRefMatcher implement this interface and reference
// EncapsulatedRevocationTokenIdentifier[R] unqualified, so both are defined here to match the
// shape those files were written against.
package spi

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// EncapsulatedRevocationTokenIdentifier is the Go-only interface counterpart of the dss-model
// class eu.europa.esig.dss.model.identifier.EncapsulatedRevocationTokenIdentifier<R>, ported as
// the concrete generic struct model.EncapsulatedRevocationTokenIdentifier[R]. Java's
// abstract-class polymorphism lets OfflineRevocationSource and this matcher hold either a
// CRLBinary or an OCSPResponseBinary through one supertype; Go structs cannot be referenced
// polymorphically that way, so this package-local interface - satisfied structurally by any type
// embedding model.EncapsulatedRevocationTokenIdentifier[R], such as CRLBinary (crlparser
// package) and OCSPResponseBinary (this package) - stands in for it. It is declared here, next
// to the matcher whose match(EncapsulatedRevocationTokenIdentifier<R>, RevocationRef<R>)
// overload is the most direct Java consumer of the type, though OfflineRevocationSource.java and
// ListRevocationSource.java both use it just as heavily.
type EncapsulatedRevocationTokenIdentifier[R revocation.Revocation] interface {
	model.IdentifierBasedObject

	// AsXmlID returns an ID conformant to XML Id. Port of asXmlId().
	AsXmlID() string
	// DigestValue returns the digest value of the binaries for the given algorithm. Port of
	// getDigestValue(DigestAlgorithm).
	DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error)
}

// RevocationTokenRefMatcher validates a revocation reference against a revocation token.
type RevocationTokenRefMatcher[R revocation.Revocation] interface {
	// Match returns true if the reference is related to the provided token. Port of
	// match(RevocationToken, RevocationRef).
	Match(token RevocationToken[R], reference RevocationRef[R]) (bool, error)

	// MatchBinary returns true if the reference is related to the encapsulated identifier.
	// Port of the match(EncapsulatedRevocationTokenIdentifier, RevocationRef) overload, which
	// Go's lack of overloading gives a distinct name from Match.
	MatchBinary(identifier EncapsulatedRevocationTokenIdentifier[R], reference RevocationRef[R]) (bool, error)
}
