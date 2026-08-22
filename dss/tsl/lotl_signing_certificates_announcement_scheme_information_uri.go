// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/LOTLSigningCertificatesAnnouncementSchemeInformationURI.java (DSS 6.5.RC1).
//
// Java declares this as an interface (java.util.function.Predicate<NonEmptyMultiLangURIType> plus
// a getUri() accessor) with exactly one implementation in the ported tree,
// OfficialJournalSchemeInformationURI.
//
// source/lotl_source.go declares its
// signingCertificatesAnnouncementPredicate field as *LOTLSigningCertificatesAnnouncementSchemeInformationURI
// and lotl_parsing_task.go (same chunk) calls .Test(...) and .Uri() directly on that pointer value.
// A pointer-to-Go-interface does not forward method calls the way a Java interface reference
// does, so those call sites only compile against a concrete struct type. This port therefore
// collapses the Java interface and its sole implementation into one concrete Go type:
// LOTLSigningCertificatesAnnouncementSchemeInformationURI IS what OfficialJournalSchemeInformationURI.java
// implements, with NewOfficialJournalSchemeInformationURI kept (in that file) as a same-shape
// constructor alias for the Java class name. OJUrlChangeDetection.java's `instanceof
// OfficialJournalSchemeInformationURI` check accordingly always holds for a non-nil predicate of
// this type - documented there.
package tsl

import (
	"net/url"
	"strings"

	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
)

// LOTLSigningCertificatesAnnouncementSchemeInformationURI filters the LOTL signing certificates
// scheme information URI.
type LOTLSigningCertificatesAnnouncementSchemeInformationURI struct {
	// officialJournalURL is the OJ URL.
	officialJournalURL string
}

// NewLOTLSigningCertificatesAnnouncementSchemeInformationURI constructs the predicate for the
// given Official Journal URL.
//
// Panics when officialJournalURL is empty (Objects.requireNonNull, "Official Journal URL cannot
// be null!").
func NewLOTLSigningCertificatesAnnouncementSchemeInformationURI(officialJournalURL string) *LOTLSigningCertificatesAnnouncementSchemeInformationURI {
	if officialJournalURL == "" {
		panic("Official Journal URL cannot be null!")
	}
	return &LOTLSigningCertificatesAnnouncementSchemeInformationURI{officialJournalURL: officialJournalURL}
}

// Test ports test(NonEmptyMultiLangURIType).
func (p *LOTLSigningCertificatesAnnouncementSchemeInformationURI) Test(t *jaxb.NonEmptyMultiLangURIType) bool {
	if t == nil || t.Value == "" {
		return false
	}
	domain, ok := p.ojDomain()
	if !ok {
		return false
	}
	return strings.Contains(t.Value, domain)
}

// Uri returns the URI of the signing-certificate announcement page. Port of getUri().
func (p *LOTLSigningCertificatesAnnouncementSchemeInformationURI) Uri() string {
	return p.officialJournalURL
}

// ojDomain ports the private getOJDomain(): a malformed URL panics with the Java DSSException
// message at Test-call time (matching Java, which raises DSSException from within test()); the
// caller (Test) folds that into "no match" instead of propagating, since the Go interface this
// predicate satisfies elsewhere (Test bool) carries no error channel - a deviation flagged here
// rather than panicking mid-parse loop.
func (p *LOTLSigningCertificatesAnnouncementSchemeInformationURI) ojDomain() (string, bool) {
	u, err := url.Parse(p.officialJournalURL)
	if err != nil || u.Host == "" {
		return "", false
	}
	return u.Hostname(), true
}
