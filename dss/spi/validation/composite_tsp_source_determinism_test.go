package validation

import (
	"errors"
	"reflect"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// compositeTSPSourceDeterminismFakeSource records that it was tried (into a shared, ordered
// log) and always errors, so CompositeTSPSource.TimeStampResponse tries every configured
// source and the log records the full try order.
type compositeTSPSourceDeterminismFakeSource struct {
	name string
	log  *[]string
}

func (f *compositeTSPSourceDeterminismFakeSource) TimeStampResponse(enumerations.DigestAlgorithm, []byte) (*model.TimestampBinary, error) {
	*f.log = append(*f.log, f.name)
	return nil, errors.New("no timestamp available")
}

var _ TSPSource = (*compositeTSPSourceDeterminismFakeSource)(nil)

// TestCompositeTSPSourceNilSourceKeysDeterministic guards against defect #2 of the Phase 2b
// audit ("Nondeterministic output ordering"): when SetTspSources is called with a nil
// sourceKeys (no explicit try order requested), the fallback used to range directly over the
// tspSources map - randomized by Go on every run, unlike Java's HashMap (arbitrary but stable
// within a JVM run) - so the fallback now sorts lexically instead.
func TestCompositeTSPSourceNilSourceKeysDeterministic(t *testing.T) {
	build := func() []string {
		var log []string
		source := NewCompositeTSPSource()
		source.SetTspSources(map[string]TSPSource{
			"delta":   &compositeTSPSourceDeterminismFakeSource{name: "delta", log: &log},
			"alpha":   &compositeTSPSourceDeterminismFakeSource{name: "alpha", log: &log},
			"charlie": &compositeTSPSourceDeterminismFakeSource{name: "charlie", log: &log},
			"echo":    &compositeTSPSourceDeterminismFakeSource{name: "echo", log: &log},
			"bravo":   &compositeTSPSourceDeterminismFakeSource{name: "bravo", log: &log},
		}, nil)
		_, _ = source.TimeStampResponse(enumerations.DigestAlgorithm_SHA256, []byte("digest"))
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
