// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/cache/access/TLCacheAccessByKey.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see xml_download_result.go's header for the wider job.* convention):
// extends eu.europa.esig.dss.validation.job.cache.access.AbstractCacheAccessByKey<D, P, V>,
// specialized here to (job.DownloadCacheDTO, TLParsingCacheDTO, job.ValidationCacheDTO). Go has
// no generic-class inheritance with covariant type-parameter binding the way Java's
// "extends AbstractCacheAccessByKey<DownloadCacheDTO, TLParsingCacheDTO, ValidationCacheDTO>"
// does, so this embeds the assumed generic
// job.AbstractCacheAccessByKey[job.DownloadCacheDTO, TLParsingCacheDTO, job.ValidationCacheDTO]
// instantiation by value, the same pattern tl_parsing_cache_dto.go (TSLCORE chunk) already
// established for job.AbstractParsingCacheDTO.
package tsl

import "github.com/ryftcore/dss-go/dss/validation/job"

// TLCacheAccessByKey accesses cache information for a Trusted List by key.
type TLCacheAccessByKey struct {
	job.AbstractCacheAccessByKey[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO]
}

// NewTLCacheAccessByKey is the default constructor.
func NewTLCacheAccessByKey(key job.CacheKey, downloadCache *job.DownloadCache, parsingCache *job.ParsingCache,
	validationCache *job.ValidationCache) *TLCacheAccessByKey {
	return &TLCacheAccessByKey{
		AbstractCacheAccessByKey: job.NewAbstractCacheAccessByKey[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO](
			key, downloadCache, parsingCache, validationCache, NewTLReadOnlyCacheAccess(downloadCache, parsingCache, validationCache)),
	}
}
