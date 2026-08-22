// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/CacheType.java (DSS 6.5.RC1).
package job

// CacheType defines a list of possible Cache Types.
type CacheType string

const (
	// CacheTypeDownload is the download task cache.
	CacheTypeDownload CacheType = "DOWNLOAD"
	// CacheTypeParsing is the parsing task cache.
	CacheTypeParsing CacheType = "PARSING"
	// CacheTypeValidation is the validation task cache.
	CacheTypeValidation CacheType = "VALIDATION"
)
