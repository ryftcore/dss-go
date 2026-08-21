// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DocumentInfo.java (DSS 6.5.RC1).
package job

import "github.com/utain/esig/dss/model"

// DocumentInfo contains a validation result for a document. P is the parent DocumentInfo
// type, mirroring the Java self-bound type parameter "P extends DocumentInfo<P>"; Go has no
// F-bounded polymorphism for interfaces, so the bound is documented rather than enforced.
type DocumentInfo[P any] interface {
	model.IdentifierBasedObject

	// DownloadCacheInfo returns Download Cache Info. Port of getDownloadCacheInfo().
	DownloadCacheInfo() DownloadInfoRecord
	// ParsingCacheInfo returns Parsing Cache Info. Port of getParsingCacheInfo().
	ParsingCacheInfo() ParsingInfoRecord
	// ValidationCacheInfo returns Validation Cache Info. Port of getValidationCacheInfo().
	ValidationCacheInfo() ValidationInfoRecord
	// Url returns a URL that was used to download the remote file. Port of getUrl().
	//
	// INTEGRATION FIX: named Url (not the more Go-idiomatic URL) to match the naming this
	// codebase's sole real implementer (model/tsl.TLInfo, and every one of its ~28 callers
	// across dss/tsl and dss/validation) already uses throughout.
	Url() string
	// Parent returns the DocumentInfo referencing the current Trusted List. Port of
	// getParent().
	Parent() P
	// DSSIDAsString returns the String representation of the identifier. Port of
	// getDSSIdAsString().
	DSSIDAsString() string
}
