package aov

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/internal/corpustest"
)

// The cc KAT: every row of testdata/oracle/aov_cc.jsonl is the XmlCC upstream's
// SignatureAlgorithmCryptographicChecker or DigestAlgorithmCryptographicChecker
// produces for one (algorithm, key length, validation time) triple under the
// default ETSI policy's signature cryptographic suite - see
// testdata/gen/AovOracle.java. The matrix probes each algorithm at the validation
// time and at the exact instant its expiration date names, one millisecond before
// it and one millisecond after it, so that the reliability and
// at-validation-time comparisons are pinned on their boundary and not only well
// inside a branch.

func TestAovCryptographicCheckersAgainstJavaOracle(t *testing.T) {
	rows := loadAovRows(t, corpustest.Path(t, "oracle/aov_cc.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty oracle")
	}
	suite := aovDefaultPolicy(t).SignatureCryptographicConstraint(enumerations.Context_SIGNATURE)
	position := i18n.MessageTag_ACCM_POS_SIG_SIG

	// per check class, the statuses seen - every one must have an OK and a NOT OK
	statuses := map[string]map[string]int{}
	matched := 0

	for _, row := range rows {
		fields := strings.Split(row.Token, "|")
		var produced *jaxb.XmlCC
		switch row.Block {
		case "DigestAlgorithmCryptographicChecker":
			if len(fields) != 3 {
				t.Fatalf("malformed token %q", row.Token)
			}
			digestAlgorithm := aovDigestAlgorithm(fields[1])
			validationDate := aovParseDate(t, fields[2])
			produced = NewDigestAlgorithmCryptographicChecker(aovI18n(), digestAlgorithm,
				validationDate, position, suite).Execute()
		case "SignatureAlgorithmCryptographicChecker":
			if len(fields) != 4 {
				t.Fatalf("malformed token %q", row.Token)
			}
			signatureAlgorithm := aovSignatureAlgorithm(fields[1])
			keyLength := fields[2]
			if keyLength == "null" {
				// Java's null String; the Go port carries a plain string, whose
				// empty value stands in for it.
				keyLength = ""
			}
			validationDate := aovParseDate(t, fields[3])
			produced = NewSignatureAlgorithmCryptographicChecker(aovI18n(), signatureAlgorithm, keyLength,
				validationDate, position, suite).Execute()
		default:
			t.Fatalf("unknown block %q", row.Block)
		}

		got := &aovRow{Token: row.Token, Block: row.Block}
		toAovRowBody(got, &produced.XmlConstraintsConclusionContent, produced.Title)
		got.CryptographicValidation = toAovCryptographicValidation(produced.CryptographicValidation)
		matched++

		for _, constraint := range row.Constraints {
			if constraint.Name == nil || constraint.Name.Key == nil || constraint.Status == nil {
				continue
			}
			if _, seen := statuses[*constraint.Name.Key]; !seen {
				statuses[*constraint.Name.Key] = map[string]int{}
			}
			statuses[*constraint.Name.Key][*constraint.Status]++
		}

		if !reflect.DeepEqual(row, got) {
			t.Errorf("%s / %s mismatch\nexpected: %s\nactual:   %s",
				row.Block, row.Token, mustAovJSON(t, row), mustAovJSON(t, got))
		}
	}

	if matched != len(rows) {
		t.Errorf("replayed %d rows, oracle has %d", matched, len(rows))
	}
	// Every leaf check the two checkers wire must appear with both an OK and a
	// NOT OK row, so that a regression cannot hide behind a corpus that only
	// ever takes one branch. The six classes share five constraint tags
	// (SignatureAlgorithmAtValidationTimeCheck and DigestAlgorithmAtValidationTimeCheck
	// both raise ASCCM_AR).
	for _, tag := range []string{"ASCCM_AR", "ASCCM_APKSA", "ASCCM_CAA", "ASCCM_DAA", "ASCCM_PKSK"} {
		seen := statuses[tag]
		if seen["OK"] == 0 || seen["NOT OK"] == 0 {
			t.Errorf("cc check %s is one-sided in the corpus: OK=%d NOT OK=%d", tag, seen["OK"], seen["NOT OK"])
		}
	}
	if len(statuses) != 5 {
		t.Errorf("%d distinct cc checks in the corpus, expected 5", len(statuses))
	}
}

func aovDigestAlgorithm(name string) enumerations.DigestAlgorithm {
	if name == "null" {
		return ""
	}
	return enumerations.DigestAlgorithm(name)
}

func aovSignatureAlgorithm(name string) enumerations.SignatureAlgorithm {
	if name == "null" {
		return ""
	}
	return enumerations.SignatureAlgorithm(name)
}

func aovParseDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(aovDateFormat, value)
	if err != nil {
		t.Fatalf("parse date %q: %v", value, err)
	}
	return parsed.UTC()
}
