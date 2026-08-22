// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CacheStateEnum.java (DSS 6.5.RC1).
package job

import "fmt"

// CacheStateEnum contains the states for a cache record. Java enum constant values are the
// name() strings, per the enum porting convention.
type CacheStateEnum string

const (
	// CacheStateEnumRefreshNeeded: nothing / expired content is stored in the cache.
	CacheStateEnumRefreshNeeded CacheStateEnum = "REFRESH_NEEDED"
	// CacheStateEnumDesynchronized: the cache content is not synchronized with the
	// application.
	CacheStateEnumDesynchronized CacheStateEnum = "DESYNCHRONIZED"
	// CacheStateEnumSynchronized: the application and the cache content are synchronized.
	CacheStateEnumSynchronized CacheStateEnum = "SYNCHRONIZED"
	// CacheStateEnumError: the data cannot be downloaded / parsed / validated.
	CacheStateEnumError CacheStateEnum = "ERROR"
	// CacheStateEnumToBEDeleted: the cache content needs to be deleted.
	// NOTE: URL may become available again if not cleaned!
	CacheStateEnumToBEDeleted CacheStateEnum = "TO_BE_DELETED"
)

// notAllowedTransition mirrors NOT_ALLOWED_TRANSITION.
const notAllowedTransition = "Transition from '%s' to '%s' is not allowed"

// Sync ports the per-constant sync(CacheContext) overrides (DESYNCHRONIZED, SYNCHRONIZED)
// falling back to the interface default (panic) for the other constants.
func (s CacheStateEnum) Sync(cacheContext CacheContext) {
	switch s {
	case CacheStateEnumDesynchronized, CacheStateEnumSynchronized:
		cacheContext.State(CacheStateEnumSynchronized)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnumSynchronized))
	}
}

// Desync ports the per-constant desync(CacheContext) overrides (REFRESH_NEEDED,
// SYNCHRONIZED, ERROR, TO_BE_DELETED) falling back to the interface default (panic) for
// DESYNCHRONIZED itself.
func (s CacheStateEnum) Desync(cacheContext CacheContext) {
	switch s {
	case CacheStateEnumRefreshNeeded, CacheStateEnumSynchronized, CacheStateEnumError, CacheStateEnumToBEDeleted:
		cacheContext.State(CacheStateEnumDesynchronized)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnumDesynchronized))
	}
}

// RefreshNeeded ports the per-constant refreshNeeded(CacheContext) overrides
// (REFRESH_NEEDED, SYNCHRONIZED, ERROR, TO_BE_DELETED) falling back to the interface default
// (panic) for DESYNCHRONIZED.
func (s CacheStateEnum) RefreshNeeded(cacheContext CacheContext) {
	switch s {
	case CacheStateEnumRefreshNeeded, CacheStateEnumSynchronized, CacheStateEnumError, CacheStateEnumToBEDeleted:
		cacheContext.State(CacheStateEnumRefreshNeeded)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnumRefreshNeeded))
	}
}

// ToBeDeleted ports the per-constant toBeDeleted(CacheContext) overrides (REFRESH_NEEDED,
// SYNCHRONIZED, ERROR) falling back to the interface default (panic) for DESYNCHRONIZED and
// TO_BE_DELETED (TO_BE_DELETED does not override toBeDeleted() in Java).
func (s CacheStateEnum) ToBeDeleted(cacheContext CacheContext) {
	switch s {
	case CacheStateEnumRefreshNeeded, CacheStateEnumSynchronized, CacheStateEnumError:
		cacheContext.State(CacheStateEnumToBEDeleted)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnumToBEDeleted))
	}
}

// Error ports the per-constant error(CacheContext, CachedExceptionWrapper) override
// (REFRESH_NEEDED only) falling back to the interface default (panic "Cannot store error")
// for every other constant, including ERROR itself.
func (s CacheStateEnum) Error(cacheContext CacheContext, exception *CachedExceptionWrapper) {
	switch s {
	case CacheStateEnumRefreshNeeded:
		cacheContext.Error(exception)
	default:
		panic("Cannot store error")
	}
}
