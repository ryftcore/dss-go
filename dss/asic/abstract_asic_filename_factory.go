// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/AbstractASiCFilenameFactory.java
// (DSS 6.5.RC1).
//
// Java's protected members are exported here: the concrete subclasses live in the
// dss/asic/{cades,xades} packages, so Go visibility rules leave no narrower option.
//
// Every self-called protected method is routed through AbstractASiCFilenameFactoryOverrides (the
// Overrides + Init pattern used throughout the ported tree), because static Go dispatch would
// otherwise drop a subclass override silently.
//
// java.io.Serializable is dropped silently (no Go counterpart).
package asic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// AbstractASiCFilenameFactoryOverrides captures the protected methods of
// AbstractASiCFilenameFactory that the class calls on itself and a subclass may override.
type AbstractASiCFilenameFactoryOverrides interface {
	// IsAvailableName verifies whether the filename is not present within restrictedNames.
	// Port of the protected isAvailableName(String, Collection).
	IsAvailableName(filename string, restrictedNames []string) bool

	// WithMetaInfFolder appends a "META-INF/" string to the filename, when required. Port of
	// the protected getWithMetaInfFolder(String).
	WithMetaInfFolder(filename string) string

	// AssertFilenameValid verifies whether the given filename represents a valid document name
	// within a container. Port of the protected assertFilenameValid(String, List).
	AssertFilenameValid(filename string, documentsOfType []model.DSSDocument) error
}

// AbstractASiCFilenameFactory contains common methods for ASiC filename factory.
type AbstractASiCFilenameFactory struct {
	// overrides points back at the concrete factory; see InitAbstractASiCFilenameFactory.
	overrides AbstractASiCFilenameFactoryOverrides
}

// InitAbstractASiCFilenameFactory wires the self-call dispatch to the concrete factory. Port of
// the protected default constructor, extended with the overrides argument Go needs.
func (f *AbstractASiCFilenameFactory) InitAbstractASiCFilenameFactory(overrides AbstractASiCFilenameFactoryOverrides) {
	f.overrides = overrides
}

// abstractASiCFilenameFactoryOverrides returns the dispatch target, defaulting to this struct's own
// implementations when a subclass never called InitAbstractASiCFilenameFactory.
func (f *AbstractASiCFilenameFactory) abstractASiCFilenameFactoryOverrides() AbstractASiCFilenameFactoryOverrides {
	if f.overrides != nil {
		return f.overrides
	}
	return f
}

// NextAvailableDocumentName returns the next available suffix (i.e. number) for a new file entry
// across given document names.
//
// Ex.: For provided two signature documents, the method will return "003" suffix, to be used for a
// signature file creation with a name "META-INF/signature003.xml".
//
// Port of the protected getNextAvailableDocumentName(String, Collection).
func (f *AbstractASiCFilenameFactory) NextAvailableDocumentName(nameTemplate string, existingDocumentNames []string) string {
	// Use set to exclude duplicated names
	seen := make(map[string]struct{}, len(existingDocumentNames))
	distinct := make([]string, 0, len(existingDocumentNames))
	for _, name := range existingDocumentNames {
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			distinct = append(distinct, name)
		}
	}
	return f.documentNameRecursively(nameTemplate, distinct)
}

// documentNameRecursively is the port of the private getDocumentNameRecursively(String,
// Collection).
func (f *AbstractASiCFilenameFactory) documentNameRecursively(nameTemplate string, existingDocumentNames []string) string {
	number := len(existingDocumentNames) + 1
	numberStr := strconv.Itoa(number)
	zeroPad := "000"
	candidateSuffix := numberStr
	if len(numberStr) < len(zeroPad) {
		candidateSuffix = zeroPad[len(numberStr):] + numberStr // 2 -> 002
	}
	candidateName := strings.ReplaceAll(nameTemplate, "001", candidateSuffix)
	if !f.abstractASiCFilenameFactoryOverrides().IsAvailableName(candidateName, existingDocumentNames) {
		// to list to allow increment; increase the amount of entries within the list
		grown := make([]string, 0, len(existingDocumentNames)+1)
		grown = append(grown, existingDocumentNames...)
		grown = append(grown, candidateName)
		return f.documentNameRecursively(nameTemplate, grown)
	}
	return candidateName
}

// IsAvailableName verifies whether the filename is not present within restrictedNames. Port of the
// protected isAvailableName(String, Collection).
func (f *AbstractASiCFilenameFactory) IsAvailableName(filename string, restrictedNames []string) bool {
	for _, restrictedName := range restrictedNames {
		if restrictedName == filename {
			return false
		}
	}
	return true
}

// WithMetaInfFolder appends a "META-INF/" string to the filename, when required. Port of the
// protected getWithMetaInfFolder(String).
func (f *AbstractASiCFilenameFactory) WithMetaInfFolder(filename string) string {
	if !strings.HasPrefix(filename, ASiCUtilsMetaInfFolder) {
		filename = ASiCUtilsMetaInfFolder + filename
	}
	return filename
}

// AssertASiCContentIsValid verifies whether the provided asicContent contains all the required
// information for a new filename determination.
//
// Panics with the Java messages when asicContent or its container type is missing
// (Objects.requireNonNull); returns an error for any other container type
// (IllegalArgumentException).
//
// Port of the protected assertASiCContentIsValid(ASiCContent).
func (f *AbstractASiCFilenameFactory) AssertASiCContentIsValid(asicContent *Content) error {
	if asicContent == nil {
		panic("ASiCContent shall be provided!")
	}
	if asicContent.ContainerType() == "" {
		panic("Type of ASiC Container shall be defined!")
	}
	if enumerations.ASiCContainerTypeASiCS != asicContent.ContainerType() &&
		enumerations.ASiCContainerTypeASiCE != asicContent.ContainerType() {
		return fmt.Errorf("The type of the ASiCContent shall be one of ASiC-S or ASiC-E!")
	}
	return nil
}

// AssertFilenameValid verifies whether the given filename represents a valid document name within a
// container. Port of the protected assertFilenameValid(String, List).
func (f *AbstractASiCFilenameFactory) AssertFilenameValid(filename string, documentsOfType []model.DSSDocument) error {
	if !f.abstractASiCFilenameFactoryOverrides().IsAvailableName(filename, spi.DSSUtilsDocumentNames(documentsOfType)) {
		return exception.NewIllegalInputException(fmt.Sprintf(
			"The filename '%s' cannot be used, as a document of the same name is already present within the container!",
			filename))
	}
	return nil
}

// ValidDataPackageFilename returns a valid data package filename. Port of the protected
// getValidDataPackageFilename(String, Content).
func (f *AbstractASiCFilenameFactory) ValidDataPackageFilename(dataPackageFilename string, asicContent *Content) (string, error) {
	if err := f.abstractASiCFilenameFactoryOverrides().AssertFilenameValid(dataPackageFilename, asicContent.SignedDocuments()); err != nil {
		return "", err
	}
	if strings.Contains(dataPackageFilename, "/") {
		return "", fmt.Errorf("A data package file within ASiC container shall be on the root level!")
	} else if !strings.HasSuffix(strings.ToLower(dataPackageFilename), ".zip") {
		return "", fmt.Errorf("A data package filename within ASiC container shall ends with '.zip'!")
	}
	return dataPackageFilename, nil
}

// ValidEvidenceRecordManifestFilename returns a valid evidence record manifest filename.
//
// NOTE: The name of the file shall be: "META-INF/ASiCEvidenceRecordManifest*.xml". "META-INF/" is
// optional.
//
// Port of the protected getValidEvidenceRecordManifestFilename(String, ASiCContent).
func (f *AbstractASiCFilenameFactory) ValidEvidenceRecordManifestFilename(evidenceRecordManifestFilename string, asicContent *Content) (string, error) {
	evidenceRecordManifestFilename = f.abstractASiCFilenameFactoryOverrides().WithMetaInfFolder(evidenceRecordManifestFilename)
	if err := f.abstractASiCFilenameFactoryOverrides().AssertFilenameValid(evidenceRecordManifestFilename, asicContent.EvidenceRecordManifestDocuments()); err != nil {
		return "", err
	}
	if !strings.HasPrefix(evidenceRecordManifestFilename, ASiCUtilsMetaInfFolder+ASiCUtilsASiCEvidenceRecordManifestFilename) ||
		!strings.HasSuffix(evidenceRecordManifestFilename, ASiCUtilsXMLExtension) {
		return "", fmt.Errorf("ASiC evidence record manifest file within ASiC container shall match the template '%s'!",
			ASiCUtilsMetaInfFolder+ASiCUtilsASiCEvidenceRecordManifestFilename+"*"+ASiCUtilsXMLExtension)
	}
	return evidenceRecordManifestFilename, nil
}
