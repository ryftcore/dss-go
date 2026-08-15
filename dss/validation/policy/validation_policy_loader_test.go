package policy

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dsspolicy "github.com/utain/esig/dss/policy"
	cryptojson "github.com/utain/esig/dss/policy/crypto/json"
	cryptoxml "github.com/utain/esig/dss/policy/crypto/xml"
	"github.com/utain/esig/dss/policy/jaxb"
)

// Test-local ServiceLoader wiring: this chunk's manifest does not include (and PORTING.md
// forbids editing) the frozen factory packages, so the concrete implementations are registered
// here rather than via their own init() - see validation_policy_loader.go's file header. This
// mirrors, for a test, exactly what a composing application would do once at startup.
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
		enumerations.DigestAlgorithm_MD5: true, enumerations.DigestAlgorithm_SHA1: true,
		enumerations.DigestAlgorithm_SHA224: true, enumerations.DigestAlgorithm_SHA256: true,
		enumerations.DigestAlgorithm_SHA384: true, enumerations.DigestAlgorithm_SHA512: true,
		enumerations.DigestAlgorithm_SHA3_256: true, enumerations.DigestAlgorithm_SHA3_384: true,
		enumerations.DigestAlgorithm_SHA3_512: true, enumerations.DigestAlgorithm_RIPEMD160: true,
		enumerations.DigestAlgorithm_WHIRLPOOL: true,
	}

	policy := FromDefaultValidationPolicy().WithDefaultCryptographicSuite().Create()
	got := digestAlgorithmKeys(t, policy.SignatureCryptographicConstraint(enumerations.Context_SIGNATURE).AcceptableDigestAlgorithms())
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
		if context == enumerations.Context_EVIDENCE_RECORD {
			continue
		}
		sigAlgos := digestAlgorithmKeys(t, overridden.SignatureCryptographicConstraint(context).AcceptableDigestAlgorithms())
		if len(sigAlgos) != 1 || !sigAlgos[enumerations.DigestAlgorithm_SHA256] {
			t.Fatalf("context %s: expected {SHA256}, got %v", context, sigAlgos)
		}
	}
	erAlgos := digestAlgorithmKeys(t, overridden.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms())
	if len(erAlgos) != 1 || !erAlgos[enumerations.DigestAlgorithm_SHA256] {
		t.Fatalf("EvidenceRecord: expected {SHA256}, got %v", erAlgos)
	}
}

func TestValidationPolicyLoader_WithCryptographicSuiteForContextScopesOnlyThatContext(t *testing.T) {
	altCrypto := dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())

	for _, context := range enumerations.ContextValues() {
		vp := FromDefaultValidationPolicy().WithCryptographicSuiteForContext(altCrypto, context).Create()

		for _, currentContext := range enumerations.ContextValues() {
			if currentContext == enumerations.Context_EVIDENCE_RECORD {
				erAlgos := digestAlgorithmKeys(t, vp.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms())
				wantSHA256Only := context == currentContext
				gotSHA256Only := len(erAlgos) == 1 && erAlgos[enumerations.DigestAlgorithm_SHA256]
				if gotSHA256Only != wantSHA256Only {
					t.Fatalf("context %s / current %s: EvidenceRecord algos = %v, want SHA256-only=%v", context, currentContext, erAlgos, wantSHA256Only)
				}
				continue
			}
			sigAlgos := digestAlgorithmKeys(t, vp.SignatureCryptographicConstraint(currentContext).AcceptableDigestAlgorithms())
			wantSHA256Only := context == currentContext
			gotSHA256Only := len(sigAlgos) == 1 && sigAlgos[enumerations.DigestAlgorithm_SHA256]
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
		AndLevel(enumerations.Level_WARN).
		Create()

	cs := vp.SignatureCryptographicConstraint(enumerations.Context_SIGNATURE)
	if cs.AcceptableDigestAlgorithmsLevel() != enumerations.Level_WARN {
		t.Fatalf("AcceptableDigestAlgorithmsLevel = %s, want WARN", cs.AcceptableDigestAlgorithmsLevel())
	}
	if cs.AcceptableSignatureAlgorithmsLevel() != enumerations.Level_WARN {
		t.Fatalf("AcceptableSignatureAlgorithmsLevel = %s, want WARN", cs.AcceptableSignatureAlgorithmsLevel())
	}

	vp2 := FromValidationPolicy(dsspolicy.NewEtsiValidationPolicy(cp)).
		WithCryptographicSuite(dsspolicy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint())).
		AndLevel(enumerations.Level_WARN).
		AndAcceptableDigestAlgorithmsLevel(enumerations.Level_FAIL).
		Create()
	cs2 := vp2.SignatureCryptographicConstraint(enumerations.Context_SIGNATURE)
	if cs2.AcceptableDigestAlgorithmsLevel() != enumerations.Level_FAIL {
		t.Fatalf("AcceptableDigestAlgorithmsLevel = %s, want FAIL (overridden after AndLevel)", cs2.AcceptableDigestAlgorithmsLevel())
	}
	if cs2.AcceptableSignatureAlgorithmsLevel() != enumerations.Level_WARN {
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
	if _, err := loader.WithDefaultCryptographicSuiteForContextAndSubContext(enumerations.Context_EVIDENCE_RECORD, enumerations.SubContext_SIGNING_CERT); err == nil {
		t.Fatal("expected an error for EVIDENCE_RECORD with a non-empty SubContext")
	}
}
