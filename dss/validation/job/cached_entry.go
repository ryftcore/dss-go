// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CachedEntry.java (DSS 6.5.RC1).
package job

import (
	"reflect"
	"sync"
	"time"

	"github.com/ryftcore/dss-go/dss/utils"
)

// CachedEntry defines a cached entry. R is the type of the cached result, mirroring Java's
// "CachedEntry<R extends CachedResult>".
//
// slf4j trace/debug logging is dropped (no observable behavior).
//
// DEVIATION (concurrency): a CachedEntry is shared between goroutines - AbstractCache.Get hands
// the same *CachedEntry to every caller of a key, and the TL/LOTL job drives one entry from
// several analysis goroutines at once (one goroutine per LOTL/TL analysis, plus one per pivot in
// tsl.LOTLWithPivotsAnalysis, which all expire the shared LOTL and preceding-pivot validation
// entries). Upstream's CachedEntry/CurrentCacheContext carry no synchronization of their own
// (plain fields, no volatile/synchronized), but a racing plain-field write on the JVM is benign,
// whereas a Go data race on the state or result interface value can tear (undefined behaviour).
// mu therefore serializes every access to cacheContext and cachedResult. Each exported method
// takes it once for its whole body, so a compound transition such as Update's "desync, then
// store the result" or Error's "test whether the error is new, then store it" is atomic, and
// delegates to the unexported ...Locked helpers, never to another locking method (sync.Mutex is
// not reentrant). Nothing called while mu is held reaches back into a cache: the CacheState
// transition table only calls back into this entry's own CurrentCacheContext, so holding mu
// across those callbacks cannot deadlock (and a transition panic unwinds through the deferred
// Unlock). The lock order in this package is AbstractCache.mu -> CachedEntry.mu (Dump reads the
// entries while holding its map mutex); no CachedEntry method ever calls into an AbstractCache.
// Sequential behaviour, including every transition panic, is identical to upstream's.
type CachedEntry[R CachedResult] struct {
	// mu guards cacheContext and cachedResult.
	mu sync.Mutex

	// cacheContext contains information about the cache entry. It is only ever reached through
	// this entry with mu held: CurrentCacheContext is not safe for concurrent use by itself.
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
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.CurrentState()
}

// LastStateTransitionTime returns the last status change time. Port of
// getLastStateTransitionTime().
func (e *CachedEntry[R]) LastStateTransitionTime() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.LastStateTransitionTime()
}

// LastSuccessSynchronizationTime returns the last synchronization time. Port of
// getLastSuccessSynchronizationTime().
func (e *CachedEntry[R]) LastSuccessSynchronizationTime() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.LastSuccessSynchronizationTime()
}

// CachedResult returns the cached result. Port of getCachedResult().
func (e *CachedEntry[R]) CachedResult() R {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cachedResult
}

// Update updates the cache record. Panics if newCachedResult is the zero value of R,
// mirroring Java's Objects.requireNonNull("Cached result cannot be overwritten with a null
// value"). Port of update(R).
func (e *CachedEntry[R]) Update(newCachedResult R) {
	if isNilCachedResult(newCachedResult) {
		panic("Cached result cannot be overwritten with a null value")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
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
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cacheContext.SyncUpdateDate()
}

// Error sets the error. Port of error(CachedExceptionWrapper).
func (e *CachedEntry[R]) Error(exception *CachedExceptionWrapper) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.isNewErrorLocked(exception) {
		e.cacheContext.Error(exception)
		var zero R
		e.cachedResult = zero // reset in case of error
	} else {
		e.cacheContext.ErrorUpdateDate(exception)
	}
}

// isNewErrorLocked checks if the wrappedException is a new one. Port of
// isNewError(CachedExceptionWrapper). The caller holds e.mu.
func (e *CachedEntry[R]) isNewErrorLocked(wrappedException *CachedExceptionWrapper) bool {
	return !e.cacheContext.IsError() ||
		!utils.AreStringsEqual(e.exceptionStackTraceLocked(), wrappedException.StackTrace())
}

// Expire expires the cache entry. Port of expire().
func (e *CachedEntry[R]) Expire() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cacheContext.RefreshNeeded()
}

// Sync synchronizes the cache entry. Port of sync().
func (e *CachedEntry[R]) Sync() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cacheContext.Sync()
}

// ToBeDeleted sets 'toBeDeleted' status for the cache entry. Port of toBeDeleted().
func (e *CachedEntry[R]) ToBeDeleted() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cacheContext.ToBeDeleted()
}

// IsToBeDeleted checks if the status 'toBeDeleted' is set for the cache entry. Port of
// isToBeDeleted().
func (e *CachedEntry[R]) IsToBeDeleted() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.IsToBeDeleted()
}

// IsDesync checks if the status 'desynchronized' is set for the cache entry. Port of
// isDesync().
func (e *CachedEntry[R]) IsDesync() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.IsDesync()
}

// IsRefreshNeeded checks if the refresh is needed for the cache entry. Port of
// isRefreshNeeded().
func (e *CachedEntry[R]) IsRefreshNeeded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.IsRefreshNeeded()
}

// IsEmpty checks if the cache record is empty. Port of isEmpty().
func (e *CachedEntry[R]) IsEmpty() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return isNilCachedResult(e.cachedResult)
}

// IsResultExist gets if a result exists under the record. Port of isResultExist().
func (e *CachedEntry[R]) IsResultExist() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return !isNilCachedResult(e.cachedResult)
}

// IsError checks if the current status of the cache is error. Port of isError().
func (e *CachedEntry[R]) IsError() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cacheContext.IsError()
}

// ExceptionMessage gets the exception message for an error status. Port of
// getExceptionMessage().
func (e *CachedEntry[R]) ExceptionMessage() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().ExceptionMessage()
	}
	return ""
}

// ExceptionStackTrace gets the exception stack trace for an error status. Port of
// getExceptionStackTrace().
func (e *CachedEntry[R]) ExceptionStackTrace() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.exceptionStackTraceLocked()
}

// exceptionStackTraceLocked is ExceptionStackTrace for a caller that already holds e.mu.
func (e *CachedEntry[R]) exceptionStackTraceLocked() string {
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().StackTrace()
	}
	return ""
}

// ExceptionFirstOccurrenceTime gets the first time when the exception occurred. The zero
// time.Time stands for Java's null. Port of getExceptionFirstOccurrenceTime().
func (e *CachedEntry[R]) ExceptionFirstOccurrenceTime() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().Date()
	}
	return time.Time{}
}

// ExceptionLastOccurrenceTime gets the last time when the exception occurred. The zero
// time.Time stands for Java's null. Port of getExceptionLastOccurrenceTime().
func (e *CachedEntry[R]) ExceptionLastOccurrenceTime() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cacheContext.Exception() != nil {
		return e.cacheContext.Exception().LastOccurrenceDate()
	}
	return time.Time{}
}
