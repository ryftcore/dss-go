// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCUtils.java
// (DSS 6.5.RC1).
//
// NAMING: the class is a static-only utility holder, so its members become package-level
// identifiers prefixed with the class name, following the DSSUtils/CAdESUtils precedent. Java's
// overloads (getZipComment x4, getMimeType x3, getContainerType x3, isASiC{S,E}Container x2) cannot
// share a name in Go and carry a suffix naming the argument they take.
//
// slf4j LOG calls are dropped per PORTING.md; the branches they sit in are preserved.
package asic

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

const (
	// ASiCUtilsManifestFilename is the manifest name. Port of MANIFEST_FILENAME.
	ASiCUtilsManifestFilename = "Manifest"

	// ASiCUtilsASiCManifestFilename is the ASiC Manifest name. Port of ASIC_MANIFEST_FILENAME.
	ASiCUtilsASiCManifestFilename = "ASiCManifest"

	// ASiCUtilsASiCArchiveManifestFilename is the ASiC Archive Manifest name. Port of
	// ASIC_ARCHIVE_MANIFEST_FILENAME.
	ASiCUtilsASiCArchiveManifestFilename = "ASiCArchiveManifest"

	// ASiCUtilsASiCEvidenceRecordManifestFilename is the ASiC Evidence Record Manifest name.
	// Port of ASIC_EVIDENCE_RECORD_MANIFEST_FILENAME.
	ASiCUtilsASiCEvidenceRecordManifestFilename = "ASiCEvidenceRecordManifest"

	// ASiCUtilsASiCXAdESManifestFilename is the ASiC-E with XAdES Manifest name. Port of
	// ASIC_XAdES_MANIFEST_FILENAME.
	ASiCUtilsASiCXAdESManifestFilename = "manifest"

	// ASiCUtilsMimeType is the mimetype filename. Port of MIME_TYPE.
	ASiCUtilsMimeType = "mimetype"

	// ASiCUtilsMimeTypeComment is the mimetype comment. Port of MIME_TYPE_COMMENT.
	ASiCUtilsMimeTypeComment = ASiCUtilsMimeType + "="

	// ASiCUtilsMetaInfFolder is the META-INF folder. Port of META_INF_FOLDER.
	ASiCUtilsMetaInfFolder = "META-INF/"

	// ASiCUtilsPackageZip is the "package.zip" filename. Port of PACKAGE_ZIP.
	ASiCUtilsPackageZip = "package.zip"

	// ASiCUtilsSignatureFilename is the signature filename. Port of SIGNATURE_FILENAME.
	ASiCUtilsSignatureFilename = "signature"

	// ASiCUtilsSignaturesFilename is the signature filename. Port of SIGNATURES_FILENAME.
	ASiCUtilsSignaturesFilename = "signatures"

	// ASiCUtilsTimestampFilename is the timestamp filename. Port of TIMESTAMP_FILENAME.
	ASiCUtilsTimestampFilename = "timestamp"

	// ASiCUtilsEvidenceRecordFilename is the evidence record filename. Port of
	// EVIDENCE_RECORD_FILENAME.
	ASiCUtilsEvidenceRecordFilename = "evidencerecord"

	// ASiCUtilsCAdESSignatureExtension is the signature file extension. Port of
	// CADES_SIGNATURE_EXTENSION.
	ASiCUtilsCAdESSignatureExtension = ".p7s"

	// ASiCUtilsTSTExtension is the timestamp file extension. Port of TST_EXTENSION.
	ASiCUtilsTSTExtension = ".tst"

	// ASiCUtilsERASN1Extension is the evidence record IETF RFC 4998 file extension. Port of
	// ER_ASN1_EXTENSION.
	ASiCUtilsERASN1Extension = ".ers"

	// ASiCUtilsXMLExtension is the XML file extension. Port of XML_EXTENSION.
	ASiCUtilsXMLExtension = ".xml"

	// ASiCUtilsSignaturesXML is the ASiC-S with XAdES signature document name
	// (META-INF/signatures.xml). Port of SIGNATURES_XML.
	ASiCUtilsSignaturesXML = ASiCUtilsMetaInfFolder + ASiCUtilsSignaturesFilename + ASiCUtilsXMLExtension

	// ASiCUtilsOpenDocumentSignatures is the ASiC-E with XAdES for OpenDocument signature file
	// name (META-INF/documentsignatures.xml). Port of OPEN_DOCUMENT_SIGNATURES.
	ASiCUtilsOpenDocumentSignatures = ASiCUtilsMetaInfFolder + "documentsignatures.xml"

	// ASiCUtilsASiCEMetaInfManifest is the default XML manifest filename
	// (META-INF/manifest.xml). Port of ASICE_METAINF_MANIFEST.
	ASiCUtilsASiCEMetaInfManifest = ASiCUtilsMetaInfFolder + ASiCUtilsASiCXAdESManifestFilename + ASiCUtilsXMLExtension

	// ASiCUtilsASiCEMetaInfXAdESSignature is the default signature filename for ASiC-E with
	// XAdES container. Port of ASICE_METAINF_XADES_SIGNATURE.
	ASiCUtilsASiCEMetaInfXAdESSignature = ASiCUtilsMetaInfFolder + "signatures001.xml"

	// ASiCUtilsASiCEMetaInfCAdESSignature is the default signature filename for ASiC-E with
	// CAdES container. Port of ASICE_METAINF_CADES_SIGNATURE.
	ASiCUtilsASiCEMetaInfCAdESSignature = ASiCUtilsMetaInfFolder + "signature001.p7s"

	// ASiCUtilsASiCEMetaInfCAdESTimestamp is the default timestamp filename for ASiC-E with
	// CAdES container. Port of ASICE_METAINF_CADES_TIMESTAMP.
	ASiCUtilsASiCEMetaInfCAdESTimestamp = ASiCUtilsMetaInfFolder + "timestamp001.tst"

	// ASiCUtilsASiCEMetaInfCAdESEvidenceRecordASN1 is the default ERS evidence record filename
	// for ASiC-E with CAdES container. Port of ASICE_METAINF_CADES_EVIDENCE_RECORD_ASN1.
	ASiCUtilsASiCEMetaInfCAdESEvidenceRecordASN1 = ASiCUtilsMetaInfFolder + "evidencerecord001.ers"

	// ASiCUtilsASiCEMetaInfCAdESEvidenceRecordXML is the default XMLERS evidence record
	// filename for ASiC-E with CAdES container. Port of ASICE_METAINF_CADES_EVIDENCE_RECORD_XML.
	ASiCUtilsASiCEMetaInfCAdESEvidenceRecordXML = ASiCUtilsMetaInfFolder + "evidencerecord001.xml"

	// ASiCUtilsASiCEMetaInfCAdESManifest is the default ASIC manifest filename for ASiC-E with
	// CAdES container. Port of ASICE_METAINF_CADES_MANIFEST.
	ASiCUtilsASiCEMetaInfCAdESManifest = ASiCUtilsMetaInfFolder + "ASiCManifest001.xml"

	// ASiCUtilsASiCEMetaInfCAdESArchiveManifest is the default ASIC archive manifest filename
	// for ASiC-E with CAdES container. Port of ASICE_METAINF_CADES_ARCHIVE_MANIFEST.
	ASiCUtilsASiCEMetaInfCAdESArchiveManifest = ASiCUtilsMetaInfFolder + "ASiCArchiveManifest001.xml"

	// ASiCUtilsASiCEMetaInfEvidenceRecordManifest is the default ASIC evidence record manifest
	// filename for ASiC-E container. Port of ASICE_METAINF_EVIDENCE_RECORD_MANIFEST.
	ASiCUtilsASiCEMetaInfEvidenceRecordManifest = ASiCUtilsMetaInfFolder + "ASiCEvidenceRecordManifest001.xml"

	// ASiCUtilsSignatureP7S is the ASiC-S with CAdES signature document name
	// (META-INF/signature.p7s). Port of SIGNATURE_P7S.
	ASiCUtilsSignatureP7S = ASiCUtilsMetaInfFolder + ASiCUtilsSignatureFilename + ASiCUtilsCAdESSignatureExtension

	// ASiCUtilsTimestampTST is the ASiC-S with CAdES timestamp document name
	// (META-INF/timestamp.tst). Port of TIMESTAMP_TST.
	ASiCUtilsTimestampTST = ASiCUtilsMetaInfFolder + ASiCUtilsTimestampFilename + ASiCUtilsTSTExtension

	// ASiCUtilsEvidenceRecordERS is the ASiC with XAdES evidence record ASN.1 document name
	// (META-INF/evidencerecord.ers). Port of EVIDENCE_RECORD_ERS.
	ASiCUtilsEvidenceRecordERS = ASiCUtilsMetaInfFolder + ASiCUtilsEvidenceRecordFilename + ASiCUtilsERASN1Extension

	// ASiCUtilsEvidenceRecordXML is the ASiC with XAdES evidence record ASN.1 document name
	// (META-INF/evidencerecord.xml). Port of EVIDENCE_RECORD_XML.
	ASiCUtilsEvidenceRecordXML = ASiCUtilsMetaInfFolder + ASiCUtilsEvidenceRecordFilename + ASiCUtilsXMLExtension
)

// asicUtilsZipPrefix identifies the first bytes of a zip archive document. Port of the private
// ZIP_PREFIX.
var asicUtilsZipPrefix = []byte{'P', 'K'}

// asicUtilsMagicDir is the zip comment identifier in the end of ZIP archive. Port of the private
// MAGIC_DIR.
var asicUtilsMagicDir = []byte{0x50, 0x4b, 0x05, 0x06}

// asicUtilsMaxToRead is the maximum number of bytes to be read in a file to extract a zip comment.
// Port of the private MAX_TO_READ.
const asicUtilsMaxToRead = 0xFFFF + 2 + 4

// ASiCUtilsIsSignature verifies if the entryName represents a signature file name. Port of
// isSignature(String).
func UtilsIsSignature(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(entryName, ASiCUtilsSignatureFilename) &&
		!strings.Contains(entryName, ASiCUtilsManifestFilename)
}

// ASiCUtilsIsTimestamp verifies if the entryName represents a timestamp file name. Port of
// isTimestamp(String).
func UtilsIsTimestamp(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(entryName, ASiCUtilsTimestampFilename) &&
		strings.HasSuffix(entryName, ASiCUtilsTSTExtension)
}

// UtilsIsEvidenceRecord verifies if the entryName represents an evidence record filename. Port
// of isEvidenceRecord(String).
func UtilsIsEvidenceRecord(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(entryName, ASiCUtilsEvidenceRecordFilename) &&
		(strings.HasSuffix(entryName, ASiCUtilsXMLExtension) || strings.HasSuffix(entryName, ASiCUtilsERASN1Extension))
}

// UtilsIsXmlEvidenceRecord verifies if the entryName represents an XMLERS evidence record
// filename. Port of isXmlEvidenceRecord(String).
func UtilsIsXmlEvidenceRecord(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(entryName, ASiCUtilsEvidenceRecordFilename) &&
		strings.HasSuffix(entryName, ASiCUtilsXMLExtension)
}

// UtilsIsAsn1EvidenceRecord verifies if the entryName represents an ERS ASN.1 evidence record
// filename. Port of isAsn1EvidenceRecord(String).
func UtilsIsAsn1EvidenceRecord(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(entryName, ASiCUtilsEvidenceRecordFilename) &&
		strings.HasSuffix(entryName, ASiCUtilsERASN1Extension)
}

// ASiCUtilsMimeTypeString returns the target MimeType string. Port of
// getMimeTypeString(Parameters).
func UtilsMimeTypeString(asicParameters *Parameters) string {
	mimeType := UtilsMimeTypeFromParameters(asicParameters)
	return mimeType.MimeTypeString()
}

// UtilsZipCommentFromParameters returns a ZIP Comment String according to the given
// parameters. Port of getZipComment(ASiCParameters).
func UtilsZipCommentFromParameters(asicParameters *Parameters) string {
	if asicParameters.IsZipComment() {
		return UtilsZipCommentFromMimeTypeString(UtilsMimeTypeString(asicParameters))
	}
	return utils.EmptyString
}

// ASiCUtilsZipCommentFromMimeType returns a ZIP Comment String from the provided MimeType. Port of
// getZipComment(MimeType).
func UtilsZipCommentFromMimeType(mimeType enumerations.MimeType) string {
	return UtilsZipCommentFromMimeTypeString(mimeType.MimeTypeString())
}

// UtilsZipCommentFromMimeTypeString returns a ZIP Comment String from the provided
// mimeTypeString. Port of getZipComment(String).
func UtilsZipCommentFromMimeTypeString(mimeTypeString string) string {
	return ASiCUtilsMimeTypeComment + mimeTypeString
}

// ASiCUtilsIsASiCMimeType checks if the given MimeType is ASiC MimeType. Port of
// isASiCMimeType(MimeType).
func UtilsIsASiCMimeType(mimeType enumerations.MimeType) bool {
	return mimeType == enumerations.MimeTypeEnumASiCS ||
		mimeType == enumerations.MimeTypeEnumASiCE
}

// ASiCUtilsIsOpenDocumentMimeType checks if the given MimeType is OpenDocument MimeType. Port of
// isOpenDocumentMimeType(MimeType).
func UtilsIsOpenDocumentMimeType(mimeType enumerations.MimeType) bool {
	return mimeType == enumerations.MimeTypeEnumODT ||
		mimeType == enumerations.MimeTypeEnumODS ||
		mimeType == enumerations.MimeTypeEnumODG ||
		mimeType == enumerations.MimeTypeEnumODP
}

// UtilsASiCContainerType returns the related ASiCContainerType for the given asicMimeType.
//
// Panics with the Java message when asicMimeType is nil (Objects.requireNonNull); returns an error
// for a mimetype that is neither ASiC nor OpenDocument (IllegalArgumentException).
//
// Port of getASiCContainerType(MimeType).
func UtilsASiCContainerType(asicMimeType enumerations.MimeType) (enumerations.ASiCContainerType, error) {
	if asicMimeType == nil {
		panic("MimeType cannot be null!")
	}
	if asicMimeType == enumerations.MimeTypeEnumASiCS {
		return enumerations.ASiCContainerTypeASiCS, nil
	} else if asicMimeType == enumerations.MimeTypeEnumASiCE || UtilsIsOpenDocumentMimeType(asicMimeType) {
		return enumerations.ASiCContainerTypeASiCE, nil
	}
	return "", fmt.Errorf("Not allowed mimetype '%s'", asicMimeType.MimeTypeString())
}

// UtilsIsASiCE checks if the parameters are configured for ASiCE creation.
//
// Panics with the Java message when the container type is not set (Objects.requireNonNull).
//
// Port of isASiCE(ASiCParameters).
func UtilsIsASiCE(asicParameters *Parameters) bool {
	if asicParameters.ContainerType() == "" {
		panic("ASiCContainerType must be defined!")
	}
	return enumerations.ASiCContainerTypeASiCE == asicParameters.ContainerType()
}

// UtilsIsASiCS checks if the parameters are configured for ASiCS creation.
//
// Panics with the Java message when the container type is not set (Objects.requireNonNull).
//
// Port of isASiCS(ASiCParameters).
func UtilsIsASiCS(asicParameters *Parameters) bool {
	if asicParameters.ContainerType() == "" {
		panic("ASiCContainerType must be defined!")
	}
	return enumerations.ASiCContainerTypeASiCS == asicParameters.ContainerType()
}

// ASiCUtilsMimeTypeFromParameters returns a relevant MimeType for the provided parameters. Port of
// getMimeType(Parameters).
func UtilsMimeTypeFromParameters(asicParameters *Parameters) enumerations.MimeType {
	if utils.IsStringNotBlank(asicParameters.MimeType()) {
		return enumerations.MimeTypeFromMimeTypeString(asicParameters.MimeType())
	}
	if UtilsIsASiCE(asicParameters) {
		return enumerations.MimeTypeEnumASiCE
	}
	return enumerations.MimeTypeEnumASiCS
}

// UtilsFilesContainMetaInfFolder checks if the list of filenames contains a document within
// the /META-INF folder. Port of filesContainMetaInfFolder(List).
func UtilsFilesContainMetaInfFolder(filenames []string) bool {
	for _, filename := range filenames {
		if strings.HasPrefix(filename, ASiCUtilsMetaInfFolder) {
			return true
		}
	}
	return false
}

// UtilsFilesContainCorrectSignatureFileWithExtension checks if the list of filenames contains
// a signature with the expected extension. Port of
// filesContainCorrectSignatureFileWithExtension(List, String).
func UtilsFilesContainCorrectSignatureFileWithExtension(filenames []string, extension string) bool {
	for _, filename := range filenames {
		if UtilsIsSignature(filename) && strings.HasSuffix(filename, extension) {
			return true
		}
	}
	return false
}

// ASiCUtilsFilesContainSignatures checks if the list of filenames contains a signature(s). Port of
// filesContainSignatures(List).
func UtilsFilesContainSignatures(filenames []string) bool {
	for _, filename := range filenames {
		if UtilsIsSignature(filename) {
			return true
		}
	}
	return false
}

// ASiCUtilsFilesContainTimestamps checks if the list of filenames contains a timestamp. Port of
// filesContainTimestamps(List).
func UtilsFilesContainTimestamps(filenames []string) bool {
	for _, filename := range filenames {
		if UtilsIsTimestamp(filename) {
			return true
		}
	}
	return false
}

// UtilsFilesContainEvidenceRecords checks if the list of filenames contains an evidence
// record. Port of filesContainEvidenceRecords(List).
func UtilsFilesContainEvidenceRecords(filenames []string) bool {
	for _, filename := range filenames {
		if UtilsIsEvidenceRecord(filename) {
			return true
		}
	}
	return false
}

// UtilsIsAsicFileContent checks if the list of filenames represents an ASiC container content.
// Port of isAsicFileContent(List).
func UtilsIsAsicFileContent(filenames []string) bool {
	return UtilsFilesContainCorrectSignatureFileWithExtension(filenames, ASiCUtilsCAdESSignatureExtension) ||
		UtilsFilesContainCorrectSignatureFileWithExtension(filenames, ASiCUtilsXMLExtension) ||
		UtilsFilesContainTimestamps(filenames) ||
		UtilsFilesContainEvidenceRecords(filenames)
}

// ASiCUtilsIsZip checks if the document is a ZIP container. Port of isZip(DSSDocument).
func UtilsIsZip(doc model.DSSDocument) (bool, error) {
	if doc == nil {
		return false, nil
	}
	if _, ok := doc.(*model.DigestDocument); ok {
		return false, nil
	}
	is, err := doc.OpenStream()
	if err != nil {
		return false, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to determine whether the document with name '%s' represents a ZIP container. Reason : %s",
			doc.Name(), err.Error()), err)
	}
	defer is.Close()
	result, err := UtilsIsZipStream(is)
	if err != nil {
		return false, model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to determine whether the document with name '%s' represents a ZIP container. Reason : %s",
			doc.Name(), err.Error()), err)
	}
	return result, nil
}

// UtilsIsZipStream checks if the given stream contains a ZIP container.
//
// Panics with the Java message when is is nil (Objects.requireNonNull).
//
// Port of isZip(InputStream).
func UtilsIsZipStream(is io.Reader) (bool, error) {
	if is == nil {
		panic("InputStream cannot be null!")
	}
	startsWith, err := utils.StartsWithStream(is, asicUtilsZipPrefix)
	if err != nil {
		return false, exception.NewIllegalInputExceptionWithCause("Unable to read the 2 first bytes", err)
	}
	return startsWith, nil
}

// ASiCUtilsIsASiC verifies whether the given document represents an ASiC container. Port of
// isASiC(DSSDocument).
func UtilsIsASiC(doc model.DSSDocument) (bool, error) {
	isZip, err := UtilsIsZip(doc)
	if err != nil {
		return false, err
	}
	if isZip {
		filenames, err := ZipUtilsInstance().ExtractEntryNames(doc)
		if err != nil {
			return false, err
		}
		return UtilsFilesContainMetaInfFolder(filenames), nil
	}
	return false, nil
}

// UtilsIsASiCWithXAdES checks if the extracted filenames represent an ASiC with XAdES content.
//
// Note: The method looks for format specific documents. It returns FALSE for shared document
// formats between XAdES and CAdES (e.g. evidence records).
//
// Port of isASiCWithXAdES(List).
func UtilsIsASiCWithXAdES(filenames []string) bool {
	return UtilsFilesContainCorrectSignatureFileWithExtension(filenames, ASiCUtilsXMLExtension)
}

// UtilsIsASiCWithCAdES checks if the extracted filenames represent an ASiC with CAdES content.
//
// Note: The method looks for format specific documents. It returns FALSE for shared document
// formats between XAdES and CAdES (e.g. evidence records).
//
// Port of isASiCWithCAdES(List).
func UtilsIsASiCWithCAdES(filenames []string) bool {
	return UtilsFilesContainCorrectSignatureFileWithExtension(filenames, ASiCUtilsCAdESSignatureExtension) ||
		UtilsFilesContainTimestamps(filenames)
}

// ASiCUtilsIsXAdES checks if the entryName is a relevant XAdES signature. Port of isXAdES(String).
func UtilsIsXAdES(entryName string) bool {
	return UtilsIsSignature(entryName) && strings.HasSuffix(entryName, ASiCUtilsXMLExtension)
}

// ASiCUtilsIsCAdES checks if the entryName is a relevant CAdES signature. Port of isCAdES(String).
func UtilsIsCAdES(entryName string) bool {
	return UtilsIsSignature(entryName) && strings.HasSuffix(entryName, ASiCUtilsCAdESSignatureExtension)
}

// ASiCUtilsIsContainerOpenDocument checks if the archive represents an OpenDocument. Port of
// isContainerOpenDocument(DSSDocument).
func UtilsIsContainerOpenDocument(archiveContainer model.DSSDocument) (bool, error) {
	mimetype, err := asicUtilsMimetypeDocument(archiveContainer)
	if err != nil {
		return false, err
	}
	if mimetype == nil {
		return false, nil
	}
	return UtilsIsOpenDocument(mimetype)
}

// asicUtilsMimetypeDocument is the port of the private getMimetypeDocument(DSSDocument).
func asicUtilsMimetypeDocument(archiveDocument model.DSSDocument) (model.DSSDocument, error) {
	documents, err := ZipUtilsInstance().ExtractContainerContent(archiveDocument)
	if err != nil {
		return nil, err
	}
	for _, doc := range documents {
		if UtilsIsMimetype(doc.Name()) {
			return doc, nil
		}
	}
	return nil, nil
}

// UtilsIsOpenDocument checks if the mimeType document defines an OpenDocument.
//
// Port of isOpenDocument(DSSDocument). Note upstream's dead null-guard: getMimeType(mimeTypeDoc) is
// evaluated before `if (mimeTypeDoc != null)`, so a null argument returns FALSE only because
// getMimeType returns null for it - reproduced here.
func UtilsIsOpenDocument(mimeTypeDoc model.DSSDocument) (bool, error) {
	mimeType, err := UtilsMimeTypeFromDocument(mimeTypeDoc)
	if err != nil {
		return false, err
	}
	if mimeTypeDoc != nil {
		return UtilsIsOpenDocumentMimeType(mimeType), nil
	}
	return false, nil
}

// UtilsAreFilesContainMimetype checks if the list of filenames contains a mimetype file. Port
// of areFilesContainMimetype(List).
func UtilsAreFilesContainMimetype(filenames []string) bool {
	for _, filename := range filenames {
		if UtilsIsMimetype(filename) {
			return true
		}
	}
	return false
}

// ASiCUtilsIsMimetype checks if the given name is a "mimetype". Port of isMimetype(String).
func UtilsIsMimetype(entryName string) bool {
	return ASiCUtilsMimeType == entryName
}

// ASiCUtilsMimeTypeFromDocument extracts and returns the MimeType from the document. Port of
// getMimeType(DSSDocument).
func UtilsMimeTypeFromDocument(mimeTypeDocument model.DSSDocument) (enumerations.MimeType, error) {
	if mimeTypeDocument == nil {
		return nil, nil
	}
	is, err := mimeTypeDocument.OpenStream()
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("Unable to read mimetype file. Reason : %s", err.Error()), err)
	}
	defer is.Close()
	byteArray, err := utils.ToByteArray(is)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("Unable to read mimetype file. Reason : %s", err.Error()), err)
	}
	mimeTypeString := string(byteArray)
	return enumerations.MimeTypeFromMimeTypeString(mimeTypeString), nil
}

// ASiCUtilsIsASiCSContainer verifies whether the given container is of ASiC-S format type. Port of
// isASiCSContainer(DSSDocument).
func UtilsIsASiCSContainer(container model.DSSDocument) (bool, error) {
	containerType, err := UtilsContainerType(container)
	if err != nil {
		return false, err
	}
	return enumerations.ASiCContainerTypeASiCS == containerType, nil
}

// ASiCUtilsIsASiCEContainer verifies whether the given container is of ASiC-E format type. Port of
// isASiCEContainer(DSSDocument).
func UtilsIsASiCEContainer(container model.DSSDocument) (bool, error) {
	containerType, err := UtilsContainerType(container)
	if err != nil {
		return false, err
	}
	return enumerations.ASiCContainerTypeASiCE == containerType, nil
}

// UtilsContainerType verifies the type of the provided container document.
//
// Panics with the Java message when archiveContainer is nil (Objects.requireNonNull).
//
// Port of getContainerType(DSSDocument).
func UtilsContainerType(archiveContainer model.DSSDocument) (enumerations.ASiCContainerType, error) {
	if archiveContainer == nil {
		panic("Archive container shall be provided!")
	}

	entryNames, err := ZipUtilsInstance().ExtractEntryNames(archiveContainer)
	if err != nil {
		return "", err
	}

	var mimetypeDocument model.DSSDocument
	if UtilsAreFilesContainMimetype(entryNames) {
		mimetypeDocument, err = asicUtilsMimetypeDocument(archiveContainer)
		if err != nil {
			return "", err
		}
	}
	zipComment, err := UtilsZipCommentFromArchiveContainer(archiveContainer)
	if err != nil {
		return "", err
	}
	signedDocumentsNumber := asicUtilsNumberOfSignedRootDocuments(entryNames)

	return asicUtilsContainerType(archiveContainer.MimeType(), mimetypeDocument, zipComment, signedDocumentsNumber)
}

// UtilsIsASiCSContainerContent verifies whether the given Content is of ASiC-S format type.
// Port of isASiCSContainer(ASiCContent).
func UtilsIsASiCSContainerContent(asicContent *Content) (bool, error) {
	containerType, err := UtilsContainerTypeOfContent(asicContent)
	if err != nil {
		return false, err
	}
	return enumerations.ASiCContainerTypeASiCS == containerType, nil
}

// UtilsIsASiCEContainerContent verifies whether the given Content is of ASiC-E format type.
// Port of isASiCEContainer(ASiCContent).
func UtilsIsASiCEContainerContent(asicContent *Content) (bool, error) {
	containerType, err := UtilsContainerTypeOfContent(asicContent)
	if err != nil {
		return false, err
	}
	return enumerations.ASiCContainerTypeASiCE == containerType, nil
}

// UtilsContainerTypeOfContent returns the container type.
//
// Panics with the Java message when asicContent is nil (Objects.requireNonNull).
//
// Port of getContainerType(ASiCContent).
func UtilsContainerTypeOfContent(asicContent *Content) (enumerations.ASiCContainerType, error) {
	if asicContent == nil {
		panic("ASiCContent shall be provided!")
	}

	if asicContent.ContainerType() != "" {
		return asicContent.ContainerType(), nil
	}
	var containerMimeType enumerations.MimeType
	if asicContent.AsicContainer() != nil {
		containerMimeType = asicContent.AsicContainer().MimeType()
	}
	return asicUtilsContainerType(containerMimeType, asicContent.MimeTypeDocument(), asicContent.ZipComment(),
		utils.CollectionSize(asicContent.RootLevelSignedDocuments()))
}

// asicUtilsNumberOfSignedRootDocuments is the port of the private
// getNumberOfSignedRootDocuments(List).
func asicUtilsNumberOfSignedRootDocuments(containerEntryNames []string) int {
	signedDocumentCounter := 0
	for _, documentName := range containerEntryNames {
		if !strings.Contains(documentName, "/") && !UtilsIsMimetype(documentName) {
			signedDocumentCounter++
		}
	}
	return signedDocumentCounter
}

// asicUtilsContainerType is the port of the private getContainerType(MimeType, DSSDocument, String,
// int).
func asicUtilsContainerType(containerMimeType enumerations.MimeType, mimetypeDocument model.DSSDocument,
	zipComment string, rootSignedDocumentsNumber int) (enumerations.ASiCContainerType, error) {
	// 1. Identify container type based on the mimetype document
	containerType, err := asicUtilsContainerTypeFromMimeTypeDocument(mimetypeDocument)
	if err != nil {
		return "", err
	}
	if containerType != "" {
		return containerType, nil
	}
	// 2. Identify container type based on the zip comment
	containerType, err = asicUtilsContainerTypeFromZipComment(zipComment)
	if err != nil {
		return "", err
	}
	if containerType != "" {
		return containerType, nil
	}
	// 3. Check if the container contains more than one document at the root level (ASiC-E)
	if rootSignedDocumentsNumber > 1 {
		return enumerations.ASiCContainerTypeASiCE, nil
	}
	// 4. Return enforced container type, when present
	containerType, err = asicUtilsContainerTypeFromMimeType(containerMimeType)
	if err != nil {
		return "", err
	}
	if containerType != "" {
		return containerType, nil
	}
	// 5. Check if the container contains one document at the root level (ASiC-S)
	// (upstream logs "Unable to define the ASiC Container type with its properties..." here)
	if rootSignedDocumentsNumber == 1 {
		containerType = enumerations.ASiCContainerTypeASiCS
	}
	// (upstream warns "The provided container does not contain signer documents on the root
	// level!" otherwise)
	return containerType, nil
}

// asicUtilsContainerTypeFromZipComment is the port of the private
// getContainerTypeFromZipComment(String).
func asicUtilsContainerTypeFromZipComment(zipComment string) (enumerations.ASiCContainerType, error) {
	if utils.IsStringNotBlank(zipComment) {
		indexOf := strings.Index(zipComment, ASiCUtilsMimeTypeComment)
		if indexOf > -1 {
			asicCommentMimeTypeString := zipComment[len(ASiCUtilsMimeTypeComment)+indexOf:]
			mimeTypeFromZipComment := enumerations.MimeTypeFromMimeTypeString(asicCommentMimeTypeString)
			return asicUtilsContainerTypeFromMimeType(mimeTypeFromZipComment)
		}
	}
	return "", nil
}

// asicUtilsContainerTypeFromMimeTypeDocument is the port of the private
// getContainerTypeFromMimeTypeDocument(DSSDocument).
func asicUtilsContainerTypeFromMimeTypeDocument(mimetype model.DSSDocument) (enumerations.ASiCContainerType, error) {
	if mimetype != nil {
		mimeTypeFromEmbeddedFile, err := UtilsMimeTypeFromDocument(mimetype)
		if err != nil {
			return "", err
		}
		return asicUtilsContainerTypeFromMimeType(mimeTypeFromEmbeddedFile)
	}
	return "", nil
}

// asicUtilsContainerTypeFromMimeType is the port of the private
// getContainerTypeFromMimeType(MimeType).
func asicUtilsContainerTypeFromMimeType(mimeType enumerations.MimeType) (enumerations.ASiCContainerType, error) {
	if UtilsIsASiCMimeType(mimeType) {
		return UtilsASiCContainerType(mimeType)
	}
	return "", nil
}

// ASiCUtilsIsManifest checks if the fileName matches a Manifest name standard. Port of
// isManifest(String).
func UtilsIsManifest(fileName string) bool {
	return strings.HasPrefix(fileName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(fileName, ASiCUtilsASiCManifestFilename) &&
		strings.HasSuffix(fileName, ASiCUtilsXMLExtension)
}

// UtilsIsArchiveManifest checks if the fileName matches an Archive Manifest name standard. Port
// of isArchiveManifest(String).
func UtilsIsArchiveManifest(fileName string) bool {
	return strings.HasPrefix(fileName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(fileName, ASiCUtilsASiCArchiveManifestFilename) &&
		strings.HasSuffix(fileName, ASiCUtilsXMLExtension)
}

// UtilsIsEvidenceRecordManifest checks if the fileName matches an Evidence Record Manifest name
// standard. Port of isEvidenceRecordManifest(String).
func UtilsIsEvidenceRecordManifest(fileName string) bool {
	return strings.HasPrefix(fileName, ASiCUtilsMetaInfFolder) &&
		strings.Contains(fileName, ASiCUtilsASiCEvidenceRecordManifestFilename) &&
		strings.HasSuffix(fileName, ASiCUtilsXMLExtension)
}

// ASiCUtilsCoversSignature checks if the manifestFile covers a signature. Port of
// coversSignature(ManifestFile).
func UtilsCoversSignature(manifestFile *model.ManifestFile) bool {
	for _, manifestEntry := range manifestFile.Entries() {
		if UtilsIsSignature(manifestEntry.Uri()) {
			return true
		}
	}
	return false
}

// UtilsAddOrReplaceDocument searches for a document in documentList with the name of
// newDocument and replaces the found entry with the updated version, or adds the document to the
// given list if no such entry has been found.
//
// Java mutates the caller's list in place; a Go slice cannot be grown in place, so the resulting
// list is returned and callers assign it back (see the mutation idiom note in asic_content.go).
//
// Port of addOrReplaceDocument(List, DSSDocument).
func UtilsAddOrReplaceDocument(documentList []model.DSSDocument, newDocument model.DSSDocument) []model.DSSDocument {
	for i := 0; i < len(documentList); i++ {
		if newDocument.Name() == documentList[i].Name() {
			documentList[i] = newDocument
			return documentList
		}
	}
	return append(documentList, newDocument)
}

// UtilsEnsureMimeTypeAndZipComment ensures the mimetype file and zip-comment are present within
// the given Content. If the entry is not defined, a new value is created from Parameters.
//
// Port of ensureMimeTypeAndZipComment(ASiCContent, ASiCParameters).
func UtilsEnsureMimeTypeAndZipComment(asicContent *Content, asicParameters *Parameters) (*Content, error) {
	if asicContent.MimeTypeDocument() == nil {
		mimeType, err := asicUtilsMimeTypeFromContent(asicContent, asicParameters)
		if err != nil {
			return nil, err
		}
		mimetypeDocument := asicUtilsCreateMimetypeDocument(mimeType)
		asicContent.SetMimeTypeDocument(mimetypeDocument)
	}
	if utils.IsStringEmpty(asicContent.ZipComment()) {
		zipComment := asicUtilsZipCommentFromContent(asicContent, asicParameters)
		asicContent.SetZipComment(zipComment)
	}
	return asicContent, nil
}

// asicUtilsMimeTypeFromContent is the port of the private getMimeType(ASiCContent, ASiCParameters).
func asicUtilsMimeTypeFromContent(asicContent *Content, asicParameters *Parameters) (enumerations.MimeType, error) {
	var mimeType enumerations.MimeType
	mimeTypeDocument := asicContent.MimeTypeDocument()
	if mimeTypeDocument != nil {
		// re-use the same mime-type when extending a container
		var err error
		mimeType, err = UtilsMimeTypeFromDocument(mimeTypeDocument)
		if err != nil {
			return nil, err
		}
	}
	if mimeType == nil {
		if asicParameters == nil {
			panic("ASiCParameters shall be present for the requested operation!")
		}
		mimeType = UtilsMimeTypeFromParameters(asicParameters)
	}
	return mimeType, nil
}

// asicUtilsZipCommentFromContent is the port of the private getZipComment(ASiCContent,
// Parameters).
func asicUtilsZipCommentFromContent(asicContent *Content, asicParameters *Parameters) string {
	zipComment := asicContent.ZipComment()
	if utils.IsStringNotEmpty(zipComment) {
		return zipComment
	} else if asicParameters != nil {
		return UtilsZipCommentFromParameters(asicParameters)
	}
	return utils.EmptyString
}

// asicUtilsCreateMimetypeDocument is the port of the private createMimetypeDocument(MimeType).
func asicUtilsCreateMimetypeDocument(mimeType enumerations.MimeType) model.DSSDocument {
	mimeTypeBytes := []byte(mimeType.MimeTypeString())
	mimetypeDocument := model.NewInMemoryDocumentWithName(mimeTypeBytes, ASiCUtilsMimeType)
	zipEntryDocument := NewContainerEntryDocument(mimetypeDocument)
	zipEntryDocument.ZipEntry().SetCompressionMethod(0) // ensure STORED compression method
	return zipEntryDocument
}

// UtilsRootLevelSignedDocuments retrieves signed documents from a root level of the container
// (used for ASiC-E container). Port of getRootLevelSignedDocuments(ASiCContent).
func UtilsRootLevelSignedDocuments(asicContent *Content) []model.DSSDocument {
	signedDocuments := asicContent.SignedDocuments()
	if utils.IsCollectionEmpty(signedDocuments) {
		return []model.DSSDocument{}
	} else if utils.CollectionSize(signedDocuments) == 1 {
		return signedDocuments
	}
	return UtilsRootLevelDocuments(signedDocuments)
}

// UtilsRootLevelDocuments returns root-level documents across the provided list of documents.
// Port of getRootLevelDocuments(List).
func UtilsRootLevelDocuments(documents []model.DSSDocument) []model.DSSDocument {
	rootDocuments := make([]model.DSSDocument, 0)
	for _, doc := range documents {
		documentName := doc.Name()
		if documentName != "" && !strings.Contains(documentName, "/") && ASiCUtilsMimeType != documentName {
			rootDocuments = append(rootDocuments, doc)
		}
	}
	return rootDocuments
}

// UtilsZipCommentFromArchiveContainer returns a zip comment from the ASiC container, or "" when
// the archive carries none (Java null).
//
// Port of getZipComment(DSSDocument), byte-scan and all - including its signed-byte commentLen
// arithmetic, whose only consequence upstream is a log line this port drops.
func UtilsZipCommentFromArchiveContainer(archiveContainer model.DSSDocument) (string, error) {
	fileLength, err := asicUtilsFileLength(archiveContainer)
	if err != nil {
		return "", err
	}
	is, err := archiveContainer.OpenStream()
	if err != nil {
		return "", model.NewDSSError(fmt.Sprintf("Unable to read content of document with name '%s'. Reason : %s",
			archiveContainer.Name(), err.Error()))
	}
	defer is.Close()

	if fileLength > asicUtilsMaxToRead {
		toSkip := fileLength - asicUtilsMaxToRead
		skipped, err := io.CopyN(io.Discard, is, toSkip)
		if err != nil || skipped != toSkip {
			return "", model.NewDSSError(fmt.Sprintf("Unable to read content of document with name '%s'. Reason : %s",
				archiveContainer.Name(), "Different amount of bytes have been skipped!"))
		}
	}

	buffer, err := spi.DSSUtilsToByteArrayFromReader(is)
	if err != nil {
		return "", model.NewDSSError(fmt.Sprintf("Unable to read content of document with name '%s'. Reason : %s",
			archiveContainer.Name(), err.Error()))
	}
	if utils.IsArrayEmpty(buffer) {
		// upstream warns "An empty container obtained! Unable to extract zip comment."
		return "", nil
	}

	length := len(buffer)

	// Check the buffer from the end
	for ii := length - 22; ii >= 0; ii-- {
		isMagicStart := true
		for jj := 0; jj < len(asicUtilsMagicDir); jj++ {
			if buffer[ii+jj] != asicUtilsMagicDir[jj] {
				isMagicStart = false
				break
			}
		}
		if isMagicStart {
			// Magic Start found!
			realLen := length - ii - 22
			// upstream compares realLen against a commentLen computed from SIGNED bytes
			// (buffer[ii+20] + buffer[ii+21] * 256), only to log a warning on mismatch;
			// the log is dropped, so the arithmetic has no observable effect here.
			if realLen == 0 {
				return "", nil
			}
			return string(buffer[ii+22 : ii+22+realLen]), nil
		}
	}
	// upstream logs "Zip comment is not found in the provided container with name '{}'"
	return "", nil
}

// asicUtilsFileLength is the port of the private getFileLength(DSSDocument).
func asicUtilsFileLength(archiveContainer model.DSSDocument) (int64, error) {
	if fileDocument, ok := archiveContainer.(*model.FileDocument); ok {
		info, err := os.Stat(fileDocument.Path())
		if err != nil {
			return 0, nil // java.io.File#length() returns 0 for an unreadable file
		}
		return info.Size(), nil
	}

	is, err := archiveContainer.OpenStream()
	if err != nil {
		return 0, model.NewDSSErrorMessageCause("Unable to compute archive size", err)
	}
	defer is.Close()
	size, err := utils.GetInputStreamSize(is)
	if err != nil {
		return 0, model.NewDSSErrorMessageCause("Unable to compute archive size", err)
	}
	return size, nil
}

// UtilsToSimpleManifestEntries transforms a list of given documents to a list of "simple" (only
// basic information) manifest entries.
//
// Panics with the Java message when a document is nil (Objects.requireNonNull).
//
// Port of toSimpleManifestEntries(List).
func UtilsToSimpleManifestEntries(documents []model.DSSDocument) []*model.ManifestEntry {
	entries := make([]*model.ManifestEntry, 0)
	for _, doc := range documents {
		if doc == nil {
			panic("DSSDocument cannot be null!")
		}
		entry := model.NewManifestEntry()
		entry.SetUri(doc.Name())
		entry.SetMimeType(doc.MimeType())
		entry.SetFound(true)
		entry.SetDocument(doc)
		entries = append(entries, entry)
	}
	return entries
}

// UtilsIsCoveredByManifest checks if a document (e.g. a signature) with the given filename is
// covered by a manifest. Port of isCoveredByManifest(List, String).
func UtilsIsCoveredByManifest(manifestDocuments []model.DSSDocument, filename string) bool {
	if utils.IsCollectionNotEmpty(manifestDocuments) {
		for _, archiveManifest := range manifestDocuments {
			manifestFile := ManifestParserGetManifestFile(archiveManifest)
			if manifestFile != nil {
				for _, entry := range manifestFile.Entries() {
					if filename != "" && filename == entry.Uri() {
						return true
					}
				}
			}
		}
	}
	return false
}
