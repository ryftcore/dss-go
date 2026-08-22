package fc

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/policy"
)

// The FC KAT: every row of testdata/oracle/fc_blocks.jsonl is the XmlFC upstream's
// SignatureFormatChecking / TimestampFormatChecking produces for one token of one
// XmlDiagnosticData dump of the marshal-parity corpus, under the default ETSI
// validation policy - see ../testdata/gen/FcSavOracle.java. This test replays the
// same tokens and compares the XmlConstraint sequence and the XmlConclusion.

// fcCorpusRoot is the module-root-relative path (inside the external
// corpus/ tree, see internal/corpustest) of the marshal-parity
// diagnostic-data corpus the oracle was run over.
const fcCorpusRoot = "diagnostic/jaxb/testdata/oracle"

type fcOracleMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type fcOracleConclusion struct {
	Indication    *string            `json:"indication"`
	SubIndication *string            `json:"subIndication"`
	Errors        []*fcOracleMessage `json:"errors"`
	Warnings      []*fcOracleMessage `json:"warnings"`
	Infos         []*fcOracleMessage `json:"infos"`
}

type fcOracleConstraint struct {
	Name           *fcOracleMessage `json:"name"`
	Status         *string          `json:"status"`
	Error          *fcOracleMessage `json:"error"`
	Warning        *fcOracleMessage `json:"warning"`
	Info           *fcOracleMessage `json:"info"`
	AdditionalInfo *string          `json:"additionalInfo"`
	Id             *string          `json:"id"`
	BlockType      *string          `json:"blockType"`
}

type fcOracleRow struct {
	File        string                `json:"file"`
	Token       string                `json:"token"`
	Context     string                `json:"context"`
	Block       string                `json:"block"`
	Title       *string               `json:"title"`
	Conclusion  *fcOracleConclusion   `json:"conclusion"`
	Constraints []*fcOracleConstraint `json:"constraints"`
}

func loadFCRows(t *testing.T, path string) []*fcOracleRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*fcOracleRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &fcOracleRow{}
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

// loadFCDiagnosticData reads a dump through jaxb.Unmarshal, which links the IDREF
// graph the wrappers navigate.
func loadFCDiagnosticData(t *testing.T, name string) *diagnostic.Data {
	t.Helper()
	data, err := os.ReadFile(corpustest.RootPath(t, filepath.Join(fcCorpusRoot, name)))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	jaxbData, err := diagnosticjaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return diagnostic.NewData(jaxbData)
}

func toFCRow(file, token string, context enumerations.Context, block string,
	content *jaxb.XmlConstraintsConclusionContent, title string) *fcOracleRow {
	row := &fcOracleRow{File: file, Token: token, Context: string(context), Block: block}
	titleValue := title
	row.Title = &titleValue
	if content.Conclusion != nil {
		row.Conclusion = &fcOracleConclusion{
			Errors:   toFCMessages(content.Conclusion.Errors),
			Warnings: toFCMessages(content.Conclusion.Warnings),
			Infos:    toFCMessages(content.Conclusion.Infos),
		}
		if indication := string(content.Conclusion.Indication.Indication()); indication != "" {
			row.Conclusion.Indication = &indication
		}
		if content.Conclusion.SubIndication != nil {
			subIndication := string(content.Conclusion.SubIndication.SubIndication())
			row.Conclusion.SubIndication = &subIndication
		}
	}
	row.Constraints = make([]*fcOracleConstraint, 0, len(content.Constraint))
	for _, constraint := range content.Constraint {
		status := string(constraint.Status)
		converted := &fcOracleConstraint{
			Name:           toFCMessage(constraint.Name),
			Status:         &status,
			Error:          toFCMessage(constraint.Error),
			Warning:        toFCMessage(constraint.Warning),
			Info:           toFCMessage(constraint.Info),
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

func toFCMessages(messages []*jaxb.XmlMessage) []*fcOracleMessage {
	converted := make([]*fcOracleMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toFCMessage(message))
	}
	return converted
}

func toFCMessage(message *jaxb.XmlMessage) *fcOracleMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &fcOracleMessage{Key: message.Key, Value: &value}
}

func mustFCJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

// safeExecute runs fn, converting a panic into an error marker so that a token the
// Java oracle itself could not process (it recorded no row) does not abort the run.
func safeExecute(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func TestFormatCheckingAgainstJavaOracle(t *testing.T) {
	rows := loadFCRows(t, corpustest.Path(t, "oracle/fc_blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	validationPolicy := policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()
	i18nProvider := i18n.NewProvider()

	byFile := map[string]map[string]*fcOracleRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]*fcOracleRow{}
		}
		byFile[row.File][row.Token] = row
	}

	statuses := map[string]int{}
	constraintNames := map[string]int{}
	matched := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadFCDiagnosticData(t, name)
			expected := byFile[name]

			compare := func(token string, context enumerations.Context, run func() *jaxb.XmlFC) {
				want, ok := expected[token]
				if !ok {
					// The Java oracle threw on this token and recorded no row.
					return
				}
				var result *jaxb.XmlFC
				if safeExecute(func() { result = run() }) {
					t.Errorf("token %s: Go panicked where Java produced a row", token)
					return
				}
				got := toFCRow(name, token, context, "FC",
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
					t.Errorf("token %s: FC mismatch\nexpected: %s\nactual:   %s",
						token, mustFCJSON(t, want), mustFCJSON(t, got))
				}
			}

			for _, signature := range diagnosticData.Signatures() {
				context := enumerations.ContextSignature
				if signature.IsCounterSignature() {
					context = enumerations.ContextCounterSignature
				}
				sig := signature
				compare(sig.Id(), context, func() *jaxb.XmlFC {
					return NewSignatureFormatChecking(i18nProvider, diagnosticData, sig, context, validationPolicy).Execute()
				})
			}
			for _, timestamp := range diagnosticData.TimestampList() {
				tst := timestamp
				compare(tst.Id(), enumerations.ContextTimestamp, func() *jaxb.XmlFC {
					return NewTimestampFormatChecking(i18nProvider, diagnosticData, tst,
						enumerations.ContextTimestamp, validationPolicy).Execute()
				})
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("matched %d of %d oracle rows", matched, len(rows))
	}
	for _, status := range []string{"OK", "NOT OK"} {
		if statuses[status] == 0 {
			t.Errorf("the corpus never produced a %q FC constraint", status)
		}
	}
	t.Logf("FC constraints exercised by the corpus (%d rows): %v", matched, constraintNames)
}
