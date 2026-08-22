package validation

import (
	"errors"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// compositeTSPSourceTestSource is a TSPSource whose answer the test dictates, recording whether
// it was asked.
type compositeTSPSourceTestSource struct {
	binary *model.TimestampBinary
	err    error
	called bool
}

// TimeStampResponse returns the configured answer.
func (s *compositeTSPSourceTestSource) TimeStampResponse(enumerations.DigestAlgorithm,
	[]byte) (*model.TimestampBinary, error) {
	s.called = true
	return s.binary, s.err
}

// TestCompositeTSPSourceTriesSourcesInOrder checks that the composite stops at the first source
// that answers and skips the ones that fail or answer nothing, in the order the caller pinned.
func TestCompositeTSPSourceTriesSourcesInOrder(t *testing.T) {
	failing := &compositeTSPSourceTestSource{err: errors.New("network down")}
	empty := &compositeTSPSourceTestSource{}
	answering := &compositeTSPSourceTestSource{binary: model.NewTimestampBinary([]byte{0x01})}
	unused := &compositeTSPSourceTestSource{binary: model.NewTimestampBinary([]byte{0x02})}

	composite := NewCompositeTSPSource()
	composite.SetTspSources(map[string]TSPSource{
		"failing": failing, "empty": empty, "answering": answering, "unused": unused,
	}, []string{"failing", "empty", "answering", "unused"})

	binary, err := composite.TimeStampResponse(enumerations.DigestAlgorithmSHA256, make([]byte, 32))
	if err != nil {
		t.Fatalf("TimeStampResponse() failed: %v", err)
	}
	if len(binary.Bytes()) != 1 || binary.Bytes()[0] != 0x01 {
		t.Errorf("TimeStampResponse() = %x, want the third source's answer", binary.Bytes())
	}
	if !failing.called || !empty.called || !answering.called {
		t.Error("the composite did not try every source up to the answering one")
	}
	if unused.called {
		t.Error("the composite kept trying after a source answered")
	}
}

// TestCompositeTSPSourceExhausted checks the DSSExternalResourceException raised when no source
// answers, message included.
func TestCompositeTSPSourceExhausted(t *testing.T) {
	composite := NewCompositeTSPSource()
	composite.SetTspSources(map[string]TSPSource{
		"first":  &compositeTSPSourceTestSource{err: errors.New("boom")},
		"second": &compositeTSPSourceTestSource{},
	}, []string{"first", "second"})

	_, err := composite.TimeStampResponse(enumerations.DigestAlgorithmSHA256, make([]byte, 32))
	if err == nil {
		t.Fatal("TimeStampResponse() succeeded, want an error")
	}
	var externalResource *exception.DSSExternalResourceException
	if !errors.As(err, &externalResource) {
		t.Errorf("error is %T, want a DSSExternalResourceException", err)
	}
	if got, want := err.Error(), "Unable to retrieve the timestamp (2 tries)"; got != want {
		t.Errorf("error = %q, want %q", got, want)
	}
}

// TestCompositeTSPSourceUnorderedSources checks that an unordered map is still fully tried, which
// is the behaviour Java's setTspSources(Map) declares.
func TestCompositeTSPSourceUnorderedSources(t *testing.T) {
	only := &compositeTSPSourceTestSource{binary: model.NewTimestampBinary([]byte{0x03})}
	composite := NewCompositeTSPSource()
	composite.SetTspSources(map[string]TSPSource{"only": only}, nil)

	binary, err := composite.TimeStampResponse(enumerations.DigestAlgorithmSHA512, make([]byte, 64))
	if err != nil {
		t.Fatalf("TimeStampResponse() failed: %v", err)
	}
	if len(binary.Bytes()) != 1 || binary.Bytes()[0] != 0x03 {
		t.Errorf("TimeStampResponse() = %x, want the only source's answer", binary.Bytes())
	}
}
