// Ported from
// dss-cms/src/main/java/eu/europa/esig/dss/cms/operator/CustomMessageDigestCalculatorProvider.java
// (DSS 6.5.RC1).
package cms

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// CustomMessageDigestCalculatorProvider represents a DigestCalculatorProvider for a
// message-digest calculation: it always answers the digest value it was constructed with,
// regardless of the algorithm identifier asked for (see enumerations.DigestAlgorithm below,
// the assumption being that the caller only ever asks with the algorithm messageDigestAlgo
// names). Port of the CustomMessageDigestCalculatorProvider class.
type CustomMessageDigestCalculatorProvider struct {
	// messageDigestValue is the message digest value.
	messageDigestValue []byte
}

// NewCustomMessageDigestCalculatorProvider is the default constructor to create an object with
// a message digest provided in a form of byte array. Port of
// CustomMessageDigestCalculatorProvider(DigestAlgorithm, byte[]); messageDigestAlgo is kept as
// a parameter for API parity with Java, though - unlike Java, which logs it - this port does
// not otherwise use it, having no logging (PORTING.md).
//
// Panics with the Java messages when messageDigestAlgo is empty or messageDigestValue is nil
// (Objects.requireNonNull).
func NewCustomMessageDigestCalculatorProvider(messageDigestAlgo enumerations.DigestAlgorithm, messageDigestValue []byte) *CustomMessageDigestCalculatorProvider {
	if messageDigestAlgo == "" {
		panic("DigestAlgorithm shall be defined!")
	}
	if messageDigestValue == nil {
		panic("Digest value shall be defined!")
	}
	return &CustomMessageDigestCalculatorProvider{messageDigestValue: messageDigestValue}
}

// Digest returns the pre-computed messageDigestValue, whatever digestAlgorithmIdentifier is
// asked for. Port of #get(AlgorithmIdentifier).
func (p *CustomMessageDigestCalculatorProvider) Digest(digestAlgorithmIdentifier *asn1ber.AlgorithmIdentifier) ([]byte, error) {
	return p.messageDigestValue, nil
}

var _ DigestCalculatorProvider = (*CustomMessageDigestCalculatorProvider)(nil)
