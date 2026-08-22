// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESLevelBaselineLTA.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithCAdESLevelBaselineLTA extends ASiC with CAdES signatures to LTA-level.
//
// Java's `extends ASiCWithCAdESSignatureExtension` becomes embedding plus the
// InitASiCWithCAdESSignatureExtension(self) registration, so the base dispatches
// extensionRequired / assertExtendSignaturePossible / getLTAExtensionProfile into this type.
type ASiCWithCAdESLevelBaselineLTA struct {
	*ASiCWithCAdESSignatureExtension

	// asicFilenameFactory defines rules for filename creation for timestamp and archive
	// manifest files.
	asicFilenameFactory ASiCWithCAdESFilenameFactory
}

var (
	_ ASiCWithCAdESSignatureExtensionOverrides = (*ASiCWithCAdESLevelBaselineLTA)(nil)
	_ ASiCWithCAdESSignatureExtender           = (*ASiCWithCAdESLevelBaselineLTA)(nil)
)

// NewASiCWithCAdESLevelBaselineLTA is the default constructor. Ports
// ASiCWithCAdESLevelBaselineLTA(CertificateVerifier, TSPSource).
func NewASiCWithCAdESLevelBaselineLTA(certificateVerifier validation.CertificateVerifier,
	tspSource validation.TSPSource) *ASiCWithCAdESLevelBaselineLTA {
	return NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory(certificateVerifier, tspSource,
		NewDefaultASiCWithCAdESFilenameFactory())
}

// NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory is the constructor with filename factory.
// Ports ASiCWithCAdESLevelBaselineLTA(CertificateVerifier, TSPSource, ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory(certificateVerifier validation.CertificateVerifier,
	tspSource validation.TSPSource,
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESLevelBaselineLTA {
	extension := &ASiCWithCAdESLevelBaselineLTA{
		ASiCWithCAdESSignatureExtension: NewASiCWithCAdESSignatureExtension(certificateVerifier, tspSource),
		asicFilenameFactory:             asicFilenameFactory,
	}
	extension.InitASiCWithCAdESSignatureExtension(extension)
	return extension
}

// Extend ports the @Override extend(ASiCContent, CAdESSignatureParameters).
func (e *ASiCWithCAdESLevelBaselineLTA) Extend(asicContent *asic.ASiCContent,
	parameters *dsscades.CAdESSignatureParameters) *asic.ASiCContent {
	// LT-level extension, if required
	asicContent = e.ASiCWithCAdESSignatureExtension.Extend(asicContent, parameters)
	// LTA-level extension
	return e.extendWithDigestAlgorithms(asicContent, e.GetReferenceDigestAlgorithmOrDefault(parameters),
		parameters.GetArchiveTimestampParameters().DigestAlgorithm())
}

// ExtendWithDigestAlgorithm extends the asicContent with an ArchiveManifest timestamp.
//
// NOTE: This method is to be used for a direct timestamping with an ArchiveManifest, without
// in-depth signature attributes (the signature extension is still applied). Use
// Extend(ASiCContent, *CAdESSignatureParameters) for a proper signature(s) extension.
//
// Ports the public extend(ASiCContent, DigestAlgorithm); Go has no overloading, so the Java
// overload set extend(ASiCContent, CAdESSignatureParameters) / extend(ASiCContent,
// DigestAlgorithm) becomes Extend / ExtendWithDigestAlgorithm.
func (e *ASiCWithCAdESLevelBaselineLTA) ExtendWithDigestAlgorithm(asicContent *asic.ASiCContent,
	digestAlgorithm enumerations.DigestAlgorithm) *asic.ASiCContent {
	// ensure the signatures are extended to LT-level, when necessary
	asicContent = e.ASiCWithCAdESSignatureExtension.Extend(asicContent, e.getEmptyLTLevelSignatureParameters())
	return e.extendWithDigestAlgorithms(asicContent, digestAlgorithm, digestAlgorithm)
}

// extendWithDigestAlgorithms extends the ASiC Container by adding a new Archive Manifest,
// time-stamp file and the necessary validation data. Ports the private
// extend(ASiCContent, DigestAlgorithm, DigestAlgorithm).
func (e *ASiCWithCAdESLevelBaselineLTA) extendWithDigestAlgorithms(asicContent *asic.ASiCContent,
	manifestDigestAlgorithm, tstDigestAlgorithm enumerations.DigestAlgorithm) *asic.ASiCContent {
	// shall be computed on the first step, before timestamp extension/creation
	timestampFilename := e.asicFilenameFactory.TimestampFilename(asicContent)

	lastManifestFile := e.getLastManifestFile(asicContent.AllManifestDocuments())

	timestampDocuments := asicContent.TimestampDocuments()
	lastTimestamp := e.getLastTimestampDocument(lastManifestFile, timestampDocuments)
	if lastTimestamp != nil {
		coveredByAnyManifest := e.IsCoveredByAnyManifest(asicContent, lastTimestamp)
		e.AssertExtendTimestampPossible(coveredByAnyManifest)

		// Extending an existing archive timestamp needs to re-validate the container's
		// signatures and detached timestamps to gather fresh validation data - that requires
		// the dss/validation engine (Phase 8, not yet ported). See
		// asic_with_cades_lta_validation_phase8.go / _nophase8.go for the two build-tagged
		// implementations of this step.
		extendedTimestamp := e.extendLastArchiveTimestampWithValidationData(asicContent, lastTimestamp)
		asicContent.SetTimestampDocuments(asic.ASiCUtilsAddOrReplaceDocument(asicContent.TimestampDocuments(), extendedTimestamp))
	}

	var lastArchiveManifest model.DSSDocument
	if lastManifestFile != nil && e.isLastArchiveManifest(lastManifestFile.Filename()) {
		lastArchiveManifest = lastManifestFile.Document()
		lastArchiveManifest.SetName(e.asicFilenameFactory.ArchiveManifestFilename(asicContent))
	}

	builder := NewASiCEWithCAdESArchiveManifestBuilder(
		asicContent, lastArchiveManifest, manifestDigestAlgorithm, timestampFilename)
	archiveManifest, err := builder.Build()
	if err != nil {
		panic(err)
	}
	asicContent.SetArchiveManifestDocuments(append(asicContent.ArchiveManifestDocuments(), archiveManifest))

	archiveManifestDigest, err := spi.DSSUtilsDigestOfDocument(tstDigestAlgorithm, archiveManifest)
	if err != nil {
		panic(err)
	}
	timeStampResponse, err := e.TspSource.TimeStampResponse(tstDigestAlgorithm, archiveManifestDigest)
	if err != nil {
		panic(err)
	}
	derEncoded, err := spi.DSSASN1UtilsDEREncodedTimestampBinary(timeStampResponse)
	if err != nil {
		panic(err)
	}
	timestamp := model.NewInMemoryDocumentWithMimeType(derEncoded, timestampFilename, enumerations.MimeTypeEnum_TST)
	asicContent.SetTimestampDocuments(append(asicContent.TimestampDocuments(), timestamp))

	return asicContent
}

// getLastManifestFile ports the private getLastManifestFile(List).
func (e *ASiCWithCAdESLevelBaselineLTA) getLastManifestFile(manifests []model.DSSDocument) *model.ManifestFile {
	lastManifest := e.getLastArchiveManifest(manifests)
	if lastManifest == nil {
		lastManifest = spi.DSSUtilsDocumentWithLastName(manifests)
	}
	if lastManifest != nil {
		return asic.ASiCManifestParserGetManifestFile(lastManifest)
	}
	return nil
}

// getLastArchiveManifest ports the private getLastArchiveManifest(List).
func (e *ASiCWithCAdESLevelBaselineLTA) getLastArchiveManifest(manifests []model.DSSDocument) model.DSSDocument {
	if utils.IsCollectionNotEmpty(manifests) {
		for _, manifest := range manifests {
			if e.isLastArchiveManifest(manifest.Name()) {
				return manifest
			}
		}
	}
	return nil
}

// isLastArchiveManifest ports the private isLastArchiveManifest(String).
func (e *ASiCWithCAdESLevelBaselineLTA) isLastArchiveManifest(fileName string) bool {
	return ASiCWithCAdESUtilsDefaultArchiveManifestFilename == fileName
}

// getLastTimestampDocument ports the private getLastTimestampDocument(ManifestFile, List).
func (e *ASiCWithCAdESLevelBaselineLTA) getLastTimestampDocument(lastManifestFile *model.ManifestFile,
	timestamps []model.DSSDocument) model.DSSDocument {
	if lastManifestFile != nil {
		return spi.DSSUtilsDocumentWithName(timestamps, lastManifestFile.SignatureFilename())
	}
	return spi.DSSUtilsDocumentWithLastName(timestamps)
}

// extendTimestamp ports the private extendTimestamp(DSSDocument, ValidationData).
func (e *ASiCWithCAdESLevelBaselineLTA) extendTimestamp(archiveTimestamp model.DSSDocument,
	validationDataForInclusion *validation.ValidationData) model.DSSDocument {
	timestampCMS, err := cms.CMSUtilsParseToCMS(archiveTimestamp)
	if err != nil {
		panic(err)
	}
	cmsBuilder := cms.NewCMSBuilder().SetOriginalCMS(timestampCMS)
	extendedCMS, err := cmsBuilder.ExtendCMSSignedData(
		validationDataForInclusion.CertificateTokens(), validationDataForInclusion.CrlTokens(),
		validationDataForInclusion.OcspTokens())
	if err != nil {
		panic(err)
	}
	return model.NewInMemoryDocumentWithMimeType(extendedCMS.DEREncoded(), archiveTimestamp.Name(), enumerations.MimeTypeEnum_TST)
}

// getEmptyLTLevelSignatureParameters ports the private getEmptyLTLevelSignatureParameters().
func (e *ASiCWithCAdESLevelBaselineLTA) getEmptyLTLevelSignatureParameters() *dsscades.CAdESSignatureParameters {
	parameters := dsscades.NewCAdESSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevel_CAdES_BASELINE_LT)
	return parameters
}

// ExtensionRequired ports the @Override protected
// extensionRequired(CAdESSignatureParameters, boolean).
func (e *ASiCWithCAdESLevelBaselineLTA) ExtensionRequired(parameters *dsscades.CAdESSignatureParameters,
	coveredByManifest bool) bool {
	return !coveredByManifest
}

// GetLTAExtensionProfile ports the @Override protected
// getLTAExtensionProfile(TSPSource, CertificateVerifier).
func (e *ASiCWithCAdESLevelBaselineLTA) GetLTAExtensionProfile(tspSource validation.TSPSource,
	certificateVerifier validation.CertificateVerifier) dsscades.CAdESSignatureExtender {
	return dsscades.NewCAdESLevelBaselineLT(tspSource, certificateVerifier)
}

// AssertExtendSignaturePossible ports the @Override protected
// assertExtendSignaturePossible(CAdESSignatureParameters, boolean).
//
// Panics with an *exception.IllegalInputException when the signature is already covered by a
// manifest file.
func (e *ASiCWithCAdESLevelBaselineLTA) AssertExtendSignaturePossible(
	parameters *dsscades.CAdESSignatureParameters, coveredByManifest bool) {
	if coveredByManifest {
		panic(exception.NewIllegalInputException(
			"Cannot extend signature to '" + string(enumerations.SignatureLevel_CAdES_BASELINE_LTA) +
				"'. The signature is already covered by a manifest file."))
	}
}

// AssertExtendTimestampPossible checks if the timestamp extension is possible. Ports the
// protected assertExtendTimestampPossible(boolean).
//
// Panics with an *exception.IllegalInputException when the timestamp is already covered by a
// manifest file.
func (e *ASiCWithCAdESLevelBaselineLTA) AssertExtendTimestampPossible(coveredByManifest bool) {
	if coveredByManifest {
		panic(exception.NewIllegalInputException(
			"Cannot extend the last timestamp. The timestamp is already covered by a manifest file."))
	}
}

// asicCAdESLTARevocationIdentifiers narrows a list of revocation binaries to the Identifiers
// ValidationData#excludeCRLTokens / #excludeOCSPTokens take, which is what Java's
// `Collection<? extends Identifier>` parameter accepts directly.
func asicCAdESLTARevocationIdentifiers[R revocation.Revocation](
	binaries []spi.EncapsulatedRevocationTokenIdentifier[R]) []model.Identifier {
	identifiers := make([]model.Identifier, 0, len(binaries))
	for _, binary := range binaries {
		identifiers = append(identifiers, binary.DSSID())
	}
	return identifiers
}
