// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/LOTLParsingResult.java (DSS 6.5.RC1).
package tsl

import tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"

// LOTLParsingResult is a parsed LOTL result.
type LOTLParsingResult struct {
	AbstractTLParsingResult

	// lotlPointers is the list of LOTL pointers.
	lotlPointers []*tslmodel.OtherTSLPointer

	// tlPointers is the list of TL pointers.
	tlPointers []*tslmodel.OtherTSLPointer

	// signingCertificateAnnouncementURL is the signing certificate announcement URL.
	signingCertificateAnnouncementURL string

	// pivotURLs is the list of pivot URLs.
	pivotURLs []string
}

// NewLOTLParsingResult is the default constructor. Port of LOTLParsingResult().
func NewLOTLParsingResult() *LOTLParsingResult {
	return &LOTLParsingResult{AbstractTLParsingResult: NewAbstractTLParsingResult()}
}

// LotlPointers gets the LOTL other TSL pointers. Port of getLotlPointers().
func (r *LOTLParsingResult) LotlPointers() []*tslmodel.OtherTSLPointer {
	return r.lotlPointers
}

// SetLotlPointers sets the LOTL other pointers. Port of setLotlPointers(List).
func (r *LOTLParsingResult) SetLotlPointers(lotlPointers []*tslmodel.OtherTSLPointer) {
	r.lotlPointers = lotlPointers
}

// TlPointers gets the TL other TSL pointers. Port of getTlPointers().
func (r *LOTLParsingResult) TlPointers() []*tslmodel.OtherTSLPointer {
	return r.tlPointers
}

// SetTlPointers sets the TL other pointers. Port of setTlPointers(List).
func (r *LOTLParsingResult) SetTlPointers(tlPointers []*tslmodel.OtherTSLPointer) {
	r.tlPointers = tlPointers
}

// SigningCertificateAnnouncementURL gets the signing certificate announcement URL. Port of
// getSigningCertificateAnnouncementURL().
func (r *LOTLParsingResult) SigningCertificateAnnouncementURL() string {
	return r.signingCertificateAnnouncementURL
}

// SetSigningCertificateAnnouncementURL sets the signing certificate announcement URL. Port of
// setSigningCertificateAnnouncementURL(String).
func (r *LOTLParsingResult) SetSigningCertificateAnnouncementURL(signingCertificateAnnouncementURL string) {
	r.signingCertificateAnnouncementURL = signingCertificateAnnouncementURL
}

// PivotURLs gets the pivot URLs. Port of getPivotURLs().
func (r *LOTLParsingResult) PivotURLs() []string {
	return r.pivotURLs
}

// SetPivotURLs sets the pivot URLs. Port of setPivotURLs(List).
func (r *LOTLParsingResult) SetPivotURLs(pivotURLs []string) {
	r.pivotURLs = pivotURLs
}
