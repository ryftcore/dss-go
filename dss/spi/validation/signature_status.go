// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/status/SignatureStatus.java (DSS 6.5.RC1).
//
// SCC flattening: Java's spi.validation.status package is flattened into this package (see
// PORTING.md / S2B_BRIEF.md); the sibling status types (TokenStatus, RevocationFreshnessStatus)
// are ported by other chunks of phase 2b into the same package.
package validation

import (
	"github.com/utain/esig/dss/alert"
	"github.com/utain/esig/dss/utils"
)

// SignatureStatus contains signatures concerned by an occurred event and corresponding
// information about them.
type SignatureStatus struct {
	*alert.ObjectStatus

	// relatedSignatureMap maps signatures concerned by the check to the corresponding error
	// explanation messages. Java's Map<AdvancedSignature, String> keyed by object identity has
	// no comparable Go representation for an interface value pointing at arbitrary
	// implementations; the map is keyed by the signature's Id() (matching the convention used
	// for AdvancedSignature elsewhere in this package) while relatedSignatures preserves the
	// concrete signatures for GetRelatedSignatures().
	relatedSignatureMap map[string]string
	relatedSignatures   map[string]AdvancedSignature
}

// NewSignatureStatus is the default constructor initializing an empty map.
func NewSignatureStatus() *SignatureStatus {
	return &SignatureStatus{
		ObjectStatus:        alert.NewObjectStatus(),
		relatedSignatureMap: make(map[string]string),
		relatedSignatures:   make(map[string]AdvancedSignature),
	}
}

// AddRelatedTokenAndErrorMessage adds a concerned signature and information about the
// occurred event. Port of addRelatedTokenAndErrorMessage(AdvancedSignature, String).
func (s *SignatureStatus) AddRelatedTokenAndErrorMessage(signature AdvancedSignature, errorMessage string) {
	s.ObjectStatus.AddRelatedObjectIdentifierAndErrorMessage(signature.ID(), errorMessage)
	s.relatedSignatureMap[signature.ID()] = errorMessage
	s.relatedSignatures[signature.ID()] = signature
}

// RelatedSignatures returns a collection of signatures concerned by failure of the processed
// check. Port of getRelatedSignatures().
func (s *SignatureStatus) RelatedSignatures() []AdvancedSignature {
	signatures := make([]AdvancedSignature, 0, len(s.relatedSignatures))
	for _, signature := range s.relatedSignatures {
		signatures = append(signatures, signature)
	}
	return signatures
}

// MessageForSignature returns the error message for the given signature. Port of
// getMessageForSignature(AdvancedSignature).
func (s *SignatureStatus) MessageForSignature(signature AdvancedSignature) string {
	return s.relatedSignatureMap[signature.ID()]
}

// IsEmpty reports whether the status has no related object nor related signature.
func (s *SignatureStatus) IsEmpty() bool {
	return s.ObjectStatus.IsEmpty() && utils.IsMapEmpty(s.relatedSignatureMap)
}
