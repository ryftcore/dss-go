// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/CacheCleaner.java (DSS 6.5.RC1).
package job

import "github.com/utain/esig/dss/spi/client/http"

// CacheCleaner is used to clean outdated cache entries.
//
// slf4j trace/info/warn logging is dropped (no observable behavior).
type CacheCleaner struct {
	// cleanMemory: if true, removes all map entries with status TO_BE_DELETED. Default: true.
	cleanMemory bool

	// cleanFileSystem: if true, removes files from the file system for each entry with
	// status TO_BE_DELETED. Default: false.
	cleanFileSystem bool

	// dssFileLoader is the DSSCacheFileLoader used to remove files from the File System.
	dssFileLoader http.DSSCacheFileLoader
}

// NewCacheCleaner creates a CacheCleaner with default configuration (cleanMemory = true,
// cleanFileSystem = false) and no file loader. Port of the default constructor.
func NewCacheCleaner() *CacheCleaner {
	return &CacheCleaner{cleanMemory: true}
}

// SetCleanMemory sets the cleanMemory property. Port of setCleanMemory(boolean).
func (c *CacheCleaner) SetCleanMemory(cleanMemory bool) {
	c.cleanMemory = cleanMemory
}

// SetCleanFileSystem sets the cleanFileSystem property. Port of setCleanFileSystem(boolean).
func (c *CacheCleaner) SetCleanFileSystem(cleanFileSystem bool) {
	c.cleanFileSystem = cleanFileSystem
}

// SetDSSFileLoader sets the DSSFileLoader that will be used for file removing. Port of
// setDSSFileLoader(DSSCacheFileLoader).
func (c *CacheCleaner) SetDSSFileLoader(dssFileLoader http.DSSCacheFileLoader) {
	c.dssFileLoader = dssFileLoader
}

// Clean cleans the given entry. Panics if cleanFileSystem is enabled but no DSSFileLoader was
// configured, mirroring Java's Objects.requireNonNull. Port of clean(CacheAccessByKey).
func (c *CacheCleaner) Clean(cacheAccess CacheAccessByKey) {
	fileNeedToBeDeleted := cacheAccess.IsFileNeedToBeDeleted()
	if c.cleanMemory {
		cacheAccess.DeleteDownloadCacheIfNeeded()
		cacheAccess.DeleteParsingCacheIfNeeded()
		cacheAccess.DeleteValidationCacheIfNeeded()
	}

	if c.cleanFileSystem {
		if c.dssFileLoader == nil {
			panic("Cannot remove files from the file system. The DSSFileLoader must be defined!")
		}

		if fileNeedToBeDeleted {
			// The Java catch(Exception) around dssFileLoader.remove(...) only logs a warning
			// on failure; Remove() here returns a bool with no error, so there is nothing to
			// catch (the underlying DSSCacheFileLoader implementation is responsible for not
			// panicking on removal failure).
			_ = c.dssFileLoader.Remove(cacheAccess.GetCacheKey().Key())
		}
	}
}
