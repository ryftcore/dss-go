// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/EncapsulatedTimestampScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EncapsulatedTimestampScopeFinder is used to find a signature scope for an embedded
// timestamp from a collection of SignatureScope candidates, extracted from a signature.
type EncapsulatedTimestampScopeFinder struct {
	AbstractSignatureScopeFinder

	// Signature is the AdvancedSignature embedding the timestamp.
	Signature validation.AdvancedSignature
}

// NewEncapsulatedTimestampScopeFinder is the default constructor instantiating object with
// null signature. Port of EncapsulatedTimestampScopeFinder().
func NewEncapsulatedTimestampScopeFinder() *EncapsulatedTimestampScopeFinder {
	return &EncapsulatedTimestampScopeFinder{AbstractSignatureScopeFinder: NewAbstractSignatureScopeFinder()}
}

// SetSignature sets an encapsulating AdvancedSignature. Port of
// setSignature(AdvancedSignature).
func (f *EncapsulatedTimestampScopeFinder) SetSignature(signature validation.AdvancedSignature) {
	f.Signature = signature
}

// FindTimestampScope returns a timestamp scope for the given TimestampToken. Port of
// findTimestampScope(TimestampToken).
func (f *EncapsulatedTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() {
		return f.FilterCoveredSignatureScopes(timestampToken)
	}
	return []mscope.SignatureScope{}
}

// FilterCoveredSignatureScopes filters and returns covered SignatureScopes by the current
// timestamp. Port of filterCoveredSignatureScopes(TimestampToken).
func (f *EncapsulatedTimestampScopeFinder) FilterCoveredSignatureScopes(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	// return all by default
	return f.Signature.SignatureScopes()
}

// compile-time interface assertion.
var _ TimestampScopeFinder = (*EncapsulatedTimestampScopeFinder)(nil)
