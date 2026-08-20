// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CacheStateEnum.java (DSS 6.5.RC1).
package job

import "fmt"

// CacheStateEnum contains the states for a cache record. Java enum constant values are the
// name() strings, per the enum porting convention.
type CacheStateEnum string

const (
	// CacheStateEnum_REFRESH_NEEDED: nothing / expired content is stored in the cache.
	CacheStateEnum_REFRESH_NEEDED CacheStateEnum = "REFRESH_NEEDED"
	// CacheStateEnum_DESYNCHRONIZED: the cache content is not synchronized with the
	// application.
	CacheStateEnum_DESYNCHRONIZED CacheStateEnum = "DESYNCHRONIZED"
	// CacheStateEnum_SYNCHRONIZED: the application and the cache content are synchronized.
	CacheStateEnum_SYNCHRONIZED CacheStateEnum = "SYNCHRONIZED"
	// CacheStateEnum_ERROR: the data cannot be downloaded / parsed / validated.
	CacheStateEnum_ERROR CacheStateEnum = "ERROR"
	// CacheStateEnum_TO_BE_DELETED: the cache content needs to be deleted.
	// NOTE: URL may become available again if not cleaned!
	CacheStateEnum_TO_BE_DELETED CacheStateEnum = "TO_BE_DELETED"
)

// notAllowedTransition mirrors NOT_ALLOWED_TRANSITION.
const notAllowedTransition = "Transition from '%s' to '%s' is not allowed"

// Sync ports the per-constant sync(CacheContext) overrides (DESYNCHRONIZED, SYNCHRONIZED)
// falling back to the interface default (panic) for the other constants.
func (s CacheStateEnum) Sync(cacheContext CacheContext) {
	switch s {
	case CacheStateEnum_DESYNCHRONIZED, CacheStateEnum_SYNCHRONIZED:
		cacheContext.State(CacheStateEnum_SYNCHRONIZED)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnum_SYNCHRONIZED))
	}
}

// Desync ports the per-constant desync(CacheContext) overrides (REFRESH_NEEDED,
// SYNCHRONIZED, ERROR, TO_BE_DELETED) falling back to the interface default (panic) for
// DESYNCHRONIZED itself.
func (s CacheStateEnum) Desync(cacheContext CacheContext) {
	switch s {
	case CacheStateEnum_REFRESH_NEEDED, CacheStateEnum_SYNCHRONIZED, CacheStateEnum_ERROR, CacheStateEnum_TO_BE_DELETED:
		cacheContext.State(CacheStateEnum_DESYNCHRONIZED)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnum_DESYNCHRONIZED))
	}
}

// RefreshNeeded ports the per-constant refreshNeeded(CacheContext) overrides
// (REFRESH_NEEDED, SYNCHRONIZED, ERROR, TO_BE_DELETED) falling back to the interface default
// (panic) for DESYNCHRONIZED.
func (s CacheStateEnum) RefreshNeeded(cacheContext CacheContext) {
	switch s {
	case CacheStateEnum_REFRESH_NEEDED, CacheStateEnum_SYNCHRONIZED, CacheStateEnum_ERROR, CacheStateEnum_TO_BE_DELETED:
		cacheContext.State(CacheStateEnum_REFRESH_NEEDED)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnum_REFRESH_NEEDED))
	}
}

// ToBeDeleted ports the per-constant toBeDeleted(CacheContext) overrides (REFRESH_NEEDED,
// SYNCHRONIZED, ERROR) falling back to the interface default (panic) for DESYNCHRONIZED and
// TO_BE_DELETED (TO_BE_DELETED does not override toBeDeleted() in Java).
func (s CacheStateEnum) ToBeDeleted(cacheContext CacheContext) {
	switch s {
	case CacheStateEnum_REFRESH_NEEDED, CacheStateEnum_SYNCHRONIZED, CacheStateEnum_ERROR:
		cacheContext.State(CacheStateEnum_TO_BE_DELETED)
	default:
		panic(fmt.Sprintf(notAllowedTransition, cacheContext.CurrentState(), CacheStateEnum_TO_BE_DELETED))
	}
}

// Error ports the per-constant error(CacheContext, CachedExceptionWrapper) override
// (REFRESH_NEEDED only) falling back to the interface default (panic "Cannot store error")
// for every other constant, including ERROR itself.
func (s CacheStateEnum) Error(cacheContext CacheContext, exception *CachedExceptionWrapper) {
	switch s {
	case CacheStateEnum_REFRESH_NEEDED:
		cacheContext.Error(exception)
	default:
		panic("Cannot store error")
	}
}
