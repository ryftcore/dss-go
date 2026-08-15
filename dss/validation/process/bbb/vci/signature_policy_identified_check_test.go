package vci

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignaturePolicyIdentifiedCheck: happy and failure path against the Java oracle
// (testdata/oracle/vci_direct.jsonl).
func TestSignaturePolicyIdentifiedCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario  string
		signature *diagnostic.SignatureWrapper
	}{
		{"policy-identified-ok", policySignature(ptr("1.2.3.4"), ptr(true), ptr(true), ptr(false), true)},
		{"policy-identified-ko", policySignature(ptr("1.2.3.4"), ptr(false), ptr(true), ptr(false), true)},
	} {
		signature := tc.signature
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI] {
				return NewSignaturePolicyIdentifiedCheck(i18nProviderForTests, result, signature, rule)
			})
	}
}
