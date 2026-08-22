# Migrating from Java DSS

If you already know Java DSS, you know this library. It is a port, file by file,
class by class: the same services, the same parameter objects, the same
validation process, the same reports. This page is the translation layer.

## The five rules that cover most of it

**1. Class names survive.** `CAdESService` is `cades.CAdESService`.
`XAdESSignatureParameters` is `xades.XAdESSignatureParameters`. Java packages
become Go packages by the mapping in the table below.

**2. Getters and setters lose their prefix — sometimes.** Java's
`getSignatureAlgorithm()` becomes `SignatureAlgorithm()` where the `get` reads
as noise. Where upstream's own naming is load-bearing, it is kept. When in
doubt, `go doc` the package.

**3. Enum constants are `TypeJavaName` in MixedCaps, and their value is still
Java's `name()`.** Java's class-scoped constants share one namespace in Go, so
the type name is a prefix; the underscores of the Java name are dropped and each
part is MixedCased, with acronyms kept upper-case:

```java
SignatureLevel.XAdES_BASELINE_LTA   // Java
```
```go
enumerations.SignatureLevelXAdESBaselineLTA  // Go
```

The *value* is unchanged — the string `"XAdES_BASELINE_LTA"`, exactly Java's
`name()` — which is what makes serialized output interchangeable. Only the Go
spelling of the identifier differs from Java's. Some more mappings:

| Java | Go |
| --- | --- |
| `DigestAlgorithm.SHA256` | `enumerations.DigestAlgorithmSHA256` |
| `SignaturePackaging.ENVELOPED` | `enumerations.SignaturePackagingEnveloped` |
| `AdditionalServiceInformation.FOR_ESIGNATURES` | `enumerations.AdditionalServiceInformationForESignatures` |
| `MessageTag.BBB_XCV_CCCBB` | `i18n.MessageTagBBBXCVCCCBB` |

`valueOf` becomes `SignatureLevelValueOf(string) (SignatureLevel, error)`;
`forOid` / `forUri` become `…ForOID` / `…ForURI`.

**4. Exceptions become errors — except where they become panics.** Checked and
meaningful exceptions are `error` returns, matchable with `errors.As` against
types like `model.DSSError`. Constructors that throw become
`New…() (T, error)`. But the ported services follow Java in raising unchecked
exceptions from deep inside, which the port turns into panics; **the `dss`
facade recovers those and returns them as errors**. Below the facade, recover
them yourself.

**5. Static utility classes become prefixed package functions.**
`DSSUtils.loadCertificate(…)` is `spi.DSSUtilsLoadCertificate(…)`. The class
name becomes part of the function name, because Go has no class scope.

## Two things that genuinely differ

**No `ServiceLoader`.** Java discovers format validators, validation policies
and cryptographic suites at runtime. Go has no equivalent, so **importing the
root `dss` package registers all six format families, the ETSI policy and the
XML cryptographic suite** — which is what makes `dss.Validate`'s format
auto-detection work. If you use the underlying packages directly, without the
facade, you must register them yourself:

```go
validation.RegisterDocumentValidatorFactory(…)
policy.RegisterValidationPolicyFactory(…)
```

Blank-importing the format packages you need is usually the tidiest way.

**BouncyCastle is gone.** Its role is filled by `crypto/*`,
`encoding/asn1`, `golang.org/x/crypto/cryptobyte` and — for what the standard
library cannot do byte-exactly — internal packages written for this port: a
BER/DER engine, a CMS and RFC 3161 core, seven canonicalization algorithms, a
PDF engine, a JOSE serializer. If your Java code touched BouncyCastle types
directly, that part does not translate; the port's own types replace them.

## The facade shortcut

Before translating class by class, check whether you need to. The root `dss`
package collapses the common signing and validation flows into a handful of
calls:

```java
// Java
CAdESService service = new CAdESService(commonCertificateVerifier);
service.setTspSource(tspSource);
CAdESSignatureParameters parameters = new CAdESSignatureParameters();
parameters.setSignatureLevel(SignatureLevel.CAdES_BASELINE_T);
parameters.setSigningCertificate(privateKey.getCertificate());
parameters.setCertificateChain(privateKey.getCertificateChain());
ToBeSigned dataToSign = service.getDataToSign(document, parameters);
SignatureValue signatureValue = token.sign(dataToSign, digestAlgorithm, privateKey);
DSSDocument signed = service.signDocument(document, parameters, signatureValue);
```

```go
// Go, through the facade
signed, err := dss.Sign(doc, signer, dss.SignOptions{
    Format:    dss.FormatCAdES,
    Level:     dss.LevelT,
    TSPSource: tsa,
})
```

The facade delegates to exactly the services above and holds no logic of its
own — so anything it does not cover (counter-signatures, visible appearances,
XAdES references and transforms, ASiC filename factories, policy stores) you
reach by dropping to the same packages your Java code already names. Mixing the
two levels needs no conversion: the facade's types are Go type *aliases* of the
underlying ones, not wrappers.

## Java class → Go symbol

Harvested from the `// Ported from` header on each Go file. Import paths are
relative to `github.com/ryftcore/dss-go/dss`.

### Signing services and parameters

| Java class | Go package | Go symbol |
|---|---|---|
| `CAdESService` | `cades` | `CAdESService`, `NewCAdESService` |
| `XAdESService` | `xades` | `XAdESService`, `NewXAdESService` |
| `PAdESService` | `pades` | `PAdESService`, `NewPAdESService` |
| `JAdESService` | `jades` | `JAdESService`, `NewJAdESService` |
| `ASiCWithCAdESService` | `asic/cades` | `ASiCWithCAdESService`, `NewASiCWithCAdESService` |
| `ASiCWithXAdESService` | `asic/xades` | `ASiCWithXAdESService`, `NewASiCWithXAdESService` |
| `CAdESSignatureParameters` | `cades` | `CAdESSignatureParameters` |
| `XAdESSignatureParameters` | `xades` | `XAdESSignatureParameters` |
| `PAdESSignatureParameters` | `pades` | `PAdESSignatureParameters` |
| `JAdESSignatureParameters` | `jades` | `JAdESSignatureParameters` |
| `ASiCWithCAdESSignatureParameters` | `asic/cades` | `ASiCWithCAdESSignatureParameters` |
| `ASiCWithXAdESSignatureParameters` | `asic/xades` | `ASiCWithXAdESSignatureParameters` |

### Keys and tokens

| Java class | Go package | Go symbol |
|---|---|---|
| `SignatureTokenConnection` | `token` | `SignatureTokenConnection` |
| `Pkcs12SignatureToken` | `token` | `Pkcs12SignatureToken`, `NewPkcs12SignatureTokenFromFilepath` |
| `JKSSignatureToken` | `token` | `JKSSignatureToken` — **constructors return "not supported"** |
| `Pkcs11SignatureToken` | `token` | `Pkcs11SignatureToken` — **constructors return "not supported"** |
| `DSSPrivateKeyEntry` | `token` | `DSSPrivateKeyEntry` |
| `TSPSource` | `spi/validation` | `TSPSource` |
| `KeyEntityTSPSource` | `spi/validation` | `KeyEntityTSPSource`, `NewKeyEntityTSPSource` |
| `OnlineTSPSource` | — | **not ported**; see [Known gaps](../compatibility/known-gaps.md) |

### Documents and model types

| Java class | Go package | Go symbol |
|---|---|---|
| `DSSDocument` | `model` | `DSSDocument` (interface) |
| `InMemoryDocument` | `model` | `InMemoryDocument`, `NewInMemoryDocumentWithName` |
| `FileDocument` | `model` | `FileDocument`, `NewFileDocument` |
| `DigestDocument` | `model` | `DigestDocument`, `NewDigestDocument` |
| `CertificateToken` | `model` | `CertificateToken`, `NewCertificateToken` |
| `SignatureValue` | `model` | `SignatureValue`, `NewSignatureValue` |
| `ToBeSigned` | `model` | `ToBeSigned`, `NewToBeSigned` |
| `DSSException` | `model` | `DSSError`, `NewDSSError` — match with `errors.As` |
| `DSSUtils` | `spi` | package functions prefixed `DSSUtils…` |

### Validation

| Java class | Go package | Go symbol |
|---|---|---|
| `SignedDocumentValidator` | `validation` | `SignedDocumentValidator`, `SignedDocumentValidatorFromDocument` |
| `DocumentValidator` | `validation` | `DocumentValidator` (interface) |
| `CertificateVerifier` | `spi/validation` | `CertificateVerifier` (interface) |
| `CommonCertificateVerifier` | `spi/validation` | `CommonCertificateVerifier`, `NewCommonCertificateVerifier` |
| `CertificateSource` | `spi` | `CertificateSource` (interface) |
| `CommonTrustedCertificateSource` | `spi` | `CommonTrustedCertificateSource` |
| `AdvancedSignature` | `spi/validation` | `AdvancedSignature` (interface) |
| `DefaultAdvancedSignature` | `spi/validation` | `DefaultAdvancedSignature` |
| `TimestampToken` | `spi/validation` | `TimestampToken` |
| `RevocationToken` | `spi` | `RevocationToken` |
| `ValidationPolicy` | `model/policy` | `ValidationPolicy` (interface) |
| `EtsiValidationPolicy` | `policy` | `EtsiValidationPolicy` |

### Reports

| Java class | Go package | Go symbol |
|---|---|---|
| `Reports` | `validation/reports` | `Reports` — also embedded in `dss.Reports` |
| `CertificateReports` | `validation/reports` | `CertificateReports` |
| `SimpleReport` | `simplereport` | `SimpleReport` |
| `DetailedReport` | `detailedreport` | `DetailedReport` |
| `DiagnosticData` | `diagnostic` | `DiagnosticData` |

### Trusted lists

| Java class | Go package | Go symbol |
|---|---|---|
| `TLValidationJob` | `tsl` | `TLValidationJob`, `NewTLValidationJob` |
| `LOTLSource` | `tsl` | `LOTLSource`, `NewLOTLSource` |
| `TLSource` | `tsl` | `TLSource`, `NewTLSource` |
| `TrustedListsCertificateSource` | `spi/tsl` | `TrustedListsCertificateSource` |

### Containers and enumerations

| Java class | Go package | Go symbol |
|---|---|---|
| `ASiCContainerExtractor` | `asic` | `ASiCContainerExtractor` |
| `SignatureLevel` | `enumerations` | `SignatureLevel` |
| `SignaturePackaging` | `enumerations` | `SignaturePackaging` |
| `DigestAlgorithm` | `enumerations` | `DigestAlgorithm` |
| `SignatureAlgorithm` | `enumerations` | `SignatureAlgorithm` |
| `EncryptionAlgorithm` | `enumerations` | `EncryptionAlgorithm` |
| `ASiCContainerType` | `enumerations` | `ASiCContainerType` |
| `ValidationLevel` | `enumerations` | `ValidationLevel` |
| `Indication` | `enumerations` | `Indication` |
| `SubIndication` | `enumerations` | `SubIndication` |
| `SignatureQualification` | `enumerations` | `SignatureQualification` |
| `MimeType` | `enumerations` | `MimeType` |

Anything not listed: the Go file's header names its Java source, so
`grep -r "YourClass.java" dss/` finds it. Full API reference on
[pkg.go.dev](https://pkg.go.dev/github.com/ryftcore/dss-go/dss).

## Maven module → Go package

| Java module(s) | Go package |
|---|---|
| `dss-enumerations` | `enumerations` |
| `dss-alert` | `alert` |
| `dss-utils` (+ commons/guava implementations) | `utils` |
| `dss-model` | `model` |
| `dss-spi`, `dss-crl-parser*`, `dss-token`, `dss-document` | `spi`, `crlparser`, `token`, `document` |
| `dss-cms`, `dss-cms-object`, `dss-cms-stream` | `cms` |
| `dss-cades` | `cades` |
| `dss-xml-common`, `dss-xml-utils` | `xml`, plus internal canonicalization packages |
| `dss-xades` | `xades` |
| `dss-pades` (+ the PDF engine) | `pades` |
| `dss-jades` | `jades` |
| `dss-asic-*` | `asic`, `asic/cades`, `asic/xades` |
| `dss-validation` (+ `*-report-jaxb`, `dss-policy-*`) | `validation`, `validation/reports`, `validation/policy`, `policy`, `diagnostic` |
| `dss-tsl-validation`, `dss-validation-job` | `tsl` |
| `dss-i18n` | `i18n` |

## A migration checklist

- [ ] **Check whether the facade covers your flow.** Most applications do only
      "sign this" and "validate that", and the facade is a much smaller surface
      to maintain.
- [ ] **Replace `ServiceLoader` assumptions** with explicit registration — or
      import the root `dss` package and let it register for you.
- [ ] **Replace exception handling with error handling**, and remember the
      facade recovers panics from the layers below.
- [ ] **Find your `OnlineTSPSource` / `OnlineCRLSource` / `OnlineOCSPSource`
      usage** and decide what replaces it. This is the change most likely to
      need real work.
- [ ] **Check for JKS or PKCS#11 key stores.** Neither is supported; PKCS#12 is,
      fully.
- [ ] **Check for visible PDF signature appearances.** Not supported.
- [ ] **Diff the reports.** Run the same document through both implementations
      with pinned trust anchors, policy, validation time and locale, and diff
      the DetailedReports. See
      [Verifying Java DSS interop](../guides/verifying-java-dss-interop.md).

## Next

- [Getting started](../getting-started.md) — the facade in twenty lines.
- [Known gaps](../compatibility/known-gaps.md) — read before committing to the
  migration.
