package diagnostic

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// The certificate-builder KAT: testdata/oracle/<name>.xml is the marshalled
// XmlDiagnosticData upstream's CertificateDiagnosticDataBuilder produces for the single
// certificate spi/testdata/certificate_extensions/<name>.der, at a fixed validation date -
// see testdata/gen/CertificateDiagnosticDataOracle.java. This test builds the same document
// in Go and byte-compares the marshalled XML, which is what pins the element order of every
// certificate / certificate-extension / QcStatements builder in this package.
//
// One document per certificate keeps the comparison independent of the Java HashSet
// iteration order over the used-certificate set, which Go cannot reproduce.

// katValidationDate is the validation time the oracle ran with, 2024-01-01T00:00:00Z.
var katValidationDate = time.Unix(1704067200, 0).UTC()

// katOrderDeviation lists the certificates whose document differs from the Java one only
// in the order of the <OtherExtension> children. That order is the documented, deliberate
// deviation of spi.CertificateExtensionsUtilsCertificateExtensions: Java walks the JCA's
// getCriticalExtensionOIDs()/getNonCriticalExtensionOIDs() HashSets, whose hash order Go
// cannot reproduce, so the port walks the certificate's extensions in DER order instead.
// These documents are therefore compared as a multiset of lines, which still pins every
// element, attribute and value - only the sequence of the sibling extensions is relaxed.
var katOrderDeviation = map[string]bool{
	"cert_01": true,
	"cert_07": true,
}

// katKnownDeviation lists the certificates the port cannot yet reproduce at all, with the
// reason. Each is a pre-existing frozen-package gap, not a phase 8c regression; the list is
// asserted to be exactly this set so a new failure cannot hide inside it.
var katKnownDeviation = map[string]string{
	// RIPEMD160withRSA: crypto/x509 has no signature-algorithm binding for it, so
	// CertificateToken.IsSelfSigned() (which verifies the signature) answers false where
	// the JCA answers true. That suppresses IssuerEntityKey, SignatureIntact/Valid and
	// KeyLengthUsedToSignThisToken. The same certificate also exposes a TeletexString
	// distinguished-name decoding difference in model's DN rendering.
	"cert_05": "RIPEMD160withRSA unsupported by crypto/x509 + TeletexString DN rendering (model, frozen)",
	// crypto/x509 rejects the DER outright where the JCA parses it.
	"cert_16": "crypto/x509 rejects the DER",
}

func TestCertificateDiagnosticDataBuilderKAT(t *testing.T) {
	oracleDir := corpustest.Path(t, "oracle")
	entries, err := os.ReadDir(oracleDir)
	if err != nil {
		t.Fatalf("read oracle dir: %v", err)
	}
	derDir := filepath.Join("..", "..", "..", "spi", "testdata", "certificate_extensions")

	compared := 0
	orderCompared := 0
	seenDeviation := map[string]bool{}
	for _, entry := range entries {
		name := entry.Name()
		if filepath.Ext(name) != ".xml" {
			continue
		}
		base := strings.TrimSuffix(name, ".xml")
		if reason, known := katKnownDeviation[base]; known {
			seenDeviation[base] = true
			t.Logf("%s: known deviation, not compared (%s)", base, reason)
			continue
		}
		t.Run(base, func(t *testing.T) {
			expected, err := os.ReadFile(filepath.Join(oracleDir, name))
			if err != nil {
				t.Fatalf("read oracle: %v", err)
			}
			der, err := os.ReadFile(filepath.Join(derDir, base+".der"))
			if err != nil {
				t.Fatalf("read der: %v", err)
			}
			cert, err := x509.ParseCertificate(der)
			if err != nil {
				t.Fatalf("crypto/x509 rejects %s (add it to katKnownDeviation if that is expected): %v",
					base, err)
			}
			token, err := model.NewCertificateToken(cert)
			if err != nil {
				t.Fatalf("build token: %v", err)
			}

			diagnosticData := NewCertificateDiagnosticDataBuilder().
				UsedCertificates([]*model.CertificateToken{token}).
				UsedRevocations([]validation.AnyRevocationToken{}).
				DefaultDigestAlgorithm(enumerations.DigestAlgorithm_SHA256).
				TokenExtractionStrategy(enumerations.TokenExtractionStrategy_NONE).
				ValidationDate(katValidationDate).
				Build()

			actual, err := jaxb.Marshal(diagnosticData)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(actual) == string(expected) {
				compared++
				if katOrderDeviation[base] {
					t.Errorf("%s: listed in katOrderDeviation but now byte-identical - drop the entry",
						base)
				}
				return
			}
			if katOrderDeviation[base] {
				seenDeviation[base] = true
				if diff := lineMultisetDifference(string(expected), string(actual)); diff != "" {
					t.Errorf("%s: differs from the Java oracle in more than sibling order\n%s",
						base, diff)
					return
				}
				orderCompared++
				return
			}
			t.Errorf("%s: marshalled XmlDiagnosticData differs from the Java oracle\n%s",
				base, firstDifference(string(expected), string(actual)))
		})
	}
	if compared == 0 {
		t.Fatal("no certificate compared against the oracle")
	}
	for base := range katKnownDeviation {
		if !seenDeviation[base] {
			t.Errorf("katKnownDeviation lists %q, but the oracle has no such document", base)
		}
	}
	for base := range katOrderDeviation {
		if !seenDeviation[base] {
			t.Errorf("katOrderDeviation lists %q, but it did not deviate", base)
		}
	}
	t.Logf("certificate diagnostic-data documents byte-identical to Java: %d "+
		"(+%d equal up to the documented extension order, %d known deviations)",
		compared, orderCompared, len(katKnownDeviation))
}

// lineMultisetDifference reports the lines that differ between two documents once both
// are sorted, i.e. it ignores the order of sibling elements but nothing else. It returns
// the empty string when the two documents carry exactly the same lines.
func lineMultisetDifference(expected, actual string) string {
	expectedLines := append([]string(nil), strings.Split(expected, "\n")...)
	actualLines := append([]string(nil), strings.Split(actual, "\n")...)
	sort.Strings(expectedLines)
	sort.Strings(actualLines)
	if len(expectedLines) != len(actualLines) {
		return "documents carry a different number of lines"
	}
	var b strings.Builder
	for i := range expectedLines {
		if expectedLines[i] != actualLines[i] {
			b.WriteString("\n  java: ")
			b.WriteString(expectedLines[i])
			b.WriteString("\n  go:   ")
			b.WriteString(actualLines[i])
		}
	}
	return b.String()
}

// firstDifference renders the first differing line of two documents with a little context.
func firstDifference(expected, actual string) string {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")
	for i := 0; i < len(expectedLines) || i < len(actualLines); i++ {
		var e, a string
		if i < len(expectedLines) {
			e = expectedLines[i]
		}
		if i < len(actualLines) {
			a = actualLines[i]
		}
		if e != a {
			var b strings.Builder
			b.WriteString("first difference at line ")
			b.WriteString(itoa(i + 1))
			b.WriteString("\n  java: ")
			b.WriteString(e)
			b.WriteString("\n  go:   ")
			b.WriteString(a)
			return b.String()
		}
	}
	return "documents differ only in length"
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
