package xcv

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// The rac check KAT: every row of testdata/oracle/xcva_direct.jsonl is the XmlRAC
// upstream produces when one of the twelve
// eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks classes is run alone,
// at Level.FAIL, over one (certificate, certificate revocation data) pair of the
// marshal-parity corpus or of the synthetic dumps in testdata/dd - see
// testdata/gen/XcvaDirectOracle.java. This test replays the same pairs.

// xcvaFailLevel is the Level.FAIL rule the direct oracle drove every check with.
var xcvaFailLevel = process.GetLevelRule(enumerations.Level_FAIL)

// singleRACChain is a chain of exactly one item over an XmlRAC, the Go form of
// XcvaDirectOracle.SingleRACChain.
type singleRACChain struct {
	*process.ChainBase[*jaxb.XmlRAC]
	factory func(result *process.Result[*jaxb.XmlRAC], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRAC]
}

func newSingleRACChain(factory func(result *process.Result[*jaxb.XmlRAC],
	constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRAC]) *singleRACChain {
	xmlRAC := &jaxb.XmlRAC{}
	c := &singleRACChain{
		ChainBase: process.NewChainBase(xcvaI18n(), process.NewResult(xmlRAC,
			&xmlRAC.XmlConstraintsConclusionContent, &xmlRAC.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleRACChain) InitChain() {
	c.FirstItem = c.factory(c.Result, xcvaFailLevel)
}

func TestRacChecksAgainstJavaOracle(t *testing.T) {
	rows := loadXcvaRows(t, corpustest.Path(t, "oracle/xcva_direct.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}

	// index: file -> token -> check
	byFile := map[string]map[string]map[string]*xcvaRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]map[string]*xcvaRow{}
		}
		if _, seen := byFile[row.File][row.Token]; !seen {
			byFile[row.File][row.Token] = map[string]*xcvaRow{}
		}
		byFile[row.File][row.Token][row.Check] = row
	}

	statuses := map[string]map[string]int{}
	matched := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadXcvaDiagnosticData(t, name)
			expected := byFile[name]

			compare := func(token, check string,
				factory func(result *process.Result[*jaxb.XmlRAC], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRAC]) {
				want, ok := expected[token][check]
				if !ok {
					return // the Java oracle threw on this input and recorded no row
				}
				if want.Title == nil {
					// Known mapping (same as the phase 8c direct corpora): a
					// single-item chain defines no title MessageTag, so Java leaves
					// the attribute null where the generated non-pointer Go member
					// spells it "". This one field is normalised and nothing else.
					empty := ""
					want.Title = &empty
				}
				var result *jaxb.XmlRAC
				if xcvaSafeExecute(func() { result = newSingleRACChain(factory).Execute() }) {
					t.Errorf("%s / %s: Go panicked where Java produced a row", token, check)
					return
				}
				got := &xcvaRow{File: name, Token: token, Check: check,
					Context: want.Context, Block: want.Block}
				toXcvaBody(&got.xcvaNode, &result.XmlConstraintsConclusionContent, result.Title)
				matched++
				if statuses[check] == nil {
					statuses[check] = map[string]int{}
				}
				for _, constraint := range want.Constraints {
					if constraint.Status != nil {
						statuses[check][*constraint.Status]++
					}
				}
				if !reflect.DeepEqual(want, got) {
					t.Errorf("%s / %s: mismatch\nexpected: %s\nactual:   %s",
						token, check, mustXcvaJSON(t, want), mustXcvaJSON(t, got))
				}
			}

			for _, certificate := range diagnosticData.UsedCertificates() {
				cert := certificate
				compare(cert.Id(), "RevocationIssuerRevocationDataAvailableCheck",
					func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
						return NewRevocationIssuerRevocationDataAvailableCheck(xcvaI18n(), r, cert, l)
					})

				for _, revocation := range cert.CertificateRevocationData() {
					rev := revocation
					base := &rev.RevocationWrapper
					token := cert.Id() + "|" + rev.Id()

					compare(token, "RevocationDataKnownCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationDataKnownCheck(xcvaI18n(), r, rev, l)
						})
					compare(token, "RevocationIssuerKnownCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationIssuerKnownCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "ThisUpdatePresenceCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewThisUpdatePresenceCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "RevocationIssuerValidAtProductionTimeCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationIssuerValidAtProductionTimeCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "RevocationResponderIdMatchCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationResponderIdMatchCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "RevocationCertHashPresenceCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationCertHashPresenceCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "RevocationCertHashMatchCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationCertHashMatchCheck(xcvaI18n(), r, base, l)
						})
					compare(token, "SelfIssuedOCSPCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewSelfIssuedOCSPCheck(xcvaI18n(), r, cert, base, l)
						})
					compare(token, "RevocationAfterCertificateIssuanceCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationAfterCertificateIssuanceCheck(xcvaI18n(), r, cert, base, l)
						})
					compare(token, "RevocationHasInformationAboutCertificateCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationHasInformationAboutCertificateCheck(xcvaI18n(), r, cert, base, l)
						})

					// The oracle feeds this one the XmlRAC a real
					// RevocationAcceptanceChecker run produced for the same pair.
					var racResult *jaxb.XmlRAC
					if xcvaSafeExecute(func() {
						racResult = NewRevocationAcceptanceChecker(xcvaI18n(), cert, rev, xcvaCurrentTime,
							xcvaDefaultPolicy(t), map[string]struct{}{}).Execute()
					}) {
						continue
					}
					compare(token, "RevocationAcceptanceCheckerResultCheck",
						func(r *process.Result[*jaxb.XmlRAC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRAC] {
							return NewRevocationAcceptanceCheckerResultCheck(xcvaI18n(), r, racResult, l)
						})
				}
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}

	// Every check class must keep both an OK and a NOT OK row, so that a
	// regression cannot hide behind a corpus that only ever exercises one side.
	var checks []string
	for check := range statuses {
		checks = append(checks, check)
	}
	sort.Strings(checks)
	for _, check := range checks {
		if statuses[check]["OK"] == 0 || statuses[check]["NOT OK"] == 0 {
			t.Errorf("check %s lost its OK/NOT OK pair: %v", check, statuses[check])
		}
	}
	if len(checks) != 12 {
		t.Errorf("expected 12 rac check classes in the corpus, found %d: %v", len(checks), checks)
	}
}
