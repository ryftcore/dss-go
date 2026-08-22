// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteEvaluation.java (DSS 6.5.RC1).
package policy

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestCryptographicSuiteEvaluation_RoundTrip(t *testing.T) {
	e := NewCryptographicSuiteEvaluation()
	minSize := 2048
	param := NewCryptographicSuiteParameter()
	param.SetName("modulusLength")
	param.SetMin(&minSize)
	e.SetParameterList([]*CryptographicSuiteParameter{param})

	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	e.SetValidityStart(&start)
	e.SetValidityEnd(&end)
	e.SetAlgorithmUsage([]enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA,
	})
	e.SetRecommendation(enumerations.CryptographicSuiteRecommendation_RECOMMENDED)

	c := CryptographicSuiteEvaluationCopy(e)
	if !e.Equals(c) {
		t.Fatalf("copy not equal to original")
	}
	if len(c.ParameterList()) != 1 || c.ParameterList()[0].Name() != "modulusLength" {
		t.Fatalf("copy parameter list mismatch: %v", c.ParameterList())
	}

	// deep copy: mutating the copy's parameter list must not affect the
	// original.
	c.ParameterList()[0].SetName("mutated")
	if e.ParameterList()[0].Name() != "modulusLength" {
		t.Fatalf("original mutated via copy: %v", e.ParameterList()[0].Name())
	}

	if CryptographicSuiteEvaluationCopy(nil) != nil {
		t.Fatalf("copy of nil must be nil")
	}
}
