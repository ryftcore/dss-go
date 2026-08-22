package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AtLeastOneReferenceDataObjectFoundCheck: happy and failure path against the
// Java oracle (testdata/oracle/cv_direct.jsonl).
func TestAtLeastOneReferenceDataObjectFoundCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario       string
		digestMatchers []*diagnosticjaxb.XmlDigestMatcher
	}{
		{"at-least-one-found-ok", evidenceRecordDigestMatchers(true, false)},
		{"at-least-one-found-ko", evidenceRecordDigestMatchers(false, false)},
	} {
		digestMatchers := tc.digestMatchers
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
				return NewAtLeastOneReferenceDataObjectFoundCheck(i18nProviderForTests, result, digestMatchers, rule)
			})
	}
}
