// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESDataToSignHelperBuilder.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithCAdESDataToSignHelperBuilderOverrides declares the operation Java's abstract
// ASiCWithCAdESDataToSignHelperBuilder leaves to its subclasses and calls back into from
// Build(). Per S7_BRIEF.md's virtual-dispatch warning, every concrete builder must call
// InitASiCWithCAdESDataToSignHelperBuilder with itself before use.
type ASiCWithCAdESDataToSignHelperBuilderOverrides interface {
	// GetManifestBuilder returns an AbstractASiCManifestBuilder to be used for a
	// signed/timestamped manifest creation. Port of the protected abstract
	// getManifestBuilder(ASiCContent, ASiCWithCAdESCommonParameters).
	//
	// Java's covariant overrides narrow the return type to ASiCEWithCAdESManifestBuilder; in Go
	// the subclasses hand back the *asic.AbstractASiCManifestBuilder their concrete builder
	// embeds, which already carries the concrete builder as its registered overrides - so
	// Build() on it dispatches to the concrete getSigReferenceMimeType()/getManifestFilename()
	// exactly as Java does.
	GetManifestBuilder(asicContent *asic.ASiCContent, parameters ASiCWithCAdESCommonParameters) *asic.AbstractASiCManifestBuilder
}

// ASiCWithCAdESDataToSignHelperBuilder builds a relevant GetDataToSignASiCWithCAdESHelper for
// ASiC with CAdES dataToSign creation.
type ASiCWithCAdESDataToSignHelperBuilder struct {
	AbstractASiCWithCAdESDataToSignHelperBuilder

	// overrides points back at the concrete builder; see
	// InitASiCWithCAdESDataToSignHelperBuilder.
	overrides ASiCWithCAdESDataToSignHelperBuilderOverrides
}

// NewASiCWithCAdESDataToSignHelperBuilder builds the base state a subclass embeds. Ports the
// protected ASiCWithCAdESDataToSignHelperBuilder(ASiCWithCAdESFilenameFactory) constructor; the
// concrete leaf constructor must follow it with InitASiCWithCAdESDataToSignHelperBuilder and
// InitAbstractASiCDataToSignHelperBuilder.
func NewASiCWithCAdESDataToSignHelperBuilder(asicFilenameFactory ASiCWithCAdESFilenameFactory) ASiCWithCAdESDataToSignHelperBuilder {
	return ASiCWithCAdESDataToSignHelperBuilder{
		AbstractASiCWithCAdESDataToSignHelperBuilder: NewAbstractASiCWithCAdESDataToSignHelperBuilder(asicFilenameFactory),
	}
}

// InitASiCWithCAdESDataToSignHelperBuilder registers the concrete builder with this base so
// Build() can dispatch GetManifestBuilder. Every concrete builder constructor must call this
// once.
func (b *ASiCWithCAdESDataToSignHelperBuilder) InitASiCWithCAdESDataToSignHelperBuilder(overrides ASiCWithCAdESDataToSignHelperBuilderOverrides) {
	b.overrides = overrides
}

func (b *ASiCWithCAdESDataToSignHelperBuilder) requireOverrides() ASiCWithCAdESDataToSignHelperBuilderOverrides {
	if b.overrides == nil {
		panic("ASiCWithCAdESDataToSignHelperBuilder was not initialised: the concrete builder must call InitASiCWithCAdESDataToSignHelperBuilder in its constructor")
	}
	return b.overrides
}

// Build creates a GetDataToSignASiCWithCAdESHelper from an ASiCContent. Ports
// build(ASiCContent, ASiCWithCAdESCommonParameters).
func (b *ASiCWithCAdESDataToSignHelperBuilder) Build(asicContent *asic.ASiCContent,
	parameters ASiCWithCAdESCommonParameters) GetDataToSignASiCWithCAdESHelper {
	asicContent, err := asic.ASiCUtilsEnsureMimeTypeAndZipComment(asicContent, parameters.ASiC())
	if err != nil {
		panic(err)
	}
	if b.IsASiCArchive(asicContent) {
		return b.fromArchive(asicContent, parameters)
	}
	return b.fromFiles(asicContent, parameters)
}

// fromArchive ports the private fromArchive(ASiCContent, ASiCWithCAdESCommonParameters).
//
// Panics with Java's UnsupportedOperationException message on a container type mismatch.
func (b *ASiCWithCAdESDataToSignHelperBuilder) fromArchive(asicContent *asic.ASiCContent,
	parameters ASiCWithCAdESCommonParameters) GetDataToSignASiCWithCAdESHelper {
	currentContainerType := asicContent.ContainerType()

	asice := asic.ASiCUtilsIsASiCE(parameters.ASiC())
	switch {
	case asice && enumerations.ASiCContainerType_ASiC_E == currentContainerType:
		manifestDocument := b.createManifestDocument(asicContent, parameters)
		return NewDataToSignASiCEWithCAdESHelper(asicContent, manifestDocument)

	case !asice && enumerations.ASiCContainerType_ASiC_S == currentContainerType:
		return NewDataToSignASiCSWithCAdESFromArchive(asicContent)

	default:
		panic(fmt.Sprintf("Original container type '%s' vs parameter : '%s'",
			currentContainerType, parameters.ASiC().ContainerType()))
	}
}

// fromFiles ports the private fromFiles(ASiCContent, ASiCWithCAdESCommonParameters).
func (b *ASiCWithCAdESDataToSignHelperBuilder) fromFiles(asicContent *asic.ASiCContent,
	parameters ASiCWithCAdESCommonParameters) GetDataToSignASiCWithCAdESHelper {
	if asic.ASiCUtilsIsASiCE(parameters.ASiC()) {
		asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_E)
		manifestDocument := b.createManifestDocument(asicContent, parameters)
		return NewDataToSignASiCEWithCAdESHelper(asicContent, manifestDocument)
	}

	asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_S)
	asicsSignedDocument := b.GetASiCSSignedDocument(asicContent.SignedDocuments(), parameters.ZipCreationDate())
	asicContent.SetSignedDocuments([]model.DSSDocument{asicsSignedDocument})
	return NewDataToSignASiCSWithCAdESFromFiles(asicContent)
}

// createManifestDocument ports the private
// createManifestDocument(ASiCContent, ASiCWithCAdESCommonParameters).
func (b *ASiCWithCAdESDataToSignHelperBuilder) createManifestDocument(asicContent *asic.ASiCContent,
	parameters ASiCWithCAdESCommonParameters) model.DSSDocument {
	manifestDocument, err := b.requireOverrides().GetManifestBuilder(asicContent, parameters).Build()
	if err != nil {
		panic(err)
	}
	return manifestDocument
}
