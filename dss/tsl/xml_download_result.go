// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/download/XmlDownloadResult.java (DSS 6.5.RC1).
//
// Implements job.DownloadResult (Java
// eu.europa.esig.dss.validation.job.download.DownloadResult): GetDSSDocument()/GetDigest()/
// GetSha2ErrorMessages() are spelled DSSDocument()/Digest()/Sha2ErrorMessages() here, per
// PORTING.md's get-prefix-dropped naming.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// XmlDownloadResult defines the download result.
type XmlDownloadResult struct {
	// dssDocument is the downloaded document.
	dssDocument model.DSSDocument

	// digest is the digest of the canonicalized document.
	digest model.Digest
}

var _ job.DownloadResult = (*XmlDownloadResult)(nil)

// NewXmlDownloadResult is the default constructor. Port of XmlDownloadResult(DSSDocument, Digest).
func NewXmlDownloadResult(dssDocument model.DSSDocument, digest model.Digest) *XmlDownloadResult {
	return &XmlDownloadResult{dssDocument: dssDocument, digest: digest}
}

// DSSDocument gets the downloaded document. Port of getDSSDocument().
func (r *XmlDownloadResult) DSSDocument() model.DSSDocument {
	return r.dssDocument
}

// Digest gets digest of the canonicalized document. Port of getDigest().
func (r *XmlDownloadResult) Digest() model.Digest {
	return r.digest
}

// Sha2ErrorMessages returns error messages occurred during sha2 processing, if applicable. Port
// of getSha2ErrorMessages().
func (r *XmlDownloadResult) Sha2ErrorMessages() []string {
	if documentWithSha2, ok := r.dssDocument.(*DocumentWithSha2); ok {
		return documentWithSha2.Errors()
	}
	return nil
}
