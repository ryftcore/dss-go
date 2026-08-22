// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/EvidenceRecordScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	"slices"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// EvidenceRecordScopeFinder extracts evidence record scopes representing the covered
// archival data objects.
type EvidenceRecordScopeFinder struct {
	AbstractSignatureScopeFinder

	// EvidenceRecord is the associated evidence record.
	EvidenceRecord validation.EvidenceRecord
}

// NewEvidenceRecordScopeFinder is the default constructor. Port of
// EvidenceRecordScopeFinder(EvidenceRecord). Panics if evidenceRecord is nil, mirroring
// Java's Objects.requireNonNull(evidenceRecord, "EvidenceRecord shall be provided!").
func NewEvidenceRecordScopeFinder(evidenceRecord validation.EvidenceRecord) *EvidenceRecordScopeFinder {
	if evidenceRecord == nil {
		panic("EvidenceRecord shall be provided!")
	}
	return &EvidenceRecordScopeFinder{
		AbstractSignatureScopeFinder: NewAbstractSignatureScopeFinder(),
		EvidenceRecord:               evidenceRecord,
	}
}

// FindEvidenceRecordScope returns an evidence record scope for the associated
// EvidenceRecord. Port of findEvidenceRecordScope().
func (f *EvidenceRecordScopeFinder) FindEvidenceRecordScope() []mscope.SignatureScope {
	evidenceRecordScopes := append([]mscope.SignatureScope{}, f.FindEvidenceRecordScopeForReferences(f.EvidenceRecord.ReferenceValidation())...)
	if f.IsSignatureEmbeddedAndValid(f.EvidenceRecord) {
		evidenceRecordScopes = f.enrichRecursively(evidenceRecordScopes, f.EvidenceRecord.MasterSignature().SignatureScopes())
	}
	return evidenceRecordScopes
}

func (f *EvidenceRecordScopeFinder) enrichRecursively(evidenceRecordScopes []mscope.SignatureScope, signatureScopes []mscope.SignatureScope) []mscope.SignatureScope {
	for _, signatureScope := range signatureScopes {
		if !containsSignatureScope(evidenceRecordScopes, signatureScope) {
			evidenceRecordScopes = append(evidenceRecordScopes, signatureScope)
		} else if utils.IsCollectionNotEmpty(signatureScope.Children()) {
			evidenceRecordScopes = f.enrichRecursively(evidenceRecordScopes, signatureScope.Children())
		}
	}
	return evidenceRecordScopes
}

func containsSignatureScope(scopes []mscope.SignatureScope, target mscope.SignatureScope) bool {
	for _, s := range scopes {
		if s.Equals(target) {
			return true
		}
	}
	return false
}

// IsSignatureEmbeddedAndValid verifies whether the signature is embedded and covers the
// master signature. Port of isSignatureEmbeddedAndValid(EvidenceRecord).
func (f *EvidenceRecordScopeFinder) IsSignatureEmbeddedAndValid(evidenceRecord validation.EvidenceRecord) bool {
	if evidenceRecord.IsEmbedded() {
		for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
			if enumerations.DigestMatcherTypeEvidenceRecordMasterSignature == referenceValidation.Type() && referenceValidation.IsIntact() {
				return true
			}
		}
	}
	return false
}

// FindEvidenceRecordScopeForReferences extracts evidence record scopes for the provided list
// of reference validation and detached content. Port of
// findEvidenceRecordScope(List<ReferenceValidation>).
func (f *EvidenceRecordScopeFinder) FindEvidenceRecordScopeForReferences(referenceValidations []*model.ReferenceValidation) []mscope.SignatureScope {
	signatureScopes := make([]mscope.SignatureScope, 0)

	detachedContents := f.EvidenceRecord.DetachedContents()

	coveredDocuments := make([]model.DSSDocument, 0)
	for _, referenceValidation := range referenceValidations {
		if referenceValidation.IsIntact() {
			switch referenceValidation.Type() {
			case enumerations.DigestMatcherTypeEvidenceRecordArchiveObject:
				var detachedDocument model.DSSDocument
				if utils.CollectionSize(detachedContents) == 1 {
					detachedDocument = detachedContents[0]
				} else {
					detachedDocument = f.getDetachedDocument(referenceValidation, detachedContents)
				}
				if detachedDocument != nil && !slices.ContainsFunc(coveredDocuments, func(d model.DSSDocument) bool {
					return documentsEqual(d, detachedDocument)
				}) {
					fileName := detachedDocument.Name()
					signatureScopes = append(signatureScopes, NewFullSignatureScope(fileName, detachedDocument))
					coveredDocuments = append(coveredDocuments, detachedDocument) // do not add documents with the same digests
				}
			case enumerations.DigestMatcherTypeEvidenceRecordMasterSignature:
				masterSignature := f.EvidenceRecord.MasterSignature()
				signatureScopes = append(signatureScopes, NewEvidenceRecordMasterSignatureScope(masterSignature,
					f.CreateDigestDocument(referenceValidation.Digest())))
			default:
				// skip
			}
		}
	}

	return signatureScopes
}

func (f *EvidenceRecordScopeFinder) getDetachedDocument(referenceValidation *model.ReferenceValidation, detachedDocuments []model.DSSDocument) model.DSSDocument {
	if utils.IsCollectionNotEmpty(detachedDocuments) {
		for _, document := range detachedDocuments {
			if document.Name() == "" {
				panic("Name shall be defined when multiple documents provided!")
			}
			if referenceValidation.Document() != nil && referenceValidation.Document().Name() == document.Name() {
				return document
			}
		}
	}
	return nil
}

// documentsEqual reports whether a and b represent the same DSSDocument, mirroring Java's
// content-based DSSDocument#equals overrides (InMemoryDocument, DigestDocument, FileDocument).
// Other DSSDocument implementations fall back to Go's interface equality, the counterpart of
// Java's default reference-equality Object#equals.
func documentsEqual(a, b model.DSSDocument) bool {
	switch av := a.(type) {
	case *model.InMemoryDocument:
		bv, ok := b.(*model.InMemoryDocument)
		return ok && av.Equals(bv)
	case *model.DigestDocument:
		bv, ok := b.(*model.DigestDocument)
		return ok && av.Equals(bv)
	case *model.FileDocument:
		bv, ok := b.(*model.FileDocument)
		return ok && av.Equals(bv)
	default:
		return a == b
	}
}
