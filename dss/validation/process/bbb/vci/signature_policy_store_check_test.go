package vci

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignaturePolicyStoreCheck: happy and failure path against the Java oracle
// (testdata/oracle/vci_direct.jsonl).
func TestSignaturePolicyStoreCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario  string
		signature *diagnostic.SignatureWrapper
	}{
		{"policy-store-ok", policySignature(ptr("1.2.3.4"), ptr(true), ptr(true), ptr(false), true)},
		{"policy-store-ko", policySignature(ptr("1.2.3.4"), ptr(true), ptr(true), ptr(false), false)},
	} {
		signature := tc.signature
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI] {
				return NewSignaturePolicyStoreCheck(i18nProviderForTests, result, signature, rule)
			})
	}
}
