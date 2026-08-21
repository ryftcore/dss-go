package isc

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningCertificateRecognitionCheck: happy and failure path against the Java oracle
// (testdata/oracle/isc_direct.jsonl).
func TestSigningCertificateRecognitionCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		token    func() *diagnostic.SignatureWrapper
	}{
		{"recognition-ok", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, true, true) }},
		{"recognition-ko", func() *diagnostic.SignatureWrapper { return syntheticSignature(true, true, true, true, false) }},
	} {
		token := tc.token()
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC] {
				return NewSigningCertificateRecognitionCheck(i18nProviderForTests, result, token, rule)
			})
	}
}
