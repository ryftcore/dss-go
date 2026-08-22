// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/scope/JAdESSignatureScopeFinder.java (DSS 6.5.RC1).
//
// SCC flattening: Java's eu.europa.esig.dss.jades.validation.scope package folds into this one Go
// package per the phase-6 package layout (S6_BRIEF.md).
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's file header. This
// file additionally needs:
//
//	func (s *JAdESSignature) ReferenceValidations() []*model.ReferenceValidation // getReferenceValidations()
//
// (IsCounterSignature/IsKeyBindingSignature/MasterSignature/EAA are already part of
// validation.AdvancedSignature, which *JAdESSignature is assumed to implement.)
//
// FORWARD DEPENDENCY: HTTPHeader / HTTPHeaderDigest / HTTPHeaderSignatureScope /
// HTTPHeaderMessageBodySignatureScope (Java eu.europa.esig.dss.jades package - the "jades root"
// SCC per S6_BRIEF.md - HTTPHeader.java, HTTPHeaderDigest.java, HTTPHeaderSignatureScope.java,
// HTTPHeaderMessageBodySignatureScope.java) are not in this manifest. Assumed shapes, inferred
// from every call this file makes to them and from http_headers_payload_builder.go's own
// httpHeaderDocument interface (already landed, same package):
//
//	type HTTPHeader struct { ... } // implements model.DSSDocument, Value()/SetValue() string
//	type HTTPHeaderDigest struct { HTTPHeader; ... }
//	func (h *HTTPHeaderDigest) MessageBodyDocument() model.DSSDocument // getMessageBodyDocument()
//
//	type HTTPHeaderSignatureScope struct { scope.SignatureScope-implementing base; ... }
//	func NewHTTPHeaderSignatureScope(document model.DSSDocument) *HTTPHeaderSignatureScope
//
//	type HTTPHeaderMessageBodySignatureScope struct { ... }
//	func NewHTTPHeaderMessageBodySignatureScope(document model.DSSDocument) *HTTPHeaderMessageBodySignatureScope
//
// DSSJsonUtils.HTTP_HEADER_DIGEST is already ported as DSSJsonUtilsHTTPHeaderDigest in
// dss_json_utils.go (a sibling chunk's file, landed before this one); reused here rather than
// redeclared.
package jades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESSignatureScopeFinder finds a SignatureScope for a JAdES signature. Port of the class
// JAdESSignatureScopeFinder, extending scope.AbstractSignatureScopeFinder and implementing
// scope.SignatureScopeFinder[*JAdESSignature].
type JAdESSignatureScopeFinder struct {
	scope.AbstractSignatureScopeFinder
}

// NewJAdESSignatureScopeFinder is the port of the default constructor.
func NewJAdESSignatureScopeFinder() *JAdESSignatureScopeFinder {
	return &JAdESSignatureScopeFinder{AbstractSignatureScopeFinder: scope.NewAbstractSignatureScopeFinder()}
}

// FindSignatureScope implements scope.SignatureScopeFinder. Port of
// findSignatureScope(JAdESSignature).
func (f *JAdESSignatureScopeFinder) FindSignatureScope(jadesSignature *JAdESSignature) []modelscope.SignatureScope {
	var result []modelscope.SignatureScope

	originalDocuments := f.getOriginalDocuments(jadesSignature)
	if utils.IsCollectionEmpty(originalDocuments) {
		return result
	}

	referenceValidations := jadesSignature.ReferenceValidations()
	for _, referenceValidation := range referenceValidations {
		if !referenceValidation.IsIntact() {
			continue
		}

		if _, ok := originalDocuments[0].(*HTTPHeader); ok {
			// only http header documents shall be present
			return f.getHttpHeaderSignatureScope(originalDocuments)

		} else if len(originalDocuments) == 1 {
			if jadesSignature.IsCounterSignature() {
				// only one document shall be present
				return []modelscope.SignatureScope{scope.NewCounterSignatureScope(jadesSignature.MasterSignature(), originalDocuments[0])}
			} else if jadesSignature.IsKeyBindingSignature() {
				// only one document shall be present
				return []modelscope.SignatureScope{scope.NewKeyBindingSignatureScope(jadesSignature.EAA(), originalDocuments[0])}
			} else if jadesSignature.EAA() != nil {
				// only one document shall be present
				return []modelscope.SignatureScope{scope.NewEAASignatureScope(jadesSignature.EAA(), originalDocuments[0])}
			}
			return []modelscope.SignatureScope{f.getSignatureScopeFromOriginalDocument(originalDocuments[0], referenceValidation)}

		} else if referenceValidation.Uri() != "" {
			document := referenceValidation.Document()
			result = append(result, f.getSignatureScopeFromOriginalDocument(document, referenceValidation))

		} else if len(referenceValidations) == 1 {
			return f.getSignatureScopeFromOriginalDocuments(originalDocuments)
		}
	}

	return result
}

// getOriginalDocuments returns original documents for the given JAdES signature. Port of the
// protected getOriginalDocuments(JAdESSignature).
//
// The DSSException Java catches (logging "A JAdES signer's original document is not found
// [{}].") is swallowed the same way here, since slf4j logging is dropped per PORTING.md.
func (f *JAdESSignatureScopeFinder) getOriginalDocuments(jadesSignature *JAdESSignature) []model.DSSDocument {
	documents, err := jadesSignature.OriginalDocuments()
	if err != nil {
		return nil
	}
	return documents
}

// getSignatureScopeFromOriginalDocument returns a SignatureScope for the given originalDocument.
// Port of the protected getSignatureScopeFromOriginalDocument(DSSDocument, ReferenceValidation).
func (f *JAdESSignatureScopeFinder) getSignatureScopeFromOriginalDocument(originalDocument model.DSSDocument,
	referenceValidation *model.ReferenceValidation) modelscope.SignatureScope {
	var documentName string
	if originalDocument != nil && originalDocument.Name() != "" {
		documentName = originalDocument.Name()
	} else if referenceValidation.Uri() != "" {
		documentName = referenceValidation.Uri()
	} else if len(referenceValidation.DataObjectReferences()) == 1 {
		documentName = referenceValidation.DataObjectReferences()[0]
	}
	if digestDocument, ok := originalDocument.(*model.DigestDocument); ok {
		return scope.NewDigestSignatureScope(documentName, digestDocument)
	}
	return scope.NewFullSignatureScope(documentName, originalDocument)
}

// getSignatureScopeFromOriginalDocuments extracts a SignatureScope list from a list of original
// documents. Port of the protected getSignatureScopeFromOriginalDocuments(List).
func (f *JAdESSignatureScopeFinder) getSignatureScopeFromOriginalDocuments(originalDocuments []model.DSSDocument) []modelscope.SignatureScope {
	var result []modelscope.SignatureScope
	if utils.IsCollectionEmpty(originalDocuments) {
		return result
	}

	for _, originalDocument := range originalDocuments {
		documentName := originalDocument.Name()
		if _, ok := originalDocument.(*HTTPHeader); ok {
			// only http header documents shall be present
			return f.getHttpHeaderSignatureScope(originalDocuments)

		} else if digestDocument, ok := originalDocument.(*model.DigestDocument); ok {
			result = append(result, scope.NewDigestSignatureScope(documentName, digestDocument))

		} else {
			result = append(result, scope.NewFullSignatureScope(documentName, originalDocument))
		}
	}

	return result
}

// getHttpHeaderSignatureScope ports the private getHttpHeaderSignatureScope(List).
func (f *JAdESSignatureScopeFinder) getHttpHeaderSignatureScope(originalDocuments []model.DSSDocument) []modelscope.SignatureScope {
	var httpHeadersSignatureScopes []modelscope.SignatureScope

	httpHeadersPayloadSignatureScope := f.getHttpHeadersPayloadSignatureScope(originalDocuments)
	httpHeadersSignatureScopes = append(httpHeadersSignatureScopes, httpHeadersPayloadSignatureScope)

	for _, document := range originalDocuments {
		if httpHeaderDigest, ok := document.(*HTTPHeaderDigest); ok && DSSJsonUtilsHTTPHeaderDigest == document.Name() {
			if httpHeaderDigestSignatureScope := f.getHttpHeaderDigestSignatureScope(httpHeaderDigest); httpHeaderDigestSignatureScope != nil {
				httpHeadersSignatureScopes = append(httpHeadersSignatureScopes, httpHeaderDigestSignatureScope)
			}
			break // only one shall be present
		}
	}

	return httpHeadersSignatureScopes
}

// getHttpHeadersPayloadSignatureScope ports the private getHttpHeadersPayloadSignatureScope(List).
//
// Panics on a build failure (Java's HttpHeadersPayloadBuilder#build() throws unchecked; this
// port's Build returns an error - see http_headers_payload_builder.go).
func (f *JAdESSignatureScopeFinder) getHttpHeadersPayloadSignatureScope(originalDocuments []model.DSSDocument) modelscope.SignatureScope {
	httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(originalDocuments, false)
	payload, err := httpHeadersPayloadBuilder.Build()
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return NewHTTPHeaderSignatureScope(f.CreateInMemoryDocument(payload))
}

// getHttpHeaderDigestSignatureScope ports the private getHttpHeaderDigestSignatureScope(HTTPHeader).
func (f *JAdESSignatureScopeFinder) getHttpHeaderDigestSignatureScope(digestHttpHeader *HTTPHeaderDigest) modelscope.SignatureScope {
	digest := f.getDigest(digestHttpHeader.Value())
	if digest == nil {
		return nil
	}
	return NewHTTPHeaderMessageBodySignatureScope(digestHttpHeader.MessageBodyDocument())
}

// getDigest ports the private getDigest(String), swallowing any parse failure and returning nil,
// matching the Java catch-all try/catch(Exception).
func (f *JAdESSignatureScopeFinder) getDigest(digestHeaderValue string) (result *model.Digest) {
	defer func() {
		if recover() != nil {
			// Upstream logs "Unable to extract Digest HTTP Header value. Reason : {}".
			result = nil
		}
	}()

	valueParts := strings.Split(digestHeaderValue, "=")
	if len(valueParts) != 2 {
		// Upstream logs "Not conformant value of 'Digest' header : '{}'!".
		return nil
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForHttpHeader(valueParts[0])
	if err != nil {
		return nil
	}
	digestValue := utils.FromBase64(valueParts[1])
	digest := model.NewDigest(digestAlgorithm, digestValue)
	return &digest
}

// compile-time assertions: *JAdESSignatureScopeFinder implements scope.SignatureScopeFinder.
var _ scope.SignatureScopeFinder[*JAdESSignature] = (*JAdESSignatureScopeFinder)(nil)
