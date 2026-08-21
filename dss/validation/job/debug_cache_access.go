// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/access/DebugCacheAccess.java (DSS 6.5.RC1).
package job

// DebugCacheAccess generates and prints a report of the current cache state.
//
// slf4j info logging is dropped, but the rendered report string is still produced and
// returned by Dump so callers that want it (e.g. for their own logging) can use it.
type DebugCacheAccess struct {
	// downloadCache is the global download Cache.
	downloadCache *DownloadCache

	// parsingCache is the global parsing Cache.
	parsingCache *ParsingCache

	// validationCache is the global validation Cache.
	validationCache *ValidationCache
}

// NewDebugCacheAccess creates a DebugCacheAccess. Port of the default constructor.
func NewDebugCacheAccess(downloadCache *DownloadCache, parsingCache *ParsingCache, validationCache *ValidationCache) *DebugCacheAccess {
	return &DebugCacheAccess{downloadCache: downloadCache, parsingCache: parsingCache, validationCache: validationCache}
}

// Dump renders the report for the current cache state. Port of dump(), returning the
// rendered string instead of only logging it (see the type's note on dropped logging).
func (d *DebugCacheAccess) Dump() string {
	var sb []byte
	sb = append(sb, "Cache contents"...)
	sb = append(sb, '\n')
	sb = append(sb, d.downloadCache.Dump()...)
	sb = append(sb, '\n')
	sb = append(sb, d.parsingCache.Dump()...)
	sb = append(sb, '\n')
	sb = append(sb, d.validationCache.Dump()...)
	return string(sb)
}
