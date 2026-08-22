package aia

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// compositeAIASourceDeterminismFakeSource records that it was tried (into a shared, ordered
// log) and always reports "no data" (nil), so CompositeAIASource.CertificatesByAIA tries every
// configured source and the log records the full try order.
type compositeAIASourceDeterminismFakeSource struct {
	name string
	log  *[]string
}

func (f *compositeAIASourceDeterminismFakeSource) CertificatesByAIA(_ *model.CertificateToken) []*model.CertificateToken {
	*f.log = append(*f.log, f.name)
	return nil
}

var _ AIASource = (*compositeAIASourceDeterminismFakeSource)(nil)

// TestCompositeAIASourceOrderedKeysDeterministic verifies that compositeAIASourceOrderedKeys'
// try order is stable across runs: it used to range directly over the aiaSources map -
// randomized by Go on every run, unlike Java's HashMap (arbitrary but stable within a JVM run)
// - so the try order is now sorted lexically instead.
func TestCompositeAIASourceOrderedKeysDeterministic(t *testing.T) {
	build := func() []string {
		var log []string
		source := NewCompositeAIASource()
		source.SetAIASources(map[string]AIASource{
			"delta":   &compositeAIASourceDeterminismFakeSource{name: "delta", log: &log},
			"alpha":   &compositeAIASourceDeterminismFakeSource{name: "alpha", log: &log},
			"charlie": &compositeAIASourceDeterminismFakeSource{name: "charlie", log: &log},
			"echo":    &compositeAIASourceDeterminismFakeSource{name: "echo", log: &log},
			"bravo":   &compositeAIASourceDeterminismFakeSource{name: "bravo", log: &log},
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
