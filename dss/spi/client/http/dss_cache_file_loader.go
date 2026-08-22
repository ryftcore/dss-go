// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/DSSCacheFileLoader.java (DSS 6.5.RC1).
package http

import "github.com/ryftcore/dss-go/dss/model"

// DSSCacheFileLoader implements a file loader implementing a caching
// mechanism, allowing to remove cache externally (to be used within a
// CacheCleaner).
type DSSCacheFileLoader interface {
	DSSFileLoader

	// GetDocumentRefresh downloads a DSSDocument from the specified url with
	// a custom setting indicating whether the refresh of the document's
	// cache shall be enforced, when applicable.
	GetDocumentRefresh(url string, refresh bool) (model.DSSDocument, error)

	// GetDocumentFromCache loads a document for a given url from the cache
	// folder. If the document is not found in the cache, returns nil.
	GetDocumentFromCache(url string) model.DSSDocument

	// Remove removes the file from cache with the given url. Returns true
	// when the file was successfully deleted, false otherwise.
	Remove(url string) bool
}
