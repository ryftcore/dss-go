// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/ValidationJob.java (DSS 6.5.RC1).
package job

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/model"
	modeljob "github.com/ryftcore/dss-go/dss/model/job"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/utils"
)

// Runnable is the Go stand-in for java.lang.Runnable, used by the (abstract)
// GetDocumentAnalysis/GetDocumentListAnalysis hooks and satisfied by
// *AbstractRunnableAnalysis.
type Runnable interface {
	Run()
}

// ExecutorService is the Go stand-in for the subset of java.util.concurrent.ExecutorService
// used by ValidationJob: submitting Runnable tasks for asynchronous execution, and orderly
// or immediate shutdown. See PORTING.md: "Java concurrency (ExecutorService) → goroutines
// with the same task decomposition."
type ExecutorService interface {
	// Submit schedules task for asynchronous execution. Port of submit(Runnable).
	Submit(task Runnable)
	// ShutdownNow requests an immediate shutdown. Port of shutdownNow().
	ShutdownNow()
	// IsShutdown reports whether shutdown has been requested. Port of isShutdown().
	IsShutdown() bool
}

// cachedThreadPoolExecutorService is a minimal ExecutorService that runs every submitted
// task on its own goroutine, mirroring Java's Executors.newCachedThreadPool() (an unbounded
// thread-per-task pool). Port of the default field initializer
// "Executors.newCachedThreadPool()".
//
// DEVIATION: Java's ExecutorService.shutdownNow() also attempts to interrupt and cancel
// already-running tasks; the Go port only marks the service as shut down (Submit after
// shutdown is still accepted, since ValidationJob never calls Submit after ShutdownNow in
// this module — SetExecutorService replaces the field wholesale).
type cachedThreadPoolExecutorService struct {
	mu       sync.Mutex
	shutdown bool
}

// NewCachedThreadPoolExecutorService creates the default ExecutorService. Port of
// Executors.newCachedThreadPool().
func NewCachedThreadPoolExecutorService() ExecutorService {
	return &cachedThreadPoolExecutorService{}
}

func (e *cachedThreadPoolExecutorService) Submit(task Runnable) {
	go task.Run()
}

func (e *cachedThreadPoolExecutorService) ShutdownNow() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.shutdown = true
}

func (e *cachedThreadPoolExecutorService) IsShutdown() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.shutdown
}

// ValidationJobOverrides declares the operations Java's abstract ValidationJob<D,L,C> class
// leaves abstract, standing in for the virtual dispatch the base needs to reach the concrete
// job. A concrete job registers itself with ValidationJob.InitValidationJob.
type ValidationJobOverrides[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D], C CacheAccessFactory] interface {
	// GetValidationJobSummaryBuilder instantiates a builder to create a
	// modeljob.ValidationJobSummary. Port of the abstract protected
	// getValidationJobSummaryBuilder().
	GetValidationJobSummaryBuilder() ValidationJobSummaryBuilder[D, L]

	// GetDocumentAnalysis returns a Runnable to perform an analysis of a document. Port of
	// the abstract protected getDocumentAnalysis(DocumentSource, DSSFileLoader,
	// CountDownLatch).
	GetDocumentAnalysis(documentSource *DocumentSource, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) Runnable

	// GetDocumentListAnalysis returns a Runnable to perform an analysis of a document list.
	// Port of the abstract protected getDocumentListAnalysis(DocumentSource, DSSFileLoader,
	// CountDownLatch).
	GetDocumentListAnalysis(documentSource *DocumentSource, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) Runnable

	// ExtractOtherDocumentSources extracts other document sources from a collection of
	// document lists. Port of the abstract protected extractOtherDocumentSources().
	ExtractOtherDocumentSources() []*DocumentSource

	// HandleDocumentChanges executes logic on occurred document changes. Port of the
	// abstract protected handleDocumentChanges(Map, Map).
	HandleDocumentChanges(oldParsingValues, newParsingValues map[CacheKey]modeljob.ParsingInfoRecord)

	// SynchronizeCertificateSources synchronizes the certificate source, if applicable.
	// Port of the abstract protected synchronizeCertificateSources().
	SynchronizeCertificateSources()
}

// ValidationJob is the main class performing the document download / parsing / validation
// tasks. D is the current modeljob.DocumentInfo, L the parent modeljob.DocumentListInfo, C
// the CacheAccessFactory, mirroring Java's "ValidationJob<D extends DocumentInfo<L>, L
// extends DocumentListInfo<L, D>, C extends CacheAccessFactory>".
//
// slf4j logging is dropped (no observable behavior).
//
// DEVIATION: Java's getSummary()/offlineRefresh()/onlineRefresh() are `synchronized`
// (reentrant monitor) methods, and offlineRefresh/onlineRefresh call the private refresh()
// which itself calls getSummary() — legal in Java because the monitor is reentrant. Go's
// sync.Mutex is not reentrant, so the mutex here guards only the three public entry points
// (GetSummary, OfflineRefresh, OnlineRefresh); refresh() calls the unexported buildSummary()
// directly instead of the locking GetSummary(), preserving the same critical section
// semantics without deadlocking.
type ValidationJob[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D], C CacheAccessFactory] struct {
	// overrides points back at the concrete job; see InitValidationJob.
	overrides ValidationJobOverrides[D, L, C]

	mu sync.Mutex

	// cacheAccessFactory contains all caches for the current validation job.
	cacheAccessFactory C

	// executorService provides methods to manage the asynchronous behaviour.
	executorService ExecutorService

	// documentSources is an array of zero, one or more document sources not referenced in a
	// document list(s).
	documentSources []*DocumentSource

	// documentListSources is an array of zero, one or more document list sources.
	documentListSources []*DocumentSource

	// offlineLoader is used for offline data loading from a local source.
	offlineLoader http.DSSFileLoader

	// onlineLoader is used for online data loading from a remote source.
	onlineLoader http.DSSFileLoader

	// cacheCleaner is used to clean the cache.
	cacheCleaner *CacheCleaner

	// synchronizationStrategy is the strategy to follow to synchronize the certificates.
	// Default: all documents and document lists are synchronized.
	synchronizationStrategy SynchronizationStrategy[D, L]

	// debug allows printing the cache content before and after the synchronization
	// (default: false).
	debug bool

	// documentListAlerts holds document list info alerts.
	documentListAlerts []alert.Alert[L]

	// documentAlerts holds document info alerts.
	documentAlerts []alert.Alert[D]
}

// NewValidationJob creates a ValidationJob with the given cache access factory. Port of the
// protected constructor.
func NewValidationJob[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D], C CacheAccessFactory](cacheAccessFactory C) ValidationJob[D, L, C] {
	return ValidationJob[D, L, C]{
		cacheAccessFactory:      cacheAccessFactory,
		executorService:         NewCachedThreadPoolExecutorService(),
		synchronizationStrategy: NewAcceptAllStrategy[D, L](),
	}
}

// InitValidationJob registers the concrete job with its base so that the base can dispatch
// to the abstract hooks. It must be called exactly once, by the concrete job's constructor,
// before any other method.
func (j *ValidationJob[D, L, C]) InitValidationJob(overrides ValidationJobOverrides[D, L, C]) {
	j.overrides = overrides
}

func (j *ValidationJob[D, L, C]) validationJobOverrides() ValidationJobOverrides[D, L, C] {
	if j.overrides == nil {
		panic("ValidationJob was not initialised: the concrete job must call InitValidationJob in its constructor")
	}
	return j.overrides
}

// GetCacheAccessFactory returns the cache access factory. Port of getCacheAccessFactory().
func (j *ValidationJob[D, L, C]) GetCacheAccessFactory() C {
	return j.cacheAccessFactory
}

// GetDocumentSources returns the independent document sources. Port of
// getDocumentSources().
func (j *ValidationJob[D, L, C]) GetDocumentSources() []*DocumentSource {
	return j.documentSources
}

// SetDocumentSources sets the additional document sources. Port of
// setDocumentSources(DocumentSource...).
func (j *ValidationJob[D, L, C]) SetDocumentSources(documentSources ...*DocumentSource) {
	j.documentSources = documentSources
}

// GetDocumentListSources returns the document list sources. Port of
// getDocumentListSources().
func (j *ValidationJob[D, L, C]) GetDocumentListSources() []*DocumentSource {
	return j.documentListSources
}

// SetDocumentListSources sets the document lists sources. Port of
// setDocumentListSources(DocumentSource...).
func (j *ValidationJob[D, L, C]) SetDocumentListSources(documentListSources ...*DocumentSource) {
	j.documentListSources = documentListSources
}

// SetExecutorService sets the execution service to manage the asynchronous behaviour,
// shutting down the previous one first if it was not already shut down. Port of
// setExecutorService(ExecutorService).
func (j *ValidationJob[D, L, C]) SetExecutorService(executorService ExecutorService) {
	if j.executorService != nil && !j.executorService.IsShutdown() {
		j.executorService.ShutdownNow()
	}
	j.executorService = executorService
}

// SetOfflineDataLoader sets the offline DSSFileLoader used for data loading from the local
// source. Port of setOfflineDataLoader(DSSFileLoader).
func (j *ValidationJob[D, L, C]) SetOfflineDataLoader(offlineLoader http.DSSFileLoader) {
	j.offlineLoader = offlineLoader
}

// SetOnlineDataLoader sets the online DSSFileLoader used for data loading from a remote
// source. Port of setOnlineDataLoader(DSSFileLoader).
func (j *ValidationJob[D, L, C]) SetOnlineDataLoader(onlineLoader http.DSSFileLoader) {
	j.onlineLoader = onlineLoader
}

// SetCacheCleaner sets the cacheCleaner. Port of setCacheCleaner(CacheCleaner).
func (j *ValidationJob[D, L, C]) SetCacheCleaner(cacheCleaner *CacheCleaner) {
	j.cacheCleaner = cacheCleaner
}

// GetSynchronizationStrategy returns the synchronization strategy. Port of
// getSynchronizationStrategy().
func (j *ValidationJob[D, L, C]) GetSynchronizationStrategy() SynchronizationStrategy[D, L] {
	return j.synchronizationStrategy
}

// SetSynchronizationStrategy sets the strategy to follow for the certificate
// synchronization. Panics if synchronizationStrategy is nil, mirroring Java's
// Objects.requireNonNull("The SynchronizationStrategy cannot be null"). Port of
// setSynchronizationStrategy(SynchronizationStrategy).
func (j *ValidationJob[D, L, C]) SetSynchronizationStrategy(synchronizationStrategy SynchronizationStrategy[D, L]) {
	if synchronizationStrategy == nil {
		panic("The SynchronizationStrategy cannot be null")
	}
	j.synchronizationStrategy = synchronizationStrategy
}

// SetDebug sets the debug mode (print the cache contents before and after the
// synchronization). Port of setDebug(boolean).
func (j *ValidationJob[D, L, C]) SetDebug(debug bool) {
	j.debug = debug
}

// SetDocumentListAlerts sets the document list alerts to be processed. Port of
// setDocumentListAlerts(List).
func (j *ValidationJob[D, L, C]) SetDocumentListAlerts(documentListAlerts []alert.Alert[L]) {
	j.documentListAlerts = documentListAlerts
}

// SetDocumentAlerts sets the document alerts to be processed. Port of
// setDocumentAlerts(List).
func (j *ValidationJob[D, L, C]) SetDocumentAlerts(documentAlerts []alert.Alert[D]) {
	j.documentAlerts = documentAlerts
}

// GetSummary returns the validation job summary for all processed documents / document
// lists. Port of the synchronized getSummary(). See the type's DEVIATION note on the
// locking strategy.
func (j *ValidationJob[D, L, C]) GetSummary() modeljob.ValidationJobSummary[D, L] {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.buildSummary()
}

// buildSummary builds the summary without locking mu; callers already holding mu (refresh)
// call this directly instead of GetSummary.
func (j *ValidationJob[D, L, C]) buildSummary() modeljob.ValidationJobSummary[D, L] {
	return j.validationJobOverrides().GetValidationJobSummaryBuilder().Build()
}

// OfflineRefresh executes the refresh in offline mode (no data from remote sources will be
// downloaded). By default used on initialization. Panics if no offline loader was
// configured, mirroring Java's Objects.requireNonNull("The offlineLoader must be
// defined!"). Port of the synchronized offlineRefresh().
func (j *ValidationJob[D, L, C]) OfflineRefresh() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.offlineLoader == nil {
		panic("The offlineLoader must be defined!")
	}
	return j.refresh(j.offlineLoader)
}

// OnlineRefresh executes the refresh in online mode (all data will be updated from remote
// sources). Used as default database update. Panics if no online loader was configured,
// mirroring Java's Objects.requireNonNull("The onlineLoader must be defined!"). Port of the
// synchronized onlineRefresh().
func (j *ValidationJob[D, L, C]) OnlineRefresh() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.onlineLoader == nil {
		panic("The onlineLoader must be defined!")
	}
	return j.refresh(j.onlineLoader)
}

// refresh ports the private refresh(DSSFileLoader). Errors returned here are ported from
// Java's unchecked DSSException raised by checkNoDuplicateUrls, propagated to the caller
// instead of being thrown, per the checked/unchecked-exception-to-error porting convention.
func (j *ValidationJob[D, L, C]) refresh(dssFileLoader http.DSSFileLoader) error {
	var currentDocSources []*DocumentSource
	if j.documentSources != nil {
		currentDocSources = append(currentDocSources, j.documentSources...)
	}

	// Execute all document lists
	if utils.IsArrayNotEmpty(j.documentListSources) {
		docLists := j.documentListSources

		if err := j.executeDocumentListSourcesAnalysis(docLists, dssFileLoader); err != nil {
			return err
		}

		// extract document sources from cached document lists
		currentDocSources = append(currentDocSources, j.validationJobOverrides().ExtractOtherDocumentSources()...)
	}

	// And then, execute all TLs (manual configs + documents from document lists)
	if err := j.executeDocumentSourcesAnalysis(currentDocSources, dssFileLoader); err != nil {
		return err
	}

	// alerts()
	if utils.IsCollectionNotEmpty(j.documentListAlerts) || utils.IsCollectionNotEmpty(j.documentAlerts) {
		jobSummary := j.buildSummary()
		alerter := NewValidationJobAlerter[D, L](j.documentListAlerts, j.documentAlerts)
		alerter.DetectChanges(jobSummary)
	}

	if j.debug {
		j.cacheAccessFactory.GetDebugCacheAccess().Dump()
	}

	// certificate source sync + cache sync if needed
	j.validationJobOverrides().SynchronizeCertificateSources()

	j.executeCacheCleaner()

	if j.debug {
		j.cacheAccessFactory.GetDebugCacheAccess().Dump()
	}

	return nil
}

// executeDocumentListSourcesAnalysis ports the private
// executeDocumentListSourcesAnalysis(List, DSSFileLoader).
//
// DEVIATION: Java's latch.await() catches InterruptedException, restores the interrupt flag
// and logs; sync.WaitGroup.Wait() has no interruptible/error-returning form, so this simply
// blocks until every submitted analysis has called Done().
func (j *ValidationJob[D, L, C]) executeDocumentListSourcesAnalysis(docListSources []*DocumentSource, dssFileLoader http.DSSFileLoader) error {
	nbDocListSources := len(docListSources)
	if nbDocListSources == 0 {
		return nil
	}

	if err := j.checkNoDuplicateUrls(docListSources); err != nil {
		return err
	}

	oldParsingValues := j.extractParsingCache(docListSources)

	var latch sync.WaitGroup
	latch.Add(nbDocListSources)
	for _, docListSource := range docListSources {
		j.executorService.Submit(j.validationJobOverrides().GetDocumentListAnalysis(docListSource, dssFileLoader, &latch))
	}
	latch.Wait()

	newParsingValues := j.extractParsingCache(docListSources)

	j.validationJobOverrides().HandleDocumentChanges(oldParsingValues, newParsingValues)
	return nil
}

// extractParsingCache ports the private extractParsingCache(List<DocumentSource>).
func (j *ValidationJob[D, L, C]) extractParsingCache(docListSources []*DocumentSource) map[CacheKey]modeljob.ParsingInfoRecord {
	readOnlyCacheAccess := j.cacheAccessFactory.GetReadOnlyCacheAccess()
	result := make(map[CacheKey]modeljob.ParsingInfoRecord, len(docListSources))
	for _, s := range docListSources {
		result[s.CacheKey()] = readOnlyCacheAccess.GetParsingInfoRecord(s.CacheKey())
	}
	return result
}

// executeDocumentSourcesAnalysis ports the private executeDocumentSourcesAnalysis(List,
// DSSFileLoader). See executeDocumentListSourcesAnalysis for the latch-interruption
// DEVIATION note (applies identically here).
func (j *ValidationJob[D, L, C]) executeDocumentSourcesAnalysis(docSources []*DocumentSource, dssFileLoader http.DSSFileLoader) error {
	nbDocSources := len(docSources)
	if nbDocSources == 0 {
		return nil
	}

	if err := j.checkNoDuplicateUrls(docSources); err != nil {
		return err
	}

	var latch sync.WaitGroup
	latch.Add(nbDocSources)
	for _, docSource := range docSources {
		j.executorService.Submit(j.validationJobOverrides().GetDocumentAnalysis(docSource, dssFileLoader, &latch))
	}
	latch.Wait()

	return nil
}

// executeCacheCleaner ports the private executeCacheCleaner().
func (j *ValidationJob[D, L, C]) executeCacheCleaner() {
	if j.cacheCleaner == nil {
		return
	}

	cacheKeys := j.cacheAccessFactory.GetReadOnlyCacheAccess().GetAllCacheKeys()
	for cacheKey := range cacheKeys {
		cacheAccess := j.cacheAccessFactory.GetCacheAccess(cacheKey)
		j.cacheCleaner.Clean(cacheAccess)
	}
}

// checkNoDuplicateUrls ports the private checkNoDuplicateUrls(List<DocumentSource>).
// Duplicate urls mean cache conflict. Returns a *model.DSSError, mirroring Java's thrown
// (unchecked) DSSException.
func (j *ValidationJob[D, L, C]) checkNoDuplicateUrls(sources []*DocumentSource) error {
	allUrls := make([]string, 0, len(sources))
	for _, s := range sources {
		allUrls = append(allUrls, s.Url())
	}
	uniqueUrls := make(map[string]struct{}, len(allUrls))
	for _, u := range allUrls {
		uniqueUrls[u] = struct{}{}
	}
	if len(allUrls) > len(uniqueUrls) {
		return model.NewDSSError(fmt.Sprintf("Duplicate urls found : %s", javaListString(allUrls)))
	}
	return nil
}

// javaListString renders a []string the way java.util.List#toString() does: "[a, b, c]".
func javaListString(values []string) string {
	return "[" + strings.Join(values, ", ") + "]"
}
