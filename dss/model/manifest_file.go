// Ported from dss-model/.../ManifestFile.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// ManifestFile represents a parsed Manifest File object.
type ManifestFile struct {
	// document is the DSSDocument represented by the ManifestFile.
	document DSSDocument

	// signatureFilename is the name of a signature or timestamp associated
	// to the ManifestFile.
	signatureFilename string

	// entries is the list of entries present in the document.
	entries []*ManifestEntry

	// manifestType defines the type of the manifest document.
	manifestType enumerations.ASiCManifestTypeEnum
}

// NewManifestFile instantiates the object with null values. Ports the
// default constructor.
func NewManifestFile() *ManifestFile {
	return &ManifestFile{}
}

// Document gets the DSSDocument representing the manifest.
func (m *ManifestFile) Document() DSSDocument { return m.document }

// SetDocument sets the manifest document.
func (m *ManifestFile) SetDocument(document DSSDocument) { m.document = document }

// Filename gets the manifest document's filename. Panics if no document is
// set (Java calls document.getName() unguarded, i.e. NullPointerException).
func (m *ManifestFile) Filename() string {
	return m.document.Name()
}

// SignatureFilename gets the signature filename.
func (m *ManifestFile) SignatureFilename() string { return m.signatureFilename }

// SetSignatureFilename sets the signature filename.
func (m *ManifestFile) SetSignatureFilename(signatureFilename string) {
	m.signatureFilename = signatureFilename
}

// DigestValue gets the digest value of the manifest document for
// digestAlgorithm.
func (m *ManifestFile) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	return m.document.DigestValue(digestAlgorithm)
}

// Entries gets the list of ManifestEntry values, lazily initializing it.
// Ports ManifestFile#getEntries.
func (m *ManifestFile) Entries() []*ManifestEntry {
	if m.entries == nil {
		m.entries = []*ManifestEntry{}
	}
	return m.entries
}

// SetEntries sets the list of ManifestEntry values.
func (m *ManifestFile) SetEntries(entries []*ManifestEntry) { m.entries = entries }

// ManifestType gets the type of the ASiC Manifest file.
func (m *ManifestFile) ManifestType() enumerations.ASiCManifestTypeEnum { return m.manifestType }

// SetManifestType sets the type of the ASiC Manifest file.
func (m *ManifestFile) SetManifestType(manifestType enumerations.ASiCManifestTypeEnum) {
	m.manifestType = manifestType
}

// RootFile returns the ManifestEntry with Rootfile()==true, or nil if none
// is found. Ports ManifestFile#getRootFile.
func (m *ManifestFile) RootFile() *ManifestEntry {
	for _, entry := range m.Entries() {
		if entry.IsRootfile() {
			return entry
		}
	}
	return nil
}

// IsDocumentCovered checks if the document with documentName is covered by
// the Manifest.
func (m *ManifestFile) IsDocumentCovered(documentName string) bool {
	if documentName != "" {
		for _, entry := range m.Entries() {
			if documentName == entry.Uri() {
				return true
			}
		}
	}
	return false
}
