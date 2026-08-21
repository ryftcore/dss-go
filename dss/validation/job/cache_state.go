// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CacheState.java (DSS 6.5.RC1).
package job

// CacheState defines the different possible transitions from a CacheContext to another one.
//
// Java models this as an enum (CacheStateEnum) implementing this interface, with each
// constant overriding a subset of the methods to define which transitions are legal from
// that state; methods left un-overridden fall back to the interface default, which panics.
// See CacheStateEnum for the Go port of the per-constant transition table.
type CacheState interface {
	// Sync marks the cache entry as Synchronized. Port of sync(CacheContext).
	Sync(cacheContext CacheContext)
	// Desync marks the cache entry as Desynchronized. Port of desync(CacheContext).
	Desync(cacheContext CacheContext)
	// RefreshNeeded marks the cache entry as needing a refresh. Port of
	// refreshNeeded(CacheContext).
	RefreshNeeded(cacheContext CacheContext)
	// ToBeDeleted marks the cache entry as to be deleted. Port of toBeDeleted(CacheContext).
	ToBeDeleted(cacheContext CacheContext)
	// Error marks the cache entry in error state with the related exception. Port of
	// error(CacheContext, CachedExceptionWrapper).
	Error(cacheContext CacheContext, exception *CachedExceptionWrapper)
}
