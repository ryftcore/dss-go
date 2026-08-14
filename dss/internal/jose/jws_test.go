package jose

import (
	"bufio"
	"crypto/x509"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

// jwsCase is one row of testdata/jws_oracle.tsv: a JWS that jose4j really signed, together with
// the public key and the signing input it computed.
type jwsCase struct {
	name         string
	alg          string
	publicKeyDER []byte
	compact      string
	signingInput []byte
}

func loadJWSOracle(t *testing.T) []jwsCase {
	t.Helper()
	f, err := os.Open("testdata/jws_oracle.tsv")
	if err != nil {
		t.Fatalf("cannot open the signed-JWS oracle: %v", err)
	}
	defer f.Close()

	var cases []jwsCase
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 6 || cols[0] != "JWS" {
			t.Fatalf("line %d: malformed row", line)
		}
		der, err := hex.DecodeString(cols[3])
		if err != nil {
			t.Fatalf("line %d: public key is not hex: %v", line, err)
		}
		input, err := hex.DecodeString(cols[5])
		if err != nil {
			t.Fatalf("line %d: signing input is not hex: %v", line, err)
		}
		cases = append(cases, jwsCase{
			name:         cols[1],
			alg:          cols[2],
			publicKeyDER: der,
			compact:      cols[4],
			signingInput: input,
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading the signed-JWS oracle: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("the signed-JWS oracle is empty")
	}
	return cases
}

// TestKnownAnswersSignedJWS checks the whole chain - compact deserialization, header decoding,
// signing-input construction and signature verification - against signatures jose4j actually
// produced, for every algorithm family and for the RFC 7797 unencoded-payload form.
//
// The signing-input comparison is the part that matters most: a Go implementation that computed
// a different input would still "verify" nothing, but one that computed the right input and
// the wrong verdict is a security bug, so both are asserted separately.
func TestKnownAnswersSignedJWS(t *testing.T) {
	for _, c := range loadJWSOracle(t) {
		t.Run(c.name, func(t *testing.T) {
			publicKey, err := x509.ParsePKIXPublicKey(c.publicKeyDER)
			if err != nil {
				t.Fatalf("parsing the public key: %v", err)
			}

			jws := NewJWS()
			// The DSS subclass disables the crit check; do the same so the fixtures that carry
			// a 'crit' header exercise the same path dss-jades does.
			jws.CheckCritOverride = func(*JWS) error { return nil }
			parts := CompactDeserialize(c.compact)
			if len(parts) != CompactSerializationParts {
				t.Fatalf("got %d compact parts, want %d", len(parts), CompactSerializationParts)
			}
			if err := jws.SetCompactSerializationParts(parts); err != nil {
				t.Fatalf("SetCompactSerializationParts: %v", err)
			}

			if got := jws.AlgorithmHeaderValue(); got != c.alg {
				t.Errorf("alg header = %q, want %q", got, c.alg)
			}
			if got := jws.EncodedHeader(); got != parts[0] {
				t.Errorf("EncodedHeader = %q, want the original %q", got, parts[0])
			}
			if got := jws.SigningInputBytes(); string(got) != string(c.signingInput) {
				t.Errorf("signing input\n got %q\nwant %q", got, c.signingInput)
			}

			jws.SetKey(publicKey)
			jws.SetDoKeyValidation(false)
			valid, err := jws.VerifySignature()
			if err != nil {
				t.Fatalf("VerifySignature: %v", err)
			}
			if !valid {
				t.Fatal("a signature jose4j produced did not verify")
			}

			// A single flipped bit in the signature must be rejected, and rejected as "false"
			// rather than as an error.
			tampered := NewJWS()
			tampered.CheckCritOverride = func(*JWS) error { return nil }
			if err := tampered.SetCompactSerializationParts(parts); err != nil {
				t.Fatalf("SetCompactSerializationParts: %v", err)
			}
			signature := append([]byte(nil), tampered.Signature()...)
			signature[len(signature)-1] ^= 0x01
			tampered.SetSignature(signature)
			tampered.SetKey(publicKey)
			tampered.SetDoKeyValidation(false)
			valid, err = tampered.VerifySignature()
			if err != nil {
				t.Fatalf("VerifySignature on a tampered signature: %v", err)
			}
			if valid {
				t.Error("a tampered signature verified")
			}

			// So must a flipped bit in the payload.
			modified := NewJWS()
			modified.CheckCritOverride = func(*JWS) error { return nil }
			if err := modified.SetCompactSerializationParts(parts); err != nil {
				t.Fatalf("SetCompactSerializationParts: %v", err)
			}
			payload := append([]byte(nil), modified.UnverifiedPayloadBytes()...)
			if len(payload) > 0 {
				payload[0] ^= 0x01
				modified.SetPayload(string(payload))
				modified.SetKey(publicKey)
				modified.SetDoKeyValidation(false)
				valid, err = modified.VerifySignature()
				if err != nil {
					t.Fatalf("VerifySignature on a modified payload: %v", err)
				}
				if valid {
					t.Error("a modified payload verified")
				}
			}
		})
	}
}

// TestUnencodedPayloadSigningInput states the RFC 7797 rule on its own, without a signature in
// the way, because it is the single place where the payload reaches the signing input as raw
// bytes instead of as base64url text.
func TestUnencodedPayloadSigningInput(t *testing.T) {
	// Bytes that are not valid UTF-8: routing them through a string would replace them with
	// U+FFFD and change what gets signed.
	payload := []byte{0xff, 0xfe, 0x00, 0x41}

	jws := NewJWS()
	jws.SetHeader(HeaderAlgorithm, AlgorithmRS256)
	jws.SetHeader(HeaderBase64URLEncodePayload, false)
	jws.SetPayloadBytes(payload)

	if !jws.IsRfc7797UnencodedPayload() {
		t.Fatal("b64=false was not recognised")
	}
	want := append(append([]byte(nil), ASCIIBytes(jws.EncodedHeader())...), '.')
	want = append(want, payload...)
	if got := jws.SigningInputBytes(); string(got) != string(want) {
		t.Errorf("signing input\n got %x\nwant %x", got, want)
	}
}

// TestIsRfc7797UnencodedPayloadIsStrict pins the narrow reading of "b64" that upstream uses:
// only a JSON boolean false counts.
func TestIsRfc7797UnencodedPayloadIsStrict(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"absent", nil, false},
		{"boolean false", false, true},
		{"boolean true", true, false},
		{"string false", "false", false},
		{"number zero", NewLong(0), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jws := NewJWS()
			if tc.value != nil {
				jws.SetHeader(HeaderBase64URLEncodePayload, tc.value)
			}
			if got := jws.IsRfc7797UnencodedPayload(); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestEncodedHeaderIsPreservedAcrossParsing is the caching contract of Headers: a header that
// arrived base64url-encoded must be handed back byte for byte, because re-serializing its
// members could produce a different member order or different escaping - and therefore a
// signing input that was never signed.
func TestEncodedHeaderIsPreservedAcrossParsing(t *testing.T) {
	// Members deliberately out of the order a fresh serialization would produce, with a space
	// after the colon that a re-serialization would drop.
	original := `{"typ": "JOSE", "alg":"RS256"}`
	encoded := Base64URLEncodeUTF8(original)

	jws := NewJWS()
	if err := jws.SetEncodedHeader(encoded); err != nil {
		t.Fatalf("SetEncodedHeader: %v", err)
	}
	if got := jws.EncodedHeader(); got != encoded {
		t.Errorf("EncodedHeader = %q, want the original %q", got, encoded)
	}
	if got := jws.Headers().FullHeaderAsJSONString(); got != original {
		t.Errorf("FullHeaderAsJSONString = %q, want the original %q", got, original)
	}
	if got := jws.AlgorithmHeaderValue(); got != AlgorithmRS256 {
		t.Errorf("alg = %q, want %q", got, AlgorithmRS256)
	}

	// Writing a member invalidates the cache, so from then on the header is rendered fresh.
	jws.SetHeader(HeaderKeyID, "k1")
	if got := jws.Headers().FullHeaderAsJSONString(); got != `{"typ":"JOSE","alg":"RS256","kid":"k1"}` {
		t.Errorf("after a write, FullHeaderAsJSONString = %q", got)
	}
}

// TestSetEncodedHeaderRejectsEmpty covers checkNotEmptyPart, which is what turns a "..signature"
// document into a parse failure rather than a signature over an empty header.
func TestSetEncodedHeaderRejectsEmpty(t *testing.T) {
	if err := NewJWS().SetEncodedHeader(""); err == nil {
		t.Fatal("an empty encoded header was accepted")
	}
}

// TestSetCompactSerializationPartsRejectsWrongCount covers the three-part precondition.
func TestSetCompactSerializationPartsRejectsWrongCount(t *testing.T) {
	for _, parts := range [][]string{{}, {"a"}, {"a", "b"}, {"a", "b", "c", "d"}} {
		if err := NewJWS().SetCompactSerializationParts(parts); err == nil {
			t.Errorf("%d parts were accepted", len(parts))
		}
	}
}

// TestCheckCritOverride is the virtual-dispatch guard: without the hook a Go base type would
// keep running its own checkCrit and reject every JAdES signature whose 'crit' names a JAdES
// header, which is exactly what the DSS subclass overrides away.
func TestCheckCritOverride(t *testing.T) {
	jws := NewJWS()
	jws.SetHeader(HeaderCritical, []any{"sigT"})

	if err := jws.CheckCrit(); err == nil {
		t.Fatal("an unknown critical header was accepted by the default check")
	}

	jws.SetKnownCriticalHeaders([]string{"sigT"})
	if err := jws.CheckCrit(); err != nil {
		t.Fatalf("a declared critical header was rejected: %v", err)
	}

	other := NewJWS()
	other.SetHeader(HeaderCritical, []any{"sigT"})
	other.CheckCritOverride = func(*JWS) error { return nil }
	if err := other.CheckCrit(); err != nil {
		t.Fatalf("the override did not take effect: %v", err)
	}
}

// TestVerifySignatureRejectsNoneAlgorithm covers the DISALLOW_NONE constraint that
// JsonWebSignature's constructor installs and DSS never relaxes.
func TestVerifySignatureRejectsNoneAlgorithm(t *testing.T) {
	jws := NewJWS()
	jws.SetHeader(HeaderAlgorithm, AlgorithmNone)
	jws.SetPayload("x")
	jws.SetSignature(nil)
	if _, err := jws.VerifySignature(); err == nil {
		t.Fatal("the \"none\" algorithm was accepted")
	}
}

// TestSetPayloadBytesKeepsStaleEncodedPayload documents a jose4j quirk rather than endorsing it:
// setPayloadBytes does not clear the cached encoded payload, and JWSJsonSerializationParser's
// call order depends on that not having been "fixed".
func TestSetPayloadBytesKeepsStaleEncodedPayload(t *testing.T) {
	jws := NewJWS()
	jws.SetEncodedPayload(Base64URLEncode([]byte("first")))
	jws.SetPayloadBytes([]byte("second"))

	if got := jws.UnverifiedPayload(); got != "second" {
		t.Errorf("payload = %q, want %q", got, "second")
	}
	if got := jws.EncodedPayload(); got != Base64URLEncode([]byte("first")) {
		t.Errorf("EncodedPayload = %q, want the stale %q", got, Base64URLEncode([]byte("first")))
	}
	// setPayload, by contrast, does clear it.
	jws.SetPayload("third")
	if got := jws.EncodedPayload(); got != Base64URLEncode([]byte("third")) {
		t.Errorf("after SetPayload, EncodedPayload = %q", got)
	}
}
