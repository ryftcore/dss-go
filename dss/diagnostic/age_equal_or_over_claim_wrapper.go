// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/AgeEqualOrOverClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// AgeEqualOrOverClaimWrapper wraps a jaxb.XmlAgeEqualOrOverClaim claim.
type AgeEqualOrOverClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlAgeEqualOrOverClaim, shadowing the promoted (synthetic) field
	// of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlAgeEqualOrOverClaim
}

// NewAgeEqualOrOverClaimWrapper is the default constructor. Port of
// AgeEqualOrOverClaimWrapper(XmlAgeEqualOrOverClaim).
func NewAgeEqualOrOverClaimWrapper(wrapped *jaxb.XmlAgeEqualOrOverClaim) *AgeEqualOrOverClaimWrapper {
	return NewAgeEqualOrOverClaimWrapperWithParent(wrapped, nil)
}

// NewAgeEqualOrOverClaimWrapperWithParent is the constructor with a parent claim provided. Port
// of AgeEqualOrOverClaimWrapper(XmlAgeEqualOrOverClaim, ClaimWrapper).
func NewAgeEqualOrOverClaimWrapperWithParent(wrapped *jaxb.XmlAgeEqualOrOverClaim, parent *ClaimWrapper) *AgeEqualOrOverClaimWrapper {
	return &AgeEqualOrOverClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// AgeEqualOrOverList gets a list of age specific claims embedded within the map. Port of
// getAgeEqualOrOverList().
func (w *AgeEqualOrOverClaimWrapper) AgeEqualOrOverList() []*AgeOverNNClaimWrapper {
	ageOverNN := w.wrapped.AgeOverNNClaim
	if len(ageOverNN) == 0 {
		return nil
	}
	result := make([]*AgeOverNNClaimWrapper, 0, len(ageOverNN))
	for _, item := range ageOverNN {
		result = append(result, NewAgeOverNNClaimWrapper(item))
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *AgeEqualOrOverClaimWrapper) Wrapped() *jaxb.XmlAgeEqualOrOverClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type; isList/getList/isMap/getMap are not
// overridden here. See the package note in claim_wrapper.go.
func (w *AgeEqualOrOverClaimWrapper) AsClaim() *ClaimWrapper { return &w.ClaimWrapper }
