package spi

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
)

// dssASN1UtilsTestKAT loads the known-answer file produced against BouncyCastle 1.78.1 and
// OpenJDK 21 (testdata/asn1/kat.txt): one "key|value" pair per line.
func dssASN1UtilsTestKAT(t *testing.T) map[string]string {
	t.Helper()
	file, err := os.Open(corpustest.Path(t, filepath.Join("asn1", "kat.txt")))
	if err != nil {
		t.Fatalf("unable to open the known-answer file: %v", err)
	}
	defer func() { _ = file.Close() }()

	answers := make(map[string]string)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, "|")
		if index < 0 {
			continue
		}
		answers[line[:index]] = line[index+1:]
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unable to read the known-answer file: %v", err)
	}
	return answers
}

// dssASN1UtilsTestHex decodes a hex string from the known-answer file.
func dssASN1UtilsTestHex(t *testing.T, answers map[string]string, key string) []byte {
	t.Helper()
	value, ok := answers[key]
	if !ok {
		t.Fatalf("the known-answer file has no entry %q", key)
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("entry %q is not hex: %v", key, err)
	}
	return decoded
}

// dssASN1UtilsTestCertificate loads one of the DER certificates of testdata/asn1.
func dssASN1UtilsTestCertificate(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("testdata", "asn1", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// TestDSSASN1UtilsAlgorithmIdentifierForDigest checks every digest algorithm identifier
// against BouncyCastle's DefaultDigestAlgorithmIdentifierFinder (bcpkix 1.78.1), plus the
// SHAKE256-512 special case DSS-3651 introduced.
func TestDSSASN1UtilsAlgorithmIdentifierForDigest(t *testing.T) {
	expected := map[enumerations.DigestAlgorithm]string{
		enumerations.DigestAlgorithm_SHA1:         "300906052b0e03021a0500",
		enumerations.DigestAlgorithm_SHA224:       "300b0609608648016503040204",
		enumerations.DigestAlgorithm_SHA256:       "300b0609608648016503040201",
		enumerations.DigestAlgorithm_SHA384:       "300b0609608648016503040202",
		enumerations.DigestAlgorithm_SHA512:       "300b0609608648016503040203",
		enumerations.DigestAlgorithm_SHA3_224:     "300b0609608648016503040207",
		enumerations.DigestAlgorithm_SHA3_256:     "300b0609608648016503040208",
		enumerations.DigestAlgorithm_SHA3_384:     "300b0609608648016503040209",
		enumerations.DigestAlgorithm_SHA3_512:     "300b060960864801650304020a",
		enumerations.DigestAlgorithm_SHAKE128:     "300b060960864801650304020b",
		enumerations.DigestAlgorithm_SHAKE256:     "300b060960864801650304020c",
		enumerations.DigestAlgorithm_RIPEMD160:    "300906052b240302010500",
		enumerations.DigestAlgorithm_MD2:          "300c06082a864886f70d02020500",
		enumerations.DigestAlgorithm_MD5:          "300c06082a864886f70d02050500",
		enumerations.DigestAlgorithm_WHIRLPOOL:    "3008060628cf06030037",
		enumerations.DigestAlgorithm_SHAKE256_512: "300f060960864801650304021202020200",
	}
	for _, digestAlgorithm := range enumerations.DigestAlgorithmValues() {
		identifier, err := DSSASN1UtilsAlgorithmIdentifierForDigest(digestAlgorithm)
		if err != nil {
			t.Fatalf("%s: %v", digestAlgorithm, err)
		}
		if got := hex.EncodeToString(identifier.DER()); got != expected[digestAlgorithm] {
			t.Errorf("%s: got %s, want %s", digestAlgorithm, got, expected[digestAlgorithm])
		}
		// The structure must survive a parse/encode round trip.
		parsed, err := ParseAlgorithmIdentifier(identifier.DER())
		if err != nil {
			t.Fatalf("%s: unable to parse back: %v", digestAlgorithm, err)
		}
		if !parsed.Equals(identifier) {
			t.Errorf("%s: the AlgorithmIdentifier did not survive a round trip", digestAlgorithm)
		}
	}
}

// TestDSSASN1UtilsEncodings checks the BER-to-DER/DL/BER re-encodings against the answers
// BouncyCastle's ASN1Primitive#getEncoded(String) produces for the same inputs.
func TestDSSASN1UtilsEncodings(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	compared := 0
	for index := 0; index <= 8; index++ {
		suffix := strconv.Itoa(index)
		input, ok := answers["ber_in_"+suffix]
		if !ok {
			t.Fatalf("the known-answer file has no entry ber_in_%s", suffix)
		}
		binaries, err := hex.DecodeString(input)
		if err != nil {
			t.Fatalf("ber_in_%s is not hex: %v", suffix, err)
		}
		for _, encoding := range []struct {
			key     string
			convert func([]byte) ([]byte, error)
		}{
			{"ber_der_" + suffix, DSSASN1UtilsDEREncoded},
			{"ber_dl_" + suffix, DSSASN1UtilsDLEncoded},
			{"ber_ber_" + suffix, DSSASN1UtilsBEREncoded},
		} {
			want, ok := answers[encoding.key]
			if !ok {
				continue
			}
			got, err := encoding.convert(binaries)
			if err != nil {
				t.Fatalf("%s: %v", encoding.key, err)
			}
			if hex.EncodeToString(got) != want {
				t.Errorf("%s:\n got %s\nwant %s", encoding.key, hex.EncodeToString(got), want)
			}
			compared++
		}
	}
	if compared < 20 {
		t.Fatalf("only %d encodings were compared; the known-answer file looks truncated", compared)
	}
}

// TestDSSASN1UtilsDEREncodedIsIdempotent checks that a DER input is returned unchanged, so
// that re-encoding a certificate or a timestamp never perturbs a digest.
func TestDSSASN1UtilsDEREncodedIsIdempotent(t *testing.T) {
	for _, name := range []string{"subject.der", "issuer.der", "ocsp_ca.der", "ocsp_resp_byname.der"} {
		der, err := os.ReadFile(filepath.Join("testdata", "asn1", name))
		if err != nil {
			t.Fatalf("unable to read %s: %v", name, err)
		}
		encoded, err := DSSASN1UtilsDEREncoded(der)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !bytes.Equal(der, encoded) {
			t.Errorf("%s: the DER encoding changed the binaries", name)
		}
	}
}

// TestDSSASN1UtilsDEREncodedRejectsTrailingData checks the "extra data" guard of the
// re-encoding entry points.
func TestDSSASN1UtilsDEREncodedRejectsTrailingData(t *testing.T) {
	if _, err := DSSASN1UtilsDEREncoded([]byte{0x05, 0x00, 0x05, 0x00}); err == nil {
		t.Error("two concatenated objects must be rejected")
	}
	if _, err := DSSASN1UtilsToASN1Primitive([]byte{0x05, 0x00, 0x2A}); err == nil {
		t.Error("trailing bytes must be rejected")
	}
	if _, err := DSSASN1UtilsToASN1Primitive([]byte{0x30, 0x05, 0x02}); err == nil {
		t.Error("a truncated object must be rejected")
	}
}

// TestDSSASN1UtilsDSASignatureValue checks the ASN.1 <-> plain (R || S) conversions against
// BouncyCastle's PlainDSAEncoding.
func TestDSSASN1UtilsDSASignatureValue(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	for index := 0; index < 3; index++ {
		suffix := strconv.Itoa(index)
		asn1Value := dssASN1UtilsTestHex(t, answers, "dsa_asn1_"+suffix)
		wantPlain := dssASN1UtilsTestHex(t, answers, "dsa_plain_"+suffix)
		wantOrder, ok := new(big.Int).SetString(answers["dsa_order_"+suffix], 10)
		if !ok {
			t.Fatalf("dsa_order_%s is not a number", suffix)
		}

		if !DSSASN1UtilsIsAsn1EncodedSignatureValue(asn1Value) {
			t.Errorf("dsa_asn1_%s must be recognised as an ASN.1 signature value", suffix)
		}
		order, err := DSSASN1UtilsOrderFromSignatureValue(asn1Value)
		if err != nil {
			t.Fatalf("dsa_asn1_%s: %v", suffix, err)
		}
		if order.Cmp(wantOrder) != 0 {
			t.Errorf("dsa_order_%s: got %s, want %s", suffix, order, wantOrder)
		}
		plain, err := DSSASN1UtilsToPlainDSASignatureValue(asn1Value)
		if err != nil {
			t.Fatalf("dsa_plain_%s: %v", suffix, err)
		}
		if !bytes.Equal(plain, wantPlain) {
			t.Errorf("dsa_plain_%s:\n got %x\nwant %x", suffix, plain, wantPlain)
		}
		// The plain form must be sized to the order.
		if len(plain)%2 != 0 || len(plain)/2 != (wantOrder.BitLen()+7)/8 {
			t.Errorf("dsa_plain_%s: unexpected length %d for an order of %d bits",
				suffix, len(plain), wantOrder.BitLen())
		}
		bitLength, err := DSSASN1UtilsSignatureValueBitLength(asn1Value)
		if err != nil {
			t.Fatalf("dsa bit length %s: %v", suffix, err)
		}
		if bitLength != (wantOrder.BitLen()+7)/8*8 {
			t.Errorf("dsa bit length %s: got %d", suffix, bitLength)
		}
		// The plain form must resolve to the same order.
		plainOrder, err := DSSASN1UtilsOrderFromSignatureValue(plain)
		if err != nil {
			t.Fatalf("dsa plain order %s: %v", suffix, err)
		}
		if plainOrder.Cmp(wantOrder) != 0 {
			t.Errorf("dsa plain order %s: got %s, want %s", suffix, plainOrder, wantOrder)
		}
	}
}

// TestDSSASN1UtilsToStandardDSASignatureValue checks the plain-to-ASN.1 conversion.
func TestDSSASN1UtilsToStandardDSASignatureValue(t *testing.T) {
	plain, err := hex.DecodeString("0102")
	if err != nil {
		t.Fatal(err)
	}
	standard, err := DSSASN1UtilsToStandardDSASignatureValue(plain)
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(standard); got != "3006020101020102" {
		t.Errorf("got %s, want 3006020101020102", got)
	}
	// A value whose high bit is set must gain a leading zero octet, as an ASN.1 INTEGER is
	// signed. 0x80 || 0x01 is r = 128, s = 1.
	standard, err = DSSASN1UtilsToStandardDSASignatureValue([]byte{0x80, 0x01})
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(standard); got != "300702020080020101" {
		t.Errorf("got %s, want 300702020080020101", got)
	}
	if _, err := DSSASN1UtilsToStandardDSASignatureValue(nil); err == nil {
		t.Error("an empty signature value must be rejected")
	}
}

// TestDSSASN1UtilsEnsurePlainSignatureValue checks that only the DSA-family algorithms are
// converted, and only when the value really is ASN.1 encoded.
func TestDSSASN1UtilsEnsurePlainSignatureValue(t *testing.T) {
	asn1Value, err := hex.DecodeString("3006020101020102")
	if err != nil {
		t.Fatal(err)
	}
	for _, algorithm := range []enumerations.EncryptionAlgorithm{
		enumerations.EncryptionAlgorithm_ECDSA,
		enumerations.EncryptionAlgorithm_PLAIN_ECDSA,
		enumerations.EncryptionAlgorithm_DSA,
	} {
		value, err := DSSASN1UtilsEnsurePlainSignatureValue(algorithm, asn1Value)
		if err != nil {
			t.Fatalf("%s: %v", algorithm, err)
		}
		if got := hex.EncodeToString(value); got != "0102" {
			t.Errorf("%s: got %s, want 0102", algorithm, got)
		}
	}
	value, err := DSSASN1UtilsEnsurePlainSignatureValue(enumerations.EncryptionAlgorithm_RSA, asn1Value)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(value, asn1Value) {
		t.Error("an RSA signature value must be left untouched")
	}
	// A plain value stays plain, even for ECDSA.
	plain := []byte{0x01, 0x02}
	value, err = DSSASN1UtilsEnsurePlainSignatureValue(enumerations.EncryptionAlgorithm_ECDSA, plain)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(value, plain) {
		t.Error("a plain signature value must be left untouched")
	}
}

// TestDSSASN1UtilsComputeSkiFromCert checks the SHA-1 of the subjectPublicKey against the
// value Java computes from SubjectPublicKeyInfo#getPublicKeyData().
func TestDSSASN1UtilsComputeSkiFromCert(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	for _, entry := range []struct{ file, key string }{
		{"issuer.der", "ski_issuer"},
		{"subject.der", "ski_subject"},
	} {
		token := dssASN1UtilsTestCertificate(t, entry.file)
		ski, err := DSSASN1UtilsComputeSkiFromCert(token)
		if err != nil {
			t.Fatalf("%s: %v", entry.file, err)
		}
		want := dssASN1UtilsTestHex(t, answers, entry.key)
		if !bytes.Equal(ski, want) {
			t.Errorf("%s:\n got %x\nwant %x", entry.file, ski, want)
		}
		if !DSSASN1UtilsIsSkiEqual(want, token) {
			t.Errorf("%s: the SKI must compare equal", entry.file)
		}
		if DSSASN1UtilsIsSkiEqual([]byte{0x00}, token) {
			t.Errorf("%s: a foreign SKI must not compare equal", entry.file)
		}
	}
}

// TestDSSASN1UtilsIssuerSerial checks the IssuerSerial built from a certificate against
// BouncyCastle's IssuerSerial(GeneralNames, BigInteger), and the parsing back.
func TestDSSASN1UtilsIssuerSerial(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	token := dssASN1UtilsTestCertificate(t, "subject.der")

	issuerSerial := DSSASN1UtilsIssuerSerialForCertificate(token)
	want := dssASN1UtilsTestHex(t, answers, "issuer_serial")
	if !bytes.Equal(issuerSerial.DER(), want) {
		t.Errorf("IssuerSerial:\n got %x\nwant %x", issuerSerial.DER(), want)
	}

	parsed := DSSASN1UtilsIssuerSerial(want)
	if parsed == nil {
		t.Fatal("the IssuerSerial could not be parsed back")
	}
	if !bytes.Equal(parsed.DER(), want) {
		t.Errorf("the IssuerSerial did not survive a round trip:\n got %x\nwant %x", parsed.DER(), want)
	}
	if parsed.Serial.Cmp(token.SerialNumber()) != 0 {
		t.Errorf("serial: got %s, want %s", parsed.Serial, token.SerialNumber())
	}
	if len(parsed.Issuer) != 1 || parsed.Issuer[0].TagNo != 4 {
		t.Fatalf("expected a single directoryName GeneralName, got %+v", parsed.Issuer)
	}
	if !bytes.Equal(parsed.Issuer[0].Name, token.Certificate().RawIssuer) {
		t.Error("the parsed issuer name does not match the certificate's issuer")
	}

	signerIdentifier := DSSASN1UtilsToSignerIdentifierFromIssuerSerial(parsed)
	if signerIdentifier == nil {
		t.Fatal("the SignerIdentifier could not be built")
	}
	if DSSASN1UtilsToSignerIdentifierFromIssuerSerial(nil) != nil {
		t.Error("a nil IssuerSerial must yield a nil SignerIdentifier")
	}
	if DSSASN1UtilsIssuerSerial([]byte{0x05, 0x00}) != nil {
		t.Error("binaries that are not an IssuerSerial must yield nil")
	}
}

// TestDSSASN1UtilsString checks getString against BouncyCastle's IETFUtils#valueToString
// followed by String#trim, for a value of every relevant ASN.1 type.
func TestDSSASN1UtilsString(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	for index := 0; index < 10; index++ {
		suffix := strconv.Itoa(index)
		input := dssASN1UtilsTestHex(t, answers, "v2s_in_"+suffix)
		want := answers["v2s_trim_"+suffix]
		if got := DSSASN1UtilsString(input); got != want {
			t.Errorf("v2s_%s: got %q, want %q", suffix, got, want)
		}
	}
	if got := DSSASN1UtilsString(nil); got != "" {
		t.Errorf("a nil value must yield the empty string, got %q", got)
	}
	if got := DSSASN1UtilsString([]byte{0x30, 0x05}); got != "" {
		t.Errorf("a malformed value must yield the empty string, got %q", got)
	}
}

// TestDSSASN1UtilsExtractAttributeFromX500Principal checks the attribute extraction and the
// human-readable-name helpers against BouncyCastle's X500Name#getRDNs.
func TestDSSASN1UtilsExtractAttributeFromX500Principal(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	token := dssASN1UtilsTestCertificate(t, "subject.der")

	for _, entry := range []struct {
		key string
		oid asn1.ObjectIdentifier
	}{
		{"attr_CN", asn1.ObjectIdentifier{2, 5, 4, 3}},
		{"attr_GIVENNAME", asn1.ObjectIdentifier{2, 5, 4, 42}},
		{"attr_SURNAME", asn1.ObjectIdentifier{2, 5, 4, 4}},
		{"attr_O", asn1.ObjectIdentifier{2, 5, 4, 10}},
		{"attr_OU", asn1.ObjectIdentifier{2, 5, 4, 11}},
		{"attr_L", asn1.ObjectIdentifier{2, 5, 4, 7}},
		{"attr_C", asn1.ObjectIdentifier{2, 5, 4, 6}},
	} {
		want := answers[entry.key]
		if want == "<null>" {
			want = ""
		}
		got := DSSASN1UtilsExtractAttributeFromX500Principal(entry.oid, token.Subject())
		if got != want {
			t.Errorf("%s: got %q, want %q", entry.key, got, want)
		}
	}

	if got := DSSASN1UtilsSubjectCommonName(token); got != answers["attr_CN"] {
		t.Errorf("subject common name: got %q, want %q", got, answers["attr_CN"])
	}
	// CN is present, so the pretty-printed name is the CN.
	if got := DSSASN1UtilsHumanReadableName(token); got != answers["attr_CN"] {
		t.Errorf("human readable name: got %q, want %q", got, answers["attr_CN"])
	}
}

// TestDSSASN1UtilsAttributeMapAndEquality checks the attribute map and the X500Principal
// comparison it backs.
func TestDSSASN1UtilsAttributeMapAndEquality(t *testing.T) {
	subject := dssASN1UtilsTestCertificate(t, "subject.der")
	issuer := dssASN1UtilsTestCertificate(t, "issuer.der")

	attributes := DSSASN1UtilsAttributeMap(subject.Subject().Principal())
	if len(attributes) != 7 {
		t.Errorf("expected 7 attributes, got %d: %v", len(attributes), attributes)
	}
	// UPSTREAM QUIRK: the keys are the hex of the attribute type's DER, not dotted OIDs.
	// 2.5.4.3 (CN) encodes as 06 03 55 04 03.
	commonName, ok := attributes["#0603550403"]
	if !ok {
		t.Fatalf("the CN key must be the hex of its DER encoding, got the keys %v", attributes)
	}
	if commonName != "Ján Procháska" {
		t.Errorf("CN: got %q", commonName)
	}

	if !DSSASN1UtilsX500PrincipalAreEquals(subject.Subject().Principal(), subject.Subject().Principal()) {
		t.Error("a principal must equal itself")
	}
	if DSSASN1UtilsX500PrincipalAreEquals(subject.Subject().Principal(), issuer.Subject().Principal()) {
		t.Error("two different principals must not be equal")
	}
	if DSSASN1UtilsX500PrincipalAreEquals(nil, subject.Subject().Principal()) {
		t.Error("a nil principal must never be equal")
	}
	// The subject of the certificate and its own issuer's field are different names, but a
	// principal rebuilt from the same DER must compare equal.
	rebuilt, err := DSSASN1UtilsToX500Principal(subject.Certificate().RawSubject)
	if err != nil {
		t.Fatal(err)
	}
	if !DSSASN1UtilsX500PrincipalAreEquals(subject.Subject().Principal(), rebuilt) {
		t.Error("a principal rebuilt from the same DER must be equal")
	}
	if _, err := DSSASN1UtilsToX500Principal([]byte{0x05, 0x00}); err == nil {
		t.Error("binaries that are not a Name must be rejected")
	}
}

// TestDSSASN1UtilsBuildSPDocSpecificationID checks the OID/URI CHOICE against the DER
// BouncyCastle produces for the same strings.
func TestDSSASN1UtilsBuildSPDocSpecificationID(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	for _, value := range []string{"1.2.3.4.5", "2.999", "http://nowina.lu/policy", "1.2.03.4", "0.4.0.19122.2.1"} {
		want := dssASN1UtilsTestHex(t, answers, "spdoc_"+value)
		got, err := DSSASN1UtilsBuildSPDocSpecificationID(value)
		if err != nil {
			t.Fatalf("%s: %v", value, err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s:\n got %x\nwant %x", value, got, want)
		}
	}
	// DSSUtils#isOidCode accepts a single arc, ASN.1 does not: upstream throws here too.
	if _, err := DSSASN1UtilsBuildSPDocSpecificationID("0"); err == nil {
		t.Error("a single-arc OID must be rejected")
	}
}

// TestDSSASN1UtilsDates checks the UTCTime and GeneralizedTime readers against the instants
// BouncyCastle's Time#getDate produces, including the two-digit-year rule of UTCTime.
func TestDSSASN1UtilsDates(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	for key, value := range answers {
		var derKey string
		switch {
		case strings.HasPrefix(key, "gtime_der_"), strings.HasPrefix(key, "utctime_der_"):
			continue
		case strings.HasPrefix(key, "gtime_"):
			derKey = "gtime_der_" + strings.TrimPrefix(key, "gtime_")
		case strings.HasPrefix(key, "utctime_"):
			derKey = "utctime_der_" + strings.TrimPrefix(key, "utctime_")
		default:
			continue
		}
		millis, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			t.Fatalf("%s is not a number: %v", key, err)
		}
		want := time.UnixMilli(millis).UTC()

		// getDate(ASN1Encodable) handles both time types.
		got := DSSASN1UtilsDate(dssASN1UtilsTestHex(t, answers, derKey))
		if !got.Equal(want) {
			t.Errorf("%s: got %s, want %s", key, got.UTC(), want)
		}

		if strings.HasPrefix(key, "gtime_") {
			// toDate(ASN1GeneralizedTime) only handles GeneralizedTime.
			got, err := DSSASN1UtilsToDate(dssASN1UtilsTestHex(t, answers, derKey))
			if err != nil {
				t.Fatalf("%s: %v", key, err)
			}
			if !got.Equal(want) {
				t.Errorf("%s (toDate): got %s, want %s", key, got.UTC(), want)
			}
		}
	}

	// A UTCTime is not a GeneralizedTime.
	if _, err := DSSASN1UtilsToDate(dssASN1UtilsTestHex(t, answers, "utctime_der_200102030405Z")); err == nil {
		t.Error("toDate must reject a UTCTime")
	}
	// getDate returns the zero time instead of raising.
	if got := DSSASN1UtilsDate([]byte{0x05, 0x00}); !got.IsZero() {
		t.Errorf("getDate must yield the zero time for a non-time value, got %s", got)
	}
}

// TestDSSASN1UtilsAsn1SignaturePolicyDigest checks the TS 101 733 5.8.1 policy digest.
func TestDSSASN1UtilsAsn1SignaturePolicyDigest(t *testing.T) {
	answers := dssASN1UtilsTestKAT(t)
	policyBytes := dssASN1UtilsTestHex(t, answers, "policy_bytes")
	want := dssASN1UtilsTestHex(t, answers, "policy_digest_sha256")
	got, err := DSSASN1UtilsAsn1SignaturePolicyDigest(enumerations.DigestAlgorithm_SHA256, policyBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("\n got %x\nwant %x", got, want)
	}
	if _, err := DSSASN1UtilsAsn1SignaturePolicyDigest(enumerations.DigestAlgorithm_SHA256, []byte{0x05, 0x00}); err == nil {
		t.Error("a policy that is not a SEQUENCE must be rejected")
	}
}

// TestDSSASN1UtilsOctetStringHelpers checks the DEROctetString accessors.
func TestDSSASN1UtilsOctetStringHelpers(t *testing.T) {
	// 04 02 05 00 : an OCTET STRING encapsulating DERNull.
	if !DSSASN1UtilsIsDEROctetStringNull([]byte{0x04, 0x02, 0x05, 0x00}) {
		t.Error("an OCTET STRING holding NULL must be recognised")
	}
	if DSSASN1UtilsIsDEROctetStringNull([]byte{0x04, 0x03, 0x02, 0x01, 0x07}) {
		t.Error("an OCTET STRING holding an INTEGER is not NULL")
	}
	if DSSASN1UtilsIsDEROctetStringNull([]byte{0x05, 0x00}) {
		t.Error("a value that is not an OCTET STRING must answer false")
	}

	// 04 08 30 06 02 01 07 02 01 09 : an OCTET STRING encapsulating a SEQUENCE.
	encapsulated := []byte{0x04, 0x08, 0x30, 0x06, 0x02, 0x01, 0x07, 0x02, 0x01, 0x09}
	sequence, err := DSSASN1UtilsAsn1SequenceFromDerOctetString(encapsulated)
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(sequence); got != "3006020107020109" {
		t.Errorf("got %s", got)
	}
	if _, err := DSSASN1UtilsAsn1SequenceFromDerOctetString([]byte{0x04, 0x02, 0x05, 0x00}); err == nil {
		t.Error("an encapsulated NULL is not a SEQUENCE")
	}

	// 04 03 02 01 07 : an OCTET STRING encapsulating INTEGER 7.
	integer, err := DSSASN1UtilsAsn1IntegerFromDerOctetString([]byte{0x04, 0x03, 0x02, 0x01, 0x07})
	if err != nil {
		t.Fatal(err)
	}
	if integer.Int64() != 7 {
		t.Errorf("got %s, want 7", integer)
	}
	// A negative INTEGER must decode as two's complement.
	integer, err = DSSASN1UtilsAsn1IntegerFromDerOctetString([]byte{0x04, 0x03, 0x02, 0x01, 0xFF})
	if err != nil {
		t.Fatal(err)
	}
	if integer.Int64() != -1 {
		t.Errorf("got %s, want -1", integer)
	}

	value, err := DSSASN1UtilsToString([]byte{0x04, 0x03, 'a', 'b', 'c'})
	if err != nil {
		t.Fatal(err)
	}
	if value != "abc" {
		t.Errorf("got %q, want abc", value)
	}

	// 30 09 04 02 01 02 04 03 03 04 05 : a SEQUENCE of two OCTET STRINGs.
	octetStrings, err := DSSASN1UtilsDEROctetStrings([]byte{0x30, 0x09, 0x04, 0x02, 0x01, 0x02, 0x04, 0x03, 0x03, 0x04, 0x05})
	if err != nil {
		t.Fatal(err)
	}
	if len(octetStrings) != 2 || !bytes.Equal(octetStrings[0], []byte{0x01, 0x02}) ||
		!bytes.Equal(octetStrings[1], []byte{0x03, 0x04, 0x05}) {
		t.Errorf("got %x", octetStrings)
	}
	if list, err := DSSASN1UtilsDEROctetStrings(nil); err != nil || len(list) != 0 {
		t.Errorf("a nil sequence must yield an empty list, got %v %v", list, err)
	}
	if _, err := DSSASN1UtilsDEROctetStrings([]byte{0x30, 0x03, 0x02, 0x01, 0x07}); err == nil {
		t.Error("a sequence holding a non-OCTET STRING must be rejected")
	}
}

// TestDSSASN1UtilsTagAndEncodingChecks covers the small predicates.
func TestDSSASN1UtilsTagAndEncodingChecks(t *testing.T) {
	if !DSSASN1UtilsIsASN1SequenceTag(0x30) {
		t.Error("0x30 is the SEQUENCE tag")
	}
	if DSSASN1UtilsIsASN1SequenceTag(0x31) {
		t.Error("0x31 is the SET tag, not the SEQUENCE tag")
	}
	if DSSASN1UtilsIsAsn1Encoded(nil) || DSSASN1UtilsIsAsn1Encoded([]byte{}) {
		t.Error("empty binaries are not ASN.1 encoded")
	}
	if !DSSASN1UtilsIsAsn1Encoded([]byte{0x05, 0x00}) {
		t.Error("DERNull is ASN.1 encoded")
	}
	if DSSASN1UtilsIsAsn1Encoded([]byte{0x30, 0x0F, 0x01}) {
		t.Error("a truncated object is not ASN.1 encoded")
	}
	if DSSASN1UtilsIsAsn1EncodedSignatureValue([]byte{0x30, 0x03, 0x02, 0x01, 0x07}) {
		t.Error("a one-element SEQUENCE is not an ASN.1 signature value")
	}
	if !DSSASN1UtilsIsAsn1EncodedSignatureValue([]byte{0x30, 0x06, 0x02, 0x01, 0x07, 0x02, 0x01, 0x09}) {
		t.Error("a two-element SEQUENCE is an ASN.1 signature value")
	}
}

// TestDSSASN1UtilsAlgorithmIdentifierFromATSHashIndex covers the ats-hash-index accessor,
// including the upstream case where the first element is a bare OBJECT IDENTIFIER.
func TestDSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(t *testing.T) {
	identifier, err := DSSASN1UtilsAlgorithmIdentifierForDigest(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	empty := []byte{0x30, 0x00}
	// A table of four elements, the first being the AlgorithmIdentifier.
	body := append([]byte{}, identifier.DER()...)
	for index := 0; index < 3; index++ {
		body = append(body, empty...)
	}
	table := append([]byte{0x30, byte(len(body))}, body...)
	got := DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(table)
	if got == nil || !got.Equals(identifier) {
		t.Errorf("expected the SHA-256 AlgorithmIdentifier, got %+v", got)
	}

	// The same table with a bare OID in the first slot.
	oid := []byte{0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01}
	body = append([]byte{}, oid...)
	for index := 0; index < 3; index++ {
		body = append(body, empty...)
	}
	table = append([]byte{0x30, byte(len(body))}, body...)
	got = DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(table)
	if got == nil || got.Parameters != nil || !got.Algorithm.Equal(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}) {
		t.Errorf("expected a parameterless SHA-256 identifier, got %+v", got)
	}

	// A table of three elements carries no algorithm.
	body = nil
	for index := 0; index < 3; index++ {
		body = append(body, empty...)
	}
	table = append([]byte{0x30, byte(len(body))}, body...)
	if got := DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(table); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
	if got := DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(nil); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

// TestDSSASN1UtilsToSignerIdentifier checks the issuer/serial/ski SignerIdentifier builder.
func TestDSSASN1UtilsToSignerIdentifier(t *testing.T) {
	token := dssASN1UtilsTestCertificate(t, "subject.der")
	ski, err := DSSASN1UtilsComputeSkiFromCert(token)
	if err != nil {
		t.Fatal(err)
	}
	signerIdentifier := DSSASN1UtilsToSignerIdentifier(token.Issuer().Principal(), token.SerialNumber(), ski)
	if signerIdentifier == nil {
		t.Fatal("the SignerIdentifier must not be nil")
	}
}

// TestDSSASN1UtilsDirectoryStringValue covers the DirectoryString reader.
func TestDSSASN1UtilsDirectoryStringValue(t *testing.T) {
	// 0C 05 "Hello" : a UTF8String.
	if got := DSSASN1UtilsDirectoryStringValue([]byte{0x0C, 0x05, 'H', 'e', 'l', 'l', 'o'}); got != "Hello" {
		t.Errorf("got %q, want Hello", got)
	}
	// 13 03 "abc" : a PrintableString.
	if got := DSSASN1UtilsDirectoryStringValue([]byte{0x13, 0x03, 'a', 'b', 'c'}); got != "abc" {
		t.Errorf("got %q, want abc", got)
	}
	// An INTEGER is not a DirectoryString.
	if got := DSSASN1UtilsDirectoryStringValue([]byte{0x02, 0x01, 0x07}); got != "" {
		t.Errorf("got %q, want the empty string", got)
	}
	if got := DSSASN1UtilsDirectoryStringValue(nil); got != "" {
		t.Errorf("got %q, want the empty string", got)
	}
}
