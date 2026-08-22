// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/ValidationCache.java (DSS 6.5.RC1).
package job

// ValidationCache stores validation information for processed files.
type ValidationCache struct {
	AbstractCache[ValidationResult]
}

// NewValidationCache creates an empty ValidationCache. Port of the default constructor.
func NewValidationCache() *ValidationCache {
	c := &ValidationCache{AbstractCache: NewAbstractCache[ValidationResult]()}
	c.InitAbstractCache(c)
	return c
}

// CacheType returns CacheTypeValidation. Port of getCacheType().
func (c *ValidationCache) CacheType() CacheType {
	return CacheTypeValidation
}
