// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/cache/CacheCleaner.java (DSS 6.5.RC1).
//
// Deprecated since DSS 6.5 in favor of eu.europa.esig.dss.validation.job.cache.CacheCleaner
// (job.CacheCleaner - see xml_download_result.go's header for the wider job.* convention this
// port follows). Java expresses "deprecated wrapper with no added
// behaviour" via subclassing; Go has no inheritance, so this type embeds job.CacheCleaner
// directly, promoting its exported behaviour unchanged.
package tsl

import "github.com/ryftcore/dss-go/dss/validation/job"

// CacheCleaner is used to clean outdated cache entries.
//
// Deprecated: since DSS 6.5. Use github.com/ryftcore/dss-go/dss/validation/job.CacheCleaner instead.
type CacheCleaner struct {
	*job.CacheCleaner
}

// NewCacheCleaner instantiates a CacheCleaner with default configuration and a nil file loader.
// Port of the default constructor.
//
// Deprecated: since DSS 6.5. Use job.NewCacheCleaner instead.
func NewCacheCleaner() *CacheCleaner {
	return &CacheCleaner{CacheCleaner: job.NewCacheCleaner()}
}
