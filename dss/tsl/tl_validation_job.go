// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/job/TLValidationJob.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY: TLValidationJob extends the generic
// github.com/utain/esig/dss/validation/job.ValidationJob[D, L, C] (dss-validation-job, VALJOB
// chunk, now landed - see validation_job.go, abstract_analysis.go and
// abstract_runnable_analysis.go there for the base's exact API this file relies on).
//
// JUDGMENT CALL - recovering TLSource/LOTLSource identity: job.ValidationJob stores document
// sources as []*job.DocumentSource (the concrete base struct TLSource/LOTLSource embed by value,
// not a polymorphic reference the way Java's covariant TLSource[]/LOTLSource[] arrays preserve
// object identity through the base DocumentSource[] type). GetDocumentAnalysis/
// GetDocumentListAnalysis therefore only ever receive a *job.DocumentSource back, with no way to
// recover the original *TLSource/*LOTLSource's own fields (predicates, TL versions, pivot
// support, ...) from the pointer alone. This port keeps side maps keyed by job.CacheKey
// (tlSourcesByCacheKey/lotlSourcesByCacheKey), populated whenever a TLSource/LOTLSource is
// registered (SetTrustedListSources, SetListOfTrustedListSources, ExtractOtherDocumentSources),
// and looks the concrete source back up by the *job.DocumentSource's own CacheKey() - which is
// exactly the join key job.ValidationJob itself uses internally (extractParsingCache,
// checkNoDuplicateUrls), so this preserves the Java behaviour without needing an interface
// abstraction over DocumentSource.
package tsl

import (
	"sync"

	"github.com/utain/esig/dss/alert"
	modeljob "github.com/utain/esig/dss/model/job"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/validation/job"
)

// TLValidationJob is the main class performing the TL/LOTL download / parsing / validation
// tasks.
type TLValidationJob struct {
	job.ValidationJob[*tslmodel.TLInfo, *tslmodel.LOTLInfo, *TLCacheAccessFactory]

	// trustPropertiesCertificateSource is the certificate source to be synchronized.
	trustPropertiesCertificateSource tslmodel.TrustPropertiesCertificateSource

	// tlSourcesByCacheKey/lotlSourcesByCacheKey recover the concrete *TLSource/*LOTLSource for a
	// *job.DocumentSource the base hands back - see this file's header JUDGMENT CALL.
	tlSourcesByCacheKey   map[job.CacheKey]*TLSource
	lotlSourcesByCacheKey map[job.CacheKey]*LOTLSource
}

var _ job.ValidationJobOverrides[*tslmodel.TLInfo, *tslmodel.LOTLInfo, *TLCacheAccessFactory] = (*TLValidationJob)(nil)

// NewTLValidationJob is the default constructor instantiating the object with null
// configuration.
func NewTLValidationJob() *TLValidationJob {
	j := &TLValidationJob{
		ValidationJob:         job.NewValidationJob[*tslmodel.TLInfo, *tslmodel.LOTLInfo, *TLCacheAccessFactory](NewTLCacheAccessFactory()),
		tlSourcesByCacheKey:   make(map[job.CacheKey]*TLSource),
		lotlSourcesByCacheKey: make(map[job.CacheKey]*LOTLSource),
	}
	j.InitValidationJob(j)
	return j
}

// TLSources ports the covariant override protected TLSource[] getDocumentSources().
func (j *TLValidationJob) TLSources() []*TLSource {
	sources := j.GetDocumentSources()
	result := make([]*TLSource, 0, len(sources))
	for _, s := range sources {
		if tlSource, ok := j.tlSourcesByCacheKey[s.CacheKey()]; ok {
			result = append(result, tlSource)
		}
	}
	return result
}

// SetTrustedListSources sets the additional TL Sources. Port of setTrustedListSources(TLSource...).
func (j *TLValidationJob) SetTrustedListSources(trustedListSources ...*TLSource) {
	refs := make([]*job.DocumentSource, 0, len(trustedListSources))
	for _, s := range trustedListSources {
		j.tlSourcesByCacheKey[s.CacheKey()] = s
		refs = append(refs, &s.DocumentSource)
	}
	j.SetDocumentSources(refs...)
}

// LOTLSources ports the covariant override protected LOTLSource[] getDocumentListSources().
func (j *TLValidationJob) LOTLSources() []*LOTLSource {
	sources := j.GetDocumentListSources()
	result := make([]*LOTLSource, 0, len(sources))
	for _, s := range sources {
		if lotlSource, ok := j.lotlSourcesByCacheKey[s.CacheKey()]; ok {
			result = append(result, lotlSource)
		}
	}
	return result
}

// SetListOfTrustedListSources sets the LOTL Sources. Port of setListOfTrustedListSources(LOTLSource...).
func (j *TLValidationJob) SetListOfTrustedListSources(listOfTrustedListSources ...*LOTLSource) {
	refs := make([]*job.DocumentSource, 0, len(listOfTrustedListSources))
	for _, s := range listOfTrustedListSources {
		j.lotlSourcesByCacheKey[s.CacheKey()] = s
		refs = append(refs, &s.DocumentSource)
	}
	j.SetDocumentListSources(refs...)
}

// SetTrustedListCertificateSource sets the TrustedListsCertificateSource to be filled with the
// job. Port of setTrustedListCertificateSource(TrustPropertiesCertificateSource).
func (j *TLValidationJob) SetTrustedListCertificateSource(trustPropertiesCertificateSource tslmodel.TrustPropertiesCertificateSource) {
	j.trustPropertiesCertificateSource = trustPropertiesCertificateSource
}

// SetLOTLAlerts sets the LOTL alerts to be processed. Port of setLOTLAlerts(List).
func (j *TLValidationJob) SetLOTLAlerts(lotlAlerts []*LOTLAlert) {
	converted := make([]alert.Alert[tslmodel.LOTLInfo], len(lotlAlerts))
	for i, a := range lotlAlerts {
		converted[i] = a
	}
	j.SetDocumentListAlerts(converted)
}

// SetTLAlerts sets the TL alerts to be processed. Port of setTLAlerts(List).
func (j *TLValidationJob) SetTLAlerts(tlAlerts []*TLAlert) {
	converted := make([]alert.Alert[tslmodel.TLInfo], len(tlAlerts))
	for i, a := range tlAlerts {
		converted[i] = a
	}
	j.SetDocumentAlerts(converted)
}

// Summary returns the validation job summary for all processed LOTL / TLs. Port of the
// synchronized getSummary(). The base's GetSummary() returns the generic
// modeljob.ValidationJobSummary[TLInfo, LOTLInfo] interface (a *tlValidationJobSummaryAdapter,
// see tl_validation_job_summary_builder.go); this recovers the concrete summary through it.
func (j *TLValidationJob) Summary() *tslmodel.TLValidationJobSummary {
	adapter := j.GetSummary().(*tlValidationJobSummaryAdapter)
	return adapter.TLValidationJobSummary
}

// GetValidationJobSummaryBuilder ports the protected getValidationJobSummaryBuilder() override.
func (j *TLValidationJob) GetValidationJobSummaryBuilder() job.ValidationJobSummaryBuilder[*tslmodel.TLInfo, *tslmodel.LOTLInfo] {
	return NewTLValidationJobSummaryBuilder(j.GetCacheAccessFactory().ReadOnlyCacheAccess(), j.TLSources(), j.LOTLSources()).BuildTyped()
}

// GetDocumentAnalysis ports the protected getDocumentAnalysis(DocumentSource, DSSFileLoader,
// CountDownLatch) override. Panics with the Java message when documentSource does not correspond
// to a registered TLSource.
func (j *TLValidationJob) GetDocumentAnalysis(documentSource *job.DocumentSource, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) job.Runnable {
	tlSource, ok := j.tlSourcesByCacheKey[documentSource.CacheKey()]
	if !ok {
		panic("The provided document source is not a TLSource!")
	}
	cacheAccess := j.GetCacheAccessFactory().CacheAccess(tlSource.CacheKey())
	return NewTLAnalysis(tlSource, cacheAccess, dssFileLoader, latch)
}

// GetDocumentListAnalysis ports the protected getDocumentListAnalysis(DocumentSource,
// DSSFileLoader, CountDownLatch) override. Panics with the Java message when documentSource does
// not correspond to a registered LOTLSource.
func (j *TLValidationJob) GetDocumentListAnalysis(documentSource *job.DocumentSource, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) job.Runnable {
	lotlSource, ok := j.lotlSourcesByCacheKey[documentSource.CacheKey()]
	if !ok {
		panic("The provided document source is not a LOTLSource!")
	}
	cacheAccess := j.GetCacheAccessFactory().CacheAccess(documentSource.CacheKey())
	if lotlSource.IsPivotSupport() {
		return NewLOTLWithPivotsAnalysis(lotlSource, cacheAccess, dssFileLoader, j.GetCacheAccessFactory(), latch)
	}
	return NewLOTLAnalysis(lotlSource, cacheAccess, dssFileLoader, latch)
}

// ExtractOtherDocumentSources ports the protected extractOtherDocumentSources() override.
func (j *TLValidationJob) ExtractOtherDocumentSources() []*job.DocumentSource {
	lotlSources := j.LOTLSources()
	if len(lotlSources) == 0 {
		return nil
	}
	tlSourceBuilder := NewTLSourceBuilder(lotlSources, j.extractParsingCache(lotlSources))
	built := tlSourceBuilder.Build()
	refs := make([]*job.DocumentSource, 0, len(built))
	for _, s := range built {
		j.tlSourcesByCacheKey[s.CacheKey()] = s
		refs = append(refs, &s.DocumentSource)
	}
	return refs
}

// extractParsingCache ports the private extractParsingCache(List<LOTLSource>).
func (j *TLValidationJob) extractParsingCache(lotlSources []*LOTLSource) map[job.CacheKey]*TLParsingCacheDTO {
	readOnlyCacheAccess := j.GetCacheAccessFactory().ReadOnlyCacheAccess()
	result := make(map[job.CacheKey]*TLParsingCacheDTO, len(lotlSources))
	for _, s := range lotlSources {
		result[s.CacheKey()] = readOnlyCacheAccess.GetParsingInfoRecordTyped(s.CacheKey())
	}
	return result
}

// SynchronizeCertificateSources ports the protected synchronizeCertificateSources() override.
func (j *TLValidationJob) SynchronizeCertificateSources() {
	if j.trustPropertiesCertificateSource == nil {
		return
	}

	synchronizer := NewTrustedListCertificateSourceSynchronizer(
		j.TLSources(), j.LOTLSources(), j.trustPropertiesCertificateSource, j.GetSynchronizationStrategy(),
		j.GetCacheAccessFactory().GetSynchronizerCacheAccess(), j.GetCacheAccessFactory().ReadOnlyCacheAccess())
	synchronizer.Sync()
}

// HandleDocumentChanges ports the protected handleDocumentChanges(Map, Map) override.
func (j *TLValidationJob) HandleDocumentChanges(oldParsingValues, newParsingValues map[job.CacheKey]modeljob.ParsingInfoRecord) {
	lotlChangeApplier := NewLOTLChangeApplier(j.GetCacheAccessFactory().GetDocumentChangesCacheAccess(), oldParsingValues, newParsingValues)
	lotlChangeApplier.AnalyzeAndApply()
}
