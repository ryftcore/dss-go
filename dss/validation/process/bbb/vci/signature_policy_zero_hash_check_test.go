package vci

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignaturePolicyZeroHashCheck: happy and failure path against the Java oracle
// (testdata/oracle/vci_direct.jsonl).
func TestSignaturePolicyZeroHashCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario  string
		signature *diagnostic.SignatureWrapper
	}{
		{"policy-zero-hash-ok", policySignature(ptr("1.2.3.4"), ptr(true), ptr(false), ptr(true), true)},
		{"policy-zero-hash-ko", policySignature(ptr("1.2.3.4"), ptr(true), ptr(false), ptr(false), true)},
	} {
		signature := tc.signature
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI] {
				return NewSignaturePolicyZeroHashCheck(i18nProviderForTests, result, signature, rule)
			})
	}
}
