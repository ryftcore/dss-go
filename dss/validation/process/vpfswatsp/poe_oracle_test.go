package vpfswatsp

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
)

// The POE KAT: every row of testdata/oracle/poe.jsonl is what upstream's
// eu.europa.esig.dss.validation.process.vpfswatsp POE core (POE, TimestampPOE,
// EvidenceRecordPOE, POEComparator, POEExtraction) produces for one
// diagnostic-data dump at control time 2024-01-01T00:00:00Z - see
// testdata/gen/POEOracle.java. Inputs are the 54-dump marshal-parity corpus plus
// this package's seven synthetic dumps (testdata/gen/POESyntheticDumps.java),
// which supply the shapes the corpus has none of: evidence records, time-stamps
// sharing a production time, broken message imprints, orphan tokens.

const (
	poeCorpusDir = "../../../diagnostic/jaxb/testdata/oracle"
	poeDumpDir   = "testdata/dd"
	poeOracleLog = "testdata/oracle/poe.jsonl"
)

var poeControlTime = time.Unix(1704067200, 0).UTC()

// poeProbes are the isPOEExists probe times of the oracle, in its order.
var poeProbes = []time.Time{
	time.Unix(946684800, 0).UTC(),
	time.Unix(1464739200, 0).UTC(),
	time.Unix(1577836800, 0).UTC(),
	time.Unix(1704067200, 0).UTC(),
}

type poeSnapshot struct {
	Time          int64   `json:"time"`
	ProviderId    *string `json:"providerId"`
	TokenProvided bool    `json:"tokenProvided"`
	Objects       int     `json:"objects"`
}

type poeTokenRow struct {
	Id                   string       `json:"id"`
	AfterInit            *poeSnapshot `json:"afterInit"`
	AfterTimestamps      *poeSnapshot `json:"afterTimestamps"`
	AfterEvidenceRecords *poeSnapshot `json:"afterEvidenceRecords"`
	Exists               []bool       `json:"exists"`
	InValidityRange      *bool        `json:"inValidityRange"`
}

type poeCompareRow struct {
	A       string `json:"a"`
	B       string `json:"b"`
	Compare int    `json:"compare"`
	Before  bool   `json:"before"`
}

type poeSignatureRow struct {
	Id     string       `json:"id"`
	Lowest *poeSnapshot `json:"lowest"`
}

type poeRow struct {
	Dump         string             `json:"dump"`
	Tokens       []*poeTokenRow     `json:"tokens"`
	ERExtraction []string           `json:"erExtraction"`
	Compare      []*poeCompareRow   `json:"compare"`
	SignaturePOE []*poeSignatureRow `json:"signaturePOE"`
}

func TestPOEOracle(t *testing.T) {
	rows := readPOEOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty POE oracle corpus")
	}
	for _, row := range rows {
		row := row
		t.Run(row.Dump, func(t *testing.T) {
			dd := loadPOEDump(t, row.Dump)

			// The SET of token ids Init seeds. It is compared unordered on
			// purpose: two of the eleven lists Init walks - getAllSignatures()
			// and getAllRevocationData() - are java.util.HashSets upstream, so
			// their Java iteration order is String-hash order, which the Go
			// wrappers' deterministic slices neither can nor should reproduce.
			// Nothing downstream of init() depends on that order (it only seeds
			// a map keyed by token id), so the difference is invisible; see the
			// batch notes.
			seeded := poeSeededTokenIds(dd)
			if len(seeded) != len(row.Tokens) {
				t.Fatalf("seeded token count = %d, want %d", len(seeded), len(row.Tokens))
			}
			seededSet := make(map[string]struct{}, len(seeded))
			for _, id := range seeded {
				seededSet[id] = struct{}{}
			}
			for _, token := range row.Tokens {
				if _, ok := seededSet[token.Id]; !ok {
					t.Fatalf("token %q not seeded by Init", token.Id)
				}
			}

			poe := NewPOEExtraction()
			poe.Init(dd, poeControlTime)
			assertPOESnapshots(t, "afterInit", poe, row.Tokens, func(r *poeTokenRow) *poeSnapshot {
				return r.AfterInit
			})

			poe.CollectAllPOE(dd.TimestampList())
			assertPOESnapshots(t, "afterTimestamps", poe, row.Tokens, func(r *poeTokenRow) *poeSnapshot {
				return r.AfterTimestamps
			})

			evidenceRecords := dd.EvidenceRecords()
			if len(evidenceRecords) != len(row.ERExtraction) {
				t.Fatalf("evidence record count = %d, want %d", len(evidenceRecords), len(row.ERExtraction))
			}
			for i, evidenceRecord := range evidenceRecords {
				threw := extractEvidenceRecordPOERecovered(poe, evidenceRecord)
				want := row.ERExtraction[i] == "throws"
				if threw != want {
					t.Fatalf("evidence record %q extraction panicked = %v, want %v",
						evidenceRecord.Id(), threw, want)
				}
			}
			assertPOESnapshots(t, "afterEvidenceRecords", poe, row.Tokens, func(r *poeTokenRow) *poeSnapshot {
				return r.AfterEvidenceRecords
			})

			for _, token := range row.Tokens {
				for i, want := range token.Exists {
					if got := poe.IsPOEExists(token.Id, poeProbes[i]); got != want {
						t.Errorf("IsPOEExists(%q, %s) = %v, want %v",
							token.Id, poeProbes[i].Format(time.RFC3339), got, want)
					}
				}
				if token.InValidityRange != nil {
					certificate := poeCertificateById(dd, token.Id)
					got := poe.IsPOEExistInRange(token.Id, certificate.NotBefore(), certificate.NotAfter())
					if got != *token.InValidityRange {
						t.Errorf("IsPOEExistInRange(%q) = %v, want %v", token.Id, got, *token.InValidityRange)
					}
				}
			}

			// POEComparator over every ordered pair, exactly as the oracle built them.
			labels, poes := poeComparablePOEs(dd)
			comparator := NewPOEComparator()
			i := 0
			for a := range poes {
				for b := range poes {
					if i >= len(row.Compare) {
						t.Fatalf("more POE pairs than oracle rows (%d)", len(row.Compare))
					}
					want := row.Compare[i]
					if labels[a] != want.A || labels[b] != want.B {
						t.Fatalf("POE pair[%d] = (%q,%q), want (%q,%q)",
							i, labels[a], labels[b], want.A, want.B)
					}
					if got := poeSign(comparator.Compare(poes[a], poes[b])); got != want.Compare {
						t.Errorf("Compare(%s, %s) = %d, want %d", want.A, want.B, got, want.Compare)
					}
					if got := comparator.Before(poes[a], poes[b]); got != want.Before {
						t.Errorf("Before(%s, %s) = %v, want %v", want.A, want.B, got, want.Before)
					}
					i++
				}
			}
			if i != len(row.Compare) {
				t.Fatalf("POE pair count = %d, want %d", i, len(row.Compare))
			}

			// AddSignaturePOE, as ValidationProcessForSignaturesWithArchivalData step 4) does.
			sigPoe := NewPOEExtraction()
			sigPoe.Init(dd, poeControlTime)
			sigPoe.CollectAllPOE(dd.TimestampList())
			signatures := dd.AllSignatures()
			if len(signatures) != len(row.SignaturePOE) {
				t.Fatalf("signature count = %d, want %d", len(signatures), len(row.SignaturePOE))
			}
			// Matched by id, not by position: getAllSignatures() is a HashSet
			// upstream (see the seeded-token note above). Each signature's POE
			// is independent of the order they are added in.
			for _, signature := range signatures {
				sigPoe.AddSignaturePOE(signature, NewPOE(poeProbes[1]))
			}
			for _, want := range row.SignaturePOE {
				assertPOESnapshot(t, "signaturePOE:"+want.Id, sigPoe.GetLowestPOE(want.Id), want.Lowest)
			}
		})
	}
}

// extractEvidenceRecordPOERecovered runs ExtractEvidenceRecordPOE and reports
// whether it panicked - the Go form of upstream's NullPointerException for an
// evidence record whose time-stamp reference was never resolved.
func extractEvidenceRecordPOERecovered(poe *POEExtraction, evidenceRecord *diagnostic.EvidenceRecordWrapper) (threw bool) {
	defer func() {
		if recover() != nil {
			threw = true
		}
	}()
	poe.ExtractEvidenceRecordPOE(evidenceRecord)
	return false
}

func assertPOESnapshots(t *testing.T, stage string, poe *POEExtraction, tokens []*poeTokenRow,
	pick func(*poeTokenRow) *poeSnapshot) {
	t.Helper()
	for _, token := range tokens {
		assertPOESnapshot(t, stage+":"+token.Id, poe.GetLowestPOE(token.Id), pick(token))
	}
}

func assertPOESnapshot(t *testing.T, label string, got POE, want *poeSnapshot) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Errorf("%s: lowest POE = %v, want nil", label, got)
		}
		return
	}
	if got == nil {
		t.Errorf("%s: lowest POE = nil, want %+v", label, *want)
		return
	}
	if gotMillis := got.Time().UnixMilli(); gotMillis != want.Time {
		t.Errorf("%s: time = %d, want %d", label, gotMillis, want.Time)
	}
	gotProviderId := got.POEProviderId()
	switch {
	case want.ProviderId == nil && gotProviderId != nil:
		t.Errorf("%s: providerId = %q, want nil", label, *gotProviderId)
	case want.ProviderId != nil && gotProviderId == nil:
		t.Errorf("%s: providerId = nil, want %q", label, *want.ProviderId)
	case want.ProviderId != nil && *gotProviderId != *want.ProviderId:
		t.Errorf("%s: providerId = %q, want %q", label, *gotProviderId, *want.ProviderId)
	}
	if got.IsTokenProvided() != want.TokenProvided {
		t.Errorf("%s: tokenProvided = %v, want %v", label, got.IsTokenProvided(), want.TokenProvided)
	}
	if len(got.POEObjects()) != want.Objects {
		t.Errorf("%s: covered objects = %d, want %d", label, len(got.POEObjects()), want.Objects)
	}
}

// poeSeededTokenIds mirrors the id walk of POEExtraction#init, de-duplicated the
// way the oracle does it.
func poeSeededTokenIds(dd *diagnostic.DiagnosticData) []string {
	var ids []string
	for _, w := range dd.AllSignatures() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.TimestampList() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.EvidenceRecords() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.EAAs() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.UsedCertificates() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllRevocationData() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllSignerDocuments() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllOrphanCertificateObjects() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllOrphanCertificateReferences() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllOrphanRevocationObjects() {
		ids = append(ids, w.Id())
	}
	for _, w := range dd.AllOrphanRevocationReferences() {
		ids = append(ids, w.Id())
	}
	var unique []string
	seen := make(map[string]struct{})
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}

// poeComparablePOEs mirrors the oracle's comparablePOEs.
func poeComparablePOEs(dd *diagnostic.DiagnosticData) ([]string, []POE) {
	labels := []string{"control-time"}
	poes := []POE{NewPOE(poeControlTime)}
	for _, timestamp := range dd.TimestampList() {
		if timestamp.ProductionTime() != nil {
			labels = append(labels, "tst:"+timestamp.Id())
			poes = append(poes, NewTimestampPOE(timestamp))
		}
	}
	for _, evidenceRecord := range dd.EvidenceRecords() {
		if poeEvidenceRecordUsable(evidenceRecord) {
			labels = append(labels, "er:"+evidenceRecord.Id())
			poes = append(poes, NewEvidenceRecordPOE(evidenceRecord))
		}
	}
	return labels, poes
}

// poeEvidenceRecordUsable is the Go form of the oracle's guarded
// "er.getFirstTimestamp() != null && ...getProductionTime() != null", whose
// try/catch covers the unresolved-time-stamp-reference NullPointerException.
func poeEvidenceRecordUsable(evidenceRecord *diagnostic.EvidenceRecordWrapper) (usable bool) {
	defer func() {
		if recover() != nil {
			usable = false
		}
	}()
	firstTimestamp := evidenceRecord.FirstTimestamp()
	return firstTimestamp != nil && firstTimestamp.ProductionTime() != nil
}

func poeCertificateById(dd *diagnostic.DiagnosticData, id string) *diagnostic.CertificateWrapper {
	for _, certificate := range dd.UsedCertificates() {
		if certificate.Id() == id {
			return certificate
		}
	}
	return nil
}

func poeSign(value int) int {
	switch {
	case value < 0:
		return -1
	case value > 0:
		return 1
	default:
		return 0
	}
}

func readPOEOracle(t *testing.T) []*poeRow {
	t.Helper()
	file, err := os.Open(poeOracleLog)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*poeRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<26)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		row := &poeRow{}
		if err := json.Unmarshal(line, row); err != nil {
			t.Fatalf("parse oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	return rows
}

func loadPOEDump(t *testing.T, name string) *diagnostic.DiagnosticData {
	t.Helper()
	for _, dir := range []string{poeCorpusDir, poeDumpDir} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		jaxbData, err := diagnosticjaxb.Unmarshal(data)
		if err != nil {
			t.Fatalf("unmarshal %s: %v", path, err)
		}
		return diagnostic.NewDiagnosticData(jaxbData)
	}
	t.Fatalf("dump %s not found", name)
	return nil
}
