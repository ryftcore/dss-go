package cv

import (
	"reflect"
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// The direct-check KAT machinery: testdata/oracle/cv_direct.jsonl holds, for
// each scenario named below, the XmlCV upstream produces when the check is run
// alone at Level.FAIL - see ../testdata/gen/BbbBlocksOracle.java. It covers the
// three checks CryptographicVerification never wires (they belong to the
// evidence-record and archival blocks of later phases) plus the paths the
// corpus KAT leaves untouched.

// singleCheckChain is a chain of exactly one item, mirroring the oracle's
// SingleCheckChain.
type singleCheckChain struct {
	*process.ChainBase[*jaxb.XmlCV]
	factory func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV]
}

func newSingleCheckChain(
	factory func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV],
) *singleCheckChain {
	xmlCV := &jaxb.XmlCV{}
	c := &singleCheckChain{
		ChainBase: process.NewChainBase(i18nProviderForTests, process.NewResult(xmlCV,
			&xmlCV.XmlConstraintsConclusionContent, &xmlCV.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleCheckChain) InitChain() {
	c.FirstItem = c.factory(c.Result, process.GetLevelRule(enumerations.Level_FAIL))
}

// assertDirectRow runs the single-item chain and compares it against the named
// row of the direct oracle.
func assertDirectRow(t *testing.T, name string,
	factory func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV]) {
	t.Helper()
	rows := loadRows(t, corpustest.Path(t, "oracle/cv_direct.jsonl"))
	var expected *oracleRow
	for _, row := range rows {
		if row.Token == name {
			expected = row
		}
	}
	if expected == nil {
		t.Fatalf("no oracle row named %q", name)
	}
	if expected.Title == nil {
		// A chain with no title MessageTag: Java leaves the attribute null,
		// which the generated non-pointer Go member spells as "".
		empty := ""
		expected.Title = &empty
	}
	result := newSingleCheckChain(factory).Execute()
	actual := toRow(expected.File, expected.Token, enumerations.Context(expected.Context), "CV",
		&result.XmlConstraintsConclusionContent, result.Title)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: mismatch\nexpected: %s\nactual:   %s", name, mustJSON(t, expected), mustJSON(t, actual))
	}
}

// evidenceRecordDigestMatchers builds two EVIDENCE_RECORD_ARCHIVE_OBJECT
// matchers plus, optionally, an orphan reference - the oracle's digestMatchers.
func evidenceRecordDigestMatchers(dataFound, withOrphan bool) []*diagnosticjaxb.XmlDigestMatcher {
	var digestMatchers []*diagnosticjaxb.XmlDigestMatcher
	first := &diagnosticjaxb.XmlDigestMatcher{DataFound: dataFound, DataIntact: dataFound}
	setDigestMatcherType(first, enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT)
	first.DocumentName = ptr("doc.xml")
	digestMatchers = append(digestMatchers, first)
	second := &diagnosticjaxb.XmlDigestMatcher{DataFound: false, DataIntact: false}
	setDigestMatcherType(second, enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT)
	second.DocumentName = ptr("other.xml")
	digestMatchers = append(digestMatchers, second)
	if withOrphan {
		orphan := &diagnosticjaxb.XmlDigestMatcher{DataFound: true, DataIntact: true}
		setDigestMatcherType(orphan, enumerations.DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE)
		digestMatchers = append(digestMatchers, orphan)
	}
	return digestMatchers
}

// manifestEntries builds two MANIFEST_ENTRY matchers, the second one optionally
// not found - the oracle's manifestEntries.
func manifestEntries(firstFound, secondFound bool) []*diagnosticjaxb.XmlDigestMatcher {
	first := &diagnosticjaxb.XmlDigestMatcher{DataFound: firstFound, DataIntact: firstFound}
	setDigestMatcherType(first, enumerations.DigestMatcherType_MANIFEST_ENTRY)
	first.Uri = ptr("doc.xml")
	first.DocumentName = ptr("doc.xml")
	second := &diagnosticjaxb.XmlDigestMatcher{DataFound: secondFound, DataIntact: secondFound}
	setDigestMatcherType(second, enumerations.DigestMatcherType_MANIFEST_ENTRY)
	second.Uri = ptr("other.xml")
	second.DocumentName = ptr("renamed.xml")
	return []*diagnosticjaxb.XmlDigestMatcher{first, second}
}

func setDigestMatcherType(digestMatcher *diagnosticjaxb.XmlDigestMatcher, t enumerations.DigestMatcherType) {
	value := diagnosticjaxb.DigestMatcherTypeValue(t)
	digestMatcher.Type = &value
}

// syntheticSignature builds the token the oracle drives SignatureIntactCheck
// with in the CERTIFICATE-context scenario: a signature with no basic-signature
// element, hence not intact.
func syntheticSignature() *diagnostic.SignatureWrapper {
	xmlSignature := &diagnosticjaxb.XmlSignature{}
	xmlSignature.Id = diagnosticjaxb.NewCollapsedString("S-SYNTHETIC")
	return diagnostic.NewSignatureWrapper(xmlSignature)
}

// corpusSignature loads the signature of the corpus dump the oracle drove
// SignatureIntactWithIdCheck with.
func corpusSignature(t *testing.T, file, id string) *diagnostic.SignatureWrapper {
	t.Helper()
	for _, signature := range loadDiagnosticData(t, file).Signatures() {
		if signature.Id() == id {
			return signature
		}
	}
	t.Fatalf("signature %s not found in %s", id, file)
	return nil
}

func ptr[T any](value T) *T { return &value }
