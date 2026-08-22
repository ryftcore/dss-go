// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/PAdESProfileParameters.java
// (DSS 6.5.RC1).
//
// This class is used to accelerate signature creation process for PAdES. The cache is set within
// PAdESService#getDataToSign(...) method and used in PAdESService#signDocument(...) method.
//
// Java's `extends ProfileParameters` becomes embedding.
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
)

// ProfileParameters is used to accelerate the signature creation process for PAdES.
type ProfileParameters struct {
	document.ProfileParameters

	// pdfToBeSignedCache is the internal cache used to accelerate the signature creation
	// process.
	pdfToBeSignedCache *PdfSignatureCache
}

// NewProfileParameters is the default constructor.
func NewProfileParameters() *ProfileParameters {
	return &ProfileParameters{ProfileParameters: *document.NewProfileParameters()}
}

// PdfToBeSignedCache gets the PDF signature cache, lazily instantiating it. Port of
// #getPdfToBeSignedCache.
func (p *ProfileParameters) PdfToBeSignedCache() *PdfSignatureCache {
	if p.pdfToBeSignedCache == nil {
		p.pdfToBeSignedCache = NewPdfSignatureCache()
	}
	return p.pdfToBeSignedCache
}

// SetPdfToBeSignedCache sets the PDF signature cache. Port of #setPdfToBeSignedCache.
func (p *ProfileParameters) SetPdfToBeSignedCache(pdfToBeSignedCache *PdfSignatureCache) {
	p.pdfToBeSignedCache = pdfToBeSignedCache
}

// String ports #toString.
func (p *ProfileParameters) String() string {
	return fmt.Sprintf("PAdESProfileParameters [pdfToBeSignedCache=%v] %s", p.pdfToBeSignedCache, p.ProfileParameters.String())
}

// Equals ports #equals.
func (p *ProfileParameters) Equals(other *ProfileParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.ProfileParameters.Equals(&other.ProfileParameters) {
		return false
	}
	return p.pdfToBeSignedCache == other.pdfToBeSignedCache
}
