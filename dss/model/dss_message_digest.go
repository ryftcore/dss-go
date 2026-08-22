// Ported from dss-model/.../DSSMessageDigest.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// DSSMessageDigest holds a digest algorithm and digest value for
// message-digest computation.
//
// Digest (the embedded/extended Java class) is outside this manifest;
// assumed to already exist in this package with a NewDigest(algorithm,
// value) constructor and Algorithm()/Value()/String() accessors mirroring
// its Java getAlgorithm()/getValue()/toString().
type DSSMessageDigest struct {
	Digest
}

// NewDSSMessageDigest creates an empty message-digest object. Ports the
// empty constructor.
func NewDSSMessageDigest() DSSMessageDigest {
	return DSSMessageDigest{}
}

// NewDSSMessageDigestWithValue creates a message-digest with the provided
// digest algorithm and the corresponding hash value. Ports
// DSSMessageDigest(DigestAlgorithm, byte[]).
func NewDSSMessageDigestWithValue(algorithm enumerations.DigestAlgorithm, value []byte) DSSMessageDigest {
	return DSSMessageDigest{Digest: NewDigest(algorithm, value)}
}

// NewDSSMessageDigestFromDigest creates a message-digest from the provided
// Digest object. Ports DSSMessageDigest(Digest).
func NewDSSMessageDigestFromDigest(digest Digest) DSSMessageDigest {
	return DSSMessageDigest{Digest: NewDigest(digest.Algorithm(), digest.Value())}
}

// CreateEmptyDSSMessageDigest creates an empty message-digest object with
// empty values. Ports DSSMessageDigest#createEmptyDigest.
func CreateEmptyDSSMessageDigest() DSSMessageDigest {
	return DSSMessageDigest{}
}

// String ports DSSMessageDigest#toString.
func (d DSSMessageDigest) String() string {
	return "MessageDigest [" + d.Digest.String() + "]"
}
