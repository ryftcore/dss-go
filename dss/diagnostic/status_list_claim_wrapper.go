// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/StatusListClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

// StatusListClaimWrapper wraps a jaxb.XmlStatusListClaim.
type StatusListClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlStatusListClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlStatusListClaim
}

// NewStatusListClaimWrapper is the default constructor. Port of
// StatusListClaimWrapper(XmlStatusListClaim).
func NewStatusListClaimWrapper(wrapped *jaxb.XmlStatusListClaim) *StatusListClaimWrapper {
	return NewStatusListClaimWrapperWithParent(wrapped, nil)
}

// NewStatusListClaimWrapperWithParent is the constructor with a parent provided. Port of
// StatusListClaimWrapper(XmlStatusListClaim, ClaimWrapper).
func NewStatusListClaimWrapperWithParent(wrapped *jaxb.XmlStatusListClaim, parent *ClaimWrapper) *StatusListClaimWrapper {
	w := &StatusListClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// Index gets the status's unique index identifier. Port of getIndex().
func (w *StatusListClaimWrapper) Index() *ClaimWrapper {
	if w.wrapped.Index != nil {
		return NewClaimWrapperWithParent(w.wrapped.Index, &w.ClaimWrapper)
	}
	return nil
}

// Uri gets the status's uri. Port of getUri().
func (w *StatusListClaimWrapper) Uri() *ClaimWrapper {
	if w.wrapped.Uri != nil {
		return NewClaimWrapperWithParent(w.wrapped.Uri, &w.ClaimWrapper)
	}
	return nil
}

// Certificate gets the certificate containing the public key that signed or sealed the
// top-level certificate in the x5chain element in the MSO revocation list structure. Port of
// getCertificate().
func (w *StatusListClaimWrapper) Certificate() *ClaimWrapper {
	if w.wrapped.Certificate != nil {
		return NewClaimWrapperWithParent(w.wrapped.Certificate, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: a StatusListClaimWrapper is unconditionally a map claim. Port of the
// overridden isMap().
func (w *StatusListClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated index/uri/certificate child claims
// rather than the generic Entry list. Port of the overridden getMap().
func (w *StatusListClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if index := w.Index(); index != nil {
		result[index.Name()] = index
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
func (w *StatusListClaimWrapper) Wrapped() *jaxb.XmlStatusListClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *StatusListClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
