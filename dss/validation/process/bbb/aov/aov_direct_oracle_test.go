package aov

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// The aov/checks KAT: every row of testdata/oracle/aov_direct.jsonl is the XmlSAV
// upstream produces when one of the four aov/checks result checks is driven alone
// through a one-item chain at Level.FAIL over a hand-built XmlAOV / XmlCC shape -
// see testdata/gen/AovOracle.java. These four are the checks the AOV blocks
// themselves never wire (they are consumed by fc / sav / xcv), so the block
// corpus cannot reach them; the shapes cover every branch their process() and
// buildAdditionalInfo() take.

// singleAovSAVChain is a chain of exactly one item, mirroring the oracle's
// SingleSAVChain.
type singleAovSAVChain struct {
	*process.ChainBase[*jaxb.XmlSAV]
	factory aovCheckFactory
}

type aovCheckFactory func(result *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV]

func newSingleAovSAVChain(i18nProvider *i18n.Provider, factory aovCheckFactory) *singleAovSAVChain {
	xmlSAV := &jaxb.XmlSAV{}
	c := &singleAovSAVChain{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlSAV,
			&xmlSAV.XmlConstraintsConclusionContent, &xmlSAV.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleAovSAVChain) InitChain() {
	c.FirstItem = c.factory(c.Result)
}

func TestAovDirectChecksAgainstJavaOracle(t *testing.T) {
	rows := loadAovRows(t, corpustest.Path(t, "oracle/aov_direct.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	fail := process.GetLevelRule(enumerations.LevelFail)
	position := i18n.MessageTagACCMPosSigSig

	classes := map[string]map[string]int{}
	matched := 0

	for _, row := range rows {
		fields := strings.Split(row.Token, "|")
		class, shape := fields[0], fields[len(fields)-1]
		certificateContext := len(fields) == 3 && fields[1] == "cert"

		var factory aovCheckFactory
		switch class {
		case "AlgorithmObsolescenceValidationCheck":
			aovResult := aovOfShape(shape)
			factory = func(r *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
				return NewAlgorithmObsolescenceValidationCheck(aovI18n(), r, aovResult, aovCurrentTime, position, "T-1")
			}
		case "AlgorithmObsolescenceValidationCheckWithId":
			aovResult := aovOfShape(shape)
			factory = func(r *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
				return NewAlgorithmObsolescenceValidationCheckWithId(aovI18n(), r, aovResult, aovCurrentTime, position, "T-1")
			}
		case "SignatureAlgorithmCryptographicCheckerResultCheck":
			ccResult := ccOfShape(shape)
			if certificateContext {
				factory = func(r *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
					return NewSignatureAlgorithmCryptographicCheckerResultCheckWithContext(aovI18n(), r, aovCurrentTime,
						enumerations.ContextCertificate, position, ccResult, fail, "C-1")
				}
			} else {
				factory = func(r *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
					return NewSignatureAlgorithmCryptographicCheckerResultCheck(aovI18n(), r, aovCurrentTime,
						position, ccResult, fail)
				}
			}
		case "DigestAlgorithmCryptographicCheckerResultCheck":
			ccResult := ccOfShape(shape)
			factory = func(r *process.Result[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
				return NewDigestAlgorithmCryptographicCheckerResultCheck(aovI18n(), r, aovCurrentTime,
					position, ccResult, fail)
			}
		default:
			t.Fatalf("unknown check class %q", class)
		}

		var produced *jaxb.XmlSAV
		if aovSafeExecute(func() { produced = newSingleAovSAVChain(aovI18n(), factory).Execute() }) {
			t.Errorf("%s: Go panicked where Java produced a row", row.Token)
			continue
		}
		got := &aovRow{Token: row.Token, Block: row.Block}
		toAovRowBody(got, &produced.XmlConstraintsConclusionContent, produced.Title)
		// Known mapping: a one-item chain defines no title MessageTag, so
		// Java leaves the Title attribute null where the generated Go model
		// carries a plain string. This one field is normalised and nothing
		// else.
		if row.Title == nil && got.Title != nil && *got.Title == "" {
			got.Title = nil
		}
		matched++

		if _, seen := classes[class]; !seen {
			classes[class] = map[string]int{}
		}
		for _, constraint := range row.Constraints {
			if constraint.Status != nil {
				classes[class][*constraint.Status]++
			}
		}

		if !reflect.DeepEqual(row, got) {
			t.Errorf("%s mismatch\nexpected: %s\nactual:   %s",
				row.Token, mustAovJSON(t, row), mustAovJSON(t, got))
		}
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}
	for _, class := range []string{
		"AlgorithmObsolescenceValidationCheck",
		"AlgorithmObsolescenceValidationCheckWithId",
		"SignatureAlgorithmCryptographicCheckerResultCheck",
		"DigestAlgorithmCryptographicCheckerResultCheck",
	} {
		seen := classes[class]
		if seen["OK"] == 0 || seen["NOT OK"] == 0 {
			t.Errorf("check %s is one-sided in the corpus: OK=%d NOT OK=%d", class, seen["OK"], seen["NOT OK"])
		}
	}
}

// aovOfShape mirrors AovOracle#aovOfShape.
func aovOfShape(shape string) *jaxb.XmlAOV {
	result := &jaxb.XmlAOV{}
	conclusion := &jaxb.XmlConclusion{}
	switch shape {
	case "passed":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
	case "passed-with-algo":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		keyLength := "2048"
		validation := &jaxb.XmlCryptographicValidation{
			Algorithm: &jaxb.XmlCryptographicAlgorithm{
				Name:      "RSA",
				Uri:       "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
				KeyLength: &keyLength,
			},
			Conclusion: &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(enumerations.IndicationPassed)},
		}
		result.SignatureCryptographicValidation = validation
	case "error":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationCryptoConstraintsFailureNoPOE)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors,
			aovXmlMessage("BBB_SAV_ASCCM_ANS", "the algorithm is no longer reliable"))
	case "warning":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		conclusion.Warnings = append(conclusion.Warnings,
			aovXmlMessage("BBB_SAV_ASCCM_ANS", "the algorithm is about to expire"))
	case "info":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		conclusion.Infos = append(conclusion.Infos,
			aovXmlMessage("BBB_SAV_ASCCM_ANS", "informative note"))
	default:
		panic("unknown shape " + shape)
	}
	result.Conclusion = conclusion
	return result
}

// ccOfShape mirrors AovOracle#ccOfShape.
func ccOfShape(shape string) *jaxb.XmlCC {
	result := &jaxb.XmlCC{}
	conclusion := &jaxb.XmlConclusion{}
	algorithm := &jaxb.XmlCryptographicAlgorithm{
		Name: "SHA256",
		Uri:  "http://www.w3.org/2001/04/xmlenc#sha256",
	}
	switch shape {
	case "passed":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		keyLength := "2048"
		algorithm.KeyLength = &keyLength
	case "passed-no-keylength":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
	case "failed":
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationCryptoConstraintsFailureNoPOE)
		conclusion.SubIndication = &subIndication
		conclusion.Errors = append(conclusion.Errors, aovXmlMessage("ASCCM_AR_ANS_ANR", "SHA1 is not reliable"))
	default:
		panic("unknown shape " + shape)
	}
	validationConclusion := &jaxb.XmlConclusion{
		Indication:    conclusion.Indication,
		SubIndication: conclusion.SubIndication,
	}
	validationConclusion.Errors = append(validationConclusion.Errors, conclusion.Errors...)
	result.CryptographicValidation = &jaxb.XmlCryptographicValidation{
		Algorithm:  algorithm,
		Conclusion: validationConclusion,
	}
	result.Conclusion = conclusion
	return result
}

func aovXmlMessage(key, value string) *jaxb.XmlMessage {
	id := key
	return &jaxb.XmlMessage{Key: &id, Value: value}
}

var _ policy.LevelRule = process.GetLevelRule(enumerations.LevelFail)
