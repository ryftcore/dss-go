package job

import "testing"

// TestCacheStateEnum_TransitionTable verifies every (state, transition) pair against the
// Java oracle table encoded in CacheStateEnum.java's per-constant method overrides (see the
// file header comment in cache_state_enum.go for the derivation).
func TestCacheStateEnum_TransitionTable(t *testing.T) {
	allStates := []CacheStateEnum{
		CacheStateEnum_REFRESH_NEEDED, CacheStateEnum_DESYNCHRONIZED, CacheStateEnum_SYNCHRONIZED,
		CacheStateEnum_ERROR, CacheStateEnum_TO_BE_DELETED,
	}

	// wantPanic[method] is the set of starting states for which the transition panics.
	wantPanic := map[string]map[CacheStateEnum]bool{
		"sync":          {CacheStateEnum_REFRESH_NEEDED: true, CacheStateEnum_ERROR: true, CacheStateEnum_TO_BE_DELETED: true},
		"desync":        {CacheStateEnum_DESYNCHRONIZED: true},
		"refreshNeeded": {CacheStateEnum_DESYNCHRONIZED: true},
		"toBeDeleted":   {CacheStateEnum_DESYNCHRONIZED: true, CacheStateEnum_TO_BE_DELETED: true},
		"error":         {CacheStateEnum_DESYNCHRONIZED: true, CacheStateEnum_SYNCHRONIZED: true, CacheStateEnum_ERROR: true, CacheStateEnum_TO_BE_DELETED: true},
	}

	for _, start := range allStates {
		t.Run(string(start)+"/sync", func(t *testing.T) {
			checkTransition(t, start, wantPanic["sync"][start], CacheStateEnum_SYNCHRONIZED, func(ctx *CurrentCacheContext) { start.Sync(ctx) })
		})
		t.Run(string(start)+"/desync", func(t *testing.T) {
			checkTransition(t, start, wantPanic["desync"][start], CacheStateEnum_DESYNCHRONIZED, func(ctx *CurrentCacheContext) { start.Desync(ctx) })
		})
		t.Run(string(start)+"/refreshNeeded", func(t *testing.T) {
			checkTransition(t, start, wantPanic["refreshNeeded"][start], CacheStateEnum_REFRESH_NEEDED, func(ctx *CurrentCacheContext) { start.RefreshNeeded(ctx) })
		})
		t.Run(string(start)+"/toBeDeleted", func(t *testing.T) {
			checkTransition(t, start, wantPanic["toBeDeleted"][start], CacheStateEnum_TO_BE_DELETED, func(ctx *CurrentCacheContext) { start.ToBeDeleted(ctx) })
		})
		t.Run(string(start)+"/error", func(t *testing.T) {
			ctx := forceCacheState(start)
			defer func() {
				r := recover()
				if wantPanic["error"][start] {
					if r == nil {
						t.Fatalf("state=%s: expected panic on error(), got none", start)
					}
					return
				}
				if r != nil {
					t.Fatalf("state=%s: unexpected panic on error(): %v", start, r)
				}
				if ctx.CurrentState() != CacheStateEnum_ERROR {
					t.Fatalf("state=%s: after error() want ERROR, got %s", start, ctx.CurrentState())
				}
			}()
			start.Error(ctx, NewCachedExceptionWrapper(errTest))
		})
	}
}

func checkTransition(t *testing.T, start CacheStateEnum, wantPanic bool, target CacheStateEnum, call func(ctx *CurrentCacheContext)) {
	t.Helper()
	ctx := forceCacheState(start)
	defer func() {
		r := recover()
		if wantPanic {
			if r == nil {
				t.Fatalf("state=%s: expected panic, got none (ended in state %s)", start, ctx.CurrentState())
			}
			return
		}
		if r != nil {
			t.Fatalf("state=%s: unexpected panic: %v", start, r)
		}
		if ctx.CurrentState() != target {
			t.Fatalf("state=%s: want %s, got %s", start, target, ctx.CurrentState())
		}
	}()
	call(ctx)
}

// forceCacheState builds a CurrentCacheContext directly in the given state, bypassing the
// transition table (which is exactly what is under test).
func forceCacheState(state CacheStateEnum) *CurrentCacheContext {
	ctx := NewCurrentCacheContext()
	ctx.state = state
	return ctx
}

type testError string

func (e testError) Error() string { return string(e) }

const errTest = testError("boom")
