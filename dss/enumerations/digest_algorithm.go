// Ported from dss-enumerations/.../DigestAlgorithm.java (DSS 6.5.RC1).
//
// getMessageDigest()/getMessageDigest(Provider) are omitted: they return a
// java.security.MessageDigest, which has no Go stdlib equivalent tied to
// this package; digest computation belongs in the module that ports
// java.security usage, keyed off JavaName().
package enumerations

import "fmt"

// DigestAlgorithm represents the supported digest (hash) algorithms.
// Implements OidAndUriBasedEnum.
type DigestAlgorithm string

const (
	// DigestAlgorithm_SHA1 is SHA-1.
	DigestAlgorithm_SHA1 DigestAlgorithm = "SHA1"
	// DigestAlgorithm_SHA224 is SHA-224.
	DigestAlgorithm_SHA224 DigestAlgorithm = "SHA224"
	// DigestAlgorithm_SHA256 is SHA-256.
	DigestAlgorithm_SHA256 DigestAlgorithm = "SHA256"
	// DigestAlgorithm_SHA384 is SHA-384.
	DigestAlgorithm_SHA384 DigestAlgorithm = "SHA384"
	// DigestAlgorithm_SHA512 is SHA-512.
	DigestAlgorithm_SHA512 DigestAlgorithm = "SHA512"
	// DigestAlgorithm_SHA3_224 is SHA3-224. See https://tools.ietf.org/html/rfc6931
	DigestAlgorithm_SHA3_224 DigestAlgorithm = "SHA3_224"
	// DigestAlgorithm_SHA3_256 is SHA3-256.
	DigestAlgorithm_SHA3_256 DigestAlgorithm = "SHA3_256"
	// DigestAlgorithm_SHA3_384 is SHA3-384.
	DigestAlgorithm_SHA3_384 DigestAlgorithm = "SHA3_384"
	// DigestAlgorithm_SHA3_512 is SHA3-512.
	DigestAlgorithm_SHA3_512 DigestAlgorithm = "SHA3_512"
	// DigestAlgorithm_SHAKE128 is SHAKE-128.
	DigestAlgorithm_SHAKE128 DigestAlgorithm = "SHAKE128"
	// DigestAlgorithm_SHAKE256 is SHAKE-256.
	DigestAlgorithm_SHAKE256 DigestAlgorithm = "SHAKE256"
	// DigestAlgorithm_SHAKE256_512 is SHAKE-256 with output 512 bits.
	DigestAlgorithm_SHAKE256_512 DigestAlgorithm = "SHAKE256_512"
	// DigestAlgorithm_RIPEMD160 is RIPEMD160.
	DigestAlgorithm_RIPEMD160 DigestAlgorithm = "RIPEMD160"
	// DigestAlgorithm_MD2 is MD2.
	DigestAlgorithm_MD2 DigestAlgorithm = "MD2"
	// DigestAlgorithm_MD5 is MD5.
	DigestAlgorithm_MD5 DigestAlgorithm = "MD5"
	// DigestAlgorithm_WHIRLPOOL is WHIRLPOOL.
	DigestAlgorithm_WHIRLPOOL DigestAlgorithm = "WHIRLPOOL"
)

// digestAlgorithmFields holds all per-constant attributes of DigestAlgorithm,
// mirroring the Java enum's instance fields.
type digestAlgorithmFields struct {
	name          string
	javaName      string
	oid           string
	xmlID         string
	jadesID       string
	httpHeaderID  string
	sdJwtID       string
	srIntegrityID string
	coseID        *int64
	msoID         string
	saltLength    int
}

func int64p(v int64) *int64 { return &v }

// digestAlgorithmData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var digestAlgorithmData = map[DigestAlgorithm]digestAlgorithmFields{
	DigestAlgorithm_SHA1: {
		name: "SHA1", javaName: "SHA-1", oid: "1.3.14.3.2.26", xmlID: "http://www.w3.org/2000/09/xmldsig#sha1",
		jadesID: "", httpHeaderID: "SHA", sdJwtID: "", srIntegrityID: "", coseID: int64p(-14), msoID: "", saltLength: 20,
	},
	DigestAlgorithm_SHA224: {
		name: "SHA224", javaName: "SHA-224", oid: "2.16.840.1.101.3.4.2.4", xmlID: "http://www.w3.org/2001/04/xmldsig-more#sha224",
		jadesID: "S224", saltLength: 28,
	},
	DigestAlgorithm_SHA256: {
		name: "SHA256", javaName: "SHA-256", oid: "2.16.840.1.101.3.4.2.1", xmlID: "http://www.w3.org/2001/04/xmlenc#sha256",
		jadesID: "S256", httpHeaderID: "SHA-256", sdJwtID: "sha-256", srIntegrityID: "sha256", coseID: int64p(-16), msoID: "SHA-256", saltLength: 32,
	},
	DigestAlgorithm_SHA384: {
		name: "SHA384", javaName: "SHA-384", oid: "2.16.840.1.101.3.4.2.2", xmlID: "http://www.w3.org/2001/04/xmldsig-more#sha384",
		jadesID: "S384", httpHeaderID: "", sdJwtID: "sha-384", srIntegrityID: "sha384", coseID: int64p(-43), msoID: "SHA-384", saltLength: 48,
	},
	DigestAlgorithm_SHA512: {
		name: "SHA512", javaName: "SHA-512", oid: "2.16.840.1.101.3.4.2.3", xmlID: "http://www.w3.org/2001/04/xmlenc#sha512",
		jadesID: "S512", httpHeaderID: "SHA-512", sdJwtID: "sha-512", srIntegrityID: "sha512", coseID: int64p(-44), msoID: "SHA-512", saltLength: 64,
	},
	DigestAlgorithm_SHA3_224: {
		name: "SHA3-224", javaName: "SHA3-224", oid: "2.16.840.1.101.3.4.2.7", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-224",
		jadesID: "", httpHeaderID: "", sdJwtID: "sha3-224", saltLength: 28,
	},
	DigestAlgorithm_SHA3_256: {
		name: "SHA3-256", javaName: "SHA3-256", oid: "2.16.840.1.101.3.4.2.8", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-256",
		jadesID: "S3-256", httpHeaderID: "", sdJwtID: "sha3-256", saltLength: 32,
	},
	DigestAlgorithm_SHA3_384: {
		name: "SHA3-384", javaName: "SHA3-384", oid: "2.16.840.1.101.3.4.2.9", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-384",
		jadesID: "S3-384", httpHeaderID: "", sdJwtID: "sha3-384", saltLength: 48,
	},
	DigestAlgorithm_SHA3_512: {
		name: "SHA3-512", javaName: "SHA3-512", oid: "2.16.840.1.101.3.4.2.10", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-512",
		jadesID: "S3-512", httpHeaderID: "", sdJwtID: "sha3-512", saltLength: 64,
	},
	DigestAlgorithm_SHAKE128: {
		name: "SHAKE-128", javaName: "SHAKE-128", oid: "2.16.840.1.101.3.4.2.11", xmlID: "",
		coseID: int64p(-18),
	},
	DigestAlgorithm_SHAKE256: {
		name: "SHAKE-256", javaName: "SHAKE-256", oid: "2.16.840.1.101.3.4.2.12", xmlID: "",
	},
	DigestAlgorithm_SHAKE256_512: {
		name: "SHAKE256-512", javaName: "SHAKE256-512", oid: "2.16.840.1.101.3.4.2.18", xmlID: "",
		coseID: int64p(-45),
	},
	DigestAlgorithm_RIPEMD160: {
		name: "RIPEMD160", javaName: "RIPEMD160", oid: "1.3.36.3.2.1", xmlID: "http://www.w3.org/2001/04/xmlenc#ripemd160",
	},
	DigestAlgorithm_MD2: {
		name: "MD2", javaName: "MD2", oid: "1.2.840.113549.2.2", xmlID: "http://www.w3.org/2001/04/xmldsig-more#md2",
	},
	DigestAlgorithm_MD5: {
		name: "MD5", javaName: "MD5", oid: "1.2.840.113549.2.5", xmlID: "http://www.w3.org/2001/04/xmldsig-more#md5",
		jadesID: "", httpHeaderID: "MD5",
	},
	DigestAlgorithm_WHIRLPOOL: {
		name: "WHIRLPOOL", javaName: "WHIRLPOOL", oid: "1.0.10118.3.0.55", xmlID: "http://www.w3.org/2007/05/xmldsig-more#whirlpool",
	},
}

// DigestAlgorithmValues returns all constants in declaration order.
func DigestAlgorithmValues() []DigestAlgorithm {
	return []DigestAlgorithm{
		DigestAlgorithm_SHA1,
		DigestAlgorithm_SHA224,
		DigestAlgorithm_SHA256,
		DigestAlgorithm_SHA384,
		DigestAlgorithm_SHA512,
		DigestAlgorithm_SHA3_224,
		DigestAlgorithm_SHA3_256,
		DigestAlgorithm_SHA3_384,
		DigestAlgorithm_SHA3_512,
		DigestAlgorithm_SHAKE128,
		DigestAlgorithm_SHAKE256,
		DigestAlgorithm_SHAKE256_512,
		DigestAlgorithm_RIPEMD160,
		DigestAlgorithm_MD2,
		DigestAlgorithm_MD5,
		DigestAlgorithm_WHIRLPOOL,
	}
}

// Name gets the algorithm name.
func (d DigestAlgorithm) Name() string { return digestAlgorithmData[d].name }

// JavaName gets the JCE algorithm name.
func (d DigestAlgorithm) JavaName() string { return digestAlgorithmData[d].javaName }

// OID gets the algorithm OID. Implements OidBasedEnum.
func (d DigestAlgorithm) OID() string { return digestAlgorithmData[d].oid }

// URI gets the algorithm uri. Implements UriBasedEnum.
func (d DigestAlgorithm) URI() string { return digestAlgorithmData[d].xmlID }

// JAdESID gets the algorithm id used in JAdES Signatures.
//
// TS 119-182 Annex E (normative): Digest algorithms identifiers for JAdES
// signatures.
func (d DigestAlgorithm) JAdESID() string { return digestAlgorithmData[d].jadesID }

// HttpHeaderAlgo gets the algorithm name according to RFC 5843.
func (d DigestAlgorithm) HttpHeaderAlgo() string { return digestAlgorithmData[d].httpHeaderID }

// SDJWTID gets the algorithm id used for claims integrity definition
// within SD-JWT. The allowed values are available at the IANA Named
// Information registry.
func (d DigestAlgorithm) SDJWTID() string { return digestAlgorithmData[d].sdJwtID }

// SubresourceIntegrityID gets the algorithm id used for subresource
// integrity calculation, per W3C Subresource Integrity.
func (d DigestAlgorithm) SubresourceIntegrityID() string {
	return digestAlgorithmData[d].srIntegrityID
}

// CoseID gets the algorithm Id used in COSE Signatures (IANA COSE
// Algorithms registry). Returns nil if the algorithm has no COSE id,
// mirroring Java's null Long.
func (d DigestAlgorithm) CoseID() *int64 { return digestAlgorithmData[d].coseID }

// MSOID gets the algorithm Id used in MobileSecurityObject structure of
// mdoc, per ISO/IEC 18013-5 "9.1.2.5 Message digest function".
func (d DigestAlgorithm) MSOID() string { return digestAlgorithmData[d].msoID }

// SaltLength gets the salt length (PSS).
func (d DigestAlgorithm) SaltLength() int { return digestAlgorithmData[d].saltLength }

func digestAlgorithmUnsupported(value string) error {
	return fmt.Errorf("Unsupported algorithm: %s", value)
}

// DigestAlgorithmForName returns the digest algorithm associated to the
// given name.
func DigestAlgorithmForName(name string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].name == name {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(name)
}

// DigestAlgorithmForNameDefault returns the digest algorithm associated to
// the given name, or defaultValue if not found.
func DigestAlgorithmForNameDefault(name string, defaultValue DigestAlgorithm) DigestAlgorithm {
	v, err := DigestAlgorithmForName(name)
	if err != nil {
		return defaultValue
	}
	return v
}

// DigestAlgorithmIsSupportedAlgorithm returns an indication if the
// algorithm with the given name is supported.
func DigestAlgorithmIsSupportedAlgorithm(name string) bool {
	_, err := DigestAlgorithmForName(name)
	return err == nil
}

// DigestAlgorithmForJavaName returns the digest algorithm associated to the
// given JCE name.
func DigestAlgorithmForJavaName(javaName string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].javaName == javaName {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(javaName)
}

// DigestAlgorithmForOID returns the digest algorithm associated to the
// given OID.
func DigestAlgorithmForOID(oid string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].oid == oid {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(oid)
}

// DigestAlgorithmForXML returns the digest algorithm associated to the
// given XML url.
func DigestAlgorithmForXML(xmlName string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].xmlID == xmlName {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(xmlName)
}

// DigestAlgorithmForJAdES returns the digest algorithm associated with the
// given identifier, according to TS 119 182-1, Annex E (Digest algorithms
// identifiers for JAdES signatures).
func DigestAlgorithmForJAdES(algoID string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].jadesID == algoID {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(algoID)
}

// DigestAlgorithmForHttpHeader returns the digest algorithm associated to
// the given JWS Http Header Hash Algorithm. See RFC 5843.
func DigestAlgorithmForHttpHeader(hashName string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].httpHeaderID == hashName {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(hashName)
}

// DigestAlgorithmForSdJwtId returns the digest algorithm associated to the
// algorithm identifiers used within SD-JWT tokens. See IANA "Named
// Information Hash Algorithm" registry.
func DigestAlgorithmForSdJwtId(sdJwtID string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].sdJwtID == sdJwtID {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(sdJwtID)
}

// DigestAlgorithmForSrIntegrityId returns the digest algorithm associated
// to a subresource integrity definition as defined in W3C Subresource
// Integrity.
func DigestAlgorithmForSrIntegrityId(srIntegrityID string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].srIntegrityID == srIntegrityID {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(srIntegrityID)
}

// DigestAlgorithmForCOSE returns the digest algorithm associated with the
// given identifier, according to IANA CBOR Object Signing and Encryption
// (COSE).
func DigestAlgorithmForCOSE(algoID int64) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		id := digestAlgorithmData[v].coseID
		if id != nil && *id == algoID {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(fmt.Sprintf("%d", algoID))
}

// DigestAlgorithmForMSO returns the digest algorithm associated with the
// given identifier, according to ISO 18013-5 "9.1.2.5 Message digest
// function".
func DigestAlgorithmForMSO(algoID string) (DigestAlgorithm, error) {
	for _, v := range DigestAlgorithmValues() {
		if digestAlgorithmData[v].msoID == algoID {
			return v, nil
		}
	}
	return "", digestAlgorithmUnsupported(algoID)
}
