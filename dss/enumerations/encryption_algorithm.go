// Ported from dss-enumerations/.../EncryptionAlgorithm.java (DSS 6.5.RC1).
//
// forKey(java.security.Key) is omitted: java.security.Key has no Go stdlib
// equivalent tied to this package. Callers should determine the JCE-style
// algorithm name themselves and call EncryptionAlgorithmForName.
package enumerations

import "fmt"

// EncryptionAlgorithm represents the supported signature encryption
// algorithms. Implements OidBasedEnum.
type EncryptionAlgorithm string

const (
	// EncryptionAlgorithm_RSA is RSA.
	EncryptionAlgorithm_RSA EncryptionAlgorithm = "RSA"
	// EncryptionAlgorithm_RSASSA_PSS is RSASSA-PSS.
	EncryptionAlgorithm_RSASSA_PSS EncryptionAlgorithm = "RSASSA_PSS"
	// EncryptionAlgorithm_DSA is DSA.
	EncryptionAlgorithm_DSA EncryptionAlgorithm = "DSA"
	// EncryptionAlgorithm_ECDSA is ECDSA.
	EncryptionAlgorithm_ECDSA EncryptionAlgorithm = "ECDSA"
	// EncryptionAlgorithm_PLAIN_ECDSA is PLAIN-ECDSA.
	EncryptionAlgorithm_PLAIN_ECDSA EncryptionAlgorithm = "PLAIN_ECDSA"
	// EncryptionAlgorithm_X25519 is X25519.
	EncryptionAlgorithm_X25519 EncryptionAlgorithm = "X25519"
	// EncryptionAlgorithm_X448 is X448.
	EncryptionAlgorithm_X448 EncryptionAlgorithm = "X448"
	// EncryptionAlgorithm_EDDSA is EdDSA.
	EncryptionAlgorithm_EDDSA EncryptionAlgorithm = "EDDSA"
	// EncryptionAlgorithm_HMAC is HMAC.
	EncryptionAlgorithm_HMAC EncryptionAlgorithm = "HMAC"
)

// encryptionAlgorithmFields holds the (name, oid, padding) tuple for each constant.
type encryptionAlgorithmFields struct {
	name    string
	oid     string
	padding string
}

// encryptionAlgorithmData holds the full field tuple for each constant,
// copied verbatim from the Java enum constructors.
var encryptionAlgorithmData = map[EncryptionAlgorithm]encryptionAlgorithmFields{
	EncryptionAlgorithm_RSA:         {"RSA", "1.2.840.113549.1.1.1", "RSA/ECB/PKCS1Padding"},
	EncryptionAlgorithm_RSASSA_PSS:  {"RSASSA-PSS", "1.2.840.113549.1.1.10", "RSA/ECB/OAEPPadding"},
	EncryptionAlgorithm_DSA:         {"DSA", "1.2.840.10040.4.1", "DSA"},
	EncryptionAlgorithm_ECDSA:       {"ECDSA", "1.2.840.10045.2.1", "ECDSA"},
	EncryptionAlgorithm_PLAIN_ECDSA: {"PLAIN-ECDSA", "0.4.0.127.0.7.1.1.4.1", "PLAIN-ECDSA"},
	EncryptionAlgorithm_X25519:      {"X25519", "1.3.101.110", "X25519"},
	EncryptionAlgorithm_X448:        {"X448", "1.3.101.111", "X448"},
	EncryptionAlgorithm_EDDSA:       {"EdDSA", "", "EdDSA"},
	EncryptionAlgorithm_HMAC:        {"HMAC", "", ""},
}

// encryptionAlgorithmOIDLookup maps OIDs to their EncryptionAlgorithm,
// mirroring Java's Registry.OID_ALGORITHMS.
var encryptionAlgorithmOIDLookup = func() map[string]EncryptionAlgorithm {
	m := make(map[string]EncryptionAlgorithm)
	for _, v := range EncryptionAlgorithmValues() {
		m[encryptionAlgorithmData[v].oid] = v
	}
	return m
}()

// EncryptionAlgorithmValues returns all constants in declaration order.
func EncryptionAlgorithmValues() []EncryptionAlgorithm {
	return []EncryptionAlgorithm{
		EncryptionAlgorithm_RSA,
		EncryptionAlgorithm_RSASSA_PSS,
		EncryptionAlgorithm_DSA,
		EncryptionAlgorithm_ECDSA,
		EncryptionAlgorithm_PLAIN_ECDSA,
		EncryptionAlgorithm_X25519,
		EncryptionAlgorithm_X448,
		EncryptionAlgorithm_EDDSA,
		EncryptionAlgorithm_HMAC,
	}
}

// EncryptionAlgorithmForOID returns the encryption algorithm associated to
// the given OID.
func EncryptionAlgorithmForOID(oid string) (EncryptionAlgorithm, error) {
	if algorithm, ok := encryptionAlgorithmOIDLookup[oid]; ok {
		return algorithm, nil
	}
	return "", fmt.Errorf("unsupported algorithm: %s", oid)
}

// EncryptionAlgorithmForName returns the encryption algorithm associated to
// the given JCE name.
func EncryptionAlgorithmForName(name string) (EncryptionAlgorithm, error) {
	// To be checked if ECC exists also.
	if name == "EC" || name == "ECC" {
		return EncryptionAlgorithm_ECDSA, nil
	}

	// Since JDK 15.
	if name == "Ed25519" || name == "Ed448" {
		return EncryptionAlgorithm_EDDSA, nil
	}

	for _, v := range EncryptionAlgorithmValues() {
		if v.Name() == name || string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("unsupported algorithm: %s", name)
}

// EncryptionAlgorithmForNameDefault returns the encryption algorithm
// associated to the given JCE name, or defaultValue if the name doesn't
// match any algorithm.
func EncryptionAlgorithmForNameDefault(name string, defaultValue EncryptionAlgorithm) EncryptionAlgorithm {
	v, err := EncryptionAlgorithmForName(name)
	if err != nil {
		return defaultValue
	}
	return v
}

// Name gets the algorithm name.
func (e EncryptionAlgorithm) Name() string {
	return encryptionAlgorithmData[e].name
}

// OID gets the ASN1 algorithm OID. Implements OidBasedEnum.
func (e EncryptionAlgorithm) OID() string {
	return encryptionAlgorithmData[e].oid
}

// Padding gets the algorithm padding.
func (e EncryptionAlgorithm) Padding() string {
	return encryptionAlgorithmData[e].padding
}

// IsEquivalent verifies if the provided encryptionAlgorithm is equivalent to
// the current one. Equivalent means the same token key can be used for
// signature creation with both algorithms.
func (e EncryptionAlgorithm) IsEquivalent(encryptionAlgorithm EncryptionAlgorithm) bool {
	if encryptionAlgorithm == "" {
		return false
	}
	if e == encryptionAlgorithm {
		return true
	}
	if e.isRSAFamily() && encryptionAlgorithm.isRSAFamily() {
		return true
	}
	if e.isEcDSAFamily() && encryptionAlgorithm.isEcDSAFamily() {
		return true
	}
	if e.isEdDSAFamily() && encryptionAlgorithm.isEdDSAFamily() {
		return true
	}
	return false
}

func (e EncryptionAlgorithm) isRSAFamily() bool {
	return e == EncryptionAlgorithm_RSA || e == EncryptionAlgorithm_RSASSA_PSS
}

func (e EncryptionAlgorithm) isEcDSAFamily() bool {
	return e == EncryptionAlgorithm_ECDSA || e == EncryptionAlgorithm_PLAIN_ECDSA
}

func (e EncryptionAlgorithm) isEdDSAFamily() bool {
	return e == EncryptionAlgorithm_X25519 || e == EncryptionAlgorithm_X448 || e == EncryptionAlgorithm_EDDSA
}

// compile-time interface assertion.
var _ OidBasedEnum = EncryptionAlgorithm("")
