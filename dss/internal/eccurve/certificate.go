package eccurve

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// oidPublicKeyECDSA is id-ecPublicKey, the only algorithm whose ECParameters can name one of
// this package's curves.
var oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}

// init restores crypto/x509's pre-Go-1.23 tolerance for a CertificateSerialNumber whose DER
// encoding is technically a negative INTEGER (a leading content byte >= 0x80 with no 0x00 pad -
// RFC 5280 4.1.2.2 requires serial numbers to be non-negative, but does not require an encoder
// to add the pad that keeps a naturally-MSB-set magnitude positive, and several real-world CAs
// do not). Go 1.23 made crypto/x509.ParseCertificate reject such certificates outright unless
// the GODEBUG setting x509negativeserial=1 is present; BouncyCastle (what every upstream DSS
// signature this port cross-validates against was verified with) has no such restriction and
// parses the certificate's ASN.1 INTEGER exactly as encoded, sign included.
//
// A signature-validation library must be able to read every certificate a real CA issued,
// however imperfectly DER-encoded, in order to validate signatures made against it - this is
// not a cryptographic weakening (the certificate's own signature is still fully verified; only
// the serial number's sign convention is relaxed to match what it was actually encoded as) and
// merely brings Go's parser back in line with the BC behaviour every fixture here was produced
// against. internal/godebug (which crypto/x509 reads x509negativeserial through) re-parses
// GODEBUG lazily on first use, so setting the environment variable here - before this program's
// own code ever parses a certificate, since Go runs an imported package's init() before its
// importer's - takes effect for every crypto/x509.ParseCertificate call in the process,
// matching this package's other ParseCertificate accommodation (spliceParsablePublicKey) in
// spirit: never fail to read a certificate stdlib is merely being stricter than necessary about.
//
// Confirmed necessary by pades/testdata/upstream/validation/PAdES-LT.pdf in the PAdES
// cross-validation harness, whose /DSS dictionary carries three certificates with exactly this
// encoding; without this, all three are silently dropped from the certificate pool (Go's
// ParseCertificate error is swallowed the same way every other unparsable-certificate error in
// this codebase's DSS-dictionary/VRI extraction is, per PORTING.md), and one of them is the
// signing certificate for the PDF's oldest signature.
func init() {
	godebug := os.Getenv("GODEBUG")
	for _, setting := range strings.Split(godebug, ",") {
		if strings.HasPrefix(setting, "x509negativeserial=") {
			// The environment (or an earlier init) already has an explicit opinion; respect it.
			return
		}
	}
	if godebug != "" {
		godebug += ","
	}
	os.Setenv("GODEBUG", godebug+"x509negativeserial=1")
}

// ParseCertificate decodes an X.509 certificate, accepting the elliptic curves crypto/x509 does
// not know about in addition to the four it does.
//
// It always tries crypto/x509.ParseCertificate first and returns its result unchanged whenever
// that succeeds, so nothing about the common path changes. Only when stdlib fails does it look
// for the one cause it can repair - a SubjectPublicKeyInfo naming a curve from CurveForOID - and
// even then stdlib remains the only decoder of the certificate's structure: this hands it a
// THROWAWAY copy in which the SubjectPublicKeyInfo has been swapped for a well-formed P-256 one
// purely so the parse completes, then puts the caller's own bytes back on the result and
// replaces PublicKey with the real key on the real curve.
//
// The substitution is required, rather than a mere OID swap, because crypto/x509 rejects the
// public key twice over: once for the unrecognised curve OID, and again because the encoded
// point is not on P-256. Since Raw, RawTBSCertificate and RawSubjectPublicKeyInfo are all
// restored before returning, no re-encoded byte is ever visible to a caller, and in particular
// signature verification still runs over the original RawTBSCertificate.
func ParseCertificate(der []byte) (*x509.Certificate, error) {
	certificate, err := x509.ParseCertificate(der)
	if err == nil {
		return certificate, nil
	}

	patched, originalTBS, originalSPKI, ok := spliceParsablePublicKey(der)
	if !ok {
		return nil, err
	}
	publicKey, keyErr := parseUnsupportedCurvePublicKey(originalSPKI)
	if keyErr != nil {
		return nil, err
	}
	certificate, patchedErr := x509.ParseCertificate(patched)
	if patchedErr != nil {
		// The unknown curve was not the only thing wrong with this certificate; the caller
		// should see the failure stdlib reported for the bytes it actually supplied.
		return nil, err
	}

	certificate.Raw = der
	certificate.RawTBSCertificate = originalTBS
	certificate.RawSubjectPublicKeyInfo = originalSPKI
	certificate.PublicKey = publicKey
	return certificate, nil
}

// spliceParsablePublicKey returns a copy of a Certificate whose SubjectPublicKeyInfo has been
// replaced by a valid P-256 one, together with the original TBSCertificate and
// SubjectPublicKeyInfo elements it was derived from. ok is false whenever der is not a
// structurally plausible certificate, in which case the caller keeps stdlib's own error.
func spliceParsablePublicKey(der []byte) (patched, originalTBS, originalSPKI []byte, ok bool) {
	outer := cryptobyte.String(der)
	var certificate cryptobyte.String
	if !outer.ReadASN1(&certificate, cryptobyte_asn1.SEQUENCE) || !outer.Empty() {
		return nil, nil, nil, false
	}
	rest := certificate
	var tbsCertificate cryptobyte.String
	if !rest.ReadASN1Element(&tbsCertificate, cryptobyte_asn1.SEQUENCE) {
		return nil, nil, nil, false
	}
	originalTBS = []byte(tbsCertificate)
	trailing := []byte(rest) // signatureAlgorithm + signatureValue, copied verbatim

	var tbsContent cryptobyte.String
	inner := cryptobyte.String(originalTBS)
	if !inner.ReadASN1(&tbsContent, cryptobyte_asn1.SEQUENCE) {
		return nil, nil, nil, false
	}
	content := []byte(tbsContent)

	// TBSCertificate ::= SEQUENCE { [0] version OPTIONAL, serialNumber, signature, issuer,
	// validity, subject, subjectPublicKeyInfo, ... }. Walk to the SPKI, keeping everything
	// before and after it byte-for-byte.
	cursor := cryptobyte.String(content)
	var skipped cryptobyte.String
	versionTag := cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()
	if cursor.PeekASN1Tag(versionTag) && !cursor.ReadASN1Element(&skipped, versionTag) {
		return nil, nil, nil, false
	}
	for _, tag := range []cryptobyte_asn1.Tag{
		cryptobyte_asn1.INTEGER,  // serialNumber
		cryptobyte_asn1.SEQUENCE, // signature
		cryptobyte_asn1.SEQUENCE, // issuer
		cryptobyte_asn1.SEQUENCE, // validity
		cryptobyte_asn1.SEQUENCE, // subject
	} {
		if !cursor.ReadASN1Element(&skipped, tag) {
			return nil, nil, nil, false
		}
	}
	prefixLen := len(content) - len(cursor)
	var spki cryptobyte.String
	if !cursor.ReadASN1Element(&spki, cryptobyte_asn1.SEQUENCE) {
		return nil, nil, nil, false
	}
	originalSPKI = []byte(spki)
	suffix := content[prefixLen+len(originalSPKI):]

	placeholder := placeholderSubjectPublicKeyInfo()
	if placeholder == nil {
		return nil, nil, nil, false
	}

	var patchedTBS cryptobyte.Builder
	patchedTBS.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
		child.AddBytes(content[:prefixLen])
		child.AddBytes(placeholder)
		child.AddBytes(suffix)
	})
	patchedTBSBytes, err := patchedTBS.Bytes()
	if err != nil {
		return nil, nil, nil, false
	}

	var patchedCertificate cryptobyte.Builder
	patchedCertificate.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
		child.AddBytes(patchedTBSBytes)
		child.AddBytes(trailing)
	})
	patchedBytes, err := patchedCertificate.Bytes()
	if err != nil {
		return nil, nil, nil, false
	}
	return patchedBytes, originalTBS, originalSPKI, true
}

// placeholderSubjectPublicKeyInfo is a SubjectPublicKeyInfo crypto/x509 is guaranteed to accept:
// the P-256 generator, encoded by the standard library itself so that no DER is hand-written
// here. It never reaches a caller - ParseCertificate overwrites both PublicKey and
// RawSubjectPublicKeyInfo before returning.
func placeholderSubjectPublicKeyInfo() []byte {
	p256 := elliptic.P256().Params()
	encoded, err := x509.MarshalPKIXPublicKey(&ecdsa.PublicKey{Curve: elliptic.P256(), X: p256.Gx, Y: p256.Gy})
	if err != nil {
		return nil
	}
	return encoded
}

// parseUnsupportedCurvePublicKey decodes a SubjectPublicKeyInfo that names one of this package's
// curves, answering an error for anything else - including a well-formed SPKI on a curve
// crypto/x509 already handles, which cannot be the reason stdlib failed.
func parseUnsupportedCurvePublicKey(spkiDER []byte) (*ecdsa.PublicKey, error) {
	var info struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}
	if _, err := asn1.Unmarshal(spkiDER, &info); err != nil {
		return nil, err
	}
	if !info.Algorithm.Algorithm.Equal(oidPublicKeyECDSA) {
		return nil, errors.New("eccurve: SubjectPublicKeyInfo does not carry an EC public key")
	}
	var namedCurveOID asn1.ObjectIdentifier
	if _, err := asn1.Unmarshal(info.Algorithm.Parameters.FullBytes, &namedCurveOID); err != nil {
		return nil, errors.New("eccurve: ECParameters do not name a curve")
	}
	curve := CurveForOID(namedCurveOID.String())
	if curve == nil {
		return nil, errors.New("eccurve: unknown named curve " + namedCurveOID.String())
	}
	if info.PublicKey.BitLength%8 != 0 {
		return nil, errors.New("eccurve: EC point is not a whole number of octets")
	}
	// elliptic.Unmarshal performs both the length check and the on-curve check against the
	// curve's own IsOnCurve, so a point off the group is rejected here rather than reaching
	// ECDSA verification.
	x, y := elliptic.Unmarshal(curve, info.PublicKey.RightAlign())
	if x == nil {
		return nil, errors.New("eccurve: failed to unmarshal elliptic curve point")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}
