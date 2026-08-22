// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/DownloadCache.java (DSS 6.5.RC1).
package job

import (
	"slices"

	"github.com/ryftcore/dss-go/dss/utils"
)

// DownloadCache stores downloaded files.
//
// slf4j trace/warn logging is dropped (no observable behavior).
type DownloadCache struct {
	AbstractCache[DownloadResult]
}

// NewDownloadCache creates an empty DownloadCache. Port of the default constructor.
func NewDownloadCache() *DownloadCache {
	c := &DownloadCache{AbstractCache: NewAbstractCache[DownloadResult]()}
	c.InitAbstractCache(c)
	return c
}

// IsUpToDate checks if the file with the given cacheKey is up to date. Port of
// isUpToDate(CacheKey, DownloadResult).
func (c *DownloadCache) IsUpToDate(cacheKey CacheKey, downloadedResult DownloadResult) bool {
	if c.IsToBeDeleted(cacheKey) {
		return false
	}

	cachedFileEntry := c.Get(cacheKey)
	if !cachedFileEntry.IsEmpty() {
		cachedResult := cachedFileEntry.CachedResult()
		digestMatch := cachedResult.Digest().Equals(downloadedResult.Digest())
		sha2ContentMatch := isSHA2ContentMatch(cachedResult, downloadedResult)
		upToDate := digestMatch && sha2ContentMatch
		if upToDate {
			cachedFileEntry.SyncUpdateDate()
		}
		return upToDate
	}
	return false
}

// isSHA2ContentMatch mirrors isSHA2ContentMatch(DownloadResult, DownloadResult) verbatim,
// including its asymmetric condition (both operands check cachedResult's messages are empty;
// the Java source never inspects downloadedResult's list on the equals() branch).
func isSHA2ContentMatch(cachedResult DownloadResult, downloadedResult DownloadResult) bool {
	return (utils.IsCollectionEmpty(cachedResult.Sha2ErrorMessages()) && utils.IsCollectionEmpty(downloadedResult.Sha2ErrorMessages())) ||
		(utils.IsCollectionEmpty(cachedResult.Sha2ErrorMessages()) && slices.Equal(cachedResult.Sha2ErrorMessages(), downloadedResult.Sha2ErrorMessages()))
}

// CacheType returns CacheTypeDownload. Port of getCacheType().
func (c *DownloadCache) CacheType() CacheType {
	return CacheTypeDownload
}
