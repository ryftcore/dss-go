// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/ChangesCacheAccess.java (DSS 6.5.RC1).
package job

// ChangesCacheAccess forces an update of a document cache information.
type ChangesCacheAccess struct {
	// downloadCache is the global download Cache.
	downloadCache *DownloadCache

	// parsingCache is the global parsing Cache.
	parsingCache *ParsingCache

	// validationCache is the global validation Cache.
	validationCache *ValidationCache
}

// NewChangesCacheAccess creates a ChangesCacheAccess. Port of the default constructor.
func NewChangesCacheAccess(downloadCache *DownloadCache, parsingCache *ParsingCache, validationCache *ValidationCache) *ChangesCacheAccess {
	return &ChangesCacheAccess{downloadCache: downloadCache, parsingCache: parsingCache, validationCache: validationCache}
}

// ToBeDeleted sets 'toBeDeleted' status for all records with the given key. Port of
// toBeDeleted(CacheKey).
func (c *ChangesCacheAccess) ToBeDeleted(cacheKey CacheKey) {
	c.downloadCache.ToBeDeleted(cacheKey)
	c.parsingCache.ToBeDeleted(cacheKey)
	c.validationCache.ToBeDeleted(cacheKey)
}

// ExpireSignatureValidation sets the expired status for the validation record for the
// cacheKey. Port of expireSignatureValidation(CacheKey).
func (c *ChangesCacheAccess) ExpireSignatureValidation(cacheKey CacheKey) {
	c.validationCache.Expire(cacheKey)
}
