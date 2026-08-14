// Package eccurve supplies the elliptic curves crypto/x509 does not recognise, and the
// certificate parser that puts them to use.
//
// WHY THIS EXISTS. crypto/x509 hardcodes exactly four named curves (P-224/256/384/521, see
// namedCurveFromOID in the standard library's x509.go) and fails the WHOLE certificate parse -
// not just the key - with "x509: unsupported elliptic curve" for anything else. Upstream DSS
// reaches java.security.cert.CertificateFactory with BouncyCastle registered, which knows the
// RFC 5639 Brainpool curves, so every Brainpool-signed certificate that upstream reads
// perfectly well was, before this package, entirely unparseable anywhere in this port - not
// XAdES-specific, but every path that loads a certificate from bytes. Such certificates are
// routine in German eIDAS / gematik telematics deployments; the dss-xades cross-validation
// corpus contains one (testdata/upstream/xades-ecc-brainpool.xml), where upstream identifies the
// signing certificate and reports the signature intact and this port reported neither.
//
// SECURITY NOTE. The arithmetic below is the textbook Jacobian-coordinate short-Weierstrass
// implementation over math/big, and it is NOT constant-time. That is acceptable here, and only
// here, because these curves are used exclusively to VERIFY signatures with PUBLIC keys read out
// of certificates: no secret scalar is ever multiplied on them. Nothing in this package may be
// used for key generation, signing, or ECDH; crypto/ecdsa's own Sign refuses custom curves under
// FIPS-140 mode for the same reason.
package eccurve

import (
	"crypto/elliptic"
	"math/big"
	"sync"
)

// weierstrassCurve is a short-Weierstrass curve y^2 = x^3 + a*x + b over GF(p) with an
// ARBITRARY a. It exists because crypto/elliptic.CurveParams implements only the a = -3 special
// case its four NIST curves share, which no Brainpool "r1" curve satisfies.
type weierstrassCurve struct {
	params elliptic.CurveParams
	a      *big.Int
	// oid is the dotted named-curve OID this curve is identified by inside an
	// AlgorithmIdentifier's ECParameters.
	oid string
}

// Params implements elliptic.Curve.
func (c *weierstrassCurve) Params() *elliptic.CurveParams { return &c.params }

// IsOnCurve implements elliptic.Curve: y^2 == x^3 + a*x + b (mod p).
func (c *weierstrassCurve) IsOnCurve(x, y *big.Int) bool {
	p := c.params.P
	if x.Sign() < 0 || y.Sign() < 0 || x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		return false
	}
	left := new(big.Int).Mul(y, y)
	left.Mod(left, p)

	right := new(big.Int).Mul(x, x)
	right.Mod(right, p)
	right.Mul(right, x)
	right.Mod(right, p)

	ax := new(big.Int).Mul(c.a, x)
	right.Add(right, ax)
	right.Add(right, c.params.B)
	right.Mod(right, p)

	return left.Cmp(right) == 0
}

// affineFromJacobian converts (x, y, z) back to affine coordinates, answering (0, 0) for the
// point at infinity - the same representation crypto/elliptic uses.
func (c *weierstrassCurve) affineFromJacobian(x, y, z *big.Int) (*big.Int, *big.Int) {
	if z.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}
	p := c.params.P
	zInv := new(big.Int).ModInverse(z, p)
	if zInv == nil {
		return new(big.Int), new(big.Int)
	}
	zInv2 := new(big.Int).Mul(zInv, zInv)
	zInv2.Mod(zInv2, p)

	outX := new(big.Int).Mul(x, zInv2)
	outX.Mod(outX, p)

	zInv3 := new(big.Int).Mul(zInv2, zInv)
	zInv3.Mod(zInv3, p)
	outY := new(big.Int).Mul(y, zInv3)
	outY.Mod(outY, p)
	return outX, outY
}

// addJacobian is the "add-2007-bl" style Jacobian point addition, independent of a.
func (c *weierstrassCurve) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	p := c.params.P
	if z1.Sign() == 0 {
		return new(big.Int).Set(x2), new(big.Int).Set(y2), new(big.Int).Set(z2)
	}
	if z2.Sign() == 0 {
		return new(big.Int).Set(x1), new(big.Int).Set(y1), new(big.Int).Set(z1)
	}

	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, p)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, p)

	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, p)
	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, p)

	s1 := new(big.Int).Mul(y1, z2z2)
	s1.Mod(s1, p)
	s1.Mul(s1, z2)
	s1.Mod(s1, p)
	s2 := new(big.Int).Mul(y2, z1z1)
	s2.Mod(s2, p)
	s2.Mul(s2, z1)
	s2.Mod(s2, p)

	if u1.Cmp(u2) == 0 {
		if s1.Cmp(s2) != 0 {
			// P + (-P): the point at infinity.
			return new(big.Int), new(big.Int), new(big.Int)
		}
		return c.doubleJacobian(x1, y1, z1)
	}

	h := new(big.Int).Sub(u2, u1)
	h.Mod(h, p)
	r := new(big.Int).Sub(s2, s1)
	r.Mod(r, p)

	h2 := new(big.Int).Mul(h, h)
	h2.Mod(h2, p)
	h3 := new(big.Int).Mul(h2, h)
	h3.Mod(h3, p)
	u1h2 := new(big.Int).Mul(u1, h2)
	u1h2.Mod(u1h2, p)

	x3 := new(big.Int).Mul(r, r)
	x3.Sub(x3, h3)
	x3.Sub(x3, u1h2)
	x3.Sub(x3, u1h2)
	x3.Mod(x3, p)

	y3 := new(big.Int).Sub(u1h2, x3)
	y3.Mul(y3, r)
	s1h3 := new(big.Int).Mul(s1, h3)
	y3.Sub(y3, s1h3)
	y3.Mod(y3, p)

	z3 := new(big.Int).Mul(z1, z2)
	z3.Mod(z3, p)
	z3.Mul(z3, h)
	z3.Mod(z3, p)

	return x3, y3, z3
}

// doubleJacobian is the generic-a Jacobian doubling: M = 3*X^2 + a*Z^4, where crypto/elliptic's
// own version folds a = -3 into M = 3*(X - Z^2)*(X + Z^2).
func (c *weierstrassCurve) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	p := c.params.P
	if z.Sign() == 0 || y.Sign() == 0 {
		return new(big.Int), new(big.Int), new(big.Int)
	}

	yy := new(big.Int).Mul(y, y)
	yy.Mod(yy, p)

	s := new(big.Int).Mul(x, yy)
	s.Lsh(s, 2)
	s.Mod(s, p)

	zz := new(big.Int).Mul(z, z)
	zz.Mod(zz, p)
	z4 := new(big.Int).Mul(zz, zz)
	z4.Mod(z4, p)

	m := new(big.Int).Mul(x, x)
	m.Mod(m, p)
	m.Mul(m, big.NewInt(3))
	az4 := new(big.Int).Mul(c.a, z4)
	m.Add(m, az4)
	m.Mod(m, p)

	x3 := new(big.Int).Mul(m, m)
	x3.Sub(x3, s)
	x3.Sub(x3, s)
	x3.Mod(x3, p)

	yyyy := new(big.Int).Mul(yy, yy)
	yyyy.Mod(yyyy, p)
	yyyy.Lsh(yyyy, 3)
	yyyy.Mod(yyyy, p)

	y3 := new(big.Int).Sub(s, x3)
	y3.Mul(y3, m)
	y3.Sub(y3, yyyy)
	y3.Mod(y3, p)

	z3 := new(big.Int).Mul(y, z)
	z3.Lsh(z3, 1)
	z3.Mod(z3, p)

	return x3, y3, z3
}

// Add implements elliptic.Curve.
func (c *weierstrassCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	z1 := jacobianZFor(x1, y1)
	z2 := jacobianZFor(x2, y2)
	return c.affineFromJacobian(c.addJacobian(x1, y1, z1, x2, y2, z2))
}

// Double implements elliptic.Curve.
func (c *weierstrassCurve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	z1 := jacobianZFor(x1, y1)
	return c.affineFromJacobian(c.doubleJacobian(x1, y1, z1))
}

// ScalarMult implements elliptic.Curve by plain left-to-right double-and-add. See the package
// note: NOT constant-time, and only ever driven by public scalars.
func (c *weierstrassCurve) ScalarMult(bx, by *big.Int, k []byte) (*big.Int, *big.Int) {
	x, y, z := new(big.Int), new(big.Int), new(big.Int)
	bz := jacobianZFor(bx, by)
	for _, b := range k {
		for bit := 0; bit < 8; bit++ {
			x, y, z = c.doubleJacobian(x, y, z)
			if b&0x80 == 0x80 {
				x, y, z = c.addJacobian(bx, by, bz, x, y, z)
			}
			b <<= 1
		}
	}
	return c.affineFromJacobian(x, y, z)
}

// ScalarBaseMult implements elliptic.Curve.
func (c *weierstrassCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return c.ScalarMult(c.params.Gx, c.params.Gy, k)
}

// jacobianZFor answers the Jacobian Z for an affine point: 0 for the point at infinity, which
// crypto/elliptic spells (0, 0), and 1 otherwise.
func jacobianZFor(x, y *big.Int) *big.Int {
	if x.Sign() == 0 && y.Sign() == 0 {
		return new(big.Int)
	}
	return big.NewInt(1)
}

// registry maps a dotted named-curve OID to the curve it names.
var (
	registryOnce sync.Once
	registry     = map[string]*weierstrassCurve{}
)

func registerWeierstrassCurve(c *weierstrassCurve) { registry[c.oid] = c }

func hexInt(s string) *big.Int {
	value, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("eccurve: malformed curve constant " + s)
	}
	return value
}

// CurveForOID answers the curve a dotted named-curve OID identifies, or nil when this package
// does not know it. The returned curve is shared and must not be mutated.
func CurveForOID(oid string) elliptic.Curve {
	registryOnce.Do(initRegistry)
	if c, ok := registry[oid]; ok {
		return c
	}
	return nil
}

// KnownCurveOIDs lists every OID CurveForOID answers, in no particular order. Test-facing.
func KnownCurveOIDs() []string {
	registryOnce.Do(initRegistry)
	oids := make([]string, 0, len(registry))
	for oid := range registry {
		oids = append(oids, oid)
	}
	return oids
}
