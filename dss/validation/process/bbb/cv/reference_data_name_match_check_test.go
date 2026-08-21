package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ReferenceDataNameMatchCheck's failure path against the Java oracle
// (testdata/oracle/cv_direct.jsonl), pinning the MANIFEST_ENTRY message tags and
// the additional-info message. Its happy path is exercised by the corpus KAT
// (BBB_CV_DMENMND / BBB_CV_DRNMND, status OK).
func TestReferenceDataNameMatchCheckAgainstJavaOracle(t *testing.T) {
	digestMatcher := manifestEntries(true, false)[1]
	assertDirectRow(t, "reference-name-match-ko",
		func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
			return NewReferenceDataNameMatchCheck(i18nProviderForTests, result, digestMatcher, rule)
		})
}
