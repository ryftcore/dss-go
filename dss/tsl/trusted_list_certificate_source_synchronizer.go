// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sync/TrustedListCertificateSourceSynchronizer.java (DSS 6.5.RC1).
//
// Uses job.CacheKey and job.SynchronizerCacheAccess (see xml_download_result.go's header for
// the wider job.* convention).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// TrustedListCertificateSourceSynchronizer loads the trusted certificate source.
type TrustedListCertificateSourceSynchronizer struct {
	// tlSources is the list of TLSources to extract summary for.
	tlSources []*TLSource

	// lotlSources is the list of LOTLSource to extract summary for.
	lotlSources []*LOTLSource

	// synchronizationStrategy is the strategy to follow for the certificate synchronization.
	synchronizationStrategy job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo]

	// certificateSource is the certificate source to be synchronized.
	certificateSource tslmodel.TrustPropertiesCertificateSource

	// syncCacheAccess is the cache access.
	syncCacheAccess *job.SynchronizerCacheAccess

	// readOnlyCacheAccess is the cache access.
	readOnlyCacheAccess *TLReadOnlyCacheAccess
}

// NewTrustedListCertificateSourceSynchronizer is the default constructor.
func NewTrustedListCertificateSourceSynchronizer(tlSources []*TLSource, lotlSources []*LOTLSource,
	certificateSource tslmodel.TrustPropertiesCertificateSource, synchronizationStrategy job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo],
	syncCacheAccess *job.SynchronizerCacheAccess, readOnlyCacheAccess *TLReadOnlyCacheAccess) *TrustedListCertificateSourceSynchronizer {
	return &TrustedListCertificateSourceSynchronizer{
		tlSources: tlSources, lotlSources: lotlSources, synchronizationStrategy: synchronizationStrategy,
		certificateSource: certificateSource, syncCacheAccess: syncCacheAccess, readOnlyCacheAccess: readOnlyCacheAccess,
	}
}

// Sync synchronizes the trusted certificate source based on the validation job processing
// result. Port of sync(). Java catches and logs any Exception; this recovers a panic the same
// way, since the private helper methods below panic instead of returning an error.
func (s *TrustedListCertificateSourceSynchronizer) Sync() {
	defer func() {
		_ = recover()
	}()

	summaryBuilder := NewTLValidationJobSummaryBuilder(s.readOnlyCacheAccess, s.tlSources, s.lotlSources)

	summary := summaryBuilder.Build()
	if s.isCertificateSyncNeeded(summary) {
		s.synchronizeCertificates(summary)
	}
	s.syncCache(summary)

	// re-build summary after synchronization
	summary = summaryBuilder.Build()
	s.certificateSource.SetSummary(summary)
}

func (s *TrustedListCertificateSourceSynchronizer) isCertificateSyncNeeded(summary *tslmodel.TLValidationJobSummary) bool {
	for _, lotlInfo := range summary.LOTLInfos() {
		if s.isTLParsingDesyncOrError(&lotlInfo.TLInfo) || s.isTLParsingDesyncOrErrorList(lotlInfo.TLInfos()) {
			return true
		}
	}
	return s.isTLParsingDesyncOrErrorList(summary.OtherTLInfos())
}

func (s *TrustedListCertificateSourceSynchronizer) isTLParsingDesyncOrErrorList(tlInfos []*tslmodel.TLInfo) bool {
	for _, tlInfo := range tlInfos {
		if s.isTLParsingDesyncOrError(tlInfo) {
			return true
		}
	}
	return false
}

func (s *TrustedListCertificateSourceSynchronizer) isTLParsingDesyncOrError(tlInfo *tslmodel.TLInfo) bool {
	parsingCacheInfo := tlInfo.ParsingCacheInfo()
	return parsingCacheInfo == nil || parsingCacheInfo.IsDesynchronized() || parsingCacheInfo.IsError()
}

func (s *TrustedListCertificateSourceSynchronizer) synchronizeCertificates(summary *tslmodel.TLValidationJobSummary) {
	trustPropertiesByCerts := newSynchronizerCertificateMap[[]*tslmodel.TrustProperties]()
	trustTimeByCerts := newSynchronizerCertificateMap[[]*tslmodel.CertificateTrustTime]()
	for _, lotlInfo := range summary.LOTLInfos() {
		if s.synchronizationStrategy.CanBeSynchronizedDocumentList(lotlInfo) {
			s.addCertificatesFromTLs(trustPropertiesByCerts, trustTimeByCerts, lotlInfo.TLInfos(), lotlInfo)
		}
	}
	s.addCertificatesFromTLs(trustPropertiesByCerts, trustTimeByCerts, summary.OtherTLInfos(), nil)
	s.certificateSource.SetTrustPropertiesByCertificates(trustPropertiesByCerts.asMap())
	s.certificateSource.SetTrustTimeByCertificates(trustTimeByCerts.asMap())
}

// synchronizerCertificateMap stands in for the two java.util.HashMap<CertificateToken, List<?>>
// accumulators synchronizeCertificates() builds.
//
// A plain Go map[*model.CertificateToken]V would NOT be that map: CertificateToken#equals (and
// hashCode) is the certificate's DSS id, i.e. a digest of its bytes, so Java merges the same CA
// certificate reached through two different trusted lists - parsed into two distinct
// CertificateToken instances - into ONE entry, where a pointer-keyed Go map keeps two. Keying by
// DSS id reproduces Java's merge, and remembering the first-seen token per id keeps a single,
// stable representative for the map handed to the certificate source.
type synchronizerCertificateMap[V any] struct {
	// order lists the ids in first-encounter order, so nothing here depends on Go map
	// iteration order.
	order []string
	// tokens maps a DSS id to the first CertificateToken instance seen for it.
	tokens map[string]*model.CertificateToken
	// values maps a DSS id to its accumulated value.
	values map[string]V
}

func newSynchronizerCertificateMap[V any]() *synchronizerCertificateMap[V] {
	return &synchronizerCertificateMap[V]{
		tokens: make(map[string]*model.CertificateToken),
		values: make(map[string]V),
	}
}

// get answers the value accumulated for the given certificate (the zero value when absent),
// standing in for Map#get.
func (m *synchronizerCertificateMap[V]) get(certificate *model.CertificateToken) V {
	return m.values[certificate.DSSIDAsString()]
}

// put stores the value for the given certificate, standing in for Map#put.
func (m *synchronizerCertificateMap[V]) put(certificate *model.CertificateToken, value V) {
	id := certificate.DSSIDAsString()
	if _, seen := m.tokens[id]; !seen {
		m.tokens[id] = certificate
		m.order = append(m.order, id)
	}
	m.values[id] = value
}

// asMap materialises the accumulator as the map the TrustPropertiesCertificateSource setters
// take, one entry per distinct certificate.
func (m *synchronizerCertificateMap[V]) asMap() map[*model.CertificateToken]V {
	out := make(map[*model.CertificateToken]V, len(m.order))
	for _, id := range m.order {
		out[m.tokens[id]] = m.values[id]
	}
	return out
}

func (s *TrustedListCertificateSourceSynchronizer) addCertificatesFromTLs(trustPropertiesByCerts *synchronizerCertificateMap[[]*tslmodel.TrustProperties],
	trustTimeByCerts *synchronizerCertificateMap[[]*tslmodel.CertificateTrustTime], tlInfos []*tslmodel.TLInfo, relatedLOTL *tslmodel.LOTLInfo) {

	for _, tlInfo := range tlInfos {
		if !s.synchronizationStrategy.CanBeSynchronizedDocument(tlInfo) {
			continue
		}
		parsingCacheInfo, ok := tlInfo.TLParsingCacheInfo()
		if !ok || !parsingCacheInfo.IsResultExist() {
			continue
		}
		trustServiceProviders := parsingCacheInfo.TrustServiceProviders()
		if !utils.IsCollectionNotEmpty(trustServiceProviders) {
			continue
		}
		trustAnchorValidityPredicate := s.trustAnchorValidityPredicate(tlInfo, relatedLOTL)
		for _, original := range trustServiceProviders {
			detached := s.detached(original)
			for _, trustService := range original.Services() {
				statusAndInformationExtensions := trustService.StatusAndInformationExtensions()
				trustProperties := s.trustProperties(relatedLOTL, tlInfo, detached, statusAndInformationExtensions)
				certificateTrustTimes := s.certificateTrustTimes(statusAndInformationExtensions, trustAnchorValidityPredicate)
				for _, certificate := range trustService.Certificates() {
					s.addCertificate(trustPropertiesByCerts, trustTimeByCerts, certificate, trustProperties, certificateTrustTimes)
				}
			}
		}
	}
}

func (s *TrustedListCertificateSourceSynchronizer) addCertificate(trustPropertiesByCerts *synchronizerCertificateMap[[]*tslmodel.TrustProperties],
	trustTimeByCerts *synchronizerCertificateMap[[]*tslmodel.CertificateTrustTime], certificate *model.CertificateToken,
	trustProperties *tslmodel.TrustProperties, certificateTrustTimes []*tslmodel.CertificateTrustTime) {

	trustPropertiesList := trustPropertiesByCerts.get(certificate)
	if !trustPropertiesListContains(trustPropertiesList, trustProperties) {
		trustPropertiesList = append(trustPropertiesList, trustProperties)
	}
	trustPropertiesByCerts.put(certificate, trustPropertiesList)
	certificateTrustTimeList := trustTimeByCerts.get(certificate)
	for _, certificateTrustTime := range certificateTrustTimes {
		if !certificateTrustTimeListContains(certificateTrustTimeList, certificateTrustTime) {
			certificateTrustTimeList = append(certificateTrustTimeList, certificateTrustTime)
		}
	}
	trustTimeByCerts.put(certificate, certificateTrustTimeList)
}

func (s *TrustedListCertificateSourceSynchronizer) detached(original *tslmodel.TrustServiceProvider) *tslmodel.TrustServiceProvider {
	builder := NewTrustServiceProviderBuilderFromOriginal(original)
	builder.SetServices(nil)
	return builder.Build()
}

func (s *TrustedListCertificateSourceSynchronizer) syncCache(summary *tslmodel.TLValidationJobSummary) {
	for _, lotlInfo := range summary.LOTLInfos() {
		s.syncTLInfosCache(lotlInfo.TLInfos())
		s.syncPivotsCache(lotlInfo.PivotInfos())
		s.syncCacheAccess.Sync(job.NewCacheKey(lotlInfo.Url()))
	}
	s.syncTLInfosCache(summary.OtherTLInfos())
}

func (s *TrustedListCertificateSourceSynchronizer) syncPivotsCache(pivotInfos []*tslmodel.PivotInfo) {
	for _, pivotInfo := range pivotInfos {
		s.syncCacheAccess.Sync(job.NewCacheKey(pivotInfo.Url()))
	}
}

func (s *TrustedListCertificateSourceSynchronizer) syncTLInfosCache(tlInfos []*tslmodel.TLInfo) {
	for _, tlInfo := range tlInfos {
		s.syncCacheAccess.Sync(job.NewCacheKey(tlInfo.Url()))
	}
}

func (s *TrustedListCertificateSourceSynchronizer) trustProperties(relatedLOTL *tslmodel.LOTLInfo, tlInfo *tslmodel.TLInfo, detached *tslmodel.TrustServiceProvider,
	statusAndInformationExtensions *timedependent.TimeDependentValues[*tslmodel.TrustServiceStatusAndInformationExtensions]) *tslmodel.TrustProperties {
	if relatedLOTL != nil {
		return tslmodel.NewTrustPropertiesWithLOTL(relatedLOTL, tlInfo, detached, statusAndInformationExtensions)
	}
	return tslmodel.NewTrustProperties(tlInfo, detached, statusAndInformationExtensions)
}

func (s *TrustedListCertificateSourceSynchronizer) certificateTrustTimes(
	statusAndInformationExtensions *timedependent.TimeDependentValues[*tslmodel.TrustServiceStatusAndInformationExtensions],
	trustAnchorValidityPredicate TrustAnchorPeriodPredicate) []*tslmodel.CertificateTrustTime {
	if trustAnchorValidityPredicate == nil {
		// return empty instance (always valid), when no predicate is defined
		return []*tslmodel.CertificateTrustTime{tslmodel.NewCertificateTrustTime(true)}
	}

	var result []*tslmodel.CertificateTrustTime
	for trustServiceStatusAndInformation := range statusAndInformationExtensions.Iterator() {
		// TODO : add handling of MRA ?
		if trustAnchorValidityPredicate.Test(trustServiceStatusAndInformation) {
			result = append(result, tslmodel.NewCertificateTrustTimeWithRange(trustServiceStatusAndInformation.StartDate(), trustServiceStatusAndInformation.EndDate()))
		} else {
			result = append(result, tslmodel.NewCertificateTrustTime(false)) // not trusted
		}
	}
	return result
}

func (s *TrustedListCertificateSourceSynchronizer) trustAnchorValidityPredicate(tlInfo *tslmodel.TLInfo, relatedLOTLInfo *tslmodel.LOTLInfo) TrustAnchorPeriodPredicate {
	tlSource := s.relatedTLSource(tlInfo, relatedLOTLInfo)
	if tlSource != nil {
		return tlSource.TrustAnchorValidityPredicate()
	}
	return nil
}

func (s *TrustedListCertificateSourceSynchronizer) relatedTLSource(tlInfo *tslmodel.TLInfo, relatedLOTLInfo *tslmodel.LOTLInfo) *TLSource {
	if relatedLOTLInfo != nil {
		for _, lotlSource := range s.lotlSources {
			if lotlSource.Url() == relatedLOTLInfo.Url() {
				return &lotlSource.TLSource
			}
		}
	}
	for _, tlSource := range s.tlSources {
		if tlSource.Url() == tlInfo.Url() {
			return tlSource
		}
	}
	return nil
}

// trustPropertiesListContains reports whether list contains value.
//
// DEVIATION: model/tsl.TrustProperties exposes no Equals method, so this compares by pointer
// identity rather than Java's List#contains (which delegates to TrustProperties#equals()).
// Within one addCertificatesFromTLs pass the same *TrustProperties instance is reused for every
// certificate of one trust service (see trustProperties's single call site above the
// certificates loop), so pointer identity still dedups the common case exactly; it only
// under-dedups two structurally-equal but separately-constructed TrustProperties values, which
// addCertificatesFromTLs never produces for the same certificate within one summary.
func trustPropertiesListContains(list []*tslmodel.TrustProperties, value *tslmodel.TrustProperties) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// certificateTrustTimeListContains reports whether list contains value, per CertificateTrustTime's
// own Equals (matching Java's List#contains).
func certificateTrustTimeListContains(list []*tslmodel.CertificateTrustTime, value *tslmodel.CertificateTrustTime) bool {
	for _, v := range list {
		if v.Equals(value) {
			return true
		}
	}
	return false
}
