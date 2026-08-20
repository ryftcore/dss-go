// Oracle test for TrustAnchorVerifierFactory and RevocationDataVerifierFactory.
//
// testdata/oracle/verifier_factories.tsv is a pure Java dump, produced by
// testdata/oracle/gen/VerifierFactoryOracle.java, which drives the upstream
// classes eu.europa.esig.dss.validation.TrustAnchorVerifierFactory and
// ...RevocationDataVerifierFactory over the four validation policies DSS ships
// (the default one plus the certificate/QWAC/EAA constraint files this package
// embeds under resources/).
//
// The dump carries two lines per policy. The "tav." line is asserted here in
// full: TrustAnchorVerifier exposes every value it configures. The "rdv." line
// records what the Java generator can only read by reflection - upstream's
// RevocationDataVerifier keeps those fields package-private and offers no
// getters, and so does the port - so it is shipped as evidence of what the
// factory must produce, and the Go side asserts what is reachable (that the
// factory builds a verifier for every policy without failing). Closing that gap
// needs behavioural assertions through IsRevocationDataSkip /
// IsRevocationDataFresh, which in turn need real revocation tokens; see the
// phase-8f porter notes.

package validation

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/model"
	modelpolicy "github.com/utain/esig/dss/model/policy"
	dsspolicy "github.com/utain/esig/dss/policy"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

// verifierFactoryOracleTime is the validation time the Java dump used.
var verifierFactoryOracleTime = time.UnixMilli(1700000000000).UTC()

func TestVerifierFactoriesOracle(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())

	data, err := os.ReadFile(filepath.Join("testdata", "oracle", "verifier_factories.tsv"))
	if err != nil {
		t.Fatalf("reading the oracle dump: %v", err)
	}

	asserted := 0
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			t.Fatalf("malformed oracle row: %q", line)
		}
		policyName := fields[0]
		if !strings.HasPrefix(fields[1], "tav.") {
			// The reflection-only RevocationDataVerifier row; see the file header.
			continue
		}
		t.Run(policyName, func(t *testing.T) {
			validationPolicy := loadOraclePolicy(t, policyName)

			trustAnchorVerifier := NewTrustAnchorVerifierFactory(validationPolicy).Create()
			values := map[string]bool{
				"tav.acceptRevocationUntrusted": trustAnchorVerifier.IsAcceptRevocationUntrustedCertificateChains(),
				"tav.acceptTimestampUntrusted":  trustAnchorVerifier.IsAcceptTimestampUntrustedCertificateChains(),
				"tav.useSunsetDate":             trustAnchorVerifier.IsUseSunsetDate(),
			}
			for _, field := range fields[1:] {
				name, want, found := strings.Cut(field, "=")
				if !found {
					t.Fatalf("malformed oracle field: %q", field)
				}
				got, known := values[name]
				if !known {
					t.Fatalf("unknown oracle field: %q", name)
				}
				wantBool, err := strconv.ParseBool(want)
				if err != nil {
					t.Fatalf("malformed oracle value for %s: %q", name, want)
				}
				if got != wantBool {
					t.Errorf("%s = %v, want %v", name, got, wantBool)
				}
			}

			// RevocationDataVerifierFactory has no readable configuration (see
			// the file header); assert it builds for every shipped policy.
			revocationDataVerifier := NewRevocationDataVerifierFactory(validationPolicy).
				SetValidationTime(verifierFactoryOracleTime).Create()
			if revocationDataVerifier == nil {
				t.Fatal("RevocationDataVerifierFactory.Create() returned nil")
			}
		})
		asserted++
	}
	if asserted != 4 {
		t.Fatalf("asserted %d policies, want 4", asserted)
	}
}

// loadOraclePolicy resolves a policy name of the Java dump ("default" or a
// "/policy/<name>.xml" classpath path) to a ValidationPolicy.
func loadOraclePolicy(t *testing.T, name string) modelpolicy.ValidationPolicy {
	t.Helper()
	if name == "default" {
		return validationpolicy.FromDefaultValidationPolicy().Create()
	}
	data, err := os.ReadFile(filepath.Join("resources", filepath.Base(name)))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return validationpolicy.FromValidationPolicyDocument(model.NewInMemoryDocument(data)).Create()
}
