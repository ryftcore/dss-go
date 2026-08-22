// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESDataToSignHelperBuilder.java
// (DSS 6.5.RC1).
package xades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithXAdESDataToSignHelperBuilder builds a relevant GetDataToSignASiCWithXAdESHelper for
// ASiC with XAdES dataToSign creation.
//
// Java's `extends AbstractASiCDataToSignHelperBuilder` becomes embedding plus the
// InitAbstractASiCDataToSignHelperBuilder(self) registration; the base dispatches
// GetDataPackageName (implemented right here, since Java's class is concrete - there is no
// intermediate abstract subclass the way dss-asic-cades has one) through
// asic.AbstractASiCDataToSignHelperBuilderOverrides.
type ASiCWithXAdESDataToSignHelperBuilder struct {
	asic.AbstractASiCDataToSignHelperBuilder

	// asicFilenameFactory defines rules for filename creation for new data package file. Java
	// declares the field protected final; the whole dss-asic-xades module is one Go package, so
	// the field stays unexported and is read directly within this package.
	asicFilenameFactory ASiCWithXAdESFilenameFactory
}

var _ asic.AbstractASiCDataToSignHelperBuilderOverrides = (*ASiCWithXAdESDataToSignHelperBuilder)(nil)

// NewASiCWithXAdESDataToSignHelperBuilder is the default constructor. Ports
// ASiCWithXAdESDataToSignHelperBuilder(ASiCWithXAdESFilenameFactory).
func NewASiCWithXAdESDataToSignHelperBuilder(asicFilenameFactory ASiCWithXAdESFilenameFactory) *ASiCWithXAdESDataToSignHelperBuilder {
	b := &ASiCWithXAdESDataToSignHelperBuilder{
		AbstractASiCDataToSignHelperBuilder: asic.NewAbstractASiCDataToSignHelperBuilderBase(),
		asicFilenameFactory:                 asicFilenameFactory,
	}
	b.InitAbstractASiCDataToSignHelperBuilder(b)
	return b
}

// Build creates a GetDataToSignASiCWithXAdESHelper from an ASiCContent. Ports
// build(ASiCContent, ASiCWithXAdESSignatureParameters).
//
// Panics with Java's UnsupportedOperationException message on a container type mismatch, and
// re-raises the error from ASiCUtilsEnsureMimeTypeAndZipComment/ASiCUtilsIsOpenDocument as a
// panic (Java's underlying DSSException propagates the same way).
func (b *ASiCWithXAdESDataToSignHelperBuilder) Build(asicContent *asic.ASiCContent,
	parameters *ASiCWithXAdESSignatureParameters) GetDataToSignASiCWithXAdESHelper {
	asicContent, err := asic.ASiCUtilsEnsureMimeTypeAndZipComment(asicContent, parameters.ASiC())
	if err != nil {
		panic(err)
	}

	isOpenDocument, err := asic.ASiCUtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		panic(err)
	}
	if isOpenDocument {
		return NewDataToSignOpenDocumentHelper(asicContent)
	}

	// if ASiC with XAdES (no detached timestamps are allowed)
	if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) {
		currentContainerType := asicContent.ContainerType()

		asice := asic.ASiCUtilsIsASiCE(parameters.ASiC())
		switch {
		case asice && enumerations.ASiCContainerTypeASiCE == currentContainerType:
			return NewDataToSignASiCEWithXAdESHelper(asicContent)
		case !asice && enumerations.ASiCContainerTypeASiCS == currentContainerType:
			return NewDataToSignASiCSWithXAdESHelper(asicContent)
		default:
			panic(fmt.Sprintf("Original container type '%s' vs parameter : '%s'",
				currentContainerType, parameters.ASiC().ContainerType()))
		}
	}

	return b.fromFiles(asicContent, parameters)
}

// fromFiles ports the private fromFiles(ASiCContent, ASiCWithXAdESSignatureParameters).
func (b *ASiCWithXAdESDataToSignHelperBuilder) fromFiles(asicContent *asic.ASiCContent,
	parameters *ASiCWithXAdESSignatureParameters) GetDataToSignASiCWithXAdESHelper {
	if asic.ASiCUtilsIsASiCE(parameters.ASiC()) {
		asicManifest := b.createASiCManifest(asicContent)
		asicContent.SetManifestDocuments(append(asicContent.ManifestDocuments(), asicManifest))
		return NewDataToSignASiCEWithXAdESHelper(asicContent)
	}

	var signingDate time.Time
	if bLevelSigningDate := parameters.BLevel().SigningDate(); bLevelSigningDate != nil {
		signingDate = *bLevelSigningDate
	}
	asicsSignedDocument := b.GetASiCSSignedDocument(asicContent.SignedDocuments(), signingDate)
	asicContent.SetSignedDocuments([]model.DSSDocument{asicsSignedDocument})
	return NewDataToSignASiCSWithXAdESHelper(asicContent)
}

// createASiCManifest returns the ASiC Manifest. Ports the private
// createASiCManifest(ASiCContent).
func (b *ASiCWithXAdESDataToSignHelperBuilder) createASiCManifest(asicContent *asic.ASiCContent) model.DSSDocument {
	manifestDocument, err := NewASiCEWithXAdESManifestBuilder().
		SetDocuments(asicContent.SignedDocuments()).
		SetManifestFilename(b.asicFilenameFactory.ManifestFilename(asicContent)).
		Build()
	if err != nil {
		panic(err)
	}
	return manifestDocument
}

// GetDataPackageName ports the @Override protected getDataPackageName(ASiCContent).
func (b *ASiCWithXAdESDataToSignHelperBuilder) GetDataPackageName(asicContent *asic.ASiCContent) string {
	return b.asicFilenameFactory.DataPackageFilename(asicContent)
}
