// Command bytecmp is the Go half of the CAdES-B byte-exactness harness
// (cades_byte_exactness_test.go): it signs a fixed payload with this package's own Service,
// with every input pinned - key store, signing certificate and chain, content, signing time,
// digest algorithm - so the resulting CMS can be compared byte for byte against the one upstream
// DSS's own CAdESService produces from the identical inputs (ByteExactnessFixtures.java).
//
// RSA PKCS#1 v1.5 is deterministic, so two implementations that agree on every byte of the signed
// attributes necessarily agree on the signature value too; any remaining difference is in how the
// SignedData around it is assembled.
//
// Usage: go run . <keystore.p12> <password> <signing year> <output directory> <prefix>
//
// Writes <outdir>/<prefix>-datatosign.bin (the DER SET OF signed attributes) and
// <outdir>/<prefix>.p7m (the complete CMS).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/token"
)

// content is the fixed payload, byte-identical to ByteExactnessFixtures.java's.
var content = []byte("byte-exactness probe content")

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintln(os.Stderr, "usage: bytecmp <keystore.p12> <password> <signing year> <outdir> <prefix>")
		os.Exit(2)
	}
	keyStorePath, password, yearText, outDir, prefix := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]

	year, err := strconv.Atoi(yearText)
	check(err)

	keyStoreBytes, err := os.ReadFile(keyStorePath)
	check(err)
	signatureToken, err := token.NewPkcs12SignatureTokenFromBytes(keyStoreBytes,
		token.NewPasswordProtection([]byte(password)))
	check(err)
	keys, err := signatureToken.Keys()
	check(err)
	if len(keys) == 0 {
		check(fmt.Errorf("%s carries no key entries", keyStorePath))
	}
	entry, ok := keys[0].(*token.KSPrivateKeyEntry)
	if !ok {
		check(fmt.Errorf("%s: unexpected key entry type %T", keyStorePath, keys[0]))
	}
	signerEntry, err := signatureToken.KeyWithPassword(entry.Alias(), token.NewPasswordProtection([]byte(password)))
	check(err)

	// A fixed instant inside the signing certificate's validity window; the signing-time
	// attribute is part of the signed attributes, so it has to match Java's exactly.
	signingDate := time.Date(year, 1, 15, 10, 30, 45, 0, time.UTC)

	parameters := cades.NewSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevelCAdESBaselineB)
	parameters.SetSignaturePackaging(enumerations.SignaturePackagingEnveloping)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	parameters.BLevel().SetSigningDate(&signingDate)

	service := cades.NewService(validation.NewCommonCertificateVerifier())
	document := model.NewInMemoryDocumentWithName(content, "probe.bin")

	dataToSign := service.GetDataToSign(document, parameters)
	check(os.WriteFile(filepath.Join(outDir, prefix+"-datatosign.bin"), dataToSign.Bytes(), 0o644))

	signatureValue, err := signatureToken.Sign(dataToSign, parameters.DigestAlgorithm(), signerEntry)
	check(err)
	signedDocument := service.SignDocument(document, parameters, signatureValue)
	check(signedDocument.Save(filepath.Join(outDir, prefix+".p7m")))

	fmt.Println("wrote", prefix+"-datatosign.bin", "and", prefix+".p7m")
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "bytecmp:", err)
		os.Exit(1)
	}
}
