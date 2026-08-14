// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/IntegrityClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// IntegrityClaimWrapper represents an integrity claim for a certain claim attribute.
type IntegrityClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlIntegrityClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlIntegrityClaim
}

// NewIntegrityClaimWrapper is the default constructor. Port of IntegrityClaimWrapper(XmlIntegrityClaim).
func NewIntegrityClaimWrapper(wrapped *jaxb.XmlIntegrityClaim) *IntegrityClaimWrapper {
	return NewIntegrityClaimWrapperWithParent(wrapped, nil)
}

// NewIntegrityClaimWrapperWithParent is the constructor with a parent claim provided. Port of
// IntegrityClaimWrapper(XmlIntegrityClaim, ClaimWrapper).
func NewIntegrityClaimWrapperWithParent(wrapped *jaxb.XmlIntegrityClaim, parent *ClaimWrapper) *IntegrityClaimWrapper {
	return &IntegrityClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// DigestAlgorithm gets the digest algorithm used for the claim hash computation. Port of
// getDigestAlgorithm().
func (w *IntegrityClaimWrapper) DigestAlgorithm() enumerations.DigestAlgorithm {
	if w.wrapped.DigestMethod != nil {
		return enumerations.DigestAlgorithm(*w.wrapped.DigestMethod)
	}
	return ""
}

// DigestValue gets the digest value of the computed claim integrity hash. Port of
// getDigestValue().
func (w *IntegrityClaimWrapper) DigestValue() []byte {
	if w.wrapped.DigestValue != nil {
		return []byte(*w.wrapped.DigestValue)
	}
	return nil
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *IntegrityClaimWrapper) Wrapped() *jaxb.XmlIntegrityClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type; isList/getList/isMap/getMap are not
// overridden here. See the package note in claim_wrapper.go.
func (w *IntegrityClaimWrapper) AsClaim() *ClaimWrapper { return &w.ClaimWrapper }
