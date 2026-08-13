// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/scope/CAdESEvidenceRecordScopeFinder.java (DSS 6.5.RC1).
//
// SCC flattening: see cades_signature_scope_finder.go's header - the same applies here.
package cades

import (
	"fmt"

	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
	spiscope "github.com/utain/esig/dss/spi/validation/scope"
)

// CAdESEvidenceRecordScopeFinder builds a list of covered scopes for a CAdES embedded
// Evidence Record. Port of the class CAdESEvidenceRecordScopeFinder, extending
// spiscope.EvidenceRecordScopeFinder.
type CAdESEvidenceRecordScopeFinder struct {
	spiscope.EvidenceRecordScopeFinder

	// signature is the signature to cover.
	signature validation.AdvancedSignature
}

// NewCAdESEvidenceRecordScopeFinder is the port of the constructor
// CAdESEvidenceRecordScopeFinder(EvidenceRecord, AdvancedSignature).
func NewCAdESEvidenceRecordScopeFinder(evidenceRecord validation.EvidenceRecord, signature validation.AdvancedSignature) *CAdESEvidenceRecordScopeFinder {
	return &CAdESEvidenceRecordScopeFinder{
		EvidenceRecordScopeFinder: *spiscope.NewEvidenceRecordScopeFinder(evidenceRecord),
		signature:                 signature,
	}
}

// FindEvidenceRecordScope returns a list of covered scopes for the CAdES embedded Evidence
// Record. Port of the findEvidenceRecordScope() override.
func (f *CAdESEvidenceRecordScopeFinder) FindEvidenceRecordScope() []mscope.SignatureScope {
	evidenceRecordScopes := f.EvidenceRecordScopeFinder.FindEvidenceRecordScope()
	if f.IsSignatureEmbeddedAndValid(f.EvidenceRecord) && f.isSignatureCovered(f.EvidenceRecord, f.signature) {
		evidenceRecordScopes = append(evidenceRecordScopes,
			newEvidenceRecordCAdESSignatureScope(f.signature, f.getCAdESSignatureDocument(f.signature)))
	}
	return evidenceRecordScopes
}

// isSignatureCovered ports the private isSignatureCovered(EvidenceRecord, AdvancedSignature).
func (f *CAdESEvidenceRecordScopeFinder) isSignatureCovered(evidenceRecord validation.EvidenceRecord, signature validation.AdvancedSignature) bool {
	masterSignature := evidenceRecord.MasterSignature().(*CAdESSignature)
	cadesSignature := signature.(*CAdESSignature)
	return masterSignature.CMS() == cadesSignature.CMS()
}

// getCAdESSignatureDocument ports the private getCAdESSignatureDocument(AdvancedSignature).
// TODO (upstream): improve?
func (f *CAdESEvidenceRecordScopeFinder) getCAdESSignatureDocument(signature validation.AdvancedSignature) model.DSSDocument {
	cadesSignature := signature.(*CAdESSignature)
	derEncoded := cadesSignature.SignerInformation().DER()
	return f.CreateInMemoryDocument(derEncoded)
}

// evidenceRecordCAdESSignatureScope is used for an evidence record scope definition, covering
// the same CMS signature. Port of the private static class EvidenceRecordCAdESSignatureScope,
// extending spiscope.EvidenceRecordMasterSignatureScope.
type evidenceRecordCAdESSignatureScope struct {
	spiscope.EvidenceRecordMasterSignatureScope
}

// newEvidenceRecordCAdESSignatureScope is the port of the constructor
// EvidenceRecordCAdESSignatureScope(AdvancedSignature, DSSDocument).
func newEvidenceRecordCAdESSignatureScope(masterSignature validation.AdvancedSignature, originalDocument model.DSSDocument) *evidenceRecordCAdESSignatureScope {
	return &evidenceRecordCAdESSignatureScope{
		EvidenceRecordMasterSignatureScope: *spiscope.NewEvidenceRecordMasterSignatureScope(masterSignature, originalDocument),
	}
}

// Description shadows spiscope.EvidenceRecordMasterSignatureScope.Description() (via
// spiscope.CounterSignatureScope). Port of the getDescription(TokenIdentifierProvider) override.
func (s *evidenceRecordCAdESSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("Signature with Id : %s", tokenIdentifierProvider.IDAsString(s.MasterSignature))
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*evidenceRecordCAdESSignatureScope)(nil)
