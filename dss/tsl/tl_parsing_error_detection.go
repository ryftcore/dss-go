// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/detections/TLParsingErrorDetection.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/alert"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TLParsingErrorDetection detects an error on TL parsing or structure validation.
type TLParsingErrorDetection struct{}

var _ alert.AlertDetector[*tslmodel.TLInfo] = (*TLParsingErrorDetection)(nil)

// NewTLParsingErrorDetection is the default constructor.
func NewTLParsingErrorDetection() *TLParsingErrorDetection {
	return &TLParsingErrorDetection{}
}

// Detect ports detect(TLInfo).
func (d *TLParsingErrorDetection) Detect(info *tslmodel.TLInfo) bool {
	parsingCacheInfo := info.ParsingCacheInfo()
	return parsingCacheInfo != nil && (parsingCacheInfo.IsError() || utils.IsCollectionNotEmpty(parsingCacheInfo.StructureValidationMessages()))
}
