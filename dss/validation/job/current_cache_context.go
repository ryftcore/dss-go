// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CurrentCacheContext.java (DSS 6.5.RC1).
package job

import "time"

// CurrentCacheContext contains information for a cache record state.
//
// slf4j trace/debug logging is dropped (no observable behavior).
type CurrentCacheContext struct {
	// state is the current state.
	state CacheState

	// lastStateTransitionTime is the last time when the state of the current cache context
	// has been changed. The zero time.Time stands for Java's null.
	lastStateTransitionTime time.Time

	// lastSuccessSynchronizationTime is the last time when the cache had been synchronized.
	// The zero time.Time stands for Java's null.
	lastSuccessSynchronizationTime time.Time

	// exception is the wrapped exception, when in error state.
	exception *CachedExceptionWrapper
}

// NewCurrentCacheContext creates a CurrentCacheContext initialized to REFRESH_NEEDED. Port
// of the default constructor.
func NewCurrentCacheContext() *CurrentCacheContext {
	c := &CurrentCacheContext{}
	c.State(CacheStateEnum_REFRESH_NEEDED)
	return c
}

// CurrentState returns the current state. Port of getCurrentState().
func (c *CurrentCacheContext) CurrentState() CacheStateEnum {
	return c.state.(CacheStateEnum)
}

// LastStateTransitionTime returns the last state transition time. Port of
// getLastStateTransitionTime().
func (c *CurrentCacheContext) LastStateTransitionTime() time.Time {
	return c.lastStateTransitionTime
}

// LastSuccessSynchronizationTime returns the last successful synchronization time. Port of
// getLastSuccessSynchronizationTime().
func (c *CurrentCacheContext) LastSuccessSynchronizationTime() time.Time {
	return c.lastSuccessSynchronizationTime
}

// State operates a state change: no-op if newState equals the current state, otherwise
// updates the state, records the transition time and clears any stored exception. Port of
// state(CacheState).
func (c *CurrentCacheContext) State(newState CacheState) {
	if c.state == newState {
		return
	}
	c.state = newState
	c.lastStateTransitionTime = time.Now()
	c.exception = nil
}

// SyncUpdateDate updates the lastSuccessSynchronization date. Port of syncUpdateDate().
func (c *CurrentCacheContext) SyncUpdateDate() {
	c.lastSuccessSynchronizationTime = time.Now()
}

// Error marks the context as ERROR and stores the exception, bypassing the CacheState
// transition table. Port of error(CachedExceptionWrapper).
func (c *CurrentCacheContext) Error(cachedException *CachedExceptionWrapper) {
	c.state = CacheStateEnum_ERROR
	c.exception = cachedException
}

// ErrorUpdateDate updates the last occurrence date of the currently stored exception. Port
// of errorUpdateDate(CachedExceptionWrapper).
func (c *CurrentCacheContext) ErrorUpdateDate(updatedException *CachedExceptionWrapper) {
	c.exception.SetLastOccurrenceDate(updatedException.Date())
}

// Desync delegates the DESYNCHRONIZED transition to the current state. Port of desync().
func (c *CurrentCacheContext) Desync() {
	c.state.Desync(c)
}

// Sync delegates the SYNCHRONIZED transition to the current state, then updates the sync
// date. Port of sync().
func (c *CurrentCacheContext) Sync() {
	c.state.Sync(c)
	c.SyncUpdateDate()
}

// RefreshNeeded delegates the REFRESH_NEEDED transition to the current state. Port of
// refreshNeeded().
func (c *CurrentCacheContext) RefreshNeeded() {
	c.state.RefreshNeeded(c)
}

// ToBeDeleted delegates the TO_BE_DELETED transition to the current state. Port of
// toBeDeleted().
func (c *CurrentCacheContext) ToBeDeleted() {
	c.state.ToBeDeleted(c)
}

// IsRefreshNeeded returns true when the current state is REFRESH_NEEDED. Port of
// isRefreshNeeded().
func (c *CurrentCacheContext) IsRefreshNeeded() bool {
	return CacheStateEnum_REFRESH_NEEDED == c.state
}

// IsError returns true when the current state is ERROR. Port of isError().
func (c *CurrentCacheContext) IsError() bool {
	return CacheStateEnum_ERROR == c.state
}

// Exception returns the stored exception, if any. Port of getException().
func (c *CurrentCacheContext) Exception() *CachedExceptionWrapper {
	return c.exception
}

// IsToBeDeleted returns true when the current state is TO_BE_DELETED. Port of
// isToBeDeleted().
func (c *CurrentCacheContext) IsToBeDeleted() bool {
	return CacheStateEnum_TO_BE_DELETED == c.state
}

// IsDesync returns true when the current state is DESYNCHRONIZED. Port of isDesync().
func (c *CurrentCacheContext) IsDesync() bool {
	return CacheStateEnum_DESYNCHRONIZED == c.state
}
