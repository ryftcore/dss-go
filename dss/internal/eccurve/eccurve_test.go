package eccurve

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"
)

// TestBrainpoolCurveConstants re-derives every registered curve's parameters from the curve
// arithmetic itself, so that a typo in any of the numbers in brainpool.go cannot survive. For a
// correct short-Weierstrass group of prime order n with generator G:
//
//   - G lies on the curve;
//   - n*G is the point at infinity, spelled (0, 0);
//   - (n-1)*G == -G, i.e. the same x with y negated mod p;
//   - 2*G == G+G, tying the (a-dependent) doubling formula to the generic addition one, which
//     is the single most likely place for a wrong "a" to hide.
func TestBrainpoolCurveConstants(t *testing.T) {
	oids := KnownCurveOIDs()
	if len(oids) != 14 {
		t.Fatalf("expected the 14 RFC 5639 curves, got %d", len(oids))
	}
	for _, oid := range oids {
		oid := oid
		curve := CurveForOID(oid)
		params := curve.Params()
		t.Run(params.Name, func(t *testing.T) {
			if !curve.IsOnCurve(params.Gx, params.Gy) {
				t.Fatalf("%s: the generator is not on the curve", params.Name)
			}
			if params.BitSize != params.P.BitLen() {
				t.Errorf("%s: BitSize %d but p has %d bits", params.Name, params.BitSize, params.P.BitLen())
			}

			x, y := curve.ScalarBaseMult(params.N.Bytes())
			if x.Sign() != 0 || y.Sign() != 0 {
				t.Errorf("%s: n*G = (%s, %s), want the point at infinity", params.Name, x, y)
			}

			nMinusOne := new(big.Int).Sub(params.N, big.NewInt(1))
			x, y = curve.ScalarBaseMult(nMinusOne.Bytes())
			if x.Cmp(params.Gx) != 0 {
				t.Errorf("%s: (n-1)*G has x = %s, want Gx", params.Name, x)
			}
			negGy := new(big.Int).Sub(params.P, params.Gy)
			if y.Cmp(negGy) != 0 {
				t.Errorf("%s: (n-1)*G has y = %s, want -Gy", params.Name, y)
			}

			doubleX, doubleY := curve.Double(params.Gx, params.Gy)
			addX, addY := curve.Add(params.Gx, params.Gy, params.Gx, params.Gy)
			if doubleX.Cmp(addX) != 0 || doubleY.Cmp(addY) != 0 {
				t.Errorf("%s: Double(G) = (%s, %s) but Add(G, G) = (%s, %s)",
					params.Name, doubleX, doubleY, addX, addY)
			}
			if !curve.IsOnCurve(doubleX, doubleY) {
				t.Errorf("%s: 2*G is not on the curve", params.Name)
			}

			// A scalar with no structure, to exercise the double-and-add path rather than just
			// its endpoints: 7*G computed two independent ways.
			sevenX, sevenY := curve.ScalarBaseMult([]byte{7})
			accX, accY := curve.Double(doubleX, doubleY) // 4G
			accX, accY = curve.Add(accX, accY, doubleX, doubleY)
			accX, accY = curve.Add(accX, accY, params.Gx, params.Gy)
			if sevenX.Cmp(accX) != 0 || sevenY.Cmp(accY) != 0 {
				t.Errorf("%s: 7*G disagrees between ScalarBaseMult and repeated Add/Double", params.Name)
			}
		})
	}
}

// TestBrainpoolPointOffCurveRejected guards the IsOnCurve predicate itself: with a wrong "a" or
// a permissive check, ECDSA verification would accept points that are not group members.
func TestBrainpoolPointOffCurveRejected(t *testing.T) {
	curve := CurveForOID(BrainpoolP256r1OID)
	params := curve.Params()
	offY := new(big.Int).Add(params.Gy, big.NewInt(1))
	if curve.IsOnCurve(params.Gx, offY) {
		t.Error("a point one off the generator was accepted as on-curve")
	}
	if curve.IsOnCurve(new(big.Int).Neg(big.NewInt(1)), params.Gy) {
		t.Error("a negative coordinate was accepted")
	}
	if curve.IsOnCurve(params.P, params.Gy) {
		t.Error("a coordinate equal to p was accepted")
	}
}

// TestCurveForOIDUnknown keeps CurveForOID from inventing curves.
func TestCurveForOIDUnknown(t *testing.T) {
	if c := CurveForOID("1.2.840.10045.3.1.7"); c != nil {
		t.Errorf("CurveForOID answered %v for prime256v1, which crypto/x509 handles itself", c)
	}
	if c := CurveForOID("9.9.9"); c != nil {
		t.Errorf("CurveForOID answered %v for a nonsense OID", c)
	}
}

// TestBrainpoolECDSAVerifyInteropWithStdlib proves the whole point of the package: a signature
// made on a Brainpool key verifies through crypto/ecdsa, which routes custom curves to its
// legacy big.Int path and therefore exercises ScalarBaseMult, ScalarMult and Add together.
func TestBrainpoolECDSAVerifyInteropWithStdlib(t *testing.T) {
	curve := CurveForOID(BrainpoolP256r1OID)
	// A deterministic key: d = 12345, Q = d*G. Signing uses this package's own arithmetic via
	// crypto/ecdsa's legacy path, which is exactly the code under test.
	d := big.NewInt(12345)
	qx, qy := curve.ScalarBaseMult(d.Bytes())
	priv := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: curve, X: qx, Y: qy}, D: d}

	digest := make([]byte, 32)
	for i := range digest {
		digest[i] = byte(i)
	}
	signature, err := ecdsa.SignASN1(newDeterministicReader(), priv, digest)
	if err != nil {
		t.Fatalf("SignASN1 on a Brainpool key: %v", err)
	}
	if !ecdsa.VerifyASN1(&priv.PublicKey, digest, signature) {
		t.Fatal("a Brainpool signature did not verify against its own public key")
	}
	digest[0] ^= 0xff
	if ecdsa.VerifyASN1(&priv.PublicKey, digest, signature) {
		t.Fatal("a Brainpool signature verified against the wrong digest")
	}
}

// TestBrainpoolCurvesAreNotStdlibCurves documents why the package cannot simply delegate:
// crypto/elliptic's four curves are all a = -3, and none of them shares a Brainpool prime.
func TestBrainpoolCurvesAreNotStdlibCurves(t *testing.T) {
	stdlib := []elliptic.Curve{elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()}
	for _, oid := range KnownCurveOIDs() {
		curve := CurveForOID(oid)
		for _, std := range stdlib {
			if curve.Params().P.Cmp(std.Params().P) == 0 {
				t.Errorf("%s shares its prime with %s", curve.Params().Name, std.Params().Name)
			}
		}
	}
}

// newDeterministicReader answers a reader of a fixed byte pattern, so the test never depends on
// system entropy. It is only ever used to pick a nonce for a throwaway test key.
func newDeterministicReader() *deterministicReader { return &deterministicReader{} }

type deterministicReader struct{ n byte }

func (r *deterministicReader) Read(p []byte) (int, error) {
	for i := range p {
		r.n += 0x1f
		p[i] = r.n
	}
	return len(p), nil
}
