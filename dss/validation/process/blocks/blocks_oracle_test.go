package blocks

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy"
)

// The BasicBuildingBlocks KAT: every row of testdata/oracle/blocks.jsonl is what
// upstream's eu.europa.esig.dss.validation.process.bbb.BasicBuildingBlocks
// produces for one token of one diagnostic-data dump, at validation time
// 2024-01-01T00:00:00Z under the default ETSI policy - see
// testdata/gen/BlocksOracle.java. The row carries what the dispatcher itself
// decides: which of the seven sub-blocks it instantiated, each one's conclusion,
// the aggregated conclusion updateFinalConclusion() built, the certificate chain
// it copied off the ISC block, and the cross/equivalent certificate lists
// addAdditionalInfo() hung off each XmlSubXCV.

const (
	// blocksCorpusRoot is the module-root-relative path (inside the
	// external corpus/ tree, see internal/corpustest) of the
	// marshal-parity diagnostic-data corpus.
	blocksCorpusRoot = "diagnostic/jaxb/testdata/oracle"
	blocksDumpDir    = "../bbb/xcv/testdata/dd"
	blocksOwnDumpDir = "testdata/dd"
)

var blocksCurrentTime = time.Unix(1704067200, 0).UTC()

type blocksMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type blocksConclusion struct {
	Indication    *string          `json:"indication"`
	SubIndication *string          `json:"subIndication"`
	Errors        []*blocksMessage `json:"errors"`
	Warnings      []*blocksMessage `json:"warnings"`
	Infos         []*blocksMessage `json:"infos"`
}

type blocksBlock struct {
	Title      *string           `json:"title"`
	Conclusion *blocksConclusion `json:"conclusion"`
}

type blocksSubXCV struct {
	Id                     *string           `json:"id"`
	CrossCertificates      []string          `json:"crossCertificates"`
	EquivalentCertificates []string          `json:"equivalentCertificates"`
	Conclusion             *blocksConclusion `json:"conclusion"`
}

type blocksRow struct {
	File             string            `json:"file"`
	Token            string            `json:"token"`
	Context          string            `json:"context"`
	Id               *string           `json:"id"`
	Type             *string           `json:"type"`
	Conclusion       *blocksConclusion `json:"conclusion"`
	FC               *blocksBlock      `json:"fc"`
	ISC              *blocksBlock      `json:"isc"`
	VCI              *blocksBlock      `json:"vci"`
	AOV              *blocksBlock      `json:"aov"`
	XCV              *blocksBlock      `json:"xcv"`
	CV               *blocksBlock      `json:"cv"`
	SAV              *blocksBlock      `json:"sav"`
	CertificateChain []string          `json:"certificateChain"`
	SubXCV           []*blocksSubXCV   `json:"subXCV"`
}

func TestBasicBuildingBlocksAgainstJavaOracle(t *testing.T) {
	rows := loadBlocksRows(t, corpustest.Path(t, "oracle/blocks.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	validationPolicy := blocksDefaultPolicy(t)

	byFile := map[string]map[string]*blocksRow{}
	var files []string
	for _, row := range rows {
		if _, seen := byFile[row.File]; !seen {
			files = append(files, row.File)
			byFile[row.File] = map[string]*blocksRow{}
		}
		byFile[row.File][row.Token] = row
	}

	contexts := map[string]int{}
	indications := map[string]int{}
	matched := 0

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			diagnosticData := loadBlocksDiagnosticData(t, name)
			// The dispatcher shares one bbbs map across every token of a dump,
			// exactly as the oracle does.
			bbbs := map[string]*jaxb.XmlBasicBuildingBlocks{}

			compare := func(token string, proxy diagnostic.TokenProxy, context enumerations.Context) {
				want, ok := byFile[name][token]
				if !ok {
					return // the Java oracle threw on this token and recorded no row
				}
				var produced *jaxb.XmlBasicBuildingBlocks
				if blocksSafeExecute(func() {
					produced = NewBasicBuildingBlocks(blocksI18n(), diagnosticData, proxy,
						blocksCurrentTime, bbbs, validationPolicy, context).Execute()
				}) {
					t.Errorf("%s: Go panicked where Java produced a row", token)
					return
				}
				got := toBlocksRow(name, token, want.Context, produced)
				matched++
				contexts[want.Context]++
				if want.Conclusion != nil && want.Conclusion.Indication != nil {
					indications[*want.Conclusion.Indication]++
				}
				if !reflect.DeepEqual(want, got) {
					t.Errorf("%s mismatch\nexpected: %s\nactual:   %s",
						token, mustBlocksJSON(t, want), mustBlocksJSON(t, got))
				}
			}

			for _, signature := range diagnosticData.Signatures() {
				context := enumerations.ContextSignature
				if signature.IsCounterSignature() {
					context = enumerations.ContextCounterSignature
				}
				compare("SIG|"+signature.Id(), signature, context)
			}
			for _, timestamp := range diagnosticData.TimestampList() {
				compare("TST|"+timestamp.Id(), timestamp, enumerations.ContextTimestamp)
			}
			revocations := append([]*diagnostic.RevocationWrapper(nil), diagnosticData.AllRevocationData()...)
			sort.SliceStable(revocations, func(i, j int) bool { return revocations[i].Id() < revocations[j].Id() })
			for _, revocation := range revocations {
				compare("REV|"+revocation.Id(), revocation, enumerations.ContextRevocation)
			}
			for _, certificate := range diagnosticData.UsedCertificates() {
				compare("CERT|"+certificate.Id(), certificate, enumerations.ContextCertificate)
			}
		})
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}
	for _, context := range []string{"SIGNATURE", "TIMESTAMP", "REVOCATION", "CERTIFICATE"} {
		if contexts[context] == 0 {
			t.Errorf("no %s row replayed", context)
		}
	}
	for _, indication := range []string{"PASSED", "INDETERMINATE"} {
		if indications[indication] == 0 {
			t.Errorf("no %s conclusion in the replayed corpus", indication)
		}
	}
}

// ------------------------------------------------------------------ plumbing

func loadBlocksRows(t *testing.T, path string) []*blocksRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*blocksRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &blocksRow{}
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

func loadBlocksDiagnosticData(t *testing.T, name string) *diagnostic.DiagnosticData {
	t.Helper()
	var path string
	if rest, ok := strings.CutPrefix(name, "dd/"); ok {
		path = filepath.Join(blocksDumpDir, rest)
	} else if rest, ok := strings.CutPrefix(name, "own/"); ok {
		path = filepath.Join(blocksOwnDumpDir, rest)
	} else {
		path = corpustest.RootPath(t, filepath.Join(blocksCorpusRoot, name))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	jaxbData, err := diagnosticjaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return diagnostic.NewDiagnosticData(jaxbData)
}

var blocksI18nProvider = i18n.NewI18nProvider()

func blocksI18n() *i18n.I18nProvider { return blocksI18nProvider }

var blocksDefaultPolicyValue modelpolicy.ValidationPolicy

func blocksDefaultPolicy(t *testing.T) modelpolicy.ValidationPolicy {
	t.Helper()
	if blocksDefaultPolicyValue == nil {
		blocksDefaultPolicyValue = policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()
	}
	return blocksDefaultPolicyValue
}

func blocksSafeExecute(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func mustBlocksJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

// ---------------------------------------------------------------- conversion

func toBlocksConclusion(conclusion *jaxb.XmlConclusion) *blocksConclusion {
	if conclusion == nil {
		return nil
	}
	converted := &blocksConclusion{
		Errors:   toBlocksMessages(conclusion.Errors),
		Warnings: toBlocksMessages(conclusion.Warnings),
		Infos:    toBlocksMessages(conclusion.Infos),
	}
	if indication := string(conclusion.Indication.Indication()); indication != "" {
		converted.Indication = &indication
	}
	if conclusion.SubIndication != nil {
		subIndication := string(conclusion.SubIndication.SubIndication())
		converted.SubIndication = &subIndication
	}
	return converted
}

func toBlocksMessages(messages []*jaxb.XmlMessage) []*blocksMessage {
	converted := make([]*blocksMessage, 0, len(messages))
	for _, message := range messages {
		value := message.Value
		converted = append(converted, &blocksMessage{Key: message.Key, Value: &value})
	}
	return converted
}

func toBlocksBlock(content *jaxb.XmlConstraintsConclusionContent, title string, present bool) *blocksBlock {
	if !present {
		return nil
	}
	titleValue := title
	return &blocksBlock{Title: &titleValue, Conclusion: toBlocksConclusion(content.Conclusion)}
}

func toBlocksRow(file, token, context string, result *jaxb.XmlBasicBuildingBlocks) *blocksRow {
	id := result.Id
	row := &blocksRow{File: file, Token: token, Context: context, Id: &id}
	if typeValue := string(result.Type.Context()); typeValue != "" {
		row.Type = &typeValue
	}
	row.Conclusion = toBlocksConclusion(result.Conclusion)
	if result.FC != nil {
		row.FC = toBlocksBlock(&result.FC.XmlConstraintsConclusionContent, result.FC.Title, true)
	}
	if result.ISC != nil {
		row.ISC = toBlocksBlock(&result.ISC.XmlConstraintsConclusionContent, result.ISC.Title, true)
	}
	if result.VCI != nil {
		row.VCI = toBlocksBlock(&result.VCI.XmlConstraintsConclusionContent, result.VCI.Title, true)
	}
	if result.AOV != nil {
		row.AOV = toBlocksBlock(&result.AOV.XmlConstraintsConclusionContent, result.AOV.Title, true)
	}
	if result.XCV != nil {
		row.XCV = toBlocksBlock(&result.XCV.XmlConstraintsConclusionContent, result.XCV.Title, true)
	}
	if result.CV != nil {
		row.CV = toBlocksBlock(&result.CV.XmlConstraintsConclusionContent, result.CV.Title, true)
	}
	if result.SAV != nil {
		row.SAV = toBlocksBlock(&result.SAV.XmlConstraintsConclusionContent, result.SAV.Title, true)
	}
	if result.CertificateChain != nil {
		row.CertificateChain = make([]string, 0, len(result.CertificateChain.ChainItem))
		for _, item := range result.CertificateChain.ChainItem {
			row.CertificateChain = append(row.CertificateChain, item.Id)
		}
	}
	if result.XCV != nil {
		row.SubXCV = make([]*blocksSubXCV, 0, len(result.XCV.SubXCV))
		for _, sub := range result.XCV.SubXCV {
			id := sub.Id
			item := &blocksSubXCV{Id: &id, Conclusion: toBlocksConclusion(sub.Conclusion)}
			item.CrossCertificates = []string{}
			if sub.CrossCertificate != nil {
				item.CrossCertificates = append(item.CrossCertificates, *sub.CrossCertificate...)
			}
			item.EquivalentCertificates = []string{}
			if sub.EquivalentCertificate != nil {
				item.EquivalentCertificates = append(item.EquivalentCertificates, *sub.EquivalentCertificate...)
			}
			row.SubXCV = append(row.SubXCV, item)
		}
	}
	return row
}
