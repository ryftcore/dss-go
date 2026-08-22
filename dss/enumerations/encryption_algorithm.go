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
	// EncryptionAlgorithmRSA is RSA.
	EncryptionAlgorithmRSA EncryptionAlgorithm = "RSA"
	// EncryptionAlgorithmRSASSAPSS is RSASSA-PSS.
	EncryptionAlgorithmRSASSAPSS EncryptionAlgorithm = "RSASSA_PSS"
	// EncryptionAlgorithmDSA is DSA.
	EncryptionAlgorithmDSA EncryptionAlgorithm = "DSA"
	// EncryptionAlgorithmECDSA is ECDSA.
	EncryptionAlgorithmECDSA EncryptionAlgorithm = "ECDSA"
	// EncryptionAlgorithmPlainECDSA is PLAIN-ECDSA.
	EncryptionAlgorithmPlainECDSA EncryptionAlgorithm = "PLAIN_ECDSA"
	// EncryptionAlgorithmX25519 is X25519.
	EncryptionAlgorithmX25519 EncryptionAlgorithm = "X25519"
	// EncryptionAlgorithmX448 is X448.
	EncryptionAlgorithmX448 EncryptionAlgorithm = "X448"
	// EncryptionAlgorithmEDDSA is EdDSA.
	EncryptionAlgorithmEDDSA EncryptionAlgorithm = "EDDSA"
	// EncryptionAlgorithmHMAC is HMAC.
	EncryptionAlgorithmHMAC EncryptionAlgorithm = "HMAC"
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
	EncryptionAlgorithmRSA:        {"RSA", "1.2.840.113549.1.1.1", "RSA/ECB/PKCS1Padding"},
	EncryptionAlgorithmRSASSAPSS:  {"RSASSA-PSS", "1.2.840.113549.1.1.10", "RSA/ECB/OAEPPadding"},
	EncryptionAlgorithmDSA:        {"DSA", "1.2.840.10040.4.1", "DSA"},
	EncryptionAlgorithmECDSA:      {"ECDSA", "1.2.840.10045.2.1", "ECDSA"},
	EncryptionAlgorithmPlainECDSA: {"PLAIN-ECDSA", "0.4.0.127.0.7.1.1.4.1", "PLAIN-ECDSA"},
	EncryptionAlgorithmX25519:     {"X25519", "1.3.101.110", "X25519"},
	EncryptionAlgorithmX448:       {"X448", "1.3.101.111", "X448"},
	EncryptionAlgorithmEDDSA:      {"EdDSA", "", "EdDSA"},
	EncryptionAlgorithmHMAC:       {"HMAC", "", ""},
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
		EncryptionAlgorithmRSA,
		EncryptionAlgorithmRSASSAPSS,
		EncryptionAlgorithmDSA,
		EncryptionAlgorithmECDSA,
		EncryptionAlgorithmPlainECDSA,
		EncryptionAlgorithmX25519,
		EncryptionAlgorithmX448,
		EncryptionAlgorithmEDDSA,
		EncryptionAlgorithmHMAC,
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
		return EncryptionAlgorithmECDSA, nil
	}

	// Since JDK 15.
	if name == "Ed25519" || name == "Ed448" {
		return EncryptionAlgorithmEDDSA, nil
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
	return e == EncryptionAlgorithmRSA || e == EncryptionAlgorithmRSASSAPSS
}

func (e EncryptionAlgorithm) isEcDSAFamily() bool {
	return e == EncryptionAlgorithmECDSA || e == EncryptionAlgorithmPlainECDSA
}

func (e EncryptionAlgorithm) isEdDSAFamily() bool {
	return e == EncryptionAlgorithmX25519 || e == EncryptionAlgorithmX448 || e == EncryptionAlgorithmEDDSA
}

// compile-time interface assertion.
var _ OidBasedEnum = EncryptionAlgorithm("")
