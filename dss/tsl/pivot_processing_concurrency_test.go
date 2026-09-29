package tsl

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// Regression tests for the concurrent pivot fan-out of LOTLWithPivotsAnalysis
// (downloadAndParseAllPivots runs one goroutine per pivot; every one of them expires the SAME
// LOTL validation entry and every preceding pivot's validation entry):
//
//   - T18A-SEC-001 / T31-SEC-001: the shared validation entries were mutated by several
//     goroutines with no synchronization - a data race, reported by `go test -race`.
//   - T18A-SEC-002: a pivot whose document lacks <SchemeInformation> makes LOTLParsingTask
//     dereference nil (upstream's NullPointerException); on the pivot goroutine that unrecovered
//     panic used to kill the whole process, where upstream records a parsing error on the
//     pivot's cache entry (AbstractAnalysis.parsing's catch) and carries on.

const (
	pivotFanOutLOTLURL = "https://example.org/lotl.xml"

	// pivotWithoutSchemeInformation is a well-formed XML document with a TrustServiceStatusList
	// root and NO SchemeInformation child: it downloads and parses as a TL document, and then
	// LOTLParsingTask dereferences the missing SchemeInformation.
	pivotWithoutSchemeInformation = `<?xml version="1.0" encoding="UTF-8"?>` +
		`<TrustServiceStatusList xmlns="http://uri.etsi.org/02231/v2#" TSLTag="http://uri.etsi.org/19612/TSLTag"/>`

	// pivotWithEmptySchemeInformation carries an (empty) SchemeInformation, so it parses without
	// raising, and - having no pointer to another TSL - yields no PivotProcessingResult.
	pivotWithEmptySchemeInformation = `<?xml version="1.0" encoding="UTF-8"?>` +
		`<TrustServiceStatusList xmlns="http://uri.etsi.org/02231/v2#" TSLTag="http://uri.etsi.org/19612/TSLTag">` +
		`<SchemeInformation/></TrustServiceStatusList>`
)

// pivotFanOutLoader is an offline file loader answering the same XML for every pivot URL.
type pivotFanOutLoader struct {
	xml string
}

func (l pivotFanOutLoader) GetDocument(url string) (model.DSSDocument, error) {
	return model.NewInMemoryDocument([]byte(l.xml)), nil
}

// pivotFanOutValidationResult is a minimal job.ValidationResult to arm validation entries with.
type pivotFanOutValidationResult struct{}

func (pivotFanOutValidationResult) Indication() enumerations.Indication {
	return enumerations.IndicationTotalPassed
}
func (pivotFanOutValidationResult) SubIndication() enumerations.SubIndication { return "" }
func (pivotFanOutValidationResult) SigningTime() time.Time                    { return time.Time{} }
func (pivotFanOutValidationResult) SigningCertificate() *model.CertificateToken {
	return nil
}
func (pivotFanOutValidationResult) PotentialSigners() []*model.CertificateToken { return nil }

// newPivotFanOut builds the LOTL analysis over a fresh cache factory and returns it together with
// the pivot URLs. When arm is true, the LOTL's and every pivot's validation entry is put in the
// SYNCHRONIZED state (as after a previous refresh), so that expiring them really writes state.
func newPivotFanOut(t *testing.T, pivotXML string, pivots int, arm bool) (*LOTLWithPivotsAnalysis, *TLCacheAccessFactory, []string) {
	t.Helper()
	factory := NewTLCacheAccessFactory()

	lotlSource := NewLOTLSource()
	lotlSource.SetUrl(pivotFanOutLOTLURL)
	lotlSource.SetPivotSupport(true)

	pivotURLs := make([]string, pivots)
	for i := range pivotURLs {
		pivotURLs[i] = fmt.Sprintf("https://example.org/pivot-%d.xml", i)
	}

	if arm {
		validationCache := factory.ValidationCache()
		for _, url := range append([]string{pivotFanOutLOTLURL}, pivotURLs...) {
			key := job.NewCacheKey(url)
			validationCache.Update(key, pivotFanOutValidationResult{})
			validationCache.Sync(key)
		}
	}

	var latch sync.WaitGroup
	analysis := NewLOTLWithPivotsAnalysis(lotlSource, factory.CacheAccess(job.NewCacheKey(pivotFanOutLOTLURL)),
		pivotFanOutLoader{xml: pivotXML}, factory, &latch)
	return analysis, factory, pivotURLs
}

// TestDownloadAndParseAllPivots_ConcurrentExpireOfSharedEntries drives the fan-out with several
// pivots over armed (SYNCHRONIZED) validation entries. Every pivot goroutine downloads a new
// document and so expires the shared LOTL entry and the preceding pivots' entries: with the
// race detector on, unsynchronized entry state is reported here.
func TestDownloadAndParseAllPivots_ConcurrentExpireOfSharedEntries(t *testing.T) {
	const rounds = 25
	for round := 0; round < rounds; round++ {
		analysis, factory, pivotURLs := newPivotFanOut(t, pivotWithEmptySchemeInformation, 6, true)

		results := analysis.downloadAndParseAllPivots(pivotURLs)

		// Each pivot parses (no pointer to another TSL, hence no result) and expired the
		// shared entries.
		for _, url := range pivotURLs {
			if results[url] != nil {
				t.Fatalf("round %d: pivot %s unexpectedly produced a processing result", round, url)
			}
			if factory.ParsingCache().Get(job.NewCacheKey(url)).IsError() {
				t.Fatalf("round %d: pivot %s: unexpected parsing error: %s", round, url,
					factory.ParsingCache().Get(job.NewCacheKey(url)).ExceptionMessage())
			}
		}
		if !factory.ValidationCache().IsRefreshNeeded(job.NewCacheKey(pivotFanOutLOTLURL)) {
			t.Fatalf("round %d: the LOTL validation entry should have been expired by the pivots", round)
		}
		for _, url := range pivotURLs[:len(pivotURLs)-1] {
			if !factory.ValidationCache().IsRefreshNeeded(job.NewCacheKey(url)) {
				t.Fatalf("round %d: the validation entry of the preceding pivot %s should have been expired", round, url)
			}
		}
	}
}

// TestDownloadAndParseAllPivots_PivotWithoutSchemeInformationDoesNotCrash feeds every pivot a
// document that makes LOTLParsingTask dereference nil. The process must survive, each pivot must
// end in a recorded PARSING error (upstream's AbstractAnalysis.parsing catch), and no pivot may
// yield a processing result.
func TestDownloadAndParseAllPivots_PivotWithoutSchemeInformationDoesNotCrash(t *testing.T) {
	analysis, factory, pivotURLs := newPivotFanOut(t, pivotWithoutSchemeInformation, 3, false)

	results := analysis.downloadAndParseAllPivots(pivotURLs)

	for _, url := range pivotURLs {
		key := job.NewCacheKey(url)
		if results[url] != nil {
			t.Errorf("pivot %s: a pivot that cannot be parsed must not yield a processing result", url)
		}
		if !factory.DownloadCache().Get(key).IsResultExist() {
			t.Errorf("pivot %s: the download itself succeeded and must stay cached", url)
		}
		if !factory.ParsingCache().Get(key).IsError() {
			t.Errorf("pivot %s: expected the parsing failure to be recorded on the parsing cache entry, got state %s",
				url, factory.ParsingCache().Get(key).CurrentState())
		}
	}
}

// TestPivotProcessingCall_PanicIsRecordedAsParsingError pins the same behaviour on the single
// pivot: Call() returns normally, with no result, and the pivot's parsing entry is in ERROR.
func TestPivotProcessingCall_PanicIsRecordedAsParsingError(t *testing.T) {
	factory := NewTLCacheAccessFactory()
	pivotSource := NewLOTLSource()
	pivotSource.SetUrl("https://example.org/pivot.xml")
	pivotKey := job.NewCacheKey("https://example.org/pivot.xml")
	lotlKey := job.NewCacheKey(pivotFanOutLOTLURL)

	processing := NewPivotProcessing(pivotSource, factory.CacheAccess(pivotKey), factory.CacheAccess(lotlKey), nil,
		pivotFanOutLoader{xml: pivotWithoutSchemeInformation})

	result, err := processing.Call()
	if err != nil || result != nil {
		t.Fatalf("Call() = (%v, %v), want (nil, nil)", result, err)
	}
	entry := factory.ParsingCache().Get(pivotKey)
	if !entry.IsError() {
		t.Fatalf("expected a parsing error on the pivot, got state %s", entry.CurrentState())
	}
	if entry.ExceptionMessage() == "" {
		t.Error("the recorded parsing error must carry the panic's message")
	}
}

// TestLOTLAnalysisRun_ParsingPanicIsRecordedAsParsingError covers the main analysis path, which
// shares AbstractAnalysis.parsing with upstream: the nil dereference LOTLParsingTask reproduces
// from upstream's NullPointerException on a LOTL lacking <SchemeInformation> is caught there and
// recorded as a PARSING error on the cache entry (upstream: catch (Exception e) ->
// cacheAccess.parsingError(e)). Before, Run()'s blanket recover swallowed the panic without
// recording anything, so the entry stayed REFRESH_NEEDED and the parsing-error alert could never
// fire for such a document.
func TestLOTLAnalysisRun_ParsingPanicIsRecordedAsParsingError(t *testing.T) {
	factory := NewTLCacheAccessFactory()
	lotlSource := NewLOTLSource()
	lotlSource.SetUrl(pivotFanOutLOTLURL)
	certificateSource := spi.NewCommonCertificateSource()
	lotlSource.SetCertificateSource(&certificateSource)
	lotlKey := job.NewCacheKey(pivotFanOutLOTLURL)

	var latch sync.WaitGroup
	latch.Add(1)
	analysis := NewLOTLAnalysis(lotlSource, factory.CacheAccess(lotlKey),
		pivotFanOutLoader{xml: pivotWithoutSchemeInformation}, &latch)
	analysis.Run()
	latch.Wait()

	if !factory.DownloadCache().Get(lotlKey).IsResultExist() {
		t.Error("the LOTL download itself succeeded and must be cached")
	}
	entry := factory.ParsingCache().Get(lotlKey)
	if !entry.IsError() {
		t.Fatalf("expected the parsing failure to be recorded, got parsing state %s", entry.CurrentState())
	}
	if entry.ExceptionMessage() == "" {
		t.Error("the recorded parsing error must carry the panic's message")
	}
}
