package diagnostic

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

func timestampedObjectTypeP(v enumerations.TimestampedObjectType) *jaxb.TimestampedObjectTypeValue {
	tv := jaxb.TimestampedObjectTypeValue(v)
	return &tv
}

func TestTimestampWrapper_TimestampedObjectCategories(t *testing.T) {
	sig := &jaxb.XmlSignature{}
	sig.Id = jaxb.NewCollapsedString("sig-1")

	cert := &jaxb.XmlCertificate{}
	cert.Id = jaxb.NewCollapsedString("cert-1")

	rev := &jaxb.XmlRevocation{}
	rev.Id = jaxb.NewCollapsedString("rev-1")

	ts := &jaxb.XmlTimestamp{
		TimestampedObjects: &jaxb.TimestampedObjectsWrapper{
			Items: []*jaxb.XmlTimestampedObject{
				{Token: jaxb.NewXmlTokenRef(sig), Category: timestampedObjectTypeP(enumerations.TimestampedObjectTypeSignature)},
				{Token: jaxb.NewXmlTokenRef(cert), Category: timestampedObjectTypeP(enumerations.TimestampedObjectTypeCertificate)},
				{Token: jaxb.NewXmlTokenRef(rev), Category: timestampedObjectTypeP(enumerations.TimestampedObjectTypeRevocation)},
			},
		},
	}
	ts.Id = jaxb.NewCollapsedString("ts-1")

	w := NewTimestampWrapper(ts)

	sigs := w.TimestampedSignatures()
	if len(sigs) != 1 || sigs[0].Id() != "sig-1" {
		t.Fatalf("unexpected TimestampedSignatures: %+v", sigs)
	}

	certs := w.TimestampedCertificates()
	if len(certs) != 1 || certs[0].Id() != "cert-1" {
		t.Fatalf("unexpected TimestampedCertificates: %+v", certs)
	}

	revs := w.TimestampedRevocations()
	if len(revs) != 1 || revs[0].Id() != "rev-1" {
		t.Fatalf("unexpected TimestampedRevocations: %+v", revs)
	}

	// No timestamps/evidence records/signed data/orphan tokens registered.
	if got := w.TimestampedTimestamps(); got != nil {
		t.Fatalf("expected no timestamped timestamps, got %+v", got)
	}
	if got := w.AllTimestampedOrphanTokens(); got != nil {
		t.Fatalf("expected no orphan tokens, got %+v", got)
	}
}

func TestTimestampWrapper_TimestampedObjectsWrongTypePanics(t *testing.T) {
	cert := &jaxb.XmlCertificate{}
	cert.Id = jaxb.NewCollapsedString("cert-1")

	ts := &jaxb.XmlTimestamp{
		TimestampedObjects: &jaxb.TimestampedObjectsWrapper{
			Items: []*jaxb.XmlTimestampedObject{
				// Mismatched category vs actual token type - Java throws
				// IllegalArgumentException here; the Go port panics.
				{Token: jaxb.NewXmlTokenRef(cert), Category: timestampedObjectTypeP(enumerations.TimestampedObjectTypeSignature)},
			},
		},
	}
	ts.Id = jaxb.NewCollapsedString("ts-1")
	w := NewTimestampWrapper(ts)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected a panic for a category/type mismatch")
		}
	}()
	w.TimestampedSignatures()
}

func TestTimestampWrapper_TypeAndArchiveType(t *testing.T) {
	tsType := enumerations.TimestampTypeArchiveTimestamp
	tsTypeValue := jaxb.TimestampTypeValue(tsType)
	archiveType := enumerations.ArchiveTimestampTypeCAdESV3
	archiveTypeValue := jaxb.ArchiveTimestampTypeValue(archiveType)

	ts := &jaxb.XmlTimestamp{
		Type:                 &tsTypeValue,
		ArchiveTimestampType: &archiveTypeValue,
	}
	w := NewTimestampWrapper(ts)
	if w.Type() != enumerations.TimestampTypeArchiveTimestamp {
		t.Fatalf("unexpected Type(): %v", w.Type())
	}
	if w.ArchiveTimestampType() != enumerations.ArchiveTimestampTypeCAdESV3 {
		t.Fatalf("unexpected ArchiveTimestampType(): %v", w.ArchiveTimestampType())
	}
}

func TestTimestampWrapper_MessageImprint(t *testing.T) {
	digestMatcherType := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherTypeMessageImprint)
	dataFound := true
	ts := &jaxb.XmlTimestamp{
		DigestMatcher: []*jaxb.XmlDigestMatcher{
			{Type: &digestMatcherType, DataFound: dataFound, DataIntact: true},
		},
	}
	w := NewTimestampWrapper(ts)
	mi := w.MessageImprint()
	if mi == nil || !mi.DataFound || !mi.DataIntact {
		t.Fatalf("unexpected MessageImprint(): %+v", mi)
	}
	if !w.IsMessageImprintDataFound() || !w.IsMessageImprintDataIntact() {
		t.Fatalf("expected message imprint found/intact true")
	}
}
