// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/summary/TLValidationJobSummaryBuilder.java (DSS 6.5.RC1).
//
// FLAG (pre-existing, model/tsl frozen package - out of this manifest): the header of
// model/tsl/tl_validation_job_summary.go claims *TLValidationJobSummary "implements the assumed
// job.ValidationJobSummary[TLInfo, LOTLInfo] interface (getDocumentListInfos/getOtherDocumentInfos,
// ported as DocumentListInfos()/OtherDocumentInfos())", but the type's actual accessors are named
// LOTLInfos()/OtherTLInfos() - they do NOT satisfy modeljob.ValidationJobSummary[TLInfo, LOTLInfo]
// (github.com/utain/esig/dss/model/job, requiring exactly DocumentListInfos()/OtherDocumentInfos())
// by name, which the now-landed validation/job package's ValidationJobSummaryBuilder[D, L]
// interface (Build() modeljob.ValidationJobSummary[D, L]) requires verbatim. Rather than editing
// the frozen model/tsl file, tlValidationJobSummaryAdapter below (this file, package tsl) wraps a
// built *TLValidationJobSummary and supplies the two missing names by delegation; BuildTyped()
// returns it as the interface job.ValidationJobSummaryBuilder[TLInfo, LOTLInfo].Build() needs,
// while Build() keeps returning the concrete *TLValidationJobSummary Java's covariant build()
// returns, for every other caller (TrustedListCertificateSourceSynchronizer.Sync(), etc.) that
// wants the concrete LOTLInfos()/OtherTLInfos() accessors.
package tsl

import (
	"github.com/utain/esig/dss/model"
	modeljob "github.com/utain/esig/dss/model/job"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/utils"
	validationjob "github.com/utain/esig/dss/validation/job"
)

// tlValidationJobSummaryAdapter adapts *tslmodel.TLValidationJobSummary to
// modeljob.ValidationJobSummary[*tslmodel.TLInfo, *tslmodel.LOTLInfo] - see this file's FLAG header.
type tlValidationJobSummaryAdapter struct {
	*tslmodel.TLValidationJobSummary
}

var _ modeljob.ValidationJobSummary[*tslmodel.TLInfo, *tslmodel.LOTLInfo] = (*tlValidationJobSummaryAdapter)(nil)

// DocumentListInfos ports the generic getDocumentListInfos(), delegating to LOTLInfos() directly:
// D/L are instantiated as *TLInfo/*LOTLInfo throughout this package (see tl_validation_job.go's
// header), and LOTLInfos() already returns []*LOTLInfo, matching without conversion.
func (a *tlValidationJobSummaryAdapter) DocumentListInfos() []*tslmodel.LOTLInfo {
	return a.LOTLInfos()
}

// OtherDocumentInfos ports the generic getOtherDocumentInfos(), delegating to OtherTLInfos()
// directly (see DocumentListInfos above).
func (a *tlValidationJobSummaryAdapter) OtherDocumentInfos() []*tslmodel.TLInfo {
	return a.OtherTLInfos()
}

// TLValidationJobSummaryBuilder builds a TLValidationJobSummary.
type TLValidationJobSummaryBuilder struct {
	// readOnlyCacheAccess is a read-only access for the cache of the current Validation Job.
	readOnlyCacheAccess *TLReadOnlyCacheAccess

	// tlSources is the list of TLSources to extract summary for.
	tlSources []*TLSource

	// lotlSources is the list of LOTLSource to extract summary for.
	lotlSources []*LOTLSource
}

var _ validationjob.ValidationJobSummaryBuilder[*tslmodel.TLInfo, *tslmodel.LOTLInfo] = tlValidationJobSummaryBuilderTyped{}

// tlValidationJobSummaryBuilderTyped adapts *TLValidationJobSummaryBuilder to
// validationjob.ValidationJobSummaryBuilder[TLInfo, LOTLInfo] - see this file's FLAG header.
type tlValidationJobSummaryBuilderTyped struct {
	builder *TLValidationJobSummaryBuilder
}

// Build satisfies validationjob.ValidationJobSummaryBuilder[TLInfo, LOTLInfo].
func (t tlValidationJobSummaryBuilderTyped) Build() modeljob.ValidationJobSummary[*tslmodel.TLInfo, *tslmodel.LOTLInfo] {
	return &tlValidationJobSummaryAdapter{TLValidationJobSummary: t.builder.Build()}
}

// NewTLValidationJobSummaryBuilder is the default constructor.
func NewTLValidationJobSummaryBuilder(readOnlyCacheAccess *TLReadOnlyCacheAccess, tlSources []*TLSource, lotlSources []*LOTLSource) *TLValidationJobSummaryBuilder {
	return &TLValidationJobSummaryBuilder{readOnlyCacheAccess: readOnlyCacheAccess, tlSources: tlSources, lotlSources: lotlSources}
}

// BuildTyped returns this builder adapted to validationjob.ValidationJobSummaryBuilder[TLInfo,
// LOTLInfo], for TLValidationJob.GetValidationJobSummaryBuilder (the
// validationjob.ValidationJobOverrides hook) to return directly.
func (b *TLValidationJobSummaryBuilder) BuildTyped() validationjob.ValidationJobSummaryBuilder[*tslmodel.TLInfo, *tslmodel.LOTLInfo] {
	return tlValidationJobSummaryBuilderTyped{builder: b}
}

// Build builds the TLValidationJobSummary. Port of build().
//
// Java's IllegalArgumentException, raised by the TLValidationJobSummary constructor when both
// lists end up empty, becomes NewTLValidationJobSummary's returned error; this port panics with
// that error's message to keep Build()'s signature aligned with the assumed
// job.ValidationJobSummaryBuilder[D, L] contract's non-erroring Build() D (see this file's FLAG
// header) - a data-dependent condition upstream never actually allowed a caller to recover from
// either (the constructor threw an unchecked exception).
func (b *TLValidationJobSummaryBuilder) Build() *tslmodel.TLValidationJobSummary {
	var otherTLInfos []*tslmodel.TLInfo
	if utils.IsArrayNotEmpty(b.tlSources) {
		for _, tlSource := range b.tlSources {
			otherTLInfos = append(otherTLInfos, b.buildTLInfo(tlSource))
		}
	}

	var lotlList []*tslmodel.LOTLInfo
	if utils.IsArrayNotEmpty(b.lotlSources) {
		for _, lotlSource := range b.lotlSources {
			lotlParsingResult := b.readOnlyCacheAccess.GetParsingInfoRecordTyped(lotlSource.CacheKey())

			lotlInfo := b.buildLOTLInfo(lotlSource)

			var tlInfos []*tslmodel.TLInfo
			currentTLSources := b.extractTLSources(lotlParsingResult)
			for _, tlSource := range currentTLSources {
				otherTSLPointer := b.getOtherTSLPointer(lotlParsingResult.TlOtherPointers(), tlSource.Url())
				tlInfos = append(tlInfos, b.buildTLInfoWithParent(tlSource, &lotlInfo, otherTSLPointer))
			}
			lotlInfo.SetTlInfos(tlInfos)

			if lotlSource.IsPivotSupport() {
				var pivotInfos []*tslmodel.PivotInfo

				currentCertificates := b.getLOTLKeystoreCertificates(lotlSource)

				pivotSources := b.extractPivotSources(lotlParsingResult)
				for _, pivotSource := range pivotSources {
					pivotParsingCacheDTO := b.readOnlyCacheAccess.GetParsingInfoRecordTyped(pivotSource.CacheKey())
					pivotCertificateTokens := b.getPivotCertificateTokens(pivotParsingCacheDTO)
					certificateChangesMap := b.getCertificateChangesMap(pivotCertificateTokens, currentCertificates)
					associatedLOTLLocation := b.getAssociatedLOTLLocation(pivotParsingCacheDTO)
					pivotInfos = append(pivotInfos, b.buildPivotInfo(pivotSource, certificateChangesMap, associatedLOTLLocation))

					currentCertificates = pivotCertificateTokens
				}
				lotlInfo.SetPivotInfos(pivotInfos)

			} else {
				lotlInfo.SetPivotInfos(nil)
			}

			lotlList = append(lotlList, &lotlInfo)
		}
	}

	summary, err := tslmodel.NewTLValidationJobSummary(lotlList, otherTLInfos)
	if err != nil {
		panic(err)
	}
	return summary
}

func (b *TLValidationJobSummaryBuilder) buildLOTLInfo(lotlSource *LOTLSource) tslmodel.LOTLInfo {
	cacheKey := lotlSource.CacheKey()
	return tslmodel.NewLOTLInfo(
		b.readOnlyCacheAccess.GetDownloadInfoRecord(cacheKey),
		b.readOnlyCacheAccess.GetParsingInfoRecordTyped(cacheKey),
		b.readOnlyCacheAccess.GetValidationInfoRecord(cacheKey), lotlSource.Url())
}

func (b *TLValidationJobSummaryBuilder) buildTLInfo(tlSource *TLSource) *tslmodel.TLInfo {
	cacheKey := tlSource.CacheKey()
	return tslmodel.NewTLInfo(
		b.readOnlyCacheAccess.GetDownloadInfoRecord(cacheKey),
		b.readOnlyCacheAccess.GetParsingInfoRecordTyped(cacheKey),
		b.readOnlyCacheAccess.GetValidationInfoRecord(cacheKey), tlSource.Url())
}

func (b *TLValidationJobSummaryBuilder) buildTLInfoWithParent(tlSource *TLSource, lotlInfo *tslmodel.LOTLInfo, otherTSLPointer *tslmodel.OtherTSLPointer) *tslmodel.TLInfo {
	cacheKey := tlSource.CacheKey()
	return tslmodel.NewTLInfoFull(
		b.readOnlyCacheAccess.GetDownloadInfoRecord(cacheKey),
		b.readOnlyCacheAccess.GetParsingInfoRecordTyped(cacheKey),
		b.readOnlyCacheAccess.GetValidationInfoRecord(cacheKey), tlSource.Url(), lotlInfo, otherTSLPointer)
}

func (b *TLValidationJobSummaryBuilder) buildPivotInfo(pivotSource *LOTLSource, certificateChangesMap map[*model.CertificateToken]tslmodel.CertificatePivotStatus,
	associatedLOTLLocation string) *tslmodel.PivotInfo {
	cacheKey := pivotSource.CacheKey()
	return tslmodel.NewPivotInfo(
		b.readOnlyCacheAccess.GetDownloadInfoRecord(cacheKey),
		b.readOnlyCacheAccess.GetParsingInfoRecordTyped(cacheKey),
		b.readOnlyCacheAccess.GetValidationInfoRecord(cacheKey), pivotSource.Url(),
		certificateChangesMap, associatedLOTLLocation)
}

func (b *TLValidationJobSummaryBuilder) getOtherTSLPointer(tlOtherPointers []*tslmodel.OtherTSLPointer, tslPointerLocation string) *tslmodel.OtherTSLPointer {
	for _, otherTSLPointer := range tlOtherPointers {
		if tslPointerLocation == otherTSLPointer.TSLLocation() {
			return otherTSLPointer
		}
	}
	return nil
}

func (b *TLValidationJobSummaryBuilder) extractTLSources(lotlParsingResult *TLParsingCacheDTO) []*TLSource {
	var result []*TLSource
	if lotlParsingResult != nil && lotlParsingResult.IsResultExist() {
		for _, otherTSLPointerDTO := range lotlParsingResult.TlOtherPointers() {
			tlSource := NewTLSource()
			tlSource.SetUrl(otherTSLPointerDTO.TSLLocation())
			result = append(result, tlSource)
		}
	}
	return result
}

func (b *TLValidationJobSummaryBuilder) getLOTLKeystoreCertificates(lotlSource *LOTLSource) []*model.CertificateToken {
	certificateSource := lotlSource.CertificateSource()
	if certificateSource != nil {
		return certificateSource.Certificates()
	}
	return nil
}

func (b *TLValidationJobSummaryBuilder) extractPivotSources(lotlParsingResult *TLParsingCacheDTO) []*LOTLSource {
	var result []*LOTLSource
	if lotlParsingResult != nil && lotlParsingResult.IsResultExist() {
		for _, pivotUrl := range lotlParsingResult.PivotUrls() {
			pivotSource := NewLOTLSource()
			pivotSource.SetUrl(pivotUrl)
			result = append(result, pivotSource)
		}
	}
	return utils.ReverseList(result)
}

func (b *TLValidationJobSummaryBuilder) getPivotCertificateTokens(parsingCacheDTO *TLParsingCacheDTO) []*model.CertificateToken {
	lotlOtherPointers := parsingCacheDTO.LotlOtherPointers()
	if len(lotlOtherPointers) == 1 {
		return lotlOtherPointers[0].SdiCertificates()
	}
	return nil
}

func (b *TLValidationJobSummaryBuilder) getCertificateChangesMap(pivotSourceCertificates, currentCertificates []*model.CertificateToken) map[*model.CertificateToken]tslmodel.CertificatePivotStatus {
	certificateChangesMap := make(map[*model.CertificateToken]tslmodel.CertificatePivotStatus)

	var commonCertificates []*model.CertificateToken
	for _, c := range pivotSourceCertificates {
		if certificateTokenListContains(currentCertificates, c) {
			commonCertificates = append(commonCertificates, c)
		}
	}

	// added certificates
	for _, certificateToken := range pivotSourceCertificates {
		if !certificateTokenListContains(commonCertificates, certificateToken) {
			certificateChangesMap[certificateToken] = tslmodel.CertificatePivotStatus_ADDED
		}
	}

	// common certificates
	for _, certificateToken := range commonCertificates {
		certificateChangesMap[certificateToken] = tslmodel.CertificatePivotStatus_NOT_CHANGED
	}

	// removed certificates
	for _, certificateToken := range currentCertificates {
		if !certificateTokenListContains(commonCertificates, certificateToken) {
			certificateChangesMap[certificateToken] = tslmodel.CertificatePivotStatus_REMOVED
		}
	}

	return certificateChangesMap
}

func (b *TLValidationJobSummaryBuilder) getAssociatedLOTLLocation(parsingCacheDTO *TLParsingCacheDTO) string {
	xmllotlPointer := ParsingUtilsXMLLOTLPointer(parsingCacheDTO)
	if xmllotlPointer != nil {
		return xmllotlPointer.TSLLocation()
	}
	return ""
}

// certificateTokenListContains reports whether list contains token, per CertificateToken's own
// Equals (matching Java's List#contains, which delegates to equals()).
func certificateTokenListContains(list []*model.CertificateToken, token *model.CertificateToken) bool {
	for _, c := range list {
		if c.Equals(token) {
			return true
		}
	}
	return false
}
