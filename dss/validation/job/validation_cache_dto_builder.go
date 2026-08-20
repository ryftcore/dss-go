// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/builder/ValidationCacheDTOBuilder.java (DSS 6.5.RC1).
package job

// ValidationCacheDTOBuilder builds a ValidationCacheDTO.
type ValidationCacheDTOBuilder struct {
	AbstractCacheDTOBuilder[ValidationResult]
}

// NewValidationCacheDTOBuilder creates a ValidationCacheDTOBuilder for the given validation
// cache entry. Port of the public constructor.
func NewValidationCacheDTOBuilder(cachedEntry *CachedEntry[ValidationResult]) *ValidationCacheDTOBuilder {
	return &ValidationCacheDTOBuilder{AbstractCacheDTOBuilder: NewAbstractCacheDTOBuilder(cachedEntry)}
}

// Build builds the ValidationCacheDTO. Port of build().
func (b *ValidationCacheDTOBuilder) Build() *ValidationCacheDTO {
	validationCacheDTO := NewValidationCacheDTOFrom(b.AbstractCacheDTOBuilder.Build())
	if b.IsResultExist() {
		validationCacheDTO.SetIndication(b.Result().Indication())
		validationCacheDTO.SetSubIndication(b.Result().SubIndication())
		validationCacheDTO.SetSigningTime(b.Result().SigningTime())
		validationCacheDTO.SetSigningCertificate(b.Result().SigningCertificate())
		validationCacheDTO.SetPotentialSigners(b.Result().PotentialSigners())
	}
	return validationCacheDTO
}
