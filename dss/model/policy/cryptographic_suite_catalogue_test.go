// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteCatalogue.java (DSS 6.5.RC1).
package policy

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func newTestCatalogue() *CryptographicSuiteCatalogue {
	globalAlgo := NewCryptographicSuiteAlgorithm()
	globalAlgo.SetAlgorithmIdentifierOIDs([]string{"2.16.840.1.101.3.4.2.1"})
	// no evaluation list -> unconditionally included regardless of usage

	certAlgo := NewCryptographicSuiteAlgorithm()
	certAlgo.SetAlgorithmIdentifierOIDs([]string{"1.2.840.113549.1.1.11"})
	certEval := NewCryptographicSuiteEvaluation()
	certEval.SetAlgorithmUsage([]enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsageSignCertificates,
		enumerations.CryptographicSuiteAlgorithmUsageValidateCertificates,
	})
	certAlgo.SetEvaluationList([]*CryptographicSuiteEvaluation{certEval})

	return NewCryptographicSuiteCatalogue(
		func() *CryptographicSuiteMetadata {
			m := NewCryptographicSuiteMetadata()
			m.SetPolicyName("test-catalogue")
			return m
		},
		func() []*CryptographicSuiteAlgorithm {
			return []*CryptographicSuiteAlgorithm{globalAlgo, certAlgo}
		},
	)
}

func TestCryptographicSuiteCatalogue_GlobalExcludesCertificateOnlyUsage(t *testing.T) {
	catalogue := newTestCatalogue()
	suite := catalogue.CryptographicSuite()
	if suite.PolicyName() != "test-catalogue" {
		t.Fatalf("PolicyName() = %q", suite.PolicyName())
	}
	digests := suite.AcceptableDigestAlgorithms()
	if len(digests) != 1 {
		t.Fatalf("AcceptableDigestAlgorithms() = %v, want the unconditional SHA256 entry only", digests)
	}
	sigs := suite.AcceptableSignatureAlgorithms()
	// RSA_SHA256 is SIGN_CERTIFICATES/VALIDATE_CERTIFICATES scoped only,
	// so it must be filtered out of the global (SIGN_DATA/VALIDATE_DATA)
	// suite.
	if _, ok := sigs[enumerations.SignatureAlgorithmRSASHA256]; ok {
		t.Fatalf("global suite unexpectedly includes certificate-only RSA_SHA256: %v", sigs)
	}
}

func TestCryptographicSuiteCatalogue_SignatureCertificatesIncludesCertUsage(t *testing.T) {
	catalogue := newTestCatalogue()
	suite := catalogue.SignatureCertificatesCryptographicSuite()
	sigs := suite.AcceptableSignatureAlgorithms()
	if _, ok := sigs[enumerations.SignatureAlgorithmRSASHA256]; !ok {
		t.Fatalf("SignatureCertificatesCryptographicSuite() missing RSA_SHA256: %v", sigs)
	}
}

func TestCryptographicSuiteCatalogue_DelegatesToGlobal(t *testing.T) {
	catalogue := newTestCatalogue()
	if catalogue.SignatureCryptographicSuite() != catalogue.CryptographicSuite() {
		t.Fatalf("SignatureCryptographicSuite() must return the same cached instance as CryptographicSuite()")
	}
	if catalogue.CounterSignatureCertificatesCryptographicSuite() != catalogue.SignatureCertificatesCryptographicSuite() {
		t.Fatalf("CounterSignatureCertificatesCryptographicSuite() must delegate to SignatureCertificatesCryptographicSuite()")
	}
}

func TestCryptographicSuiteCatalogue_MetadataAndAlgorithmListAreLazyAndCached(t *testing.T) {
	calls := 0
	catalogue := NewCryptographicSuiteCatalogue(
		func() *CryptographicSuiteMetadata {
			calls++
			return NewCryptographicSuiteMetadata()
		},
		func() []*CryptographicSuiteAlgorithm { return nil },
	)
	catalogue.Metadata()
	catalogue.Metadata()
	if calls != 1 {
		t.Fatalf("buildMetadata called %d times, want 1 (lazy + cached)", calls)
	}
}
