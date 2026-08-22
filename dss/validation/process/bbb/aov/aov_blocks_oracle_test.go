package aov

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// The AOV block KAT: every row of testdata/oracle/aov_blocks.jsonl is the XmlAOV
// one of upstream's ten concrete AlgorithmObsolescenceValidation subclasses
// produces for one token of one diagnostic-data dump, at validation time
// 2024-01-01T00:00:00Z under the default ETSI policy - see
// testdata/gen/AovOracle.java. This test replays the same tokens and compares the
// whole produced tree: the constraint sequence, the conclusion, and all four
// XmlCryptographicValidation members.

func TestAovBlocksAgainstJavaOracle(t *testing.T) {
	rows := loadAovRows(t, corpustest.Path(t, "oracle/aov_blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	validationPolicy := aovDefaultPolicy(t)

	// index: file -> token -> block
	byFile := map[string]map[string]map[string]*aovRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]map[string]*aovRow{}
		}
		if _, seen := byFile[row.File][row.Token]; !seen {
			byFile[row.File][row.Token] = map[string]*aovRow{}
		}
		byFile[row.File][row.Token][row.Block] = row
	}

	blocks := map[string]int{}
	indications := map[string]int{}
	matched := 0
	deviations := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadAovDiagnosticData(t, name)

			compare := func(token, block string, run func() *jaxb.XmlAOV) {
				want, ok := byFile[name][token][block]
				if !ok {
					return // the Java oracle threw on this token and recorded no row
				}
				var produced *jaxb.XmlAOV
				if aovSafeExecute(func() { produced = run() }) {
					t.Errorf("%s / %s: Go panicked where Java produced a row", token, block)
					return
				}
				got := toAovBlockRow(name, token, block, want.Context, produced)
				matched++
				blocks[block]++
				if want.Conclusion != nil && want.Conclusion.Indication != nil {
					indications[*want.Conclusion.Indication]++
				}
				if !reflect.DeepEqual(want, got) {
					if aovKnownSignCertRefOrderDeviation[name+" "+token+" "+block] &&
						aovSameConstraintMultiset(t, want, got) {
						deviations++
						return
					}
					t.Errorf("%s / %s mismatch\nexpected: %s\nactual:   %s",
						token, block, mustAovJSON(t, want), mustAovJSON(t, got))
				}
			}

			for _, signature := range diagnosticData.Signatures() {
				context := enumerations.ContextSignature
				if signature.IsCounterSignature() {
					context = enumerations.ContextCounterSignature
				}
				sig, ctx := signature, context
				id := "SIG|" + sig.Id()
				compare(id, "SignatureAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewSignatureAlgorithmObsolescenceValidation(aovI18n(), sig, ctx, aovCurrentTime, validationPolicy).Execute()
				})
				compare(id, "SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewSignatureValueAndSignedAttributesAlgorithmObsolescenceValidation(aovI18n(), sig, ctx, aovCurrentTime, validationPolicy).Execute()
				})
				compare(id, "SignatureSignedDataAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewSignatureSignedDataAlgorithmObsolescenceValidation(aovI18n(), sig, ctx, aovCurrentTime, validationPolicy).Execute()
				})
				compare(id, "TokenCertificateChainAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewTokenCertificateChainAlgorithmObsolescenceValidation(aovI18n(), sig, ctx, aovCurrentTime, validationPolicy).Execute()
				})
			}
			for _, timestamp := range diagnosticData.TimestampList() {
				tst := timestamp
				id := "TST|" + tst.Id()
				compare(id, "TimestampAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewTimestampAlgorithmObsolescenceValidation(aovI18n(), tst, aovCurrentTime, validationPolicy).Execute()
				})
				compare(id, "TokenCertificateChainAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewTokenCertificateChainAlgorithmObsolescenceValidation(aovI18n(), tst, enumerations.ContextTimestamp, aovCurrentTime, validationPolicy).Execute()
				})
			}
			for _, revocation := range aovRevocationsSortedById(diagnosticData) {
				rev := revocation
				compare("REV|"+rev.Id(), "RevocationDataAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewRevocationDataAlgorithmObsolescenceValidation(aovI18n(), rev, aovCurrentTime, validationPolicy).Execute()
				})
			}
			for _, evidenceRecord := range diagnosticData.EvidenceRecords() {
				er := evidenceRecord
				compare("ER|"+er.Id(), "EvidenceRecordAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewEvidenceRecordAlgorithmObsolescenceValidation(aovI18n(), er, aovCurrentTime, validationPolicy).Execute()
				})
			}
			for _, certificate := range diagnosticData.UsedCertificates() {
				cert := certificate
				for _, subContext := range []enumerations.SubContext{
					enumerations.SubContextSigningCert, enumerations.SubContextCACertificate} {
					sub := subContext
					compare("CERT|"+cert.Id()+"|"+string(sub), "CertificateAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
						return NewCertificateAlgorithmObsolescenceValidation(aovI18n(), cert,
							enumerations.ContextCertificate, sub, aovCurrentTime, validationPolicy).Execute()
					})
				}
				compare("CERT|"+cert.Id(), "CertificateAndChainAlgorithmObsolescenceValidation", func() *jaxb.XmlAOV {
					return NewCertificateAndChainAlgorithmObsolescenceValidation(aovI18n(), cert,
						enumerations.ContextCertificate, aovCurrentTime, validationPolicy).Execute()
				})
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}
	// The known deviation must stay exactly as wide as it is documented: a row
	// that stops deviating, or one that starts, both fail here.
	if deviations != len(aovKnownSignCertRefOrderDeviation) {
		t.Errorf("%d rows took the known signing-certificate-reference order deviation, %d are listed",
			deviations, len(aovKnownSignCertRefOrderDeviation))
	}
	for _, block := range []string{
		"SignatureAlgorithmObsolescenceValidation",
		"SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation",
		"SignatureSignedDataAlgorithmObsolescenceValidation",
		"TokenCertificateChainAlgorithmObsolescenceValidation",
		"TimestampAlgorithmObsolescenceValidation",
		"RevocationDataAlgorithmObsolescenceValidation",
		"CertificateAlgorithmObsolescenceValidation",
		"CertificateAndChainAlgorithmObsolescenceValidation",
	} {
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

// aovKnownSignCertRefOrderDeviation lists the "<file> <token> <block>" rows whose
// constraint sequence deviates from Java's, and only in its order.
//
// It is EMPTY: utils.JavaHashMapComputeIfAbsentKeyOrder reproduces the real
// java.util.HashMap key order that
// SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation#buildSignedAttributesValidationChain
// relies on when grouping signing-certificate references by certificate id
// (computeIfAbsent PREPENDS colliding keys, unlike put), so the constraint
// sequence - and the reported signed-attributes digest algorithm, which the
// loop's "first result wins" rule ties to the same order - match exactly.
var aovKnownSignCertRefOrderDeviation = map[string]bool{}

// aovSameConstraintMultiset reports whether two rows differ only in the order of
// their constraint list: everything outside the list must be equal, and the list
// must be a permutation of the other.
func aovSameConstraintMultiset(t *testing.T, want, got *aovRow) bool {
	t.Helper()
	if len(want.Constraints) != len(got.Constraints) {
		return false
	}
	wantOutside, gotOutside := *want, *got
	wantOutside.Constraints, gotOutside.Constraints = nil, nil
	if !reflect.DeepEqual(&wantOutside, &gotOutside) {
		return false
	}
	rendered := func(constraints []*aovConstraint) []string {
		out := make([]string, 0, len(constraints))
		for _, constraint := range constraints {
			out = append(out, mustAovJSON(t, constraint))
		}
		sort.Strings(out)
		return out
	}
	return reflect.DeepEqual(rendered(want.Constraints), rendered(got.Constraints))
}

// aovRevocationsSortedById mirrors the driver, which sorts the Java Set of
// revocation wrappers by id before iterating it.
func aovRevocationsSortedById(data *diagnostic.Data) []*diagnostic.RevocationWrapper {
	revocations := append([]*diagnostic.RevocationWrapper(nil), data.AllRevocationData()...)
	sort.SliceStable(revocations, func(i, j int) bool { return revocations[i].Id() < revocations[j].Id() })
	return revocations
}

func toAovBlockRow(file, token, block, context string, result *jaxb.XmlAOV) *aovRow {
	row := &aovRow{File: file, Token: token, Block: block, Context: context}
	row.ValidationTime = aovDate(result.ValidationTime)
	toAovRowBody(row, &result.XmlConstraintsConclusionContent, result.Title)
	row.SignatureCryptographicValidation = toAovCryptographicValidation(result.SignatureCryptographicValidation)
	row.SignedAttributesValidation = toAovCryptographicValidation(result.SignedAttributesValidation)
	row.DigestMatchersValidation = toAovCryptographicValidation(result.DigestMatchersValidation)
	if result.CertificateChainCryptographicValidation != nil {
		items := result.CertificateChainCryptographicValidation.CertificateCryptographicValidation
		row.CertificateChainCryptographicValidation = make([]*aovCryptographicValidation, 0, len(items))
		for _, item := range items {
			row.CertificateChainCryptographicValidation = append(row.CertificateChainCryptographicValidation,
				toAovCryptographicValidation(item))
		}
	}
	return row
}
