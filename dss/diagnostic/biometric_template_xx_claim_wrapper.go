// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/BiometricTemplateXXClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// BiometricTemplateXXClaimWrapper wraps a jaxb.XmlBiometricTemplateXXClaim.
type BiometricTemplateXXClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlBiometricTemplateXXClaim, shadowing the promoted (synthetic)
	// field of the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlBiometricTemplateXXClaim
}

// NewBiometricTemplateXXClaimWrapper is the default constructor. Port of
// BiometricTemplateXXClaimWrapper(XmlBiometricTemplateXXClaim).
func NewBiometricTemplateXXClaimWrapper(wrapped *jaxb.XmlBiometricTemplateXXClaim) *BiometricTemplateXXClaimWrapper {
	return NewBiometricTemplateXXClaimWrapperWithParent(wrapped, nil)
}

// NewBiometricTemplateXXClaimWrapperWithParent is the constructor with a parent claim provided.
// Port of BiometricTemplateXXClaimWrapper(XmlBiometricTemplateXXClaim, ClaimWrapper).
func NewBiometricTemplateXXClaimWrapperWithParent(wrapped *jaxb.XmlBiometricTemplateXXClaim, parent *ClaimWrapper) *BiometricTemplateXXClaimWrapper {
	return &BiometricTemplateXXClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// Type gets the type of the corresponding biometric template information as defined in the
// claim. Port of getType().
func (w *BiometricTemplateXXClaimWrapper) Type() string {
	if w.wrapped.Type != nil {
		return *w.wrapped.Type
	}
	return ""
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *BiometricTemplateXXClaimWrapper) Wrapped() *jaxb.XmlBiometricTemplateXXClaim {
	return w.wrapped
}

// AsClaim views the wrapper as its ClaimWrapper base type; isList/getList/isMap/getMap are not
// overridden here. See the package note in claim_wrapper.go.
func (w *BiometricTemplateXXClaimWrapper) AsClaim() *ClaimWrapper { return &w.ClaimWrapper }
