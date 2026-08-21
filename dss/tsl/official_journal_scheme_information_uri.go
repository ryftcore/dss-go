// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/OfficialJournalSchemeInformationURI.java (DSS 6.5.RC1).
//
// See lotl_signing_certificates_announcement_scheme_information_uri.go's header: Java's
// OfficialJournalSchemeInformationURI implements LOTLSigningCertificatesAnnouncementSchemeInformationURI
// and is that interface's only implementation in the ported tree, so this port collapses the two
// into one concrete type and keeps this constructor as a same-shape alias for the Java class name.
package tsl

// NewOfficialJournalSchemeInformationURI is the default constructor. Port of
// OfficialJournalSchemeInformationURI(String).
func NewOfficialJournalSchemeInformationURI(officialJournalURL string) *LOTLSigningCertificatesAnnouncementSchemeInformationURI {
	return NewLOTLSigningCertificatesAnnouncementSchemeInformationURI(officialJournalURL)
}
