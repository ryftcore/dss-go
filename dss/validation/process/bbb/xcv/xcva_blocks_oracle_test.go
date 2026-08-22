package xcv

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// The XCVA block KAT: every row of testdata/oracle/xcva_blocks.jsonl is the
// XmlXCV / XmlCRS / XmlRAC upstream's X509CertificateValidation,
// CertificateRevocationSelector or RevocationAcceptanceChecker produces for one
// token of one diagnostic-data dump, at validation time 2024-01-01T00:00:00Z -
// see testdata/gen/XcvaOracle.java. This test replays the same tokens and
// compares the whole produced tree: the constraint sequence, the conclusion, and
// the nested SubXCV / RAC / CRS results.
//
// X509CertificateValidation is driven the way BasicBuildingBlocks drives it, with
// the usage time that dispatcher passes for each kind of token.

func TestXcvaBlocksAgainstJavaOracle(t *testing.T) {
	rows := loadXcvaRows(t, corpustest.Path(t, "oracle/xcva_blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	policies := xcvaPolicies(t)

	// index: file -> policy -> token
	byFile := map[string]map[string]map[string]*xcvaRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]map[string]*xcvaRow{}
		}
		if _, seen := byFile[row.File][row.Policy]; !seen {
			byFile[row.File][row.Policy] = map[string]*xcvaRow{}
		}
		byFile[row.File][row.Policy][row.Token] = row
	}

	blocks := map[string]int{}
	indications := map[string]int{}
	matched := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadXcvaDiagnosticData(t, name)

			var policyNames []string
			for policyName := range byFile[name] {
				policyNames = append(policyNames, policyName)
			}
			sort.Strings(policyNames)

			for _, policyName := range policyNames {
				expected := byFile[name][policyName]
				validationPolicy := policies[policyName]
				if validationPolicy == nil {
					t.Fatalf("unknown policy %q in the oracle", policyName)
				}

				compare := func(token string, run func() *xcvaNode) {
					want, ok := expected[token]
					if !ok {
						return // the Java oracle threw on this token and recorded no row
					}
					var node *xcvaNode
					if xcvaSafeExecute(func() { node = run() }) {
						t.Errorf("%s / %s: Go panicked where Java produced a row", policyName, token)
						return
					}
					got := &xcvaRow{File: name, Policy: policyName, Token: token,
						Context: want.Context, Block: want.Block, xcvaNode: *node}
					matched++
					blocks[want.Block]++
					if want.Conclusion != nil && want.Conclusion.Indication != nil {
						indications[*want.Conclusion.Indication]++
					}
					if !reflect.DeepEqual(want, got) {
						t.Errorf("%s / %s: %s mismatch\nexpected: %s\nactual:   %s",
							policyName, token, want.Block, mustXcvaJSON(t, want), mustXcvaJSON(t, got))
					}
				}

				for _, signature := range diagnosticData.Signatures() {
					context := enumerations.Context_SIGNATURE
					if signature.IsCounterSignature() {
						context = enumerations.Context_COUNTER_SIGNATURE
					}
					signingCertificate := signature.SigningCertificate()
					if signingCertificate == nil {
						continue
					}
					cert, ctx := signingCertificate, context
					compare("SIG|"+signature.Id(), func() *xcvaNode {
						return toXcvaXCV(NewX509CertificateValidation(xcvaI18n(), cert, xcvaCurrentTime,
							cert.NotBefore(), ctx, passedAOV(), validationPolicy).Execute())
					})
				}
				for _, timestamp := range diagnosticData.TimestampList() {
					signingCertificate := timestamp.SigningCertificate()
					if signingCertificate == nil {
						continue
					}
					cert, usageTime := signingCertificate, timestamp.ProductionTime()
					compare("TST|"+timestamp.Id(), func() *xcvaNode {
						return toXcvaXCV(NewX509CertificateValidation(xcvaI18n(), cert, xcvaCurrentTime,
							usageTime, enumerations.Context_TIMESTAMP, passedAOV(), validationPolicy).Execute())
					})
				}
				for _, revocation := range xcvaRevocationsSortedById(diagnosticData) {
					signingCertificate := revocation.SigningCertificate()
					if signingCertificate == nil {
						continue
					}
					cert, usageTime := signingCertificate, revocation.ProductionDate()
					compare("REV|"+revocation.Id(), func() *xcvaNode {
						return toXcvaXCV(NewX509CertificateValidation(xcvaI18n(), cert, xcvaCurrentTime,
							usageTime, enumerations.Context_REVOCATION, passedAOV(), validationPolicy).Execute())
					})
				}
				for _, certificate := range diagnosticData.UsedCertificates() {
					cert := certificate
					compare("CERT|"+cert.Id(), func() *xcvaNode {
						return toXcvaXCV(NewX509CertificateValidation(xcvaI18n(), cert, xcvaCurrentTime,
							cert.NotBefore(), enumerations.Context_CERTIFICATE, passedAOV(), validationPolicy).Execute())
					})
					compare("CRS|"+cert.Id(), func() *xcvaNode {
						return toXcvaCRS(NewCertificateRevocationSelector(xcvaI18n(), cert,
							xcvaCurrentTime, validationPolicy).Execute())
					})
					for _, revocation := range cert.CertificateRevocationData() {
						rev := revocation
						compare("RAC|"+cert.Id()+"|"+rev.Id(), func() *xcvaNode {
							return toXcvaRAC(NewRevocationAcceptanceChecker(xcvaI18n(), cert, rev,
								xcvaCurrentTime, validationPolicy, map[string]struct{}{}).Execute())
						})
					}
				}
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}
	for _, block := range []string{"XCV", "CRS", "RAC"} {
		if blocks[block] == 0 {
			t.Errorf("no %s row replayed", block)
		}
	}
	// The corpus must keep both a passing and a failing verdict, so that a
	// regression cannot hide behind inputs that only ever fail.
	for _, indication := range []string{"PASSED", "INDETERMINATE"} {
		if indications[indication] == 0 {
			t.Errorf("no %s conclusion in the replayed corpus", indication)
		}
	}
}

// xcvaRevocationsSortedById mirrors the driver, which sorts the Java Set of
// revocation wrappers by id before iterating it.
func xcvaRevocationsSortedById(data *diagnostic.DiagnosticData) []*diagnostic.RevocationWrapper {
	revocations := append([]*diagnostic.RevocationWrapper(nil), data.AllRevocationData()...)
	sort.SliceStable(revocations, func(i, j int) bool { return revocations[i].Id() < revocations[j].Id() })
	return revocations
}
