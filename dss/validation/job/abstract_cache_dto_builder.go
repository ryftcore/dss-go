// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/builder/AbstractCacheDTOBuilder.java (DSS 6.5.RC1).
package job

// AbstractCacheDTOBuilder is an abstract builder of a Cache DTO. R is the type of the cache
// result, mirroring Java's "AbstractCacheDTOBuilder<R extends CachedResult>".
//
// slf4j trace/debug logging is dropped (no observable behavior).
type AbstractCacheDTOBuilder[R CachedResult] struct {
	// cachedEntry is the cached entry.
	cachedEntry *CachedEntry[R]
}

// NewAbstractCacheDTOBuilder creates an AbstractCacheDTOBuilder for cachedEntry. Port of the
// protected constructor.
func NewAbstractCacheDTOBuilder[R CachedResult](cachedEntry *CachedEntry[R]) AbstractCacheDTOBuilder[R] {
	return AbstractCacheDTOBuilder[R]{cachedEntry: cachedEntry}
}

// Build builds the DTO. Port of build().
func (b *AbstractCacheDTOBuilder[R]) Build() *AbstractCacheDTO {
	abstractCacheDTO := NewAbstractCacheDTO()
	abstractCacheDTO.SetCacheState(b.cachedEntry.CurrentState())
	abstractCacheDTO.SetLastStateTransitionTime(b.cachedEntry.LastStateTransitionTime())
	abstractCacheDTO.SetLastSuccessSynchronizationTime(b.cachedEntry.LastSuccessSynchronizationTime())
	abstractCacheDTO.SetExceptionMessage(b.cachedEntry.ExceptionMessage())
	abstractCacheDTO.SetExceptionStackTrace(b.cachedEntry.ExceptionStackTrace())
	abstractCacheDTO.SetExceptionFirstOccurrenceTime(b.cachedEntry.ExceptionFirstOccurrenceTime())
	abstractCacheDTO.SetExceptionLastOccurrenceTime(b.cachedEntry.ExceptionLastOccurrenceTime())
	abstractCacheDTO.SetResultExist(b.IsResultExist())
	return abstractCacheDTO
}

// Result returns the cached result. Port of the protected final getResult().
func (b *AbstractCacheDTOBuilder[R]) Result() R {
	return b.cachedEntry.CachedResult()
}

// IsResultExist returns true if the result exists. Port of isResultExist().
func (b *AbstractCacheDTOBuilder[R]) IsResultExist() bool {
	return !isNilCachedResult(b.Result())
}
