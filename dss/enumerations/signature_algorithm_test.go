package enumerations

import (
	"encoding/hex"
	"strconv"
	"testing"
)

// TestSignatureAlgorithmFields exhaustively checks the (encryptionAlgorithm,
// digestAlgorithm) pair for every SignatureAlgorithm against the Java enum
// constructor arguments.
func TestSignatureAlgorithmFields(t *testing.T) {
	cases := []struct {
		v      SignatureAlgorithm
		enc    EncryptionAlgorithm
		digest DigestAlgorithm
	}{
		{SignatureAlgorithm_RSA_RAW, EncryptionAlgorithm_RSA, ""},
		{SignatureAlgorithm_RSA_SHA1, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_RSA_SHA224, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_RSA_SHA256, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_RSA_SHA384, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_RSA_SHA512, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_RSA_SHA3_224, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_RSA_SHA3_256, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_RSA_SHA3_384, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_RSA_SHA3_512, EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1, EncryptionAlgorithm_RSASSA_PSS, ""},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1, EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_RSA_RIPEMD160, EncryptionAlgorithm_RSA, DigestAlgorithm_RIPEMD160},
		{SignatureAlgorithm_RSA_MD2, EncryptionAlgorithm_RSA, DigestAlgorithm_MD2},
		{SignatureAlgorithm_RSA_MD5, EncryptionAlgorithm_RSA, DigestAlgorithm_MD5},
		{SignatureAlgorithm_ECDSA_RAW, EncryptionAlgorithm_ECDSA, ""},
		{SignatureAlgorithm_ECDSA_SHA1, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_ECDSA_SHA224, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_ECDSA_SHA256, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_ECDSA_SHA384, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_ECDSA_SHA512, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_ECDSA_SHA3_224, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_ECDSA_SHA3_256, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_ECDSA_SHA3_384, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_ECDSA_SHA3_512, EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_ECDSA_RIPEMD160, EncryptionAlgorithm_ECDSA, DigestAlgorithm_RIPEMD160},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA1, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA224, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA256, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA384, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA512, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA3_224, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA3_256, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA3_384, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA3_512, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160, EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_RIPEMD160},
		{SignatureAlgorithm_DSA_RAW, EncryptionAlgorithm_DSA, ""},
		{SignatureAlgorithm_DSA_SHA1, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_DSA_SHA224, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_DSA_SHA256, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_DSA_SHA384, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_DSA_SHA512, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_DSA_SHA3_224, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_DSA_SHA3_256, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_DSA_SHA3_384, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_DSA_SHA3_512, EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_HMAC_SHA1, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA1},
		{SignatureAlgorithm_HMAC_SHA224, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA224},
		{SignatureAlgorithm_HMAC_SHA256, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA256},
		{SignatureAlgorithm_HMAC_SHA384, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA384},
		{SignatureAlgorithm_HMAC_SHA512, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_HMAC_SHA3_224, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_224},
		{SignatureAlgorithm_HMAC_SHA3_256, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_256},
		{SignatureAlgorithm_HMAC_SHA3_384, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_384},
		{SignatureAlgorithm_HMAC_SHA3_512, EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_512},
		{SignatureAlgorithm_HMAC_RIPEMD160, EncryptionAlgorithm_HMAC, DigestAlgorithm_RIPEMD160},
		{SignatureAlgorithm_ED25519, EncryptionAlgorithm_EDDSA, DigestAlgorithm_SHA512},
		{SignatureAlgorithm_ED448, EncryptionAlgorithm_EDDSA, DigestAlgorithm_SHAKE256_512},
	}
	if len(cases) != len(SignatureAlgorithmValues()) {
		t.Fatalf("test case count = %d, want %d (SignatureAlgorithmValues length)", len(cases), len(SignatureAlgorithmValues()))
	}
	for _, c := range cases {
		if got := c.v.EncryptionAlgorithm(); got != c.enc {
			t.Errorf("%v.EncryptionAlgorithm() = %q, want %q", c.v, got, c.enc)
		}
		if got := c.v.DigestAlgorithm(); got != c.digest {
			t.Errorf("%v.DigestAlgorithm() = %q, want %q", c.v, got, c.digest)
		}
	}
}

func TestSignatureAlgorithmValueOf(t *testing.T) {
	for _, v := range SignatureAlgorithmValues() {
		got, err := SignatureAlgorithmValueOf(string(v))
		if err != nil {
			t.Fatalf("SignatureAlgorithmValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignatureAlgorithmValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := SignatureAlgorithmValueOf("bogus"); err == nil {
		t.Error("SignatureAlgorithmValueOf(\"bogus\") expected error, got nil")
	}
}

// TestSignatureAlgorithmForXML exhaustively round-trips every XML URI alias declared in
// the Java source (including the duplicate RSA_RIPEMD160 alias) through
// SignatureAlgorithmForXML.
func TestSignatureAlgorithmForXML(t *testing.T) {
	cases := []struct {
		uri  string
		want SignatureAlgorithm
	}{
		{"http://www.w3.org/2000/09/xmldsig#rsa-sha1", SignatureAlgorithm_RSA_SHA1},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-sha224", SignatureAlgorithm_RSA_SHA224},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", SignatureAlgorithm_RSA_SHA256},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-sha384", SignatureAlgorithm_RSA_SHA384},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-sha512", SignatureAlgorithm_RSA_SHA512},
		{"http://www.w3.org/2007/05/xmldsig-more#sha1-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha224-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha256-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha384-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha512-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha3-224-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha3-256-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha3-384-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1},
		{"http://www.w3.org/2007/05/xmldsig-more#sha3-512-rsa-MGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-ripemd160", SignatureAlgorithm_RSA_RIPEMD160},
		{"http://www.w3.org/2001/04/xmldsig-more/rsa-ripemd160", SignatureAlgorithm_RSA_RIPEMD160},
		{"http://www.w3.org/2001/04/xmldsig-more#rsa-md5", SignatureAlgorithm_RSA_MD5},
		{"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha1", SignatureAlgorithm_ECDSA_SHA1},
		{"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha224", SignatureAlgorithm_ECDSA_SHA224},
		{"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256", SignatureAlgorithm_ECDSA_SHA256},
		{"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha384", SignatureAlgorithm_ECDSA_SHA384},
		{"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha512", SignatureAlgorithm_ECDSA_SHA512},
		{"http://www.w3.org/2021/04/xmldsig-more#ecdsa-sha3-224", SignatureAlgorithm_ECDSA_SHA3_224},
		{"http://www.w3.org/2021/04/xmldsig-more#ecdsa-sha3-256", SignatureAlgorithm_ECDSA_SHA3_256},
		{"http://www.w3.org/2021/04/xmldsig-more#ecdsa-sha3-384", SignatureAlgorithm_ECDSA_SHA3_384},
		{"http://www.w3.org/2021/04/xmldsig-more#ecdsa-sha3-512", SignatureAlgorithm_ECDSA_SHA3_512},
		{"http://www.w3.org/2007/05/xmldsig-more#ecdsa-ripemd160", SignatureAlgorithm_ECDSA_RIPEMD160},
		{"http://www.w3.org/2021/04/xmldsig-more#eddsa-ed25519", SignatureAlgorithm_ED25519},
		{"http://www.w3.org/2021/04/xmldsig-more#eddsa-ed448", SignatureAlgorithm_ED448},
		{"http://www.w3.org/2000/09/xmldsig#dsa-sha1", SignatureAlgorithm_DSA_SHA1},
		{"http://www.w3.org/2009/xmldsig11#dsa-sha256", SignatureAlgorithm_DSA_SHA256},
		{"http://www.w3.org/2000/09/xmldsig#hmac-sha1", SignatureAlgorithm_HMAC_SHA1},
		{"http://www.w3.org/2001/04/xmldsig-more#hmac-sha224", SignatureAlgorithm_HMAC_SHA224},
		{"http://www.w3.org/2001/04/xmldsig-more#hmac-sha256", SignatureAlgorithm_HMAC_SHA256},
		{"http://www.w3.org/2001/04/xmldsig-more#hmac-sha384", SignatureAlgorithm_HMAC_SHA384},
		{"http://www.w3.org/2001/04/xmldsig-more#hmac-sha512", SignatureAlgorithm_HMAC_SHA512},
		{"http://www.w3.org/2001/04/xmldsig-more#hmac-ripemd160", SignatureAlgorithm_HMAC_RIPEMD160},
	}
	if len(cases) != len(signatureAlgorithmXMLPairs) {
		t.Fatalf("test case count = %d, want %d (signatureAlgorithmXMLPairs length)", len(cases), len(signatureAlgorithmXMLPairs))
	}
	for _, c := range cases {
		got, err := SignatureAlgorithmForXML(c.uri)
		if err != nil {
			t.Fatalf("SignatureAlgorithmForXML(%q) returned error: %v", c.uri, err)
		}
		if got != c.want {
			t.Errorf("SignatureAlgorithmForXML(%q) = %q, want %q", c.uri, got, c.want)
		}
	}

	if _, err := SignatureAlgorithmForXML("bogus"); err == nil {
		t.Error("SignatureAlgorithmForXML(\"bogus\") expected error, got nil")
	}
	if got := SignatureAlgorithmForXMLDefault("bogus", SignatureAlgorithm_RSA_SHA256); got != SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("SignatureAlgorithmForXMLDefault(\"bogus\", RSA_SHA256) = %q, want RSA_SHA256", got)
	}
	if got := SignatureAlgorithmForXMLDefault("http://www.w3.org/2000/09/xmldsig#rsa-sha1", SignatureAlgorithm_RSA_SHA256); got != SignatureAlgorithm_RSA_SHA1 {
		t.Errorf("SignatureAlgorithmForXMLDefault(known) = %q, want RSA_SHA1", got)
	}
}

// TestSignatureAlgorithmURI checks the canonical (reverse) XML URI for every algorithm
// that has one, including the PLAIN_ECDSA algorithms which inherit their ECDSA
// counterpart's URI via ensurePlainECDSA.
func TestSignatureAlgorithmURI(t *testing.T) {
	cases := []struct {
		v    SignatureAlgorithm
		want string
	}{
		{SignatureAlgorithm_RSA_SHA1, "http://www.w3.org/2000/09/xmldsig#rsa-sha1"},
		// RSA_RIPEMD160 has two XML URI aliases; the "/rsa-ripemd160" spelling is the one
		// that wins upstream's HashMap-order reverse mapping, despite being declared second.
		{SignatureAlgorithm_RSA_RIPEMD160, "http://www.w3.org/2001/04/xmldsig-more/rsa-ripemd160"},
		{SignatureAlgorithm_ECDSA_SHA256, "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"},
		{SignatureAlgorithm_PLAIN_ECDSA_SHA256, "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"},
		{SignatureAlgorithm_ECDSA_RIPEMD160, "http://www.w3.org/2007/05/xmldsig-more#ecdsa-ripemd160"},
		{SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160, "http://www.w3.org/2007/05/xmldsig-more#ecdsa-ripemd160"},
		{SignatureAlgorithm_ED25519, "http://www.w3.org/2021/04/xmldsig-more#eddsa-ed25519"},
		{SignatureAlgorithm_ED448, "http://www.w3.org/2021/04/xmldsig-more#eddsa-ed448"},
		{SignatureAlgorithm_RSA_RAW, ""},
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.want {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.want)
		}
	}
}

// TestSignatureAlgorithmForOID exhaustively round-trips every OID alias declared in the
// Java source, including the duplicate RSA_SHA1/DSA_SHA1 aliases and the
// ECDSA_RIPEMD160/PLAIN_ECDSA_RIPEMD160 OID collision.
func TestSignatureAlgorithmForOID(t *testing.T) {
	cases := []struct {
		oid  string
		want SignatureAlgorithm
	}{
		{"1.2.840.113549.1.1.5", SignatureAlgorithm_RSA_SHA1},
		{"1.3.14.3.2.29", SignatureAlgorithm_RSA_SHA1},
		{"1.2.840.113549.1.1.14", SignatureAlgorithm_RSA_SHA224},
		{"1.2.840.113549.1.1.11", SignatureAlgorithm_RSA_SHA256},
		{"1.2.840.113549.1.1.12", SignatureAlgorithm_RSA_SHA384},
		{"1.2.840.113549.1.1.13", SignatureAlgorithm_RSA_SHA512},
		{"1.3.36.3.3.1.2", SignatureAlgorithm_RSA_RIPEMD160},
		{"2.16.840.1.101.3.4.3.13", SignatureAlgorithm_RSA_SHA3_224},
		{"2.16.840.1.101.3.4.3.14", SignatureAlgorithm_RSA_SHA3_256},
		{"2.16.840.1.101.3.4.3.15", SignatureAlgorithm_RSA_SHA3_384},
		{"2.16.840.1.101.3.4.3.16", SignatureAlgorithm_RSA_SHA3_512},
		{"1.2.840.113549.1.1.4", SignatureAlgorithm_RSA_MD5},
		{"1.2.840.113549.1.1.2", SignatureAlgorithm_RSA_MD2},
		{"1.2.840.10045.4.1", SignatureAlgorithm_ECDSA_SHA1},
		{"1.2.840.10045.4.3.1", SignatureAlgorithm_ECDSA_SHA224},
		{"1.2.840.10045.4.3.2", SignatureAlgorithm_ECDSA_SHA256},
		{"1.2.840.10045.4.3.3", SignatureAlgorithm_ECDSA_SHA384},
		{"1.2.840.10045.4.3.4", SignatureAlgorithm_ECDSA_SHA512},
		// Overwritten by PLAIN_ECDSA_RIPEMD160 registered later at the same OID.
		{"0.4.0.127.0.7.1.1.4.1.6", SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160},
		{"2.16.840.1.101.3.4.3.9", SignatureAlgorithm_ECDSA_SHA3_224},
		{"2.16.840.1.101.3.4.3.10", SignatureAlgorithm_ECDSA_SHA3_256},
		{"2.16.840.1.101.3.4.3.11", SignatureAlgorithm_ECDSA_SHA3_384},
		{"2.16.840.1.101.3.4.3.12", SignatureAlgorithm_ECDSA_SHA3_512},
		{"0.4.0.127.0.7.1.1.4.1.1", SignatureAlgorithm_PLAIN_ECDSA_SHA1},
		{"0.4.0.127.0.7.1.1.4.1.2", SignatureAlgorithm_PLAIN_ECDSA_SHA224},
		{"0.4.0.127.0.7.1.1.4.1.3", SignatureAlgorithm_PLAIN_ECDSA_SHA256},
		{"0.4.0.127.0.7.1.1.4.1.4", SignatureAlgorithm_PLAIN_ECDSA_SHA384},
		{"0.4.0.127.0.7.1.1.4.1.5", SignatureAlgorithm_PLAIN_ECDSA_SHA512},
		{"0.4.0.127.0.7.1.1.4.1.8", SignatureAlgorithm_PLAIN_ECDSA_SHA3_224},
		{"0.4.0.127.0.7.1.1.4.1.9", SignatureAlgorithm_PLAIN_ECDSA_SHA3_256},
		{"0.4.0.127.0.7.1.1.4.1.10", SignatureAlgorithm_PLAIN_ECDSA_SHA3_384},
		{"0.4.0.127.0.7.1.1.4.1.11", SignatureAlgorithm_PLAIN_ECDSA_SHA3_512},
		{"1.3.101.112", SignatureAlgorithm_ED25519},
		{"1.3.101.113", SignatureAlgorithm_ED448},
		{"1.2.840.10040.4.3", SignatureAlgorithm_DSA_SHA1},
		{"1.2.14888.3.0.1", SignatureAlgorithm_DSA_SHA1},
		{"2.16.840.1.101.3.4.3.1", SignatureAlgorithm_DSA_SHA224},
		{"2.16.840.1.101.3.4.3.2", SignatureAlgorithm_DSA_SHA256},
		{"2.16.840.1.101.3.4.3.3", SignatureAlgorithm_DSA_SHA384},
		{"2.16.840.1.101.3.4.3.4", SignatureAlgorithm_DSA_SHA512},
		{"2.16.840.1.101.3.4.3.5", SignatureAlgorithm_DSA_SHA3_224},
		{"2.16.840.1.101.3.4.3.6", SignatureAlgorithm_DSA_SHA3_256},
		{"2.16.840.1.101.3.4.3.7", SignatureAlgorithm_DSA_SHA3_384},
		{"2.16.840.1.101.3.4.3.8", SignatureAlgorithm_DSA_SHA3_512},
		{"1.2.840.113549.2.7", SignatureAlgorithm_HMAC_SHA1},
		{"1.2.840.113549.2.8", SignatureAlgorithm_HMAC_SHA224},
		{"1.2.840.113549.2.9", SignatureAlgorithm_HMAC_SHA256},
		{"1.2.840.113549.2.10", SignatureAlgorithm_HMAC_SHA384},
		{"1.2.840.113549.2.11", SignatureAlgorithm_HMAC_SHA512},
		{"1.3.6.1.5.5.8.1.4", SignatureAlgorithm_HMAC_RIPEMD160},
		{"2.16.840.1.101.3.4.2.13", SignatureAlgorithm_HMAC_SHA3_224},
		{"2.16.840.1.101.3.4.2.14", SignatureAlgorithm_HMAC_SHA3_256},
		{"2.16.840.1.101.3.4.2.15", SignatureAlgorithm_HMAC_SHA3_384},
		{"2.16.840.1.101.3.4.2.16", SignatureAlgorithm_HMAC_SHA3_512},
		{"1.2.840.113549.1.1.10", SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1},
	}
	if len(cases) != len(signatureAlgorithmOIDPairs)-1 {
		// -1: the two entries at OID "...4.1.6" collapse into a single test case.
		t.Fatalf("test case count = %d, want %d (signatureAlgorithmOIDPairs length - 1)", len(cases), len(signatureAlgorithmOIDPairs)-1)
	}
	for _, c := range cases {
		got, err := SignatureAlgorithmForOID(c.oid)
		if err != nil {
			t.Fatalf("SignatureAlgorithmForOID(%q) returned error: %v", c.oid, err)
		}
		if got != c.want {
			t.Errorf("SignatureAlgorithmForOID(%q) = %q, want %q", c.oid, got, c.want)
		}
	}

	if _, err := SignatureAlgorithmForOID("9.9.9"); err == nil {
		t.Error("SignatureAlgorithmForOID(\"9.9.9\") expected error, got nil")
	}

	// ECDSA_RIPEMD160's OID entry was overwritten by PLAIN_ECDSA_RIPEMD160's later
	// registration at the same OID, exactly as in the upstream HashMap; it must not
	// have a resolvable OID.
	if got := SignatureAlgorithm_ECDSA_RIPEMD160.OID(); got != "" {
		t.Errorf("ECDSA_RIPEMD160.OID() = %q, want \"\" (overwritten by PLAIN_ECDSA_RIPEMD160)", got)
	}
	if got := SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160.OID(); got != "0.4.0.127.0.7.1.1.4.1.6" {
		t.Errorf("PLAIN_ECDSA_RIPEMD160.OID() = %q, want \"0.4.0.127.0.7.1.1.4.1.6\"", got)
	}
}

func TestSignatureAlgorithmURIBasedOnOID(t *testing.T) {
	if got, want := SignatureAlgorithm_RSA_SHA256.URIBasedOnOID(), "urn:oid:1.2.840.113549.1.1.11"; got != want {
		t.Errorf("RSA_SHA256.URIBasedOnOID() = %q, want %q", got, want)
	}
	// Upstream concatenates a Java null for algorithms without an OID, producing the
	// literal "urn:oid:null". Reproduced verbatim; see the note on URIBasedOnOID.
	for _, v := range []SignatureAlgorithm{
		SignatureAlgorithm_RSA_RAW,
		SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1,
		SignatureAlgorithm_ECDSA_RAW,
		SignatureAlgorithm_ECDSA_RIPEMD160,
		SignatureAlgorithm_DSA_RAW,
	} {
		if got, want := v.URIBasedOnOID(), "urn:oid:null"; got != want {
			t.Errorf("%s.URIBasedOnOID() = %q, want %q", string(v), got, want)
		}
	}
}

// TestSignatureAlgorithmCanonicalAliasesMatchJVM pins every reverse-lookup accessor for
// every algorithm to the values produced by running the upstream Java enum on a JDK. The
// interesting cases are algorithms with more than one alias, where the winner is decided
// by java.util.HashMap iteration order rather than declaration order.
func TestSignatureAlgorithmCanonicalAliasesMatchJVM(t *testing.T) {
	// name -> {URI, OID, JCEID, JWAID, COSEID}; "" means Java returned null.
	want := map[SignatureAlgorithm][5]string{
		SignatureAlgorithm_RSA_RIPEMD160:   {"http://www.w3.org/2001/04/xmldsig-more/rsa-ripemd160", "1.3.36.3.3.1.2", "RIPEMD160withRSA", "", ""},
		SignatureAlgorithm_RSA_SHA1:        {"http://www.w3.org/2000/09/xmldsig#rsa-sha1", "1.2.840.113549.1.1.5", "SHA1withRSA", "", ""},
		SignatureAlgorithm_DSA_SHA1:        {"http://www.w3.org/2000/09/xmldsig#dsa-sha1", "1.2.840.10040.4.3", "SHA1withDSA", "", ""},
		SignatureAlgorithm_ECDSA_SHA256:    {"http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256", "1.2.840.10045.4.3.2", "SHA256withECDSA", "ES256", "-7"},
		SignatureAlgorithm_ED25519:         {"http://www.w3.org/2021/04/xmldsig-more#eddsa-ed25519", "1.3.101.112", "Ed25519", "EdDSA", "-8"},
		SignatureAlgorithm_ED448:           {"http://www.w3.org/2021/04/xmldsig-more#eddsa-ed448", "1.3.101.113", "Ed448", "EdDSA", "-8"},
		SignatureAlgorithm_ECDSA_RIPEMD160: {"http://www.w3.org/2007/05/xmldsig-more#ecdsa-ripemd160", "", "RIPEMD160withECDSA", "", ""},
	}
	for v, w := range want {
		cose := ""
		if c, ok := v.COSEID(); ok {
			cose = strconv.FormatInt(c, 10)
		}
		got := [5]string{v.URI(), v.OID(), v.JCEID(), v.JWAID(), cose}
		if got != w {
			t.Errorf("%s: got %v, want %v", string(v), got, w)
		}
	}
}

// TestSignatureAlgorithmForJAVA exhaustively round-trips every JCE name declared in the
// Java source through SignatureAlgorithmForJAVA.
func TestSignatureAlgorithmForJAVA(t *testing.T) {
	for _, p := range signatureAlgorithmJCEPairs {
		got, err := SignatureAlgorithmForJAVA(p.key)
		if err != nil {
			t.Fatalf("SignatureAlgorithmForJAVA(%q) returned error: %v", p.key, err)
		}
		if got != p.algo {
			t.Errorf("SignatureAlgorithmForJAVA(%q) = %q, want %q", p.key, got, p.algo)
		}
		if jceID := got.JCEID(); jceID == "" {
			t.Errorf("%v.JCEID() is empty", got)
		}
	}
	if _, err := SignatureAlgorithmForJAVA("bogus"); err == nil {
		t.Error("SignatureAlgorithmForJAVA(\"bogus\") expected error, got nil")
	}
	if got, want := SignatureAlgorithm_RSA_SHA256.JCEID(), "SHA256withRSA"; got != want {
		t.Errorf("RSA_SHA256.JCEID() = %q, want %q", got, want)
	}
}

// TestSignatureAlgorithmForJWA exhaustively round-trips every JWA identifier declared in
// the Java source through SignatureAlgorithmForJWA.
func TestSignatureAlgorithmForJWA(t *testing.T) {
	cases := []struct {
		id   string
		want SignatureAlgorithm
	}{
		{"HS256", SignatureAlgorithm_HMAC_SHA256},
		{"HS384", SignatureAlgorithm_HMAC_SHA384},
		{"HS512", SignatureAlgorithm_HMAC_SHA512},
		{"RS256", SignatureAlgorithm_RSA_SHA256},
		{"RS384", SignatureAlgorithm_RSA_SHA384},
		{"RS512", SignatureAlgorithm_RSA_SHA512},
		{"ES256", SignatureAlgorithm_ECDSA_SHA256},
		{"ES384", SignatureAlgorithm_ECDSA_SHA384},
		{"ES512", SignatureAlgorithm_ECDSA_SHA512},
		{"PS256", SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1},
		{"PS384", SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
		{"PS512", SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
		{"EdDSA", SignatureAlgorithm_ED25519},
	}
	if len(cases) != len(signatureAlgorithmJWAPairs) {
		t.Fatalf("test case count = %d, want %d (signatureAlgorithmJWAPairs length)", len(cases), len(signatureAlgorithmJWAPairs))
	}
	for _, c := range cases {
		got, err := SignatureAlgorithmForJWA(c.id)
		if err != nil {
			t.Fatalf("SignatureAlgorithmForJWA(%q) returned error: %v", c.id, err)
		}
		if got != c.want {
			t.Errorf("SignatureAlgorithmForJWA(%q) = %q, want %q", c.id, got, c.want)
		}
	}

	if _, err := SignatureAlgorithmForJWA("bogus"); err == nil {
		t.Error("SignatureAlgorithmForJWA(\"bogus\") expected error, got nil")
	}
	if got := SignatureAlgorithmForJWADefault("bogus", SignatureAlgorithm_RSA_SHA256); got != SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("SignatureAlgorithmForJWADefault(\"bogus\", RSA_SHA256) = %q, want RSA_SHA256", got)
	}

	// PLAIN_ECDSA_SHA256 inherits ES256 via ensurePlainECDSA (reverse-only); it is not
	// reachable from the forward lookup.
	if got, want := SignatureAlgorithm_PLAIN_ECDSA_SHA256.JWAID(), "ES256"; got != want {
		t.Errorf("PLAIN_ECDSA_SHA256.JWAID() = %q, want %q", got, want)
	}
	// ED448 explicitly reports "EdDSA" as its JWA id even though "EdDSA" resolves to
	// ED25519 on the forward lookup.
	if got, want := SignatureAlgorithm_ED448.JWAID(), "EdDSA"; got != want {
		t.Errorf("ED448.JWAID() = %q, want %q", got, want)
	}
}

// TestSignatureAlgorithmForCOSE exhaustively round-trips every COSE algorithm key
// declared in the Java source through SignatureAlgorithmForCOSEDefault.
func TestSignatureAlgorithmForCOSE(t *testing.T) {
	cases := []struct {
		key  int64
		want SignatureAlgorithm
	}{
		{-257, SignatureAlgorithm_RSA_SHA256},
		{-258, SignatureAlgorithm_RSA_SHA384},
		{-259, SignatureAlgorithm_RSA_SHA512},
		{-37, SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1},
		{-38, SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
		{-39, SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
		{-7, SignatureAlgorithm_ECDSA_SHA256},
		{-35, SignatureAlgorithm_ECDSA_SHA384},
		{-36, SignatureAlgorithm_ECDSA_SHA512},
		{-8, SignatureAlgorithm_ED25519},
	}
	if len(cases) != len(signatureAlgorithmCOSEPairs) {
		t.Fatalf("test case count = %d, want %d (signatureAlgorithmCOSEPairs length)", len(cases), len(signatureAlgorithmCOSEPairs))
	}
	const sentinel = SignatureAlgorithm("")
	for _, c := range cases {
		got := SignatureAlgorithmForCOSEDefault(c.key, sentinel)
		if got != c.want {
			t.Errorf("SignatureAlgorithmForCOSEDefault(%d) = %q, want %q", c.key, got, c.want)
		}
	}

	if got := SignatureAlgorithmForCOSEDefault(999999, sentinel); got != sentinel {
		t.Errorf("SignatureAlgorithmForCOSEDefault(999999) = %q, want sentinel", got)
	}

	// PLAIN_ECDSA_SHA256 inherits COSE key -7 via ensurePlainECDSA (reverse-only).
	if got, ok := SignatureAlgorithm_PLAIN_ECDSA_SHA256.COSEID(); !ok || got != -7 {
		t.Errorf("PLAIN_ECDSA_SHA256.COSEID() = (%d, %v), want (-7, true)", got, ok)
	}
	// ED448 explicitly reports -8 as its COSE key even though -8 resolves to ED25519
	// on the forward lookup.
	if got, ok := SignatureAlgorithm_ED448.COSEID(); !ok || got != -8 {
		t.Errorf("ED448.COSEID() = (%d, %v), want (-8, true)", got, ok)
	}
	if _, ok := SignatureAlgorithm_RSA_MD2.COSEID(); ok {
		t.Error("RSA_MD2.COSEID() ok = true, want false (no COSE mapping)")
	}
}

func TestSignatureAlgorithmGetAlgorithm(t *testing.T) {
	for _, v := range SignatureAlgorithmValues() {
		f := signatureAlgorithmData[v]
		got := SignatureAlgorithmGetAlgorithm(f.encryption, f.digest)
		if got != v {
			t.Errorf("SignatureAlgorithmGetAlgorithm(%q, %q) = %q, want %q", f.encryption, f.digest, got, v)
		}
	}
	if got := SignatureAlgorithmGetAlgorithm(EncryptionAlgorithm("BOGUS"), ""); got != "" {
		t.Errorf("SignatureAlgorithmGetAlgorithm(bogus) = %q, want \"\"", got)
	}
}

// signatureAlgorithmPSSParams holds RSASSA-PSS-params DER encodings produced by
// BouncyCastle 1.78.1's org.bouncycastle.asn1.pkcs.RSASSAPSSparams, one per hash algorithm.
// BouncyCastle omits the DEFAULT fields, so the SHA-1 entry carries no hashAlgorithm at all -
// which is exactly the case the DEFAULT sha1 rule has to cover.
var signatureAlgorithmPSSParams = []struct {
	name string
	der  string
	want SignatureAlgorithm
}{
	{"sha1 (defaulted)", "3005a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1},
	{"all defaults", "3000", SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1},
	{"sha224", "3034a00f300d06096086480165030402040500a11c301a06092a864886f70d010108300d06096086480165030402040500a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1},
	{"sha256", "3034a00f300d06096086480165030402010500a11c301a06092a864886f70d010108300d06096086480165030402010500a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1},
	{"sha384", "3034a00f300d06096086480165030402020500a11c301a06092a864886f70d010108300d06096086480165030402020500a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
	{"sha512", "3034a00f300d06096086480165030402030500a11c301a06092a864886f70d010108300d06096086480165030402030500a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
	{"sha3-256", "3034a00f300d06096086480165030402080500a11c301a06092a864886f70d010108300d06096086480165030402080500a203020120", SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1},
	// Produced by the JDK itself: AlgorithmParameters.getInstance("PSS").init(new
	// PSSParameterSpec(...)).getEncoded() on OpenJDK 21, i.e. exactly the byte sequence
	// upstream's AlgorithmParameters#init(byte[]) round-trips. The salt lengths differ from
	// the BouncyCastle rows above, which is what makes them worth pinning separately.
	{"jdk sha224", "3034a00f300d06096086480165030402040500a11c301a06092a864886f70d010108300d06096086480165030402040500a20302011c", SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1},
	{"jdk sha384", "3034a00f300d06096086480165030402020500a11c301a06092a864886f70d010108300d06096086480165030402020500a203020130", SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
	{"jdk sha512", "3034a00f300d06096086480165030402030500a11c301a06092a864886f70d010108300d06096086480165030402030500a203020140", SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
	{"jdk sha3-512", "3034a00f300d060960864801650304020a0500a11c301a06092a864886f70d010108300d060960864801650304020a0500a203020140", SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1},
}

// TestSignatureAlgorithmForOIDAndParamsPSS checks that the RSASSA-PSS OID resolves to the
// digest named by the RSASSA-PSS-params, which is what upstream reads out of a
// PSSParameterSpec.
func TestSignatureAlgorithmForOIDAndParamsPSS(t *testing.T) {
	const pssOID = "1.2.840.113549.1.1.10"
	for _, entry := range signatureAlgorithmPSSParams {
		params, err := hex.DecodeString(entry.der)
		if err != nil {
			t.Fatalf("%s: %v", entry.name, err)
		}
		got, err := SignatureAlgorithmForOIDAndParams(pssOID, params)
		if err != nil {
			t.Fatalf("%s: %v", entry.name, err)
		}
		if got != entry.want {
			t.Errorf("%s: got %q, want %q", entry.name, got, entry.want)
		}
	}

	// Without parameters the OID keeps its nominal algorithm, as upstream does.
	got, err := SignatureAlgorithmForOIDAndParams(pssOID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1 {
		t.Errorf("nil params: got %q, want %q", got, SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1)
	}
	if got, err := SignatureAlgorithmForOID(pssOID); err != nil || got != SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1 {
		t.Errorf("forOID: got (%q, %v)", got, err)
	}

	// Parameters are only read for the RSASSA-PSS encryption algorithm.
	got, err = SignatureAlgorithmForOIDAndParams("1.2.840.113549.1.1.11", []byte{0x05, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if got != SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("RSA-SHA256: got %q", got)
	}
}

// TestSignatureAlgorithmForOIDAndParamsPSSErrors checks the failures upstream reports as
// IllegalArgumentException("Unable to initialize PSS").
func TestSignatureAlgorithmForOIDAndParamsPSSErrors(t *testing.T) {
	const pssOID = "1.2.840.113549.1.1.10"
	for _, entry := range []struct {
		name string
		der  string
	}{
		{"not a sequence", "0500"},
		{"trailing data", "30000500"},
		{"hashAlgorithm is not an AlgorithmIdentifier", "3004a0020500"},
		{"hashAlgorithm without an OID", "3004a0023000"},
		{"unknown digest OID", "300ea00c300a06082a864886f70d0203"},
		{"truncated", "3005a0"},
		// sun.security.rsa.PSSParameters rejects these two the same way, which upstream
		// reports as IllegalArgumentException("Unable to initialize PSS").
		{"trailerField other than 1", "3005a303020102"},
		{"maskGenAlgorithm that is not MGF1", "300da10b3009060706052b0e03021a"},
	} {
		params, err := hex.DecodeString(entry.der)
		if err != nil {
			t.Fatalf("%s: %v", entry.name, err)
		}
		if _, err := SignatureAlgorithmForOIDAndParams(pssOID, params); err == nil {
			t.Errorf("%s: expected an error", entry.name)
		}
	}
	// An unknown OID is rejected before the parameters are looked at.
	if _, err := SignatureAlgorithmForOIDAndParams("1.2.3.4", nil); err == nil {
		t.Error("an unknown OID must be rejected")
	}
}
