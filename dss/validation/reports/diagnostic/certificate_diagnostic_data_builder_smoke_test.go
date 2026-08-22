package diagnostic

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// certificateExtensionsSmokeUnparseable mirrors spi.certificateExtensionsKATUnparseable: these
// DER files crypto/x509 rejects but the corpus still ships, so they are skipped rather than
// failing the smoke test.
var certificateExtensionsSmokeUnparseable = map[string]bool{
	"cert_16.der": true,
	"cert_19.der": true,
}

// loadSmokeCertificates parses every certificate in spi/testdata/certificate_extensions,
// reusing the corpus the frozen spi package's own KAT test already exercises (see
// certificate_extensions_utils_kat_test.go).
func loadSmokeCertificates(t *testing.T) []*model.CertificateToken {
	t.Helper()
	dir := filepath.Join("..", "..", "..", "spi", "testdata", "certificate_extensions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read testdata dir: %v", err)
	}
	var tokens []*model.CertificateToken
	for _, entry := range entries {
		name := entry.Name()
		if filepath.Ext(name) != ".der" || certificateExtensionsSmokeUnparseable[name] {
			continue
		}
		der, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		token, err := model.NewCertificateToken(cert)
		if err != nil {
			t.Fatalf("build token for %s: %v", name, err)
		}
		tokens = append(tokens, token)
	}
	if len(tokens) == 0 {
		t.Fatal("no certificates loaded from testdata corpus")
	}
	return tokens
}

// TestCertificateDiagnosticDataBuilderSmoke exercises BuildDetachedXmlCertificate (and, through
// it, every certificate-extension builder in diagnostic_data_builder.go) over the real
// certificate corpus spi's own KAT test uses, verifying the builder never panics and that the
// resulting XmlDiagnosticData round-trips through Marshal without an error - a structural smoke
// test standing in for a full byte-compare oracle harness, which does not exist yet for this
// builder.
func TestCertificateDiagnosticDataBuilderSmoke(t *testing.T) {
	tokens := loadSmokeCertificates(t)

	builder := NewCertificateDiagnosticDataBuilder()
	builder.UsedCertificates(tokens)
	builder.AllCertificateSources(spi.NewListCertificateSource())

	dd := builder.Build()
	if dd == nil {
		t.Fatal("Build returned nil")
	}
	if dd.UsedCertificates == nil || len(dd.UsedCertificates.All()) != len(tokens) {
		t.Fatalf("expected %d used certificates, got %v", len(tokens), dd.UsedCertificates)
	}

	out, err := jaxb.Marshal(dd)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("Marshal produced no output")
	}

	// Round-trip: unmarshal what we just marshaled and marshal again; the two byte slices must
	// be identical (idempotent marshaling), catching wrapper/order mistakes even without a Java
	// oracle to diff against.
	dd2, err := jaxb.Unmarshal(out)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	out2, err := jaxb.Marshal(dd2)
	if err != nil {
		t.Fatalf("re-Marshal: %v", err)
	}
	if string(out) != string(out2) {
		t.Fatalf("marshal is not idempotent through a round-trip:\n--- first ---\n%s\n--- second ---\n%s", out, out2)
	}
}

// TestCertificateDiagnosticDataBuilderSmokeSingle drills into a single certificate's
// BuildDetachedXmlCertificate output to catch obviously wrong field mappings (id, subject/
// issuer DN entries, self-signed flag) that an idempotent-marshal check alone would not.
func TestCertificateDiagnosticDataBuilderSmokeSingle(t *testing.T) {
	tokens := loadSmokeCertificates(t)
	token := tokens[0]

	builder := NewCertificateDiagnosticDataBuilder()
	builder.UsedCertificates([]*model.CertificateToken{token})
	builder.AllCertificateSources(spi.NewListCertificateSource())
	builder.UsedRevocations(nil)

	xmlCert := builder.BuildDetachedXmlCertificate(token)
	if xmlCert.Id == nil || xmlCert.Id.String() == "" {
		t.Fatal("Id not set")
	}
	if len(xmlCert.SubjectDistinguishedName) != 2 {
		t.Fatalf("expected 2 subject DN entries (CANONICAL+RFC2253), got %d", len(xmlCert.SubjectDistinguishedName))
	}
	if len(xmlCert.IssuerDistinguishedName) != 2 {
		t.Fatalf("expected 2 issuer DN entries (CANONICAL+RFC2253), got %d", len(xmlCert.IssuerDistinguishedName))
	}
	if xmlCert.SelfSigned != token.IsSelfSigned() {
		t.Fatalf("SelfSigned = %v, want %v", xmlCert.SelfSigned, token.IsSelfSigned())
	}
	if xmlCert.EntityKey == nil || *xmlCert.EntityKey != token.EntityKey().AsXmlID() {
		t.Fatalf("EntityKey = %v, want %v", xmlCert.EntityKey, token.EntityKey().AsXmlID())
	}
}

// TestSignedDocumentDiagnosticDataBuilderSmokeAssertConfigurationValid exercises the
// requireNonNull-derived panic path (Java's assertConfigurationValid()), matching the "requireNonNull
// -> panic(Java message)" PORTING.md rule.
func TestSignedDocumentDiagnosticDataBuilderSmokeAssertConfigurationValid(t *testing.T) {
	builder := NewSignedDocumentDiagnosticDataBuilder()
	builder.AllCertificateSources(spi.NewListCertificateSource())

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic when signedDocument is not set")
		}
	}()
	builder.Build()
}

// TestQWACCertificateDiagnosticDataBuilderSmokeAssertConfigurationValid exercises the QWAC
// override of assertConfigurationValid(), confirming the SignedDocumentDiagnosticDataBuilder
// overrides-interface dispatch (InitSignedDocumentDiagnosticDataBuilder) actually reaches the
// subclass override rather than the base's default.
func TestQWACCertificateDiagnosticDataBuilderSmokeAssertConfigurationValid(t *testing.T) {
	builder := NewQWACCertificateDiagnosticDataBuilder()
	builder.AllCertificateSources(spi.NewListCertificateSource())
	// signedDocument left nil: the base's own AssertConfigurationValid would panic on that
	// first; set a stub document so only the QWAC override's websiteUrl check is exercised.
	builder.Document(model.NewInMemoryDocument([]byte("stub")))

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic when websiteUrl is not set")
		}
	}()
	builder.Build()
}
