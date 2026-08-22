package spi

import (
	"bytes"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// dssMessageDigestCalculatorTestVectors holds the digest of "abc" for every algorithm the Go
// port supports, as produced by BouncyCastle 1.78.1 / OpenJDK 21 through
// MessageDigest.getInstance(DigestAlgorithm#getJavaName()).
var dssMessageDigestCalculatorTestVectors = map[enumerations.DigestAlgorithm]string{
	enumerations.DigestAlgorithmSHA1:      "a9993e364706816aba3e25717850c26c9cd0d89d",
	enumerations.DigestAlgorithmSHA224:    "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7",
	enumerations.DigestAlgorithmSHA256:    "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
	enumerations.DigestAlgorithmSHA384:    "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7",
	enumerations.DigestAlgorithmSHA512:    "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
	enumerations.DigestAlgorithmSHA3224:   "e642824c3f8cf24ad09234ee7d3c766fc9a3a5168d0c94ad73b46fdf",
	enumerations.DigestAlgorithmSHA3256:   "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532",
	enumerations.DigestAlgorithmSHA3384:   "ec01498288516fc926459f58e2c6ad8df9b473cb0fc08c2596da7cf0e49be4b298d88cea927ac7f539f1edf228376d25",
	enumerations.DigestAlgorithmSHA3512:   "b751850b1a57168a5693cd924b6b096e08f621827444f70d884f5d0240d2712e10e116e9192af3c91a7ec57647e3934057340b4cf408d5a56592f8274eec53f0",
	enumerations.DigestAlgorithmRIPEMD160: "8eb208f7e05d987a9b044a8e98c6b087f15a0bfc",
	enumerations.DigestAlgorithmMD5:       "900150983cd24fb0d6963f7d28e17f72",
	// BouncyCastle registers "SHAKE256-512" as SHAKE-256 squeezed to 512 bits.
	enumerations.DigestAlgorithmSHAKE256512: "483366601360a8771c6863080cc4114d8db44530f8f1e1ee4f94ea37e78b5739d5a15bef186a5386c75744c0527e1faa9f8726e462a12a4feb06bd8801e751e4",
}

// TestDSSMessageDigestCalculatorKnownAnswers checks every supported algorithm against the
// JCA digest of "abc".
func TestDSSMessageDigestCalculatorKnownAnswers(t *testing.T) {
	for digestAlgorithm, want := range dssMessageDigestCalculatorTestVectors {
		calculator, err := NewDSSMessageDigestCalculator(digestAlgorithm)
		if err != nil {
			t.Fatalf("%s: %v", digestAlgorithm, err)
		}
		calculator.Update([]byte("abc"))
		messageDigest := calculator.MessageDigest(digestAlgorithm)
		if messageDigest.Algorithm() != digestAlgorithm {
			t.Errorf("%s: the message digest reports %s", digestAlgorithm, messageDigest.Algorithm())
		}
		if got := hex.EncodeToString(messageDigest.Value()); got != want {
			t.Errorf("%s:\n got %s\nwant %s", digestAlgorithm, got, want)
		}
	}
}

// TestDSSMessageDigestCalculatorUnsupportedAlgorithms records which algorithms have no
// implementation. SHAKE-128 and SHAKE-256 are unavailable upstream too - BouncyCastle
// registers no JCA MessageDigest under those names - whereas MD2 and WHIRLPOOL are a
// documented Go-side gap.
func TestDSSMessageDigestCalculatorUnsupportedAlgorithms(t *testing.T) {
	for _, digestAlgorithm := range []enumerations.DigestAlgorithm{
		enumerations.DigestAlgorithmSHAKE128,
		enumerations.DigestAlgorithmSHAKE256,
		enumerations.DigestAlgorithmMD2,
		enumerations.DigestAlgorithmWHIRLPOOL,
	} {
		if _, err := NewDSSMessageDigestCalculator(digestAlgorithm); err == nil {
			t.Errorf("%s: expected an error", digestAlgorithm)
		}
	}
	// Every other algorithm must be constructible.
	for _, digestAlgorithm := range enumerations.DigestAlgorithmValues() {
		if _, ok := dssMessageDigestCalculatorTestVectors[digestAlgorithm]; !ok {
			continue
		}
		if _, err := NewDSSMessageDigestCalculator(digestAlgorithm); err != nil {
			t.Errorf("%s: %v", digestAlgorithm, err)
		}
	}
}

// TestDSSMessageDigestCalculatorUpdateVariants checks that the byte, array, range and reader
// update paths agree, and that reading a digest resets the state.
func TestDSSMessageDigestCalculatorUpdateVariants(t *testing.T) {
	want := dssMessageDigestCalculatorTestVectors[enumerations.DigestAlgorithmSHA256]

	byteWise, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range []byte("abc") {
		byteWise.UpdateByte(b)
	}
	if got := hex.EncodeToString(byteWise.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("byte-wise: got %s", got)
	}
	// getMessageDigest resets, so the same calculator restarts from an empty digest.
	byteWise.Update([]byte("abc"))
	if got := hex.EncodeToString(byteWise.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("after reset: got %s", got)
	}

	ranged, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	ranged.UpdateRange([]byte("xxabcxx"), 2, 3)
	if got := hex.EncodeToString(ranged.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("ranged: got %s", got)
	}

	fromReader, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if err := fromReader.UpdateReader(strings.NewReader("abc")); err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(fromReader.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("from reader: got %s", got)
	}
	// A nil input is ignored, as a null array or stream is upstream.
	fromReader.Update(nil)
	if err := fromReader.UpdateReader(nil); err != nil {
		t.Fatal(err)
	}
	empty, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	emptyDigest := hex.EncodeToString(empty.MessageDigest(enumerations.DigestAlgorithmSHA256).Value())
	if got := hex.EncodeToString(fromReader.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != emptyDigest {
		t.Errorf("nil updates must not change the digest: got %s", got)
	}

	// A reader is closed once consumed, matching the try-with-resources upstream.
	closing := &dssMessageDigestCalculatorTestReader{Reader: strings.NewReader("abc")}
	if err := fromReader.UpdateReader(closing); err != nil {
		t.Fatal(err)
	}
	if !closing.closed {
		t.Error("UpdateReader must close a reader that is also an io.Closer")
	}
}

// dssMessageDigestCalculatorTestReader records whether it was closed.
type dssMessageDigestCalculatorTestReader struct {
	io.Reader
	closed bool
}

func (r *dssMessageDigestCalculatorTestReader) Close() error {
	r.closed = true
	return nil
}

// TestDSSMessageDigestCalculatorMultipleAlgorithms checks that one calculator can serve
// several algorithms at once, as the CAdES archive-timestamp code needs.
func TestDSSMessageDigestCalculatorMultipleAlgorithms(t *testing.T) {
	algorithms := []enumerations.DigestAlgorithm{
		enumerations.DigestAlgorithmSHA512,
		enumerations.DigestAlgorithmSHA1,
		enumerations.DigestAlgorithmSHA256,
	}
	calculator, err := NewDSSMessageDigestCalculatorForAlgorithms(algorithms)
	if err != nil {
		t.Fatal(err)
	}
	calculator.Update([]byte("abc"))
	for _, digestAlgorithm := range algorithms {
		got := hex.EncodeToString(calculator.MessageDigest(digestAlgorithm).Value())
		if want := dssMessageDigestCalculatorTestVectors[digestAlgorithm]; got != want {
			t.Errorf("%s:\n got %s\nwant %s", digestAlgorithm, got, want)
		}
	}
	// Asking for an algorithm that was not part of the computation is a programmer error.
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Error("expected a panic for an unused algorithm")
		}
	}()
	calculator.MessageDigest(enumerations.DigestAlgorithmMD5)
}

// TestDSSMessageDigestCalculatorConstructorGuards checks the requireNonNull-style panics.
func TestDSSMessageDigestCalculatorConstructorGuards(t *testing.T) {
	for _, entry := range []struct {
		name       string
		algorithms []enumerations.DigestAlgorithm
	}{
		{"nil", nil},
		{"empty", []enumerations.DigestAlgorithm{}},
		{"nil element", []enumerations.DigestAlgorithm{enumerations.DigestAlgorithmSHA256, ""}},
	} {
		func() {
			defer func() {
				if recovered := recover(); recovered == nil {
					t.Errorf("%s: expected a panic", entry.name)
				}
			}()
			_, _ = NewDSSMessageDigestCalculatorForAlgorithms(entry.algorithms)
		}()
	}
}

// TestDSSMessageDigestCalculatorWriter checks the on-the-fly digest writers.
func TestDSSMessageDigestCalculatorWriter(t *testing.T) {
	want := dssMessageDigestCalculatorTestVectors[enumerations.DigestAlgorithmSHA256]

	discarding, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	writer := discarding.Writer()
	if _, err := writer.Write([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(discarding.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("Writer: got %s", got)
	}

	mirroring, err := NewDSSMessageDigestCalculator(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	writer = mirroring.WriterFor(&buffer)
	if _, err := writer.Write([]byte("a")); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("bc")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if buffer.String() != "abc" {
		t.Errorf("the wrapped writer received %q", buffer.String())
	}
	if got := hex.EncodeToString(mirroring.MessageDigest(enumerations.DigestAlgorithmSHA256).Value()); got != want {
		t.Errorf("WriterFor: got %s", got)
	}
}
