// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/cache/access/TLCacheAccessFactory.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/validation/job"

// TLCacheAccessFactory accesses the cache for the Trusted Lists validation job.
type TLCacheAccessFactory struct {
	job.AbstractCacheAccessFactory[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO]
}

var _ job.AbstractCacheAccessFactoryOverrides[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO] = (*TLCacheAccessFactory)(nil)

// NewTLCacheAccessFactory is the default constructor.
func NewTLCacheAccessFactory() *TLCacheAccessFactory {
	f := &TLCacheAccessFactory{
		AbstractCacheAccessFactory: job.NewAbstractCacheAccessFactory[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO](),
	}
	f.InitAbstractCacheAccessFactory(f)
	return f
}

// GetCacheAccess ports the covariant override getCacheAccess(CacheKey), satisfying
// job.AbstractCacheAccessFactoryOverrides.
func (f *TLCacheAccessFactory) GetCacheAccess(key job.CacheKey) job.CacheAccessByKey {
	return f.CacheAccess(key)
}

// CacheAccess is the concrete-typed form of GetCacheAccess, used by TL-specific callers that
// need TLCacheAccessByKey's own read methods (e.g. GetParsingReadOnlyResult typed as
// *TLParsingCacheDTO through TLReadOnlyCacheAccess, reached via ReadOnlyCacheAccess() below).
func (f *TLCacheAccessFactory) CacheAccess(key job.CacheKey) *TLCacheAccessByKey {
	return NewTLCacheAccessByKey(key, f.DownloadCache(), f.ParsingCache(), f.ValidationCache())
}

// GetReadOnlyCacheAccessTyped ports the covariant override getReadOnlyCacheAccess(), satisfying
// job.AbstractCacheAccessFactoryOverrides.
func (f *TLCacheAccessFactory) GetReadOnlyCacheAccessTyped() job.ParametrizedReadOnlyCacheAccess[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO] {
	return f.ReadOnlyCacheAccess()
}

// ReadOnlyCacheAccess is the concrete-typed form of GetReadOnlyCacheAccessTyped.
func (f *TLCacheAccessFactory) ReadOnlyCacheAccess() *TLReadOnlyCacheAccess {
	return NewTLReadOnlyCacheAccess(f.DownloadCache(), f.ParsingCache(), f.ValidationCache())
}
