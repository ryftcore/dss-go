// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/ParsingCache.java (DSS 6.5.RC1).
package job

// ParsingCache stores results of document parsings.
type ParsingCache struct {
	AbstractCache[ParsingResult]
}

// NewParsingCache creates an empty ParsingCache. Port of the default constructor.
func NewParsingCache() *ParsingCache {
	c := &ParsingCache{AbstractCache: NewAbstractCache[ParsingResult]()}
	c.InitAbstractCache(c)
	return c
}

// CacheType returns CacheTypeParsing. Port of getCacheType().
func (c *ParsingCache) CacheType() CacheType {
	return CacheTypeParsing
}
