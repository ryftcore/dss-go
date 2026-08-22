package sav

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/policy"
)

// The SAV KAT: every row of testdata/oracle/sav_blocks.jsonl is the XmlSAV
// upstream's SignatureAcceptanceValidation / TimestampAcceptanceValidation /
// RevocationAcceptanceValidation produces for one token of one XmlDiagnosticData
// dump of the marshal-parity corpus, under the default ETSI validation policy and
// a PASSED XmlAOV - see ../testdata/gen/FcSavOracle.java. This test replays the
// same tokens and compares the XmlConstraint sequence and the XmlConclusion.

// savCorpusDir is the marshal-parity diagnostic-data corpus the oracle was run over.
const savCorpusDir = "diagnostic/jaxb/testdata/oracle" // corpus/-relative (internal/corpustest)

// savCurrentTime is the fixed validation time the oracle ran with,
// 2024-01-01T00:00:00Z.
var savCurrentTime = time.Unix(1704067200, 0).UTC()

type savOracleMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type savOracleConclusion struct {
	Indication    *string             `json:"indication"`
	SubIndication *string             `json:"subIndication"`
	Errors        []*savOracleMessage `json:"errors"`
	Warnings      []*savOracleMessage `json:"warnings"`
	Infos         []*savOracleMessage `json:"infos"`
}

type savOracleConstraint struct {
	Name           *savOracleMessage `json:"name"`
	Status         *string           `json:"status"`
	Error          *savOracleMessage `json:"error"`
	Warning        *savOracleMessage `json:"warning"`
	Info           *savOracleMessage `json:"info"`
	AdditionalInfo *string           `json:"additionalInfo"`
	Id             *string           `json:"id"`
	BlockType      *string           `json:"blockType"`
}

type savOracleRow struct {
	File        string                 `json:"file"`
	Token       string                 `json:"token"`
	Context     string                 `json:"context"`
	Block       string                 `json:"block"`
	Title       *string                `json:"title"`
	Conclusion  *savOracleConclusion   `json:"conclusion"`
	Constraints []*savOracleConstraint `json:"constraints"`
}

func loadSAVRows(t *testing.T, path string) []*savOracleRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*savOracleRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &savOracleRow{}
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

func loadSAVDiagnosticData(t *testing.T, name string) *diagnostic.Data {
	t.Helper()
	data, err := os.ReadFile(corpustest.RootPath(t, filepath.Join(savCorpusDir, name)))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	jaxbData, err := diagnosticjaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return diagnostic.NewData(jaxbData)
}

// passedAOV is the PASSED Algorithm Obsolescence Validation result the oracle fed
// the chains: AlgorithmObsolescenceValidationCheck reads its conclusion to pick
// its own Level and to decide process().
func passedAOV() *jaxb.XmlAOV {
	aov := &jaxb.XmlAOV{}
	aov.Conclusion = &jaxb.XmlConclusion{
		Indication: jaxb.IndicationValue(enumerations.IndicationPassed),
	}
	return aov
}

func toSAVRow(file, token string, context enumerations.Context,
	content *jaxb.XmlConstraintsConclusionContent, title string) *savOracleRow {
	row := &savOracleRow{File: file, Token: token, Context: string(context), Block: "SAV"}
	titleValue := title
	row.Title = &titleValue
	if content.Conclusion != nil {
		row.Conclusion = &savOracleConclusion{
			Errors:   toSAVMessages(content.Conclusion.Errors),
			Warnings: toSAVMessages(content.Conclusion.Warnings),
			Infos:    toSAVMessages(content.Conclusion.Infos),
		}
		if indication := string(content.Conclusion.Indication.Indication()); indication != "" {
			row.Conclusion.Indication = &indication
		}
		if content.Conclusion.SubIndication != nil {
			subIndication := string(content.Conclusion.SubIndication.SubIndication())
			row.Conclusion.SubIndication = &subIndication
		}
	}
	row.Constraints = make([]*savOracleConstraint, 0, len(content.Constraint))
	for _, constraint := range content.Constraint {
		status := string(constraint.Status)
		converted := &savOracleConstraint{
			Name:           toSAVMessage(constraint.Name),
			Status:         &status,
			Error:          toSAVMessage(constraint.Error),
			Warning:        toSAVMessage(constraint.Warning),
			Info:           toSAVMessage(constraint.Info),
			AdditionalInfo: constraint.AdditionalInfo,
			Id:             constraint.Id,
		}
		if constraint.BlockType != nil {
			blockType := string(*constraint.BlockType)
			converted.BlockType = &blockType
		}
		row.Constraints = append(row.Constraints, converted)
	}
	return row
}

func toSAVMessages(messages []*jaxb.XmlMessage) []*savOracleMessage {
	converted := make([]*savOracleMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toSAVMessage(message))
	}
	return converted
}

func toSAVMessage(message *jaxb.XmlMessage) *savOracleMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &savOracleMessage{Key: message.Key, Value: &value}
}

func mustSAVJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

// savSafeExecute runs fn, converting a panic into a marker so that a token the
// Java oracle itself could not process (it recorded no row) does not abort the run.
func savSafeExecute(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func TestAcceptanceValidationAgainstJavaOracle(t *testing.T) {
	rows := loadSAVRows(t, corpustest.Path(t, "oracle/sav_blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	validationPolicy := policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()
	i18nProvider := i18n.NewProvider()

	byFile := map[string]map[string]*savOracleRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]*savOracleRow{}
		}
		byFile[row.File][row.Token] = row
	}

	statuses := map[string]int{}
	constraintNames := map[string]int{}
	matched := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadSAVDiagnosticData(t, name)
			expected := byFile[name]

			compare := func(token string, context enumerations.Context, run func() *jaxb.XmlSAV) {
				want, ok := expected[token]
				if !ok {
					return // the Java oracle threw on this token and recorded no row
				}
				var result *jaxb.XmlSAV
				if savSafeExecute(func() { result = run() }) {
					t.Errorf("token %s: Go panicked where Java produced a row", token)
					return
				}
				got := toSAVRow(name, token, context,
					&result.XmlConstraintsConclusionContent, result.Title)
				matched++
				for _, constraint := range want.Constraints {
					if constraint.Status != nil {
						statuses[*constraint.Status]++
					}
					if constraint.Name != nil && constraint.Name.Key != nil {
						constraintNames[*constraint.Name.Key]++
					}
				}
				if !reflect.DeepEqual(want, got) {
					t.Errorf("token %s: SAV mismatch\nexpected: %s\nactual:   %s",
						token, mustSAVJSON(t, want), mustSAVJSON(t, got))
				}
			}

			for _, signature := range diagnosticData.Signatures() {
				context := enumerations.ContextSignature
				if signature.IsCounterSignature() {
					context = enumerations.ContextCounterSignature
				}
				sig := signature
				compare(sig.Id(), context, func() *jaxb.XmlSAV {
					return NewSignatureAcceptanceValidation(i18nProvider, diagnosticData, savCurrentTime,
						sig, context, map[string]*jaxb.XmlBasicBuildingBlocks{}, passedAOV(), validationPolicy).Execute()
				})
			}
			for _, timestamp := range diagnosticData.TimestampList() {
				tst := timestamp
				compare(tst.Id(), enumerations.ContextTimestamp, func() *jaxb.XmlSAV {
					return NewTimestampAcceptanceValidation(i18nProvider, savCurrentTime, tst,
						passedAOV(), validationPolicy).Execute()
				})
			}
			revocations := diagnosticData.AllRevocationData()
			sort.SliceStable(revocations, func(i, j int) bool { return revocations[i].Id() < revocations[j].Id() })
			for _, revocation := range revocations {
				rev := revocation
				compare(rev.Id(), enumerations.ContextRevocation, func() *jaxb.XmlSAV {
					return NewRevocationAcceptanceValidation(i18nProvider, savCurrentTime, rev,
						passedAOV(), validationPolicy).Execute()
				})
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("matched %d of %d oracle rows", matched, len(rows))
	}
	for _, status := range []string{"OK", "NOT OK"} {
		if statuses[status] == 0 {
			t.Errorf("the corpus never produced a %q SAV constraint", status)
		}
	}
	t.Logf("SAV constraints exercised by the corpus (%d rows): %v", matched, constraintNames)
}
