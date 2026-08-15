package isc

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// DigestValuePresentCheck: happy and failure path against the Java oracle
// (testdata/oracle/isc_direct.jsonl).
func TestDigestValuePresentCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		token    func() *diagnostic.SignatureWrapper
	}{
		{"digest-value-present-ok", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, true, true) }},
		{"digest-value-present-ko", func() *diagnostic.SignatureWrapper { return syntheticSignature(false, false, true, true, true) }},
	} {
		token := tc.token()
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC] {
				return NewDigestValuePresentCheck(i18nProviderForTests, result, token, rule)
			})
	}
}
