// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/scope/CAdESSignatureScopeFinder.java (DSS 6.5.RC1).
//
// SCC flattening: Java eu.europa.esig.dss.cades.validation.scope is flattened into this one Go
// package.
package cades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESSignatureScopeFinder finds SignatureScopes for a CAdES signature. Port of the class
// SignatureScopeFinder, extending spiscope.AbstractSignatureScopeFinder and implementing
// spiscope.SignatureScopeFinder[*Signature].
type SignatureScopeFinder struct {
	spiscope.AbstractSignatureScopeFinder
}

// NewCAdESSignatureScopeFinder is the port of the default constructor.
func NewSignatureScopeFinder() *SignatureScopeFinder {
	return &SignatureScopeFinder{AbstractSignatureScopeFinder: spiscope.NewAbstractSignatureScopeFinder()}
}

// FindSignatureScope returns a list of SignatureScopes from a signature. Port of
// findSignatureScope(Signature).
func (f *SignatureScopeFinder) FindSignatureScope(cadesSignature *Signature) []scope.SignatureScope {
	originalDocument := f.getOriginalDocument(cadesSignature)
	if originalDocument == nil {
		return []scope.SignatureScope{}
	}

	result := make([]scope.SignatureScope, 0)
	if f.isASiCSArchive(cadesSignature) {
		containerSignatureScope := spiscope.NewContainerSignatureScope(originalDocument)
		result = append(result, containerSignatureScope)
		for _, archivedDocument := range cadesSignature.ContainerContents() {
			containerSignatureScope.AddChildSignatureScope(spiscope.NewContainerContentSignatureScope(archivedDocument))
		}

	} else if f.IsASiCEArchive(cadesSignature) {
		manifestFile := cadesSignature.ManifestFile()
		manifestSignatureScope := spiscope.NewManifestSignatureScope(manifestFile)
		result = append(result, manifestSignatureScope)

		for _, manifestEntry := range manifestFile.Entries() {
			if manifestEntry.IsIntact() {
				referencedDocument := f.getReferencedDocument(manifestEntry, cadesSignature.ContainerContents())
				manifestSignatureScope.AddChildSignatureScope(spiscope.NewFullSignatureScope(manifestEntry.Uri(), referencedDocument))
			}
		}

	} else {
		referenceValidations := cadesSignature.ReferenceValidations()
		if utils.IsCollectionNotEmpty(referenceValidations) {
			reference := referenceValidations[0] // only one Reference is allowed in CAdES
			if reference.IsIntact() {
				return f.getSignatureScopeFromOriginalDocument(cadesSignature, originalDocument)
			} else if reference.IsFound() {
				return f.getSignatureScopeFromReferenceValidation(reference)
			}
		}
	}
	return result
}

// getSignatureScopeFromOriginalDocument returns a list of SignatureScopes from the signed
// document. Port of the protected getSignatureScopeFromOriginalDocument(CAdESSignature,
// DSSDocument).
func (f *SignatureScopeFinder) getSignatureScopeFromOriginalDocument(cadesSignature *Signature,
	originalDocument model.DSSDocument) []scope.SignatureScope {
	result := make([]scope.SignatureScope, 0)
	if originalDocument == nil {
		return result
	}

	fileName := originalDocument.Name()
	if cadesSignature.IsCounterSignature() {
		return []scope.SignatureScope{spiscope.NewCounterSignatureScope(cadesSignature.MasterSignature(), originalDocument)}

	} else if digestDocument, ok := originalDocument.(*model.DigestDocument); ok {
		result = append(result, spiscope.NewDigestSignatureScope(fileName, digestDocument))

	} else {
		result = append(result, spiscope.NewFullSignatureScope(fileName, originalDocument))
	}

	return result
}

// getSignatureScopeFromReferenceValidation gets a list of SignatureScopes from a
// ReferenceValidation. Port of the protected getSignatureScopeFromReferenceValidation(
// ReferenceValidation).
func (f *SignatureScopeFinder) getSignatureScopeFromReferenceValidation(reference *model.ReferenceValidation) []scope.SignatureScope {
	result := make([]scope.SignatureScope, 0)
	digestDocument := f.CreateDigestDocument(reference.Digest())
	if digestDocument != nil {
		var fileName string
		if reference.Document() != nil {
			fileName = reference.Document().Name()
		}
		result = append(result, spiscope.NewFullSignatureScope(fileName, digestDocument))
	}
	return result
}

// getOriginalDocument returns the original document for the given CAdES signature. Port of the
// protected getOriginalDocument(CAdESSignature); the DSSException Java catches and logs around
// is swallowed here the same way, since slf4j logging is dropped per PORTING.md.
func (f *SignatureScopeFinder) getOriginalDocument(cadesSignature *Signature) model.DSSDocument {
	document, err := cadesSignature.OriginalDocument()
	if err != nil {
		return nil
	}
	return document
}

// isASiCSArchive shadows spiscope.AbstractSignatureScopeFinder.IsASiCSArchive(): a CAdES
// signature is only an ASiC-S archive when it is not also an ASiC-E archive. Port of the
// protected isASiCSArchive(AdvancedSignature) override.
func (f *SignatureScopeFinder) isASiCSArchive(advancedSignature validation.AdvancedSignature) bool {
	return f.AbstractSignatureScopeFinder.IsASiCSArchive(advancedSignature) && !f.AbstractSignatureScopeFinder.IsASiCEArchive(advancedSignature)
}

// getReferencedDocument returns a document referenced from manifestEntry. Port of the protected
// getReferencedDocument(ManifestEntry, List).
func (f *SignatureScopeFinder) getReferencedDocument(manifestEntry *model.ManifestEntry, detachedDocuments []model.DSSDocument) model.DSSDocument {
	document := spi.DSSUtilsDocumentWithName(detachedDocuments, manifestEntry.Uri())
	if document == nil {
		document = f.CreateDigestDocument(manifestEntry.Digest())
	}
	return document
}

// compile-time interface assertion.
var _ spiscope.SignatureScopeFinder[*Signature] = (*SignatureScopeFinder)(nil)
