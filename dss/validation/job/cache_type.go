// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/CacheType.java (DSS 6.5.RC1).
package job

// CacheType defines a list of possible Cache Types.
type CacheType string

const (
	// CacheType_DOWNLOAD is the download task cache.
	CacheType_DOWNLOAD CacheType = "DOWNLOAD"
	// CacheType_PARSING is the parsing task cache.
	CacheType_PARSING CacheType = "PARSING"
	// CacheType_VALIDATION is the validation task cache.
	CacheType_VALIDATION CacheType = "VALIDATION"
)
