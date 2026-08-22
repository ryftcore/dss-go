package isc

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/policy"
)

// The ISC KAT: every row of testdata/oracle/isc_blocks.jsonl is the XmlISC
// upstream's IdentificationOfTheSigningCertificate produces for one token of one
// XmlDiagnosticData dump of the marshal-parity corpus, under the default ETSI
// validation policy - see ../testdata/gen/BbbBlocksOracle.java. This test replays
// the same tokens in the same order and compares.

// corpusDir is the marshal-parity diagnostic-data corpus the oracle was run over.
const corpusDir = "diagnostic/jaxb/testdata/oracle" // corpus/-relative (internal/corpustest)

// i18nProviderForTests is the provider both KATs run with, as the oracle does.
var i18nProviderForTests = i18n.NewProvider()

type oracleMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type oracleConclusion struct {
	Indication    *string          `json:"indication"`
	SubIndication *string          `json:"subIndication"`
	Errors        []*oracleMessage `json:"errors"`
	Warnings      []*oracleMessage `json:"warnings"`
	Infos         []*oracleMessage `json:"infos"`
}

type oracleConstraint struct {
	Name           *oracleMessage `json:"name"`
	Status         *string        `json:"status"`
	Error          *oracleMessage `json:"error"`
	Warning        *oracleMessage `json:"warning"`
	Info           *oracleMessage `json:"info"`
	AdditionalInfo *string        `json:"additionalInfo"`
	Id             *string        `json:"id"`
	BlockType      *string        `json:"blockType"`
}

type oracleChainItem struct {
	Id     *string `json:"id"`
	Source *string `json:"source"`
}

type oracleRow struct {
	File             string              `json:"file"`
	Token            string              `json:"token"`
	Context          string              `json:"context"`
	Block            string              `json:"block"`
	Title            *string             `json:"title"`
	Conclusion       *oracleConclusion   `json:"conclusion"`
	Constraints      []*oracleConstraint `json:"constraints"`
	CertificateChain []*oracleChainItem  `json:"certificateChain"`
}

func loadRows(t *testing.T, path string) []*oracleRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*oracleRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &oracleRow{}
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

// loadDiagnosticData reads a dump through jaxb.Unmarshal, not through
// DataFacade.Unmarshal: only the former links the IDREF graph, and
// without that link the chain items carry no certificate (see notes).
func loadDiagnosticData(t *testing.T, name string) *diagnostic.Data {
	t.Helper()
	data, err := os.ReadFile(corpustest.RootPath(t, filepath.Join(corpusDir, name)))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	jaxbData, err := diagnosticjaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return diagnostic.NewData(jaxbData)
}

func toRow(file, token string, context enumerations.Context, block string,
	content *jaxb.XmlConstraintsConclusionContent, title string) *oracleRow {
	row := &oracleRow{File: file, Token: token, Context: string(context), Block: block}
	titleValue := title
	row.Title = &titleValue
	if content.Conclusion != nil {
		row.Conclusion = &oracleConclusion{
			Errors:   toMessages(content.Conclusion.Errors),
			Warnings: toMessages(content.Conclusion.Warnings),
			Infos:    toMessages(content.Conclusion.Infos),
		}
		if indication := string(content.Conclusion.Indication.Indication()); indication != "" {
			row.Conclusion.Indication = &indication
		}
		if content.Conclusion.SubIndication != nil {
			subIndication := string(content.Conclusion.SubIndication.SubIndication())
			row.Conclusion.SubIndication = &subIndication
		}
	}
	row.Constraints = make([]*oracleConstraint, 0, len(content.Constraint))
	for _, constraint := range content.Constraint {
		status := string(constraint.Status)
		converted := &oracleConstraint{
			Name:           toMessage(constraint.Name),
			Status:         &status,
			Error:          toMessage(constraint.Error),
			Warning:        toMessage(constraint.Warning),
			Info:           toMessage(constraint.Info),
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

func toMessages(messages []*jaxb.XmlMessage) []*oracleMessage {
	converted := make([]*oracleMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toMessage(message))
	}
	return converted
}

func toMessage(message *jaxb.XmlMessage) *oracleMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &oracleMessage{Key: message.Key, Value: &value}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

func TestIdentificationOfTheSigningCertificateAgainstJavaOracle(t *testing.T) {
	rows := loadRows(t, corpustest.Path(t, "oracle/isc_blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	i18nProvider := i18nProviderForTests
	validationPolicy := policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()

	// Group the rows by dump, keeping the oracle's order within a dump: the Go
	// enumeration below has to reproduce it token for token.
	byFile := map[string][]*oracleRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
		}
		byFile[row.File] = append(byFile[row.File], row)
	}

	statuses := map[string]int{}

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadDiagnosticData(t, name)
			expected := byFile[name]

			var produced []*oracleRow
			appendRow := func(token diagnostic.TokenProxy, context enumerations.Context) {
				block := NewIdentificationOfTheSigningCertificate(i18nProvider, token, context, validationPolicy)
				result := block.Execute()
				row := toRow(name, token.Id(), context, "ISC",
					&result.XmlConstraintsConclusionContent, result.Title)
				if result.CertificateChain != nil {
					row.CertificateChain = make([]*oracleChainItem, 0, len(result.CertificateChain.ChainItem))
					for _, chainItem := range result.CertificateChain.ChainItem {
						id := chainItem.Id
						source := string(chainItem.Source.CertificateSourceType())
						row.CertificateChain = append(row.CertificateChain,
							&oracleChainItem{Id: &id, Source: &source})
					}
				}
				produced = append(produced, row)
			}

			for _, signature := range diagnosticData.Signatures() {
				context := enumerations.ContextSignature
				if signature.IsCounterSignature() {
					context = enumerations.ContextCounterSignature
				}
				appendRow(signature, context)
			}
			for _, timestamp := range diagnosticData.TimestampList() {
				appendRow(timestamp, enumerations.ContextTimestamp)
			}
			revocations := diagnosticData.AllRevocationData()
			sort.SliceStable(revocations, func(i, j int) bool { return revocations[i].Id() < revocations[j].Id() })
			for _, revocation := range revocations {
				appendRow(revocation, enumerations.ContextRevocation)
			}

			if len(produced) != len(expected) {
				t.Fatalf("produced %d ISC results, the oracle holds %d", len(produced), len(expected))
			}
			for i, want := range expected {
				got := produced[i]
				for _, constraint := range want.Constraints {
					if constraint.Status != nil {
						statuses[*constraint.Status]++
					}
				}
				if !reflect.DeepEqual(want, got) {
					t.Errorf("token %s: ISC mismatch\nexpected: %s\nactual:   %s",
						want.Token, mustJSON(t, want), mustJSON(t, got))
				}
			}
		})
	}

	// The corpus has to exercise both outcomes of the block, not just the happy path.
	for _, status := range []string{"OK", "NOT OK"} {
		if statuses[status] == 0 {
			t.Errorf("the corpus never produced a %q ISC constraint", status)
		}
	}
}
