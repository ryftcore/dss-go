// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CachedEntry.java (DSS 6.5.RC1).
package job

import (
	"reflect"
	"time"

	"github.com/utain/esig/dss/utils"
)

// CachedEntry defines a cached entry. R is the type of the cached result, mirroring Java's
// "CachedEntry<R extends CachedResult>".
//
// slf4j trace/debug logging is dropped (no observable behavior).
type CachedEntry[R CachedResult] struct {
	// cacheContext contains information about the cache entry.
	cacheContext CacheContext

	// cachedResult is the cached result. The zero value of R stands for Java's null.
	cachedResult R
}

// NewCachedEntry creates an empty CachedEntry. Port of the empty constructor.
func NewCachedEntry[R CachedResult]() *CachedEntry[R] {
	return &CachedEntry[R]{cacheContext: NewCurrentCacheContext()}
}

// NewCachedEntryWithResult creates a CachedEntry initialized with cachedResult. Port of
// CachedEntry(R).
func NewCachedEntryWithResult[R CachedResult](cachedResult R) *CachedEntry[R] {
	e := NewCachedEntry[R]()
	e.Update(cachedResult)
	return e
}

// CurrentState returns the state of the cache. Port of getCurrentState().
func (e *CachedEntry[R]) CurrentState() CacheStateEnum {
	return e.cacheContext.CurrentState()
}

// LastStateTransitionTime returns the last status change time. Port of
// getLastStateTransitionTime().
func (e *CachedEntry[R]) LastStateTransitionTime() time.Time {
	return e.cacheContext.LastStateTransitionTime()
}

// LastSuccessSynchronizationTime returns the last synchronization time. Port of
// getLastSuccessSynchronizationTime().
func (e *CachedEntry[R]) LastSuccessSynchronizationTime() time.Time {
	return e.cacheContext.LastSuccessSynchronizationTime()
}

// CachedResult returns the cached result. Port of getCachedResult().
func (e *CachedEntry[R]) CachedResult() R {
	return e.cachedResult
}

// Update updates the cache record. Panics if newCachedResult is the zero value of R,
// mirroring Java's Objects.requireNonNull("Cached result cannot be overwritten with a null
// value"). Port of update(R).
func (e *CachedEntry[R]) Update(newCachedResult R) {
	if isNilCachedResult(newCachedResult) {
		panic("Cached result cannot be overwritten with a null value")
	}
	e.cacheContext.Desync() // if transition is not allowed, cached object is not updated
	e.cachedResult = newCachedResult
}

// isNilCachedResult reports whether a CachedResult value is Java-null, i.e. a nil interface
// or a nil value stored in a non-nil interface (covers R being a pointer or interface type).
func isNilCachedResult[R CachedResult](v R) bool {
	var iface CachedResult = v
	if iface == nil {
		return true
	}
	rv := reflect.ValueOf(iface)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}

// SyncUpdateDate synchronizes the update date. Port of syncUpdateDate().
func (e *CachedEntry[R]) SyncUpdateDate() {
	e.cacheContext.SyncUpdateDate()
}

// Error sets the error. Port of error(CachedExceptionWrapper).
func (e *CachedEntry[R]) Error(exception *CachedExceptionWrapper) {
	if e.isNewError(exception) {
		e.cacheContext.Error(exception)
		var zero R
		e.cachedResult = zero // reset in case of error
	} else {
		e.cacheContext.ErrorUpdateDate(exception)
	}
}

// isNewError checks if the wrappedException is a new one. Port of
// isNewError(CachedExceptionWrapper).
func (e *CachedEntry[R]) isNewError(wrappedException *CachedExceptionWrapper) bool {
	return !e.IsError() || !utils.AreStringsEqual(e.ExceptionStackTrace(), wrappedException.StackTrace())
}

// Expire expires the cache entry. Port of expire().
func (e *CachedEntry[R]) Expire() {
	e.cacheContext.RefreshNeeded()
}

// Sync synchronizes the cache entry. Port of sync().
func (e *CachedEntry[R]) Sync() {
	e.cacheContext.Sync()
}

// ToBeDeleted sets 'toBeDeleted' status for the cache entry. Port of toBeDeleted().
func (e *CachedEntry[R]) ToBeDeleted() {
	e.cacheContext.ToBeDeleted()
}

// IsToBeDeleted checks if the status 'toBeDeleted' is set for the cache entry. Port of
// isToBeDeleted().
func (e *CachedEntry[R]) IsToBeDeleted() bool {
	return e.cacheContext.IsToBeDeleted()
}

// IsDesync checks if the status 'desynchronized' is set for the cache entry. Port of
// isDesync().
func (e *CachedEntry[R]) IsDesync() bool {
	return e.cacheContext.IsDesync()
}

// IsRefreshNeeded checks if the refresh is needed for the cache entry. Port of
// isRefreshNeeded().
func (e *CachedEntry[R]) IsRefreshNeeded() bool {
	return e.cacheContext.IsRefreshNeeded()
}

// IsEmpty checks if the cache record is empty. Port of isEmpty().
func (e *CachedEntry[R]) IsEmpty() bool {
	return isNilCachedResult(e.cachedResult)
}

// IsResultExist gets if a result exists under the record. Port of isResultExist().
func (e *CachedEntry[R]) IsResultExist() bool {
	return !e.IsEmpty()
}

// IsError checks if the current status of the cache is error. Port of isError().
func (e *CachedEntry[R]) IsError() bool {
	return e.cacheContext.IsError()
}

// ExceptionMessage gets the exception message for an error status. Port of
// getExceptionMessage().
func (e *CachedEntry[R]) ExceptionMessage() string {
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().ExceptionMessage()
	}
	return ""
}

// ExceptionStackTrace gets the exception stack trace for an error status. Port of
// getExceptionStackTrace().
func (e *CachedEntry[R]) ExceptionStackTrace() string {
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().StackTrace()
	}
	return ""
}

// ExceptionFirstOccurrenceTime gets the first time when the exception occurred. The zero
// time.Time stands for Java's null. Port of getExceptionFirstOccurrenceTime().
func (e *CachedEntry[R]) ExceptionFirstOccurrenceTime() time.Time {
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().Date()
	}
	return time.Time{}
}

// ExceptionLastOccurrenceTime gets the last time when the exception occurred. The zero
// time.Time stands for Java's null. Port of getExceptionLastOccurrenceTime().
func (e *CachedEntry[R]) ExceptionLastOccurrenceTime() time.Time {
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().LastOccurrenceDate()
	}
	return time.Time{}
}
