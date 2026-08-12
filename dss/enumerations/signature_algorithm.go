// Ported from dss-enumerations/.../SignatureAlgorithm.java (DSS 6.5.RC1).
//
// Upstream builds its reverse lookup tables (the OID/URI/JCE-name/JWA/COSE "for this
// algorithm" getters) by iterating java.util.HashMap entry sets and doing put(value, key).
// Where more than one alias maps to the same SignatureAlgorithm (e.g. two XML
// RSA-RIPEMD160 URIs, two RSA_SHA1/DSA_SHA1 OIDs), the alias that wins is therefore the
// one appearing LAST in java.util.HashMap iteration order. That order is not arbitrary:
// String/Long hashCode and HashMap's bucket layout are specified by the Java platform and
// the insertion sequence is fixed in the upstream source, so it is fully reproducible.
// javaHashMapOrder below reproduces it, which makes URI()/OID()/JCEID()/JWAID()/COSEID()
// byte-identical to Java (verified against a JDK 21 run of the upstream enum for all 66
// algorithms). Picking the first alias in declaration order instead would be wrong for
// RSA_RIPEMD160, whose canonical URI is the ".../xmldsig-more/rsa-ripemd160" spelling and
// not the "#rsa-ripemd160" one declared first. The forward lookups
// (SignatureAlgorithmForXML/ForOID/ForJWA/ForCOSE) accept every alias, matching Java.
// One upstream OID collision is preserved verbatim and is NOT a porting error:
// "0.4.0.127.0.7.1.1.4.1.6" is registered for both ECDSA_RIPEMD160 and
// PLAIN_ECDSA_RIPEMD160 in that OID table order, so the later entry
// (PLAIN_ECDSA_RIPEMD160) wins both the forward OID lookup and OID(); ECDSA_RIPEMD160
// therefore has no OID (OID() returns "") exactly as in the Java HashMap.put overwrite.

package enumerations

import (
	encoding_asn1 "encoding/asn1"
	"errors"
	"fmt"
	"unicode/utf16"

	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// SignatureAlgorithm lists supported signature algorithms.
//
// Implements OidAndUriBasedEnum.
type SignatureAlgorithm string

const (
	// SignatureAlgorithm_RSA_RAW is RSA without digest algorithm.
	SignatureAlgorithm_RSA_RAW SignatureAlgorithm = "RSA_RAW"
	// SignatureAlgorithm_RSA_SHA1 is RSA with SHA-1.
	SignatureAlgorithm_RSA_SHA1 SignatureAlgorithm = "RSA_SHA1"
	// SignatureAlgorithm_RSA_SHA224 is RSA with SHA-224.
	SignatureAlgorithm_RSA_SHA224 SignatureAlgorithm = "RSA_SHA224"
	// SignatureAlgorithm_RSA_SHA256 is RSA with SHA-256.
	SignatureAlgorithm_RSA_SHA256 SignatureAlgorithm = "RSA_SHA256"
	// SignatureAlgorithm_RSA_SHA384 is RSA with SHA-384.
	SignatureAlgorithm_RSA_SHA384 SignatureAlgorithm = "RSA_SHA384"
	// SignatureAlgorithm_RSA_SHA512 is RSA with SHA-512.
	SignatureAlgorithm_RSA_SHA512 SignatureAlgorithm = "RSA_SHA512"
	// SignatureAlgorithm_RSA_SHA3_224 is RSA with SHA3-224.
	SignatureAlgorithm_RSA_SHA3_224 SignatureAlgorithm = "RSA_SHA3_224"
	// SignatureAlgorithm_RSA_SHA3_256 is RSA with SHA3-256.
	SignatureAlgorithm_RSA_SHA3_256 SignatureAlgorithm = "RSA_SHA3_256"
	// SignatureAlgorithm_RSA_SHA3_384 is RSA with SHA3-384.
	SignatureAlgorithm_RSA_SHA3_384 SignatureAlgorithm = "RSA_SHA3_384"
	// SignatureAlgorithm_RSA_SHA3_512 is RSA with SHA3-512.
	SignatureAlgorithm_RSA_SHA3_512 SignatureAlgorithm = "RSA_SHA3_512"
	// SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1 is RSA with MGF1 without digest algorithm.
	SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_RAW_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1 is RSA with MGF1 with SHA-1.
	SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA1_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1 is RSA with MGF1 with SHA-224.
	SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA224_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1 is RSA with MGF1 with SHA-256.
	SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA256_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1 is RSA with MGF1 with SHA-384.
	SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA384_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1 is RSA with MGF1 with SHA-512.
	SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA512_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1 is RSA with MGF1 with SHA3-224.
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA3_224_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1 is RSA with MGF1 with SHA3-256.
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA3_256_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1 is RSA with MGF1 with SHA3-384.
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA3_384_MGF1"
	// SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1 is RSA with MGF1 with SHA3-512.
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1 SignatureAlgorithm = "RSA_SSA_PSS_SHA3_512_MGF1"
	// SignatureAlgorithm_RSA_RIPEMD160 is RSA with RIPEMD160.
	SignatureAlgorithm_RSA_RIPEMD160 SignatureAlgorithm = "RSA_RIPEMD160"
	// SignatureAlgorithm_RSA_MD2 is RSA with MD2.
	SignatureAlgorithm_RSA_MD2 SignatureAlgorithm = "RSA_MD2"
	// SignatureAlgorithm_RSA_MD5 is RSA with MD5.
	SignatureAlgorithm_RSA_MD5 SignatureAlgorithm = "RSA_MD5"
	// SignatureAlgorithm_ECDSA_RAW is ECDSA without digest algorithm.
	SignatureAlgorithm_ECDSA_RAW SignatureAlgorithm = "ECDSA_RAW"
	// SignatureAlgorithm_ECDSA_SHA1 is ECDSA with SHA-1.
	SignatureAlgorithm_ECDSA_SHA1 SignatureAlgorithm = "ECDSA_SHA1"
	// SignatureAlgorithm_ECDSA_SHA224 is ECDSA with SHA-224.
	SignatureAlgorithm_ECDSA_SHA224 SignatureAlgorithm = "ECDSA_SHA224"
	// SignatureAlgorithm_ECDSA_SHA256 is ECDSA with SHA-256.
	SignatureAlgorithm_ECDSA_SHA256 SignatureAlgorithm = "ECDSA_SHA256"
	// SignatureAlgorithm_ECDSA_SHA384 is ECDSA with SHA-384.
	SignatureAlgorithm_ECDSA_SHA384 SignatureAlgorithm = "ECDSA_SHA384"
	// SignatureAlgorithm_ECDSA_SHA512 is ECDSA with SHA-512.
	SignatureAlgorithm_ECDSA_SHA512 SignatureAlgorithm = "ECDSA_SHA512"
	// SignatureAlgorithm_ECDSA_SHA3_224 is ECDSA with SHA3-224.
	SignatureAlgorithm_ECDSA_SHA3_224 SignatureAlgorithm = "ECDSA_SHA3_224"
	// SignatureAlgorithm_ECDSA_SHA3_256 is ECDSA with SHA3-256.
	SignatureAlgorithm_ECDSA_SHA3_256 SignatureAlgorithm = "ECDSA_SHA3_256"
	// SignatureAlgorithm_ECDSA_SHA3_384 is ECDSA with SHA3-384.
	SignatureAlgorithm_ECDSA_SHA3_384 SignatureAlgorithm = "ECDSA_SHA3_384"
	// SignatureAlgorithm_ECDSA_SHA3_512 is ECDSA with SHA3-512.
	SignatureAlgorithm_ECDSA_SHA3_512 SignatureAlgorithm = "ECDSA_SHA3_512"
	// SignatureAlgorithm_ECDSA_RIPEMD160 is ECDSA with RIPEMD160.
	SignatureAlgorithm_ECDSA_RIPEMD160 SignatureAlgorithm = "ECDSA_RIPEMD160"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA1 is PLAIN-ECDSA with SHA-1.
	SignatureAlgorithm_PLAIN_ECDSA_SHA1 SignatureAlgorithm = "PLAIN_ECDSA_SHA1"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA224 is PLAIN-ECDSA with SHA-224.
	SignatureAlgorithm_PLAIN_ECDSA_SHA224 SignatureAlgorithm = "PLAIN_ECDSA_SHA224"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA256 is PLAIN-ECDSA with SHA-256.
	SignatureAlgorithm_PLAIN_ECDSA_SHA256 SignatureAlgorithm = "PLAIN_ECDSA_SHA256"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA384 is PLAIN-ECDSA with SHA-384.
	SignatureAlgorithm_PLAIN_ECDSA_SHA384 SignatureAlgorithm = "PLAIN_ECDSA_SHA384"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA512 is PLAIN-ECDSA with SHA-512.
	SignatureAlgorithm_PLAIN_ECDSA_SHA512 SignatureAlgorithm = "PLAIN_ECDSA_SHA512"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA3_224 is PLAIN-ECDSA with SHA3-224.
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_224 SignatureAlgorithm = "PLAIN_ECDSA_SHA3_224"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA3_256 is PLAIN-ECDSA with SHA3-256.
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_256 SignatureAlgorithm = "PLAIN_ECDSA_SHA3_256"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA3_384 is PLAIN-ECDSA with SHA3-384.
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_384 SignatureAlgorithm = "PLAIN_ECDSA_SHA3_384"
	// SignatureAlgorithm_PLAIN_ECDSA_SHA3_512 is PLAIN-ECDSA with SHA3-512.
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_512 SignatureAlgorithm = "PLAIN_ECDSA_SHA3_512"
	// SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160 is PLAIN-ECDSA with RIPEMD160.
	SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160 SignatureAlgorithm = "PLAIN_ECDSA_RIPEMD160"
	// SignatureAlgorithm_DSA_RAW is DSA without digest algorithm.
	SignatureAlgorithm_DSA_RAW SignatureAlgorithm = "DSA_RAW"
	// SignatureAlgorithm_DSA_SHA1 is DSA with SHA-1.
	SignatureAlgorithm_DSA_SHA1 SignatureAlgorithm = "DSA_SHA1"
	// SignatureAlgorithm_DSA_SHA224 is DSA with SHA-224.
	SignatureAlgorithm_DSA_SHA224 SignatureAlgorithm = "DSA_SHA224"
	// SignatureAlgorithm_DSA_SHA256 is DSA with SHA-256.
	SignatureAlgorithm_DSA_SHA256 SignatureAlgorithm = "DSA_SHA256"
	// SignatureAlgorithm_DSA_SHA384 is DSA with SHA-384.
	SignatureAlgorithm_DSA_SHA384 SignatureAlgorithm = "DSA_SHA384"
	// SignatureAlgorithm_DSA_SHA512 is DSA with SHA-512.
	SignatureAlgorithm_DSA_SHA512 SignatureAlgorithm = "DSA_SHA512"
	// SignatureAlgorithm_DSA_SHA3_224 is DSA with SHA3-224.
	SignatureAlgorithm_DSA_SHA3_224 SignatureAlgorithm = "DSA_SHA3_224"
	// SignatureAlgorithm_DSA_SHA3_256 is DSA with SHA3-256.
	SignatureAlgorithm_DSA_SHA3_256 SignatureAlgorithm = "DSA_SHA3_256"
	// SignatureAlgorithm_DSA_SHA3_384 is DSA with SHA3-384.
	SignatureAlgorithm_DSA_SHA3_384 SignatureAlgorithm = "DSA_SHA3_384"
	// SignatureAlgorithm_DSA_SHA3_512 is DSA with SHA3-512.
	SignatureAlgorithm_DSA_SHA3_512 SignatureAlgorithm = "DSA_SHA3_512"
	// SignatureAlgorithm_HMAC_SHA1 is HMAC with SHA-1.
	SignatureAlgorithm_HMAC_SHA1 SignatureAlgorithm = "HMAC_SHA1"
	// SignatureAlgorithm_HMAC_SHA224 is HMAC with SHA-224.
	SignatureAlgorithm_HMAC_SHA224 SignatureAlgorithm = "HMAC_SHA224"
	// SignatureAlgorithm_HMAC_SHA256 is HMAC with SHA-256.
	SignatureAlgorithm_HMAC_SHA256 SignatureAlgorithm = "HMAC_SHA256"
	// SignatureAlgorithm_HMAC_SHA384 is HMAC with SHA-384.
	SignatureAlgorithm_HMAC_SHA384 SignatureAlgorithm = "HMAC_SHA384"
	// SignatureAlgorithm_HMAC_SHA512 is HMAC with SHA-512.
	SignatureAlgorithm_HMAC_SHA512 SignatureAlgorithm = "HMAC_SHA512"
	// SignatureAlgorithm_HMAC_SHA3_224 is HMAC with SHA3-224.
	SignatureAlgorithm_HMAC_SHA3_224 SignatureAlgorithm = "HMAC_SHA3_224"
	// SignatureAlgorithm_HMAC_SHA3_256 is HMAC with SHA3-256.
	SignatureAlgorithm_HMAC_SHA3_256 SignatureAlgorithm = "HMAC_SHA3_256"
	// SignatureAlgorithm_HMAC_SHA3_384 is HMAC with SHA3-384.
	SignatureAlgorithm_HMAC_SHA3_384 SignatureAlgorithm = "HMAC_SHA3_384"
	// SignatureAlgorithm_HMAC_SHA3_512 is HMAC with SHA3-512.
	SignatureAlgorithm_HMAC_SHA3_512 SignatureAlgorithm = "HMAC_SHA3_512"
	// SignatureAlgorithm_HMAC_RIPEMD160 is HMAC with RIPEMD160.
	SignatureAlgorithm_HMAC_RIPEMD160 SignatureAlgorithm = "HMAC_RIPEMD160"
	// SignatureAlgorithm_ED25519 is EDDSA with SHA512 (RFC 8419 section 3.1).
	SignatureAlgorithm_ED25519 SignatureAlgorithm = "ED25519"
	// SignatureAlgorithm_ED448 is EDDSA with SHAKE256-512.
	SignatureAlgorithm_ED448 SignatureAlgorithm = "ED448"
)

// signatureAlgorithmOIDNamespacePrefix is the OID URI prefix (RFC 3061).
const signatureAlgorithmOIDNamespacePrefix = "urn:oid:"

// signatureAlgorithmFields holds the (encryptionAlgorithm, digestAlgorithm) pair for a
// SignatureAlgorithm. digest == "" represents Java's null (a "RAW"/no-digest variant).
type signatureAlgorithmFields struct {
	encryption EncryptionAlgorithm
	digest     DigestAlgorithm
}

// signatureAlgorithmData holds the (encryptionAlgorithm, digestAlgorithm) pair for every
// SignatureAlgorithm, in the exact combinations declared by the Java enum constructors.
var signatureAlgorithmData = map[SignatureAlgorithm]signatureAlgorithmFields{
	SignatureAlgorithm_RSA_RAW:                   {EncryptionAlgorithm_RSA, ""},
	SignatureAlgorithm_RSA_SHA1:                  {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA1},
	SignatureAlgorithm_RSA_SHA224:                {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA224},
	SignatureAlgorithm_RSA_SHA256:                {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA256},
	SignatureAlgorithm_RSA_SHA384:                {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA384},
	SignatureAlgorithm_RSA_SHA512:                {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA512},
	SignatureAlgorithm_RSA_SHA3_224:              {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_RSA_SHA3_256:              {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_RSA_SHA3_384:              {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_RSA_SHA3_512:              {EncryptionAlgorithm_RSA, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1:      {EncryptionAlgorithm_RSASSA_PSS, ""},
	SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1:     {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA1},
	SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1:   {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA224},
	SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1:   {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA256},
	SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1:   {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA384},
	SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1:   {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA512},
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1: {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1: {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1: {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1: {EncryptionAlgorithm_RSASSA_PSS, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_RSA_RIPEMD160:             {EncryptionAlgorithm_RSA, DigestAlgorithm_RIPEMD160},
	SignatureAlgorithm_RSA_MD2:                   {EncryptionAlgorithm_RSA, DigestAlgorithm_MD2},
	SignatureAlgorithm_RSA_MD5:                   {EncryptionAlgorithm_RSA, DigestAlgorithm_MD5},
	SignatureAlgorithm_ECDSA_RAW:                 {EncryptionAlgorithm_ECDSA, ""},
	SignatureAlgorithm_ECDSA_SHA1:                {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA1},
	SignatureAlgorithm_ECDSA_SHA224:              {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA224},
	SignatureAlgorithm_ECDSA_SHA256:              {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA256},
	SignatureAlgorithm_ECDSA_SHA384:              {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA384},
	SignatureAlgorithm_ECDSA_SHA512:              {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA512},
	SignatureAlgorithm_ECDSA_SHA3_224:            {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_ECDSA_SHA3_256:            {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_ECDSA_SHA3_384:            {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_ECDSA_SHA3_512:            {EncryptionAlgorithm_ECDSA, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_ECDSA_RIPEMD160:           {EncryptionAlgorithm_ECDSA, DigestAlgorithm_RIPEMD160},
	SignatureAlgorithm_PLAIN_ECDSA_SHA1:          {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA1},
	SignatureAlgorithm_PLAIN_ECDSA_SHA224:        {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA224},
	SignatureAlgorithm_PLAIN_ECDSA_SHA256:        {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA256},
	SignatureAlgorithm_PLAIN_ECDSA_SHA384:        {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA384},
	SignatureAlgorithm_PLAIN_ECDSA_SHA512:        {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA512},
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_224:      {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_256:      {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_384:      {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_PLAIN_ECDSA_SHA3_512:      {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160:     {EncryptionAlgorithm_PLAIN_ECDSA, DigestAlgorithm_RIPEMD160},
	SignatureAlgorithm_DSA_RAW:                   {EncryptionAlgorithm_DSA, ""},
	SignatureAlgorithm_DSA_SHA1:                  {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA1},
	SignatureAlgorithm_DSA_SHA224:                {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA224},
	SignatureAlgorithm_DSA_SHA256:                {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA256},
	SignatureAlgorithm_DSA_SHA384:                {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA384},
	SignatureAlgorithm_DSA_SHA512:                {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA512},
	SignatureAlgorithm_DSA_SHA3_224:              {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_DSA_SHA3_256:              {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_DSA_SHA3_384:              {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_DSA_SHA3_512:              {EncryptionAlgorithm_DSA, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_HMAC_SHA1:                 {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA1},
	SignatureAlgorithm_HMAC_SHA224:               {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA224},
	SignatureAlgorithm_HMAC_SHA256:               {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA256},
	SignatureAlgorithm_HMAC_SHA384:               {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA384},
	SignatureAlgorithm_HMAC_SHA512:               {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA512},
	SignatureAlgorithm_HMAC_SHA3_224:             {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_224},
	SignatureAlgorithm_HMAC_SHA3_256:             {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_256},
	SignatureAlgorithm_HMAC_SHA3_384:             {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_384},
	SignatureAlgorithm_HMAC_SHA3_512:             {EncryptionAlgorithm_HMAC, DigestAlgorithm_SHA3_512},
	SignatureAlgorithm_HMAC_RIPEMD160:            {EncryptionAlgorithm_HMAC, DigestAlgorithm_RIPEMD160},
	SignatureAlgorithm_ED25519:                   {EncryptionAlgorithm_EDDSA, DigestAlgorithm_SHA512},
	SignatureAlgorithm_ED448:                     {EncryptionAlgorithm_EDDSA, DigestAlgorithm_SHAKE256_512},
}

// SignatureAlgorithmValues returns all SignatureAlgorithm constants in declaration order.
func SignatureAlgorithmValues() []SignatureAlgorithm {
	return []SignatureAlgorithm{
		SignatureAlgorithm_RSA_RAW,
		SignatureAlgorithm_RSA_SHA1,
		SignatureAlgorithm_RSA_SHA224,
		SignatureAlgorithm_RSA_SHA256,
		SignatureAlgorithm_RSA_SHA384,
		SignatureAlgorithm_RSA_SHA512,
		SignatureAlgorithm_RSA_SHA3_224,
		SignatureAlgorithm_RSA_SHA3_256,
		SignatureAlgorithm_RSA_SHA3_384,
		SignatureAlgorithm_RSA_SHA3_512,
		SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1,
		SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1,
		SignatureAlgorithm_RSA_RIPEMD160,
		SignatureAlgorithm_RSA_MD2,
		SignatureAlgorithm_RSA_MD5,
		SignatureAlgorithm_ECDSA_RAW,
		SignatureAlgorithm_ECDSA_SHA1,
		SignatureAlgorithm_ECDSA_SHA224,
		SignatureAlgorithm_ECDSA_SHA256,
		SignatureAlgorithm_ECDSA_SHA384,
		SignatureAlgorithm_ECDSA_SHA512,
		SignatureAlgorithm_ECDSA_SHA3_224,
		SignatureAlgorithm_ECDSA_SHA3_256,
		SignatureAlgorithm_ECDSA_SHA3_384,
		SignatureAlgorithm_ECDSA_SHA3_512,
		SignatureAlgorithm_ECDSA_RIPEMD160,
		SignatureAlgorithm_PLAIN_ECDSA_SHA1,
		SignatureAlgorithm_PLAIN_ECDSA_SHA224,
		SignatureAlgorithm_PLAIN_ECDSA_SHA256,
		SignatureAlgorithm_PLAIN_ECDSA_SHA384,
		SignatureAlgorithm_PLAIN_ECDSA_SHA512,
		SignatureAlgorithm_PLAIN_ECDSA_SHA3_224,
		SignatureAlgorithm_PLAIN_ECDSA_SHA3_256,
		SignatureAlgorithm_PLAIN_ECDSA_SHA3_384,
		SignatureAlgorithm_PLAIN_ECDSA_SHA3_512,
		SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160,
		SignatureAlgorithm_DSA_RAW,
		SignatureAlgorithm_DSA_SHA1,
		SignatureAlgorithm_DSA_SHA224,
		SignatureAlgorithm_DSA_SHA256,
		SignatureAlgorithm_DSA_SHA384,
		SignatureAlgorithm_DSA_SHA512,
		SignatureAlgorithm_DSA_SHA3_224,
		SignatureAlgorithm_DSA_SHA3_256,
		SignatureAlgorithm_DSA_SHA3_384,
		SignatureAlgorithm_DSA_SHA3_512,
		SignatureAlgorithm_HMAC_SHA1,
		SignatureAlgorithm_HMAC_SHA224,
		SignatureAlgorithm_HMAC_SHA256,
		SignatureAlgorithm_HMAC_SHA384,
		SignatureAlgorithm_HMAC_SHA512,
		SignatureAlgorithm_HMAC_SHA3_224,
		SignatureAlgorithm_HMAC_SHA3_256,
		SignatureAlgorithm_HMAC_SHA3_384,
		SignatureAlgorithm_HMAC_SHA3_512,
		SignatureAlgorithm_HMAC_RIPEMD160,
		SignatureAlgorithm_ED25519,
		SignatureAlgorithm_ED448,
	}
}

// SignatureAlgorithmValueOf returns the SignatureAlgorithm matching the given Java enum name.
func SignatureAlgorithmValueOf(name string) (SignatureAlgorithm, error) {
	for _, v := range SignatureAlgorithmValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignatureAlgorithm.%s", name)
}

// EncryptionAlgorithm returns the encryption algorithm.
func (s SignatureAlgorithm) EncryptionAlgorithm() EncryptionAlgorithm {
	return signatureAlgorithmData[s].encryption
}

// DigestAlgorithm returns the digest algorithm, or "" for the RAW/no-digest variants
// (RSA_RAW, RSA_SSA_PSS_RAW_MGF1, ECDSA_RAW, DSA_RAW — Java's null).
func (s SignatureAlgorithm) DigestAlgorithm() DigestAlgorithm {
	return signatureAlgorithmData[s].digest
}

// Name returns a user-friendly name for the signature algorithm.
func (s SignatureAlgorithm) Name() string {
	name := s.EncryptionAlgorithm().Name()
	if d := s.DigestAlgorithm(); d != "" {
		name += " with " + d.Name()
	}
	return name
}

// signatureAlgorithmPair is a (lookup key, algorithm) pair, used to build both the
// forward alias->algorithm maps and, deterministically, their reverse.
type signatureAlgorithmPair[K comparable] struct {
	key  K
	algo SignatureAlgorithm
}

// signatureAlgorithmForward builds a forward lookup map from an ordered pair list,
// replicating java.util.HashMap.put semantics: a later pair with an equal key
// overwrites an earlier one.
func signatureAlgorithmForward[K comparable](pairs []signatureAlgorithmPair[K]) map[K]SignatureAlgorithm {
	m := make(map[K]SignatureAlgorithm, len(pairs))
	for _, p := range pairs {
		m[p.key] = p.algo
	}
	return m
}

// javaStringHashCode reproduces java.lang.String.hashCode(): s[0]*31^(n-1) + ... + s[n-1],
// computed over UTF-16 code units with int32 (two's complement) overflow.
func javaStringHashCode(s string) int32 {
	var h int32
	for _, u := range utf16.Encode([]rune(s)) {
		h = 31*h + int32(u)
	}
	return h
}

// javaLongHashCode reproduces java.lang.Long.hashCode(): (int)(value ^ (value >>> 32)).
func javaLongHashCode(v int64) int32 {
	u := uint64(v)
	return int32(uint32(u ^ (u >> 32)))
}

// javaHashMapOrder returns keys in java.util.HashMap iteration order, given the order in
// which they were inserted. HashMap iterates its bucket table from index 0 upwards and,
// within a bucket, follows the chain in insertion order (resize splits preserve relative
// order, and none of these tables is large enough to treeify). The table capacity is the
// smallest power of two >= 16 whose 0.75 load-factor threshold accommodates the entry
// count, and a key's bucket is spread(hash) & (capacity-1) where spread(h) = h ^ (h>>>16).
//
// Upstream builds every "for this algorithm" reverse map by iterating a HashMap entry set
// and doing put(value, key), so the winning alias for an algorithm with several aliases is
// whichever comes LAST in this order. Reproducing the order exactly is what makes URI(),
// OID(), JCEID(), JWAID() and COSEID() byte-identical to Java.
func javaHashMapOrder[K comparable](keys []K, hash func(K) int32) []K {
	distinct := make([]K, 0, len(keys))
	seen := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		distinct = append(distinct, k)
	}
	capacity := 16
	for float64(len(distinct)) > 0.75*float64(capacity) {
		capacity *= 2
	}
	buckets := make([][]K, capacity)
	for _, k := range distinct {
		h := uint32(hash(k))
		idx := int((h ^ (h >> 16)) & uint32(capacity-1))
		buckets[idx] = append(buckets[idx], k)
	}
	out := make([]K, 0, len(distinct))
	for _, b := range buckets {
		out = append(out, b...)
	}
	return out
}

// signatureAlgorithmReverse builds the reverse (algorithm->canonical alias) map, mirroring
// upstream's register*ForKey loops: iterate the forward HashMap's entries in HashMap order
// and put(value, key), so for an algorithm with several aliases the LAST one in that order
// wins. When plainECDSA is set it also mirrors upstream's ensurePlainECDSA call inside the
// same loop iteration, which maps each ECDSA_* alias onto the PLAIN_ECDSA_* algorithm with
// the same digest.
func signatureAlgorithmReverse[K comparable](pairs []signatureAlgorithmPair[K], forward map[K]SignatureAlgorithm, hash func(K) int32, plainECDSA bool) map[SignatureAlgorithm]K {
	keys := make([]K, 0, len(pairs))
	for _, p := range pairs {
		keys = append(keys, p.key)
	}
	m := make(map[SignatureAlgorithm]K, len(pairs))
	for _, k := range javaHashMapOrder(keys, hash) {
		algo := forward[k]
		m[algo] = k
		if !plainECDSA {
			continue
		}
		if signatureAlgorithmData[algo].encryption != EncryptionAlgorithm_ECDSA {
			continue
		}
		if plain := signatureAlgorithmGetAlgorithm(EncryptionAlgorithm_PLAIN_ECDSA, signatureAlgorithmData[algo].digest); plain != "" {
			m[plain] = k
		}
	}
	return m
}

// signatureAlgorithmGetAlgorithm finds the SignatureAlgorithm for the given
// (encryption, digest) pair, or "" if none matches (Java's null).
func signatureAlgorithmGetAlgorithm(encryptionAlgorithm EncryptionAlgorithm, digestAlgorithm DigestAlgorithm) SignatureAlgorithm {
	for _, v := range SignatureAlgorithmValues() {
		f := signatureAlgorithmData[v]
		if f.encryption == encryptionAlgorithm && f.digest == digestAlgorithm {
			return v
		}
	}
	return ""
}

// SignatureAlgorithmGetAlgorithm returns, for the given encryption and digest
// algorithm, the corresponding SignatureAlgorithm, or "" if no combination matches.
func SignatureAlgorithmGetAlgorithm(encryptionAlgorithm EncryptionAlgorithm, digestAlgorithm DigestAlgorithm) SignatureAlgorithm {
	return signatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
}

// signatureAlgorithmXMLPairs lists the XML algorithm URIs
// (http://www.w3.org/TR/2013/NOTE-xmlsec-algorithms-20130411/) in upstream declaration order.
var signatureAlgorithmXMLPairs = []signatureAlgorithmPair[string]{
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
	// Support of not standard AT algorithm name; see
	// http://www.rfc-editor.org/errata_search.php?rfc=4051
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

var signatureAlgorithmXMLForward = signatureAlgorithmForward(signatureAlgorithmXMLPairs)
var signatureAlgorithmXMLReverse = func() map[SignatureAlgorithm]string {
	return signatureAlgorithmReverse(signatureAlgorithmXMLPairs, signatureAlgorithmXMLForward, javaStringHashCode, true)
}()

// signatureAlgorithmOIDPairs lists the signature algorithm OIDs in upstream declaration
// order, including the upstream ECDSA_RIPEMD160/PLAIN_ECDSA_RIPEMD160 OID collision
// documented in the file header.
var signatureAlgorithmOIDPairs = []signatureAlgorithmPair[string]{
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
	{"0.4.0.127.0.7.1.1.4.1.6", SignatureAlgorithm_ECDSA_RIPEMD160},
	{"2.16.840.1.101.3.4.3.9", SignatureAlgorithm_ECDSA_SHA3_224},
	{"2.16.840.1.101.3.4.3.10", SignatureAlgorithm_ECDSA_SHA3_256},
	{"2.16.840.1.101.3.4.3.11", SignatureAlgorithm_ECDSA_SHA3_384},
	{"2.16.840.1.101.3.4.3.12", SignatureAlgorithm_ECDSA_SHA3_512},
	{"0.4.0.127.0.7.1.1.4.1.1", SignatureAlgorithm_PLAIN_ECDSA_SHA1},
	{"0.4.0.127.0.7.1.1.4.1.2", SignatureAlgorithm_PLAIN_ECDSA_SHA224},
	{"0.4.0.127.0.7.1.1.4.1.3", SignatureAlgorithm_PLAIN_ECDSA_SHA256},
	{"0.4.0.127.0.7.1.1.4.1.4", SignatureAlgorithm_PLAIN_ECDSA_SHA384},
	{"0.4.0.127.0.7.1.1.4.1.5", SignatureAlgorithm_PLAIN_ECDSA_SHA512},
	{"0.4.0.127.0.7.1.1.4.1.6", SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160},
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

var signatureAlgorithmOIDForward = signatureAlgorithmForward(signatureAlgorithmOIDPairs)
var signatureAlgorithmOIDReverse = signatureAlgorithmReverse(signatureAlgorithmOIDPairs, signatureAlgorithmOIDForward, javaStringHashCode, false)

// signatureAlgorithmJCEPairs lists the JAVA JCE signature algorithm names in upstream
// declaration order.
var signatureAlgorithmJCEPairs = []signatureAlgorithmPair[string]{
	{"NONEwithRSA", SignatureAlgorithm_RSA_RAW},
	{"SHA1withRSA", SignatureAlgorithm_RSA_SHA1},
	{"SHA224withRSA", SignatureAlgorithm_RSA_SHA224},
	{"SHA256withRSA", SignatureAlgorithm_RSA_SHA256},
	{"SHA384withRSA", SignatureAlgorithm_RSA_SHA384},
	{"SHA512withRSA", SignatureAlgorithm_RSA_SHA512},
	{"SHA3-224withRSA", SignatureAlgorithm_RSA_SHA3_224},
	{"SHA3-256withRSA", SignatureAlgorithm_RSA_SHA3_256},
	{"SHA3-384withRSA", SignatureAlgorithm_RSA_SHA3_384},
	{"SHA3-512withRSA", SignatureAlgorithm_RSA_SHA3_512},
	{"NONEwithRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_RAW_MGF1},
	{"SHA1withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA1_MGF1},
	{"SHA224withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA224_MGF1},
	{"SHA256withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1},
	{"SHA384withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1},
	{"SHA512withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1},
	{"SHA3-224withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_224_MGF1},
	{"SHA3-256withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_256_MGF1},
	{"SHA3-384withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_384_MGF1},
	{"SHA3-512withRSAandMGF1", SignatureAlgorithm_RSA_SSA_PSS_SHA3_512_MGF1},
	{"RIPEMD160withRSA", SignatureAlgorithm_RSA_RIPEMD160},
	{"MD5withRSA", SignatureAlgorithm_RSA_MD5},
	{"MD2withRSA", SignatureAlgorithm_RSA_MD2},
	{"NONEwithECDSA", SignatureAlgorithm_ECDSA_RAW},
	{"SHA1withECDSA", SignatureAlgorithm_ECDSA_SHA1},
	{"SHA224withECDSA", SignatureAlgorithm_ECDSA_SHA224},
	{"SHA256withECDSA", SignatureAlgorithm_ECDSA_SHA256},
	{"SHA384withECDSA", SignatureAlgorithm_ECDSA_SHA384},
	{"SHA512withECDSA", SignatureAlgorithm_ECDSA_SHA512},
	{"RIPEMD160withECDSA", SignatureAlgorithm_ECDSA_RIPEMD160},
	{"SHA3-224withECDSA", SignatureAlgorithm_ECDSA_SHA3_224},
	{"SHA3-256withECDSA", SignatureAlgorithm_ECDSA_SHA3_256},
	{"SHA3-384withECDSA", SignatureAlgorithm_ECDSA_SHA3_384},
	{"SHA3-512withECDSA", SignatureAlgorithm_ECDSA_SHA3_512},
	{"SHA1withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA1},
	{"SHA224withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA224},
	{"SHA256withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA256},
	{"SHA384withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA384},
	{"SHA512withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA512},
	{"RIPEMD160withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160},
	{"SHA3-224withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA3_224},
	{"SHA3-256withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA3_256},
	{"SHA3-384withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA3_384},
	{"SHA3-512withPLAIN-ECDSA", SignatureAlgorithm_PLAIN_ECDSA_SHA3_512},
	{"Ed25519", SignatureAlgorithm_ED25519},
	{"Ed448", SignatureAlgorithm_ED448},
	{"NONEwithDSA", SignatureAlgorithm_DSA_RAW},
	{"SHA1withDSA", SignatureAlgorithm_DSA_SHA1},
	{"SHA224withDSA", SignatureAlgorithm_DSA_SHA224},
	{"SHA256withDSA", SignatureAlgorithm_DSA_SHA256},
	{"SHA384withDSA", SignatureAlgorithm_DSA_SHA384},
	{"SHA512withDSA", SignatureAlgorithm_DSA_SHA512},
	{"SHA3-224withDSA", SignatureAlgorithm_DSA_SHA3_224},
	{"SHA3-256withDSA", SignatureAlgorithm_DSA_SHA3_256},
	{"SHA3-384withDSA", SignatureAlgorithm_DSA_SHA3_384},
	{"SHA3-512withDSA", SignatureAlgorithm_DSA_SHA3_512},
	{"SHA1withHMAC", SignatureAlgorithm_HMAC_SHA1},
	{"SHA224withHMAC", SignatureAlgorithm_HMAC_SHA224},
	{"SHA256withHMAC", SignatureAlgorithm_HMAC_SHA256},
	{"SHA384withHMAC", SignatureAlgorithm_HMAC_SHA384},
	{"SHA512withHMAC", SignatureAlgorithm_HMAC_SHA512},
	{"SHA3-224withHMAC", SignatureAlgorithm_HMAC_SHA3_224},
	{"SHA3-256withHMAC", SignatureAlgorithm_HMAC_SHA3_256},
	{"SHA3-384withHMAC", SignatureAlgorithm_HMAC_SHA3_384},
	{"SHA3-512withHMAC", SignatureAlgorithm_HMAC_SHA3_512},
	{"RIPEMD160withHMAC", SignatureAlgorithm_HMAC_RIPEMD160},
}

var signatureAlgorithmJCEForward = signatureAlgorithmForward(signatureAlgorithmJCEPairs)
var signatureAlgorithmJCEReverse = signatureAlgorithmReverse(signatureAlgorithmJCEPairs, signatureAlgorithmJCEForward, javaStringHashCode, false)

// signatureAlgorithmJWAPairs lists the JWA algorithm identifiers (RFC 7518 section 3.1)
// in upstream declaration order.
var signatureAlgorithmJWAPairs = []signatureAlgorithmPair[string]{
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

var signatureAlgorithmJWAForward = signatureAlgorithmForward(signatureAlgorithmJWAPairs)
var signatureAlgorithmJWAReverse = func() map[SignatureAlgorithm]string {
	r := signatureAlgorithmReverse(signatureAlgorithmJWAPairs, signatureAlgorithmJWAForward, javaStringHashCode, true)
	// Explicit upstream override: ED448 also reports "EdDSA" as its JWA id, though
	// "EdDSA" only resolves back to ED25519 on the forward lookup.
	r[SignatureAlgorithm_ED448] = "EdDSA"
	return r
}()

// signatureAlgorithmCOSEPairs lists the COSE algorithm keys (https://www.iana.org/assignments/cose/cose.xml)
// in upstream declaration order.
var signatureAlgorithmCOSEPairs = []signatureAlgorithmPair[int64]{
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

var signatureAlgorithmCOSEForward = signatureAlgorithmForward(signatureAlgorithmCOSEPairs)
var signatureAlgorithmCOSEReverse = func() map[SignatureAlgorithm]int64 {
	r := signatureAlgorithmReverse(signatureAlgorithmCOSEPairs, signatureAlgorithmCOSEForward, javaLongHashCode, true)
	// Explicit upstream override: ED448 also reports -8 as its COSE key, though -8
	// only resolves back to ED25519 on the forward lookup.
	r[SignatureAlgorithm_ED448] = -8
	return r
}()

// SignatureAlgorithmForXML returns the corresponding SignatureAlgorithm for the given XML URI.
func SignatureAlgorithmForXML(xmlName string) (SignatureAlgorithm, error) {
	if a, ok := signatureAlgorithmXMLForward[xmlName]; ok {
		return a, nil
	}
	return "", fmt.Errorf("unsupported algorithm: %s", xmlName)
}

// SignatureAlgorithmForXMLDefault returns the SignatureAlgorithm for the given XML URI,
// or defaultValue if the algorithm is unknown.
func SignatureAlgorithmForXMLDefault(xmlName string, defaultValue SignatureAlgorithm) SignatureAlgorithm {
	if a, ok := signatureAlgorithmXMLForward[xmlName]; ok {
		return a
	}
	return defaultValue
}

// SignatureAlgorithmForOID returns the corresponding SignatureAlgorithm for the given OID.
func SignatureAlgorithmForOID(oid string) (SignatureAlgorithm, error) {
	return SignatureAlgorithmForOIDAndParams(oid, nil)
}

// SignatureAlgorithmForOIDAndParams returns the corresponding SignatureAlgorithm for the
// given OID and signature algorithm parameters.
//
// The single RSASSA-PSS OID (1.2.840.113549.1.1.10) names the padding, not the digest, so
// for that algorithm upstream reads the digest out of the RSASSA-PSS-params structure -
// AlgorithmParameters.getInstance("PSS").init(sigAlgParams), then
// PSSParameterSpec#getDigestAlgorithm() and DigestAlgorithm#forJavaName. This port decodes
// the same field directly (RFC 4055 hashAlgorithm [0], DEFAULT sha1) and resolves it by OID,
// which is equivalent for every digest the JDK names and additionally covers the SHA-3
// OIDs. Java's IllegalArgumentException("Unable to initialize PSS") becomes an error.
//
// Passing nil parameters keeps the OID's nominal algorithm, exactly as upstream does; note
// that the nominal algorithm of the PSS OID is RSA_SSA_PSS_SHA1_MGF1.
//
// As upstream, a combination with no matching constant (say RSASSA-PSS over WHIRLPOOL)
// yields the empty SignatureAlgorithm and no error, mirroring the null Java returns.
func SignatureAlgorithmForOIDAndParams(oid string, sigAlgParams []byte) (SignatureAlgorithm, error) {
	algorithm, ok := signatureAlgorithmOIDForward[oid]
	if !ok {
		return "", fmt.Errorf("unsupported algorithm: %s", oid)
	}
	if EncryptionAlgorithm_RSASSA_PSS == algorithm.EncryptionAlgorithm() && sigAlgParams != nil {
		digestAlgorithm, err := signatureAlgorithmPSSDigestAlgorithm(sigAlgParams)
		if err != nil {
			return "", fmt.Errorf("Unable to initialize PSS: %w", err)
		}
		algorithm = signatureAlgorithmGetAlgorithm(algorithm.EncryptionAlgorithm(), digestAlgorithm)
	}
	return algorithm, nil
}

// signatureAlgorithmPSSDigestAlgorithm decodes the hashAlgorithm of an RFC 4055
// RSASSA-PSS-params structure:
//
//	RSASSA-PSS-params ::= SEQUENCE {
//	    hashAlgorithm     [0] HashAlgorithm    DEFAULT sha1,
//	    maskGenAlgorithm  [1] MaskGenAlgorithm DEFAULT mgf1SHA1,
//	    saltLength        [2] INTEGER          DEFAULT 20,
//	    trailerField      [3] TrailerField     DEFAULT trailerFieldBC }
//
// Only the hashAlgorithm is read, but maskGenAlgorithm and trailerField are validated the way
// the JDK's sun.security.rsa.PSSParameters#engineInit(byte[]) does - it rejects a
// maskGenAlgorithm that is not MGF1 ("Only MGF1 mgf is supported") and a trailerField other
// than 1 ("Unsupported trailerField value") - so that parameters upstream refuses to decode
// are refused here too. The JDK additionally restricts the MGF1 digest to a fixed name list;
// that whitelist is provider-specific and is NOT reproduced.
func signatureAlgorithmPSSDigestAlgorithm(sigAlgParams []byte) (DigestAlgorithm, error) {
	input := cryptobyte.String(sigAlgParams)
	var params cryptobyte.String
	if !input.ReadASN1(&params, cryptobyte_asn1.SEQUENCE) {
		return "", errors.New("the parameters are not an RSASSA-PSS-params SEQUENCE")
	}
	if !input.Empty() {
		return "", errors.New("extra data found after the RSASSA-PSS-params")
	}

	var hashAlgorithm cryptobyte.String
	var hashAlgorithmPresent bool
	if !params.ReadOptionalASN1(&hashAlgorithm, &hashAlgorithmPresent,
		cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
		return "", errors.New("malformed RSASSA-PSS-params hashAlgorithm")
	}

	var maskGenAlgorithm cryptobyte.String
	var maskGenAlgorithmPresent bool
	if !params.ReadOptionalASN1(&maskGenAlgorithm, &maskGenAlgorithmPresent,
		cryptobyte_asn1.Tag(1).Constructed().ContextSpecific()) {
		return "", errors.New("malformed RSASSA-PSS-params maskGenAlgorithm")
	}
	if maskGenAlgorithmPresent {
		var maskGenAlgorithmIdentifier cryptobyte.String
		if !maskGenAlgorithm.ReadASN1(&maskGenAlgorithmIdentifier, cryptobyte_asn1.SEQUENCE) {
			return "", errors.New("the RSASSA-PSS-params maskGenAlgorithm is not an AlgorithmIdentifier")
		}
		var maskGenOID encoding_asn1.ObjectIdentifier
		if !maskGenAlgorithmIdentifier.ReadASN1ObjectIdentifier(&maskGenOID) {
			return "", errors.New("the RSASSA-PSS-params maskGenAlgorithm has no OBJECT IDENTIFIER")
		}
		// id-mgf1 OBJECT IDENTIFIER ::= { pkcs-1 8 }
		if maskGenOID.String() != "1.2.840.113549.1.1.8" {
			return "", errors.New("Only MGF1 mgf is supported")
		}
	}

	// saltLength [2] INTEGER DEFAULT 20 is skipped over, its value being unused.
	var saltLength cryptobyte.String
	var saltLengthPresent bool
	if !params.ReadOptionalASN1(&saltLength, &saltLengthPresent,
		cryptobyte_asn1.Tag(2).Constructed().ContextSpecific()) {
		return "", errors.New("malformed RSASSA-PSS-params saltLength")
	}

	var trailerField cryptobyte.String
	var trailerFieldPresent bool
	if !params.ReadOptionalASN1(&trailerField, &trailerFieldPresent,
		cryptobyte_asn1.Tag(3).Constructed().ContextSpecific()) {
		return "", errors.New("malformed RSASSA-PSS-params trailerField")
	}
	if trailerFieldPresent {
		var trailerFieldValue int
		if !trailerField.ReadASN1Integer(&trailerFieldValue) {
			return "", errors.New("the RSASSA-PSS-params trailerField is not an INTEGER")
		}
		// TrailerField ::= INTEGER { trailerFieldBC(1) }
		if trailerFieldValue != 1 {
			return "", fmt.Errorf("Unsupported trailerField value %d", trailerFieldValue)
		}
	}

	if !hashAlgorithmPresent {
		// hashAlgorithm DEFAULT sha1
		return DigestAlgorithm_SHA1, nil
	}
	var algorithmIdentifier cryptobyte.String
	if !hashAlgorithm.ReadASN1(&algorithmIdentifier, cryptobyte_asn1.SEQUENCE) {
		return "", errors.New("the RSASSA-PSS-params hashAlgorithm is not an AlgorithmIdentifier")
	}
	var oid encoding_asn1.ObjectIdentifier
	if !algorithmIdentifier.ReadASN1ObjectIdentifier(&oid) {
		return "", errors.New("the RSASSA-PSS-params hashAlgorithm has no OBJECT IDENTIFIER")
	}
	return DigestAlgorithmForOID(oid.String())
}

// SignatureAlgorithmForJWA returns the corresponding SignatureAlgorithm for the given JWA name.
func SignatureAlgorithmForJWA(jsonWebAlgorithm string) (SignatureAlgorithm, error) {
	if a, ok := signatureAlgorithmJWAForward[jsonWebAlgorithm]; ok {
		return a, nil
	}
	return "", fmt.Errorf("unsupported algorithm: %s", jsonWebAlgorithm)
}

// SignatureAlgorithmForJWADefault returns the SignatureAlgorithm for the given JWA name,
// or defaultValue if the algorithm is unknown.
func SignatureAlgorithmForJWADefault(jsonWebAlgorithm string, defaultValue SignatureAlgorithm) SignatureAlgorithm {
	if a, ok := signatureAlgorithmJWAForward[jsonWebAlgorithm]; ok {
		return a
	}
	return defaultValue
}

// SignatureAlgorithmForCOSEDefault returns the SignatureAlgorithm for the given COSE
// algorithm key, or defaultValue if the algorithm is unknown.
func SignatureAlgorithmForCOSEDefault(algorithmKey int64, defaultValue SignatureAlgorithm) SignatureAlgorithm {
	if a, ok := signatureAlgorithmCOSEForward[algorithmKey]; ok {
		return a
	}
	return defaultValue
}

// SignatureAlgorithmForJAVA returns the corresponding SignatureAlgorithm for the given
// Java JCE signature algorithm name.
func SignatureAlgorithmForJAVA(javaName string) (SignatureAlgorithm, error) {
	if a, ok := signatureAlgorithmJCEForward[javaName]; ok {
		return a, nil
	}
	return "", fmt.Errorf("unsupported algorithm: %s", javaName)
}

// URI returns the XML ID (URI) of the signature algorithm. Implements OidAndUriBasedEnum / UriBasedEnum.
func (s SignatureAlgorithm) URI() string {
	return signatureAlgorithmXMLReverse[s]
}

// OID returns the OID of the signature algorithm. Implements OidAndUriBasedEnum / OidBasedEnum.
func (s SignatureAlgorithm) OID() string {
	return signatureAlgorithmOIDReverse[s]
}

// URIBasedOnOID returns the URI of the signature algorithm generated from its OID:
//
//	Ex.: OID = 1.2.4.5.6.8 becomes URI = urn:oid:1.2.4.5.6.8
//
// Note: see RFC 3061 "A URN Namespace of Object Identifiers".
// For an algorithm with no registered OID upstream concatenates a Java null, yielding the
// literal "urn:oid:null" (e.g. RSA_RAW, the RSA_SSA_PSS_* family, ECDSA_RAW,
// ECDSA_RIPEMD160, DSA_RAW). That string is reproduced verbatim rather than collapsed to
// "urn:oid:", because this value is embedded in emitted signatures and the port's contract
// is byte-identical output with upstream DSS.
func (s SignatureAlgorithm) URIBasedOnOID() string {
	oid := s.OID()
	if oid == "" {
		oid = "null"
	}
	return signatureAlgorithmOIDNamespacePrefix + oid
}

// JCEID returns the algorithm identifier corresponding to the JAVA JCE class name.
func (s SignatureAlgorithm) JCEID() string {
	return signatureAlgorithmJCEReverse[s]
}

// JWAID returns the algorithm identifier corresponding to JWA accepted algorithms (RFC 7518).
func (s SignatureAlgorithm) JWAID() string {
	return signatureAlgorithmJWAReverse[s]
}

// COSEID returns the algorithm identifier corresponding to COSE accepted algorithms
// (RFC 9053), and false if the signature algorithm has no COSE identifier.
func (s SignatureAlgorithm) COSEID() (int64, bool) {
	v, ok := signatureAlgorithmCOSEReverse[s]
	return v, ok
}
