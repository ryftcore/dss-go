// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/AbstractASiCSignatureService.java (DSS 6.5.RC1).
package asic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractASiCSignatureServiceOverrides declares the operations AbstractASiCSignatureService
// calls back into virtually. Per S7_BRIEF.md's virtual-dispatch warning: the single-document
// convenience wrappers (GetContentTimestamp, GetDataToSign, SignDocument, Timestamp) each
// delegate to their multi-document counterpart, which is what the concrete, format-specific
// ASiC service (CADSIGN/XADSIGN chunks) implements to satisfy
// document.MultipleDocumentsSignatureService; GetArchiveExtractor is the format-specific
// container extractor lookup. Every concrete service must call
// InitAbstractASiCSignatureService with itself before use.
type AbstractASiCSignatureServiceOverrides[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters] interface {
	// GetContentTimestampMultiple creates a content-timestamp attribute for multiple
	// documents. Port of document.MultipleDocumentsSignatureService's getContentTimestamp(List, SP).
	GetContentTimestampMultiple(toSignDocuments []model.DSSDocument, parameters SP) *validation.TimestampToken

	// GetDataToSignMultiple retrieves the bytes of the data that need to be signed for
	// multiple documents. Port of getDataToSign(List, SP).
	GetDataToSignMultiple(toSignDocuments []model.DSSDocument, parameters SP) *model.ToBeSigned

	// SignDocumentMultiple signs multiple documents with the provided signatureValue. Port of
	// signDocument(List, SP, SignatureValue).
	SignDocumentMultiple(toSignDocuments []model.DSSDocument, parameters SP, signatureValue *model.SignatureValue) model.DSSDocument

	// TimestampMultiple timestamps multiple documents with the provided parameters. Port of
	// timestamp(List, TP).
	TimestampMultiple(toTimestampDocuments []model.DSSDocument, parameters TP) model.DSSDocument

	// GetArchiveExtractor returns a relevant ASiC container extractor for the given format.
	// Port of the protected abstract getArchiveExtractor(DSSDocument).
	GetArchiveExtractor(archive model.DSSDocument) *DefaultASiCContainerExtractor

	// AddContainerEvidenceRecordMultiple creates a new ASiC container with the
	// evidenceRecordDocument applied to documents. Port of the public abstract
	// addContainerEvidenceRecord(List, DSSDocument, ASiCContainerEvidenceRecordParameters).
	AddContainerEvidenceRecordMultiple(documents []model.DSSDocument, evidenceRecordDocument model.DSSDocument, parameters *ASiCContainerEvidenceRecordParameters) model.DSSDocument
}

// AbstractASiCSignatureService contains the main methods for ASiC signature creation/extension,
// generic over SP (signature parameters), TP (timestamp parameters), CSP (counter-signature
// parameters) and ERP (evidence record incorporation parameters).
//
// Java's class extends AbstractSignatureService<SP, TP> and implements
// MultipleDocumentsSignatureService<SP, TP>, CounterSignatureService<CSP> and
// EvidenceRecordIncorporationService<ERP> (the last a local forward declaration, see
// evidence_record_incorporation_service.go); like AbstractSignatureService itself (see its Go
// port's doc comment), this type does not assert those interface conformances - it is meant to
// be embedded by a concrete, format-specific service that supplies the remaining methods.
type AbstractASiCSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters, CSP model.SerializableCounterSignatureParameters, ERP model.SerializableEvidenceRecordIncorporationParameters] struct {
	document.AbstractSignatureService[SP, TP]

	// overrides points back at the concrete service; see InitAbstractASiCSignatureService.
	overrides AbstractASiCSignatureServiceOverrides[SP, TP]
}

// NewAbstractASiCSignatureService is the default constructor. Ports
// AbstractASiCSignatureService(CertificateVerifier). The subclass constructor must follow it
// with InitAbstractASiCSignatureService.
func NewAbstractASiCSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters, CSP model.SerializableCounterSignatureParameters, ERP model.SerializableEvidenceRecordIncorporationParameters](certificateVerifier validation.CertificateVerifier) AbstractASiCSignatureService[SP, TP, CSP, ERP] {
	return AbstractASiCSignatureService[SP, TP, CSP, ERP]{
		AbstractSignatureService: document.NewAbstractSignatureService[SP, TP](certificateVerifier),
	}
}

// InitAbstractASiCSignatureService registers the concrete service with its base so the base
// can dispatch the multi-document operations and GetArchiveExtractor. Every concrete service
// constructor must call this once.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) InitAbstractASiCSignatureService(overrides AbstractASiCSignatureServiceOverrides[SP, TP]) {
	s.overrides = overrides
}

func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) requireOverrides() AbstractASiCSignatureServiceOverrides[SP, TP] {
	if s.overrides == nil {
		panic("AbstractASiCSignatureService was not initialised: the concrete service must call InitAbstractASiCSignatureService in its constructor")
	}
	return s.overrides
}

// GetContentTimestamp ports the @Override getContentTimestamp(DSSDocument, SP).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) GetContentTimestamp(toSignDocument model.DSSDocument, parameters SP) *validation.TimestampToken {
	return s.requireOverrides().GetContentTimestampMultiple([]model.DSSDocument{toSignDocument}, parameters)
}

// GetDataToSign ports the @Override getDataToSign(DSSDocument, SP).
//
// Panics with Java's message when toSignDocument is nil (Objects.requireNonNull).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) GetDataToSign(toSignDocument model.DSSDocument, parameters SP) *model.ToBeSigned {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	return s.requireOverrides().GetDataToSignMultiple([]model.DSSDocument{toSignDocument}, parameters)
}

// SignDocument ports the @Override signDocument(DSSDocument, SP, SignatureValue).
//
// Panics with Java's message when toSignDocument is nil.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) SignDocument(toSignDocument model.DSSDocument, parameters SP, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	return s.requireOverrides().SignDocumentMultiple([]model.DSSDocument{toSignDocument}, parameters, signatureValue)
}

// SignatureTimestamp ports the @Override timestamp(DSSDocument, TP). Named SignatureTimestamp
// (not Timestamp) to avoid colliding with the embedded
// document.AbstractSignatureService[SP,TP].Timestamp(DSSDocument, TP) method, which this
// override shadows for callers going through *AbstractASiCSignatureService.
//
// Panics with Java's message when toTimestampDocument is nil.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) SignatureTimestamp(toTimestampDocument model.DSSDocument, parameters TP) model.DSSDocument {
	if toTimestampDocument == nil {
		panic("toTimestampDocument cannot be null!")
	}
	return s.requireOverrides().TimestampMultiple([]model.DSSDocument{toTimestampDocument}, parameters)
}

// AddContainerEvidenceRecordSingle creates a new ASiC container with the
// evidenceRecordDocument applied to the document. If the provided original document is an
// existing ASiC container, then the evidenceRecordDocument will be evaluated against the
// container files and places within the container. Ports
// addContainerEvidenceRecord(DSSDocument, DSSDocument, ASiCContainerEvidenceRecordParameters).
//
// Named AddContainerEvidenceRecordSingle (not AddContainerEvidenceRecord) since Go cannot
// overload by parameter type against the abstract multi-document
// AddContainerEvidenceRecord(documents, evidenceRecordDocument, parameters) every concrete ASiC
// service must implement directly (there is no base implementation to embed - Java's abstract
// method has no body).
//
// Panics with Java's message when document is nil.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) AddContainerEvidenceRecordSingle(doc model.DSSDocument, evidenceRecordDocument model.DSSDocument, parameters *ASiCContainerEvidenceRecordParameters) model.DSSDocument {
	if doc == nil {
		panic("Document cannot be null!")
	}
	return s.requireOverrides().AddContainerEvidenceRecordMultiple([]model.DSSDocument{doc}, evidenceRecordDocument, parameters)
}

// ExtractCurrentArchive extracts the content (documents) of the ASiC container. Ports the
// protected extractCurrentArchive(DSSDocument).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) ExtractCurrentArchive(archive model.DSSDocument) *ASiCContent {
	extractor := s.requireOverrides().GetArchiveExtractor(archive)
	content, err := extractor.Extract()
	if err != nil {
		panic(err)
	}
	return content
}

// BuildASiCContainer creates a ZIP-Archive by copying the provided documents to the new
// container using the current time as ZIP creation time. Ports the protected
// buildASiCContainer(ASiCContent).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) BuildASiCContainer(asicContent *ASiCContent) model.DSSDocument {
	return s.BuildASiCContainerAt(asicContent, time.Now())
}

// BuildASiCContainerAt creates a ZIP-Archive by copying the provided documents to the new
// container. Ports the protected buildASiCContainer(ASiCContent, Date).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) BuildASiCContainerAt(asicContent *ASiCContent, creationTime time.Time) model.DSSDocument {
	zipArchive, err := ZipUtilsInstance().CreateZipArchiveAt(asicContent, creationTime)
	if err != nil {
		panic(err)
	}
	mimeType, err := ASiCUtilsMimeTypeFromDocument(asicContent.MimeTypeDocument())
	if err != nil {
		panic(err)
	}
	zipArchive.SetMimeType(mimeType)
	return zipArchive
}

// AssertSignaturePossible verifies whether the signature creation is possible with the
// provided documents. Ports the protected assertSignaturePossible(List).
//
// Panics with Java's IllegalArgumentException messages on invalid input.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) AssertSignaturePossible(documentsToSign []model.DSSDocument) {
	if utils.IsCollectionEmpty(documentsToSign) {
		panic("List of documents to sign cannot be empty!")
	}
	for _, doc := range documentsToSign {
		if _, ok := doc.(*model.DigestDocument); ok {
			panic("ASiC container creation is not possible with DigestDocument!")
		}
	}
}

// AssertCounterSignatureParametersValid verifies a validity of counter signature parameters.
// Ports the protected assertCounterSignatureParametersValid(CSP).
//
// Panics with Java's message when the signature id to counter sign is not defined.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) AssertCounterSignatureParametersValid(parameters CSP) {
	if parameters.SignatureIdToCounterSign() == "" {
		panic("The Id of a signature to be counter signed shall be defined! " +
			"Please use SerializableCounterSignatureParameters.setSignatureIdToCounterSign(signatureId) method.")
	}
}

// AssertAddSignaturePolicyStorePossible verifies if incorporation of a SignaturePolicyStore is
// possible. Ports the protected assertAddSignaturePolicyStorePossible(ASiCContent).
//
// Panics with Java's UnsupportedOperationException message when no matching signature
// documents are found.
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) AssertAddSignaturePolicyStorePossible(asicContent *ASiCContent) {
	if utils.IsCollectionEmpty(asicContent.SignatureDocuments()) {
		panic("Signature documents of the expected format are not found in the provided ASiC Container! " +
			"Add a SignaturePolicyStore is not possible!")
	}
}

// GetFinalArchiveName generates and returns a final name for the archive to create. Ports the
// protected getFinalArchiveName(DSSDocument, SigningOperation, MimeType).
func (s *AbstractASiCSignatureService[SP, TP, CSP, ERP]) GetFinalArchiveName(originalFile model.DSSDocument, operation enumerations.SigningOperation, containerMimeType enumerations.MimeType) (string, error) {
	return s.GetFinalDocumentNameWithMimeType(originalFile, operation, "", containerMimeType)
}
