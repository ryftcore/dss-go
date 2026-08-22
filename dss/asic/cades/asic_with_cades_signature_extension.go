// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESSignatureExtension.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart), as is serialVersionUID.
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// ASiCWithCAdESSignatureExtensionOverrides declares the operations
// ASiCWithCAdESSignatureExtension calls back into virtually and ASiCWithCAdESLevelBaselineLTA
// overrides. Every instance must be registered with
// InitASiCWithCAdESSignatureExtension before use - NewASiCWithCAdESSignatureExtension does it
// for the base itself, and the LTA constructor re-registers with the subclass.
type ASiCWithCAdESSignatureExtensionOverrides interface {
	// ExtensionRequired checks whether the signature extension is required for the particular
	// document. Port of the protected extensionRequired(CAdESSignatureParameters, boolean).
	ExtensionRequired(parameters *dsscades.CAdESSignatureParameters, coveredByManifest bool) bool

	// AssertExtendSignaturePossible checks if the signature extension is possible. Port of the
	// protected assertExtendSignaturePossible(CAdESSignatureParameters, boolean).
	AssertExtendSignaturePossible(parameters *dsscades.CAdESSignatureParameters, coveredByManifest bool)

	// GetLTAExtensionProfile returns the profile required for an LTA-level signature
	// augmentation according to the given container type. Port of the protected
	// getLTAExtensionProfile(TSPSource, CertificateVerifier).
	GetLTAExtensionProfile(tspSource validation.TSPSource,
		certificateVerifier validation.CertificateVerifier) dsscades.CAdESSignatureExtender
}

// ASiCWithCAdESSignatureExtender is the slice of ASiCWithCAdESSignatureExtension that
// ASiCWithCAdESService#getExtensionProfile's local variable is typed with upstream: Java assigns
// either an ASiCWithCAdESSignatureExtension or an ASiCWithCAdESLevelBaselineLTA (which overrides
// extend) to an ASiCWithCAdESSignatureExtension variable, so the Go switch needs an interface
// value. Both types satisfy it.
type ASiCWithCAdESSignatureExtender interface {
	// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating
	// with internal objects during the signature creation procedure. Port of
	// #setResourcesHandlerBuilder.
	SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder)

	// Extend extends signatures within the ASiCContent. Port of
	// #extend(ASiCContent, CAdESSignatureParameters).
	Extend(asicContent *asic.ASiCContent, parameters *dsscades.CAdESSignatureParameters) *asic.ASiCContent
}

// ASiCWithCAdESSignatureExtension extends an ASiC with CAdES signature.
type ASiCWithCAdESSignatureExtension struct {
	// CertificateVerifier is the CertificateVerifier used for a certificate chain validation.
	// Java declares the field protected final.
	CertificateVerifier validation.CertificateVerifier

	// TspSource is the TSPSource to use for timestamp requests. Java declares the field
	// protected final.
	TspSource validation.TSPSource

	// ResourcesHandlerBuilder is used to create data container objects such as an OutputStream
	// or a DSSDocument.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder

	// overrides points back at the concrete extension; see InitASiCWithCAdESSignatureExtension.
	overrides ASiCWithCAdESSignatureExtensionOverrides
}

var (
	_ ASiCWithCAdESSignatureExtensionOverrides = (*ASiCWithCAdESSignatureExtension)(nil)
	_ ASiCWithCAdESSignatureExtender           = (*ASiCWithCAdESSignatureExtension)(nil)
)

// NewASiCWithCAdESSignatureExtension is the default constructor. Ports
// ASiCWithCAdESSignatureExtension(CertificateVerifier, TSPSource).
func NewASiCWithCAdESSignatureExtension(certificateVerifier validation.CertificateVerifier,
	tspSource validation.TSPSource) *ASiCWithCAdESSignatureExtension {
	extension := &ASiCWithCAdESSignatureExtension{
		CertificateVerifier: certificateVerifier,
		TspSource:           tspSource,
	}
	extension.InitASiCWithCAdESSignatureExtension(extension)
	return extension
}

// InitASiCWithCAdESSignatureExtension registers the concrete extension with this base so the
// base can dispatch ExtensionRequired / AssertExtendSignaturePossible / GetLTAExtensionProfile.
func (e *ASiCWithCAdESSignatureExtension) InitASiCWithCAdESSignatureExtension(overrides ASiCWithCAdESSignatureExtensionOverrides) {
	e.overrides = overrides
}

func (e *ASiCWithCAdESSignatureExtension) requireOverrides() ASiCWithCAdESSignatureExtensionOverrides {
	if e.overrides == nil {
		panic("ASiCWithCAdESSignatureExtension was not initialised: the concrete extension must call InitASiCWithCAdESSignatureExtension in its constructor")
	}
	return e.overrides
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure. Ports
// setResourcesHandlerBuilder(DSSResourcesHandlerBuilder).
func (e *ASiCWithCAdESSignatureExtension) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	e.ResourcesHandlerBuilder = resourcesHandlerBuilder
}

// Extend extends signatures within the ASiCContent. Ports
// extend(ASiCContent, CAdESSignatureParameters).
//
// Panics with an *exception.IllegalInputException when the container type cannot be extracted.
func (e *ASiCWithCAdESSignatureExtension) Extend(asicContent *asic.ASiCContent,
	parameters *dsscades.CAdESSignatureParameters) *asic.ASiCContent {
	signatureDocuments := asicContent.SignatureDocuments()

	containerType := asicContent.ContainerType()
	if containerType == "" {
		panic(exception.NewIllegalInputException(
			"The container type of the provided document is not supported or cannot be extracted!"))
	}

	for _, signature := range signatureDocuments {
		coveredByArchiveManifest := e.IsCoveredByArchiveManifest(asicContent, signature)
		if e.requireOverrides().ExtensionRequired(parameters, coveredByArchiveManifest) {
			// to not extend a signature covered by any other manifest
			coveredByAnyManifest := e.IsCoveredByAnyManifest(asicContent, signature)
			e.requireOverrides().AssertExtendSignaturePossible(parameters, coveredByAnyManifest)

			extendedSignature := e.extendSignatureDocument(signature, asicContent, parameters)
			signatureDocuments = asic.ASiCUtilsAddOrReplaceDocument(signatureDocuments, extendedSignature)
			asicContent.SetSignatureDocuments(signatureDocuments)
		}
	}

	return asicContent
}

// extendSignatureDocument ports the private
// extendSignatureDocument(DSSDocument, ASiCContent, CAdESSignatureParameters).
func (e *ASiCWithCAdESSignatureExtension) extendSignatureDocument(signature model.DSSDocument,
	asicContent *asic.ASiCContent, cadesParameters *dsscades.CAdESSignatureParameters) model.DSSDocument {
	detachedContents := e.getDetachedContents(signature, asicContent)
	cadesParameters.GetContext().SetDetachedContents(detachedContents)

	originalName := signature.Name()
	extendDocument, err := e.getExtensionProfile(cadesParameters).ExtendSignatures(signature, cadesParameters)
	if err != nil {
		panic(err)
	}
	extendDocument.SetName(originalName)
	return extendDocument
}

// getDetachedContents ports the private getDetachedContents(DSSDocument, ASiCContent).
//
// NOTE: Java builds Collections.singletonList(linkedManifest) even when the linked manifest was
// not found, i.e. a one-element list holding null; the Go port keeps that shape (a one-element
// slice holding a nil DSSDocument) rather than "fixing" it.
func (e *ASiCWithCAdESSignatureExtension) getDetachedContents(signatureDocument model.DSSDocument,
	asicContent *asic.ASiCContent) []model.DSSDocument {
	if enumerations.ASiCContainerTypeASiCE == asicContent.ContainerType() {
		manifests := asicContent.ManifestDocuments()
		linkedManifest := asic.ASiCManifestParserGetLinkedManifest(manifests, signatureDocument.Name())
		return []model.DSSDocument{linkedManifest}
	}
	return asicContent.SignedDocuments()
}

// GetReferenceDigestAlgorithmOrDefault returns params.referenceDigestAlgorithm if it exists,
// params.digestAlgorithm otherwise. Ports the protected
// getReferenceDigestAlgorithmOrDefault(CAdESSignatureParameters).
func (e *ASiCWithCAdESSignatureExtension) GetReferenceDigestAlgorithmOrDefault(
	params *dsscades.CAdESSignatureParameters) enumerations.DigestAlgorithm {
	if params.ReferenceDigestAlgorithm() != "" {
		return params.ReferenceDigestAlgorithm()
	}
	return params.DigestAlgorithm()
}

// getExtensionProfile returns the extension profile to be used for a CAdES signature
// augmentation. Ports the private getExtensionProfile(CAdESSignatureParameters).
//
// Panics with Java's messages when the signature level is undefined or unsupported.
func (e *ASiCWithCAdESSignatureExtension) getExtensionProfile(
	parameters *dsscades.CAdESSignatureParameters) dsscades.CAdESSignatureExtender {
	signatureLevel := parameters.SignatureLevel()
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}
	var cadesSignatureExtension dsscades.CAdESSignatureExtender
	switch signatureLevel {
	case enumerations.SignatureLevelCAdESBaselineT:
		cadesSignatureExtension = dsscades.NewCAdESLevelBaselineT(e.TspSource, e.CertificateVerifier)
	case enumerations.SignatureLevelCAdESBaselineLT:
		cadesSignatureExtension = dsscades.NewCAdESLevelBaselineLT(e.TspSource, e.CertificateVerifier)
	case enumerations.SignatureLevelCAdESBaselineLTA:
		cadesSignatureExtension = e.requireOverrides().GetLTAExtensionProfile(e.TspSource, e.CertificateVerifier)
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", signatureLevel))
	}
	cadesSignatureExtension.SetResourcesHandlerBuilder(e.ResourcesHandlerBuilder)
	return cadesSignatureExtension
}

// GetLTAExtensionProfile ports the protected getLTAExtensionProfile(TSPSource,
// CertificateVerifier).
func (e *ASiCWithCAdESSignatureExtension) GetLTAExtensionProfile(tspSource validation.TSPSource,
	certificateVerifier validation.CertificateVerifier) dsscades.CAdESSignatureExtender {
	return dsscades.NewCAdESLevelBaselineLTA(tspSource, certificateVerifier)
}

// ExtensionRequired checks whether the signature extension is required for the particular
// document. Ports the protected extensionRequired(CAdESSignatureParameters, boolean).
func (e *ASiCWithCAdESSignatureExtension) ExtensionRequired(parameters *dsscades.CAdESSignatureParameters,
	coveredByManifest bool) bool {
	signatureLevel := parameters.SignatureLevel()
	return enumerations.SignatureLevelCAdESBaselineT == signatureLevel ||
		enumerations.SignatureLevelCAdESBaselineLT == signatureLevel || !coveredByManifest
}

// AssertExtendSignaturePossible checks if the signature extension is possible. Ports the
// protected assertExtendSignaturePossible(CAdESSignatureParameters, boolean).
//
// Panics with an *exception.IllegalInputException when the signature is already covered by a
// manifest file.
func (e *ASiCWithCAdESSignatureExtension) AssertExtendSignaturePossible(
	parameters *dsscades.CAdESSignatureParameters, coveredByManifest bool) {
	signatureLevel := parameters.SignatureLevel()
	if (enumerations.SignatureLevelCAdESBaselineT == signatureLevel ||
		enumerations.SignatureLevelCAdESBaselineLT == signatureLevel) && coveredByManifest {
		panic(exception.NewIllegalInputException(fmt.Sprintf(
			"Cannot extend signature to '%s'. The signature is already covered by a manifest file.", signatureLevel)))
	}
}

// IsCoveredByArchiveManifest verifies whether the signature document is covered by an Archive
// Manifest. Ports the protected isCoveredByArchiveManifest(ASiCContent, DSSDocument).
func (e *ASiCWithCAdESSignatureExtension) IsCoveredByArchiveManifest(asicContent *asic.ASiCContent,
	signature model.DSSDocument) bool {
	return asic.ASiCUtilsIsCoveredByManifest(asicContent.ArchiveManifestDocuments(), signature.Name())
}

// IsCoveredByAnyManifest verifies whether the signature document is covered by any ASiC Manifest
// file. Ports the protected isCoveredByAnyManifest(ASiCContent, DSSDocument).
func (e *ASiCWithCAdESSignatureExtension) IsCoveredByAnyManifest(asicContent *asic.ASiCContent,
	signature model.DSSDocument) bool {
	return asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), signature.Name())
}
