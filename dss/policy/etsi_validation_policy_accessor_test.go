// Accessor behavior tests vs Java-dumped answers, per PORTING.md/S8A_BRIEF.md's
// "POLICY: ... ValidationPolicy accessor behavior tests vs Java-dumped
// answers" requirement.
//
// testdata/oracle/constraint.accessors.txt is a tab-separated "key<TAB>value"
// dump produced by running upstream's real EtsiValidationPolicy (loaded
// through ValidationPolicyFacade) over jaxb/testdata's own copy of upstream's
// src/main/resources/policy/constraint.xml, and calling this package's own
// PolicyAccessorOracle.java (not checked into this repository - see
// jaxb/doc.go for the sibling PolicyOracle.java precedent) against a broad
// sample of ValidationPolicy accessors across every enumerations.Context and
// enumerations.SubContext value. Each test below replays the identical call
// sequence against the Go EtsiValidationPolicy and asserts byte-identical
// answers.
//
// java.util.Set<CryptographicSuiteEvaluation> has no defined iteration
// order; both the oracle and this file canonicalize each evaluation set by
// rendering every element to a string and sorting those strings before
// comparing, so the golden file does not depend on java.util.HashSet's
// bucket order (see the oracle's cryptoSuiteImpl for the matching logic).
// The same applies to the AcceptableDigestAlgorithms/AcceptableSignatureAlgorithms
// map key sets: both sides render an alphabetically-sorted key list rather
// than relying on any particular map/TreeMap iteration order.
package policy

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	modelpolicy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/policy/jaxb"
)

// loadOraclePolicy unmarshals jaxb/testdata/policy/<name>.xml and wraps it as
// an EtsiValidationPolicy, mirroring ValidationPolicyFacade.getValidationPolicy.
func loadOraclePolicy(t *testing.T, name string) *EtsiValidationPolicy {
	t.Helper()
	data, err := os.ReadFile(corpustest.RootPath(t, "policy/jaxb/testdata/policy/"+name+".xml"))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	cp, err := jaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return NewEtsiValidationPolicy(cp)
}

// readGolden parses a "key<TAB>value" oracle dump into a key->value map,
// undoing the escaping PolicyAccessorOracle#line applies to embedded
// backslashes/newlines/tabs.
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
		key := line[:idx]
		if _, dup := golden[key]; dup {
			t.Fatalf("duplicate golden key %q", key)
		}
		golden[key] = unescapeOracleValue(line[idx+1:])
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan %s: %v", path, err)
	}
	return golden
}

func unescapeOracleValue(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				out.WriteByte('\n')
				i++
				continue
			case 't':
				out.WriteByte('\t')
				i++
				continue
			case '\\':
				out.WriteByte('\\')
				i++
				continue
			}
		}
		out.WriteByte(s[i])
	}
	return out.String()
}

// TestEtsiValidationPolicyAccessorsMatchOracle replays the oracle's exact
// accessor sequence against the Go port and asserts each answer matches.
func TestEtsiValidationPolicyAccessorsMatchOracle(t *testing.T) {
	policy := loadOraclePolicy(t, "constraint")
	golden := readGolden(t, corpustest.Path(t, "oracle/constraint.accessors.txt"))
	seen := make(map[string]bool, len(golden))

	check(t, golden, seen, "PolicyName", policy.PolicyName())
	check(t, golden, seen, "PolicyDescription", policy.PolicyDescription())
	check(t, golden, seen, "ValidationModel", string(policy.ValidationModel()))
	check(t, golden, seen, "EIDASConstraintPresent", strconv.FormatBool(policy.EIDASConstraintPresent()))

	for _, ctx := range allContexts {
		label := string(ctx)

		checkPanicking(t, golden, seen, "StructuralValidationConstraint("+label+")", func() string {
			return levelRuleString(policy.StructuralValidationConstraint(ctx))
		})
		checkPanicking(t, golden, seen, "SigningTimeConstraint("+label+")", func() string {
			return levelRuleString(policy.SigningTimeConstraint(ctx))
		})
		checkPanicking(t, golden, seen, "SignaturePolicyConstraint("+label+")", func() string {
			return multiValuesRuleString(policy.SignaturePolicyConstraint(ctx))
		})
		checkPanicking(t, golden, seen, "ContentTypeConstraint("+label+")", func() string {
			return multiValuesRuleString(policy.ContentTypeConstraint(ctx))
		})
		checkPanicking(t, golden, seen, "ProspectiveCertificateChainConstraint("+label+")", func() string {
			return levelRuleString(policy.ProspectiveCertificateChainConstraint(ctx))
		})

		checkCryptoSuitePanicking(t, golden, seen, "SignatureCryptographicConstraint("+label+")", func() modelpolicy.CryptographicSuite {
			return policy.SignatureCryptographicConstraint(ctx)
		})

		for _, sub := range allSubContexts {
			subLabel := string(sub)
			checkPanicking(t, golden, seen, fmt.Sprintf("CertificateCAConstraint(%s,%s)", label, subLabel), func() string {
				return levelRuleString(policy.CertificateCAConstraint(ctx, sub))
			})
			checkPanicking(t, golden, seen, fmt.Sprintf("CertificateNotExpiredConstraint(%s,%s)", label, subLabel), func() string {
				return levelRuleString(policy.CertificateNotExpiredConstraint(ctx, sub))
			})
			checkPanicking(t, golden, seen, fmt.Sprintf("CertificateMinQcEuLimitValueConstraint(%s,%s)", label, subLabel), func() string {
				return numericValueRuleString(policy.CertificateMinQcEuLimitValueConstraint(ctx, sub))
			})
			checkPanicking(t, golden, seen, fmt.Sprintf("CertificateCountryConstraint(%s,%s)", label, subLabel), func() string {
				return multiValuesRuleString(policy.CertificateCountryConstraint(ctx, sub))
			})
			checkCryptoSuitePanicking(t, golden, seen, fmt.Sprintf("CertificateCryptographicConstraint(%s,%s)", label, subLabel), func() modelpolicy.CryptographicSuite {
				return policy.CertificateCryptographicConstraint(ctx, sub)
			})
		}
	}

	checkDuration(t, golden, seen, "TimestampDelayConstraint", policy.TimestampDelayConstraint())
	check(t, golden, seen, "TimestampValidConstraint", levelRuleString(policy.TimestampValidConstraint()))
	check(t, golden, seen, "RevocationDataAvailableConstraint(SIGNATURE,SIGNING_CERT)",
		levelRuleString(policy.RevocationDataAvailableConstraint(enumerations.Context_SIGNATURE, enumerations.SubContext_SIGNING_CERT)))
	checkDuration(t, golden, seen, "RevocationFreshnessConstraint(SIGNATURE,SIGNING_CERT)",
		policy.RevocationFreshnessConstraint(enumerations.Context_SIGNATURE, enumerations.SubContext_SIGNING_CERT))
	check(t, golden, seen, "AcceptedContainerTypesConstraint", multiValuesRuleString(policy.AcceptedContainerTypesConstraint()))
	check(t, golden, seen, "PDFACompliantConstraint", levelRuleString(policy.PDFACompliantConstraint()))
	check(t, golden, seen, "AcceptablePDFAProfilesConstraint", multiValuesRuleString(policy.AcceptablePDFAProfilesConstraint()))

	checkCryptoSuite(t, golden, seen, "EAACryptographicConstraint", policy.EAACryptographicConstraint())
	checkCryptoSuite(t, golden, seen, "EvidenceRecordCryptographicConstraint", policy.EvidenceRecordCryptographicConstraint())

	// Every golden key must have been exercised: guards against the test
	// silently drifting out of sync with the oracle dump (e.g. after a
	// regeneration that adds keys this file doesn't yet replay).
	for key := range golden {
		if !seen[key] {
			t.Errorf("golden key %q was never checked by this test", key)
		}
	}
}

// TestEtsiValidationPolicyNameAndDescriptionMatchOracle spot-checks the
// per-variant PolicyName for the three sample policy documents whose
// accessor answers are otherwise identical to constraint.xml's (the
// eaa/certificate/qwac sample policies share the same constraint body -
// verified by diffing the full oracle dumps against constraint.xml's - so a
// full golden replay would be redundant with
// TestEtsiValidationPolicyAccessorsMatchOracle above).
func TestEtsiValidationPolicyNameAndDescriptionMatchOracle(t *testing.T) {
	cases := []struct{ name, wantName string }{
		{"constraint", "QES AES/QC AES TL based"},
		{"eaa-constraint", "Certificate policy TL based"},
		{"certificate-constraint", "Certificate policy TL based"},
		{"qwac-constraint", "QWAC policy TL based"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			policy := loadOraclePolicy(t, c.name)
			if got := policy.PolicyName(); got != c.wantName {
				t.Errorf("PolicyName() = %q, want %q", got, c.wantName)
			}
		})
	}
}

// allContexts mirrors Java's Context.values() enumeration (declaration)
// order, which PolicyAccessorOracle iterates with `for (Context ctx :
// Context.values())`. Order does not affect this test's assertions (each
// key is looked up by name), but is kept faithful to the oracle for
// readability.
var allContexts = []enumerations.Context{
	enumerations.Context_SIGNATURE,
	enumerations.Context_COUNTER_SIGNATURE,
	enumerations.Context_KEY_BINDING_SIGNATURE,
	enumerations.Context_TIMESTAMP,
	enumerations.Context_EVIDENCE_RECORD,
	enumerations.Context_REVOCATION,
	enumerations.Context_CERTIFICATE,
	enumerations.Context_EAA,
	enumerations.Context_EAA_REVOCATION,
}

// allSubContexts mirrors Java's SubContext.values() enumeration order.
var allSubContexts = []enumerations.SubContext{
	enumerations.SubContext_SIGNING_CERT,
	enumerations.SubContext_CA_CERTIFICATE,
}

func levelRuleString(r modelpolicy.LevelRule) string {
	if r == nil {
		return "<nilrule>"
	}
	return string(r.Level())
}

func multiValuesRuleString(r modelpolicy.MultiValuesRule) string {
	if r == nil {
		return "<nilrule>"
	}
	values := r.Values()
	if len(values) == 0 {
		return "[]"
	}
	return "[" + strings.Join(values, ", ") + "]"
}

func numericValueRuleString(r modelpolicy.NumericValueRule) string {
	if r == nil {
		return "<nilrule>"
	}
	v := r.Value()
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// check asserts golden[key] == got, marking key as seen. Fails the test if
// key is missing from golden entirely (a sign the replay sequence and the
// oracle have drifted apart).
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

func checkDuration(t *testing.T, golden map[string]string, seen map[string]bool, key string, rule modelpolicy.DurationRule) {
	t.Helper()
	got := "<nilrule>"
	if rule != nil {
		got = strconv.FormatInt(rule.Duration(), 10)
	}
	check(t, golden, seen, key, got)
}

// checkPanicking ports the oracle's per-call try/catch: if the golden value
// is the sentinel emitted for a caught UnsupportedOperationException, this
// asserts the Go call panics instead (matching this port's use of panic() to
// carry Java's unchecked exception - see etsi_validation_policy.go's
// basicSignatureConstraintsByContext).
func checkPanicking(t *testing.T, golden map[string]string, seen map[string]bool, key string, fn func() string) {
	t.Helper()
	want, ok := golden[key]
	if !ok {
		t.Fatalf("golden file missing key %q", key)
	}
	seen[key] = true
	if strings.HasPrefix(want, "<exception:") {
		assertPanics(t, key, want, func() { fn() })
		return
	}
	if got := fn(); got != want {
		t.Errorf("%s = %q, want %q", key, got, want)
	}
}

func assertPanics(t *testing.T, key, want string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: want panic (oracle raised %s), got none", key, want)
		}
	}()
	fn()
}

func checkCryptoSuite(t *testing.T, golden map[string]string, seen map[string]bool, prefix string, suite modelpolicy.CryptographicSuite) {
	t.Helper()
	check(t, golden, seen, prefix+".Level", levelRuleString(suite))
	check(t, golden, seen, prefix+".PolicyName", suite.PolicyName())
	check(t, golden, seen, prefix+".AcceptableDigestAlgorithms.keys", sortedDigestAlgorithmNames(suite.AcceptableDigestAlgorithms()))
	check(t, golden, seen, prefix+".AcceptableSignatureAlgorithms.keys", sortedSignatureAlgorithmNames(suite.AcceptableSignatureAlgorithms()))

	sigAlgos := suite.AcceptableSignatureAlgorithms()
	for _, sa := range []enumerations.SignatureAlgorithm{
		enumerations.SignatureAlgorithm_RSA_SHA256,
		enumerations.SignatureAlgorithm_RSA_SHA1,
		enumerations.SignatureAlgorithm_ECDSA_SHA256,
	} {
		check(t, golden, seen, prefix+"."+string(sa)+".evaluations", evaluationsCanonical(sigAlgos[sa]))
	}

	check(t, golden, seen, prefix+".AcceptableDigestAlgorithmsLevel", string(suite.AcceptableDigestAlgorithmsLevel()))
	check(t, golden, seen, prefix+".AcceptableSignatureAlgorithmsLevel", string(suite.AcceptableSignatureAlgorithmsLevel()))
	check(t, golden, seen, prefix+".AcceptableSignatureAlgorithmsMiniKeySizeLevel", string(suite.AcceptableSignatureAlgorithmsMiniKeySizeLevel()))
	check(t, golden, seen, prefix+".AlgorithmsExpirationDateLevel", string(suite.AlgorithmsExpirationDateLevel()))
	check(t, golden, seen, prefix+".AlgorithmsExpirationDateAfterUpdateLevel", string(suite.AlgorithmsExpirationDateAfterUpdateLevel()))

	updateDateStr := "<null>"
	if d := suite.CryptographicSuiteUpdateDate(); d != nil {
		updateDateStr = d.UTC().Format("2006-01-02")
	}
	check(t, golden, seen, prefix+".CryptographicSuiteUpdateDate", updateDateStr)
}

func checkCryptoSuitePanicking(t *testing.T, golden map[string]string, seen map[string]bool, key string, get func() modelpolicy.CryptographicSuite) {
	t.Helper()
	// The oracle's cryptoSuiteImpl always begins by emitting "<key>.Level" -
	// use that sub-key to decide whether the whole group was a caught
	// exception (the oracle's safe() wrapper emits a single "<key>" line,
	// with no ".Level" suffix, when the accessor call itself throws).
	if want, ok := golden[key]; ok && strings.HasPrefix(want, "<exception:") {
		seen[key] = true
		assertPanics(t, key, want, func() { get() })
		return
	}
	checkCryptoSuite(t, golden, seen, key, get())
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

// evaluationsCanonical renders evaluations the same way the oracle's
// dumpFirstEvaluation does: each element as
// "[minSizes=name=min,;validityEnd=YYYY-MM-DD]" (or "none"), the strings
// then sorted lexicographically and concatenated - see this file's header
// for why (Set has no defined order).
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
