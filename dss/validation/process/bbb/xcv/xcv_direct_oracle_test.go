package xcv

import (
	"bufio"
	"encoding/json"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// The 75-check KAT: every row of testdata/oracle/xcv_direct.jsonl is the
// XmlSubXCV/XmlXCV/XmlRFC upstream produces when one of the 75 directly
// instantiable eu.europa.esig.dss.validation.process.bbb.xcv.checks /
// sub.checks / sub.checks.pseudo / rfc.checks classes of the phase 8d XCVB
// manifest is run alone, at Level.FAIL, through a chain of exactly one item -
// see testdata/gen/XcvOracle.java. This test replays the same drive.
//
// Unlike xcva_direct_oracle_test.go's twelve rac/checks classes (all sharing
// XmlRAC), these 75 classes are split across three result types
// (XmlSubXCV/XmlXCV/XmlRFC), so three small single-item chains stand in for
// XcvOracle's SingleSubXCVChain/SingleXCVChain/SingleRFCChain.

// -------------------------------------------------------------- row shape

// xcvDirectRow is one line of testdata/oracle/xcv_direct.jsonl: just the
// title/conclusion/constraints body XcvOracle's row() writes (no nesting - the
// 83-file manifest's checks/sub.checks/rfc.checks classes produce a single
// flat ConstraintsConclusion, never a tree).
type xcvDirectRow struct {
	File        string            `json:"file"`
	Token       string            `json:"token"`
	Check       string            `json:"check"`
	Block       string            `json:"block"`
	Title       *string           `json:"title"`
	Conclusion  *xcvaConclusion   `json:"conclusion"`
	Constraints []*xcvaConstraint `json:"constraints"`
}

func loadXcvDirectRows(t *testing.T, path string) []*xcvDirectRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*xcvDirectRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<25)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &xcvDirectRow{}
		if err := json.Unmarshal(scanner.Bytes(), row); err != nil {
			t.Fatalf("parse oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	return rows
}

// fillXcvDirectBody fills title/conclusion/constraints from a produced
// ConstraintsConclusion, the way toXcvaBody does for the nested xcvaNode.
func fillXcvDirectBody(row *xcvDirectRow, content *jaxb.XmlConstraintsConclusionContent, title string) {
	titleValue := title
	row.Title = &titleValue
	row.Conclusion = toXcvaConclusion(content.Conclusion)
	row.Constraints = toXcvaConstraints(content.Constraint)
}

// ------------------------------------------------------------- single chains

// singleSubXCVChain is a chain of exactly one item over an XmlSubXCV.
type singleSubXCVChain struct {
	*process.ChainBase[*jaxb.XmlSubXCV]
	factory func(result *process.Result[*jaxb.XmlSubXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV]
}

func newSingleSubXCVChain(factory func(result *process.Result[*jaxb.XmlSubXCV],
	constraint policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV]) *singleSubXCVChain {
	xmlSubXCV := &jaxb.XmlSubXCV{}
	c := &singleSubXCVChain{
		ChainBase: process.NewChainBase(xcvaI18n(), process.NewResult(xmlSubXCV,
			&xmlSubXCV.XmlConstraintsConclusionContent, &xmlSubXCV.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleSubXCVChain) InitChain() { c.FirstItem = c.factory(c.Result, xcvaFailLevel) }

// singleXCVDirectChain is a chain of exactly one item over an XmlXCV.
type singleXCVDirectChain struct {
	*process.ChainBase[*jaxb.XmlXCV]
	factory func(result *process.Result[*jaxb.XmlXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlXCV]
}

func newSingleXCVDirectChain(factory func(result *process.Result[*jaxb.XmlXCV],
	constraint policy.LevelRule) process.ChainItem[*jaxb.XmlXCV]) *singleXCVDirectChain {
	xmlXCV := &jaxb.XmlXCV{}
	c := &singleXCVDirectChain{
		ChainBase: process.NewChainBase(xcvaI18n(), process.NewResult(xmlXCV,
			&xmlXCV.XmlConstraintsConclusionContent, &xmlXCV.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleXCVDirectChain) InitChain() { c.FirstItem = c.factory(c.Result, xcvaFailLevel) }

// singleRFCChain is a chain of exactly one item over an XmlRFC.
type singleRFCChain struct {
	*process.ChainBase[*jaxb.XmlRFC]
	factory func(result *process.Result[*jaxb.XmlRFC], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRFC]
}

func newSingleRFCChain(factory func(result *process.Result[*jaxb.XmlRFC],
	constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRFC]) *singleRFCChain {
	xmlRFC := &jaxb.XmlRFC{}
	c := &singleRFCChain{
		ChainBase: process.NewChainBase(xcvaI18n(), process.NewResult(xmlRFC,
			&xmlRFC.XmlConstraintsConclusionContent, &xmlRFC.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleRFCChain) InitChain() { c.FirstItem = c.factory(c.Result, xcvaFailLevel) }

// ----------------------------------------------------------- literal rules

type literalMultiValuesRule struct{ value string }

func (r literalMultiValuesRule) Level() enumerations.Level { return enumerations.LevelFail }
func (r literalMultiValuesRule) Values() []string          { return []string{r.value} }

type literalNumericValueRule struct{ value float64 }

func (r literalNumericValueRule) Level() enumerations.Level { return enumerations.LevelFail }
func (r literalNumericValueRule) Value() float64            { return r.value }

type literalValueRule struct{ value string }

func (r literalValueRule) Level() enumerations.Level { return enumerations.LevelFail }
func (r literalValueRule) Value() string             { return r.value }

type literalDurationRule struct{ millis int64 }

func (r literalDurationRule) Level() enumerations.Level { return enumerations.LevelFail }
func (r literalDurationRule) Duration() int64           { return r.millis }

type literalApplicabilityRule struct{}

func (r literalApplicabilityRule) Level() enumerations.Level                     { return enumerations.LevelFail }
func (r literalApplicabilityRule) CertificateExtensions() policy.MultiValuesRule { return nil }
func (r literalApplicabilityRule) CertificatePolicies() policy.MultiValuesRule   { return nil }

var (
	xcvDirectAny           = literalMultiValuesRule{"*"}
	xcvDirectNone          = literalMultiValuesRule{"NO-SUCH-VALUE"}
	xcvDirectNumLow        = literalNumericValueRule{0}
	xcvDirectNumHigh       = literalNumericValueRule{999999999}
	xcvDirectValueAny      = literalValueRule{"*"}
	xcvDirectValueNone     = literalValueRule{"NO-SUCH-VALUE"}
	xcvDirectAppliesNever  = literalApplicabilityRule{}
	xcvDirectDurationTight = literalDurationRule{1000}
	xcvDirectDurationLoose = literalDurationRule{1000 * 60 * 60 * 24 * 365 * 50}
)

// xcvDirectCurrentTime is the fixed validation/current time XcvOracle used:
// 2024-01-01T00:00:00Z (same instant as xcvaCurrentTime).
var xcvDirectCurrentTime = xcvaCurrentTime

// -------------------------------------------------------------- the test

func TestXcvChecksAgainstJavaOracle(t *testing.T) {
	rows := loadXcvDirectRows(t, corpustest.Path(t, "oracle/xcv_direct.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}

	// index: file -> token -> check -> rows, in file order. Almost every
	// (file, token, check) triple is unique, but RevocationDataFreshCheck is
	// driven twice per token (a loose then a tight DurationRule, to produce
	// both an OK and a NOT OK row over the very same revocation - see
	// XcvOracle's emitDump), so the two rows share a key and must both
	// survive; every compare* call below consumes one entry off the front of
	// its list, in the same loose-then-tight order Java wrote them.
	byFile := map[string]map[string]map[string][]*xcvDirectRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]map[string][]*xcvDirectRow{}
		}
		if _, seen := byFile[row.File][row.Token]; !seen {
			byFile[row.File][row.Token] = map[string][]*xcvDirectRow{}
		}
		byFile[row.File][row.Token][row.Check] = append(byFile[row.File][row.Token][row.Check], row)
	}

	statuses := map[string]map[string]int{}
	matched := 0

	// popExpected removes and returns the first row queued for check, so a
	// check driven more than once per token (RevocationDataFreshCheck) is
	// matched against successive rows in the order Java wrote them.
	popExpected := func(expected map[string][]*xcvDirectRow, check string) (*xcvDirectRow, bool) {
		queue := expected[check]
		if len(queue) == 0 {
			return nil, false
		}
		expected[check] = queue[1:]
		return queue[0], true
	}

	compareSub := func(t *testing.T, expected map[string][]*xcvDirectRow, file, token, check string,
		factory func(result *process.Result[*jaxb.XmlSubXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV]) {
		t.Helper()
		want, ok := popExpected(expected, check)
		if !ok {
			return // the Java oracle threw on this input and recorded no row
		}
		var result *jaxb.XmlSubXCV
		if xcvaSafeExecute(func() { result = newSingleSubXCVChain(factory).Execute() }) {
			t.Errorf("%s / %s: Go panicked where Java produced a row", token, check)
			return
		}
		got := &xcvDirectRow{File: file, Token: token, Check: check, Block: want.Block}
		fillXcvDirectBody(got, &result.XmlConstraintsConclusionContent, result.Title)
		recordMatch(statuses, &matched, want, got, token, check, t)
	}

	compareXCV := func(t *testing.T, expected map[string][]*xcvDirectRow, file, token, check string,
		factory func(result *process.Result[*jaxb.XmlXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlXCV]) {
		t.Helper()
		want, ok := popExpected(expected, check)
		if !ok {
			return
		}
		var result *jaxb.XmlXCV
		if xcvaSafeExecute(func() { result = newSingleXCVDirectChain(factory).Execute() }) {
			t.Errorf("%s / %s: Go panicked where Java produced a row", token, check)
			return
		}
		got := &xcvDirectRow{File: file, Token: token, Check: check, Block: want.Block}
		fillXcvDirectBody(got, &result.XmlConstraintsConclusionContent, result.Title)
		recordMatch(statuses, &matched, want, got, token, check, t)
	}

	compareRFC := func(t *testing.T, expected map[string][]*xcvDirectRow, file, token, check string,
		factory func(result *process.Result[*jaxb.XmlRFC], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlRFC]) {
		t.Helper()
		want, ok := popExpected(expected, check)
		if !ok {
			return
		}
		var result *jaxb.XmlRFC
		if xcvaSafeExecute(func() { result = newSingleRFCChain(factory).Execute() }) {
			t.Errorf("%s / %s: Go panicked where Java produced a row", token, check)
			return
		}
		got := &xcvDirectRow{File: file, Token: token, Check: check, Block: want.Block}
		fillXcvDirectBody(got, &result.XmlConstraintsConclusionContent, result.Title)
		recordMatch(statuses, &matched, want, got, token, check, t)
	}

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			expectedByToken := byFile[name]

			var diagnosticData *diagnostic.DiagnosticData
			if name != "synthetic" {
				diagnosticData = loadXcvaDiagnosticData(t, name)
			}

			compareCertificateChecks := func(certificate *diagnostic.CertificateWrapper) {
				id := certificate.Id()
				expected := expectedByToken[id]
				if expected == nil {
					expected = map[string][]*xcvDirectRow{}
				}

				compareSub(t, expected, name, id, "BasicConstraintsCACheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewBasicConstraintsCACheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "BasicConstraintsMaxPathLengthCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewBasicConstraintsMaxPathLengthCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "AuthorityInfoAccessPresentCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewAuthorityInfoAccessPresentCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "AuthorityKeyIdentifierPresentCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewAuthorityKeyIdentifierPresentCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "SubjectKeyIdentifierPresentCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewSubjectKeyIdentifierPresentCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "RevocationInfoAccessPresentCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewRevocationInfoAccessPresentCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "NoRevAvailCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewNoRevAvailCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateSelfSignedCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateSelfSignedCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateNotSelfSignedCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateNotSelfSignedCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateSignatureValidCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateSignatureValidCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateIssuerNameCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateIssuerNameCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "OtherTrustAnchorExistsCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewOtherTrustAnchorExistsCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificatePolicyTreeCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificatePolicyTreeCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateNameConstraintsCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateNameConstraintsCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "SerialNumberCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewSerialNumberCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "PseudoUsageCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewPseudoUsageCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificatePolicyQualifiedIdsCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificatePolicyQualifiedIdsCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificatePolicySupportedByQSCDIdsCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificatePolicySupportedByQSCDIdsCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateQcComplianceCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateQcComplianceCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateQcSSCDCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateQcSSCDCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateIssuedToNaturalPersonCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateIssuedToNaturalPersonCheck(xcvaI18n(), r, certificate, l)
				})
				compareSub(t, expected, name, id, "CertificateIssuedToLegalPersonCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateIssuedToLegalPersonCheck(xcvaI18n(), r, certificate, l)
				})

				for _, variant := range []string{"any", "none"} {
					mv := xcvDirectAny
					if variant == "none" {
						mv = xcvDirectNone
					}
					compareSub(t, expected, name, id, "CommonNameCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCommonNameCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CountryCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCountryCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "EmailCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewEmailCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "GivenNameCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewGivenNameCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "LocalityCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewLocalityCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "OrganizationIdentifierCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewOrganizationIdentifierCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "OrganizationNameCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewOrganizationNameCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "OrganizationUnitCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewOrganizationUnitCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "StateCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewStateCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "SurnameCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewSurnameCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "TitleCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewTitleCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "PseudonymCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewPseudonymCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificatePolicyIdsCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificatePolicyIdsCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateForbiddenExtensionsCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateForbiddenExtensionsCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateSupportedCriticalExtensionsCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateSupportedCriticalExtensionsCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcEuPDSLocationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcEuPDSLocationCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcTypeCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcTypeCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcCCLegislationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcCCLegislationCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcQSCDLegislationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcQSCDLegislationCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcIdentificationMethodCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcIdentificationMethodCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateSemanticsIdentifierCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateSemanticsIdentifierCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificatePS2DQcCompetentAuthorityIdCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificatePS2DQcCompetentAuthorityIdCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificatePS2DQcCompetentAuthorityNameCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificatePS2DQcCompetentAuthorityNameCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificatePS2DQcRolesOfPSPCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificatePS2DQcRolesOfPSPCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcPSBCountryOfLegislationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcPSBCountryOfLegislationCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcPSBAuthSourceIdentificationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcPSBAuthSourceIdentificationCheck(xcvaI18n(), r, certificate, mv)
					})
					compareSub(t, expected, name, id, "CertificateQcPSBLegislationIdentificationCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcPSBLegislationIdentificationCheck(xcvaI18n(), r, certificate, mv)
					})
				}

				for _, variant := range []string{"low", "high"} {
					nv := xcvDirectNumLow
					if variant == "high" {
						nv = xcvDirectNumHigh
					}
					compareSub(t, expected, name, id, "CertificateMinQcEuRetentionPeriodCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateMinQcEuRetentionPeriodCheck(xcvaI18n(), r, certificate, nv)
					})
					compareSub(t, expected, name, id, "CertificateMinQcTransactionLimitCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateMinQcTransactionLimitCheck(xcvaI18n(), r, certificate, nv)
					})
				}
				for _, variant := range []string{"any", "none"} {
					vr := xcvDirectValueAny
					if variant == "none" {
						vr = xcvDirectValueNone
					}
					compareSub(t, expected, name, id, "CertificateQcEuLimitValueCurrencyCheck-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateQcEuLimitValueCurrencyCheck(xcvaI18n(), r, certificate, vr)
					})
				}

				for _, subContext := range []enumerations.SubContext{enumerations.SubContextSigningCert, enumerations.SubContextCACertificate} {
					sc := subContext
					suffix := string(sc)
					for _, variant := range []string{"any", "none"} {
						mv := xcvDirectAny
						if variant == "none" {
							mv = xcvDirectNone
						}
						compareSub(t, expected, name, id, "KeyUsageCheck-"+suffix+"-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
							return NewKeyUsageCheck(xcvaI18n(), r, certificate, enumerations.ContextSignature, sc, mv)
						})
						compareSub(t, expected, name, id, "ExtendedKeyUsageCheck-"+suffix+"-"+variant, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
							return NewExtendedKeyUsageCheck(xcvaI18n(), r, certificate, enumerations.ContextSignature, sc, mv)
						})
					}
					compareSub(t, expected, name, id, "CertificateValidationBeforeSunsetDateCheck-"+suffix, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateValidationBeforeSunsetDateCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, xcvDirectCurrentTime, l)
					})
					compareSub(t, expected, name, id, "CertificateValidationBeforeSunsetDateWithIdCheck-"+suffix, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateValidationBeforeSunsetDateWithIdCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, xcvDirectCurrentTime, l)
					})
					compareSub(t, expected, name, id, "RevocationDataRequiredCheck-"+suffix, func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewRevocationDataRequiredCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, xcvDirectCurrentTime, xcvaFailLevel, xcvDirectAppliesNever)
					})
				}

				compareSub(t, expected, name, id, "CertificateValidityRangeCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewCertificateValidityRangeCheckMinimal[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, xcvDirectCurrentTime, l)
				})

				compareXCV(t, expected, name, id, "ProspectiveCertificateChainCheck", func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
					return NewProspectiveCertificateChainCheck[*jaxb.XmlXCV](xcvaI18n(), r, certificate, enumerations.ContextSignature, l)
				})
				compareXCV(t, expected, name, id, "ProspectiveCertificateChainAtValidationTimeCheck", func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
					return NewProspectiveCertificateChainAtValidationTimeCheck(xcvaI18n(), r, certificate, xcvDirectCurrentTime, l)
				})

				if usageTime := certificate.NotBefore(); usageTime != nil {
					for _, variant := range []string{"any", "none"} {
						mv := xcvDirectAny
						if variant == "none" {
							mv = xcvDirectNone
						}
						compareXCV(t, expected, name, id, "TrustServiceStatusCheck-"+variant, func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
							return NewTrustServiceStatusCheck(xcvaI18n(), r, certificate, usageTime, enumerations.ContextSignature, mv)
						})
						compareXCV(t, expected, name, id, "TrustServiceTypeIdentifierCheck-"+variant, func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
							return NewTrustServiceTypeIdentifierCheck(xcvaI18n(), r, certificate, usageTime, enumerations.ContextSignature, mv)
						})
					}
				}

				// --- checks consuming a synthetic wrapper result keyed to this real certificate.
				for _, shape := range []string{"passed", "passed-with-algo", "error", "warning"} {
					aov := aovOfShape(shape, id)
					shapeToken := id + "/" + shape
					shapeExpected := expectedByToken[shapeToken]
					if shapeExpected == nil {
						shapeExpected = map[string][]*xcvDirectRow{}
					}
					compareSub(t, shapeExpected, name, shapeToken, "CertificateAlgorithmObsolescenceValidationCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateAlgorithmObsolescenceValidationCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, aov, xcvDirectCurrentTime, i18n.MessageTag_SIGNING_CERTIFICATE, id)
					})
				}
				for _, indication := range []enumerations.Indication{enumerations.IndicationPassed, enumerations.IndicationIndeterminate} {
					indToken := id + "/" + string(indication)
					indExpected := expectedByToken[indToken]
					if indExpected == nil {
						indExpected = map[string][]*xcvDirectRow{}
					}

					crs := crsOfIndication(indication, id)
					compareSub(t, indExpected, name, indToken, "CertificateRevocationSelectorResultCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateRevocationSelectorResultCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, crs, l)
					})

					rfc := rfcOfIndication(indication, id)
					compareSub(t, indExpected, name, indToken, "RevocationFreshnessCheckerResultCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewRevocationFreshnessCheckerResultCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, rfc, l)
					})

					subXCV := subXCVOfIndication(indication, id)
					compareXCV(t, indExpected, name, indToken, "CheckSubXCVResult", func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
						return NewCheckSubXCVResult(xcvaI18n(), r, subXCV, l)
					})
				}
			}

			compareRevocationChecks := func(certificate *diagnostic.CertificateWrapper, revocation *diagnostic.CertificateRevocationWrapper) {
				token := certificate.Id() + "|" + revocation.Id()
				expected := expectedByToken[token]
				if expected == nil {
					expected = map[string][]*xcvDirectRow{}
				}
				base := &revocation.RevocationWrapper

				compareRFC(t, expected, name, token, "NextUpdateCheck", func(r *process.Result[*jaxb.XmlRFC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRFC] {
					return NewNextUpdateCheck(xcvaI18n(), r, base, l)
				})
				// RevocationDataFreshCheck is driven twice per token in
				// XcvOracle (a loose then a tight DurationRule), so it queues
				// two rows under the same check key; popExpected consumes them
				// loose-first, matching the order Java wrote them.
				compareRFC(t, expected, name, token, "RevocationDataFreshCheck", func(r *process.Result[*jaxb.XmlRFC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRFC] {
					return NewRevocationDataFreshCheck(xcvaI18n(), r, base, xcvDirectCurrentTime, xcvDirectDurationLoose)
				})
				compareRFC(t, expected, name, token, "RevocationDataFreshCheck", func(r *process.Result[*jaxb.XmlRFC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRFC] {
					return NewRevocationDataFreshCheck(xcvaI18n(), r, base, xcvDirectCurrentTime, xcvDirectDurationTight)
				})
				compareRFC(t, expected, name, token, "RevocationDataFreshCheckWithNullConstraint", func(r *process.Result[*jaxb.XmlRFC], l policy.LevelRule) process.ChainItem[*jaxb.XmlRFC] {
					return NewRevocationDataFreshCheckWithNullConstraint(xcvaI18n(), r, base, xcvDirectCurrentTime, l)
				})
				compareSub(t, expected, name, token, "AcceptableRevocationDataAvailableCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewAcceptableRevocationDataAvailableCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, base, l)
				})
				compareSub(t, expected, name, token, "RevocationIssuerTrustedCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewRevocationIssuerTrustedCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, base.SigningCertificate(), xcvDirectCurrentTime, xcvaFailLevel, l)
				})
				compareSub(t, expected, name, token, "RevocationIssuerValidityRangeCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewRevocationIssuerValidityRangeCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, base, xcvDirectCurrentTime, l)
				})

				for _, subContext := range []enumerations.SubContext{enumerations.SubContextSigningCert, enumerations.SubContextCACertificate} {
					sc := subContext
					subToken := token + "|" + string(sc)
					subExpected := expectedByToken[subToken]
					if subExpected == nil {
						subExpected = map[string][]*xcvDirectRow{}
					}
					compareSub(t, subExpected, name, subToken, "CertificateNotOnHoldCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateNotOnHoldCheck(xcvaI18n(), r, revocation, xcvDirectCurrentTime, l)
					})
					compareSub(t, subExpected, name, subToken, "CertificateNotRevokedCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
						return NewCertificateNotRevokedCheck(xcvaI18n(), r, revocation, xcvDirectCurrentTime, l, sc)
					})
				}
			}

			if name == "synthetic" {
				runSyntheticChecks(t, expectedByToken, compareSub, compareXCV)
				return
			}

			for _, certificate := range diagnosticData.UsedCertificates() {
				compareCertificateChecks(certificate)

				for _, revocation := range certificate.CertificateRevocationData() {
					compareRevocationChecks(certificate, revocation)
				}

				id := certificate.Id()
				expected := expectedByToken[id]
				if expected == nil {
					expected = map[string][]*xcvDirectRow{}
				}
				compareSub(t, expected, name, id, "AcceptableRevocationDataAvailableCheck-none", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewAcceptableRevocationDataAvailableCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, nil, l)
				})
				compareSub(t, expected, name, id, "RevocationDataAvailableCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
					return NewRevocationDataAvailableCheck[*jaxb.XmlSubXCV](xcvaI18n(), r, certificate, l)
				})
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}

	// Every one of the 75 checks must keep both an OK and a NOT OK row.
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
	if len(checks) != 75 {
		t.Errorf("expected 75 xcv check classes in the corpus, found %d: %v", len(checks), checks)
	}
}

// recordMatch compares got against want, recording the OK/NOT OK statuses the
// row's constraints carry (bucketed by the underlying check class, not the
// -any/-none/-low/-high/-SIGNING_CERT/-CA_CERTIFICATE variant string XcvOracle
// suffixes the "check" field with) and failing t if the two bodies differ.
func recordMatch(statuses map[string]map[string]int, matched *int, want, got *xcvDirectRow, token, check string, t *testing.T) {
	t.Helper()
	*matched++
	class := xcvCheckClass(check)
	if statuses[class] == nil {
		statuses[class] = map[string]int{}
	}
	for _, constraint := range want.Constraints {
		if constraint.Status != nil {
			statuses[class][*constraint.Status]++
		}
	}
	if want.Title == nil {
		// Known mapping (same as the phase 8c direct corpora and XCVA's
		// xcva_direct_oracle_test.go): a single-item chain defines no title
		// MessageTag, so Java leaves the attribute null where the generated
		// non-pointer Go member spells it "".
		empty := ""
		want.Title = &empty
	}
	want.File, want.Token, want.Check = got.File, got.Token, got.Check
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s / %s: mismatch\nexpected: %s\nactual:   %s",
			token, check, mustXcvaJSON(t, want), mustXcvaJSON(t, got))
	}
}

// xcvDirectVariantSuffixes are the trailing "-"-joined tokens XcvOracle
// appends to a check's Java class name to distinguish the ANY/NONE,
// low/high or SIGNING_CERT/CA_CERTIFICATE row it drove; stripped to recover
// the underlying check class for the OK/NOT OK pairing invariant (every one
// of the 75 check *classes*, not every variant string, must see both).
var xcvDirectVariantSuffixes = map[string]bool{
	"any": true, "none": true, "low": true, "high": true,
	"SIGNING_CERT": true, "CA_CERTIFICATE": true,
}

func xcvCheckClass(check string) string {
	parts := strings.Split(check, "-")
	end := len(parts)
	for end > 1 && xcvDirectVariantSuffixes[parts[end-1]] {
		end--
	}
	return strings.Join(parts[:end], "-")
}

// ----------------------------------------------------------------- synthetic

// runSyntheticChecks replays emitSynthetic(): the branches the real corpus
// (plus XCVA's dd/*.xml dumps) never reaches.
func runSyntheticChecks(t *testing.T, expectedByToken map[string]map[string][]*xcvDirectRow,
	compareSub func(t *testing.T, expected map[string][]*xcvDirectRow, file, token, check string,
		factory func(result *process.Result[*jaxb.XmlSubXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV]),
	compareXCV func(t *testing.T, expected map[string][]*xcvDirectRow, file, token, check string,
		factory func(result *process.Result[*jaxb.XmlXCV], constraint policy.LevelRule) process.ChainItem[*jaxb.XmlXCV])) {

	expectedFor := func(token string) map[string][]*xcvDirectRow {
		if e := expectedByToken[token]; e != nil {
			return e
		}
		return map[string][]*xcvDirectRow{}
	}

	ncLeaf := nameConstraintViolationLeaf()
	compareSub(t, expectedFor(ncLeaf.Id()), "synthetic", ncLeaf.Id(), "CertificateNameConstraintsCheck",
		func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
			return NewCertificateNameConstraintsCheck(xcvaI18n(), r, ncLeaf, l)
		})

	noRevAvailViolation := noRevAvailViolationCertificate()
	compareSub(t, expectedFor(noRevAvailViolation.Id()), "synthetic", noRevAvailViolation.Id(), "NoRevAvailCheck",
		func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
			return NewNoRevAvailCheck(xcvaI18n(), r, noRevAvailViolation, l)
		})

	blankSerial := blankSerialNumberCertificate()
	compareSub(t, expectedFor(blankSerial.Id()), "synthetic", blankSerial.Id(), "SerialNumberCheck",
		func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
			return NewSerialNumberCheck(xcvaI18n(), r, blankSerial, l)
		})

	qc := qcStatementCertificate()
	qcExpected := expectedFor(qc.Id())
	compareSub(t, qcExpected, "synthetic", qc.Id(), "PseudonymCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewPseudonymCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcCCLegislationCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcCCLegislationCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcQSCDLegislationCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcQSCDLegislationCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcIdentificationMethodCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcIdentificationMethodCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificatePS2DQcCompetentAuthorityIdCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificatePS2DQcCompetentAuthorityIdCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificatePS2DQcCompetentAuthorityNameCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificatePS2DQcCompetentAuthorityNameCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificatePS2DQcRolesOfPSPCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificatePS2DQcRolesOfPSPCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcPSBCountryOfLegislationCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcPSBCountryOfLegislationCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcPSBAuthSourceIdentificationCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcPSBAuthSourceIdentificationCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})
	compareSub(t, qcExpected, "synthetic", qc.Id(), "CertificateQcPSBLegislationIdentificationCheck-any", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateQcPSBLegislationIdentificationCheck(xcvaI18n(), r, qc, xcvDirectAny)
	})

	usageTime := time.Unix(1700000000, 0).UTC()
	trustService := trustServiceCertificate(usageTime)
	trustExpected := expectedFor(trustService.Id())
	compareXCV(t, trustExpected, "synthetic", trustService.Id(), "TrustServiceStatusCheck-any", func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
		return NewTrustServiceStatusCheck(xcvaI18n(), r, trustService, &usageTime, enumerations.ContextSignature, xcvDirectAny)
	})
	compareXCV(t, trustExpected, "synthetic", trustService.Id(), "TrustServiceTypeIdentifierCheck-any", func(r *process.Result[*jaxb.XmlXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlXCV] {
		return NewTrustServiceTypeIdentifierCheck(xcvaI18n(), r, trustService, &usageTime, enumerations.ContextSignature, xcvDirectAny)
	})

	onHold := onHoldRevocation()
	compareSub(t, expectedFor(onHold.Id()), "synthetic", onHold.Id(), "CertificateNotOnHoldCheck", func(r *process.Result[*jaxb.XmlSubXCV], l policy.LevelRule) process.ChainItem[*jaxb.XmlSubXCV] {
		return NewCertificateNotOnHoldCheck(xcvaI18n(), r, onHold, xcvDirectCurrentTime, l)
	})
}

// onHoldRevocation is a certificate revocation whose reason is
// CERTIFICATE_HOLD, at a past revocation date.
func onHoldRevocation() *diagnostic.CertificateRevocationWrapper {
	revocation := &diagjaxb.XmlRevocation{}
	revocationID := diagjaxb.CollapsedString("R-SYNTH-ON-HOLD")
	revocation.XmlAbstractTokenAttrs.Id = &revocationID
	status := diagjaxb.CertificateStatusValue(enumerations.CertificateStatusRevoked)
	reason := diagjaxb.RevocationReasonValue(enumerations.RevocationReasonCertificateHold)
	revocationDate := diagjaxb.XSDateTime(xcvDirectCurrentTime.Add(-24 * time.Hour))
	certRevocation := &diagjaxb.XmlCertificateRevocation{
		Revocation:     revocation,
		Status:         &status,
		Reason:         &reason,
		RevocationDate: &revocationDate,
	}
	return diagnostic.NewCertificateRevocationWrapper(certRevocation)
}

// nameConstraintViolationLeaf is a leaf certificate whose DN falls outside its
// issuing CA's permitted name-constraint subtree.
func nameConstraintViolationLeaf() *diagnostic.CertificateWrapper {
	ca := baseCertificate("C-SYNTH-NC-CA")
	permitted := &diagjaxb.XmlGeneralSubtree{
		XmlGeneralNameContent: diagjaxb.XmlGeneralNameContent{Value: "O=Allowed"},
	}
	nameType := diagjaxb.GeneralNameTypeValue(enumerations.GeneralNameTypeDirectoryName)
	permitted.XmlGeneralNameAttrs.Type = &nameType
	oid := enumerations.CertificateExtensionEnumNameConstraints.OID()
	nameConstraints := &diagjaxb.XmlNameConstraints{
		PermittedSubtree: []*diagjaxb.XmlGeneralSubtree{permitted},
	}
	nameConstraints.OID = &oid
	addExtension(ca, nameConstraints)

	leaf := baseCertificate("C-SYNTH-NC-LEAF")
	format := "RFC2253"
	leaf.SubjectDistinguishedName = []*diagjaxb.XmlDistinguishedName{
		{Value: "CN=Leaf,O=Excluded", Format: &format},
	}
	leaf.CertificateChain = &diagjaxb.CertificateChainWrapper{
		Items: []*diagjaxb.XmlChainItem{{Certificate: ca}},
	}
	return diagnostic.NewCertificateWrapper(leaf)
}

// noRevAvailViolationCertificate declares noRevAvail while still publishing a
// conflicting OCSP access point (RFC 9608 violation).
func noRevAvailViolationCertificate() *diagnostic.CertificateWrapper {
	xml := baseCertificate("C-SYNTH-NORA-VIOLATION")
	present := true
	noRevAvail := &diagjaxb.XmlNoRevAvail{Present: &present}
	noRevAvailOID := enumerations.CertificateExtensionEnumNoRevocationAvailable.OID()
	noRevAvail.OID = &noRevAvailOID
	addExtension(xml, noRevAvail)

	aia := &diagjaxb.XmlAuthorityInformationAccess{OcspUrl: []string{"http://example.org/ocsp"}}
	aiaOID := enumerations.CertificateExtensionEnumAuthorityInformationAccess.OID()
	aia.OID = &aiaOID
	addExtension(xml, aia)
	return diagnostic.NewCertificateWrapper(xml)
}

// blankSerialNumberCertificate carries no serial number.
func blankSerialNumberCertificate() *diagnostic.CertificateWrapper {
	return diagnostic.NewCertificateWrapper(baseCertificate("C-SYNTH-NO-SERIAL"))
}

// qcStatementCertificate carries every QC/PSD2/pseudonym attribute the corpus
// has none of.
func qcStatementCertificate() *diagnostic.CertificateWrapper {
	xml := baseCertificate("C-SYNTH-QC")
	pseudonym := "Synthetic Pseudonym"
	xml.Pseudonym = &pseudonym

	euCC := "EU"
	euQSCD := "EU"
	countryOfLegislation := "EU"
	authSource := "Synthetic Authentic Source"
	legislationID := "Synthetic Legislation"
	ncaID := "EU-NCA-Synthetic"
	ncaName := "Synthetic National Competent Authority"
	roleName := "PSP_AS"

	qc := &diagjaxb.XmlQcStatements{
		QcCClegislation:   &diagjaxb.QcCClegislationWrapper{Items: []string{euCC}},
		QcQSCDlegislation: &diagjaxb.QcQSCDlegislationWrapper{Items: []string{euQSCD}},
	}
	identMethodDesc := "physical presence"
	qc.QcIdentMethod = &diagjaxb.XmlOID{
		XmlOIDContent: diagjaxb.XmlOIDContent{Value: "0.4.0.1862.1.6.1"},
		XmlOIDAttrs:   diagjaxb.XmlOIDAttrs{Description: &identMethodDesc},
	}
	semanticsDesc := "natural person"
	qc.SemanticsIdentifier = &diagjaxb.XmlOID{
		XmlOIDContent: diagjaxb.XmlOIDContent{Value: "0.4.0.194121.1.1"},
		XmlOIDAttrs:   diagjaxb.XmlOIDAttrs{Description: &semanticsDesc},
	}
	qc.QcPSB = &diagjaxb.XmlQcPSB{
		CountryOfLegislation:      &countryOfLegislation,
		AuthSourceIdentification:  &authSource,
		LegislationIdentification: &legislationID,
	}
	role := &diagjaxb.XmlRoleOfPSP{
		Name: &roleName,
		Oid:  &diagjaxb.XmlOID{XmlOIDContent: diagjaxb.XmlOIDContent{Value: "0.4.0.19495.1.1"}},
	}
	qc.PSD2QcInfo = &diagjaxb.XmlPSD2QcInfo{
		NcaId:      &ncaID,
		NcaName:    &ncaName,
		RolesOfPSP: &diagjaxb.RolesOfPSPWrapper{Items: []*diagjaxb.XmlRoleOfPSP{role}},
	}

	qcOID := enumerations.CertificateExtensionEnumQCStatements.OID()
	qc.OID = &qcOID
	addExtension(xml, qc)
	return diagnostic.NewCertificateWrapper(xml)
}

// trustServiceCertificate is associated with a TrustServiceProvider entry
// covering usageTime.
func trustServiceCertificate(usageTime time.Time) *diagnostic.CertificateWrapper {
	xml := baseCertificate("C-SYNTH-TRUST-SERVICE")

	serviceType := "http://uri.etsi.org/TrstSvc/Svctype/CA/QC"
	status := "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted"
	startDate := diagjaxb.XSDateTime(usageTime.Add(-365 * 24 * time.Hour))
	service := &diagjaxb.XmlTrustService{
		XmlTrustedEntityServiceContent: diagjaxb.XmlTrustedEntityServiceContent{
			ServiceType: &serviceType,
			Status:      &status,
			StartDate:   &startDate,
		},
		XmlTrustedEntityServiceAttrs: diagjaxb.XmlTrustedEntityServiceAttrs{
			ServiceDigitalIdentifier: xml,
		},
	}
	provider := &diagjaxb.XmlTrustServiceProvider{
		TrustServices: &diagjaxb.TrustServicesWrapper{Items: []*diagjaxb.XmlTrustService{service}},
	}
	xml.TrustServiceProviders = &diagjaxb.TrustServiceProvidersWrapper{
		Items: []*diagjaxb.XmlTrustServiceProvider{provider},
	}
	return diagnostic.NewCertificateWrapper(xml)
}

func baseCertificate(id string) *diagjaxb.XmlCertificate {
	xml := &diagjaxb.XmlCertificate{}
	idValue := diagjaxb.CollapsedString(id)
	xml.XmlAbstractTokenAttrs.Id = &idValue
	return xml
}

func addExtension(xml *diagjaxb.XmlCertificate, item diagjaxb.XmlCertificateExtensionItem) {
	if xml.CertificateExtensions == nil {
		xml.CertificateExtensions = &diagjaxb.CertificateExtensionsWrapper{}
	}
	xml.CertificateExtensions.Items = append(xml.CertificateExtensions.Items, item)
}

// -------------------------------------------------------- wrapper shapes

// aovOfShape builds the XmlAOV wrapper result CertificateAlgorithmObsolescenceValidationCheck
// consumes, at the given conclusion shape.
func aovOfShape(shape, certificateID string) *jaxb.XmlAOV {
	aov := &jaxb.XmlAOV{}
	conclusion := &jaxb.XmlConclusion{}
	switch shape {
	case "error":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationCryptoConstraintsFailure)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors, xcvDirectMessage("ASCCM_AR_ANS_ANR", "The algorithm is no longer reliable!"))
	case "warning":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		conclusion.Warnings = append(conclusion.Warnings, xcvDirectMessage("ASCCM_AR_ANS_AKSNR", "The key size is no longer reliable!"))
	default:
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
	}
	aov.Conclusion = conclusion

	algorithm := &jaxb.XmlCryptographicAlgorithm{Name: "RSA with SHA256"}
	if shape == "passed-with-algo" {
		keyLength := "2048"
		algorithm.KeyLength = &keyLength
	}
	tokenID := certificateID
	certValidation := &jaxb.XmlCryptographicValidation{
		Algorithm:  algorithm,
		Conclusion: conclusion,
		TokenId:    &tokenID,
	}
	aov.CertificateChainCryptographicValidation = &jaxb.XmlCertificateChainCryptographicValidation{
		CertificateCryptographicValidation: []*jaxb.XmlCryptographicValidation{certValidation},
	}
	return aov
}

func xcvDirectMessage(key, value string) *jaxb.XmlMessage {
	k := key
	return &jaxb.XmlMessage{Key: &k, Value: value}
}

func crsOfIndication(indication enumerations.Indication, id string) *jaxb.XmlCRS {
	crs := &jaxb.XmlCRS{}
	crsID := id
	crs.Id = &crsID
	conclusion := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
	if indication != enumerations.IndicationPassed {
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationCertificateChainGeneralFailure)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors, xcvDirectMessage("BBB_XCV_IARDPFC_ANS", "No acceptable revocation data found."))
	} else {
		latest := "R-" + id
		crs.LatestAcceptableRevocationId = &latest
	}
	crs.Conclusion = conclusion
	return crs
}

func rfcOfIndication(indication enumerations.Indication, id string) *jaxb.XmlRFC {
	rfc := &jaxb.XmlRFC{}
	rfcID := id
	rfc.Id = &rfcID
	conclusion := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
	if indication != enumerations.IndicationPassed {
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationTryLater)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors, xcvDirectMessage("BBB_RFC_IRIF_ANS", "The revocation is not considered as 'fresh'."))
	}
	rfc.Conclusion = conclusion
	return rfc
}

func subXCVOfIndication(indication enumerations.Indication, id string) *jaxb.XmlSubXCV {
	subXCV := &jaxb.XmlSubXCV{}
	subXCV.Id = id
	conclusion := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
	if indication != enumerations.IndicationPassed {
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationChainConstraintsFailure)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors, xcvDirectMessage("BBB_XCV_SUB_ANS", "The certificate validation is not conclusive!"))
	}
	subXCV.Conclusion = conclusion
	return subXCV
}

var _ = big.NewInt // silence unused import if pruned during edits
