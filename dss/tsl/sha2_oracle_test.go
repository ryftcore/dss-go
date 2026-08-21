// Java-oracle tests for the sha2 ports (document_with_sha2.go,
// abstract_trusted_list_with_sha2_predicate.go, default_trusted_list_with_sha2_predicate.go,
// sha2_file_cache_data_loader.go).
//
// Every fixture and every expected verdict/message below is transcribed verbatim from the
// upstream JUnit suites (dss-tsl-validation 6.5.RC1
// eu.europa.esig.dss.tsl.sha2.{DefaultTrustedListWithSha2Predicate,Sha2FileCacheDataLoader}Test).
// testdata/sk-tl.xml and testdata/sk-tl-sn-95.xml are the upstream
// src/test/resources/sk-tl.xml and sk-tl-sn-95.xml, copied unchanged; their SHA-256 digests are
// the two hex strings the upstream tests hard-code as .sha2 payloads.
//
// Test vectors are ported, not the JUnit code, per PORTING.md.
package tsl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/client/http"
)

const (
	// skTLSha2 is the SHA-256 of testdata/sk-tl.xml, HEX-encoded.
	skTLSha2 = "8c43cc710e6d1cc77189c6ca4ef3932e98860575aaaaab77446f167c4fb11618"
	// skTLSN95Sha2 is the SHA-256 of testdata/sk-tl-sn-95.xml, HEX-encoded.
	skTLSN95Sha2 = "c662c9f5252fa9bca9d98d038f3ae2d139f1406c63e2d8b709ba857e140229c1"
)

// sha2TestFileDocument loads one of the two Trusted List fixtures.
func sha2TestFileDocument(t *testing.T, name string) *model.FileDocument {
	t.Helper()
	document, err := model.NewFileDocument("testdata/" + name)
	if err != nil {
		t.Fatalf("unable to load %s: %v", name, err)
	}
	return document
}

// mockSha2Predicate is the Go form of the upstream test's
// MockDefaultTrustedListWithSha2Predicate: a DefaultTrustedListWithSha2Predicate whose
// getCurrentTime() answers a fixed instant.
type mockSha2Predicate struct {
	*DefaultTrustedListWithSha2Predicate
	validationTime time.Time
}

func newMockSha2Predicate(validationTime time.Time) *mockSha2Predicate {
	predicate := &mockSha2Predicate{
		DefaultTrustedListWithSha2Predicate: NewDefaultTrustedListWithSha2Predicate(),
		validationTime:                      validationTime,
	}
	// Re-register so that Test() dispatches CurrentTime to this override, the way Java reaches
	// the subclass through virtual dispatch.
	predicate.InitDefaultTrustedListWithSha2Predicate(predicate)
	return predicate
}

// CurrentTime overrides the base's getCurrentTime().
func (p *mockSha2Predicate) CurrentTime() time.Time { return p.validationTime }

// containsError reports whether the document carries exactly the given error message.
func containsError(documentWithSha2 *DocumentWithSha2, message string) bool {
	for _, e := range documentWithSha2.Errors() {
		if e == message {
			return true
		}
	}
	return false
}

// anyErrorContains reports whether any error message carries all the given substrings.
func anyErrorContains(documentWithSha2 *DocumentWithSha2, substrings ...string) bool {
	for _, e := range documentWithSha2.Errors() {
		matched := true
		for _, substring := range substrings {
			if !strings.Contains(e, substring) {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

// TestDefaultTrustedListWithSha2Predicate_Oracle ports
// DefaultTrustedListWithSha2PredicateTest#test.
func TestDefaultTrustedListWithSha2Predicate_Oracle(t *testing.T) {
	tl := sha2TestFileDocument(t, "sk-tl.xml")
	sha2Document := model.NewInMemoryDocument([]byte(skTLSha2))
	wrongTl := sha2TestFileDocument(t, "sk-tl-sn-95.xml")
	wrongSha2Document := model.NewInMemoryDocument([]byte(skTLSN95Sha2))

	currentTimePredicate := newMockSha2Predicate(time.Now())

	func() {
		defer func() {
			recovered := recover()
			if recovered == nil {
				t.Error("test(null) should have panicked")
			} else if got, want := fmt.Sprint(recovered), "Document shall be provided!"; got != want {
				t.Errorf("panic = %q, want %q", got, want)
			}
		}()
		_, _ = currentTimePredicate.Test(nil)
	}()

	documentWithSha2 := NewDocumentWithSha2(nil, nil)
	if len(documentWithSha2.Errors()) != 0 {
		t.Error("a fresh DocumentWithSha2 should carry no errors")
	}
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !containsError(documentWithSha2, "No cached document has been found.") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	documentWithSha2 = NewDocumentWithSha2(tl, nil)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !containsError(documentWithSha2, "No sha2 document has been found.") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	documentWithSha2 = NewDocumentWithSha2(nil, sha2Document)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !containsError(documentWithSha2, "No cached document has been found.") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	documentWithSha2 = NewDocumentWithSha2(tl, wrongSha2Document)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !anyErrorContains(documentWithSha2, "Digest present within sha2 file", " do not match digest of") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	// The digests match, but the TL's NextUpdate (2020-02-19T00:00:00Z) is in the past.
	documentWithSha2 = NewDocumentWithSha2(tl, sha2Document)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !containsError(documentWithSha2, "NextUpdate '2020-02-19T00:00:00Z' has been reached.") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	pastTimePredicate := newMockSha2Predicate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))

	documentWithSha2 = NewDocumentWithSha2(wrongTl, sha2Document)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !anyErrorContains(documentWithSha2, "Digest present within sha2 file", " do not match digest of") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	// A Trusted List handed in as its own "sha2" file cannot match either.
	documentWithSha2 = NewDocumentWithSha2(tl, wrongTl)
	assertRefreshRequired(t, currentTimePredicate, documentWithSha2)
	if !anyErrorContains(documentWithSha2, "Digest present within sha2 file", " do not match digest of") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}

	documentWithSha2 = NewDocumentWithSha2(tl, sha2Document)
	accepted, err := pastTimePredicate.Test(documentWithSha2)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if !accepted {
		t.Error("the document should be accepted before its NextUpdate")
	}
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none", documentWithSha2.Errors())
	}

	// A zero cache expiration time always forces a refresh, without recording any error.
	pastTimePredicate.SetCacheExpirationTime(0)
	accepted, err = pastTimePredicate.Test(documentWithSha2)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if accepted {
		t.Error("an expired cache should force a refresh")
	}
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none", documentWithSha2.Errors())
	}
}

func assertRefreshRequired(t *testing.T, predicate TrustedListWithSha2Predicate, documentWithSha2 *DocumentWithSha2) {
	t.Helper()
	accepted, err := predicate.Test(documentWithSha2)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if accepted {
		t.Error("a refresh should have been required")
	}
	if len(documentWithSha2.Errors()) == 0 {
		t.Error("an error message should have been recorded")
	}
}

// mockCacheFileLoader is the Go stand-in for the upstream test's FileCacheDataLoader wrapping a
// MockDataLoader: it serves documents from urlMap and keeps whatever it served in an in-memory
// cache. A URL absent from urlMap fails with the MockDataLoader message.
type mockCacheFileLoader struct {
	urlMap map[string]model.DSSDocument
	cache  map[string]model.DSSDocument
}

var _ http.DSSCacheFileLoader = (*mockCacheFileLoader)(nil)

func newMockCacheFileLoader(urlMap map[string]model.DSSDocument) *mockCacheFileLoader {
	return &mockCacheFileLoader{urlMap: urlMap, cache: map[string]model.DSSDocument{}}
}

func (l *mockCacheFileLoader) GetDocument(url string) (model.DSSDocument, error) {
	return l.GetDocumentRefresh(url, false)
}

func (l *mockCacheFileLoader) GetDocumentRefresh(url string, refresh bool) (model.DSSDocument, error) {
	if !refresh {
		if cached, ok := l.cache[url]; ok {
			return cached, nil
		}
	}
	document, ok := l.urlMap[url]
	if !ok {
		return nil, fmt.Errorf("Cannot retrieve data from url [%s]", url)
	}
	l.cache[url] = document
	return document, nil
}

func (l *mockCacheFileLoader) GetDocumentFromCache(url string) model.DSSDocument {
	if cached, ok := l.cache[url]; ok {
		return cached
	}
	return nil
}

func (l *mockCacheFileLoader) Remove(url string) bool {
	if _, ok := l.cache[url]; ok {
		delete(l.cache, url)
		return true
	}
	return false
}

// sha2LoaderURLMap builds the upstream test's url map.
func sha2LoaderURLMap(t *testing.T) map[string]model.DSSDocument {
	t.Helper()
	tl := sha2TestFileDocument(t, "sk-tl.xml")
	tlNew := sha2TestFileDocument(t, "sk-tl-sn-95.xml")
	return map[string]model.DSSDocument{
		"tl_ok.xml":          tl,
		"tl_ok.sha2":         model.NewInMemoryDocument([]byte(skTLSha2)),
		"tl_ko.xml":          tl,
		"tl_ko.sha2":         model.NewInMemoryDocument([]byte(skTLSN95Sha2)),
		"tl_no_sha2.xml":     tl,
		"tl_bad.ext":         tl,
		"tl_bad.sha2":        model.NewInMemoryDocument([]byte(skTLSha2)),
		"tl_no_dot_xml":      tl,
		"tl_no_dot_sha2":     model.NewInMemoryDocument([]byte(skTLSha2)),
		"tl_refresh.xml":     tl,
		"tl_new.xml":         tlNew,
		"tl_refresh.sha2":    model.NewInMemoryDocument([]byte(skTLSN95Sha2)),
		"tl_no_refresh.xml":  tl,
		"tl_no_refresh.sha2": model.NewInMemoryDocument([]byte(skTLSha2)),
	}
}

// newSha2LoaderUnderTest builds the loader the upstream tests exercise: the mock file cache
// loader plus the fixed-time predicate (1 January 2020, i.e. before the fixture's NextUpdate).
func newSha2LoaderUnderTest(t *testing.T, urlMap map[string]model.DSSDocument) (*Sha2FileCacheDataLoader, *mockCacheFileLoader) {
	t.Helper()
	fileDataLoader := newMockCacheFileLoader(urlMap)
	loader := NewSha2FileCacheDataLoaderWithDataLoader(fileDataLoader)
	loader.SetPredicate(newMockSha2Predicate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)))
	return loader, fileDataLoader
}

func sha2TestDigest(t *testing.T, document model.DSSDocument) string {
	t.Helper()
	digest, err := document.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	return fmt.Sprintf("%x", digest)
}

// TestSha2FileCacheDataLoader_OracleGoodDoc ports Sha2FileCacheDataLoaderTest#goodDocTest.
func TestSha2FileCacheDataLoader_OracleGoodDoc(t *testing.T) {
	urlMap := sha2LoaderURLMap(t)
	loader, _ := newSha2LoaderUnderTest(t, urlMap)

	if loader.GetDocumentFromCache("tl_ok.xml") != nil {
		t.Error("nothing should be cached yet")
	}

	document, err := loader.GetDocument("tl_ok.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2, ok := document.(*DocumentWithSha2)
	if !ok {
		t.Fatalf("GetDocument returned %T, want *DocumentWithSha2", document)
	}
	if got, want := sha2TestDigest(t, documentWithSha2.Document()), skTLSha2; got != want {
		t.Errorf("wrapped document digest = %s, want %s", got, want)
	}
	if got, want := sha2TestDigest(t, documentWithSha2.Sha2Document()),
		sha2TestDigest(t, urlMap["tl_ok.sha2"]); got != want {
		t.Errorf("sha2 document digest = %s, want %s", got, want)
	}
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none", documentWithSha2.Errors())
	}
	// The DocumentWithSha2 itself streams the wrapped document.
	if got, want := sha2TestDigest(t, document), skTLSha2; got != want {
		t.Errorf("DocumentWithSha2 digest = %s, want %s", got, want)
	}

	documentFromCache := loader.GetDocumentFromCache("tl_ok.xml")
	if documentFromCache == nil {
		t.Fatal("the document should now be cached")
	}
	if got, want := sha2TestDigest(t, documentFromCache), skTLSha2; got != want {
		t.Errorf("cached digest = %s, want %s", got, want)
	}

	refreshedDocument, err := loader.RefreshedDocument("tl_ok.xml")
	if err != nil {
		t.Fatalf("RefreshedDocument: %v", err)
	}
	if got, want := sha2TestDigest(t, refreshedDocument), skTLSha2; got != want {
		t.Errorf("refreshed digest = %s, want %s", got, want)
	}

	sha2File, err := loader.Sha2File("tl_ok.xml")
	if err != nil {
		t.Fatalf("Sha2File: %v", err)
	}
	if got, want := sha2TestDigest(t, sha2File), sha2TestDigest(t, urlMap["tl_ok.sha2"]); got != want {
		t.Errorf("sha2 file digest = %s, want %s", got, want)
	}

	sha2FileUrl, err := loader.Sha2FileUrl("tl_ok.xml")
	if err != nil {
		t.Fatalf("Sha2FileUrl: %v", err)
	}
	if sha2FileUrl != "tl_ok.sha2" {
		t.Errorf("Sha2FileUrl = %q, want %q", sha2FileUrl, "tl_ok.sha2")
	}
}

// TestSha2FileCacheDataLoader_OracleWrongDigest ports
// Sha2FileCacheDataLoaderTest#wrongDigestDocTest.
func TestSha2FileCacheDataLoader_OracleWrongDigest(t *testing.T) {
	loader, _ := newSha2LoaderUnderTest(t, sha2LoaderURLMap(t))

	document, err := loader.GetDocument("tl_ko.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 := document.(*DocumentWithSha2)
	if got := len(documentWithSha2.Errors()); got != 1 {
		t.Fatalf("errors = %v, want exactly one", documentWithSha2.Errors())
	}
	if !strings.Contains(documentWithSha2.Errors()[0], "do not match digest of the cached document") {
		t.Errorf("error = %q", documentWithSha2.Errors()[0])
	}
}

// TestSha2FileCacheDataLoader_OracleNoSha2 ports Sha2FileCacheDataLoaderTest#noSha2DocTest.
func TestSha2FileCacheDataLoader_OracleNoSha2(t *testing.T) {
	loader, _ := newSha2LoaderUnderTest(t, sha2LoaderURLMap(t))

	document, err := loader.GetDocument("tl_no_sha2.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 := document.(*DocumentWithSha2)
	if documentWithSha2.Sha2Document() != nil {
		t.Error("no sha2 document should have been found")
	}
	if got := len(documentWithSha2.Errors()); got != 2 {
		t.Fatalf("errors = %v, want exactly two", documentWithSha2.Errors())
	}
	if !anyErrorContains(documentWithSha2, "No sha2 document has been found") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}
	if !anyErrorContains(documentWithSha2, "Cannot retrieve data from url") {
		t.Errorf("errors = %v", documentWithSha2.Errors())
	}
}

// TestSha2FileCacheDataLoader_OracleBadExtension ports
// Sha2FileCacheDataLoaderTest#badExtensionDocTest and #extensionNoDotDocTest.
func TestSha2FileCacheDataLoader_OracleBadExtension(t *testing.T) {
	cases := []struct {
		url       string
		wantError string
	}{
		{"tl_bad.ext", "The Trusted List extension '.ext' is not supported! Shall be one of '.xml' or '.xtsl'."},
		{"tl_no_dot_xml", "The Trusted List extension '' is not supported! Shall be one of '.xml' or '.xtsl'."},
	}
	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			loader, _ := newSha2LoaderUnderTest(t, sha2LoaderURLMap(t))
			document, err := loader.GetDocument(c.url)
			if err != nil {
				t.Fatalf("GetDocument: %v", err)
			}
			documentWithSha2 := document.(*DocumentWithSha2)
			if documentWithSha2.Sha2Document() != nil {
				t.Error("no sha2 document should have been found")
			}
			if got := len(documentWithSha2.Errors()); got != 2 {
				t.Fatalf("errors = %v, want exactly two", documentWithSha2.Errors())
			}
			if !anyErrorContains(documentWithSha2, "No sha2 document has been found") {
				t.Errorf("errors = %v", documentWithSha2.Errors())
			}
			if !anyErrorContains(documentWithSha2, c.wantError) {
				t.Errorf("errors = %v, want one containing %q", documentWithSha2.Errors(), c.wantError)
			}
		})
	}
}

// TestSha2FileCacheDataLoader_OracleRefresh ports Sha2FileCacheDataLoaderTest#refreshTest and
// #noRefreshTest: a mismatching .sha2 forces a re-download, a matching one does not.
func TestSha2FileCacheDataLoader_OracleRefresh(t *testing.T) {
	urlMap := sha2LoaderURLMap(t)
	loader, _ := newSha2LoaderUnderTest(t, urlMap)

	document, err := loader.GetDocument("tl_refresh.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 := document.(*DocumentWithSha2)
	if len(documentWithSha2.Errors()) == 0 {
		t.Fatal("the mismatching sha2 should have been reported")
	}
	if !strings.Contains(documentWithSha2.Errors()[0], "do not match digest of the cached document") {
		t.Errorf("error = %q", documentWithSha2.Errors()[0])
	}

	// The published document is replaced by the one the .sha2 actually describes.
	urlMap["tl_refresh.xml"] = urlMap["tl_new.xml"]

	document, err = loader.GetDocument("tl_refresh.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 = document.(*DocumentWithSha2)
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none once the digests agree", documentWithSha2.Errors())
	}
	if got, want := sha2TestDigest(t, documentWithSha2.Document()), skTLSN95Sha2; got != want {
		t.Errorf("refreshed document digest = %s, want %s", got, want)
	}
}

// TestSha2FileCacheDataLoader_OracleNoRefresh ports Sha2FileCacheDataLoaderTest#noRefreshTest.
func TestSha2FileCacheDataLoader_OracleNoRefresh(t *testing.T) {
	urlMap := sha2LoaderURLMap(t)
	loader, _ := newSha2LoaderUnderTest(t, urlMap)

	document, err := loader.GetDocument("tl_no_refresh.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	if len(document.(*DocumentWithSha2).Errors()) != 0 {
		t.Fatalf("errors = %v, want none", document.(*DocumentWithSha2).Errors())
	}

	// Replacing the published document must NOT be picked up: the cached copy still matches
	// its .sha2, so no refresh is triggered.
	urlMap["tl_no_refresh.xml"] = urlMap["tl_new.xml"]

	document, err = loader.GetDocument("tl_no_refresh.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 := document.(*DocumentWithSha2)
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none", documentWithSha2.Errors())
	}
	if got, want := sha2TestDigest(t, documentWithSha2.Document()), skTLSha2; got != want {
		t.Errorf("document digest = %s, want the cached %s", got, want)
	}
}

// TestSha2FileCacheDataLoader_OracleErrorOnDownload ports
// Sha2FileCacheDataLoaderTest#errorOnDownloadTest.
func TestSha2FileCacheDataLoader_OracleErrorOnDownload(t *testing.T) {
	urlMap := map[string]model.DSSDocument{
		"tl_ok.xml":  sha2TestFileDocument(t, "sk-tl.xml"),
		"tl_ok.sha2": model.NewInMemoryDocument([]byte(skTLSha2)),
	}
	loader, fileDataLoader := newSha2LoaderUnderTest(t, urlMap)

	if loader.GetDocumentFromCache("tl_ok.xml") != nil {
		t.Error("nothing should be cached yet")
	}
	document, err := loader.GetDocument("tl_ok.xml")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	documentWithSha2 := document.(*DocumentWithSha2)
	if documentWithSha2.Document() == nil || documentWithSha2.Sha2Document() == nil {
		t.Fatal("both documents should have been retrieved")
	}
	if len(documentWithSha2.Errors()) != 0 {
		t.Errorf("errors = %v, want none", documentWithSha2.Errors())
	}

	// Once nothing can be downloaded any more, the refresh the missing .sha2 forces fails.
	for key := range urlMap {
		delete(urlMap, key)
	}
	for key := range fileDataLoader.cache {
		delete(fileDataLoader.cache, key)
	}

	if _, err := loader.GetDocument("tl_ok.xml"); err == nil {
		t.Error("the failed download should have been reported")
	} else if got, want := err.Error(), "Cannot retrieve data from url [tl_ok.xml]"; got != want {
		t.Errorf("error = %q, want %q", got, want)
	}
}

// TestSha2FileCacheDataLoader_OracleNullChecks ports Sha2FileCacheDataLoaderTest#nullTest: the
// Objects.requireNonNull messages, which the Go port raises as panics.
func TestSha2FileCacheDataLoader_OracleNullChecks(t *testing.T) {
	assertPanics := func(name, want string, call func()) {
		t.Helper()
		defer func() {
			recovered := recover()
			if recovered == nil {
				t.Errorf("%s should have panicked", name)
			} else if got := fmt.Sprint(recovered); got != want {
				t.Errorf("%s panic = %q, want %q", name, got, want)
			}
		}()
		call()
	}

	assertPanics("empty url", "URL cannot be null!", func() {
		_, _ = NewSha2FileCacheDataLoader().GetDocument("")
	})
	assertPanics("no data loader", "DSSCacheFileLoader shall be provided!", func() {
		_, _ = NewSha2FileCacheDataLoader().GetDocument("tl_ok.xml")
	})
	assertPanics("nil data loader", "DSSCacheFileLoader shall be provided!", func() {
		_, _ = NewSha2FileCacheDataLoaderWithDataLoader(nil).GetDocument("tl_ok.xml")
	})

	loader := NewSha2FileCacheDataLoaderWithDataLoader(newMockCacheFileLoader(sha2LoaderURLMap(t)))
	assertPanics("no predicate", "Predicate shall be provided!", func() {
		_, _ = loader.GetDocument("tl_ok.xml")
	})
}

// TestSha2FileCacheDataLoader_Factories covers the four pre-configured factory methods.
func TestSha2FileCacheDataLoader_Factories(t *testing.T) {
	fileDataLoader := newMockCacheFileLoader(map[string]model.DSSDocument{})
	for _, loader := range []*Sha2FileCacheDataLoader{
		InitSha2StrictDataLoader(fileDataLoader),
		InitSha2DailyUpdateDataLoader(fileDataLoader),
		InitSha2CustomExpirationDataLoader(fileDataLoader, 1234),
		InitSha2IgnoredDataLoader(fileDataLoader),
	} {
		if loader.DataLoader() != fileDataLoader {
			t.Error("the factory should keep the provided data loader")
		}
		if loader.predicate == nil {
			t.Error("the factory should install a predicate")
		}
		// The configuration is complete, so the assertion must not panic.
		loader.AssertConfigurationIsValid()
	}
	if got := InitSha2DailyUpdateDataLoader(fileDataLoader).predicate.(*DefaultTrustedListWithSha2Predicate).cacheExpirationTime; got != 24*60*60*1000 {
		t.Errorf("daily-update cache expiration = %d, want 86400000", got)
	}
	if got := InitSha2StrictDataLoader(fileDataLoader).predicate.(*DefaultTrustedListWithSha2Predicate).cacheExpirationTime; got != -1 {
		t.Errorf("strict cache expiration = %d, want -1", got)
	}
	if got := InitSha2IgnoredDataLoader(fileDataLoader).predicate.(*DefaultTrustedListWithSha2Predicate).cacheExpirationTime; got != 0 {
		t.Errorf("ignored cache expiration = %d, want 0", got)
	}
	if got := InitSha2CustomExpirationDataLoader(fileDataLoader, 1234).predicate.(*DefaultTrustedListWithSha2Predicate).cacheExpirationTime; got != 1234 {
		t.Errorf("custom cache expiration = %d, want 1234", got)
	}
}
