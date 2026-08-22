// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/PivotProcessingResult.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PivotProcessingResult contains the pivot and its introduced signing certificates for the LOTL
// or the next pivot.
type PivotProcessingResult struct {
	// pivot is the pivot document.
	pivot model.DSSDocument

	// certificateSource is the certificate source to use.
	certificateSource spi.CertificateSource

	// lotlLocation is the LOTL location.
	lotlLocation string
}

// NewPivotProcessingResult is the default constructor.
func NewPivotProcessingResult(pivot model.DSSDocument, certificateSource spi.CertificateSource, lotlLocation string) *PivotProcessingResult {
	return &PivotProcessingResult{pivot: pivot, certificateSource: certificateSource, lotlLocation: lotlLocation}
}

// Pivot gets the pivot document. Port of getPivot().
func (r *PivotProcessingResult) Pivot() model.DSSDocument {
	return r.pivot
}

// CertificateSource gets the certificate source. Port of getCertificateSource().
func (r *PivotProcessingResult) CertificateSource() spi.CertificateSource {
	return r.certificateSource
}

// LotlLocation gets the LOTL location. Port of getLotlLocation().
func (r *PivotProcessingResult) LotlLocation() string {
	return r.lotlLocation
}
