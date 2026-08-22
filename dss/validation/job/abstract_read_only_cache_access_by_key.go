// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/AbstractReadOnlyCacheAccessByKey.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// AbstractReadOnlyCacheAccessByKey prevents reading of other records but the one with the
// defined key. D, P, V mirror Java's "<D extends DownloadInfoRecord, P extends
// ParsingInfoRecord, V extends ValidationInfoRecord>".
//
// DEVIATION: unlike Java's covariant override (getDownloadReadOnlyResult() returns D, not
// the base DownloadInfoRecord), Go has no covariant return types, so the three read methods
// below return the base modeljob interfaces directly, which lets
// AbstractReadOnlyCacheAccessByKey satisfy the plain ReadOnlyCacheAccessByKey interface for
// any D/P/V instantiation. The underlying value is still exactly the typed D/P/V result
// returned by the ParametrizedReadOnlyCacheAccess; callers needing the concrete type type-
// assert, or use readOnlyCacheAccess.GetDownloadInfoRecordTyped(key) directly.
type AbstractReadOnlyCacheAccessByKey[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord] struct {
	// key is the key of the CacheEntry.
	key CacheKey

	// readOnlyCacheAccess reads the cache by the given key.
	readOnlyCacheAccess ParametrizedReadOnlyCacheAccess[D, P, V]
}

// NewAbstractReadOnlyCacheAccessByKey creates an AbstractReadOnlyCacheAccessByKey for the
// given key and read-only cache access. Port of the protected constructor.
func NewAbstractReadOnlyCacheAccessByKey[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord](
	key CacheKey, readOnlyCacheAccess ParametrizedReadOnlyCacheAccess[D, P, V]) AbstractReadOnlyCacheAccessByKey[D, P, V] {
	return AbstractReadOnlyCacheAccessByKey[D, P, V]{key: key, readOnlyCacheAccess: readOnlyCacheAccess}
}

// GetDownloadReadOnlyResult returns the cached read-only download result DTO. Port of
// getDownloadReadOnlyResult(). See the type DEVIATION note.
func (a *AbstractReadOnlyCacheAccessByKey[D, P, V]) GetDownloadReadOnlyResult() modeljob.DownloadInfoRecord {
	return a.readOnlyCacheAccess.GetDownloadInfoRecordTyped(a.key)
}

// GetParsingReadOnlyResult returns the cached read-only parsing result DTO. Port of
// getParsingReadOnlyResult(). See the type DEVIATION note.
func (a *AbstractReadOnlyCacheAccessByKey[D, P, V]) GetParsingReadOnlyResult() modeljob.ParsingInfoRecord {
	return a.readOnlyCacheAccess.GetParsingInfoRecordTyped(a.key)
}

// GetValidationReadOnlyResult returns the cached read-only validation result DTO. Port of
// getValidationReadOnlyResult(). See the type DEVIATION note.
func (a *AbstractReadOnlyCacheAccessByKey[D, P, V]) GetValidationReadOnlyResult() modeljob.ValidationInfoRecord {
	return a.readOnlyCacheAccess.GetValidationInfoRecordTyped(a.key)
}
