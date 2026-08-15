// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/ValidityInfoClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// ValidityInfoClaimWrapper wraps a jaxb.XmlValidityInfoClaim. Unlike most claim subtype
// wrappers, Java declares only the single-argument constructor here (no parent-taking overload).
type ValidityInfoClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlValidityInfoClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlValidityInfoClaim
}

// NewValidityInfoClaimWrapper is the default constructor. Port of
// ValidityInfoClaimWrapper(XmlValidityInfoClaim).
func NewValidityInfoClaimWrapper(wrapped *jaxb.XmlValidityInfoClaim) *ValidityInfoClaimWrapper {
	w := &ValidityInfoClaimWrapper{
		ClaimWrapper: *NewClaimWrapper(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs)),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// Signed gets the timestamp at which the MSO signature was created. Port of getSigned().
func (w *ValidityInfoClaimWrapper) Signed() *ClaimWrapper {
	if w.wrapped.Signed != nil {
		return NewClaimWrapperWithParent(w.wrapped.Signed, &w.ClaimWrapper)
	}
	return nil
}

// ValidFrom gets the timestamp before which the MSO is not yet valid. Port of getValidFrom().
func (w *ValidityInfoClaimWrapper) ValidFrom() *ClaimWrapper {
	if w.wrapped.ValidFrom != nil {
		return NewClaimWrapperWithParent(w.wrapped.ValidFrom, &w.ClaimWrapper)
	}
	return nil
}

// ValidUntil gets the timestamp after which the MSO is no longer valid. Port of getValidUntil().
func (w *ValidityInfoClaimWrapper) ValidUntil() *ClaimWrapper {
	if w.wrapped.ValidUntil != nil {
		return NewClaimWrapperWithParent(w.wrapped.ValidUntil, &w.ClaimWrapper)
	}
	return nil
}

// ExpectedUpdate gets the timestamp at which the issuing authority infrastructure expects to
// re-sign the MSO. Port of getExpectedUpdate().
func (w *ValidityInfoClaimWrapper) ExpectedUpdate() *ClaimWrapper {
	if w.wrapped.ExpectedUpdate != nil {
		return NewClaimWrapperWithParent(w.wrapped.ExpectedUpdate, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: a ValidityInfoClaimWrapper is unconditionally a map claim. Port of the
// overridden isMap().
func (w *ValidityInfoClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated validity-info child claims rather
// than the generic Entry list. Port of the overridden getMap().
func (w *ValidityInfoClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if signed := w.Signed(); signed != nil {
		result[signed.Name()] = signed
	}
	if validFrom := w.ValidFrom(); validFrom != nil {
		result[validFrom.Name()] = validFrom
	}
	if validUntil := w.ValidUntil(); validUntil != nil {
		result[validUntil.Name()] = validUntil
	}
	if expectedUpdate := w.ExpectedUpdate(); expectedUpdate != nil {
		result[expectedUpdate.Name()] = expectedUpdate
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *ValidityInfoClaimWrapper) Wrapped() *jaxb.XmlValidityInfoClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *ValidityInfoClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
