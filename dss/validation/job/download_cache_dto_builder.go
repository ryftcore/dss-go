// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/builder/DownloadCacheDTOBuilder.java (DSS 6.5.RC1).
package job

// DownloadCacheDTOBuilder builds a DownloadCacheDTO.
type DownloadCacheDTOBuilder struct {
	AbstractCacheDTOBuilder[DownloadResult]
}

// NewDownloadCacheDTOBuilder creates a DownloadCacheDTOBuilder for the given download cache
// entry. Port of the public constructor.
func NewDownloadCacheDTOBuilder(cachedEntry *CachedEntry[DownloadResult]) *DownloadCacheDTOBuilder {
	return &DownloadCacheDTOBuilder{AbstractCacheDTOBuilder: NewAbstractCacheDTOBuilder(cachedEntry)}
}

// Build builds the DownloadCacheDTO. Port of build().
func (b *DownloadCacheDTOBuilder) Build() *DownloadCacheDTO {
	downloadCacheDTO := NewDownloadCacheDTOFrom(b.AbstractCacheDTOBuilder.Build())
	if b.IsResultExist() {
		downloadCacheDTO.SetDocument(b.Result().DSSDocument())
		downloadCacheDTO.SetSha2ErrorMessages(b.Result().Sha2ErrorMessages())
	}
	return downloadCacheDTO
}
