// Known-answer tests for CMSCertificateSource, CMSCRLSource and CMSOCSPSource.
//
// The fixture testdata/cmssrc/cades-full.p7s and the expectation
// testdata/cmssrc/kat_cms_sources.txt are produced by testdata/cmssrc/generator/CmsSourceKatGen.java,
// which builds the CMS document with BouncyCastle 1.84 and then dumps what the REAL DSS
// 6.5.RC1 sources extract from it: the expectations below are upstream's own output, not a
// hand-written guess. See that file's header for how to regenerate them.
package spi

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// cmsSourcesKatFixture parses the fixture and returns the CMS plus its single signer.
func cmsSourcesKatFixture(t *testing.T) (*cmscore.CMS, *cmscore.SignerInfo) {
	t.Helper()
	encoded, err := os.ReadFile(corpustest.Path(t, "cmssrc/cades-full.p7s"))
	if err != nil {
		t.Fatalf("reading the fixture: %v", err)
	}
	cms, err := cmscore.ParseCMS(encoded)
	if err != nil {
		t.Fatalf("parsing the fixture: %v", err)
	}
	if len(cms.SignerInfos()) != 1 {
		t.Fatalf("the fixture holds %d signers, expected 1", len(cms.SignerInfos()))
	}
	return cms, cms.SignerInfos()[0]
}

// cmsSourcesKatExpectations reads the key|value expectation file.
func cmsSourcesKatExpectations(t *testing.T) map[string]string {
	t.Helper()
	content, err := os.ReadFile("testdata/cmssrc/kat_cms_sources.txt")
	if err != nil {
		t.Fatalf("reading the expectations: %v", err)
	}
	expectations := make(map[string]string)
	for _, line := range strings.Split(strings.TrimRight(string(content), "\n"), "\n") {
		key, value, found := strings.Cut(line, "|")
		if !found {
			t.Fatalf("malformed expectation line %q", line)
		}
		expectations[key] = value
	}
	return expectations
}

// cmsSourcesKatCheck compares one produced value against its expectation.
func cmsSourcesKatCheck(t *testing.T, expectations map[string]string, key, produced string) {
	t.Helper()
	expected, found := expectations[key]
	if !found {
		t.Fatalf("no expectation for %q", key)
	}
	if produced != expected {
		t.Errorf("%s:\n produced %q\n expected %q", key, produced, expected)
	}
}

// cmsSourcesKatJoin sorts and joins the values the way the generator's putList does.
func cmsSourcesKatJoin(values []string) string {
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	// The generator writes a java.util.TreeSet, which deduplicates.
	deduplicated := make([]string, 0, len(sorted))
	for index, value := range sorted {
		if index == 0 || value != sorted[index-1] {
			deduplicated = append(deduplicated, value)
		}
	}
	return strings.Join(deduplicated, ",")
}

// cmsSourcesKatSignerIdentifier formats a SignerIdentifier exactly as the generator does.
func cmsSourcesKatSignerIdentifier(identifier *SignerIdentifier) string {
	if identifier == nil {
		return "<null>"
	}
	issuerName := ""
	if identifier.IssuerName() != nil {
		issuerName = identifier.IssuerName().RFC2253Name()
	}
	serialNumber := ""
	if identifier.SerialNumber() != nil {
		serialNumber = identifier.SerialNumber().String()
	}
	return fmt.Sprintf("%s;%s;%s;%t", issuerName, serialNumber,
		hex.EncodeToString(identifier.Ski()), identifier.IsCurrent())
}

func TestCMSCertificateSourceKAT(t *testing.T) {
	cms, signer := cmsSourcesKatFixture(t)
	expectations := cmsSourcesKatExpectations(t)

	source, err := NewCMSCertificateSource(cms.SignerInfos(), cms.Certificates(), signer)
	if err != nil {
		t.Fatalf("NewCMSCertificateSource: %v", err)
	}

	tokenIDs := func(tokens []*model.CertificateToken) []string {
		values := make([]string, 0, len(tokens))
		for _, token := range tokens {
			values = append(values, token.DSSIDAsString())
		}
		return values
	}
	cmsSourcesKatCheck(t, expectations, "cert.all", cmsSourcesKatJoin(tokenIDs(source.Certificates())))
	cmsSourcesKatCheck(t, expectations, "cert.signedData", cmsSourcesKatJoin(tokenIDs(source.SignedDataCertificates())))
	cmsSourcesKatCheck(t, expectations, "cert.certificateValues", cmsSourcesKatJoin(tokenIDs(source.CertificateValues())))

	identifiers := make([]string, 0)
	for _, identifier := range source.AllCertificateIdentifiers() {
		identifiers = append(identifiers, cmsSourcesKatSignerIdentifier(identifier))
	}
	cmsSourcesKatCheck(t, expectations, "cert.identifiers", cmsSourcesKatJoin(identifiers))
	cmsSourcesKatCheck(t, expectations, "cert.currentIdentifier",
		cmsSourcesKatSignerIdentifier(source.CurrentCertificateIdentifier()))

	certRefs := func(refs []*CertificateRef) []string {
		values := make([]string, 0, len(refs))
		for _, ref := range refs {
			builder := &strings.Builder{}
			builder.WriteString(ref.DSSIDAsString())
			builder.WriteString(":")
			if certDigest := ref.CertDigest(); !certDigest.IsEmpty() {
				builder.WriteString(string(certDigest.Algorithm()))
				builder.WriteString(":")
				builder.WriteString(hex.EncodeToString(certDigest.Value()))
			}
			builder.WriteString(":")
			builder.WriteString(cmsSourcesKatSignerIdentifier(ref.CertificateIdentifier()))
			values = append(values, builder.String())
		}
		return values
	}
	cmsSourcesKatCheck(t, expectations, "ref.signingCertificate", cmsSourcesKatJoin(certRefs(source.SigningCertificateRefs())))
	cmsSourcesKatCheck(t, expectations, "ref.completeCertificateRefs", cmsSourcesKatJoin(certRefs(source.CompleteCertificateRefs())))
	cmsSourcesKatCheck(t, expectations, "ref.attributeCertificateRefs", cmsSourcesKatJoin(certRefs(source.AttributeCertificateRefs())))
}

func TestCMSCRLSourceKAT(t *testing.T) {
	cms, signer := cmsSourcesKatFixture(t)
	expectations := cmsSourcesKatExpectations(t)

	source, err := NewCMSCRLSource(cms.CRLs(), signer.UnsignedAttributes)
	if err != nil {
		t.Fatalf("NewCMSCRLSource: %v", err)
	}

	cmsSourcesKatCheck(t, expectations, "crl.cmsSignedData",
		cmsSourcesKatJoin(cmsSourcesKatBinaryIDs(source.CMSSignedDataRevocationBinaries())))
	cmsSourcesKatCheck(t, expectations, "crl.revocationValues",
		cmsSourcesKatJoin(cmsSourcesKatBinaryIDs(source.RevocationValuesBinaries())))

	crlRefs := func(refs []RevocationRef[revocation.CRL]) []string {
		values := make([]string, 0, len(refs))
		for _, ref := range refs {
			crlRef, ok := ref.(*CRLRef)
			if !ok {
				t.Fatalf("a CRL source produced a %T reference", ref)
			}
			issuer := ""
			if crlRef.CRLIssuer() != nil {
				issuer = crlRef.CRLIssuer().RFC2253Name()
			}
			issueTime := ""
			if !crlRef.CRLIssueTime().IsZero() {
				issueTime = fmt.Sprintf("%d", crlRef.CRLIssueTime().UnixMilli())
			}
			number := ""
			if crlRef.CRLNumber() != nil {
				number = crlRef.CRLNumber().String()
			}
			values = append(values, fmt.Sprintf("%s:%s:%s:%s:%s:%s", crlRef.DSSIDAsString(),
				string(crlRef.Digest().Algorithm()), hex.EncodeToString(crlRef.Digest().Value()),
				issuer, issueTime, number))
		}
		return values
	}
	cmsSourcesKatCheck(t, expectations, "crl.completeRefs", cmsSourcesKatJoin(crlRefs(source.CompleteRevocationRefs())))
	cmsSourcesKatCheck(t, expectations, "crl.attributeRefs", cmsSourcesKatJoin(crlRefs(source.AttributeRevocationRefs())))
}

func TestCMSOCSPSourceKAT(t *testing.T) {
	cms, signer := cmsSourcesKatFixture(t)
	expectations := cmsSourcesKatExpectations(t)

	source, err := NewCMSOCSPSource(cms.OCSPResponses(), cms.OCSPBasicResponses(),
		signer.UnsignedAttributes)
	if err != nil {
		t.Fatalf("NewCMSOCSPSource: %v", err)
	}

	fromSignedData := make([]string, 0)
	for _, binary := range source.CMSSignedDataRevocationBinaries() {
		ocspBinary, ok := binary.(*OCSPResponseBinary)
		if !ok {
			t.Fatalf("an OCSP source produced a %T binary", binary)
		}
		fromSignedData = append(fromSignedData,
			ocspBinary.AsXmlID()+":"+ocspBinary.ASN1ObjectIdentifier().String())
	}
	cmsSourcesKatCheck(t, expectations, "ocsp.cmsSignedData", cmsSourcesKatJoin(fromSignedData))
	cmsSourcesKatCheck(t, expectations, "ocsp.revocationValues",
		cmsSourcesKatJoin(cmsSourcesKatBinaryIDs(source.RevocationValuesBinaries())))

	ocspRefs := func(refs []RevocationRef[revocation.OCSP]) []string {
		values := make([]string, 0, len(refs))
		for _, ref := range refs {
			ocspRef, ok := ref.(*OCSPRef)
			if !ok {
				t.Fatalf("an OCSP source produced a %T reference", ref)
			}
			producedAt := ""
			if !ocspRef.ProducedAt().IsZero() {
				producedAt = fmt.Sprintf("%d", ocspRef.ProducedAt().UnixMilli())
			}
			responder := ""
			if ocspRef.ResponderId() != nil && ocspRef.ResponderId().X500Principal() != nil {
				responder = ocspRef.ResponderId().X500Principal().RFC2253Name()
			}
			values = append(values, fmt.Sprintf("%s:%s:%s:%s:%s", ocspRef.DSSIDAsString(),
				string(ocspRef.Digest().Algorithm()), hex.EncodeToString(ocspRef.Digest().Value()),
				producedAt, responder))
		}
		return values
	}
	cmsSourcesKatCheck(t, expectations, "ocsp.completeRefs", cmsSourcesKatJoin(ocspRefs(source.CompleteRevocationRefs())))
	cmsSourcesKatCheck(t, expectations, "ocsp.attributeRefs", cmsSourcesKatJoin(ocspRefs(source.AttributeRevocationRefs())))
}

// cmsSourcesKatBinaryIDs is the generator's binaries() helper.
func cmsSourcesKatBinaryIDs[R revocation.Revocation](binaries []EncapsulatedRevocationTokenIdentifier[R]) []string {
	values := make([]string, 0, len(binaries))
	for _, binary := range binaries {
		values = append(values, binary.AsXmlID())
	}
	return values
}

// TestCMSRevocationSourcesDeduplicateBinaries checks that a revocation datum reached through
// two origins is stored once, carrying both - the fixture's CRL sits in SignedData.crls and
// in the revocation-values attribute, and its OCSP response likewise.
func TestCMSRevocationSourcesDeduplicateBinaries(t *testing.T) {
	cms, signer := cmsSourcesKatFixture(t)

	crlSource, err := NewCMSCRLSource(cms.CRLs(), signer.UnsignedAttributes)
	if err != nil {
		t.Fatalf("NewCMSCRLSource: %v", err)
	}
	if got := len(crlSource.AllRevocationBinaries()); got != 1 {
		t.Errorf("the CRL source holds %d binaries, expected 1", got)
	}
	if got := len(crlSource.CMSSignedDataRevocationBinaries()); got != 1 {
		t.Errorf("the CRL source holds %d CMS_SIGNED_DATA binaries, expected 1", got)
	}
	if got := len(crlSource.RevocationValuesBinaries()); got != 1 {
		t.Errorf("the CRL source holds %d REVOCATION_VALUES binaries, expected 1", got)
	}

	ocspSource, err := NewCMSOCSPSource(cms.OCSPResponses(), cms.OCSPBasicResponses(),
		signer.UnsignedAttributes)
	if err != nil {
		t.Fatalf("NewCMSOCSPSource: %v", err)
	}
	// Two distinct responses come from SignedData.crls, one of which is repeated in the
	// revocation-values attribute.
	if got := len(ocspSource.AllRevocationBinaries()); got != 2 {
		t.Errorf("the OCSP source holds %d binaries, expected 2", got)
	}
	if got := len(ocspSource.CMSSignedDataRevocationBinaries()); got != 2 {
		t.Errorf("the OCSP source holds %d CMS_SIGNED_DATA binaries, expected 2", got)
	}
	if got := len(ocspSource.RevocationValuesBinaries()); got != 1 {
		t.Errorf("the OCSP source holds %d REVOCATION_VALUES binaries, expected 1", got)
	}
}
