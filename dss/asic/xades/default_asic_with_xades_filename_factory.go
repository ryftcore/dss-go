// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/DefaultASiCWithXAdESFilenameFactory.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// DefaultASiCWithXAdESFilenameFactory provides a default implementation of
// ASiCWithXAdESFilenameFactory used within basic configuration of DSS for creation of
// filenames for new container entries.
type DefaultASiCWithXAdESFilenameFactory struct {
	asic.AbstractFilenameFactory
}

var _ ASiCWithXAdESFilenameFactory = (*DefaultASiCWithXAdESFilenameFactory)(nil)

// NewDefaultASiCWithXAdESFilenameFactory is the default constructor.
func NewDefaultASiCWithXAdESFilenameFactory() *DefaultASiCWithXAdESFilenameFactory {
	return &DefaultASiCWithXAdESFilenameFactory{}
}

// SignatureFilename ports the @Override getSignatureFilename(ASiCContent).
func (f *DefaultASiCWithXAdESFilenameFactory) SignatureFilename(asicContent *asic.Content) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCS, err := asic.UtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS {
		return asic.ASiCUtilsSignaturesXML // "META-INF/signatures.xml"
	}

	isOpenDocument, err := asic.UtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		panic(err)
	}
	if isOpenDocument {
		return asic.ASiCUtilsOpenDocumentSignatures // "META-INF/documentsignatures.xml"
	}

	// ASiC-E
	existingSignatureNames := spi.DSSUtilsDocumentNames(asicContent.SignatureDocuments())
	// "META-INF/signatures*.xml"
	return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfXAdESSignature, existingSignatureNames)
}

// ManifestFilename ports the @Override getManifestFilename(ASiCContent).
func (f *DefaultASiCWithXAdESFilenameFactory) ManifestFilename(asicContent *asic.Content) string {
	return asic.ASiCUtilsASiCEMetaInfManifest // "META-INF/manifest.xml"
}

// DataPackageFilename ports the @Override getDataPackageFilename(ASiCContent).
func (f *DefaultASiCWithXAdESFilenameFactory) DataPackageFilename(asicContent *asic.Content) string {
	return asic.ASiCUtilsPackageZip // "package.zip"
}

// EvidenceRecordFilename ports the @Override getEvidenceRecordFilename(ASiCContent,
// EvidenceRecordTypeEnum).
//
// Panics with Java's NullPointerException message when evidenceRecordType is empty
// (Objects.requireNonNull), or its UnsupportedOperationException message for an unsupported
// evidenceRecordType.
func (f *DefaultASiCWithXAdESFilenameFactory) EvidenceRecordFilename(asicContent *asic.Content,
	evidenceRecordType enumerations.EvidenceRecordTypeEnum) string {
	if evidenceRecordType == "" {
		panic("EvidenceRecordType shall be defined!")
	}
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	// Same name for both ASiC-S and ASiC-E
	switch evidenceRecordType {
	case enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord:
		return asic.ASiCUtilsEvidenceRecordXML // "META-INF/evidencerecord.xml"
	case enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord:
		return asic.ASiCUtilsEvidenceRecordERS
	default:
		panic(exception.NewIllegalInputException(
			fmt.Sprintf("The Evidence Record Type '%s' is not supported!", evidenceRecordType)))
	}
}

// EvidenceRecordManifestFilename ports the @Override
// getEvidenceRecordManifestFilename(Content).
func (f *DefaultASiCWithXAdESFilenameFactory) EvidenceRecordManifestFilename(asicContent *asic.Content) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	existingManifestNames := spi.DSSUtilsDocumentNames(asicContent.EvidenceRecordManifestDocuments())
	// "META-INF/ASiCEvidenceRecordManifest*.xml"
	return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfEvidenceRecordManifest, existingManifestNames)
}
