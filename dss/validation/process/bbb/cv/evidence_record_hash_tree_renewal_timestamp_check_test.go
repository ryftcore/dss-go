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
func evidenceRecordDiagnosticData(covered bool) *diagnostic.DiagnosticData {
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

	return diagnostic.NewDiagnosticData(&diagnosticjaxb.XmlDiagnosticData{
		EvidenceRecords: &diagnosticjaxb.EvidenceRecordsWrapper{
			Items: []*diagnosticjaxb.XmlEvidenceRecord{evidenceRecord},
		},
		UsedTimestamps: &diagnosticjaxb.UsedTimestampsWrapper{
			Items: []*diagnosticjaxb.XmlTimestamp{timestamp},
		},
	})
}
