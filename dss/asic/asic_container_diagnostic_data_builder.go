//go:build phase8

// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/ASiCContainerDiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation and dss/validation/diagnostic do not exist yet. Drop the tag once Phase 8 lands
// those packages; the file needs no other change (see the cades/cades_diagnostic_data_builder.go
// precedent this file follows).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.reports.diagnostic.
// SignedDocumentDiagnosticDataBuilder, and the JAXB diagnostic-data model it builds
// (eu.europa.esig.dss.diagnostic.jaxb.{XmlContainerInfo,XmlDiagnosticData,XmlManifestFile})
// belong to dss-validation (+*-report-jaxb), assigned to the not-yet-ported
// `validation`/`validation/diagnostic` packages (Phase 8). There is no Go type to embed or build
// here yet - the method bodies below are ported 1:1 against the package path and shape the
// cades/cades_diagnostic_data_builder.go precedent implies, so this file needs no further
// changes once Phase 8 lands the package.
package asic

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// ASiCContainerDiagnosticDataBuilder is the DiagnosticDataBuilder for an ASiC container. Port
// of the class ASiCContainerDiagnosticDataBuilder, extending
// validation/diagnostic.SignedDocumentDiagnosticDataBuilder.
type ASiCContainerDiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder

	// containerInfo is the information about the validating container.
	containerInfo *model.ContainerInfo
}

// NewASiCContainerDiagnosticDataBuilder instantiates a builder with nil/zero values. Port of
// the default constructor.
func NewASiCContainerDiagnosticDataBuilder() *ASiCContainerDiagnosticDataBuilder {
	return &ASiCContainerDiagnosticDataBuilder{}
}

// ContainerInfo sets the container info (ASiC). Ports containerInfo(ContainerInfo).
func (b *ASiCContainerDiagnosticDataBuilder) ContainerInfo(containerInfo *model.ContainerInfo) *ASiCContainerDiagnosticDataBuilder {
	b.containerInfo = containerInfo
	return b
}

// Build ports the @Override build(). Shadows the embedded base's method of the same name; see
// the file header's forward-dependency note and PORTING.md's "Virtual dispatch" precedent on
// why the override has to be reproduced this way rather than relying on embedding alone.
func (b *ASiCContainerDiagnosticDataBuilder) Build() *dssdiagnostic.XmlDiagnosticData {
	diagnosticData := b.SignedDocumentDiagnosticDataBuilder.Build()
	diagnosticData.SetContainerInfo(b.getXmlContainerInfo())
	return diagnosticData
}

// getXmlContainerInfo ports the private getXmlContainerInfo().
func (b *ASiCContainerDiagnosticDataBuilder) getXmlContainerInfo() *dssdiagnostic.XmlContainerInfo {
	if b.containerInfo == nil {
		return nil
	}
	xmlContainerInfo := dssdiagnostic.NewXmlContainerInfo()
	xmlContainerInfo.SetContainerType(b.containerInfo.ContainerType())
	zipComment := b.containerInfo.ZipComment()
	if utils.IsStringNotBlank(zipComment) {
		xmlContainerInfo.SetZipComment(zipComment)
	}
	xmlContainerInfo.SetMimeTypeFilePresent(b.containerInfo.IsMimeTypeFilePresent())
	xmlContainerInfo.SetMimeTypeContent(b.containerInfo.MimeTypeContent())
	xmlContainerInfo.SetContentFiles(b.containerInfo.SignedDocumentFilenames())
	xmlContainerInfo.SetManifestFiles(b.getXmlManifests(b.containerInfo.ManifestFiles()))
	return xmlContainerInfo
}

// getXmlManifests ports the private getXmlManifests(List).
func (b *ASiCContainerDiagnosticDataBuilder) getXmlManifests(manifestFiles []*model.ManifestFile) []*dssdiagnostic.XmlManifestFile {
	if utils.IsCollectionEmpty(manifestFiles) {
		return []*dssdiagnostic.XmlManifestFile{}
	}
	xmlManifests := make([]*dssdiagnostic.XmlManifestFile, 0, len(manifestFiles))
	for _, manifestFile := range manifestFiles {
		xmlManifest := dssdiagnostic.NewXmlManifestFile()
		xmlManifest.SetFilename(manifestFile.Filename())
		xmlManifest.SetSignatureFilename(manifestFile.SignatureFilename())
		for _, entry := range manifestFile.Entries() {
			xmlManifest.Entries = append(xmlManifest.Entries, entry.Uri())
		}
		xmlManifests = append(xmlManifests, xmlManifest)
	}
	return xmlManifests
}
