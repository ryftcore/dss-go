// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/CacheAccessByKey.java (DSS 6.5.RC1).
package job

// CacheAccessByKey provides an interface for accessing a cache by key.
//
// Java overloads update(DownloadResult) / update(ParsingResult) / update(ValidationResult);
// Go has no overloading, so each overload gets a distinct name: UpdateDownloadResult,
// UpdateParsingResult, UpdateValidationResult.
type CacheAccessByKey interface {
	ReadOnlyCacheAccessByKey

	// GetCacheKey returns the CacheKey. Port of getCacheKey().
	GetCacheKey() CacheKey

	// IsUpToDate checks if the download result is up to date for the given key. Port of
	// isUpToDate(DownloadResult).
	IsUpToDate(downloadResult DownloadResult) bool

	// UpdateDownloadResult updates the download result. Port of update(DownloadResult).
	UpdateDownloadResult(result DownloadResult)

	// DownloadError sets the download error. Port of downloadError(Exception).
	DownloadError(e error)

	// IsParsingRefreshNeeded gets if the parsing refresh is needed. Port of
	// isParsingRefreshNeeded().
	IsParsingRefreshNeeded() bool

	// UpdateParsingResult updates the parsing result. Port of update(ParsingResult).
	UpdateParsingResult(parsingResult ParsingResult)

	// ExpireParsing sets the parsing record to the expired state. Port of expireParsing().
	ExpireParsing()

	// ParsingError sets the parsing error. Port of parsingError(Exception).
	ParsingError(e error)

	// IsValidationRefreshNeeded gets if the validation refresh is needed. Port of
	// isValidationRefreshNeeded().
	IsValidationRefreshNeeded() bool

	// ExpireValidation expires the validation record. Port of expireValidation().
	ExpireValidation()

	// UpdateValidationResult updates the validation record. Port of update(ValidationResult).
	UpdateValidationResult(validationResult ValidationResult)

	// ValidationError sets the validation error. Port of validationError(Exception).
	ValidationError(e error)

	// IsFileNeedToBeDeleted checks if the entry must be deleted from the file cache
	// (download cache). Port of isFileNeedToBeDeleted().
	IsFileNeedToBeDeleted() bool

	// DeleteDownloadCacheIfNeeded removes the entry from downloadCache if its value is
	// TO_BE_DELETED. Port of deleteDownloadCacheIfNeeded().
	DeleteDownloadCacheIfNeeded()

	// DeleteParsingCacheIfNeeded removes the entry from parsingCache if its value is
	// TO_BE_DELETED. Port of deleteParsingCacheIfNeeded().
	DeleteParsingCacheIfNeeded()

	// DeleteValidationCacheIfNeeded removes the entry from validationCache if its value is
	// TO_BE_DELETED. Port of deleteValidationCacheIfNeeded().
	DeleteValidationCacheIfNeeded()
}
