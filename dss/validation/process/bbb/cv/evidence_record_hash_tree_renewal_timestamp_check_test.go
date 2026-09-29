package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EvidenceRecordHashTreeRenewalTimestampCheck: happy and failure path against
// the Java oracle (testdata/oracle/cv_direct.jsonl). The corpus KAT only reaches
// the failure path, so the covering case is driven over the same synthetic
// diagnostic data the oracle builds.
func TestEvidenceRecordHashTreeRenewalTimestampCheckAgainstJavaOracle(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		covered  bool
	}{
		{"er-hash-tree-renewal-ok", true},
		{"er-hash-tree-renewal-ko", false},
	} {
		diagnosticData := evidenceRecordDiagnosticData(tc.covered)
		renewal := diagnosticData.TimestampList()[0]
		assertDirectRow(t, tc.scenario,
			func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
				return NewEvidenceRecordHashTreeRenewalTimestampCheck(i18nProviderForTests, result,
					diagnosticData, renewal, rule)
			})
	}
}

// evidenceRecordDiagnosticData builds the diagnostic data the oracle drives the
// check with: one evidence record covering "doc.xml" and one HashTree-renewal
// archive time-stamp which does, or does not, cover it too.
func evidenceRecordDiagnosticData(covered bool) *diagnostic.Data {
	erMatcher := &diagnosticjaxb.XmlDigestMatcher{DataFound: true, DataIntact: true}
	setDigestMatcherType(erMatcher, enumerations.DigestMatcherTypeEvidenceRecordArchiveObject)
	erMatcher.DocumentName = ptr("doc.xml")

	tstMatcher := &diagnosticjaxb.XmlDigestMatcher{DataFound: true, DataIntact: true}
	setDigestMatcherType(tstMatcher, enumerations.DigestMatcherTypeEvidenceRecordArchiveObject)
	if covered {
		tstMatcher.DocumentName = ptr("doc.xml")
	} else {
		tstMatcher.DocumentName = ptr("other.xml")
	}

	timestampType := diagnosticjaxb.TimestampTypeValue(enumerations.TimestampTypeEvidenceRecordTimestamp)
	erTimestampType := diagnosticjaxb.EvidenceRecordTimestampTypeValue(
		enumerations.EvidenceRecordTimestampTypeHashTreeRenewalArchiveTimestamp)
	timestamp := &diagnosticjaxb.XmlTimestamp{
		Type:                        &timestampType,
		EvidenceRecordTimestampType: &erTimestampType,
		DigestMatcher:               []*diagnosticjaxb.XmlDigestMatcher{tstMatcher},
	}
	timestamp.Id = diagnosticjaxb.NewCollapsedString("T-SYNTHETIC")

	evidenceRecord := &diagnosticjaxb.XmlEvidenceRecord{
		DigestMatchers: &diagnosticjaxb.DigestMatchersWrapper{
			Items: []*diagnosticjaxb.XmlDigestMatcher{erMatcher},
		},
		EvidenceRecordTimestamps: &diagnosticjaxb.EvidenceRecordTimestampsWrapper{
			Items: []*diagnosticjaxb.XmlFoundTimestamp{{Timestamp: timestamp}},
		},
	}
	evidenceRecord.Id = diagnosticjaxb.NewCollapsedString("ER-SYNTHETIC")

	return diagnostic.NewData(&diagnosticjaxb.XmlDiagnosticData{
		EvidenceRecords: &diagnosticjaxb.EvidenceRecordsWrapper{
			Items: []*diagnosticjaxb.XmlEvidenceRecord{evidenceRecord},
		},
		UsedTimestamps: &diagnosticjaxb.UsedTimestampsWrapper{
			Items: []*diagnosticjaxb.XmlTimestamp{timestamp},
		},
	})
}

// multisetCovers must agree with the upstream scan-and-remove loop for every
// input, duplicates included (the removal "to avoid checking duplicates" makes
// this a multiset, not a set, inclusion test).
func TestMultisetCoversMatchesScanAndRemove(t *testing.T) {
	reference := func(required, available []string) bool {
		available = append([]string(nil), available...)
		for _, name := range required {
			index := -1
			for i, candidate := range available {
				if candidate == name {
					index = i
					break
				}
			}
			if index < 0 {
				return false
			}
			available = append(available[:index], available[index+1:]...)
		}
		return true
	}
	names := []string{"", "a", "b", "c"}
	var lists [][]string
	var build func(prefix []string, depth int)
	build = func(prefix []string, depth int) {
		lists = append(lists, append([]string(nil), prefix...))
		if depth == 0 {
			return
		}
		for _, n := range names {
			build(append(prefix, n), depth-1)
		}
	}
	build(nil, 3)
	for _, required := range lists {
		for _, available := range lists {
			if got, want := multisetCovers(required, available), reference(required, available); got != want {
				t.Fatalf("multisetCovers(%q, %q) = %v, want %v", required, available, got, want)
			}
		}
	}
}
