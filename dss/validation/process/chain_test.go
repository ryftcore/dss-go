package process

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
)

// The chain-semantics KAT: every row of testdata/oracle/chain_semantics.jsonl is
// the XmlConstraintsConclusion upstream's Chain/ChainItem produce for the
// identically named scenario below (see testdata/gen/ChainSemanticsOracle.java).
// The scenario table here has to stay in sync with the Java one, entry for entry.

// spec is the specification of a single chain item, mirroring the Java
// oracle's Spec.
type spec struct {
	level               enumerations.Level // "" = constraint not defined
	valid               bool
	uninterrupted       bool
	bbbId               string // "" = null
	messageTag          i18n.MessageTag
	errorMessageTag     i18n.MessageTag
	additionalInfo      i18n.MessageTag
	blockType           jaxb.XmlBlockType
	failedIndication    enumerations.Indication
	failedSubIndication enumerations.SubIndication
	successIndication   enumerations.Indication
	successSubIndicaton enumerations.SubIndication
	previousErrors      []i18n.MessageTag
}

func newSpec(level enumerations.Level, valid bool) spec {
	return spec{
		level:               level,
		valid:               valid,
		messageTag:          i18n.MessageTagBBBICSISCI,
		errorMessageTag:     i18n.MessageTagBBBICSISCIANS,
		failedIndication:    enumerations.IndicationIndeterminate,
		failedSubIndication: enumerations.SubIndicationNoSigningCertificateFound,
	}
}

func (s spec) tags(messageTag, errorMessageTag i18n.MessageTag) spec {
	s.messageTag = messageTag
	s.errorMessageTag = errorMessageTag
	return s
}

func (s spec) withAdditionalInfo(tag i18n.MessageTag) spec {
	s.additionalInfo = tag
	return s
}

func (s spec) withBlockType(blockType jaxb.XmlBlockType) spec {
	s.blockType = blockType
	return s
}

func (s spec) withBbbId(bbbId string) spec {
	s.bbbId = bbbId
	return s
}

func (s spec) failure(indication enumerations.Indication, subIndication enumerations.SubIndication) spec {
	s.failedIndication = indication
	s.failedSubIndication = subIndication
	return s
}

func (s spec) success(indication enumerations.Indication, subIndication enumerations.SubIndication) spec {
	s.successIndication = indication
	s.successSubIndicaton = subIndication
	return s
}

func (s spec) asUninterrupted() spec {
	s.uninterrupted = true
	return s
}

func (s spec) withPreviousErrors(tags ...i18n.MessageTag) spec {
	s.previousErrors = tags
	return s
}

func (s spec) rule() policy.LevelRule {
	return GetLevelRule(s.level)
}

// specItem mirrors the Java oracle's SpecItem.
type specItem struct {
	*ChainItemBase[*jaxb.XmlConstraintsConclusion]
	spec spec
}

func newSpecItem(i18nProvider *i18n.I18nProvider, result *Result[*jaxb.XmlConstraintsConclusion],
	s spec) *specItem {
	var base *ChainItemBase[*jaxb.XmlConstraintsConclusion]
	if s.bbbId == "" {
		base = NewChainItemBase(i18nProvider, result, s.rule())
	} else {
		base = NewChainItemBaseWithId(i18nProvider, result, s.rule(), s.bbbId)
	}
	c := &specItem{ChainItemBase: base, spec: s}
	c.InitChainItem(c)
	return c
}

func (c *specItem) Process() bool                    { return c.spec.valid }
func (c *specItem) MessageTag() i18n.MessageTag      { return c.spec.messageTag }
func (c *specItem) ErrorMessageTag() i18n.MessageTag { return c.spec.errorMessageTag }
func (c *specItem) AdditionalInfo() i18n.MessageTag  { return c.spec.additionalInfo }
func (c *specItem) BlockType() jaxb.XmlBlockType     { return c.spec.blockType }
func (c *specItem) SuccessIndication() enumerations.Indication {
	return c.spec.successIndication
}
func (c *specItem) SuccessSubIndication() enumerations.SubIndication {
	return c.spec.successSubIndicaton
}
func (c *specItem) FailedIndicationForConclusion() enumerations.Indication {
	return c.spec.failedIndication
}
func (c *specItem) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.spec.failedSubIndication
}
func (c *specItem) PreviousErrors() []*jaxb.XmlMessage {
	messages := make([]*jaxb.XmlMessage, 0)
	for _, tag := range c.spec.previousErrors {
		messages = append(messages, c.BuildXmlMessage(tag))
	}
	return messages
}

// uninterruptedSpecItem mirrors the Java oracle's UninterruptedSpecItem.
type uninterruptedSpecItem struct {
	*UninterruptedChainItemBase[*jaxb.XmlConstraintsConclusion]
	spec spec
}

func newUninterruptedSpecItem(i18nProvider *i18n.I18nProvider,
	result *Result[*jaxb.XmlConstraintsConclusion], s spec) *uninterruptedSpecItem {
	var base *UninterruptedChainItemBase[*jaxb.XmlConstraintsConclusion]
	if s.bbbId == "" {
		base = NewUninterruptedChainItemBase(i18nProvider, result, s.rule())
	} else {
		base = NewUninterruptedChainItemBaseWithId(i18nProvider, result, s.rule(), s.bbbId)
	}
	c := &uninterruptedSpecItem{UninterruptedChainItemBase: base, spec: s}
	c.InitChainItem(c)
	return c
}

func (c *uninterruptedSpecItem) Process() bool                    { return c.spec.valid }
func (c *uninterruptedSpecItem) MessageTag() i18n.MessageTag      { return c.spec.messageTag }
func (c *uninterruptedSpecItem) ErrorMessageTag() i18n.MessageTag { return c.spec.errorMessageTag }
func (c *uninterruptedSpecItem) AdditionalInfo() i18n.MessageTag  { return c.spec.additionalInfo }
func (c *uninterruptedSpecItem) BlockType() jaxb.XmlBlockType     { return c.spec.blockType }
func (c *uninterruptedSpecItem) SuccessIndication() enumerations.Indication {
	return c.spec.successIndication
}
func (c *uninterruptedSpecItem) SuccessSubIndication() enumerations.SubIndication {
	return c.spec.successSubIndicaton
}
func (c *uninterruptedSpecItem) FailedIndicationForConclusion() enumerations.Indication {
	return c.spec.failedIndication
}
func (c *uninterruptedSpecItem) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.spec.failedSubIndication
}
func (c *uninterruptedSpecItem) PreviousErrors() []*jaxb.XmlMessage {
	messages := make([]*jaxb.XmlMessage, 0)
	for _, tag := range c.spec.previousErrors {
		messages = append(messages, c.BuildXmlMessage(tag))
	}
	return messages
}

// specChain mirrors the Java oracle's SpecChain.
type specChain struct {
	*ChainBase[*jaxb.XmlConstraintsConclusion]
	specs []spec
	title i18n.MessageTag
}

func newSpecChain(i18nProvider *i18n.I18nProvider, title i18n.MessageTag, specs []spec) *specChain {
	result := &jaxb.XmlConstraintsConclusion{}
	c := &specChain{
		ChainBase: NewChainBase(i18nProvider, NewResult(result,
			&result.XmlConstraintsConclusionContent, &result.XmlConstraintsConclusionAttrs)),
		specs: specs,
		title: title,
	}
	c.InitChainBase(c)
	return c
}

func (c *specChain) Title() i18n.MessageTag { return c.title }

func (c *specChain) InitChain() {
	var item ChainItem[*jaxb.XmlConstraintsConclusion]
	for _, s := range c.specs {
		var next ChainItem[*jaxb.XmlConstraintsConclusion]
		if s.uninterrupted {
			next = newUninterruptedSpecItem(c.I18nProvider, c.Result, s)
		} else {
			next = newSpecItem(c.I18nProvider, c.Result, s)
		}
		if item == nil {
			item = next
			c.FirstItem = item
		} else {
			item = item.SetNextItem(next)
		}
	}
}

// scenario mirrors the Java oracle's Scenario.
type scenario struct {
	name  string
	title i18n.MessageTag
	specs []spec
}

func chainScenarios() []scenario {
	return []scenario{
		{"fail-level-valid", i18n.MessageTagIdentificationOfTheSigningCertificate, []spec{
			newSpec(enumerations.LevelFail, true),
		}},
		{"fail-level-invalid", i18n.MessageTagIdentificationOfTheSigningCertificate, []spec{
			newSpec(enumerations.LevelFail, false),
		}},
		{"fail-level-short-circuit", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelFail, false),
			newSpec(enumerations.LevelFail, true).tags(i18n.MessageTagBBBCVISI, i18n.MessageTagBBBCVISIANS),
		}},
		{"fail-level-valid-then-invalid", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelFail, true),
			newSpec(enumerations.LevelFail, false).
				tags(i18n.MessageTagBBBCVISI, i18n.MessageTagBBBCVISIANS).
				failure(enumerations.IndicationFailed, enumerations.SubIndicationSigCryptoFailure),
		}},
		{"warn-level-invalid", i18n.MessageTagValidationContextInitialization, []spec{
			newSpec(enumerations.LevelWarn, false).
				tags(i18n.MessageTagBBBVCIIZHSP, i18n.MessageTagBBBVCIIZHSPANS),
		}},
		{"warn-level-valid", i18n.MessageTagValidationContextInitialization, []spec{
			newSpec(enumerations.LevelWarn, true).
				tags(i18n.MessageTagBBBVCIIZHSP, i18n.MessageTagBBBVCIIZHSPANS),
		}},
		{"inform-level-invalid", i18n.MessageTagValidationContextInitialization, []spec{
			newSpec(enumerations.LevelInform, false).
				tags(i18n.MessageTagBBBVCIISPSUPP, i18n.MessageTagBBBVCIISPSUPPANS),
		}},
		{"ignore-level-invalid", i18n.MessageTagValidationContextInitialization, []spec{
			newSpec(enumerations.LevelIgnore, false).withAdditionalInfo(i18n.MessageTagTokenID),
			newSpec(enumerations.LevelFail, true).tags(i18n.MessageTagBBBCVISI, i18n.MessageTagBBBCVISIANS),
		}},
		{"undefined-constraint", "", []spec{
			newSpec("", false),
			newSpec(enumerations.LevelFail, true).tags(i18n.MessageTagBBBCVISI, i18n.MessageTagBBBCVISIANS),
		}},
		{"warn-info-then-fail", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelWarn, false).
				tags(i18n.MessageTagBBBVCIIZHSP, i18n.MessageTagBBBVCIIZHSPANS),
			newSpec(enumerations.LevelInform, false).
				tags(i18n.MessageTagBBBVCIISPSUPP, i18n.MessageTagBBBVCIISPSUPPANS),
			newSpec(enumerations.LevelFail, false).
				tags(i18n.MessageTagBBBCVIRDOF, i18n.MessageTagBBBCVIRDOFANS).
				failure(enumerations.IndicationIndeterminate, enumerations.SubIndicationSignedDataNotFound),
		}},
		{"uninterrupted-continues-on-fail", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelFail, false).asUninterrupted().
				tags(i18n.MessageTagBBBCVIRDOF, i18n.MessageTagBBBCVIRDOFANS),
			newSpec(enumerations.LevelFail, false).asUninterrupted().
				tags(i18n.MessageTagBBBCVIRDOI, i18n.MessageTagBBBCVIRDOIANS).
				failure(enumerations.IndicationFailed, enumerations.SubIndicationHashFailure),
			newSpec(enumerations.LevelWarn, false).
				tags(i18n.MessageTagBBBVCIIZHSP, i18n.MessageTagBBBVCIIZHSPANS),
		}},
		{"custom-success-conclusion", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelFail, true).
				success(enumerations.IndicationPassed, enumerations.SubIndicationNoPOE),
			newSpec(enumerations.LevelFail, false).tags(i18n.MessageTagBBBCVISI, i18n.MessageTagBBBCVISIANS),
		}},
		{"previous-errors", i18n.MessageTagCryptographicVerification, []spec{
			newSpec(enumerations.LevelFail, false).
				tags(i18n.MessageTagBBBCVIRDOF, i18n.MessageTagBBBCVIRDOFANS).
				withPreviousErrors(i18n.MessageTagBBBCVISIANS, i18n.MessageTagBBBICSISCIANS),
		}},
		{"constraint-members-populated", i18n.MessageTagIdentificationOfTheSigningCertificate, []spec{
			newSpec(enumerations.LevelFail, false).withBbbId("S-1234").
				withBlockType(jaxb.XmlBlockTypeSigBBB).withAdditionalInfo(i18n.MessageTagEmpty),
		}},
		{"no-title-passed", "", []spec{
			newSpec(enumerations.LevelIgnore, true),
		}},
	}
}

// ---------------------------------------------------------------- oracle rows

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

type oracleRow struct {
	Scenario    string              `json:"scenario"`
	Title       *string             `json:"title"`
	Conclusion  *oracleConclusion   `json:"conclusion"`
	Constraints []*oracleConstraint `json:"constraints"`
}

func loadOracleRows(t *testing.T, path string) map[string]*oracleRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	rows := make(map[string]*oracleRow)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		row := &oracleRow{}
		if err := json.Unmarshal(line, row); err != nil {
			t.Fatalf("parse oracle row: %v", err)
		}
		rows[row.Scenario] = row
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	return rows
}

// toOracleRow renders a produced result in the oracle's shape.
func toOracleRow(name string, result *jaxb.XmlConstraintsConclusion) *oracleRow {
	row := &oracleRow{Scenario: name}
	title := result.Title
	row.Title = &title
	if result.Conclusion != nil {
		row.Conclusion = &oracleConclusion{
			Errors:   toOracleMessages(result.Conclusion.Errors),
			Warnings: toOracleMessages(result.Conclusion.Warnings),
			Infos:    toOracleMessages(result.Conclusion.Infos),
		}
		if indication := string(result.Conclusion.Indication.Indication()); indication != "" {
			row.Conclusion.Indication = &indication
		}
		if result.Conclusion.SubIndication != nil {
			subIndication := string(result.Conclusion.SubIndication.SubIndication())
			row.Conclusion.SubIndication = &subIndication
		}
	}
	row.Constraints = make([]*oracleConstraint, 0, len(result.Constraint))
	for _, constraint := range result.Constraint {
		status := string(constraint.Status)
		oc := &oracleConstraint{
			Name:           toOracleMessage(constraint.Name),
			Status:         &status,
			Error:          toOracleMessage(constraint.Error),
			Warning:        toOracleMessage(constraint.Warning),
			Info:           toOracleMessage(constraint.Info),
			AdditionalInfo: constraint.AdditionalInfo,
			Id:             constraint.Id,
		}
		if constraint.BlockType != nil {
			blockType := string(*constraint.BlockType)
			oc.BlockType = &blockType
		}
		row.Constraints = append(row.Constraints, oc)
	}
	return row
}

func toOracleMessages(messages []*jaxb.XmlMessage) []*oracleMessage {
	converted := make([]*oracleMessage, 0, len(messages))
	for _, message := range messages {
		converted = append(converted, toOracleMessage(message))
	}
	return converted
}

func toOracleMessage(message *jaxb.XmlMessage) *oracleMessage {
	if message == nil {
		return nil
	}
	value := message.Value
	return &oracleMessage{Key: message.Key, Value: &value}
}

func TestChainSemanticsAgainstJavaOracle(t *testing.T) {
	rows := loadOracleRows(t, corpustest.Path(t, "oracle/chain_semantics.jsonl"))
	i18nProvider := i18n.NewI18nProvider()

	for _, sc := range chainScenarios() {
		t.Run(sc.name, func(t *testing.T) {
			expected, ok := rows[sc.name]
			if !ok {
				t.Fatalf("no oracle row for scenario %q", sc.name)
			}
			if expected.Title == nil {
				// Java leaves the Title attribute null when the chain defines
				// no title MessageTag. The generated Go model carries Title as
				// a plain string, so null and "" are the same document; see
				// Result.SetTitle.
				empty := ""
				expected.Title = &empty
			}
			result := newSpecChain(i18nProvider, sc.title, sc.specs).Execute()
			actual := toOracleRow(sc.name, result)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("chain result mismatch\nexpected: %s\nactual:   %s",
					mustJSON(t, expected), mustJSON(t, actual))
			}
		})
	}

	if len(rows) != len(chainScenarios()) {
		t.Errorf("oracle holds %d rows, the scenario table %d", len(rows), len(chainScenarios()))
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}
