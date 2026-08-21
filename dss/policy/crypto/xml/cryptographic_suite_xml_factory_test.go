// Accessor behavior tests vs Java-dumped answers, per PORTING.md/S8A_BRIEF.md's
// crypto-suite KAT requirement.
//
// testdata/oracle/dss-crypto-suite.accessors.txt is a tab-separated
// "key<TAB>value" dump produced by running upstream's real
// CryptographicSuiteXmlFactory/CryptographicSuiteCatalogue over this
// package's own copy of upstream's src/main/resources/suite/dss-crypto-suite.xml
// (byte-identical to dss/policy/crypto/json's dss-crypto-suite.json content,
// so the golden values are identical too), calling this package's own
// CryptoXmlOracle.java (not checked into this repository - see
// policy/jaxb/doc.go for the sibling PolicyOracle.java precedent). See
// dss/policy/crypto/json/cryptographic_suite_json_factory_test.go's header
// for the Set<CryptographicSuiteEvaluation>-has-no-order rationale this
// file's canonicalization follows too.
package cryptoxml

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	modelpolicy "github.com/utain/esig/dss/model/policy"
)

func readGolden(t *testing.T, path string) map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	golden := make(map[string]string)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		idx := strings.IndexByte(line, '\t')
		if idx < 0 {
			t.Fatalf("malformed golden line (no tab): %q", line)
		}
		golden[line[:idx]] = line[idx+1:]
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan %s: %v", path, err)
	}
	return golden
}

func TestCryptographicSuiteXmlCatalogueMatchesOracle(t *testing.T) {
	data, err := os.ReadFile(corpustest.Path(t, "dss-crypto-suite.xml"))
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	factory := NewCryptographicSuiteXmlFactory()
	doc := model.NewInMemoryDocument(data)
	catalogue := factory.LoadCryptographicSuite(doc)

	golden := readGolden(t, corpustest.Path(t, "oracle/dss-crypto-suite.accessors.txt"))
	seen := make(map[string]bool, len(golden))

	checkSuite(t, golden, seen, "CryptographicSuite", catalogue.CryptographicSuite())
	checkSuite(t, golden, seen, "SignatureCertificatesCryptographicSuite", catalogue.SignatureCertificatesCryptographicSuite())
	checkSuite(t, golden, seen, "RevocationCryptographicSuite", catalogue.RevocationCryptographicSuite())
	checkSuite(t, golden, seen, "TimestampCryptographicSuite", catalogue.TimestampCryptographicSuite())

	for key := range golden {
		if !seen[key] {
			t.Errorf("golden key %q was never checked by this test", key)
		}
	}
}

func checkSuite(t *testing.T, golden map[string]string, seen map[string]bool, prefix string, suite modelpolicy.CryptographicSuite) {
	t.Helper()
	check(t, golden, seen, prefix+".Level", string(suite.Level()))
	check(t, golden, seen, prefix+".PolicyName", suite.PolicyName())

	updateDateStr := "<null>"
	if d := suite.CryptographicSuiteUpdateDate(); d != nil {
		updateDateStr = d.UTC().Format("2006-01-02")
	}
	check(t, golden, seen, prefix+".CryptographicSuiteUpdateDate", updateDateStr)

	check(t, golden, seen, prefix+".AcceptableDigestAlgorithms.keys", sortedDigestAlgorithmNames(suite.AcceptableDigestAlgorithms()))
	check(t, golden, seen, prefix+".AcceptableSignatureAlgorithms.keys", sortedSignatureAlgorithmNames(suite.AcceptableSignatureAlgorithms()))

	sigAlgos := suite.AcceptableSignatureAlgorithms()
	for _, sa := range []enumerations.SignatureAlgorithm{
		enumerations.SignatureAlgorithm_RSA_SHA256,
		enumerations.SignatureAlgorithm_ECDSA_SHA256,
	} {
		check(t, golden, seen, prefix+"."+string(sa)+".evaluations", evaluationsCanonical(sigAlgos[sa]))
	}
}

func check(t *testing.T, golden map[string]string, seen map[string]bool, key, got string) {
	t.Helper()
	want, ok := golden[key]
	if !ok {
		t.Fatalf("golden file missing key %q", key)
	}
	seen[key] = true
	if got != want {
		t.Errorf("%s = %q, want %q", key, got, want)
	}
}

func sortedDigestAlgorithmNames(m map[enumerations.DigestAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation) string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, string(k))
	}
	sort.Strings(names)
	return strings.Join(names, ";")
}

func sortedSignatureAlgorithmNames(m map[enumerations.SignatureAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation) string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, string(k))
	}
	sort.Strings(names)
	return strings.Join(names, ";")
}

func evaluationsCanonical(evaluations []*modelpolicy.CryptographicSuiteEvaluation) string {
	rendered := make([]string, 0, len(evaluations))
	for _, e := range evaluations {
		var sb strings.Builder
		sb.WriteString("[minSizes=")
		for _, p := range e.ParameterList() {
			sb.WriteString(p.Name())
			sb.WriteByte('=')
			if p.Min() != nil {
				sb.WriteString(strconv.Itoa(*p.Min()))
			} else {
				sb.WriteString("null")
			}
			sb.WriteByte(',')
		}
		sb.WriteString(";validityEnd=")
		if end := e.ValidityEnd(); end != nil {
			sb.WriteString(end.UTC().Format("2006-01-02"))
		} else {
			sb.WriteString("none")
		}
		sb.WriteByte(']')
		rendered = append(rendered, sb.String())
	}
	sort.Strings(rendered)
	return strings.Join(rendered, "")
}

func TestCryptographicSuiteXmlFactoryIsSupported(t *testing.T) {
	data, err := os.ReadFile(corpustest.Path(t, "dss-crypto-suite.xml"))
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	factory := NewCryptographicSuiteXmlFactory()

	if !factory.IsSupported(model.NewInMemoryDocument(data)) {
		t.Error("IsSupported(dss-crypto-suite.xml) = false, want true")
	}
	if factory.IsSupported(model.NewInMemoryDocument([]byte("not xml at all"))) {
		t.Error("IsSupported(garbage) = true, want false")
	}
}

func TestCryptographicSuiteXmlFactoryLoadDefaultCryptographicSuite(t *testing.T) {
	factory := NewCryptographicSuiteXmlFactory()
	catalogue := factory.LoadDefaultCryptographicSuite()
	if got, want := catalogue.CryptographicSuite().PolicyName(), "ETSI ESI Cryptographic Suites"; got != want {
		t.Errorf("PolicyName() = %q, want %q", got, want)
	}
}

func TestCryptographicSuiteXmlFacadeGetCryptographicSuiteNilReader(t *testing.T) {
	facade := NewCryptographicSuiteXmlFacade()
	if _, err := facade.GetCryptographicSuite(nil); err == nil {
		t.Error("GetCryptographicSuite(nil) = nil error, want an error")
	}
}
