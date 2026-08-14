// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESService.java (DSS 6.5.RC1).
//
// Java extends AbstractSignatureService<XAdESSignatureParameters, XAdESTimestampParameters> and
// implements MultipleDocumentsSignatureService, CounterSignatureService and
// EvidenceRecordIncorporationService; the Go port embeds
// document.AbstractSignatureService[*XAdESSignatureParameters, *XAdESTimestampParameters], the
// same way cades.CAdESService does.
//
// # Overloads, and the one interface Go cannot let this type carry
//
// Java overloads getContentTimestamp, getDataToSign, signDocument and timestamp on
// DSSDocument vs List<DSSDocument>: DocumentSignatureService declares the first shape,
// MultipleDocumentsSignatureService the second, and one Java class satisfies both. Go has no
// overloading, so a single type cannot carry both method sets. This port keeps the plain names
// for the single-document shape - so *XAdESService satisfies document.DocumentSignatureService,
// as cades.CAdESService does - and gives the multi-document shape distinct names:
//
//	getContentTimestamp(List, SP)          -> GetContentTimestampForDocuments
//	getDataToSign(List, SP)                -> GetDataToSignForDocuments
//	signDocument(List, SP, SignatureValue) -> SignDocuments
//	timestamp(List, TP)                    -> TimestampDocuments
//
// JUDGMENT CALL (flagged for the integrator): so that the Java contract is not simply lost,
// MultipleDocumentsService() returns a thin adapter that does satisfy
// document.MultipleDocumentsSignatureService by forwarding to the four methods above. It adds
// no behaviour; ASiC-XAdES (phase 7) is the caller that will need it.
//
// # Errors
//
// The document service interfaces return bare values, so - exactly as in cades_service.go -
// this is where the (T, error) of the layers below turns back into Java's propagating
// exception, i.e. into a panic whose value is the error itself.
//
// java.io.Serializable, the serialVersionUID and slf4j are dropped.
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// init ports the static initializer block of XAdESService.
func init() {
	xmlutils.SantuarioInitializerInit()
	DSSXMLUtilsRegisterXAdESNamespaces()
}

// XAdESSignatureExtender is the slice of SignatureExtension<XAdESSignatureParameters> that
// getExtensionProfile's local variable is typed with upstream: Java assigns an
// XAdESLevelBaselineT, -C, -X, -XL, -A, -LT or -LTA to it, and each of the seven is a distinct
// Go type embedding XAdESLevelBaselineT, so the switch below needs an interface value. Every
// XAdES extension level satisfies it through that embedded base. Same technique, and same
// reason, as cades.CAdESSignatureExtender.
type XAdESSignatureExtender interface {
	// SetTspSource sets the TSP source used when extending. Port of
	// XAdESLevelBaselineT#setTspSource.
	SetTspSource(tspSource validation.TSPSource)

	// ExtendSignaturesDocument extends the signatures of the given document. Port of
	// SignatureExtension#extendSignatures(DSSDocument, XAdESSignatureParameters); the plain
	// ExtendSignatures name belongs to the List<AdvancedSignature> overload every level
	// overrides, so the document-taking entry point carries the Document suffix.
	ExtendSignaturesDocument(document model.DSSDocument,
		params *XAdESSignatureParameters) (model.DSSDocument, error)
}

// XAdESService is the XAdES implementation of DocumentSignatureService.
type XAdESService struct {
	document.AbstractSignatureService[*XAdESSignatureParameters, *XAdESTimestampParameters]
}

// NewXAdESService creates an instance of the XAdESService. A certificate verifier must be
// provided; it supplies information on the sources to be used in the validation process in the
// context of a signature. Port of XAdESService(CertificateVerifier).
func NewXAdESService(certificateVerifier validation.CertificateVerifier) *XAdESService {
	// Upstream logs "+ XAdESService created".
	return &XAdESService{
		AbstractSignatureService: document.NewAbstractSignatureService[*XAdESSignatureParameters,
			*XAdESTimestampParameters](certificateVerifier),
	}
}

// GetContentTimestamp creates a content-timestamp covering the document to be signed.
// Port of the #getContentTimestamp(DSSDocument, XAdESSignatureParameters) overload.
func (s *XAdESService) GetContentTimestamp(toSignDocument model.DSSDocument,
	parameters *XAdESSignatureParameters) *validation.TimestampToken {
	return s.GetContentTimestampForDocuments([]model.DSSDocument{toSignDocument}, parameters)
}

// GetContentTimestampForDocuments creates a content-timestamp covering all documents to be
// signed. Port of the #getContentTimestamp(List<DSSDocument>, XAdESSignatureParameters) overload.
func (s *XAdESService) GetContentTimestampForDocuments(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters) *validation.TimestampToken {
	if s.TspSource == nil {
		panic("A TSPSource is required !")
	}
	builder := NewAllDataObjectsTimeStampBuilder(s.TspSource, parameters)
	token, err := builder.BuildForDocuments(toSignDocuments)
	if err != nil {
		panic(err)
	}
	return token
}

// GetDataToSign retrieves the data to be signed.
// Port of the #getDataToSign(DSSDocument, XAdESSignatureParameters) overload.
func (s *XAdESService) GetDataToSign(toSignDocument model.DSSDocument,
	parameters *XAdESSignatureParameters) *model.ToBeSigned {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	return s.GetDataToSignForDocuments([]model.DSSDocument{toSignDocument}, parameters)
}

// GetDataToSignForDocuments retrieves the data to be signed over several documents.
// Port of the #getDataToSign(List<DSSDocument>, XAdESSignatureParameters) overload.
func (s *XAdESService) GetDataToSignForDocuments(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters) *model.ToBeSigned {
	if toSignDocuments == nil {
		panic("toSignDocuments cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}

	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	xadesServiceAssertMultiDocumentsAllowed(toSignDocuments, parameters)
	xadesServiceAssertDocumentsValid(toSignDocuments)

	levelBaselineB := NewXAdESLevelBaselineB(s.CertificateVerifier)
	dataToSign, err := levelBaselineB.GetDataToSignForDocuments(toSignDocuments, parameters)
	if err != nil {
		panic(err)
	}
	// Upstream traces the data to sign here.
	parameters.GetContext().SetProfile(levelBaselineB)
	return model.NewToBeSignedWithBytes(dataToSign)
}

// SignDocument signs the document with the provided signature value.
// Port of the #signDocument(DSSDocument, XAdESSignatureParameters, SignatureValue) overload.
func (s *XAdESService) SignDocument(toSignDocument model.DSSDocument,
	parameters *XAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument is not defined!")
	}
	return s.SignDocuments([]model.DSSDocument{toSignDocument}, parameters, signatureValue)
}

// SignDocuments signs the documents with the provided signature value.
// Port of the #signDocument(List<DSSDocument>, XAdESSignatureParameters, SignatureValue) overload.
func (s *XAdESService) SignDocuments(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocuments == nil {
		panic("toSignDocuments are not defined!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if parameters.SignatureLevel() == "" {
		panic("SignatureLevel must be defined!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}

	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	xadesServiceAssertMultiDocumentsAllowed(toSignDocuments, parameters)
	xadesServiceAssertDocumentsValid(toSignDocuments)

	parameters.GetContext().SetOperationKind(enumerations.SigningOperation_SIGN)
	var profile XAdESSignatureProfile
	context := parameters.GetContext()
	if context.Profile() != nil {
		profile = context.Profile()
	} else {
		profile = NewXAdESLevelBaselineB(s.CertificateVerifier)
	}

	result, err := profile.SignDocuments(toSignDocuments, parameters, signatureValue.Value())
	if err != nil {
		panic(err)
	}
	extension := s.extensionProfile(parameters)
	if extension != nil {
		if enumerations.SignaturePackaging_DETACHED == parameters.SignaturePackaging() {
			parameters.GetContext().SetDetachedContents(toSignDocuments)
		}
		if result, err = extension.ExtendSignaturesDocument(result, parameters); err != nil {
			panic(err)
		}
	}

	// The internal parameters (e.g. deterministic Id) are reset between two consecutive signing
	// operations. It prevents sharing two signatures the same cached data.
	parameters.Reinit()
	name, err := s.GetFinalFileNameWithLevel(toSignDocuments[0], enumerations.SigningOperation_SIGN,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	result.SetName(name)
	return result
}

// ExtendDocument extends the level of the signatures in the document. Port of #extendDocument.
func (s *XAdESService) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *XAdESSignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument cannot be null!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}
	if parameters.SignatureLevel() == "" {
		panic("SignatureLevel must be defined!")
	}

	parameters.GetContext().SetOperationKind(enumerations.SigningOperation_EXTEND)
	extension := s.extensionProfile(parameters)
	if extension != nil {
		dssDocument, err := extension.ExtendSignaturesDocument(toExtendDocument, parameters)
		if err != nil {
			panic(err)
		}
		name, err := s.GetFinalFileNameWithLevel(toExtendDocument, enumerations.SigningOperation_EXTEND,
			parameters.SignatureLevel())
		if err != nil {
			panic(err)
		}
		dssDocument.SetName(name)
		return dssDocument
	}
	panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", parameters.SignatureLevel()))
}

// TimestampDocuments is unsupported for this file format.
// Port of the overridden #timestamp(List<DSSDocument>, XAdESTimestampParameters).
func (s *XAdESService) TimestampDocuments(toTimestampDocuments []model.DSSDocument,
	parameters *XAdESTimestampParameters) model.DSSDocument {
	panic("Unsupported operation for this file format")
}

// extensionProfile chooses the extension profile according to the passed parameters, returning
// nil for the -B level exactly as Java returns null. Port of the private getExtensionProfile.
func (s *XAdESService) extensionProfile(parameters *XAdESSignatureParameters) XAdESSignatureExtender {
	var extension XAdESSignatureExtender
	switch parameters.SignatureLevel() {
	case enumerations.SignatureLevel_XAdES_BASELINE_B:
		return nil
	case enumerations.SignatureLevel_XAdES_BASELINE_T:
		extension = NewXAdESLevelBaselineT(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_C:
		extension = NewXAdESLevelC(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_X:
		extension = NewXAdESLevelX(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_XL:
		extension = NewXAdESLevelXL(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_A:
		extension = NewXAdESLevelA(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_BASELINE_LT:
		extension = NewXAdESLevelBaselineLT(s.CertificateVerifier)
	case enumerations.SignatureLevel_XAdES_BASELINE_LTA:
		extension = NewXAdESLevelBaselineLTA(s.CertificateVerifier)
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", parameters.SignatureLevel()))
	}
	extension.SetTspSource(s.TspSource)
	return extension
}

// xadesServiceAssertMultiDocumentsAllowed checks that only DETACHED and ENVELOPING signatures
// carry several documents. Port of the private assertMultiDocumentsAllowed.
func xadesServiceAssertMultiDocumentsAllowed(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters) {
	if parameters.SignaturePackaging() == "" {
		panic("SignaturePackaging shall be defined!")
	}

	if utils.CollectionSize(toSignDocuments) == 0 {
		panic("The documents to sign must be provided!")

	} else if utils.CollectionSize(toSignDocuments) > 1 {
		signaturePackaging := parameters.SignaturePackaging()
		if signaturePackaging == "" || enumerations.SignaturePackaging_ENVELOPED == signaturePackaging {
			panic("Not supported operation (only DETACHED or ENVELOPING are allowed)")
		}
	}
}

// xadesServiceAssertDocumentsValid ports the private assertDocumentsValid.
func xadesServiceAssertDocumentsValid(toSignDocuments []model.DSSDocument) {
	documentNames := make([]string, 0)
	for _, doc := range toSignDocuments {
		if doc == nil {
			panic("Document to sign cannot be null!")
		}

		if len(toSignDocuments) > 1 && utils.IsStringBlank(doc.Name()) {
			panic("All documents in the list to be signed shall have names!")
		}
		for _, name := range documentNames {
			if name == doc.Name() {
				panic(fmt.Sprintf("The documents to be signed shall have different names! "+
					"The name '%s' appears multiple times.", doc.Name()))
			}
		}
		documentNames = append(documentNames, doc.Name())
	}
}

// AddSignaturePolicyStore incorporates a Signature Policy Store as an unsigned property into the
// XAdES Signature. Port of #addSignaturePolicyStore.
func (s *XAdESService) AddSignaturePolicyStore(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	if doc == nil {
		panic("The document cannot be null")
	}
	if signaturePolicyStore == nil {
		panic("The signaturePolicyStore cannot be null")
	}

	builder := NewSignaturePolicyStoreBuilder()
	signatureWithPolicyStore, err := builder.AddSignaturePolicyStore(doc, signaturePolicyStore)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileName(doc, enumerations.SigningOperation_ADD_SIG_POLICY_STORE)
	if err != nil {
		panic(err)
	}
	signatureWithPolicyStore.SetName(name)
	signatureWithPolicyStore.SetMimeType(doc.MimeType())
	return signatureWithPolicyStore
}

// GetDataToBeCounterSigned retrieves the data to be counter-signed.
// Port of #getDataToBeCounterSigned.
func (s *XAdESService) GetDataToBeCounterSigned(signatureDocument model.DSSDocument,
	parameters *XAdESCounterSignatureParameters) *model.ToBeSigned {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	xadesServiceVerifyAndSetCounterSignatureParameters(parameters)

	counterSignatureBuilder := NewCounterSignatureBuilder(s.CertificateVerifier)
	signatureValue, err := counterSignatureBuilder.GetCanonicalizedSignatureValue(signatureDocument, parameters)
	if err != nil {
		panic(err)
	}

	counterSignatureReference, err := counterSignatureBuilder.BuildCounterSignatureDSSReference(
		signatureDocument, parameters)
	if err != nil {
		panic(err)
	}
	parameters.SetReferences([]*DSSReference{counterSignatureReference})

	return s.GetDataToSign(signatureValue, &parameters.XAdESSignatureParameters)
}

// CounterSignSignature counter-signs the signature document with the provided signature value.
// Port of #counterSignSignature.
func (s *XAdESService) CounterSignSignature(signatureDocument model.DSSDocument,
	parameters *XAdESCounterSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("signatureValue cannot be null!")
	}
	xadesServiceVerifyAndSetCounterSignatureParameters(parameters)

	counterSignatureBuilder := NewCounterSignatureBuilder(s.CertificateVerifier)
	signatureValueToSign, err := counterSignatureBuilder.GetCanonicalizedSignatureValue(
		signatureDocument, parameters)
	if err != nil {
		panic(err)
	}
	parameters.GetContext().SetDetachedContents([]model.DSSDocument{signatureValueToSign})

	counterSignatureReference, err := counterSignatureBuilder.BuildCounterSignatureDSSReference(
		signatureDocument, parameters)
	if err != nil {
		panic(err)
	}
	parameters.SetReferences([]*DSSReference{counterSignatureReference})

	counterSignature := s.SignDocument(signatureValueToSign, &parameters.XAdESSignatureParameters, signatureValue)
	counterSigned, err := counterSignatureBuilder.BuildEmbeddedCounterSignature(signatureDocument,
		counterSignature, parameters)
	if err != nil {
		panic(err)
	}

	parameters.Reinit()
	name, err := s.GetFinalFileNameWithLevel(signatureDocument, enumerations.SigningOperation_COUNTER_SIGN,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	counterSigned.SetName(name)
	counterSigned.SetMimeType(signatureDocument.MimeType())

	return counterSigned
}

// AddSignatureEvidenceRecord incorporates an evidence record as an unsigned property into the
// XAdES Signature. Port of #addSignatureEvidenceRecord.
func (s *XAdESService) AddSignatureEvidenceRecord(signatureDocument, evidenceRecordDocument model.DSSDocument,
	parameters *XAdESEvidenceRecordIncorporationParameters) model.DSSDocument {
	if signatureDocument == nil {
		panic("The signature document cannot be null")
	}
	if evidenceRecordDocument == nil {
		panic("The evidence record document cannot be null")
	}

	builder := NewEmbeddedEvidenceRecordBuilder(s.CertificateVerifier)
	signatureWithEvidenceRecord, err := builder.AddEvidenceRecord(signatureDocument, evidenceRecordDocument,
		parameters)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileName(signatureDocument, enumerations.SigningOperation_ADD_EVIDENCE_RECORD)
	if err != nil {
		panic(err)
	}
	signatureWithEvidenceRecord.SetName(name)
	signatureWithEvidenceRecord.SetMimeType(signatureDocument.MimeType())
	return signatureWithEvidenceRecord
}

// xadesServiceVerifyAndSetCounterSignatureParameters ports the private
// verifyAndSetCounterSignatureParameters.
func xadesServiceVerifyAndSetCounterSignatureParameters(parameters *XAdESCounterSignatureParameters) {
	if parameters.SignaturePackaging() == "" {
		parameters.SetSignaturePackaging(enumerations.SignaturePackaging_DETACHED)
	} else if enumerations.SignaturePackaging_DETACHED != parameters.SignaturePackaging() {
		panic(fmt.Sprintf("The SignaturePackaging '%s' is not supported by XAdES Counter Signature!",
			parameters.SignaturePackaging()))
	}
}

// MultipleDocumentsService adapts this service to document.MultipleDocumentsSignatureService.
// See the package-level note: Java's XAdESService implements that interface directly, which Go
// cannot express on the same type because DocumentSignatureService declares the same four
// method names with single-document parameters.
func (s *XAdESService) MultipleDocumentsService() document.MultipleDocumentsSignatureService[
	*XAdESSignatureParameters, *XAdESTimestampParameters] {
	return &xadesServiceMultipleDocumentsAdapter{service: s}
}

// xadesServiceMultipleDocumentsAdapter forwards MultipleDocumentsSignatureService to the
// ...ForDocuments / SignDocuments / TimestampDocuments methods of XAdESService. It holds no
// state and adds no behaviour.
type xadesServiceMultipleDocumentsAdapter struct {
	service *XAdESService
}

func (a *xadesServiceMultipleDocumentsAdapter) GetContentTimestamp(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters) *validation.TimestampToken {
	return a.service.GetContentTimestampForDocuments(toSignDocuments, parameters)
}

func (a *xadesServiceMultipleDocumentsAdapter) GetDataToSign(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters) *model.ToBeSigned {
	return a.service.GetDataToSignForDocuments(toSignDocuments, parameters)
}

func (a *xadesServiceMultipleDocumentsAdapter) IsValidSignatureValue(toBeSigned *model.ToBeSigned,
	signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	return a.service.IsValidSignatureValue(toBeSigned, signatureValue, signingCertificate)
}

func (a *xadesServiceMultipleDocumentsAdapter) SignDocument(toSignDocuments []model.DSSDocument,
	parameters *XAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	return a.service.SignDocuments(toSignDocuments, parameters, signatureValue)
}

func (a *xadesServiceMultipleDocumentsAdapter) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *XAdESSignatureParameters) model.DSSDocument {
	return a.service.ExtendDocument(toExtendDocument, parameters)
}

func (a *xadesServiceMultipleDocumentsAdapter) Timestamp(toTimestampDocuments []model.DSSDocument,
	parameters *XAdESTimestampParameters) model.DSSDocument {
	return a.service.TimestampDocuments(toTimestampDocuments, parameters)
}

// Compile-time assertions that XAdESService satisfies the three service interfaces Java's
// XAdESService carries with single-document (or non-overloaded) signatures, and that the
// adapter above satisfies the fourth.
var (
	_ document.DocumentSignatureService[*XAdESSignatureParameters, *XAdESTimestampParameters]          = (*XAdESService)(nil)
	_ document.CounterSignatureService[*XAdESCounterSignatureParameters]                               = (*XAdESService)(nil)
	_ document.EvidenceRecordIncorporationService[*XAdESEvidenceRecordIncorporationParameters]         = (*XAdESService)(nil)
	_ document.MultipleDocumentsSignatureService[*XAdESSignatureParameters, *XAdESTimestampParameters] = (*xadesServiceMultipleDocumentsAdapter)(nil)
)
