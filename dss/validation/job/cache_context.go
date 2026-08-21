// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CacheContext.java (DSS 6.5.RC1).
package job

import "time"

// CacheContext contains information for a Cache entry.
type CacheContext interface {
	// CurrentState returns the current state in the cache. Port of getCurrentState().
	CurrentState() CacheStateEnum
	// LastStateTransitionTime returns the date of the last state transition. The zero
	// time.Time stands for Java's null. Port of getLastStateTransitionTime().
	LastStateTransitionTime() time.Time
	// LastSuccessSynchronizationTime returns the last time when the cache has been
	// synchronized successfully. NOTE: can be the zero time.Time in case if the cache has
	// never been synchronized. Port of getLastSuccessSynchronizationTime().
	LastSuccessSynchronizationTime() time.Time
	// State operates a state change. Port of state(CacheState).
	State(newState CacheState)
	// Desync sets the context as DESYNCHRONIZED. Port of desync().
	Desync()
	// Sync sets the context as SYNCHRONIZED. Port of sync().
	Sync()
	// RefreshNeeded sets the context as REFRESH_NEEDED. Port of refreshNeeded().
	RefreshNeeded()
	// ToBeDeleted sets the context as TO_BE_DELETED. Port of toBeDeleted().
	ToBeDeleted()
	// IsRefreshNeeded returns TRUE is a refresh is needed (missing / expired data). Port of
	// isRefreshNeeded().
	IsRefreshNeeded() bool
	// SyncUpdateDate updates the lastSuccessSynchronization date. Port of syncUpdateDate().
	SyncUpdateDate()
	// IsError returns TRUE if the cache is in an error status. Port of isError().
	IsError() bool
	// Error stores the exception for its occurrence time. Port of
	// error(CachedExceptionWrapper).
	Error(exception *CachedExceptionWrapper)
	// ErrorUpdateDate stores the last occurrence of this exception. Port of
	// errorUpdateDate(CachedExceptionWrapper).
	ErrorUpdateDate(updatedException *CachedExceptionWrapper)
	// Exception returns the met exception. Port of getException().
	Exception() *CachedExceptionWrapper
	// IsToBeDeleted returns TRUE if the cache is in TO_BE_DELETED status. Port of
	// isToBeDeleted().
	IsToBeDeleted() bool
	// IsDesync returns TRUE if the entry is desynchronized. Port of isDesync().
	IsDesync() bool
}
