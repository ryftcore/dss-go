package cmscore

import (
	"bytes"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// TestSignedAttributeOrderMatchesTheProducer is the known-answer test for the DER SET OF
// ordering: the attributes of a fixture are fed back to the writer in reverse order, and the
// result has to be the producer's own bytes. That pins this package's ordering to OpenSSL's
// and BouncyCastle's rather than merely to itself.
func TestSignedAttributeOrderMatchesTheProducer(t *testing.T) {
	for _, name := range []string{"rsa-sha256-attached.p7s", "ocsp-crl.p7s", "timestamp-token.tst"} {
		t.Run(name, func(t *testing.T) {
			cms, err := ParseCMS(readFixture(t, name))
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			signerInfo := cms.SignerInfos()[0]
			reversed := make(Attributes, 0, len(signerInfo.SignedAttributes))
			for index := len(signerInfo.SignedAttributes) - 1; index >= 0; index-- {
				reversed = append(reversed, signerInfo.SignedAttributes[index])
			}
			if bytes.Equal(reversed[0].Encoded(), signerInfo.SignedAttributes[0].Encoded()) {
				t.Fatalf("the fixture has too few attributes for the order to matter")
			}
			if got := reversed.DERImplicitTagged(0); !bytes.Equal(got, signerInfo.SignedAttributesRaw()) {
				t.Errorf("re-ordering changed the encoding:\n%s",
					firstDifference(signerInfo.SignedAttributesRaw(), got))
			}
			if got := reversed.DERSetEncoded(); !bytes.Equal(got, signerInfo.SignedAttributesDER()) {
				t.Errorf("re-ordering changed the signed bytes")
			}
		})
	}
}

// TestDERSetOfOrdersByEncoding checks the ordering rule itself, including the case X.690
// clause 11.6 settles by padding: an encoding that is a prefix of another sorts first.
func TestDERSetOfOrdersByEncoding(t *testing.T) {
	members := [][]byte{
		{0x04, 0x02, 0x01, 0x02},
		{0x04, 0x01, 0x01},
		{0x02, 0x01, 0x7F},
		{0x04, 0x01, 0x00},
	}
	want := []byte{
		0x31, 0x0D,
		0x02, 0x01, 0x7F,
		0x04, 0x01, 0x00,
		0x04, 0x01, 0x01,
		0x04, 0x02, 0x01, 0x02,
	}
	if got := derSetOf([]byte{setIdentifier}, members); !bytes.Equal(got, want) {
		t.Errorf("derSetOf = %x, want %x", got, want)
	}
}

// TestAttributeLookup checks the AttributeTable accessors against a fixture that carries the
// four attributes OpenSSL writes.
func TestAttributeLookup(t *testing.T) {
	cms, err := ParseCMS(readFixture(t, "rsa-sha256-attached.p7s"))
	if err != nil {
		t.Fatalf("ParseCMS: %v", err)
	}
	attributes := cms.SignerInfos()[0].SignedAttributes
	if attributes.Get(OIDContentType) == nil {
		t.Errorf("the content-type attribute is not found")
	}
	if got := len(attributes.GetAll(OIDContentType)); got != 1 {
		t.Errorf("GetAll returned %d content-type attributes, want 1", got)
	}
	absent := asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if attributes.Get(absent) != nil {
		t.Errorf("an absent attribute was found")
	}
	if attributes.GetAll(absent) != nil {
		t.Errorf("GetAll found an absent attribute")
	}
}

// TestAttributeValuesAreDERSorted checks that an attribute's own values SET is ordered too.
func TestAttributeValuesAreDERSorted(t *testing.T) {
	attribute := NewAttribute(OIDContentType,
		[]byte{0x04, 0x01, 0x02},
		[]byte{0x04, 0x01, 0x01},
	)
	encoded := attribute.DER()
	values := encoded[len(encoded)-8:]
	want := []byte{0x31, 0x06, 0x04, 0x01, 0x01, 0x04, 0x01, 0x02}
	if !bytes.Equal(values, want) {
		t.Errorf("the attribute values are encoded as %x, want %x", values, want)
	}
}

// TestBuiltAttributeRoundTrips checks that a built attribute parses back to the same values.
func TestBuiltAttributeRoundTrips(t *testing.T) {
	value := encodeOctetString([]byte{0xDE, 0xAD})
	built := NewAttribute(OIDMessageDigest, value)
	element, err := parseOne(built.DER(), "Attribute")
	if err != nil {
		t.Fatalf("the built attribute does not parse back: %v", err)
	}
	parsed, err := AttributeFromElement(element)
	if err != nil {
		t.Fatalf("AttributeFromElement: %v", err)
	}
	if !parsed.Type.Equal(OIDMessageDigest) {
		t.Errorf("attrType = %s, want id-messageDigest", parsed.Type)
	}
	if len(parsed.Values) != 1 || !bytes.Equal(parsed.Values[0].Encoded(), value) {
		t.Errorf("the attribute value did not survive the round trip")
	}
	if !bytes.Equal(parsed.DER(), built.DER()) {
		t.Errorf("the re-encoding differs from the built one")
	}
}

// TestGeneralNameAlternativesRoundTrip checks the parser against the writer of
// asn1ber.GeneralName for every alternative a tsa field may use: parsing the encoding a
// producer wrote and writing it back has to give the same bytes.
func TestGeneralNameAlternativesRoundTrip(t *testing.T) {
	directoryName := []byte{0x30, 0x03, 0x02, 0x01, 0x07}
	// OtherName ::= SEQUENCE { type-id OBJECT IDENTIFIER, value [0] EXPLICIT ANY }. The
	// implicit [0] tag replaces its SEQUENCE tag on the wire, so only the body travels.
	otherNameBody := []byte{0x06, 0x03, 0x2A, 0x03, 0x04, 0xA0, 0x02, 0x05, 0x00}
	cases := []struct {
		name string
		// encoded is the GeneralName as it appears inside the CHOICE.
		encoded []byte
		tagNo   int
		// value is the expected Name, i.e. the alternative's value on its own.
		value []byte
	}{
		{
			name:    "otherName",
			encoded: append([]byte{0xA0, 0x09}, otherNameBody...),
			tagNo:   0, value: append([]byte{0x30, 0x09}, otherNameBody...),
		},
		{
			name:    "rfc822Name",
			encoded: []byte{0x81, 0x03, 'a', '@', 'b'},
			tagNo:   1, value: []byte{0x16, 0x03, 'a', '@', 'b'},
		},
		{
			name:    "dNSName",
			encoded: []byte{0x82, 0x03, 'd', 's', 's'},
			tagNo:   2, value: []byte{0x16, 0x03, 'd', 's', 's'},
		},
		{
			name:    "directoryName",
			encoded: append([]byte{0xA4, 0x05}, directoryName...),
			tagNo:   4, value: directoryName,
		},
		{
			name:    "uniformResourceIdentifier",
			encoded: []byte{0x86, 0x02, 'h', 'i'},
			tagNo:   6, value: []byte{0x16, 0x02, 'h', 'i'},
		},
		{
			name:    "iPAddress",
			encoded: []byte{0x87, 0x04, 10, 0, 0, 1},
			tagNo:   7, value: []byte{0x04, 0x04, 10, 0, 0, 1},
		},
		{
			name:    "registeredID",
			encoded: []byte{0x88, 0x03, 0x2A, 0x03, 0x04},
			tagNo:   8, value: []byte{0x06, 0x03, 0x2A, 0x03, 0x04},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			element, err := parseOne(testCase.encoded, "GeneralName")
			if err != nil {
				t.Fatalf("the input does not parse: %v", err)
			}
			generalName, err := ParseGeneralName(element)
			if err != nil {
				t.Fatalf("ParseGeneralName: %v", err)
			}
			if generalName.TagNo != testCase.tagNo {
				t.Errorf("tagNo = %d, want %d", generalName.TagNo, testCase.tagNo)
			}
			if !bytes.Equal(generalName.Name, testCase.value) {
				t.Errorf("Name = %x, want %x", generalName.Name, testCase.value)
			}
			if got := generalName.DER(); !bytes.Equal(got, testCase.encoded) {
				t.Errorf("the round trip gives %x, want %x", got, testCase.encoded)
			}
		})
	}
}

// TestGeneralNameRejectsUnknownAlternatives checks the guards of the CHOICE parser.
func TestGeneralNameRejectsUnknownAlternatives(t *testing.T) {
	universal, err := parseOne([]byte{0x04, 0x01, 0x00}, "value")
	if err != nil {
		t.Fatalf("the input does not parse: %v", err)
	}
	if _, err := ParseGeneralName(universal); err == nil {
		t.Errorf("a universal element was accepted as a GeneralName")
	}
	unknown, err := parseOne([]byte{0x89, 0x01, 0x00}, "value")
	if err != nil {
		t.Fatalf("the input does not parse: %v", err)
	}
	if _, err := ParseGeneralName(unknown); err == nil {
		t.Errorf("the unknown alternative [9] was accepted")
	}
}

// TestEncapsulatedContentInfoBuild checks both forms of the built encapContentInfo.
func TestEncapsulatedContentInfoBuild(t *testing.T) {
	attached := NewEncapsulatedContentInfo(OIDData, []byte{0x01, 0x02})
	if attached.IsDetached() {
		t.Errorf("an attached encapContentInfo reports itself detached")
	}
	want := []byte{0x30, 0x11,
		0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x07, 0x01,
		0xA0, 0x04, 0x04, 0x02, 0x01, 0x02}
	if got := attached.DER(); !bytes.Equal(got, want) {
		t.Errorf("attached DER = %x, want %x", got, want)
	}

	detached := NewEncapsulatedContentInfo(OIDData, nil)
	if !detached.IsDetached() {
		t.Errorf("a detached encapContentInfo reports itself attached")
	}
	if got := detached.DER(); len(got) != 13 || got[len(got)-1] != 0x01 {
		t.Errorf("detached DER = %x, want the eContentType alone", got)
	}
	// An empty - but present - content is not a detached signature.
	empty := NewEncapsulatedContentInfo(OIDData, []byte{})
	if empty.IsDetached() {
		t.Errorf("an empty eContent was taken for a detached signature")
	}
}

// TestOtherRevocationInfoImplicitTagging checks that the [1] IMPLICIT tag of a built
// OtherRevocationInfoFormat replaces the SEQUENCE tag, and that parsing it back gives the same
// format and value.
func TestOtherRevocationInfoImplicitTagging(t *testing.T) {
	value := encodeOctetString([]byte{0xAB})
	choice := NewOtherRevocationInfoChoice(NewOtherRevocationInfoFormat(OIDRIOCSPResponse, value))
	encoded := choice.DER()
	if encoded[0] != asn1ber.ClassContextSpecific|asn1ber.Constructed|1 {
		t.Fatalf("the member does not carry the [1] IMPLICIT tag: %x", encoded[:1])
	}

	set := derSetOf([]byte{contextIdentifier(1)}, [][]byte{encoded})
	element, err := parseOne(set, "crls")
	if err != nil {
		t.Fatalf("the built crls field does not parse back: %v", err)
	}
	parsed, err := revocationInfoChoicesFromElement(element)
	if err != nil {
		t.Fatalf("revocationInfoChoicesFromElement: %v", err)
	}
	if !parsed.HasOtherFormat() {
		t.Errorf("the other format is not reported")
	}
	responses := parsed.OCSPResponses()
	if len(responses) != 1 || !bytes.Equal(responses[0], value) {
		t.Errorf("OCSPResponses = %x, want %x", responses, value)
	}
	if len(parsed.CRLs()) != 0 {
		t.Errorf("an OtherRevocationInfoFormat was reported as a CRL")
	}
	if got := parsed.DER(); !bytes.Equal(got, set) {
		t.Errorf("the crls round trip gives %x, want %x", got, set)
	}
}

// TestBuilderRejectsIncompleteInput checks the builder guards.
func TestBuilderRejectsIncompleteInput(t *testing.T) {
	digest := asn1ber.NewAlgorithmIdentifierWithParameters(oidSHA256, asn1ber.DERNull)
	signature := asn1ber.NewAlgorithmIdentifierWithParameters(oidRSA, asn1ber.DERNull)
	sid := NewSubjectKeyIdentifierSID([]byte{0x01})
	cases := []struct {
		name    string
		builder *SignerInfoBuilder
	}{
		{"no sid", &SignerInfoBuilder{DigestAlgorithm: digest, SignatureAlgorithm: signature, Signature: []byte{1}}},
		{"no digest algorithm", &SignerInfoBuilder{SID: sid, SignatureAlgorithm: signature, Signature: []byte{1}}},
		{"no signature algorithm", &SignerInfoBuilder{SID: sid, DigestAlgorithm: digest, Signature: []byte{1}}},
		{"no signature", &SignerInfoBuilder{SID: sid, DigestAlgorithm: digest, SignatureAlgorithm: signature}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := testCase.builder.Build(); err == nil {
				t.Errorf("the incomplete builder was accepted")
			}
		})
	}
	if _, err := (&SignedDataBuilder{}).Build(); err == nil {
		t.Errorf("a SignedData builder without an encapContentInfo was accepted")
	}
}

// TestSignedDataBuilderDefaultsDigestAlgorithms checks that an unset digestAlgorithms field is
// filled from the signers, without duplicates, as CMSSignedDataGenerator does.
func TestSignedDataBuilderDefaultsDigestAlgorithms(t *testing.T) {
	sha256Algorithm := asn1ber.NewAlgorithmIdentifierWithParameters(oidSHA256, asn1ber.DERNull)
	sha384Algorithm := asn1ber.NewAlgorithmIdentifierWithParameters(oidSHA384, asn1ber.DERNull)
	signature := asn1ber.NewAlgorithmIdentifierWithParameters(oidRSA, asn1ber.DERNull)
	signer := func(digest *asn1ber.AlgorithmIdentifier) *SignerInfo {
		built, err := (&SignerInfoBuilder{
			SID:                NewIssuerAndSerialNumberSID([]byte{0x30, 0x00}, bigOne()),
			DigestAlgorithm:    digest,
			SignatureAlgorithm: signature,
			Signature:          []byte{0x01},
		}).Build()
		if err != nil {
			t.Fatalf("SignerInfoBuilder.Build: %v", err)
		}
		return built
	}
	builder := &SignedDataBuilder{
		EncapContentInfo: NewEncapsulatedContentInfo(OIDData, []byte{0x01}),
		SignerInfos: []*SignerInfo{
			signer(sha256Algorithm), signer(sha384Algorithm), signer(sha256Algorithm),
		},
	}
	signedData, err := builder.Build()
	if err != nil {
		t.Fatalf("SignedDataBuilder.Build: %v", err)
	}
	if len(signedData.DigestAlgorithms) != 2 {
		t.Fatalf("%d digest algorithms, want 2", len(signedData.DigestAlgorithms))
	}
	if !signedData.DigestAlgorithms[0].Equals(sha256Algorithm) ||
		!signedData.DigestAlgorithms[1].Equals(sha384Algorithm) {
		t.Errorf("digestAlgorithms = %v, want SHA-256 then SHA-384", signedData.DigestAlgorithms)
	}
	if signedData.Version != CMSVersion1 {
		t.Errorf("version = %d, want 1", signedData.Version)
	}
}

// bigOne returns the serial number used by the synthetic SignerIdentifiers above.
func bigOne() *big.Int { return big.NewInt(1) }
