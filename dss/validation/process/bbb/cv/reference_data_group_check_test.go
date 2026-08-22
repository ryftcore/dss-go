package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ReferenceDataGroupCheck: happy and failure path against the Java oracle
// (testdata/oracle/cv_direct.jsonl).
func TestReferenceDataGroupCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario       string
		digestMatchers []*diagnosticjaxb.XmlDigestMatcher
	}{
		{"reference-data-group-ok", evidenceRecordDigestMatchers(true, false)},
		{"reference-data-group-ko", evidenceRecordDigestMatchers(true, true)},
	} {
		digestMatchers := tc.digestMatchers
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
				return NewReferenceDataGroupCheck(i18nProviderForTests, result, digestMatchers, rule)
			})
	}
}
