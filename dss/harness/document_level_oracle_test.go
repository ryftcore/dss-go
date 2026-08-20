// THE PHASE 8 EXIT CRITERION - document-level end-to-end parity, item (B) of
// the s8f harness brief.
//
// testdata/oracle/document_level.jsonl is a pure Java dump, produced by
// testdata/oracle/gen/DocumentLevelOracle.java, for the 60 real signed
// fixtures listed in testdata/oracle/document_level_manifest.tsv - 10 per
// format family (CAdES, XAdES, PAdES, JAdES, ASiC-CAdES, ASiC-XAdES), the
// fixtures already vendored under each format package's testdata/upstream
// tree (the same trees the six format packages' *_smoke_test.go files read).
//
// Both sides run the SAME pipeline as the smoke tests -
// SignedDocumentValidator.fromDocument's format-autodetection dispatch, a
// permissive CertificateVerifier (CommonCertificateVerifier(true) with the
// seven Alert setters silenced - identical to every smoke test's
// permissiveCertificateVerifier(), since these are offline fixtures with no
// live revocation data or trust anchors), locale en - and additionally pin
// ValidationLevel.ARCHIVAL_DATA (the smoke tests use BASIC_SIGNATURES; this
// harness exercises the fuller LTA/archival path per the brief, which is also
// what item (A)'s and (C)'s oracles use).
//
// For every fixture, this test compares:
//   - the diagnostic-data core: signatures count, used-certificates count, container type
//   - every top-level Signature/Timestamp/EvidenceRecord Indication/SubIndication
//     in the DetailedReport
//   - every BasicBuildingBlocks conclusion
//   - every signature's SimpleReport qualification
//
// This is document-level, not diagnostic-data-level: both sides parse the
// ORIGINAL signed document (CMS/XML/PDF/JWS/ZIP) from scratch, so it also
// exercises each format's own diagnostic-data builder - the part item (A)'s
// diag-data-corpus oracle cannot reach, since that one starts from
// already-built diagnostic data.
//
// # F2 and F3, found live by this harness and FIXED in this pass
//
// Building fresh from a real document (as this test does, and item (A)'s
// oracle never does) surfaced two real defects item (A) cannot see. Both are
// now fixed, so this test compares every fixture with NO tolerance at all -
// exact BasicBuildingBlocks Ids included:
//
//   - F2, in spi/validation/timestamp/signature_timestamp_source.go's
//     getAttributeOrder: it compared the carrying SignatureAttribute by Go
//     POINTER identity, where Java calls signatureAttribute.equals(property),
//     which every concrete attribute class (CAdESAttribute, XAdESAttribute,
//     JAdESAttribute, CBAdESAttribute) overrides as identifier equality. Since
//     SignatureProperties.Attributes() rebuilds its list on every call, the
//     pointer never matched, the order silently came back nil, and the
//     "-OOA-<n>" component vanished from the position string every encapsulated
//     TimestampToken's T-... identifier hashes - wrong ids on 30 of these 60
//     fixtures (every CAdES/XAdES/PAdES/ASiC format carrying an embedded
//     signature- or archive-timestamp).
//
//   - F3, in cms/cms_utils.go's CMSUtilsWriteContentInfoEncoded: it wrote the
//     OUTER CMS ContentInfo instead of SignedData.encapContentInfo, feeding the
//     whole signature into every CAdES archive-timestamp-v2 message imprint, so
//     such timestamps verified FAILED/HASH_FAILURE where Java found them intact
//     (cades/baseline-lta, Signature-C-B-LTA-10.p7m).
package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/utain/esig/dss/alert"
	detailedreportjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dsspolicy "github.com/utain/esig/dss/policy"
	cryptoxml "github.com/utain/esig/dss/policy/crypto/xml"
	simplereportjaxb "github.com/utain/esig/dss/simplereport/jaxb"
	"github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
	"github.com/utain/esig/dss/validation/reports"

	// Blank-imported purely for their init()-registered
	// dssvalidation.RegisterDocumentValidatorFactory side effects, which
	// SignedDocumentValidatorFromDocument's format-autodetection dispatch
	// needs for all six format families this corpus exercises.
	_ "github.com/utain/esig/dss/asic/cades"
	_ "github.com/utain/esig/dss/asic/xades"
	_ "github.com/utain/esig/dss/cades"
	_ "github.com/utain/esig/dss/jades"
	_ "github.com/utain/esig/dss/pades"
	_ "github.com/utain/esig/dss/xades"
)

func init() {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	validationpolicy.RegisterCryptographicSuiteFactory(cryptoxml.NewCryptographicSuiteXmlFactory())
}

// permissiveCertificateVerifier matches every format package's own
// *_smoke_test.go helper of the same name byte-for-byte in behavior.
func permissiveCertificateVerifier() *validation.CommonCertificateVerifier {
	v := validation.NewCommonCertificateVerifierSimple(true)
	v.SetAlertOnMissingRevocationData(alert.NewSilentOnStatusAlert())
	v.SetAlertOnRevokedCertificate(alert.NewSilentOnStatusAlert())
	v.SetAlertOnInvalidSignature(alert.NewSilentOnStatusAlert())
	v.SetAlertOnInvalidTimestamp(alert.NewSilentOnStatusAlert())
	v.SetAlertOnUncoveredPOE(alert.NewSilentOnStatusAlert())
	v.SetAlertOnExpiredCertificate(alert.NewSilentOnStatusAlert())
	v.SetAlertOnNotYetValidCertificate(alert.NewSilentOnStatusAlert())
	return v
}

type dlToken struct {
	Kind          string `json:"kind"`
	ID            string `json:"id"`
	Indication    string `json:"indication"`
	SubIndication string `json:"subIndication"`
}

type dlBBB struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Indication    string `json:"indication"`
	SubIndication string `json:"subIndication"`
}

type dlQualification struct {
	ID            string `json:"id"`
	Qualification string `json:"qualification"`
}

type dlCore struct {
	SignaturesCount   int    `json:"signaturesCount"`
	CertificatesCount int    `json:"certificatesCount"`
	ContainerType     string `json:"containerType"`
}

type documentLevelRow struct {
	Format        string            `json:"format"`
	Name          string            `json:"name"`
	Error         string            `json:"error"`
	Core          *dlCore           `json:"core"`
	Tokens        []dlToken         `json:"tokens"`
	BBB           []dlBBB           `json:"bbb"`
	Qualification []dlQualification `json:"qualification"`
}

type manifestEntry struct {
	Format string
	Name   string
	Path   string
}

func TestDocumentLevelOracle(t *testing.T) {
	rows := readDocumentLevelOracle(t)
	manifest := readDocumentLevelManifest(t)
	if len(rows) == 0 || len(rows) != len(manifest) {
		t.Fatalf("oracle/manifest row count mismatch: %d rows, %d manifest entries", len(rows), len(manifest))
	}

	byKey := make(map[string]documentLevelRow)
	for _, row := range rows {
		byKey[row.Format+"/"+row.Name] = row
	}

	var totalMismatches int
	perFormat := make(map[string]int)
	for _, entry := range manifest {
		entry := entry
		key := entry.Format + "/" + entry.Name
		row, ok := byKey[key]
		if !ok {
			t.Fatalf("no oracle row for manifest entry %s", key)
		}
		t.Run(key, func(t *testing.T) {
			doc, err := model.NewFileDocument(entry.Path)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", entry.Path, err)
			}

			gotReports, runErr := runDocumentValidatorSafely(doc)
			if runErr != nil {
				if row.Error == "" {
					t.Errorf("Go errored validating %s (%v) but the Java oracle succeeded - DEFECT", entry.Path, runErr)
				} else {
					t.Logf("Go also errored (%v); java: %s", runErr, row.Error)
				}
				return
			}
			if row.Error != "" {
				t.Errorf("Go succeeded validating %s but the Java oracle errored (%s) - DEFECT", entry.Path, row.Error)
				return
			}

			diagnosticData := gotReports.GetDiagnosticDataJaxb()
			gotSigCount := len(diagnosticData.Signatures.All())
			gotCertCount := len(diagnosticData.UsedCertificates.All())
			gotContainerType := ""
			if diagnosticData.ContainerInfo != nil && diagnosticData.ContainerInfo.ContainerType != nil {
				gotContainerType = string(*diagnosticData.ContainerInfo.ContainerType)
			}
			if row.Core == nil {
				t.Fatalf("oracle row %s has no core section", key)
			}
			if gotSigCount != row.Core.SignaturesCount {
				t.Errorf("%s: signaturesCount go=%d java=%d", key, gotSigCount, row.Core.SignaturesCount)
				totalMismatches++
			}
			if gotCertCount != row.Core.CertificatesCount {
				t.Errorf("%s: certificatesCount go=%d java=%d", key, gotCertCount, row.Core.CertificatesCount)
				totalMismatches++
			}
			if gotContainerType != row.Core.ContainerType {
				t.Errorf("%s: containerType go=%s java=%s", key, gotContainerType, row.Core.ContainerType)
				totalMismatches++
			}

			mismatches := compareDocumentLevelTokens(t, key, row.Tokens, gotReports.GetDetailedReportJaxb())
			mismatches += compareDocumentLevelBBB(t, key, row.BBB, gotReports.GetDetailedReportJaxb().BasicBuildingBlocks)
			mismatches += compareDocumentLevelQualification(t, key, row.Qualification, gotReports.GetSimpleReportJaxb())
			totalMismatches += mismatches
			perFormat[entry.Format] += mismatches
		})
	}
	t.Logf("document-level oracle: %d fixtures, %d field mismatches, by format: %v", len(manifest), totalMismatches, perFormat)
}

// runDocumentValidatorSafely recovers from a panic so one malformed fixture
// reports as a mismatch on its own subtest rather than aborting the whole
// corpus run - a panic is exactly the kind of Go/Java divergence this oracle
// exists to surface.
func runDocumentValidatorSafely(doc model.DSSDocument) (r *reportsResult, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()
	documentValidator, ferr := dssvalidation.SignedDocumentValidatorFromDocument(doc)
	if ferr != nil {
		return nil, ferr
	}
	documentValidator.SetCertificateVerifier(permissiveCertificateVerifier())
	documentValidator.SetValidationLevel(enumerations.ValidationLevel_ARCHIVAL_DATA)
	documentValidator.SetLocale("en")
	reports, verr := documentValidator.ValidateDocument()
	if verr != nil {
		return nil, verr
	}
	return &reportsResult{reports: reports}, nil
}

// reportsResult adapts *reports.Reports's three JAXB accessors so the
// comparison helpers below don't need to import the reports package just for
// this one indirection.
type reportsResult struct {
	reports *reports.Reports
}

func (r *reportsResult) GetDiagnosticDataJaxb() *diagnosticjaxb.XmlDiagnosticData {
	return r.reports.GetDiagnosticDataJaxb()
}

func (r *reportsResult) GetDetailedReportJaxb() *detailedreportjaxb.XmlDetailedReport {
	return r.reports.GetDetailedReportJaxb()
}

func (r *reportsResult) GetSimpleReportJaxb() *simplereportjaxb.XmlSimpleReport {
	return r.reports.GetSimpleReportJaxb()
}

// compareDocumentLevelTokens compares the Java-dumped top-level
// Signature/Timestamp/EvidenceRecord verdicts against the Go DetailedReport.
func compareDocumentLevelTokens(t *testing.T, key string, want []dlToken, detailed *detailedreportjaxb.XmlDetailedReport) int {
	t.Helper()
	mismatches := 0
	gotByKey := make(map[string]dlToken)
	var gotOrder []string
	for _, item := range detailed.SignatureOrTimestampOrEvidenceRecord {
		var kind, id string
		var conclusion *detailedreportjaxb.XmlConclusion
		switch v := item.(type) {
		case *detailedreportjaxb.XmlSignature:
			kind, conclusion = "Signature", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		case *detailedreportjaxb.XmlTimestamp:
			kind, conclusion = "Timestamp", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		case *detailedreportjaxb.XmlEvidenceRecord:
			kind, conclusion = "EvidenceRecord", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		default:
			continue
		}
		tok := dlToken{Kind: kind, ID: id}
		if conclusion != nil {
			tok.Indication = string(conclusion.Indication)
			if conclusion.SubIndication != nil {
				tok.SubIndication = string(*conclusion.SubIndication)
			}
		}
		mapKey := kind + "/" + id
		gotByKey[mapKey] = tok
		gotOrder = append(gotOrder, mapKey)
	}
	wantKeys := make(map[string]bool)
	for _, w := range want {
		mapKey := w.Kind + "/" + w.ID
		wantKeys[mapKey] = true
		g, ok := gotByKey[mapKey]
		if !ok {
			t.Errorf("%s: missing token %s (java: indication=%s subIndication=%s)", key, mapKey, w.Indication, w.SubIndication)
			mismatches++
			continue
		}
		if g.Indication != w.Indication || g.SubIndication != w.SubIndication {
			t.Errorf("%s: token %s indication/subIndication go=%s/%s java=%s/%s",
				key, mapKey, g.Indication, g.SubIndication, w.Indication, w.SubIndication)
			mismatches++
		}
	}
	for _, mapKey := range gotOrder {
		if !wantKeys[mapKey] {
			t.Errorf("%s: unexpected extra token %s in the Go report", key, mapKey)
			mismatches++
		}
	}
	return mismatches
}

// compareDocumentLevelBBB compares a Java-dumped BasicBuildingBlocks list
// against the Go one, keyed by Id, reporting every difference via t.Errorf.
func compareDocumentLevelBBB(t *testing.T, key string, want []dlBBB, got []*detailedreportjaxb.XmlBasicBuildingBlocks) int {
	t.Helper()
	diffs := diffDocumentLevelBBB(want, got)
	for _, d := range diffs {
		t.Errorf("%s: %s", key, d)
	}
	return len(diffs)
}

// diffDocumentLevelBBB is the comparison compareDocumentLevelBBB runs,
// returning one human-readable line per difference found.
func diffDocumentLevelBBB(want []dlBBB, got []*detailedreportjaxb.XmlBasicBuildingBlocks) []string {
	var diffs []string
	gotByID := make(map[string]dlBBB)
	var gotIDs []string
	for _, bbb := range got {
		entry := dlBBB{ID: bbb.Id, Type: string(bbb.Type)}
		if bbb.Conclusion != nil {
			entry.Indication = string(bbb.Conclusion.Indication)
			if bbb.Conclusion.SubIndication != nil {
				entry.SubIndication = string(*bbb.Conclusion.SubIndication)
			}
		}
		gotByID[bbb.Id] = entry
		gotIDs = append(gotIDs, bbb.Id)
	}
	wantIDs := make(map[string]bool)
	for _, w := range want {
		wantIDs[w.ID] = true
		g, ok := gotByID[w.ID]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("missing BasicBuildingBlocks %s (java: type=%s indication=%s subIndication=%s)",
				w.ID, w.Type, w.Indication, w.SubIndication))
			continue
		}
		if g.Type != w.Type || g.Indication != w.Indication || g.SubIndication != w.SubIndication {
			diffs = append(diffs, fmt.Sprintf("BasicBuildingBlocks %s go=(%s,%s,%s) java=(%s,%s,%s)",
				w.ID, g.Type, g.Indication, g.SubIndication, w.Type, w.Indication, w.SubIndication))
		}
	}
	for _, id := range gotIDs {
		if !wantIDs[id] {
			diffs = append(diffs, fmt.Sprintf("unexpected extra BasicBuildingBlocks %s in the Go report", id))
		}
	}
	return diffs
}

// compareDocumentLevelQualification compares Java-dumped signature
// qualifications against the Go SimpleReport.
func compareDocumentLevelQualification(t *testing.T, key string, want []dlQualification, simple *simplereportjaxb.XmlSimpleReport) int {
	t.Helper()
	mismatches := 0
	gotQual := make(map[string]string)
	for _, item := range simple.SignatureOrTimestampOrEvidenceRecord {
		sig, ok := item.(*simplereportjaxb.XmlSignature)
		if !ok {
			continue
		}
		q := ""
		if sig.SignatureLevel != nil {
			q = string(sig.SignatureLevel.Value)
		}
		gotQual[sig.Id] = q
	}
	seen := make(map[string]bool)
	for _, w := range want {
		seen[w.ID] = true
		g, ok := gotQual[w.ID]
		if !ok {
			t.Errorf("%s: missing simple-report signature %s (java qualification=%s)", key, w.ID, w.Qualification)
			mismatches++
			continue
		}
		if g != w.Qualification {
			t.Errorf("%s: signature %s qualification go=%s java=%s", key, w.ID, g, w.Qualification)
			mismatches++
		}
	}
	for id := range gotQual {
		if !seen[id] {
			t.Errorf("%s: unexpected extra simple-report signature %s in the Go report", key, id)
			mismatches++
		}
	}
	return mismatches
}

// readDocumentLevelOracle loads testdata/oracle/document_level.jsonl.
func readDocumentLevelOracle(t *testing.T) []documentLevelRow {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "oracle", "document_level.jsonl"))
	if err != nil {
		t.Fatalf("reading the document-level oracle dump: %v", err)
	}
	var rows []documentLevelRow
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row documentLevelRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatalf("parsing the document-level oracle dump: %v", err)
		}
		rows = append(rows, row)
	}
	return rows
}

// readDocumentLevelManifest loads testdata/oracle/document_level_manifest.tsv,
// sorted by format then name so the subtest order is stable.
func readDocumentLevelManifest(t *testing.T) []manifestEntry {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "oracle", "document_level_manifest.tsv"))
	if err != nil {
		t.Fatalf("reading the document-level manifest: %v", err)
	}
	var entries []manifestEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 3 {
			t.Fatalf("malformed manifest line: %q", line)
		}
		entries = append(entries, manifestEntry{Format: parts[0], Name: parts[1], Path: parts[2]})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Format != entries[j].Format {
			return entries[i].Format < entries[j].Format
		}
		return entries[i].Name < entries[j].Name
	})
	return entries
}
