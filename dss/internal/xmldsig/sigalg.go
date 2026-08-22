// Ported from org.apache.xml.security.algorithms.SignatureAlgorithm and its implementations
// SignatureBaseRSA, SignatureBaseRSAPSS, SignatureECDSA, SignatureDSA and SignatureEDDSA
// (Apache Santuario xmlsec 3.0.6), plus the JCEMapper entries DSS installs.
package xmldsig

import (
	"crypto"
	"crypto/dsa" //nolint:staticcheck // XMLDSIG defines a DSAwithSHA1 method; verifying it is not deprecated.
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
)

// ErrUnsupportedSignatureAlgorithm reports a ds:SignatureMethod this build cannot verify.
// It is the counterpart of Santuario's XMLSignatureException("algorithms.NoSuchAlgorithm"),
// and it is deliberately distinct from "the signature does not verify".
var ErrUnsupportedSignatureAlgorithm = errors.New("xmldsig: unsupported signature algorithm")

// ErrKeyMismatch reports a public key whose type does not match the signature method - an
// RSA key under an ECDSA method, say. Santuario surfaces it as an InvalidKeyException wrapped
// in XMLSignatureException, i.e. as a failure, never as "invalid signature".
var ErrKeyMismatch = errors.New("xmldsig: the public key does not match the signature algorithm")

// ErrMalformedSignatureValue reports a ds:SignatureValue that cannot be a signature under
// the key at all - the wrong length for the RSA modulus. The JCE throws SignatureException
// for it rather than answering false, and so must this: "this is not a signature" and "this
// signature is forged" are different findings, and a validator that reports the second for
// the first has silently decided the document was signed and merely tampered with.
var ErrMalformedSignatureValue = errors.New("xmldsig: malformed ds:SignatureValue")

// cryptoHashFor resolves the crypto.Hash IDENTITY of a DigestAlgorithm, which is what RSA
// PKCS#1 v1.5 and PSS verification need (the OID goes into the padding); the octets themselves
// are hashed through spi.DSSUtilsDigest, so the two never disagree about what SHA-256 is.
func cryptoHashFor(alg enumerations.DigestAlgorithm) (crypto.Hash, error) {
	var h crypto.Hash
	switch alg {
	case enumerations.DigestAlgorithmMD5:
		h = crypto.MD5
	case enumerations.DigestAlgorithmSHA1:
		h = crypto.SHA1
	case enumerations.DigestAlgorithmSHA224:
		h = crypto.SHA224
	case enumerations.DigestAlgorithmSHA256:
		h = crypto.SHA256
	case enumerations.DigestAlgorithmSHA384:
		h = crypto.SHA384
	case enumerations.DigestAlgorithmSHA512:
		h = crypto.SHA512
	case enumerations.DigestAlgorithmSHA3224:
		h = crypto.SHA3_224
	case enumerations.DigestAlgorithmSHA3256:
		h = crypto.SHA3_256
	case enumerations.DigestAlgorithmSHA3384:
		h = crypto.SHA3_384
	case enumerations.DigestAlgorithmSHA3512:
		h = crypto.SHA3_512
	case enumerations.DigestAlgorithmRIPEMD160:
		h = crypto.RIPEMD160
	default:
		return 0, fmt.Errorf("%w: digest %s", ErrUnsupportedSignatureAlgorithm, alg)
	}
	if !h.Available() {
		return 0, fmt.Errorf("%w: digest %s is not linked into this build", ErrUnsupportedSignatureAlgorithm, alg)
	}
	return h, nil
}

// verifySignature checks signatureValue over signedContent under the ds:SignatureMethod URI.
// Port of SignatureAlgorithm#verify as XMLSignature#checkSignatureValue drives it.
//
// The (bool, error) split is Santuario's boolean/exception split and callers depend on it:
// false means the cryptography said no, an error means the algorithm, the key or the encoding
// made the question unanswerable. Collapsing the two would turn "this signature is forged"
// into the same report as "this build cannot check brainpool curves".
func verifySignature(uri string, pub crypto.PublicKey, signedContent, signatureValue []byte) (bool, error) {
	alg, err := enumerations.SignatureAlgorithmForXML(uri)
	if err != nil {
		return false, fmt.Errorf("%w: %q", ErrUnsupportedSignatureAlgorithm, uri)
	}
	digestAlg := alg.DigestAlgorithm()
	encryption := alg.EncryptionAlgorithm()

	switch encryption {
	case enumerations.EncryptionAlgorithmRSA, enumerations.EncryptionAlgorithmRSASSAPSS:
		key, ok := pub.(*rsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("%w: RSA method with a %T key", ErrKeyMismatch, pub)
		}
		// A ds:SignatureValue whose length is not the modulus length is a MALFORMED
		// signature, not a wrong one, and the two are reported differently. The JCE's
		// RSA verifier throws SignatureException("Signature length not correct: got N
		// but was expecting M") before any arithmetic, which Santuario propagates;
		// crypto/rsa folds the same condition into ErrVerification, which cryptoVerdict
		// would then turn into a plain "false".
		//
		// Two upstream fixtures land here and both must fail rather than answer:
		// validation/xades-not-base64-sigValue.xml, whose ds:SignatureValue is the text
		// "Hello World!", and validation/dss1334/simple-test.signed-only-detached-
		// LuxTrustCA3.xml, whose 256-byte value does not match its key.
		if len(signatureValue) != key.Size() {
			return false, fmt.Errorf("%w: ds:SignatureValue is %d bytes, the RSA modulus is %d",
				ErrMalformedSignatureValue, len(signatureValue), key.Size())
		}
		h, err := cryptoHashFor(digestAlg)
		if err != nil {
			return false, err
		}
		digest, err := spi.DSSUtilsDigest(digestAlg, signedContent)
		if err != nil {
			return false, err
		}
		if encryption == enumerations.EncryptionAlgorithmRSASSAPSS {
			// SignatureBaseRSAPSS builds a PSSParameterSpec whose salt length is the digest
			// length and whose trailer field is 1, which is exactly PSSSaltLengthEqualsHash
			// with MGF1 over the same digest - the RSASSA-PSS profile XMLDSIG 1.1 defines.
			err = rsa.VerifyPSS(key, h, digest, signatureValue, &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthEqualsHash,
				Hash:       h,
			})
		} else {
			err = rsa.VerifyPKCS1v15(key, h, digest, signatureValue)
		}
		return cryptoVerdict(err)

	case enumerations.EncryptionAlgorithmECDSA, enumerations.EncryptionAlgorithmPlainECDSA:
		key, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("%w: ECDSA method with a %T key", ErrKeyMismatch, pub)
		}
		r, s, err := decodeXMLDSigECDSA(signatureValue, key)
		if err != nil {
			return false, err
		}
		digest, err := spi.DSSUtilsDigest(digestAlg, signedContent)
		if err != nil {
			return false, err
		}
		return ecdsa.Verify(key, digest, r, s), nil

	case enumerations.EncryptionAlgorithmDSA:
		key, ok := pub.(*dsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("%w: DSA method with a %T key", ErrKeyMismatch, pub)
		}
		// SignatureDSA uses the same fixed-width r||s encoding as ECDSA, sized by the
		// subgroup order q rather than by the curve.
		size := (key.Q.BitLen() + 7) / 8
		if len(signatureValue) != 2*size {
			return false, fmt.Errorf("xmldsig: DSA signature length %d, expected %d",
				len(signatureValue), 2*size)
		}
		digest, err := spi.DSSUtilsDigest(digestAlg, signedContent)
		if err != nil {
			return false, err
		}
		r := new(big.Int).SetBytes(signatureValue[:size])
		s := new(big.Int).SetBytes(signatureValue[size:])
		return dsa.Verify(key, digest, r, s), nil

	case enumerations.EncryptionAlgorithmEDDSA:
		key, ok := pub.(ed25519.PublicKey)
		if !ok {
			return false, fmt.Errorf("%w: EdDSA method with a %T key", ErrKeyMismatch, pub)
		}
		if alg != enumerations.SignatureAlgorithmED25519 {
			return false, fmt.Errorf("%w: %s", ErrUnsupportedSignatureAlgorithm, alg)
		}
		return ed25519.Verify(key, signedContent, signatureValue), nil
	}
	return false, fmt.Errorf("%w: %s", ErrUnsupportedSignatureAlgorithm, alg)
}

// cryptoVerdict maps a crypto/rsa verification outcome onto Santuario's boolean.
// rsa.ErrVerification is "the signature is wrong"; anything else is a structural failure.
func cryptoVerdict(err error) (bool, error) {
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, rsa.ErrVerification):
		return false, nil
	default:
		return false, err
	}
}

// decodeXMLDSigECDSA splits an XMLDSIG ECDSA signature into (r, s). Port of
// SignatureECDSA#convertXMLDSIGtoASN1, read backwards.
//
// XMLDSIG 6.4.3 encodes an ECDSA signature as the raw pair r||s, each padded to the byte
// length of the curve's order - the IEEE P1363 form - and NOT as the DER SEQUENCE that
// crypto/x509 and most other Go APIs expect. Getting this wrong does not produce an error, it
// produces a signature that never verifies, which is why it has its own function and its own
// length check.
func decodeXMLDSigECDSA(sig []byte, key *ecdsa.PublicKey) (r, s *big.Int, err error) {
	if key.Curve == nil {
		return nil, nil, fmt.Errorf("%w: the EC key names no curve", ErrKeyMismatch)
	}
	size := (key.Curve.Params().N.BitLen() + 7) / 8
	if len(sig) != 2*size {
		return nil, nil, fmt.Errorf("xmldsig: ECDSA signature length %d, expected %d for %s",
			len(sig), 2*size, key.Curve.Params().Name)
	}
	return new(big.Int).SetBytes(sig[:size]), new(big.Int).SetBytes(sig[size:]), nil
}
