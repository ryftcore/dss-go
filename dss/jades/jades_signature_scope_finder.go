// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/scope/JAdESSignatureScopeFinder.java (DSS 6.5.RC1).
//
// SCC flattening: Java's eu.europa.esig.dss.jades.validation.scope package is flattened into
// this one Go package.
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
// SignatureScopeFinder, extending scope.AbstractSignatureScopeFinder and implementing
// scope.SignatureScopeFinder[*Signature].
type SignatureScopeFinder struct {
	scope.AbstractSignatureScopeFinder
}

// NewJAdESSignatureScopeFinder is the port of the default constructor.
func NewSignatureScopeFinder() *SignatureScopeFinder {
	return &SignatureScopeFinder{AbstractSignatureScopeFinder: scope.NewAbstractSignatureScopeFinder()}
}

// FindSignatureScope implements scope.SignatureScopeFinder. Port of
// findSignatureScope(Signature).
func (f *SignatureScopeFinder) FindSignatureScope(jadesSignature *Signature) []modelscope.SignatureScope {
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
// protected getOriginalDocuments(Signature).
//
// The DSSException Java catches (logging "A JAdES signer's original document is not found
// [{}].") is swallowed the same way here, since slf4j logging is dropped per PORTING.md.
func (f *SignatureScopeFinder) getOriginalDocuments(jadesSignature *Signature) []model.DSSDocument {
	documents, err := jadesSignature.OriginalDocuments()
	if err != nil {
		return nil
	}
	return documents
}

// getSignatureScopeFromOriginalDocument returns a SignatureScope for the given originalDocument.
// Port of the protected getSignatureScopeFromOriginalDocument(DSSDocument, ReferenceValidation).
func (f *SignatureScopeFinder) getSignatureScopeFromOriginalDocument(originalDocument model.DSSDocument,
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
func (f *SignatureScopeFinder) getSignatureScopeFromOriginalDocuments(originalDocuments []model.DSSDocument) []modelscope.SignatureScope {
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
func (f *SignatureScopeFinder) getHttpHeaderSignatureScope(originalDocuments []model.DSSDocument) []modelscope.SignatureScope {
	var httpHeadersSignatureScopes []modelscope.SignatureScope

	httpHeadersPayloadSignatureScope := f.getHttpHeadersPayloadSignatureScope(originalDocuments)
	httpHeadersSignatureScopes = append(httpHeadersSignatureScopes, httpHeadersPayloadSignatureScope)

	for _, document := range originalDocuments {
		// Java: DSSJsonUtils.HTTP_HEADER_DIGEST.equals(document.getName()) && document instanceof HTTPHeader
		// (HTTPHeader or its HTTPHeaderDigest subclass).
		if httpHeader, ok := httpHeadersPayloadBuilderAsHTTPHeader(document); ok && DSSJsonUtilsHTTPHeaderDigest == document.Name() {
			if httpHeaderDigestSignatureScope := f.getHttpHeaderDigestSignatureScope(httpHeader); httpHeaderDigestSignatureScope != nil {
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
func (f *SignatureScopeFinder) getHttpHeadersPayloadSignatureScope(originalDocuments []model.DSSDocument) modelscope.SignatureScope {
	httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(originalDocuments, false)
	payload, err := httpHeadersPayloadBuilder.Build()
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return NewHTTPHeaderSignatureScope(f.CreateInMemoryDocument(payload))
}

// getHttpHeaderDigestSignatureScope ports the private getHttpHeaderDigestSignatureScope(HTTPHeader).
//
// A HTTPHeaderDigest carries its message body document; a plain HTTPHeader named 'Digest' only
// carries the digest value, which is exposed as a digest document.
func (f *SignatureScopeFinder) getHttpHeaderDigestSignatureScope(digestHttpHeader httpHeaderDocument) modelscope.SignatureScope {
	digest := f.getDigest(digestHttpHeader.Value())
	if digest == nil {
		return nil
	}
	if httpHeaderDigest, ok := digestHttpHeader.(*HTTPHeaderDigest); ok {
		return NewHTTPHeaderMessageBodySignatureScope(httpHeaderDigest.MessageBodyDocument())
	}
	return NewHTTPHeaderMessageBodySignatureScope(f.CreateDigestDocument(*digest))
}

// javaSplit reproduces java.lang.String#split(String) for a literal single-character separator:
// trailing empty strings are removed from the result (so "a=b=" yields ["a", "b"]), unless the
// separator does not occur, in which case the input itself is the only element. Go's
// strings.Split keeps the trailing empty strings.
func javaSplit(s, sep string) []string {
	if !strings.Contains(s, sep) {
		return []string{s}
	}
	parts := strings.Split(s, sep)
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// getDigest ports the private getDigest(String), swallowing any parse failure and returning nil,
// matching the Java catch-all try/catch(Exception).
func (f *SignatureScopeFinder) getDigest(digestHeaderValue string) (result *model.Digest) {
	defer func() {
		if recover() != nil {
			// Upstream logs "Unable to extract Digest HTTP Header value. Reason : {}".
			result = nil
		}
	}()

	// The RFC 3230 instance digest is 'algo=' + padded base64, so it usually ends with '='; Java's
	// String#split drops those trailing empty strings, strings.Split does not.
	valueParts := javaSplit(digestHeaderValue, "=")
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

// compile-time assertions: *SignatureScopeFinder implements scope.SignatureScopeFinder.
var _ scope.SignatureScopeFinder[*Signature] = (*SignatureScopeFinder)(nil)
