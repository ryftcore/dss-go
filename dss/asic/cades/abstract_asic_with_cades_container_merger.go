// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/merge/AbstractASiCWithCAdESContainerMerger.java (DSS 6.5.RC1).
package cades

import (
	"bytes"
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractASiCWithCAdESContainerMerger contains common code for ASiC with CAdES container merger
// classes.
type AbstractASiCWithCAdESContainerMerger struct {
	*asic.DefaultContainerMerger

	// AsicFilenameFactory defines rules for filename creation for new ZIP entries (e.g.
	// signature files, etc.). Java declares the field protected.
	AsicFilenameFactory ASiCWithCAdESFilenameFactory

	// ResourcesHandlerBuilder is used to write a created CMS into a defined implementation of
	// an OutputStream or a DSSDocument. Java declares the field protected.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

// NewAbstractASiCWithCAdESContainerMergerBase is the empty constructor. Port of the
// package-private empty constructor. The subclass constructor must follow it with
// InitDefaultContainerMerger and then InitFromDocuments/InitFromASiCContents.
func NewAbstractASiCWithCAdESContainerMergerBase() AbstractASiCWithCAdESContainerMerger {
	return AbstractASiCWithCAdESContainerMerger{
		DefaultContainerMerger:  asic.NewDefaultContainerMergerBase(),
		AsicFilenameFactory:     NewDefaultASiCWithCAdESFilenameFactory(),
		ResourcesHandlerBuilder: dsscades.CAdESUtilsDefaultResourcesHandlerBuilder,
	}
}

// SetAsicFilenameFactory sets the ASiCWithCAdESFilenameFactory defining a set of rules for
// naming of newly created ZIP entries, such as signature files.
//
// Panics with the Java message when asicFilenameFactory is nil (Objects.requireNonNull).
func (m *AbstractASiCWithCAdESContainerMerger) SetAsicFilenameFactory(asicFilenameFactory ASiCWithCAdESFilenameFactory) {
	if asicFilenameFactory == nil {
		panic("ASiCWithCAdESFilenameFactory cannot be null!")
	}
	m.AsicFilenameFactory = asicFilenameFactory
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal CMS objects during the signature creation procedure.
//
// NOTE: The DSSResourcesHandlerBuilder is supported only within the 'dss-cms-stream' module!
func (m *AbstractASiCWithCAdESContainerMerger) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	m.ResourcesHandlerBuilder = cms.UtilsResourcesHandlerBuilder(resourcesHandlerBuilder)
}

// IsSupportedDocument ports the @Override protected isSupported(DSSDocument).
func (m *AbstractASiCWithCAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	return NewASiCWithCAdESFormatDetector().IsSupportedZip(container)
}

// IsSupportedContent ports the @Override protected isSupported(ASiCContent).
func (m *AbstractASiCWithCAdESContainerMerger) IsSupportedContent(asicContent *asic.Content) bool {
	return NewASiCWithCAdESFormatDetector().IsSupportedZipContent(asicContent)
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(DSSDocument).
func (m *AbstractASiCWithCAdESContainerMerger) GetContainerExtractor(container model.DSSDocument) *asic.DefaultContainerExtractor {
	return &NewASiCWithCAdESContainerExtractor(container).DefaultContainerExtractor
}

// MergeCmsSignatures merges signature documents representing CMS signatures into a single CMS
// signature document. Ports the protected mergeCmsSignatures(List).
//
// Panics with Java's DSSException message on a merge failure.
func (m *AbstractASiCWithCAdESContainerMerger) MergeCmsSignatures(signatureDocuments []model.DSSDocument) model.DSSDocument {
	cmsList := m.getCMSList(signatureDocuments)

	originalCMS := cmsList[0] // getFirstCMS

	signerInformationStore := m.getSignerInformationStore(cmsList)
	mergedCMS, err := cms.UtilsReplaceSigners(originalCMS, signerInformationStore)
	if err != nil {
		panic(fmt.Sprintf("Unable to merge ASiC-S with CAdES container. Reason : %s", err.Error()))
	}

	certificatesStore := m.getCertificatesStore(cmsList)
	certAttributeStore := m.getCertAttributeStore(cmsList)
	crlStore := m.getCRLStore(cmsList)
	ocspResponsesStore := m.getOCSPResponsesStore(cmsList)
	ocspBasicStore := m.getOCSPBasicStore(cmsList)
	mergedCMS, err = cms.UtilsReplaceCertificatesAndCRLs(mergedCMS,
		certificatesStore, certAttributeStore, crlStore, ocspResponsesStore, ocspBasicStore)
	if err != nil {
		panic(fmt.Sprintf("Unable to merge ASiC-S with CAdES container. Reason : %s", err.Error()))
	}

	digestAlgorithms := m.getDigestAlgorithms(cmsList)
	mergedCMS, err = cms.UtilsPopulateDigestAlgorithmSet(mergedCMS, digestAlgorithms)
	if err != nil {
		panic(fmt.Sprintf("Unable to merge ASiC-S with CAdES container. Reason : %s", err.Error()))
	}

	cmsDocument, err := cms.UtilsWriteToDSSDocument(mergedCMS, m.ResourcesHandlerBuilder)
	if err != nil {
		panic(fmt.Sprintf("Unable to merge ASiC-S with CAdES container. Reason : %s", err.Error()))
	}
	cmsDocument.SetName(m.getSignatureDocumentName(signatureDocuments))
	return cmsDocument
}

// getCMSList ports the private getCMSList(List).
func (m *AbstractASiCWithCAdESContainerMerger) getCMSList(signatureDocuments []model.DSSDocument) []*cms.CMS {
	signedDataList := make([]*cms.CMS, 0, len(signatureDocuments))
	for _, signatureDocument := range signatureDocuments {
		documentAnalyzer, err := dsscades.NewCMSDocumentAnalyzerFromDocument(signatureDocument)
		if err != nil {
			panic(err)
		}
		signedDataList = append(signedDataList, documentAnalyzer.CMS())
	}
	return signedDataList
}

// getSignerInformationStore ports the private getSignerInformationStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getSignerInformationStore(cmsList []*cms.CMS) []*cmscore.SignerInfo {
	signerInformations := make([]*cmscore.SignerInfo, 0)
	for _, signedData := range cmsList {
		signerInformations = append(signerInformations, signedData.SignerInfos()...)
	}
	return signerInformations
}

// getCertificatesStore ports the private getCertificatesStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getCertificatesStore(cmsList []*cms.CMS) [][]byte {
	var result [][]byte
	for _, signedData := range cmsList {
		result = appendUniqueBytes(result, signedData.Certificates()...)
	}
	return result
}

// getCertAttributeStore ports the private getCertAttributeStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getCertAttributeStore(cmsList []*cms.CMS) [][]byte {
	var result [][]byte
	for _, signedData := range cmsList {
		result = appendUniqueBytes(result, signedData.AttributeCertificates()...)
	}
	return result
}

// getCRLStore ports the private getCRLStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getCRLStore(cmsList []*cms.CMS) [][]byte {
	var result [][]byte
	for _, c := range cmsList {
		result = appendUniqueBytes(result, c.CRLs()...)
	}
	return result
}

// getOCSPResponsesStore ports the private getOCSPResponsesStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getOCSPResponsesStore(cmsList []*cms.CMS) [][]byte {
	var result [][]byte
	for _, c := range cmsList {
		result = appendUniqueBytes(result, c.OcspResponseStore()...)
	}
	return result
}

// getOCSPBasicStore ports the private getOCSPBasicStore(List).
func (m *AbstractASiCWithCAdESContainerMerger) getOCSPBasicStore(cmsList []*cms.CMS) [][]byte {
	var result [][]byte
	for _, c := range cmsList {
		result = appendUniqueBytes(result, c.OcspBasicStore()...)
	}
	return result
}

// getDigestAlgorithms ports the private getDigestAlgorithms(List).
func (m *AbstractASiCWithCAdESContainerMerger) getDigestAlgorithms(cmsList []*cms.CMS) []*asn1ber.AlgorithmIdentifier {
	result := make([]*asn1ber.AlgorithmIdentifier, 0)
	for _, c := range cmsList {
		result = append(result, c.DigestAlgorithmIDs()...)
	}
	return result
}

// getSignatureDocumentName ports the private getSignatureDocumentName(List).
//
// Panics with an *exception.IllegalInputException when no signature document was provided.
func (m *AbstractASiCWithCAdESContainerMerger) getSignatureDocumentName(signatureDocuments []model.DSSDocument) string {
	if utils.IsCollectionNotEmpty(signatureDocuments) {
		return signatureDocuments[0].Name()
	}
	panic(exception.NewIllegalInputException("At least one signature file shall be provided for merging!"))
}

// GetAllSignatureDocuments returns all signature documents extracted from the given Content
// containers. Ports the protected getAllSignatureDocuments(ASiCContent...).
func (m *AbstractASiCWithCAdESContainerMerger) GetAllSignatureDocuments(asicContents ...*asic.Content) []model.DSSDocument {
	signatureDocuments := make([]model.DSSDocument, 0)
	for _, asicContent := range asicContents {
		signatureDocuments = append(signatureDocuments, asicContent.SignatureDocuments()...)
	}
	return signatureDocuments
}

// appendUniqueBytes appends each of items to result, skipping an item already present (matching
// java.util.LinkedHashSet's dedup-by-equality with insertion order preserved). Not a cross-file
// shared helper, per PORTING.md - kept local to this file (used by the five getXStore methods
// above).
func appendUniqueBytes(result [][]byte, items ...[]byte) [][]byte {
	for _, item := range items {
		found := false
		for _, existing := range result {
			if bytes.Equal(existing, item) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}
	return result
}
