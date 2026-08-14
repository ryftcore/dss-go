// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/IdentifierListClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// IdentifierListClaimWrapper wraps a jaxb.XmlIdentifierListClaim.
type IdentifierListClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlIdentifierListClaim, shadowing the promoted (synthetic) field
	// of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlIdentifierListClaim
}

// NewIdentifierListClaimWrapper is the default constructor. Port of
// IdentifierListClaimWrapper(XmlIdentifierListClaim).
func NewIdentifierListClaimWrapper(wrapped *jaxb.XmlIdentifierListClaim) *IdentifierListClaimWrapper {
	return NewIdentifierListClaimWrapperWithParent(wrapped, nil)
}

// NewIdentifierListClaimWrapperWithParent is the constructor with a parent provided. Port of
// IdentifierListClaimWrapper(XmlIdentifierListClaim, ClaimWrapper).
func NewIdentifierListClaimWrapperWithParent(wrapped *jaxb.XmlIdentifierListClaim, parent *ClaimWrapper) *IdentifierListClaimWrapper {
	return &IdentifierListClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// Identifier gets the status's unique identifier. Port of getIdentifier().
func (w *IdentifierListClaimWrapper) Identifier() *ClaimWrapper {
	if w.wrapped.Identifier != nil {
		return NewClaimWrapperWithParent(w.wrapped.Identifier, &w.ClaimWrapper)
	}
	return nil
}

// Uri gets the status's uri. Port of getUri().
func (w *IdentifierListClaimWrapper) Uri() *ClaimWrapper {
	if w.wrapped.Uri != nil {
		return NewClaimWrapperWithParent(w.wrapped.Uri, &w.ClaimWrapper)
	}
	return nil
}

// Certificate gets the certificate containing the public key that signed or sealed the
// top-level certificate in the x5chain element in the MSO revocation list structure. Port of
// getCertificate().
func (w *IdentifierListClaimWrapper) Certificate() *ClaimWrapper {
	if w.wrapped.Certificate != nil {
		return NewClaimWrapperWithParent(w.wrapped.Certificate, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: an IdentifierListClaimWrapper is unconditionally a map claim. Port of
// the overridden isMap().
func (w *IdentifierListClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated identifier/uri/certificate child
// claims rather than the generic Entry list. Port of the overridden getMap().
func (w *IdentifierListClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if identifier := w.Identifier(); identifier != nil {
		result[identifier.Name()] = identifier
	}
	if uri := w.Uri(); uri != nil {
		result[uri.Name()] = uri
	}
	if certificate := w.Certificate(); certificate != nil {
		result[certificate.Name()] = certificate
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *IdentifierListClaimWrapper) Wrapped() *jaxb.XmlIdentifierListClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *IdentifierListClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
