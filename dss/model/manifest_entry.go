// Ported from dss-model/.../ManifestEntry.java (DSS 6.5.RC1).
package model

import "github.com/utain/esig/dss/enumerations"

// ManifestEntry defines a referenced document entry of a ManifestFile.
type ManifestEntry struct {
	// uri is the reference URI.
	uri string

	// mimeType is the mimetype of the entry.
	mimeType enumerations.MimeType

	// digest is the digest of the referenced entry.
	digest Digest

	// document is the matching document, when found.
	document DSSDocument

	// found defines if the referenced data is found (used for reference
	// validation).
	found bool

	// intact defines if the referenced data is intact (digest matches)
	// (used for reference validation).
	intact bool

	// rootfile defines if it is the root file.
	rootfile bool
}

// NewManifestEntry instantiates the object with null values. Ports the
// default constructor.
func NewManifestEntry() *ManifestEntry {
	return &ManifestEntry{}
}

// Uri gets the filename.
func (m *ManifestEntry) Uri() string { return m.uri }

// SetUri sets the filename.
func (m *ManifestEntry) SetUri(uri string) { m.uri = uri }

// MimeType gets the mimetype.
func (m *ManifestEntry) MimeType() enumerations.MimeType { return m.mimeType }

// SetMimeType sets the mimetype.
func (m *ManifestEntry) SetMimeType(mimeType enumerations.MimeType) { m.mimeType = mimeType }

// Digest gets the manifest entry digest.
func (m *ManifestEntry) Digest() Digest { return m.digest }

// SetDigest sets the manifest entry digest.
func (m *ManifestEntry) SetDigest(digest Digest) { m.digest = digest }

// SetDocumentName is deprecated: it no longer sets anything. Ports the
// deprecated ManifestEntry#setDocumentName, which upstream logs a warning
// and skips processing. Use SetDocument instead.
//
// upstream logs "Use of deprecated method #setDocumentName. Please switch
// to #setDocument. Current method processing is skipped" here.
func (m *ManifestEntry) SetDocumentName(documentName string) {
	// intentionally a no-op, matching upstream's deprecated behavior
}

// Document gets the corresponding document.
func (m *ManifestEntry) Document() DSSDocument { return m.document }

// SetDocument sets the corresponding document.
func (m *ManifestEntry) SetDocument(document DSSDocument) { m.document = document }

// IsFound gets if the referenced document has been found.
func (m *ManifestEntry) IsFound() bool { return m.found }

// SetFound sets if the referenced document has been found.
func (m *ManifestEntry) SetFound(found bool) { m.found = found }

// IsIntact gets if the digest of the reference document matches.
func (m *ManifestEntry) IsIntact() bool { return m.intact }

// SetIntact sets if the digest of the reference document matches.
func (m *ManifestEntry) SetIntact(intact bool) { m.intact = intact }

// IsRootfile checks if it is a rootfile.
func (m *ManifestEntry) IsRootfile() bool { return m.rootfile }

// SetRootfile sets if the entry is a rootfile.
func (m *ManifestEntry) SetRootfile(rootfile bool) { m.rootfile = rootfile }
