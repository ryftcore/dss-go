// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/PivotProcessingResultFromCacheAccessBuilder.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
)

// PivotProcessingResultFromCacheAccessBuilder creates a PivotProcessingResult from a given
// TLCacheAccessByKey.
type PivotProcessingResultFromCacheAccessBuilder struct {
	// cacheAccessByKey is the cache access to the given pivot.
	cacheAccessByKey *TLCacheAccessByKey
}

// NewPivotProcessingResultFromCacheAccessBuilder is the default constructor.
func NewPivotProcessingResultFromCacheAccessBuilder(cacheAccessByKey *TLCacheAccessByKey) *PivotProcessingResultFromCacheAccessBuilder {
	return &PivotProcessingResultFromCacheAccessBuilder{cacheAccessByKey: cacheAccessByKey}
}

// Build builds the PivotProcessingResult. Port of build().
func (b *PivotProcessingResultFromCacheAccessBuilder) Build() *PivotProcessingResult {
	parsingCacheEntry, _ := b.cacheAccessByKey.GetParsingReadOnlyResult().(*TLParsingCacheDTO)
	xmlLotlPointer := ParsingUtilsXMLLOTLPointer(parsingCacheEntry)
	return NewPivotProcessingResult(b.getDocument(), b.getCertificateSource(xmlLotlPointer), b.getLotlLocation(xmlLotlPointer))
}

func (b *PivotProcessingResultFromCacheAccessBuilder) getDocument() model.DSSDocument {
	downloadResult := b.cacheAccessByKey.GetDownloadReadOnlyResult()
	if downloadResult != nil && downloadResult.IsResultExist() {
		return downloadResult.Document()
	}
	return nil
}

func (b *PivotProcessingResultFromCacheAccessBuilder) getCertificateSource(xmlLotlPointer *tslmodel.OtherTSLPointer) spi.CertificateSource {
	if xmlLotlPointer != nil {
		return ParsingUtilsLOTLAnnouncedCertificateSource(xmlLotlPointer)
	}
	return nil
}

func (b *PivotProcessingResultFromCacheAccessBuilder) getLotlLocation(xmlLotlPointer *tslmodel.OtherTSLPointer) string {
	if xmlLotlPointer != nil {
		return xmlLotlPointer.TSLLocation()
	}
	return ""
}
