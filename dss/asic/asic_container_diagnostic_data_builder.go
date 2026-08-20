// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/ASiCContainerDiagnosticDataBuilder.java (DSS 6.5.RC1).
package asic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	dssdiagnostic "github.com/utain/esig/dss/validation/reports/diagnostic"
)

// ASiCContainerDiagnosticDataBuilder is the DiagnosticDataBuilder for an ASiC container. Port
// of the class ASiCContainerDiagnosticDataBuilder, extending
// validation/reports/diagnostic.SignedDocumentDiagnosticDataBuilder.
type ASiCContainerDiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder

	// containerInfo is the information about the validating container.
	containerInfo *model.ContainerInfo
}

// NewASiCContainerDiagnosticDataBuilder instantiates a builder with nil/zero values. Port of
// the default constructor.
func NewASiCContainerDiagnosticDataBuilder() *ASiCContainerDiagnosticDataBuilder {
	b := &ASiCContainerDiagnosticDataBuilder{
		SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// ContainerInfo sets the container info (ASiC). Ports containerInfo(ContainerInfo).
func (b *ASiCContainerDiagnosticDataBuilder) ContainerInfo(containerInfo *model.ContainerInfo) *ASiCContainerDiagnosticDataBuilder {
	b.containerInfo = containerInfo
	return b
}

// Build ports the @Override build().
func (b *ASiCContainerDiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	diagnosticData := b.SignedDocumentDiagnosticDataBuilder.Build()
	diagnosticData.ContainerInfo = b.getXmlContainerInfo()
	return diagnosticData
}

// getXmlContainerInfo ports the private getXmlContainerInfo().
func (b *ASiCContainerDiagnosticDataBuilder) getXmlContainerInfo() *jaxb.XmlContainerInfo {
	if b.containerInfo == nil {
		return nil
	}
	xmlContainerInfo := &jaxb.XmlContainerInfo{}
	containerType := jaxb.ASiCContainerTypeValue(b.containerInfo.ContainerType())
	xmlContainerInfo.ContainerType = &containerType
	if zipComment := b.containerInfo.ZipComment(); utils.IsStringNotBlank(zipComment) {
		xmlContainerInfo.ZipComment = &zipComment
	}
	mimeTypeFilePresent := b.containerInfo.IsMimeTypeFilePresent()
	xmlContainerInfo.MimeTypeFilePresent = &mimeTypeFilePresent
	if mimeTypeContent := b.containerInfo.MimeTypeContent(); mimeTypeContent != "" {
		xmlContainerInfo.MimeTypeContent = &mimeTypeContent
	}
	xmlContainerInfo.ContentFiles = &jaxb.ContentFilesWrapper{Items: b.containerInfo.SignedDocumentFilenames()}
	xmlContainerInfo.ManifestFiles = &jaxb.ManifestFilesWrapper{Items: b.getXmlManifests(b.containerInfo.ManifestFiles())}
	return xmlContainerInfo
}

// getXmlManifests ports the private getXmlManifests(List).
func (b *ASiCContainerDiagnosticDataBuilder) getXmlManifests(manifestFiles []*model.ManifestFile) []*jaxb.XmlManifestFile {
	xmlManifests := make([]*jaxb.XmlManifestFile, 0, len(manifestFiles))
	if utils.IsCollectionNotEmpty(manifestFiles) {
		for _, manifestFile := range manifestFiles {
			xmlManifest := &jaxb.XmlManifestFile{}
			if filename := manifestFile.Filename(); filename != "" {
				xmlManifest.Filename = &filename
			}
			if signatureFilename := manifestFile.SignatureFilename(); signatureFilename != "" {
				xmlManifest.SignatureFilename = &signatureFilename
			}
			entries := make([]string, 0, len(manifestFile.Entries()))
			for _, entry := range manifestFile.Entries() {
				entries = append(entries, entry.Uri())
			}
			xmlManifest.Entries = &jaxb.EntriesWrapper{Items: entries}
			xmlManifests = append(xmlManifests, xmlManifest)
		}
	}
	return xmlManifests
}
