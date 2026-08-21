// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSPKUtils.java (DSS 6.5.RC1).
package spi

import (
	"crypto/dsa" //nolint:staticcheck // DSA keys still occur in legacy signatures being validated.
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/asn1"
	"strconv"

	"github.com/ryftcore/dss-go/dss/model"
	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// dssPKUtilsUnknownKeySize defines the string returned in case of unknown key size.
const dssPKUtilsUnknownKeySize = "?"

// Algorithm OIDs of the EdDSA/XDH keys Go's crypto/x509 cannot parse into a typed key.
var (
	dssPKUtilsOIDEd25519 = asn1.ObjectIdentifier{1, 3, 101, 112}
	dssPKUtilsOIDEd448   = asn1.ObjectIdentifier{1, 3, 101, 113}
	dssPKUtilsOIDX25519  = asn1.ObjectIdentifier{1, 3, 101, 110}
	dssPKUtilsOIDX448    = asn1.ObjectIdentifier{1, 3, 101, 111}
)

// DSSPKUtilsStringPublicKeySizeOfToken returns the key length used to sign the given token.
// Port of getStringPublicKeySize(Token).
func DSSPKUtilsStringPublicKeySizeOfToken(token model.Token) string {
	keyLength := dssPKUtilsUnknownKeySize
	var issuerPublicKey *model.PublicKey
	if token.PublicKeyOfTheSigner() != nil {
		issuerPublicKey = token.PublicKeyOfTheSigner()
	} else if token.IsSelfSigned() {
		// Only a CertificateToken can be self-signed; the Java cast would throw otherwise.
		issuerPublicKey = token.(*model.CertificateToken).PublicKey()
	}
	if issuerPublicKey != nil {
		keyLength = DSSPKUtilsStringPublicKeySize(issuerPublicKey)
	}
	return keyLength
}

// DSSPKUtilsStringPublicKeySize returns the key length extracted from the public key, or
// "?" when it cannot be determined. Port of getStringPublicKeySize(PublicKey).
func DSSPKUtilsStringPublicKeySize(publicKey *model.PublicKey) string {
	keyLength := dssPKUtilsUnknownKeySize
	if publicKey != nil {
		publicKeySize := DSSPKUtilsPublicKeySize(publicKey)
		if publicKeySize > 0 {
			keyLength = strconv.Itoa(publicKeySize)
		}
	}
	return keyLength
}

// DSSPKUtilsPublicKeySize returns the public key size extracted from the public key
// infrastructure: -1 for an unrecognised key type and 0 for a supported key whose length
// cannot be established. Port of getPublicKeySize(PublicKey).
//
// UPSTREAM QUIRK, reproduced verbatim: the EdDSA and XDH branches return a BYTE count
// (the SubjectPublicKeyInfo length minus a fixed 12-byte prefix), not a bit length, unlike
// every other branch. An Ed25519 key therefore reports 32 and an Ed448 key 57.
//
// DEVIATIONS: BouncyCastle's JCEECPublicKey branch - an EC key carrying explicit domain
// parameters, sized by the subgroup order n - has no Go counterpart, because crypto/x509
// only decodes named curves; such a key is reported as unrecognised (-1) rather than 0.
// Ed448/X448 keys, which crypto/x509 cannot parse at all, are recognised from their
// SubjectPublicKeyInfo algorithm OID so that they keep reporting a size.
func DSSPKUtilsPublicKeySize(publicKey *model.PublicKey) int {
	publicKeySize := -1
	switch key := publicKey.Key().(type) {
	case *rsa.PublicKey:
		publicKeySize = key.N.BitLen()
	case *ecdsa.PublicKey:
		if key.Curve != nil {
			publicKeySize = key.Curve.Params().BitSize
		} else {
			publicKeySize = 0
		}
	case *dsa.PublicKey:
		publicKeySize = key.Parameters.P.BitLen()
	case ed25519.PublicKey:
		return dssPKUtilsEncodedKeySize(publicKey)
	case *ecdh.PublicKey:
		// crypto/x509 decodes an X25519 SubjectPublicKeyInfo into an ecdh key; the EC
		// curves are returned as *ecdsa.PublicKey and never reach this branch.
		return dssPKUtilsEncodedKeySize(publicKey)
	default:
		switch {
		case dssPKUtilsHasAlgorithmOID(publicKey, dssPKUtilsOIDEd25519),
			dssPKUtilsHasAlgorithmOID(publicKey, dssPKUtilsOIDEd448),
			dssPKUtilsHasAlgorithmOID(publicKey, dssPKUtilsOIDX25519),
			dssPKUtilsHasAlgorithmOID(publicKey, dssPKUtilsOIDX448):
			return dssPKUtilsEncodedKeySize(publicKey)
		}
		// Upstream logs "Unknown public key infrastructure: {}" here.
	}
	return publicKeySize
}

// dssPKUtilsEncodedKeySize reproduces the EdDSA/XDH sizing of getPublicKeySize: the length
// of the encoded SubjectPublicKeyInfo minus the 12-byte prefix Java measures as
// Utils.fromHex("3043300506032b6571033a00").length (Ed448) resp.
// Utils.fromHex("3042300506032b656f033900").length (X448) - both 12 bytes.
func dssPKUtilsEncodedKeySize(publicKey *model.PublicKey) int {
	const prefixSize = 12
	return len(publicKey.Encoded()) - prefixSize
}

// dssPKUtilsHasAlgorithmOID reports whether the key's SubjectPublicKeyInfo names the given
// algorithm OID. It lets the port recognise the key types crypto/x509 refuses to parse.
func dssPKUtilsHasAlgorithmOID(publicKey *model.PublicKey, oid asn1.ObjectIdentifier) bool {
	var spki, algorithm cryptobyte.String
	var keyOID asn1.ObjectIdentifier
	input := cryptobyte.String(publicKey.Encoded())
	if !input.ReadASN1(&spki, cbasn1.SEQUENCE) {
		return false
	}
	if !spki.ReadASN1(&algorithm, cbasn1.SEQUENCE) {
		return false
	}
	if !algorithm.ReadASN1ObjectIdentifier(&keyOID) {
		return false
	}
	return oid.Equal(keyOID)
}
