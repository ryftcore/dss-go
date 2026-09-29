// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/AbstractCache.java (DSS 6.5.RC1).
package job

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractCacheOverrides declares the operation Java's abstract AbstractCache<R> class
// leaves abstract, standing in for the virtual dispatch the base needs to reach the concrete
// cache's type. A concrete cache registers itself with AbstractCache.InitAbstractCache.
type AbstractCacheOverrides interface {
	// CacheType returns the type of the current Cache. Port of the abstract protected
	// getCacheType().
	CacheType() CacheType
}

// AbstractCache contains basic methods for handling the CachedResult implementations. R is
// the type of the cached result, mirroring Java's "AbstractCache<R extends CachedResult>".
//
// Java's Map<CacheKey, CachedEntry<R>> is a ConcurrentHashMap; the Go port uses a plain map
// guarded by a mutex, since get-or-create must be atomic (Java's ConcurrentHashMap has no
// direct single-call equivalent for "get, else insert-and-return" without additional
// synchronization either — the Java code itself is not atomic across the two map accesses in
// get(); the mutex here is at least as safe).
//
// slf4j trace/debug/info/warn logging is dropped (no observable behavior).
type AbstractCache[R CachedResult] struct {
	// overrides points back at the concrete cache; see InitAbstractCache.
	overrides AbstractCacheOverrides

	mu sync.Mutex

	// cachedEntriesMap maps a CacheKey to the related result wrapper CachedEntry<R>.
	cachedEntriesMap map[CacheKey]*CachedEntry[R]
}

// NewAbstractCache creates an AbstractCache with an empty map. Port of the protected default
// constructor.
func NewAbstractCache[R CachedResult]() AbstractCache[R] {
	return AbstractCache[R]{cachedEntriesMap: make(map[CacheKey]*CachedEntry[R])}
}

// InitAbstractCache registers the concrete cache with its base so that the base can dispatch
// to CacheType. It must be called exactly once, by the concrete cache's constructor, before
// any other method.
func (a *AbstractCache[R]) InitAbstractCache(overrides AbstractCacheOverrides) {
	a.overrides = overrides
}

// abstractCacheOverrides returns the registered overrides, panicking when the concrete cache
// forgot to call InitAbstractCache.
func (a *AbstractCache[R]) abstractCacheOverrides() AbstractCacheOverrides {
	if a.overrides == nil {
		panic("AbstractCache was not initialised: the concrete cache must call InitAbstractCache in its constructor")
	}
	return a.overrides
}

// Keys returns all current keys. Port of getKeys().
func (a *AbstractCache[R]) Keys() []CacheKey {
	a.mu.Lock()
	defer a.mu.Unlock()
	keys := make([]CacheKey, 0, len(a.cachedEntriesMap))
	for k := range a.cachedEntriesMap {
		keys = append(keys, k)
	}
	return keys
}

// Get returns the CachedEntry for the related cacheKey. Returns a new empty entry if no
// result is found for the key. Port of get(CacheKey).
//
// mu only guards the map: the returned *CachedEntry is shared by every caller of the key and
// is safe for concurrent use because it serializes its own state (see CachedEntry).
func (a *AbstractCache[R]) Get(cacheKey CacheKey) *CachedEntry[R] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if cacheWrapper, ok := a.cachedEntriesMap[cacheKey]; ok {
		return cacheWrapper
	}
	emptyEntry := NewCachedEntry[R]()
	a.cachedEntriesMap[cacheKey] = emptyEntry
	return emptyEntry
}

// Update updates in the cache the value for cacheKey with the given result. Port of
// update(CacheKey, R).
func (a *AbstractCache[R]) Update(cacheKey CacheKey, result R) {
	a.Get(cacheKey).Update(result)
}

// Expire updates the state for a CachedEntry matching the given key to EXPIRED. Port of
// expire(CacheKey).
func (a *AbstractCache[R]) Expire(cacheKey CacheKey) {
	a.Get(cacheKey).Expire()
}

// Remove removes the requested entry with the given cacheKey. Port of remove(CacheKey).
func (a *AbstractCache[R]) Remove(cacheKey CacheKey) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.cachedEntriesMap, cacheKey)
}

// Sync updates the state for a CachedEntry matching the given key to SYNCHRONIZED. Port of
// sync(CacheKey).
func (a *AbstractCache[R]) Sync(cacheKey CacheKey) {
	a.Get(cacheKey).Sync()
}

// IsRefreshNeeded checks if a CachedEntry for the given key is not up to date. Port of
// isRefreshNeeded(CacheKey).
func (a *AbstractCache[R]) IsRefreshNeeded(cacheKey CacheKey) bool {
	return a.Get(cacheKey).IsRefreshNeeded()
}

// IsDesync checks if a CachedEntry for the given key is desynchronized. Port of
// isDesync(CacheKey).
func (a *AbstractCache[R]) IsDesync(cacheKey CacheKey) bool {
	return a.Get(cacheKey).IsDesync()
}

// IsEmpty checks if a CachedEntry for the given key is empty (has no result). Port of
// isEmpty(CacheKey).
func (a *AbstractCache[R]) IsEmpty(cacheKey CacheKey) bool {
	return a.Get(cacheKey).IsEmpty()
}

// Error updates entry status to ERROR value. Port of error(CacheKey, Exception).
func (a *AbstractCache[R]) Error(cacheKey CacheKey, e error) {
	wrappedException := NewCachedExceptionWrapper(e)
	a.Get(cacheKey).Error(wrappedException)
}

// ToBeDeleted updates entry status to TO_BE_DELETED value. Port of toBeDeleted(CacheKey).
func (a *AbstractCache[R]) ToBeDeleted(cacheKey CacheKey) {
	a.Get(cacheKey).ToBeDeleted()
}

// IsToBeDeleted checks if the requested cacheKey has TO_BE_DELETED value. Port of
// isToBeDeleted(CacheKey).
func (a *AbstractCache[R]) IsToBeDeleted(cacheKey CacheKey) bool {
	return a.Get(cacheKey).IsToBeDeleted()
}

// Dump produces a report of the current cache state. Port of dump().
//
// DEVIATION: Java iterates the ConcurrentHashMap's entrySet() in its own unspecified hash
// order; the Go port sorts entries by cache key so the dump is deterministic (this method
// only produces human-readable diagnostic text, never a value compared byte-for-byte against
// Java, but determinism is kept for reproducible debug output).
func (a *AbstractCache[R]) Dump() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("Cache ")
	sb.WriteString(string(a.abstractCacheOverrides().CacheType()))
	if utils.IsMapEmpty(a.cachedEntriesMap) {
		sb.WriteString(" : EMPTY")
	} else {
		sb.WriteString(fmt.Sprintf(" : (nb entries : %d)\n", len(a.cachedEntriesMap)))

		keys := make([]CacheKey, 0, len(a.cachedEntriesMap))
		for k := range a.cachedEntriesMap {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i].Key() < keys[j].Key() })

		for _, key := range keys {
			value := a.cachedEntriesMap[key]
			currentKey := key.Key()
			currentState := value.CurrentState()
			date := "?"
			lastStateTransitionTime := value.LastStateTransitionTime()
			if !lastStateTransitionTime.IsZero() {
				date = lastStateTransitionTime.Format("2006/01/02 15:04:05")
			}
			sb.WriteString(fmt.Sprintf("%-70.70s -> %-25.25s @ %.20s", currentKey, currentState, date))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
