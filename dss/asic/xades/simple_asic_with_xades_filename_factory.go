// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/SimpleASiCWithXAdESFilenameFactory.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package xades

import (
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SimpleASiCWithXAdESFilenameFactory provides a simple way to define custom names for file
// entries created within an ASiC with XAdES container, by using set and get methods.
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
type SimpleASiCWithXAdESFilenameFactory struct {
	DefaultASiCWithXAdESFilenameFactory

	// signatureFilename defines a name of a creating signature file (e.g. "signatures.xml").
	signatureFilename string

	// manifestFilename defines a name of a creating manifest file (e.g. "manifest.xml").
	manifestFilename string

	// dataPackageFilename defines a name of a creating ZIP archive, containing multiple
	// signer documents (in case of ASiC-S container).
	dataPackageFilename string

	// evidenceRecordManifestFilename defines a new name for the last evidence record manifest
	// file to be created (e.g. "META-INF/ASiCEvidenceRecordManifest001.xml").
	evidenceRecordManifestFilename string
}

var _ ASiCWithXAdESFilenameFactory = (*SimpleASiCWithXAdESFilenameFactory)(nil)

// NewSimpleASiCWithXAdESFilenameFactory instantiates a factory with null values. Port of the
// default constructor.
func NewSimpleASiCWithXAdESFilenameFactory() *SimpleASiCWithXAdESFilenameFactory {
	return &SimpleASiCWithXAdESFilenameFactory{}
}

// SignatureFilename ports the @Override getSignatureFilename(ASiCContent).
func (f *SimpleASiCWithXAdESFilenameFactory) SignatureFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.signatureFilename) {
		return f.getValidSignatureFilename(f.signatureFilename, asicContent)
	}
	return f.DefaultASiCWithXAdESFilenameFactory.SignatureFilename(asicContent)
}

// SetSignatureFilename sets a filename for a new signature document (when applicable).
//
// NOTE: The name of the signature file shall be:
//   - ASiC-S with XAdES : "META-INF/signatures.xml";
//   - ASiC-E with XAdES : "META-INF/signatures*.xml";
//   - OpenDocument : "META-INF/documentsignatures.xml".
//
// "META-INF/" is optional.
func (f *SimpleASiCWithXAdESFilenameFactory) SetSignatureFilename(signatureFilename string) {
	f.signatureFilename = signatureFilename
}

// ManifestFilename ports the @Override getManifestFilename(ASiCContent).
func (f *SimpleASiCWithXAdESFilenameFactory) ManifestFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.manifestFilename) {
		return f.getValidManifestFilename(f.manifestFilename, asicContent)
	}
	return f.DefaultASiCWithXAdESFilenameFactory.ManifestFilename(asicContent)
}

// SetManifestFilename sets a filename for a new manifest document (when applicable).
//
// NOTE: The name of the manifest file shall be:
//   - ASiC-E with XAdES : "META-INF/manifest.xml".
//
// "META-INF/" is optional.
func (f *SimpleASiCWithXAdESFilenameFactory) SetManifestFilename(manifestFilename string) {
	f.manifestFilename = manifestFilename
}

// DataPackageFilename ports the @Override getDataPackageFilename(ASiCContent).
func (f *SimpleASiCWithXAdESFilenameFactory) DataPackageFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.dataPackageFilename) {
		validFilename, err := f.ValidDataPackageFilename(f.dataPackageFilename, asicContent)
		if err != nil {
			panic(err)
		}
		return validFilename
	}
	return f.DefaultASiCWithXAdESFilenameFactory.DataPackageFilename(asicContent)
}

// SetDataPackageFilename sets a filename for a new ZIP data package (when applicable).
//
// NOTE: The name of the data package file shall be:
//   - ASiC-S with XAdES : "*.zip".
func (f *SimpleASiCWithXAdESFilenameFactory) SetDataPackageFilename(dataPackageFilename string) {
	f.dataPackageFilename = dataPackageFilename
}

// EvidenceRecordManifestFilename ports the @Override
// getEvidenceRecordManifestFilename(ASiCContent).
func (f *SimpleASiCWithXAdESFilenameFactory) EvidenceRecordManifestFilename(asicContent *asic.ASiCContent) string {
	if utils.IsStringNotEmpty(f.evidenceRecordManifestFilename) {
		validFilename, err := f.ValidEvidenceRecordManifestFilename(f.evidenceRecordManifestFilename, asicContent)
		if err != nil {
			panic(err)
		}
		return validFilename
	}
	return f.DefaultASiCWithXAdESFilenameFactory.EvidenceRecordManifestFilename(asicContent)
}

// SetEvidenceRecordManifestFilename sets a new filename for the ASiC evidence record manifest
// document (when applicable).
func (f *SimpleASiCWithXAdESFilenameFactory) SetEvidenceRecordManifestFilename(evidenceRecordManifestFilename string) {
	f.evidenceRecordManifestFilename = evidenceRecordManifestFilename
}

// getValidSignatureFilename returns a valid signature filename. Ports the protected
// getValidSignatureFilename(String, ASiCContent).
//
// Panics with Java's IllegalArgumentException messages on an invalid filename.
func (f *SimpleASiCWithXAdESFilenameFactory) getValidSignatureFilename(signatureFilename string, asicContent *asic.ASiCContent) string {
	signatureFilename = f.WithMetaInfFolder(signatureFilename)
	if err := f.AssertFilenameValid(signatureFilename, asicContent.SignatureDocuments()); err != nil {
		panic(err)
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	isOpenDocument, err := asic.ASiCUtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		panic(err)
	}
	if isASiCS && asic.ASiCUtilsSignaturesXML != signatureFilename {
		panic(exception.NewIllegalInputException(fmt.Sprintf("A signature file within ASiC-S with XAdES container "+
			"shall have name '%s'!", asic.ASiCUtilsSignaturesXML)))

	} else if isOpenDocument && asic.ASiCUtilsOpenDocumentSignatures != signatureFilename {
		panic(exception.NewIllegalInputException(fmt.Sprintf("A signature file within OpenDocument container "+
			"shall have name '%s'!", asic.ASiCUtilsOpenDocumentSignatures)))

	} else if !strings.HasPrefix(signatureFilename, asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsSignaturesFilename) ||
		!strings.HasSuffix(signatureFilename, asic.ASiCUtilsXMLExtension) { // ASiC-E
		panic(exception.NewIllegalInputException(fmt.Sprintf("A signature file within ASiC-E with XAdES container "+
			"shall match the template '%s'!", asic.ASiCUtilsMetaInfFolder+asic.ASiCUtilsSignaturesFilename+"*"+
			asic.ASiCUtilsXMLExtension)))
	}
	return signatureFilename
}

// getValidManifestFilename returns a valid manifest filename. Ports the protected
// getValidManifestFilename(String, ASiCContent).
func (f *SimpleASiCWithXAdESFilenameFactory) getValidManifestFilename(manifestFilename string, asicContent *asic.ASiCContent) string {
	manifestFilename = f.WithMetaInfFolder(manifestFilename)
	if err := f.AssertFilenameValid(manifestFilename, asicContent.ManifestDocuments()); err != nil {
		panic(err)
	}
	if asic.ASiCUtilsASiCEMetaInfManifest != manifestFilename {
		panic(exception.NewIllegalInputException(fmt.Sprintf("A manifest file within ASiC with XAdES container "+
			"shall have name '%s'!", asic.ASiCUtilsASiCEMetaInfManifest)))
	}
	return manifestFilename
}
