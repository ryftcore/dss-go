// Behaviour tests for the three CMS sources and the ESS/ESF structures they read, covering
// the branches the cades-full.p7s known-answer fixture cannot reach (a subjectKeyIdentifier
// signer, a signer absent from SignerInfos, the DEFAULT of ESSCertIDv2.hashAlgorithm, a
// multi-valued attribute) plus the candidate resolution of
// ExtractCandidatesForSigningCertificate.
package spi

import (
	"bytes"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
)

// cmsSourcesTestSignerInfo builds a SignerInfo carrying nothing but a SignerIdentifier, which
// is all extractCertificateIdentifiers reads.
func cmsSourcesTestSignerInfo(sid *cmscore.SignerIdentifier) *cmscore.SignerInfo {
	return &cmscore.SignerInfo{SID: sid}
}

func TestCMSCertificateSourceIdentifierFromSubjectKeyIdentifier(t *testing.T) {
	cms, fixtureSigner := cmsSourcesKatFixture(t)
	issuer := fixtureSigner.SID.IssuerAndSerialNumber.Issuer
	serialNumber := fixtureSigner.SID.IssuerAndSerialNumber.SerialNumber

	ski := []byte{0xfc, 0x6c, 0x33, 0x88, 0xb0, 0xc2, 0x80, 0x64}
	byIssuerAndSerial := cmsSourcesTestSignerInfo(cmscore.NewIssuerAndSerialNumberSID(issuer, serialNumber))
	bySubjectKeyIdentifier := cmsSourcesTestSignerInfo(cmscore.NewSubjectKeyIdentifierSID(ski))

	source, err := NewCMSCertificateSource(
		[]*cmscore.SignerInfo{byIssuerAndSerial, bySubjectKeyIdentifier},
		cms.Certificates(), bySubjectKeyIdentifier)
	if err != nil {
		t.Fatalf("NewCMSCertificateSource: %v", err)
	}

	identifiers := source.AllCertificateIdentifiers()
	if len(identifiers) != 2 {
		t.Fatalf("got %d certificate identifiers, want 2", len(identifiers))
	}
	if identifiers[0].IssuerName() == nil || identifiers[0].SerialNumber() == nil {
		t.Errorf("the issuerAndSerialNumber identifier lost its issuer or serial number")
	}
	if identifiers[0].Ski() != nil {
		t.Errorf("the issuerAndSerialNumber identifier carries a SKI: %x", identifiers[0].Ski())
	}
	if identifiers[0].IsCurrent() {
		t.Errorf("the issuerAndSerialNumber identifier is marked current")
	}

	current := source.CurrentCertificateIdentifier()
	if current == nil {
		t.Fatal("no current certificate identifier")
	}
	if !bytes.Equal(current.Ski(), ski) {
		t.Errorf("current SKI = %x, want %x", current.Ski(), ski)
	}
	if current.IssuerName() != nil || current.SerialNumber() != nil {
		t.Errorf("the subjectKeyIdentifier identifier carries an issuer or a serial number")
	}
}

func TestCMSCertificateSourceCurrentSignerNotInSignerInfos(t *testing.T) {
	cms, fixtureSigner := cmsSourcesKatFixture(t)
	current := cmsSourcesTestSignerInfo(cmscore.NewIssuerAndSerialNumberSID(
		fixtureSigner.SID.IssuerAndSerialNumber.Issuer, big.NewInt(0x4242)))

	// An empty SignerInformationStore makes the lookup fail, so upstream adds the current
	// identifier itself, flagged as current.
	source, err := NewCMSCertificateSource(nil, cms.Certificates(), current)
	if err != nil {
		t.Fatalf("NewCMSCertificateSource: %v", err)
	}
	identifiers := source.AllCertificateIdentifiers()
	if len(identifiers) != 1 {
		t.Fatalf("got %d certificate identifiers, want 1", len(identifiers))
	}
	if !identifiers[0].IsCurrent() {
		t.Error("the fallback identifier is not marked current")
	}
	if identifiers[0].SerialNumber().Cmp(big.NewInt(0x4242)) != 0 {
		t.Errorf("serial number = %s, want 16962", identifiers[0].SerialNumber())
	}
}

func TestCMSCertificateSourceExtractCandidatesForSigningCertificate(t *testing.T) {
	cms, signer := cmsSourcesKatFixture(t)
	source, err := NewCMSCertificateSource(cms.SignerInfos(), cms.Certificates(), signer)
	if err != nil {
		t.Fatalf("NewCMSCertificateSource: %v", err)
	}

	// No external source is needed: the signing certificate is in SignedData.certificates.
	candidates := source.CandidatesForSigningCertificate(nil)
	validity := candidates.TheCertificateValidity()
	if validity == nil {
		t.Fatal("no certificate validity was selected")
	}
	if validity.CertificateToken() == nil {
		t.Fatal("the selected validity carries no certificate")
	}
	if !validity.IsDigestPresent() {
		t.Error("digestPresent = false, want true")
	}
	if !validity.IsDigestEqual() {
		t.Error("digestEqual = false, want true")
	}
	if !validity.IsIssuerSerialPresent() {
		t.Error("issuerSerialPresent = false, want true")
	}
	if !validity.IsSerialNumberEqual() {
		t.Error("serialNumberEqual = false, want true")
	}
	if !validity.IsDistinguishedNameEqual() {
		t.Error("distinguishedNameEqual = false, want true")
	}
	if !validity.IsSignerIdMatch() {
		t.Error("signerIdMatch = false, want true")
	}
	if !validity.IsValid() {
		t.Error("the selected validity is not valid")
	}

	// The result is cached, as getCandidatesForSigningCertificate promises.
	if source.CandidatesForSigningCertificate(nil) != candidates {
		t.Error("the candidates were recomputed instead of being cached")
	}
}

func TestCMSCertificateSourceExtractCandidatesUnresolvedCertificate(t *testing.T) {
	_, fixtureSigner := cmsSourcesKatFixture(t)
	// A signer whose certificate is not in the CMS: the validity is built from the
	// identifier, and the reference comparison runs against it rather than a token.
	current := cmsSourcesTestSignerInfo(cmscore.NewIssuerAndSerialNumberSID(
		fixtureSigner.SID.IssuerAndSerialNumber.Issuer, big.NewInt(0x4242)))
	source, err := NewCMSCertificateSource(nil, nil, current)
	if err != nil {
		t.Fatalf("NewCMSCertificateSource: %v", err)
	}

	candidates := source.CandidatesForSigningCertificate(nil)
	validity := candidates.TheCertificateValidity()
	if validity == nil {
		t.Fatal("no certificate validity was selected")
	}
	if validity.CertificateToken() != nil {
		t.Error("a certificate token was resolved, though none is available")
	}
	if validity.SignerInfo() == nil {
		t.Error("the validity carries no signer identifier")
	}
	// No signing-certificate reference was extracted, so the reference-driven flags stay unset.
	if validity.IsDigestPresent() || validity.IsIssuerSerialPresent() {
		t.Error("reference-driven flags were set without a signing certificate reference")
	}
}

func TestParseESSCertIDv2DefaultHashAlgorithm(t *testing.T) {
	certHash := make([]byte, 32)
	for index := range certHash {
		certHash[index] = byte(index + 1)
	}
	// ESSCertIDv2 ::= SEQUENCE { certHash OCTET STRING }, the hashAlgorithm DEFAULT omitted.
	encoded := asn1ber.WriteSequence(append([]byte{0x04, 0x20}, certHash...))

	essCertID, err := ParseESSCertIDv2(encoded)
	if err != nil {
		t.Fatalf("ParseESSCertIDv2: %v", err)
	}
	sha256 := asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	if !essCertID.HashAlgorithm.Algorithm.Equal(sha256) {
		t.Errorf("hashAlgorithm = %s, want id-sha256", essCertID.HashAlgorithm.Algorithm)
	}
	if !bytes.Equal(essCertID.CertHash, certHash) {
		t.Errorf("certHash = %x, want %x", essCertID.CertHash, certHash)
	}
	if essCertID.IssuerSerial != nil {
		t.Error("an issuerSerial was decoded from a one-member sequence")
	}
}

func TestParseESSCertIDv2BadSequenceSize(t *testing.T) {
	if _, err := ParseESSCertIDv2(asn1ber.WriteSequence(nil)); err == nil {
		t.Error("an empty ESSCertIDv2 was accepted")
	}
}

func TestDSSASN1UtilsAsn1EncodableOnlyAcceptsASingleValue(t *testing.T) {
	oid := asn1.ObjectIdentifier{1, 2, 3, 4}
	single := cmscore.NewAttribute(oid, []byte{0x05, 0x00})
	parsed, err := cmscore.AttributeFromElement(mustParseElement(t, single.DER()))
	if err != nil {
		t.Fatalf("AttributeFromElement: %v", err)
	}
	if DSSASN1UtilsAsn1Encodable(parsed) == nil {
		t.Error("a single-valued attribute yielded no value")
	}

	two := cmscore.NewAttribute(oid, []byte{0x05, 0x00}, []byte{0x01, 0x01, 0xff})
	parsedTwo, err := cmscore.AttributeFromElement(mustParseElement(t, two.DER()))
	if err != nil {
		t.Fatalf("AttributeFromElement: %v", err)
	}
	if DSSASN1UtilsAsn1Encodable(parsedTwo) != nil {
		t.Error("a two-valued attribute yielded a value")
	}

	if DSSASN1UtilsAsn1Encodable(nil) != nil {
		t.Error("a nil attribute yielded a value")
	}
}

func TestDSSASN1UtilsAsn1AttributesIsNeverNil(t *testing.T) {
	attributes := DSSASN1UtilsAsn1Attributes(nil, asn1.ObjectIdentifier{1, 2, 3, 4})
	if attributes == nil {
		t.Fatal("a nil attribute table yielded a nil slice, want an empty one")
	}
	if len(attributes) != 0 {
		t.Errorf("got %d attributes, want 0", len(attributes))
	}
}

func TestDSSASN1UtilsRevocationValuesRejectsGarbage(t *testing.T) {
	if DSSASN1UtilsRevocationValues(nil) != nil {
		t.Error("a nil encodable yielded a RevocationValues")
	}
	if DSSASN1UtilsRevocationValues([]byte{0x05, 0x00}) != nil {
		t.Error("a NULL yielded a RevocationValues")
	}
	// [3] is not one of the three defined fields.
	unknown := asn1ber.WriteSequence(asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|3, nil))
	if DSSASN1UtilsRevocationValues(unknown) != nil {
		t.Error("an unknown tag yielded a RevocationValues")
	}
}

func TestParseRevocationValuesOcspValsOnly(t *testing.T) {
	// RevocationValues ::= SEQUENCE { ocspVals [1] SEQUENCE OF BasicOCSPResponse }
	inner := asn1ber.WriteSequence(asn1ber.WriteSequence(nil))
	encoded := asn1ber.WriteSequence(asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, inner))

	values, err := ParseRevocationValues(encoded)
	if err != nil {
		t.Fatalf("ParseRevocationValues: %v", err)
	}
	if len(values.CrlVals) != 0 {
		t.Errorf("got %d crlVals, want 0", len(values.CrlVals))
	}
	if len(values.OcspVals) != 1 {
		t.Fatalf("got %d ocspVals, want 1", len(values.OcspVals))
	}
	if values.OtherRevVals != nil {
		t.Error("an otherRevVals was decoded from a sequence without one")
	}
}

func TestParseCrlOcspRefOcspidsOnly(t *testing.T) {
	// CrlOcspRef ::= SEQUENCE { ocspids [1] OcspListID }, OcspListID being a SEQUENCE holding
	// the SEQUENCE OF OcspResponsesID.
	ocspListID := asn1ber.WriteSequence(asn1ber.WriteSequence(asn1ber.WriteSequence(nil)))
	encoded := asn1ber.WriteSequence(asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, ocspListID))

	ref, err := ParseCrlOcspRef(encoded)
	if err != nil {
		t.Fatalf("ParseCrlOcspRef: %v", err)
	}
	if ref.Crlids != nil {
		t.Error("a CRLListID was decoded from a sequence without one")
	}
	if ref.Ocspids == nil {
		t.Fatal("the OcspListID was not decoded")
	}
	if _, err := ref.Ocspids.OcspResponses(); err == nil {
		t.Error("an empty SEQUENCE was accepted as an OcspResponsesID")
	}
}

func TestNewCMSCertificateSourcePanicsWithoutCurrentSigner(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("no panic")
		}
		if recovered != "currentSignerInformation is null, it must be provided!" {
			t.Errorf("panic message = %v", recovered)
		}
	}()
	_, _ = NewCMSCertificateSource(nil, nil, nil)
}

// mustParseElement parses a DER value or fails the test.
func mustParseElement(t *testing.T, encoded []byte) *asn1ber.Element {
	t.Helper()
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		t.Fatalf("parsing %x: %v", encoded, err)
	}
	if len(rest) != 0 {
		t.Fatalf("extra data after %x", encoded)
	}
	return element
}
