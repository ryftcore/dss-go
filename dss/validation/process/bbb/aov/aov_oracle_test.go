package aov

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy"
)

// Shared plumbing of the three AOV oracle tests (aov_blocks_oracle_test.go,
// aov_cc_oracle_test.go and aov_direct_oracle_test.go): the JSONL row shape the
// Java driver in testdata/gen emits, the corpus loaders, and the conversion of a
// produced detailed-report object into that shape.
//
// See testdata/README.md for what the corpora are and how to regenerate them.

// aovCorpusRoot is the module-root-relative path of the marshal-parity
// diagnostic-data corpus the oracle ran over; it lives in the external
// corpus/ tree (see internal/corpustest).
const aovCorpusRoot = "diagnostic/jaxb/testdata/oracle"

// aovDumpDir holds XCVA's synthetic dumps, which this oracle replays too.
const aovDumpDir = "../xcv/testdata/dd"

// aovCurrentTime is the fixed validation time the oracle ran with,
// 2024-01-01T00:00:00Z.
var aovCurrentTime = time.Unix(1704067200, 0).UTC()

// aovDateFormat is the lexical form AovOracle prints dates in (millisecond
// precision, because the cc corpus probes expiration boundaries one
// millisecond apart).
const aovDateFormat = "2006-01-02T15:04:05.000Z"

type aovMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type aovConclusion struct {
	Indication    *string       `json:"indication"`
	SubIndication *string       `json:"subIndication"`
	Errors        []*aovMessage `json:"errors"`
	Warnings      []*aovMessage `json:"warnings"`
	Infos         []*aovMessage `json:"infos"`
}

type aovConstraint struct {
	Name           *aovMessage `json:"name"`
	Status         *string     `json:"status"`
	Error          *aovMessage `json:"error"`
	Warning        *aovMessage `json:"warning"`
	Info           *aovMessage `json:"info"`
	AdditionalInfo *string     `json:"additionalInfo"`
	Id             *string     `json:"id"`
	BlockType      *string     `json:"blockType"`
}

type aovAlgorithm struct {
	Name      *string `json:"name"`
	Uri       *string `json:"uri"`
	KeyLength *string `json:"keyLength"`
}

type aovCryptographicValidation struct {
	Algorithm                    *aovAlgorithm  `json:"algorithm"`
	NotAfter                     *string        `json:"notAfter"`
	ConcernedMaterialDescription *string        `json:"concernedMaterialDescription"`
	TokenId                      *string        `json:"tokenId"`
	Conclusion                   *aovConclusion `json:"conclusion"`
}

// aovRow is one line of any of the three corpora.
type aovRow struct {
	File           string  `json:"file"`
	Token          string  `json:"token"`
	Block          string  `json:"block"`
	Context        string  `json:"context"`
	ValidationTime *string `json:"validationTime"`

	Title       *string          `json:"title"`
	Conclusion  *aovConclusion   `json:"conclusion"`
	Constraints []*aovConstraint `json:"constraints"`

	SignatureCryptographicValidation        *aovCryptographicValidation   `json:"signatureCryptographicValidation"`
	SignedAttributesValidation              *aovCryptographicValidation   `json:"signedAttributesValidation"`
	DigestMatchersValidation                *aovCryptographicValidation   `json:"digestMatchersValidation"`
	CertificateChainCryptographicValidation []*aovCryptographicValidation `json:"certificateChainCryptographicValidation"`

	CryptographicValidation *aovCryptographicValidation `json:"cryptographicValidation"`
}

func loadAovRows(t *testing.T, path string) []*aovRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*aovRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &aovRow{}
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

// loadAovDiagnosticData reads the dump a row names: "dd/<name>" is one of XCVA's
// synthetic dumps, anything else a member of the marshal-parity corpus.
func loadAovDiagnosticData(t *testing.T, name string) *diagnostic.Data {
	t.Helper()
	var path string
	if rest, ok := strings.CutPrefix(name, "dd/"); ok {
		path = filepath.Join(aovDumpDir, rest)
	} else {
		path = corpustest.RootPath(t, filepath.Join(aovCorpusRoot, name))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	jaxbData, err := diagnosticjaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return diagnostic.NewData(jaxbData)
}

// aovI18nProvider is the single Provider all three replays use.
var aovI18nProvider = i18n.NewProvider()

func aovI18n() *i18n.Provider { return aovI18nProvider }

var aovDefaultPolicyValue modelpolicy.ValidationPolicy

func aovDefaultPolicy(t *testing.T) modelpolicy.ValidationPolicy {
	t.Helper()
	if aovDefaultPolicyValue == nil {
		aovDefaultPolicyValue = policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()
	}
	return aovDefaultPolicyValue
}

// ---------------------------------------------------------------- conversion

func toAovConclusion(conclusion *jaxb.XmlConclusion) *aovConclusion {
	if conclusion == nil {
		return nil
	}
	converted := &aovConclusion{
		Errors:   toAovMessages(conclusion.Errors),
		Warnings: toAovMessages(conclusion.Warnings),
		Infos:    toAovMessages(conclusion.Infos),
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

func toAovMessages(messages []*jaxb.XmlMessage) []*aovMessage {
	converted := make([]*aovMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toAovMessage(message))
	}
	return converted
}

func toAovMessage(message *jaxb.XmlMessage) *aovMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &aovMessage{Key: message.Key, Value: &value}
}

func toAovConstraints(constraints []*jaxb.XmlConstraint) []*aovConstraint {
	converted := make([]*aovConstraint, 0, len(constraints))
	for _, constraint := range constraints {
		status := string(constraint.Status)
		item := &aovConstraint{
			Name:           toAovMessage(constraint.Name),
			Status:         &status,
			Error:          toAovMessage(constraint.Error),
			Warning:        toAovMessage(constraint.Warning),
			Info:           toAovMessage(constraint.Info),
			AdditionalInfo: constraint.AdditionalInfo,
			Id:             constraint.Id,
		}
		if constraint.BlockType != nil {
			blockType := string(*constraint.BlockType)
			item.BlockType = &blockType
		}
		converted = append(converted, item)
	}
	return converted
}

func toAovCryptographicValidation(validation *jaxb.XmlCryptographicValidation) *aovCryptographicValidation {
	if validation == nil {
		return nil
	}
	converted := &aovCryptographicValidation{
		NotAfter:                     aovDate(validation.NotAfter),
		ConcernedMaterialDescription: validation.ConcernedMaterialDescription,
		TokenId:                      validation.TokenId,
		Conclusion:                   toAovConclusion(validation.Conclusion),
	}
	if validation.Algorithm != nil {
		name, uri := validation.Algorithm.Name, validation.Algorithm.Uri
		converted.Algorithm = &aovAlgorithm{Name: &name, Uri: &uri, KeyLength: validation.Algorithm.KeyLength}
	}
	return converted
}

func aovDate(value *jaxb.XSDateTime) *string {
	if value == nil {
		return nil
	}
	rendered := time.Time(*value).UTC().Format(aovDateFormat)
	return &rendered
}

// toAovRowBody fills the title/conclusion/constraints the Java "body" writer emits.
func toAovRowBody(row *aovRow, content *jaxb.XmlConstraintsConclusionContent, title string) {
	titleValue := title
	row.Title = &titleValue
	row.Conclusion = toAovConclusion(content.Conclusion)
	row.Constraints = toAovConstraints(content.Constraint)
}

// aovSafeExecute runs fn, converting a panic into a marker so that a token the
// Java oracle itself could not process does not abort the run.
func aovSafeExecute(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func mustAovJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}
