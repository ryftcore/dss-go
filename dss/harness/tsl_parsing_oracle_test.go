// Phase 9 cross-validation harness, contract item (A): LOTL/TL/pivot/MRA parse parity.
//
// testdata/oracle/tsl/tsl_parsing.jsonl is a pure Java dump, produced by
// testdata/oracle/tsl/gen/TSLParsingOracle.java against DSS 6.5.RC1's
// dss-tsl-validation TLParsingTask/LOTLParsingTask, for every fixture listed in
// testdata/oracle/tsl/tsl_parsing_manifest.tsv (kind<TAB>name<TAB>path). The fixtures are the
// upstream dss-tsl-validation src/test/resources files of the same names, copied unchanged into
// testdata/oracle/tsl/: two real country TLs (de-tl.xml, fr.xml) plus dk_tl-sn21.xml, sk-tl.xml,
// four LOTL variants (eu-lotl-250.xml the full real LOTL, eu-lotl.xml, eu-lotl-pivot.xml,
// eu-lotl-no-sig.xml), the two error-classification fixtures (eu-lotl-broken-sig.xml - parses
// fine, signature is what's broken; eu-lotl-not-parseable.xml - must fail to parse on both
// sides), and the MRA pair (mra-lotl.xml, mra-zz-tl.xml).
//
// This test parses the SAME bytes with NewTLParsingTask/NewLOTLParsingTask over a default
// TLSource/LOTLSource (LOTLSource additionally gets SetMraSupport(true) for the two mra-*
// fixtures, matching the generator), then compares every field TLParsingResult/LOTLParsingResult
// expose against the Java dump: TSLType, sequence number, version, territory, issue/next-update
// dates, distribution points; every TrustServiceProvider (names, trade names, registration
// identifiers, postal/electronic addresses, territory) and every TrustService's full
// time-dependent history (names, type, status, additional service info URIs, service supply
// points, expired-certs-revocation-info, certificate SHA-256 digests, and every
// ConditionForQualifiers rendered through Condition.ToString("") - already a byte-exact port of
// Java's Condition#toString(String), so this is a genuine cross-language comparison, not a
// self-check); for LOTLs, every OtherTSLPointer (lotlPointers/tlPointers) including its MRA
// block (technical type, version, both legislations, and every ServiceEquivalence in the
// equivalence history: legal info identifier, status, type-ASi/qualifier maps, status-equivalence
// pairs, certificate-content equivalences with their Condition tree and QCStatementOids).
//
// Map-shaped fields (names/tradeNames/electronicAddresses/postalAddresses per TSP or service,
// and the two MRA equivalence maps) are dumped through Go's encoding/json, which sorts map
// string keys, and the Java generator dumps them through a TreeMap for the same reason: Java
// HashMap iteration order is unspecified and unrelated to any real defect, so comparing it would
// manufacture false positives rather than catch anything. Every list-shaped field (TSP order,
// service history order, distribution points, pivot URLs, positional certificate lists) is left
// in encounter order on both sides and compared positionally - both engines unmarshal the same
// XML bytes in document order, so a positional mismatch there IS a real defect (e.g. the kind
// abstract_parsing_task_oracle_test.go's map-iteration-order lessons warn about).
package harness

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/tsl"
)

// tpNullInt is the sentinel both the Go conversion and the Java generator use for a null/nil
// Integer field, so the JSON never needs a separate null representation for this int-only value.
const tpNullInt = -2147483648

type tpRecord struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	ParseError bool   `json:"parseError"`

	TSLType            string   `json:"tslType"`
	SequenceNumber     int      `json:"sequenceNumber"`
	Version            int      `json:"version"`
	Territory          string   `json:"territory"`
	IssueDate          int64    `json:"issueDate"`
	NextUpdateDate     int64    `json:"nextUpdateDate"`
	DistributionPoints []string `json:"distributionPoints"`

	TSPs []tpTSP `json:"tsps"`

	LotlPointers               []tpPointer `json:"lotlPointers"`
	TlPointers                 []tpPointer `json:"tlPointers"`
	SigningCertAnnouncementURL string      `json:"signingCertAnnouncementURL"`
	PivotURLs                  []string    `json:"pivotURLs"`
}

type tpTSP struct {
	Names                   map[string][]string `json:"names"`
	TradeNames              map[string][]string `json:"tradeNames"`
	RegistrationIdentifiers []string            `json:"registrationIdentifiers"`
	PostalAddresses         map[string]string   `json:"postalAddresses"`
	ElectronicAddresses     map[string][]string `json:"electronicAddresses"`
	Information             map[string]string   `json:"information"`
	Territory               string              `json:"territory"`
	Services                []tpService         `json:"services"`
}

type tpService struct {
	Start                      int64               `json:"start"`
	End                        int64               `json:"end"`
	Names                      map[string][]string `json:"names"`
	Type                       string              `json:"type"`
	Status                     string              `json:"status"`
	AdditionalServiceInfoUris  []string            `json:"additionalServiceInfoUris"`
	ServiceSupplyPoints        []string            `json:"serviceSupplyPoints"`
	ExpiredCertsRevocationInfo int64               `json:"expiredCertsRevocationInfo"`
	CertificatesSha256         []string            `json:"certificatesSha256"`
	Conditions                 []tpCondition       `json:"conditions"`
}

type tpCondition struct {
	Qualifiers []string `json:"qualifiers"`
	Critical   bool     `json:"critical"`
	Condition  string   `json:"condition"`
}

type tpPointer struct {
	Location                 string              `json:"location"`
	TSLLocation              string              `json:"tslLocation"`
	SchemeTerritory          string              `json:"schemeTerritory"`
	TSLType                  string              `json:"tslType"`
	MimeType                 string              `json:"mimeType"`
	SchemeOperatorNames      map[string][]string `json:"schemeOperatorNames"`
	SchemeTypeCommunityRules map[string][]string `json:"schemeTypeCommunityRules"`
	SdiCertificatesSha256    []string            `json:"sdiCertificatesSha256"`
	MRA                      *tpMRA              `json:"mra"`
}

type tpMRA struct {
	TechnicalType       string                 `json:"technicalType"`
	Version             string                 `json:"version"`
	PointingLegislation string                 `json:"pointingLegislation"`
	PointedLegislation  string                 `json:"pointedLegislation"`
	ServiceEquivalences []tpServiceEquivalence `json:"serviceEquivalences"`
}

type tpServiceEquivalence struct {
	Start                          int64                        `json:"start"`
	End                            int64                        `json:"end"`
	LegalInfoIdentifier            string                       `json:"legalInfoIdentifier"`
	Status                         string                       `json:"status"`
	TypeAsiEquivalence             map[string]string            `json:"typeAsiEquivalence"`
	StatusEquivalence              []tpStatusEquivalenceMapping `json:"statusEquivalence"`
	CertificateContentEquivalences []tpCertContentEquivalence   `json:"certificateContentEquivalences"`
	QualifierEquivalence           map[string]string            `json:"qualifierEquivalence"`
}

type tpStatusEquivalenceMapping struct {
	PointedStatuses  []string `json:"pointedStatuses"`
	PointingStatuses []string `json:"pointingStatuses"`
}

type tpCertContentEquivalence struct {
	Context            string             `json:"context"`
	Condition          string             `json:"condition"`
	ContentReplacement *tpQCStatementOids `json:"contentReplacement"`
}

type tpQCStatementOids struct {
	QcStatementIds           []string `json:"qcStatementIds"`
	QcTypeIds                []string `json:"qcTypeIds"`
	QcCClegislations         []string `json:"qcCClegislations"`
	QcStatementIdsToRemove   []string `json:"qcStatementIdsToRemove"`
	QcTypeIdsToRemove        []string `json:"qcTypeIdsToRemove"`
	QcCClegislationsToRemove []string `json:"qcCClegislationsToRemove"`
}

// tpTSLTypeURI renders a TSLType as its URI (the literal XML value both engines parse from), or
// "" when unset. TSLType is an interface (enumerations.UriBasedEnum), not a plain string type.
func tpTSLTypeURI(t interface{ URI() string }) string {
	if t == nil {
		return ""
	}
	return t.URI()
}

func tpEpochMillis(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

func tpIntOrNull(p *int) int {
	if p == nil {
		return tpNullInt
	}
	return *p
}

// tpSorted returns a sorted copy of ss, never nil (Go's encoding/json renders a nil slice as
// JSON null but an empty slice as [], and the Java generator always emits [] for an empty/null
// List - see dumpStrList in TSLParsingOracle.java - so nil here would be a harness-only false
// mismatch, not a real defect).
func tpSorted(ss []string) []string {
	out := make([]string, 0, len(ss))
	out = append(out, ss...)
	sort.Strings(out)
	return out
}

// tpNonNilStrings mirrors tpSorted's nil-vs-empty normalization for fields kept in encounter
// order (not sorted).
func tpNonNilStrings(ss []string) []string {
	if ss == nil {
		return []string{}
	}
	return ss
}

func tpSortedMapLists(m map[string][]string) map[string][]string {
	out := make(map[string][]string, len(m))
	for k, v := range m {
		out[k] = tpSorted(v)
	}
	return out
}

func tpCertSha256(certs []*model.CertificateToken) []string {
	out := make([]string, len(certs))
	for i, c := range certs {
		sum := sha256.Sum256(c.Encoded())
		out[i] = hex.EncodeToString(sum[:])
	}
	return out
}

func tpConditionsForQualifiers(cfqs []*tslmodel.ConditionForQualifiers) []tpCondition {
	out := make([]tpCondition, len(cfqs))
	for i, c := range cfqs {
		condStr := ""
		if c.Condition() != nil {
			condStr = c.Condition().ToString("")
		}
		out[i] = tpCondition{
			Qualifiers: tpSorted(c.Qualifiers()),
			Critical:   c.IsCritical(),
			Condition:  condStr,
		}
	}
	return out
}

func tpServices(services []*tslmodel.TrustService) []tpService {
	var out []tpService
	for _, svc := range services {
		for entry := range svc.StatusAndInformationExtensions().Iterator() {
			out = append(out, tpService{
				Start:                      tpEpochMillis(entry.StartDate()),
				End:                        tpEpochMillis(entry.EndDate()),
				Names:                      tpSortedMapLists(entry.Names()),
				Type:                       entry.Type(),
				Status:                     entry.Status(),
				AdditionalServiceInfoUris:  tpSorted(entry.AdditionalServiceInfoUris()),
				ServiceSupplyPoints:        tpSorted(entry.ServiceSupplyPoints()),
				ExpiredCertsRevocationInfo: tpEpochMillis(entry.ExpiredCertsRevocationInfo()),
				CertificatesSha256:         tpCertSha256(svc.Certificates()),
				Conditions:                 tpConditionsForQualifiers(entry.ConditionsForQualifiers()),
			})
		}
	}
	if out == nil {
		out = []tpService{}
	}
	return out
}

func tpTSPs(tsps []*tslmodel.TrustServiceProvider) []tpTSP {
	out := make([]tpTSP, len(tsps))
	for i, t := range tsps {
		out[i] = tpTSP{
			Names:                   tpSortedMapLists(t.Names()),
			TradeNames:              tpSortedMapLists(t.TradeNames()),
			RegistrationIdentifiers: tpSorted(t.RegistrationIdentifiers()),
			PostalAddresses:         t.PostalAddresses(),
			ElectronicAddresses:     tpSortedMapLists(t.ElectronicAddresses()),
			Information:             t.Information(),
			Territory:               t.Territory(),
			Services:                tpServices(t.Services()),
		}
	}
	return out
}

func tpQCStatementOidsFrom(q *tslmodel.QCStatementOids) *tpQCStatementOids {
	if q == nil {
		return nil
	}
	return &tpQCStatementOids{
		QcStatementIds:           tpSorted(q.QcStatementIds()),
		QcTypeIds:                tpSorted(q.QcTypeIds()),
		QcCClegislations:         tpSorted(q.QcCClegislations()),
		QcStatementIdsToRemove:   tpSorted(q.QcStatementIdsToRemove()),
		QcTypeIdsToRemove:        tpSorted(q.QcTypeIdsToRemove()),
		QcCClegislationsToRemove: tpSorted(q.QcCClegislationsToRemove()),
	}
}

func tpCertContentEquivalences(ccs []*tslmodel.CertificateContentEquivalence) []tpCertContentEquivalence {
	out := make([]tpCertContentEquivalence, len(ccs))
	for i, c := range ccs {
		condStr := ""
		if c.Condition() != nil {
			condStr = c.Condition().ToString("")
		}
		out[i] = tpCertContentEquivalence{
			Context:            string(c.Context()),
			Condition:          condStr,
			ContentReplacement: tpQCStatementOidsFrom(c.ContentReplacement()),
		}
	}
	return out
}

func tpStatusEquivalence(ms []tslmodel.StatusEquivalenceMapping) []tpStatusEquivalenceMapping {
	out := make([]tpStatusEquivalenceMapping, len(ms))
	for i, m := range ms {
		out[i] = tpStatusEquivalenceMapping{
			PointedStatuses:  tpSorted(m.PointedStatuses),
			PointingStatuses: tpSorted(m.PointingStatuses),
		}
	}
	return out
}

func tpMRAFrom(mra *tslmodel.MRA) *tpMRA {
	if mra == nil {
		return nil
	}
	out := &tpMRA{
		TechnicalType:       mra.TechnicalType(),
		Version:             mra.Version(),
		PointingLegislation: mra.PointingContractingPartyLegislation(),
		PointedLegislation:  mra.PointedContractingPartyLegislation(),
	}
	for _, mtdv := range mra.ServiceEquivalence() {
		if mtdv == nil {
			continue
		}
		for _, se := range mtdv.List() {
			out.ServiceEquivalences = append(out.ServiceEquivalences, tpServiceEquivalence{
				Start:                          tpEpochMillis(se.StartDate()),
				End:                            tpEpochMillis(se.EndDate()),
				LegalInfoIdentifier:            se.LegalInfoIdentifier(),
				Status:                         string(se.Status()),
				TypeAsiEquivalence:             tpServiceTypeASiMap(se.TypeAsiEquivalence()),
				StatusEquivalence:              tpStatusEquivalence(se.StatusEquivalence()),
				CertificateContentEquivalences: tpCertContentEquivalences(se.CertificateContentEquivalences()),
				QualifierEquivalence:           se.QualifierEquivalence(),
			})
		}
	}
	if out.ServiceEquivalences == nil {
		out.ServiceEquivalences = []tpServiceEquivalence{}
	}
	return out
}

func tpServiceTypeASiMap(m map[tslmodel.ServiceTypeASi]tslmodel.ServiceTypeASi) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		kk, vv := k, v
		out[kk.Type()+"|"+kk.Asi()] = vv.Type() + "|" + vv.Asi()
	}
	return out
}

func tpPointers(ptrs []*tslmodel.OtherTSLPointer) []tpPointer {
	out := make([]tpPointer, len(ptrs))
	for i, p := range ptrs {
		out[i] = tpPointer{
			Location:                 p.Location(),
			TSLLocation:              p.TSLLocation(),
			SchemeTerritory:          p.SchemeTerritory(),
			TSLType:                  p.TslType(),
			MimeType:                 p.MimeType(),
			SchemeOperatorNames:      tpSortedMapLists(p.SchemeOperatorNames()),
			SchemeTypeCommunityRules: tpSortedMapLists(p.SchemeTypeCommunityRules()),
			SdiCertificatesSha256:    tpCertSha256(p.SdiCertificates()),
			MRA:                      tpMRAFrom(p.Mra()),
		}
	}
	return out
}

// tpParse parses one fixture with the Go port, mirroring exactly what
// testdata/oracle/tsl/gen/TSLParsingOracle.java does with TLParsingTask/LOTLParsingTask.
func tpParse(t *testing.T, kind, path string) tpRecord {
	t.Helper()
	document, err := model.NewFileDocument(path)
	if err != nil {
		t.Fatalf("NewFileDocument(%s): %v", path, err)
	}

	rec := tpRecord{Kind: kind}

	if kind == "LOTL" {
		source := tsl.NewLOTLSource()
		if filepath.Base(path) == "mra-lotl.xml" {
			source.SetMraSupport(true)
		}
		result, err := tsl.NewLOTLParsingTask(document, source).Get()
		if err != nil {
			rec.ParseError = true
			return rec
		}
		rec.TSLType = tpTSLTypeURI(result.TSLType())
		rec.SequenceNumber = tpIntOrNull(result.SequenceNumber())
		rec.Version = tpIntOrNull(result.Version())
		rec.Territory = result.Territory()
		rec.IssueDate = tpEpochMillis(result.IssueDate())
		rec.NextUpdateDate = tpEpochMillis(result.NextUpdateDate())
		rec.DistributionPoints = tpNonNilStrings(result.DistributionPoints())
		rec.TSPs = []tpTSP{}
		rec.LotlPointers = tpPointers(result.LotlPointers())
		rec.TlPointers = tpPointers(result.TlPointers())
		rec.SigningCertAnnouncementURL = result.SigningCertificateAnnouncementURL()
		rec.PivotURLs = tpNonNilStrings(result.PivotURLs())
		return rec
	}

	source := tsl.NewTLSource()
	result, err := tsl.NewTLParsingTask(document, source).Get()
	if err != nil {
		rec.ParseError = true
		return rec
	}
	rec.TSLType = tpTSLTypeURI(result.TSLType())
	rec.SequenceNumber = tpIntOrNull(result.SequenceNumber())
	rec.Version = tpIntOrNull(result.Version())
	rec.Territory = result.Territory()
	rec.IssueDate = tpEpochMillis(result.IssueDate())
	rec.NextUpdateDate = tpEpochMillis(result.NextUpdateDate())
	rec.DistributionPoints = tpNonNilStrings(result.DistributionPoints())
	rec.TSPs = tpTSPs(result.TrustServiceProviders())
	rec.LotlPointers = []tpPointer{}
	rec.TlPointers = []tpPointer{}
	rec.PivotURLs = []string{}
	return rec
}

// tpManifestRow is one row of tsl_parsing_manifest.tsv (kind, name, path).
type tpManifestRow struct {
	Kind string
	Name string
	Path string
}

func tpReadManifest(t *testing.T) []tpManifestRow {
	t.Helper()
	f, err := os.Open("testdata/oracle/tsl/tsl_parsing_manifest.tsv")
	if err != nil {
		t.Fatalf("open manifest: %v", err)
	}
	defer f.Close()

	var rows []tpManifestRow
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := tpSplitTab(line)
		if len(parts) != 3 {
			t.Fatalf("malformed manifest line: %q", line)
		}
		rows = append(rows, tpManifestRow{Kind: parts[0], Name: parts[1], Path: parts[2]})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan manifest: %v", err)
	}
	return rows
}

func tpSplitTab(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\t' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

func tpReadOracle(t *testing.T) map[string]tpRecord {
	t.Helper()
	f, err := os.Open("testdata/oracle/tsl/tsl_parsing.jsonl")
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()

	out := make(map[string]tpRecord)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec tpRecord
		if err := json.Unmarshal(line, &rec); err != nil {
			t.Fatalf("unmarshal oracle line: %v\n%s", err, line)
		}
		out[rec.Name] = rec
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan oracle: %v", err)
	}
	return out
}

// TestTSLParsingOracle is the Phase 9 harness contract item (A).
func TestTSLParsingOracle(t *testing.T) {
	rows := tpReadManifest(t)
	oracle := tpReadOracle(t)

	if len(rows) == 0 {
		t.Fatal("empty manifest")
	}

	for _, row := range rows {
		row := row
		t.Run(row.Name, func(t *testing.T) {
			want, ok := oracle[row.Name]
			if !ok {
				t.Fatalf("no oracle record for %s", row.Name)
			}
			got := tpParse(t, row.Kind, row.Path)
			got.Name = row.Name

			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("marshal got: %v", err)
			}
			// Re-marshal want through the same Go struct (round-tripped from the oracle JSON) so
			// both sides go through identical map-key sorting and field ordering before the
			// byte comparison.
			wantJSON, err := json.Marshal(want)
			if err != nil {
				t.Fatalf("marshal want: %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("mismatch for %s (kind=%s):\n GOT: %s\nWANT: %s", row.Name, row.Kind, gotJSON, wantJSON)
			}
		})
	}
}
