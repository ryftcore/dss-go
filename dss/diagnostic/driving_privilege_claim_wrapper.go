// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/DrivingPrivilegeClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// DrivingPrivilegeClaimWrapper provides user-friendly access to the information present within a
// driving privilege claim.
type DrivingPrivilegeClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlDrivingPrivilegeClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlDrivingPrivilegeClaim
}

// NewDrivingPrivilegeClaimWrapper is the default constructor. Port of
// DrivingPrivilegeClaimWrapper(XmlDrivingPrivilegeClaim).
func NewDrivingPrivilegeClaimWrapper(wrapped *jaxb.XmlDrivingPrivilegeClaim) *DrivingPrivilegeClaimWrapper {
	return NewDrivingPrivilegeClaimWrapperWithParent(wrapped, nil)
}

// NewDrivingPrivilegeClaimWrapperWithParent is the constructor with a parent provided. Port of
// DrivingPrivilegeClaimWrapper(XmlDrivingPrivilegeClaim, ClaimWrapper).
func NewDrivingPrivilegeClaimWrapperWithParent(wrapped *jaxb.XmlDrivingPrivilegeClaim, parent *ClaimWrapper) *DrivingPrivilegeClaimWrapper {
	w := &DrivingPrivilegeClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// VehicleCategoryCode gets the vehicle category code. Port of getVehicleCategoryCode().
func (w *DrivingPrivilegeClaimWrapper) VehicleCategoryCode() *ClaimWrapper {
	if w.wrapped.VehicleCategoryCode != nil {
		return NewClaimWrapperWithParent(w.wrapped.VehicleCategoryCode, &w.ClaimWrapper)
	}
	return nil
}

// IssueDate gets the issuance date of the driving privilege. Port of getIssueDate().
func (w *DrivingPrivilegeClaimWrapper) IssueDate() *ClaimWrapper {
	if w.wrapped.IssueDate != nil {
		return NewClaimWrapperWithParent(w.wrapped.IssueDate, &w.ClaimWrapper)
	}
	return nil
}

// ExpiryDate gets the expiration date of the driving privilege. Port of getExpiryDate().
func (w *DrivingPrivilegeClaimWrapper) ExpiryDate() *ClaimWrapper {
	if w.wrapped.ExpiryDate != nil {
		return NewClaimWrapperWithParent(w.wrapped.ExpiryDate, &w.ClaimWrapper)
	}
	return nil
}

// Codes gets the vehicle category code. Port of getCodes().
func (w *DrivingPrivilegeClaimWrapper) Codes() *DrivingPrivilegeCodesClaimWrapper {
	if w.wrapped.Codes != nil {
		return NewDrivingPrivilegeCodesClaimWrapperWithParent(w.wrapped.Codes, &w.ClaimWrapper)
	}
	return nil
}

// IsMap is the override: a DrivingPrivilegeClaimWrapper is unconditionally a map claim. Port of
// the overridden isMap().
func (w *DrivingPrivilegeClaimWrapper) IsMap() bool { return true }

// Map is the override, assembling the map from the dedicated driving-privilege child claims
// rather than the generic Entry list. Port of the overridden getMap().
func (w *DrivingPrivilegeClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if vehicleCategoryCode := w.VehicleCategoryCode(); vehicleCategoryCode != nil {
		result[vehicleCategoryCode.Name()] = vehicleCategoryCode
	}
	if issueDate := w.IssueDate(); issueDate != nil {
		result[issueDate.Name()] = issueDate
	}
	if expiryDate := w.ExpiryDate(); expiryDate != nil {
		result[expiryDate.Name()] = expiryDate
	}
	if codes := w.Codes(); codes != nil {
		result[codes.Name()] = codes.AsClaim()
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *DrivingPrivilegeClaimWrapper) Wrapped() *jaxb.XmlDrivingPrivilegeClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the Map() override; see the
// package note in claim_wrapper.go.
func (w *DrivingPrivilegeClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	isMap := true
	cw.isMapOverride = &isMap
	return &cw
}
