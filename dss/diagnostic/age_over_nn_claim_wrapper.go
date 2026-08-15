// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/AgeOverNNClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// AgeOverNNClaimWrapper wraps an jaxb.XmlAgeOverNNClaim.
type AgeOverNNClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlAgeOverNNClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlAgeOverNNClaim
}

// NewAgeOverNNClaimWrapper is the default constructor. Port of AgeOverNNClaimWrapper(XmlAgeOverNNClaim).
func NewAgeOverNNClaimWrapper(wrapped *jaxb.XmlAgeOverNNClaim) *AgeOverNNClaimWrapper {
	return NewAgeOverNNClaimWrapperWithParent(wrapped, nil)
}

// NewAgeOverNNClaimWrapperWithParent is the constructor with a parent claim provided. Port of
// AgeOverNNClaimWrapper(XmlAgeOverNNClaim, ClaimWrapper).
func NewAgeOverNNClaimWrapperWithParent(wrapped *jaxb.XmlAgeOverNNClaim, parent *ClaimWrapper) *AgeOverNNClaimWrapper {
	return &AgeOverNNClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// Age gets the age value used for a verification within the claim. Port of getAge().
func (w *AgeOverNNClaimWrapper) Age() int {
	if w.wrapped.Age != nil {
		return *w.wrapped.Age
	}
	return 0
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *AgeOverNNClaimWrapper) Wrapped() *jaxb.XmlAgeOverNNClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type; isList/getList/isMap/getMap are not
// overridden here, so the embedded ClaimWrapper's own generic behaviour already applies and no
// override needs to be baked in. See the package note in claim_wrapper.go.
func (w *AgeOverNNClaimWrapper) AsClaim() *ClaimWrapper { return &w.ClaimWrapper }
