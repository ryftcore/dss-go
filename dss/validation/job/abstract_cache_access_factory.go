// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/AbstractCacheAccessFactory.java (DSS 6.5.RC1).
package job

import modeljob "github.com/utain/esig/dss/model/job"

// AbstractCacheAccessFactoryOverrides declares the operations Java's abstract
// AbstractCacheAccessFactory<D,P,V> class leaves abstract, standing in for the virtual
// dispatch the base needs to reach the concrete factory. A concrete factory registers itself
// with AbstractCacheAccessFactory.InitAbstractCacheAccessFactory.
type AbstractCacheAccessFactoryOverrides[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord] interface {
	// GetCacheAccess loads a class to deal with a cache by the key records. Port of the
	// abstract getCacheAccess(CacheKey) (declared to return AbstractCacheAccessByKey<D,P,V>
	// in Java; kept as the plain CacheAccessByKey here since every concrete
	// AbstractCacheAccessByKey[D,P,V] already satisfies it).
	GetCacheAccess(key CacheKey) CacheAccessByKey

	// GetReadOnlyCacheAccessTyped loads the typed read-only cache access. Port of the
	// abstract getReadOnlyCacheAccess() (renamed to avoid a name clash with
	// AbstractCacheAccessFactory's own base-interface-returning GetReadOnlyCacheAccess, see
	// its DEVIATION note).
	GetReadOnlyCacheAccessTyped() ParametrizedReadOnlyCacheAccess[D, P, V]
}

// AbstractCacheAccessFactory builds the classes to deal with the cache. D, P, V mirror
// Java's "<D extends DownloadInfoRecord, P extends ParsingInfoRecord, V extends
// ValidationInfoRecord>". It implements CacheAccessFactory.
//
// DEVIATION: Java's getReadOnlyCacheAccess() covariantly overrides CacheAccessFactory's
// return type to ParametrizedReadOnlyCacheAccess<D,P,V>; Go has no covariant return types,
// so GetReadOnlyCacheAccess here keeps the base ReadOnlyCacheAccess return type (the typed
// value returned by GetReadOnlyCacheAccessTyped satisfies it structurally, since
// ParametrizedReadOnlyCacheAccess embeds ReadOnlyCacheAccess).
type AbstractCacheAccessFactory[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord] struct {
	// overrides points back at the concrete factory; see InitAbstractCacheAccessFactory.
	overrides AbstractCacheAccessFactoryOverrides[D, P, V]

	// downloadCache is the global download Cache.
	downloadCache *DownloadCache

	// parsingCache is the global parsing Cache.
	parsingCache *ParsingCache

	// validationCache is the global validation Cache.
	validationCache *ValidationCache
}

// NewAbstractCacheAccessFactory creates an AbstractCacheAccessFactory with fresh caches.
// Port of the protected default constructor.
func NewAbstractCacheAccessFactory[D modeljob.DownloadInfoRecord, P modeljob.ParsingInfoRecord, V modeljob.ValidationInfoRecord]() AbstractCacheAccessFactory[D, P, V] {
	return AbstractCacheAccessFactory[D, P, V]{
		downloadCache:   NewDownloadCache(),
		parsingCache:    NewParsingCache(),
		validationCache: NewValidationCache(),
	}
}

// InitAbstractCacheAccessFactory registers the concrete factory with its base so that the
// base can dispatch to GetCacheAccess/GetReadOnlyCacheAccessTyped. It must be called exactly
// once, by the concrete factory's constructor, before any other method.
func (a *AbstractCacheAccessFactory[D, P, V]) InitAbstractCacheAccessFactory(overrides AbstractCacheAccessFactoryOverrides[D, P, V]) {
	a.overrides = overrides
}

// abstractCacheAccessFactoryOverrides returns the registered overrides, panicking when the
// concrete factory forgot to call InitAbstractCacheAccessFactory.
func (a *AbstractCacheAccessFactory[D, P, V]) abstractCacheAccessFactoryOverrides() AbstractCacheAccessFactoryOverrides[D, P, V] {
	if a.overrides == nil {
		panic("AbstractCacheAccessFactory was not initialised: the concrete factory must call InitAbstractCacheAccessFactory in its constructor")
	}
	return a.overrides
}

// DownloadCache returns the global download Cache (protected field access for subclasses).
func (a *AbstractCacheAccessFactory[D, P, V]) DownloadCache() *DownloadCache {
	return a.downloadCache
}

// ParsingCache returns the global parsing Cache (protected field access for subclasses).
func (a *AbstractCacheAccessFactory[D, P, V]) ParsingCache() *ParsingCache {
	return a.parsingCache
}

// ValidationCache returns the global validation Cache (protected field access for
// subclasses).
func (a *AbstractCacheAccessFactory[D, P, V]) ValidationCache() *ValidationCache {
	return a.validationCache
}

// GetCacheAccess loads a class to deal with a cache by the key records. Port of the abstract
// getCacheAccess(CacheKey), dispatched to the concrete factory.
func (a *AbstractCacheAccessFactory[D, P, V]) GetCacheAccess(key CacheKey) CacheAccessByKey {
	return a.abstractCacheAccessFactoryOverrides().GetCacheAccess(key)
}

// GetDocumentChangesCacheAccess loads a class for document updates. Port of
// getDocumentChangesCacheAccess().
func (a *AbstractCacheAccessFactory[D, P, V]) GetDocumentChangesCacheAccess() *ChangesCacheAccess {
	return NewChangesCacheAccess(a.downloadCache, a.parsingCache, a.validationCache)
}

// GetReadOnlyCacheAccess loads a read-only cache access. Port of getReadOnlyCacheAccess().
// See the type's DEVIATION note.
func (a *AbstractCacheAccessFactory[D, P, V]) GetReadOnlyCacheAccess() ReadOnlyCacheAccess {
	return a.abstractCacheAccessFactoryOverrides().GetReadOnlyCacheAccessTyped()
}

// GetSynchronizerCacheAccess loads a cache access to synchronize records. Port of
// getSynchronizerCacheAccess().
func (a *AbstractCacheAccessFactory[D, P, V]) GetSynchronizerCacheAccess() *SynchronizerCacheAccess {
	return NewSynchronizerCacheAccess(a.downloadCache, a.parsingCache, a.validationCache)
}

// GetDebugCacheAccess loads a cache access to load the information about the current cache
// state. Port of getDebugCacheAccess().
func (a *AbstractCacheAccessFactory[D, P, V]) GetDebugCacheAccess() *DebugCacheAccess {
	return NewDebugCacheAccess(a.downloadCache, a.parsingCache, a.validationCache)
}

var _ CacheAccessFactory = (*AbstractCacheAccessFactory[modeljob.DownloadInfoRecord, modeljob.ParsingInfoRecord, modeljob.ValidationInfoRecord])(nil)
