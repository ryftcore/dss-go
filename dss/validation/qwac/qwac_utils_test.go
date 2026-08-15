package qwac

import (
	"testing"

	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

func certWithId(id string) *diagnostic.CertificateWrapper {
	cid := jaxb.CollapsedString(id)
	return diagnostic.NewCertificateWrapper(&jaxb.XmlCertificate{
		XmlAbstractTokenAttrs: jaxb.XmlAbstractTokenAttrs{Id: &cid},
	})
}

func sigWithDigestMatchers(matchers ...*jaxb.XmlDigestMatcher) *diagnostic.SignatureWrapper {
	return diagnostic.NewSignatureWrapper(&jaxb.XmlSignature{
		DigestMatchers: &jaxb.DigestMatchersWrapper{Items: matchers},
	})
}

func sigDEntryFor(documentName string, found, intact bool) *jaxb.XmlDigestMatcher {
	t := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_SIG_D_ENTRY)
	name := documentName
	return &jaxb.XmlDigestMatcher{
		DataFound:    found,
		DataIntact:   intact,
		Type:         &t,
		DocumentName: &name,
	}
}

func TestGetIdentifiedTLSCertificates_MatchesFoundIntactSigDEntry(t *testing.T) {
	cert := certWithId("C-1")
	sig := sigWithDigestMatchers(sigDEntryFor("C-1", true, true))

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert})
	if len(got) != 1 || got[0] != cert {
		t.Fatalf("expected [cert], got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_NoMatchOnIdMismatch(t *testing.T) {
	cert := certWithId("C-1")
	sig := sigWithDigestMatchers(sigDEntryFor("C-OTHER", true, true))

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert})
	if len(got) != 0 {
		t.Fatalf("expected no matches, got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_SkipsWhenNotFoundOrNotIntact(t *testing.T) {
	cert := certWithId("C-1")

	notFound := sigWithDigestMatchers(sigDEntryFor("C-1", false, true))
	if got := GetIdentifiedTLSCertificates(notFound, []*diagnostic.CertificateWrapper{cert}); len(got) != 0 {
		t.Fatalf("expected no matches when DataFound=false, got %v", got)
	}

	notIntact := sigWithDigestMatchers(sigDEntryFor("C-1", true, false))
	if got := GetIdentifiedTLSCertificates(notIntact, []*diagnostic.CertificateWrapper{cert}); len(got) != 0 {
		t.Fatalf("expected no matches when DataIntact=false, got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_SkipsWrongType(t *testing.T) {
	cert := certWithId("C-1")
	other := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_MESSAGE_DIGEST)
	name := "C-1"
	sig := sigWithDigestMatchers(&jaxb.XmlDigestMatcher{
		DataFound: true, DataIntact: true, Type: &other, DocumentName: &name,
	})

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert})
	if len(got) != 0 {
		t.Fatalf("expected no matches for non-SIG_D_ENTRY type, got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_SkipsNilType(t *testing.T) {
	cert := certWithId("C-1")
	name := "C-1"
	sig := sigWithDigestMatchers(&jaxb.XmlDigestMatcher{
		DataFound: true, DataIntact: true, DocumentName: &name,
	})

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert})
	if len(got) != 0 {
		t.Fatalf("expected no matches for nil type, got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_SkipsNilDocumentName(t *testing.T) {
	cert := certWithId("C-1")
	tp := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_SIG_D_ENTRY)
	sig := sigWithDigestMatchers(&jaxb.XmlDigestMatcher{
		DataFound: true, DataIntact: true, Type: &tp,
	})

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert})
	if len(got) != 0 {
		t.Fatalf("expected no matches for nil DocumentName, got %v", got)
	}
}

func TestGetIdentifiedTLSCertificates_MultipleMatchersAndCandidates(t *testing.T) {
	cert1 := certWithId("C-1")
	cert2 := certWithId("C-2")
	sig := sigWithDigestMatchers(
		sigDEntryFor("C-2", true, true),
		sigDEntryFor("C-1", true, true),
	)

	got := GetIdentifiedTLSCertificates(sig, []*diagnostic.CertificateWrapper{cert1, cert2})
	if len(got) != 2 || got[0] != cert2 || got[1] != cert1 {
		t.Fatalf("expected [cert2, cert1] in digest-matcher order, got %v", got)
	}
}
