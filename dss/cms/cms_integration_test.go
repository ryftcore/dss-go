// End-to-end round-trip test exercising the whole CMSAPI wrap-layer against internal/cmscore:
// CustomContentSignerBuilder -> SignerInfoGeneratorBuilder -> Builder -> Generator's
// native Generate, mirroring the two-step DSS signing flow (empty-signature "data to sign",
// then a real signature) CAdESService drives in the cades package. Not a KAT (no Java oracle
// output is compared byte for byte here - that is cades' BUILD chunk's job for the CAdES
// signed-attribute construction specifically); this proves the cms package's own plumbing
// produces a well-formed, self-consistent, cryptographically valid CMS SignedData.
package cms

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

func generateTestCertificate(t *testing.T) (*rsa.PrivateKey, *model.CertificateToken) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %s", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject:      pkix.Name{CommonName: "cms test signer"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %s", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %s", err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %s", err)
	}
	return key, token
}

func TestCMSBuilderRoundTrip(t *testing.T) {
	key, signingCertificate := generateTestCertificate(t)
	document := model.NewInMemoryDocument([]byte("hello CMS"))

	buildSignerInfoGenerator := func(contentSigner *CustomContentSigner) *SignerInfoGenerator {
		generator, err := NewCMSSignerInfoGeneratorBuilder().
			SetSigningCertificate(signingCertificate).
			SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256).
			Build(document, contentSigner)
		if err != nil {
			t.Fatalf("Build: %s", err)
		}
		return generator
	}

	// Step 1: data to sign, an empty-signature ContentSigner capturing the bytes-to-be-signed.
	dataToSignSigner, err := NewCustomContentSignerBuilder().Build(enumerations.SignatureAlgorithmRSASHA256)
	if err != nil {
		t.Fatalf("Build content signer: %s", err)
	}
	generator1 := buildSignerInfoGenerator(dataToSignSigner)

	builder := NewCMSBuilder().
		SetSigningCertificate(signingCertificate).
		SetTrustAnchorBPPolicy(false)
	if _, err := builder.CreateCMS(generator1, document); err != nil {
		t.Fatalf("CreateCMS (data to sign): %s", err)
	}
	dataToSign := dataToSignSigner.OutputStream().Bytes()
	if len(dataToSign) == 0 {
		t.Fatal("no data to sign captured")
	}

	// Step 2: sign the captured bytes for real and rebuild the CMS with the actual signature.
	digest := sha256.Sum256(dataToSign)
	signatureValue, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign: %s", err)
	}

	realSigner, err := NewCustomContentSignerBuilder().BuildWithSignatureValue(
		enumerations.SignatureAlgorithmRSASHA256, model.NewSignatureValueWithValue(enumerations.SignatureAlgorithmRSASHA256, signatureValue))
	if err != nil {
		t.Fatalf("BuildWithSignatureValue: %s", err)
	}
	generator2 := buildSignerInfoGenerator(realSigner)

	signedCMS, err := builder.CreateCMS(generator2, document)
	if err != nil {
		t.Fatalf("CreateCMS (final): %s", err)
	}

	// The final CMS parses back cleanly and holds one signer, one certificate, and an
	// encapsulated (non-detached) content, per Builder's SetEncapsulate default.
	if signedCMS.IsDetachedSignature() {
		t.Error("expected an encapsulated signature (SetEncapsulate defaults to true)")
	}
	if got, want := signedCMS.SignedContent().(*model.InMemoryDocument).Bytes(), []byte("hello CMS"); string(got) != string(want) {
		t.Errorf("signed content = %q, want %q", got, want)
	}
	if len(signedCMS.Certificates()) != 1 {
		t.Fatalf("got %d certificates, want 1", len(signedCMS.Certificates()))
	}
	if string(signedCMS.Certificates()[0]) != string(signingCertificate.Encoded()) {
		t.Error("the embedded certificate does not match the signing certificate")
	}

	reparsed, err := UtilsParseToCMSBinaries(signedCMS.DEREncoded())
	if err != nil {
		t.Fatalf("re-parse: %s", err)
	}
	signerInfos := reparsed.SignerInfos()
	if len(signerInfos) != 1 {
		t.Fatalf("got %d signer infos, want 1", len(signerInfos))
	}
	signerInfo := signerInfos[0]

	// The three injected attributes are present.
	if signerInfo.SignedAttributes.Get(cmscore.OIDContentType) == nil {
		t.Error("missing content-type attribute")
	}
	messageDigestAttribute := signerInfo.SignedAttributes.Get(cmscore.OIDMessageDigest)
	if messageDigestAttribute == nil {
		t.Fatal("missing message-digest attribute")
	}
	wantDigest := sha256.Sum256([]byte("hello CMS"))
	gotDigestElement := messageDigestAttribute.Values[0]
	if string(gotDigestElement.Octets()) != string(wantDigest[:]) {
		t.Errorf("message-digest = %x, want %x", gotDigestElement.Octets(), wantDigest)
	}
	if signerInfo.SignedAttributes.Get(OIDIdAaCmsAlgorithmProtect) == nil {
		t.Error("missing cms-algorithm-protection attribute")
	}

	// The signature verifies against the DER SET OF the signed attributes, exactly as a real
	// verifier would compute it (RFC 5652 clause 5.4).
	signedAttrsDigest := sha256.Sum256(signerInfo.SignedAttributesDER())
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, signedAttrsDigest[:], signerInfo.Signature); err != nil {
		t.Errorf("signature does not verify: %s", err)
	}
}

// TestCMSBuilderDetached checks SetEncapsulate(false): the signed content travels out of band,
// but the message-digest is still computed over it.
func TestCMSBuilderDetached(t *testing.T) {
	_, signingCertificate := generateTestCertificate(t)
	document := model.NewInMemoryDocument([]byte("detached content"))

	contentSigner, err := NewCustomContentSignerBuilder().BuildWithSignatureValue(enumerations.SignatureAlgorithmRSASHA256,
		model.NewSignatureValueWithValue(enumerations.SignatureAlgorithmRSASHA256, []byte{1, 2, 3, 4}))
	if err != nil {
		t.Fatalf("Build content signer: %s", err)
	}
	generator, err := NewCMSSignerInfoGeneratorBuilder().
		SetSigningCertificate(signingCertificate).
		SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256).
		Build(document, contentSigner)
	if err != nil {
		t.Fatalf("Build: %s", err)
	}

	cms, err := NewCMSBuilder().
		SetSigningCertificate(signingCertificate).
		SetTrustAnchorBPPolicy(false).
		SetEncapsulate(false).
		CreateCMS(generator, document)
	if err != nil {
		t.Fatalf("CreateCMS: %s", err)
	}
	if !cms.IsDetachedSignature() {
		t.Error("expected a detached signature")
	}
	if cms.SignedContent() != nil {
		t.Error("expected no signed content for a detached signature")
	}
}
