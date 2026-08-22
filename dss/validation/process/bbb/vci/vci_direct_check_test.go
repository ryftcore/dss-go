package vci

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// The direct-check KAT machinery: testdata/oracle/vci_direct.jsonl holds, for
// each scenario named below, the XmlVCI upstream produces when the check is run
// alone at Level.FAIL over a synthetic signature policy - see
// ../testdata/gen/BbbBlocksOracle.java. The corpus KAT never drives a VCI check
// to a failure, so both paths of every check are covered here.

// singleCheckChain is a chain of exactly one item, mirroring the oracle's
// SingleVCIChain.
type singleCheckChain struct {
	*process.ChainBase[*jaxb.XmlVCI]
	factory func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI]
}

func newSingleCheckChain(
	factory func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI],
) *singleCheckChain {
	xmlVCI := &jaxb.XmlVCI{}
	c := &singleCheckChain{
		ChainBase: process.NewChainBase(i18nProviderForTests, process.NewResult(xmlVCI,
			&xmlVCI.XmlConstraintsConclusionContent, &xmlVCI.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleCheckChain) InitChain() {
	c.FirstItem = c.factory(c.Result, process.GetLevelRule(enumerations.LevelFail))
}

// assertDirectRow runs the single-item chain and compares it against the named
// row of the direct oracle.
func assertDirectRow(t *testing.T, name string,
	factory func(result *process.Result[*jaxb.XmlVCI], rule policy.LevelRule) process.ChainItem[*jaxb.XmlVCI]) {
	t.Helper()
	rows := loadRows(t, corpustest.Path(t, "oracle/vci_direct.jsonl"))
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
	actual := toRow("synthetic", name, enumerations.ContextSignature, "VCI",
		&result.XmlConstraintsConclusionContent, result.Title)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: mismatch\nexpected: %s\nactual:   %s", name, mustJSON(t, expected), mustJSON(t, actual))
	}
}

// anyPolicyRule is the MultiValuesRule the oracle drives
// SignaturePolicyIdentifierCheck with.
type anyPolicyRule struct{}

func (anyPolicyRule) Level() enumerations.Level { return enumerations.LevelFail }
func (anyPolicyRule) Values() []string          { return []string{"ANY_POLICY"} }

// policySignature builds a signature carrying a signature policy with the
// requested flags, mirroring the oracle's builder. A nil policyId is Java's null
// Id.
func policySignature(policyId *string, identified, digestMatch, zeroHash *bool,
	storePresent bool) *diagnostic.SignatureWrapper {
	xmlSignature := &diagnosticjaxb.XmlSignature{}
	xmlSignature.Id = diagnosticjaxb.NewCollapsedString("S-SYNTHETIC-POLICY")

	digestMethod := diagnosticjaxb.DigestAlgorithmValue(enumerations.DigestAlgorithmSHA256)
	digestValue := diagnosticjaxb.Base64Binary{7, 8, 9}
	digest := &diagnosticjaxb.XmlPolicyDigestAlgoAndValue{}
	digest.DigestMethod = &digestMethod
	digest.DigestValue = &digestValue
	digest.Match = digestMatch
	digest.ZeroHash = zeroHash

	xmlSignature.Policy = &diagnosticjaxb.XmlPolicy{
		Id:                 policyId,
		Identified:         identified,
		DigestAlgoAndValue: digest,
	}

	if storePresent {
		store := &diagnosticjaxb.XmlSignaturePolicyStore{}
		store.Id = ptr("SPS-SYNTHETIC")
		xmlSignature.SignaturePolicyStore = store
	}
	return diagnostic.NewSignatureWrapper(xmlSignature)
}

func ptr[T any](value T) *T { return &value }
