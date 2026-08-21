package xcv

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
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy"
)

// Shared plumbing of the two XCVA oracle tests (xcva_direct_oracle_test.go and
// xcva_blocks_oracle_test.go): the JSONL row shape the Java drivers in
// testdata/gen emit, the corpus loaders, and the conversion of a produced
// detailed-report object into that shape.
//
// See testdata/README.md for what the corpora are and how to regenerate them.

// xcvaCorpusRoot is the module-root-relative path (inside the external
// corpus/ tree, see internal/corpustest) of the marshal-parity
// diagnostic-data corpus the oracles were run over; xcvaDumpDir holds the
// synthetic dumps XcvaSyntheticDumps wrote next to them, checked in locally.
// A row's "file" is the plain dump name for the first and "dd/<name>" for
// the second.
const (
	xcvaCorpusRoot = "diagnostic/jaxb/testdata/oracle"
	xcvaDumpDir    = "testdata/dd"
)

// xcvaCurrentTime is the fixed validation time the oracles ran with,
// 2024-01-01T00:00:00Z.
var xcvaCurrentTime = time.Unix(1704067200, 0).UTC()

// xcvaDateFormat is the lexical form XcvaOracle prints the RAC dates in.
const xcvaDateFormat = "2006-01-02T15:04:05Z"

type xcvaMessage struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type xcvaConclusion struct {
	Indication    *string        `json:"indication"`
	SubIndication *string        `json:"subIndication"`
	Errors        []*xcvaMessage `json:"errors"`
	Warnings      []*xcvaMessage `json:"warnings"`
	Infos         []*xcvaMessage `json:"infos"`
}

type xcvaConstraint struct {
	Name           *xcvaMessage `json:"name"`
	Status         *string      `json:"status"`
	Error          *xcvaMessage `json:"error"`
	Warning        *xcvaMessage `json:"warning"`
	Info           *xcvaMessage `json:"info"`
	AdditionalInfo *string      `json:"additionalInfo"`
	Id             *string      `json:"id"`
	BlockType      *string      `json:"blockType"`
}

// xcvaNode is one XmlConstraintsConclusion of the produced tree: the XCV itself,
// or one of its SubXCV children, or a CRS with its RAC children, or a RAC with
// its CRS child. Only the members the emitting block actually carries are filled.
type xcvaNode struct {
	Id                           *string           `json:"id"`
	LatestAcceptableRevocationId *string           `json:"latestAcceptableRevocationId"`
	AcceptableRevocationId       []string          `json:"acceptableRevocationId"`
	RevocationThisUpdate         *string           `json:"revocationThisUpdate"`
	RevocationProductionDate     *string           `json:"revocationProductionDate"`
	TrustAnchor                  *bool             `json:"trustAnchor"`
	SelfSigned                   *bool             `json:"selfSigned"`
	Title                        *string           `json:"title"`
	Conclusion                   *xcvaConclusion   `json:"conclusion"`
	Constraints                  []*xcvaConstraint `json:"constraints"`
	RAC                          []*xcvaNode       `json:"rac"`
	CRS                          *xcvaNode         `json:"crs"`
	SubXCV                       []*xcvaNode       `json:"subXCV"`
}

// xcvaRow is one line of either corpus.
type xcvaRow struct {
	File    string `json:"file"`
	Policy  string `json:"policy"`
	Token   string `json:"token"`
	Check   string `json:"check"`
	Context string `json:"context"`
	Block   string `json:"block"`
	xcvaNode
}

func loadXcvaRows(t *testing.T, path string) []*xcvaRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*xcvaRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &xcvaRow{}
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

// loadXcvaDiagnosticData reads the dump a row names: "dd/<name>" is one of the
// synthetic dumps in testdata/dd, anything else a member of the marshal-parity
// corpus.
func loadXcvaDiagnosticData(t *testing.T, name string) *diagnostic.DiagnosticData {
	t.Helper()
	var path string
	if rest, ok := strings.CutPrefix(name, "dd/"); ok {
		path = filepath.Join(xcvaDumpDir, rest)
	} else {
		path = corpustest.RootPath(t, filepath.Join(xcvaCorpusRoot, name))
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

// xcvaPolicies returns the three validation policies the block oracle used, keyed
// by the name its rows carry: the default ETSI one, the CHAIN-model variant and the
// OCSP-certHash variant, both committed in testdata/policy.
func xcvaPolicies(t *testing.T) map[string]modelpolicy.ValidationPolicy {
	t.Helper()
	variant, err := os.ReadFile("testdata/policy/constraint-chain.xml")
	if err != nil {
		t.Fatalf("read chain policy: %v", err)
	}
	certHash, err := os.ReadFile("testdata/policy/constraint-certhash.xml")
	if err != nil {
		t.Fatalf("read certhash policy: %v", err)
	}
	factory := policy.NewEtsiValidationPolicyFactory()
	return map[string]modelpolicy.ValidationPolicy{
		"default":  xcvaDefaultPolicy(t),
		"chain":    factory.LoadValidationPolicy(model.NewInMemoryDocument(variant)),
		"certhash": factory.LoadValidationPolicy(model.NewInMemoryDocument(certHash)),
	}
}

// xcvaI18nProvider is the single I18nProvider both replays use, built once.
var xcvaI18nProvider = i18n.NewI18nProvider()

func xcvaI18n() *i18n.I18nProvider { return xcvaI18nProvider }

// xcvaDefaultPolicyValue is the default ETSI policy, loaded once.
var xcvaDefaultPolicyValue modelpolicy.ValidationPolicy

func xcvaDefaultPolicy(t *testing.T) modelpolicy.ValidationPolicy {
	t.Helper()
	if xcvaDefaultPolicyValue == nil {
		xcvaDefaultPolicyValue = policy.NewEtsiValidationPolicyFactory().LoadDefaultValidationPolicy()
	}
	return xcvaDefaultPolicyValue
}

// passedAOV is the PASSED Algorithm Obsolescence Validation result the oracle fed
// the blocks, the way the phase 8c sav corpus does.
func passedAOV() *jaxb.XmlAOV {
	aov := &jaxb.XmlAOV{}
	aov.Conclusion = &jaxb.XmlConclusion{
		Indication: jaxb.IndicationValue(enumerations.Indication_PASSED),
	}
	return aov
}

// ---------------------------------------------------------------- conversion

func toXcvaConclusion(conclusion *jaxb.XmlConclusion) *xcvaConclusion {
	if conclusion == nil {
		return nil
	}
	converted := &xcvaConclusion{
		Errors:   toXcvaMessages(conclusion.Errors),
		Warnings: toXcvaMessages(conclusion.Warnings),
		Infos:    toXcvaMessages(conclusion.Infos),
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

func toXcvaMessages(messages []*jaxb.XmlMessage) []*xcvaMessage {
	converted := make([]*xcvaMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toXcvaMessage(message))
	}
	return converted
}

func toXcvaMessage(message *jaxb.XmlMessage) *xcvaMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &xcvaMessage{Key: message.Key, Value: &value}
}

func toXcvaConstraints(constraints []*jaxb.XmlConstraint) []*xcvaConstraint {
	converted := make([]*xcvaConstraint, 0, len(constraints))
	for _, constraint := range constraints {
		status := string(constraint.Status)
		item := &xcvaConstraint{
			Name:           toXcvaMessage(constraint.Name),
			Status:         &status,
			Error:          toXcvaMessage(constraint.Error),
			Warning:        toXcvaMessage(constraint.Warning),
			Info:           toXcvaMessage(constraint.Info),
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

// toXcvaBody fills the title/conclusion/constraints the Java "body" writer emits.
func toXcvaBody(node *xcvaNode, content *jaxb.XmlConstraintsConclusionContent, title string) {
	titleValue := title
	node.Title = &titleValue
	node.Conclusion = toXcvaConclusion(content.Conclusion)
	node.Constraints = toXcvaConstraints(content.Constraint)
}

func toXcvaXCV(result *jaxb.XmlXCV) *xcvaNode {
	node := &xcvaNode{}
	toXcvaBody(node, &result.XmlConstraintsConclusionContent, result.Title)
	node.SubXCV = make([]*xcvaNode, 0, len(result.SubXCV))
	for _, sub := range result.SubXCV {
		node.SubXCV = append(node.SubXCV, toXcvaSubXCV(sub))
	}
	return node
}

func toXcvaSubXCV(sub *jaxb.XmlSubXCV) *xcvaNode {
	id := sub.Id
	node := &xcvaNode{Id: &id, TrustAnchor: sub.TrustAnchor, SelfSigned: sub.SelfSigned}
	toXcvaBody(node, &sub.XmlConstraintsConclusionContent, sub.Title)
	return node
}

func toXcvaCRS(result *jaxb.XmlCRS) *xcvaNode {
	node := &xcvaNode{Id: result.Id, LatestAcceptableRevocationId: result.LatestAcceptableRevocationId}
	node.AcceptableRevocationId = []string{}
	if result.AcceptableRevocationId != nil {
		node.AcceptableRevocationId = append(node.AcceptableRevocationId, *result.AcceptableRevocationId...)
	}
	toXcvaBody(node, &result.XmlConstraintsConclusionContent, result.Title)
	node.RAC = make([]*xcvaNode, 0, len(result.RAC))
	for _, rac := range result.RAC {
		node.RAC = append(node.RAC, toXcvaRAC(rac))
	}
	return node
}

func toXcvaRAC(result *jaxb.XmlRAC) *xcvaNode {
	node := &xcvaNode{
		Id:                       result.Id,
		RevocationThisUpdate:     xcvaDate(result.RevocationThisUpdate),
		RevocationProductionDate: xcvaDate(result.RevocationProductionDate),
	}
	toXcvaBody(node, &result.XmlConstraintsConclusionContent, result.Title)
	if result.CRS != nil {
		node.CRS = toXcvaCRS(result.CRS)
	}
	return node
}

// xcvaDate renders an XmlRAC date the way XcvaOracle does. Known mapping: the
// generated Go model carries RevocationThisUpdate and RevocationProductionDate as
// plain XSDateTime members where the JAXB class has nullable Dates (see
// revocation_acceptance_checker.go), so the zero time stands in for Java's null
// and is rendered as one here.
func xcvaDate(value jaxb.XSDateTime) *string {
	instant := time.Time(value)
	if instant.IsZero() {
		return nil
	}
	rendered := instant.UTC().Format(xcvaDateFormat)
	return &rendered
}

// xcvaSafeExecute runs fn, converting a panic into a marker so that a token the
// Java oracle itself could not process (it recorded no row) does not abort the run.
func xcvaSafeExecute(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func mustXcvaJSON(t *testing.T, value interface{}) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}
