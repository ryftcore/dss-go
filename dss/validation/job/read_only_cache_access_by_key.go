// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/ReadOnlyCacheAccessByKey.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// ReadOnlyCacheAccessByKey provides a read only interface for cache by key.
type ReadOnlyCacheAccessByKey interface {
	// GetDownloadReadOnlyResult returns the cached read-only download result DTO. Port of
	// getDownloadReadOnlyResult().
	GetDownloadReadOnlyResult() modeljob.DownloadInfoRecord
	// GetParsingReadOnlyResult returns the cached read-only parsing result DTO. Port of
	// getParsingReadOnlyResult().
	GetParsingReadOnlyResult() modeljob.ParsingInfoRecord
	// GetValidationReadOnlyResult returns the cached read-only validation result DTO. Port
	// of getValidationReadOnlyResult().
	GetValidationReadOnlyResult() modeljob.ValidationInfoRecord
}
