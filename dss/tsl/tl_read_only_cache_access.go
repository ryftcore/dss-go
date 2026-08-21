// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/job/TLReadOnlyCacheAccess.java (DSS 6.5.RC1).
package tsl

import (
	modeljob "github.com/utain/esig/dss/model/job"
	"github.com/utain/esig/dss/validation/job"
)

// TLReadOnlyCacheAccess accesses the Trusted List cache in a read-only mode.
type TLReadOnlyCacheAccess struct {
	job.AbstractCacheAccess
}

var _ job.ParametrizedReadOnlyCacheAccess[*job.DownloadCacheDTO, *TLParsingCacheDTO, *job.ValidationCacheDTO] = (*TLReadOnlyCacheAccess)(nil)

// NewTLReadOnlyCacheAccess is the default constructor.
func NewTLReadOnlyCacheAccess(fileCache *job.DownloadCache, parsingCache *job.ParsingCache, validationCache *job.ValidationCache) *TLReadOnlyCacheAccess {
	return &TLReadOnlyCacheAccess{AbstractCacheAccess: job.NewAbstractCacheAccess(fileCache, parsingCache, validationCache)}
}

// GetDownloadInfoRecordTyped ports the covariant getDownloadInfoRecord(CacheKey) override.
func (a *TLReadOnlyCacheAccess) GetDownloadInfoRecordTyped(key job.CacheKey) *job.DownloadCacheDTO {
	downloadCacheEntry := a.GetDownloadCacheEntry(key)
	return job.NewDownloadCacheDTOBuilder(downloadCacheEntry).Build()
}

// GetParsingInfoRecordTyped ports the covariant getParsingInfoRecord(CacheKey) override.
func (a *TLReadOnlyCacheAccess) GetParsingInfoRecordTyped(key job.CacheKey) *TLParsingCacheDTO {
	parsingCacheEntry := a.GetParsingCacheEntry(key)
	return NewTLParsingCacheDTOBuilder(parsingCacheEntry).Build()
}

// GetValidationInfoRecordTyped ports the covariant getValidationInfoRecord(CacheKey) override.
func (a *TLReadOnlyCacheAccess) GetValidationInfoRecordTyped(key job.CacheKey) *job.ValidationCacheDTO {
	validationCacheEntry := a.GetValidationCacheEntry(key)
	return job.NewValidationCacheDTOBuilder(validationCacheEntry).Build()
}

// GetDownloadInfoRecord satisfies job.ReadOnlyCacheAccess with the base interface return type
// (Go has no covariant return types - see AbstractReadOnlyCacheAccessByKey's DEVIATION note in
// the job package for the identical rationale).
func (a *TLReadOnlyCacheAccess) GetDownloadInfoRecord(key job.CacheKey) modeljob.DownloadInfoRecord {
	return a.GetDownloadInfoRecordTyped(key)
}

// GetParsingInfoRecord satisfies job.ReadOnlyCacheAccess with the base interface return type.
func (a *TLReadOnlyCacheAccess) GetParsingInfoRecord(key job.CacheKey) modeljob.ParsingInfoRecord {
	return a.GetParsingInfoRecordTyped(key)
}

// GetValidationInfoRecord satisfies job.ReadOnlyCacheAccess with the base interface return type.
func (a *TLReadOnlyCacheAccess) GetValidationInfoRecord(key job.CacheKey) modeljob.ValidationInfoRecord {
	return a.GetValidationInfoRecordTyped(key)
}
