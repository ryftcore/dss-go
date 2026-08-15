package isc

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// IssuerSerialMatchCheck: happy and failure path against the Java oracle
// (testdata/oracle/isc_direct.jsonl).
func TestIssuerSerialMatchCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		token    func() *diagnostic.SignatureWrapper
	}{
		{"issuer-serial-match-ok", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, true, true) }},
		{"issuer-serial-match-ko", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, false, true) }},
	} {
		token := tc.token()
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC] {
				return NewIssuerSerialMatchCheck(i18nProviderForTests, result, token, rule)
			})
	}
}
