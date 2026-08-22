// Checks signatureAlgorithmIdentifierHex (custom_content_signer.go) against
// testdata/signature-algorithm-identifiers.txt, a real BouncyCastle 1.84
// DefaultSignatureAlgorithmIdentifierFinder run (see testdata/gen/GenSignatureAlgorithmIdentifiers.java):
// every successfully resolved JCE name must produce byte-identical DER, and every JCE name
// BouncyCastle itself rejects must be absent from the Go table too.
//
// The oracle file is the generator's output verbatim, comment lines included, and carries three
// sections. Earlier revisions of this test could only parse the first one, so the file had been
// hand-trimmed down to it - which quietly dropped the ground truth for the other two (the
// CMSAlgorithmProtection tag shape and DERTaggedObject's default explicitness), leaving both
// claims asserted nowhere despite being cited in the source comments of
// cms_signed_attribute_table_generator.go and signer_attribute_v2.go. All three are checked here
// now, and the file is byte-for-byte reproducible by re-running the generator.
package cms

import (
	"bufio"
	"encoding/asn1"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// signatureAlgorithmOracle parses testdata/signature-algorithm-identifiers.txt into its sections,
// keyed by the "# <title>" comment that introduces each one. Entries keep the file's order.
func signatureAlgorithmOracle(t *testing.T) map[string]map[string]string {
	t.Helper()
	file, err := os.Open("testdata/signature-algorithm-identifiers.txt")
	if err != nil {
		t.Fatalf("open testdata: %s", err)
	}
	defer file.Close()

	sections := make(map[string]map[string]string)
	section := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			section = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			sections[section] = make(map[string]string)
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			t.Fatalf("malformed line: %q", line)
		}
		if section == "" {
			t.Fatalf("entry %q appears before any section header", line)
		}
		sections[section][key] = value
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan testdata: %s", err)
	}
	return sections
}

// oracleSection returns one section, failing when the file no longer carries it - so a
// regenerated oracle that drops a section is caught rather than silently reducing coverage.
func oracleSection(t *testing.T, sections map[string]map[string]string, name string) map[string]string {
	t.Helper()
	section, found := sections[name]
	if !found {
		t.Fatalf("the oracle carries no %q section", name)
	}
	if len(section) == 0 {
		t.Fatalf("the oracle's %q section is empty", name)
	}
	return section
}

func TestSignatureAlgorithmIdentifiersMatchBouncyCastle(t *testing.T) {
	entries := oracleSection(t, signatureAlgorithmOracle(t),
		"signatureAlgorithmIdentifierHex (DefaultSignatureAlgorithmIdentifierFinder#find)")

	for jceID, want := range entries {
		got, ok := signatureAlgorithmIdentifierHex[jceID]
		if strings.HasPrefix(want, "ERROR:") {
			if ok {
				t.Errorf("%s: BouncyCastle rejects it (%s), but the Go table has an entry: %s", jceID, want, got)
			}
			continue
		}
		if !ok {
			t.Errorf("%s: missing from signatureAlgorithmIdentifierHex, want %s", jceID, want)
			continue
		}
		if got != want {
			t.Errorf("%s: got %s, want %s", jceID, got, want)
		}
	}

	// Every table entry must in turn be exercised by the KAT file, so the table cannot grow
	// stale entries the KAT no longer covers.
	for jceID := range signatureAlgorithmIdentifierHex {
		if _, covered := entries[jceID]; !covered {
			t.Errorf("%s: in the Go table but not exercised by the KAT file", jceID)
		}
	}
}

// TestCMSAlgorithmProtectionMatchesBouncyCastle checks cmsAlgorithmProtectionDER against the
// encoding BouncyCastle's own org.bouncycastle.asn1.cms.CMSAlgorithmProtection produces. The
// point is the [1] signatureAlgorithm tag: the ASN.1 module alone does not say whether it is
// IMPLICIT or EXPLICIT, and CMS modules mix the two by field, so the answer has to come from a
// real BouncyCastle run rather than from reading the grammar.
func TestCMSAlgorithmProtectionMatchesBouncyCastle(t *testing.T) {
	entries := oracleSection(t, signatureAlgorithmOracle(t), "CMSAlgorithmProtection tag shape (SIGNATURE alternative)")

	sha256 := asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	sha256WithRSA := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}

	t.Run("no-params", func(t *testing.T) {
		got := cmsAlgorithmProtectionDER(
			&asn1ber.AlgorithmIdentifier{Algorithm: sha256},
			&asn1ber.AlgorithmIdentifier{Algorithm: sha256WithRSA})
		hexEqual(t, got, entries["no-params"])
	})

	t.Run("digest-null-param", func(t *testing.T) {
		got := cmsAlgorithmProtectionDER(
			&asn1ber.AlgorithmIdentifier{Algorithm: sha256, Parameters: []byte{0x05, 0x00}},
			&asn1ber.AlgorithmIdentifier{Algorithm: sha256WithRSA})
		hexEqual(t, got, entries["digest-null-param"])
	})
}

// TestDERTaggedObjectDefaultExplicitness pins the fact signer_attribute_v2.go's DER() relies on:
// BouncyCastle's two-argument `new DERTaggedObject(tagNo, obj)` tags EXPLICITLY, so a tagged
// field keeps the inner value's own tag rather than replacing it. The oracle records all three
// spellings; two-arg must equal explicit and differ from implicit.
func TestDERTaggedObjectDefaultExplicitness(t *testing.T) {
	entries := oracleSection(t, signatureAlgorithmOracle(t),
		"DERTaggedObject(int, ASN1Encodable) two-argument default explicitness")

	if entries["two-arg"] != entries["explicit"] {
		t.Errorf("two-arg = %s, explicit = %s: the two-argument form is not explicit",
			entries["two-arg"], entries["explicit"])
	}
	if entries["two-arg"] == entries["implicit"] {
		t.Errorf("two-arg = %s equals implicit: the oracle cannot tell the two apart", entries["two-arg"])
	}

	// The Go port's own explicit tagging must reproduce it: wrapping the same inner SEQUENCE the
	// generator used ([0] over SEQUENCE { OID 1.2.3.4 }) has to give the oracle's bytes.
	inner := asn1ber.WriteSequence(asn1ber.EncodeOID(asn1.ObjectIdentifier{1, 2, 3, 4}))
	got := asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|0, inner) //nolint:staticcheck // mirrors upstream SignerAttributeV2#toASN1Primitive: `new DERTaggedObject(0, ...)` - the `0` is the ASN.1 context tag number the oracle entry was generated with.
	hexEqual(t, got, entries["two-arg"])
}

// TestFindSignatureAlgorithmIdentifierEveryDSSSignatureAlgorithm exercises
// findSignatureAlgorithmIdentifier for every enumerations.SignatureAlgorithm this port produces
// a JCEID() for, checking that families outside the KAT's scope (HMAC, GOST, SM2, ...) simply
// come back as errors rather than panicking or producing garbage.
func TestFindSignatureAlgorithmIdentifierNeverPanics(t *testing.T) {
	for _, signatureAlgorithm := range allSignatureAlgorithmsForTest() {
		jceID := signatureAlgorithm.JCEID()
		if jceID == "" {
			continue
		}
		identifier, err := findSignatureAlgorithmIdentifier(jceID)
		if err == nil && identifier == nil {
			t.Errorf("%s (%s): nil identifier and nil error", signatureAlgorithm, jceID)
		}
	}
}

func allSignatureAlgorithmsForTest() []enumerations.SignatureAlgorithm {
	return []enumerations.SignatureAlgorithm{
		enumerations.SignatureAlgorithmRSASHA1,
		enumerations.SignatureAlgorithmRSASHA256,
		enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1,
		enumerations.SignatureAlgorithmECDSASHA256,
		enumerations.SignatureAlgorithmPlainECDSASHA256,
		enumerations.SignatureAlgorithmED25519,
		enumerations.SignatureAlgorithmED448,
		enumerations.SignatureAlgorithmDSASHA256,
		enumerations.SignatureAlgorithmHMACSHA256,
	}
}

func hexEqual(t *testing.T, got []byte, wantHex string) {
	t.Helper()
	want, err := hex.DecodeString(wantHex)
	if err != nil {
		t.Fatalf("bad hex fixture: %s", err)
	}
	if string(got) != string(want) {
		t.Fatalf("got %x, want %s", got, wantHex)
	}
}
