// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/CompositeTSPSource.java (DSS 6.5.RC1).
package validation

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
)

// CompositeTSPSource allows retrieving a timestamp with different sources. The composite tries
// all sources until it gets a non-empty response.
//
// Be careful, all given tspSources MUST accept the same digest algorithm.
type CompositeTSPSource struct {
	// tspSources maps source keys to their corresponding TSP sources.
	tspSources map[string]TSPSource

	// sourceOrder preserves insertion order, standing in for Java's LinkedHashMap-like
	// iteration order guarantee that Map.entrySet() otherwise would not provide; SetTspSources
	// is the only mutator and replaces the whole map, so recomputing the order there is enough.
	// Same treatment as CompositeRevocationSource in the spi package.
	sourceOrder []string
}

// NewCompositeTSPSource instantiates the object with nil values. Port of the default constructor.
func NewCompositeTSPSource() *CompositeTSPSource {
	return &CompositeTSPSource{}
}

// SetTspSources allows providing multiple tspSources, keyed by a label. Be careful, all given
// tspSources MUST accept the same digest algorithm.
// Port of setTspSources(Map<String, TSPSource>).
//
// Go maps have no defined iteration order (Java's HashMap does not either, but callers in
// practice pass a LinkedHashMap); sourceKeys lets callers pin down the try order the way a
// LinkedHashMap argument would. When sourceKeys is nil the iteration order is unspecified,
// matching Java's setTspSources(Map) as declared.
func (s *CompositeTSPSource) SetTspSources(tspSources map[string]TSPSource, sourceKeys []string) {
	s.tspSources = tspSources
	if sourceKeys != nil {
		s.sourceOrder = sourceKeys
	} else {
		s.sourceOrder = nil
		for sourceKey := range tspSources {
			s.sourceOrder = append(s.sourceOrder, sourceKey)
		}
	}
}

// TimeStampResponse tries all sources until one returns a timestamp, and fails with a
// DSSExternalResourceException when none does.
// Port of getTimeStampResponse(DigestAlgorithm, byte[]).
//
// Java catches and logs any Exception a source raises and keeps trying the remaining ones; the
// Go port skips a source that returns an error (the equivalent of that catch) and lets a source
// that panics propagate, exactly as CompositeRevocationSource does.
func (s *CompositeTSPSource) TimeStampResponse(digestAlgorithm enumerations.DigestAlgorithm,
	digestValue []byte) (*model.TimestampBinary, error) {
	for _, sourceKey := range s.sourceOrder {
		source, found := s.tspSources[sourceKey]
		if !found {
			continue
		}
		// Upstream logs "Trying to get timestamp with TSPSource '{}'" at debug level.
		timestampBinary, err := source.TimeStampResponse(digestAlgorithm, digestValue)
		if err != nil {
			// Upstream logs "Unable to retrieve the timestamp with TSPSource '{}' : {}".
			continue
		}
		if timestampBinary != nil {
			// Upstream logs "Successfully retrieved timestamp with TSPSource '{}'".
			return timestampBinary, nil
		}
	}
	return nil, exception.NewDSSExternalResourceException(
		fmt.Sprintf("Unable to retrieve the timestamp (%d tries)", len(s.tspSources)))
}

// compile-time assertion: a CompositeTSPSource is a TSPSource.
var _ TSPSource = (*CompositeTSPSource)(nil)
