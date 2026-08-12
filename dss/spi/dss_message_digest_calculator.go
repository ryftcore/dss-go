// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSMessageDigestCalculator.java (DSS 6.5.RC1).
package spi

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// DSSMessageDigestCalculator computes model.DSSMessageDigest values over the provided
// input, for one or several digest algorithms at once.
type DSSMessageDigestCalculator struct {
	// messageDigestMap is the map of DigestAlgorithm and corresponding message-digest.
	messageDigestMap map[enumerations.DigestAlgorithm]hash.Hash
	// algorithms keeps the algorithms in the order Java's EnumMap iterates them, i.e. in
	// DigestAlgorithm declaration order, so that updates are applied deterministically.
	algorithms []enumerations.DigestAlgorithm
}

// NewDSSMessageDigestCalculator creates a calculator for a single digest algorithm.
// Port of DSSMessageDigestCalculator(DigestAlgorithm).
func NewDSSMessageDigestCalculator(digestAlgorithm enumerations.DigestAlgorithm) (*DSSMessageDigestCalculator, error) {
	return NewDSSMessageDigestCalculatorForAlgorithms([]enumerations.DigestAlgorithm{digestAlgorithm})
}

// NewDSSMessageDigestCalculatorForAlgorithms creates a calculator for several digest
// algorithms. Port of DSSMessageDigestCalculator(Collection<DigestAlgorithm>).
//
// Panics with the Java message when the collection is nil ("DigestAlgorithms shall be
// defined!"), empty ("DigestAlgorithms collection cannot be empty!") or holds an empty
// algorithm ("DigestAlgorithm cannot be null!"): those are Objects.requireNonNull /
// IllegalArgumentException programmer errors upstream. An algorithm with no available
// implementation is data-dependent and is returned as an error, matching the DSSException
// upstream raises for a JCA NoSuchAlgorithmException.
func NewDSSMessageDigestCalculatorForAlgorithms(digestAlgorithms []enumerations.DigestAlgorithm) (*DSSMessageDigestCalculator, error) {
	if digestAlgorithms == nil {
		panic("DigestAlgorithms shall be defined!")
	}
	if len(digestAlgorithms) == 0 {
		panic("DigestAlgorithms collection cannot be empty!")
	}
	calculator := &DSSMessageDigestCalculator{
		messageDigestMap: make(map[enumerations.DigestAlgorithm]hash.Hash, len(digestAlgorithms)),
	}
	for _, digestAlgorithm := range digestAlgorithms {
		if digestAlgorithm == "" {
			panic("DigestAlgorithm cannot be null!")
		}
		messageDigest, err := dssMessageDigestCalculatorMessageDigest(digestAlgorithm)
		if err != nil {
			return nil, err
		}
		calculator.messageDigestMap[digestAlgorithm] = messageDigest
	}
	// EnumMap iterates in enum declaration order; reproduce it for a deterministic order.
	for _, digestAlgorithm := range enumerations.DigestAlgorithmValues() {
		if _, ok := calculator.messageDigestMap[digestAlgorithm]; ok {
			calculator.algorithms = append(calculator.algorithms, digestAlgorithm)
		}
	}
	return calculator, nil
}

// UpdateByte updates the digests with the provided byte. Port of update(byte).
func (c *DSSMessageDigestCalculator) UpdateByte(byteToAdd byte) {
	c.UpdateRange([]byte{byteToAdd}, 0, 1)
}

// Update updates the digests with the provided array of bytes. Port of update(byte[]);
// a nil slice is ignored, as a null array is upstream.
func (c *DSSMessageDigestCalculator) Update(bytes []byte) {
	if bytes != nil {
		c.UpdateRange(bytes, 0, len(bytes))
	}
}

// UpdateRange updates the digests with length bytes starting at offset.
// Port of update(byte[], int, int); a nil slice is ignored.
func (c *DSSMessageDigestCalculator) UpdateRange(bytes []byte, offset, length int) {
	if bytes == nil {
		return
	}
	chunk := bytes[offset : offset+length]
	for _, digestAlgorithm := range c.algorithms {
		// hash.Hash never reports an error from Write.
		_, _ = c.messageDigestMap[digestAlgorithm].Write(chunk)
	}
}

// UpdateReader updates the digests by reading the provided reader to EOF. When the reader
// is also an io.Closer it is closed afterwards, matching the try-with-resources upstream.
// Port of update(InputStream); a nil reader is ignored.
func (c *DSSMessageDigestCalculator) UpdateReader(reader io.Reader) error {
	if reader == nil {
		return nil
	}
	if closer, ok := reader.(io.Closer); ok {
		defer func() { _ = closer.Close() }()
	}
	buffer := make([]byte, 4096)
	for {
		count, err := reader.Read(buffer)
		if count > 0 {
			c.UpdateRange(buffer, 0, count)
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// MessageDigest returns the model.DSSMessageDigest for the given algorithm and resets the
// state of that message-digest. Port of getMessageDigest(DigestAlgorithm).
//
// Panics with the Java message when the algorithm was not part of the computation
// (IllegalArgumentException upstream).
func (c *DSSMessageDigestCalculator) MessageDigest(digestAlgorithm enumerations.DigestAlgorithm) model.DSSMessageDigest {
	messageDigest, ok := c.messageDigestMap[digestAlgorithm]
	if !ok {
		panic("The DigestAlgorithm was not used on message-digest computation!")
	}
	// java.security.MessageDigest#digest() returns the digest and resets; hash.Hash#Sum
	// does not reset, so the two steps are explicit here.
	value := messageDigest.Sum(nil)
	messageDigest.Reset()
	return model.NewDSSMessageDigestWithValue(digestAlgorithm, value)
}

// Writer returns an io.WriteCloser that updates this calculator as it is written to,
// discarding the bytes. Port of getOutputStream().
func (c *DSSMessageDigestCalculator) Writer() io.WriteCloser {
	return c.WriterFor(io.Discard)
}

// WriterFor returns an io.WriteCloser that writes to the provided writer and updates this
// calculator with the same bytes. Closing it closes the wrapped writer when that writer is
// an io.Closer. Port of getOutputStream(OutputStream).
func (c *DSSMessageDigestCalculator) WriterFor(writer io.Writer) io.WriteCloser {
	return &dssMessageDigestCalculatorWriter{wrapped: writer, calculator: c}
}

// dssMessageDigestCalculatorWriter is the anonymous OutputStream subclass upstream returns
// from getOutputStream.
type dssMessageDigestCalculatorWriter struct {
	wrapped    io.Writer
	calculator *DSSMessageDigestCalculator
}

// Write forwards the bytes to the wrapped writer and updates the calculator. The digest is
// only updated with the bytes the wrapped writer accepted, so that a short write cannot
// silently desynchronise the digest from the output.
func (w *dssMessageDigestCalculatorWriter) Write(p []byte) (int, error) {
	n, err := w.wrapped.Write(p)
	if n > 0 {
		w.calculator.UpdateRange(p, 0, n)
	}
	return n, err
}

// Close closes the wrapped writer when it is an io.Closer.
func (w *dssMessageDigestCalculatorWriter) Close() error {
	if closer, ok := w.wrapped.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// dssMessageDigestCalculatorMessageDigest resolves the hash.Hash of a DigestAlgorithm,
// standing in for DigestAlgorithm#getMessageDigest(), i.e.
// MessageDigest.getInstance(javaName) against the BouncyCastle provider.
//
// The availability matches upstream except for two documented gaps: MD2 and WHIRLPOOL have
// no Go implementation, where BouncyCastle/the JDK provide one. SHAKE-128 and SHAKE-256 are
// unavailable in BOTH ports - BouncyCastle registers no JCA MessageDigest under those names
// either, so upstream raises the same failure.
func dssMessageDigestCalculatorMessageDigest(digestAlgorithm enumerations.DigestAlgorithm) (hash.Hash, error) {
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
	case enumerations.DigestAlgorithm_SHAKE256_512:
		// BouncyCastle registers "SHAKE256-512" as SHAKE-256 squeezed to 512 bits.
		return &dssMessageDigestCalculatorShake{shake: sha3.NewShake256(), size: 64}, nil
	}
	return nil, model.NewDSSErrorMessageCause(
		fmt.Sprintf("Unable to build MessageDigest for the algorithm '%s'", digestAlgorithm.Name()),
		fmt.Errorf("NoSuchAlgorithmException: %s", digestAlgorithm.JavaName()))
}

// dssMessageDigestCalculatorShake adapts a SHAKE extendable-output function to hash.Hash by
// fixing the output length, the way BouncyCastle's SHAKEDigest exposes a fixed digest size.
type dssMessageDigestCalculatorShake struct {
	shake sha3.ShakeHash
	size  int
}

// Write absorbs more input.
func (s *dssMessageDigestCalculatorShake) Write(p []byte) (int, error) { return s.shake.Write(p) }

// Sum appends the fixed-length squeezed output to b, leaving this hash's state untouched.
func (s *dssMessageDigestCalculatorShake) Sum(b []byte) []byte {
	// Squeezing consumes the sponge, so the digest is taken from a clone.
	clone := s.shake.Clone()
	out := make([]byte, s.size)
	_, _ = io.ReadFull(clone, out)
	return append(b, out...)
}

// Reset discards the absorbed input.
func (s *dssMessageDigestCalculatorShake) Reset() { s.shake.Reset() }

// Size returns the fixed output length in bytes.
func (s *dssMessageDigestCalculatorShake) Size() int { return s.size }

// BlockSize returns the sponge rate of SHAKE-256.
func (s *dssMessageDigestCalculatorShake) BlockSize() int { return s.shake.BlockSize() }

// compile-time interface assertion.
var _ hash.Hash = (*dssMessageDigestCalculatorShake)(nil)
