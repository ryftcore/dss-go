// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/AbstractCacheAccess.java (DSS 6.5.RC1).
package job

// AbstractCacheAccess reads the relevant cache by the given key.
//
// slf4j trace logging is dropped (no observable behavior).
type AbstractCacheAccess struct {
	// downloadCache is the global download Cache.
	downloadCache *DownloadCache

	// parsingCache is the global parsing Cache.
	parsingCache *ParsingCache

	// validationCache is the global validation Cache.
	validationCache *ValidationCache
}

// NewAbstractCacheAccess creates an AbstractCacheAccess. Port of the protected constructor.
func NewAbstractCacheAccess(fileCache *DownloadCache, parsingCache *ParsingCache, validationCache *ValidationCache) AbstractCacheAccess {
	return AbstractCacheAccess{downloadCache: fileCache, parsingCache: parsingCache, validationCache: validationCache}
}

// GetDownloadCacheEntry returns the download cache DTO result. Port of
// getDownloadCacheEntry(CacheKey).
func (a *AbstractCacheAccess) GetDownloadCacheEntry(key CacheKey) *CachedEntry[DownloadResult] {
	return a.downloadCache.Get(key)
}

// GetParsingCacheEntry returns the parsing cache DTO result. Port of
// getParsingCacheEntry(CacheKey).
func (a *AbstractCacheAccess) GetParsingCacheEntry(key CacheKey) *CachedEntry[ParsingResult] {
	return a.parsingCache.Get(key)
}

// GetValidationCacheEntry returns the validation cache DTO result. Port of
// getValidationCacheEntry(CacheKey).
func (a *AbstractCacheAccess) GetValidationCacheEntry(key CacheKey) *CachedEntry[ValidationResult] {
	return a.validationCache.Get(key)
}

// GetAllCacheKeys returns all found keys in any cache. Port of getAllCacheKeys().
func (a *AbstractCacheAccess) GetAllCacheKeys() map[CacheKey]struct{} {
	keys := make(map[CacheKey]struct{})
	for _, k := range a.downloadCache.Keys() {
		keys[k] = struct{}{}
	}
	for _, k := range a.parsingCache.Keys() {
		keys[k] = struct{}{}
	}
	for _, k := range a.validationCache.Keys() {
		keys[k] = struct{}{}
	}
	return keys
}
