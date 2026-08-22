// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESAnalyzer.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
//
// NOT gated behind `phase8`, unlike its CAdES sibling. The CAdES container analyzer is forced to
// be: asic/cades/asic_container_with_cades_analyzer.go needs ASiCWithCAdESTimestampAnalyzer, which
// extends dss/validation's DetachedTimestampAnalyzer - a package Phase 8 has not landed, so that
// file cannot compile at all today. Nothing in the XAdES container analyzer reaches outside the
// already-ported tree: asic.AbstractASiCContainerAnalyzer, spi/validation{,/analyzer} and the
// frozen dss/xades analyzer are all live, and the type is directly usable through the ported
// analyzer.DocumentAnalyzer interface. Gating it too would have hidden a working surface -
// per-signature analysis of ASiC containers - from the build and from cross-validation, so it is
// live and pinned against upstream by asic/broad_corpus_cross_validation_test.go's
// TestBroadCorpusXAdESSignatureAnalysisMatchesUpstream over the whole fixture corpus. What
// genuinely awaits Phase 8 is the *validator* wrapper (asic_container_with_xades_validator.go),
// which does import dss/validation.
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
// ASiCContainerWithXAdESAnalyzer(ASiCContent).
func NewASiCContainerWithXAdESAnalyzerFromContent(asicContent *asic.ASiCContent) *ASiCContainerWithXAdESAnalyzer {
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
func (a *ASiCContainerWithXAdESAnalyzer) IsSupportedASiCContent(asicContent *asic.ASiCContent) bool {
	return NewASiCWithXAdESFormatDetector().IsSupportedASiCContent(asicContent)
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithXAdESAnalyzer) GetContainerExtractor() *asic.DefaultASiCContainerExtractor {
	return &NewASiCWithXAdESContainerExtractor(a.Document()).DefaultASiCContainerExtractor
}

// GetSignatureAnalyzers ports the @Override protected getSignatureAnalyzers(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
//
// FLAGGED CROSS-CHUNK GAP (mirrors asic/cades/asic_container_with_cades_analyzer.go's identical
// note): Java's getSignatureAnalyzers() forwards `this.getSignaturePolicyProvider()` (a protected
// accessor on the frozen analyzer.DefaultDocumentAnalyzer, unexported in the Go port as
// signaturePolicyProviderOrDefault and not reachable from another package) into each nested
// XMLDocumentAnalyzer. No exported equivalent exists today, so this propagation is dropped here:
// each nested XMLDocumentAnalyzer instead lazily instantiates its own default
// SignaturePolicyProvider. This only differs observably when a caller has set a *custom*
// SignaturePolicyProvider on the outer analyzer via SetSignaturePolicyProvider.
func (a *ASiCContainerWithXAdESAnalyzer) GetSignatureAnalyzers() []analyzer.DocumentAnalyzer {
	if a.SignatureValidators == nil {
		a.SignatureValidators = make([]analyzer.DocumentAnalyzer, 0)
		for _, signature := range a.GetSignatureDocuments() {
			documentAnalyzer, err := dssxades.NewXMLDocumentAnalyzer(signature)
			if err != nil {
				panic(err)
			}
			documentAnalyzer.SetCertificateVerifier(a.CertificateVerifier())

			isOpenDocument, err := asic.ASiCUtilsIsOpenDocument(a.GetMimeTypeDocument())
			if err == nil && isOpenDocument {
				documentAnalyzer.SetDetachedContents(OpenDocumentSupportUtilsGetOpenDocumentCoverage(a.AsicContent))
			} else if enumerations.ASiCContainerType_ASiC_S == a.GetContainerType() {
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
		manifestFile := asic.ASiCManifestParserGetManifestFile(manifestDocument)
		if manifestFile != nil {
			manifestFile.SetManifestType(enumerations.ASiCManifestTypeEnum_EVIDENCE_RECORD)
			manifestValidator := asic.NewASiCManifestValidator(manifestFile, a.GetAllDocuments())
			manifestValidator.ValidateEntries()
			descriptions = append(descriptions, manifestFile)
		}
	}

	return descriptions
}

// OriginalDocumentsForSignature ports the @Override getOriginalDocuments(AdvancedSignature).
func (a *ASiCContainerWithXAdESAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	xadesSignature := advancedSignature.(*dssxades.XAdESSignature)
	retrievedDocs := dssxades.XAdESSignatureUtilsGetSignerDocuments(xadesSignature)
	if utils.IsCollectionNotEmpty(retrievedDocs) {
		return a.extractArchiveDocuments(retrievedDocs)
	}
	return []model.DSSDocument{}
}

func (a *ASiCContainerWithXAdESAnalyzer) extractArchiveDocuments(retrievedDocs []model.DSSDocument) []model.DSSDocument {
	if enumerations.ASiCContainerType_ASiC_S == a.GetContainerType() {
		return a.GetSignedDocumentsASiCS(retrievedDocs)
	}
	return retrievedDocs
}
