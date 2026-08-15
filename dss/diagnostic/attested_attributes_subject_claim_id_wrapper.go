// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/AttestedAttributesSubjectClaimIdWrapper.java (DSS 6.5.RC1).
//
// Java's constructor and getWrapped() declare the base type XmlClaim, and getFamilyName()/
// getGivenName()/getDocumentNumber()/isMap()/getMap() runtime-check
// `wrapped instanceof XmlAttestedAttributesSubjectIdClaim` before using the subtype-specific
// fields. Every call site that constructs this wrapper (AttestedAttributesSubjectClaimWrapper's
// getSubjectId()) passes a value whose generated-JAXB field is already concretely typed
// XmlAttestedAttributesSubjectIdClaim (see jaxb_claim.go), so the instanceof check is always true
// on this side; the wrapper is ported directly against the concrete type, dropping the
// always-true runtime check but keeping its behaviour (here the check gates no further data
// condition, so isMap() is unconditionally true, matching AddressClaimWrapper's pattern).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// AttestedAttributesSubjectClaimIdWrapper wraps a jaxb.XmlAttestedAttributesSubjectIdClaim.
type AttestedAttributesSubjectClaimIdWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlAttestedAttributesSubjectIdClaim, shadowing the promoted
	// (synthetic) field of the embedded ClaimWrapper; see the covariant-getWrapped note in
	// claim_wrapper.go. NOTE: Java's getWrapped() is not itself overridden here (unlike its
	// siblings) - it keeps returning the base XmlClaim type - so no Wrapped() override is added
	// on this type either.
	wrapped *jaxb.XmlAttestedAttributesSubjectIdClaim
}

// NewAttestedAttributesSubjectClaimIdWrapper is the default constructor. Port of
// AttestedAttributesSubjectClaimIdWrapper(XmlClaim).
func NewAttestedAttributesSubjectClaimIdWrapper(wrapped *jaxb.XmlAttestedAttributesSubjectIdClaim) *AttestedAttributesSubjectClaimIdWrapper {
	return NewAttestedAttributesSubjectClaimIdWrapperWithParent(wrapped, nil)
}

// NewAttestedAttributesSubjectClaimIdWrapperWithParent is the constructor with a parent
// provided. Port of AttestedAttributesSubjectClaimIdWrapper(XmlClaim, ClaimWrapper).
func NewAttestedAttributesSubjectClaimIdWrapperWithParent(wrapped *jaxb.XmlAttestedAttributesSubjectIdClaim, parent *ClaimWrapper) *AttestedAttributesSubjectClaimIdWrapper {
	w := &AttestedAttributesSubjectClaimIdWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// FamilyName gets the family name of the attribute subject. Port of getFamilyName() (the
// wrapped instanceof XmlAttestedAttributesSubjectIdClaim branch is always taken here, see the
// file note above).
func (w *AttestedAttributesSubjectClaimIdWrapper) FamilyName() *ClaimWrapper {
	if w.wrapped.FamilyName != nil {
		return NewClaimWrapperWithParent(w.wrapped.FamilyName, &w.ClaimWrapper)
	}
	return nil
}

// GivenName gets the given name of the attribute subject. Port of getGivenName() (the wrapped
// instanceof XmlAttestedAttributesSubjectIdClaim branch is always taken here, see the file note
// above).
func (w *AttestedAttributesSubjectClaimIdWrapper) GivenName() *ClaimWrapper {
	if w.wrapped.GivenName != nil {
		return NewClaimWrapperWithParent(w.wrapped.GivenName, &w.ClaimWrapper)
	}
	return nil
}

// DocumentNumber gets the given name of the attribute subject (sic, the Java Javadoc is
// copy-pasted from getGivenName(); reproduced as-is). Port of getDocumentNumber() (the wrapped
// instanceof XmlAttestedAttributesSubjectIdClaim branch is always taken here, see the file note
// above).
func (w *AttestedAttributesSubjectClaimIdWrapper) DocumentNumber() *ClaimWrapper {
	if w.wrapped.DocumentNumber != nil {
		return NewClaimWrapperWithParent(w.wrapped.DocumentNumber, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: an AttestedAttributesSubjectClaimIdWrapper is unconditionally a map
// claim. Port of the overridden isMap() (the wrapped instanceof
// XmlAttestedAttributesSubjectIdClaim conjunct is always true here, see the file note above).
func (w *AttestedAttributesSubjectClaimIdWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated family/given-name/document-number
// child claims rather than the generic Entry list. Port of the overridden getMap() (the wrapped
// instanceof XmlAttestedAttributesSubjectIdClaim branch is always taken here, see the file note
// above).
func (w *AttestedAttributesSubjectClaimIdWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if familyName := w.FamilyName(); familyName != nil {
		result[familyName.Name()] = familyName
	}
	if givenName := w.GivenName(); givenName != nil {
		result[givenName.Name()] = givenName
	}
	if documentNumber := w.DocumentNumber(); documentNumber != nil {
		result[documentNumber.Name()] = documentNumber
	}
	return result
}

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *AttestedAttributesSubjectClaimIdWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
