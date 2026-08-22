// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/XAdESSignatureUtils.java (DSS 6.5.RC1).
//
// org.apache.xml.security.signature.Reference is xmldsig.Reference per internal/xmldsig's doc.go
// mapping table; Reference#getReferencedBytes is Reference.ReferencedBytes.
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug call site is a Go
// catch-and-continue (best effort) instead.
package xades

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESSignatureUtilsGetSignerDocuments returns the list of original signed documents. Ports
// the final class XAdESSignatureUtils's static getSignerDocuments(Signature).
func SignatureUtilsGetSignerDocuments(sig *Signature) []model.DSSDocument {
	result := []model.DSSDocument{}

	signatureCryptographicVerification := sig.SignatureCryptographicVerification()
	if signatureCryptographicVerification == nil || !signatureCryptographicVerification.IsSignatureValid() {
		return result
	}
	references := sig.References()
	if utils.IsCollectionNotEmpty(references) {
		for _, reference := range references {
			if !DSSXMLUtilsIsSignedProperties(reference, sig.XAdESPaths()) {
				referenceDocument := xadesSignatureUtilsGetReferenceDocument(reference, sig)
				if referenceDocument != nil {
					result = append(result, referenceDocument)
				}
			}
			// Upstream logs "Not able to extract an original content for a reference with name
			// '{}' and URI '{}'. Reason : {}" on a DSSException from getReferenceDocument; this
			// port's helper does not throw, so nothing is caught here.
		}
	}
	return result
}

// xadesSignatureUtilsGetReferenceDocument ports the private getReferenceDocument(Reference,
// Signature).
func xadesSignatureUtilsGetReferenceDocument(reference *xmldsig.Reference, sig *Signature) model.DSSDocument {
	if document := xadesSignatureUtilsGetDSObject(reference, sig); document != nil {
		return document
	}
	if document := xadesSignatureUtilsGetDSManifest(reference, sig); document != nil {
		return document
	}

	// if not an object or object has not been found
	referencedBytes, err := reference.ReferencedBytes()
	if err != nil {
		// Upstream logs "Unable to retrieve reference {}. Reason : {}".
		return nil
	}
	if referencedBytes != nil {
		return model.NewInMemoryDocumentWithName(referencedBytes, reference.URI())
	}
	// Upstream logs "Reference bytes returned null value : {}" then
	// "A referenced document not found for a reference with Id : [{}]".
	return nil
}

// xadesSignatureUtilsGetDSObject ports the private getDSObject(Reference, XAdESSignature).
func xadesSignatureUtilsGetDSObject(reference *xmldsig.Reference, sig *Signature) model.DSSDocument {
	defer func() {
		// Upstream catches a broad Exception here and logs "An error occurred during an attempt
		// to extract signed object. Reason : {}"; a panic from a malformed same-document URI is
		// the Go analogue and is likewise swallowed.
		recover()
	}()
	if !reference.TypeIsReferenceToObject() && reference.Type() != "" {
		return nil
	}
	objectId := xmlutils.DomUtilsGetId(reference.URI())
	objectById := DSSXMLUtilsGetObjectById(sig.SignatureElement(), objectId)
	if objectById == nil || objectById.FirstChild == nil {
		return nil
	}
	var buf bytes.Buffer
	for child := objectById.FirstChild; child != nil; child = child.NextSibling {
		nodeBytes, err := xmlutils.DomUtilsGetNodeBytes(child)
		if err == nil && nodeBytes != nil {
			buf.Write(nodeBytes)
		}
	}
	return model.NewInMemoryDocumentWithMimeType(buf.Bytes(), objectId, enumerations.MimeTypeEnumXML)
}

// xadesSignatureUtilsGetDSManifest ports the private getDSManifest(Reference, XAdESSignature).
func xadesSignatureUtilsGetDSManifest(reference *xmldsig.Reference, sig *Signature) model.DSSDocument {
	defer func() {
		// Upstream catches a broad Exception here and logs "An error occurred during an attempt
		// to extract signed manifest. Reason : {}".
		recover()
	}()
	if !reference.TypeIsReferenceToManifest() && reference.Type() != "" {
		return nil
	}
	manifestId := xmlutils.DomUtilsGetId(reference.URI())
	manifestById := DSSXMLUtilsGetManifestById(sig.SignatureElement(), manifestId)
	if manifestById == nil {
		return nil
	}
	bytesValue, err := xmlutils.DomUtilsGetNodeBytes(manifestById)
	if err != nil || bytesValue == nil {
		return nil
	}
	return model.NewInMemoryDocumentWithMimeType(bytesValue, manifestId, enumerations.MimeTypeEnumXML)
}

// SignatureUtilsIsKeyInfoCovered verifies whether the ds:KeyInfo element is signed by the
// signature. Ports the static isKeyInfoCovered(XAdESSignature).
func SignatureUtilsIsKeyInfoCovered(sig *Signature) bool {
	referenceValidations := sig.ReferenceValidations()
	if utils.IsCollectionNotEmpty(referenceValidations) {
		for _, referenceValidation := range referenceValidations {
			if enumerations.DigestMatcherTypeKeyInfo == referenceValidation.Type() &&
				referenceValidation.IsFound() && referenceValidation.IsIntact() {
				return true
			}
		}
	}
	return false
}

// SignatureUtilsGetLastSealingEvidenceRecordAttribute returns the latest
// "SealingEvidenceRecords" unsigned property, when present. Ports the static
// getLastSealingEvidenceRecordAttribute(UnsignedSigProperties).
func SignatureUtilsGetLastSealingEvidenceRecordAttribute(unsignedSigProperties *UnsignedSigProperties) *Attribute {
	// Execute in reverse order in order to change only last evidence-record, when applicable.
	attributes := unsignedSigProperties.Attributes()
	for i := len(attributes) - 1; i >= 0; i-- {
		attribute := attributes[i]
		if definition.XAdESEvidencerecordNamespaceElementSealingEvidenceRecords.IsSameTagName(attribute.Name()) {
			return attribute
		}
	}
	return nil
}
