package pfx

import (
	"crypto"
	"crypto/dsa" //nolint:staticcheck // DSA keys still occur in legacy PKCS#12 fixtures being loaded.
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// parsePrivateKeyInfo parses a cleartext PKCS#8 PrivateKeyInfo (RFC 5958), the shape both a
// keyBag and a decrypted pkcs8ShroudedKeyBag hold. crypto/x509 parses every algorithm it knows
// (RSA, EC, Ed25519); DSA - which neither x509.ParsePKCS8PrivateKey nor crypto/dsa itself parses
// from PKCS#8 - is parsed by hand in parseDSAPrivateKeyInfo.
func parsePrivateKeyInfo(der []byte) (crypto.PrivateKey, error) {
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err == nil {
		return key, nil
	}
	if dsaKey, dsaErr := parseDSAPrivateKeyInfo(der); dsaErr == nil {
		return dsaKey, nil
	}
	return nil, err
}

// dsaPrivateKeyInfo is the PKCS#8 PrivateKeyInfo shape a DSA key uses (RFC 3279 section 2.3.2):
// the algorithm's parameters carry the domain parameters (P, Q, G), and the OCTET STRING
// privateKey field wraps a bare INTEGER holding X. Unlike RSA/EC/Ed25519 there is no public key
// component to extract, so Y is recomputed as G^X mod P.
func parseDSAPrivateKeyInfo(der []byte) (*dsa.PrivateKey, error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, errors.New("pfx: not a PrivateKeyInfo SEQUENCE")
	}
	children := element.Children()
	if len(children) < 3 {
		return nil, fmt.Errorf("pfx: PrivateKeyInfo holds %d components, at least 3 expected", len(children))
	}

	algorithm, err := asn1ber.AlgorithmIdentifierFromElement(children[1])
	if err != nil {
		return nil, err
	}
	if !algorithm.Algorithm.Equal(oidDSA) {
		return nil, fmt.Errorf("pfx: unsupported PKCS#8 private key algorithm %s", algorithm.Algorithm)
	}
	if algorithm.Parameters == nil {
		return nil, errors.New("pfx: DSA PrivateKeyInfo is missing its domain parameters")
	}
	p, q, g, err := parseDSADomainParameters(algorithm.Parameters)
	if err != nil {
		return nil, err
	}

	privateKeyField := children[2]
	if !privateKeyField.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: PrivateKeyInfo.privateKey is not an OCTET STRING")
	}
	xElement, xRest, err := asn1ber.Parse(privateKeyField.Octets())
	if err != nil || len(xRest) != 0 || !xElement.IsUniversal(asn1ber.TagInteger) {
		return nil, errors.New("pfx: DSA PrivateKeyInfo.privateKey does not wrap an INTEGER")
	}
	x := xElement.Integer()
	if x.Sign() <= 0 || x.Cmp(q) >= 0 {
		return nil, errors.New("pfx: DSA private key X is out of range")
	}

	y := new(big.Int).Exp(g, x, p)
	return &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{P: p, Q: q, G: g},
			Y:          y,
		},
		X: x,
	}, nil
}

// parseDSADomainParameters decodes the Dss-Parms SEQUENCE { p, q, g } RFC 3279 section 2.3.2
// defines as id-dsa's AlgorithmIdentifier.parameters.
func parseDSADomainParameters(der []byte) (p, q, g *big.Int, err error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(rest) != 0 || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, nil, nil, errors.New("pfx: DSA algorithm parameters are not a SEQUENCE")
	}
	children := element.Children()
	if len(children) != 3 {
		return nil, nil, nil, fmt.Errorf("pfx: Dss-Parms holds %d components, 3 expected", len(children))
	}
	for i, child := range children {
		if !child.IsUniversal(asn1ber.TagInteger) {
			return nil, nil, nil, fmt.Errorf("pfx: Dss-Parms component %d is not an INTEGER", i)
		}
		// A zero modulus would turn the Y = G^X mod P computation into an unbounded
		// exponentiation (big.Int#Exp with m == 0 computes G^X in full).
		if child.Integer().Sign() <= 0 {
			return nil, nil, nil, fmt.Errorf("pfx: Dss-Parms component %d is not positive", i)
		}
	}
	return children[0].Integer(), children[1].Integer(), children[2].Integer(), nil
}
