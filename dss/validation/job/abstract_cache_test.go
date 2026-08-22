package job

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// fakeDownloadResult is a minimal DownloadResult for testing DownloadCache/AbstractCache.
type fakeDownloadResult struct {
	doc               model.DSSDocument
	digest            model.Digest
	sha2ErrorMessages []string
}

func (r *fakeDownloadResult) DSSDocument() model.DSSDocument { return r.doc }
func (r *fakeDownloadResult) Digest() model.Digest           { return r.digest }
func (r *fakeDownloadResult) Sha2ErrorMessages() []string    { return r.sha2ErrorMessages }

func newFakeDownloadResult(content string, sha2ErrorMessages []string) *fakeDownloadResult {
	return &fakeDownloadResult{
		doc:               model.NewInMemoryDocument([]byte(content)),
		digest:            model.NewDigest(enumerations.DigestAlgorithmSHA256, []byte(content)),
		sha2ErrorMessages: sha2ErrorMessages,
	}
}

func TestAbstractCache_GetCreatesEmptyEntryAndIsStable(t *testing.T) {
	c := NewDownloadCache()
	key := NewCacheKey("https://example.org/tl.xml")

	e1 := c.Get(key)
	if !e1.IsEmpty() {
		t.Fatal("a freshly created entry should be empty")
	}
	e2 := c.Get(key)
	if e1 != e2 {
		t.Fatal("Get() should return the same *CachedEntry on repeated calls for the same key")
	}
}

func TestAbstractCache_UpdateExpireRemoveToBeDeleted(t *testing.T) {
	c := NewParsingCache()
	key := NewCacheKey("https://example.org/tl.xml")

	if !c.IsRefreshNeeded(key) {
		t.Fatal("a fresh key should need a refresh")
	}

	c.ToBeDeleted(key)
	if !c.IsToBeDeleted(key) {
		t.Fatal("expected TO_BE_DELETED after ToBeDeleted()")
	}

	c.Remove(key)
	if !c.IsEmpty(key) {
		t.Fatal("after Remove, Get should hand back a fresh empty entry")
	}
	if !c.IsRefreshNeeded(key) {
		t.Fatal("the re-created entry after Remove should need a refresh again")
	}
}

func TestAbstractCache_Error(t *testing.T) {
	c := NewValidationCache()
	key := NewCacheKey("https://example.org/tl.xml")

	c.Error(key, testError("boom"))
	if !c.Get(key).IsError() {
		t.Fatal("expected ERROR state after Error()")
	}
}

func TestDownloadCache_IsUpToDate(t *testing.T) {
	c := NewDownloadCache()
	key := NewCacheKey("https://example.org/tl.xml")

	first := newFakeDownloadResult("hello", nil)
	if c.IsUpToDate(key, first) {
		t.Fatal("an empty cache entry is never up to date")
	}
	c.Update(key, first)
	c.Sync(key) // DESYNCHRONIZED -> SYNCHRONIZED, so ToBeDeleted() below is a legal transition.

	// Same digest, both sha2 message lists empty => up to date.
	same := newFakeDownloadResult("hello", nil)
	if !c.IsUpToDate(key, same) {
		t.Fatal("same content digest and no sha2 errors should be up to date")
	}

	// Different content => different digest => not up to date.
	different := newFakeDownloadResult("world", nil)
	if c.IsUpToDate(key, different) {
		t.Fatal("different content digest should not be up to date")
	}

	// TO_BE_DELETED forces an update regardless of digest.
	c.ToBeDeleted(key)
	if c.IsUpToDate(key, same) {
		t.Fatal("a TO_BE_DELETED entry should never be reported up to date")
	}
}

func TestDownloadCache_IsSHA2ContentMatch_JavaVerbatimQuirk(t *testing.T) {
	// Ports the exact (asymmetric) condition from DownloadCache#isSHA2ContentMatch: the
	// downloadedResult's own sha2ErrorMessages emptiness is never checked on the
	// list-equality branch, only the cachedResult's.
	cached := newFakeDownloadResult("hello", nil)
	downloadedWithErrors := newFakeDownloadResult("hello", []string{"sha2 mismatch"})

	// cachedResult.Sha2ErrorMessages() is empty, downloadedResult's is not: the first
	// disjunct requires both empty (fails, since downloaded has an error); the second
	// disjunct only requires cachedResult empty AND list-equality — nil != []string{...},
	// so it also fails. isSHA2ContentMatch is false, matching the Java oracle for this case.
	if isSHA2ContentMatch(cached, downloadedWithErrors) {
		t.Fatal("cached has no sha2 errors but downloaded does: should not match")
	}

	bothEmpty := newFakeDownloadResult("hello", nil)
	if !isSHA2ContentMatch(cached, bothEmpty) {
		t.Fatal("both empty sha2 error lists should match")
	}
}
