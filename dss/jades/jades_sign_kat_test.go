// Known-answer test for the JAdES signing chunk against a Java oracle.
//
// testdata/jades-sign-oracle.json is produced by testdata/gen/JAdESSignOracle.java, which drives
// upstream DSS 6.5.RC1's own JAdESLevelBaselineB, JAdESCompactBuilder, JAdESSerializationBuilder
// and HttpHeadersPayloadBuilder (see that file's header for the exact command). Nothing here is
// hand-derived.
//
// For every parameter combination the golden holds four strings this port has to reproduce byte
// for byte:
//
//   - protectedHeader: the serialized JWS protected header, i.e. the insertion order of the
//     JAdESLevelBaselineB properties and the exact JSON jose4j writes for them. It is checked
//     against the first segment of the signing input, so a divergence in a single member order,
//     a stripped "application/" prefix or an int-vs-string rendering fails here.
//   - payloadBase64: the JWS payload each SignaturePackaging / 'sigD' mechanism / 'b64'
//     combination yields, including the empty payload of ObjectIdByURIHash and the concatenated
//     octets of ObjectIdByURI.
//   - dataToBeSigned: BASE64URL(header) '.' payload-per-'b64' - the bytes that actually get
//     signed.
//   - signature: the finished compact / flattened-JSON / complete-JSON serialization for a FIXED
//     signature value.
//
// The 'sigD' HttpHeaders payload builder is checked separately, in both its signature and its
// timestamp (message-body substituting) shapes.
package jades

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// jadesSignKATOracle is testdata/jades-sign-oracle.json.
type jadesSignKATOracle struct {
	SigningCertificate   string             `json:"signingCertificate"`
	SignatureValue       string             `json:"signatureValue"`
	Payload              string             `json:"payload"`
	SigningDateMillis    int64              `json:"signingDateMillis"`
	ExpirationDateMillis int64              `json:"expirationDateMillis"`
	Cases                []jadesSignKATCase `json:"cases"`
	HTTPHeaders          map[string]string  `json:"httpHeaders"`
}

// jadesSignKATCase is one row of the oracle.
type jadesSignKATCase struct {
	Name            string `json:"name"`
	ProtectedHeader string `json:"protectedHeader"`
	PayloadBase64   string `json:"payloadBase64"`
	DataToBeSigned  string `json:"dataToBeSigned"`
	MimeType        string `json:"mimeType"`
	Signature       string `json:"signature"`
}

// jadesSignKATFixture carries everything the per-case parameter builders need.
type jadesSignKATFixture struct {
	oracle             *jadesSignKATOracle
	signingCertificate *model.CertificateToken
	signatureValue     *model.SignatureValue
	payload            []byte
	signingDate        time.Time
	expirationDate     time.Time
	contentTimestamp   *validation.TimestampToken
}

func jadesSignKATLoad(t *testing.T) *jadesSignKATFixture {
	t.Helper()

	raw, err := os.ReadFile("testdata/jades-sign-oracle.json")
	if err != nil {
		t.Fatalf("cannot read the oracle: %v", err)
	}
	oracle := &jadesSignKATOracle{}
	if err := json.Unmarshal(raw, oracle); err != nil {
		t.Fatalf("cannot parse the oracle: %v", err)
	}

	der, err := base64.StdEncoding.DecodeString(oracle.SigningCertificate)
	if err != nil {
		t.Fatalf("cannot decode the signing certificate: %v", err)
	}
	signingCertificate, err := spi.DSSUtilsLoadCertificateFromBinary(der)
	if err != nil {
		t.Fatalf("cannot load the signing certificate: %v", err)
	}

	signatureValueBytes, err := base64.StdEncoding.DecodeString(oracle.SignatureValue)
	if err != nil {
		t.Fatalf("cannot decode the signature value: %v", err)
	}

	timestampBinaries, err := os.ReadFile("testdata/content-timestamp.tst")
	if err != nil {
		t.Fatalf("cannot read the content timestamp: %v", err)
	}
	contentTimestamp, err := validation.NewTimestampToken(timestampBinaries,
		enumerations.TimestampType_CONTENT_TIMESTAMP)
	if err != nil {
		t.Fatalf("cannot load the content timestamp: %v", err)
	}

	return &jadesSignKATFixture{
		oracle:             oracle,
		signingCertificate: signingCertificate,
		signatureValue: model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256,
			signatureValueBytes),
		payload:          []byte(oracle.Payload),
		signingDate:      time.UnixMilli(oracle.SigningDateMillis).UTC(),
		expirationDate:   time.UnixMilli(oracle.ExpirationDateMillis).UTC(),
		contentTimestamp: contentTimestamp,
	}
}

// baseParameters mirrors JAdESSignOracle#baseParameters.
func (f *jadesSignKATFixture) baseParameters() *JAdESSignatureParameters {
	parameters := NewJAdESSignatureParameters()
	parameters.SetSigningCertificate(f.signingCertificate)
	parameters.SetCertificateChainFromTokens(f.signingCertificate)
	parameters.SetSignatureLevel(enumerations.SignatureLevel_JAdES_BASELINE_B)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	parameters.SetSigningCertificateDigestMethod(enumerations.DigestAlgorithm_SHA256)
	signingDate := f.signingDate
	parameters.BLevel().SetSigningDate(&signingDate)
	return parameters
}

func (f *jadesSignKATFixture) enveloping(
	serializationType enumerations.JWSSerializationType) *JAdESSignatureParameters {
	parameters := f.baseParameters()
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_ENVELOPING)
	parameters.SetJwsSerializationType(serializationType)
	return parameters
}

func (f *jadesSignKATFixture) detached(serializationType enumerations.JWSSerializationType,
	mechanism enumerations.SigDMechanism) *JAdESSignatureParameters {
	parameters := f.baseParameters()
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_DETACHED)
	parameters.SetJwsSerializationType(serializationType)
	parameters.SetSigDMechanism(mechanism)
	return parameters
}

// payloadDocument mirrors JAdESSignOracle#payloadDocument.
func (f *jadesSignKATFixture) payloadDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType(f.payload, "payload.txt", enumerations.MimeTypeEnum_TEXT)
}

// detachedDocuments mirrors JAdESSignOracle#detachedDocuments.
func jadesSignKATDetachedDocuments() []model.DSSDocument {
	return []model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("first detached"), "first.txt",
			enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("second detached"), "second.bin",
			enumerations.MimeTypeEnum_BINARY),
	}
}

// jadesSignKATHTTPHeaderDocuments mirrors JAdESSignOracle#httpHeaderDocuments.
func jadesSignKATHTTPHeaderDocuments(t *testing.T) []model.DSSDocument {
	t.Helper()
	digestHeader := NewHTTPHeaderDigest(
		model.NewInMemoryDocument([]byte(`{"hello":"world"}`)), enumerations.DigestAlgorithm_SHA256)
	return []model.DSSDocument{
		NewHTTPHeader("content-type", "application/json"),
		NewHTTPHeader("x-example", " leading and trailing "),
		NewHTTPHeader("x-example", "second value"),
		digestHeader,
	}
}

// full mirrors JAdESSignOracle#full.
func (f *jadesSignKATFixture) full(t *testing.T) *JAdESSignatureParameters {
	t.Helper()
	parameters := f.enveloping(enumerations.JWSSerializationType_JSON_SERIALIZATION)

	signerLocation := model.NewSignerLocation()
	signerLocation.SetCountry("LU")
	signerLocation.SetLocality("Luxembourg")
	signerLocation.SetStateOrProvince("Luxembourg")
	signerLocation.SetPostOfficeBoxNumber("PO-42")
	signerLocation.SetPostalCode("L-1234")
	signerLocation.SetStreetAddress("1 rue de la Gare")
	parameters.BLevel().SetSignerLocation(signerLocation)

	parameters.BLevel().SetClaimedSignerRoles([]string{"head-of-department", "auditor"})
	parameters.BLevel().SetSignedAssertions([]string{`{"assertion":"value"}`})

	commitment := model.NewCommonCommitmentType()
	commitment.SetUri("http://uri.etsi.org/01903/v1.2.2#ProofOfOrigin")
	commitment.SetDescription("Proof of origin")
	commitment.SetDocumentationReferences("https://example.org/doc1", "https://example.org/doc2")
	jsonQualifier := model.NewCommitmentQualifier()
	jsonQualifier.SetContent(model.NewInMemoryDocument([]byte(`{"qualifier":"json","nested":{"a":1}}`)))
	textQualifier := model.NewCommitmentQualifier()
	textQualifier.SetContent(model.NewInMemoryDocument([]byte("plain qualifier")))
	commitment.SetCommitmentTypeQualifiers(jsonQualifier, textQualifier)
	parameters.BLevel().SetCommitmentTypeIndications([]enumerations.CommitmentType{
		enumerations.CommitmentTypeEnum_ProofOfReceipt, commitment,
	})

	policyDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA256, []byte("policy"))
	if err != nil {
		t.Fatalf("cannot compute the policy digest: %v", err)
	}
	policy := model.NewPolicy()
	policy.SetId("urn:oid:1.2.3.4.5")
	policy.SetDescription("Oracle test policy")
	policy.SetDocumentationReferences("https://example.org/policy.pdf")
	policy.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	policy.SetDigestValue(policyDigest)
	policy.SetSpuri("https://example.org/policy")
	userNotice := model.NewUserNotice()
	userNotice.SetOrganization("DSS Go Port")
	userNotice.SetNoticeNumbers(1, 2, 3)
	userNotice.SetExplicitText("Explicit notice text")
	policy.SetUserNotice(userNotice)
	spDocSpecification := model.NewSpDocSpecification()
	spDocSpecification.SetId("urn:oid:1.2.3.4.6")
	spDocSpecification.SetDescription("Policy document specification")
	spDocSpecification.SetDocumentationReferences("https://example.org/spec")
	spDocSpecification.SetQualifier(enumerations.ObjectIdentifierQualifier_OID_AS_URN)
	policy.SetSpDocSpecification(spDocSpecification)
	parameters.BLevel().SetSignaturePolicy(policy)

	parameters.SetContentTimestamps([]*validation.TimestampToken{f.contentTimestamp})
	return parameters
}

// jadesSignKATConfigure rebuilds, by oracle case name, the parameters and documents
// JAdESSignOracle#cases() used. A name the Go side does not know fails the test rather than
// silently skipping a golden row.
func (f *jadesSignKATFixture) configure(t *testing.T,
	name string) (*JAdESSignatureParameters, []model.DSSDocument) {
	t.Helper()
	switch name {
	case "compact-enveloping-default":
		return f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION),
			[]model.DSSDocument{f.payloadDocument()}

	case "compact-enveloping-b64false":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetBase64UrlEncodedPayload(false)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "json-enveloping-default":
		return f.enveloping(enumerations.JWSSerializationType_JSON_SERIALIZATION),
			[]model.DSSDocument{f.payloadDocument()}

	case "flattened-enveloping-default":
		return f.enveloping(enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION),
			[]model.DSSDocument{f.payloadDocument()}

	case "flattened-enveloping-b64false":
		parameters := f.enveloping(enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION)
		parameters.SetBase64UrlEncodedPayload(false)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "compact-detached-nosigd":
		return f.detached(enumerations.JWSSerializationType_COMPACT_SERIALIZATION,
			enumerations.SigDMechanism_NO_SIG_D), []model.DSSDocument{f.payloadDocument()}

	case "json-detached-objectidbyuri":
		return f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_OBJECT_ID_BY_URI), jadesSignKATDetachedDocuments()

	case "json-detached-objectidbyuri-b64false":
		parameters := f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_OBJECT_ID_BY_URI)
		parameters.SetBase64UrlEncodedPayload(false)
		return parameters, jadesSignKATDetachedDocuments()

	case "json-detached-objectidbyurihash":
		return f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH), jadesSignKATDetachedDocuments()

	case "json-detached-objectidbyurihash-b64false":
		parameters := f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH)
		parameters.SetBase64UrlEncodedPayload(false)
		return parameters, jadesSignKATDetachedDocuments()

	case "json-detached-objectidbyurihash-sha512":
		parameters := f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH)
		parameters.SetReferenceDigestAlgorithm(enumerations.DigestAlgorithm_SHA512)
		return parameters, jadesSignKATDetachedDocuments()

	case "json-detached-httpheaders":
		parameters := f.detached(enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.SigDMechanism_HTTP_HEADERS)
		parameters.SetBase64UrlEncodedPayload(false)
		return parameters, jadesSignKATHTTPHeaderDocuments(t)

	case "compact-enveloping-sigt":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetJadesSigningTimeType(JAdESSigningTimeType_SIG_T)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "compact-enveloping-no-signing-time":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetJadesSigningTimeType(JAdESSigningTimeType_NONE)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "compact-enveloping-x5to":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetSigningCertificateDigestMethod(enumerations.DigestAlgorithm_SHA512)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "compact-enveloping-minimal":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetIncludeKeyIdentifier(false)
		parameters.SetIncludeCertificateChain(false)
		parameters.SetIncludeSignatureType(false)
		parameters.SetContentType("application/vnd.oracle.test+json")
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "compact-enveloping-kid-x5u-typ-exp":
		parameters := f.enveloping(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		parameters.SetKeyIdentifier("my-key-identifier")
		parameters.SetX509Url("https://example.org/certs/signer.pem")
		parameters.SetSignatureType("application/jose+json")
		expirationDate := f.expirationDate
		parameters.SetExpirationTime(&expirationDate)
		return parameters, []model.DSSDocument{f.payloadDocument()}

	case "json-enveloping-full":
		return f.full(t), []model.DSSDocument{f.payloadDocument()}

	default:
		t.Fatalf("unknown oracle case %q - the Go side and testdata/gen/JAdESSignOracle.java "+
			"have drifted apart", name)
		return nil, nil
	}
}

// jadesSignKATBuilder mirrors JAdESSignOracle#builder.
func jadesSignKATBuilder(t *testing.T, parameters *JAdESSignatureParameters,
	documents []model.DSSDocument) JAdESBuilder {
	t.Helper()
	certificateVerifier := validation.NewCommonCertificateVerifier()
	var builder JAdESBuilder
	var err error
	if parameters.JwsSerializationType() == enumerations.JWSSerializationType_COMPACT_SERIALIZATION {
		builder, err = NewJAdESCompactBuilder(certificateVerifier, parameters, documents)
	} else {
		builder, err = NewJAdESSerializationBuilder(certificateVerifier, parameters, documents)
	}
	if err != nil {
		t.Fatalf("cannot build the JAdESBuilder: %v", err)
	}
	return builder
}

func TestJAdESSignKATProtectedHeaderAndPayload(t *testing.T) {
	fixture := jadesSignKATLoad(t)
	if len(fixture.oracle.Cases) == 0 {
		t.Fatal("the oracle holds no cases")
	}

	for _, oracleCase := range fixture.oracle.Cases {
		t.Run(oracleCase.Name, func(t *testing.T) {
			parameters, documents := fixture.configure(t, oracleCase.Name)
			certificateVerifier := validation.NewCommonCertificateVerifier()

			levelBaselineB, err := NewJAdESLevelBaselineB(certificateVerifier, parameters, documents)
			if err != nil {
				t.Fatalf("cannot build JAdESLevelBaselineB: %v", err)
			}
			if _, err := levelBaselineB.SignedProperties(); err != nil {
				t.Fatalf("cannot build the signed properties: %v", err)
			}
			payloadBytes, err := levelBaselineB.PayloadBytes()
			if err != nil {
				t.Fatalf("cannot build the payload: %v", err)
			}
			if got := base64.StdEncoding.EncodeToString(payloadBytes); got != oracleCase.PayloadBase64 {
				t.Errorf("payload mismatch\n got: %s\nwant: %s", got, oracleCase.PayloadBase64)
			}

			// The protected header is checked through the signing input rather than through a
			// JSON re-serialization, so the comparison covers exactly the bytes that get signed.
			builder := jadesSignKATBuilder(t, parameters, documents)
			toBeSigned, err := builder.BuildDataToBeSigned()
			if err != nil {
				t.Fatalf("cannot build the data to be signed: %v", err)
			}
			dataToBeSigned := string(toBeSigned.Bytes())
			if dataToBeSigned != oracleCase.DataToBeSigned {
				t.Errorf("dataToBeSigned mismatch\n got: %s\nwant: %s",
					dataToBeSigned, oracleCase.DataToBeSigned)
			}

			wantEncodedHeader := base64.RawURLEncoding.EncodeToString([]byte(oracleCase.ProtectedHeader))
			gotEncodedHeader := strings.SplitN(dataToBeSigned, ".", 2)[0]
			if gotEncodedHeader != wantEncodedHeader {
				gotHeader, decodeErr := base64.RawURLEncoding.DecodeString(gotEncodedHeader)
				if decodeErr != nil {
					t.Fatalf("the produced protected header is not base64url: %v", decodeErr)
				}
				t.Errorf("protected header mismatch\n got: %s\nwant: %s",
					string(gotHeader), oracleCase.ProtectedHeader)
			}
		})
	}
}

func TestJAdESSignKATSerializations(t *testing.T) {
	fixture := jadesSignKATLoad(t)

	for _, oracleCase := range fixture.oracle.Cases {
		t.Run(oracleCase.Name, func(t *testing.T) {
			parameters, documents := fixture.configure(t, oracleCase.Name)
			builder := jadesSignKATBuilder(t, parameters, documents)

			if got := builder.MimeType().MimeTypeString(); got != oracleCase.MimeType {
				t.Errorf("mime type mismatch: got %s, want %s", got, oracleCase.MimeType)
			}

			signed, err := builder.Build(fixture.signatureValue)
			if err != nil {
				t.Fatalf("cannot build the signature: %v", err)
			}
			binaries, err := spi.DSSUtilsToByteArrayOfDocument(signed)
			if err != nil {
				t.Fatalf("cannot read the signature: %v", err)
			}
			if got := string(binaries); got != oracleCase.Signature {
				t.Errorf("signature mismatch\n got: %s\nwant: %s", got, oracleCase.Signature)
			}
		})
	}
}

func TestJAdESSignKATHttpHeadersPayload(t *testing.T) {
	fixture := jadesSignKATLoad(t)

	for _, testCase := range []struct {
		key         string
		isTimestamp bool
	}{
		{key: "signature", isTimestamp: false},
		{key: "timestamp", isTimestamp: true},
	} {
		t.Run(testCase.key, func(t *testing.T) {
			want, found := fixture.oracle.HTTPHeaders[testCase.key]
			if !found {
				t.Fatalf("the oracle holds no %q httpHeaders payload", testCase.key)
			}
			payload, err := NewHttpHeadersPayloadBuilder(jadesSignKATHTTPHeaderDocuments(t),
				testCase.isTimestamp).Build()
			if err != nil {
				t.Fatalf("cannot build the HttpHeaders payload: %v", err)
			}
			if got := base64.StdEncoding.EncodeToString(payload); got != want {
				gotDecoded, _ := base64.StdEncoding.DecodeString(got)
				wantDecoded, _ := base64.StdEncoding.DecodeString(want)
				t.Errorf("HttpHeaders payload mismatch\n got: %q\nwant: %q",
					string(gotDecoded), string(wantDecoded))
			}
		})
	}
}
