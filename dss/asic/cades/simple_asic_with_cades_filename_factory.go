// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/SimpleASiCWithCAdESFilenameFactory.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package cades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SimpleASiCWithCAdESFilenameFactory provides a simple way to define custom names for file
// entries created within an ASiC with CAdES container, by using setter/getter methods.
//
// This factory adds "META-INF/" prefix to the filename, when required.
//
// When a target filename for a particular document type is not specified, then the default
// processing will take precedence.
//
// NOTE: This factory shall be modified when consequently signing/extending a single container.
//
// WARN: The class does not verify the conformance of the defined filenames to the EN 319 162-1
// standard.
type SimpleASiCWithCAdESFilenameFactory struct {
	DefaultASiCWithCAdESFilenameFactory

	// signatureFilename defines a name of a creating signature file (e.g. "signature001.p7s").
	signatureFilename string

	// timestampFilename defines a name of a creating timestamp file (e.g. "timestamp001.tst").
	timestampFilename string

	// evidenceRecordFilename defines a name of a creating evidence record file (e.g.
	// "evidencerecord001.ers" or "evidencerecord001.xml").
	evidenceRecordFilename string

	// manifestFilename defines a name of a creating manifest file (e.g. "ASiCManifest001.xml").
	manifestFilename string

	// archiveManifestFilename defines a new name for the last archive manifest file to be moved
	// (e.g. "ASiCArchiveManifest001.xml").
	archiveManifestFilename string

	// evidenceRecordManifestFilename defines a new name for the last evidence record manifest
	// file to be created (e.g. "META-INF/ASiCEvidenceRecordManifest001.xml").
	evidenceRecordManifestFilename string

	// dataPackageFilename defines a name of a creating ZIP archive, containing multiple signer
	// documents (in case of ASiC-S container).
	dataPackageFilename string
}

var _ ASiCWithCAdESFilenameFactory = (*SimpleASiCWithCAdESFilenameFactory)(nil)

// NewSimpleASiCWithCAdESFilenameFactory instantiates a factory with null values. Port of the
// default constructor.
func NewSimpleASiCWithCAdESFilenameFactory() *SimpleASiCWithCAdESFilenameFactory {
	return &SimpleASiCWithCAdESFilenameFactory{}
}

// SignatureFilename ports the @Override getSignatureFilename(ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) SignatureFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.signatureFilename) {
		return f.getValidSignatureFilename(f.signatureFilename, asicContent)
	}
	return f.DefaultASiCWithCAdESFilenameFactory.SignatureFilename(asicContent)
}

// SetSignatureFilename sets a filename for a new signature document (when applicable).
//
// NOTE: The name of the signature file shall be:
//   - ASiC-S with CAdES : "META-INF/signature.p7s";
//   - ASiC-E with CAdES : "META-INF/signature*.p7s".
//
// "META-INF/" is optional.
func (f *SimpleASiCWithCAdESFilenameFactory) SetSignatureFilename(signatureFilename string) {
	f.signatureFilename = signatureFilename
}

// TimestampFilename ports the @Override getTimestampFilename(ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) TimestampFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.timestampFilename) {
		return f.getValidTimestampFilename(f.timestampFilename, asicContent)
	}
	return f.DefaultASiCWithCAdESFilenameFactory.TimestampFilename(asicContent)
}

// SetTimestampFilename sets a filename for a new timestamp document (when applicable).
//
// NOTE: The name of the timestamp file shall be:
//   - ASiC-S with CAdES : "META-INF/timestamp.tst";
//   - ASiC-E with CAdES : "META-INF/timestamp*.tst".
//
// "META-INF/" is optional.
func (f *SimpleASiCWithCAdESFilenameFactory) SetTimestampFilename(timestampFilename string) {
	f.timestampFilename = timestampFilename
}

// EvidenceRecordFilename ports the @Override getEvidenceRecordFilename(ASiCContent,
// EvidenceRecordTypeEnum).
//
// Panics with Java's NullPointerException message when evidenceRecordType is empty
// (Objects.requireNonNull).
func (f *SimpleASiCWithCAdESFilenameFactory) EvidenceRecordFilename(asicContent *asic.ASiCContent, evidenceRecordType enumerations.EvidenceRecordTypeEnum) string {
	if evidenceRecordType == "" {
		panic("EvidenceRecordType shall be defined!")
	}
	if utils.IsStringNotEmpty(f.evidenceRecordFilename) {
		return f.getValidEvidenceRecordFilename(f.evidenceRecordFilename, asicContent, evidenceRecordType)
	}
	return f.DefaultASiCWithCAdESFilenameFactory.EvidenceRecordFilename(asicContent, evidenceRecordType)
}

// SetEvidenceRecordFilename sets a filename for a new evidence record document.
//
// NOTE: The name of the evidence record file shall be:
//   - ASiC-S with CAdES : "META-INF/evidencerecord.xml" or "META-INF/evidencerecord.ers"
//   - ASiC-E with CAdES : "META-INF/*evidencerecord*.xml" or "META-INF/*evidencerecord*.ers"
//
// "META-INF/" is optional.
func (f *SimpleASiCWithCAdESFilenameFactory) SetEvidenceRecordFilename(evidenceRecordFilename string) {
	f.evidenceRecordFilename = evidenceRecordFilename
}

// ManifestFilename ports the @Override getManifestFilename(ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) ManifestFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.manifestFilename) {
		return f.getValidManifestFilename(f.manifestFilename, asicContent)
	}
	return f.DefaultASiCWithCAdESFilenameFactory.ManifestFilename(asicContent)
}

// SetManifestFilename sets a filename for a new manifest document (when applicable).
//
// NOTE: The name of the manifest file shall be:
//   - ASiC-E with CAdES : "META-INF/ASiCManifest*.xml".
//
// "META-INF/" is optional.
func (f *SimpleASiCWithCAdESFilenameFactory) SetManifestFilename(manifestFilename string) {
	f.manifestFilename = manifestFilename
}

// ArchiveManifestFilename ports the @Override getArchiveManifestFilename(ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) ArchiveManifestFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.archiveManifestFilename) {
		return f.getValidArchiveManifestFilename(f.archiveManifestFilename, asicContent)
	}
	return f.DefaultASiCWithCAdESFilenameFactory.ArchiveManifestFilename(asicContent)
}

// SetArchiveManifestFilename sets a new filename for the last archive manifest document (when
// applicable).
func (f *SimpleASiCWithCAdESFilenameFactory) SetArchiveManifestFilename(archiveManifestFilename string) {
	f.archiveManifestFilename = archiveManifestFilename
}

// DataPackageFilename ports the @Override getDataPackageFilename(ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) DataPackageFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.dataPackageFilename) {
		validFilename, err := f.ValidDataPackageFilename(f.dataPackageFilename, asicContent)
		if err != nil {
			panic(err)
		}
		return validFilename
	}
	return f.DefaultASiCWithCAdESFilenameFactory.DataPackageFilename(asicContent)
}

// SetDataPackageFilename sets a filename for a new ZIP data package (when applicable).
func (f *SimpleASiCWithCAdESFilenameFactory) SetDataPackageFilename(dataPackageFilename string) {
	f.dataPackageFilename = dataPackageFilename
}

// EvidenceRecordManifestFilename ports the @Override getEvidenceRecordManifestFilename(
// ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) EvidenceRecordManifestFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.evidenceRecordManifestFilename) {
		validFilename, err := f.ValidEvidenceRecordManifestFilename(f.evidenceRecordManifestFilename, asicContent)
		if err != nil {
			panic(err)
		}
		return validFilename
	}
	return f.DefaultASiCWithCAdESFilenameFactory.EvidenceRecordManifestFilename(asicContent)
}

// SetEvidenceRecordManifestFilename sets a new filename for the ASiC evidence record manifest
// document (when applicable).
func (f *SimpleASiCWithCAdESFilenameFactory) SetEvidenceRecordManifestFilename(evidenceRecordManifestFilename string) {
	f.evidenceRecordManifestFilename = evidenceRecordManifestFilename
}

// getValidSignatureFilename returns a valid signature filename. Ports the protected
// getValidSignatureFilename(String, ASiCContent).
//
// Panics with Java's IllegalArgumentException messages on an invalid filename.
func (f *SimpleASiCWithCAdESFilenameFactory) getValidSignatureFilename(signatureFilename string, asicContent *asic.ASiCContent) string {
	signatureFilename = f.WithMetaInfFolder(signatureFilename)
	if err := f.AssertFilenameValid(signatureFilename, asicContent.SignatureDocuments()); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS && asic.ASiCUtilsSignatureP7S != signatureFilename {
		panic(exception.NewIllegalInputException("A signature file within ASiC-S with CAdES container " +
			"shall have name '" + asic.ASiCUtilsSignatureP7S + "'!"))

	} else if !strings.HasPrefix(signatureFilename, asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsSignatureFilename) ||
		!strings.HasSuffix(signatureFilename, asic.ASiCUtilsCAdESSignatureExtension) { // ASiC-E
		panic(exception.NewIllegalInputException("A signature file within ASiC-E with CAdES container " +
			"shall match the template '" + asic.ASiCUtilsMetaInfFolder + asic.ASiCUtilsSignatureFilename + "*" +
			asic.ASiCUtilsCAdESSignatureExtension + "'!"))
	}
	return signatureFilename
}

// getValidTimestampFilename returns a valid timestamp filename. Ports the protected
// getValidTimestampFilename(String, ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) getValidTimestampFilename(timestampFilename string, asicContent *asic.ASiCContent) string {
	timestampFilename = f.WithMetaInfFolder(timestampFilename)
	if err := f.AssertFilenameValid(timestampFilename, asicContent.TimestampDocuments()); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS && utils.IsCollectionEmpty(asicContent.TimestampDocuments()) &&
		asic.ASiCUtilsTimestampTST != timestampFilename {
		panic(exception.NewIllegalInputException("A timestamp file within ASiC-S with CAdES container " +
			"shall have name '" + asic.ASiCUtilsTimestampTST + "'!"))

	} else if !strings.HasPrefix(timestampFilename, asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsTimestampFilename) ||
		!strings.HasSuffix(timestampFilename, asic.ASiCUtilsTSTExtension) { // ASiC-E
		panic(exception.NewIllegalInputException("A timestamp file within ASiC-E with CAdES container " +
			"shall match the template '" + asic.ASiCUtilsMetaInfFolder + asic.ASiCUtilsTimestampFilename + "*" +
			asic.ASiCUtilsTSTExtension + "'!"))
	}
	return timestampFilename
}

// getValidEvidenceRecordFilename returns a valid evidence record filename. Ports the protected
// getValidEvidenceRecordFilename(String, ASiCContent, EvidenceRecordTypeEnum).
func (f *SimpleASiCWithCAdESFilenameFactory) getValidEvidenceRecordFilename(evidenceRecordFilename string, asicContent *asic.ASiCContent, evidenceRecordType enumerations.EvidenceRecordTypeEnum) string {
	if err := f.AssertASiCContentIsValid(asicContent); err != nil {
		panic(err)
	}
	evidenceRecordFilename = f.WithMetaInfFolder(evidenceRecordFilename)
	if err := f.AssertFilenameValid(evidenceRecordFilename, asicContent.EvidenceRecordDocuments()); err != nil {
		panic(err)
	}
	switch evidenceRecordType {
	case enumerations.EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD:
		if !strings.HasSuffix(evidenceRecordFilename, asic.ASiCUtilsXMLExtension) {
			panic(exception.NewIllegalInputException("An XMLERS evidence record file within " +
				"ASiC container shall end with '" + asic.ASiCUtilsXMLExtension + "' extension!"))
		}
	case enumerations.EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD:
		if !strings.HasSuffix(evidenceRecordFilename, asic.ASiCUtilsERASN1Extension) {
			panic(exception.NewIllegalInputException("An ERS evidence record file within " +
				"ASiC container shall end with '" + asic.ASiCUtilsERASN1Extension + "' extension!"))
		}
	default:
		panic(exception.NewIllegalInputException(
			"The Evidence Record Type '" + string(evidenceRecordType) + "' is not supported!"))
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	if isASiCS && utils.IsCollectionEmpty(asicContent.EvidenceRecordDocuments()) &&
		!(asic.ASiCUtilsEvidenceRecordERS == evidenceRecordFilename || asic.ASiCUtilsEvidenceRecordXML == evidenceRecordFilename) {
		panic(exception.NewIllegalInputException("An evidence record file within ASiC-S with CAdES container " +
			"shall have name '" + asic.ASiCUtilsEvidenceRecordERS + "' or '" + asic.ASiCUtilsEvidenceRecordXML + "'!"))

	} else if !strings.HasPrefix(evidenceRecordFilename, asic.ASiCUtilsMetaInfFolder) ||
		!strings.Contains(evidenceRecordFilename, asic.ASiCUtilsEvidenceRecordFilename) {
		// ASiC-E
		panic(exception.NewIllegalInputException("An evidence record file within ASiC-E with CAdES container " +
			"shall match the template '" + asic.ASiCUtilsMetaInfFolder + "*" + asic.ASiCUtilsEvidenceRecordFilename + "*" +
			"(" + asic.ASiCUtilsERASN1Extension + "||" + asic.ASiCUtilsXMLExtension + ")'!"))
	}
	return evidenceRecordFilename
}

// getValidManifestFilename returns a valid manifest filename. Ports the protected
// getValidManifestFilename(String, ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) getValidManifestFilename(manifestFilename string, asicContent *asic.ASiCContent) string {
	manifestFilename = f.WithMetaInfFolder(manifestFilename)
	if err := f.AssertFilenameValid(manifestFilename, asicContent.ManifestDocuments()); err != nil {
		panic(err)
	}
	if !strings.HasPrefix(manifestFilename, asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsASiCManifestFilename) ||
		!strings.HasSuffix(manifestFilename, asic.ASiCUtilsXMLExtension) {
		panic(exception.NewIllegalInputException("A manifest file within ASiC with CAdES container " +
			"shall match the template '" + asic.ASiCUtilsMetaInfFolder + asic.ASiCUtilsASiCManifestFilename + "*" +
			asic.ASiCUtilsXMLExtension + "'!"))
	}
	return manifestFilename
}

// getValidArchiveManifestFilename returns a valid archive manifest filename.
//
// NOTE: The name of the archive manifest file shall be:
//   - ASiC-E with CAdES : "META-INF/ASiCArchiveManifest*.xml".
//
// "META-INF/" is optional.
//
// Ports the protected getValidArchiveManifestFilename(String, ASiCContent).
func (f *SimpleASiCWithCAdESFilenameFactory) getValidArchiveManifestFilename(archiveManifestFilename string, asicContent *asic.ASiCContent) string {
	archiveManifestFilename = f.WithMetaInfFolder(archiveManifestFilename)
	if err := f.AssertFilenameValid(archiveManifestFilename, asicContent.ArchiveManifestDocuments()); err != nil {
		panic(err)
	}
	if !strings.HasPrefix(archiveManifestFilename, asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsASiCArchiveManifestFilename) ||
		!strings.HasSuffix(archiveManifestFilename, asic.ASiCUtilsXMLExtension) {
		panic(exception.NewIllegalInputException("An archive manifest file within ASiC with CAdES container " +
			"shall match the template '" + asic.ASiCUtilsMetaInfFolder + asic.ASiCUtilsASiCArchiveManifestFilename + "*" +
			asic.ASiCUtilsXMLExtension + "'!"))

	} else if ASiCWithCAdESUtilsDefaultArchiveManifestFilename == archiveManifestFilename {
		panic(exception.NewIllegalInputException("An archive manifest file within ASiC with CAdES container " +
			"cannot be moved to a file with name '" + ASiCWithCAdESUtilsDefaultArchiveManifestFilename + "'!"))
	}
	return archiveManifestFilename
}
