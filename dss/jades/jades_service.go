// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESService.java (DSS 6.5.RC1).
//
// Java extends AbstractSignatureService<JAdESSignatureParameters, JAdESTimestampParameters> and
// implements MultipleDocumentsSignatureService and CounterSignatureService; the Go port embeds
// document.AbstractSignatureService[*JAdESSignatureParameters, *JAdESTimestampParameters], the
// same way xades.XAdESService and cades.CAdESService do.
//
// # Overloads, and the one interface Go cannot let this type carry
//
// Java overloads getContentTimestamp, getDataToSign, signDocument and timestamp on DSSDocument vs
// List<DSSDocument>: DocumentSignatureService declares the first shape,
// MultipleDocumentsSignatureService the second, and one Java class satisfies both. Go has no
// overloading, so a single type cannot carry both method sets. This port keeps the plain names
// for the single-document shape - so *JAdESService satisfies document.DocumentSignatureService -
// and gives the multi-document shape distinct names, exactly as xades_service.go does:
//
//	getContentTimestamp(List, SP)          -> GetContentTimestampForDocuments
//	getDataToSign(List, SP)                -> GetDataToSignForDocuments
//	signDocument(List, SP, SignatureValue) -> SignDocuments
//	timestamp(List, TP)                    -> TimestampDocuments
//
// MultipleDocumentsService() returns a thin adapter that does satisfy
// document.MultipleDocumentsSignatureService by forwarding to those methods, so the Java contract
// is not lost; ASiC-JAdES (phase 7) is the caller that will need it.
//
// # Errors
//
// The document service interfaces return bare values, so - exactly as in xades_service.go - this
// is where the (T, error) of the layers below turns back into Java's propagating exception, i.e.
// into a panic whose value is the error itself.
//
// java.io.Serializable, the serialVersionUID and slf4j are dropped.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESSignatureExtender is the static type of getExtensionProfile's local variable upstream:
// Java assigns a JAdESLevelBaselineT, -LT or -LTA to a JAdESLevelBaselineExtension and then calls
// setTspSource on it, which the Java interface does not declare (the three concrete classes do).
// Go needs one interface value for the switch, so this adds SetTspSource on top of the ported
// JAdESLevelBaselineExtension; all three levels satisfy it through their embedded
// JAdESLevelBaselineT. Same technique, and same reason, as xades.XAdESSignatureExtender.
type JAdESSignatureExtender interface {
	JAdESLevelBaselineExtension

	// SetTspSource sets the TSP source used when extending. Port of
	// JAdESLevelBaselineT#setTspSource.
	SetTspSource(tspSource validation.TSPSource)
}

// JAdESService contains the methods for JAdES signature creation/extension.
type JAdESService struct {
	document.AbstractSignatureService[*JAdESSignatureParameters, *JAdESTimestampParameters]
}

// NewJAdESService creates an instance of the JAdESService. A certificate verifier must be
// provided; it supplies information on the sources to be used in the validation process in the
// context of a signature. Port of JAdESService(CertificateVerifier).
func NewJAdESService(certificateVerifier validation.CertificateVerifier) *JAdESService {
	// Upstream logs "+ JAdESService created".
	return &JAdESService{
		AbstractSignatureService: document.NewAbstractSignatureService[*JAdESSignatureParameters,
			*JAdESTimestampParameters](certificateVerifier),
	}
}

// GetContentTimestamp creates a content-timestamp covering the document to be signed.
// Port of the #getContentTimestamp(DSSDocument, JAdESSignatureParameters) overload.
func (s *JAdESService) GetContentTimestamp(toSignDocument model.DSSDocument,
	parameters *JAdESSignatureParameters) *validation.TimestampToken {
	return s.GetContentTimestampForDocuments([]model.DSSDocument{toSignDocument}, parameters)
}

// GetContentTimestampForDocuments creates a TimestampToken for a detached JAdES (with a 'sigD'
// parameter).
//
// NOTE: the toSignDocuments must be present in the same order they will be passed to the
// signature computation process.
//
// Port of the #getContentTimestamp(List<DSSDocument>, JAdESSignatureParameters) overload.
func (s *JAdESService) GetContentTimestampForDocuments(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters) *validation.TimestampToken {
	if s.TspSource == nil {
		panic("A TSPSource is required!")
	}
	if err := jadesServiceAssertContentTimestampCreationPossible(toSignDocuments); err != nil {
		panic(err)
	}

	var messageImprint []byte
	var err error
	if enumerations.SigDMechanismHTTPHeaders == parameters.SigDMechanism() {
		httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(toSignDocuments, true)
		messageImprint, err = httpHeadersPayloadBuilder.Build()
	} else {
		messageImprint, err = DSSJsonUtilsConcatenateDSSDocuments(toSignDocuments,
			parameters.IsBase64UrlEncodedPayload())
	}
	if err != nil {
		panic(err)
	}

	digestAlgorithm := parameters.GetContentTimestampParameters().DigestAlgorithm()
	digest, err := spi.DSSUtilsDigest(digestAlgorithm, messageImprint)
	if err != nil {
		panic(err)
	}
	timeStampResponse, err := s.TspSource.TimeStampResponse(digestAlgorithm, digest)
	if err != nil {
		panic(err)
	}
	timestampToken, err := validation.NewTimestampToken(timeStampResponse.Bytes(),
		enumerations.TimestampTypeContentTimestamp)
	if err != nil {
		panic(fmt.Errorf("Cannot create a content TimestampToken: %w", err))
	}
	return timestampToken
}

// jadesServiceAssertContentTimestampCreationPossible ports the private
// assertContentTimestampCreationPossible.
func jadesServiceAssertContentTimestampCreationPossible(documents []model.DSSDocument) error {
	if utils.IsCollectionEmpty(documents) {
		return fmt.Errorf("Original documents must be provided to generate a content timestamp!")
	}
	for _, doc := range documents {
		if _, isDigestDocument := doc.(*model.DigestDocument); isDigestDocument {
			return fmt.Errorf("Content timestamp creation is not possible with DigestDocument!")
		}
	}
	return nil
}

// GetDataToSign retrieves the data to be signed.
// Port of the #getDataToSign(DSSDocument, JAdESSignatureParameters) overload.
func (s *JAdESService) GetDataToSign(toSignDocument model.DSSDocument,
	parameters *JAdESSignatureParameters) *model.ToBeSigned {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	return s.GetDataToSignForDocuments([]model.DSSDocument{toSignDocument}, parameters)
}

// GetDataToSignForDocuments retrieves the data to be signed over several documents.
// Port of the #getDataToSign(List<DSSDocument>, JAdESSignatureParameters) overload.
func (s *JAdESService) GetDataToSignForDocuments(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters) *model.ToBeSigned {
	if toSignDocuments == nil {
		panic("toSignDocuments cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}

	if err := jadesServiceAssertMultiDocumentsAllowed(toSignDocuments, parameters); err != nil {
		panic(err)
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	jadesBuilder, err := s.JAdESBuilder(parameters, toSignDocuments)
	if err != nil {
		panic(err)
	}
	toBeSigned, err := jadesBuilder.BuildDataToBeSigned()
	if err != nil {
		panic(err)
	}
	return toBeSigned
}

// jadesServiceAssertMultiDocumentsAllowed checks that only DETACHED signatures carry several
// documents. Port of the private assertMultiDocumentsAllowed.
func jadesServiceAssertMultiDocumentsAllowed(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters) error {
	if parameters.SignaturePackaging() == "" {
		panic("SignaturePackaging shall be defined!")
	}

	if utils.IsCollectionEmpty(toSignDocuments) {
		return fmt.Errorf("The documents to sign must be provided!")
	}
	signaturePackaging := parameters.SignaturePackaging()
	if enumerations.SignaturePackagingDetached != signaturePackaging && len(toSignDocuments) > 1 {
		return fmt.Errorf("Not supported operation (only DETACHED are allowed for multiple document signing)!")
	}
	if enumerations.SignaturePackagingDetached == signaturePackaging &&
		enumerations.SigDMechanismNoSigD == parameters.SigDMechanism() && len(toSignDocuments) > 1 {
		return fmt.Errorf("NO_SIG_D mechanism is not allowed for multiple documents!")
	}
	return nil
}

// SignDocument signs the document with the provided signature value.
// Port of the #signDocument(DSSDocument, JAdESSignatureParameters, SignatureValue) overload.
func (s *JAdESService) SignDocument(toSignDocument model.DSSDocument,
	parameters *JAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	return s.SignDocuments([]model.DSSDocument{toSignDocument}, parameters, signatureValue)
}

// SignDocuments signs the documents with the provided signature value.
// Port of the #signDocument(List<DSSDocument>, JAdESSignatureParameters, SignatureValue) overload.
func (s *JAdESService) SignDocuments(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocuments == nil {
		panic("toSignDocuments cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}
	if err := jadesServiceAssertMultiDocumentsAllowed(toSignDocuments, parameters); err != nil {
		panic(err)
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	jadesBuilder, err := s.JAdESBuilder(parameters, toSignDocuments)
	if err != nil {
		panic(err)
	}
	signedDocument, err := jadesBuilder.Build(signatureValue)
	if err != nil {
		panic(err)
	}

	signatureExtension := s.extensionProfile(parameters)
	if signatureExtension != nil {
		if enumerations.SignaturePackagingDetached == parameters.SignaturePackaging() &&
			utils.IsCollectionEmpty(parameters.DetachedContents()) {
			parameters.GetContext().SetDetachedContents(toSignDocuments)
		}
		signatureExtension.SetOperationKind(enumerations.SigningOperationSign)
		if signedDocument, err = signatureExtension.ExtendSignaturesDocument(signedDocument,
			parameters); err != nil {
			panic(err)
		}
	}

	parameters.Reinit()
	name, err := s.GetFinalFileNameWithLevel(toSignDocuments[0], enumerations.SigningOperationSign,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	signedDocument.SetName(name)
	signedDocument.SetMimeType(jadesBuilder.MimeType())
	return signedDocument
}

// JAdESBuilder returns the JAdESBuilder to be used. Port of the protected #getJAdESBuilder.
func (s *JAdESService) JAdESBuilder(parameters *JAdESSignatureParameters,
	documentsToSign []model.DSSDocument) (JAdESBuilder, error) {
	jwsJsonSerializationObject := jadesServiceJWSJsonSerializationObjectToSign(documentsToSign)
	if jadesServiceContainsSignatures(jwsJsonSerializationObject) {
		if !jwsJsonSerializationObject.IsValid() {
			return nil, exception.NewIllegalInputException(fmt.Sprintf(
				"Parallel signing is not supported for invalid RFC 7515 signatures. Reason(s) : %s",
				jwsJsonSerializationObject.StructuralValidationErrors()))
		}
		// return a builder for parallel signing
		return NewJAdESSerializationBuilderFromSignature(s.CertificateVerifier, parameters,
			jwsJsonSerializationObject)
	}

	switch parameters.JwsSerializationType() {
	case enumerations.JWSSerializationTypeCompactSerialization:
		return NewJAdESCompactBuilder(s.CertificateVerifier, parameters, documentsToSign)
	case enumerations.JWSSerializationTypeJSONSerialization,
		enumerations.JWSSerializationTypeFlattenedJSONSerialization:
		return NewJAdESSerializationBuilder(s.CertificateVerifier, parameters, documentsToSign)
	default:
		return nil, fmt.Errorf("The requested JWS Serialization Type '%s' is not supported!",
			parameters.JwsSerializationType())
	}
}

// jadesServiceJWSJsonSerializationObjectToSign ports the private
// getJWSJsonSerializationObjectToSign.
func jadesServiceJWSJsonSerializationObjectToSign(
	documentsToSign []model.DSSDocument) *JWSJsonSerializationObject {
	if utils.IsCollectionNotEmpty(documentsToSign) && len(documentsToSign) == 1 {
		doc := documentsToSign[0]
		documentAnalyzerFactory := NewJWSDocumentAnalyzerFactory()
		if documentAnalyzerFactory.IsSupported(doc) {
			documentAnalyzer := jwsDocumentAnalyzerBase(documentAnalyzerFactory.Create(doc))
			return documentAnalyzer.JwsJsonSerializationObject()
		}
	}
	return nil
}

// jadesServiceContainsSignatures ports the private containsSignatures.
func jadesServiceContainsSignatures(jwsJsonSerializationObject *JWSJsonSerializationObject) bool {
	return jwsJsonSerializationObject != nil &&
		utils.IsCollectionNotEmpty(jwsJsonSerializationObject.Signatures())
}

// ExtendDocument extends the level of the signatures in the document. Port of #extendDocument.
func (s *JAdESService) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *JAdESSignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument cannot be null!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}
	if parameters.SignatureLevel() == "" {
		panic("SignatureLevel must be defined!")
	}
	if err := jadesServiceAssertExtensionPossible(parameters); err != nil {
		panic(err)
	}

	signatureExtension := s.extensionProfile(parameters)
	if signatureExtension != nil {
		signatureExtension.SetOperationKind(enumerations.SigningOperationExtend)
		dssDocument, err := signatureExtension.ExtendSignaturesDocument(toExtendDocument, parameters)
		if err != nil {
			panic(err)
		}
		name, err := s.GetFinalFileNameWithLevel(toExtendDocument, enumerations.SigningOperationExtend,
			parameters.SignatureLevel())
		if err != nil {
			panic(err)
		}
		dssDocument.SetName(name)
		dssDocument.SetMimeType(enumerations.MimeTypeEnumJOSEJSON)
		return dssDocument
	}
	panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", parameters.SignatureLevel()))
}

// jadesServiceAssertExtensionPossible ports the private assertExtensionPossible.
func jadesServiceAssertExtensionPossible(parameters *JAdESSignatureParameters) error {
	if enumerations.JWSSerializationTypeJSONSerialization != parameters.JwsSerializationType() &&
		enumerations.JWSSerializationTypeFlattenedJSONSerialization != parameters.JwsSerializationType() {
		return fmt.Errorf("The type '%s' does not support signature extension!",
			parameters.JwsSerializationType())
	}
	return nil
}

// extensionProfile chooses the extension profile according to the passed parameters, returning
// nil for the -B level exactly as Java returns null. Port of the private getExtensionProfile.
func (s *JAdESService) extensionProfile(parameters *JAdESSignatureParameters) JAdESSignatureExtender {
	var extension JAdESSignatureExtender
	switch parameters.SignatureLevel() {
	case enumerations.SignatureLevelJAdESBaselineB:
		return nil
	case enumerations.SignatureLevelJAdESBaselineT:
		extension = NewJAdESLevelBaselineT(s.CertificateVerifier)
	case enumerations.SignatureLevelJAdESBaselineLT:
		extension = NewJAdESLevelBaselineLT(s.CertificateVerifier)
	case enumerations.SignatureLevelJAdESBaselineLTA:
		extension = NewJAdESLevelBaselineLTA(s.CertificateVerifier)
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", parameters.SignatureLevel()))
	}
	extension.SetTspSource(s.TspSource)
	return extension
}

// TimestampDocuments is unsupported for this file format.
// Port of the overridden #timestamp(List<DSSDocument>, JAdESTimestampParameters).
func (s *JAdESService) TimestampDocuments(toTimestampDocuments []model.DSSDocument,
	parameters *JAdESTimestampParameters) model.DSSDocument {
	panic("Unsupported operation for this file format")
}

// AddSignaturePolicyStore incorporates a Signature Policy Store as a base64Url-encoded unsigned
// property into the JAdES Signature.
// Port of the #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore) overload.
func (s *JAdESService) AddSignaturePolicyStore(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	return s.AddSignaturePolicyStoreWithEncoding(doc, signaturePolicyStore, true)
}

// AddSignaturePolicyStoreWithEncoding incorporates a Signature Policy Store as an unsigned
// property into the JAdES Signature. base64UrlInstance defines whether the SignaturePolicyStore
// shall be incorporated in its corresponding base64Url representation; when FALSE it is
// incorporated in the clear JSON representation.
// Port of the #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore, boolean) overload.
func (s *JAdESService) AddSignaturePolicyStoreWithEncoding(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore, base64UrlInstance bool) model.DSSDocument {
	if doc == nil {
		panic("The document cannot be null")
	}
	if signaturePolicyStore == nil {
		panic("The signaturePolicyStore cannot be null")
	}

	builder := NewJAdESSignaturePolicyStoreBuilder()
	signatureWithPolicyStore, err := builder.AddSignaturePolicyStore(doc, signaturePolicyStore,
		base64UrlInstance)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileName(doc, enumerations.SigningOperationAddSigPolicyStore)
	if err != nil {
		panic(err)
	}
	signatureWithPolicyStore.SetName(name)
	signatureWithPolicyStore.SetMimeType(doc.MimeType())
	return signatureWithPolicyStore
}

// GetDataToBeCounterSigned retrieves the data to be counter-signed.
// Port of #getDataToBeCounterSigned.
func (s *JAdESService) GetDataToBeCounterSigned(signatureDocument model.DSSDocument,
	parameters *JAdESCounterSignatureParameters) *model.ToBeSigned {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if err := jadesServiceVerifyAndSetCounterSignatureParameters(parameters); err != nil {
		panic(err)
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	counterSignatureBuilder := NewJAdESCounterSignatureBuilder()
	signatureValueToSign, err := counterSignatureBuilder.GetSignatureValueToBeSigned(signatureDocument,
		parameters)
	if err != nil {
		panic(err)
	}

	return s.GetDataToSign(signatureValueToSign, &parameters.JAdESSignatureParameters)
}

// CounterSignSignature counter-signs the signature document with the provided signature value.
// Port of #counterSignSignature.
func (s *JAdESService) CounterSignSignature(signatureDocument model.DSSDocument,
	parameters *JAdESCounterSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("signatureValue cannot be null!")
	}
	if err := jadesServiceVerifyAndSetCounterSignatureParameters(parameters); err != nil {
		panic(err)
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	counterSignatureBuilder := NewJAdESCounterSignatureBuilder()
	signatureValueToSign, err := counterSignatureBuilder.GetSignatureValueToBeSigned(signatureDocument,
		parameters)
	if err != nil {
		panic(err)
	}

	counterSignature := s.SignDocument(signatureValueToSign, &parameters.JAdESSignatureParameters,
		signatureValue)

	counterSigned, err := counterSignatureBuilder.BuildEmbeddedCounterSignature(signatureDocument,
		counterSignature, parameters)
	if err != nil {
		panic(err)
	}

	parameters.Reinit()
	name, err := s.GetFinalFileNameWithLevel(signatureDocument, enumerations.SigningOperationCounterSign,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	counterSigned.SetName(name)
	counterSigned.SetMimeType(signatureDocument.MimeType())

	return counterSigned
}

// jadesServiceVerifyAndSetCounterSignatureParameters ports the private
// verifyAndSetCounterSignatureParameters.
func jadesServiceVerifyAndSetCounterSignatureParameters(
	parameters *JAdESCounterSignatureParameters) error {
	if parameters.SignaturePackaging() == "" {
		// attached counter signature is created by default
		parameters.SetSignaturePackaging(enumerations.SignaturePackagingEnveloping)
	}

	switch parameters.SignaturePackaging() {
	case enumerations.SignaturePackagingEnveloping:
		// nothing to do
	case enumerations.SignaturePackagingDetached:
		if parameters.SigDMechanism() == "" {
			parameters.SetSigDMechanism(enumerations.SigDMechanismNoSigD)
		} else if enumerations.SigDMechanismNoSigD != parameters.SigDMechanism() {
			return fmt.Errorf("The SigDMechanism '%s' is not supported by JAdES Counter Signature!",
				parameters.SigDMechanism())
		}
	default:
		return fmt.Errorf("The SignaturePackaging '%s' is not supported by JAdES Counter Signature!",
			parameters.SignaturePackaging())
	}

	if enumerations.JWSSerializationTypeJSONSerialization == parameters.JwsSerializationType() {
		return fmt.Errorf("The JWSSerializationType.JSON_SERIALIZATION parameter " +
			"is not supported for a JAdES Counter Signature!")
	}
	if parameters.ContentType() != "" {
		return fmt.Errorf("Content Type protected header shall not be present " +
			"for a JAdES Counter Signature!")
	}
	return nil
}

// IsValidSignatureValue verifies the signature value against a ToBeSigned and a CertificateToken,
// additionally enforcing the ECDSA curve/digest pairing RFC 7518 requires for a JWS.
// Port of the overridden #isValidSignatureValue.
func (s *JAdESService) IsValidSignatureValue(toBeSigned *model.ToBeSigned,
	signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	if !s.AbstractSignatureService.IsValidSignatureValue(toBeSigned, signatureValue, signingCertificate) {
		return false
	}

	if err := jadesServiceAssertSigningCertificateValidForAlgorithm(signatureValue.Algorithm(),
		signingCertificate); err != nil {
		// Upstream warns "Invalid signature value : {}".
		return false
	}
	if err := jadesServiceAssertSignatureValueValid(signatureValue.Algorithm(), signatureValue); err != nil {
		// Upstream warns "Invalid signature value : {}".
		return false
	}
	return true
}

// AssertSigningCertificateValid raises a panic if the signing rules forbid the use of the
// certificate, adding the JAdES-specific ECDSA checks on top of the base implementation.
// Port of the protected, overridden #assertSigningCertificateValid(AbstractSignatureParameters).
func (s *JAdESService) AssertSigningCertificateValid(
	parameters *document.AbstractSignatureParameters[*JAdESTimestampParameters]) {
	s.AbstractSignatureService.AssertSigningCertificateValid(parameters)
	if err := jadesServiceAssertSigningCertificateValidForAlgorithm(parameters.SignatureAlgorithm(),
		parameters.SigningCertificate()); err != nil {
		panic(err)
	}
}

// jadesServiceAssertSigningCertificateValidForAlgorithm ports the private
// assertSigningCertificateValid(SignatureAlgorithm, CertificateToken).
func jadesServiceAssertSigningCertificateValidForAlgorithm(
	signatureAlgorithm enumerations.SignatureAlgorithm,
	signingCertificate *model.CertificateToken) error {
	if signatureAlgorithm.EncryptionAlgorithm() == "" || signatureAlgorithm.DigestAlgorithm() == "" ||
		!signatureAlgorithm.EncryptionAlgorithm().IsEquivalent(enumerations.EncryptionAlgorithmECDSA) ||
		signingCertificate == nil {
		return nil
	}
	errorMessage := "For ECDSA with %s a key with P-%s curve shall be used for a JWS! See RFC 7518."
	keySize := spi.DSSPKUtilsPublicKeySize(signingCertificate.PublicKey())
	switch signatureAlgorithm.DigestAlgorithm() {
	case enumerations.DigestAlgorithmSHA256:
		if keySize != 256 {
			return fmt.Errorf(errorMessage, signatureAlgorithm.DigestAlgorithm(), "256")
		}
	case enumerations.DigestAlgorithmSHA384:
		if keySize != 384 {
			return fmt.Errorf(errorMessage, signatureAlgorithm.DigestAlgorithm(), "384")
		}
	case enumerations.DigestAlgorithmSHA512:
		if keySize != 521 {
			return fmt.Errorf(errorMessage, signatureAlgorithm.DigestAlgorithm(), "521")
		}
	default:
		return fmt.Errorf("ECDSA with %s is not supported for JWS!", signatureAlgorithm.DigestAlgorithm())
	}
	return nil
}

// jadesServiceAssertSignatureValueValid ports the private assertSignatureValueValid.
func jadesServiceAssertSignatureValueValid(targetSignatureAlgorithm enumerations.SignatureAlgorithm,
	signatureValue *model.SignatureValue) error {
	if !targetSignatureAlgorithm.EncryptionAlgorithm().IsEquivalent(enumerations.EncryptionAlgorithmECDSA) {
		return nil
	}
	errorMessage := "Invalid SignatureValue obtained! " +
		"For ECDSA with %s a key with P-%s curve shall be used for a JWS. See RFC 7518."
	bitLength, err := spi.DSSASN1UtilsSignatureValueBitLength(signatureValue.Value())
	if err != nil {
		return err
	}
	switch targetSignatureAlgorithm.DigestAlgorithm() {
	case enumerations.DigestAlgorithmSHA256:
		if bitLength != 256 {
			return exception.NewIllegalInputException(fmt.Sprintf(errorMessage,
				targetSignatureAlgorithm.DigestAlgorithm(), "256"))
		}
	case enumerations.DigestAlgorithmSHA384:
		if bitLength != 384 {
			return fmt.Errorf(errorMessage, targetSignatureAlgorithm.DigestAlgorithm(), "384")
		}
	case enumerations.DigestAlgorithmSHA512:
		if bitLength != 520 && bitLength != 528 {
			return fmt.Errorf(errorMessage, targetSignatureAlgorithm.DigestAlgorithm(), "521")
		}
	default:
		return fmt.Errorf("ECDSA with %s is not supported for JWS!",
			targetSignatureAlgorithm.DigestAlgorithm())
	}
	return nil
}

// MultipleDocumentsService adapts this service to document.MultipleDocumentsSignatureService.
// See the package-level note: Java's JAdESService implements that interface directly, which Go
// cannot express on the same type because DocumentSignatureService declares the same four method
// names with single-document parameters.
func (s *JAdESService) MultipleDocumentsService() document.MultipleDocumentsSignatureService[
	*JAdESSignatureParameters, *JAdESTimestampParameters] {
	return &jadesServiceMultipleDocumentsAdapter{service: s}
}

// jadesServiceMultipleDocumentsAdapter forwards MultipleDocumentsSignatureService to the
// ...ForDocuments / SignDocuments / TimestampDocuments methods of JAdESService. It holds no state
// and adds no behaviour.
type jadesServiceMultipleDocumentsAdapter struct {
	service *JAdESService
}

func (a *jadesServiceMultipleDocumentsAdapter) GetContentTimestamp(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters) *validation.TimestampToken {
	return a.service.GetContentTimestampForDocuments(toSignDocuments, parameters)
}

func (a *jadesServiceMultipleDocumentsAdapter) GetDataToSign(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters) *model.ToBeSigned {
	return a.service.GetDataToSignForDocuments(toSignDocuments, parameters)
}

func (a *jadesServiceMultipleDocumentsAdapter) IsValidSignatureValue(toBeSigned *model.ToBeSigned,
	signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	return a.service.IsValidSignatureValue(toBeSigned, signatureValue, signingCertificate)
}

func (a *jadesServiceMultipleDocumentsAdapter) SignDocument(toSignDocuments []model.DSSDocument,
	parameters *JAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	return a.service.SignDocuments(toSignDocuments, parameters, signatureValue)
}

func (a *jadesServiceMultipleDocumentsAdapter) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *JAdESSignatureParameters) model.DSSDocument {
	return a.service.ExtendDocument(toExtendDocument, parameters)
}

func (a *jadesServiceMultipleDocumentsAdapter) Timestamp(toTimestampDocuments []model.DSSDocument,
	parameters *JAdESTimestampParameters) model.DSSDocument {
	return a.service.TimestampDocuments(toTimestampDocuments, parameters)
}

// Compile-time assertions that JAdESService satisfies the service interfaces Java's JAdESService
// carries with single-document (or non-overloaded) signatures, and that the adapter satisfies the
// multi-document one.
var (
	_ document.DocumentSignatureService[*JAdESSignatureParameters, *JAdESTimestampParameters]          = (*JAdESService)(nil)
	_ document.CounterSignatureService[*JAdESCounterSignatureParameters]                               = (*JAdESService)(nil)
	_ document.MultipleDocumentsSignatureService[*JAdESSignatureParameters, *JAdESTimestampParameters] = (*jadesServiceMultipleDocumentsAdapter)(nil)
)
