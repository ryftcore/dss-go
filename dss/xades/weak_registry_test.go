package xades

import (
	"runtime"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// eventuallyAtMost collects garbage until size() drops to max or five seconds pass. The
// registry cleanups run on a separate goroutine after a collection, hence the polling.
func eventuallyAtMost(t *testing.T, what string, size func() int, max int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		runtime.GC()
		if size() <= max {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s still holds %d entries, want at most %d: the registry pins unreachable values", what, size(), max)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestSignaturePolicyRegistryDoesNotPinPolicies is the T40-PERF-001 regression: the registry that
// recovers *SignaturePolicy from its embedded base used to be a strong sync.Map with no
// eviction, so every policy read from every validated signature stayed reachable forever.
func TestSignaturePolicyRegistryDoesNotPinPolicies(t *testing.T) {
	keep := NewSignaturePolicyWithIdentifier("urn:keep")
	before := xadesSignaturePolicyRegistry.size()

	for i := 0; i < 2000; i++ {
		NewSignaturePolicy()
	}
	eventuallyAtMost(t, "xadesSignaturePolicyRegistry", xadesSignaturePolicyRegistry.size, before)

	// A caller holding only the base pointer still recovers the concrete value.
	base := &keep.Policy
	keep = nil
	runtime.GC()
	got, ok := SignaturePolicyFor(base)
	if !ok || got == nil || got.Identifier() != "urn:keep" {
		t.Fatalf("SignaturePolicyFor(base) = (%v, %v), want the policy built for urn:keep", got, ok)
	}
	runtime.KeepAlive(base)
}

func TestSignaturePolicyForUnknownAndNil(t *testing.T) {
	if _, ok := SignaturePolicyFor(nil); ok {
		t.Error("SignaturePolicyFor(nil) reported a policy")
	}
	if _, ok := ReferenceValidationFor(nil); ok {
		t.Error("ReferenceValidationFor(nil) reported a validation")
	}
	stranger := NewSignaturePolicy().Policy // a copy: its address was never registered
	if _, ok := SignaturePolicyFor(&stranger); ok {
		t.Error("SignaturePolicyFor(copy) resolved a policy that was never registered at that address")
	}
}

// TestPolicyTransformsRegistryDoesNotPinPolicies covers the third registry, which backs
// SignatureBuilderRegisterPolicyTransforms.
func TestPolicyTransformsRegistryDoesNotPinPolicies(t *testing.T) {
	size := func() int {
		n := 0
		xadesSignatureBuilderPolicyTransformsRegistry.Range(func(_, _ any) bool { n++; return true })
		return n
	}

	keep := NewXmlPolicyWithTransforms()
	keep.SetTransforms([]DSSTransform{nil})
	registered := SignatureBuilderRegisterPolicyTransforms(keep)
	before := size()

	for i := 0; i < 2000; i++ {
		p := NewXmlPolicyWithTransforms()
		p.SetTransforms([]DSSTransform{nil})
		SignatureBuilderRegisterPolicyTransforms(p)
	}
	eventuallyAtMost(t, "xadesSignatureBuilderPolicyTransformsRegistry", size, before)

	transforms, ok := SignatureBuilderPolicyTransforms(registered)
	if !ok || len(transforms) != 1 {
		t.Fatalf("SignatureBuilderPolicyTransforms(registered) = (%v, %v), want the one registered transform", transforms, ok)
	}
	if _, ok := SignatureBuilderPolicyTransforms(model.NewPolicy()); ok {
		t.Error("a plain model.Policy resolved registered transforms")
	}
	runtime.KeepAlive(keep)
}

// TestWeakRegistryResolvesThroughBaseOnly pins the property the registries rely on: the base is
// embedded by value, so holding only its address keeps the enclosing value resolvable.
func TestWeakRegistryResolvesThroughBaseOnly(t *testing.T) {
	type base struct{ n int }
	type outer struct {
		base
		tag string
	}
	var r weakRegistry[base, outer]
	o := &outer{tag: "x"}
	r.store(&o.base, o)
	b := &o.base
	o = nil
	runtime.GC()
	got, ok := r.load(b)
	if !ok || got.tag != "x" {
		t.Fatalf("load through the base pointer = (%v, %v)", got, ok)
	}
	runtime.KeepAlive(b)
}
