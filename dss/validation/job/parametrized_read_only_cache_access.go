// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/ParametrizedReadOnlyCacheAccess.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// ParametrizedReadOnlyCacheAccess is the parametrized interface to access the cache in
// read-only mode. D, P, V mirror Java's "<D extends DownloadInfoRecord, P extends
// ParsingInfoRecord, V extends ValidationInfoRecord>".
type ParametrizedReadOnlyCacheAccess[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord] interface {
	ReadOnlyCacheAccess

	// GetDownloadInfoRecordTyped returns the download cache DTO result, typed as D.
	//
	// Java's covariant override narrows getDownloadInfoRecord(CacheKey)'s return type from
	// DownloadInfoRecord to D; Go has no covariant return types, so ReadOnlyCacheAccess's
	// GetDownloadInfoRecord (returning the base modeljob.DownloadInfoRecord) is kept as-is
	// to satisfy that interface, and this differently-named method offers the typed D
	// result for callers that need it.
	GetDownloadInfoRecordTyped(key CacheKey) D

	// GetParsingInfoRecordTyped returns the parsing cache DTO result, typed as P. See
	// GetDownloadInfoRecordTyped for the covariance-substitute naming rationale.
	GetParsingInfoRecordTyped(key CacheKey) P

	// GetValidationInfoRecordTyped returns the validation cache DTO result, typed as V. See
	// GetDownloadInfoRecordTyped for the covariance-substitute naming rationale.
	GetValidationInfoRecordTyped(key CacheKey) V
}
