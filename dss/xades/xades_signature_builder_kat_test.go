package xades

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// The KAT for the XAdES signature-building core. Every expectation in
// testdata/sign-a-builder.txt is upstream DSS 6.5.RC1's own output, dumped by
// testdata/gen/SignABuilderOracle.java against the same fixed inputs this file feeds the Go
// port: testdata/signer.crt, testdata/content-timestamp.tst, the signing date
// 2021-01-01T00:00:00Z and a constant 256-byte dummy ds:SignatureValue. Nothing here is
// hand-derived.
//
// Four keys are compared per signing case:
//
//	signedinfo-c14n  the canonicalized ds:SignedInfo Build() returns - the data-to-be-signed
//	signedinfo-raw   the ds:SignedInfo element before canonicalization
//	signedprops      the xades:SignedProperties element before canonicalization
//	document         the complete signature XML SignDocument produces
//
// Together they pin attribute order, namespace placement, prefix conventions, Id generation and
// the absence of pretty-printing, which is what "byte-compatible with upstream" means for this
// chunk.

// xadesSignABuilderSigningDate is the frozen signing date every oracle case uses.
var xadesSignABuilderSigningDate = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

// xadesSignABuilderOracle is the parsed testdata/sign-a-builder.txt: case name -> key -> bytes.
type xadesSignABuilderOracle map[string]map[string][]byte

// loadXAdESSignABuilderOracle parses the oracle dump.
func loadXAdESSignABuilderOracle(t *testing.T) xadesSignABuilderOracle {
	t.Helper()
	raw, err := os.ReadFile(corpustest.Path(t, "sign-a-builder.txt"))
	if err != nil {
		t.Fatalf("reading the oracle: %v", err)
	}
	oracle := xadesSignABuilderOracle{}
	current := ""
	for _, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "case ") {
			current = strings.TrimPrefix(line, "case ")
			if _, ok := oracle[current]; !ok {
				oracle[current] = map[string][]byte{}
			}
			continue
		}
		key, value, ok := strings.Cut(line, " ")
		if !ok {
			t.Fatalf("malformed oracle line: %q", line)
		}
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			t.Fatalf("case %s key %s: %v", current, key, err)
		}
		oracle[current][key] = decoded
	}
	if len(oracle) == 0 {
		t.Fatal("the oracle is empty")
	}
	return oracle
}

// expect returns the oracle bytes for a case/key pair, failing when they are absent.
func (o xadesSignABuilderOracle) expect(t *testing.T, caseName, key string) []byte {
	t.Helper()
	record, ok := o[caseName]
	if !ok {
		t.Fatalf("the oracle has no case %q", caseName)
	}
	value, ok := record[key]
	if !ok {
		t.Fatalf("the oracle case %q has no key %q", caseName, key)
	}
	return value
}

// xadesSignABuilderSigner loads testdata/signer.crt.
func xadesSignABuilderSigner(t *testing.T) *model.CertificateToken {
	t.Helper()
	pemBytes, err := os.ReadFile("testdata/signer.crt")
	if err != nil {
		t.Fatalf("reading the signer certificate: %v", err)
	}
	certificate, err := spi.DSSUtilsLoadCertificateFromBinary(pemBytes)
	if err != nil {
		t.Fatalf("parsing the signer certificate: %v", err)
	}
	return certificate
}

// xadesSignABuilderBaseParams mirrors SignABuilderOracle.baseParams().
func xadesSignABuilderBaseParams(t *testing.T) *SignatureParameters {
	t.Helper()
	params := NewSignatureParameters()
	signer := xadesSignABuilderSigner(t)
	params.SetSigningCertificate(signer)
	params.SetCertificateChain([]*model.CertificateToken{signer})
	signingDate := xadesSignABuilderSigningDate
	params.BLevel().SetSigningDate(&signingDate)
	params.SetSignatureLevel(enumerations.SignatureLevelXAdESBaselineB)
	params.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	return params
}

func xadesSignABuilderEnvelopingParams(t *testing.T) *SignatureParameters {
	params := xadesSignABuilderBaseParams(t)
	params.SetSignaturePackaging(enumerations.SignaturePackagingEnveloping)
	return params
}

func xadesSignABuilderEnvelopedParams(t *testing.T) *SignatureParameters {
	params := xadesSignABuilderBaseParams(t)
	params.SetSignaturePackaging(enumerations.SignaturePackagingEnveloped)
	return params
}

func xadesSignABuilderDetachedParams(t *testing.T) *SignatureParameters {
	params := xadesSignABuilderBaseParams(t)
	params.SetSignaturePackaging(enumerations.SignaturePackagingDetached)
	return params
}

func xadesSignABuilderInternallyDetachedParams(t *testing.T) *SignatureParameters {
	params := xadesSignABuilderBaseParams(t)
	params.SetSignaturePackaging(enumerations.SignaturePackagingInternallyDetached)
	return params
}

const (
	xadesSignABuilderTextContent = "Hello World!"
	xadesSignABuilderXMLContent  = `<?xml version="1.0" encoding="UTF-8"?><root xmlns="http://sample.com" Id="root-id">` +
		`<child Id="child-id">text</child></root>`
)

func xadesSignABuilderTextDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesSignABuilderTextContent), "hello.txt",
		enumerations.MimeTypeEnumText)
}

func xadesSignABuilderXMLDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesSignABuilderXMLContent), "sample.xml",
		enumerations.MimeTypeEnumXML)
}

// xadesSignABuilderSignatureValue is the constant dummy ds:SignatureValue the oracle uses:
// 256 bytes 0x00..0xFF.
func xadesSignABuilderSignatureValue() []byte {
	value := make([]byte, 256)
	for i := range value {
		value[i] = byte(i)
	}
	return value
}

func xadesSignABuilderVerifier() validation.CertificateVerifier {
	return validation.NewCommonCertificateVerifier()
}

// xadesSignABuilderReadAll drains a DSSDocument.
func xadesSignABuilderReadAll(t *testing.T, document model.DSSDocument) []byte {
	t.Helper()
	reader, err := document.OpenStream()
	if err != nil {
		t.Fatalf("opening the document: %v", err)
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	return content
}

// assertBytes compares one produced blob against the oracle, reporting the two as text when they
// differ (every blob in this KAT is XML or canonicalized XML).
func xadesSignABuilderAssertBytes(t *testing.T, caseName, key string, want, got []byte) {
	t.Helper()
	if !bytes.Equal(want, got) {
		t.Errorf("case %s / %s: bytes differ from the Java oracle\nwant: %s\ngot:  %s",
			caseName, key, want, got)
	}
}

// runBuildCase is the body every signing case shares: build, compare the three intermediate
// blobs, sign, compare the document.
func xadesSignABuilderRunBuildCase(t *testing.T, oracle xadesSignABuilderOracle, caseName string,
	params *SignatureParameters, documents []model.DSSDocument) {
	t.Helper()

	builderRef, err := SignatureBuilderGetSignatureBuilderForDocuments(params, documents,
		xadesSignABuilderVerifier())
	if err != nil {
		t.Fatalf("case %s: building the signature builder: %v", caseName, err)
	}
	canonicalizedSignedInfo, err := builderRef.Build()
	if err != nil {
		t.Fatalf("case %s: Build: %v", caseName, err)
	}
	xadesSignABuilderAssertBytes(t, caseName, "signedinfo-c14n",
		oracle.expect(t, caseName, "signedinfo-c14n"), canonicalizedSignedInfo)

	base := xadesSignABuilderBaseOf(t, builderRef)

	signedInfoRaw, err := xmlutils.DomUtilsSerializeNode(base.SignedInfoDom)
	if err != nil {
		t.Fatalf("case %s: serializing ds:SignedInfo: %v", caseName, err)
	}
	xadesSignABuilderAssertBytes(t, caseName, "signedinfo-raw",
		oracle.expect(t, caseName, "signedinfo-raw"), signedInfoRaw)

	if _, hasSignedProps := oracle[caseName]["signedprops"]; hasSignedProps {
		signedProps, err := xmlutils.DomUtilsSerializeNode(base.SignedPropertiesDom)
		if err != nil {
			t.Fatalf("case %s: serializing xades:SignedProperties: %v", caseName, err)
		}
		xadesSignABuilderAssertBytes(t, caseName, "signedprops",
			oracle.expect(t, caseName, "signedprops"), signedProps)
	}

	signed, err := builderRef.SignDocument(xadesSignABuilderSignatureValue())
	if err != nil {
		t.Fatalf("case %s: SignDocument: %v", caseName, err)
	}
	xadesSignABuilderAssertBytes(t, caseName, "document", oracle.expect(t, caseName, "document"),
		xadesSignABuilderReadAll(t, signed))
}

// xadesSignABuilderBaseOf recovers the embedded AbstractSignatureBuilder of whichever packaging
// builder the factory returned, so the test can read the cached DOM elements the oracle dumps.
func xadesSignABuilderBaseOf(t *testing.T, ref SignatureBuilderRef) *AbstractSignatureBuilder {
	t.Helper()
	switch builder := ref.(type) {
	case *EnvelopingSignatureBuilder:
		return &builder.AbstractSignatureBuilder
	case *EnvelopedSignatureBuilder:
		return &builder.AbstractSignatureBuilder
	case *DetachedSignatureBuilder:
		return &builder.AbstractSignatureBuilder
	case *InternallyDetachedSignatureBuilder:
		return &builder.AbstractSignatureBuilder
	}
	t.Fatalf("unexpected builder type %T", ref)
	return nil
}

func TestXAdESSignatureBuilderAgainstJavaOraclePackaging(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-b", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-b", xadesSignABuilderEnvelopingParams(t),
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
	t.Run("enveloped-b", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "enveloped-b", xadesSignABuilderEnvelopedParams(t),
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})
	t.Run("detached-b", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "detached-b", xadesSignABuilderDetachedParams(t),
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
	t.Run("internally-detached-b", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "internally-detached-b",
			xadesSignABuilderInternallyDetachedParams(t),
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})

	second := model.NewInMemoryDocumentWithMimeType([]byte("second"), "second.txt",
		enumerations.MimeTypeEnumText)
	t.Run("enveloping-two-documents", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-two-documents",
			xadesSignABuilderEnvelopingParams(t),
			[]model.DSSDocument{xadesSignABuilderTextDocument(), second})
	})
	t.Run("detached-two-documents", func(t *testing.T) {
		xadesSignABuilderRunBuildCase(t, oracle, "detached-two-documents",
			xadesSignABuilderDetachedParams(t),
			[]model.DSSDocument{xadesSignABuilderTextDocument(), second})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOracleEnvelopingObjects(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-embed-xml", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetEmbedXML(true)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-embed-xml", params,
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})

	t.Run("enveloping-manifest", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetManifestSignature(true)
		manifestBuilder, err := NewManifestBuilder(enumerations.DigestAlgorithmSHA256,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
		if err != nil {
			t.Fatalf("NewManifestBuilder: %v", err)
		}
		manifest, err := manifestBuilder.Build()
		if err != nil {
			t.Fatalf("building the manifest: %v", err)
		}
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-manifest", params,
			[]model.DSSDocument{manifest})
	})

	t.Run("enveloping-custom-objects", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)

		xmlObject := NewDSSObject()
		xmlObject.SetContent(model.NewInMemoryDocument(
			[]byte(`<data xmlns="http://sample.com">object-content</data>`)))
		xmlObject.SetId("custom-xml-object")
		xmlObject.SetMimeType(enumerations.MimeTypeEnumXML.MimeTypeString())

		textObject := NewDSSObject()
		textObject.SetContent(model.NewInMemoryDocument([]byte("plain object")))
		textObject.SetId("custom-text-object")
		textObject.SetMimeType(enumerations.MimeTypeEnumText.MimeTypeString())
		textObject.SetEncodingAlgorithm("http://www.w3.org/2000/09/xmldsig#base64")

		params.SetObjects([]*DSSObject{xmlObject, textObject})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-custom-objects", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOracleXPathPlacement(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloped-xpath-after", func(t *testing.T) {
		params := xadesSignABuilderEnvelopedParams(t)
		params.SetXPathLocationString("//*[local-name()='child']")
		params.SetXPathElementPlacement(XPathElementPlacementXPathAfter)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloped-xpath-after", params,
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})

	t.Run("enveloped-xpath-first-child", func(t *testing.T) {
		params := xadesSignABuilderEnvelopedParams(t)
		params.SetXPathLocationString("//*[local-name()='root']")
		params.SetXPathElementPlacement(XPathElementPlacementXPathFirstChildOf)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloped-xpath-first-child", params,
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})

	t.Run("internally-detached-xpath", func(t *testing.T) {
		params := xadesSignABuilderInternallyDetachedParams(t)
		params.SetXPathLocationString("//*[local-name()='root']")
		xadesSignABuilderRunBuildCase(t, oracle, "internally-detached-xpath", params,
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOracleNamespaces(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-default-xmldsig-prefix", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetXmldsigNamespace(common.NewDSSNamespace(common.XMLDSigNS.Uri(), ""))
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-default-xmldsig-prefix", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-same-prefix", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetXmldsigNamespace(common.NewDSSNamespace(common.XMLDSigNS.Uri(), "ns"))
		params.SetXadesNamespace(common.NewDSSNamespace(definition.XAdESNamespaceXAdES132.Uri(), "ns"))
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-same-prefix", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-xades111", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetXadesNamespace(common.NewDSSNamespace(definition.XAdESNamespaceXAdES111.Uri(), "xades111"))
		params.SetEn319132(false)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-xades111", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-xades122", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetXadesNamespace(common.NewDSSNamespace(definition.XAdESNamespaceXAdES122.Uri(), "xades122"))
		params.SetEn319132(false)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-xades122", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOracleKeyInfo(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-sign-key-info", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetSignKeyInfo(true)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-sign-key-info", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-x509-subject-name", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetAddX509SubjectName(true)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-x509-subject-name", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-no-signing-certificate", func(t *testing.T) {
		params := NewSignatureParameters()
		signingDate := xadesSignABuilderSigningDate
		params.BLevel().SetSigningDate(&signingDate)
		params.SetSignatureLevel(enumerations.SignatureLevelXAdESBaselineB)
		params.SetSignaturePackaging(enumerations.SignaturePackagingEnveloping)
		params.SetGenerateTBSWithoutCertificate(true)
		params.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
		params.SetEncryptionAlgorithm(enumerations.EncryptionAlgorithmRSA)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-no-signing-certificate", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-signing-cert-sha1", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetSigningCertificateDigestMethod(enumerations.DigestAlgorithmSHA1)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-signing-cert-sha1", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-signing-certificate-v1", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetEn319132(false)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-signing-certificate-v1", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOracleSignedProperties(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-policy-implied", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.BLevel().SetSignaturePolicy(model.NewPolicy())
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-policy-implied", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-policy-explicit", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		policy := model.NewPolicy()
		policy.SetId("http://spuri.test/policy")
		policy.SetDescription("Test policy description")
		policy.SetDocumentationReferences("http://doc.ref/1", "http://doc.ref/2")
		policy.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
		policy.SetDigestValue([]byte{1, 2, 3, 4})
		policy.SetSpuri("http://spuri.test")
		userNotice := model.NewUserNotice()
		userNotice.SetOrganization("ACME Ltd")
		userNotice.SetNoticeNumbers(1, 2)
		userNotice.SetExplicitText("This is a test policy")
		policy.SetUserNotice(userNotice)
		spDocSpecification := model.NewSpDocSpecification()
		spDocSpecification.SetId("1.2.3.4.6")
		spDocSpecification.SetDescription("Spec description")
		spDocSpecification.SetDocumentationReferences("http://spec.ref/1")
		policy.SetSpDocSpecification(spDocSpecification)
		params.BLevel().SetSignaturePolicy(policy)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-policy-explicit", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-policy-urn-oid", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		policy := model.NewPolicy()
		policy.SetId("1.2.3.4.5")
		policy.SetQualifier(enumerations.ObjectIdentifierQualifierOIDAsURN)
		policy.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
		policy.SetDigestValue(make([]byte, 32))
		params.BLevel().SetSignaturePolicy(policy)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-policy-urn-oid", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-signer-role-v2", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.BLevel().SetClaimedSignerRoles([]string{"Manager", "Director"})
		params.BLevel().SetSignedAssertions([]string{
			`<Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion">assertion-body</Assertion>`})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-signer-role-v2", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-signer-role-v1", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetEn319132(false)
		params.BLevel().SetClaimedSignerRoles([]string{"Manager", "Director"})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-signer-role-v1", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-production-place-v2", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.BLevel().SetSignerLocation(xadesSignABuilderSignerLocation())
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-production-place-v2", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-production-place-v1", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetEn319132(false)
		params.BLevel().SetSignerLocation(xadesSignABuilderSignerLocation())
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-production-place-v1", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-commitment-enum", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.BLevel().SetCommitmentTypeIndications(
			[]enumerations.CommitmentType{enumerations.CommitmentTypeEnumProofOfOrigin})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-commitment-enum", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-commitment-qualified", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		commitmentType := model.NewCommonCommitmentType()
		commitmentType.SetOid("1.2.840.113549.1.9.16.6.1")
		commitmentType.SetQualifier(enumerations.ObjectIdentifierQualifierOIDAsURN)
		commitmentType.SetDescription("Commitment description")
		commitmentType.SetDocumentationReferences("http://commitment.ref/1")
		commitmentType.SetSignedDataObjects("r-id-1", "#r-id-2")
		xmlQualifier := model.NewCommitmentQualifier()
		xmlQualifier.SetContent(model.NewInMemoryDocument(
			[]byte(`<qualifier xmlns="http://sample.com">value</qualifier>`)))
		textQualifier := model.NewCommitmentQualifier()
		textQualifier.SetContent(model.NewInMemoryDocument([]byte("plain qualifier")))
		commitmentType.SetCommitmentTypeQualifiers(xmlQualifier, textQualifier)
		params.BLevel().SetCommitmentTypeIndications(
			[]enumerations.CommitmentType{commitmentType})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-commitment-qualified", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-data-object-format", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		dataObjectFormat := NewDSSDataObjectFormat()
		dataObjectFormat.SetDescription("Custom description")
		dataObjectFormat.SetMimeType(enumerations.MimeTypeEnumText.MimeTypeString())
		dataObjectFormat.SetEncoding("http://www.w3.org/2000/09/xmldsig#base64")
		dataObjectFormat.SetObjectReference("#r-id-1")
		params.SetDataObjectFormatList([]*DSSDataObjectFormat{dataObjectFormat})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-data-object-format", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
}

func xadesSignABuilderSignerLocation() *model.SignerLocation {
	signerLocation := model.NewSignerLocation()
	signerLocation.SetLocality("Luxembourg")
	signerLocation.SetStreetAddress("1 Main Street")
	signerLocation.SetStateOrProvince("Luxembourg")
	signerLocation.SetPostalCode("L-1234")
	signerLocation.SetCountry("LU")
	return signerLocation
}

func TestXAdESSignatureBuilderAgainstJavaOracleContentTimestamps(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	der, err := os.ReadFile("testdata/content-timestamp.tst")
	if err != nil {
		t.Fatalf("reading the content timestamp: %v", err)
	}

	t.Run("enveloping-content-timestamp", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		token, err := validation.NewTimestampToken(der,
			enumerations.TimestampTypeAllDataObjectsTimestamp)
		if err != nil {
			t.Fatalf("parsing the content timestamp: %v", err)
		}
		token.SetCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#")
		params.SetContentTimestamps([]*validation.TimestampToken{token})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-content-timestamp", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloping-individual-content-timestamp", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		token, err := validation.NewTimestampToken(der,
			enumerations.TimestampTypeIndividualDataObjectsTimestamp)
		if err != nil {
			t.Fatalf("parsing the content timestamp: %v", err)
		}
		token.SetCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#")
		token.SetTimestampIncludes([]*validation.TimestampInclude{
			validation.NewTimestampIncludeWithURI("r-id-1", true)})
		params.SetContentTimestamps([]*validation.TimestampToken{token})
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-individual-content-timestamp", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})
}

func TestXAdESSignatureBuilderAgainstJavaOraclePrettyPrint(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	t.Run("enveloping-pretty-print", func(t *testing.T) {
		params := xadesSignABuilderEnvelopingParams(t)
		params.SetPrettyPrint(true)
		params.GetContext().SetOperationKind(enumerations.SigningOperationSign)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloping-pretty-print", params,
			[]model.DSSDocument{xadesSignABuilderTextDocument()})
	})

	t.Run("enveloped-pretty-print", func(t *testing.T) {
		params := xadesSignABuilderEnvelopedParams(t)
		params.SetPrettyPrint(true)
		params.GetContext().SetOperationKind(enumerations.SigningOperationSign)
		xadesSignABuilderRunBuildCase(t, oracle, "enveloped-pretty-print", params,
			[]model.DSSDocument{xadesSignABuilderXMLDocument()})
	})
}

// TestXAdESSignatureBuilderExplicitSignedData covers the two "bring your own bytes" parameters -
// setSignedData (an externally produced ds:SignedInfo) and setSignedAdESObject (an externally
// produced ds:Object) - which short-circuit reference building entirely.
func TestXAdESSignatureBuilderExplicitSignedData(t *testing.T) {
	oracle := loadXAdESSignABuilderOracle(t)

	sourceParams := xadesSignABuilderEnvelopingParams(t)
	sourceRef, err := SignatureBuilderGetSignatureBuilder(sourceParams,
		xadesSignABuilderTextDocument(), xadesSignABuilderVerifier())
	if err != nil {
		t.Fatalf("building the source builder: %v", err)
	}
	if _, err := sourceRef.Build(); err != nil {
		t.Fatalf("building the source signature: %v", err)
	}
	sourceBase := xadesSignABuilderBaseOf(t, sourceRef)
	signedInfoBytes, err := xmlutils.DomUtilsSerializeNode(sourceBase.SignedInfoDom)
	if err != nil {
		t.Fatalf("serializing the source ds:SignedInfo: %v", err)
	}
	signedAdESObject, err := xmlutils.DomUtilsSerializeNode(sourceBase.QualifyingPropertiesDom.Parent)
	if err != nil {
		t.Fatalf("serializing the source ds:Object: %v", err)
	}

	params := xadesSignABuilderEnvelopingParams(t)
	params.SetSignedData(signedInfoBytes)
	params.SetSignedAdESObject(signedAdESObject)
	xadesSignABuilderRunBuildCase(t, oracle, "enveloping-explicit-signed-data", params,
		[]model.DSSDocument{xadesSignABuilderTextDocument()})
}
