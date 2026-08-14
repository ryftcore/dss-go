// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentAnalyzer.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES (sibling chunks of this same phase; assumed shapes, following the
// cms_document_analyzer.go / pades_signature.go precedent for how such a header documents
// them):
//
//	func PAdESUtilsIsPDFDocument(document model.DSSDocument) bool  // PAdESUtils.isPDFDocument
//	func PAdESUtilsGetOriginalPDF(padesSignature *PAdESSignature) model.DSSDocument
//
//	// pades/dss (flattened into this package)
//	type PdfDocDssRevision struct { ... }                    // implements PdfRevision
//	func (r *PdfDocDssRevision) CertificateSource() *PdfDssDictCertificateSource
//	func (r *PdfDocDssRevision) CRLSource() *PdfDssDictCRLSource
//	func (r *PdfDocDssRevision) OCSPSource() *PdfDssDictOCSPSource
//	func (r *PdfDocDssRevision) DssDictionary() PdfDssDict
//
//	// pades root (flattened into this package)
//	type PdfDocTimestampRevision struct { ... }              // implements PdfRevision
//	func (r *PdfDocTimestampRevision) TimestampToken() *PdfTimestampToken
//	func (r *PdfDocTimestampRevision) ByteRange() *ByteRange
//
//	// pades/timestamp (flattened into this package)
//	// INTEGRATION CORRECTION: the actually-landed pdf_timestamp_token.go embeds
//	// *validation.TimestampToken (by pointer), not validation.TimestampToken (by value) as
//	// speculated below - so timestampToken.TimestampToken already IS *validation.TimestampToken;
//	// this file's two call sites that used to read &timestampToken.TimestampToken were fixed to
//	// drop the &, and the append() feeding a []*validation.TimestampToken was fixed to append
//	// timestampToken.TimestampToken instead of the *PdfTimestampToken wrapper.
//	type PdfTimestampToken struct { *validation.TimestampToken }
//	func NewPAdESTimestampScopeFinder() *PAdESTimestampScopeFinder
//	func (f *PAdESTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []scope.SignatureScope
//	func NewPdfRevisionTimestampSource(pdfRevision PdfRevision, certificateSource *spi.ListCertificateSource,
//	    crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) *PdfRevisionTimestampSource
//	func (s *PdfRevisionTimestampSource) IncorporatedReferences() []*validation.TimestampedReference
//
// java.util.Objects.requireNonNull("Document to be validated cannot be null!") stays a panic
// (matching CMSDocumentAnalyzer's NewCMSDocumentAnalyzerFromDocument precedent), while the
// IllegalInputException("Not supported document") thrown for an unsupported document becomes a
// returned error, following the same precedent.
package pades

import (
	"unsafe"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
	timestampsrc "github.com/utain/esig/dss/spi/validation/timestamp"
	"github.com/utain/esig/dss/utils"
)

// PDFDocumentAnalyzer is the validation of a PDF document. Port of the class
// PDFDocumentAnalyzer, extending analyzer.DefaultDocumentAnalyzer.
type PDFDocumentAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// pdfObjectFactory loads the relevant implementation for PDF document reading.
	pdfObjectFactory IPdfObjFactory

	// documentRevisions is the list of PDF document revisions.
	documentRevisions []PdfRevision

	// passwordProtection is the PDF document password (for protected documents).
	passwordProtection []byte
}

// newPDFDocumentAnalyzer is the port of the protected empty constructor.
func newPDFDocumentAnalyzer() *PDFDocumentAnalyzer {
	a := &PDFDocumentAnalyzer{
		DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase(),
		pdfObjectFactory:        NewDefaultPdfObjFactory(),
	}
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

// NewPDFDocumentAnalyzer creates a PDFDocumentAnalyzer from a DSSDocument. Port of the default
// constructor PDFDocumentAnalyzer(DSSDocument).
//
// Panics when document is nil (Objects.requireNonNull("Document to be validated cannot be
// null!")). IllegalInputException is an unchecked (Java RuntimeException) exception, so per
// PORTING.md the "Not supported document" it throws for an unsupported document panics too,
// carrying the exception.IllegalInputException value - matching the landed call sites
// (pades_level_baseline_t.go's PDFDocumentValidator, pades_extension_service.go's
// pdfDocumentAnalyzer) which use NewPDFDocumentAnalyzer(document) as a single-value expression.
func NewPDFDocumentAnalyzer(document model.DSSDocument) *PDFDocumentAnalyzer {
	if document == nil {
		panic("Document to be validated cannot be null!")
	}
	a := newPDFDocumentAnalyzer()
	if !a.IsSupported(document) {
		panic(exception.NewIllegalInputException("Not supported document"))
	}
	a.SetDocument(document)
	return a
}

// IsSupported checks if the document is supported by the current validator.
// Port of isSupported(DSSDocument).
func (a *PDFDocumentAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	return PAdESUtilsIsPDFDocument(dssDocument)
}

// SetPdfObjFactory sets the IPdfObjFactory. Allow to set the used implementation. Cannot be nil.
// Port of setPdfObjFactory(IPdfObjFactory).
func (a *PDFDocumentAnalyzer) SetPdfObjFactory(pdfObjFactory IPdfObjFactory) {
	if pdfObjFactory == nil {
		panic("PdfObjFactory is null")
	}
	a.pdfObjectFactory = pdfObjFactory
}

// SetPasswordProtection specifies the used password for the encrypted document.
// Port of setPasswordProtection(char[]).
func (a *PDFDocumentAnalyzer) SetPasswordProtection(passwordProtection []byte) {
	a.passwordProtection = passwordProtection
}

// PrepareValidationContext initializes and fills a ValidationContext with the necessary data
// sources. Port of the protected prepareValidationContext(Collection, Collection, Collection,
// CertificateVerifier) override.
func (a *PDFDocumentAnalyzer) PrepareValidationContext(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken, detachedEvidenceRecords []validation.EvidenceRecord,
	certificateVerifier validation.CertificateVerifier) validation.ValidationContext {
	validationContext := a.DefaultDocumentAnalyzer.PrepareValidationContext(signatures, detachedTimestamps,
		detachedEvidenceRecords, certificateVerifier)
	dssRevisions := a.DssRevisions()
	a.prepareDssDictionaryValidationContext(validationContext, dssRevisions)
	return validationContext
}

// prepareDssDictionaryValidationContext fills the validationContext with certificate tokens
// from dssRevisions. Port of the protected prepareDssDictionaryValidationContext(ValidationContext, List).
func (a *PDFDocumentAnalyzer) prepareDssDictionaryValidationContext(validationContext validation.ValidationContext,
	dssRevisions []*PdfDocDssRevision) {
	for _, dssRevision := range dssRevisions {
		validationContext.AddDocumentCertificateSource(dssRevision.CertificateSource())
		validationContext.AddDocumentCRLSource(dssRevision.CRLSource())
		validationContext.AddDocumentOCSPSource(dssRevision.OCSPSource())
	}
}

// GetAllSignatures returns a list of all signatures from the validating document, post-processed
// for shadow-attack detection. Port of the protected getAllSignatures() override.
func (a *PDFDocumentAnalyzer) GetAllSignatures() []validation.AdvancedSignature {
	allSignatures := a.DefaultDocumentAnalyzer.GetAllSignatures()
	a.postProcessing(allSignatures)
	return allSignatures
}

// postProcessing post-processes the extracted signatures.
// NOTE: the method shall be used only for the document validation.
// Port of the protected postProcessing(List) override.
func (a *PDFDocumentAnalyzer) postProcessing(signatures []validation.AdvancedSignature) {
	pdfSignatureService := a.pdfObjectFactory.NewPAdESSignatureService()
	pdfSignatureService.AnalyzePdfModifications(a.Document(), signatures, a.passwordProtection)
}

// DetachedTimestamps retrieves the detached timestamps found in the document, post-processed
// for shadow-attack detection. Port of the getDetachedTimestamps() override.
func (a *PDFDocumentAnalyzer) DetachedTimestamps() []*validation.TimestampToken {
	detachedTimestamps := a.DefaultDocumentAnalyzer.DetachedTimestamps()
	a.timestampPostProcessing(detachedTimestamps)
	return detachedTimestamps
}

// timestampPostProcessing post-processes the extracted detached timestamps.
// NOTE: the method shall be used only for the document validation.
// Port of the protected timestampPostProcessing(List) override.
func (a *PDFDocumentAnalyzer) timestampPostProcessing(timestampTokens []*validation.TimestampToken) {
	pdfSignatureService := a.pdfObjectFactory.NewPAdESSignatureService()
	pdfSignatureService.AnalyzeTimestampPdfModifications(a.Document(), timestampTokens, a.passwordProtection)
}

// BuildSignatures builds a list of signatures to be extracted from a document.
// Port of the protected buildSignatures() override.
func (a *PDFDocumentAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	signatures := make([]validation.AdvancedSignature, 0)

	dssCertificateSource := spi.NewListCertificateSource()
	dssCRLSource := spi.NewListRevocationSource[revocation.CRL]()
	dssOCSPSource := spi.NewListRevocationSource[revocation.OCSP]()

	for _, pdfRevision := range a.Revisions() {
		switch typedRevision := pdfRevision.(type) {
		case *PdfDocDssRevision:
			dssCertificateSource.Add(typedRevision.CertificateSource())
			dssCRLSource.Add(typedRevision.CRLSource())
			dssOCSPSource.Add(typedRevision.OCSPSource())

		case *PdfSignatureRevision:
			padesSignature := NewPAdESSignature(typedRevision, a.documentRevisions)
			padesSignature.SetFilename(a.Document().Name())
			padesSignature.SetSigningCertificateSource(a.SigningCertificateSource())

			listCertificateSource := spi.NewListCertificateSource()
			listCertificateSource.AddAll(dssCertificateSource)
			padesSignature.SetDssCertificateSource(listCertificateSource)

			listCRLSource := spi.NewListRevocationSource[revocation.CRL]()
			listCRLSource.AddAllFrom(dssCRLSource)
			padesSignature.SetDssCRLSource(listCRLSource)

			listOCSPSource := spi.NewListRevocationSource[revocation.OCSP]()
			listOCSPSource.AddAllFrom(dssOCSPSource)
			padesSignature.SetDssOCSPSource(listOCSPSource)

			if a.CertificateVerifier() != nil {
				padesSignature.InitBaselineRequirementsChecker(a.CertificateVerifier())
			}
			a.ValidateSignaturePolicy(padesSignature)

			signatures = append(signatures, padesSignature)
		}
	}
	return utils.ReverseList(signatures)
}

// BuildDetachedTimestamps builds a list of detached TimestampTokens extracted from the document.
// Port of the protected buildDetachedTimestamps() override.
func (a *PDFDocumentAnalyzer) BuildDetachedTimestamps() []*validation.TimestampToken {
	timestamps := make([]*validation.TimestampToken, 0)
	var coveredReferences []*validation.TimestampedReference

	certificateSource := spi.NewListCertificateSource()
	crlSource := spi.NewListRevocationSource[revocation.CRL]()
	ocspSource := spi.NewListRevocationSource[revocation.OCSP]()

	for _, pdfRevision := range utils.ReverseList(a.Revisions()) {
		switch typedRevision := pdfRevision.(type) {
		case *PdfDocTimestampRevision:
			timestampToken := a.createPdfTimestampToken(typedRevision)
			certificateSource.Add(timestampToken.CertificateSource())
			crlSource.Add(timestampToken.CRLSource())
			ocspSource.Add(timestampToken.OCSPSource())

			tokenReferences := timestampToken.TimestampedReferences()
			spi.DSSUtilsEnrichCollection(&tokenReferences, coveredReferences)
			timestampToken.SetTimestampedReferences(tokenReferences)
			timestamps = append(timestamps, timestampToken.TimestampToken)

		case *PdfDocDssRevision:
			certificateSource.Add(typedRevision.CertificateSource())
			crlSource.Add(typedRevision.CRLSource())
			ocspSource.Add(typedRevision.OCSPSource())

		case *PdfSignatureRevision:
			return timestamps
		}

		// return refs for timestamps and DSS dictionaries
		pdfRevisionTimestampSource := NewPdfRevisionTimestampSource(pdfRevision, certificateSource, crlSource, ocspSource)
		coveredReferences = append(coveredReferences, pdfRevisionTimestampSource.IncorporatedReferences()...)
	}
	return timestamps
}

// createPdfTimestampToken ports the private createPdfTimestampToken(PdfDocTimestampRevision).
//
// Panics with the underlying error's message when the timestamp cannot be created, matching the
// Java DSSException wrapping every exception raised inside the try block.
func (a *PDFDocumentAnalyzer) createPdfTimestampToken(pdfDocTimestampRevision *PdfDocTimestampRevision) *PdfTimestampToken {
	timestampToken := pdfDocTimestampRevision.TimestampToken()
	timestampToken.SetFilename(a.Document().Name())

	timestampScopes := a.PAdESTimestampScopeFinder().FindTimestampScope(timestampToken.TimestampToken)
	timestampToken.SetTimestampScopes(timestampScopes)
	timestampToken.SetTimestampedReferences(append(timestampToken.TimestampedReferences(),
		a.timestampedReferences(timestampScopes)...))
	a.appendExternalEvidenceRecordsToTimestamp(timestampToken.TimestampToken)

	return timestampToken
}

// timestampedReferences returns a list of timestamped references from the given list of
// SignatureScopes.
//
// DEVIATION: reproduces analyzer.DefaultDocumentAnalyzer's unexported
// getTimestampedReferences(List), unreachable from this package (see the file header's
// cross-package note); AddReference is the promoted exported default (accept all), unless a
// future override changes it.
func (a *PDFDocumentAnalyzer) timestampedReferences(signatureScopes []mscope.SignatureScope) []*validation.TimestampedReference {
	timestampedReferences := make([]*validation.TimestampedReference, 0, len(signatureScopes))
	for _, signatureScope := range signatureScopes {
		if a.AddReference(signatureScope) {
			timestampedReferences = append(timestampedReferences,
				validation.NewTimestampedReference(signatureScope.DSSIDAsString(), enumerations.TimestampedObjectType_SIGNED_DATA))
		}
	}
	return timestampedReferences
}

// appendExternalEvidenceRecordsToTimestamp appends the detached evidence records covering the
// time-stamp.
//
// DEVIATION: reproduces analyzer.DefaultDocumentAnalyzer's unexported
// appendExternalEvidenceRecordsToTimestamp(TimestampToken), unreachable from this package (see
// the file header's cross-package note).
func (a *PDFDocumentAnalyzer) appendExternalEvidenceRecordsToTimestamp(timestampToken *validation.TimestampToken) {
	detachedTimestampSource := timestampsrc.NewDetachedTimestampSourceWithTimestamp(timestampToken)
	for _, evidenceRecord := range a.DetachedEvidenceRecords() {
		if a.IsTimestampCoveredByEvidenceRecord(timestampToken, evidenceRecord) {
			timestampToken.AddDetachedEvidenceRecord(evidenceRecord)
			_ = detachedTimestampSource.AddExternalEvidenceRecord(evidenceRecord)
		}
	}
}

// PAdESTimestampScopeFinder returns a PDF timestamp scope finder.
// Port of the protected getPAdESTimestampScopeFinder().
func (a *PDFDocumentAnalyzer) PAdESTimestampScopeFinder() *PAdESTimestampScopeFinder {
	return NewPAdESTimestampScopeFinder()
}

// DssDictionaries returns a list of found DSS Dictionaries across different revisions.
// Port of getDssDictionaries().
func (a *PDFDocumentAnalyzer) DssDictionaries() []PdfDssDict {
	dssDicts := make([]PdfDssDict, 0)
	for _, dssRevision := range a.DssRevisions() {
		dssDicts = append(dssDicts, dssRevision.DssDictionary())
	}
	return dssDicts
}

// DssRevisions returns a list of DSS revisions. Port of the protected getDssRevisions().
func (a *PDFDocumentAnalyzer) DssRevisions() []*PdfDocDssRevision {
	dssRevisions := make([]*PdfDocDssRevision, 0)
	for _, pdfRevision := range a.Revisions() {
		if dssRevision, ok := pdfRevision.(*PdfDocDssRevision); ok {
			dssRevisions = append(dssRevisions, dssRevision)
		}
	}
	return utils.ReverseList(dssRevisions)
}

// Revisions gets the list of PDF document revisions. Port of the protected getRevisions().
func (a *PDFDocumentAnalyzer) Revisions() []PdfRevision {
	if a.documentRevisions == nil {
		pdfSignatureService := a.pdfObjectFactory.NewPAdESSignatureService()
		a.documentRevisions = pdfSignatureService.GetRevisions(a.Document(), a.passwordProtection)
	}
	return a.documentRevisions
}

// InstantiateValidationDataContainer creates a new instance of ValidationDataContainer.
// Port of the protected instantiateValidationDataContainer() override.
func (a *PDFDocumentAnalyzer) InstantiateValidationDataContainer() *validation.ValidationDataContainer {
	return &NewPdfValidationDataContainer(a.DssRevisions()).ValidationDataContainer
}

// ValidationData extracts a validation data for the provided collection of signatures and/or
// timestamps, narrowed to the PDF-specific container. Port of the
// getValidationData(Collection, Collection) override, which in Java narrows the return type to
// PdfValidationDataContainer via a checked cast of the value the base class's virtually
// dispatched instantiateValidationDataContainer() (overridden just above) already built.
//
// The public analyzer.DocumentAnalyzer interface fixes GetValidationDataWithTimestamps's return
// type to *validation.ValidationDataContainer (no covariant return in Go), so this method is a
// distinct, PDF-only entry point - matching the already-landed SIGN chunk call sites
// (pades_level_baseline_lt.go, pades_extension_service.go) that call
// pdfDocumentAnalyzer.ValidationData(signatures, detachedTimestamps) and use the result's
// PdfValidationDataContainer-only IsEmpty() / signatureService.AddDssDictionary directly.
//
// Recovering the *PdfValidationDataContainer Java's cast reaches, from the
// *validation.ValidationDataContainer GetValidationDataWithTimestamps is contractually fixed to
// return, needs unsafe.Pointer: InstantiateValidationDataContainer above always builds the
// value by allocating a *PdfValidationDataContainer and returning the address of its embedded
// validation.ValidationDataContainer field, which - because that field is declared first in
// PdfValidationDataContainer's field list (pdf_validation_data_container.go) - the Go spec
// guarantees sits at the same address as the enclosing *PdfValidationDataContainer itself. This
// is the same "first embedded field, cast back" idiom Go's own composite literals rely on for
// interface satisfaction; it panics loudly instead of miscompiling if that struct layout
// invariant is ever broken (both types live in this file's own package, so the pairing cannot
// silently drift). Panics carry the same message Java's DSSException("At least one signature or
// a timestamp shall be provided...") would.
func (a *PDFDocumentAnalyzer) ValidationData(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken) *PdfValidationDataContainer {
	container, err := a.GetValidationDataWithTimestamps(signatures, detachedTimestamps)
	if err != nil {
		panic(err)
	}
	return (*PdfValidationDataContainer)(unsafe.Pointer(container))
}

// OriginalDocumentsForSignature returns the signed document(s) without their signature(s).
// Port of the getOriginalDocuments(AdvancedSignature) override.
func (a *PDFDocumentAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	padesSignature := advancedSignature.(*PAdESSignature)
	result := make([]model.DSSDocument, 0)
	originalPDF := PAdESUtilsGetOriginalPDF(padesSignature)
	if originalPDF != nil {
		isEmpty, err := spi.DSSUtilsIsEmpty(originalPDF)
		if err == nil && !isEmpty {
			result = append(result, originalPDF)
		}
	}
	return result
}
