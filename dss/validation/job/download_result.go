// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/download/DownloadResult.java (DSS 6.5.RC1).
package job

import "github.com/utain/esig/dss/model"

// DownloadResult provides methods to extract information about a download job.
type DownloadResult interface {
	CachedResult

	// DSSDocument returns the downloaded document. Port of getDSSDocument().
	DSSDocument() model.DSSDocument
	// Digest returns the digest of a canonicalized document. Port of getDigest().
	Digest() model.Digest
	// Sha2ErrorMessages returns error messages occurred during sha2 processing, if
	// applicable: an empty list if none occurred. Port of getSha2ErrorMessages().
	//
	// TODO : remove from the interface (carried over from the Java TODO).
	Sha2ErrorMessages() []string
}
