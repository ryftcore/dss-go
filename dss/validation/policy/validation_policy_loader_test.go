package policy

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	cryptojson "github.com/ryftcore/dss-go/dss/policy/crypto/json"
	cryptoxml "github.com/ryftcore/dss-go/dss/policy/crypto/xml"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// Test-local ServiceLoader wiring: the concrete factory packages don't
// self-register - only the top-level dss package's init() (validate.go)
// does, and this package can't import that without an import cycle - so
// the test registers them here directly. This mirrors, for a test, exactly
// what a composing application would do once at startup.
func init() {
	RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	RegisterCryptographicSuiteFactory(cryptoxml.NewCryptographicSuiteXmlFactory())
	RegisterCryptographicSuiteFactory(cryptojson.NewCryptographicSuiteJsonFactory())
}

func TestValidationPolicyLoader_LoadDefault(t *testing.T) {
	vp := FromDefaultValidationPolicy().Create()
	if vp == nil {
		t.Fatal("expected a non-nil default validation policy")
	}

	vpWithCrypto := FromDefaultValidationPolicy().WithDefaultCryptographicSuite().Create()
	if vpWithCrypto == nil {
		t.Fatal("expected a non-nil validation policy")
	}
	if _, ok := vpWithCrypto.(*ValidationPolicyWithCryptographicSuite); !ok {
		t.Fatalf("expected *ValidationPolicyWithCryptographicSuite, got %T", vpWithCrypto)
	}
}

func TestValidationPolicyLoader_OverrideDefaultPolicyWithCryptographicSuite(t *testing.T) {
	fullSet := map[enumerations.DigestAlgorithm]bool{
		enumerations.DigestAlgorithmMD5: true, enumerations.DigestAlgorithmSHA1: true,
		enumerations.DigestAlgorithmSHA224: true, enumerations.DigestAlgorithmSHA256: true,
		enumerations.DigestAlgorithmSHA384: true, enumerations.DigestAlgorithmSHA512: true,
		enumerations.DigestAlgorithmSHA3256: true, enumerations.DigestAlgorithmSHA3384: true,
		enumerations.DigestAlgorithmSHA3512: true, enumerations.DigestAlgorithmRIPEMD160: true,
		enumerations.DigestAlgorithmWHIRLPOOL: true,
	}

	policy := FromDefaultValidationPolicy().WithDefaultCryptographicSuite().Create()
	got := digestAlgorithmKeys(t, policy.SignatureCryptographicConstraint(enumerations.ContextSignature).AcceptableDigestAlgorithms())
	if len(got) != len(fullSet) {
		t.Fatalf("SIGNATURE digest algorithms = %v, want the full default set (%d entries)", got, len(fullSet))
	}
	for algo := range fullSet {
		if !got[algo] {
			t.Fatalf("SIGNATURE digest algorithms missing %s: got %v", algo, got)
		}
	}

	// A global CryptographicSuite overrides all scopes, including EvidenceRecord, except for
	// SIGNATURE POLICY constraints (unaffected here).
	altCrypto := dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())
	overridden := FromDefaultValidationPolicy().WithCryptographicSuite(altCrypto).Create()

	for _, context := range enumerations.ContextValues() {
		if context == enumerations.ContextEvidenceRecord {
			continue
		}
		sigAlgos := digestAlgorithmKeys(t, overridden.SignatureCryptographicConstraint(context).AcceptableDigestAlgorithms())
		if len(sigAlgos) != 1 || !sigAlgos[enumerations.DigestAlgorithmSHA256] {
			t.Fatalf("context %s: expected {SHA256}, got %v", context, sigAlgos)
		}
	}
	erAlgos := digestAlgorithmKeys(t, overridden.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms())
	if len(erAlgos) != 1 || !erAlgos[enumerations.DigestAlgorithmSHA256] {
		t.Fatalf("EvidenceRecord: expected {SHA256}, got %v", erAlgos)
	}
}

func TestValidationPolicyLoader_WithCryptographicSuiteForContextScopesOnlyThatContext(t *testing.T) {
	altCrypto := dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())

	for _, context := range enumerations.ContextValues() {
		vp := FromDefaultValidationPolicy().WithCryptographicSuiteForContext(altCrypto, context).Create()

		for _, currentContext := range enumerations.ContextValues() {
			if currentContext == enumerations.ContextEvidenceRecord {
				erAlgos := digestAlgorithmKeys(t, vp.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms())
				wantSHA256Only := context == currentContext
				gotSHA256Only := len(erAlgos) == 1 && erAlgos[enumerations.DigestAlgorithmSHA256]
				if gotSHA256Only != wantSHA256Only {
					t.Fatalf("context %s / current %s: EvidenceRecord algos = %v, want SHA256-only=%v", context, currentContext, erAlgos, wantSHA256Only)
				}
				continue
			}
			sigAlgos := digestAlgorithmKeys(t, vp.SignatureCryptographicConstraint(currentContext).AcceptableDigestAlgorithms())
			wantSHA256Only := context == currentContext
			gotSHA256Only := len(sigAlgos) == 1 && sigAlgos[enumerations.DigestAlgorithmSHA256]
			if gotSHA256Only != wantSHA256Only {
				t.Fatalf("context %s / current %s: Signature algos = %v, want SHA256-only=%v", context, currentContext, sigAlgos, wantSHA256Only)
			}
		}
	}
}

func TestValidationPolicyLoader_WithLevels(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	vp := FromValidationPolicy(dsspolicy.NewEtsiValidationPolicy(cp)).
		WithCryptographicSuite(dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())).
		AndLevel(enumerations.LevelWarn).
		Create()

	cs := vp.SignatureCryptographicConstraint(enumerations.ContextSignature)
	if cs.AcceptableDigestAlgorithmsLevel() != enumerations.LevelWarn {
		t.Fatalf("AcceptableDigestAlgorithmsLevel = %s, want WARN", cs.AcceptableDigestAlgorithmsLevel())
	}
	if cs.AcceptableSignatureAlgorithmsLevel() != enumerations.LevelWarn {
		t.Fatalf("AcceptableSignatureAlgorithmsLevel = %s, want WARN", cs.AcceptableSignatureAlgorithmsLevel())
	}

	vp2 := FromValidationPolicy(dsspolicy.NewEtsiValidationPolicy(cp)).
		WithCryptographicSuite(dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())).
		AndLevel(enumerations.LevelWarn).
		AndAcceptableDigestAlgorithmsLevel(enumerations.LevelFail).
		Create()
	cs2 := vp2.SignatureCryptographicConstraint(enumerations.ContextSignature)
	if cs2.AcceptableDigestAlgorithmsLevel() != enumerations.LevelFail {
		t.Fatalf("AcceptableDigestAlgorithmsLevel = %s, want FAIL (overridden after AndLevel)", cs2.AcceptableDigestAlgorithmsLevel())
	}
	if cs2.AcceptableSignatureAlgorithmsLevel() != enumerations.LevelWarn {
		t.Fatalf("AcceptableSignatureAlgorithmsLevel = %s, want WARN (untouched by the digest-only override)", cs2.AcceptableSignatureAlgorithmsLevel())
	}
}

func TestValidationPolicyLoader_NotSupportedPolicyDocumentPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for an unsupported validation policy document")
		}
	}()
	// A cryptographic-suite JSON document is not a valid validation-policy XML document.
	FromValidationPolicyDocument(model.NewInMemoryDocument([]byte(`{"not":"a policy"}`)))
}

func TestValidationPolicyLoader_NilInputsPanic(t *testing.T) {
	cases := []struct {
		name string
		fn   func()
	}{
		{"FromValidationPolicyDocument(nil)", func() { FromValidationPolicyDocument(nil) }},
		{"FromValidationPolicyReader(nil)", func() { FromValidationPolicyReader(nil) }},
		{"FromValidationPolicy(nil)", func() { FromValidationPolicy(nil) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected %s to panic", c.name)
				}
			}()
			c.fn()
		})
	}
}

func TestValidationPolicyLoader_CryptoSuiteNilInputsPanic(t *testing.T) {
	loader := FromDefaultValidationPolicy()
	cases := []struct {
		name string
		fn   func()
	}{
		{"WithCryptographicSuiteDocument(nil)", func() { loader.WithCryptographicSuiteDocument(nil) }},
		{"WithCryptographicSuiteReader(nil)", func() { loader.WithCryptographicSuiteReader(nil) }},
		{"WithCryptographicSuite(nil)", func() { loader.WithCryptographicSuite(nil) }},
		{"WithCryptographicSuiteCatalogue(nil)", func() { loader.WithCryptographicSuiteCatalogue(nil) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected %s to panic", c.name)
				}
			}()
			c.fn()
		})
	}
}

func TestValidationPolicyLoader_DefaultCryptographicSuiteForEvidenceRecordWithSubContextErrors(t *testing.T) {
	loader := FromDefaultValidationPolicy()
	if _, err := loader.WithDefaultCryptographicSuiteForContextAndSubContext(enumerations.ContextEvidenceRecord, enumerations.SubContextSigningCert); err == nil {
		t.Fatal("expected an error for EVIDENCE_RECORD with a non-empty SubContext")
	}
}
