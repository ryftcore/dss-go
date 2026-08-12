// Ported from dss-model/.../CommonDocument.java (DSS 6.5.RC1).
package model

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/utain/esig/dss/enumerations"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// CommonDocument is embedded by concrete DSSDocument implementations
// (InMemoryDocument, FileDocument, DigestDocument) and provides ports of
// DSSDocument's Java default methods (save, writeTo, getMimeType,
// setMimeType, getName, setName, getDigest, getDigestValue).
//
// OpenStream is NOT implemented here (it is abstract in Java); embedders
// must implement it. Because Go embedding has no virtual dispatch, the
// default methods that need OpenStream (WriteTo, Save, Digest, DigestValue)
// are exposed here as free functions taking the owning DSSDocument, and
// each concrete type wires up a trivial forwarding method.
type CommonDocument struct {
	// digestMap caches previously computed digests, keyed by algorithm.
	digestMap map[enumerations.DigestAlgorithm][]byte

	// mimeType is the MimeType of the document.
	mimeType enumerations.MimeType

	// name is the document name.
	name string
}

// Name returns the document name. Ports CommonDocument#getName.
func (c *CommonDocument) Name() string { return c.name }

// SetName sets the document name. Ports CommonDocument#setName.
func (c *CommonDocument) SetName(name string) { c.name = name }

// MimeType returns the document's MimeType. Ports CommonDocument#getMimeType.
func (c *CommonDocument) MimeType() enumerations.MimeType { return c.mimeType }

// SetMimeType sets the document's MimeType. Ports CommonDocument#setMimeType.
func (c *CommonDocument) SetMimeType(mimeType enumerations.MimeType) { c.mimeType = mimeType }

// String ports CommonDocument#toString.
func (c *CommonDocument) String() string {
	mimeTypeString := ""
	if c.mimeType != nil {
		mimeTypeString = c.mimeType.MimeTypeString()
	}
	return "Name: " + c.name + " / MimeType: " + mimeTypeString
}

// commonDocumentEquals ports CommonDocument#equals: compares the mimeType
// and name fields shared by every CommonDocument-embedding type.
func commonDocumentEquals(a, b *CommonDocument) bool {
	return a.mimeType == b.mimeType && a.name == b.name
}

// commonDocumentSave ports CommonDocument#save.
func commonDocumentSave(doc DSSDocument, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = doc.WriteTo(f)
	return err
}

// commonDocumentWriteTo ports CommonDocument#writeTo.
func commonDocumentWriteTo(doc DSSDocument, w io.Writer) (int64, error) {
	rc, err := doc.OpenStream()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

// commonDocumentDigest ports CommonDocument#getDigest.
func commonDocumentDigest(doc DSSDocument, c *CommonDocument, digestAlgorithm enumerations.DigestAlgorithm) (Digest, error) {
	digestBytes, err := commonDocumentDigestValue(doc, c, digestAlgorithm)
	if err != nil {
		return Digest{}, err
	}
	return NewDigest(digestAlgorithm, digestBytes), nil
}

// commonDocumentDigestValue ports CommonDocument#getDigestValue.
func commonDocumentDigestValue(doc DSSDocument, c *CommonDocument, digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	if c.digestMap == nil {
		c.digestMap = make(map[enumerations.DigestAlgorithm][]byte)
	}
	if digest, ok := c.digestMap[digestAlgorithm]; ok {
		return digest, nil
	}
	h, err := commonDocumentHash(digestAlgorithm)
	if err != nil {
		return nil, &DSSError{Message: "Unable to compute the digest", Cause: err}
	}
	rc, err := doc.OpenStream()
	if err != nil {
		return nil, &DSSError{Message: "Unable to compute the digest", Cause: err}
	}
	defer rc.Close()
	if _, err := io.Copy(h, rc); err != nil {
		return nil, &DSSError{Message: "Unable to compute the digest", Cause: err}
	}
	digest := h.Sum(nil)
	c.digestMap[digestAlgorithm] = digest
	return digest, nil
}

// commonDocumentHash resolves the stdlib/x-crypto hash.Hash for
// digestAlgorithm, mirroring java.security.MessageDigest.getInstance keyed
// off DigestAlgorithm.JavaName(). Algorithms with no available Go
// implementation (MD2, WHIRLPOOL, SHAKE128, SHAKE256, SHAKE256_512) return
// an error, matching Java's NoSuchAlgorithmException path.
func commonDocumentHash(digestAlgorithm enumerations.DigestAlgorithm) (hash.Hash, error) {
	switch digestAlgorithm {
	case enumerations.DigestAlgorithm_MD5:
		return md5.New(), nil
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
	case enumerations.DigestAlgorithm_RIPEMD160:
		return ripemd160.New(), nil
	default:
		return nil, fmt.Errorf("no such algorithm: %s", digestAlgorithm.JavaName())
	}
}
