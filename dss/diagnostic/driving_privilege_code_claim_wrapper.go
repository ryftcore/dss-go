// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/DrivingPrivilegeCodeClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// DrivingPrivilegeCodeClaimWrapper represents a code information of a driving privilege.
type DrivingPrivilegeCodeClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlDrivingPrivilegeCodeClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlDrivingPrivilegeCodeClaim
}

// NewDrivingPrivilegeCodeClaimWrapper is the default constructor. Port of
// DrivingPrivilegeCodeClaimWrapper(XmlDrivingPrivilegeCodeClaim).
func NewDrivingPrivilegeCodeClaimWrapper(wrapped *jaxb.XmlDrivingPrivilegeCodeClaim) *DrivingPrivilegeCodeClaimWrapper {
	return NewDrivingPrivilegeCodeClaimWrapperWithParent(wrapped, nil)
}

// NewDrivingPrivilegeCodeClaimWrapperWithParent is the constructor with a parent provided. Port
// of DrivingPrivilegeCodeClaimWrapper(XmlDrivingPrivilegeCodeClaim, ClaimWrapper).
func NewDrivingPrivilegeCodeClaimWrapperWithParent(wrapped *jaxb.XmlDrivingPrivilegeCodeClaim, parent *ClaimWrapper) *DrivingPrivilegeCodeClaimWrapper {
	return &DrivingPrivilegeCodeClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// Code gets the code. Port of getCode().
func (w *DrivingPrivilegeCodeClaimWrapper) Code() *ClaimWrapper {
	if w.wrapped.Code != nil {
		return NewClaimWrapperWithParent(w.wrapped.Code, &w.ClaimWrapper)
	}
	return nil
}

// Sign gets the sign. Port of getSign().
func (w *DrivingPrivilegeCodeClaimWrapper) Sign() *ClaimWrapper {
	if w.wrapped.Sign != nil {
		return NewClaimWrapperWithParent(w.wrapped.Sign, &w.ClaimWrapper)
	}
	return nil
}

// Value gets the value. Port of getValue().
func (w *DrivingPrivilegeCodeClaimWrapper) Value() *ClaimWrapper {
	if w.wrapped.Value != nil {
		return NewClaimWrapperWithParent(w.wrapped.Value, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: a DrivingPrivilegeCodeClaimWrapper is unconditionally a map claim. Port
// of the overridden isMap().
func (w *DrivingPrivilegeCodeClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated code/sign/value child claims
// rather than the generic Entry list. Port of the overridden getMap().
func (w *DrivingPrivilegeCodeClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if code := w.Code(); code != nil {
		result[code.Name()] = code
	}
	if sign := w.Sign(); sign != nil {
		result[sign.Name()] = sign
	}
	if value := w.Value(); value != nil {
		result[value.Name()] = value
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *DrivingPrivilegeCodeClaimWrapper) Wrapped() *jaxb.XmlDrivingPrivilegeCodeClaim {
	return w.wrapped
}

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *DrivingPrivilegeCodeClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
