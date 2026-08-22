// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/Identifier.java (DSS 6.5.RC1).
package model

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"golang.org/x/crypto/ripemd160"
)

// IdentifierDigestAlgorithm is the DigestAlgorithm every Identifier is computed with.
// Port of the protected constant Identifier.DIGEST_ALGO.
const IdentifierDigestAlgorithm = enumerations.DigestAlgorithm_SHA256

// Identifier is the contract of the abstract Java class
// eu.europa.esig.dss.model.identifier.Identifier. Java's covariant overrides
// (Token#getDSSId returning a TokenIdentifier, EntityIdentifierBuilder#build returning an
// EntityIdentifier) cannot be expressed in Go, so every polymorphic position uses this
// interface and callers type-assert to the concrete identifier when they need more.
type Identifier interface {
	// AsXmlID returns an ID conformant to XML Id: the identifier prefix followed by the
	// hex value of the digest.
	AsXmlID() string
	// DigestID returns the Digest the identifier is built from.
	DigestID() Digest
	// Equals reports whether the two identifiers have the same concrete type and digest.
	Equals(other Identifier) bool
	// String returns the Java toString() form, "<SimpleClassName>:<digest>".
	String() string
}

// IdentifierBase carries the state and behaviour of the abstract Java class Identifier.
// Every identifier embeds it; the Java class name each subclass reports from
// getClass().getSimpleName() is passed to the constructor because Java's equals() compares
// classes and toString() prints the simple name.
type IdentifierBase struct {
	// className is the Java simple class name of the concrete identifier.
	className string
	// prefix is prepended to the hex value, e.g. "C-" + HEX.
	prefix string
	// id is the digest identifier.
	id Digest
	// xmlID caches the value returned by AsXmlID.
	xmlID string
}

// NewIdentifierBase computes an identifier from the binaries with a defined prefix.
// Port of the protected Identifier(String, byte[]) constructor.
//
// Panics with the Java message when data is missing (Objects.requireNonNull("Data binaries
// cannot be null!")). Java's companion check Objects.requireNonNull(prefix, "Prefix cannot
// be null!") has no counterpart: prefix is a Go string, and every caller in dss-model
// passes a constant literal.
func NewIdentifierBase(className, prefix string, data []byte) IdentifierBase {
	if data == nil {
		panic("Data binaries cannot be null!")
	}
	digest, err := identifierMessageDigest(IdentifierDigestAlgorithm)
	if err != nil {
		// SHA-256 is always available in the Go standard library, so this branch is
		// unreachable; upstream would raise a DSSException here.
		panic(err)
	}
	digest.Write(data)
	return IdentifierBase{
		className: className,
		prefix:    prefix,
		id:        NewDigest(IdentifierDigestAlgorithm, digest.Sum(nil)),
	}
}

// NewIdentifierBaseFromDigest builds an identifier over an already computed digest with a
// defined prefix. Port of the protected Identifier(String, Digest) constructor.
//
// Digest is a Go value type, so the Objects.requireNonNull("Digest cannot be null!") check
// upstream performs has no counterpart; the zero Digest is the port of Java's empty
// Digest(), not of null.
func NewIdentifierBaseFromDigest(className, prefix string, digest Digest) IdentifierBase {
	return IdentifierBase{className: className, prefix: prefix, id: digest}
}

// MessageDigest returns a hash for the given DigestAlgorithm. Port of the protected
// Identifier#getMessageDigest(DigestAlgorithm); Java's DSSException on an unavailable
// algorithm becomes the returned error.
func (i *IdentifierBase) MessageDigest(digestAlgorithm enumerations.DigestAlgorithm) (hash.Hash, error) {
	return identifierMessageDigest(digestAlgorithm)
}

// DigestID returns the Digest Id. Port of the package-private Identifier#getDigestId().
func (i *IdentifierBase) DigestID() Digest {
	return i.id
}

// Prefix returns the identifier prefix, e.g. "C-" for a CertificateTokenIdentifier.
// Java keeps this field private; it is exposed here because it is part of the identifier's
// documented, report-visible string form.
func (i *IdentifierBase) Prefix() string {
	return i.prefix
}

// AsXmlID returns an ID conformant to XML Id. Port of Identifier#asXmlId().
func (i *IdentifierBase) AsXmlID() string {
	if i.xmlID == "" {
		i.xmlID = i.prefix + i.id.HexValue()
	}
	return i.xmlID
}

// String returns "<SimpleClassName>:<digest>". Port of Identifier#toString().
func (i *IdentifierBase) String() string {
	return i.className + ":" + i.id.String()
}

// Equals reports whether both identifiers are of the same Java class and carry an equal
// digest. Port of Identifier#equals(Object), including its getClass() check - a
// DataIdentifier never equals an EntityIdentifier even when the digests match.
func (i *IdentifierBase) Equals(other Identifier) bool {
	if other == nil {
		return false
	}
	otherBase := identifierBaseOf(other)
	if otherBase == nil {
		return false
	}
	if i == otherBase {
		return true
	}
	if i.className != otherBase.className {
		return false
	}
	return i.id.Equals(otherBase.id)
}

// identifierBaseOf extracts the embedded IdentifierBase of any identifier, which is how the
// Go port reaches the fields Java's Identifier#equals compares directly.
func identifierBaseOf(identifier Identifier) *IdentifierBase {
	type identifierBaseHolder interface{ identifierBase() *IdentifierBase }
	if holder, ok := identifier.(identifierBaseHolder); ok {
		return holder.identifierBase()
	}
	return nil
}

// identifierBase implements the unexported accessor identifierBaseOf looks for; it is
// promoted to every type embedding IdentifierBase.
func (i *IdentifierBase) identifierBase() *IdentifierBase {
	return i
}

// identifierMessageDigest maps a DigestAlgorithm onto a Go hash, standing in for
// DigestAlgorithm#getMessageDigest(), which the dss-enumerations port deliberately omits.
//
// DEVIATION: upstream reaches every algorithm through a JCE provider (BouncyCastle in
// practice). MD2, WHIRLPOOL and the SHAKE extendable-output functions have no Go
// equivalent and are reported as unavailable; dss-spi is expected to widen this once it
// brings in a full digest layer. Every other algorithm must stay in step with the table
// CommonDocument uses, so that a digest an algorithm can produce for a document can also
// be produced for a token or an identifier, as it can in Java.
func identifierMessageDigest(digestAlgorithm enumerations.DigestAlgorithm) (hash.Hash, error) {
	switch digestAlgorithm {
	case enumerations.DigestAlgorithm_RIPEMD160:
		return ripemd160.New(), nil
	case enumerations.DigestAlgorithm_SHA1:
		return sha1.New(), nil
	case enumerations.DigestAlgorithm_SHA224:
		return sha256.New224(), nil
	case enumerations.DigestAlgorithm_SHA256:
		return sha256.New(), nil
	case enumerations.DigestAlgorithm_SHA384:
		return sha512.New384(), nil
	case enumerations.DigestAlgorithm_SHA512:
		return sha512.New(), nil
	case enumerations.DigestAlgorithm_SHA3_224:
		return sha3.New224(), nil
	case enumerations.DigestAlgorithm_SHA3_256:
		return sha3.New256(), nil
	case enumerations.DigestAlgorithm_SHA3_384:
		return sha3.New384(), nil
	case enumerations.DigestAlgorithm_SHA3_512:
		return sha3.New512(), nil
	case enumerations.DigestAlgorithm_MD5:
		return md5.New(), nil
	}
	// Upstream chains the JCA NoSuchAlgorithmException as the DSSException cause.
	return nil, NewDSSErrorMessageCause(
		fmt.Sprintf("Unable to create a MessageDigest for algorithm %s", digestAlgorithm),
		fmt.Errorf("NoSuchAlgorithmException: %s", digestAlgorithm.JavaName()))
}
