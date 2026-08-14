// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/PlaceOfBirthClaimWrapper.java (DSS 6.5.RC1).
//
// Java's constructor and getWrapped() declare the base type XmlClaim, and getCity()/getRegion()/
// getCountry()/isMap() runtime-check `wrapped instanceof XmlPlaceOfBirthClaim` before using the
// place-of-birth-specific fields. Unlike its siblings (BirthdateClaimWrapper,
// AttestedAttributesSubjectClaimIdWrapper), that check is NOT always true on this side: the two
// call sites are backed by genuinely different generated-JAXB field types (see jaxb_claim.go/
// jaxb_eaa.go) - XmlEAAPayload.PlaceOfBirth is concretely XmlPlaceOfBirthClaim, but
// XmlCredentialSubjectClaim.PlaceOfBirth is the generic XmlClaim, matching Java's own
// `XmlClaim placeOfBirth = getWrapped().getPlaceOfBirth();` local in
// CredentialSubjectClaimWrapper.getPlaceOfBirth(). Go has no upcast between the two unrelated
// jaxb struct types, so both call sites are supported directly: the exported constructors take
// the concrete XmlPlaceOfBirthClaim (matching every external caller, e.g. eaa_payload_proxy.go),
// while newPlaceOfBirthClaimWrapperFromClaim (unexported, used only by
// credential_subject_claim_wrapper.go) takes the generic XmlClaim and leaves `wrapped` nil - the
// Go equivalent of the instanceof check being false.
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// PlaceOfBirthClaimWrapper wraps a jaxb.XmlPlaceOfBirthClaim, or a generic jaxb.XmlClaim from a
// schema position not concretely typed XmlPlaceOfBirthClaim; see the file note above.
type PlaceOfBirthClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlPlaceOfBirthClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go. nil when
	// the wrapper was built from a generic XmlClaim (Java's `instanceof` false case; see the
	// file note above).
	wrapped *jaxb.XmlPlaceOfBirthClaim
}

// NewPlaceOfBirthClaimWrapper is the default constructor. Port of PlaceOfBirthClaimWrapper(XmlClaim).
func NewPlaceOfBirthClaimWrapper(wrapped *jaxb.XmlPlaceOfBirthClaim) *PlaceOfBirthClaimWrapper {
	return NewPlaceOfBirthClaimWrapperWithParent(wrapped, nil)
}

// NewPlaceOfBirthClaimWrapperWithParent is the constructor with a parent claim provided. Port of
// PlaceOfBirthClaimWrapper(XmlClaim, ClaimWrapper).
func NewPlaceOfBirthClaimWrapperWithParent(wrapped *jaxb.XmlPlaceOfBirthClaim, parent *ClaimWrapper) *PlaceOfBirthClaimWrapper {
	return &PlaceOfBirthClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// newPlaceOfBirthClaimWrapperFromClaim builds a PlaceOfBirthClaimWrapper over a generic XmlClaim,
// for schema positions where the JAXB field is not concretely typed XmlPlaceOfBirthClaim; see the
// file note above.
func newPlaceOfBirthClaimWrapperFromClaim(wrapped *jaxb.XmlClaim, parent *ClaimWrapper) *PlaceOfBirthClaimWrapper {
	return &PlaceOfBirthClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(wrapped, parent),
	}
}

// City gets the user's city or locality address, when present. Port of getCity() (the wrapped
// instanceof XmlPlaceOfBirthClaim branch; see the file note above).
func (w *PlaceOfBirthClaimWrapper) City() *ClaimWrapper {
	if w.wrapped != nil && w.wrapped.City != nil {
		return NewClaimWrapperWithParent(w.wrapped.City, &w.ClaimWrapper)
	}
	return nil
}

// Region gets the user's zip code or postal code address, when present. Port of getRegion()
// (the wrapped instanceof XmlPlaceOfBirthClaim branch; see the file note above).
func (w *PlaceOfBirthClaimWrapper) Region() *ClaimWrapper {
	if w.wrapped != nil && w.wrapped.Region != nil {
		return NewClaimWrapperWithParent(w.wrapped.Region, &w.ClaimWrapper)
	}
	return nil
}

// Country gets the user's country address, when present. Port of getCountry() (the wrapped
// instanceof XmlPlaceOfBirthClaim branch; see the file note above).
func (w *PlaceOfBirthClaimWrapper) Country() *ClaimWrapper {
	if w.wrapped != nil && w.wrapped.Country != nil {
		return NewClaimWrapperWithParent(w.wrapped.Country, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override, true when the claim was built from a concrete XmlPlaceOfBirthClaim and
// is not itself a bare Text leaf value. Port of isMap() (the wrapped instanceof
// XmlPlaceOfBirthClaim conjunct; see the file note above).
func (w *PlaceOfBirthClaimWrapper) IsMap() bool {
	return w.wrapped != nil && !w.IsText()
}

// Map is the override, assembling the map from the dedicated place-of-birth child claims when
// IsMap applies, falling back to the base's generic Entry-derived behaviour otherwise. Port of
// getMap().
func (w *PlaceOfBirthClaimWrapper) Map() map[string]*ClaimWrapper {
	if !w.IsMap() {
		return w.ClaimWrapper.Map()
	}
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if city := w.City(); city != nil {
		result[city.Name()] = city
	}
	if region := w.Region(); region != nil {
		result[region.Name()] = region
	}
	if country := w.Country(); country != nil {
		result[country.Name()] = country
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant
// getWrapped(): Java's unconditional `(XmlPlaceOfBirthClaim) super.getWrapped()` cast throws
// ClassCastException when the wrapped value is not actually an XmlPlaceOfBirthClaim; this panics
// the same way when the wrapper was built from a generic XmlClaim (see the file note above).
func (w *PlaceOfBirthClaimWrapper) Wrapped() *jaxb.XmlPlaceOfBirthClaim {
	if w.wrapped == nil {
		panic("XmlClaim cannot be cast to XmlPlaceOfBirthClaim")
	}
	return w.wrapped
}

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override only
// when IsMap applies. See the package note in claim_wrapper.go.
func (w *PlaceOfBirthClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	if w.IsMap() {
		cw.mapOverride = w.Map()
		isMap := true
		cw.isMapOverride = &isMap
	}
	return &cw
}
