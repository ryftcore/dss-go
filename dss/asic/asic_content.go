// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCContent.java
// (DSS 6.5.RC1).
//
// MUTATION IDIOM (relevant to every package building on this one): upstream mutates the container
// in place through its getters - `asicContent.getSignatureDocuments().add(doc)`,
// `...addAll(docs)`. A Go slice is a value, so a slice returned by a getter cannot be appended to
// in place. The Go spelling of those statements is therefore
//
//	asicContent.SetSignatureDocuments(append(asicContent.SignatureDocuments(), doc))
//
// which has the same effect. No `Add*` convenience methods are introduced, so that the ported call
// sites stay a line-for-line match against the Java they came from.
//
// java.io.Serializable is dropped silently (no Go counterpart).
package asic

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCContent contains grouped documents representing an ASiC container's content.
type ASiCContent struct {
	// asicContainer is the original ASiC container.
	asicContainer model.DSSDocument

	// containerType is the container type.
	containerType enumerations.ASiCContainerType

	// zipComment is the zip comment.
	zipComment string

	// mimeTypeDocument is the mimetype document.
	mimeTypeDocument model.DSSDocument

	// signedDocuments is the list of originally signed documents embedded into the container.
	signedDocuments []model.DSSDocument

	// signatureDocuments is the list of signature documents embedded into the container.
	signatureDocuments []model.DSSDocument

	// manifestDocuments is the list of manifest documents embedded into the container.
	manifestDocuments []model.DSSDocument

	// archiveManifestDocuments is the list of archive manifest documents embedded into the
	// container (ASiC with CAdES).
	archiveManifestDocuments []model.DSSDocument

	// evidenceRecordManifestDocuments is the list of evidence record manifest documents
	// embedded into the container.
	evidenceRecordManifestDocuments []model.DSSDocument

	// timestampDocuments is the list of timestamp documents embedded into the container (ASiC
	// with CAdES).
	timestampDocuments []model.DSSDocument

	// evidenceRecordDocuments is the list of evidence record documents embedded into the
	// container.
	evidenceRecordDocuments []model.DSSDocument

	// unsupportedDocuments is the list of unsupported documents embedded into the container.
	unsupportedDocuments []model.DSSDocument

	// folders is the list of folders embedded into the container.
	folders []model.DSSDocument

	// containerDocuments is the list of "package.zip" documents (ASiC-S).
	containerDocuments []model.DSSDocument
}

// NewASiCContent instantiates an object with null values and empty lists of documents. Port of the
// default constructor (Java's field initializers included).
func NewASiCContent() *ASiCContent {
	return &ASiCContent{
		signedDocuments:                 []model.DSSDocument{},
		signatureDocuments:              []model.DSSDocument{},
		manifestDocuments:               []model.DSSDocument{},
		archiveManifestDocuments:        []model.DSSDocument{},
		evidenceRecordManifestDocuments: []model.DSSDocument{},
		timestampDocuments:              []model.DSSDocument{},
		evidenceRecordDocuments:         []model.DSSDocument{},
		unsupportedDocuments:            []model.DSSDocument{},
		folders:                         []model.DSSDocument{},
		containerDocuments:              []model.DSSDocument{},
	}
}

// AsicContainer gets the original ASiC container. Port of getAsicContainer().
func (c *ASiCContent) AsicContainer() model.DSSDocument { return c.asicContainer }

// SetAsicContainer sets the original ASiC container. Port of setAsicContainer(DSSDocument).
func (c *ASiCContent) SetAsicContainer(asicContainer model.DSSDocument) {
	c.asicContainer = asicContainer
}

// ContainerType gets the container type. Port of getContainerType().
func (c *ASiCContent) ContainerType() enumerations.ASiCContainerType { return c.containerType }

// SetContainerType sets the container type. Port of setContainerType(ASiCContainerType).
func (c *ASiCContent) SetContainerType(containerType enumerations.ASiCContainerType) {
	c.containerType = containerType
}

// ZipComment gets the zip comment. Port of getZipComment().
func (c *ASiCContent) ZipComment() string { return c.zipComment }

// SetZipComment sets the zip comment. Port of setZipComment(String).
func (c *ASiCContent) SetZipComment(zipComment string) { c.zipComment = zipComment }

// MimeTypeDocument gets the mimetype document. Port of getMimeTypeDocument().
func (c *ASiCContent) MimeTypeDocument() model.DSSDocument { return c.mimeTypeDocument }

// SetMimeTypeDocument sets the mimetype document. Port of setMimeTypeDocument(DSSDocument).
func (c *ASiCContent) SetMimeTypeDocument(mimeTypeDocument model.DSSDocument) {
	c.mimeTypeDocument = mimeTypeDocument
}

// SignatureDocuments gets the signature documents. Port of getSignatureDocuments().
func (c *ASiCContent) SignatureDocuments() []model.DSSDocument { return c.signatureDocuments }

// SetSignatureDocuments sets the signature documents. Port of setSignatureDocuments(List).
func (c *ASiCContent) SetSignatureDocuments(signatureDocuments []model.DSSDocument) {
	c.signatureDocuments = signatureDocuments
}

// ManifestDocuments gets the manifest documents. Port of getManifestDocuments().
func (c *ASiCContent) ManifestDocuments() []model.DSSDocument { return c.manifestDocuments }

// SetManifestDocuments sets the manifest documents. Port of setManifestDocuments(List).
func (c *ASiCContent) SetManifestDocuments(manifestDocuments []model.DSSDocument) {
	c.manifestDocuments = manifestDocuments
}

// ArchiveManifestDocuments gets the archive manifest documents (ASiC with CAdES only). Port of
// getArchiveManifestDocuments().
func (c *ASiCContent) ArchiveManifestDocuments() []model.DSSDocument {
	return c.archiveManifestDocuments
}

// SetArchiveManifestDocuments sets the archive manifest documents (ASiC with CAdES only). Port of
// setArchiveManifestDocuments(List).
func (c *ASiCContent) SetArchiveManifestDocuments(archiveManifestDocuments []model.DSSDocument) {
	c.archiveManifestDocuments = archiveManifestDocuments
}

// EvidenceRecordManifestDocuments gets the evidence record manifest documents. Port of
// getEvidenceRecordManifestDocuments().
func (c *ASiCContent) EvidenceRecordManifestDocuments() []model.DSSDocument {
	return c.evidenceRecordManifestDocuments
}

// SetEvidenceRecordManifestDocuments sets a list of evidence record manifest documents. Port of
// setEvidenceRecordManifestDocuments(List).
func (c *ASiCContent) SetEvidenceRecordManifestDocuments(evidenceRecordManifestDocuments []model.DSSDocument) {
	c.evidenceRecordManifestDocuments = evidenceRecordManifestDocuments
}

// TimestampDocuments gets the timestamp documents (ASiC with CAdES only). Port of
// getTimestampDocuments().
func (c *ASiCContent) TimestampDocuments() []model.DSSDocument { return c.timestampDocuments }

// SetTimestampDocuments sets the timestamp documents (ASiC with CAdES only). Port of
// setTimestampDocuments(List).
func (c *ASiCContent) SetTimestampDocuments(timestampDocuments []model.DSSDocument) {
	c.timestampDocuments = timestampDocuments
}

// EvidenceRecordDocuments gets the evidence record documents. Port of
// getEvidenceRecordDocuments().
func (c *ASiCContent) EvidenceRecordDocuments() []model.DSSDocument {
	return c.evidenceRecordDocuments
}

// SetEvidenceRecordDocuments sets a list of evidence record documents. Port of
// setEvidenceRecordDocuments(List).
func (c *ASiCContent) SetEvidenceRecordDocuments(evidenceRecordDocuments []model.DSSDocument) {
	c.evidenceRecordDocuments = evidenceRecordDocuments
}

// SignedDocuments gets the signed documents. Port of getSignedDocuments().
func (c *ASiCContent) SignedDocuments() []model.DSSDocument { return c.signedDocuments }

// SetSignedDocuments sets the signed documents. Port of setSignedDocuments(List).
func (c *ASiCContent) SetSignedDocuments(signedDocuments []model.DSSDocument) {
	c.signedDocuments = signedDocuments
}

// UnsupportedDocuments gets the unsupported documents. Port of getUnsupportedDocuments().
func (c *ASiCContent) UnsupportedDocuments() []model.DSSDocument { return c.unsupportedDocuments }

// SetUnsupportedDocuments sets the unsupported documents. Port of setUnsupportedDocuments(List).
func (c *ASiCContent) SetUnsupportedDocuments(unsupportedDocuments []model.DSSDocument) {
	c.unsupportedDocuments = unsupportedDocuments
}

// Folders returns a list of folders present within the container. Port of getFolders().
func (c *ASiCContent) Folders() []model.DSSDocument { return c.folders }

// SetFolders sets a list of folders present within an archive. Port of setFolders(List).
func (c *ASiCContent) SetFolders(folders []model.DSSDocument) { c.folders = folders }

// ContainerDocuments gets the "package.zip" documents. Port of getContainerDocuments().
func (c *ASiCContent) ContainerDocuments() []model.DSSDocument { return c.containerDocuments }

// SetContainerDocuments sets the "package.zip" documents. Port of setContainerDocuments(List).
func (c *ASiCContent) SetContainerDocuments(containerDocuments []model.DSSDocument) {
	c.containerDocuments = containerDocuments
}

// RootLevelSignedDocuments returns a list of documents at the root level within the container.
// Port of getRootLevelSignedDocuments().
func (c *ASiCContent) RootLevelSignedDocuments() []model.DSSDocument {
	if utils.IsCollectionEmpty(c.SignedDocuments()) {
		return []model.DSSDocument{}
	}
	rootLevelDocuments := make([]model.DSSDocument, 0)
	for _, doc := range c.SignedDocuments() {
		if doc.Name() != "" && !strings.Contains(doc.Name(), "/") && !strings.Contains(doc.Name(), "\\") {
			rootLevelDocuments = append(rootLevelDocuments, doc)
		}
	}
	return rootLevelDocuments
}

// AllManifestDocuments returns a list of all found manifest documents. Port of
// getAllManifestDocuments().
func (c *ASiCContent) AllManifestDocuments() []model.DSSDocument {
	allManifestsList := make([]model.DSSDocument, 0)
	allManifestsList = append(allManifestsList, c.ManifestDocuments()...)
	allManifestsList = append(allManifestsList, c.ArchiveManifestDocuments()...)
	allManifestsList = append(allManifestsList, c.EvidenceRecordManifestDocuments()...)
	return allManifestsList
}

// AllDocuments gets all documents. Port of getAllDocuments().
func (c *ASiCContent) AllDocuments() []model.DSSDocument {
	allDocuments := make([]model.DSSDocument, 0)
	// "mimetype" shall be the first file in the ASiC container;
	if c.mimeTypeDocument != nil {
		allDocuments = append(allDocuments, c.mimeTypeDocument)
	}
	if utils.IsCollectionNotEmpty(c.signedDocuments) {
		allDocuments = append(allDocuments, c.signedDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.signatureDocuments) {
		allDocuments = append(allDocuments, c.signatureDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.manifestDocuments) {
		allDocuments = append(allDocuments, c.manifestDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.archiveManifestDocuments) {
		allDocuments = append(allDocuments, c.archiveManifestDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.evidenceRecordManifestDocuments) {
		allDocuments = append(allDocuments, c.evidenceRecordManifestDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.timestampDocuments) {
		allDocuments = append(allDocuments, c.timestampDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.evidenceRecordDocuments) {
		allDocuments = append(allDocuments, c.evidenceRecordDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.unsupportedDocuments) {
		allDocuments = append(allDocuments, c.unsupportedDocuments...)
	}
	if utils.IsCollectionNotEmpty(c.folders) {
		allDocuments = append(allDocuments, c.folders...)
	}

	return allDocuments
}
