package cmscore

import (
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// signerInfoOfVersion returns a SignerInfo carrying only the version, which is all
// ComputeSignedDataVersion looks at.
func signerInfoOfVersion(version int) *SignerInfo {
	return &SignerInfo{Version: version}
}

// certificateSetOf builds a CertificateSet holding one member per given alternative.
func certificateSetOf(tagNumbers ...int) *CertificateSet {
	set := &CertificateSet{}
	for _, tagNo := range tagNumbers {
		if tagNo == CertificateChoiceCertificate {
			set.Choices = append(set.Choices, NewCertificateChoice([]byte{0x30, 0x00}))
			continue
		}
		set.Choices = append(set.Choices, NewTaggedCertificateChoice(tagNo, encodeContextTagged(tagNo, nil)))
	}
	return set
}

// TestComputeSignedDataVersion walks the decision tree of RFC 5652 clause 5.1 branch by branch.
func TestComputeSignedDataVersion(t *testing.T) {
	otherFormat := asn1.ObjectIdentifier{1, 2, 3, 4}
	crlOnly := &RevocationInfoChoices{Choices: []RevocationInfoChoice{
		NewCRLRevocationInfoChoice([]byte{0x30, 0x00}),
	}}
	withOther := &RevocationInfoChoices{Choices: []RevocationInfoChoice{
		NewCRLRevocationInfoChoice([]byte{0x30, 0x00}),
		NewOtherRevocationInfoChoice(NewOtherRevocationInfoFormat(otherFormat, asn1ber.DERNull)),
	}}

	cases := []struct {
		name         string
		eContentType asn1.ObjectIdentifier
		certificates *CertificateSet
		crls         *RevocationInfoChoices
		signerInfos  []*SignerInfo
		want         int
	}{
		{
			name: "plain id-data signature", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion1,
		},
		{
			name: "no optional field at all", eContentType: OIDData,
			signerInfos: []*SignerInfo{signerInfoOfVersion(1)},
			want:        CMSVersion1,
		},
		{
			name: "CRLs but no other format", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate), crls: crlOnly,
			signerInfos: []*SignerInfo{signerInfoOfVersion(1)},
			want:        CMSVersion1,
		},
		{
			name: "eContentType other than id-data", eContentType: OIDCTTSTInfo,
			certificates: certificateSetOf(CertificateChoiceCertificate),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion3,
		},
		{
			name: "a version 3 SignerInfo", eContentType: OIDData,
			signerInfos: []*SignerInfo{signerInfoOfVersion(1), signerInfoOfVersion(3)},
			want:        CMSVersion3,
		},
		{
			name: "a version 1 attribute certificate", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate, CertificateChoiceV1AttrCert),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion3,
		},
		{
			name: "a version 2 attribute certificate", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate, CertificateChoiceV2AttrCert),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion4,
		},
		{
			name: "version 2 beats version 1 attribute certificates", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceV1AttrCert, CertificateChoiceV2AttrCert),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion4,
		},
		{
			name: "an other certificate format", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate, CertificateChoiceOther),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion5,
		},
		{
			name: "an other revocation info format", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceCertificate), crls: withOther,
			signerInfos: []*SignerInfo{signerInfoOfVersion(1)},
			want:        CMSVersion5,
		},
		{
			name: "an other format beats a version 2 attribute certificate", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceV2AttrCert), crls: withOther,
			signerInfos: []*SignerInfo{signerInfoOfVersion(3)},
			want:        CMSVersion5,
		},
		{
			// The obsolete extendedCertificate alternative is not one of the cases the
			// rule singles out, so it does not raise the version by itself.
			name: "an extended certificate", eContentType: OIDData,
			certificates: certificateSetOf(CertificateChoiceExtendedCertificate),
			signerInfos:  []*SignerInfo{signerInfoOfVersion(1)},
			want:         CMSVersion1,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ComputeSignedDataVersion(testCase.eContentType, testCase.certificates,
				testCase.crls, testCase.signerInfos)
			if got != testCase.want {
				t.Errorf("version = %d, want %d", got, testCase.want)
			}
		})
	}
}

// TestComputeSignerInfoVersion checks the rule of RFC 5652 clause 5.3.
func TestComputeSignerInfoVersion(t *testing.T) {
	issuerAndSerial := NewIssuerAndSerialNumberSID([]byte{0x30, 0x00}, big.NewInt(42))
	if got := ComputeSignerInfoVersion(issuerAndSerial); got != CMSVersion1 {
		t.Errorf("an issuerAndSerialNumber sid gives version %d, want 1", got)
	}
	subjectKeyIdentifier := NewSubjectKeyIdentifierSID([]byte{0x01, 0x02, 0x03})
	if got := ComputeSignerInfoVersion(subjectKeyIdentifier); got != CMSVersion3 {
		t.Errorf("a subjectKeyIdentifier sid gives version %d, want 3", got)
	}
}

// TestSignerIdentifierEncodings checks both alternatives of the CHOICE, the implicitly tagged
// one being primitive.
func TestSignerIdentifierEncodings(t *testing.T) {
	subjectKeyIdentifier := NewSubjectKeyIdentifierSID([]byte{0xAA, 0xBB})
	want := []byte{0x80, 0x02, 0xAA, 0xBB}
	if got := subjectKeyIdentifier.DER(); string(got) != string(want) {
		t.Errorf("subjectKeyIdentifier sid = %x, want %x", got, want)
	}
	if !subjectKeyIdentifier.IsSubjectKeyIdentifier() {
		t.Errorf("the [0] alternative does not report itself")
	}

	issuer := []byte{0x30, 0x03, 0x02, 0x01, 0x07}
	issuerAndSerial := NewIssuerAndSerialNumberSID(issuer, big.NewInt(1))
	if issuerAndSerial.IsSubjectKeyIdentifier() {
		t.Errorf("the issuerAndSerialNumber alternative reports itself as [0]")
	}
	encoded := issuerAndSerial.DER()
	if encoded[0] != 0x30 {
		t.Errorf("issuerAndSerialNumber sid = %x, want a SEQUENCE", encoded)
	}
	// Re-parsing has to give back the very same issuer bytes.
	element, err := parseOne(encoded, "sid")
	if err != nil {
		t.Fatalf("the encoding does not parse back: %v", err)
	}
	parsed, err := signerIdentifierFromElement(element)
	if err != nil {
		t.Fatalf("signerIdentifierFromElement: %v", err)
	}
	if string(parsed.IssuerAndSerialNumber.Issuer) != string(issuer) {
		t.Errorf("issuer = %x, want %x", parsed.IssuerAndSerialNumber.Issuer, issuer)
	}
	if parsed.IssuerAndSerialNumber.SerialNumber.Int64() != 1 {
		t.Errorf("serialNumber = %s, want 1", parsed.IssuerAndSerialNumber.SerialNumber)
	}
}
