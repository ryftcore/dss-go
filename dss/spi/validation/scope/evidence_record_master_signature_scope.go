// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/EvidenceRecordMasterSignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EvidenceRecordMasterSignatureScope defines a master signature scope covered by an embedded
// evidence record.
type EvidenceRecordMasterSignatureScope struct {
	CounterSignatureScope
}

// NewEvidenceRecordMasterSignatureScope is the default constructor. Port of
// EvidenceRecordMasterSignatureScope(AdvancedSignature, DSSDocument).
func NewEvidenceRecordMasterSignatureScope(masterSignature validation.AdvancedSignature, originalDocument model.DSSDocument) *EvidenceRecordMasterSignatureScope {
	return &EvidenceRecordMasterSignatureScope{CounterSignatureScope: *NewCounterSignatureScope(masterSignature, originalDocument)}
}

// Type returns the type of the signature scope. Port of getType().
func (s *EvidenceRecordMasterSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeTypeSignature
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*EvidenceRecordMasterSignatureScope)(nil)
