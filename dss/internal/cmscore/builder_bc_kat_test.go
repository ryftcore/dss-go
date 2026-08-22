package cmscore

import (
	"encoding/asn1"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// The build-side differential test against BouncyCastle. testdata/build holds SignedData
// documents assembled by testdata/gen/BuildOracle.java out of BouncyCastle's low-level
// org.bouncycastle.asn1.cms classes, which derive both CMSVersions themselves; the cases
// below rebuild each of them through this package's public builders from the same inputs and
// require the DER to be byte-identical. That pins the DER SET OF ordering of the certificates,
// the CRLs, the signer infos and the attribute values, the RFC 5652 clause 5.1 and 5.3 version
// computation, and the NULL-versus-absent handling of algorithm parameters.

// buildInputs holds the material the BouncyCastle build oracle used, rebuilt here from the
// same checked-in fixtures and the same constants.
type buildInputs struct {
	cert1, cert2   []byte
	crl            []byte
	ocsp           []byte
	issuer         []byte
	serial         *big.Int
	ski            []byte
	signature      []byte
	content        []byte
	sha256, sha512 *asn1ber.AlgorithmIdentifier
	rsa, ed25519   *asn1ber.AlgorithmIdentifier
	reversedAttrs  Attributes
}

func loadBuildInputs(t *testing.T) *buildInputs {
	t.Helper()
	chain, err := ParseCMS(readFixture(t, "rsa-sha256-chain.p7s"))
	if err != nil {
		t.Fatal(err)
	}
	certs := chain.SignedData().Certificates.Choices
	oc, err := ParseCMS(readFixture(t, "ocsp-crl.p7s"))
	if err != nil {
		t.Fatal(err)
	}
	inputs := &buildInputs{
		cert1:     certs[0].Encoded(),
		cert2:     certs[1].Encoded(),
		crl:       oc.CRLs()[0],
		ocsp:      oc.OCSPResponses()[0],
		issuer:    chain.SignerInfos()[0].SID.IssuerAndSerialNumber.Issuer,
		serial:    chain.SignerInfos()[0].SID.IssuerAndSerialNumber.SerialNumber,
		ski:       []byte{0x0a, 0x0b, 0x0c, 0x0d, 0x0e},
		signature: make([]byte, 64),
		content:   []byte("adversarial content"),
		sha256:    asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}),
		sha512:    asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}),
		rsa: asn1ber.NewAlgorithmIdentifierWithParameters(
			asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}, asn1ber.DERNull),
		ed25519: asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{1, 3, 101, 112}),
	}
	for index := range inputs.signature {
		inputs.signature[index] = byte(index + 1)
	}
	digest := make([]byte, 32)
	for index := range digest {
		digest[index] = byte(0x40 + index)
	}
	var capabilities []byte
	for index := 0; index < 6; index++ {
		capabilities = append(capabilities, asn1ber.WriteSequence(
			asn1ber.EncodeOID(asn1.ObjectIdentifier{1, 2, 840, 113549, 3, index + 2}))...)
	}
	inputs.reversedAttrs = Attributes{
		NewAttribute(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 15},
			asn1ber.WriteSequence(capabilities)),
		NewAttribute(OIDMessageDigest, asn1ber.WriteTLV(asn1ber.TagOctetString, digest)),
		NewAttribute(OIDSigningTime, asn1ber.WriteTLV(asn1ber.TagUTCTime, []byte("240102030405Z"))),
		NewAttribute(OIDContentType, asn1ber.EncodeOID(OIDData)),
	}
	return inputs
}

// TestBuilderMatchesBouncyCastle rebuilds every testdata/build document and compares the DER.
func TestBuilderMatchesBouncyCastle(t *testing.T) {
	dir := filepath.Join("testdata", "build")
	in := loadBuildInputs(t)
	iasSID := NewIssuerAndSerialNumberSID(in.issuer, in.serial)
	skiSID := NewSubjectKeyIdentifierSID(in.ski)

	signer := func(sid *SignerIdentifier, digest, signatureAlgorithm *asn1ber.AlgorithmIdentifier,
		signed, unsigned Attributes) *SignerInfo {
		builder := &SignerInfoBuilder{SID: sid, DigestAlgorithm: digest,
			SignedAttributes: signed, SignatureAlgorithm: signatureAlgorithm,
			Signature: in.signature, UnsignedAttributes: unsigned}
		built, err := builder.Build()
		if err != nil {
			t.Fatal(err)
		}
		return built
	}
	attached := NewEncapsulatedContentInfo(OIDData, in.content)
	base := signer(iasSID, in.sha256, in.rsa, in.reversedAttrs, nil)
	cert1 := NewCertificateChoice(in.cert1)
	fake := asn1ber.WriteSequence(append(asn1ber.EncodeInteger(big.NewInt(1)),
		asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{9, 9, 9})...))
	tagged := func(tagNo int) CertificateChoice {
		return NewTaggedCertificateChoice(tagNo,
			asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|byte(tagNo), fake[2:]))
	}

	cases := []struct {
		name    string
		builder *SignedDataBuilder
		version int
	}{
		{"b1-basic", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{base}}, 1},
		{"b2-ski", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{signer(skiSID, in.sha256, in.rsa, in.reversedAttrs, nil)}}, 3},
		{"b3-detached", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: NewEncapsulatedContentInfo(OIDData, nil),
			Certificates:     []CertificateChoice{cert1}, SignerInfos: []*SignerInfo{base}}, 1},
		{"b4-two-certs", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached,
			Certificates:     []CertificateChoice{NewCertificateChoice(in.cert2), cert1},
			SignerInfos:      []*SignerInfo{base}}, 1},
		{"b5-crl-other", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			CRLs: []RevocationInfoChoice{
				NewOtherRevocationInfoChoice(NewOtherRevocationInfoFormat(OIDRIOCSPResponse, in.ocsp)),
				NewCRLRevocationInfoChoice(in.crl)},
			SignerInfos: []*SignerInfo{base}}, 5},
		{"b6-attr-cert-v2", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1, tagged(2)},
			SignerInfos: []*SignerInfo{base}}, 4},
		{"b7-attr-cert-v1", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1, tagged(1)},
			SignerInfos: []*SignerInfo{base}}, 3},
		{"b8-other-cert", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1, tagged(3)},
			SignerInfos: []*SignerInfo{base}}, 5},
		{"b9-tstinfo-type", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: NewEncapsulatedContentInfo(OIDCTTSTInfo, in.content),
			Certificates:     []CertificateChoice{cert1}, SignerInfos: []*SignerInfo{base}}, 3},
		{"b10-unsigned-attrs", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{signer(iasSID, in.sha256, in.rsa, in.reversedAttrs, Attributes{
				NewAttribute(OIDCounterSignature, asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{1, 2, 3})),
				NewAttribute(asn1.ObjectIdentifier{1, 2, 3, 4, 5}, asn1ber.EncodeInteger(big.NewInt(42))),
			})}}, 1},
		{"b11-multivalue", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{signer(iasSID, in.sha256, in.rsa, Attributes{
				NewAttribute(OIDContentType, asn1ber.EncodeOID(OIDData)),
				NewAttribute(asn1.ObjectIdentifier{1, 2, 3, 4, 99},
					asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{0xff, 0xff}),
					asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{0x00}),
					asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{0x00, 0x01})),
			}, nil)}}, 1},
		{"b12-two-signers", &SignedDataBuilder{
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{base,
				signer(skiSID, in.sha512, in.ed25519, in.reversedAttrs, nil)}}, 3},
		{"b13-ed25519", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha512},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{signer(iasSID, in.sha512, in.ed25519, in.reversedAttrs, nil)}}, 1},
		{"b14-no-signed-attrs", &SignedDataBuilder{DigestAlgorithms: []*asn1ber.AlgorithmIdentifier{in.sha256},
			EncapContentInfo: attached, Certificates: []CertificateChoice{cert1},
			SignerInfos: []*SignerInfo{signer(iasSID, in.sha256, in.rsa, nil, nil)}}, 1},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			expected, err := os.ReadFile(filepath.Join(dir, testCase.name+".der"))
			if err != nil {
				t.Fatal(err)
			}
			cms, err := testCase.builder.BuildCMS()
			if err != nil {
				t.Fatal(err)
			}
			if got := cms.Version(); got != testCase.version {
				t.Errorf("version = %d, BouncyCastle computed %d", got, testCase.version)
			}
			if got := cms.ContentInfo().DER(); string(got) != string(expected) {
				t.Errorf("DER differs from BouncyCastle\n go %x\n bc %x", got, expected)
			}
		})
	}
}
