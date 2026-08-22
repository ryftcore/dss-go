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
	// DigestAlgorithmSHA1 is SHA-1.
	DigestAlgorithmSHA1 DigestAlgorithm = "SHA1"
	// DigestAlgorithmSHA224 is SHA-224.
	DigestAlgorithmSHA224 DigestAlgorithm = "SHA224"
	// DigestAlgorithmSHA256 is SHA-256.
	DigestAlgorithmSHA256 DigestAlgorithm = "SHA256"
	// DigestAlgorithmSHA384 is SHA-384.
	DigestAlgorithmSHA384 DigestAlgorithm = "SHA384"
	// DigestAlgorithmSHA512 is SHA-512.
	DigestAlgorithmSHA512 DigestAlgorithm = "SHA512"
	// DigestAlgorithmSHA3224 is SHA3-224. See https://tools.ietf.org/html/rfc6931
	DigestAlgorithmSHA3224 DigestAlgorithm = "SHA3_224"
	// DigestAlgorithmSHA3256 is SHA3-256.
	DigestAlgorithmSHA3256 DigestAlgorithm = "SHA3_256"
	// DigestAlgorithmSHA3384 is SHA3-384.
	DigestAlgorithmSHA3384 DigestAlgorithm = "SHA3_384"
	// DigestAlgorithmSHA3512 is SHA3-512.
	DigestAlgorithmSHA3512 DigestAlgorithm = "SHA3_512"
	// DigestAlgorithmSHAKE128 is SHAKE-128.
	DigestAlgorithmSHAKE128 DigestAlgorithm = "SHAKE128"
	// DigestAlgorithmSHAKE256 is SHAKE-256.
	DigestAlgorithmSHAKE256 DigestAlgorithm = "SHAKE256"
	// DigestAlgorithmSHAKE256512 is SHAKE-256 with output 512 bits.
	DigestAlgorithmSHAKE256512 DigestAlgorithm = "SHAKE256_512"
	// DigestAlgorithmRIPEMD160 is RIPEMD160.
	DigestAlgorithmRIPEMD160 DigestAlgorithm = "RIPEMD160"
	// DigestAlgorithmMD2 is MD2.
	DigestAlgorithmMD2 DigestAlgorithm = "MD2"
	// DigestAlgorithmMD5 is MD5.
	DigestAlgorithmMD5 DigestAlgorithm = "MD5"
	// DigestAlgorithmWHIRLPOOL is WHIRLPOOL.
	DigestAlgorithmWHIRLPOOL DigestAlgorithm = "WHIRLPOOL"
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
	DigestAlgorithmSHA1: {
		name: "SHA1", javaName: "SHA-1", oid: "1.3.14.3.2.26", xmlID: "http://www.w3.org/2000/09/xmldsig#sha1",
		jadesID: "", httpHeaderID: "SHA", sdJwtID: "", srIntegrityID: "", coseID: int64p(-14), msoID: "", saltLength: 20,
	},
	DigestAlgorithmSHA224: {
		name: "SHA224", javaName: "SHA-224", oid: "2.16.840.1.101.3.4.2.4", xmlID: "http://www.w3.org/2001/04/xmldsig-more#sha224",
		jadesID: "S224", saltLength: 28,
	},
	DigestAlgorithmSHA256: {
		name: "SHA256", javaName: "SHA-256", oid: "2.16.840.1.101.3.4.2.1", xmlID: "http://www.w3.org/2001/04/xmlenc#sha256",
		jadesID: "S256", httpHeaderID: "SHA-256", sdJwtID: "sha-256", srIntegrityID: "sha256", coseID: int64p(-16), msoID: "SHA-256", saltLength: 32,
	},
	DigestAlgorithmSHA384: {
		name: "SHA384", javaName: "SHA-384", oid: "2.16.840.1.101.3.4.2.2", xmlID: "http://www.w3.org/2001/04/xmldsig-more#sha384",
		jadesID: "S384", httpHeaderID: "", sdJwtID: "sha-384", srIntegrityID: "sha384", coseID: int64p(-43), msoID: "SHA-384", saltLength: 48,
	},
	DigestAlgorithmSHA512: {
		name: "SHA512", javaName: "SHA-512", oid: "2.16.840.1.101.3.4.2.3", xmlID: "http://www.w3.org/2001/04/xmlenc#sha512",
		jadesID: "S512", httpHeaderID: "SHA-512", sdJwtID: "sha-512", srIntegrityID: "sha512", coseID: int64p(-44), msoID: "SHA-512", saltLength: 64,
	},
	DigestAlgorithmSHA3224: {
		name: "SHA3-224", javaName: "SHA3-224", oid: "2.16.840.1.101.3.4.2.7", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-224",
		jadesID: "", httpHeaderID: "", sdJwtID: "sha3-224", saltLength: 28,
	},
	DigestAlgorithmSHA3256: {
		name: "SHA3-256", javaName: "SHA3-256", oid: "2.16.840.1.101.3.4.2.8", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-256",
		jadesID: "S3-256", httpHeaderID: "", sdJwtID: "sha3-256", saltLength: 32,
	},
	DigestAlgorithmSHA3384: {
		name: "SHA3-384", javaName: "SHA3-384", oid: "2.16.840.1.101.3.4.2.9", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-384",
		jadesID: "S3-384", httpHeaderID: "", sdJwtID: "sha3-384", saltLength: 48,
	},
	DigestAlgorithmSHA3512: {
		name: "SHA3-512", javaName: "SHA3-512", oid: "2.16.840.1.101.3.4.2.10", xmlID: "http://www.w3.org/2007/05/xmldsig-more#sha3-512",
		jadesID: "S3-512", httpHeaderID: "", sdJwtID: "sha3-512", saltLength: 64,
	},
	DigestAlgorithmSHAKE128: {
		name: "SHAKE-128", javaName: "SHAKE-128", oid: "2.16.840.1.101.3.4.2.11", xmlID: "",
		coseID: int64p(-18),
	},
	DigestAlgorithmSHAKE256: {
		name: "SHAKE-256", javaName: "SHAKE-256", oid: "2.16.840.1.101.3.4.2.12", xmlID: "",
	},
	DigestAlgorithmSHAKE256512: {
		name: "SHAKE256-512", javaName: "SHAKE256-512", oid: "2.16.840.1.101.3.4.2.18", xmlID: "",
		coseID: int64p(-45),
	},
	DigestAlgorithmRIPEMD160: {
		name: "RIPEMD160", javaName: "RIPEMD160", oid: "1.3.36.3.2.1", xmlID: "http://www.w3.org/2001/04/xmlenc#ripemd160",
	},
	DigestAlgorithmMD2: {
		name: "MD2", javaName: "MD2", oid: "1.2.840.113549.2.2", xmlID: "http://www.w3.org/2001/04/xmldsig-more#md2",
	},
	DigestAlgorithmMD5: {
		name: "MD5", javaName: "MD5", oid: "1.2.840.113549.2.5", xmlID: "http://www.w3.org/2001/04/xmldsig-more#md5",
		jadesID: "", httpHeaderID: "MD5",
	},
	DigestAlgorithmWHIRLPOOL: {
		name: "WHIRLPOOL", javaName: "WHIRLPOOL", oid: "1.0.10118.3.0.55", xmlID: "http://www.w3.org/2007/05/xmldsig-more#whirlpool",
	},
}

// DigestAlgorithmValues returns all constants in declaration order.
func DigestAlgorithmValues() []DigestAlgorithm {
	return []DigestAlgorithm{
		DigestAlgorithmSHA1,
		DigestAlgorithmSHA224,
		DigestAlgorithmSHA256,
		DigestAlgorithmSHA384,
		DigestAlgorithmSHA512,
		DigestAlgorithmSHA3224,
		DigestAlgorithmSHA3256,
		DigestAlgorithmSHA3384,
		DigestAlgorithmSHA3512,
		DigestAlgorithmSHAKE128,
		DigestAlgorithmSHAKE256,
		DigestAlgorithmSHAKE256512,
		DigestAlgorithmRIPEMD160,
		DigestAlgorithmMD2,
		DigestAlgorithmMD5,
		DigestAlgorithmWHIRLPOOL,
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
