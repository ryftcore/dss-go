// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteAlgorithm.java (DSS 6.5.RC1).
package policy

import "fmt"

// CryptographicSuiteAlgorithm provides a representation of an "Algorithm"
// element extracted from an ETSI TS 119 322 cryptographic suite
// catalogue.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type CryptographicSuiteAlgorithm struct {
	// algorithmIdentifierName is the value of the
	// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:Name element.
	algorithmIdentifierName string

	// algorithmIdentifierOIDs is a list of values from the
	// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:ObjectIdentifier
	// elements.
	algorithmIdentifierOIDs []string

	// algorithmIdentifierURIs is a list of values from the
	// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:URI elements.
	algorithmIdentifierURIs []string

	// evaluationList is a list of the /dssc:Algorithm/dssc:Evaluation
	// elements.
	evaluationList []*CryptographicSuiteEvaluation

	// informationTextList is a list of values from the
	// /dssc:Algorithm/dssc:Information/dssc:Text elements.
	informationTextList []string
}

// NewCryptographicSuiteAlgorithm is the default constructor.
func NewCryptographicSuiteAlgorithm() *CryptographicSuiteAlgorithm {
	return &CryptographicSuiteAlgorithm{}
}

// AlgorithmIdentifierName gets the algorithm identifier name.
func (c *CryptographicSuiteAlgorithm) AlgorithmIdentifierName() string {
	return c.algorithmIdentifierName
}

// SetAlgorithmIdentifierName sets the value of the
// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:Name element.
func (c *CryptographicSuiteAlgorithm) SetAlgorithmIdentifierName(algorithmIdentifierName string) {
	c.algorithmIdentifierName = algorithmIdentifierName
}

// AlgorithmIdentifierOIDs gets the algorithm identifier OIDs list.
func (c *CryptographicSuiteAlgorithm) AlgorithmIdentifierOIDs() []string {
	return c.algorithmIdentifierOIDs
}

// SetAlgorithmIdentifierOIDs sets a list of values from the
// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:ObjectIdentifier
// elements.
func (c *CryptographicSuiteAlgorithm) SetAlgorithmIdentifierOIDs(algorithmIdentifierOIDs []string) {
	c.algorithmIdentifierOIDs = algorithmIdentifierOIDs
}

// AlgorithmIdentifierURIs gets the algorithm identifier URIs list.
func (c *CryptographicSuiteAlgorithm) AlgorithmIdentifierURIs() []string {
	return c.algorithmIdentifierURIs
}

// SetAlgorithmIdentifierURIs sets a list of values from the
// /dssc:Algorithm/dssc:AlgorithmIdentifier/dssc:URI elements.
func (c *CryptographicSuiteAlgorithm) SetAlgorithmIdentifierURIs(algorithmIdentifierURIs []string) {
	c.algorithmIdentifierURIs = algorithmIdentifierURIs
}

// EvaluationList gets a collection of algorithm evaluation requirements.
func (c *CryptographicSuiteAlgorithm) EvaluationList() []*CryptographicSuiteEvaluation {
	return c.evaluationList
}

// SetEvaluationList sets a list of the/dssc:Algorithm/dssc:Evaluation
// elements.
func (c *CryptographicSuiteAlgorithm) SetEvaluationList(evaluationList []*CryptographicSuiteEvaluation) {
	c.evaluationList = evaluationList
}

// InformationTextList gets a list of information text strings.
func (c *CryptographicSuiteAlgorithm) InformationTextList() []string {
	return c.informationTextList
}

// SetInformationTextList sets a list of values from the
// /dssc:Algorithm/dssc:Information/dssc:Text elements.
func (c *CryptographicSuiteAlgorithm) SetInformationTextList(informationTextList []string) {
	c.informationTextList = informationTextList
}

// CryptographicSuiteAlgorithmCopy instantiates a new
// CryptographicSuiteAlgorithm by copying the values of algorithm. Ports
// CryptographicSuiteAlgorithm#copy.
func CryptographicSuiteAlgorithmCopy(algorithm *CryptographicSuiteAlgorithm) *CryptographicSuiteAlgorithm {
	if algorithm == nil {
		return nil
	}
	copyValue := NewCryptographicSuiteAlgorithm()
	copyValue.algorithmIdentifierName = algorithm.algorithmIdentifierName
	if algorithm.algorithmIdentifierOIDs != nil {
		copyValue.algorithmIdentifierOIDs = append([]string(nil), algorithm.algorithmIdentifierOIDs...)
	}
	if algorithm.algorithmIdentifierURIs != nil {
		copyValue.algorithmIdentifierURIs = append([]string(nil), algorithm.algorithmIdentifierURIs...)
	}
	if algorithm.evaluationList != nil {
		copyValue.evaluationList = make([]*CryptographicSuiteEvaluation, len(algorithm.evaluationList))
		for i, e := range algorithm.evaluationList {
			copyValue.evaluationList[i] = CryptographicSuiteEvaluationCopy(e)
		}
	}
	if algorithm.informationTextList != nil {
		copyValue.informationTextList = append([]string(nil), algorithm.informationTextList...)
	}
	return copyValue
}

// Equals ports CryptographicSuiteAlgorithm#equals.
func (c *CryptographicSuiteAlgorithm) Equals(other *CryptographicSuiteAlgorithm) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if c.algorithmIdentifierName != other.algorithmIdentifierName {
		return false
	}
	if !stringSliceEqual(c.algorithmIdentifierOIDs, other.algorithmIdentifierOIDs) {
		return false
	}
	if !stringSliceEqual(c.algorithmIdentifierURIs, other.algorithmIdentifierURIs) {
		return false
	}
	if len(c.evaluationList) != len(other.evaluationList) {
		return false
	}
	for i, e := range c.evaluationList {
		if !e.Equals(other.evaluationList[i]) {
			return false
		}
	}
	return stringSliceEqual(c.informationTextList, other.informationTextList)
}

// stringSliceEqual compares two possibly-nil string slices for content
// equality, mirroring java.util.Objects#equals over java.util.List.
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// String ports CryptographicSuiteAlgorithm#toString.
func (c *CryptographicSuiteAlgorithm) String() string {
	return fmt.Sprintf("CryptographicSuiteAlgorithm [algorithmIdentifierName='%s', algorithmIdentifierOIDs=%v, "+
		"algorithmIdentifierURIs=%v, evaluationList=%v, informationTextList=%v]",
		c.algorithmIdentifierName, c.algorithmIdentifierOIDs, c.algorithmIdentifierURIs, c.evaluationList, c.informationTextList)
}
