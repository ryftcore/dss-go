// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/DrivingPrivilegeCodesClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// DrivingPrivilegeCodesClaimWrapper represents an array of codes information for the
// corresponding driving privilege.
type DrivingPrivilegeCodesClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlDrivingPrivilegeCodesClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlDrivingPrivilegeCodesClaim
}

// NewDrivingPrivilegeCodesClaimWrapper is the default constructor. Port of
// DrivingPrivilegeCodesClaimWrapper(XmlDrivingPrivilegeCodesClaim).
func NewDrivingPrivilegeCodesClaimWrapper(wrapped *jaxb.XmlDrivingPrivilegeCodesClaim) *DrivingPrivilegeCodesClaimWrapper {
	return NewDrivingPrivilegeCodesClaimWrapperWithParent(wrapped, nil)
}

// NewDrivingPrivilegeCodesClaimWrapperWithParent is the constructor with a parent provided. Port
// of DrivingPrivilegeCodesClaimWrapper(XmlDrivingPrivilegeCodesClaim, ClaimWrapper).
func NewDrivingPrivilegeCodesClaimWrapperWithParent(wrapped *jaxb.XmlDrivingPrivilegeCodesClaim, parent *ClaimWrapper) *DrivingPrivilegeCodesClaimWrapper {
	w := &DrivingPrivilegeCodesClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// Codes gets a list of codes information for the given driving privilege. Port of getCodes():
// the Java stream filters out null elements before mapping; the same nil-skip is applied here.
func (w *DrivingPrivilegeCodesClaimWrapper) Codes() []*DrivingPrivilegeCodeClaimWrapper {
	xmlDrivingPrivilegeCodeClaims := w.wrapped.Code
	if len(xmlDrivingPrivilegeCodeClaims) == 0 {
		return nil
	}
	result := make([]*DrivingPrivilegeCodeClaimWrapper, 0, len(xmlDrivingPrivilegeCodeClaims))
	for _, x := range xmlDrivingPrivilegeCodeClaims {
		if x != nil {
			result = append(result, NewDrivingPrivilegeCodeClaimWrapperWithParent(x, &w.ClaimWrapper))
		}
	}
	return result
}

// IsList is the override: a DrivingPrivilegeCodesClaimWrapper is unconditionally a list claim.
// Port of the overridden isList().
func (w *DrivingPrivilegeCodesClaimWrapper) IsList() bool { return true }

// List is the override, assembling the list from the dedicated Codes() child claims rather than
// the generic Item list. Port of the overridden getList().
func (w *DrivingPrivilegeCodesClaimWrapper) List() []*ClaimWrapper {
	codes := w.Codes()
	if len(codes) == 0 {
		return nil
	}
	result := make([]*ClaimWrapper, 0, len(codes))
	for _, c := range codes {
		result = append(result, c.AsClaim())
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *DrivingPrivilegeCodesClaimWrapper) Wrapped() *jaxb.XmlDrivingPrivilegeCodesClaim {
	return w.wrapped
}

// AsClaim views the wrapper as its ClaimWrapper base type, baking in the List() override; see
// the package note in claim_wrapper.go.
func (w *DrivingPrivilegeCodesClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.listOverride = w.List()
	return &cw
}
