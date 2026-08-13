// Ported from
// dss-cms/src/main/java/eu/europa/esig/dss/cms/operator/PrecomputedDigestCalculatorProvider.java
// (DSS 6.5.RC1).
package cms

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/model"
)

// PrecomputedDigestCalculatorProvider allows providing digest values without the original
// document, resolving the digest algorithm to ask digestDocument for from the
// AlgorithmIdentifier it is given. Port of the PrecomputedDigestCalculatorProvider class.
type PrecomputedDigestCalculatorProvider struct {
	// digestDocument is the DSSDocument to be signed.
	digestDocument model.DSSDocument
}

// NewPrecomputedDigestCalculatorProvider is the default constructor.
// Port of PrecomputedDigestCalculatorProvider(DSSDocument).
func NewPrecomputedDigestCalculatorProvider(dssDocument model.DSSDocument) *PrecomputedDigestCalculatorProvider {
	return &PrecomputedDigestCalculatorProvider{digestDocument: dssDocument}
}

// Digest resolves digestAlgorithmIdentifier to a DigestAlgorithm and returns digestDocument's
// value for it, an empty slice when the algorithm cannot be named or digestDocument does not
// carry a value for it. Port of #get(AlgorithmIdentifier) and its private getDigestBase64,
// whose catch-all around both failure modes (DigestAlgorithm.forOID's
// IllegalArgumentException and DigestDocument#getDigestValue's IllegalArgumentException) is
// reproduced by simply falling back to DSSUtils.EMPTY_BYTE_ARRAY (upstream logs a warning,
// dropped per PORTING.md).
func (p *PrecomputedDigestCalculatorProvider) Digest(digestAlgorithmIdentifier *asn1ber.AlgorithmIdentifier) ([]byte, error) {
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(digestAlgorithmIdentifier.Algorithm.String())
	if err != nil {
		return []byte{}, nil
	}
	value, err := p.digestDocument.DigestValue(digestAlgorithm)
	if err != nil {
		return []byte{}, nil
	}
	return value, nil
}

var _ DigestCalculatorProvider = (*PrecomputedDigestCalculatorProvider)(nil)
