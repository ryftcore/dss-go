// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCWithCAdESUtils.java (DSS 6.5.RC1).
//
// slf4j LOG calls are dropped per PORTING.md; the branch it sits in is preserved.
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithCAdESUtilsDefaultArchiveManifestFilename is the default Archive Manifest filename.
// Port of DEFAULT_ARCHIVE_MANIFEST_FILENAME.
const ASiCWithCAdESUtilsDefaultArchiveManifestFilename = asic.ASiCUtilsMetaInfFolder + asic.ASiCUtilsASiCArchiveManifestFilename + asic.ASiCUtilsXMLExtension

// ASiCWithCAdESUtilsGetSignedDocument returns a list of signed documents by a signature with a
// given signatureFilename. Ports the static getSignedDocument(ASiCContent, String).
func ASiCWithCAdESUtilsGetSignedDocument(extractResult *asic.ASiCContent, signatureFilename string) model.DSSDocument {
	containerType := extractResult.ContainerType()
	if enumerations.ASiCContainerType_ASiC_S == containerType && len(extractResult.RootLevelSignedDocuments()) == 1 {
		return extractResult.RootLevelSignedDocuments()[0] // Collection size should be equal 1

	} else if enumerations.ASiCContainerType_ASiC_E == containerType {
		// the manifest file is signed
		manifestDocuments := extractResult.ManifestDocuments()
		signatureDocuments := extractResult.SignatureDocuments()
		if utils.CollectionSize(manifestDocuments) == 1 && utils.CollectionSize(signatureDocuments) == 1 {
			return manifestDocuments[0]
		}
		// we need to check the manifest file and its digest
		linkedManifest := asic.ASiCManifestParserGetLinkedManifest(extractResult.ManifestDocuments(), signatureFilename)
		if linkedManifest != nil {
			return linkedManifest
		}
		return nil // related manifest not found
	}
	// (upstream warns "Unable to extract a signed document. Reason : Unknown asic container type.")
	return nil
}
