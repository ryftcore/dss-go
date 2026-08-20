// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/alerts/detections/OJUrlChangeDetection.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/alert"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// OJUrlChangeDetection detects a change of the OJ URL.
type OJUrlChangeDetection struct {
	// lotlSource is the LOTL source.
	lotlSource *LOTLSource
}

var _ alert.AlertDetector[*tslmodel.LOTLInfo] = (*OJUrlChangeDetection)(nil)

// NewOJUrlChangeDetection is the default constructor.
func NewOJUrlChangeDetection(lotlSource *LOTLSource) *OJUrlChangeDetection {
	return &OJUrlChangeDetection{lotlSource: lotlSource}
}

// Detect ports detect(LOTLInfo).
func (d *OJUrlChangeDetection) Detect(info *tslmodel.LOTLInfo) bool {
	if d.lotlSource.Url() != info.Url() {
		return false
	}

	parsingCacheInfo, ok := info.TLParsingCacheInfo()
	if ok && parsingCacheInfo.IsDesynchronized() {
		signingCertificatesAnnouncementPredicate := d.lotlSource.SigningCertificatesAnnouncementPredicate()
		// JUDGMENT CALL: Java's instanceof OfficialJournalSchemeInformationURI check always
		// holds here - see lotl_signing_certificates_announcement_scheme_information_uri.go's
		// header, which collapses that interface and its only implementer into one type.
		if signingCertificatesAnnouncementPredicate != nil {
			officialJournalURL := signingCertificatesAnnouncementPredicate.Uri()
			signingCertificateAnnouncementUrl := parsingCacheInfo.SigningCertificateAnnouncementUrl()

			if officialJournalURL != signingCertificateAnnouncementUrl {
				return true
			}
		}
	}

	return false
}
