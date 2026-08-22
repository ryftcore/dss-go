// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XAdESSignatureScopeFinder.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// SignatureScopeFinder performs operations in order to find all signed data for a XAdES
// Signature. Port of the class XAdESSignatureScopeFinder, extending
// spiscope.AbstractSignatureScopeFinder and implementing spiscope.SignatureScopeFinder[*Signature].
type SignatureScopeFinder struct {
	spiscope.AbstractSignatureScopeFinder
}

// NewXAdESSignatureScopeFinder is the port of the default constructor.
func NewSignatureScopeFinder() *SignatureScopeFinder {
	return &SignatureScopeFinder{AbstractSignatureScopeFinder: spiscope.NewAbstractSignatureScopeFinder()}
}

// FindSignatureScope returns a list of SignatureScopes from a signature. Port of
// findSignatureScope(Signature).
func (f *SignatureScopeFinder) FindSignatureScope(xadesSignature *Signature) []mscope.SignatureScope {
	result := make([]mscope.SignatureScope, 0)

	for _, xadesReferenceValidation := range xadesSignature.XAdESReferenceValidations() {
		switch xadesReferenceValidation.Type() {
		case enumerations.DigestMatcherTypeSignedProperties, enumerations.DigestMatcherTypeKeyInfo,
			enumerations.DigestMatcherTypeSignatureProperties:
			// not a subject for the Signature Scope
			continue
		}

		uri := xadesReferenceValidation.Uri()
		xmlIdOfSignedElement := xmlutils.DomUtilsGetId(uri)
		transformations := xadesReferenceValidation.TransformationNames()

		switch {
		case xadesReferenceValidation.IsFound() && enumerations.DigestMatcherTypeXPointer == xadesReferenceValidation.Type():
			result = append(result, newXPointerSignatureScope(uri,
				f.CreateInMemoryDocument(xadesReferenceValidation.OriginalContentBytes()), transformations))

		case xadesReferenceValidation.IsFound() && enumerations.DigestMatcherTypeObject == xadesReferenceValidation.Type():
			objectById := DSSXMLUtilsGetObjectById(xadesSignature.SignatureElement(), uri)
			if objectById != nil && objectById.FirstChild != nil {
				referencedObject := objectById.FirstChild
				result = append(result, newXmlElementSignatureScope(xmlIdOfSignedElement,
					f.CreateInMemoryDocument(xadesSignatureScopeFinderMustGetNodeBytes(referencedObject)), transformations))
			}

		case xadesReferenceValidation.IsFound() && enumerations.DigestMatcherTypeManifest == xadesReferenceValidation.Type():
			manifestSignatureScope := spiscope.NewManifestSignatureScopeWithNameAndTransformations(
				xadesSignatureScopeFinderReferenceName(&xadesReferenceValidation.ReferenceValidation),
				f.CreateDigestDocument(xadesReferenceValidation.Digest()), xadesReferenceValidation.TransformationNames())
			result = append(result, manifestSignatureScope)

			for _, manifestEntry := range xadesReferenceValidation.DependentValidations() {
				manifestEntryReferenceName := xadesSignatureScopeFinderReferenceName(manifestEntry)
				if manifestEntryReferenceName != "" && manifestEntry.IsFound() {
					// try to get document digest from list of detached contents
					detachedSignatureScopeResult := f.getFromDetachedContent(xadesSignature, transformations, manifestEntry)
					if detachedSignatureScopeResult != nil {
						manifestSignatureScope.AddChildSignatureScope(detachedSignatureScopeResult)
					} else if !manifestEntry.Digest().IsEmpty() {
						// if the relative detached content is not found, store the reference value
						manifestSignatureScope.AddChildSignatureScope(NewManifestEntrySignatureScope(
							manifestEntryReferenceName, f.CreateDigestDocument(manifestEntry.Digest()),
							xadesSignatureScopeFinderReferenceName(&xadesReferenceValidation.ReferenceValidation),
							manifestEntry.TransformationNames()))
					}
				}
			}

		case xadesReferenceValidation.IsFound() && enumerations.DigestMatcherTypeCounterSignature == xadesReferenceValidation.Type() &&
			xadesSignature.MasterSignature() != nil:
			result = append(result, spiscope.NewCounterSignatureScope(xadesSignature.MasterSignature(),
				f.CreateInMemoryDocument(xadesReferenceValidation.OriginalContentBytes())))

		case xadesReferenceValidation.IsFound() && utils.EmptyString == uri:
			originalContentBytes := xadesReferenceValidation.OriginalContentBytes()
			if originalContentBytes != nil {
				// self contained document
				result = append(result, newXmlRootSignatureScope(f.CreateInMemoryDocument(originalContentBytes), transformations))
			}

		case xadesReferenceValidation.IsFound() && xmlutils.DomUtilsIsElementReference(uri):
			signedElement := xmlutils.XPathUtilsGetElementById(xadesSignature.SignatureElement().OwnerDocument(), xmlutils.DomUtilsGetId(uri))
			if signedElement != nil {
				if f.isEverythingCovered(xadesSignature, xmlIdOfSignedElement) {
					result = append(result, newXmlRootSignatureScope(
						f.CreateInMemoryDocument(xadesSignatureScopeFinderMustGetNodeBytes(signedElement)), transformations))
				} else {
					result = append(result, newXmlElementSignatureScope(xmlIdOfSignedElement,
						f.CreateInMemoryDocument(xadesSignatureScopeFinderMustGetNodeBytes(signedElement)), transformations))
				}
			}

		case xadesReferenceValidation.IsIntact() && utils.IsCollectionNotEmpty(xadesSignature.DetachedContents()):
			// detached file (the signature must intact in order to be sure in the correctness of
			// the provided file)
			signatureScope := f.getFromDetachedContent(xadesSignature, transformations, &xadesReferenceValidation.ReferenceValidation)
			if signatureScope != nil {
				result = append(result, signatureScope)
			}

		case utils.IsCollectionEmpty(transformations):
			// if a matching file was not found around the detached contents and transformations
			// are not defined, use the original reference data
			result = append(result, spiscope.NewFullSignatureScope(uri, f.CreateDigestDocument(xadesReferenceValidation.Digest())))
		}
	}
	return result
}

// xadesSignatureScopeFinderReferenceName ports the private getReferenceName(ReferenceValidation).
func xadesSignatureScopeFinderReferenceName(referenceValidation *model.ReferenceValidation) string {
	if referenceValidation.Document() != nil && utils.IsStringNotEmpty(referenceValidation.Document().Name()) {
		return referenceValidation.Document().Name()
	} else if utils.IsStringNotEmpty(referenceValidation.Id()) {
		return xmlutils.DomUtilsGetId(referenceValidation.Id())
	}
	// URI can be empty "".
	return xmlutils.DomUtilsGetId(referenceValidation.Uri())
}

// getFromDetachedContent ports the private getFromDetachedContent(XAdESSignature, List,
// ReferenceValidation).
func (f *SignatureScopeFinder) getFromDetachedContent(xadesSignature *Signature,
	transformations []string, xadesReferenceValidation *model.ReferenceValidation) mscope.SignatureScope {
	detachedDocument := xadesReferenceValidation.Document()
	if detachedDocument == nil {
		return nil
	}

	fileName := detachedDocument.Name()
	if fileName == "" {
		fileName = xadesSignatureScopeFinderReferencedDocumentName(xadesReferenceValidation)
	}

	switch {
	case xadesSignatureScopeFinderIsDigestDocument(detachedDocument):
		return spiscope.NewDigestSignatureScope(fileName, detachedDocument)

	case utils.IsCollectionNotEmpty(transformations):
		return newXmlFullSignatureScope(fileName, detachedDocument, transformations)

	case f.IsASiCSArchive(xadesSignature):
		containerSignatureScope := spiscope.NewContainerSignatureScopeWithName(fileName, detachedDocument)
		for _, archivedDocument := range xadesSignature.ContainerContents() {
			containerSignatureScope.AddChildSignatureScope(spiscope.NewContainerContentSignatureScope(archivedDocument))
		}
		return containerSignatureScope

	default:
		return spiscope.NewFullSignatureScope(fileName, detachedDocument)
	}
}

// xadesSignatureScopeFinderIsDigestDocument reports whether document is a *model.DigestDocument,
// the Go counterpart of Java's `detachedDocument instanceof DigestDocument`.
func xadesSignatureScopeFinderIsDigestDocument(document model.DSSDocument) bool {
	_, ok := document.(*model.DigestDocument)
	return ok
}

// xadesSignatureScopeFinderReferencedDocumentName ports the private
// getReferencedDocumentName(ReferenceValidation).
func xadesSignatureScopeFinderReferencedDocumentName(referenceValidation *model.ReferenceValidation) string {
	if referenceValidation.Document() != nil && utils.IsStringNotEmpty(referenceValidation.Document().Name()) {
		return referenceValidation.Document().Name()
	} else if utils.IsStringNotEmpty(referenceValidation.Uri()) {
		// used for a document extraction, empty URI is not acceptable
		return xmlutils.DomUtilsGetId(referenceValidation.Uri())
	}
	return ""
}

// isEverythingCovered ports the private isEverythingCovered(XAdESSignature, String).
func (f *SignatureScopeFinder) isEverythingCovered(signature *Signature, coveredObjectId string) bool {
	parent := signature.SignatureElement().OwnerDocument().DocumentElement()
	return parent != nil && xadesSignatureScopeFinderIsRelatedToUri(parent, coveredObjectId)
}

// xadesSignatureScopeFinderIsRelatedToUri ports the private isRelatedToUri(Node, String).
func xadesSignatureScopeFinderIsRelatedToUri(currentNode *xmldom.Node, id string) bool {
	idValue := DSSXMLUtilsGetIDIdentifier(currentNode)
	if idValue == "" {
		return utils.IsStringBlank(id)
	}
	return id == idValue || id == utils.EmptyString
}

// xadesSignatureScopeFinderMustGetNodeBytes serializes node, panicking with a model.DSSError on
// failure: Java's DomUtils.getNodeBytes(Node) is called here with no surrounding try/catch, so a
// failure propagates unchecked exactly like every other uncaught DSSException in this port.
func xadesSignatureScopeFinderMustGetNodeBytes(node *xmldom.Node) []byte {
	b, err := xmlutils.DomUtilsGetNodeBytes(node)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return b
}

// compile-time interface assertion.
var _ spiscope.SignatureScopeFinder[*Signature] = (*SignatureScopeFinder)(nil)

// compile-time assertion: Signature satisfies validation.AdvancedSignature, matching the
// MasterSignature()/DetachedContents()/ContainerContents() calls above.
var _ validation.AdvancedSignature = (*Signature)(nil)
