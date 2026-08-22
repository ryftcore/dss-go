// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sha2/Sha2FileCacheDataLoader.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// sha2FileCacheDataLoaderOneDayMillis defines a one day constraint in milliseconds. Port of the
// private ONE_DAY_MILLIS constant.
const sha2FileCacheDataLoaderOneDayMillis int64 = 24 * 60 * 60 * 1000 // 24 hours

// Sha2FileCacheDataLoaderOverrides captures every member Sha2FileCacheDataLoader reaches through
// virtual dispatch: the four protected hooks below plus the public DSSCacheFileLoader surface it
// calls on itself (GetDocumentRefresh from GetDocument, GetDocumentFromCache from
// GetDocumentRefresh), all of which a subclass may override.
type Sha2FileCacheDataLoaderOverrides interface {
	http.DSSCacheFileLoader

	// AssertConfigurationIsValid verifies whether the configuration is complete. Port of the
	// protected assertConfigurationIsValid().
	AssertConfigurationIsValid()

	// RefreshedDocument re-downloads the document at documentUrl. Port of the protected
	// getRefreshedDocument(String).
	RefreshedDocument(documentUrl string) (model.DSSDocument, error)

	// Sha2File returns the sha2 file for documentUrl. Port of the protected
	// getSha2File(String).
	Sha2File(documentUrl string) (model.DSSDocument, error)

	// MergeDocumentWithSha2 pairs a cached document with its sha2 document. Port of the
	// protected mergeDocumentWithSha2(DSSDocument, DSSDocument).
	MergeDocumentWithSha2(cachedDocument, sha2Document model.DSSDocument) *DocumentWithSha2

	// CheckRefreshRequired reports whether the cached document should be refreshed. Port of
	// the protected checkRefreshRequired(DocumentWithSha2).
	CheckRefreshRequired(documentWithSha2 *DocumentWithSha2) (bool, error)
}

// Sha2FileCacheDataLoader implements the document loading logic defined within ETSI TS 119 612
// "6.1 TL publication" for Trusted Lists. It tries to access a corresponding .sha2 file for
// every requested document available in the cache, compares its digest, and enforces a document
// update when the document has expired.
//
// Constructors allow manual configuration of the object, and the package-level factories build
// pre-configured objects for Trusted Lists validation:
//   - InitSha2StrictDataLoader enforces refresh of a Trusted List only when a new .sha2 document
//     is obtained or NextUpdate has been reached;
//   - InitSha2DailyUpdateDataLoader enforces refresh when a new .sha2 document is obtained,
//     NextUpdate has been reached, or when the document has not been updated for at least
//     24 hours;
//   - InitSha2CustomExpirationDataLoader enforces refresh when a new .sha2 document is obtained,
//     NextUpdate has been reached, or when the document has not been updated for the indicated
//     time period;
//   - InitSha2IgnoredDataLoader enforces refresh in all cases.
//
// java.io.Serializable has no Go counterpart and is dropped.
type Sha2FileCacheDataLoader struct {
	// overrides points back at the concrete loader; see InitSha2FileCacheDataLoader.
	overrides Sha2FileCacheDataLoaderOverrides

	// dataLoader is the file cache data loader used to load the document.
	dataLoader http.DSSCacheFileLoader

	// predicate checks whether the sha2 document matched the digest of the original document
	// or whether refresh of the document shall be enforced.
	predicate TrustedListWithSha2Predicate
}

var (
	_ http.DSSCacheFileLoader          = (*Sha2FileCacheDataLoader)(nil)
	_ Sha2FileCacheDataLoaderOverrides = (*Sha2FileCacheDataLoader)(nil)
)

// NewSha2FileCacheDataLoader creates an object with an empty configuration (which shall be
// provided with a setter). Port of the default constructor.
func NewSha2FileCacheDataLoader() *Sha2FileCacheDataLoader {
	loader := &Sha2FileCacheDataLoader{}
	loader.InitSha2FileCacheDataLoader(loader)
	return loader
}

// NewSha2FileCacheDataLoaderWithDataLoader creates an object with a defined DSSCacheFileLoader.
// The predicate shall be provided with a setter. Port of
// Sha2FileCacheDataLoader(DSSCacheFileLoader).
func NewSha2FileCacheDataLoaderWithDataLoader(dataLoader http.DSSCacheFileLoader) *Sha2FileCacheDataLoader {
	loader := &Sha2FileCacheDataLoader{dataLoader: dataLoader}
	loader.InitSha2FileCacheDataLoader(loader)
	return loader
}

// InitSha2FileCacheDataLoader registers the concrete loader with its base so that
// GetDocumentRefresh can dispatch the overridable members the way Java reaches an overridden
// method through virtual dispatch. It must be called exactly once, by the concrete loader's
// constructor, before any other method.
func (l *Sha2FileCacheDataLoader) InitSha2FileCacheDataLoader(overrides Sha2FileCacheDataLoaderOverrides) {
	l.overrides = overrides
}

// sha2FileCacheDataLoaderOverrides returns the registered overrides, panicking when the concrete
// loader forgot to call InitSha2FileCacheDataLoader.
func (l *Sha2FileCacheDataLoader) sha2FileCacheDataLoaderOverrides() Sha2FileCacheDataLoaderOverrides {
	if l.overrides == nil {
		panic("Sha2FileCacheDataLoader was not initialised: the concrete loader must call InitSha2FileCacheDataLoader in its constructor")
	}
	return l.overrides
}

// InitSha2StrictDataLoader instantiates a Sha2FileCacheDataLoader with a pre-configured
// predicate, forcing a Trusted List refresh in case of an updated .sha2 document, or when a
// NextUpdate has been reached. The created object does not enforce refresh after a specific time
// period. Port of the static initSha2StrictDataLoader(DSSCacheFileLoader).
func InitSha2StrictDataLoader(dataLoader http.DSSCacheFileLoader) *Sha2FileCacheDataLoader {
	sha2DataLoader := NewSha2FileCacheDataLoaderWithDataLoader(dataLoader)

	sha2Predicate := NewDefaultTrustedListWithSha2Predicate()
	sha2Predicate.SetCacheExpirationTime(-1) // cache do not expire
	sha2DataLoader.SetPredicate(sha2Predicate)

	return sha2DataLoader
}

// InitSha2DailyUpdateDataLoader instantiates a Sha2FileCacheDataLoader with a pre-configured
// predicate, forcing a Trusted List refresh in case of an updated .sha2 document, when a
// NextUpdate has been reached, or when the document has not been re-downloaded for at least a
// day. Port of the static initSha2DailyUpdateDataLoader(DSSCacheFileLoader).
func InitSha2DailyUpdateDataLoader(dataLoader http.DSSCacheFileLoader) *Sha2FileCacheDataLoader {
	sha2DataLoader := NewSha2FileCacheDataLoaderWithDataLoader(dataLoader)

	sha2Predicate := NewDefaultTrustedListWithSha2Predicate()
	sha2Predicate.SetCacheExpirationTime(sha2FileCacheDataLoaderOneDayMillis)
	sha2DataLoader.SetPredicate(sha2Predicate)

	return sha2DataLoader
}

// InitSha2CustomExpirationDataLoader instantiates a Sha2FileCacheDataLoader with a
// pre-configured predicate, forcing a Trusted List refresh in case of an updated .sha2 document,
// when a NextUpdate has been reached, or when the cached document expired according to the
// provided cacheExpirationTime value (in milliseconds). Port of the static
// initSha2CustomExpirationDataLoader(DSSCacheFileLoader, long).
func InitSha2CustomExpirationDataLoader(dataLoader http.DSSCacheFileLoader, cacheExpirationTime int64) *Sha2FileCacheDataLoader {
	sha2DataLoader := NewSha2FileCacheDataLoaderWithDataLoader(dataLoader)

	sha2Predicate := NewDefaultTrustedListWithSha2Predicate()
	sha2Predicate.SetCacheExpirationTime(cacheExpirationTime)
	sha2DataLoader.SetPredicate(sha2Predicate)

	return sha2DataLoader
}

// InitSha2IgnoredDataLoader instantiates a Sha2FileCacheDataLoader with a pre-configured
// predicate, forcing a Trusted List refresh in all cases despite the .sha2 file document
// content. Port of the static initSha2IgnoredDataLoader(DSSCacheFileLoader).
func InitSha2IgnoredDataLoader(dataLoader http.DSSCacheFileLoader) *Sha2FileCacheDataLoader {
	sha2DataLoader := NewSha2FileCacheDataLoaderWithDataLoader(dataLoader)

	sha2Predicate := NewDefaultTrustedListWithSha2Predicate()
	sha2Predicate.SetCacheExpirationTime(0) // cache is always updated
	sha2DataLoader.SetPredicate(sha2Predicate)

	return sha2DataLoader
}

// DataLoader returns the file cache data loader used to load the documents. Port of
// getDataLoader().
func (l *Sha2FileCacheDataLoader) DataLoader() http.DSSCacheFileLoader {
	return l.dataLoader
}

// SetDataLoader sets the file cache data loader to be used to load the documents. Port of
// setDataLoader(DSSCacheFileLoader).
func (l *Sha2FileCacheDataLoader) SetDataLoader(dataLoader http.DSSCacheFileLoader) {
	l.dataLoader = dataLoader
}

// SetPredicate sets a predicate evaluating a condition for a document to be refreshed. The
// predicate returns true when the condition is valid and no document refresh is required, false
// otherwise. Port of setPredicate(Predicate<DocumentWithSha2>).
func (l *Sha2FileCacheDataLoader) SetPredicate(predicate TrustedListWithSha2Predicate) {
	l.predicate = predicate
}

// GetDocument downloads a document from url. Port of the getDocument(String) override.
func (l *Sha2FileCacheDataLoader) GetDocument(url string) (model.DSSDocument, error) {
	return l.sha2FileCacheDataLoaderOverrides().GetDocumentRefresh(url, false)
}

// GetDocumentRefresh downloads a document from url, refreshing the cache when requested. Port of
// the getDocument(String, boolean) override.
//
// Panics with the Java message when url is empty (Objects.requireNonNull(url, "URL cannot be
// null!"); the empty string is the Go stand-in for a null String).
func (l *Sha2FileCacheDataLoader) GetDocumentRefresh(url string, refresh bool) (model.DSSDocument, error) {
	if url == "" {
		panic("URL cannot be null!")
	}

	overrides := l.sha2FileCacheDataLoaderOverrides()
	overrides.AssertConfigurationIsValid()

	var sha2Document model.DSSDocument
	sha2ExtractionStatus := ""
	sha2Document, err := overrides.Sha2File(url)
	if err != nil {
		sha2ExtractionStatus = err.Error()
		// Upstream logs "No sha2 document has been found : %s".
		sha2Document = nil
		refresh = true // force the refresh
	}

	var cachedDocument model.DSSDocument
	var refreshedDocument model.DSSDocument
	if refresh {
		// Upstream logs "Refresh has been requested for a document with URL '{}'" at debug
		// level.
		refreshedDocument, err = overrides.RefreshedDocument(url)
		if err != nil {
			return nil, err
		}

	} else {
		cachedDocument = overrides.GetDocumentFromCache(url)
		if cachedDocument == nil {
			// Upstream logs "No cached document found for URL '{}'" at debug level.
			refreshedDocument, err = overrides.RefreshedDocument(url)
			if err != nil {
				return nil, err
			}
		}
	}

	var documentWithSha2 *DocumentWithSha2
	if cachedDocument != nil {
		documentWithSha2 = overrides.MergeDocumentWithSha2(cachedDocument, sha2Document)
		refreshRequired, err := overrides.CheckRefreshRequired(documentWithSha2)
		if err != nil {
			return nil, err
		}
		if refreshRequired {
			// Upstream logs "Refresh the document from URL '{}'...".
			refreshedDocument, err = overrides.RefreshedDocument(url)
			if err != nil {
				return nil, err
			}
		}
		// Otherwise upstream logs "Sha2 document condition match. Return cached document for
		// URL '{}'" at debug level.
	}

	if refreshedDocument != nil {
		documentWithSha2 = overrides.MergeDocumentWithSha2(refreshedDocument, sha2Document)
		// verify if new document matches sha2
		if _, err := overrides.CheckRefreshRequired(documentWithSha2); err != nil {
			return nil, err
		}
	}

	if documentWithSha2 != nil && utils.IsStringNotEmpty(sha2ExtractionStatus) {
		documentWithSha2.AddErrorMessage(sha2ExtractionStatus)
	}

	// NOTE: Java returns the DocumentWithSha2 reference, which is null when neither a cached
	// nor a refreshed document was produced; a nil *DocumentWithSha2 must not be handed back
	// as a non-nil model.DSSDocument interface value, hence the explicit nil return.
	if documentWithSha2 == nil {
		return nil, nil
	}
	return documentWithSha2, nil
}

// RefreshedDocument re-downloads the document at documentUrl. Port of the protected
// getRefreshedDocument(String).
//
// NOTE: upstream's Javadoc on this method ("returns a document from cache, when applicable")
// describes getDocumentFromCache, not this method; the implementation forces a refresh.
func (l *Sha2FileCacheDataLoader) RefreshedDocument(documentUrl string) (model.DSSDocument, error) {
	return l.dataLoader.GetDocumentRefresh(documentUrl, true)
}

// Sha2File returns the sha2 file for the given documentUrl, or nil when no sha2 document is
// found. Port of the protected getSha2File(String).
func (l *Sha2FileCacheDataLoader) Sha2File(documentUrl string) (model.DSSDocument, error) {
	sha2FileUrl, err := l.Sha2FileUrl(documentUrl)
	if err != nil {
		return nil, err
	}
	return l.dataLoader.GetDocumentRefresh(sha2FileUrl, true)
}

// Sha2FileUrl transforms a documentUrl into the corresponding URL location containing a sha2
// document. Port of the protected getSha2FileUrl(String).
func (l *Sha2FileCacheDataLoader) Sha2FileUrl(documentUrl string) (string, error) {
	// "." is ignored in processing, to allow processing for some exotic cases
	fileExtension := utils.GetFileNameExtension(documentUrl)
	if fileExtension != "" && len(documentUrl) > len(fileExtension) {
		fileExtension = documentUrl[len(documentUrl)-1-len(fileExtension):] // add dot
		documentUrl = documentUrl[:len(documentUrl)-len(fileExtension)]     // remove extension
	}
	if err := l.AssertExtensionIsSupported(fileExtension); err != nil {
		return "", err
	}
	return documentUrl + ".sha2", nil
}

// AssertExtensionIsSupported verifies whether the remote document's fileExtension is supported
// by the implementation: the Trusted Lists distribution points shall end with ".xml" or
// ".xtsl". Port of the protected assertExtensionIsSupported(String); Java's thrown
// DSSExternalResourceException becomes a returned error, per PORTING.md.
func (l *Sha2FileCacheDataLoader) AssertExtensionIsSupported(fileExtension string) error {
	if fileExtension != ".xml" && fileExtension != ".xtsl" {
		return exception.NewDSSExternalResourceException(fmt.Sprintf(
			"The Trusted List extension '%s' is not supported! Shall be one of '.xml' or '.xtsl'.",
			fileExtension))
	}
	return nil
}

// MergeDocumentWithSha2 creates a DocumentWithSha2 by merging a cachedDocument and a
// sha2Document together. Port of the protected mergeDocumentWithSha2(DSSDocument, DSSDocument).
func (l *Sha2FileCacheDataLoader) MergeDocumentWithSha2(cachedDocument, sha2Document model.DSSDocument) *DocumentWithSha2 {
	return NewDocumentWithSha2(cachedDocument, sha2Document)
}

// CheckRefreshRequired reports whether the cached document should be refreshed. Port of the
// protected checkRefreshRequired(DocumentWithSha2).
func (l *Sha2FileCacheDataLoader) CheckRefreshRequired(documentWithSha2 *DocumentWithSha2) (bool, error) {
	accepted, err := l.predicate.Test(documentWithSha2)
	if err != nil {
		return false, err
	}
	return !accepted, nil
}

// GetDocumentFromCache loads a document for a given url from the cache folder, answering nil
// when the document is not found. Port of the getDocumentFromCache(String) override.
//
// Java wraps the delegate call in a try/catch that logs "An error occurred on cached document
// extraction from URL '{}' : {}" and answers null; the Go DSSCacheFileLoader interface declares
// GetDocumentFromCache with no error return at all (a miss is already nil), so there is nothing
// left to catch.
func (l *Sha2FileCacheDataLoader) GetDocumentFromCache(url string) model.DSSDocument {
	l.sha2FileCacheDataLoaderOverrides().AssertConfigurationIsValid()
	return l.dataLoader.GetDocumentFromCache(url)
}

// Remove removes the file with the given url from the cache. Port of the remove(String)
// override.
func (l *Sha2FileCacheDataLoader) Remove(url string) bool {
	l.sha2FileCacheDataLoaderOverrides().AssertConfigurationIsValid()
	return l.dataLoader.Remove(url)
}

// AssertConfigurationIsValid verifies whether the configuration of the loader is complete to
// proceed with execution. Port of the protected assertConfigurationIsValid(); Java's
// Objects.requireNonNull becomes a panic carrying the same message.
func (l *Sha2FileCacheDataLoader) AssertConfigurationIsValid() {
	if l.dataLoader == nil {
		panic("DSSCacheFileLoader shall be provided!")
	}
	if l.predicate == nil {
		panic("Predicate shall be provided!")
	}
}
