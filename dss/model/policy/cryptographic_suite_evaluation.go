// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteEvaluation.java (DSS 6.5.RC1).
package policy

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CryptographicSuiteEvaluation provides a representation of an
// "Evaluation" element extracted from an ETSI TS 119 322 cryptographic
// suite catalogue.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type CryptographicSuiteEvaluation struct {
	// parameterList is a list of the /dssc:Evaluation/dssc:Parameter
	// elements.
	parameterList []*CryptographicSuiteParameter

	// validityStart is the value of the
	// /dssc:Evaluation/dssc:Validity/dssc:Start element.
	validityStart *time.Time

	// validityEnd is the value of the
	// /dssc:Evaluation/dssc:Validity/dssc:End element.
	validityEnd *time.Time

	// algorithmUsage is a list of the
	// /dssc:Evaluation/etsi19322:MoreDetails/etsi19322:AlgorithmUsage
	// elements.
	algorithmUsage []enumerations.CryptographicSuiteAlgorithmUsage

	// recommendation is the value of the
	// /dssc:Evaluation/dssc:Validity/etsi19322:MoreDetails/etsi19322:Recommendation
	// element.
	recommendation enumerations.CryptographicSuiteRecommendation
}

// NewCryptographicSuiteEvaluation is the default constructor.
func NewCryptographicSuiteEvaluation() *CryptographicSuiteEvaluation {
	return &CryptographicSuiteEvaluation{}
}

// ParameterList gets the list of algorithm evaluation parameters.
func (c *CryptographicSuiteEvaluation) ParameterList() []*CryptographicSuiteParameter {
	return c.parameterList
}

// SetParameterList sets list of the /dssc:Evaluation/dssc:Parameter
// elements.
func (c *CryptographicSuiteEvaluation) SetParameterList(parameterList []*CryptographicSuiteParameter) {
	c.parameterList = parameterList
}

// ValidityStart gets the algorithm evaluation validity start date.
func (c *CryptographicSuiteEvaluation) ValidityStart() *time.Time {
	return c.validityStart
}

// SetValidityStart sets a value of the
// /dssc:Evaluation/dssc:Validity/dssc:Start element.
func (c *CryptographicSuiteEvaluation) SetValidityStart(validityStart *time.Time) {
	c.validityStart = validityStart
}

// ValidityEnd gets the algorithm evaluation validity end date.
func (c *CryptographicSuiteEvaluation) ValidityEnd() *time.Time {
	return c.validityEnd
}

// SetValidityEnd sets a value of the
// /dssc:Evaluation/dssc:Validity/dssc:End element.
func (c *CryptographicSuiteEvaluation) SetValidityEnd(validityEnd *time.Time) {
	c.validityEnd = validityEnd
}

// AlgorithmUsage gets the algorithm evaluation's usage scope.
func (c *CryptographicSuiteEvaluation) AlgorithmUsage() []enumerations.CryptographicSuiteAlgorithmUsage {
	return c.algorithmUsage
}

// SetAlgorithmUsage sets a list of the
// /dssc:Evaluation/etsi19322:MoreDetails/etsi19322:AlgorithmUsage
// elements.
func (c *CryptographicSuiteEvaluation) SetAlgorithmUsage(algorithmUsage []enumerations.CryptographicSuiteAlgorithmUsage) {
	c.algorithmUsage = algorithmUsage
}

// Recommendation gets the algorithm evaluation recommendation.
func (c *CryptographicSuiteEvaluation) Recommendation() enumerations.CryptographicSuiteRecommendation {
	return c.recommendation
}

// SetRecommendation sets a value of the
// /dssc:Evaluation/dssc:Validity/etsi19322:MoreDetails/etsi19322:Recommendation
// element.
func (c *CryptographicSuiteEvaluation) SetRecommendation(recommendation enumerations.CryptographicSuiteRecommendation) {
	c.recommendation = recommendation
}

// CryptographicSuiteEvaluationCopy instantiates a new
// CryptographicSuiteEvaluation by copying the values of evaluation. Ports
// CryptographicSuiteEvaluation#copy.
func CryptographicSuiteEvaluationCopy(evaluation *CryptographicSuiteEvaluation) *CryptographicSuiteEvaluation {
	if evaluation == nil {
		return nil
	}
	copyValue := NewCryptographicSuiteEvaluation()
	if evaluation.parameterList != nil {
		copyValue.parameterList = make([]*CryptographicSuiteParameter, len(evaluation.parameterList))
		for i, p := range evaluation.parameterList {
			copyValue.parameterList[i] = CryptographicSuiteParameterCopy(p)
		}
	}
	copyValue.validityStart = evaluation.validityStart
	copyValue.validityEnd = evaluation.validityEnd
	if evaluation.algorithmUsage != nil {
		copyValue.algorithmUsage = append([]enumerations.CryptographicSuiteAlgorithmUsage(nil), evaluation.algorithmUsage...)
	}
	copyValue.recommendation = evaluation.recommendation
	return copyValue
}

// Equals ports CryptographicSuiteEvaluation#equals.
func (c *CryptographicSuiteEvaluation) Equals(other *CryptographicSuiteEvaluation) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if len(c.parameterList) != len(other.parameterList) {
		return false
	}
	for i, p := range c.parameterList {
		if !p.Equals(other.parameterList[i]) {
			return false
		}
	}
	if !cryptographicSuiteEvaluationTimeEqual(c.validityStart, other.validityStart) {
		return false
	}
	if !cryptographicSuiteEvaluationTimeEqual(c.validityEnd, other.validityEnd) {
		return false
	}
	if len(c.algorithmUsage) != len(other.algorithmUsage) {
		return false
	}
	for i, u := range c.algorithmUsage {
		if u != other.algorithmUsage[i] {
			return false
		}
	}
	return c.recommendation == other.recommendation
}

// cryptographicSuiteEvaluationTimeEqual compares two possibly-nil
// *time.Time values.
func cryptographicSuiteEvaluationTimeEqual(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

// String ports CryptographicSuiteEvaluation#toString.
func (c *CryptographicSuiteEvaluation) String() string {
	return fmt.Sprintf("CryptographicSuiteEvaluation [parameterList=%v, validityStart=%v, validityEnd=%v, algorithmUsage=%v, recommendation=%v]",
		c.parameterList, c.validityStart, c.validityEnd, c.algorithmUsage, c.recommendation)
}
