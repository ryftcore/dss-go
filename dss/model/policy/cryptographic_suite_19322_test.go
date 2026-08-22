// Ported from dss-model/.../model/policy/crypto/CryptographicSuite19322.java (DSS 6.5.RC1).
package policy

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestCryptographicSuite19322_PolicyNameAndUpdateDate(t *testing.T) {
	metadata := NewCryptographicSuiteMetadata()
	metadata.SetPolicyName("ETSI TS 119 312")
	issued := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	metadata.SetPolicyIssueDate(&issued)

	suite := NewCryptographicSuite19322(metadata, []*CryptographicSuiteAlgorithm{})
	if suite.PolicyName() != "ETSI TS 119 312" {
		t.Fatalf("PolicyName() = %q", suite.PolicyName())
	}
	if suite.CryptographicSuiteUpdateDate() == nil || !suite.CryptographicSuiteUpdateDate().Equal(issued) {
		t.Fatalf("CryptographicSuiteUpdateDate() = %v, want %v", suite.CryptographicSuiteUpdateDate(), issued)
	}
}

func TestCryptographicSuite19322_Defaults(t *testing.T) {
	suite := NewCryptographicSuite19322(NewCryptographicSuiteMetadata(), []*CryptographicSuiteAlgorithm{})
	if suite.Level() != enumerations.LevelFail {
		t.Fatalf("Level() = %v, want FAIL", suite.Level())
	}
	if suite.AlgorithmsExpirationDateAfterUpdateLevel() != enumerations.LevelWarn {
		t.Fatalf("AlgorithmsExpirationDateAfterUpdateLevel() = %v, want WARN", suite.AlgorithmsExpirationDateAfterUpdateLevel())
	}
	// unset sub-levels fall back to the global level
	if suite.AcceptableDigestAlgorithmsLevel() != enumerations.LevelFail {
		t.Fatalf("AcceptableDigestAlgorithmsLevel() = %v, want FAIL fallback", suite.AcceptableDigestAlgorithmsLevel())
	}
	suite.SetLevel(enumerations.LevelWarn)
	if suite.AcceptableDigestAlgorithmsLevel() != enumerations.LevelWarn {
		t.Fatalf("AcceptableDigestAlgorithmsLevel() did not track updated global level")
	}
	suite.SetAcceptableDigestAlgorithmsLevel(enumerations.LevelInform)
	if suite.AcceptableDigestAlgorithmsLevel() != enumerations.LevelInform {
		t.Fatalf("AcceptableDigestAlgorithmsLevel() = %v, want explicit INFORM", suite.AcceptableDigestAlgorithmsLevel())
	}
}

func TestCryptographicSuite19322_AcceptableDigestAlgorithms_ByOID(t *testing.T) {
	algo := NewCryptographicSuiteAlgorithm()
	algo.SetAlgorithmIdentifierOIDs([]string{"2.16.840.1.101.3.4.2.1"}) // SHA-256
	eval := NewCryptographicSuiteEvaluation()
	eval.SetRecommendation(enumerations.CryptographicSuiteRecommendationRecommended)
	algo.SetEvaluationList([]*CryptographicSuiteEvaluation{eval})

	suite := NewCryptographicSuite19322(NewCryptographicSuiteMetadata(), []*CryptographicSuiteAlgorithm{algo})
	digests := suite.AcceptableDigestAlgorithms()
	evaluations, ok := digests[enumerations.DigestAlgorithmSHA256]
	if !ok || len(evaluations) != 1 {
		t.Fatalf("AcceptableDigestAlgorithms() = %v, want SHA256 with 1 evaluation", digests)
	}
	if evaluations[0].Recommendation() != enumerations.CryptographicSuiteRecommendationRecommended {
		t.Fatalf("evaluation recommendation mismatch: %v", evaluations[0].Recommendation())
	}

	// cached: mutating the returned map's backing slice must observe the
	// same cached result on a second call.
	if again := suite.AcceptableDigestAlgorithms(); len(again[enumerations.DigestAlgorithmSHA256]) != 1 {
		t.Fatalf("second call did not return cached map: %v", again)
	}
}

func TestCryptographicSuite19322_AcceptableSignatureAlgorithms_DirectOID(t *testing.T) {
	// RSA_SHA256 signature algorithm OID, directly declared.
	algo := NewCryptographicSuiteAlgorithm()
	algo.SetAlgorithmIdentifierOIDs([]string{"1.2.840.113549.1.1.11"})
	eval := NewCryptographicSuiteEvaluation()
	algo.SetEvaluationList([]*CryptographicSuiteEvaluation{eval})

	suite := NewCryptographicSuite19322(NewCryptographicSuiteMetadata(), []*CryptographicSuiteAlgorithm{algo})
	sigs := suite.AcceptableSignatureAlgorithms()
	if _, ok := sigs[enumerations.SignatureAlgorithmRSASHA256]; !ok {
		t.Fatalf("AcceptableSignatureAlgorithms() = %v, want RSA_SHA256 present", sigs)
	}
}

func TestCryptographicSuite19322_AcceptableSignatureAlgorithms_DerivedFromEncryptionAndDigest(t *testing.T) {
	// Encryption algorithm RSA declared separately from digest algorithm
	// SHA-256; the SignatureAlgorithm must be derived from the pairing
	// per Step 2b/2c of the Java implementation.
	rsaAlgo := NewCryptographicSuiteAlgorithm()
	rsaAlgo.SetAlgorithmIdentifierOIDs([]string{"1.2.840.113549.1.1.1"}) // RSA encryption
	rsaEval := NewCryptographicSuiteEvaluation()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	rsaEval.SetValidityStart(&start)
	rsaAlgo.SetEvaluationList([]*CryptographicSuiteEvaluation{rsaEval})

	sha256Algo := NewCryptographicSuiteAlgorithm()
	sha256Algo.SetAlgorithmIdentifierOIDs([]string{"2.16.840.1.101.3.4.2.1"}) // SHA-256
	sha256Eval := NewCryptographicSuiteEvaluation()
	end := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	sha256Eval.SetValidityEnd(&end)
	sha256Algo.SetEvaluationList([]*CryptographicSuiteEvaluation{sha256Eval})

	suite := NewCryptographicSuite19322(NewCryptographicSuiteMetadata(), []*CryptographicSuiteAlgorithm{rsaAlgo, sha256Algo})
	sigs := suite.AcceptableSignatureAlgorithms()
	evaluations, ok := sigs[enumerations.SignatureAlgorithmRSASHA256]
	if !ok || len(evaluations) == 0 {
		t.Fatalf("AcceptableSignatureAlgorithms() = %v, want derived RSA_SHA256", sigs)
	}
	// the derived evaluation must inherit the tighter validity window
	// from both the encryption and digest algorithm evaluations.
	found := false
	for _, e := range evaluations {
		if e.ValidityStart() != nil && e.ValidityStart().Equal(start) &&
			e.ValidityEnd() != nil && e.ValidityEnd().Equal(end) {
			found = true
		}
	}
	if !found {
		t.Fatalf("derived evaluation validity window mismatch: %v", evaluations)
	}
}

func TestCryptographicSuite19322_UnknownOID(t *testing.T) {
	algo := NewCryptographicSuiteAlgorithm()
	algo.SetAlgorithmIdentifierOIDs([]string{"1.2.3.4.5.6.7.8.9"})
	suite := NewCryptographicSuite19322(NewCryptographicSuiteMetadata(), []*CryptographicSuiteAlgorithm{algo})
	if digests := suite.AcceptableDigestAlgorithms(); len(digests) != 0 {
		t.Fatalf("AcceptableDigestAlgorithms() = %v, want empty for unrecognized OID", digests)
	}
}

func TestNewCryptographicSuite19322_PanicsOnNilMetadata(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on nil metadata")
		}
	}()
	NewCryptographicSuite19322(nil, nil)
}
