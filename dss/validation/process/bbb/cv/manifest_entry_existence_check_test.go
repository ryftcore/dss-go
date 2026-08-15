package cv

import (
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ManifestEntryExistenceCheck: happy and failure path against the Java oracle
// (testdata/oracle/cv_direct.jsonl). The failure path also pins the
// ISMEC_ANS_2 error-message branch (manifest entries present, none found).
func TestManifestEntryExistenceCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario       string
		digestMatchers []*diagnosticjaxb.XmlDigestMatcher
	}{
		{"manifest-entry-existence-ok", manifestEntries(true, true)},
		{"manifest-entry-existence-ko", manifestEntries(false, false)},
	} {
		digestMatchers := tc.digestMatchers
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
				return NewManifestEntryExistenceCheck(i18nProviderForTests, result, digestMatchers, rule)
			})
	}
}
