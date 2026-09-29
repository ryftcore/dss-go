package xades

import (
	"runtime"
	"sync"
	"weak"
)

// weakRegistry recovers the enclosing value of an embedded base struct without keeping either
// alive. It backs the "no instanceof" workaround the XAdES port needs wherever Java downcasts:
// ReferenceValidation and SignaturePolicy embed model types by value, callers hand the base
// pointer around, and Java's `(ReferenceValidation) referenceValidation` has no Go equivalent
// through embedding.
//
// The registry is keyed by a weak pointer to the embedded base and holds a weak pointer to the
// enclosing value, and a cleanup attached to the enclosing value removes the entry once it is
// unreachable. A strong map (the sync.Map these replaced) pinned every *xmldsig.Reference and
// DOM subtree of every signature ever validated in the process for the life of the process -
// unbounded growth on a long-lived validation server, where Java's garbage collector would have
// reclaimed them.
//
// The base is embedded by value, so a pointer to it keeps the enclosing allocation alive: a
// caller that holds only the base pointer still resolves to a live enclosing value.
type weakRegistry[B, C any] struct {
	m sync.Map // map[weak.Pointer[B]]weak.Pointer[C]
}

// store records that base is the embedded base of outer. base must point into outer.
func (r *weakRegistry[B, C]) store(base *B, outer *C) {
	key := weak.Make(base)
	r.m.Store(key, weak.Make(outer))
	runtime.AddCleanup(outer, func(k weak.Pointer[B]) { r.m.Delete(k) }, key)
}

// load returns the enclosing value registered for base, if it is still alive.
func (r *weakRegistry[B, C]) load(base *B) (*C, bool) {
	v, ok := r.m.Load(weak.Make(base))
	if !ok {
		return nil, false
	}
	outer := v.(weak.Pointer[C]).Value()
	return outer, outer != nil
}

// size counts the entries, for tests.
func (r *weakRegistry[B, C]) size() int {
	n := 0
	r.m.Range(func(_, _ any) bool { n++; return true })
	return n
}
