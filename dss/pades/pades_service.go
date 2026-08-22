// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESService.java (DSS 6.5.RC1).
//
// Java extends AbstractSignatureService<PAdESSignatureParameters, PAdESTimestampParameters>; the
// Go port embeds document.AbstractSignatureService[*SignatureParameters, *cades.TimestampParameters]
// (NOT [...,*TimestampParameters] - an INTEGRATION CORRECTION, see below) and satisfies
// document.SignatureService[*SignatureParameters, *TimestampParameters] with
// the methods below (see the compile-time assertion at the end of the file); those two generic
// parameterizations are independent; nothing requires them to agree.
//
// INTEGRATION CORRECTION: originally embedded with TP=*TimestampParameters, matching
// SignatureService's TP. But pades_signature_parameters.go's header (see its "Context /
// timestamp-parameters storage" section) documents a deliberate design: SignatureParameters
// embeds cades.SignatureParameters by a single, unshadowed field, so its only promoted
// AbstractSignatureParameters is document.AbstractSignatureParameters[*cades.TimestampParameters]
// - there is no [*TimestampParameters]-instantiated one to take the address of. That made
// GetDataToSign/SignDocument's `s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)`
// calls (which need *document.AbstractSignatureParameters[TP] for the embedded base's own TP) a
// type error. Of AbstractSignatureService[SP, TP]'s two TP-typed methods, Timestamp is shadowed
// by Service's own directly-declared Timestamp(*TimestampParameters) below (so the
// embedded, unreachable default's TP is inconsequential) and AssertSigningCertificateValid is
// the only one actually invoked through the embedded base - so re-pointing the embed's TP at
// *cades.TimestampParameters (matching what parameters.AbstractSignatureParameters actually
// is) resolves the mismatch with no change to pades_signature_parameters.go's chosen structure.
//
// As in cades/cades_service.go, this is where the (T, error) of the layers below turns back into
// Java's propagating unchecked exception, i.e. into a panic whose value is the error itself, so a
// recovering caller can still inspect it with errors.As.
//
// ServiceLoaderPdfObjFactory has no Go counterpart: there is exactly one native backend, so the
// default factory is NewDefaultPdfObjFactory (internal/pdf/DESIGN.md §0.2).
//
// java.io.Serializable, the serialVersionUID and slf4j are dropped (PORTING.md).
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// Service is the PAdES implementation of SignatureService.
type Service struct {
	document.AbstractSignatureService[*SignatureParameters, *cades.TimestampParameters]

	// cmsForPAdESGenerationService builds the CMS signed data.
	cmsForPAdESGenerationService *ExternalCMSService

	// pdfObjFactory loads a relevant implementation for signature creation/extension.
	pdfObjFactory IPdfObjFactory
}

// NewService creates an instance of the Service. A certificate verifier must be
// provided: it gives information on the sources to be used in the validation process in the
// context of a signature. Port of PAdESService(CertificateVerifier).
func NewService(certificateVerifier validation.CertificateVerifier) *Service {
	// Upstream logs "+ PAdESService created".
	return &Service{
		AbstractSignatureService: document.NewAbstractSignatureService[*SignatureParameters,
			*cades.TimestampParameters](certificateVerifier),
		cmsForPAdESGenerationService: NewExternalCMSService(certificateVerifier),
		pdfObjFactory:                NewDefaultPdfObjFactory(),
	}
}

// SetPdfObjFactory sets the IPdfObjFactory, i.e. the implementation to be used. Cannot be nil.
// Port of #setPdfObjFactory.
func (s *Service) SetPdfObjFactory(pdfObjFactory IPdfObjFactory) {
	if pdfObjFactory == nil {
		panic("PdfObjFactory is null")
	}
	s.pdfObjFactory = pdfObjFactory
}

// SetTspSource defines the TSP source, propagating it to the CMS generation service.
// Port of the #setTspSource override.
func (s *Service) SetTspSource(tspSource validation.TSPSource) {
	s.AbstractSignatureService.SetTspSource(tspSource)
	s.cmsForPAdESGenerationService.SetTspSource(tspSource)
}

// extensionProfile ports the private #getExtensionProfile. It returns nil for PAdES-BASELINE-B,
// which upstream's null return models.
func (s *Service) extensionProfile(signatureLevel enumerations.SignatureLevel) document.SignatureExtension[*SignatureParameters] {
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}
	switch signatureLevel {
	case enumerations.SignatureLevelPAdESBaselineB:
		return nil
	case enumerations.SignatureLevelPAdESBaselineT:
		return NewLevelBaselineT(s.TspSource, s.CertificateVerifier, s.pdfObjFactory)
	case enumerations.SignatureLevelPAdESBaselineLT:
		return NewLevelBaselineLT(s.TspSource, s.CertificateVerifier, s.pdfObjFactory)
	case enumerations.SignatureLevelPAdESBaselineLTA:
		return NewLevelBaselineLTA(s.TspSource, s.CertificateVerifier, s.pdfObjFactory)
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", signatureLevel))
	}
}

// GetContentTimestamp requests a content time-stamp computed on the PDF revision to be signed.
// Port of #getContentTimestamp.
func (s *Service) GetContentTimestamp(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) *validation.TimestampToken {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	UtilsAssertPdfDocument(toSignDocument)
	padesServiceAssertContentTimestampParametersValid(parameters)

	pdfSignatureService := s.ContentTimestampService()
	messageDigest := pdfSignatureService.MessageDigest(toSignDocument, parameters)
	timeStampResponse, err := s.TspSource.TimeStampResponse(parameters.DigestAlgorithm(), messageDigest.Value())
	if err != nil {
		panic(err)
	}
	timestampToken, err := validation.NewTimestampToken(timeStampResponse.Bytes(),
		enumerations.TimestampTypeContentTimestamp)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot obtain the content timestamp", err))
	}
	return timestampToken
}

// padesServiceAssertContentTimestampParametersValid ports the private
// #assertContentTimestampParametersValid.
func padesServiceAssertContentTimestampParametersValid(parameters *SignatureParameters) {
	if parameters.DigestAlgorithm() != parameters.ContentTimestampParameters().DigestAlgorithm() {
		panic("DigestAlgorithm for content timestamp creation shall be " +
			"the same as the one defined in PAdESSignatureParameters!")
	}
}

// PreviewPageWithVisualSignature returns a page preview with the visual signature.
// Port of #previewPageWithVisualSignature.
func (s *Service) PreviewPageWithVisualSignature(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	UtilsAssertPdfDocument(toSignDocument)

	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.PreviewPageWithVisualSignature(toSignDocument, parameters)
}

// PreviewSignatureField returns a preview of the signature field.
// Port of #previewSignatureField.
func (s *Service) PreviewSignatureField(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	UtilsAssertPdfDocument(toSignDocument)

	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.PreviewSignatureField(toSignDocument, parameters)
}

// GetDataToSign retrieves the data to be signed, i.e. the signed attributes of the CMS built over
// the digest of the PDF revision's ByteRange. Port of #getDataToSign.
func (s *Service) GetDataToSign(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) *model.ToBeSigned {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}

	UtilsAssertPdfDocument(toSignDocument)
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	messageDigest := s.ComputeDocumentDigest(toSignDocument, parameters)
	toBeSigned, err := s.cmsForPAdESGenerationService.BuildToBeSignedData(messageDigest, parameters)
	if err != nil {
		panic(err)
	}
	return toBeSigned
}

// ComputeDocumentDigest computes the digest of the document to be signed.
// Port of the protected #computeDocumentDigest.
func (s *Service) ComputeDocumentDigest(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSMessageDigest {
	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.MessageDigest(toSignDocument, parameters)
}

// SignDocument signs the document with the provided signature value. Port of #signDocument.
func (s *Service) SignDocument(toSignDocument model.DSSDocument, parameters *SignatureParameters,
	signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}

	UtilsAssertPdfDocument(toSignDocument)
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	signatureValue, err := s.EnsureSignatureValue(parameters.SignatureAlgorithm(), signatureValue)
	if err != nil {
		panic(err)
	}

	signatureLevel := parameters.SignatureLevel()
	cmsSignedData := s.GenerateCMSSignedData(toSignDocument, parameters, signatureValue)

	pdfSignatureService := s.PAdESSignatureService()
	signature := pdfSignatureService.Sign(toSignDocument, cmsSignedData, parameters)

	extension := s.extensionProfile(signatureLevel)
	if signatureLevel != enumerations.SignatureLevelPAdESBaselineB &&
		signatureLevel != enumerations.SignatureLevelPAdESBaselineT && extension != nil {
		signature = extension.ExtendSignatures(signature, parameters)
	}

	parameters.Reinit()
	name, err := s.GetFinalFileNameWithLevel(toSignDocument, enumerations.SigningOperationSign,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	signature.SetName(name)
	return signature
}

// GenerateCMSSignedData generates the DER-encoded CMS signed data enveloped by the PDF signature
// dictionary. Port of the protected #generateCMSSignedData.
func (s *Service) GenerateCMSSignedData(toSignDocument model.DSSDocument,
	parameters *SignatureParameters, signatureValue *model.SignatureValue) []byte {
	signatureAlgorithm := parameters.SignatureAlgorithm()
	signatureLevel := parameters.SignatureLevel()
	if signatureAlgorithm == "" {
		panic("SignatureAlgorithm cannot be null!")
	}
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}

	messageDigest := s.ComputeDocumentDigest(toSignDocument, parameters)
	signedCMS, err := s.cmsForPAdESGenerationService.BuildCMS(messageDigest, parameters, signatureValue)
	if err != nil {
		panic(err)
	}
	return signedCMS.DEREncoded()
}

// ExtendDocument extends the signatures of the given PDF document. Port of #extendDocument.
func (s *Service) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument is not defined!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}

	UtilsAssertPdfDocument(toExtendDocument)
	padesServiceAssertExtensionParametersValid(parameters)

	extension := s.extensionProfile(parameters.SignatureLevel())
	if extension != nil {
		extended := extension.ExtendSignatures(toExtendDocument, parameters)
		name, err := s.GetFinalFileNameWithLevel(toExtendDocument, enumerations.SigningOperationExtend,
			parameters.SignatureLevel())
		if err != nil {
			panic(err)
		}
		extended.SetName(name)
		return extended
	}
	return toExtendDocument
}

// padesServiceAssertExtensionParametersValid ports the private #assertExtensionParametersValid.
func padesServiceAssertExtensionParametersValid(parameters *SignatureParameters) {
	if enumerations.SignatureLevelPAdESBaselineB == parameters.SignatureLevel() {
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", parameters.SignatureLevel()))
	}
}

// GetAvailableSignatureFields returns the not signed signature fields of the document.
// Port of getAvailableSignatureFields(DSSDocument).
func (s *Service) GetAvailableSignatureFields(document model.DSSDocument) []string {
	return s.GetAvailableSignatureFieldsWithPassword(document, nil)
}

// GetAvailableSignatureFieldsWithPassword returns the not signed signature fields of an
// encrypted document. Port of getAvailableSignatureFields(DSSDocument, char[]).
func (s *Service) GetAvailableSignatureFieldsWithPassword(document model.DSSDocument,
	passwordProtection []byte) []string {
	if document == nil {
		panic("DSSDocument is not defined!")
	}
	UtilsAssertPdfDocument(document)

	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.GetAvailableSignatureFields(document, passwordProtection)
}

// AddNewSignatureField adds a new signature field to an existing PDF document.
// Port of addNewSignatureField(DSSDocument, SignatureFieldParameters).
func (s *Service) AddNewSignatureField(document model.DSSDocument,
	parameters *SignatureFieldParameters) model.DSSDocument {
	return s.AddNewSignatureFieldWithPassword(document, parameters, nil)
}

// AddNewSignatureFieldWithPassword adds a new signature field to an existing encrypted PDF
// document. Port of addNewSignatureField(DSSDocument, SignatureFieldParameters, char[]).
func (s *Service) AddNewSignatureFieldWithPassword(document model.DSSDocument,
	parameters *SignatureFieldParameters, passwordProtection []byte) model.DSSDocument {
	if document == nil {
		panic("DSSDocument is not defined!")
	}
	if parameters == nil {
		panic("SignatureFieldParameters cannot be null!")
	}
	UtilsAssertPdfDocument(document)

	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.AddNewSignatureField(document, parameters, passwordProtection)
}

// Timestamp adds a document timestamp to an unsigned document, incorporating the validation data
// first. Port of #timestamp.
func (s *Service) Timestamp(toTimestampDocument model.DSSDocument,
	parameters *TimestampParameters) model.DSSDocument {
	if toTimestampDocument == nil {
		panic("Document to be timestamped is not defined!")
	}
	if parameters == nil {
		panic("PAdESTimestampParameters cannot be null!")
	}
	UtilsAssertPdfDocument(toTimestampDocument)

	extensionService := NewExtensionServiceWithFactory(s.CertificateVerifier, s.pdfObjFactory)
	extendedDocument := extensionService.IncorporateValidationDataWithPassword(toTimestampDocument,
		parameters.PasswordProtection())

	timestampService := NewTimestampServiceWithPDFService(s.TspSource, s.SignatureTimestampService())
	timestampedDocument := timestampService.TimestampDocument(extendedDocument, parameters)
	name, err := s.GetFinalFileName(toTimestampDocument, enumerations.SigningOperationTimestamp)
	if err != nil {
		panic(err)
	}
	timestampedDocument.SetName(name)
	return timestampedDocument
}

// PAdESSignatureService returns a new PDFSignatureService for a signature creation.
// Port of the protected #getPAdESSignatureService.
func (s *Service) PAdESSignatureService() PDFSignatureService {
	return s.pdfObjFactory.NewPAdESSignatureService()
}

// ContentTimestampService returns a new PDFSignatureService for a content timestamp creation.
// Port of the protected #getContentTimestampService.
func (s *Service) ContentTimestampService() PDFSignatureService {
	return s.pdfObjFactory.NewContentTimestampService()
}

// SignatureTimestampService returns a new PDFSignatureService for a timestamp creation.
// Port of the protected #getSignatureTimestampService.
func (s *Service) SignatureTimestampService() PDFSignatureService {
	return s.pdfObjFactory.NewSignatureTimestampService()
}

// Compile-time interface assertion, standing in for Java's
// "extends AbstractSignatureService<SignatureParameters, TimestampParameters>".
var _ document.SignatureService[*SignatureParameters, *TimestampParameters] = (*Service)(nil)
