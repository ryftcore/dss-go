package cades

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/token"
)

// TestCounterSignatureCarriesNoMimeType drives a whole counter-signature through
// Service.GetDataToBeCounterSigned / CounterSignSignature. CAdESLevelBaselineB#addMimeType returns
// early for CAdESCounterSignatureParameters, so upstream's counter-signature has no mimeType
// signed attribute (a counter-signature signs a SignerInfo, not a document), whereas the
// signature it counter-signs has one. CounterSignatureBuilder used to leave the
// counter-signature flag of CMSForCAdESBuilderHelper unset, so the port added the attribute to
// every counter-signature; TestCAdESLevelBaselineBCounterSignatureSuppressesMimeType only covered
// the profile in isolation.
func TestCounterSignatureCarriesNoMimeType(t *testing.T) {
	keyStoreBytes, err := os.ReadFile(filepath.Join("testdata", "bytecmp", "signer_rsa.p12"))
	if err != nil {
		t.Fatal(err)
	}
	signatureToken, err := token.NewPkcs12SignatureTokenFromBytes(keyStoreBytes,
		token.NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		t.Fatal(err)
	}
	keys, err := signatureToken.Keys()
	if err != nil || len(keys) == 0 {
		t.Fatalf("keys: %v, %d entries", err, len(keys))
	}
	entry, ok := keys[0].(*token.KSPrivateKeyEntry)
	if !ok {
		t.Fatalf("unexpected key entry type %T", keys[0])
	}
	signerEntry, err := signatureToken.KeyWithPassword(entry.Alias(), token.NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		t.Fatal(err)
	}
	signingDate := time.Date(2027, 1, 15, 10, 30, 45, 0, time.UTC)
	service := NewService(validation.NewCommonCertificateVerifier())

	parameters := NewSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevelCAdESBaselineB)
	parameters.SetSignaturePackaging(enumerations.SignaturePackagingEnveloping)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	parameters.BLevel().SetSigningDate(&signingDate)

	document := model.NewInMemoryDocumentWithName([]byte("counter-signature probe content"), "probe.bin")
	dataToSign := service.GetDataToSign(document, parameters)
	signatureValue, err := signatureToken.Sign(dataToSign, parameters.DigestAlgorithm(), signerEntry)
	if err != nil {
		t.Fatal(err)
	}
	signed := service.SignDocument(document, parameters, signatureValue)

	analyzer, err := NewCMSDocumentAnalyzerFromDocument(signed)
	if err != nil {
		t.Fatal(err)
	}
	signatures := analyzer.Signatures()
	if len(signatures) != 1 {
		t.Fatalf("%d signatures, want 1", len(signatures))
	}
	master, ok := signatures[0].(*Signature)
	if !ok {
		t.Fatalf("unexpected signature type %T", signatures[0])
	}
	if master.SignerInformation().SignedAttributes.Get(spi.OIDIdAaEtsMimeType) == nil {
		t.Fatal("the master signature is expected to carry a mimeType attribute")
	}

	counterParameters := NewCounterSignatureParameters()
	counterParameters.SetSignatureLevel(enumerations.SignatureLevelCAdESBaselineB)
	counterParameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	counterParameters.SetSigningCertificate(signerEntry.Certificate())
	counterParameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	counterParameters.BLevel().SetSigningDate(&signingDate)
	counterParameters.SetSignatureIdToCounterSign(master.ID())

	counterDataToSign := service.GetDataToBeCounterSigned(signed, counterParameters)
	counterSignatureValue, err := signatureToken.Sign(counterDataToSign, counterParameters.DigestAlgorithm(), signerEntry)
	if err != nil {
		t.Fatal(err)
	}
	counterSigned := service.CounterSignSignature(signed, counterParameters, counterSignatureValue)

	analyzer, err = NewCMSDocumentAnalyzerFromDocument(counterSigned)
	if err != nil {
		t.Fatal(err)
	}
	signatures = analyzer.Signatures()
	if len(signatures) != 1 {
		t.Fatalf("%d signatures after counter-signing, want 1", len(signatures))
	}
	counterSignatures := signatures[0].CounterSignatures()
	if len(counterSignatures) != 1 {
		t.Fatalf("%d counter-signatures, want 1", len(counterSignatures))
	}
	counter, ok := counterSignatures[0].(*Signature)
	if !ok {
		t.Fatalf("unexpected counter-signature type %T", counterSignatures[0])
	}
	if counter.SignerInformation().SignedAttributes.Get(spi.OIDIdAaEtsMimeType) != nil {
		t.Error("a counter-signature shall carry no mimeType signed attribute")
	}
	if counter.SignerInformation().SignedAttributes.Get(spi.OIDIdAaSigningCertificateV2) == nil &&
		counter.SignerInformation().SignedAttributes.Get(spi.OIDIdAaSigningCertificate) == nil {
		t.Error("the counter-signature lost its signing-certificate attribute")
	}
}
