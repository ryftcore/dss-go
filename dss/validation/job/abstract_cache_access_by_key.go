// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/AbstractCacheAccessByKey.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// AbstractCacheAccessByKey is the abstract implementation to access a cache record by a
// specified key. D, P, V mirror Java's "<D extends DownloadInfoRecord, P extends
// ParsingInfoRecord, V extends ValidationInfoRecord>". It implements CacheAccessByKey (via
// the base-interface-returning read methods inherited from AbstractReadOnlyCacheAccessByKey,
// see that type's DEVIATION note).
type AbstractCacheAccessByKey[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord] struct {
	AbstractReadOnlyCacheAccessByKey[D, P, V]

	// downloadCache is the global download Cache.
	downloadCache *DownloadCache

	// parsingCache is the global parsing Cache.
	parsingCache *ParsingCache

	// validationCache is the global validation Cache.
	validationCache *ValidationCache
}

// NewAbstractCacheAccessByKey creates an AbstractCacheAccessByKey. Port of the protected
// constructor.
func NewAbstractCacheAccessByKey[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord](
	key CacheKey, downloadCache *DownloadCache, parsingCache *ParsingCache, validationCache *ValidationCache,
	readOnlyCacheAccess ParametrizedReadOnlyCacheAccess[D, P, V]) AbstractCacheAccessByKey[D, P, V] {
	return AbstractCacheAccessByKey[D, P, V]{
		AbstractReadOnlyCacheAccessByKey: NewAbstractReadOnlyCacheAccessByKey(key, readOnlyCacheAccess),
		downloadCache:                    downloadCache,
		parsingCache:                     parsingCache,
		validationCache:                  validationCache,
	}
}

// GetCacheKey returns the CacheKey. Port of getCacheKey().
func (a *AbstractCacheAccessByKey[D, P, V]) GetCacheKey() CacheKey {
	return a.key
}

// IsUpToDate checks if the download result is up to date for the given key. Port of
// isUpToDate(DownloadResult).
func (a *AbstractCacheAccessByKey[D, P, V]) IsUpToDate(downloadResult DownloadResult) bool {
	return a.downloadCache.IsUpToDate(a.key, downloadResult)
}

// UpdateDownloadResult updates the download result. Port of update(DownloadResult).
func (a *AbstractCacheAccessByKey[D, P, V]) UpdateDownloadResult(result DownloadResult) {
	a.downloadCache.Update(a.key, result)
}

// DownloadError sets the download error. Port of downloadError(Exception).
func (a *AbstractCacheAccessByKey[D, P, V]) DownloadError(e error) {
	a.downloadCache.Error(a.key, e)
}

// IsParsingRefreshNeeded gets if the parsing refresh is needed. Port of
// isParsingRefreshNeeded().
func (a *AbstractCacheAccessByKey[D, P, V]) IsParsingRefreshNeeded() bool {
	return a.parsingCache.IsRefreshNeeded(a.key)
}

// UpdateParsingResult updates the parsing result. Port of update(ParsingResult).
func (a *AbstractCacheAccessByKey[D, P, V]) UpdateParsingResult(parsingResult ParsingResult) {
	a.parsingCache.Update(a.key, parsingResult)
}

// ExpireParsing sets the parsing record to the expired state. Port of expireParsing().
func (a *AbstractCacheAccessByKey[D, P, V]) ExpireParsing() {
	a.parsingCache.Expire(a.key)
}

// ParsingError sets the parsing error. Port of parsingError(Exception).
func (a *AbstractCacheAccessByKey[D, P, V]) ParsingError(e error) {
	a.parsingCache.Error(a.key, e)
}

// IsValidationRefreshNeeded gets if the validation refresh is needed. Port of
// isValidationRefreshNeeded().
func (a *AbstractCacheAccessByKey[D, P, V]) IsValidationRefreshNeeded() bool {
	return a.validationCache.IsRefreshNeeded(a.key)
}

// ExpireValidation expires the validation record. Port of expireValidation().
func (a *AbstractCacheAccessByKey[D, P, V]) ExpireValidation() {
	a.validationCache.Expire(a.key)
}

// UpdateValidationResult updates the validation record. Port of update(ValidationResult).
func (a *AbstractCacheAccessByKey[D, P, V]) UpdateValidationResult(validationResult ValidationResult) {
	a.validationCache.Update(a.key, validationResult)
}

// ValidationError sets the validation error. Port of validationError(Exception).
func (a *AbstractCacheAccessByKey[D, P, V]) ValidationError(e error) {
	a.validationCache.Error(a.key, e)
}

// IsFileNeedToBeDeleted checks if the entry must be deleted from the file cache (download
// cache). Port of isFileNeedToBeDeleted().
func (a *AbstractCacheAccessByKey[D, P, V]) IsFileNeedToBeDeleted() bool {
	return a.downloadCache.IsToBeDeleted(a.key)
}

// DeleteDownloadCacheIfNeeded removes the entry from downloadCache if its value is
// TO_BE_DELETED. Port of deleteDownloadCacheIfNeeded().
func (a *AbstractCacheAccessByKey[D, P, V]) DeleteDownloadCacheIfNeeded() {
	if a.downloadCache.IsToBeDeleted(a.key) {
		a.downloadCache.Remove(a.key)
	}
}

// DeleteParsingCacheIfNeeded removes the entry from parsingCache if its value is
// TO_BE_DELETED. Port of deleteParsingCacheIfNeeded().
func (a *AbstractCacheAccessByKey[D, P, V]) DeleteParsingCacheIfNeeded() {
	if a.parsingCache.IsToBeDeleted(a.key) {
		a.parsingCache.Remove(a.key)
	}
}

// DeleteValidationCacheIfNeeded removes the entry from validationCache if its value is
// TO_BE_DELETED. Port of deleteValidationCacheIfNeeded() (Java's doc comment says
// "parsingCache" but the implementation acts on validationCache; ported as implemented).
func (a *AbstractCacheAccessByKey[D, P, V]) DeleteValidationCacheIfNeeded() {
	if a.validationCache.IsToBeDeleted(a.key) {
		a.validationCache.Remove(a.key)
	}
}

var _ CacheAccessByKey = (*AbstractCacheAccessByKey[modeljob.DownloadInfoRecord, modeljob.ParsingInfoRecord, modeljob.ValidationInfoRecord])(nil)
