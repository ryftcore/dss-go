package vci

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignaturePolicyHashValidCheck: happy and failure path against the Java oracle
// (testdata/oracle/vci_direct.jsonl).
func TestSignaturePolicyHashValidCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario  string
		signature *diagnostic.SignatureWrapper
	}{
		{"policy-hash-valid-ok", policySignature(ptr("1.2.3.4"), ptr(true), ptr(true), ptr(false), true)},
		{"policy-hash-valid-ko", policySignature(ptr("1.2.3.4"), ptr(true), ptr(false), ptr(false), true)},
	} {
		signature := tc.signature
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI] {
				return NewSignaturePolicyHashValidCheck(i18nProviderForTests, result, signature, rule)
			})
	}
}
