package job

import "testing"

// fakeResult is a minimal CachedResult for testing CachedEntry[R] generically.
type fakeResult struct{ v int }

func TestCachedEntry_NewIsRefreshNeededAndEmpty(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()
	if !e.IsRefreshNeeded() {
		t.Fatal("new entry should be REFRESH_NEEDED")
	}
	if !e.IsEmpty() {
		t.Fatal("new entry should be empty")
	}
	if e.IsResultExist() {
		t.Fatal("new entry should not have a result")
	}
}

func TestCachedEntry_UpdatePanicsOnNil(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic updating with a nil result")
		}
	}()
	e.Update(nil)
}

func TestCachedEntry_UpdateThenSync(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()
	e.Update(&fakeResult{v: 1})
	if e.CurrentState() != CacheStateEnum_DESYNCHRONIZED {
		t.Fatalf("after update want DESYNCHRONIZED, got %s", e.CurrentState())
	}
	if e.IsEmpty() {
		t.Fatal("entry should no longer be empty")
	}
	e.Sync()
	if e.CurrentState() != CacheStateEnum_SYNCHRONIZED {
		t.Fatalf("after sync want SYNCHRONIZED, got %s", e.CurrentState())
	}
	if e.LastSuccessSynchronizationTime().IsZero() {
		t.Fatal("sync should set LastSuccessSynchronizationTime")
	}
}

func TestCachedEntry_ErrorResetsResultAndIsSticky(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()
	e.Update(&fakeResult{v: 1})
	e.Sync()

	e.Error(NewCachedExceptionWrapper(testError("first failure")))
	if !e.IsError() {
		t.Fatal("expected error state")
	}
	if !e.IsEmpty() {
		t.Fatal("error() should reset the cached result")
	}
	if e.ExceptionMessage() != "first failure" {
		t.Fatalf("ExceptionMessage() = %q", e.ExceptionMessage())
	}
	firstOccurrence := e.ExceptionFirstOccurrenceTime()
	if firstOccurrence.IsZero() {
		t.Fatal("expected a first occurrence time")
	}

	// Same stack trace (same error value/type) => not a new error, only the last-occurrence
	// date moves; the first-occurrence date and message are unchanged. Mirrors Java's
	// CachedEntry#isNewError comparing stack traces.
	e.Error(NewCachedExceptionWrapper(testError("first failure")))
	if e.ExceptionFirstOccurrenceTime() != firstOccurrence {
		t.Fatal("re-recording the same error should not move the first occurrence time")
	}

	// A different error message is a new error. CurrentCacheContext.Error() (called here)
	// sets the ERROR state directly, bypassing the CacheState transition table entirely
	// (mirroring Java's CurrentCacheContext#error(CachedExceptionWrapper), which never
	// delegates to CacheState#error) — so, unlike the ERROR->ERROR case exercised through
	// the transition table in cache_state_enum_test.go, this does not panic.
	e.Error(NewCachedExceptionWrapper(testError("second, different failure")))
	if !e.IsError() {
		t.Fatal("expected error state")
	}
	if e.ExceptionMessage() != "second, different failure" {
		t.Fatalf("ExceptionMessage() = %q, want the new message", e.ExceptionMessage())
	}
	if e.ExceptionFirstOccurrenceTime() == firstOccurrence {
		t.Fatal("a genuinely new error should move the first occurrence time")
	}
}

func TestCachedEntry_ExpireAndToBeDeleted(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()
	e.Update(&fakeResult{v: 1})
	e.Sync()

	e.Expire()
	if !e.IsRefreshNeeded() {
		t.Fatal("expire() should move to REFRESH_NEEDED")
	}

	e.ToBeDeleted()
	if !e.IsToBeDeleted() {
		t.Fatal("expected TO_BE_DELETED state")
	}
}
