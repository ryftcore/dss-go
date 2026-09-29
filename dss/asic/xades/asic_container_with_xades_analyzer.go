// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESAnalyzer.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.validation lands in this same Go
// package (dss/asic/xades).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/utils"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
)

// ASiCContainerWithXAdESAnalyzer is an implementation to validate ASiC containers with XAdES
// signature(s).
type ASiCContainerWithXAdESAnalyzer struct {
	*asic.AbstractASiCContainerAnalyzer
}

var _ asic.AbstractASiCContainerAnalyzerOverrides = (*ASiCContainerWithXAdESAnalyzer)(nil)
var _ analyzer.DocumentAnalyzer = (*ASiCContainerWithXAdESAnalyzer)(nil)

// newASiCContainerWithXAdESAnalyzer is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithXAdESAnalyzer() *ASiCContainerWithXAdESAnalyzer {
	a := &ASiCContainerWithXAdESAnalyzer{
		AbstractASiCContainerAnalyzer: asic.NewAbstractASiCContainerAnalyzerBase(),
	}
	a.InitAbstractASiCContainerAnalyzer(a)
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

// NewASiCContainerWithXAdESAnalyzer is the default constructor. Ports
// ASiCContainerWithXAdESAnalyzer(DSSDocument).
func NewASiCContainerWithXAdESAnalyzer(asicContainer model.DSSDocument) *ASiCContainerWithXAdESAnalyzer {
	a := newASiCContainerWithXAdESAnalyzer()
	a.InitFromDocument(asicContainer)
	return a
}

// NewASiCContainerWithXAdESAnalyzerFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithXAdESAnalyzer(Content).
func NewASiCContainerWithXAdESAnalyzerFromContent(asicContent *asic.Content) *ASiCContainerWithXAdESAnalyzer {
	a := newASiCContainerWithXAdESAnalyzer()
	a.InitFromContent(asicContent)
	return a
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (a *ASiCContainerWithXAdESAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	return NewASiCWithXAdESFormatDetector().IsSupportedASiC(dssDocument)
}

// IsSupportedASiCContent ports the @Override isSupported(ASiCContent), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithXAdESAnalyzer) IsSupportedASiCContent(asicContent *asic.Content) bool {
	return NewASiCWithXAdESFormatDetector().IsSupportedASiCContent(asicContent)
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithXAdESAnalyzer) GetContainerExtractor() *asic.DefaultContainerExtractor {
	return &NewASiCWithXAdESContainerExtractor(a.Document()).DefaultContainerExtractor
}

// GetSignatureAnalyzers ports the @Override protected getSignatureAnalyzers(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
//
// Java's getSignatureAnalyzers() forwards `this.getSignaturePolicyProvider()` into each nested
// XMLDocumentAnalyzer. That protected getter creates, when none was set, a default provider whose
// NativeHTTPDataLoader downloads the policy from the URL named in the signature - i.e. from the
// (untrusted) container - with no timeout, a side effect plain XAdES validation does not have.
//
// DIVERGENCE, deliberate: ASiCContainerWithXAdESAnalyzer.getSignatureAnalyzers - only a provider
// the caller explicitly set on the outer analyzer (SetSignaturePolicyProvider) is forwarded, so
// a custom provider reaches the nested analyzers as it does upstream, whereas the default
// provider and the container-driven network fetch are not created. Same as
// asic/cades/asic_container_with_cades_analyzer.go; see PORTING.md.
func (a *ASiCContainerWithXAdESAnalyzer) GetSignatureAnalyzers() []analyzer.DocumentAnalyzer {
	if a.SignatureValidators == nil {
		a.SignatureValidators = make([]analyzer.DocumentAnalyzer, 0)
		for _, signature := range a.GetSignatureDocuments() {
			documentAnalyzer, err := dssxades.NewXMLDocumentAnalyzer(signature)
			if err != nil {
				panic(err)
			}
			documentAnalyzer.SetCertificateVerifier(a.CertificateVerifier())
			if signaturePolicyProvider := a.ConfiguredSignaturePolicyProvider(); signaturePolicyProvider != nil {
				documentAnalyzer.SetSignaturePolicyProvider(signaturePolicyProvider)
			}

			isOpenDocument, err := asic.UtilsIsOpenDocument(a.GetMimeTypeDocument())
			if err == nil && isOpenDocument {
				documentAnalyzer.SetDetachedContents(OpenDocumentSupportUtilsGetOpenDocumentCoverage(a.AsicContent))
			} else if enumerations.ASiCContainerTypeASiCS == a.GetContainerType() {
				documentAnalyzer.SetDetachedContents(a.GetSignedDocuments())
				documentAnalyzer.SetContainerContents(a.GetArchiveDocuments())
			} else {
				documentAnalyzer.SetDetachedContents(a.GetAllDocuments())
			}

			a.SignatureValidators = append(a.SignatureValidators, documentAnalyzer)
		}
	}
	return a.SignatureValidators
}

// GetManifestFilesDescriptions ports the @Override protected getManifestFilesDescriptions(),
// implementing asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithXAdESAnalyzer) GetManifestFilesDescriptions() []*model.ManifestFile {
	descriptions := make([]*model.ManifestFile, 0)

	signatureDocuments := a.GetSignatureDocuments()
	manifestDocuments := a.GetManifestDocuments()
	// All signatures use the same file : manifest.xml
	for _, signatureDoc := range signatureDocuments {
		for _, manifestDoc := range manifestDocuments {
			manifestParser := NewASiCEWithXAdESManifestParserWithSignature(signatureDoc, manifestDoc)
			descriptions = append(descriptions, manifestParser.Manifest())
		}
	}

	for _, manifestDocument := range a.GetEvidenceRecordManifestDocuments() {
		manifestFile := asic.ManifestParserGetManifestFile(manifestDocument)
		if manifestFile != nil {
			manifestFile.SetManifestType(enumerations.ASiCManifestTypeEnumEvidenceRecord)
			manifestValidator := asic.NewManifestValidator(manifestFile, a.GetAllDocuments())
			manifestValidator.ValidateEntries()
			descriptions = append(descriptions, manifestFile)
		}
	}

	return descriptions
}

// OriginalDocumentsForSignature ports the @Override getOriginalDocuments(AdvancedSignature).
func (a *ASiCContainerWithXAdESAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	xadesSignature := advancedSignature.(*dssxades.Signature)
	retrievedDocs := dssxades.SignatureUtilsGetSignerDocuments(xadesSignature)
	if utils.IsCollectionNotEmpty(retrievedDocs) {
		return a.extractArchiveDocuments(retrievedDocs)
	}
	return []model.DSSDocument{}
}

func (a *ASiCContainerWithXAdESAnalyzer) extractArchiveDocuments(retrievedDocs []model.DSSDocument) []model.DSSDocument {
	if enumerations.ASiCContainerTypeASiCS == a.GetContainerType() {
		return a.GetSignedDocumentsASiCS(retrievedDocs)
	}
	return retrievedDocs
}
