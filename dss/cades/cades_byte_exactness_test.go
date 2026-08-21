// CAdES-B byte-exactness against upstream: the BUILD chunk's core contract.
//
// testdata/bytecmp holds two generators - main.go (this port's CAdESService) and
// ByteExactnessFixtures.java (upstream DSS 6.5.RC1's CAdESService) - fed byte-identical inputs:
// the same PKCS#12 key store, signing certificate and chain, payload, signing time and digest
// algorithm. RSA PKCS#1 v1.5 is deterministic, so the two CMS documents must agree byte for byte
// wherever the two implementations agree at all, and every difference is a real finding rather
// than crypto nondeterminism.
//
// Two assertions come out of that, and they are deliberately different in strength:
//
//   - The DER SET OF signed attributes (getDataToSign's output, the bytes the signature is
//     computed over) must be byte-identical for both fixtures. This is the CAdESLevelBaselineB
//     contract, checked here end-to-end through the real service rather than attribute by
//     attribute the way cades_level_baseline_b_kat_test.go does it.
//   - The complete CMS must be byte-identical too, except that the SignedData certificates SET
//     may be ordered differently. That single exception is internal/cmscore's one documented
//     deviation (DER-sorted where BouncyCastle keeps insertion order; both are valid DER for a
//     SET OF, and the SignerInfo is unaffected). The test does not merely tolerate it: it pins
//     it, by asserting the two certificate SETs hold exactly the same certificates and that
//     normalising only that one SET makes the documents byte-identical - so any *other*
//     difference, anywhere in the document, fails.
//
// The single-certificate fixture additionally shows the deviation is confined to ordering: with
// nothing to reorder, the raw bytes already match with no normalisation at all.
//
// Skipped under -short and when Java/Maven or a built upstream DSS checkout are unavailable,
// exactly like cades_downstream_cross_validation_test.go (whose detection and classpath helpers
// this reuses).
package cades

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// byteExactnessFixture is one (key store, password, signing year) triple both generators run on.
type byteExactnessFixture struct {
	name string
	// keyStore is a file name inside testdata/bytecmp.
	keyStore string
	password string
	// signingYear must fall inside the signing certificate's validity window, since
	// CAdESService refuses to sign with a certificate that is not yet valid or has expired.
	signingYear string
	// certificateCount is how many certificates the key entry's chain carries, and therefore how
	// many the SignedData certificates SET holds.
	certificateCount int
}

// byteExactnessFixtures covers both the trivial case (one certificate, nothing to reorder) and
// the case that actually exercises the certificates-SET deviation (a three-certificate chain).
var byteExactnessFixtures = []byteExactnessFixture{
	{"single-certificate", "signer_rsa.p12", "testpassword", "2027", 1},
	{"three-certificate-chain", "good-user.p12", "ks-password", "2018", 3},
}

func TestCAdESLevelBaselineBBytesMatchUpstream(t *testing.T) {
	if testing.Short() {
		t.Skip("byte-exactness against upstream DSS shells out to go run/mvn/javac/java; skipped under -short")
	}

	upstreamHome := os.Getenv("DSS_UPSTREAM_HOME")
	if skipReason := detectJavaAndUpstreamDSS(upstreamHome); skipReason != "" {
		t.Skip(skipReason)
	}

	classpath, err := buildUpstreamDSSClasspath(upstreamHome)
	if err != nil {
		t.Fatalf("building upstream DSS classpath: %v", err)
	}
	// dss-token (Pkcs12SignatureToken, the signing side) is not a dependency of dss-validation,
	// so buildUpstreamDSSClasspath - which resolves that module's runtime classpath - does not
	// bring it in. The cross-validation harness never needed it because it only validates.
	tokenClasses := filepath.Join(upstreamHome, "dss-token", "target", "classes")
	if info, err := os.Stat(tokenClasses); err != nil || !info.IsDir() {
		t.Skipf("upstream dss-token is not built at %s (run `mvn -o -pl dss-token -am install -DskipTests`): %v",
			tokenClasses, err)
	}
	classpath = tokenClasses + string(os.PathListSeparator) + classpath

	javaClasses := t.TempDir()
	javac := exec.Command("javac", "-nowarn", "-cp", classpath, "-d", javaClasses, "ByteExactnessFixtures.java")
	javac.Dir = filepath.Join("testdata", "bytecmp")
	if output, err := javac.CombinedOutput(); err != nil {
		t.Fatalf("javac ByteExactnessFixtures.java failed: %v\n%s", err, output)
	}

	for _, fixture := range byteExactnessFixtures {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			outDir := t.TempDir()
			keyStore, err := filepath.Abs(filepath.Join("testdata", "bytecmp", fixture.keyStore))
			if err != nil {
				t.Fatal(err)
			}

			// A failure in either generator is a hard failure: both are running code this test
			// exists to compare, not optional tooling.
			goRun := exec.Command("go", "run", ".", keyStore, fixture.password, fixture.signingYear, outDir, "go")
			goRun.Dir = filepath.Join("testdata", "bytecmp")
			if output, err := goRun.CombinedOutput(); err != nil {
				t.Fatalf("the Go generator failed: %v\n%s", err, output)
			}
			javaRun := exec.Command("java", "-cp", classpath+string(os.PathListSeparator)+javaClasses,
				"ByteExactnessFixtures", keyStore, fixture.password, fixture.signingYear, outDir, "java")
			if output, err := javaRun.CombinedOutput(); err != nil {
				t.Fatalf("the Java generator failed: %v\n%s", err, stripJavaToolOptionsNoise(string(output)))
			}

			goDataToSign := readByteExactnessFile(t, outDir, "go-datatosign.bin")
			javaDataToSign := readByteExactnessFile(t, outDir, "java-datatosign.bin")
			if !bytes.Equal(goDataToSign, javaDataToSign) {
				t.Errorf("the signed attributes differ from upstream's:\n go   (%d bytes) %x\n java (%d bytes) %x",
					len(goDataToSign), goDataToSign, len(javaDataToSign), javaDataToSign)
			}

			goCMS := readByteExactnessFile(t, outDir, "go.p7m")
			javaCMS := readByteExactnessFile(t, outDir, "java.p7m")

			goCertificates, goNormalised := normaliseCertificatesSet(t, goCMS)
			javaCertificates, javaNormalised := normaliseCertificatesSet(t, javaCMS)

			if len(goCertificates) != fixture.certificateCount {
				t.Errorf("the Go CMS carries %d certificates, want %d", len(goCertificates), fixture.certificateCount)
			}
			if len(javaCertificates) != fixture.certificateCount {
				t.Errorf("upstream's CMS carries %d certificates, want %d", len(javaCertificates), fixture.certificateCount)
			}

			// Same certificates, whatever the order: the deviation is ordering only, never a
			// certificate this port adds or drops.
			if !equalCertificateSets(goCertificates, javaCertificates) {
				t.Errorf("the certificates SET holds different certificates:\n go   %x\n java %x",
					goCertificates, javaCertificates)
			}

			// Everything else, byte for byte. Only the certificates SET was normalised, so a
			// difference anywhere else - signed attributes, signature value, SignerInfo,
			// encapContentInfo, version numbers - fails here.
			if !bytes.Equal(goNormalised, javaNormalised) {
				t.Errorf("the CMS differs from upstream's outside the certificates SET"+
					" (%d vs %d bytes)", len(goNormalised), len(javaNormalised))
			}

			// With a single certificate there is nothing to reorder, so the deviation must not
			// show up at all: the raw bytes have to match with no normalisation.
			if fixture.certificateCount == 1 && !bytes.Equal(goCMS, javaCMS) {
				t.Error("a single-certificate CMS must already be byte-identical to upstream's")
			}
		})
	}
}

// readByteExactnessFile reads one generator output.
func readByteExactnessFile(t *testing.T, dir, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return data
}

// normaliseCertificatesSet returns the DER of every certificate in a CMS document's SignedData
// certificates SET, plus a copy of the document in which that SET's members have been sorted by
// their encoding. Sorting only that one SET is what lets the rest of the document be compared
// byte for byte while still allowing internal/cmscore's documented ordering deviation.
func normaliseCertificatesSet(t *testing.T, document []byte) ([][]byte, []byte) {
	t.Helper()

	contentInfo, _, err := asn1ber.Parse(document)
	if err != nil {
		t.Fatalf("parsing the CMS: %v", err)
	}
	// ContentInfo ::= SEQUENCE { contentType OID, content [0] EXPLICIT ANY }
	content := contentInfo.Children()[1]
	// SignedData ::= SEQUENCE { version, digestAlgorithms, encapContentInfo,
	//                           certificates [0] IMPLICIT OPTIONAL, ... }
	signedData := content.Children()[0]

	var certificatesSet *asn1ber.Element
	for _, child := range signedData.Children() {
		if child.Class() == asn1ber.ClassContextSpecific && child.TagNumber() == 0 {
			certificatesSet = child
			break
		}
	}
	if certificatesSet == nil {
		t.Fatal("the SignedData carries no certificates SET")
	}

	certificates := make([][]byte, 0, len(certificatesSet.Children()))
	for _, certificate := range certificatesSet.Children() {
		certificates = append(certificates, certificate.Encoded())
	}

	sorted := make([][]byte, len(certificates))
	copy(sorted, certificates)
	sort.Slice(sorted, func(a, b int) bool { return bytes.Compare(sorted[a], sorted[b]) < 0 })

	var body []byte
	for _, certificate := range sorted {
		body = append(body, certificate...)
	}
	rebuilt := asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|0, body)

	normalised := bytes.Replace(document, certificatesSet.Encoded(), rebuilt, 1)
	if len(normalised) != len(document) {
		t.Fatalf("normalising the certificates SET changed the document length (%d -> %d)",
			len(document), len(normalised))
	}
	return certificates, normalised
}

// equalCertificateSets reports whether two certificate lists hold the same certificates,
// regardless of order.
func equalCertificateSets(first, second [][]byte) bool {
	if len(first) != len(second) {
		return false
	}
	sortedFirst := make([][]byte, len(first))
	copy(sortedFirst, first)
	sortedSecond := make([][]byte, len(second))
	copy(sortedSecond, second)
	sort.Slice(sortedFirst, func(a, b int) bool { return bytes.Compare(sortedFirst[a], sortedFirst[b]) < 0 })
	sort.Slice(sortedSecond, func(a, b int) bool { return bytes.Compare(sortedSecond[a], sortedSecond[b]) < 0 })
	for i := range sortedFirst {
		if !bytes.Equal(sortedFirst[i], sortedSecond[i]) {
			return false
		}
	}
	return true
}
