// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/DrivingPrivilegesClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// DrivingPrivilegesClaimWrapper provides user-friendly access to the information present within
// driving privileges claim.
type DrivingPrivilegesClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlDrivingPrivilegesClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlDrivingPrivilegesClaim
}

// NewDrivingPrivilegesClaimWrapper is the default constructor. Port of
// DrivingPrivilegesClaimWrapper(XmlDrivingPrivilegesClaim).
func NewDrivingPrivilegesClaimWrapper(wrapped *jaxb.XmlDrivingPrivilegesClaim) *DrivingPrivilegesClaimWrapper {
	return NewDrivingPrivilegesClaimWrapperWithParent(wrapped, nil)
}

// NewDrivingPrivilegesClaimWrapperWithParent is the constructor with a parent provided. Port of
// DrivingPrivilegesClaimWrapper(XmlDrivingPrivilegesClaim, ClaimWrapper).
func NewDrivingPrivilegesClaimWrapperWithParent(wrapped *jaxb.XmlDrivingPrivilegesClaim, parent *ClaimWrapper) *DrivingPrivilegesClaimWrapper {
	w := &DrivingPrivilegesClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// DrivingPrivileges gets a list of all driving privileges defined within the claim. Port of
// getDrivingPrivileges(): the Java stream filters out null elements before mapping; the same
// nil-skip is applied here.
func (w *DrivingPrivilegesClaimWrapper) DrivingPrivileges() []*DrivingPrivilegeClaimWrapper {
	xmlDrivingPrivileges := w.wrapped.DrivingPrivilege
	if len(xmlDrivingPrivileges) == 0 {
		return nil
	}
	result := make([]*DrivingPrivilegeClaimWrapper, 0, len(xmlDrivingPrivileges))
	for _, x := range xmlDrivingPrivileges {
		if x != nil {
			result = append(result, NewDrivingPrivilegeClaimWrapperWithParent(x, &w.ClaimWrapper))
		}
	}
	return result
}

// IsList is the override: a DrivingPrivilegesClaimWrapper is unconditionally a list claim. Port
// of the overridden isList().
func (w *DrivingPrivilegesClaimWrapper) IsList() bool { return true }

// List is the override, assembling the list from the dedicated DrivingPrivileges() child claims
// rather than the generic Item list. Port of the overridden getList().
func (w *DrivingPrivilegesClaimWrapper) List() []*ClaimWrapper {
	drivingPrivileges := w.DrivingPrivileges()
	if len(drivingPrivileges) == 0 {
		return nil
	}
	result := make([]*ClaimWrapper, 0, len(drivingPrivileges))
	for _, dp := range drivingPrivileges {
		result = append(result, dp.AsClaim())
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *DrivingPrivilegesClaimWrapper) Wrapped() *jaxb.XmlDrivingPrivilegesClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the List() override; see
// the package note in claim_wrapper.go.
func (w *DrivingPrivilegesClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.listOverride = w.List()
	return &cw
}
