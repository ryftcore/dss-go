// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/LOTLWithPivotsAnalysis.java (DSS 6.5.RC1).
//
// GOROUTINE MAPPING: downloadAndParseAllPivots' Java ExecutorService.newFixedThreadPool(n) +
// Future<PivotProcessingResult> per pivot becomes one goroutine per pivot needing a refresh, each
// writing its PivotProcessingResult into a shared map guarded by a sync.Mutex (the map/Future
// collection Java performs after every submission completes); a sync.WaitGroup replaces the
// explicit shutdown()/awaitTermination() shutdown sequence (Go has no thread pool to shut down -
// each goroutine simply returns when done, so there is nothing to await beyond the WaitGroup).
// Iteration order over pivotURLs (sequential, for the "needs update?" decision and the reversed
// walk that actually applies results) is unchanged and remains deterministic; only the
// in-flight download/parse work is concurrent, exactly as in Java.
package tsl

import (
	"sync"

	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/job"
)

// LOTLWithPivotsAnalysis runs the job for a LOTL with pivots analysis.
type LOTLWithPivotsAnalysis struct {
	LOTLAnalysis

	// cacheAccessFactory loads a relevant cache access object.
	cacheAccessFactory *TLCacheAccessFactory

	// dssFileLoader is the file loader.
	dssFileLoader http.DSSFileLoader
}

var _ job.Runnable = (*LOTLWithPivotsAnalysis)(nil)
var _ job.AbstractRunnableAnalysisOverrides = (*LOTLWithPivotsAnalysis)(nil)

// NewLOTLWithPivotsAnalysis is the default constructor.
func NewLOTLWithPivotsAnalysis(source *LOTLSource, cacheAccess job.CacheAccessByKey, dssFileLoader http.DSSFileLoader,
	cacheAccessFactory *TLCacheAccessFactory, latch *sync.WaitGroup) *LOTLWithPivotsAnalysis {
	a := &LOTLWithPivotsAnalysis{
		cacheAccessFactory: cacheAccessFactory,
		dssFileLoader:      dssFileLoader,
	}
	a.LOTLAnalysis = *NewLOTLAnalysis(source, cacheAccess, dssFileLoader, latch)
	a.InitAbstractRunnableAnalysis(a)
	return a
}

// GetCurrentCertificateSource ports the protected getCurrentCertificateSource() override.
func (a *LOTLWithPivotsAnalysis) GetCurrentCertificateSource() spi.CertificateSource {
	initialCertificateSource := a.GetCurrentCertificateSourceDefault()

	parsingCacheEntry, _ := a.CacheAccessByKey().GetParsingReadOnlyResult().(tslmodel.TLParsingInfoRecord)
	if parsingCacheEntry != nil && parsingCacheEntry.IsResultExist() {
		pivotURLs := parsingCacheEntry.PivotUrls()
		if utils.IsCollectionEmpty(pivotURLs) {
			return initialCertificateSource
		}
		return a.getCurrentCertificateSourceFromPivots(initialCertificateSource, pivotURLs)
	}
	return initialCertificateSource
}

func (a *LOTLWithPivotsAnalysis) getCurrentCertificateSourceFromPivots(initialCertificateSource spi.CertificateSource, pivotURLs []string) spi.CertificateSource {
	/*-
	 * current 																						-> Signed with pivot 226 certificates
	 * https://ec.europa.eu/information_society/policy/esignature/trusted-list/tl-pivot-226-mp.xml	-> Signed with pivot 191 certificates
	 * https://ec.europa.eu/information_society/policy/esignature/trusted-list/tl-pivot-191-mp.xml	-> Signed with pivot 172 certificates
	 * https://ec.europa.eu/information_society/policy/esignature/trusted-list/tl-pivot-172-mp.xml 	-> Signed with OJ Certs
	 * http://eur-lex.europa.eu/legal-content/EN/TXT/?uri=uriserv:OJ.C_.2016.233.01.0001.01.ENG		-> OJ
	 */

	processingResults := a.downloadAndParseAllPivots(pivotURLs)

	readOnlyCacheAccess := a.cacheAccessFactory.ReadOnlyCacheAccess()

	pivotUrlsReversed := utils.ReverseList(pivotURLs) // -> 172, 191,..

	currentCertificateSource := initialCertificateSource
	for _, pivotUrl := range pivotUrlsReversed {
		pivotCacheKey := job.NewCacheKey(pivotUrl)

		pivotProcessingResult := processingResults[pivotUrl]
		if pivotProcessingResult != nil {
			pivotCacheAccess := a.cacheAccessFactory.CacheAccess(pivotCacheKey)
			a.validationPivot(pivotCacheAccess, pivotProcessingResult.Pivot(), currentCertificateSource)

			validationResult := readOnlyCacheAccess.GetValidationInfoRecordTyped(pivotCacheKey)
			if validationResult != nil && validationResult.IsResultExist() {
				if validationResult.IsValid() {
					currentCertificateSource = pivotProcessingResult.CertificateSource()
				}
			}
		}
	}

	return currentCertificateSource
}

func (a *LOTLWithPivotsAnalysis) validationPivot(pivotCacheAccess *TLCacheAccessByKey, document model.DSSDocument, certificateSource spi.CertificateSource) {
	// True if EMPTY / EXPIRED by TL/LOTL
	if pivotCacheAccess.IsValidationRefreshNeeded() {
		validationTask := NewTLValidatorTask(document, certificateSource)
		result, err := validationTask.Get()
		if err != nil {
			a.assertOriginalDocumentIsAccessible(pivotCacheAccess)
			pivotCacheAccess.ValidationError(err)
			return
		}
		pivotCacheAccess.UpdateValidationResult(result)
	}
}

func (a *LOTLWithPivotsAnalysis) assertOriginalDocumentIsAccessible(pivotCacheAccess *TLCacheAccessByKey) {
	// set the exception in order to avoid potential deadlock (file does not exist, but download result is present)
	downloadResult := pivotCacheAccess.GetDownloadReadOnlyResult()
	if downloadResult == nil || !downloadResult.IsResultExist() {
		return
	}
	isEmpty, err := spi.DSSUtilsIsEmpty(downloadResult.Document())
	if err != nil || !isEmpty {
		return
	}
	dssErr := model.NewDSSError("Empty content file is obtained!")
	pivotCacheAccess.DownloadError(dssErr)
	pivotCacheAccess.ParsingError(dssErr)
}

func (a *LOTLWithPivotsAnalysis) downloadAndParseAllPivots(pivotURLs []string) map[string]*PivotProcessingResult {
	processingResults := make(map[string]*PivotProcessingResult)

	lotlSource := a.source
	lotlCacheAccessByKey := a.CacheAccessByKey().(*TLCacheAccessByKey)
	pivotProcessingMap := make(map[string]*PivotProcessing)
	var pivotCacheAccessByKeyList []*TLCacheAccessByKey
	for _, pivotUrl := range pivotURLs {
		pivotCacheAccess := a.cacheAccessFactory.CacheAccess(job.NewCacheKey(pivotUrl))

		if lotlCacheAccessByKey.IsValidationRefreshNeeded() || pivotCacheAccess.IsValidationRefreshNeeded() ||
			!downloadResultExists(pivotCacheAccess) {
			pivotSource := NewLOTLSource()
			pivotSource.SetUrl(pivotUrl)
			pivotSource.SetLotlPredicate(lotlSource.LotlPredicate())
			pivotSource.SetTlPredicate(lotlSource.TlPredicate())
			pivotSource.SetPivotSupport(lotlSource.IsPivotSupport())

			// .sha2 is not supported by pivot
			dataLoader := a.dssFileLoader
			if sha2Loader, ok := a.dssFileLoader.(*Sha2FileCacheDataLoader); ok {
				dataLoader = sha2Loader.DataLoader()
			}
			pivotProcessingMap[pivotUrl] = NewPivotProcessing(pivotSource, pivotCacheAccess, lotlCacheAccessByKey,
				append([]*TLCacheAccessByKey(nil), pivotCacheAccessByKeyList...), dataLoader)

		} else {
			// if exists and no update is required
			processingResults[pivotUrl] = NewPivotProcessingResultFromCacheAccessBuilder(pivotCacheAccess).Build()
		}

		pivotCacheAccessByKeyList = append(pivotCacheAccessByKeyList, pivotCacheAccess)
	}

	if len(pivotProcessingMap) > 0 {
		var mu sync.Mutex
		var wg sync.WaitGroup
		for url, processing := range pivotProcessingMap {
			wg.Add(1)
			go func(url string, processing *PivotProcessing) {
				defer wg.Done()
				result, err := processing.Call()
				if err != nil {
					return
				}
				mu.Lock()
				processingResults[url] = result
				mu.Unlock()
			}(url, processing)
		}
		wg.Wait()
	}

	return processingResults
}

// downloadResultExists reports whether pivotCacheAccess's download read-only result exists,
// mirroring `!pivotCacheAccess.getDownloadReadOnlyResult().isResultExist()`.
func downloadResultExists(pivotCacheAccess *TLCacheAccessByKey) bool {
	downloadResult := pivotCacheAccess.GetDownloadReadOnlyResult()
	return downloadResult != nil && downloadResult.IsResultExist()
}
