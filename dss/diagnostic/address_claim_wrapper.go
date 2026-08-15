// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/AddressClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// AddressClaimWrapper wraps a jaxb.XmlAddressClaim.
type AddressClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlAddressClaim, shadowing the promoted (synthetic) field of the
	// embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlAddressClaim
}

// NewAddressClaimWrapper is the default constructor. Port of AddressClaimWrapper(XmlAddressClaim).
func NewAddressClaimWrapper(wrapped *jaxb.XmlAddressClaim) *AddressClaimWrapper {
	return NewAddressClaimWrapperWithParent(wrapped, nil)
}

// NewAddressClaimWrapperWithParent is the constructor with a parent provided. Port of
// AddressClaimWrapper(XmlAddressClaim, ClaimWrapper).
func NewAddressClaimWrapperWithParent(wrapped *jaxb.XmlAddressClaim, parent *ClaimWrapper) *AddressClaimWrapper {
	w := &AddressClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// PostalAddress gets the user's full postal or mailing address, formatted, when present. Port
// of getPostalAddress(). Uses &w.ClaimWrapper rather than w.AsClaim() as the parent: AsClaim()
// itself computes Map(), which reaches this getter, so calling AsClaim() here would recurse
// (see the Parent()/recursion note in claim_wrapper.go's package comment).
func (w *AddressClaimWrapper) PostalAddress() *ClaimWrapper {
	if w.wrapped.PostalAddress != nil {
		return NewClaimWrapperWithParent(w.wrapped.PostalAddress, &w.ClaimWrapper)
	}
	return nil
}

// StreetAddress gets the user's street address, when present. The component may include a
// house number, street name, Post Office Box, and multi-line extended street address
// information. Port of getStreetAddress(). See the PostalAddress() note on the parent argument.
func (w *AddressClaimWrapper) StreetAddress() *ClaimWrapper {
	if w.wrapped.StreetAddress != nil {
		return NewClaimWrapperWithParent(w.wrapped.StreetAddress, &w.ClaimWrapper)
	}
	return nil
}

// City gets the user's city or locality address, when present. Port of getCity(). See the
// PostalAddress() note on the parent argument.
func (w *AddressClaimWrapper) City() *ClaimWrapper {
	if w.wrapped.City != nil {
		return NewClaimWrapperWithParent(w.wrapped.City, &w.ClaimWrapper)
	}
	return nil
}

// StateOrProvince gets the user's state or region address, when present. Port of
// getStateOrProvince(). See the PostalAddress() note on the parent argument.
func (w *AddressClaimWrapper) StateOrProvince() *ClaimWrapper {
	if w.wrapped.StateOrProvince != nil {
		return NewClaimWrapperWithParent(w.wrapped.StateOrProvince, &w.ClaimWrapper)
	}
	return nil
}

// PostalCode gets the user's zip code or postal code address, when present. Port of
// getPostalCode(). See the PostalAddress() note on the parent argument.
func (w *AddressClaimWrapper) PostalCode() *ClaimWrapper {
	if w.wrapped.PostalCode != nil {
		return NewClaimWrapperWithParent(w.wrapped.PostalCode, &w.ClaimWrapper)
	}
	return nil
}

// Country gets the user's country address, when present. Port of getCountry(). See the
// PostalAddress() note on the parent argument.
func (w *AddressClaimWrapper) Country() *ClaimWrapper {
	if w.wrapped.CountryName != nil {
		return NewClaimWrapperWithParent(w.wrapped.CountryName, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: an AddressClaimWrapper is unconditionally a map claim. Port of the
// overridden isMap().
func (w *AddressClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated address child claims rather than
// the generic Entry list. Port of the overridden getMap().
func (w *AddressClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if postalAddress := w.PostalAddress(); postalAddress != nil {
		result[postalAddress.Name()] = postalAddress
	}
	if streetAddress := w.StreetAddress(); streetAddress != nil {
		result[streetAddress.Name()] = streetAddress
	}
	if city := w.City(); city != nil {
		result[city.Name()] = city
	}
	if stateOrProvince := w.StateOrProvince(); stateOrProvince != nil {
		result[stateOrProvince.Name()] = stateOrProvince
	}
	if postalCode := w.PostalCode(); postalCode != nil {
		result[postalCode.Name()] = postalCode
	}
	if country := w.Country(); country != nil {
		result[country.Name()] = country
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *AddressClaimWrapper) Wrapped() *jaxb.XmlAddressClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override so that
// consumers that only hold the *ClaimWrapper view (e.g. after being placed in a []*ClaimWrapper)
// still observe it; Go has no virtual dispatch back from an embedded base to the embedding type.
// See the package note in claim_wrapper.go.
func (w *AddressClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
