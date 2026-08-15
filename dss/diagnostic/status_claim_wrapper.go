// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/StatusClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// StatusClaimWrapper wraps a jaxb.XmlStatusClaim. Unlike most claim subtype wrappers, Java
// declares only the single-argument constructor here (no parent-taking overload).
type StatusClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlStatusClaim, shadowing the promoted (synthetic) field of the
	// embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlStatusClaim
}

// NewStatusClaimWrapper is the default constructor. Port of StatusClaimWrapper(XmlStatusClaim).
func NewStatusClaimWrapper(wrapped *jaxb.XmlStatusClaim) *StatusClaimWrapper {
	w := &StatusClaimWrapper{
		ClaimWrapper: *NewClaimWrapper(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs)),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// StatusList gets the status list. Port of getStatusList().
func (w *StatusClaimWrapper) StatusList() *StatusListClaimWrapper {
	if w.wrapped.StatusList != nil {
		return NewStatusListClaimWrapperWithParent(w.wrapped.StatusList, &w.ClaimWrapper)
	}
	return nil
}

// IdentifierList gets the identifier list. Port of getIdentifierList().
func (w *StatusClaimWrapper) IdentifierList() *IdentifierListClaimWrapper {
	if w.wrapped.IdentifierList != nil {
		return NewIdentifierListClaimWrapperWithParent(w.wrapped.IdentifierList, &w.ClaimWrapper)
	}
	return nil
}

// Index gets the status's unique index identifier. Port of getIndex().
func (w *StatusClaimWrapper) Index() *ClaimWrapper {
	if w.wrapped.Index != nil {
		return NewClaimWrapperWithParent(w.wrapped.Index, &w.ClaimWrapper)
	}
	return nil
}

// Uri gets the status's uri. Port of getUri().
func (w *StatusClaimWrapper) Uri() *ClaimWrapper {
	if w.wrapped.Uri != nil {
		return NewClaimWrapperWithParent(w.wrapped.Uri, &w.ClaimWrapper)
	}
	return nil
}

// Type gets the status's type. Port of getType().
func (w *StatusClaimWrapper) Type() *ClaimWrapper {
	if w.wrapped.Type != nil {
		return NewClaimWrapperWithParent(w.wrapped.Type, &w.ClaimWrapper)
	}
	return nil
}

// Purpose gets the status's purpose. Port of getPurpose().
func (w *StatusClaimWrapper) Purpose() *ClaimWrapper {
	if w.wrapped.Purpose != nil {
		return NewClaimWrapperWithParent(w.wrapped.Purpose, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: a StatusClaimWrapper is unconditionally a map claim. Port of the
// overridden isMap().
func (w *StatusClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated status child claims rather than
// the generic Entry list. Port of the overridden getMap().
func (w *StatusClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if statusList := w.StatusList(); statusList != nil {
		result[statusList.Name()] = statusList.AsClaim()
	}
	if identifierList := w.IdentifierList(); identifierList != nil {
		result[identifierList.Name()] = identifierList.AsClaim()
	}
	if index := w.Index(); index != nil {
		result[index.Name()] = index
	}
	if uri := w.Uri(); uri != nil {
		result[uri.Name()] = uri
	}
	if typ := w.Type(); typ != nil {
		result[typ.Name()] = typ
	}
	if purpose := w.Purpose(); purpose != nil {
		result[purpose.Name()] = purpose
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *StatusClaimWrapper) Wrapped() *jaxb.XmlStatusClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *StatusClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
