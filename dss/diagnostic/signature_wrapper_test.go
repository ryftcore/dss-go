package diagnostic

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// Compile-time checks that every DIAGWRAP_B concrete wrapper satisfies TokenProxy, the way
// AbstractTokenProxyBase's InitTokenProxy contract requires.
var (
	_ TokenProxy = (*SignatureWrapper)(nil)
	_ TokenProxy = (*TimestampWrapper)(nil)
	_ TokenProxy = (*RevocationWrapper)(nil)
)

func TestSignatureWrapper_TimestampListByType(t *testing.T) {
	sigTsType := jaxb.TimestampTypeValue(enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	archiveTsType := jaxb.TimestampTypeValue(enumerations.TimestampType_ARCHIVE_TIMESTAMP)

	sigTs := &jaxb.XmlTimestamp{Type: &sigTsType}
	sigTs.Id = jaxb.NewCollapsedString("sig-ts")
	archiveTs := &jaxb.XmlTimestamp{Type: &archiveTsType}
	archiveTs.Id = jaxb.NewCollapsedString("archive-ts")

	sig := &jaxb.XmlSignature{
		FoundTimestamps: &jaxb.FoundTimestampsWrapper{
			Items: []*jaxb.XmlFoundTimestamp{
				{Timestamp: sigTs},
				{Timestamp: archiveTs},
			},
		},
	}
	sig.Id = jaxb.NewCollapsedString("sig-1")

	w := NewSignatureWrapper(sig)

	all := w.TimestampList()
	if len(all) != 2 {
		t.Fatalf("expected 2 timestamps, got %d", len(all))
	}

	sigLevel := w.SignatureTimestamps()
	if len(sigLevel) != 1 || sigLevel[0].Id() != "sig-ts" {
		t.Fatalf("unexpected SignatureTimestamps(): %+v", sigLevel)
	}

	archiveLevel := w.ArchiveTimestamps()
	if len(archiveLevel) != 1 || archiveLevel[0].Id() != "archive-ts" {
		t.Fatalf("unexpected ArchiveTimestamps(): %+v", archiveLevel)
	}

	if !w.IsThereTLevel() {
		t.Fatalf("expected IsThereTLevel true")
	}
	if !w.IsThereALevel() {
		t.Fatalf("expected IsThereALevel true")
	}
	if w.IsThereXLevel() {
		t.Fatalf("expected IsThereXLevel false")
	}
}

func TestSignatureWrapper_MessageDigest(t *testing.T) {
	msgDigestType := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_MESSAGE_DIGEST)
	otherType := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_REFERENCE)
	sig := &jaxb.XmlSignature{
		DigestMatchers: &jaxb.DigestMatchersWrapper{
			Items: []*jaxb.XmlDigestMatcher{
				{Type: &otherType},
				{Type: &msgDigestType, DataFound: true, DataIntact: true},
			},
		},
	}
	w := NewSignatureWrapper(sig)
	md := w.MessageDigest()
	if md == nil || !md.DataFound || !md.DataIntact {
		t.Fatalf("unexpected MessageDigest(): %+v", md)
	}
}

func TestRevocationWrapper_Basics(t *testing.T) {
	origin := jaxb.RevocationOriginValue(enumerations.RevocationOrigin_INPUT_DOCUMENT)
	rev := &jaxb.XmlRevocation{Origin: &origin}
	rev.Id = jaxb.NewCollapsedString("rev-1")

	w := NewRevocationWrapper(rev)
	if w.Id() != "rev-1" {
		t.Fatalf("unexpected Id(): %q", w.Id())
	}
	if w.Origin() != enumerations.RevocationOrigin_INPUT_DOCUMENT {
		t.Fatalf("unexpected Origin(): %v", w.Origin())
	}
	if !w.IsInternalRevocationOrigin() {
		t.Fatalf("expected IsInternalRevocationOrigin true for INPUT_DOCUMENT origin")
	}

	other := NewRevocationWrapper(&jaxb.XmlRevocation{})
	other.revocation.Id = jaxb.NewCollapsedString("rev-1")
	if !w.Equals(other) {
		t.Fatalf("expected two RevocationWrapper values with the same Id to be Equals")
	}
}

func TestSignerDataWrapper_Equals(t *testing.T) {
	a := NewSignerDataWrapper(&jaxb.XmlSignerData{})
	a.signerData.Id = jaxb.NewCollapsedString("sd-1")
	b := NewSignerDataWrapper(&jaxb.XmlSignerData{})
	b.signerData.Id = jaxb.NewCollapsedString("sd-1")
	if !a.Equals(b) {
		t.Fatalf("expected a.Equals(b) true")
	}
	if a.String() != "SignerData Id='sd-1'" {
		t.Fatalf("unexpected String(): %q", a.String())
	}
}
