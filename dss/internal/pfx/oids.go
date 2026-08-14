package pfx

import "encoding/asn1"

// Content and bag type OIDs (RFC 7292).
var (
	oidData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidEncryptedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 6}

	oidKeyBag                  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 1}
	oidPKCS8ShroudedKeyBag     = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 2}
	oidCertBag                 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 3}
	oidCertTypeX509Certificate = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 22, 1}

	oidFriendlyName = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 20}
	oidLocalKeyID   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 21}
)

// Legacy RFC 7292 Appendix B (id-pkcs12-PbeIds) PBE scheme OIDs.
var (
	oidPbeWithSHAAnd128BitRC4     = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 1}
	oidPbeWithSHAAnd40BitRC4      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 2}
	oidPbeWithSHAAnd3KeyTripleDES = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 3}
	oidPbeWithSHAAnd2KeyTripleDES = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 4}
	oidPbeWithSHAAnd128BitRC2CBC  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 5}
	oidPbeWithSHAAnd40BitRC2CBC   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 6}
)

// PKCS#5 v2.0 (RFC 8018) PBES2/PBKDF2 OIDs.
var (
	oidPBES2  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13}
	oidPBKDF2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12}

	oidHMACWithSHA1   = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 7}
	oidHMACWithSHA224 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 8}
	oidHMACWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 9}
	oidHMACWithSHA384 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 10}
	oidHMACWithSHA512 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 11}
)

// Encryption scheme OIDs PBES2 wraps.
var (
	oidDESEDE3CBC = asn1.ObjectIdentifier{1, 2, 840, 113549, 3, 7}
	oidAES128CBC  = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 2}
	oidAES192CBC  = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 22}
	oidAES256CBC  = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 42}
)

// Plain digest OIDs, used both as a MacData.mac.algorithm and to identify a PBKDF2 PRF's hash.
var (
	oidSHA1   = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
	oidSHA224 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 4}
	oidSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	oidSHA384 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}
	oidSHA512 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}
)

// oidDSA is id-dsa (X9.57/RFC 3279), the one PKCS#8 key algorithm crypto/x509 does not parse.
var oidDSA = asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 1}
