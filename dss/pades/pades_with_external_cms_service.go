// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESWithExternalCMSService.java (DSS 6.5.RC1).
//
// BouncyCastle's CMSSignedData is cms.CMS (PORTING.md), so DSSUtils.toCMSSignedData becomes
// cms.UtilsParseToCMS and DSSASN1Utils.getDEREncoded(CMSSignedData) becomes CMS.DEREncoded.
// java.io.Serializable, the serialVersionUID and slf4j are dropped.
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// WithExternalCMSService creates a PAdES signature using an external CMS provider.
//
// To create a signature with this service, follow the algorithm:
//  1. create a message-digest computed on the PDF ByteRange:
//     messageDigest := service.GetMessageDigest(toSignDocument, parameters);
//  2. create the CMS signature signing the message-digest (e.g. using a remote-signing solution);
//  3. OPTIONAL: verify the validity of the obtained CMS signature with
//     IsValidCMSSignedData (cryptographic validity) and IsValidPAdESBaselineCMSSignedData
//     (CMS applicability rules for a PAdES signature creation);
//  4. create the PAdES signature by incorporating the obtained CMS signature into the PDF:
//     signedDocument := service.SignDocument(toSignDocument, parameters, cmsDocument).
//
// NOTES:
//   - unlike Service, the SignatureParameters given to this service do not need a
//     signing certificate and certificate chain when using external signing;
//   - signature extension to -T level with this service never leads to a signature-timestamp
//     inside the CMS signed data; it always creates a new revision with a document timestamp;
//   - a content timestamp is not supported by this service.
type WithExternalCMSService struct {
	// certificateVerifier is used for a certificate chain validation.
	certificateVerifier validation.CertificateVerifier

	// tspSource is used for timestamp requests.
	tspSource validation.TSPSource

	// pdfObjFactory loads a relevant implementation for signature creation/extension.
	pdfObjFactory IPdfObjFactory
}

// NewPAdESWithExternalCMSService is the default constructor.
// Port of the no-arg PAdESWithExternalCMSService() constructor.
func NewPAdESWithExternalCMSService() *WithExternalCMSService {
	return &WithExternalCMSService{pdfObjFactory: NewDefaultPdfObjFactory()}
}

// SetCertificateVerifier defines the CertificateVerifier used for a signature extension and on
// the CMS creation method. It is not required for B-level remote-signing solutions.
// Port of #setCertificateVerifier.
func (s *WithExternalCMSService) SetCertificateVerifier(certificateVerifier validation.CertificateVerifier) {
	s.certificateVerifier = certificateVerifier
}

// SetTspSource defines the TSP (timestamp provider) source. Port of #setTspSource.
func (s *WithExternalCMSService) SetTspSource(tspSource validation.TSPSource) {
	s.tspSource = tspSource
}

// SetPdfObjFactory sets the IPdfObjFactory, i.e. the implementation to be used. Cannot be nil.
// Port of #setPdfObjFactory.
func (s *WithExternalCMSService) SetPdfObjFactory(pdfObjFactory IPdfObjFactory) {
	if pdfObjFactory == nil {
		panic("PdfObjFactory is null")
	}
	s.pdfObjFactory = pdfObjFactory
}

// GetMessageDigest computes the message-digest of the signature ByteRange to be used for the CMS
// signed data creation. Port of #getMessageDigest.
func (s *WithExternalCMSService) GetMessageDigest(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSMessageDigest {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	UtilsAssertPdfDocument(toSignDocument)

	pdfSignatureService := s.PAdESSignatureService()
	return pdfSignatureService.MessageDigest(toSignDocument, parameters)
}

// SignDocument embeds the provided external cmsDocument into toSignDocument within a new
// signature revision. Port of #signDocument.
func (s *WithExternalCMSService) SignDocument(toSignDocument model.DSSDocument,
	parameters *SignatureParameters, cmsDocument model.DSSDocument) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if parameters.SignatureLevel() == "" {
		panic("SignatureLevel shall be defined within parameters!")
	}
	if cmsDocument == nil {
		panic("CMSDocument cannot be null!")
	}
	UtilsAssertPdfDocument(toSignDocument)
	padesWithExternalCMSServiceAssertNotDigestDocument(cmsDocument)

	cmsSignedData := padesWithExternalCMSServiceToCMSSignedData(cmsDocument)
	derEncodedCMS := cmsSignedData.DEREncoded()

	pdfSignatureService := s.PAdESSignatureService()
	signatureDocument := pdfSignatureService.Sign(toSignDocument, derEncodedCMS, parameters)

	if enumerations.SignatureLevelPAdESBaselineB != parameters.SignatureLevel() &&
		padesWithExternalCMSServiceIsExtensionRequired(cmsSignedData, parameters) {
		parameters.GetContext().SetDetachedContents([]model.DSSDocument{toSignDocument})
		padesService := s.PAdESService()
		signatureDocument = padesService.ExtendDocument(signatureDocument, parameters)
	}

	signatureDocument.SetName(s.FinalDocumentName(toSignDocument, parameters.SignatureLevel()))
	parameters.Reinit()
	return signatureDocument
}

// padesWithExternalCMSServiceToCMSSignedData ports the private #toCMSSignedData.
func padesWithExternalCMSServiceToCMSSignedData(document model.DSSDocument) *cms.CMS {
	parsed, err := cms.UtilsParseToCMS(document)
	if err != nil {
		panic(exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("A CMS file is expected : %s", err.Error()), err))
	}
	return parsed
}

// PAdESSignatureService returns a new PDFSignatureService for a signature creation.
// Port of the protected #getPAdESSignatureService.
func (s *WithExternalCMSService) PAdESSignatureService() PDFSignatureService {
	return s.pdfObjFactory.NewPAdESSignatureService()
}

// Service creates an instance of a Service to be used for a signature extension.
// Port of the protected #getPAdESService.
func (s *WithExternalCMSService) PAdESService() *Service {
	if s.certificateVerifier == nil {
		panic("CertificateVerifier shall be provided for PAdES extension!")
	}
	if s.tspSource == nil {
		panic("TSPSource shall be provided for PAdES extension!")
	}

	padesService := NewPAdESService(s.certificateVerifier)
	padesService.SetTspSource(s.tspSource)
	padesService.SetPdfObjFactory(s.pdfObjFactory)
	return padesService
}

// FinalDocumentName generates and returns a final name for the document to be created.
// Port of the protected #getFinalDocumentName.
func (s *WithExternalCMSService) FinalDocumentName(originalFile model.DSSDocument,
	level enumerations.SignatureLevel) string {
	name, err := validation.NewFileNameBuilder().SetOriginalFilename(originalFile.Name()).
		SetSigningOperation(enumerations.SigningOperationSign).SetSignatureLevel(level).
		SetSignaturePackaging(enumerations.SignaturePackagingEnveloped).
		SetMimeType(enumerations.MimeTypeEnumPDF).Build()
	if err != nil {
		panic(err)
	}
	return name
}

// padesWithExternalCMSServiceAssertNotDigestDocument ports the private #assertNotDigestDocument.
func padesWithExternalCMSServiceAssertNotDigestDocument(document model.DSSDocument) {
	if _, isDigestDocument := document.(*model.DigestDocument); isDigestDocument {
		panic("DigestDocument is not allowed for current operation!")
	}
}

// padesWithExternalCMSServiceIsExtensionRequired ports the private #isExtensionRequired.
func padesWithExternalCMSServiceIsExtensionRequired(cmsSignedData *cms.CMS,
	parameters *SignatureParameters) bool {
	if enumerations.SignatureLevelPAdESBaselineT == parameters.SignatureLevel() {
		// only first SignerInformation is considered.
		signerInformation := spi.DSSASN1UtilsFirstSignerInformation(cmsSignedData.SignerInfos())
		unsignedAttributes := cades.UnsignedAttributesBuild(signerInformation)
		for _, attribute := range unsignedAttributes.Attributes() {
			if cades.OIDIdAaSignatureTimeStampToken.Equal(attribute.ASN1Oid()) {
				// Upstream logs "The CMS signature already contains a signature-time-stamp
				// attribute! The extension to '%s' level is skipped."
				return false
			}
		}
	}
	return true
}

// IsValidCMSSignedData verifies whether the given CMS document is cryptographically valid
// against the message-digest computed on the PDF signature ByteRange.
// Port of #isValidCMSSignedData.
func (s *WithExternalCMSService) IsValidCMSSignedData(messageDigest model.DSSMessageDigest,
	cmsDocument model.DSSDocument) bool {
	if messageDigest.Value() == nil {
		panic("messageDigest shall be provided!")
	}
	if cmsDocument == nil {
		panic("cmsDocument shall be provided!")
	}

	parsedCMS, err := cms.UtilsParseToCMS(cmsDocument)
	if err != nil {
		// Upstream logs "Unable to decode the provided CMS document : {}".
		return false
	}

	signerInfos := parsedCMS.SignerInfos()
	if len(signerInfos) != 1 {
		// Upstream logs "CMSSignedData shall contain one and only one SignerInformation for
		// signature signing process!".
		return false
	}

	cadesSignature := padesWithExternalCMSServiceToCAdESSignature(parsedCMS, messageDigest)
	scv := cadesSignature.SignatureCryptographicVerification()
	if !scv.IsSignatureValid() {
		// Upstream logs "CMSSignedData signature is not valid!".
		return false
	}
	return true
}

// IsValidPAdESBaselineCMSSignedData verifies whether the given CMS signature is compliant with
// the PAdES format. Port of #isValidPAdESBaselineCMSSignedData.
func (s *WithExternalCMSService) IsValidPAdESBaselineCMSSignedData(messageDigest model.DSSMessageDigest,
	cmsDocument model.DSSDocument) bool {
	if messageDigest.Value() == nil {
		panic("messageDigest shall be provided!")
	}
	if cmsDocument == nil {
		panic("cmsDocument shall be provided!")
	}

	parsedCMS, err := cms.UtilsParseToCMS(cmsDocument)
	if err != nil {
		// Upstream logs "Unable to decode the provided CMS document : {}".
		return false
	}

	cadesSignature := padesWithExternalCMSServiceToCAdESSignature(parsedCMS, messageDigest)
	cmsRequirementsChecker := NewCMSForPAdESBaselineRequirementsChecker(cadesSignature)
	return cmsRequirementsChecker.IsValidForPAdESBaselineBProfile()
}

// padesWithExternalCMSServiceToCAdESSignature ports the private #toCAdESSignature.
func padesWithExternalCMSServiceToCAdESSignature(parsedCMS *cms.CMS,
	messageDigest model.DSSMessageDigest) *cades.Signature {
	signature := cades.NewCAdESSignature(parsedCMS,
		spi.DSSASN1UtilsFirstSignerInformation(parsedCMS.SignerInfos()))
	signature.SetDetachedContents([]model.DSSDocument{spi.DSSUtilsToDigestDocument(messageDigest.Digest)})
	return signature
}
