// Ported from dss-model/.../ContainerInfo.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// ContainerInfo contains information about an ASiC container.
type ContainerInfo struct {
	// containerType is the ASiC container type.
	containerType enumerations.ASiCContainerType

	// zipComment is the ZIP comment of the archive.
	zipComment string

	// mimeTypeContent is the mimetype file content.
	mimeTypeContent string

	// signedDocumentFilenames is the list of signed document filenames.
	signedDocumentFilenames []string

	// manifestFiles is the list of embedded manifest files.
	manifestFiles []*ManifestFile
}

// NewContainerInfo instantiates the object with null values. Ports the
// default constructor.
func NewContainerInfo() *ContainerInfo {
	return &ContainerInfo{}
}

// ContainerType gets the ASiCContainerType.
func (c *ContainerInfo) ContainerType() enumerations.ASiCContainerType { return c.containerType }

// SetContainerType sets the ASiCContainerType.
func (c *ContainerInfo) SetContainerType(containerType enumerations.ASiCContainerType) {
	c.containerType = containerType
}

// ZipComment gets the zip comment.
func (c *ContainerInfo) ZipComment() string { return c.zipComment }

// SetZipComment sets the zip comment.
func (c *ContainerInfo) SetZipComment(zipComment string) { c.zipComment = zipComment }

// MimeTypeContent gets the mimetype file content.
func (c *ContainerInfo) MimeTypeContent() string { return c.mimeTypeContent }

// SetMimeTypeContent sets the mimetype file content.
func (c *ContainerInfo) SetMimeTypeContent(mimeTypeContent string) {
	c.mimeTypeContent = mimeTypeContent
}

// IsMimeTypeFilePresent returns whether the mimetype file is present.
func (c *ContainerInfo) IsMimeTypeFilePresent() bool {
	return c.mimeTypeContent != ""
}

// SignedDocumentFilenames gets the list of signed document filenames.
func (c *ContainerInfo) SignedDocumentFilenames() []string { return c.signedDocumentFilenames }

// SetSignedDocumentFilenames sets the signed document filenames.
func (c *ContainerInfo) SetSignedDocumentFilenames(signedDocumentFilenames []string) {
	c.signedDocumentFilenames = signedDocumentFilenames
}

// ManifestFiles gets the list of manifest files.
func (c *ContainerInfo) ManifestFiles() []*ManifestFile { return c.manifestFiles }

// SetManifestFiles sets the list of manifest files.
func (c *ContainerInfo) SetManifestFiles(manifestFiles []*ManifestFile) {
	c.manifestFiles = manifestFiles
}
