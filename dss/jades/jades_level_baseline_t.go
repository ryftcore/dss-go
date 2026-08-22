// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineT.java (DSS 6.5.RC1).
//
// # The extension chain, and how Java's virtual extendSignatures survives the port
//
// Upstream builds the JAdES augmentation ladder out of one overridden protected method:
//
//	JAdESLevelBaselineT.extendSignatures(List<AdvancedSignature>, JAdESSignatureParameters)
//	  <- JAdESLevelBaselineLT  <- JAdESLevelBaselineLTA
//
// each override calling super.extendSignatures(...) first. Go has no method overriding across
// embedding, so - per the TokenBase.InitToken(self) convention of PORTING.md, and exactly as
// xades.XAdESLevelBaselineT does - every level embeds the level below it, registers itself with
// InitJAdESLevelBaselineT(self, certificateVerifier), and the public entry point
// ExtendSignaturesDocument dispatches into the most-derived override via t.overrides. A "super"
// call is then the plain, explicit embedded-field call, e.g.
// lt.JAdESLevelBaselineT.ExtendSignatures.
//
// Java's two extendSignatures overloads get two Go names:
//
//	extendSignatures(List<AdvancedSignature>, params)   -> ExtendSignatures (the virtual one)
//	extendSignatures(DSSDocument, params)               -> ExtendSignaturesDocument
//
// # Errors
//
// Objects.requireNonNull becomes a panic carrying the Java message; every other Java throw
// becomes a returned error. slf4j logging is dropped.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESSignatureExtensionOverrides declares the operation every JAdES extension level overrides
// and that JAdESLevelBaselineT.ExtendSignaturesDocument dispatches into. Every level satisfies it
// through its embedded ancestors.
type JAdESSignatureExtensionOverrides interface {
	// ExtendSignatures extends the given signatures to the level of the concrete
	// implementation. Port of the protected, overridden #extendSignatures(List, params).
	ExtendSignatures(signatures []validation.AdvancedSignature, params *JAdESSignatureParameters) error
}

// JAdESLevelBaselineT creates a T-level of a JAdES signature.
type JAdESLevelBaselineT struct {
	JAdESExtensionBuilder

	// CertificateVerifier is the CertificateVerifier to use. Port of the protected final
	// `certificateVerifier`.
	CertificateVerifier validation.CertificateVerifier

	// TspSource is the object encapsulating the Time Stamp Protocol needed to create the
	// level -T of the signature. Port of the protected `tspSource`.
	TspSource validation.TSPSource

	// DocumentAnalyzer is the cached instance of a document analyzer. Port of the protected
	// `documentAnalyzer`.
	DocumentAnalyzer *AbstractJWSDocumentAnalyzer

	// operationKind defines the current signing procedure (used in signature
	// creation/extension). Port of the private `operationKind`.
	operationKind enumerations.SigningOperation

	// overrides points back at the most-derived level; see InitJAdESLevelBaselineT.
	overrides JAdESSignatureExtensionOverrides
}

// NewJAdESLevelBaselineT is the default constructor.
// Port of JAdESLevelBaselineT(CertificateVerifier).
func NewJAdESLevelBaselineT(certificateVerifier validation.CertificateVerifier) *JAdESLevelBaselineT {
	extension := &JAdESLevelBaselineT{}
	extension.InitJAdESLevelBaselineT(extension, certificateVerifier)
	return extension
}

// InitJAdESLevelBaselineT registers the concrete extension level with this base and stores the
// CertificateVerifier. Port of the JAdESLevelBaselineT(CertificateVerifier) constructor body.
func (t *JAdESLevelBaselineT) InitJAdESLevelBaselineT(self JAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	t.overrides = self
	t.CertificateVerifier = certificateVerifier
}

// SetTspSource sets the TSP source to be used when extending the digital signature.
// Port of #setTspSource.
func (t *JAdESLevelBaselineT) SetTspSource(tspSource validation.TSPSource) {
	t.TspSource = tspSource
}

// SetOperationKind sets the signing operation. Port of the overridden #setOperationKind.
func (t *JAdESLevelBaselineT) SetOperationKind(signingOperation enumerations.SigningOperation) {
	t.operationKind = signingOperation
}

// ExtendSignaturesDocument extends every signature of the given document to the level of the
// concrete implementation.
// Port of the overridden #extendSignatures(DSSDocument, JAdESSignatureParameters).
func (t *JAdESLevelBaselineT) ExtendSignaturesDocument(doc model.DSSDocument,
	params *JAdESSignatureParameters) (model.DSSDocument, error) {
	if doc == nil {
		panic("The document cannot be null")
	}
	if t.TspSource == nil {
		panic("The TSPSource cannot be null")
	}

	documentAnalyzerFactory := NewJWSDocumentAnalyzerFactory()
	t.DocumentAnalyzer = jwsDocumentAnalyzerBase(documentAnalyzerFactory.Create(doc))
	t.DocumentAnalyzer.SetCertificateVerifier(t.CertificateVerifier)
	t.DocumentAnalyzer.SetDetachedContents(params.DetachedContents())
	t.DocumentAnalyzer.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)

	jwsJsonSerializationObject := t.DocumentAnalyzer.JwsJsonSerializationObject()
	if err := t.AssertJWSJsonSerializationObjectValid(jwsJsonSerializationObject); err != nil {
		return nil, err
	}

	signatures := t.DocumentAnalyzer.Signatures()
	if utils.IsCollectionEmpty(signatures) {
		return nil, exception.NewIllegalInputException("No signatures found to be extended!")
	}

	signaturesToExtend := signatures
	// this method allows extension of only the current signature on creation
	if enumerations.SigningOperationSign == t.operationKind {
		signaturesToExtend = []validation.AdvancedSignature{signatures[len(signatures)-1]}
	}

	if err := t.overrides.ExtendSignatures(signaturesToExtend, params); err != nil {
		return nil, err
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject, params.JwsSerializationType())
	return generator.Generate()
}

// ExtendSignatures extends the signatures to the -T level: for every signature that still needs
// it, a signature time-stamp is obtained from the TSP source over the signature timestamp data
// and appended to the 'etsiU' unsigned header as a 'sigTst' component.
// Port of the protected #extendSignatures(List, JAdESSignatureParameters).
func (t *JAdESLevelBaselineT) ExtendSignatures(signatures []validation.AdvancedSignature,
	params *JAdESSignatureParameters) error {
	signaturesToExtend := t.extendToTLevelSignatures(signatures, params)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	signatureRequirementsChecker := t.SignatureRequirementsChecker(params)
	signatureRequirementsChecker.AssertExtendToTLevelPossible(signaturesToExtend)

	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertSigningCertificatesAreValid(signaturesToExtend)

	for _, signature := range signaturesToExtend {
		jadesSignature, ok := signature.(*JAdESSignature)
		if !ok {
			// Java's (JAdESSignature) cast; a non-JAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}

		if err := t.AssertEtsiUComponentsConsistent(jadesSignature.Jws(), params); err != nil {
			return err
		}

		signatureTimestampParameters := params.GetSignatureTimestampParameters()
		timestampDigestAlgorithm := signatureTimestampParameters.DigestAlgorithm()

		timestampSource, err := jadesLevelBaselineTTimestampSource(jadesSignature.TimestampSource())
		if err != nil {
			return err
		}
		messageDigest := timestampSource.GetSignatureTimestampData(timestampDigestAlgorithm)
		timeStampResponse, err := t.TspSource.TimeStampResponse(timestampDigestAlgorithm,
			messageDigest.Value())
		if err != nil {
			return err
		}

		tstContainer, err := DSSJsonUtilsTstContainer([]*model.TimestampBinary{timeStampResponse}, "")
		if err != nil {
			return err
		}

		etsiUHeader := jadesSignature.EtsiUHeader()
		if err := etsiUHeader.AddComponent(JAdESHeaderParameterNamesSigTst, tstContainer,
			utils.IsTrue(params.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return err
		}
	}
	return nil
}

// SignatureRequirementsChecker instantiates a SignatureRequirementsChecker.
// Port of the protected #getSignatureRequirementsChecker.
func (t *JAdESLevelBaselineT) SignatureRequirementsChecker(
	parameters *JAdESSignatureParameters) *document.SignatureRequirementsChecker[*JAdESTimestampParameters] {
	return document.NewSignatureRequirementsChecker[*JAdESTimestampParameters](t.CertificateVerifier,
		&parameters.AbstractSignatureParameters)
}

// extendToTLevelSignatures ports the private getExtendToTLevelSignatures.
func (t *JAdESLevelBaselineT) extendToTLevelSignatures(signatures []validation.AdvancedSignature,
	parameters *JAdESSignatureParameters) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if jadesLevelBaselineTTLevelExtensionRequired(signature, parameters) {
			toBeExtended = append(toBeExtended, signature)
		}
	}
	return toBeExtended
}

// jadesLevelBaselineTTLevelExtensionRequired ports the private tLevelExtensionRequired.
func jadesLevelBaselineTTLevelExtensionRequired(jadesSignature validation.AdvancedSignature,
	parameters *JAdESSignatureParameters) bool {
	return enumerations.SignatureLevelJAdESBaselineT == parameters.SignatureLevel() ||
		!jadesSignature.HasTProfile()
}

// jadesLevelBaselineTTimestampSource narrows the signature's timestamp source to the JAdES one.
// Java's JAdESSignature#getTimestampSource() is a covariant override returning
// JAdESTimestampSource, which Go cannot express: AdvancedSignature already declares
// TimestampSource() validation.TimestampSource, so the JAdES source is recovered by assertion
// here - the same technique, and the same reason, as xades.xadesLevelBaselineTTimestampSource.
// jades_level_baseline_lta.go reuses this helper rather than duplicating the assertion.
func jadesLevelBaselineTTimestampSource(source any) (*JAdESTimestampSource, error) {
	if timestampSource, ok := source.(*JAdESTimestampSource); ok {
		return timestampSource, nil
	}
	return nil, fmt.Errorf("the signature timestamp source is not a JAdESTimestampSource, but %T", source)
}

// Compile-time assertion that *JAdESLevelBaselineT satisfies the extension contract.
var _ JAdESLevelBaselineExtension = (*JAdESLevelBaselineT)(nil)
