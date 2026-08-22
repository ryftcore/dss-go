package policy

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

func sha1CryptographicConstraint() *jaxb.CryptographicConstraint {
	cc := &jaxb.CryptographicConstraint{
		AcceptableDigestAlgo: &jaxb.ListAlgo{Algos: []*jaxb.Algo{{Value: "SHA1"}}},
	}
	cc.Level = jaxb.LevelValue(enumerations.LevelFail)
	return cc
}

func sha256CryptographicConstraint() *jaxb.CryptographicConstraint {
	cc := &jaxb.CryptographicConstraint{
		AcceptableDigestAlgo: &jaxb.ListAlgo{Algos: []*jaxb.Algo{{Value: "SHA256"}}},
	}
	cc.Level = jaxb.LevelValue(enumerations.LevelFail)
	return cc
}

func digestAlgorithmKeys(t *testing.T, m map[enumerations.DigestAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation) map[enumerations.DigestAlgorithm]bool {
	t.Helper()
	keys := make(map[enumerations.DigestAlgorithm]bool, len(m))
	for k := range m {
		keys[k] = true
	}
	return keys
}

func TestValidationPolicyWithCryptographicSuite_DelegatesPolicyName(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	name := "HelloWorld"
	cp.Name = &name
	cp.Description = "Test policy"

	p := NewValidationPolicyWithCryptographicSuite(policy.NewEtsiValidationPolicy(cp))

	if got := p.PolicyName(); got != "HelloWorld" {
		t.Fatalf("PolicyName() = %q, want HelloWorld", got)
	}
	if got := p.PolicyDescription(); got != "Test policy" {
		t.Fatalf("PolicyDescription() = %q, want %q", got, "Test policy")
	}
}

func TestValidationPolicyWithCryptographicSuite_GlobalCryptoSuiteAppliesEverywhere(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	etsi := policy.NewEtsiValidationPolicy(cp)
	p := NewValidationPolicyWithCryptographicSuite(etsi)

	for _, context := range enumerations.ContextValues() {
		if context == enumerations.ContextEvidenceRecord {
			continue
		}
		if got := len(p.SignatureCryptographicConstraint(context).AcceptableDigestAlgorithms()); got != 0 {
			t.Fatalf("context %s: expected empty digest algorithms before setting a suite, got %d", context, got)
		}
	}

	p.SetCryptographicSuite(policy.NewCryptographicConstraintWrapper(sha1CryptographicConstraint()))

	for _, context := range enumerations.ContextValues() {
		if context == enumerations.ContextEvidenceRecord {
			continue
		}
		digestAlgos := p.SignatureCryptographicConstraint(context).AcceptableDigestAlgorithms()
		if _, ok := digestAlgos[enumerations.DigestAlgorithmSHA1]; !ok || len(digestAlgos) != 1 {
			t.Fatalf("context %s: expected {SHA1}, got %v", context, digestAlgorithmKeys(t, digestAlgos))
		}
		for _, subContext := range enumerations.SubContextValues() {
			certDigestAlgos := p.CertificateCryptographicConstraint(context, subContext).AcceptableDigestAlgorithms()
			if _, ok := certDigestAlgos[enumerations.DigestAlgorithmSHA1]; !ok || len(certDigestAlgos) != 1 {
				t.Fatalf("context %s/%s: expected {SHA1}, got %v", context, subContext, digestAlgorithmKeys(t, certDigestAlgos))
			}
		}
	}
	erDigestAlgos := p.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms()
	if _, ok := erDigestAlgos[enumerations.DigestAlgorithmSHA1]; !ok || len(erDigestAlgos) != 1 {
		t.Fatalf("EvidenceRecordCryptographicConstraint: expected {SHA1}, got %v", digestAlgorithmKeys(t, erDigestAlgos))
	}
}

func TestValidationPolicyWithCryptographicSuite_PerContextOverridesOnlyThatScope(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	etsi := policy.NewEtsiValidationPolicy(cp)

	for _, context := range enumerations.ContextValues() {
		p := NewValidationPolicyWithCryptographicSuite(etsi)
		p.SetCryptographicSuiteForContext(policy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint()), context)

		for _, currentContext := range enumerations.ContextValues() {
			if currentContext == enumerations.ContextEvidenceRecord {
				erDigestAlgos := p.EvidenceRecordCryptographicConstraint().AcceptableDigestAlgorithms()
				wantSHA256 := context == currentContext
				_, hasSHA256 := erDigestAlgos[enumerations.DigestAlgorithmSHA256]
				if hasSHA256 != wantSHA256 {
					t.Fatalf("context %s / current %s: EvidenceRecord SHA256 present=%v, want %v", context, currentContext, hasSHA256, wantSHA256)
				}
				continue
			}
			digestAlgos := p.SignatureCryptographicConstraint(currentContext).AcceptableDigestAlgorithms()
			wantSHA256 := context == currentContext
			_, hasSHA256 := digestAlgos[enumerations.DigestAlgorithmSHA256]
			if hasSHA256 != wantSHA256 {
				t.Fatalf("context %s / current %s: Signature SHA256 present=%v, want %v", context, currentContext, hasSHA256, wantSHA256)
			}
		}
	}
}

func TestValidationPolicyWithCryptographicSuite_SetCryptographicSuiteForContextAndSubContext(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	etsi := policy.NewEtsiValidationPolicy(cp)
	p := NewValidationPolicyWithCryptographicSuite(etsi)

	if err := p.SetCryptographicSuiteForContextAndSubContext(policy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint()),
		enumerations.ContextCertificate, enumerations.SubContextSigningCert); err != nil {
		t.Fatalf("SetCryptographicSuiteForContextAndSubContext returned error: %v", err)
	}

	digestAlgos := p.CertificateCryptographicConstraint(enumerations.ContextCertificate, enumerations.SubContextSigningCert).AcceptableDigestAlgorithms()
	if _, ok := digestAlgos[enumerations.DigestAlgorithmSHA256]; !ok {
		t.Fatalf("expected SHA256 in SIGNING_CERT scope, got %v", digestAlgorithmKeys(t, digestAlgos))
	}
	otherDigestAlgos := p.CertificateCryptographicConstraint(enumerations.ContextCertificate, enumerations.SubContextCACertificate).AcceptableDigestAlgorithms()
	if len(otherDigestAlgos) != 0 {
		t.Fatalf("expected CA_CERTIFICATE scope untouched, got %v", digestAlgorithmKeys(t, otherDigestAlgos))
	}
}

func TestValidationPolicyWithCryptographicSuite_EvidenceRecordSubContextRejected(t *testing.T) {
	cp := &jaxb.ConstraintsParameters{}
	etsi := policy.NewEtsiValidationPolicy(cp)
	p := NewValidationPolicyWithCryptographicSuite(etsi)

	err := p.SetCryptographicSuiteForContextAndSubContext(policy.NewCryptographicConstraintWrapper(sha256CryptographicConstraint()),
		enumerations.ContextEvidenceRecord, enumerations.SubContextSigningCert)
	if err == nil {
		t.Fatal("expected an error for EVIDENCE_RECORD with a non-empty SubContext")
	}
}
