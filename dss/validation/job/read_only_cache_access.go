// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/ReadOnlyCacheAccess.java (DSS 6.5.RC1).
package job

import modeljob "github.com/utain/esig/dss/model/job"

// ReadOnlyCacheAccess accesses the cache in read-only mode.
type ReadOnlyCacheAccess interface {
	// GetDownloadInfoRecord returns the download cache DTO result. Port of
	// getDownloadInfoRecord(CacheKey).
	GetDownloadInfoRecord(key CacheKey) modeljob.DownloadInfoRecord
	// GetParsingInfoRecord returns the parsing cache DTO result. Port of
	// getParsingInfoRecord(CacheKey).
	GetParsingInfoRecord(key CacheKey) modeljob.ParsingInfoRecord
	// GetValidationInfoRecord returns the validation cache DTO result. Port of
	// getValidationInfoRecord(CacheKey).
	GetValidationInfoRecord(key CacheKey) modeljob.ValidationInfoRecord
	// GetAllCacheKeys returns all found keys in any cache. Port of getAllCacheKeys().
	GetAllCacheKeys() map[CacheKey]struct{}
}
