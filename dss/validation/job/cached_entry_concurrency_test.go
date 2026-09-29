package job

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// The tests in this file are the regression tests for the data race on the shared
// *CachedEntry state (T31-SEC-001 / T18A-SEC-001): AbstractCache.Get hands the very same
// entry to every goroutine, and the TL/LOTL job expires and reads shared LOTL/pivot validation
// entries from several goroutines at once. They assert nothing beyond "no panic and a coherent
// final state" - their value is under `go test -race`, where any unsynchronized access to the
// entry's state, exception or result is reported. Without CachedEntry's mutex they fail there.

// concurrencyWorkers and concurrencyRounds size the hammering; small enough to stay quick under
// the race detector, large enough that the detector sees the overlapping accesses reliably.
const (
	concurrencyWorkers = 8
	concurrencyRounds  = 200
)

// runConcurrently runs fn(worker) on concurrencyWorkers goroutines released together, and waits.
func runConcurrently(fn func(worker int)) {
	var ready, done sync.WaitGroup
	start := make(chan struct{})
	ready.Add(concurrencyWorkers)
	done.Add(concurrencyWorkers)
	for w := 0; w < concurrencyWorkers; w++ {
		go func(worker int) {
			defer done.Done()
			ready.Done()
			<-start
			fn(worker)
		}(w)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

// TestCachedEntry_ConcurrentAccessIsRaceFree drives one shared entry with the operations that
// are legal from every state reachable through them (Expire from REFRESH_NEEDED / ERROR, Error
// from anywhere) plus every read accessor, from several goroutines at once.
func TestCachedEntry_ConcurrentAccessIsRaceFree(t *testing.T) {
	e := NewCachedEntry[*fakeResult]()

	runConcurrently(func(worker int) {
		for i := 0; i < concurrencyRounds; i++ {
			switch (worker + i) % 4 {
			case 0:
				e.Expire()
			case 1:
				// Alternating messages: some calls store a new error (state + exception +
				// result written), some only move the last-occurrence date of the stored one.
				e.Error(NewCachedExceptionWrapper(testError(fmt.Sprintf("failure %d", i%3))))
			case 2:
				e.SyncUpdateDate()
			default:
				// Pure readers, all of them.
				_ = e.CurrentState()
				_ = e.LastStateTransitionTime()
				_ = e.LastSuccessSynchronizationTime()
				_ = e.CachedResult()
				_ = e.IsRefreshNeeded()
				_ = e.IsDesync()
				_ = e.IsToBeDeleted()
				_ = e.IsError()
				_ = e.IsEmpty()
				_ = e.IsResultExist()
				_ = e.ExceptionMessage()
				_ = e.ExceptionStackTrace()
				_ = e.ExceptionFirstOccurrenceTime()
				_ = e.ExceptionLastOccurrenceTime()
			}
		}
	})

	if s := e.CurrentState(); s != CacheStateEnumRefreshNeeded && s != CacheStateEnumError {
		t.Fatalf("unexpected final state %s", s)
	}
	if e.IsError() != (e.ExceptionMessage() != "") {
		t.Fatalf("error state and stored exception disagree: state %s, message %q", e.CurrentState(), e.ExceptionMessage())
	}
}

// TestAbstractCache_ConcurrentSharedEntryAccess is the shape of the pivot fan-out: several
// goroutines reach the same keys through AbstractCache (Get under the map mutex, then the
// entry), some expiring, some recording errors, some reading. Each round re-arms the shared
// entries (SYNCHRONIZED, so that Expire really writes the state) before releasing the workers.
func TestAbstractCache_ConcurrentSharedEntryAccess(t *testing.T) {
	c := NewValidationCache()
	keys := []CacheKey{
		NewCacheKey("https://example.org/lotl.xml"),
		NewCacheKey("https://example.org/pivot-1.xml"),
		NewCacheKey("https://example.org/pivot-2.xml"),
	}

	for round := 0; round < concurrencyRounds; round++ {
		for _, key := range keys {
			// Arm: whatever state a previous round left, end SYNCHRONIZED with a result.
			c.Remove(key)
			c.Update(key, &fakeValidationResult{})
			c.Sync(key)
		}

		runConcurrently(func(worker int) {
			for _, key := range keys {
				switch worker % 3 {
				case 0:
					// What every pivot goroutine does to the shared LOTL / preceding pivots.
					c.Expire(key)
				case 1:
					_ = c.IsRefreshNeeded(key)
					_ = c.IsEmpty(key)
					_ = c.Get(key).CurrentState()
				default:
					_ = c.Keys()
					_ = c.Dump()
				}
			}
		})

		for _, key := range keys {
			if !c.IsRefreshNeeded(key) {
				t.Fatalf("round %d: %v should have been expired by the workers, got %s", round, key.Key(), c.Get(key).CurrentState())
			}
		}
	}
}

// fakeValidationResult is a minimal ValidationResult for the shared-entry tests.
type fakeValidationResult struct{}

func (*fakeValidationResult) Indication() enumerations.Indication {
	return enumerations.IndicationTotalPassed
}
func (*fakeValidationResult) SubIndication() enumerations.SubIndication { return "" }
func (*fakeValidationResult) SigningTime() time.Time                    { return time.Time{} }
func (*fakeValidationResult) SigningCertificate() *model.CertificateToken {
	return nil
}
func (*fakeValidationResult) PotentialSigners() []*model.CertificateToken { return nil }
