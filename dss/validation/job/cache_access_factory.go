// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/CacheAccessFactory.java (DSS 6.5.RC1).
package job

// CacheAccessFactory is a factory used to create objects to interact with the cache.
type CacheAccessFactory interface {
	// GetCacheAccess loads a class to deal with a cache by the key records. Port of
	// getCacheAccess(CacheKey).
	GetCacheAccess(key CacheKey) CacheAccessByKey

	// GetDocumentChangesCacheAccess loads a class for document updates. Port of
	// getDocumentChangesCacheAccess().
	GetDocumentChangesCacheAccess() *ChangesCacheAccess

	// GetReadOnlyCacheAccess loads a read-only cache access. Port of
	// getReadOnlyCacheAccess().
	GetReadOnlyCacheAccess() ReadOnlyCacheAccess

	// GetSynchronizerCacheAccess loads a cache access to synchronize records. Port of
	// getSynchronizerCacheAccess().
	GetSynchronizerCacheAccess() *SynchronizerCacheAccess

	// GetDebugCacheAccess loads a cache access to load the information about the current
	// cache state. Port of getDebugCacheAccess().
	GetDebugCacheAccess() *DebugCacheAccess
}
