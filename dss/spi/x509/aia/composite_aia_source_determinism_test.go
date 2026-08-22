package aia

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// compositeSourceDeterminismFakeSource records that it was tried (into a shared, ordered
// log) and always reports "no data" (nil), so CompositeSource.CertificatesByAIA tries every
// configured source and the log records the full try order.
type compositeSourceDeterminismFakeSource struct {
	name string
	log  *[]string
}

func (f *compositeSourceDeterminismFakeSource) CertificatesByAIA(_ *model.CertificateToken) []*model.CertificateToken {
	*f.log = append(*f.log, f.name)
	return nil
}

var _ Source = (*compositeSourceDeterminismFakeSource)(nil)

// TestCompositeAIASourceOrderedKeysDeterministic verifies that compositeSourceOrderedKeys'
// try order is stable across runs: it used to range directly over the aiaSources map -
// randomized by Go on every run, unlike Java's HashMap (arbitrary but stable within a JVM run)
// - so the try order is now sorted lexically instead.
func TestCompositeAIASourceOrderedKeysDeterministic(t *testing.T) {
	build := func() []string {
		var log []string
		source := NewCompositeSource()
		source.SetAIASources(map[string]Source{
			"delta":   &compositeSourceDeterminismFakeSource{name: "delta", log: &log},
			"alpha":   &compositeSourceDeterminismFakeSource{name: "alpha", log: &log},
			"charlie": &compositeSourceDeterminismFakeSource{name: "charlie", log: &log},
			"echo":    &compositeSourceDeterminismFakeSource{name: "echo", log: &log},
			"bravo":   &compositeSourceDeterminismFakeSource{name: "bravo", log: &log},
		})
		func() {
			defer func() { recover() }() // every source answers nil, so the composite panics after trying all of them
			source.CertificatesByAIA(nil)
		}()
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
