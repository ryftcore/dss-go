package spi

import (
	"reflect"
	"testing"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// compositeRevocationSourceDeterminismFakeRevocation is a stand-in for a concrete Revocation
// (revocation.Revocation is a marker interface with no methods, satisfied by any type).
type compositeRevocationSourceDeterminismFakeRevocation struct{}

// compositeRevocationSourceDeterminismFakeSource records that it was tried (into a shared,
// ordered log) and always reports "no data", so CompositeRevocationSource.RevocationToken
// tries every configured source and the log records the full try order.
type compositeRevocationSourceDeterminismFakeSource struct {
	name string
	log  *[]string
}

func (f *compositeRevocationSourceDeterminismFakeSource) RevocationToken(_, _ *model.CertificateToken) RevocationToken[compositeRevocationSourceDeterminismFakeRevocation] {
	*f.log = append(*f.log, f.name)
	return nil
}

var _ RevocationSource[compositeRevocationSourceDeterminismFakeRevocation] = (*compositeRevocationSourceDeterminismFakeSource)(nil)

// TestCompositeRevocationSourceNilSourceKeysDeterministic guards against defect #2 of the
// Phase 2b audit ("Nondeterministic output ordering"): when SetSources is called with a nil
// sourceKeys (no explicit try order requested), the fallback used to range directly over the
// sources map, which - unlike Java's HashMap (arbitrary but stable within a JVM run) - is
// randomized by Go on every run. The fallback now sorts lexically, so the try order (and, for
// sources that could each answer differently, the winning result) is stable across runs.
func TestCompositeRevocationSourceNilSourceKeysDeterministic(t *testing.T) {
	build := func() []string {
		var log []string
		source := NewCompositeRevocationSource[compositeRevocationSourceDeterminismFakeRevocation]()
		sources := map[string]RevocationSource[compositeRevocationSourceDeterminismFakeRevocation]{
			"delta":   &compositeRevocationSourceDeterminismFakeSource{name: "delta", log: &log},
			"alpha":   &compositeRevocationSourceDeterminismFakeSource{name: "alpha", log: &log},
			"charlie": &compositeRevocationSourceDeterminismFakeSource{name: "charlie", log: &log},
			"echo":    &compositeRevocationSourceDeterminismFakeSource{name: "echo", log: &log},
			"bravo":   &compositeRevocationSourceDeterminismFakeSource{name: "bravo", log: &log},
		}
		source.SetSources(sources, nil)
		source.RevocationToken(nil, nil)
		return log
	}
	want := build()
	if len(want) != 5 {
		t.Fatalf("got %d tries, want 5", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: try order = %v, want %v", i, got, want)
		}
	}
}

var _ = revocation.Revocation(compositeRevocationSourceDeterminismFakeRevocation{})
