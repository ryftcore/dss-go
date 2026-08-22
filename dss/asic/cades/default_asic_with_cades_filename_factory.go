// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/DefaultASiCWithCAdESFilenameFactory.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package cades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DefaultASiCWithCAdESFilenameFactory provides a default implementation of
// ASiCWithCAdESFilenameFactory used within basic configuration of DSS for creation of filenames
// for new container entries.
type DefaultASiCWithCAdESFilenameFactory struct {
	asic.AbstractASiCFilenameFactory
}

var _ ASiCWithCAdESFilenameFactory = (*DefaultASiCWithCAdESFilenameFactory)(nil)

// NewDefaultASiCWithCAdESFilenameFactory is the default constructor.
func NewDefaultASiCWithCAdESFilenameFactory() *DefaultASiCWithCAdESFilenameFactory {
	return &DefaultASiCWithCAdESFilenameFactory{}
}

// SignatureFilename ports the @Override getSignatureFilename(ASiCContent).
func (f *DefaultASiCWithCAdESFilenameFactory) SignatureFilename(asicContent *asic.ASiCContent) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS {
		return asic.ASiCUtilsSignatureP7S // "META-INF/signature.p7s"
	}
	existingSignatureNames := spi.DSSUtilsDocumentNames(asicContent.SignatureDocuments())
	// "META-INF/signature*.p7s"
	return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfCAdESSignature, existingSignatureNames)
}

// TimestampFilename ports the @Override getTimestampFilename(ASiCContent).
func (f *DefaultASiCWithCAdESFilenameFactory) TimestampFilename(asicContent *asic.ASiCContent) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS && utils.IsCollectionEmpty(asicContent.TimestampDocuments()) {
		return asic.ASiCUtilsTimestampTST // "META-INF/timestamp.tst"
	}
	existingTimestampNames := spi.DSSUtilsDocumentNames(asicContent.TimestampDocuments())
	// "META-INF/timestamp*.tst"
	return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfCAdESTimestamp, existingTimestampNames)
}

// ManifestFilename ports the @Override getManifestFilename(ASiCContent).
//
// Panics with Java's UnsupportedOperationException message for an ASiC-S container.
func (f *DefaultASiCWithCAdESFilenameFactory) ManifestFilename(asicContent *asic.ASiCContent) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCE, err := asic.ASiCUtilsIsASiCEContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCE {
		existingManifestNames := spi.DSSUtilsDocumentNames(asicContent.ManifestDocuments())
		// "META-INF/ASiCManifest*.xml"
		return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfCAdESManifest, existingManifestNames)
	}
	panic("Manifest is not applicable for ASiC-S with CAdES container!")
}

// ArchiveManifestFilename ports the @Override getArchiveManifestFilename(ASiCContent).
//
// Panics with Java's UnsupportedOperationException message for an ASiC-S container.
func (f *DefaultASiCWithCAdESFilenameFactory) ArchiveManifestFilename(asicContent *asic.ASiCContent) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCE, err := asic.ASiCUtilsIsASiCEContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCE || utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) {
		existingArchiveManifestNames := spi.DSSUtilsDocumentNames(asicContent.ArchiveManifestDocuments())
		existingArchiveManifestNames = removeString(existingArchiveManifestNames, ASiCWithCAdESUtilsDefaultArchiveManifestFilename)
		// "META-INF/ASiCArchiveManifest*.xml"
		return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfCAdESArchiveManifest, existingArchiveManifestNames)
	}
	panic("Manifest is not applicable for ASiC-S with CAdES container!")
}

// DataPackageFilename ports the @Override getDataPackageFilename(ASiCContent).
func (f *DefaultASiCWithCAdESFilenameFactory) DataPackageFilename(asicContent *asic.ASiCContent) string {
	return asic.ASiCUtilsPackageZip // "package.zip"
}

// EvidenceRecordFilename ports the @Override getEvidenceRecordFilename(ASiCContent,
// EvidenceRecordTypeEnum).
//
// Panics with Java's NullPointerException message when evidenceRecordType is empty
// (Objects.requireNonNull), or its UnsupportedOperationException message for an unsupported
// evidenceRecordType.
func (f *DefaultASiCWithCAdESFilenameFactory) EvidenceRecordFilename(asicContent *asic.ASiCContent, evidenceRecordType enumerations.EvidenceRecordTypeEnum) string {
	if evidenceRecordType == "" {
		panic("EvidenceRecordType shall be defined!")
	}
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS {
		switch evidenceRecordType {
		case enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord:
			return asic.ASiCUtilsEvidenceRecordXML // "META-INF/evidencerecord.xml"
		case enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord:
			return asic.ASiCUtilsEvidenceRecordERS
		default:
			panic(exception.NewIllegalInputException(
				"The Evidence Record Type '" + string(evidenceRecordType) + "' is not supported!"))
		}
	}

	// ASiC-E
	existingEvidenceRecordNames := spi.DSSUtilsDocumentNames(asicContent.EvidenceRecordDocuments())
	var targetEvidenceRecordName string
	switch evidenceRecordType {
	case enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord:
		targetEvidenceRecordName = asic.ASiCUtilsASiCEMetaInfCAdESEvidenceRecordXML // "META-INF/evidencerecord*.xml"
		existingEvidenceRecordNames = filterStringsHasSuffix(existingEvidenceRecordNames, asic.ASiCUtilsXMLExtension)
	case enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord:
		targetEvidenceRecordName = asic.ASiCUtilsASiCEMetaInfCAdESEvidenceRecordASN1 // "META-INF/evidencerecord*.ers"
		existingEvidenceRecordNames = filterStringsHasSuffix(existingEvidenceRecordNames, asic.ASiCUtilsERASN1Extension)
	default:
		panic(exception.NewIllegalInputException(
			"The Evidence Record Type '" + string(evidenceRecordType) + "' is not supported!"))
	}
	return f.NextAvailableDocumentName(targetEvidenceRecordName, existingEvidenceRecordNames)
}

// EvidenceRecordManifestFilename ports the @Override getEvidenceRecordManifestFilename(
// ASiCContent).
func (f *DefaultASiCWithCAdESFilenameFactory) EvidenceRecordManifestFilename(asicContent *asic.ASiCContent) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	existingManifestNames := spi.DSSUtilsDocumentNames(asicContent.EvidenceRecordManifestDocuments())
	// "META-INF/ASiCEvidenceRecordManifest*.xml"
	return f.NextAvailableDocumentName(asic.ASiCUtilsASiCEMetaInfEvidenceRecordManifest, existingManifestNames)
}

// removeString ports the single List.remove(Object) call site above; not a cross-file shared
// helper, per PORTING.md - kept local to this file (used once, alongside filterStringsHasSuffix
// below).
func removeString(items []string, target string) []string {
	result := make([]string, 0, len(items))
	removed := false
	for _, item := range items {
		if !removed && item == target {
			removed = true
			continue
		}
		result = append(result, item)
	}
	return result
}

// filterStringsHasSuffix ports the repeated
// existingEvidenceRecordNames.stream().filter(n -> n.endsWith(...)).collect(...) idiom, local to
// this file per PORTING.md.
func filterStringsHasSuffix(items []string, suffix string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if strings.HasSuffix(item, suffix) {
			result = append(result, item)
		}
	}
	return result
}
