// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/asn1/SignedAssertions.java
// (DSS 6.5.RC1).
//
//	SignedAssertions ::= SEQUENCE OF SignedAssertion
package cms

import "github.com/utain/esig/dss/internal/asn1ber"

// SignedAssertions is a list of signed assertions. Port of the SignedAssertions class.
type SignedAssertions struct {
	// assertions is the list of signed assertions.
	assertions []*SignedAssertion
}

// NewSignedAssertions creates the SignedAssertions from a list of SignedAssertions.
// Port of SignedAssertions(List<SignedAssertion>).
func NewSignedAssertions(assertions []*SignedAssertion) *SignedAssertions {
	return &SignedAssertions{assertions: assertions}
}

// ParseSignedAssertions decodes a SignedAssertions from its encoding.
// Port of SignedAssertions.getInstance(Object).
func ParseSignedAssertions(encoded []byte) (*SignedAssertions, error) {
	children, err := cmsSequenceOfSize(encoded, "SignedAssertions", 0, 1<<31-1)
	if err != nil {
		return nil, err
	}
	assertions := make([]*SignedAssertion, 0, len(children))
	for _, child := range children {
		assertion, err := ParseSignedAssertion(child.Encoded())
		if err != nil {
			return nil, err
		}
		assertions = append(assertions, assertion)
	}
	return &SignedAssertions{assertions: assertions}, nil
}

// Assertions returns the list of SignedAssertions. Port of #getAssertions.
func (s *SignedAssertions) Assertions() []*SignedAssertion { return s.assertions }

// String ports #toString.
func (s *SignedAssertions) String() string {
	result := ""
	for _, assertion := range s.assertions {
		result += assertion.String() + "\n"
	}
	return result
}

// DER returns the DER encoding of the SignedAssertions. Port of #toASN1Primitive.
func (s *SignedAssertions) DER() []byte {
	var body []byte
	for _, assertion := range s.assertions {
		body = append(body, assertion.DER()...)
	}
	return asn1ber.WriteSequence(body)
}
