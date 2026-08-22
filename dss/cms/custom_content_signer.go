// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/operator/CustomContentSigner.java
// (DSS 6.5.RC1).
//
// Java implements org.bouncycastle.operator.ContentSigner: getAlgorithmIdentifier() names the
// signature algorithm, getOutputStream() is where a caller (BouncyCastle's
// SignerInfoGeneratorBuilder, here CMSSignerInfoGeneratorBuilder.Build) writes the bytes that
// would be signed, and getSignature() answers the signature - but always the pre-computed one
// supplied at construction, regardless of what was written to the output stream. That is DSS's
// two-step signing: the empty-signature form captures the exact bytes-to-be-sign
// (CAdESService#getDataToSign reads them back off the stream), and the second call, once an
// external signer or HSM has produced the real signature, rebuilds the same CMS with it.
//
// getAlgorithmIdentifier() is Java's `new DefaultSignatureAlgorithmIdentifierFinder().find(algorithmIdentifier)`,
// which maps a JCE algorithm name (e.g. "SHA256withRSA") to its signature AlgorithmIdentifier -
// the OID alone for most families, and for RSASSA-PSS ("...andMGF1") the explicit
// RSASSA-PSS-params SEQUENCE {hashAlgorithm, maskGenAlgorithm, saltLength, trailerField}.
// signatureAlgorithmIdentifiers below is that finder's output for every JCE name
// enumerations.SignatureAlgorithm.JCEID() can produce that BouncyCastle actually recognises,
// captured by running the exact BouncyCastle 1.84 dss-upstream depends on (see
// testdata/gen/GenSignatureAlgorithmIdentifiers.java) rather than hand re-deriving the ASN.1 -
// a hand-rolled construction of RSASSA-PSS-params or the X9.62 ECDSA/EdDSA identifiers has no
// authoritative source to check itself against, while a real BouncyCastle run does. Four
// families are absent from the map because BouncyCastle's finder itself raises
// "Unknown signature type requested" for them (verified the same way): the four NONEwith*
// "raw" algorithms (no digest to name) and RIPEMD160withECDSA.
package cms

import (
	"bytes"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// signatureAlgorithmIdentifierHex maps a JCE signature algorithm name (SignatureAlgorithm.JCEID())
// to the DER encoding of its signature AlgorithmIdentifier, i.e.
// org.bouncycastle.operator.DefaultSignatureAlgorithmIdentifierFinder#find(String)'s result for
// BouncyCastle 1.84.
var signatureAlgorithmIdentifierHex = map[string]string{
	"SHA1withRSA":              "300d06092a864886f70d0101050500",
	"SHA224withRSA":            "300d06092a864886f70d01010e0500",
	"SHA256withRSA":            "300d06092a864886f70d01010b0500",
	"SHA384withRSA":            "300d06092a864886f70d01010c0500",
	"SHA512withRSA":            "300d06092a864886f70d01010d0500",
	"SHA3-224withRSA":          "300d060960864801650304030d0500",
	"SHA3-256withRSA":          "300d060960864801650304030e0500",
	"SHA3-384withRSA":          "300d060960864801650304030f0500",
	"SHA3-512withRSA":          "300d06096086480165030403100500",
	"SHA1withRSAandMGF1":       "300d06092a864886f70d01010a3000",
	"SHA224withRSAandMGF1":     "304106092a864886f70d01010a3034a00f300d06096086480165030402040500a11c301a06092a864886f70d010108300d06096086480165030402040500a20302011c",
	"SHA256withRSAandMGF1":     "304106092a864886f70d01010a3034a00f300d06096086480165030402010500a11c301a06092a864886f70d010108300d06096086480165030402010500a203020120",
	"SHA384withRSAandMGF1":     "304106092a864886f70d01010a3034a00f300d06096086480165030402020500a11c301a06092a864886f70d010108300d06096086480165030402020500a203020130",
	"SHA512withRSAandMGF1":     "304106092a864886f70d01010a3034a00f300d06096086480165030402030500a11c301a06092a864886f70d010108300d06096086480165030402030500a203020140",
	"SHA3-224withRSAandMGF1":   "304106092a864886f70d01010a3034a00f300d06096086480165030402070500a11c301a06092a864886f70d010108300d06096086480165030402070500a20302011c",
	"SHA3-256withRSAandMGF1":   "304106092a864886f70d01010a3034a00f300d06096086480165030402080500a11c301a06092a864886f70d010108300d06096086480165030402080500a203020120",
	"SHA3-384withRSAandMGF1":   "304106092a864886f70d01010a3034a00f300d06096086480165030402090500a11c301a06092a864886f70d010108300d06096086480165030402090500a203020130",
	"SHA3-512withRSAandMGF1":   "304106092a864886f70d01010a3034a00f300d060960864801650304020a0500a11c301a06092a864886f70d010108300d060960864801650304020a0500a203020140",
	"RIPEMD160withRSA":         "300a06062b24030301020500",
	"MD5withRSA":               "300d06092a864886f70d0101040500",
	"MD2withRSA":               "300d06092a864886f70d0101020500",
	"SHA1withECDSA":            "300906072a8648ce3d0401",
	"SHA224withECDSA":          "300a06082a8648ce3d040301",
	"SHA256withECDSA":          "300a06082a8648ce3d040302",
	"SHA384withECDSA":          "300a06082a8648ce3d040303",
	"SHA512withECDSA":          "300a06082a8648ce3d040304",
	"SHA3-224withECDSA":        "300b0609608648016503040309",
	"SHA3-256withECDSA":        "300b060960864801650304030a",
	"SHA3-384withECDSA":        "300b060960864801650304030b",
	"SHA3-512withECDSA":        "300b060960864801650304030c",
	"SHA1withPLAIN-ECDSA":      "300e060a04007f000701010401010500",
	"SHA224withPLAIN-ECDSA":    "300c060a04007f00070101040102",
	"SHA256withPLAIN-ECDSA":    "300c060a04007f00070101040103",
	"SHA384withPLAIN-ECDSA":    "300c060a04007f00070101040104",
	"SHA512withPLAIN-ECDSA":    "300c060a04007f00070101040105",
	"RIPEMD160withPLAIN-ECDSA": "300e060a04007f000701010401060500",
	"SHA3-224withPLAIN-ECDSA":  "300c060a04007f00070101040108",
	"SHA3-256withPLAIN-ECDSA":  "300c060a04007f00070101040109",
	"SHA3-384withPLAIN-ECDSA":  "300c060a04007f0007010104010a",
	"SHA3-512withPLAIN-ECDSA":  "300c060a04007f0007010104010b",
	"Ed25519":                  "300506032b6570",
	"Ed448":                    "300506032b6571",
	"SHA1withDSA":              "300906072a8648ce380403",
	"SHA224withDSA":            "300b0609608648016503040301",
	"SHA256withDSA":            "300b0609608648016503040302",
	"SHA384withDSA":            "300b0609608648016503040303",
	"SHA512withDSA":            "300b0609608648016503040304",
	"SHA3-224withDSA":          "300b0609608648016503040305",
	"SHA3-256withDSA":          "300b0609608648016503040306",
	"SHA3-384withDSA":          "300b0609608648016503040307",
	"SHA3-512withDSA":          "300b0609608648016503040308",
}

// signatureAlgorithmIdentifiers is signatureAlgorithmIdentifierHex, parsed once.
var signatureAlgorithmIdentifiers = mustParseSignatureAlgorithmIdentifiers()

func mustParseSignatureAlgorithmIdentifiers() map[string]*asn1ber.AlgorithmIdentifier {
	table := make(map[string]*asn1ber.AlgorithmIdentifier, len(signatureAlgorithmIdentifierHex))
	for jceID, hexDER := range signatureAlgorithmIdentifierHex {
		der := make([]byte, len(hexDER)/2)
		if _, err := fmt.Sscanf(hexDER, "%x", &der); err != nil {
			panic(fmt.Sprintf("cms: malformed signatureAlgorithmIdentifierHex entry %q: %s", jceID, err))
		}
		algorithmIdentifier, err := asn1ber.ParseAlgorithmIdentifier(der)
		if err != nil {
			panic(fmt.Sprintf("cms: malformed signatureAlgorithmIdentifierHex entry %q: %s", jceID, err))
		}
		table[jceID] = algorithmIdentifier
	}
	return table
}

// findSignatureAlgorithmIdentifier is
// DefaultSignatureAlgorithmIdentifierFinder#find(String)'s outcome for BouncyCastle 1.84: the
// signature AlgorithmIdentifier for a JCE algorithm name, or an error carrying BouncyCastle's
// own IllegalArgumentException message for a name it does not recognise.
func findSignatureAlgorithmIdentifier(jceAlgorithmIdentifier string) (*asn1ber.AlgorithmIdentifier, error) {
	if identifier, ok := signatureAlgorithmIdentifiers[jceAlgorithmIdentifier]; ok {
		return identifier, nil
	}
	return nil, fmt.Errorf("Unknown signature type requested: %s", jceAlgorithmIdentifier)
}

// ContentSigner is the minimal replacement for org.bouncycastle.operator.ContentSigner: an
// object that names a signature algorithm, receives the bytes that would be signed through
// OutputStream, and answers the signature. See the file header for why DSS's own producers
// (CustomContentSigner, the only one) always answer a pre-computed signature.
type ContentSigner interface {
	// AlgorithmIdentifier returns the signature algorithm identifier. Port of
	// ContentSigner#getAlgorithmIdentifier.
	AlgorithmIdentifier() *asn1ber.AlgorithmIdentifier

	// OutputStream returns the buffer the bytes-to-be-signed are written to. Port of
	// ContentSigner#getOutputStream; Java returns a plain OutputStream, a *bytes.Buffer is used
	// here so that a caller building the "data to sign" can read the accumulated bytes back
	// with Bytes(), which is what CAdESService#getDataToSign needs.
	OutputStream() *bytes.Buffer

	// Signature returns the signature. Port of ContentSigner#getSignature.
	Signature() []byte
}

// CustomContentSigner is a ContentSigner using a provided pre-computed signature.
// Port of the CustomContentSigner class.
type CustomContentSigner struct {
	// algorithmIdentifier is the signature algorithm identifier.
	algorithmIdentifier *asn1ber.AlgorithmIdentifier
	// preComputedSignature is the pre-computed SignatureValue.
	preComputedSignature []byte
	// outputStream is the buffer used to write the data to.
	outputStream *bytes.Buffer
}

// NewCustomContentSigner is the default constructor for CustomContentSigner, with an absent
// signature. Port of CustomContentSigner(String), i.e.
// CustomContentSigner(algorithmIdentifier, DSSUtils.EMPTY_BYTE_ARRAY).
//
// algorithmIdentifier is the JCE algorithm identifier (SignatureAlgorithm.JCEID()); the
// returned error carries BouncyCastle's own message for an algorithm it cannot map to an
// AlgorithmIdentifier (see findSignatureAlgorithmIdentifier).
func NewCustomContentSigner(algorithmIdentifier string) (*CustomContentSigner, error) {
	return NewCustomContentSignerWithSignature(algorithmIdentifier, []byte{})
}

// NewCustomContentSignerWithSignature is the constructor for CustomContentSigner using the
// real value of the signature. Port of CustomContentSigner(String, byte[]).
func NewCustomContentSignerWithSignature(algorithmIdentifier string, preComputedSignature []byte) (*CustomContentSigner, error) {
	identifier, err := findSignatureAlgorithmIdentifier(algorithmIdentifier)
	if err != nil {
		return nil, err
	}
	return &CustomContentSigner{
		algorithmIdentifier:  identifier,
		preComputedSignature: preComputedSignature,
		outputStream:         new(bytes.Buffer),
	}, nil
}

// AlgorithmIdentifier returns the signature algorithm identifier. Port of #getAlgorithmIdentifier.
func (s *CustomContentSigner) AlgorithmIdentifier() *asn1ber.AlgorithmIdentifier {
	return s.algorithmIdentifier
}

// OutputStream returns the byteOutputStream used to write the data to. Port of #getOutputStream.
func (s *CustomContentSigner) OutputStream() *bytes.Buffer { return s.outputStream }

// Signature returns the preComputedSignature. Port of #getSignature.
func (s *CustomContentSigner) Signature() []byte { return s.preComputedSignature }

var _ ContentSigner = (*CustomContentSigner)(nil)
