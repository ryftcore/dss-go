package isc

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

// The direct-check KAT machinery: testdata/oracle/isc_direct.jsonl holds, for
// each scenario named below, the XmlISC upstream produces when the check is run
// alone at Level.FAIL over a synthetic signature - see
// ../testdata/gen/BbbBlocksOracle.java. The corpus KAT
// (identification_of_the_signing_certificate_test.go) never exercises the
// failure path of the three reference checks, so those are covered here.

// singleCheckChain is a chain of exactly one item, mirroring the oracle's
// SingleISCChain.
type singleCheckChain struct {
	*process.ChainBase[*jaxb.XmlISC]
	factory func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC]
}

func newSingleCheckChain(
	factory func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC],
) *singleCheckChain {
	xmlISC := &jaxb.XmlISC{}
	c := &singleCheckChain{
		ChainBase: process.NewChainBase(i18nProviderForTests, process.NewResult(xmlISC,
			&xmlISC.XmlConstraintsConclusionContent, &xmlISC.XmlConstraintsConclusionAttrs)),
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
	factory func(result *process.Result[*jaxb.XmlISC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlISC]) {
	t.Helper()
	rows := loadRows(t, corpustest.Path(t, "oracle/isc_direct.jsonl"))
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
	actual := toRow("synthetic", name, enumerations.Context_SIGNATURE, "ISC",
		&result.XmlConstraintsConclusionContent, result.Title)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: mismatch\nexpected: %s\nactual:   %s", name, mustJSON(t, expected), mustJSON(t, actual))
	}
}

// syntheticSignature builds a signature carrying one signing-certificate
// reference with the requested flags, mirroring the oracle's builder.
func syntheticSignature(digestPresent, digestMatch, issuerSerialPresent, issuerSerialMatch,
	withSigningCertificate bool) *diagnostic.SignatureWrapper {
	xmlSignature := &diagnosticjaxb.XmlSignature{}
	xmlSignature.Id = diagnosticjaxb.NewCollapsedString("S-SYNTHETIC")

	xmlCertificate := &diagnosticjaxb.XmlCertificate{}
	xmlCertificate.Id = diagnosticjaxb.NewCollapsedString("C-SYNTHETIC")

	origin := diagnosticjaxb.CertificateRefOriginValue(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	ref := &diagnosticjaxb.XmlCertificateRef{Origin: &origin}
	if digestPresent {
		digestMethod := diagnosticjaxb.DigestAlgorithmValue(enumerations.DigestAlgorithm_SHA256)
		digestValue := diagnosticjaxb.Base64Binary{1, 2, 3}
		match := digestMatch
		digest := &diagnosticjaxb.XmlDigestAlgoAndValue{}
		digest.DigestMethod = &digestMethod
		digest.DigestValue = &digestValue
		digest.Match = &match
		ref.DigestAlgoAndValue = digest
	}
	if issuerSerialPresent {
		match := issuerSerialMatch
		ref.IssuerSerial = &diagnosticjaxb.XmlIssuerSerial{
			Value: diagnosticjaxb.Base64Binary{4, 5, 6},
			Match: &match,
		}
	}

	related := &diagnosticjaxb.XmlRelatedCertificate{Certificate: xmlCertificate}
	related.CertificateRef = []*diagnosticjaxb.XmlCertificateRef{ref}
	xmlSignature.FoundCertificates = &diagnosticjaxb.XmlFoundCertificates{
		RelatedCertificate: []*diagnosticjaxb.XmlRelatedCertificate{related},
	}

	if withSigningCertificate {
		xmlSignature.SigningCertificate = &diagnosticjaxb.XmlSigningCertificate{Certificate: xmlCertificate}
	}
	return diagnostic.NewSignatureWrapper(xmlSignature)
}
