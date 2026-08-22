// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/AbstractASiCDataToSignHelperBuilder.java (DSS 6.5.RC1).
package asic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// AbstractASiCDataToSignHelperBuilderOverrides declares the operation
// AbstractASiCDataToSignHelperBuilder calls back into virtually from CreatePackageZip. Every
// concrete builder must call InitAbstractASiCDataToSignHelperBuilder with itself before use, or
// the base's virtual calls will not reach the override.
type AbstractASiCDataToSignHelperBuilderOverrides interface {
	// GetDataPackageName returns a name for a package zip container, containing the original
	// signer data. Port of the protected abstract getDataPackageName(ASiCContent).
	GetDataPackageName(asicContent *ASiCContent) string
}

// AbstractASiCDataToSignHelperBuilder builds a relevant GetDataToSignASiCWithCAdESHelper for
// ASiC container dataToSign creation.
type AbstractASiCDataToSignHelperBuilder struct {
	// overrides points back at the concrete builder; see InitAbstractASiCDataToSignHelperBuilder.
	overrides AbstractASiCDataToSignHelperBuilderOverrides
}

// NewAbstractASiCDataToSignHelperBuilderBase builds the base state a subclass embeds. Port of
// the protected empty constructor. The subclass constructor must follow it with
// InitAbstractASiCDataToSignHelperBuilder.
func NewAbstractASiCDataToSignHelperBuilderBase() AbstractASiCDataToSignHelperBuilder {
	return AbstractASiCDataToSignHelperBuilder{}
}

// InitAbstractASiCDataToSignHelperBuilder registers the concrete builder with its base so the
// base can dispatch GetDataPackageName. Every concrete builder constructor must call this once.
func (b *AbstractASiCDataToSignHelperBuilder) InitAbstractASiCDataToSignHelperBuilder(overrides AbstractASiCDataToSignHelperBuilderOverrides) {
	b.overrides = overrides
}

func (b *AbstractASiCDataToSignHelperBuilder) requireOverrides() AbstractASiCDataToSignHelperBuilderOverrides {
	if b.overrides == nil {
		panic("AbstractASiCDataToSignHelperBuilder was not initialised: the concrete builder must call InitAbstractASiCDataToSignHelperBuilder in its constructor")
	}
	return b.overrides
}

// GetASiCSSignedDocument returns a document to be signed in case of an ASiC-S container. Ports
// the protected getASiCSSignedDocument(List, Date).
//
// Panics with Java's IllegalArgumentException message when no file to sign was provided.
func (b *AbstractASiCDataToSignHelperBuilder) GetASiCSSignedDocument(filesToBeSigned []model.DSSDocument, signingDate time.Time) model.DSSDocument {
	switch {
	case len(filesToBeSigned) == 1:
		return filesToBeSigned[0]
	case len(filesToBeSigned) > 1:
		return b.CreatePackageZip(filesToBeSigned, signingDate)
	default:
		panic("At least one file to be signed shall be provided!")
	}
}

// CreatePackageZip creates a zip with all files to be signed. Ports
// createPackageZip(List, Date).
//
// ZipUtils.getInstance().createZipArchive(List, Date, String) (Java passes a nil zipComment
// here); ASiCContent exposes SetContainerDocuments([]model.DSSDocument). Panics on a zip
// creation error, matching the panic-on-Build-error convention used throughout this port.
func (b *AbstractASiCDataToSignHelperBuilder) CreatePackageZip(documents []model.DSSDocument, signingDate time.Time) model.DSSDocument {
	packageZip, err := ZipUtilsInstance().CreateZipArchiveFromEntriesAt(documents, signingDate, "")
	if err != nil {
		panic(err)
	}

	asicContent := NewASiCContent()
	asicContent.SetContainerDocuments(documents)
	packageZip.SetName(b.requireOverrides().GetDataPackageName(asicContent))

	packageZip.SetMimeType(enumerations.MimeTypeEnumZIP)
	return packageZip
}
