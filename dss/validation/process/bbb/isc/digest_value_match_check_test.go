package isc

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DigestValueMatchCheck: happy and failure path against the Java oracle
// (testdata/oracle/isc_direct.jsonl).
func TestDigestValueMatchCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		token    func() *diagnostic.SignatureWrapper
	}{
		{"digest-value-match-ok", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, true, true) }},
		{"digest-value-match-ko", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, false, true, true, true) }},
	} {
		token := tc.token()
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC] {
				return NewDigestValueMatchCheck(i18nProviderForTests, result, token, rule)
			})
	}
}
