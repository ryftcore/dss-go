package cv

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ManifestEntryGroupCheck: happy and failure path against the Java oracle
// (testdata/oracle/cv_direct.jsonl), including its additional-info message.
func TestManifestEntryGroupCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario       string
		digestMatchers []*diagnosticjaxb.XmlDigestMatcher
	}{
		{"manifest-entry-group-ok", manifestEntries(true, true)},
		{"manifest-entry-group-ko", manifestEntries(true, false)},
	} {
		digestMatchers := tc.digestMatchers
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
				return NewManifestEntryGroupCheck(i18nProviderForTests, result, digestMatchers, rule)
			})
	}
}
