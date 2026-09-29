package xades

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/policy"
)

// configuredProvider is implemented by the nested analyzers (analyzer.DefaultDocumentAnalyzer).
type configuredProvider interface {
	ConfiguredSignaturePolicyProvider() *policy.SignaturePolicyProvider
}

// TestASiCContainerWithXAdESAnalyzerForwardsSignaturePolicyProvider checks that a
// SignaturePolicyProvider set on the outer analyzer reaches every nested XAdES analyzer, as
// upstream's getSignatureAnalyzers() does (A16-SEC-001), and that no default provider (whose
// HTTP data loader would download the policy named in the container) is created when none was
// set (the documented divergence in the analyzer's doc comment).
func TestASiCContainerWithXAdESAnalyzerForwardsSignaturePolicyProvider(t *testing.T) {
	const fixture = "../testdata/upstream/dss-asic-xades/src/test/resources/validation/onefile-ok.asics"

	newAnalyzer := func(t *testing.T) *ASiCContainerWithXAdESAnalyzer {
		t.Helper()
		doc, err := model.NewFileDocument(asicXadesFixturePath(t, fixture))
		if err != nil {
			t.Fatalf("NewFileDocument: %v", err)
		}
		analyzer := NewASiCContainerWithXAdESAnalyzer(doc)
		analyzer.SetCertificateVerifier(permissiveCertificateVerifier())
		return analyzer
	}

	t.Run("custom provider is forwarded", func(t *testing.T) {
		custom := policy.NewSignaturePolicyProvider()
		analyzer := newAnalyzer(t)
		analyzer.SetSignaturePolicyProvider(custom)

		nested := analyzer.GetSignatureAnalyzers()
		if len(nested) == 0 {
			t.Fatal("expected at least one nested signature analyzer")
		}
		for i, signatureAnalyzer := range nested {
			if got := signatureAnalyzer.(configuredProvider).ConfiguredSignaturePolicyProvider(); got != custom {
				t.Errorf("nested analyzer %d: SignaturePolicyProvider = %p, want the custom provider %p", i, got, custom)
			}
		}
	})

	t.Run("no default provider is created", func(t *testing.T) {
		analyzer := newAnalyzer(t)
		for i, signatureAnalyzer := range analyzer.GetSignatureAnalyzers() {
			if got := signatureAnalyzer.(configuredProvider).ConfiguredSignaturePolicyProvider(); got != nil {
				t.Errorf("nested analyzer %d: unexpected SignaturePolicyProvider %p", i, got)
			}
		}
	})
}
