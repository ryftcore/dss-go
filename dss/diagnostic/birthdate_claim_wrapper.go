// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/BirthdateClaimWrapper.java (DSS 6.5.RC1).
//
// Java's constructor and getWrapped() declare the base type XmlClaim, and getBirthdate()/
// getApproximateMask()/isMap() runtime-check `wrapped instanceof XmlBirthdateClaim` before using
// the birthdate-specific fields. Every call site that constructs a BirthdateClaimWrapper (both
// here and in credential_subject_claim_wrapper.go/eaa_payload_proxy.go) passes a value whose
// generated-JAXB field is already concretely typed XmlBirthdateClaim (see jaxb_claim.go), so the
// instanceof check is always true on this side; the wrapper is ported directly against the
// concrete type, dropping the always-true runtime check but keeping its behaviour.
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// BirthdateClaimWrapper wraps a jaxb.XmlBirthdateClaim.
type BirthdateClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlBirthdateClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlBirthdateClaim
}

// NewBirthdateClaimWrapper is the default constructor. Port of BirthdateClaimWrapper(XmlClaim).
func NewBirthdateClaimWrapper(wrapped *jaxb.XmlBirthdateClaim) *BirthdateClaimWrapper {
	return NewBirthdateClaimWrapperWithParent(wrapped, nil)
}

// NewBirthdateClaimWrapperWithParent is the constructor with a parent claim provided. Port of
// BirthdateClaimWrapper(XmlClaim, ClaimWrapper).
func NewBirthdateClaimWrapperWithParent(wrapped *jaxb.XmlBirthdateClaim, parent *ClaimWrapper) *BirthdateClaimWrapper {
	return &BirthdateClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// Birthdate gets the user's birthdate. Port of getBirthdate() (the wrapped instanceof
// XmlBirthdateClaim branch is always taken here, see the file note above).
func (w *BirthdateClaimWrapper) Birthdate() *ClaimWrapper {
	if w.wrapped.Birthdate != nil {
		return NewClaimWrapperWithParent(w.wrapped.Birthdate, &w.ClaimWrapper)
	}
	if w.IsDateTime() {
		return w.AsClaim()
	}
	return nil
}

// ApproximateMask gets an 8 digit flag to denote the location of the mask in YYYYMMDD format.
// 1 denotes mask. Issuing authority should pick one exact date to be used for full-date value.
// Port of getApproximateMask() (the wrapped instanceof XmlBirthdateClaim branch is always taken
// here, see the file note above).
func (w *BirthdateClaimWrapper) ApproximateMask() *ClaimWrapper {
	if w.wrapped.ApproximateMask != nil {
		return NewClaimWrapperWithParent(w.wrapped.ApproximateMask, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override, true when the claim is not itself a bare DateTime leaf value. Port of
// isMap() (the wrapped instanceof XmlBirthdateClaim conjunct is always true here, see the file
// note above).
func (w *BirthdateClaimWrapper) IsMap() bool {
	return !w.IsDateTime()
}

// Map is the override, assembling the map from the dedicated birthdate child claims when IsMap
// applies, falling back to the base's generic Entry-derived behaviour otherwise. Port of
// getMap().
func (w *BirthdateClaimWrapper) Map() map[string]*ClaimWrapper {
	if !w.IsMap() {
		return w.ClaimWrapper.Map()
	}
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if birthdate := w.Birthdate(); birthdate != nil {
		result[birthdate.Name()] = birthdate
	}
	if approximateMask := w.ApproximateMask(); approximateMask != nil {
		result[approximateMask.Name()] = approximateMask
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *BirthdateClaimWrapper) Wrapped() *jaxb.XmlBirthdateClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override only
// when IsMap applies (otherwise the base's generic Entry-derived behaviour is what Java's
// super.getMap() fallback would also reach). See the package note in claim_wrapper.go.
func (w *BirthdateClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	if w.IsMap() {
		cw.mapOverride = w.Map()
		isMap := true
		cw.isMapOverride = &isMap
	}
	return &cw
}
