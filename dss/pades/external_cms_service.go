// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/ExternalCMSService.java (DSS 6.5.RC1).
//
// The public entry points (GetDataToSign, SignMessageDigest) turn the (T, error) of the layers
// below back into Java's propagating unchecked exception, i.e. into a panic carrying the error -
// the convention cades/cades_service.go established for a service class. The protected builders
// keep their error channel, because Service calls them directly.
//
// org.bouncycastle's ContentSigner is cms.ContentSigner (PORTING.md); slf4j is dropped.
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// ExternalCMSService generates a CMS signed data to be incorporated within a PDF document for a
// PAdES signature creation.
//
// To create a CMS with this service, follow the algorithm:
//  1. compute the DTBS from the message-digest of the PDF signature's ByteRange:
//     toBeSigned := service.GetDataToSign(messageDigest, parameters);
//  2. create the signature value by private-key signing of toBeSigned;
//  3. create the CMS signature signing the message-digest:
//     cmsSignature := service.SignMessageDigest(messageDigest, parameters, signatureValue).
//
// NOTE: this class does not create CAdES-BASELINE signatures, but CAdES-Extended signatures as
// per ETSI EN 319 122-2, suitable for a PAdES-BASELINE creation.
type ExternalCMSService struct {
	// certificateVerifier is used for a certificate chain validation.
	certificateVerifier validation.CertificateVerifier

	// tspSource is used for timestamp requests.
	tspSource validation.TSPSource

	// ResourcesHandlerBuilder writes a created CMS into a defined implementation of an
	// OutputStream or a DSSDocument.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

// NewExternalCMSService is the default constructor. The certificateVerifier provides information
// on the sources to be used in the validation process in the context of a signature.
// Port of ExternalCMSService(CertificateVerifier).
func NewExternalCMSService(certificateVerifier validation.CertificateVerifier) *ExternalCMSService {
	return &ExternalCMSService{
		certificateVerifier:     certificateVerifier,
		ResourcesHandlerBuilder: PAdESUtilsDefaultResourcesHandlerBuilder,
	}
}

// SetTspSource defines the TSP (timestamp provider) source for a T-level signature creation.
// Port of #setTspSource.
func (s *ExternalCMSService) SetTspSource(tspSource validation.TSPSource) {
	s.tspSource = tspSource
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with CMS
// object output containers during the signature creation procedure.
// NOTE: the DSSResourcesHandlerBuilder is supported only within the 'dss-cms-stream' module!
// Port of #setResourcesHandlerBuilder.
func (s *ExternalCMSService) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	s.ResourcesHandlerBuilder = cms.UtilsResourcesHandlerBuilder(resourcesHandlerBuilder)
}

// GetDataToSign computes the signed attributes of a CMS signed data to be used for a private-key
// signing. Port of #getDataToSign.
func (s *ExternalCMSService) GetDataToSign(messageDigest model.DSSMessageDigest,
	parameters *SignatureParameters) *model.ToBeSigned {
	if messageDigest.Value() == nil {
		// Java's Objects.requireNonNull(messageDigest); a DSSMessageDigest is a value type in
		// Go, so the closest analogue of "null" is the zero value, i.e. no digest value.
		panic("messageDigest cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertConfigurationValid(messageDigest, parameters)

	toBeSigned, err := s.BuildToBeSignedData(messageDigest, parameters)
	if err != nil {
		panic(err)
	}
	return toBeSigned
}

// BuildToBeSignedData builds the data to be signed without executing additional checks on the
// provided configuration. Port of the protected #buildToBeSignedData.
func (s *ExternalCMSService) BuildToBeSignedData(messageDigest model.DSSMessageDigest,
	parameters *SignatureParameters) (*model.ToBeSigned, error) {
	signatureAlgorithm := parameters.SignatureAlgorithm()
	customContentSigner, err := cms.NewCustomContentSignerBuilder().Build(signatureAlgorithm)
	if err != nil {
		return nil, err
	}

	cmsBuilderHelper := s.InitCMSBuilderHelper(messageDigest, parameters, customContentSigner)
	if _, err = cmsBuilderHelper.CreateCMS(); err != nil {
		return nil, err
	}

	return model.NewToBeSignedWithBytes(customContentSigner.OutputStream().Bytes()), nil
}

// SignMessageDigest creates a signed CMS to be incorporated within a PDF document for a PAdES
// signature creation. Port of #signMessageDigest.
func (s *ExternalCMSService) SignMessageDigest(messageDigest model.DSSMessageDigest,
	parameters *SignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if messageDigest.Value() == nil {
		// Java's Objects.requireNonNull(messageDigest); a DSSMessageDigest is a value type in
		// Go, so the closest analogue of "null" is the zero value, i.e. no digest value.
		panic("messageDigest cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}
	s.AssertConfigurationValid(messageDigest, parameters)

	signedCMS, err := s.BuildCMS(messageDigest, parameters, signatureValue)
	if err != nil {
		panic(err)
	}
	parameters.Reinit()
	signatureDocument, err := cms.UtilsWriteToDSSDocument(signedCMS, s.ResourcesHandlerBuilder)
	if err != nil {
		panic(err)
	}
	return signatureDocument
}

// BuildCMS builds a CMS without executing additional checks on the provided configuration.
// Port of the protected #buildCMS.
func (s *ExternalCMSService) BuildCMS(messageDigest model.DSSMessageDigest,
	parameters *SignatureParameters, signatureValue *model.SignatureValue) (*cms.CMS, error) {
	signatureAlgorithm := parameters.SignatureAlgorithm()
	signatureLevel := parameters.SignatureLevel()
	if signatureAlgorithm == "" {
		panic("SignatureAlgorithm cannot be null!")
	}
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}

	signatureValue, err := document.NewSignatureValueChecker().EnsureSignatureValue(signatureValue,
		parameters.SignatureAlgorithm())
	if err != nil {
		return nil, err
	}
	customContentSigner, err := cms.NewCustomContentSignerBuilder().BuildWithSignatureValue(
		signatureAlgorithm, signatureValue)
	if err != nil {
		return nil, err
	}

	cmsBuilderHelper := s.InitCMSBuilderHelper(messageDigest, parameters, customContentSigner)
	signedCMS, err := cmsBuilderHelper.CreateCMS()
	if err != nil {
		return nil, err
	}

	if enumerations.SignatureLevelPAdESBaselineB != signatureLevel {
		if s.tspSource == nil {
			panic("TSPSource shall be provided for T-level creation!")
		}
		digestDocument := spi.DSSUtilsToDigestDocument(messageDigest.Digest)
		parameters.GetContext().SetDetachedContents([]model.DSSDocument{digestDocument})

		cadesLevelBaselineT := cades.NewCAdESLevelBaselineT(s.tspSource, s.certificateVerifier)
		if signedCMS, err = cadesLevelBaselineT.ExtendCMSSignatures(signedCMS,
			&parameters.SignatureParameters); err != nil {
			return nil, err
		}
	}
	return signedCMS, nil
}

// AssertConfigurationValid verifies whether the provided parameters are valid for the external
// CMS creation process. Port of the protected #assertConfigurationValid; Java's unchecked
// IllegalArgumentException is a panic here, as the method returns nothing upstream.
func (s *ExternalCMSService) AssertConfigurationValid(messageDigest model.DSSMessageDigest,
	parameters *SignatureParameters) {
	signatureLevel := parameters.SignatureLevel()
	if signatureLevel == "" {
		panic("SignatureLevel shall be defined!")
	}
	if enumerations.SignatureLevelPAdESBaselineB != signatureLevel &&
		enumerations.SignatureLevelPAdESBaselineT != signatureLevel {
		panic(fmt.Sprintf("SignatureLevel '%s' is not supported within PAdESCMSGeneratorService!",
			signatureLevel))
	}
	s.AssertSigningCertificateValid(parameters)
	if messageDigest.Algorithm() != parameters.DigestAlgorithm() {
		panic(fmt.Sprintf("The DigestAlgorithm provided within Digest '%s' does not correspond "+
			"to the one defined in SignatureParameters '%s'!",
			messageDigest.Algorithm(), parameters.DigestAlgorithm()))
	}
}

// AssertSigningCertificateValid raises an exception if the signing rules forbid the use of the
// certificate. Port of the protected #assertSigningCertificateValid.
func (s *ExternalCMSService) AssertSigningCertificateValid(parameters *SignatureParameters) {
	signingCertificate := parameters.SigningCertificate()
	if signingCertificate == nil {
		if parameters.GenerateTBSWithoutCertificate() {
			return
		}
		panic("Signing Certificate is not defined! " +
			"Set signing certificate or use method setGenerateTBSWithoutCertificate(true).")
	}

	signatureRequirementsChecker := document.NewSignatureRequirementsChecker[*cades.TimestampParameters](
		s.certificateVerifier, &parameters.AbstractSignatureParameters)
	signatureRequirementsChecker.AssertSigningCertificateIsValid(signingCertificate)
}

// InitCMSBuilderHelper instantiates a CMSForPAdESBuilderHelper.
// Port of the protected #initCMSBuilderHelper.
func (s *ExternalCMSService) InitCMSBuilderHelper(messageDigest model.DSSMessageDigest,
	signatureParameters *SignatureParameters, contentSigner cms.ContentSigner) *CMSForPAdESBuilderHelper {
	return NewCMSForPAdESBuilderHelper(messageDigest, signatureParameters, contentSigner).
		SetTrustedCertificateSource(s.certificateVerifier.TrustedCertSources())
}
