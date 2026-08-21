// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/SynchronizerCacheAccess.java (DSS 6.5.RC1).
package job

// SynchronizerCacheAccess synchronizes all caches for the given key.
//
// slf4j debug logging is dropped (no observable behavior).
type SynchronizerCacheAccess struct {
	AbstractCacheAccess
}

// NewSynchronizerCacheAccess creates a SynchronizerCacheAccess. Port of the default
// constructor.
func NewSynchronizerCacheAccess(downloadCache *DownloadCache, parsingCache *ParsingCache, validationCache *ValidationCache) *SynchronizerCacheAccess {
	return &SynchronizerCacheAccess{AbstractCacheAccess: NewAbstractCacheAccess(downloadCache, parsingCache, validationCache)}
}

// Sync synchronizes all records for the key. Port of sync(CacheKey).
func (s *SynchronizerCacheAccess) Sync(key CacheKey) {
	if s.downloadCache.IsDesync(key) {
		s.downloadCache.Sync(key)
	}
	if s.parsingCache.IsDesync(key) {
		s.parsingCache.Sync(key)
	}
	if s.validationCache.IsDesync(key) {
		s.validationCache.Sync(key)
	}
}
